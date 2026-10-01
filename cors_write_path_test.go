package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// deployBody is a minimal, syntactically-valid /api/deploy payload. It is
// never actually deployed: the tests below assert the request is refused
// before any handler runs (or, in the allowed cases, run only against a
// stub handler).
const deployBody = `{"ip":"192.168.1.1","lnurl":"operator@wallet.example"}`

// TestWriteGateRejectsCrossOriginWrites is the TDD guard for the wizard #24
// CORS gap: the read-path allowlist only withheld Access-Control-Allow-Origin,
// so a cross-origin "simple request" POST (text/plain, no preflight) still
// reached /api/deploy — the endpoint that drives a root SSH session.
//
// Every request that carries an Origin outside the allowlist must now be
// refused before the wrapped handler runs. Same-origin/loopback callers and
// non-browser clients (no Origin header: curl, the CLI, E2E harnesses) keep
// working.
func TestWriteGateRejectsCrossOriginWrites(t *testing.T) {
	reached := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})
	h := corsMiddleware(inner)

	cases := []struct {
		name       string
		method     string
		path       string
		origin     string
		host       string
		wantStatus int
	}{
		{"cross-origin POST /api/deploy is refused", "POST", "/api/deploy", "http://evil.example", "127.0.0.1:8099", http.StatusForbidden},
		{"cross-origin POST /api/wifi-scan is refused", "POST", "/api/wifi-scan", "http://evil.example", "127.0.0.1:8099", http.StatusForbidden},
		{"cross-origin POST /api/scan is refused", "POST", "/api/scan", "http://evil.example", "127.0.0.1:8099", http.StatusForbidden},
		{"cross-origin GET read is refused", "GET", "/api/scan", "http://evil.example", "127.0.0.1:8099", http.StatusForbidden},
		{"scheme mismatch is refused", "POST", "/api/deploy", "https://localhost:8099", "127.0.0.1:8099", http.StatusForbidden},
		{"suffix trick is refused", "POST", "/api/deploy", "http://localhost:8099.evil.example", "127.0.0.1:8099", http.StatusForbidden},
		{"port mismatch is refused", "POST", "/api/deploy", "http://localhost:8098", "127.0.0.1:8099", http.StatusForbidden},
		{"DNS-rebinding Host with foreign Origin is refused", "POST", "/api/deploy", "http://evil.example", "evil.example:8099", http.StatusForbidden},
		{"null Origin is refused", "POST", "/api/deploy", "null", "127.0.0.1:8099", http.StatusForbidden},

		{"localhost origin POST is allowed", "POST", "/api/deploy", "http://localhost:8099", "127.0.0.1:8099", http.StatusOK},
		{"loopback-ip origin POST is allowed", "POST", "/api/deploy", "http://127.0.0.1:8099", "127.0.0.1:8099", http.StatusOK},
		{"non-browser POST (no Origin) is allowed", "POST", "/api/deploy", "", "127.0.0.1:8099", http.StatusOK},
		{"non-browser GET (no Origin) is allowed", "GET", "/api/scan", "", "127.0.0.1:8099", http.StatusOK},
		{"allowlisted GET read is allowed", "GET", "/api/scan", "http://localhost:8099", "127.0.0.1:8099", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reached = false

			var body io.Reader
			if tc.method == http.MethodPost {
				body = strings.NewReader(deployBody)
			}
			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Host = tc.host
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("Origin %q %s %s: status = %d, want %d",
					tc.origin, tc.method, tc.path, resp.StatusCode, tc.wantStatus)
			}
			if tc.wantStatus == http.StatusForbidden {
				if reached {
					t.Fatalf("Origin %q %s %s: wrapped handler ran for a refused request",
						tc.origin, tc.method, tc.path)
				}
				if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "" {
					t.Errorf("refused request leaked ACAO = %q, want none", got)
				}
			}
		})
	}
}

// TestWriteGateProtectsRealRouter drives the same middleware over the real
// route table so a future refactor of main()'s wiring cannot silently drop the
// gate. The deploy body is deliberately invalid ({}): in the RED state the
// request reaches handleDeploy and is answered 400 "IP required"; once the gate
// is in place the middleware refuses it with 403 before any handler — and, in
// particular, before a deployment goroutine could be spawned.
func TestWriteGateProtectsRealRouter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/scan", handleScan)
	mux.HandleFunc("/api/wifi-scan", handleWifiScan)
	mux.HandleFunc("/api/deploy", handleDeploy)
	mux.HandleFunc("/api/status/", handleStatus)
	mux.HandleFunc("/", handleIndex)
	h := corsMiddleware(mux)

	req := httptest.NewRequest("POST", "/api/deploy", strings.NewReader(`{}`))
	req.Host = "127.0.0.1:8099"
	req.Header.Set("Origin", "http://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("POST /api/deploy from foreign Origin: status = %d, want %d",
			w.Code, http.StatusForbidden)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("refusal Content-Type = %q, want application/json", ct)
	}
}

// TestWizardBindOriginExtendsAllowlist pins the WIZARD_BIND escape hatch:
// when the operator binds the wizard to a concrete non-wildcard address, that
// address's own origin joins the allowlist (so the browser UI served from it
// can still write). A wildcard or loopback bind contributes nothing — the
// static loopback entries already cover the default. The Host header is never
// trusted, so DNS rebinding cannot widen the allowlist.
func TestWizardBindOriginExtendsAllowlist(t *testing.T) {
	cases := []struct {
		bind string
		want string
	}{
		{"127.0.0.1:8099", "http://127.0.0.1:8099"},
		{"localhost:8099", "http://localhost:8099"},
		{"192.168.1.50:8099", "http://192.168.1.50:8099"},
		{"[::1]:8099", "http://[::1]:8099"},
		{"0.0.0.0:8099", ""},
		{"[::]:8099", ""},
		{"", ""},
		{"not-an-address", ""},
	}
	for _, tc := range cases {
		if got := bindOriginFor(tc.bind); got != tc.want {
			t.Errorf("bindOriginFor(%q) = %q, want %q", tc.bind, got, tc.want)
		}
	}
}
