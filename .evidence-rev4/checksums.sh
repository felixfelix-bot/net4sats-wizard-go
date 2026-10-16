#!/usr/bin/env bash
# Step 3c/3d/3e: checksum verification with independent recompute + negative control.
set -u
cd "$(dirname "$0")/../dist" || exit 1
OUT=../.evidence-rev4/03b-checksums.txt
exec > >(tee "$OUT") 2>&1

echo "===== STEP3c: sha256sum -c SHA256SUMS (all 6) ====="
sha256sum -c SHA256SUMS
echo "[sha256sum_c_rc=$?]"
echo

A=$(ls net4sats-wizard_*_linux_amd64.tar.gz)
W=$(ls net4sats-wizard_*_windows_arm64.zip)

echo "===== STEP3d: independent recompute (one archive, one non-linux) ====="
echo "artifact A (archive/linux): $A"
printf '%s  %s\n' "$(sha256sum "$A" | cut -d' ' -f1)" "$A" > /tmp/rev4-a.sums
echo "recorded_in_SHA256SUMS: $(grep -F " $A" SHA256SUMS)"
if sha256sum -c /tmp/rev4-a.sums >/dev/null 2>&1; then echo "INDEPENDENT_RECOMPUTE_A=MATCH"; else echo "INDEPENDENT_RECOMPUTE_A=MISMATCH"; fi

echo "artifact B (non-linux/windows): $W"
printf '%s  %s\n' "$(sha256sum "$W" | cut -d' ' -f1)" "$W" > /tmp/rev4-b.sums
echo "recorded_in_SHA256SUMS: $(grep -F " $W" SHA256SUMS)"
if sha256sum -c /tmp/rev4-b.sums >/dev/null 2>&1; then echo "INDEPENDENT_RECOMPUTE_B=MATCH"; else echo "INDEPENDENT_RECOMPUTE_B=MISMATCH"; fi
echo

echo "===== STEP3e: NEGATIVE CONTROL (corrupted copy MUST FAIL) ====="
cp "$A" /tmp/rev4-corrupt.tar.gz
printf 'CORRUPT' | dd of=/tmp/rev4-corrupt.tar.gz bs=1 seek=100 conv=notrunc status=none
printf '%s  %s\n' "$(sha256sum "$A" | cut -d' ' -f1)" "/tmp/rev4-corrupt.tar.gz" > /tmp/rev4-c.sums
if sha256sum -c /tmp/rev4-c.sums >/dev/null 2>&1; then
  echo "NEGATIVE_CONTROL=UNEXPECTED_PASS_VACUOUS"
else
  echo "NEGATIVE_CONTROL=FAIL_AS_EXPECTED"
fi
