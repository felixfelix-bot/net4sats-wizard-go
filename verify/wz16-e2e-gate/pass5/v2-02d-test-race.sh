#!/usr/bin/env bash
# wz16v2 STEP 2d — go test ./... -race (raw capture, real exit code)
set -u
export PATH=/home/c08r4d0r/wz16-tools:/home/c08r4d0r/wz16-tools/go122/bin:$PATH
export GOFLAGS=-mod=mod
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
LOG=/tmp/wz16v2-race.log
mkdir -p "$CAP"
cd "$REPO" || exit 9
{
  echo "### host: $(hostname) $(date -Is)"
  echo "### go: $(go version)"
  echo '# ---- STEP 2d: go test ./... -race ----'
  echo '$ go test ./... -race'
} > "$CAP/03-test-race.txt"
go test ./... -race > "$LOG" 2>&1
RC=$?
cat "$LOG" >> "$CAP/03-test-race.txt"
echo "TEST_RACE_RC=$RC" >> "$CAP/03-test-race.txt"

{
  echo "TEST_RACE_RC=$RC"
  echo "PASS      = $(grep -c -- '^--- PASS' "$LOG")"
  echo "FAIL      = $(grep -c -- '^--- FAIL' "$LOG")"
  echo "ok_pkgs   = $(grep -c -- '^ok' "$LOG")"
  echo "notest_pkgs = $(grep -c 'no test files' "$LOG")"
  echo "DATA_RACE = $(grep -c 'DATA RACE' "$LOG")"
  echo "panic     = $(grep -c 'panic:' "$LOG")"
} > "$CAP/04-test-race-summary.txt"

cat "$CAP/04-test-race-summary.txt"
echo "WROTE $CAP/03-test-race.txt $CAP/04-test-race-summary.txt"
