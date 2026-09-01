package main

import (
	"strings"
	"testing"
)

// iwUsageFixture is the real output of `iw dev scan` (INVALID syntax — iw
// has no "dev scan" subcommand without an interface name). iw prints its
// usage text to STDOUT with exit 0, so it slips past exit-code checks and
// historically reached the UI as fake "results".
const iwUsageFixture = "Usage:\tiw [options] command\nOptions:\n\t--debug\t\tenable netlink debugging\n\t--version\tshow version (6.17)\nCommands:\n\tdev <devname> ap start <SSID> ..."

// iwDevListFixture is real `iw dev` output: phy blocks whose top-level
// interface lines read "Interface <name>" (tab-indented).
const iwDevListFixture = "phy#0\n\tUnnamed/non-netdev interface\n\t\twdev 0x2\n\t\taddr e4:70:b8:c1:5a:d0\n\t\ttype P2P-device\n\tInterface wlp58s0\n\t\tifindex 4\n\t\twdev 0x1\n\t\taddr e4:70:b8:c1:5a:cf\n\t\tssid NOBIOS\n\t\ttype managed\n\t\tchannel 36 (5180 MHz), width: 80 MHz, center1: 5210 MHz\n\t\ttxpower 22.00 dBm"

// iwDevListMultiPhyFixture has two phy blocks with multiple interfaces to
// pin ordering and dedup behaviour.
const iwDevListMultiPhyFixture = "phy#0\n\tInterface wlan0\n\t\tifindex 4\n\tInterface wlan1\n\t\tifindex 5\nphy#1\n\tInterface phy1-ap0\n\t\tifindex 6"

// iwScanOutputFixture is real `iw dev wlan0 scan` output (BSS blocks).
const iwScanOutputFixture = "BSS aa:bb:cc:dd:ee:ff(on wlp58s0) -- associated\n\tfreq: 5180\n\tSSID: NOBIOS\n\tcapability: ESS Privacy SpectrumMgmt\nBSS 11:22:33:44:55:66(on wlp58s0)\n\tfreq: 2412\n\tSSID: NeighborNet\n\tcapability: ESS"

// TestScanFailedHeuristicUsage pins the regression behind the fabricated
// scan-results bug: iw usage text MUST be classified as a failure so it is
// never parsed or echoed to the UI as "results" — in every scan strategy.
func TestScanFailedHeuristicUsage(t *testing.T) {
	if !scanFailedHeuristic(iwUsageFixture) {
		t.Errorf("scanFailedHeuristic(iw usage text) = false; want true — usage text is not scan results")
	}
	if !scanFailedHeuristic("Usage: iw [options] command") {
		t.Errorf("scanFailedHeuristic(\"Usage: ...\") = false; want true")
	}
	// Sanity: real scan output must keep passing the heuristic.
	if scanFailedHeuristic(iwScanOutputFixture) {
		t.Error("scanFailedHeuristic(real iw scan output) = true; want false")
	}
	if scanFailedHeuristic(realIwinfoScanFixture) {
		t.Error("scanFailedHeuristic(real iwinfo scan output) = true; want false")
	}
}

// realIwinfoScanFixture is real `iwinfo scan` output (OpenWrt 25.12,
// GL-MT3000) for cross-parser sanity in the tests above.
const realIwinfoScanFixture = "Cell 01 - Address: AA:BB:CC:DD:EE:FF\n          ESSID: \"Net4Sats-F794\"\n          Mode: Master  Channel: 6 (2.4 GHz)\n          Signal: -45 dBm  Quality: 60/70\n          Encryption: WPA2 PSK (CCMP)\nCell 02 - Address: 11:22:33:44:55:66\n          ESSID: \"Neighbor\"\n          Mode: Master  Channel: 11 (2.4 GHz)\n          Signal: -67 dBm  Quality: 40/70\n          Encryption: WPA3 SAE (CCMP)"

// TestParseIwScanUsageText pins that usage text extracts ZERO SSIDs — the
// other half of the fabricated-results bug: even if usage text ever reached
// the parser, it must not produce network entries.
func TestParseIwScanUsageText(t *testing.T) {
	if got := parseIwScan(iwUsageFixture); len(got) != 0 {
		t.Errorf("parseIwScan(iw usage text) = %v; want 0 SSIDs", got)
	}
}

// TestParseIwScanRealOutput pins SSID extraction against real iw scan output,
// including dedup.
func TestParseIwScanRealOutput(t *testing.T) {
	got := parseIwScan(iwScanOutputFixture)
	if len(got) != 2 {
		t.Fatalf("parseIwScan(real scan) = %v; want 2 SSIDs", got)
	}
	if got[0].Name != "NOBIOS" || got[1].Name != "NeighborNet" {
		t.Errorf("parseIwScan names = %q,%q; want NOBIOS,NeighborNet (order of appearance)", got[0].Name, got[1].Name)
	}
}

// TestParseIwDevList covers the per-interface scan listing: extracting
// "Interface <name>" lines at the top level of each phy block.
func TestParseIwDevList(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"single phy one iface", iwDevListFixture, []string{"wlp58s0"}},
		{"multi phy multi iface", iwDevListMultiPhyFixture, []string{"wlan0", "wlan1", "phy1-ap0"}},
		{"empty output", "", nil},
		{"iw usage text — no ifaces", iwUsageFixture, nil},
		{"no Interface lines", "phy#0\n\taddr 00:11:22:33:44:55", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseIwDevList(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("parseIwDevList(%q) = %v; want %v", tc.in, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("parseIwDevList(%q)[%d] = %q; want %q", tc.in, i, got[i], tc.want[i])
				}
			}
		})
	}
}

// TestIwInterfaceScanCommand pins the command construction for the
// per-interface last-resort scan: a concrete interface name must be
// embedded (the old `iw dev scan` was invalid iw syntax — no such
// subcommand without an interface), stderr must be merged so real
// router-side errors reach the UI, and the interface name must be
// sanitized (it comes from `iw dev` output on the router).
func TestIwInterfaceScanCommand(t *testing.T) {
	cases := []struct {
		name    string
		iface   string
		wantSub string
		wantErr bool
	}{
		{"normal iface", "wlan0", "iw dev wlan0 scan 2>&1", false},
		{"phy-ap iface", "phy0-ap0", "iw dev phy0-ap0 scan 2>&1", false},
		{"iface with space", "wlan0; rm -rf /", "", true},
		{"empty iface", "", "", true},
		{"iface with shell metachar", "wlan0$(reboot)", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := iwInterfaceScanCommand(tc.iface)
			if tc.wantErr {
				if err == nil {
					t.Errorf("iwInterfaceScanCommand(%q) = %q, nil; want error", tc.iface, got)
				}
				return
			}
			if err != nil {
				t.Errorf("iwInterfaceScanCommand(%q) unexpected error: %v", tc.iface, err)
				return
			}
			if got != tc.wantSub {
				t.Errorf("iwInterfaceScanCommand(%q) = %q; want %q", tc.iface, got, tc.wantSub)
			}
			if strings.Contains(got, "2>/dev/null") {
				t.Errorf("iwInterfaceScanCommand(%q) = %q must merge stderr (2>&1), not discard it", tc.iface, got)
			}
		})
	}
}