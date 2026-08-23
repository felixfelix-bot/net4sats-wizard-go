// net4sats-wizard — cross-platform router onboarding wizard.
// Single binary: serves web UI + API, auto-discovers routers,
// deploys net4sats over SSH.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
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

const listenAddr = ":8099"

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

	// A fresh-reset OpenWrt router ships with radios DISABLED in UCI —
	// `iwinfo scan` would return nothing and the UI would show an empty
	// SSID list. Enable all radios, bring wifi up, and wait until every
	// radio reports up before scanning (see enableWifiAndWait).
	enableWifiAndWait(client)

	// Auto-detect wireless interfaces from `iwinfo` (no args) output.
	// OpenWrt interfaces can be named phy0-ap0, wlan0, wl0-sha0, etc —
	// not predictable, so we parse the list first.
	iwinfoOut := sshRun(client, "iwinfo 2>/dev/null")
	var wifiDevs []string
	for _, line := range strings.Split(iwinfoOut, "\n") {
		trimmed := strings.TrimSpace(line)
		// iwinfo lines look like: "phy0-ap0 ESSID: \"TollGate-F794\""
		// or "wlan0     ESSID: \"mywifi\""
		// The first token before "ESSID:" is the interface name
		if idx := strings.Index(trimmed, "ESSID:"); idx > 0 {
			iface := strings.TrimSpace(trimmed[:idx])
			if iface != "" && !strings.HasPrefix(iface, "Usage") {
				wifiDevs = append(wifiDevs, iface)
			}
		}
	}

	var scanOut string
	if len(wifiDevs) > 0 {
		// Scan each detected interface
		for _, dev := range wifiDevs {
			out := sshRun(client, "iwinfo "+dev+" scan 2>/dev/null")
			if strings.TrimSpace(out) != "" && !strings.Contains(out, "command not found") {
				scanOut += out + "\n"
			}
		}
	}

	// Fallback: try common interface names if auto-detect found nothing.
	// GL.iNet MT3000 and other modern OpenWrt devices use phy0-ap0 / phy1-ap0
	// naming — wlan0/wlan1 often don't exist. But the error text from these
	// non-existent interfaces must be detected so iw dev scan is tried next.
	if strings.TrimSpace(scanOut) == "" {
		scanOut = sshRun(client, "iwinfo wlan0 scan 2>/dev/null || iwinfo wlan1 scan 2>/dev/null")
	}

	// Detect scan failures by error signatures in the output.
	// "No such wireless device" is iwinfo's response for non-existent
	// interfaces (NOT matched by "No such device" — different string!).
	// Without this check, the wlan0/wlan1 fallback error text gets passed
	// to parseIwinfoScan which returns 0 SSIDs → {"ssids":null}.
	scanFailed := strings.TrimSpace(scanOut) == "" ||
		strings.Contains(scanOut, "command not found") ||
		strings.Contains(scanOut, "No such device") ||
		strings.Contains(scanOut, "No such wireless device") ||
		strings.Contains(scanOut, "Operation not supported") ||
		strings.Contains(scanOut, "Operation not permitted") ||
		strings.Contains(scanOut, "Device or resource busy")

	if scanFailed {
		scanOut = sshRun(client, "iw dev scan 2>/dev/null")
		if strings.TrimSpace(scanOut) != "" && !strings.Contains(scanOut, "command not found") {
			ssids := parseIwScan(scanOut)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
			return
		}
		writeError(w, 500, "WiFi scan failed — no wireless interfaces found or iwinfo/iw not available. The router may have been left in a partially-configured state by a previous deployment. Try factory resetting the router.")
		return
	}

	ssids := parseIwinfoScan(scanOut)
	w.Header().Set("Content-Type", "application/json")
	// If 0 SSIDs found despite non-empty scan output, include diagnostic
	// info so the operator can see what the router actually returned.
	// This catches format mismatches (new iwinfo output) and error text
	// that slipped past the checks above.
	if len(ssids) == 0 && strings.TrimSpace(scanOut) != "" {
		json.NewEncoder(w).Encode(map[string]any{
			"ssids": ssids,
			"debug": truncate(scanOut, 200),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"ssids": ssids})
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

// enableWifiAndWait enables all UCI wifi-devices, starts wifi, and polls
// `ubus call network.wireless status` until every radio reports up
// (~15s budget; fixed sleeps race slow driver init, polling removes that).
// Best-effort: the scan proceeds regardless after the timeout — errors
// surface in the scan step itself where they are actionable.
func enableWifiAndWait(client *ssh.Client) {
	sshRun(client, strings.Join([]string{
		`for r in $(uci -q show wireless 2>/dev/null | sed -n "s/^wireless\.\([^.]*\)=wifi-device$/\1/p"); do uci -q set wireless.$r.disabled='0'; done`,
		`uci commit wireless`,
		`wifi up 2>/dev/null || wifi 2>/dev/null || true`,
	}, " && "))
	for i := 0; i < 10; i++ {
		if allRadiosUp(sshRun(client, "ubus call network.wireless status 2>/dev/null")) {
			return
		}
		time.Sleep(1500 * time.Millisecond)
	}
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

	fmt.Printf("net4sats wizard running on http://localhost%s\n", listenAddr)
	fmt.Println("Open this URL in your browser to set up a router.")
	log.Fatal(http.ListenAndServe(listenAddr, handler))
	_ = io.Discard // keep import
}
