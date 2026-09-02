package main

import (
	"strings"
	"testing"
)

// TestEnableRadiosCommand pins the consent-gated radio-enable command builder:
// it must emit uci enables for discovered radio sections AND their owned
// wifi-iface sections, reject injection names, and no-op on empty input.
func TestEnableRadiosCommand(t *testing.T) {
	cases := []struct {
		name           string
		radios         []string
		ifaces         []string
		mustContain    []string
		mustNotContain []string
		wantEmpty      bool
	}{
		{
			name:   "discovered radios + owned ifaces",
			radios: []string{"radio0", "radio1"},
			ifaces: []string{"default_radio0", "default_radio1"},
			mustContain: []string{
				"uci set wireless.radio0.disabled='0'",
				"uci set wireless.radio1.disabled='0'",
				"uci set wireless.default_radio0.disabled='0'",
				"uci set wireless.default_radio1.disabled='0'",
			},
		},
		{
			name:        "single radio no ifaces",
			radios:      []string{"radio0"},
			mustContain: []string{"uci set wireless.radio0.disabled='0'"},
		},
		{
			name:           "injection radio name rejected",
			radios:         []string{"radio0; rm -rf /"},
			mustNotContain: []string{"rm -rf", "radio0;"},
		},
		{
			name:           "injection iface name rejected",
			ifaces:         []string{"default_radio0$(reboot)"},
			mustNotContain: []string{"reboot", "$("},
		},
		{
			name:      "empty input no-op",
			wantEmpty: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := enableRadiosCommand(tc.radios, tc.ifaces)
			if tc.wantEmpty {
				if got != "" {
					t.Errorf("enableRadiosCommand(%v, %v) = %q; want empty", tc.radios, tc.ifaces, got)
				}
				return
			}
			for _, want := range tc.mustContain {
				if !strings.Contains(got, want) {
					t.Errorf("enableRadiosCommand(%v, %v) missing %q; got %q", tc.radios, tc.ifaces, want, got)
				}
			}
			for _, notWant := range tc.mustNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("enableRadiosCommand(%v, %v) must NOT contain %q; got %q", tc.radios, tc.ifaces, notWant, got)
				}
			}
		})
	}
}

// TestRadiosDisabledResponse pins the radios_disabled scan-branch response
// shape: discovered radio names alongside radios_disabled:true.
func TestRadiosDisabledResponse(t *testing.T) {
	resp := radiosDisabledResponse([]string{"radio0", "radio1"}, "debug text")
	if resp["radios_disabled"] != true {
		t.Errorf("radios_disabled = %v; want true", resp["radios_disabled"])
	}
	radios, ok := resp["radios"].([]string)
	if !ok {
		t.Fatalf("radios = %v (%T); want []string", resp["radios"], resp["radios"])
	}
	if len(radios) != 2 || radios[0] != "radio0" || radios[1] != "radio1" {
		t.Errorf("radios = %v; want [radio0 radio1]", radios)
	}
	if resp["error"] == "" {
		t.Error("error field should be non-empty")
	}
}

// TestParseWirelessStatusIfaceSection pins that parseWirelessStatus captures
// the uci wifi-iface SECTION name (e.g. default_radio0) so the enable command
// can target owned iface sections — never hardcoded.
func TestParseWirelessStatusIfaceSection(t *testing.T) {
	st := parseWirelessStatus("", uciWirelessFixture)
	if i, ok := findIface(st, "wlan0"); !ok || i.Section != "default_radio0" {
		t.Errorf("wlan0 = %+v; want Section=default_radio0", i)
	}
	if i, ok := findIface(st, "wlan1"); !ok || i.Section != "default_radio1" {
		t.Errorf("wlan1 = %+v; want Section=default_radio1", i)
	}
}
