package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Radio is a single wireless radio section discovered via ubus/uci.
type Radio struct {
	Name     string // "radio0", "radio1", …
	Up       bool   // ubus .up (true when the radio is up)
	Disabled bool   // uci .disabled='1' OR ubus .up==false
}

// Iface is a single wireless interface discovered via ubus/uci.
type Iface struct {
	Device  string // interface device name (ifname), e.g. "wlan0"
	Radio   string // owning radio section, e.g. "radio0"
	Mode    string // "ap", "sta", "client", "adhoc", …
	Section string // uci wifi-iface section name, e.g. "default_radio0" ("" for ubus-only)
}

// wirelessStatus is the parsed result of ubus/uci wireless enumeration.
type wirelessStatus struct {
	Radios      []Radio
	Ifaces      []Iface
	AnyDisabled bool
}

// parseWirelessStatus parses `ubus call network.wireless status` JSON and
// `uci show wireless` text into a driver-agnostic enumeration of radios and
// interfaces. It is deliberately tolerant of shape variants: ubus output may
// be empty, malformed, or keyed differently; uci lines may omit ifname or
// mode. Both sources are merged — ubus is authoritative for up/disabled and
// ifname, uci fills gaps (radio sections, disabled flag, iface device/mode).
func parseWirelessStatus(ubusJSON, uciOut string) wirelessStatus {
	st := wirelessStatus{}
	radioIdx := map[string]int{} // name -> index into st.Radios
	ifaceSeen := map[string]bool{}

	addRadio := func(name string) *Radio {
		if i, ok := radioIdx[name]; ok {
			return &st.Radios[i]
		}
		radioIdx[name] = len(st.Radios)
		st.Radios = append(st.Radios, Radio{Name: name})
		return &st.Radios[len(st.Radios)-1]
	}
	addIface := func(dev, radio, mode, section string) {
		if dev == "" || ifaceSeen[dev] {
			return
		}
		ifaceSeen[dev] = true
		st.Ifaces = append(st.Ifaces, Iface{Device: dev, Radio: radio, Mode: mode, Section: section})
	}

	// --- ubus JSON ---
	if strings.TrimSpace(ubusJSON) != "" {
		var status map[string]struct {
			Up     bool `json:"up"`
			Config struct {
				Iface []struct {
					Ifname string `json:"ifname"`
					Mode   string `json:"mode"`
				} `json:"iface"`
			} `json:"config"`
		}
		if err := json.Unmarshal([]byte(ubusJSON), &status); err == nil {
			for name, radio := range status {
				r := addRadio(name)
				r.Up = radio.Up
				if !radio.Up {
					r.Disabled = true
				}
				for _, iface := range radio.Config.Iface {
					addIface(iface.Ifname, name, iface.Mode, "")
				}
			}
		}
	}

	// --- uci show wireless ---
	// radio sections: wireless.radioN=wifi-device
	// iface sections: wireless.default_radioN=wifi-iface with option device/ifname/mode
	uciDisabled := map[string]bool{}
	uciIfaceDevice := map[string]string{} // section -> device (radio name)
	uciIfaceIfname := map[string]string{} // section -> ifname
	uciIfaceMode := map[string]string{}   // section -> mode
	for _, line := range strings.Split(uciOut, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// wireless.<section>=<type>
		if eq := strings.Index(line, "="); eq > 0 {
			key := line[:eq]
			val := strings.Trim(strings.TrimSpace(line[eq+1:]), "'\"")
			parts := strings.Split(key, ".")
			if len(parts) < 2 || parts[0] != "wireless" {
				continue
			}
			section := parts[1]
			if len(parts) == 2 {
				// wireless.radioN=wifi-device
				if val == "wifi-device" {
					addRadio(section)
				}
				continue
			}
			// wireless.<section>.<option>=<value>
			option := parts[2]
			switch option {
			case "disabled":
				if val == "1" {
					uciDisabled[section] = true
				}
			case "device":
				uciIfaceDevice[section] = val
			case "ifname":
				uciIfaceIfname[section] = val
			case "mode":
				uciIfaceMode[section] = val
			}
		}
	}
	// Apply uci disabled flags to radios.
	for name := range uciDisabled {
		if i, ok := radioIdx[name]; ok {
			st.Radios[i].Disabled = true
		} else {
			r := addRadio(name)
			r.Disabled = true
		}
	}
	// Build ifaces from uci iface sections (device -> radio, ifname -> dev).
	for section, radio := range uciIfaceDevice {
		dev := uciIfaceIfname[section]
		if dev == "" {
			// No ifname option — cannot derive a concrete device; skip.
			continue
		}
		mode := uciIfaceMode[section]
		addIface(dev, radio, mode, section)
	}

	// anyDisabled: any radio disabled OR zero ifaces with a disabled radio.
	for _, r := range st.Radios {
		if r.Disabled {
			st.AnyDisabled = true
			break
		}
	}

	return st
}

// iwinfoScanCommand builds the per-device iwinfo scan command. The device
// name is validated against iwIfaceNameRe because it is interpolated into an
// SSH command string (it originates from ubus/uci output on the router).
func iwinfoScanCommand(dev string) (string, error) {
	if !iwIfaceNameRe.MatchString(dev) {
		return "", fmt.Errorf("invalid interface name %q", dev)
	}
	return "iwinfo " + dev + " scan 2>&1", nil
}

// isAPMode reports whether a wireless interface mode is an AP/Master mode
// (which cannot scan for neighbouring networks).
func isAPMode(mode string) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "ap", "master":
		return true
	}
	return false
}

// classifyScanFailure classifies a scan command's output + error into a
// human-readable failure reason, or "" when the output looks like real scan
// results. It covers: non-zero exit status, empty output, ash's "not found"
// (in addition to the existing scanFailedHeuristic strings), and the
// existing heuristic.
func classifyScanFailure(out string, err error) string {
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "exit status") || strings.Contains(msg, "exited with status") {
			return "non-zero exit status: " + msg
		}
		return "command error: " + msg
	}
	if strings.TrimSpace(out) == "" {
		return "command produced no output"
	}
	if strings.Contains(out, "not found") {
		return "command not found: " + truncate(out, 120)
	}
	if scanFailedHeuristic(out) {
		return "scan failed: " + truncate(out, 120)
	}
	return ""
}
