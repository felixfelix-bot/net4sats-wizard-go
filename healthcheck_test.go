package main

import (
	"strings"
	"testing"
)

// TestParseHealthOutput covers W1: the single-command probe must capture the
// wget exit status, the response body, and stderr as three distinct fields so
// the caller can tell "connection refused" from "timeout" from "empty body".
// The old probe (`wget ... | head -c 100 || echo 'health check failed'`) bound
// the `||` to `head` (always exit 0), so the fallback never fired.
func TestParseHealthOutput(t *testing.T) {
	cases := []struct {
		name       string
		out        string
		wantRC     int
		wantBody   string
		wantStderr string
	}{
		{
			name:       "healthy — rc 0, full advertisement body",
			out:        "HC_RC=0\nHC_BODY_BEGIN\n{\"kind\":10021,\"tags\":[[\"metric\",\"time\"],[\"price_per_step\",\"cashu\",\"1\",\"sats\",\"https://mint\",\"0\"]]}\nHC_BODY_END\nHC_ERR_BEGIN\nHC_ERR_END\n",
			wantRC:     0,
			wantBody:   `{"kind":10021,"tags":[["metric","time"],["price_per_step","cashu","1","sats","https://mint","0"]]}`,
			wantStderr: "",
		},
		{
			name:       "connection refused — rc 4, stderr captured",
			out:        "HC_RC=4\nHC_BODY_BEGIN\nHC_BODY_END\nHC_ERR_BEGIN\nwget: can't connect to remote host (127.0.0.1): Connection refused\nHC_ERR_END\n",
			wantRC:     4,
			wantBody:   "",
			wantStderr: "wget: can't connect to remote host (127.0.0.1): Connection refused",
		},
		{
			name:       "timeout — rc 4, stderr names timeout",
			out:        "HC_RC=4\nHC_BODY_BEGIN\nHC_BODY_END\nHC_ERR_BEGIN\nwget: download timed out\nHC_ERR_END\n",
			wantRC:     4,
			wantBody:   "",
			wantStderr: "wget: download timed out",
		},
		{
			name:       "empty body — rc 0 but no content",
			out:        "HC_RC=0\nHC_BODY_BEGIN\nHC_BODY_END\nHC_ERR_BEGIN\nHC_ERR_END\n",
			wantRC:     0,
			wantBody:   "",
			wantStderr: "",
		},
		{
			name:       "malformed output — rc defaults to -1",
			out:        "garbage with no markers",
			wantRC:     -1,
			wantBody:   "",
			wantStderr: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseHealthOutput(tc.out)
			if got.rc != tc.wantRC {
				t.Errorf("rc = %d; want %d", got.rc, tc.wantRC)
			}
			if got.body != tc.wantBody {
				t.Errorf("body = %q; want %q", got.body, tc.wantBody)
			}
			if got.stderr != tc.wantStderr {
				t.Errorf("stderr = %q; want %q", got.stderr, tc.wantStderr)
			}
		})
	}
}

// TestAdvertisementOK covers W5: the body check must require a full
// kind:10021 advertisement signature — "kind" AND ("price_per_step" OR
// "metric") — so the degraded-mode kind:21023 notice (which also contains
// "kind") is NOT mistaken for a healthy backend.
func TestAdvertisementOK(t *testing.T) {
	cases := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "full advertisement — kind + price_per_step + metric",
			body: `{"kind":10021,"tags":[["metric","time"],["step_size","1"],["price_per_step","cashu","1","sats","https://mint","0"]]}`,
			want: true,
		},
		{
			name: "full advertisement — kind + price_per_step only",
			body: `{"kind":10021,"tags":[["price_per_step","cashu","1","sats","https://mint","0"]]}`,
			want: true,
		},
		{
			name: "full advertisement — kind + metric only",
			body: `{"kind":10021,"tags":[["metric","time"]]}`,
			want: true,
		},
		{
			name: "degraded notice — kind 21023, no price/metric (must NOT pass)",
			body: `{"kind":21023,"pubkey":"abc","tags":[["level","warning"],["code","no-reachable-mints"]],"content":"TollGate is initializing. No reachable mints detected."}`,
			want: false,
		},
		{
			name: "empty body",
			body: "",
			want: false,
		},
		{
			name: "price_per_step but no kind",
			body: `{"tags":[["price_per_step","cashu","1","sats","https://mint","0"]]}`,
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := advertisementOK(tc.body); got != tc.want {
				t.Errorf("advertisementOK(%q) = %v; want %v", tc.body, got, tc.want)
			}
		})
	}
}

// TestPortBound covers W2: the netstat poll must detect a listener on :2121
// (BusyBox netstat -tln output, no -p flag).
func TestPortBound(t *testing.T) {
	cases := []struct {
		name string
		out  string
		want bool
	}{
		{"bound — tcp :2121", "tcp        0      0 0.0.0.0:2121            0.0.0.0:*               LISTEN", true},
		{"bound — tcp6 :2121", "tcp6       0      0 :::2121                 :::*                    LISTEN", true},
		{"not bound — other port", "tcp        0      0 0.0.0.0:8090            0.0.0.0:*               LISTEN", false},
		{"empty output", "", false},
		{"port 21210 must not match :2121", "tcp        0      0 0.0.0.0:21210           0.0.0.0:*               LISTEN", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := portBound(tc.out); got != tc.want {
				t.Errorf("portBound(%q) = %v; want %v", tc.out, got, tc.want)
			}
		})
	}
}

// TestProcessAlive covers the W4 process-liveness check (ps w | grep
// '[t]ollgate-wrt').
func TestProcessAlive(t *testing.T) {
	cases := []struct {
		name string
		out  string
		want bool
	}{
		{"alive", " 1234 root      1234 S    /usr/bin/tollgate-wrt", true},
		{"empty — dead", "", false},
		{"whitespace only — dead", "   \n  ", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := processAlive(tc.out); got != tc.want {
				t.Errorf("processAlive(%q) = %v; want %v", tc.out, got, tc.want)
			}
		})
	}
}

// TestClassifyHealthState covers W4: the four-way classification that decides
// healthy vs soft-fail (still starting) vs hard-fail (dead / http-bad).
func TestClassifyHealthState(t *testing.T) {
	cases := []struct {
		name   string
		alive  bool
		bound  bool
		httpOK bool
		want   healthState
	}{
		{"healthy", true, true, true, healthHealthy},
		{"alive + bound + http bad → httpBad", true, true, false, healthHTTPBad},
		{"alive + not bound → starting (soft-fail)", true, false, false, healthStarting},
		{"dead → dead (hard-fail)", false, false, false, healthDead},
		{"dead + bound + http bad → dead", false, true, false, healthDead},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyHealthState(tc.alive, tc.bound, tc.httpOK); got != tc.want {
				t.Errorf("classifyHealthState(%v,%v,%v) = %v; want %v", tc.alive, tc.bound, tc.httpOK, got, tc.want)
			}
		})
	}
}

// TestClockSyncLog covers W6: the clock-sync read-back log line. Drift ≤5s is
// "synced"; larger drift (or a failed read) is a warning — never a deploy
// failure.
func TestClockSyncLog(t *testing.T) {
	cases := []struct {
		name       string
		laptop     int64
		router     int64
		wantSynced bool
	}{
		{"exact match", 1700000000, 1700000000, true},
		{"2s drift", 1700000000, 1699999998, true},
		{"5s drift (boundary)", 1700000000, 1699999995, true},
		{"6s drift — warn", 1700000000, 1699999994, false},
		{"router ahead 10s — warn", 1700000000, 1700000010, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := clockSyncLog(tc.laptop, tc.router)
			hasWarn := strings.Contains(got, "FAILED")
			if tc.wantSynced && hasWarn {
				t.Errorf("clockSyncLog(%d,%d) = %q; want synced (no FAILED)", tc.laptop, tc.router, got)
			}
			if !tc.wantSynced && !hasWarn {
				t.Errorf("clockSyncLog(%d,%d) = %q; want FAILED warning", tc.laptop, tc.router, got)
			}
		})
	}
}
