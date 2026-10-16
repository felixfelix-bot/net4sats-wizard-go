#!/usr/bin/env bash
# wz16v2 STEP 3b — artifact inventory + checksum count parity + sha256sum -c (file is SHA256SUMS)
set -u
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
D="$REPO/dist"
OUT="$CAP/05b-artifacts-and-checksums.txt"
SUMS="$D/SHA256SUMS"
cd "$D" || exit 9
{
  echo "### host: $(hostname) $(date -Is)"
  echo "### dist: $D"
  echo "### checksum file: dist/SHA256SUMS  (name_template from .goreleaser.yaml)"
  echo
  echo '# ---- artifact inventory ----'
  echo '$ ls -la dist/'
  ls -la "$D"
  echo
  echo '# ---- STEP 3b-1: archives (6 expected: linux/windows/darwin x amd64/arm64) ----'
  echo '$ ls dist/*.tar.gz dist/*.zip'
  ls "$D"/*.tar.gz "$D"/*.zip 2>/dev/null
  echo "ARCHIVE_COUNT=$(ls "$D"/*.tar.gz "$D"/*.zip 2>/dev/null | wc -l)"
  echo "  linux   = $(ls "$D"/*linux*.tar.gz 2>/dev/null | wc -l)  (expect 2: amd64, arm64)"
  echo "  windows = $(ls "$D"/*windows*.zip  2>/dev/null | wc -l)  (expect 2: amd64, arm64)"
  echo "  darwin  = $(ls "$D"/*darwin*.tar.gz 2>/dev/null | wc -l)  (expect 2: amd64, arm64)"
  echo "  amd64   = $(ls "$D"/*amd64* 2>/dev/null | wc -l)  (expect 3)"
  echo "  arm64   = $(ls "$D"/*arm64* 2>/dev/null | wc -l)  (expect 3)"
  echo
  echo '# ---- STEP 3b-2: extracted per-target binaries (the 6 build targets) ----'
  echo '$ find dist -type f -name "net4sats-wizard" | sort'
  find "$D" -type f -name 'net4sats-wizard' -o -type f -name 'net4sats-wizard.exe' | sort
  echo "BINARY_COUNT=$(find "$D" -type f \( -name 'net4sats-wizard' -o -name 'net4sats-wizard.exe' \) | wc -l)  (expect 6)"
  echo
  echo '# ---- STEP 3b-3: checksum file contents ----'
  echo '$ cat dist/SHA256SUMS'
  cat "$SUMS"
  echo "CHECKSUM_LINES=$(grep -c . "$SUMS")   ARCHIVE_COUNT=$(ls "$D"/*.tar.gz "$D"/*.zip 2>/dev/null | wc -l)"
  echo
  echo '# ---- STEP 3b-4: COUNT PARITY (checksum lines == archive count) ----'
  CL=$(grep -c . "$SUMS"); AC=$(ls "$D"/*.tar.gz "$D"/*.zip 2>/dev/null | wc -l)
  if [ "$CL" = "$AC" ]; then echo "PARITY: OK  ($CL checksum lines == $AC archives)"; else echo "PARITY: MISMATCH ($CL vs $AC)"; fi
  echo
  echo '# ---- STEP 3c: sha256sum -c (ALL artifacts) ----'
  echo '$ sha256sum -c SHA256SUMS'
  sha256sum -c "$SUMS" 2>&1
} > "$OUT" 2>&1
sha256sum -c "$SUMS" > /tmp/wz16v2-sha.out 2>&1
echo "SHA256SUM_C_ALL_RC=$?" >> "$OUT"
cat /tmp/wz16v2-sha.out >> "$OUT"
echo "OK_COUNT=$(grep -c ': OK$' /tmp/wz16v2-sha.out)  FAILED_COUNT=$(grep -c 'FAILED' /tmp/wz16v2-sha.out)"
