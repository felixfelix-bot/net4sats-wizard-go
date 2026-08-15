package main

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

// expectedPinnedURLConsts is the authoritative registry of pinned download
// URL constants declared in deploy.go. TestDeployGoPinRegistry enforces that
// deploy.go declares exactly this set — a new pin cannot be added without
// registering it here, and every registered pin is exercised live by
// TestPinnedURLsAreLive (add new pins to that test's map too).
var expectedPinnedURLConsts = []string{
	"tollgatePkgURL",
	"configwizURL",
}

// forbiddenDeployIdentifiers names pins/steps removed by SW4a (Aug 2026)
// because their URLs 404'd. They must not come back:
//
//   - tollgateNftEnforceURL / 20-nds-enforce.nft: never existed as a release
//     asset. The nft enforcement rules ship INSIDE the .ipk since
//     OpenTollGate/tollgate-module-basic-go PR #283 ("fix: NDS fw4/nftables
//     enforcement bridge", merged 2026-07-24) — verified by extracting
//     data.tar.gz from tollgate-wrt_main.56.b528e1d: it contains
//     ./etc/nftables.d/20-nds-enforce.nft and ./etc/nftables.d/30-backend-firewall.nft.
//     The separate overlay download in deploy step 4 was redundant and broken.
//
//   - tollgateOSURL: dead constant (declared but never referenced) whose
//     releases.tollgate.me URL returned 404.
var forbiddenDeployIdentifiers = []string{
	"tollgateNftEnforceURL",
	"tollgateOSURL",
	"20-nds-enforce.nft",
}

// wantTollgatePkgURL is the exact release asset pinned by the wizard.
// The v0.6.1-post-merge release of felixfelix-bot/tollgate-module-basic-go
// publishes exactly one asset: tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.ipk
// (verified via the GitHub releases API on 2026-08-16). A previous pin
// referenced a main.53 asset that does not exist on that release (HTTP 404),
// which broke every fresh wizard deploy.
const wantTollgatePkgURL = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.6.1-post-merge/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.ipk"

// TestTollgatePkgURLPinsExistingAsset pins the package download URL to the
// exact asset that exists on the v0.6.1-post-merge release. Any intentional
// repin (e.g. switching to an upstream OpenTollGate release) must update this
// test in the same commit — that is what makes the pin auditable.
func TestTollgatePkgURLPinsExistingAsset(t *testing.T) {
	if tollgatePkgURL != wantTollgatePkgURL {
		t.Errorf("tollgatePkgURL =\n  %q\nwant\n  %q", tollgatePkgURL, wantTollgatePkgURL)
	}
}

// TestDeployGoPinRegistry guards the set of pinned download URLs in deploy.go:
//
//  1. every "<name>URL = \"https://..." constant declared in deploy.go must be
//     listed in expectedPinnedURLConsts (so it gets a live HTTP 200 check), and
//  2. none of the removed 404 pins/steps (forbiddenDeployIdentifiers) may be
//     reintroduced.
func TestDeployGoPinRegistry(t *testing.T) {
	src, err := os.ReadFile("deploy.go")
	if err != nil {
		t.Fatalf("reading deploy.go: %v", err)
	}
	srcStr := string(src)

	// 1. Declared URL constants must exactly match the registry.
	declRe := regexp.MustCompile(`(?m)^\t(\w+URL)\s*=\s+"https://`)
	declared := map[string]bool{}
	for _, m := range declRe.FindAllStringSubmatch(srcStr, -1) {
		declared[m[1]] = true
	}
	want := map[string]bool{}
	for _, name := range expectedPinnedURLConsts {
		want[name] = true
	}
	if !reflect.DeepEqual(declared, want) {
		t.Errorf("deploy.go declares URL constants %v, but registry expects %v\n"+
			"→ add/remove pins in BOTH expectedPinnedURLConsts and TestPinnedURLsAreLive",
			setNames(declared), setNames(want))
	}

	// 2. Removed 404 pins/steps must stay removed.
	for _, f := range forbiddenDeployIdentifiers {
		if strings.Contains(srcStr, f) {
			t.Errorf("deploy.go still references %q — this pin/step was removed because it 404s (see SW4a)", f)
		}
	}
}

// TestPinnedURLsAreLive is the regression test for SW4a: every URL the wizard
// pins must return HTTP 200 today. A pin that 404s means every fresh deploy
// breaks at the download step, exactly the outage this test guards against.
//
// It performs real HTTP requests (Range: bytes=0-0 so asset bodies are not
// downloaded). Run with -short to skip network access for offline hacking;
// the CI-parity command is plain `go test ./...`, which runs this live.
func TestPinnedURLsAreLive(t *testing.T) {
	if testing.Short() {
		t.Skip("live URL check skipped: -short mode (CI-parity runs without -short)")
	}

	// Keep in sync with expectedPinnedURLConsts (enforced by
	// TestDeployGoPinRegistry).
	pins := map[string]string{
		"tollgatePkgURL": tollgatePkgURL,
		"configwizURL":   configwizURL,
	}

	client := &http.Client{Timeout: 30 * time.Second}
	for name, url := range pins {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatalf("building request for %s: %v", name, err)
			}
			req.Header.Set("Range", "bytes=0-0") // fetch 1 byte, not the asset
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("%s = %q: request failed: %v", name, url, err)
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1))

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
				t.Errorf("%s = %q: got HTTP %d, want 200 — broken pin, fresh deploys will fail (SW4a regression)",
					name, url, resp.StatusCode)
			}
		})
	}
}

// setNames returns a deterministic sorted name list for a set map (stable
// test failure output).
func setNames(set map[string]bool) []string {
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	// small n, simple sort
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && names[j] < names[j-1]; j-- {
			names[j], names[j-1] = names[j-1], names[j]
		}
	}
	return names
}
