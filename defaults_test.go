package main

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// The defaults pinned here are the ones that reach a real router without the
// operator doing anything: the embedded UI's preseeded form values and the mint
// set the deploy step writes into /etc/tollgate/config.json. Both are policy,
// so both are guarded rather than left to review.

var (
	// devsplitInputRE captures the advanced payout-share slider tag (the
	// retained `devsplit` DOM id) so its range can be asserted independently
	// of the surrounding markup.
	devsplitInputRE = regexp.MustCompile(`(?s)<input[^>]*id="devsplit"[^>]*>`)

	// advancedOpenRE bounds the collapsed Advanced section.
	advancedOpenRE = regexp.MustCompile(`(?s)<details[^>]*class="advanced"[^>]*>`)
)

// TestDevSplitSliderIsAdvancedOnly asserts the payout-share slider is rendered
// inside the collapsed Advanced settings section, and that it still spans the
// full 0-50 range. The control is a real setting, but it is not part of the
// default happy path: moving it back into the main form would put a
// revenue-share-shaped control in front of every operator.
func TestDevSplitSliderIsAdvancedOnly(t *testing.T) {
	ui := string(indexHTML)

	open := advancedOpenRE.FindStringIndex(ui)
	if open == nil {
		t.Fatal(`advanced settings section (<details class="advanced">) is gone`)
	}
	rel := strings.Index(ui[open[1]:], "</details>")
	if rel < 0 {
		t.Fatal("advanced settings section is not closed")
	}
	closeIdx := open[1] + rel

	slider := strings.Index(ui, `id="devsplit"`)
	if slider < 0 {
		t.Fatal(`payout-share slider (id="devsplit") is gone`)
	}
	if slider < open[0] || slider > closeIdx {
		t.Errorf("payout-share slider is rendered outside the advanced section (advanced spans %d-%d, slider at %d)",
			open[0], closeIdx, slider)
	}

	input := devsplitInputRE.FindString(ui)
	if input == "" {
		t.Fatal(`could not find the payout-share slider tag (id="devsplit")`)
	}
	for _, want := range []string{`min="0"`, `max="50"`} {
		if !strings.Contains(input, want) {
			t.Errorf("payout-share slider missing %s: %s", want, input)
		}
	}
}

// productionMintURLs and testnutMintURLs are spelled out here independently of
// deploy.go so a silent reclassification of a mint (production ←→ test) cannot
// pass by mutating both sides at once.
var productionMintURLs = []string{
	"https://mint.coinos.io",
	"https://mint.minibits.cash/Bitcoin",
	"https://mint.lnserver.com",
	"https://mint.macadamia.cash",
	"https://mint.westernbtc.com",
	"https://kashu.me",
	"https://mint.cubabitcoin.org",
}

var testnutMintURLs = []string{
	"https://nofee.testnut.cashu.space",
	"https://testnut.cashu.space",
}

// mintURLs decodes a default-mint JSON array into its URLs, which also proves
// the emitted payload is well-formed JSON before it reaches jq on the router.
func mintURLs(t *testing.T, raw string) []string {
	t.Helper()
	var mints []struct {
		URL       string `json:"url"`
		MinBalanc int    `json:"min_balance"`
	}
	if err := json.Unmarshal([]byte(raw), &mints); err != nil {
		t.Fatalf("default mints are not valid JSON: %v\n%s", err, raw)
	}
	urls := make([]string, 0, len(mints))
	for _, m := range mints {
		urls = append(urls, m.URL)
	}
	return urls
}

// TestNoTestnutMintsWithoutTestMode is the core guarantee of this change: an
// ordinary deploy must not push the testnut (test-network) mints into a real
// router's config.json. Those mints issue worthless test ecash, so accepting
// them on a production router lets a client pay with value the operator cannot
// redeem. The injection used to be unconditional.
func TestNoTestnutMintsWithoutTestMode(t *testing.T) {
	urls := mintURLs(t, defaultMintsJSON(false))

	for _, u := range urls {
		if strings.Contains(u, "testnut") {
			t.Errorf("testnut mint %q is injected without --test-mode", u)
		}
	}

	// The production set must survive intact: the change drops test mints, it
	// does not empty the default mint list.
	for _, want := range productionMintURLs {
		if !containsURL(urls, want) {
			t.Errorf("production mint %q missing from the default set", want)
		}
	}
	if len(urls) != len(productionMintURLs) {
		t.Errorf("default mint count = %d, want %d (production only): %v", len(urls), len(productionMintURLs), urls)
	}
}

// TestDefaultMintsIncludeTestnutOnlyUnderFlag covers both halves of the gate on
// the single function that builds the mint payload.
func TestDefaultMintsIncludeTestnutOnlyUnderFlag(t *testing.T) {
	off := mintURLs(t, defaultMintsJSON(false))
	on := mintURLs(t, defaultMintsJSON(true))

	for _, u := range testnutMintURLs {
		if containsURL(off, u) {
			t.Errorf("test-mode off: testnut mint %q present", u)
		}
		if !containsURL(on, u) {
			t.Errorf("test-mode on: testnut mint %q missing", u)
		}
	}
	for _, u := range productionMintURLs {
		if !containsURL(on, u) {
			t.Errorf("test-mode on: production mint %q missing", u)
		}
	}
	if len(on) != len(productionMintURLs)+len(testnutMintURLs) {
		t.Errorf("test-mode on: mint count = %d, want %d", len(on), len(productionMintURLs)+len(testnutMintURLs))
	}
}

// TestConfigCommandGatesTestnut proves the gate is wired into the command that
// actually rewrites /etc/tollgate/config.json, not only into a helper that a
// caller could sidestep.
func TestConfigCommandGatesTestnut(t *testing.T) {
	off := cfgConfigCommand(0, "1.0000", "0.0000", false)
	if strings.Contains(off, "testnut") {
		t.Errorf("config command without test-mode still carries a testnut mint:\n%s", off)
	}

	on := cfgConfigCommand(0, "1.0000", "0.0000", true)
	for _, u := range testnutMintURLs {
		if !strings.Contains(on, u) {
			t.Errorf("config command with test-mode missing testnut mint %q", u)
		}
	}
}

// TestParseOptionsTestMode pins the flag itself: the only way to reach the test
// mints is an explicit --test-mode on the command line, and the default is off.
func TestParseOptionsTestMode(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{"no args", nil, false},
		{"explicit off", []string{"--test-mode=false"}, false},
		{"explicit on", []string{"--test-mode"}, true},
		{"explicit on with value", []string{"--test-mode=true"}, true},
		{"single dash", []string{"-test-mode"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := parseOptions(tc.args)
			if err != nil {
				t.Fatalf("parseOptions(%v): %v", tc.args, err)
			}
			if opts.testMode != tc.want {
				t.Errorf("parseOptions(%v).testMode = %v, want %v", tc.args, opts.testMode, tc.want)
			}
		})
	}
}

// TestDeployRequestHasNoTestModeField pins test-mode as a process-level switch.
// The wizard binds to loopback, but its API is still reachable by anything
// running on the host; a client-supplied test-mode field would let a single
// request opt itself into test mints.
func TestDeployRequestHasNoTestModeField(t *testing.T) {
	typ := reflect.TypeOf(deployRequest{})
	for i := 0; i < typ.NumField(); i++ {
		name := strings.ToLower(strings.Split(typ.Field(i).Tag.Get("json"), ",")[0])
		if strings.Contains(name, "test") {
			t.Errorf("deployRequest field %s is exposed on the wire as %q; test-mode must be a process flag",
				typ.Field(i).Name, name)
		}
	}
}

func containsURL(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}
