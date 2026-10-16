#!/usr/bin/env bash
# wz16 STEP 2b — go test ./... -race -count=1, plus a verbose run for RUN/FAIL/SKIP/RACE counts.
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
WT=/home/c08r4d0r/wz16-verify
export PATH="/home/c08r4d0r/wz16-tools/go122/bin:$PATH"
export GOTOOLCHAIN=go1.22.12
export GOCACHE="${GOCACHE:-/home/c08r4d0r/.cache/go-build-wz122}"
export CGO_ENABLED=1
cd "$WT" || exit 9

{
  echo "### host: $(hostname) $(id -un) $(date -Is)"
  echo "### go: $(go version)   cpus: $(nproc)   load: $(cut -d' ' -f1-3 /proc/loadavg)"
  echo "### gcc: $(gcc --version | head -1)   CGO_ENABLED=$CGO_ENABLED"
  echo "### worktree: $WT   HEAD: $(git rev-parse HEAD)"
  echo
  echo '# ---- 2b REQUIRED: go test ./... -race ----'
  echo '$ CGO_ENABLED=1 go test ./... -race -count=1'
  time CGO_ENABLED=1 go test ./... -race -count=1 2>&1
  echo "TEST_RACE_RC=$?"
} 2>&1 | tee "$CAP/03-test-race.txt"

{
  echo "### verbose re-run for RUN/FAIL/SKIP counters ($(date -Is))"
  echo '$ CGO_ENABLED=1 go test ./... -race -count=1 -v'
  CGO_ENABLED=1 go test ./... -race -count=1 -v 2>&1
  echo "TEST_V_RC=$?"
} > "$CAP/04-test-race-verbose.log" 2>&1

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### verbose log: $CAP/04-test-race-verbose.log"
  echo "RUN   = $(grep -c '^=== RUN' "$CAP/04-test-race-verbose.log")"
  echo "PASS  = $(grep -c '^--- PASS' "$CAP/04-test-race-verbose.log")"
  echo "FAIL  = $(grep -c '^--- FAIL' "$CAP/04-test-race-verbose.log")"
  echo "SKIP  = $(grep -c '^--- SKIP' "$CAP/04-test-race-verbose.log")"
  echo "DATA_RACE_WARNINGS = $(grep -c 'WARNING: DATA RACE' "$CAP/04-test-race-verbose.log")"
  echo "final ok line:"; grep -E '^(ok|FAIL|---)' "$CAP/04-test-race-verbose.log" | tail -3
} 2>&1 | tee "$CAP/04-test-race-summary.txt"
echo "WROTE $CAP/03-test-race.txt and 04-*"
