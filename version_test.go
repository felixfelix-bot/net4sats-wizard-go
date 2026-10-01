package main

import (
	"regexp"
	"testing"
)

// TestHandleVersionArg pins the CLI contract used by the release
// verification path (curl the asset, chmod, `net4sats-wizard --version`):
// --version and version are handled, anything else falls through to
// serving.
func TestHandleVersionArg(t *testing.T) {
	for _, arg := range []string{"--version", "version"} {
		if !handleVersionArg([]string{"net4sats-wizard", arg}) {
			t.Errorf("handleVersionArg(%q) = false, want true", arg)
		}
	}
	for _, args := range [][]string{
		{"net4sats-wizard"},
		{"net4sats-wizard", "--help"},
		{"net4sats-wizard", "-v"},
	} {
		if handleVersionArg(args) {
			t.Errorf("handleVersionArg(%q) = true, want false (only arg[1] decides)", args)
		}
	}
	// arg[1] alone decides: extra trailing args still count as handled.
	if !handleVersionArg([]string{"net4sats-wizard", "--version", "extra"}) {
		t.Error("handleVersionArg(--version extra) = false, want true")
	}
}

// TestWizardVersionShape keeps the stamped release tag well-formed so the
// binary can be matched against the GitHub release it came from.
func TestWizardVersionShape(t *testing.T) {
	re := regexp.MustCompile(`^v0\.7\.0-alpha\d+$`)
	if !re.MatchString(wizardVersion) {
		t.Errorf("wizardVersion = %q, want v0.7.0-alpha<N>", wizardVersion)
	}
}
