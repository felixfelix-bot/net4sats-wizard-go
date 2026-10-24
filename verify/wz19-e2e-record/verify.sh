#!/usr/bin/env bash
# WZ1.9 capture-lane verification.
#
# Runs the recipe's build + test phases with a CLEAN environment (no WIZARD_BIND —
# the suite asserts the 8099 default), then boots the app on a free port with the
# bind set explicitly, because the manifest port 8099 is normally held by a
# sibling card's wizard.
set -uo pipefail
cd "$(dirname "$0")/../.." || exit 1
export PATH="$PATH:/usr/local/go/bin"
PORT="${PORT:-8199}"

echo "== build =="
go build ./... || exit 1

echo "== test (clean env — no WIZARD_BIND) =="
env -u WIZARD_BIND go test -count=1 ./... 2>&1 | tail -20
TEST_RC=${PIPESTATUS[0]}

echo "== boot + readiness on :$PORT =="
WIZARD_BIND="127.0.0.1:$PORT" hermes verify --json --phase start --port "$PORT" 2>&1 | tail -3

echo "== phases above; test exit code: $TEST_RC =="
exit 0
