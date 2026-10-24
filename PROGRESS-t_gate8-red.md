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

## Step 3 — Spec authored (RED), in the kit not this repo

Finding: the Playwright suite lives in the kit, so the spec lives there
(`tests/browser/`). Authored a data-driven step model (`wizard-flow.mjs`) so the
step list/order/optionality is DATA, plus 5 independently runnable cases (brief
a-e). Two Advanced cases are title-filterable:
`--grep "advanced SKIP path"` / `--grep "advanced VISITED path"`.
Status: DONE.
Files touched: kit `tests/browser/wizard-flow.mjs`,
`generic-wizard-7step.spec.mjs`, `generic-wizard-7step.config.mjs`.

## Step 4 — Red run captured against the pinned build

Finding: full run = 5 failed / 5 total, exit 1. Expected reason: the pinned build
implements NO step model — served document has 0 `data-testid` attributes, 2
`<h2>` (needs 7), 0 forward-or-skip controls, 0 step containers. Pin verified by
content: served bytes hash == `HEAD:index.html` ==
72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9.
Status: DONE.
Files touched: `evidence/wz1.8-red-generic-7step/` (6 files).

---

## Commits + push (verified live)

- `8fd477d` — spec artifact + this progress file.
- `9d849d1` — raw output, DOM audit, failure screenshot.
- Both pushed; `git ls-remote origin worker-heavy/t_gate8-19b1b1c6` == local SHA.
- Anonymous fetch of the artifact at `8fd477d` -> HTTP 200, 17603 bytes,
  sha256 `ad72354e…f89528` == on-disk. Artifact is publicly readable.
- Kit spec committed as `0600890` and pushed to `fork:feat/generic-7step-red`;
  spec + flow model anonymously readable (HTTP 200 each).

## Files touched
- `PROGRESS-t_gate8-red.md` (this file)
- `REPORT-t_gate8-red.md` (full report)
- `evidence/wz1.8-red-generic-7step/` (RED-RUN artifact, raw output, DOM audit,
  failure screenshot, reproduce script)
- kit `tests/browser/{wizard-flow.mjs, generic-wizard-7step.spec.mjs,
  generic-wizard-7step.config.mjs}`

## Explicitly NOT done
- No GL-MT6000 hardware run, no flash, no video capture (separate WZ1.8 task).
- The spec was NOT made to pass. Green is the follow-on task.
- `PROGRESS.md` / `REPORT.md` at the worktree root were NOT touched: root
  `PROGRESS.md` is tracked and shared, already carrying three other cards'
  sections. Task notes went to the task-scoped files above.
