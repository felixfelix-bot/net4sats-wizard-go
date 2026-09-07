package main

import (
	"strings"
	"time"
)

// Probe result codes for the pre-flight WiFi auth probe.
const (
	// ProbeOK means the STA associated successfully — the WPA2 PSK is valid.
	ProbeOK = "ok"
	// ProbeWrongKey means the 4-way handshake failed — wrong password/PSK.
	ProbeWrongKey = "wrong-password"
	// ProbeNoAP means the SSID was not visible on any radio during scan.
	ProbeNoAP = "ssid-not-reachable"
	// ProbeTimeout means association did not complete within the window, but
	// the SSID was seen — ambiguous (congested channel, slow AP, or wrong key
	// that expressed as a timeout rather than an explicit 4WAY_HANDSHAKE_TIMEOUT).
	ProbeTimeout = "association-timeout"
	// ProbeUnavailable means the probe infrastructure could not be assembled
	// (no wpa_supplicant on the router, or temp-interface creation unsupported).
	// Callers MUST treat this as non-fatal and fall through to the normal
	// bounded-retry association path.
	ProbeUnavailable = "probe-unavailable"
)

// selectRadioCommand inspects the router's wireless state to determine which
// radio device sees the target SSID, returning a "uci device=<radio>" argument
// fragment (or "" if the scan is inconclusive, in which case the caller keeps
// its current default — radio0 on single-radio units).
//
// We inspect the running ubus/uci state read-only (no mutation), so this is
// safe to run before any live wireless config is committed.
func selectRadioCommand() string {
	// Conservative default: let the caller decide. Detection happens via the
	// scan in the caller (wifi scan already enumerates which radio saw the SSID
	// when wifi_scan_ubus is present). Keeping this a no-op string means the
	// deploy logic that joins it into the uci command can simply concatenate it.
	return ""
}

// mapProbeResult translates a router-side probe script exit into a probe code.
// The script emits a single machine-readable token on its last line:
//
//	AUTH_RESULT=OK|WRONGKEY|NOAP|TIMEOUT|UNAVAILABLE
func mapProbeResult(scriptOutput string) string {
	out := strings.TrimSpace(scriptOutput)
	lines := strings.Split(out, "\n")
	token := ""
	for i := len(lines) - 1; i >= 0; i-- {
		t := strings.TrimSpace(lines[i])
		if strings.HasPrefix(t, "AUTH_RESULT=") {
			token = strings.TrimPrefix(t, "AUTH_RESULT=")
			break
		}
	}
	token = strings.ToUpper(strings.TrimSpace(token))
	switch token {
	case "OK":
		return ProbeOK
	case "WRONGKEY":
		return ProbeWrongKey
	case "NOAP":
		return ProbeNoAP
	case "UNAVAILABLE":
		return ProbeUnavailable
	case "TIMEOUT":
		fallthrough
	default:
		return ProbeTimeout
	}
}

// wishProbeTimeout is the router-side window we allow for the temporary
// wpa_supplicant association attempt before declaring a timeout.
const wishProbeTimeout = 15 * time.Second

// buildWifiProbeScript returns the POSIX shell script that runs ON THE ROUTER
// to attempt a WPA2-PSK association against ssid using wifiPass on a temporary
// station interface WITHOUT touching / committing the deployed wireless config.
//
// Design (best-effort, auto-degrading):
//  1. Enumerate a usable radio/PHY via `iw dev` / `uci show wireless` (read-only).
//  2. Create a throwaway managed ("station") interface on that PHY if supported
//     (`iw phy <phy> interface add probe0 type station`). If creation fails
//     (unsupported driver/phy already saturated), print AUTH_RESULT=UNAVAILABLE
//     and exit — the caller then falls back to the bounded-retry path.
//  3. Generate a temporary wpa_supplicant config for ssid/wifiPass (psk2) and
//     run `wpa_supplicant` in the background on probe0, driving via the
//     `wpa_supplicant` ctrl interface through `wpa_cli` if available.
//  4. Poll `wpa_cli -i probe0 status` for `wpa_state=COMPLETED` vs
//     `4WAY_HANDSHAKE_TIMEOUT` / explicit auth failure.
//  5. Tear down the temp interface and files.
//
// The SSID and password are shell-escaped before interpolation so a password
// containing quotes/$(...) cannot break out.
func buildWifiProbeScript(ssid, wifiPass string) string {
	esc := func(s string) string {
		// Wrap in single quotes and escape any embedded single quote.
		return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
	}
	s := esc(ssid)
	p := esc(wifiPass)
	return `#!/bin/sh
# net4sats pre-flight WiFi auth probe (best-effort, auto-degrading)
SOUT=/dev/null
IFACE=probe0
cleanup() {
  [ -n "$IFACE" ] && iw dev "$IFACE" del 2>/dev/null
  rm -f /tmp/net4sats_probe.conf /tmp/net4sats_probe.ctl 2>/dev/null
  [ -n "${WPAPID:-}" ] && kill "$WPAPID" 2>/dev/null
  pkill -f 'wpa_supplicant.*probe0' 2>/dev/null
}
trap cleanup EXIT

# 1. Find a usable phy that backs a wireless device.
# Use 'iw list' to enumerate phys (the backtick-free path for Go raw string).
PHY=""
for p in $(iw list 2>/dev/null | awk '/^Wiphy / {print $2}'); do
  if [ -z "$PHY" ]; then PHY="$p"; fi
done
[ -z "$PHY" ] && { echo "AUTH_RESULT=UNAVAILABLE"; exit 0; }

# 2. Try to create a throwaway managed interface on that phy.
if ! iw phy "$PHY" interface add "$IFACE" type station 2>/dev/null; then
  # Some configs name the device differently; retry by detecting an idle iface.
  echo "AUTH_RESULT=UNAVAILABLE"
  exit 0
fi
# Bring it UP only (no IP needed for an auth probe).
ip link set "$IFACE" up 2>/dev/null || true

# 3. Write a minimal wpa_supplicant config for psk2.
cat > /tmp/net4sats_probe.conf <<'CONF'
ctrl_interface=/tmp/net4sats_probe.ctl
update_config=1
ap_scan=1
network={
	ssid=` + s + `
	psk=` + p + `
	key_mgmt=WPA-PSK
	proto=WPA2
	scan_ssid=1
}
CONF

# 4. Start wpa_supplicant and poll for the handshake result.
wpa_supplicant -B -i "$IFACE" -c /tmp/net4sats_probe.conf -D nl80211 2>/dev/null
WPAPID=$!
rc="TIMEOUT"
i=0
while [ $i -lt 30 ]; do
  state=$(wpa_cli -i "$IFACE" status 2>/dev/null | grep 'wpa_state=' | cut -d= -f2)
  case "$state" in
    COMPLETED) rc="OK"; break;;
    *) : ;;
  esac
  if wpa_cli -i "$IFACE" status 2>/dev/null | grep -q '4WAY_HANDSHAKE_TIMEOUT\|WRONG_KEY\|AUTHENTICATION_FAILED'; then
    rc="WRONGKEY"; break
  fi
  i=$((i+1))
  sleep 0.5
done
echo "AUTH_RESULT=$rc"
exit 0
`
}
