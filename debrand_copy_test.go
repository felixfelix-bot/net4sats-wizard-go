package main

import (
	"io/fs"
	"os"
	"regexp"
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

// bannedSplitTerms returns the lower-cased display wording of the retired
// control name. The identifiers (the devSplit wire field, the devsplit DOM ids)
// are retained, so these needles only apply to operator-visible text: the
// guard strips markup before matching.
//
// Like bannedRevenueShareTerms, the needles are assembled from fragments so a
// repo-wide grep for the retired wording stays clean.
func bannedSplitTerms() []string {
	dev := "de" + "v"
	split := "spl" + "it"
	return []string{
		dev + " " + split,
		dev + "-" + split,
		dev + split,
	}
}

// visibleText strips <script>/<style> bodies and every tag, and with them every
// attribute value, so only the copy an operator can actually read is returned.
// This is what keeps the retained `devsplit` DOM ids — and the inline JS that
// reads them by id — from tripping a display-wording guard.
func visibleText(html string) string {
	return tagRE.ReplaceAllString(scriptRE.ReplaceAllString(html, " "), " ")
}

var (
	tagRE    = regexp.MustCompile(`(?s)<[^>]*>`)
	scriptRE = regexp.MustCompile(`(?is)<(script|style)\b[^>]*>.*?</(script|style)>`)
)

// retiredNames is every retired user-facing name: the fund phrasing and the
// control's old label. Identifiers (the devSplit wire field, the devsplit DOM
// ids) are deliberately absent — they are retained.
func retiredNames() []string {
	return append(bannedRevenueShareTerms(), bannedSplitTerms()...)
}

// assertNoRetiredName fails if surface still carries any retired name.
func assertNoRetiredName(t *testing.T, surface string, body []byte) {
	t.Helper()
	lower := strings.ToLower(string(body))
	for _, term := range retiredNames() {
		if strings.Contains(lower, term) {
			t.Errorf("%s still uses retired wording %q", surface, term)
		}
	}
}

// TestNoSplitWordingInVisibleCopy pins the user-visible control name. The
// label used to read "Dev split" while the hint beneath it had already been
// neutralised, so the retired name survived in the one place every operator
// looks. Markup is stripped first: the devsplit / devsplit-val ids are
// deliberate identifiers, not copy.
func TestNoSplitWordingInVisibleCopy(t *testing.T) {
	assertNoRetiredName(t, "operator-visible wizard copy", []byte(visibleText(string(indexHTML))))
}

// TestReadmeUsesNeutralPayoutWording keeps the operator-facing README in step
// with the UI. Docs ship in the same commit as the copy, so the doc must not
// name the control by its retired name either.
func TestReadmeUsesNeutralPayoutWording(t *testing.T) {
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	assertNoRetiredName(t, "README.md", readme)
	if !strings.Contains(strings.ToLower(string(readme)), "payout share") {
		t.Error(`README.md no longer names the setting in neutral wording (want "payout share")`)
	}
}

// TestDeployLogUsesNeutralPayoutWording covers the deploy log line, which is
// shipped copy: it is streamed to the operator's terminal during a run. It used
// to echo the raw devSplit field name.
func TestDeployLogUsesNeutralPayoutWording(t *testing.T) {
	src, err := os.ReadFile("deploy.go")
	if err != nil {
		t.Fatalf("read deploy.go: %v", err)
	}
	body := string(src)
	if !strings.Contains(body, `%, payout share=`) {
		t.Error("deploy log line no longer names the setting in neutral wording (want `%, payout share=`)")
	}
	if strings.Contains(body, `devSplit="`) {
		t.Error("deploy log line still echoes the devSplit field name instead of user-facing wording")
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
