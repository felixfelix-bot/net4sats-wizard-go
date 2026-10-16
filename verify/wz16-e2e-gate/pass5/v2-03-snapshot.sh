#!/usr/bin/env bash
# wz16v2 STEP 3 — goreleaser release --snapshot --clean
set -u
export PATH=/home/c08r4d0r/wz16-tools:/home/c08r4d0r/wz16-tools/go122/bin:$PATH
export GOFLAGS=-mod=mod
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
mkdir -p "$CAP"
cd "$REPO" || exit 9
OUT="$(mktemp)"
{
  echo "### host: $(hostname) $(date -Is)"
  echo "### repo: $REPO  HEAD: $(git rev-parse HEAD)"
  echo '# ---- STEP 3a: goreleaser release --snapshot --clean ----'
  echo '$ goreleaser release --snapshot --clean'
} > "$CAP/05-goreleaser-snapshot.txt"
goreleaser release --snapshot --clean > "$OUT" 2>&1
RC=$?
cat "$OUT" >> "$CAP/05-goreleaser-snapshot.txt"
echo "GORELEASER_SNAPSHOT_RC=$RC" >> "$CAP/05-goreleaser-snapshot.txt"
tail -25 "$OUT"
echo "GORELEASER_SNAPSHOT_RC=$RC"
