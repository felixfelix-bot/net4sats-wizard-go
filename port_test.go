package main

import (
	"net"
	"strings"
	"testing"
)

// TestListenPort covers the PORT environment override: unset/whitespace
// yields the 8099 default, valid numbers pass through, and anything else
// (non-numeric, zero, out of range) is rejected with an error instead of
// silently falling back to the default (a typo'd PORT must be loud).
func TestListenPort(t *testing.T) {
	cases := []struct {
		name    string
		portEnv string
		want    int
		wantErr bool
	}{
		{"unset — default 8099", "", 8099, false},
		{"whitespace — default 8099", "   ", 8099, false},
		{"explicit 8099", "8099", 8099, false},
		{"override 8199", "8199", 8199, false},
		{"override 1", "1", 1, false},
		{"override 65535", "65535", 65535, false},
		{"non-numeric", "eighty", 0, true},
		{"zero (random) rejected", "0", 0, true},
		{"negative", "-1", 0, true},
		{"too large", "65536", 0, true},
		{"float", "80.99", 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PORT", tc.portEnv)
			got, err := listenPort()
			if tc.wantErr {
				if err == nil {
					t.Errorf("listenPort(PORT=%q) = %d, nil; want error", tc.portEnv, got)
				}
				return
			}
			if err != nil {
				t.Errorf("listenPort(PORT=%q) unexpected error: %v", tc.portEnv, err)
				return
			}
			if got != tc.want {
				t.Errorf("listenPort(PORT=%q) = %d; want %d", tc.portEnv, got, tc.want)
			}
		})
	}
}

// TestFallbackPorts covers the EADDRINUSE fallback chain: with the 8099
// default the wizard tries 8100…8109 next (the documented 8099→8109
// order), and the chain never proposes ports above 65535.
func TestFallbackPorts(t *testing.T) {
	cases := []struct {
		name      string
		preferred int
		want      []int
	}{
		{"default 8099 → 8100..8109", 8099, []int{8100, 8101, 8102, 8103, 8104, 8105, 8106, 8107, 8108, 8109}},
		{"custom 9090 → next ten", 9090, []int{9091, 9092, 9093, 9094, 9095, 9096, 9097, 9098, 9099, 9100}},
		{"near limit capped at 65535", 65531, []int{65532, 65533, 65534, 65535}},
		{"at limit — no fallbacks", 65535, []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fallbackPorts(tc.preferred)
			if len(got) != len(tc.want) {
				t.Fatalf("fallbackPorts(%d) = %v; want %v", tc.preferred, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("fallbackPorts(%d)[%d] = %d; want %d", tc.preferred, i, got[i], tc.want[i])
				}
			}
		})
	}
}

// occupy listens on an ephemeral loopback port and returns its concrete
// address plus a closer. Used by TestPickPort to hold ports deterministically.
func occupy(t *testing.T) (addr string, close func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("occupy: listen: %v", err)
	}
	return ln.Addr().String(), func() { ln.Close() }
}

// TestPickPort exercises the bind-time port selection (the fix for the
// "printed success banner, then died on bind: address already in use"
// startup bug): preferred port first, then fallbacks in order, and a
// helpful error naming the preferred port when everything is busy.
func TestPickPort(t *testing.T) {
	t.Run("preferred free — binds preferred, no fallback needed", func(t *testing.T) {
		busy, closeBusy := occupy(t)
		defer closeBusy()
		ln, addr, err := pickPort("127.0.0.1:0", []string{busy})
		if err != nil {
			t.Fatalf("pickPort: unexpected error: %v", err)
		}
		defer ln.Close()
		if addr != "127.0.0.1:0" {
			t.Errorf("pickPort addr = %q; want %q", addr, "127.0.0.1:0")
		}
	})

	t.Run("preferred busy — binds first free fallback", func(t *testing.T) {
		busy, closeBusy := occupy(t)
		defer closeBusy()
		free, closeFree := occupy(t) // reserve a guaranteed-free port…
		closeFree()                  // …then release it for pickPort
		ln, addr, err := pickPort(busy, []string{free})
		if err != nil {
			t.Fatalf("pickPort: unexpected error: %v", err)
		}
		defer ln.Close()
		if addr != free {
			t.Errorf("pickPort addr = %q; want fallback %q", addr, free)
		}
	})

	t.Run("all candidates busy — error names preferred", func(t *testing.T) {
		b1, closeB1 := occupy(t)
		defer closeB1()
		b2, closeB2 := occupy(t)
		defer closeB2()
		b3, closeB3 := occupy(t)
		defer closeB3()
		ln, addr, err := pickPort(b1, []string{b2, b3})
		if err == nil {
			ln.Close()
			t.Fatal("pickPort with all ports busy should return an error")
		}
		if ln != nil || addr != "" {
			t.Errorf("pickPort ln/addr = %v/%q; want nil/\"\"", ln, addr)
		}
		if want := "127.0.0.1:"; !strings.Contains(err.Error(), want) {
			t.Errorf("pickPort error %q should name the preferred address (contains %q)", err, want)
		}
	})
}