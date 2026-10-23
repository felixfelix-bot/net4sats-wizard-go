#!/usr/bin/env bash
# Reproduce the local equivalent of .github/workflows/ci.yml.
#
# Card t_8d1115d4 — "Verify CI parity locally: vet, test -race, lint, actionlint".
# Runs exactly the commands the workflow runs (same order), records each exit
# code, and additionally lints the workflow file itself with actionlint.
#
# Usage:  bash verify/ci-parity/run-ci-parity.sh [output-log]
#         (default log: verify/ci-parity/last-run.log)
#
# Exit code: 0 only if every step is green.
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$REPO_ROOT" || exit 99
LOG="${1:-verify/ci-parity/last-run.log}"
mkdir -p "$(dirname "$LOG")"
: > "$LOG"

# The dev box keeps the toolchain outside PATH; CI installs it via setup-go.
if ! command -v go >/dev/null 2>&1 && [ -x /usr/local/go/bin/go ]; then
  export PATH="/usr/local/go/bin:$PATH"
fi
# actionlint is not packaged here; `go install github.com/rhysd/actionlint/cmd/actionlint@v1.7.7`
# drops it in GOPATH/bin.
if ! command -v actionlint >/dev/null 2>&1 && [ -x "$(go env GOPATH)/bin/actionlint" ]; then
  export PATH="$(go env GOPATH)/bin:$PATH"
fi

FAILED=0

step() {
  local label="$1"; shift
  local out rc
  out="$("$@" 2>&1)"; rc=$?
  {
    echo "=== $label"
    echo "\$ $*"
    echo "exit_code=$rc"
    [ -n "$out" ] && echo "--- output ---" && echo "$out"
    echo
  } >> "$LOG"
  printf '%-46s exit=%d\n' "$label" "$rc"
  [ -n "$out" ] && printf '%s\n' "$out" | tail -5 | sed 's/^/    /'
  [ "$rc" -ne 0 ] && FAILED=1
  return 0
}

echo "host: $(uname -a)"            | tee -a "$LOG"
echo "go:   $(go version)"          | tee -a "$LOG"
echo "lint: $(golangci-lint --version | head -1)" | tee -a "$LOG"
echo "actionlint: $(actionlint --version | head -1)" | tee -a "$LOG"
echo                                  | tee -a "$LOG"

step "actionlint .github/workflows/ci.yml" actionlint -no-color .github/workflows/ci.yml
step "go build ./..."                      go build ./...
step "go vet ./..."                        go vet ./...
step "go test ./... -race -count=1"        go test ./... -race -count=1
step "golangci-lint run"                   golangci-lint run

if [ "$FAILED" -ne 0 ]; then
  echo "RESULT: FAIL — see $LOG"
  exit 1
fi
echo "RESULT: PASS — see $LOG"
