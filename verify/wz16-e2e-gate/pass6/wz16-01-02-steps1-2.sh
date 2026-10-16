#!/usr/bin/env bash
# wz16 e2e gate — pass 6 (2026-10-16). Steps 1 + 2.
# Run in the pass-6 worktree, detached at the real deliverable tip.
set -uo pipefail
export PATH=/usr/local/go/bin:$HOME/wz16-tools:$HOME/go/bin:$PATH
BR=ci/add-github-actions-workflow-clean
WT=/home/c08r4d0r/worktrees/wz16-pass6
CAP=/home/c08r4d0r/wz16-captures-v6
mkdir -p "$CAP"
cd "$WT"

{
echo "### wz16 e2e integration gate — PASS 6 — host $(hostname) — $(date -Iseconds)"
echo "### subject branch: $BR"
echo "### worktree:       $WT"
echo "### HEAD:           $(git rev-parse HEAD)"
echo "### HEAD subject:   $(git log -1 --format=%s)"
echo "### tools:          $(go version) | golangci-lint $(golangci-lint --version | awk '{print $4}') | $(goreleaser --version | awk '/GitVersion/{print $2}')"
echo

echo "================ STEP 1 — clean tree check ================"
echo "\$ git status --porcelain"
porc=$(git status --porcelain)
printf '%s\n' "$porc"
echo "PORCELAIN_RC=$?  LINES=$(printf '%s' "$porc" | grep -c . )"
echo
echo "\$ git push          # literal, executed in the branch's own worktree (\$WT-owner: wz16-verify)"
git -C /home/c08r4d0r/wz16-verify push 2>&1
echo "GIT_PUSH_RC=$?"
echo
echo "\$ git push origin HEAD:refs/heads/$BR     # same op, from this detached pass worktree"
git push origin "HEAD:refs/heads/$BR" 2>&1
echo "PUSH_EXPLICIT_RC=$?"
echo
echo "\$ git rev-list --left-right --count HEAD...origin/$BR"
git rev-list --left-right --count "HEAD...origin/$BR"
echo "local_HEAD            = $(git rev-parse HEAD)"
echo "origin/$BR = $(git rev-parse origin/$BR)"
echo
echo "\$ git ls-remote origin refs/heads/$BR"
git ls-remote origin "refs/heads/$BR"
echo
echo "### deliverable files present at this tip (CI workflow + goreleaser config + docs):"
git ls-tree -r --name-only HEAD | grep -E '\.github/workflows/|\.goreleaser\.yaml|^docs/'
echo

echo "================ STEP 2 — full suite ================"
echo "\$ go vet ./..."
go vet ./... 2>&1
echo "GO_VET_RC=$?"
echo
echo "\$ golangci-lint run"
golangci-lint run 2>&1
echo "GOLANGCI_RC=$?"
echo
echo "\$ goreleaser check"
goreleaser check 2>&1
echo "GORELEASER_CHECK_RC=$?"
echo
} 2>&1 | tee "$CAP/01-02-steps1-2.txt"

echo
echo "================ STEP 2b — go test ./... -race (long) ================"
{
echo "\$ go test ./... -race"
go test ./... -race 2>&1
echo "GO_TEST_RACE_RC=$?"
} 2>&1 | tee "$CAP/03-test-race.txt"

echo
echo "================ STEP 2c — go test ./... -race -v (per-test detail) ================"
{
echo "\$ go test ./... -race -v -count=1"
go test ./... -race -v -count=1 2>&1 | tee "$CAP/04-test-race-verbose.log" | tail -20
echo "GO_TEST_RACE_VERBOSE_RC=${PIPESTATUS[0]}"
echo "--- per-test tally ---"
grep -c '^=== RUN'      "$CAP/04-test-race-verbose.log" | sed 's/^/RUN_COUNT=/'
grep -c '^--- PASS'     "$CAP/04-test-race-verbose.log" | sed 's/^/PASS_COUNT=/'
grep -c '^--- FAIL'     "$CAP/04-test-race-verbose.log" | sed 's/^/FAIL_COUNT=/'
grep -c '^--- SKIP'     "$CAP/04-test-race-verbose.log" | sed 's/^/SKIP_COUNT=/'
grep -c 'DATA RACE'     "$CAP/04-test-race-verbose.log" | sed 's/^/DATA_RACE_COUNT=/'
grep -E '^(ok|FAIL|---)' "$CAP/04-test-race-verbose.log" | tail -5
} 2>&1 | tee "$CAP/04-test-race-summary.txt"
echo "PASS6 STEPS 1-2 DONE"
