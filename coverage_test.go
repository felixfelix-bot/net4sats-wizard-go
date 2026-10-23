package main

import (
	"testing"
	"time"
)

// sshProbeUnroutableIP is an RFC 5737 TEST-NET-1 documentation address. It is
// not routable on any real network, so a TCP dial to it always fails and
// sshIdentify() short-circuits to its "no connection" return path.
//
// The previous version of TestSshRunFunction dialled 192.168.1.1 — the default
// LAN gateway of an OpenWrt router, i.e. exactly the device this wizard
// targets. On a bench host that actually has such a router wired up (port 22
// open, empty root password or a key in ~/.ssh), sshConnect() really
// authenticated and sshIdentify() returned the router's real firmware/model,
// so the test asserted a property of the host's network instead of a property
// of the code and failed intermittently.
const sshProbeUnroutableIP = "192.0.2.1"

func TestSshRunFunction(t *testing.T) {
	// This tests sshRun indirectly by testing the functions that call it
	// We can't easily test sshRun directly without actual SSH connections

	// Test that sshRun is called indirectly through sshIdentify
	// This is a simple smoke test to ensure the function exists and doesn't crash
	// when called with empty credentials

	// Belt and braces: TEST-NET-1 is expected to be unroutable, but if this
	// host can somehow reach it (a fake-IP/proxy DNS setup, a netns test rig,
	// a VPN that blackhole-routes 0.0.0.0/0 into a live lab) then any host is
	// reachable and the assertion below would again be checking the network
	// rather than the code. Skip instead of reporting a bogus failure.
	if tcpProbe(sshProbeUnroutableIP, 22, 500*time.Millisecond) {
		t.Skipf("SSH port 22 on %s is reachable from this host; skipping smoke test that assumes an unreachable target", sshProbeUnroutableIP)
	}

	firmware, vendor, model := sshIdentify(sshProbeUnroutableIP, "")
	// Result will be empty strings, but the function should not panic
	if len(firmware) != 0 || len(vendor) != 0 || len(model) != 0 {
		t.Errorf("sshIdentify against unreachable %s with empty credentials should return empty strings, got firmware=%q, vendor=%q, model=%q", sshProbeUnroutableIP, firmware, vendor, model)
	}
}

func TestDiscoverRoutersNoNetwork(t *testing.T) {
	// Test discoverRouters with no network connectivity
	// This will test the basic structure without actual network calls

	routers := discoverRouters()
	// Should return empty array or routers with minimal info when no network
	// The important thing is that the function doesn't crash
	if routers == nil {
		t.Error("discoverRouters should not return nil")
	}

	// Verify each router has the expected structure
	for _, router := range routers {
		if router.IP == "" {
			t.Error("Router should have non-empty IP")
		}
		if router.Vendor == "" {
			t.Error("Router should have vendor info")
		}
	}
}

func TestEmptyJobOperations(t *testing.T) {
	// Test edge cases with empty/invalid job data

	job := newJob("")
	if job.IP != "" {
		t.Errorf("Job with empty IP should have IP = '', got %q", job.IP)
	}

	job.setStep(0, "done", "")
	if job.Step != 0 {
		t.Errorf("Job step should be 0, got %d", job.Step)
	}

	// Test setting step out of bounds
	job.setStep(999, "done", "")
	// Should not crash, should not change step out of bounds
	if job.Step > 0 {
		t.Errorf("Setting step out of bounds should not change current step, got %d", job.Step)
	}
}
