package main

import (
	"strings"
	"testing"
)

// ubusWirelessFullFixture is real `ubus call network.wireless status` output
// shape: radios keyed by "radio0"/"radio1", each with .up and .config.iface
// entries carrying {ifname, mode}.
const ubusWirelessFullFixture = `{
  "radio0": {
    "up": true,
    "config": {
      "iface": [
        {"ifname": "wlan0", "mode": "ap"},
        {"ifname": "wlan1", "mode": "sta"}
      ]
    }
  },
  "radio1": {
    "up": true,
    "config": {
      "iface": [
        {"ifname": "wlan2", "mode": "ap"}
      ]
    }
  }
}`

// ubusWirelessDownFixture has a radio reporting up=false (disabled/down).
const ubusWirelessDownFixture = `{"radio0":{"up":false,"config":{"iface":[{"ifname":"wlan0","mode":"ap"}]}}}`

// uciWirelessFixture is real `uci show wireless` output: radio sections
// (wireless.radioN=wifi-device) and iface sections (wireless.default_radioN=
// wifi-iface) with option device/ifname/mode/disabled.
const uciWirelessFixture = `wireless.radio0=wifi-device
wireless.radio0.type='mac80211'
wireless.radio0.channel='36'
wireless.radio0.disabled='0'
wireless.radio1=wifi-device
wireless.radio1.type='mac80211'
wireless.radio1.disabled='1'
wireless.default_radio0=wifi-iface
wireless.default_radio0.device='radio0'
wireless.default_radio0.ifname='wlan0'
wireless.default_radio0.mode='ap'
wireless.default_radio1=wifi-iface
wireless.default_radio1.device='radio1'
wireless.default_radio1.ifname='wlan1'
wireless.default_radio1.mode='sta'`

// uciWirelessNoIfnameFixture is a degraded variant: an iface section with no
// ifname option (must not produce a bogus empty-device iface).
const uciWirelessNoIfnameFixture = `wireless.radio0=wifi-device
wireless.default_radio0=wifi-iface
wireless.default_radio0.device='radio0'
wireless.default_radio0.mode='ap'`

func findRadio(st wirelessStatus, name string) (Radio, bool) {
	for _, r := range st.Radios {
		if r.Name == name {
			return r, true
		}
	}
	return Radio{}, false
}

func findIface(st wirelessStatus, dev string) (Iface, bool) {
	for _, i := range st.Ifaces {
		if i.Device == dev {
			return i, true
		}
	}
	return Iface{}, false
}

// TestParseWirelessStatusUbusFull pins the ubus-only path: radios and ifaces
// are extracted from the ubus JSON shape, with up/disabled and mode captured.
func TestParseWirelessStatusUbusFull(t *testing.T) {
	st := parseWirelessStatus(ubusWirelessFullFixture, "")
	if len(st.Radios) != 2 {
		t.Fatalf("radios = %v; want 2 (radio0, radio1)", st.Radios)
	}
	if len(st.Ifaces) != 3 {
		t.Fatalf("ifaces = %v; want 3 (wlan0, wlan1, wlan2)", st.Ifaces)
	}
	if st.AnyDisabled {
		t.Error("AnyDisabled = true; want false (both radios up)")
	}
	r0, ok := findRadio(st, "radio0")
	if !ok || !r0.Up || r0.Disabled {
		t.Errorf("radio0 = %+v; want up=true disabled=false", r0)
	}
	if i, ok := findIface(st, "wlan1"); !ok || i.Mode != "sta" || i.Radio != "radio0" {
		t.Errorf("wlan1 = %+v; want mode=sta radio=radio0", i)
	}
}

// TestParseWirelessStatusUciOnly pins the uci-only path (empty ubus output):
// radios and ifaces come from uci lines, disabled flag from .disabled='1'.
func TestParseWirelessStatusUciOnly(t *testing.T) {
	st := parseWirelessStatus("", uciWirelessFixture)
	if len(st.Radios) != 2 {
		t.Fatalf("radios = %v; want 2", st.Radios)
	}
	if len(st.Ifaces) != 2 {
		t.Fatalf("ifaces = %v; want 2 (wlan0, wlan1)", st.Ifaces)
	}
	if !st.AnyDisabled {
		t.Error("AnyDisabled = false; want true (radio1 disabled='1')")
	}
	r1, ok := findRadio(st, "radio1")
	if !ok || !r1.Disabled {
		t.Errorf("radio1 = %+v; want disabled=true", r1)
	}
	if i, ok := findIface(st, "wlan1"); !ok || i.Mode != "sta" || i.Radio != "radio1" {
		t.Errorf("wlan1 = %+v; want mode=sta radio=radio1", i)
	}
}

// TestParseWirelessStatusUbusDown pins the ubus up=false → disabled mapping.
func TestParseWirelessStatusUbusDown(t *testing.T) {
	st := parseWirelessStatus(ubusWirelessDownFixture, "")
	if !st.AnyDisabled {
		t.Error("AnyDisabled = false; want true (radio0 up=false)")
	}
	r0, ok := findRadio(st, "radio0")
	if !ok || !r0.Disabled {
		t.Errorf("radio0 = %+v; want disabled=true", r0)
	}
}

// TestParseWirelessStatusDegraded pins degraded variants: empty ubus output
// and a uci iface section missing its ifname must not fabricate an iface.
func TestParseWirelessStatusDegraded(t *testing.T) {
	st := parseWirelessStatus("", uciWirelessNoIfnameFixture)
	if len(st.Radios) != 1 {
		t.Fatalf("radios = %v; want 1 (radio0)", st.Radios)
	}
	if len(st.Ifaces) != 0 {
		t.Fatalf("ifaces = %v; want 0 (ifname missing)", st.Ifaces)
	}
	if st.AnyDisabled {
		t.Error("AnyDisabled = true; want false")
	}
	// Fully empty inputs must yield an empty status, not panic.
	st = parseWirelessStatus("", "")
	if len(st.Radios) != 0 || len(st.Ifaces) != 0 || st.AnyDisabled {
		t.Errorf("empty inputs = %+v; want zero status", st)
	}
	// Garbage ubus JSON must be tolerated (defensive parse).
	st = parseWirelessStatus("not json", uciWirelessFixture)
	if len(st.Radios) != 2 || len(st.Ifaces) != 2 {
		t.Errorf("garbage ubus + valid uci = %+v; want uci-derived radios/ifaces", st)
	}
}

// TestClassifyScanFailure pins the failure classifier: non-zero exit, empty
// output, ash "not found", and existing heuristic strings.
func TestClassifyScanFailure(t *testing.T) {
	cases := []struct {
		name    string
		out     string
		err     error
		wantSub string
		wantOK  bool // wantOK=true means classify returns "" (not a failure)
	}{
		{"non-zero exit", "some output", errExitStatus{}, "non-zero exit", false},
		{"empty output", "", nil, "no output", false},
		{"ash not found", "ash: iwinfo: not found", nil, "not found", false},
		{"command not found", "iwinfo: command not found", nil, "not found", false},
		{"usage text", iwUsageFixture, nil, "scan failed", false},
		{"real scan output", iwScanOutputFixture, nil, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyScanFailure(tc.out, tc.err)
			if tc.wantOK {
				if got != "" {
					t.Errorf("classifyScanFailure(%q, %v) = %q; want empty (not a failure)", tc.out, tc.err, got)
				}
				return
			}
			if !strings.Contains(got, tc.wantSub) {
				t.Errorf("classifyScanFailure(%q, %v) = %q; want substring %q", tc.out, tc.err, got, tc.wantSub)
			}
		})
	}
}

// errExitStatus is a stand-in for the *ssh.ExitError CombinedOutput returns on
// a non-zero remote exit (its Error() contains "exit status N").
type errExitStatus struct{}

func (errExitStatus) Error() string { return "Process exited with status 1" }

// TestIwinfoScanCommand pins the per-device iwinfo scan command builder: a
// concrete device name is embedded, stderr is merged, and the name is
// sanitized (it comes from ubus/uci output on the router).
func TestIwinfoScanCommand(t *testing.T) {
	cases := []struct {
		name    string
		dev     string
		wantSub string
		wantErr bool
	}{
		{"normal dev", "wlan0", "iwinfo wlan0 scan 2>&1", false},
		{"phy-ap dev", "phy0-ap0", "iwinfo phy0-ap0 scan 2>&1", false},
		{"dev with space", "wlan0; rm -rf /", "", true},
		{"empty dev", "", "", true},
		{"dev with metachar", "wlan0$(reboot)", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := iwinfoScanCommand(tc.dev)
			if tc.wantErr {
				if err == nil {
					t.Errorf("iwinfoScanCommand(%q) = %q, nil; want error", tc.dev, got)
				}
				return
			}
			if err != nil {
				t.Errorf("iwinfoScanCommand(%q) unexpected error: %v", tc.dev, err)
				return
			}
			if got != tc.wantSub {
				t.Errorf("iwinfoScanCommand(%q) = %q; want %q", tc.dev, got, tc.wantSub)
			}
			if strings.Contains(got, "2>/dev/null") {
				t.Errorf("iwinfoScanCommand(%q) = %q must merge stderr (2>&1)", tc.dev, got)
			}
		})
	}
}

// TestIsAPMode pins the AP/Master skip classification used to avoid scanning
// interfaces that cannot scan.
func TestIsAPMode(t *testing.T) {
	for _, m := range []string{"ap", "AP", "master", "Master"} {
		if !isAPMode(m) {
			t.Errorf("isAPMode(%q) = false; want true", m)
		}
	}
	for _, m := range []string{"sta", "client", "managed", "adhoc", "monitor", ""} {
		if isAPMode(m) {
			t.Errorf("isAPMode(%q) = true; want false", m)
		}
	}
}
