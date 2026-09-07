package main

import (
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// laptop_dhcp.go — after the wizard moves the router's LAN subnet (to resolve an
// upstream-WiFi conflict), the operator's laptop keeps a stale DHCP lease on the
// old subnet and can no longer reach the router. This module renews the laptop's
// DHCP lease so it lands on the new subnet — best-effort, OS-specific, with a
// graceful no-op if the laptop already renewed or has no interface on the old
// subnet. If elevation isn't available (macOS/Windows/headless Linux), the caller
// surfaces a copy-paste one-liner instead of failing the deploy.
//
// This is pure laptop-side logic. It never touches the router or blocks the
// deployment; it runs after the router LAN move is already committed.

// renewLaptopDHCP renews the laptop's DHCP lease so it re-learns the router's new
// LAN subnet. firstThreeOctets is the OLD LAN prefix (e.g. "192.168.1"); the
// interface that currently holds an address on that prefix is the stale one.
// newRouterIP is the router's NEW LAN IP (e.g. "10.51.187.1"), used to verify the
// laptop landed on the new subnet.
//
// Returns a human notice for the UI ("" = nothing to report). Never returns an
// error that should fail the deploy — worst case it logs a copy-paste one-liner.
func renewLaptopDHCP(job *Job, oldPrefix, newRouterIP string) string {
	iface, cur, err := findIfaceOnSubnet(oldPrefix)
	if err != nil || iface == "" {
		// Nothing stale to renew (laptop may already have re-leased) — no-op.
		job.addLog("Laptop DHCP: no interface on old subnet " + oldPrefix + ".0/24 — nothing to renew.")
		return ""
	}
	job.addLog("Laptop DHCP: renewing stale lease " + cur + " on " + iface + " so it can reach router at " + newRouterIP)

	cmd, args, requiresRoot, ok := renewOSCommand(iface)
	if !ok || cmd == "" {
		// Unknown OS or no renew command available — fall back to manual one-liner.
		notice := "Router LAN moved to " + newRouterIP + ". Renew your laptop's DHCP to reconnect, or use: " + manualRenewHint(iface)
		job.addLog(notice)
		return notice
	}

	if runRenew(cmd, args) {
		// Poll for the laptop to appear on the new subnet.
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			if onSubnet(newRouterIP) {
				job.addLog("Laptop renewed to " + newRouterIP + " subnet automatically.")
				return "Laptop automatically renewed to new subnet " + newRouterIP + "."
			}
			time.Sleep(1 * time.Second)
		}
		job.addLog("Laptop DHCP renewal ran but is not yet on " + newRouterIP + " subnet.")
		return "Ran DHCP renewal — if your IP still looks stale, run: " + manualRenewHint(iface)
	}

	// Command failed (likely no permission) — fall back to a manual one-liner.
	_ = requiresRoot
	notice := "Laptop DHCP renewal needs permission. Manually renew, or use: " + manualRenewHint(iface)
	job.addLog(notice)
	return notice
}

// ipInSubnet reports whether the IPv4 ipStr is contained in the /24 described by
// firstThreeOctets (e.g. "192.168.1"). Pure + testable. CIDR-matched so
// "192.168.1" never matches "192.168.10.x".
func ipInSubnet(ipStr, firstThreeOctets string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return false
	}
	_, ipNet, err := net.ParseCIDR(firstThreeOctets + ".0/24")
	if err != nil {
		return false
	}
	return ipNet.Contains(ip)
}

// findIfaceOnSubnet returns the laptop interface whose IPv4 address is inside the
// /24 described by firstThreeOctets (e.g. "192.168.1"), plus that address.
// CIDR-matched so "192.168.1" does not falsely match "192.168.10.x". Returns
// ("", "", nil) when no interface is on that subnet (treat as success/no-op).
func findIfaceOnSubnet(firstThreeOctets string) (name, addr string, err error) {
	cidr := firstThreeOctets + ".0/24"
	_, ipNet, perr := net.ParseCIDR(cidr)
	if perr != nil {
		return "", "", perr
	}
	ifaces, ierr := net.Interfaces()
	if ierr != nil {
		return "", "", ierr
	}
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 {
			continue // skip down interfaces
		}
		addrs, aerr := ifc.Addrs()
		if aerr != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.To4() == nil {
				continue // IPv6 — not a DHCP IPv4 lease
			}
			if ipNet.Contains(ip) {
				return ifc.Name, ip.String(), nil
			}
		}
	}
	return "", "", nil
}

// onSubnet reports whether the laptop currently has an IPv4 address inside the
// /24 of newRouterIP (e.g. "10.51.187.1"). This is the post-renew verification.
func onSubnet(newRouterIP string) bool {
	ip := net.ParseIP(newRouterIP)
	if ip == nil || ip.To4() == nil {
		return false
	}
	// Take the first three octets of newRouterIP (e.g. "10.51.187").
	sp := strings.Split(ip.To4().String(), ".")
	if len(sp) < 3 {
		return false
	}
	three := sp[0] + "." + sp[1] + "." + sp[2]
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, ifc := range ifaces {
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var aip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				aip = v.IP
			case *net.IPAddr:
				aip = v.IP
			}
			if aip == nil || aip.To4() == nil {
				continue
			}
			if ipInSubnet(aip.String(), three) {
				return true
			}
		}
	}
	return false
}

// renewOSCommand returns the OS-appropriate command to renew a DHCP lease on
// iface, plus whether it needs root. ok=false means no known command for this OS.
func renewOSCommand(iface string) (cmd string, args []string, requiresRoot bool, ok bool) {
	switch runtime.GOOS {
	case "linux":
		if hasCommand("nmcli") {
			// NetworkManager path — two steps: disconnect then reconnect forces a
			// fresh DHCPDISCOVER. Often allowed via polkit unprivileged.
			return "nmcli", []string{"device", "disconnect", iface}, false, true
		}
		if hasCommand("dhclient") {
			return "dhclient", []string{"-r", "-v", iface}, true, true
		}
		return "", nil, false, false
	case "darwin":
		// ipconfig set <iface> DHCP takes the BSD interface name we already have
		// (en0), avoiding networksetup's service-name mapping complexity.
		return "ipconfig", []string{"set", iface, "DHCP"}, true, true
	case "windows":
		// Named-adapter release errors on localized names; fall back is handled
		// by renewOne (bare ipconfig release|renew releases all). Start with the
		// named form; caller retries bare if it fails.
		return "ipconfig", []string{"/release", iface}, true, true
	default:
		return "", nil, false, false
	}
}

// runRenew executes the command (and, for Windows, the matching /renew), retrying
// the bare no-adapter form if the named form fails. Returns true if it believed
// to have run (best-effort; verification is the caller's subnet poll).
func runRenew(cmd string, args []string) bool {
	if err := exec.Command(cmd, args...).Run(); err != nil {
		// For Windows, retry with bare ipconfig release|renew (releases all).
		if runtime.GOOS == "windows" {
			if strings.Contains(strings.Join(args, " "), "/release") {
				_ = exec.Command("ipconfig", "/release").Run()
				_ = exec.Command("ipconfig", "/renew").Run()
				return true
			}
		}
		return false
	}
	// Complete the cycle: disconnect+reconnect (nmcli) needs the reconnect step.
	if cmd == "nmcli" && len(args) >= 1 && args[0] == "device" && strings.Contains(strings.Join(args, " "), "disconnect") {
		_ = exec.Command("nmcli", "device", "connect", args[len(args)-1]).Run()
	}
	if cmd == "dhclient" {
		_ = exec.Command("dhclient", "-v", args[len(args)-1]).Run()
	}
	if cmd == "ipconfig" && runtime.GOOS == "windows" {
		_ = exec.Command("ipconfig", "/renew").Run()
	}
	return true
}

// manualRenewHint returns the copy-paste one-liner that forces a DHCP renewal on
// iface, OS-appropriate. Used as the fallback when auto-renew lacks permission.
func manualRenewHint(iface string) string {
	switch runtime.GOOS {
	case "linux":
		if hasCommand("nmcli") {
			return "sudo nmcli device disconnect " + iface + " && sudo nmcli device connect " + iface
		}
		return "sudo dhclient -r " + iface + " && sudo dhclient " + iface
	case "darwin":
		return "sudo ipconfig set " + iface + " DHCP"
	case "windows":
		return "ipconfig /release && ipconfig /renew"
	default:
		return "renew your laptop's DHCP lease"
	}
}

// hasCommand reports whether a command is on PATH.
func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
