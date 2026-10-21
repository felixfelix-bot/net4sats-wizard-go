#!/bin/sh
# fe2b_feed_check.sh — FE2b acceptance check (TDD baseline: EXPECTED TO FAIL until the feed exists)
#
# Written by the FE2b preflight lane (board card t_18d15abe) BEFORE the feed install,
# per QGATES v3.2 (failing test first). This script encodes FE2b's definition of done:
# tollgate-wrt installed AND its origin is the FEED (not a local file), service running,
# :2121 health 2xx, portal HTTP 200.
#
# It is a HOST-SIDE, READ-ONLY probe. It does not modify the router, the repo, or the feed.
#
# Usage: sh scripts/fe2b_feed_check.sh [router_ip] [feed_url]
# Exit: 0 = all assertions pass (FE2b done) | 1 = at least one assertion FAILED
#
# Evidence captured 2026-10-21 on the FE2b preflight run: 5 of 6 assertions FAIL.
# The failing assertions are the honest baseline — they are what the feed install must satisfy.

set -u

ROUTER="${1:-192.168.1.1}"
# The exact feed URL the device's /etc/opkg/customfeeds.conf currently carries (read 2026-10-21).
FEED_URL="${2:-https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages}"
ARCH="aarch64_cortex-a53"
SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 -o BatchMode=no"
LAB_PASS="test123"

PASS=0
FAIL=0
fail() { echo "FAIL  $1"; FAIL=$((FAIL + 1)); }
pass() { echo "PASS  $1"; PASS=$((PASS + 1)); }
info() { echo "      $1"; }

rsh() {
  # read-only remote shell; lab credential per the router-testing lab convention
  if command -v sshpass >/dev/null 2>&1; then
    timeout 25 sshpass -p "$LAB_PASS" ssh $SSH_OPTS root@"$ROUTER" "$1" 2>/dev/null
  else
    timeout 25 ssh $SSH_OPTS root@"$ROUTER" "$1" 2>/dev/null
  fi
}

echo "=== fe2b_feed_check.sh — FE2b acceptance (TDD baseline; failing is the expected starting state) ==="
echo "router   : $ROUTER"
echo "feed_url : $FEED_URL"
echo "arch     : $ARCH"
echo "host     : $(hostname)  date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo

# ---------------------------------------------------------------------------
echo "--- A1: host can reach the router over SSH ---"
if rsh 'true' >/dev/null 2>&1; then
  pass "A1 ssh reachability to $ROUTER"
else
  fail "A1 ssh reachability to $ROUTER (cannot proceed with device assertions)"
fi

# ---------------------------------------------------------------------------
echo "--- A2: the published FEED INDEX is fetchable (the feed must exist to be consumable) ---"
code="$(curl -s -o /dev/null -w '%{http_code}' -L -m 20 "$FEED_URL" 2>/dev/null)"
if [ "$code" = "200" ]; then
  pass "A2 feed index fetchable: HTTP $code at $FEED_URL"
else
  fail "A2 feed index fetchable: HTTP ${code:-000} at $FEED_URL (no published index -> nothing to install from)"
fi

# ---------------------------------------------------------------------------
echo "--- A3: device carries the feed line ---"
feedline="$(rsh 'grep -h "^src/gz" /etc/opkg/customfeeds.conf /etc/apk/repositories 2>/dev/null | grep -i tollgate' | head -1)"
if [ -n "$feedline" ]; then
  pass "A3 device has a tollgate feed line"
  info "$feedline"
else
  fail "A3 device has a tollgate feed line (none found in /etc/opkg/customfeeds.conf or /etc/apk/repositories)"
fi

# ---------------------------------------------------------------------------
echo "--- A4: tollgate-wrt is installed ---"
ver="$(rsh 'opkg list-installed 2>/dev/null | grep -i "^tollgate-wrt" || apk list --installed 2>/dev/null | grep -i tollgate-wrt' | head -1)"
if [ -n "$ver" ]; then
  pass "A4 tollgate-wrt installed"
  info "$ver"
else
  fail "A4 tollgate-wrt installed (not present in opkg/apk installed list)"
fi

# ---------------------------------------------------------------------------
echo "--- A5: tollgate-wrt ORIGIN is the feed, not a local file ---"
# The discriminator: a package installed from a feed records the feed origin; a package
# installed by scp+`opkg install /tmp/<file>.ipk` does not, and a staged local file is
# still sitting in /tmp. BOTH must hold for origin==feed.
localipk="$(rsh 'ls /tmp/*.ipk /tmp/*.apk /root/*.ipk /root/*.apk 2>/dev/null | grep -i tollgate' | head -3)"
orig="$(rsh 'grep -A6 -i "^Package: tollgate-wrt$" /usr/lib/opkg/status 2>/dev/null' | grep -iE '^(Status|Source|Filename|Installed-Time)' | head -4)"
if [ -n "$localipk" ]; then
  fail "A5 origin is the feed: a LOCAL staged package file is present on the device (install was scp-based, not feed-based)"
  info "local artifact(s): $localipk"
else
  pass "A5a no local staged tollgate package file on the device"
fi
if [ -n "$orig" ]; then
  info "opkg status origin fields:"
  printf '%s\n' "$orig" | sed 's/^/      /'
  # heuristics: a feed-installed package has no 'Filename:' pointing at a local path
  if printf '%s' "$orig" | grep -qiE 'Filename: */(tmp|root)/'; then
    fail "A5 origin is the feed: opkg records a local-path Filename"
  fi
else
  fail "A5 origin is the feed: could not read opkg status for tollgate-wrt (origin unproven)"
fi

# ---------------------------------------------------------------------------
echo "--- A6: tollgate-wrt service is running ---"
if rsh 'pgrep -f "[t]ollgate-wrt" >/dev/null 2>&1 && echo RUNNING' | grep -q RUNNING; then
  pass "A6 tollgate-wrt process is running"
else
  fail "A6 tollgate-wrt process is running"
fi

# ---------------------------------------------------------------------------
echo "--- A7: :2121 API health is 2xx ---"
h2121="$(curl -s -o /dev/null -w '%{http_code}' -m 8 "http://$ROUTER:2121/" 2>/dev/null)"
case "$h2121" in
  2*) pass "A7 :2121 health HTTP $h2121" ;;
  *)  fail "A7 :2121 health HTTP ${h2121:-000} (expected 2xx)" ;;
esac

# ---------------------------------------------------------------------------
echo "--- A8: captive portal returns HTTP 200 ---"
h2051="$(curl -s -o /dev/null -w '%{http_code}' -m 8 "http://$ROUTER:2051/" 2>/dev/null)"
if [ "$h2051" = "200" ]; then
  pass "A8 portal HTTP 200"
else
  fail "A8 portal HTTP ${h2051:-000} (expected 200)"
fi

# ---------------------------------------------------------------------------
echo
echo "=== RESULT: $PASS passed, $FAIL failed ==="
if [ "$FAIL" -eq 0 ]; then
  echo "FE2b ACCEPTANCE SATISFIED"
  exit 0
fi
echo "FE2b NOT SATISFIED — this is the expected TDD baseline before the feed install."
echo "Assertions A2 (feed index) and A5 (feed origin) are the ones the feed install must flip."
exit 1
