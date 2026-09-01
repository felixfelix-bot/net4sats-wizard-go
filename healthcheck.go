package main

import (
	"fmt"
	"strconv"
	"strings"
)

// healthProbeResult is the parsed output of the single-command HTTP probe
// (W1). The probe script emits three delimited sections so the caller can
// distinguish "connection refused" from "timeout" from "empty body" — the old
// `wget ... | head -c 100 || echo 'health check failed'` bound the `||` to
// `head` (always exit 0), so the fallback never fired.
type healthProbeResult struct {
	rc     int
	body   string
	stderr string
}

// healthProbeScript returns the single-command BusyBox-sh probe that captures
// wget's exit status, the response body, and stderr as three delimited
// sections. `wget -qO-` prints the body to stdout; stderr is redirected to a
// temp file so it stays separate. The `; echo HC_RC=$?` runs unconditionally
// (not `||`), so the exit code is always captured.
func healthProbeScript() string {
	return "out=$(wget -qO- --timeout=5 http://127.0.0.1:2121/ 2>/tmp/hc_err); rc=$?; " +
		"echo HC_RC=$rc; echo HC_BODY_BEGIN; echo \"$out\"; echo HC_BODY_END; " +
		"echo HC_ERR_BEGIN; cat /tmp/hc_err 2>/dev/null; echo HC_ERR_END"
}

// parseHealthOutput parses the probe script's delimited output. A missing
// HC_RC marker yields rc == -1 (unparseable).
func parseHealthOutput(out string) healthProbeResult {
	res := healthProbeResult{rc: -1}
	if i := strings.Index(out, "HC_RC="); i >= 0 {
		rest := out[i+len("HC_RC="):]
		if j := strings.IndexAny(rest, "\n"); j >= 0 {
			if n, err := strconv.Atoi(strings.TrimSpace(rest[:j])); err == nil {
				res.rc = n
			}
		}
	}
	res.body = section(out, "HC_BODY_BEGIN", "HC_BODY_END")
	res.stderr = section(out, "HC_ERR_BEGIN", "HC_ERR_END")
	return res
}

// section extracts the text between beginMarker and endMarker in out.
func section(out, beginMarker, endMarker string) string {
	i := strings.Index(out, beginMarker)
	if i < 0 {
		return ""
	}
	i += len(beginMarker)
	// Skip the newline immediately after the begin marker.
	if i < len(out) && out[i] == '\n' {
		i++
	}
	j := strings.Index(out[i:], endMarker)
	if j < 0 {
		return strings.TrimSpace(out[i:])
	}
	return strings.TrimSpace(out[i : i+j])
}

// portBound reports whether netstat -tln output shows a listener on :2121
// (W2). BusyBox netstat lacks -p, so only -tln is used.
func portBound(netstatOut string) bool {
	for _, line := range strings.Split(netstatOut, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Match ":2121" as a local-address token. ":21210" does NOT end in
		// ":2121", so a bare HasSuffix check is sufficient to avoid it.
		for _, field := range strings.Fields(line) {
			if strings.HasSuffix(field, ":2121") {
				return true
			}
		}
	}
	return false
}

// advertisementOK reports whether body is a full kind:10021 advertisement
// (W5). The degraded-mode notice is kind:21023 and also contains "kind", so a
// bare "kind" substring match is insufficient — require "kind" AND
// ("price_per_step" OR "metric"), matching the real GetAdvertisement() JSON
// shape (kind 10021 with price_per_step tags per reachable mint, plus a
// "metric" tag).
func advertisementOK(body string) bool {
	if body == "" {
		return false
	}
	if !strings.Contains(body, "kind") {
		return false
	}
	return strings.Contains(body, "price_per_step") || strings.Contains(body, "metric")
}

// healthState is the four-way classification of the backend's health at the
// end of the step-10 probe window (W4).
type healthState int

const (
	// healthHealthy — port bound and HTTP body is a full kind:10021
	// advertisement.
	healthHealthy healthState = iota
	// healthStarting — process alive but :2121 not yet bound (cold-boot mint
	// probe still running). Soft-fail: warn, do NOT roll back wireless.
	healthStarting
	// healthDead — process gone (or respawn exhausted). Hard-fail.
	healthDead
	// healthHTTPBad — port bound but HTTP body is not a valid advertisement.
	// Hard-fail after diagnostics.
	healthHTTPBad
)

// processAlive reports whether `ps w | grep '[t]ollgate-wrt'` found a live
// process (W4).
func processAlive(psOut string) bool {
	return strings.TrimSpace(psOut) != ""
}

// classifyHealthState maps the three observed signals to a healthState (W4).
func classifyHealthState(alive, bound, httpOK bool) healthState {
	if !alive {
		return healthDead
	}
	if !bound {
		return healthStarting
	}
	if !httpOK {
		return healthHTTPBad
	}
	return healthHealthy
}

// clockSyncLog builds the W6 read-back log line. Drift ≤5s is "synced";
// larger drift is a warning (never a deploy failure).
func clockSyncLog(laptopUnix, routerUnix int64) string {
	drift := routerUnix - laptopUnix
	if drift < 0 {
		drift = -drift
	}
	if drift <= 5 {
		return fmt.Sprintf("clock synced (drift %ds)", drift)
	}
	return fmt.Sprintf("clock sync FAILED (drift %ds)", drift)
}
