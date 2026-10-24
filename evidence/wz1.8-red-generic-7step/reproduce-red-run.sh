#!/usr/bin/env bash
# Reproduce the RED run for the generic 7-step wizard spec.
#
# Captured context (2026-10-24, x280 / CobradorWave):
#   repo under test : felixfelix-bot/net4sats-wizard-go @ 867d45604de243c41406e387ea7cf0a88e78d859
#   wizard server   : /tmp/wz18-wizard (sha256 43f0ae6456a2fdc842416526835418123ee71983d6cd2c2f4c5da1478741422e)
#                     listening on 127.0.0.1:8099
#   suite           : ~/repos/physical-router-test-automation
#   tracked lane    : tests/browser/generic-wizard-7step.{spec,config}.mjs
set -euo pipefail

KIT="${KIT:-$HOME/repos/physical-router-test-automation}"
cd "$KIT/tests/browser"

# 0. Confirm the server under test is the pinned build.
echo "--- served index.html sha256 (must equal the pinned repo HEAD index.html) ---"
curl -s -m 10 "${WIZARD_URL:-http://127.0.0.1:8099}/" | sha256sum

# 1. Full red run (all 5 cases).
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs

# 2. The two Advanced cases, individually.
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs --grep "advanced SKIP path"
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs --grep "advanced VISITED path"
