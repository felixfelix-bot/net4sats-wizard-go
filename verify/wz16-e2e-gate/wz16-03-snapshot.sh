#!/usr/bin/env bash
# wz16 STEP 3 — snapshot release + authoritative artifact inventory + checksum verification.
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
WT=/home/c08r4d0r/wz16-verify
export PATH="/home/c08r4d0r/wz16-tools/go122/bin:$PATH"
export GOTOOLCHAIN=go1.22.12
export GOCACHE="${GOCACHE:-/home/c08r4d0r/.cache/go-build-wz122}"
cd "$WT" || exit 9

{
  echo "### host: $(hostname) $(id -un) $(date -Is)"
  echo "### worktree: $WT   HEAD: $(git rev-parse HEAD)"
  echo "### goreleaser: $(/home/c08r4d0r/wz16-tools/goreleaser --version 2>/dev/null | grep -iE 'version=' | head -1)"
  echo
  echo '# ---- REQUIRED: goreleaser release --snapshot --clean ----'
  echo '$ goreleaser release --snapshot --clean'
  /home/c08r4d0r/wz16-tools/goreleaser release --snapshot --clean 2>&1
  echo "SNAPSHOT_RC=$?"
} 2>&1 | tee "$CAP/05-goreleaser-snapshot.txt"

# artifact inventory + checksums -- run from dist/ so -c resolves bare filenames
cd "$WT/dist" || exit 9
{
  echo "### host: $(hostname) $(date -Is)   cwd: $(pwd)"
  echo
  echo '$ ls -1 *.tar.gz *.zip | sort'
  ls -1 *.tar.gz *.zip | sort
  echo "ARCHIVE_COUNT=$(ls -1 *.tar.gz *.zip | wc -l)"
  echo
  echo '# ---- REQUIRED: 6 artifacts, target matrix = {linux,windows,darwin} x {amd64,arm64} ----'
  echo '$ ls -1 *.tar.gz *.zip | sed -E "s/.*_(linux|windows|darwin)_(amd64|arm64)\.(tar\.gz|zip)/\1_\2/" | sort'
  ls -1 *.tar.gz *.zip | sed -E 's/.*_(linux|windows|darwin)_(amd64|arm64)\.(tar\.gz|zip)/\1_\2/' | sort
  echo "distinct_targets=$(ls -1 *.tar.gz *.zip | sed -E 's/.*_(linux|windows|darwin)_(amd64|arm64)\.(tar\.gz|zip)/\1_\2/' | sort -u | wc -l)"
  echo
  echo '$ for a in *.tar.gz *.zip; do echo "$a -> $(file -b "$a")"; done'
  for a in *.tar.gz *.zip; do echo "$a -> $(file -b "$a")"; done
  echo
  echo '# ---- REQUIRED: checksum file exists and its count matches the artifact count ----'
  echo '$ ls -la SHA256SUMS'
  ls -la SHA256SUMS
  echo '$ wc -l < SHA256SUMS'
  wc -l < SHA256SUMS
  echo "$(wc -l < SHA256SUMS)  ==  $(ls -1 *.tar.gz *.zip | wc -l)  (checksum_count == artifact_count)"
  echo
  echo '# ---- REQUIRED: sha256sum -c verification (REAL exit code, no pipe) ----'
  echo '$ sha256sum -c SHA256SUMS'
  sha256sum -c SHA256SUMS > "$CAP/sha256sum-check.out" 2>&1
  real_rc=$?
  cat "$CAP/sha256sum-check.out"
  echo "SHA256SUM_C_REAL_RC=$real_rc"
  echo "OK_LINES=$(grep -c ': OK$' "$CAP/sha256sum-check.out")  FAILED_LINES=$(grep -c 'FAILED' "$CAP/sha256sum-check.out")"
  echo
  echo '# ---- REQUIRED: two named spot checks recomputed independently (one archive, one non-linux) ----'
  LINUX=$(ls -1 *_linux_amd64.tar.gz)
  DARWIN=$(ls -1 *_darwin_arm64.tar.gz)
  WIN=$(ls -1 *_windows_amd64.zip)
  for f in "$LINUX" "$DARWIN" "$WIN"; do
    rec=$(sha256sum "$f" | cut -d' ' -f1)
    exp=$(grep -F "$f" SHA256SUMS | awk '{print $1}')
    if [ "$rec" = "$exp" ]; then verdict=MATCH; else verdict=MISMATCH; fi
    echo "artifact=$f type=$([ "${f##*.}" = zip ] && echo zip || echo tar.gz) recomputed=$rec expected=$exp $verdict"
  done
  echo "  (spot-check#1 archive   = $LINUX)"
  echo "  (spot-check#2 non-linux = $DARWIN ; plus Windows zip equivalent of Get-FileHash: $WIN)"
  echo "POWERSHELL_PRESENT=$(command -v pwsh || echo no)"
} 2>&1 | tee "$CAP/05b-artifacts-and-checksums.txt"
echo "WROTE 05-* and 05b-*"
