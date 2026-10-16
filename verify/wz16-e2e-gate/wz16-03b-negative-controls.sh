#!/usr/bin/env bash
# wz16 STEP 3b — non-vacuity controls for the checksum gate (positive + two negatives).
# Lesson carried from rev3: use EXPLICIT filenames (glob args to dd do not expand) and read the
# REAL exit code of sha256sum via a temp file (never through a pipe, which reports grep's rc).
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
DIST=/home/c08r4d0r/wz16-verify/dist
WORK=/home/c08r4d0r/wz16-tests-t_0509c320-p4/negctl
rm -rf "$WORK"; mkdir -p "$WORK"; cd "$WORK" || exit 9
cp "$DIST"/*.tar.gz "$DIST"/*.zip "$DIST/SHA256SUMS" . || exit 9
OUT="$WORK/sha256sum-check.out"
LINUX=$(ls -1 *_linux_amd64.tar.gz)
WINARM=$(ls -1 *_windows_arm64.zip)

report() {   # report <label>
  sha256sum -c SHA256SUMS > "$OUT" 2>&1
  local rc=$?
  echo "--- $1 ---"
  echo '$ sha256sum -c SHA256SUMS'
  cat "$OUT"
  echo "REAL_RC=$rc"
  echo "OK=$(grep -c ': OK$' "$OUT")  FAILED=$(grep -c 'FAILED' "$OUT")"
  echo
}

{
  echo "### host: $(hostname) $(date -Is)   scratch: $WORK"
  echo "### copies of the real snapshot artifacts + the real SHA256SUMS"
  echo
  report "CONTROL A (positive): pristine copies -- expect 6/6 OK, REAL_RC=0"
  echo "### tamper: flip one byte of the EXPLICIT path $LINUX (printf, not a dd glob)"
  printf 'X' | dd of="$LINUX" bs=1 seek=100 conv=notrunc status=none
  echo "tamper_cmd_rc=$?   size_after=$(stat -c%s "$LINUX")"
  report "CONTROL B (negative): 1 byte flipped in $LINUX -- expect that line FAILED, REAL_RC=1"
  echo "### restore, then remove a DIFFERENT artifact (explicit): $WINARM"
  cp "$DIST/$LINUX" .        # restore
  rm -f "$WINARM"
  report "CONTROL C (negative): $WINARM missing -- expect FAILED open or read, REAL_RC=1"
} 2>&1 | tee "$CAP/08-checksum-controls.txt"
echo "WROTE $CAP/08-checksum-controls.txt"
