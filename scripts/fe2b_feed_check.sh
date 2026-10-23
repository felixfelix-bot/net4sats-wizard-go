#!/bin/sh
# fe2b_feed_check.sh — FE2b acceptance gate (TDD RED baseline; EXPECTED TO FAIL today)
#
# Card: t_877ad4ef "Write and run failing feed check script, capture baseline"
#       (net4sats-mvp-v2 board). Codifies the FE2b definition of done from t_53879c67:
#
#         on GL-MT6000: add feeds line -> update -> install tollgate-wrt FROM THE FEED
#         (not a local scp'd file) -> service up -> portal serves.
#
# This is a HOST-SIDE, READ-ONLY probe. It mutates neither the router, the repo,
# nor the feed. It runs the SAME shell probes on the device as a human would, and
# only reads.
#
# WHY THE TRANSPORT IS SERIAL BY DEFAULT
#   The GL-MT6000 on this bench has NO IP path from the host: TCP 22/2121/2051 are
#   all closed (verified). Its only reachable console is USB UART (/dev/ttyUSB0),
#   identified by the device's own prompt (root@tollgate-0195), never by port name.
#   A ssh-only gate would therefore be unrunnable, and an unrunnable gate that
#   "can't reach the router" is not a RED baseline — it is a broken test.
#   FE2B_TRANSPORT=ssh still works if the box is ever put back on the LAN.
#
# Usage:
#   sh scripts/fe2b_feed_check.sh                 # auto: probe TCP, else UART
#   FE2B_TRANSPORT=serial FE2B_SERIAL=/dev/ttyUSB0 sh scripts/fe2b_feed_check.sh
#   FE2B_TRANSPORT=ssh sh scripts/fe2b_feed_check.sh 192.168.1.1
#
# Exit: 0 = every assertion passed (FE2b satisfied) | 1 = at least one FAILED
#       2 = could not establish a transport at all (no evidence either way)
#
# The five assertions the card names, plus transport + feed-index preconditions:
#   A1 tollgate-wrt is installed
#   A2 its ORIGIN is the feed, not a local file path
#   A3 the service is running
#   A4 :2121 health returns 2xx
#   A5 the portal returns HTTP 200
#   A6 the src/gz feed index is actually fetchable (precondition of A2)
#   A0 a transport to the device works (precondition of A1..A3)

set -u

ROUTER="${1:-192.168.1.1}"
FEED_URL="${2:-https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages}"
# The exact src/gz line physically present on the device (/etc/opkg/customfeeds.conf:6).
FEED_LINE="src/gz tollgate $FEED_URL"
# Device facts from the package-state capture (t_18d15abe §5, re-verified 2026-10-23).
PKG="tollgate-wrt"
SERVICE="$PKG"
INITD="/etc/init.d/$PKG"
BIN="/usr/bin/$PKG"
HEALTH_PORT="${FE2B_HEALTH_PORT:-2121}"
# The portal's served document. Root "/" is 403 BY DESIGN on this box (the document
# root has no index.html and directory listing is disabled), so asserting 200 on "/"
# would be a permanently-red check that no install could ever turn green. /splash.html
# is the captive-portal landing page a pre-auth client is actually served.
PORTAL_PORT="${FE2B_PORTAL_PORT:-2051}"
PORTAL_PATH="${FE2B_PORTAL_PATH:-/splash.html}"

TRANSPORT="${FE2B_TRANSPORT:-auto}"
SERIAL_DEV="${FE2B_SERIAL:-/dev/ttyUSB0}"
SERIAL_BAUD="${FE2B_BAUD:-115200}"
LAB_PASS="${FE2B_LAB_PASS:-test123}"
SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=8 -o BatchMode=no"

PASS=0
FAIL=0
pass() { echo "PASS  $1"; PASS=$((PASS + 1)); }
fail() { echo "FAIL  $1"; FAIL=$((FAIL + 1)); }
info() { echo "      $1"; }

die_no_transport() {
  echo
  echo "=== RESULT: $PASS passed, $FAIL failed (ABORTED) ==="
  echo "FE2b INCONCLUSIVE — no transport to $ROUTER (ssh) or $SERIAL_DEV (serial)."
  echo "This is NOT a RED baseline: nothing was measured. Fix access and re-run."
  exit 2
}

# ---------------------------------------------------------------------------
# Transport layer: rsh "<remote command>" -> remote stdout on our stdout
# ---------------------------------------------------------------------------
port_open() { # host port -> 0 if a TCP connect succeeds
  python3 - "$1" "$2" <<'PY' 2>/dev/null
import socket, sys
try:
    s = socket.create_connection((sys.argv[1], int(sys.argv[2])), 3)
    s.close()
except Exception:
    sys.exit(1)
PY
}

SERIAL_HELPER=""
write_serial_helper() {
  [ -n "$SERIAL_HELPER" ] && return 0
  SERIAL_HELPER="$(mktemp -t fe2b-serial-XXXXXX.py)"
  cat > "$SERIAL_HELPER" <<'PY'
"""Read-only serial transport: run one shell command, return its stdout.

Framing: the remote command is wrapped in unique printf markers. The markers are
written with a QUOTE SPLIT (__TGST""ART__) so the shell's echo of the command line
differs textually from the marker actually emitted, and the echo can never be
mistaken for the real start/end of output. Line-wrap at 80 cols is handled because
the markers are emitted by printf as whole lines.
"""
import sys, time
from serial import Serial

port, baud, cmd = sys.argv[1], int(sys.argv[2]), sys.argv[3]
START, END = "__TGSTART__", "__TGEND__"
# the echoed command shows these literals; the emitted output shows the real markers
START_Q, END_Q = '__TGST""ART__', '__TG""END__'
wrapped = 'printf "%%s\\n" %s; %s; printf "\\n%%s\\n" %s' % (START_Q, cmd, END_Q)


def sync_prompt(s, budget=6.0):
    """Wait until the shell prompt has actually appeared before sending anything.

    A cold serial open can land mid-prompt or mid-boot-noise; sending the wrapped
    command before the shell is listening loses it silently and looks exactly like
    "the device did not answer".
    """
    s.write(b"\n")
    s.flush()
    t0 = time.time()
    buf = b""
    while time.time() - t0 < budget:
        d = s.read(4096)
        if d:
            buf += d
            if buf.rstrip().endswith(b"#") or b"# " in buf[-200:]:
                s.read(4096)  # drain any trailing prompt bytes
                return True
    return False


def run_once(s, wait=90.0):
    s.reset_input_buffer()
    s.write((wrapped + "\n").encode())
    s.flush()
    buf = b""
    t0 = time.time()
    while time.time() - t0 < wait:
        d = s.read(4096)
        if d:
            buf += d
            if END.encode() in buf:
                break
    return buf


s = Serial(port, baud, timeout=0.5)
try:
    buf = b""
    for attempt in range(3):
        sync_prompt(s, budget=6.0)
        buf = run_once(s)
        if END.encode() in buf:
            break
        time.sleep(1.0)
finally:
    s.close()

txt = buf.decode("utf-8", "replace").replace("\r", "")
import os as _os
_dbg = _os.environ.get("FE2B_DEBUG_LOG", "")
if _dbg:
    with open(_dbg, "a") as _f:
        _f.write("PORT=%s CMD=%r\nRAW=%r\n\n" % (port, cmd, txt[:800]))
out = []
rec = False
for ln in txt.split("\n"):
    if ln.strip() == START:
        rec = True
        continue
    if ln.strip() == END:
        break
    if rec:
        out.append(ln)
# Drop trailing blank lines: the end marker is deliberately preceded by a newline
# (so a command whose last write has no trailing newline, e.g. `curl -w`, cannot
# glue itself to the marker). Leaving that newline in the payload would make
# every caller's `... | tail -1` read an empty string.
while out and not out[-1].strip():
    out.pop()
print("\n".join(out))
PY
}

TRANSPORT_READY=""
pick_transport() {
  case "$TRANSPORT" in
    ssh)
      if port_open "$ROUTER" 22; then TRANSPORT_READY=ssh; return 0; fi
      return 1 ;;
    serial)
      [ -c "$SERIAL_DEV" ] && { TRANSPORT_READY=serial; return 0; }
      return 1 ;;
    auto)
      if port_open "$ROUTER" 22; then TRANSPORT_READY=ssh; return 0; fi
      if [ -c "$SERIAL_DEV" ]; then TRANSPORT_READY=serial; write_serial_helper; return 0; fi
      return 1 ;;
    *) return 1 ;;
  esac
}

rsh() {
  _cmd="$1"
  case "${TRANSPORT_READY:-}" in
    ssh)
      if command -v sshpass >/dev/null 2>&1; then
        timeout 40 sshpass -p "$LAB_PASS" ssh $SSH_OPTS "root@$ROUTER" "$_cmd" 2>/dev/null
      else
        timeout 40 ssh $SSH_OPTS "root@$ROUTER" "$_cmd" 2>/dev/null
      fi ;;
    serial)
      write_serial_helper
      timeout 120 python3 "$SERIAL_HELPER" "$SERIAL_DEV" "$SERIAL_BAUD" "$_cmd" 2>>"${FE2B_DEBUG_LOG:-/dev/null}" ;;
    *) return 1 ;;
  esac
}

cleanup() { [ -n "$SERIAL_HELPER" ] && rm -f "$SERIAL_HELPER"; }
trap cleanup EXIT INT TERM

# ---------------------------------------------------------------------------
echo "=== fe2b_feed_check.sh — FE2b acceptance gate (TDD baseline; RED is the expected start) ==="
echo "router    : $ROUTER"
echo "feed_line : $FEED_LINE"
echo "feed_url  : $FEED_URL"
echo "transport : $TRANSPORT (requested)"
echo "host      : $(hostname 2>/dev/null || echo '?')  date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo

# ---------------------------------------------------------------------------
echo "--- A0: a transport to the device works ---"
if pick_transport; then
  echo "      using transport: $TRANSPORT_READY"
  # A liveness probe is retried: this UART is shared with other lanes on the bench,
  # so a single lost write is evidence about the write, not about the device.
  board=""
  i=1
  while [ "$i" -le 3 ]; do
    board="$(rsh 'cat /tmp/sysinfo/board_name 2>/dev/null' | tr -d '\r' | grep -i '[a-z]' | tail -1)"
    [ -n "$board" ] && break
    i=$((i + 1))
    sleep 2
  done
  if [ -n "$board" ]; then
    pass "A0 transport=$TRANSPORT_READY answered from device ($board)"
    info "board: $board"
  else
    fail "A0 transport=$TRANSPORT_READY did not answer (no board_name read back)"
    die_no_transport
  fi
else
  fail "A0 no transport: TCP $ROUTER:22 closed and $SERIAL_DEV not present"
  die_no_transport
fi

# ---------------------------------------------------------------------------
echo "--- A1: tollgate-wrt is installed ---"
installed="$(rsh "opkg list-installed 2>/dev/null | grep -i '^$PKG '" | tr -d '\r' | head -1)"
if [ -n "$installed" ]; then
  pass "A1 $PKG installed"
  info "$installed"
  INSTALLED_VERSION="$(printf '%s' "$installed" | sed -n "s/^$PKG - //p")"
else
  fail "A1 $PKG installed (absent from opkg list-installed)"
  INSTALLED_VERSION=""
fi

# ---------------------------------------------------------------------------
echo "--- A2: $PKG ORIGIN is the feed, not a local file path ---"
# opkg (1.x, as shipped on OpenWrt 24.10) records no Source/Filename field in
# /usr/lib/opkg/status, so status alone cannot prove origin. Three falsifiable
# conditions together can, and each one is independently checkable:
#   O1 no locally staged tollgate package left on the box (scp/`.ipk` path absent)
#   O2 a tollgate FEED INDEX exists in /var/opkg-lists/ (opkg update reached the feed)
#   O3 the installed version is the version that index advertises
if [ -z "$INSTALLED_VERSION" ]; then
  fail "A2 origin is the feed (undecidable: package not installed)"
else
  staged="$(rsh 'ls /tmp/*.ipk /tmp/*.apk /root/*.ipk /root/*.apk /www/*.ipk /www/*.apk 2>/dev/null | grep -i tollgate' | tr -d '\r' | head -3)"
  if [ -n "$staged" ]; then
    fail "A2[O1] origin is the feed: a LOCAL staged tollgate package is present (install was file-based, not feed-based)"
    printf '%s\n' "$staged" | while read -r l; do [ -n "$l" ] && echo "      local artifact: $l"; done
  else
    info "A2[O1] PASS — no locally staged tollgate package on the device"
  fi

  idx="$(rsh 'ls /var/opkg-lists/ 2>/dev/null | grep -i tollgate' | tr -d '\r' | head -1)"
  if [ -n "$idx" ]; then
    info "A2[O2] tollgate feed index present: /var/opkg-lists/$idx"
    idxver="$(rsh "awk '/^Package: $PKG\$/{f=1} f&&/^Version:/{print \$2; exit}' /var/opkg-lists/$idx" | tr -d '\r' | tail -1)"
    if [ -n "$idxver" ] && [ "$idxver" = "$INSTALLED_VERSION" ]; then
      info "A2[O3] installed version matches feed index: $idxver"
      if [ -z "$staged" ]; then
        pass "A2 origin is the feed (indexed version $idxver, no local artifact)"
      fi
    else
      fail "A2[O3] origin is the feed: index advertises '${idxver:-<none>}' but installed is '$INSTALLED_VERSION'"
    fi
  else
    fail "A2[O2] origin is the feed: no tollgate index in /var/opkg-lists/ (nothing to have installed from)"
  fi
fi

# ---------------------------------------------------------------------------
echo "--- A3: the $SERVICE service is running ---"
# The service's own init script is the authoritative reporter, and it parses the
# process table without the false positives of `pgrep -f` / `ps|grep`: on this box
# the installer (./tollgate-installer-linux-arm64) and the uhttpd portal processes
# both carry "tollgate" in their cmdline.
svc="$(rsh "$INITD running >/dev/null 2>&1 && echo RUNNING || echo NOTRUNNING" | tr -d '\r' | tail -1)"
enabled="$(rsh "ls /etc/rc.d/ 2>/dev/null | grep -i '$SERVICE'" | tr -d '\r' | head -1)"
info "binary: $(rsh "ls -l $BIN 2>/dev/null" | tr -d '\r' | tail -1)"
if [ "$svc" = "RUNNING" ]; then
  if [ -n "$enabled" ]; then
    pass "A3 $SERVICE running ($INITD reports RUNNING; $enabled enabled)"
  else
    fail "A3 $SERVICE running but NOT enabled in /etc/rc.d (will not survive reboot)"
  fi
else
  fail "A3 $SERVICE is not running ($INITD running -> ${svc:-<no answer>})"
fi

# ---------------------------------------------------------------------------
echo "--- A4: :$HEALTH_PORT health returns 2xx ---"
h="$(rsh "curl -s -o /dev/null -w '%{http_code}' -m 8 http://127.0.0.1:$HEALTH_PORT/" | tr -d '\r' | tail -1)"
case "$h" in
  2*) pass "A4 :$HEALTH_PORT health HTTP $h" ;;
  *)  fail "A4 :$HEALTH_PORT health HTTP ${h:-000} (expected 2xx)" ;;
esac

# ---------------------------------------------------------------------------
echo "--- A5: portal returns HTTP 200 ---"
p="$(rsh "curl -s -o /dev/null -w '%{http_code}' -m 8 http://127.0.0.1:$PORTAL_PORT$PORTAL_PATH" | tr -d '\r' | tail -1)"
if [ "$p" = "200" ]; then
  pass "A5 portal HTTP 200 at :$PORTAL_PORT$PORTAL_PATH"
else
  fail "A5 portal HTTP ${p:-000} at :$PORTAL_PORT$PORTAL_PATH (expected 200)"
fi
root="$(rsh "curl -s -o /dev/null -w '%{http_code}' -m 8 http://127.0.0.1:$PORTAL_PORT/" | tr -d '\r' | tail -1)"
info "note: :$PORTAL_PORT/ returns ${root:-000} by design (no index.html; listing disabled) — not asserted"

# ---------------------------------------------------------------------------
echo "--- A6: the src/gz feed index is fetchable (precondition of A2) ---"
code="$(curl -s -o /dev/null -w '%{http_code}' -L -m 25 "$FEED_URL" 2>/dev/null)"
if [ "$code" = "200" ]; then
  pass "A6 feed index fetchable: HTTP $code"
else
  fail "A6 feed index fetchable: HTTP ${code:-000} at $FEED_URL (no index -> nothing to install from)"
fi

# ---------------------------------------------------------------------------
echo
echo "=== RESULT: $PASS passed, $FAIL failed ==="
if [ "$FAIL" -eq 0 ]; then
  echo "FE2b ACCEPTANCE SATISFIED"
  exit 0
fi
echo "FE2b NOT SATISFIED — this is the expected TDD RED baseline before the feed install."
exit 1
