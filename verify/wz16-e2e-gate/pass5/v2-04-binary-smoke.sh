#!/usr/bin/env bash
# wz16v2 STEP 4 (corrected) — binary smoke + injection control.
# programName is pinned in main.go:55 as "tollgate-wizard", so the required
# --version line is "tollgate-wizard <version>" (NOT the repo name).
set -u
export PATH=/home/c08r4d0r/wz16-tools:/home/c08r4d0r/wz16-tools/go122/bin:$PATH
export GOFLAGS=-mod=mod
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
BIN="$REPO/dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard"
OUT="$CAP/06-version-smoke.txt"
cd "$REPO" || exit 9

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### binary under test: $BIN"
  echo "### expected format (task): 'tollgate-wizard <version>'  <- pinned by main.go:55 programName"
  echo '$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard'
  file "$BIN"
  echo '$ ls -la dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard'
  ls -la "$BIN"
  echo
  echo '# ---- STEP 4a: execute built linux/amd64 binary with --version ----'
  echo '$ ./dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version'
} > "$OUT"

"$BIN" --version > /tmp/wz16v2-ver.out 2>/tmp/wz16v2-ver.err
RC=$?
LINE=$(cat /tmp/wz16v2-ver.out)
{
  echo "STDOUT: $LINE"
  echo "STDERR: $(cat /tmp/wz16v2-ver.err)"
  echo "VERSION_EXIT_CODE=$RC"
  echo "LINE=[$LINE]"
} >> "$OUT"

case "$LINE" in
  "tollgate-wizard "*) echo "FORMAT: OK  (prints 'tollgate-wizard <version>')" >> "$OUT" ;;
  *) echo "FORMAT: FAIL (want prefix 'tollgate-wizard ')" >> "$OUT" ;;
esac
[ "$RC" = 0 ] && echo "EXIT_CODE: OK (0)" >> "$OUT" || echo "EXIT_CODE: FAIL ($RC)" >> "$OUT"
if echo "$LINE" | grep -q 'dev'; then
  echo "INJECTION: FAIL (printed 'dev' -> ldflags did NOT inject)" >> "$OUT"
else
  echo "INJECTION: OK (non-dev version printed: $(echo "$LINE" | awk '{print $2}'))" >> "$OUT"
fi

echo >> "$OUT"
echo '# ---- STEP 4b: manual ldflags injection control (value must CHANGE) ----' >> "$OUT"
echo '$ go build -ldflags "-X main.version=v9.9.9-injection-control" -o /tmp/wz16v2-ctl .' >> "$OUT"
go build -ldflags "-X main.version=v9.9.9-injection-control" -o /tmp/wz16v2-ctl . 2>>"$OUT"
echo "CTL_BUILD_RC=$?" >> "$OUT"
CTL=$(/tmp/wz16v2-ctl --version)
echo "CTL_STDOUT: $CTL" >> "$OUT"
echo "CTL_LINE=[$CTL]" >> "$OUT"
if [ "$CTL" = "tollgate-wizard v9.9.9-injection-control" ]; then
  echo "CONTROL: OK (injected value flows through; '$(echo "$CTL" | awk '{print $2}')' != snapshot '$(echo "$LINE" | awk '{print $2}')')" >> "$OUT"
else
  echo "CONTROL: FAIL (unexpected: $CTL)" >> "$OUT"
fi

echo >> "$OUT"
echo '# ---- STEP 4c: negative control (no ldflags -> must print dev) ----' >> "$OUT"
echo '$ go build -o /tmp/wz16v2-nold .' >> "$OUT"
go build -o /tmp/wz16v2-nold . 2>>"$OUT"
echo "NOLD_BUILD_RC=$?" >> "$OUT"
NOLD=$(/tmp/wz16v2-nold --version)
echo "NOLD_LINE=[$NOLD]" >> "$OUT"
if [ "$NOLD" = "tollgate-wizard dev" ]; then
  echo "NEGATIVE: OK (no injection -> 'dev', proves the snapshot's version was truly injected)" >> "$OUT"
else
  echo "NEGATIVE: FAIL (unexpected: $NOLD)" >> "$OUT"
fi

cat "$OUT"
