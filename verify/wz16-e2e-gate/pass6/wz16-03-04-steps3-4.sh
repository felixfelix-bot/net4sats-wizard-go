#!/usr/bin/env bash
# wz16 e2e gate — pass 6. Steps 3 (snapshot release) + 4 (binary smoke + injection control).
set -uo pipefail
export PATH=/usr/local/go/bin:$HOME/wz16-tools:$HOME/go/bin:$PATH
WT=/home/c08r4d0r/worktrees/wz16-pass6
CAP=/home/c08r4d0r/wz16-captures-v6
cd "$WT"
DIST="$WT/dist"

{
echo "### tools: $(goreleaser --version | awk '/GitVersion/{print "goreleaser "$2}')  $(go version)"

echo
echo "================ STEP 3 — snapshot release ================"
echo "\$ git check-ignore -v dist/   # confirm dist is not a tracked path"
git check-ignore -v dist/ || echo "(dist/ not ignored — check)"
echo
echo "\$ goreleaser release --snapshot --clean"
goreleaser release --snapshot --clean 2>&1 | tail -40
echo "GORELEASER_SNAPSHOT_RC=${PIPESTATUS[0]}"

echo
echo "================ STEP 3b — artifact + checksum inventory ================"
echo "\$ ls -l dist/*.tar.gz dist/*.zip dist/SHA256SUMS"
ls -l "$DIST"/*.tar.gz "$DIST"/*.zip "$DIST"/SHA256SUMS
echo
echo "artifact list (name  size):"
find "$DIST" -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) -printf '%f\t%s\n' | sort
ARCHIVE_COUNT=$(find "$DIST" -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) | wc -l)
echo
echo "matrix breakdown:"
printf 'linux  = %s\n' "$(find "$DIST" -maxdepth 1 -name 'net4sats-wizard_*_linux_*' | wc -l)"
printf 'darwin = %s\n' "$(find "$DIST" -maxdepth 1 -name 'net4sats-wizard_*_darwin_*' | wc -l)"
printf 'windows= %s\n' "$(find "$DIST" -maxdepth 1 -name 'net4sats-wizard_*_windows_*' | wc -l)"
echo
echo "per-target built binaries:"
find "$DIST" -maxdepth 2 -type f -name 'net4sats-wizard*' -not -name '*.tar.gz' -not -name '*.zip' -printf '%P\t%s\n' | sort
BINARY_COUNT=$(find "$DIST" -maxdepth 2 -type f -name 'net4sats-wizard*' -not -name '*.tar.gz' -not -name '*.zip' | wc -l)
echo
echo "\$ wc -l dist/SHA256SUMS"
wc -l "$DIST/SHA256SUMS"
CHECKSUM_LINES=$(wc -l < "$DIST/SHA256SUMS")
echo
echo "ARCHIVE_COUNT=$ARCHIVE_COUNT  BINARY_COUNT=$BINARY_COUNT  CHECKSUM_LINES=$CHECKSUM_LINES"
if [ "$CHECKSUM_LINES" -eq "$ARCHIVE_COUNT" ]; then echo "PARITY: OK  ($CHECKSUM_LINES checksum lines == $ARCHIVE_COUNT archives)"; else echo "PARITY: MISMATCH"; fi

echo
echo "================ STEP 3c — checksum verification ================"
echo "\$ cd dist && sha256sum -c SHA256SUMS"
( cd "$DIST" && sha256sum -c SHA256SUMS 2>&1 )
SHA_RC=$?
echo "SHA256SUM_C_RC=$SHA_RC   OK_COUNT=$( ( cd "$DIST" && sha256sum -c SHA256SUMS 2>/dev/null ) | grep -c ': OK$')  FAILED_COUNT=$( ( cd "$DIST" && sha256sum -c SHA256SUMS 2>/dev/null ) | grep -c ': FAILED')"
echo
echo "-- required: one archive (linux/amd64 .tar.gz) + one non-linux artifact (darwin/arm64 .tar.gz) --"
A=$(ls "$DIST"/net4sats-wizard_*_linux_amd64.tar.gz)
B=$(ls "$DIST"/net4sats-wizard_*_darwin_arm64.tar.gz)
echo "A (archive,    linux/amd64)  = $(basename "$A")"
echo "B (non-linux,  darwin/arm64) = $(basename "$B")"
( cd "$DIST" && sha256sum -c SHA256SUMS 2>&1 | grep -E "$(basename "$A")|$(basename "$B")" )
echo "TWO_ARTIFACT_VERIFY_RC=$?"
echo
echo "-- independent recompute (not trusting the sums file), vs the published line --"
( cd "$DIST" && sha256sum "$(basename "$A")" "$(basename "$B")" )
echo "published lines:"; grep -E "$(basename "$A")|$(basename "$B")" "$DIST/SHA256SUMS"
echo
echo "-- negative control: a corrupted copy MUST fail, proving the check is non-vacuous --"
TC=$(mktemp -d)
cp "$DIST/$(basename "$A")" "$TC/"; cp "$DIST/SHA256SUMS" "$TC/"
printf 'x' | dd of="$TC/$(basename "$A")" bs=1 seek=100 conv=notrunc status=none
( cd "$TC" && sha256sum -c SHA256SUMS 2>&1 | grep -F "$(basename "$A")" ; echo "NEGATIVE_CONTROL_RC=${PIPESTATUS[0]}  (non-zero = control is real)" )
rm -rf "$TC"
echo
echo "note: Windows 'Get-FileHash' was not executed (no Windows host available here)."
echo "      The equivalent integrity proof is the SHA-256 comparison above; both implement the same"
echo "      FIPS 180-4 SHA-256 over the file bytes. Stated rather than claimed."

echo
echo "================ STEP 4 — binary smoke (--version) ================"
BIN="$DIST/net4sats-wizard_linux_amd64_v1/net4sats-wizard"
echo "\$ file $BIN"
file "$BIN"
echo
echo "\$ $BIN --version"
OUT=$("$BIN" --version 2>/tmp/p6-verr); RC=$?
ERR=$(cat /tmp/p6-verr)
echo "STDOUT: $OUT"
echo "STDERR: ${ERR:-(empty)}"
echo "EXIT_CODE=$RC"
echo "LINE=[$OUT]"
case "$OUT" in
  tollgate-wizard\ *) echo "FORMAT:   OK   (prints 'tollgate-wizard <version>')";;
  *) echo "FORMAT:   FAIL";;
esac
[ "$RC" -eq 0 ] && echo "EXIT:     OK   (0)" || echo "EXIT:     FAIL"
echo "$OUT" | grep -q ' dev$' && echo "INJECTION: FAIL (printed 'dev')" || echo "INJECTION: OK   (non-dev version, not a fallback constant)"
echo
echo "-- manual injection control: build the same package with two different -X values --"
go build -ldflags "-X main.version=v9.9.9-pass6-control" -o /tmp/p6-ctl . && echo "CTL_OUT=[$(/tmp/p6-ctl --version)]"
go build -ldflags "-X main.version=v0.0.1-changed"      -o /tmp/p6-ctl2 . && echo "CTL2_OUT=[$(/tmp/p6-ctl2 --version)]"
go build -o /tmp/p6-nold . && echo "NOLD_OUT=[$(/tmp/p6-nold --version)]   (no ldflags -> 'dev' is the fallback)"
echo
echo "INJECTION_CHANGES_WITH_FLAG: $( [ "$(/tmp/p6-ctl --version)" != "$(/tmp/p6-ctl2 --version)" ] && echo OK || echo FAIL )"
echo "SNAPSHOT_VALUE_IS_STRIPPED_BINARY: $( [ -n "$OUT" ] && echo "$OUT" | grep -q 'snapshot' && echo OK || echo 'n/a')"

echo
echo "================ STEP 3/4 SUPPORT — code-tree identity ================"
echo "\$ git rev-parse HEAD"
git rev-parse HEAD
echo "HEAD^{tree}         = $(git rev-parse HEAD^{tree})"
echo "main.go blob        = $(git rev-parse HEAD:main.go)"
echo ".goreleaser.yaml    = $(git rev-parse HEAD:.goreleaser.yaml)"
echo ".github/workflows/ci.yml  = $(git rev-parse HEAD:.github/workflows/ci.yml)"
echo "\$ git status --porcelain   # after build/snapshot: build products must not dirty the tree"
git status --porcelain
echo "POST_BUILD_PORCELAIN_LINES=$(git status --porcelain | grep -c . )"
} 2>&1 | tee "$CAP/05-06-steps3-4.txt"
echo "PASS6 STEPS 3-4 DONE"
