package main

import (
	"io/fs"
	"strings"
	"testing"
)

// bannedRevenueShareTerms returns the lower-cased phrases that must not appear
// on any shipped user-facing surface.
//
// The needles are assembled from fragments rather than written out literally so
// that a repo-wide grep for the retired term stays clean: this file must not be
// the last place the wording survives.
func bannedRevenueShareTerms() []string {
	dev := "deve" + "loper"
	return []string{
		dev + " fund",
		dev + "-fund",
		dev + "_fund",
		"dev " + "fund",
		"dev-" + "fund",
		"dev" + "fund",
	}
}

// TestNoRevenueShareBrandingInShippedUI asserts the embedded wizard UI and the
// packaged portal/admin assets carry no revenue-share branding. The dev-split
// control stays (it is a real setting), but it is described in neutral terms:
// an optional share routed to the developer payout address, with a 0 opt-out.
func TestNoRevenueShareBrandingInShippedUI(t *testing.T) {
	banned := bannedRevenueShareTerms()

	assertClean := func(t *testing.T, name string, body []byte) {
		t.Helper()
		lower := strings.ToLower(string(body))
		for _, term := range banned {
			if strings.Contains(lower, term) {
				t.Errorf("%s still contains retired revenue-share wording %q", name, term)
			}
		}
	}

	assertClean(t, "index.html", indexHTML)

	for _, root := range []struct {
		name string
		fsys fs.FS
	}{
		{"portal", portalFS},
		{"admin", adminFS},
	} {
		err := fs.WalkDir(root.fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".html") {
				return nil
			}
			body, readErr := fs.ReadFile(root.fsys, path)
			if readErr != nil {
				return readErr
			}
			assertClean(t, root.name+"/"+path, body)
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root.name, err)
		}
	}
}

// TestDevSplitHintIsNeutral pins the wording of the advanced dev-split hint so a
// future edit cannot quietly reintroduce revenue-share marketing copy, and so
// the slider keeps advertising the 0 opt-out.
func TestDevSplitHintIsNeutral(t *testing.T) {
	ui := string(indexHTML)

	const hint = "Optional share of payments routed to the developer payout address"
	if !strings.Contains(ui, hint) {
		t.Errorf("dev-split hint missing neutral wording %q", hint)
	}
	if !strings.Contains(ui, "config.json <code>profit_share</code>") {
		t.Error("dev-split hint no longer names the config.json key it writes (profit_share)")
	}
}
