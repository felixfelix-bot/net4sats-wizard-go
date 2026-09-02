package main

import (
	"strings"
)

// enableRadiosCommand builds ONE BusyBox-safe command that enables the given
// radio sections AND their owned wifi-iface sections via uci. Radio and iface
// section names come from DISCOVERED sections (parseWirelessStatus) — never
// hardcoded radio0/radio1. Every interpolated name is injection-guarded
// against iwIfaceNameRe (the same strict charset used at wifi_scan_ubus.go);
// a name that fails the guard is silently skipped so a hostile/garbage name
// can never smuggle a second command into the SSH session. Empty input yields
// an empty command (no-op).
func enableRadiosCommand(radioNames, ifaceNames []string) string {
	var parts []string
	for _, name := range radioNames {
		if !iwIfaceNameRe.MatchString(name) {
			continue
		}
		parts = append(parts, "uci set wireless."+name+".disabled='0'")
	}
	for _, name := range ifaceNames {
		if !iwIfaceNameRe.MatchString(name) {
			continue
		}
		parts = append(parts, "uci set wireless."+name+".disabled='0'")
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " && ")
}

// radiosDisabledResponse builds the JSON response body for the
// radios_disabled scan branch: the discovered radio section names alongside
// radios_disabled:true (so the UI can render a consent-gated enable button
// labeled with the REAL radio names).
func radiosDisabledResponse(radioNames []string, debug string) map[string]any {
	return map[string]any{
		"error":           "radios disabled or down — enable them to scan for networks",
		"radios_disabled": true,
		"radios":          radioNames,
		"debug":           debug,
	}
}
