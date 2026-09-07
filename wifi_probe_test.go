package main

import "testing"

// TestMapProbeResult covers the router-side probe token → probe-code mapping.
// The probe script emits a trailing AUTH_RESULT=... token; mapProbeResult must
// find it (even with preceding log noise) and translate it correctly. Unknown /
// missing tokens degrade to ProbeTimeout (ambiguous), never a false "ok" or
// false "wrong-password".
func TestMapProbeResult(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"ok", "some log\nAUTH_RESULT=OK\n", ProbeOK},
		{"ok-upper", "AUTH_RESULT=OK", ProbeOK},
		{"wrong-key", "noise\nAUTH_RESULT=WRONGKEY\n", ProbeWrongKey},
		{"no-ap", "AUTH_RESULT=NOAP", ProbeNoAP},
		{"unavailable", "AUTH_RESULT=UNAVAILABLE", ProbeUnavailable},
		{"timeout-explicit", "AUTH_RESULT=TIMEOUT", ProbeTimeout},
		{"token-in-middle-uses-last", "AUTH_RESULT=TIMEOUT\ntail\nAUTH_RESULT=OK\n", ProbeOK},
		{"garbage-token", "AUTH_RESULT=BANANA", ProbeTimeout},
		{"no-token", "nothing relevant here", ProbeTimeout},
		{"empty", "", ProbeTimeout},
		{"lowercase", "AUTH_RESULT=ok", ProbeOK},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := mapProbeResult(c.in); got != c.want {
				t.Errorf("mapProbeResult(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// TestProbeScriptEscapesNonbreaking ensures the generated probe script shell-escapes
// SSID and password (single quotes, embedded single quote), so a password such as
// "don't" cannot break out of the wpa_supplicant network block.
func TestProbeScriptEscapesNonbreaking(t *testing.T) {
	script := buildWifiProbeScript("HotelCentral", "don'tuse'this")
	// The escaped password must appear single-quoted with the apostrophe turned
	// into the four-char sequence '\'' (shell-safe), never bare inside quotes.
	if !contains(script, `don'\''tuse'\''this`) {
		t.Errorf("expected apostrophe-escaped password in probe script, got:\n%s", script)
	}
	if contains(script, "psk=don't") {
		t.Errorf("bare password with apostrophe leaked unquoted into probe script")
	}
	// SSID must be single-quoted too.
	if !contains(script, `ssid='HotelCentral'`) {
		t.Errorf("expected quoted ssid in probe script")
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
