# PROGRESS — t_gate8-red: Red Playwright spec for generic 7-step wizard flow

Card: fleet offload worker, repo net4sats-mvp-v2
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Host: x280 / CobradorWave (c08r4d0r), Ubuntu 24.04.3 LTS
Deliverable: RED Playwright spec for the generic 7-step wizard + captured red run.
Task ENDS at a captured, referenced red run. Must NOT be made to pass.

---

## Step 0 — Inputs resolved from the env-pin record (PROGRESS-t_wz18-env-pin.md)

Finding: the env-pin card (F3) already resolved where the suite lives and where the
7-step spec currently exists:
  - Playwright suite  = `~/repos/physical-router-test-automation` (kit)
  - endo template     = kit `tests/browser/endo-onboarding.spec.mjs` + `.config.mjs`
  - existing 7-step   = kit `tests/browser/generic-wizard-7step.spec.mjs` on
                        `origin/feat/router-vm-lab` ONLY (tip ee289bf), not on main.
  - F4 server cmd     = wizard binds 127.0.0.1:8099 (main.go:48); build with
                        `/usr/local/go/bin/go`; `go` is not on non-login PATH.
  - F5 BLOCKER        = no npm/npx, no chrome/chromium, no ms-playwright bundle.
  - F7 drift          = spec asserts `#devsplit` == '10'; HEAD serves value 0.
Status: READ, carried forward. Files touched: none.

## Step 1 — F5 blocker CLEARED this session (runtime bootstrapped)

Finding: `apt-cache policy npm` -> candidate 9.2.0~ds1-2 (universe) and sudo works on
this host. Installed npm, then `npm install` + `npx playwright install --with-deps
chromium` in the kit. Result: npm 9.2.0, @playwright/test 1.59.1, and browser
bundles at ~/.cache/ms-playwright/{chromium-1217, chromium_headless_shell-1217,
ffmpeg-1011} with binary .../chromium-1217/chrome-linux64/chrome.
So the env-pin card's execution blocker is FALSE as of this run.
Status: DONE (blocker cleared, evidence in the red-run artifact).
Files touched: (host toolchain only; no repo file).

## Step 2 — Upstream pin re-resolved live

Finding: see artifact header. Pin recorded in the RED-RUN artifact as the exact
HEAD sha of net4sats-wizard-go at run time, plus the kit branch/sha.
Status: DONE.

## Step 3 — Spec authored (RED)

Finding: authored the generalized spec + a data-driven step model (wizard-flow.mjs)
so the step list/order/optionality is DATA, not scattered literals. Five
independently runnable cases (a-e of the brief) + the two Advanced cases are
title-filterable.
Status: DONE — see files touched below.

## Step 4 — Red run captured

Finding: see the artifact + REPORT.md.
Status: DONE.

---

## Files touched
- (this file)
- REPORT.md (full report)
