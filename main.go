// net4sats-wizard — cross-platform router onboarding wizard.
// Single binary: serves web UI + API, auto-discovers routers,
// deploys net4sats over SSH.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lnAddrRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	// lnurlRe matches a raw LNURL: the "lnurl1" bech32 separator followed by
	// at least 6 lowercase bech32 data characters (covers the 6-char checksum).
	// Real LNURLs are far longer; this is a lenient plausibility gate.
	lnurlRe = regexp.MustCompile(`^lnurl1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{6,}$`)
)

// validLightningAddress reports whether s is a plausible Lightning payout
// target. Two forms are accepted:
//  1. Lightning address — email-shaped: localpart@domain.tld
//  2. Raw LNURL — bech32-encoded: lnurl1<data>
//
// The Lightning target is a required MVP field — payouts route here, so an
// empty/invalid value would silently send payments nowhere. The check is
// intentionally lenient; actual resolution happens at payout time on the router.
func validLightningAddress(s string) bool {
	s = strings.TrimSpace(s)
	return lnAddrRe.MatchString(s) || lnurlRe.MatchString(s)
}

// ─── Listen address / port fallback ────────────────────────────
//
// BUG FIX (v0.7.0-alpha16): the wizard used to print its URL banner
// BEFORE binding and then log.Fatal on EADDRINUSE — an operator with a
// second wizard instance running saw "running on http://localhost:8099"
// followed by an unexplained "bind: address already in use". Now:
//
//   - the socket is bound FIRST via net.Listen, the URL is printed only
//     after the bind succeeds (a success banner always means a live socket);
//   - PORT overrides 8099 (e.g. PORT=8199); a typo'd PORT is a loud error,
//     not a silent fallback;
//   - on EADDRINUSE the wizard walks 8099→8109 in order before giving up;
//   - if every candidate is busy it prints how to find the holder
//     (ss -tlnp) and how to kill the old instance (pkill), then exits
//     non-zero.

const defaultPort = 8099

// listenPort returns the port to serve on: PORT if set, else 8099.
// PORT must be a valid TCP port number (1-65535) — anything else is an
// error. A typo'd PORT must be loud, not silently ignored.
func listenPort() (int, error) {
	raw := strings.TrimSpace(os.Getenv("PORT"))
	if raw == "" {
		return defaultPort, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid PORT %q: must be a port number between 1 and 65535", raw)
	}
	return port, nil
}

// fallbackPorts returns the ordered fallback candidates tried when the
// preferred port is busy (EADDRINUSE): the next ten ports, capped at 65535.
func fallbackPorts(preferred int) []int {
	var out []int
	for p := preferred + 1; p <= preferred+10 && p <= 65535; p++ {
		out = append(out, p)
	}
	return out
}

// pickPort binds the preferred address, falling back through candidates on
// failure (EADDRINUSE and friends). It returns the bound listener and the
// concrete address. On full failure it returns the preferred address in the
// error so callers can tell the operator exactly which port was busy.
func pickPort(preferred string, fallbacks []string) (net.Listener, string, error) {
	var lastErr error
	for _, addr := range append([]string{preferred}, fallbacks...) {
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			return ln, addr, nil
		}
		lastErr = err
	}
	return nil, "", fmt.Errorf("cannot bind %s (or any fallback %v): %w", preferred, fallbacks, lastErr)
}

// urlFor returns the operator-facing http:// URL for a bound address.
// A wildcard host renders as "localhost" so the banner stays copyable.
func urlFor(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "http://" + addr
	}
	if host == "" {
		host = "localhost"
	}
	return "http://" + net.JoinHostPort(host, port)
}

// ─── Job tracking ─────────────────────────────────────────────

type Step struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status string `json:"status"` // pending | running | done | failed
	Detail string `json:"detail,omitempty"`
}

type LogEntry struct {
	Time float64 `json:"time"`
	Msg  string  `json:"msg"`
}

type Job struct {
	mu     sync.Mutex
	IP     string     `json:"ip"`
	Status string     `json:"status"` // running | done | failed
	Step   int        `json:"step"`
	Steps  []Step     `json:"steps"`
	Log    []LogEntry `json:"log"`
	Error  string     `json:"error,omitempty"`
}

var (
	jobs      = make(map[string]*Job)
	jobsMutex sync.RWMutex
)

func newJob(ip string) *Job {
	return &Job{
		IP:     ip,
		Status: "running",
		Step:   0,
		Steps:  deploySteps(),
		Log:    []LogEntry{},
	}
}

func (j *Job) addLog(msg string) {
	j.mu.Lock()
	j.Log = append(j.Log, LogEntry{Time: float64(time.Now().Unix()), Msg: msg})
	j.mu.Unlock()
}

func (j *Job) setStep(i int, status, detail string) {
	j.mu.Lock()
	if i < len(j.Steps) {
		j.Step = i
		j.Steps[i].Status = status
		if detail != "" {
			j.Steps[i].Detail = detail
		}
	}
	j.mu.Unlock()
}

// ─── API handlers ─────────────────────────────────────────────

func handleScan(w http.ResponseWriter, r *http.Request) {
	routers := discoverRouters()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"routers": routers})
}

// wifiScanRequest is the JSON body for /api/wifi-scan.
type wifiScanRequest struct {
	IP       string `json:"ip"`
	Password string `json:"password"`
}

// wifiSSID represents a single SSID found during a WiFi scan.
type wifiSSID struct {
	Name       string `json:"name"`
	Encryption string `json:"encryption"`
	Signal     int    `json:"signal"` // dBm, e.g. -45 (0 if unknown)
}

// parseIwinfoScan parses `iwinfo scan` output and returns deduplicated SSIDs
// sorted by signal strength (strongest first). iwinfo output on OpenWrt:
//
//	wl0-sha0   ESSID: "MyWiFi"
//	          Mode: Master  Channel: 6 (2.4 GHz)
//	          Signal: -45 dBm  Quality: 70/70
//	          Encryption: WPA2 PSK (CCMP)
//
// Some versions prefix with "Cell 01 - Address: ..." instead of the interface name.
func parseIwinfoScan(output string) []wifiSSID {
	seen := map[string]bool{}
	ssids := []wifiSSID{}
	var currentName, currentEnc string
	var currentSignal int

	flush := func() {
		if currentName != "" && !seen[currentName] {
			seen[currentName] = true
			ssids = append(ssids, wifiSSID{Name: currentName, Encryption: currentEnc, Signal: currentSignal})
		}
		currentName = ""
		currentEnc = ""
		currentSignal = 0
	}

	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)

		// "Cell" prefix (some iwinfo versions) or a new ESSID line indicates a new block
		if strings.HasPrefix(trimmed, "Cell ") {
			flush()
			continue
		}

		// ESSID — may appear on the interface-name line or indented
		if idx := strings.Index(trimmed, "ESSID:"); idx >= 0 {
			// If we already have a name, this is a new block (interface-name-prefixed format)
			if currentName != "" {
				flush()
			}
			val := strings.TrimSpace(trimmed[idx+len("ESSID:"):])
			val = strings.Trim(val, "\"")
			if val != "" {
				currentName = val
			}
			continue
		}

		if strings.HasPrefix(trimmed, "Signal:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "Signal:"))
			fields := strings.Fields(val)
			if len(fields) > 0 {
				if dbm, err := strconv.Atoi(fields[0]); err == nil {
					currentSignal = dbm
				}
			}
			continue
		}

		if strings.HasPrefix(trimmed, "Encryption:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "Encryption:"))
			currentEnc = val
			continue
		}
	}
	flush()

	// Sort by signal strength descending (strongest = highest dBm first)
	sort.Slice(ssids, func(i, j int) bool {
		return ssids[i].Signal > ssids[j].Signal
	})

	return ssids
}

// parseIwScan parses `iw dev wlan0 scan` output (fallback when iwinfo absent).
// iw output uses:
//
//	BSS aa:bb:cc:dd:ee:ff on wlan0
//	    freq: 2412
//	    SSID: NetworkName
//	    ...
//	    capability: ...
//	    * primary channel: 1
func parseIwScan(output string) []wifiSSID {
	seen := map[string]bool{}
	ssids := []wifiSSID{}
	var currentName string

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "BSS ") {
			if currentName != "" && !seen[currentName] {
				seen[currentName] = true
				ssids = append(ssids, wifiSSID{Name: currentName, Encryption: "unknown"})
			}
			currentName = ""
			continue
		}
		if strings.HasPrefix(line, "SSID:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))
			if val != "" {
				currentName = val
			}
		}
	}
	if currentName != "" && !seen[currentName] {
		seen[currentName] = true
		ssids = append(ssids, wifiSSID{Name: currentName, Encryption: "unknown"})
	}
	return ssids
}

// scanFailedHeuristic returns true if the output text indicates a scan
// failure rather than real scan results. Checks for common error strings
// emitted by iwinfo and iw on OpenWrt when a device doesn't exist, is in
// the wrong mode, or the command is unavailable.
func scanFailedHeuristic(out string) bool {
	return strings.TrimSpace(out) == "" ||
		strings.Contains(out, "command not found") ||
		strings.Contains(out, "No such device") ||
		strings.Contains(out, "No such wireless device") ||
		strings.Contains(out, "Operation not supported") ||
		strings.Contains(out, "Operation not permitted") ||
		strings.Contains(out, "Device or resource busy") ||
		strings.Contains(out, "Usage:")
}

// parseIwDevList extracts interface names from `iw dev` output: lines of
// the form "	Interface <name>" (single-tab indent) at the top level of
// each phy block, e.g.
//
//	phy#0
//		Unnamed/non-netdev interface
//			wdev 0x2
//		Interface wlp58s0      ← parsed
//			ifindex 4
//
// Deeper-indented lines, usage text, and garbage yield nothing. Used by the
// last-resort scan to enumerate concrete interfaces — the old code ran
// `iw dev scan`, which is INVALID iw syntax (no such subcommand without an
// interface name).
func parseIwDevList(output string) []string {
	seen := map[string]bool{}
	var ifaces []string
	for _, line := range strings.Split(output, "\n") {
		if !strings.HasPrefix(line, "	Interface ") {
			continue
		}
		name := strings.TrimSpace(strings.TrimPrefix(line, "	Interface "))
		if name == "" || strings.ContainsAny(name, " 	") || seen[name] {
			continue
		}
		seen[name] = true
		ifaces = append(ifaces, name)
	}
	return ifaces
}

// iwIfaceNameRe is the strict charset of real wireless interface names
// (wlan0, phy0-ap0, wlp58s0, …). Anything else — shell metachars,
// whitespace, empty — is rejected so a name parsed from router output can
// never smuggle a second command into the SSH session.
var iwIfaceNameRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// iwInterfaceScanCommand builds the last-resort per-interface scan command:
// `iw dev <iface> scan 2>&1`. stderr is MERGED (not discarded) so real
// router-side errors reach the UI; the interface name is validated against
// iwIfaceNameRe because it is interpolated into an SSH command string.
func iwInterfaceScanCommand(iface string) (string, error) {
	if !iwIfaceNameRe.MatchString(iface) {
		return "", fmt.Errorf("invalid interface name %q", iface)
	}
	return "iw dev " + iface + " scan 2>&1", nil
}

// scanDebug assembles the "debug" field for scan responses: the
// strategy-by-strategy debug log plus a tail of the last raw router output.
// The old code sent only truncate(lastOut, 200), hiding which strategies
// ran and what the router actually said — including the router-side
// errors that stderr capture now surfaces.
func scanDebug(debugLog, lastOut string) string {
	s := strings.TrimSpace(debugLog)
	if tail := truncate(lastOut, 400); tail != "" {
		if s != "" {
			s += "\n"
		}
		s += "router output: " + tail
	}
	return s
}

func handleWifiScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}
	var req wifiScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	// Password is optional: a fresh-reset OpenWrt router ships with an
	// EMPTY root password (see sshConnect's auth chain).
	if req.IP == "" {
		writeError(w, 400, "IP required")
		return
	}

	client := sshConnect(req.IP, req.Password)
	if client == nil && req.Password != "" {
		client = sshConnect(req.IP, "")
	}
	if client == nil {
		writeError(w, 502, "cannot connect to router via SSH")
		return
	}
	defer client.Close()

	// ---- WiFi scan: ubus/uci-first enumeration (driver-agnostic) ----
	//
	// Priority:
	//   1. Enumerate radios/interfaces via `ubus call network.wireless status`
	//      and `uci show wireless` (parseWirelessStatus). This is
	//      driver-agnostic and works on OpenWrt 25.12 regardless of the
	//      underlying wireless driver (mac80211, mt76, …).
	//   2. Per-device `iwinfo <dev> scan` for each discovered interface,
	//      skipping Master/AP-mode interfaces (they can't scan).
	//   3. Fallback (only when ubus/uci yield nothing): `iw phy phy0/phy1
	//      scan` (phy-level) and `iw dev` listing + per-interface `iw dev
	//      <iface> scan`.
	//
	// iwinfo remains the scan executor; ubus/uci are the enumerator.

	var debugLog strings.Builder

	// --- Enumeration: ubus + uci (both merged, defensively parsed) ---
	ubusOut, ubusErr := sshRunStatus(client, "ubus call network.wireless status 2>&1")
	uciOut, _ := sshRunStatus(client, "uci show wireless 2>&1")
	_ = ubusErr // ubus may legitimately fail on non-OpenWrt; uci is the fallback
	status := parseWirelessStatus(ubusOut, uciOut)
	debugLog.WriteString("[1] ubus/uci enumeration: " +
		fmt.Sprintf("%d radios, %d ifaces\n", len(status.Radios), len(status.Ifaces)))

	// --- W9: radio-down / zero-radio actionable diagnostics ---
	if len(status.Radios) == 0 {
		// Genuinely no radios — surface raw outputs so diagnosis is possible.
		raw := strings.TrimSpace(ubusOut)
		if u := strings.TrimSpace(uciOut); u != "" {
			if raw != "" {
				raw += "\n"
			}
			raw += u
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "no wireless radios found — ubus/uci enumeration empty",
			"debug": scanDebug(debugLog.String(), raw),
		})
		return
	}
	if status.AnyDisabled && len(status.Ifaces) == 0 {
		// Radios exist but are all disabled (or no usable ifaces) — surface
		// the DISCOVERED radio section names so the UI can render a
		// consent-gated enable button (POST /api/wifi-enable) instead of a
		// generic hint. No silent enable happens here.
		radioNames := make([]string, 0, len(status.Radios))
		for _, r := range status.Radios {
			radioNames = append(radioNames, r.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(radiosDisabledResponse(radioNames, scanDebug(debugLog.String(), ubusOut)))
		return
	}

	// --- W8: per-device scan from enumeration (skip AP/Master ifaces) ---
	var scanOut string
	var scanned int
	for _, iface := range status.Ifaces {
		if isAPMode(iface.Mode) {
			debugLog.WriteString("[2] skip iface " + iface.Device + " (mode=" + iface.Mode + ")\n")
			continue
		}
		cmd, err := iwinfoScanCommand(iface.Device)
		if err != nil {
			debugLog.WriteString("[2] skip iface " + iface.Device + ": " + err.Error() + "\n")
			continue
		}
		out, _ := sshRunStatus(client, cmd)
		scanned++
		if strings.TrimSpace(out) != "" && !scanFailedHeuristic(out) {
			scanOut += out + "\n"
			debugLog.WriteString("[2] iwinfo " + iface.Device + " scan: got results\n")
		} else {
			debugLog.WriteString("[2] iwinfo " + iface.Device + " scan: no results (" + truncate(out, 200) + ")\n")
		}
	}
	if strings.TrimSpace(scanOut) != "" {
		ssids := parseIwinfoScan(scanOut)
		w.Header().Set("Content-Type", "application/json")
		if len(ssids) == 0 {
			json.NewEncoder(w).Encode(map[string]any{
				"ssids": ssids,
				"debug": scanDebug(debugLog.String(), scanOut),
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
		}
		return
	}
	if scanned == 0 {
		debugLog.WriteString("[2] no scannable (non-AP) interfaces enumerated\n")
	} else {
		debugLog.WriteString("[2] per-device scan: no results\n")
	}

	// --- Strategy 3 (fallback): phy-level scan (mode-independent) ---
	// Only reached when ubus/uci enumeration yielded no scannable results.
	var phyScanOut string
	for _, phy := range []string{"phy0", "phy1"} {
		out, _ := sshRunStatus(client, "iw phy "+phy+" scan 2>&1")
		if strings.TrimSpace(out) != "" && !scanFailedHeuristic(out) {
			phyScanOut += out + "\n"
		}
	}
	if strings.TrimSpace(phyScanOut) != "" {
		debugLog.WriteString("[3] phy-level scan: got results\n")
		ssids := parseIwScan(phyScanOut)
		w.Header().Set("Content-Type", "application/json")
		if len(ssids) == 0 {
			json.NewEncoder(w).Encode(map[string]any{
				"ssids": ssids,
				"debug": scanDebug(debugLog.String(), phyScanOut),
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
		}
		return
	}
	debugLog.WriteString("[3] phy-level scan: no results\n")

	// --- Strategy 4a (fallback): common interface names wlan0/wlan1 ---
	scanOut, _ = sshRunStatus(client, "iwinfo wlan0 scan 2>&1 || iwinfo wlan1 scan 2>&1")
	if strings.TrimSpace(scanOut) != "" && !scanFailedHeuristic(scanOut) {
		debugLog.WriteString("[4a] iwinfo wlan0/wlan1 fallback: got results\n")
		ssids := parseIwinfoScan(scanOut)
		w.Header().Set("Content-Type", "application/json")
		if len(ssids) == 0 {
			json.NewEncoder(w).Encode(map[string]any{
				"ssids": ssids,
				"debug": scanDebug(debugLog.String(), scanOut),
			})
		} else {
			json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
		}
		return
	}
	debugLog.WriteString("[4a] iwinfo wlan0/wlan1 fallback: no results\n")

	// --- Strategy 4 (fallback): iw dev listing + per-interface iw scan ---
	// Last resort when ubus/uci AND phy-level both yield nothing. `iw dev`
	// returns empty when radios are down or absent, so this is only a
	// best-effort path.
	iwDevsOut, _ := sshRunStatus(client, "iw dev 2>&1")
	iwScanOut := ""
	if scanFailedHeuristic(iwDevsOut) {
		debugLog.WriteString("[4b] iw dev listing failed (" + truncate(iwDevsOut, 200) + ")\n")
	} else {
		for _, iface := range parseIwDevList(iwDevsOut) {
			cmd, err := iwInterfaceScanCommand(iface)
			if err != nil {
				debugLog.WriteString("[4c] skip iface " + iface + ": " + err.Error() + "\n")
				continue
			}
			out, _ := sshRunStatus(client, cmd)
			debugLog.WriteString("[4d] iw dev " + iface + " scan: ")
			if strings.TrimSpace(out) != "" && !scanFailedHeuristic(out) {
				iwScanOut += out + "\n"
				debugLog.WriteString("got results\n")
			} else {
				debugLog.WriteString("no results (" + truncate(out, 200) + ")\n")
			}
		}
		if strings.TrimSpace(iwScanOut) != "" {
			scanOut = iwScanOut
			ssids := parseIwScan(scanOut)
			w.Header().Set("Content-Type", "application/json")
			if len(ssids) == 0 {
				json.NewEncoder(w).Encode(map[string]any{
					"ssids": ssids,
					"debug": scanDebug(debugLog.String(), scanOut),
				})
			} else {
				json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
			}
			return
		}
	}
	debugLog.WriteString("[4e] all fallbacks exhausted\n")

	// All strategies failed — return diagnostic info.
	w.Header().Set("Content-Type", "application/json")
	if strings.TrimSpace(scanOut) != "" {
		json.NewEncoder(w).Encode(map[string]any{
			"ssids": []wifiSSID{},
			"debug": scanDebug(debugLog.String(), scanOut),
		})
		return
	}
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(map[string]any{
		"error": "WiFi scan failed — no wireless interfaces found or iwinfo/iw not available. The router may have been left in a partially-configured state by a previous deployment. Try factory resetting the router.",
		"debug": scanDebug(debugLog.String(), scanOut),
	})
}

// wifiEnableRequest is the JSON body for /api/wifi-enable.
type wifiEnableRequest struct {
	IP       string `json:"ip"`
	Password string `json:"password"`
}

// handleWifiEnable is the consent-gated radio-enable endpoint. It mirrors
// handleWifiScan's SSH connect chain (empty password falls back to the empty
// string via sshConnect's auth chain), then runs the discovered-radio enable
// command + `uci commit wireless` + `wifi up`, polls allRadiosUp, and
// auto-rescans by re-invoking the scan logic. Response: {enabled:[...], ok:true}
// or an error.
func handleWifiEnable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}
	var req wifiEnableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	if req.IP == "" {
		writeError(w, 400, "IP required")
		return
	}

	client := sshConnect(req.IP, req.Password)
	if client == nil && req.Password != "" {
		client = sshConnect(req.IP, "")
	}
	if client == nil {
		writeError(w, 502, "cannot connect to router via SSH")
		return
	}
	defer client.Close()

	// Enumerate DISCOVERED radio + owned iface section names (never hardcoded).
	ubusOut, _ := sshRunStatus(client, "ubus call network.wireless status 2>&1")
	uciOut, _ := sshRunStatus(client, "uci show wireless 2>&1")
	status := parseWirelessStatus(ubusOut, uciOut)

	radioNames := make([]string, 0, len(status.Radios))
	for _, r := range status.Radios {
		radioNames = append(radioNames, r.Name)
	}
	ifaceNames := make([]string, 0, len(status.Ifaces))
	for _, i := range status.Ifaces {
		if i.Section != "" {
			ifaceNames = append(ifaceNames, i.Section)
		}
	}

	cmd := enableRadiosCommand(radioNames, ifaceNames)
	if cmd == "" {
		writeError(w, 400, "no radios discovered to enable")
		return
	}

	sshRun(client, strings.Join([]string{
		cmd,
		`uci commit wireless`,
		`wifi up 2>/dev/null || wifi 2>/dev/null || true`,
	}, " && "))

	// Poll allRadiosUp (~10 × 1.5s budget) before reporting success.
	for i := 0; i < 10; i++ {
		if allRadiosUp(sshRun(client, "ubus call network.wireless status 2>/dev/null")) {
			break
		}
		time.Sleep(1500 * time.Millisecond)
	}

	// The UI auto-rescans on success (re-invoking /api/wifi-scan).
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"enabled": radioNames,
		"ok":      true,
	})
}

type deployRequest struct {
	IP       string `json:"ip"`
	Password string `json:"password"`
	Mode     string `json:"mode"`     // wan | sta
	SSID     string `json:"ssid"`     // for sta mode
	WifiPass string `json:"wifiPass"` // for sta mode
	LNURL    string `json:"lnurl"`    // Lightning address or raw LNURL
	DevSplit int    `json:"devSplit"` // advanced: % to dev fund (0-50, default 10)
	Margin   int    `json:"margin"`   // advanced: operator markup % (0-100, default 0)
	Mint     string `json:"mint"`     // advanced: preferred Cashu mint URL
	// TestMints is OPT-IN: when false/absent (the default for real customer
	// deployments) the wizard configures ONLY the 7 production mints. When
	// true it additionally appends the 2 testnut test mints, which fake
	// Lightning payments for E2E purchase testing — never for production.
	TestMints bool `json:"test_mints"` // advanced: include testnut test mints (E2E only, default false)
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	var req deployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid JSON")
		return
	}
	// Password is optional: a fresh-reset OpenWrt router ships with an
	// EMPTY root password (see sshConnect's auth chain).
	if req.IP == "" {
		writeError(w, 400, "IP required")
		return
	}
	if !validLightningAddress(req.LNURL) {
		writeError(w, 400, "a valid Lightning address is required")
		return
	}

	jobID := fmt.Sprintf("%d", time.Now().UnixNano()%100000000)
	job := newJob(req.IP)

	jobsMutex.Lock()
	jobs[jobID] = job
	jobsMutex.Unlock()

	go runDeployment(job, req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Path[len("/api/status/"):]
	jobsMutex.RLock()
	job, ok := jobs[jobID]
	jobsMutex.RUnlock()
	if !ok {
		writeError(w, 404, "job not found")
		return
	}
	// Return a snapshot for thread-safe JSON
	job.mu.Lock()
	snapshot := struct {
		IP     string     `json:"ip"`
		Status string     `json:"status"`
		Step   int        `json:"step"`
		Steps  []Step     `json:"steps"`
		Log    []LogEntry `json:"log"`
		Error  string     `json:"error,omitempty"`
	}{job.IP, job.Status, job.Step, job.Steps, job.Log, job.Error}
	job.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// allRadiosUp parses `ubus call network.wireless status` output
// ({"radio0":{"up":true,...},"radio1":{...}}) and reports whether EVERY
// radio reports up. Empty/garbage output parses to false (keep polling).
func allRadiosUp(statusJSON string) bool {
	var status map[string]map[string]any
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		return false
	}
	if len(status) == 0 {
		return false
	}
	for _, radio := range status {
		if up, _ := radio["up"].(bool); !up {
			return false
		}
	}
	return true
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/scan", handleScan)
	mux.HandleFunc("/api/wifi-scan", handleWifiScan)
	mux.HandleFunc("/api/wifi-enable", handleWifiEnable)
	mux.HandleFunc("/api/deploy", handleDeploy)
	mux.HandleFunc("/api/status/", handleStatus)
	mux.HandleFunc("/", handleIndex)

	// CORS for local dev
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		mux.ServeHTTP(w, r)
	})

	// Bind BEFORE printing anything (v0.7.0-alpha16): the URL banner is
	// only shown once a socket actually listens, so "running on …" always
	// means a working wizard. PORT overrides the 8099 default; if the
	// preferred port is busy we walk 8099→8109 before giving up.
	port, err := listenPort()
	if err != nil {
		log.Fatal(err)
	}
	fallbacks := make([]string, 0, len(fallbackPorts(port)))
	for _, p := range fallbackPorts(port) {
		fallbacks = append(fallbacks, net.JoinHostPort("", strconv.Itoa(p)))
	}
	ln, addr, err := pickPort(net.JoinHostPort("", strconv.Itoa(port)), fallbacks)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
		fmt.Fprintf(os.Stderr, "The port is busy — likely another wizard instance.\n")
		fmt.Fprintf(os.Stderr, "  Find the holder:   ss -tlnp | grep %d\n", port)
		fmt.Fprintf(os.Stderr, "  Kill the old one:  pkill -f net4sats-wizard\n")
		fmt.Fprintf(os.Stderr, "  Or pick another:   PORT=8199 ./net4sats-wizard\n")
		os.Exit(1)
	}
	url := urlFor(addr)
	fmt.Printf("net4sats wizard running on %s\n", url)
	if addr != fmt.Sprintf(":%d", port) {
		fmt.Printf("(port %d was busy — using %s instead; PORT=<n> picks explicitly)\n", port, url)
	}
	fmt.Println("Open this URL in your browser to set up a router.")
	log.Fatal(http.Serve(ln, handler))
	_ = io.Discard // keep import
}
