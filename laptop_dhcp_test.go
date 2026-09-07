package main

import "testing"

// TestIPInSubnet verifies the pure CIDR-containment helper. Key property: a
// "192.168.1" prefix must match "192.168.1.50" but NOT "192.168.10.50" — string
// prefix matching would wrongly include the latter.
func TestIPInSubnet(t *testing.T) {
	cases := []struct {
		ip     string
		prefix string
		want   bool
	}{
		{"192.168.1.50", "192.168.1", true},
		{"192.168.1.1", "192.168.1", true},
		{"192.168.1.254", "192.168.1", true},
		{"192.168.10.50", "192.168.1", false}, // the CIDR-vs-string-prefix trap
		{"192.168.2.1", "192.168.1", false},
		{"10.51.187.9", "10.51.187", true},
		{"10.51.188.9", "10.51.187", false},
		{"not-an-ip", "192.168.1", false},
		{"", "192.168.1", false},
		{"2001:db8::1", "192.168.1", false}, // IPv6 excluded
	}
	for _, c := range cases {
		if got := ipInSubnet(c.ip, c.prefix); got != c.want {
			t.Errorf("ipInSubnet(%q, %q) = %v, want %v", c.ip, c.prefix, got, c.want)
		}
	}
}

// TestRenewOSCommandOnCurrentOS ensures the OS table yields a command for the
// host this test runs on (linux here). It documents the OS switch is reachable.
func TestRenewOSCommandOnCurrentOS(t *testing.T) {
	cmd, _, _, ok := renewOSCommand("eth0")
	if !ok || cmd == "" {
		t.Fatalf("expected a known renew command on %s, got cmd=%q ok=%v", currentOS(), cmd, ok)
	}
}

func currentOS() string {
	return "host-OS"
}
