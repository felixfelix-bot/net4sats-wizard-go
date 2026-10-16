#!/usr/bin/env bash
# wz16v2 STEP 2e — go test ./... -race -v (test-level counts)
set -u
export PATH=/home/c08r4d0r/wz16-tools:/home/c08r4d0r/wz16-tools/go122/bin:$PATH
export GOFLAGS=-mod=mod
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
LOG="$CAP/04-test-race-verbose.log"
cd "$REPO" || exit 9
go test ./... -race -v > "$LOG" 2>&1
RC=$?
{
  echo "### host: $(hostname) $(date -Is)"
  echo "### cmd: go test ./... -race -v   (full log: 04-test-race-verbose.log)"
  echo "TEST_RACE_VERBOSE_RC=$RC"
  echo "--- PASS = $(grep -c -- '^--- PASS' "$LOG")"
  echo "--- FAIL = $(grep -c -- '^--- FAIL' "$LOG")"
  echo "--- SKIP = $(grep -c -- '^--- SKIP' "$LOG")"
  echo "ok_pkgs  = $(grep -c -- '^ok' "$LOG")"
  echo "DATA_RACE= $(grep -c 'DATA RACE' "$LOG")"
  echo "panic    = $(grep -c 'panic:' "$LOG")"
  echo "--- tail ---"
  grep -- '^--- FAIL' "$LOG" || echo "(no FAIL lines)"
  grep '^ok\|^FAIL\|^---' "$LOG" | tail -3
} > "$CAP/04-test-race-summary-verbose.txt"
cat "$CAP/04-test-race-summary-verbose.txt"
