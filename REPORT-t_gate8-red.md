# REPORT — t_gate8-red: Red Playwright spec for the generic 7-step wizard flow

**Status: DONE.** Spec authored, RED run captured, evidence committed and pushed,
artifact publicly verifiable.

Repo: `net4sats-mvp-v2` (worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`)
Branch: `worker-heavy/t_gate8-19b1b1c6`
Host: x280 / CobradorWave (c08r4d0r), Ubuntu 24.04.3 LTS

---

## 1. What was delivered

A Playwright spec that drives the generic 7-step wizard flow
(Discover → Connect → Uplink → Payouts → Advanced (optional) → Deploy → Live)
and **fails** against the pinned upstream build. It was not made to pass —
making it green is the follow-on task.

It lives in the Playwright kit, which is where the suite actually lives
(the env-pin record, `PROGRESS-t_wz18-env-pin.md` F3, resolved this):

| File | Role |
|---|---|
| `~/repos/physical-router-test-automation/tests/browser/wizard-flow.mjs` | the step list, declared order and optionality **as data** |
| `~/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs` | 5 independently runnable cases |
| `~/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.config.mjs` | deterministic config (retries 0, workers 1, channel chromium) |

The evidence artifact lives in the task's repo, under the repo's existing
evidence convention:

```
evidence/wz1.8-red-generic-7step/
├── RED-RUN-generic-7step-t_gate8.md      # main artifact: command, pin, reason, verbatim output
├── raw-playwright-output.txt             # verbatim stdout/stderr of the full run
├── served-dom-audit.txt                  # DOM audit of the pinned build
├── screenshot-failure-advanced-visited.png
└── reproduce-red-run.sh                  # one-shot reproduction
```

The repo's evidence convention is `docs/`-plus-a-conventional-dir; this task
used `evidence/` (documented in the artifact) because it is the closest match
to the brief's suggested `evidence/` / `test/evidence/` and does not collide
with the existing `docs/evidence/*.md` measured-capture files.

## 2. The exact command

```
cd ~/repos/physical-router-test-automation/tests/browser
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs
```

Result: **5 failed / 5 total, exit code 1.**

The two Advanced cases, individually runnable — exact title filters:

```
... --config=generic-wizard-7step.config.mjs --grep "advanced SKIP path"
... --config=generic-wizard-7step.config.mjs --grep "advanced VISITED path"
```

Both select exactly 1 test and both fail for the expected reason. Verified.

## 3. Expected failure reason (stated explicitly)

**The pinned build implements no step model.** `index.html` renders one flat
`#select-view` form — Discover/Connect/Uplink/Payouts/Advanced controls are all
merged into a single card — with no per-step container, no forward-navigation
control, no per-step gating, and no terminal `Live` step carrying post-deploy
evidence.

Measured against the *served* document (`served-dom-audit.txt`):

| Probe | Pinned build | Spec requires |
|---|---|---|
| `data-testid` attributes | **0** | 7 (`step-discover` … `step-live`) |
| `<h2>` step headings | **2** (`net4sats is live!`, `Setup failed`) | 7 |
| forward / skip controls | **0** | a `Next`/`Continue` per non-terminal step + `Skip` on Advanced |
| step containers | **0** | 7 |

This is a **missing-generic-wizard-behaviour** failure. It is **not** a typo,
import error, or environment failure: the spec parses, the browser launches,
the server serves HTTP 200, and every assertion reaches a live DOM before
failing on a wizard-contract assertion.

**Why this is not a wrong-reason red.** The spec intercepts `/api/scan`,
`/api/wifi-scan`, `/api/deploy` and `/api/status/<job>` and doubles them with
deterministic fixtures. The doubled status payload is a *completed* deployment
(all 11 backend steps `done`, job `status: 'done'`). So no router, no LAN and
no network latency participates, and the only thing that can keep the flow from
reaching `Live` is missing wizard UI — precisely the property under test. This
mattered: no router is reachable from this host at all (`/api/scan` returns no
routers; no route to `192.168.1.1`), so a hardware-dependent spec would have
failed for a *network* reason.

## 4. Pin — verified by content, not assumption

| Item | Value |
|---|---|
| Repo under test | `felixfelix-bot/net4sats-wizard-go` |
| Repo HEAD at capture | `867d45604de243c41406e387ea7cf0a88e78d859` |
| Served `index.html` sha256 | `72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9` |
| Repo `HEAD:index.html` sha256 | `72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9` |
| Server binary | `/tmp/wz18-wizard`, sha256 `43f0ae6456a2fdc842416526835418123ee71983d6cd2c2f4c5da1478741422e` |
| Server bind | `127.0.0.1:8099` (`main.go:48` `defaultListenAddr`) |
| Playwright | `@playwright/test` 1.59.1 |
| Browser | `channel: 'chromium'` → `~/.cache/ms-playwright/chromium-1217/chrome-linux64/chrome` |

The two `index.html` hashes being identical is the proof that the running
server is the pinned source — not merely that a server is running.

## 5. Acceptance checklist

- [x] Spec committed with conventional message `test(wizard): ...`
      (kit commit `0600890`; repo commit `8fd477d`).
- [x] Red output artifact committed **and** referenced with its exact command
      and the stated expected failure reason.
- [x] The two Advanced cases are individually runnable; exact title filters
      recorded (`--grep "advanced SKIP path"`, `--grep "advanced VISITED path"`).
- [x] Coverage a–e, each an independently runnable case filtered by title.
- [x] Determinism: no sleeps or wall-clock timing; explicit UI-state waits only;
      role/`getByTestId` selectors; per-case reset/cleanup.
- [x] Wizard definition is **data** (`wizard-flow.mjs`), not scattered literals.
- [x] Branch pushed; remote ref verified to equal the local SHA.

## 6. Commits and verification

| Repo | Commit | Remote verification |
|---|---|---|
| `net4sats-mvp-v2` | `8fd477d` spec artifact + progress | pushed to `origin/worker-heavy/t_gate8-19b1b1c6` |
| `net4sats-mvp-v2` | `9d849d1` raw output + DOM audit + screenshot | pushed (branch tip) |
| kit | `0600890` spec + model + config | pushed to `fork:feat/generic-7step-red` |

Verified live (observed, not assumed):

- `git ls-remote origin worker-heavy/t_gate8-19b1b1c6` → `8fd477d…` == local HEAD at push time.
- `git ls-remote fork feat/generic-7step-red` → `0600890…` == local kit HEAD.
- Anonymous fetch (no auth) of the artifact at the pushed SHA →
  `raw.githubusercontent.com` **HTTP 200, 17603 bytes**, and its sha256 equals
  the on-disk file byte-for-byte (`ad72354e…f89528`).
- Anonymous fetch of the spec and the flow model at kit SHA `0600890` →
  **HTTP 200** each.

## 7. Notes, corrections, and things I did not do

**The env-pin card's F5 blocker is stale.** F5 recorded "no browser runtime on
this host — node 18 present but npx/npm absent, no chrome, no ms-playwright
bundle — the endo spec CANNOT be executed as-is." That was true then and is
false now: `apt-cache policy npm` offers 9.2.0 (universe) and sudo works here.
I installed npm, ran `npm install`, and `npx playwright install --with-deps
chromium` in the kit, giving `@playwright/test` 1.59.1 plus
`~/.cache/ms-playwright/{chromium-1217, chromium_headless_shell-1217,
ffmpeg-1011}`. The spec now runs. Anyone re-reading F5 should treat it as
superseded.

**The pre-existing `origin/feat/router-vm-lab` spec was used as a structural
reference only.** Its actual red reason is different and weaker: it asserted
`#devsplit` == `'10'` (env-pin F7 drift — HEAD serves `0`), which would be a
fixture-value mismatch, not a missing-contract failure. The new spec does not
assert a hardcoded default; the Advanced VISITED case drives `10` explicitly
through the slider and asserts the flow *reports* it, so a legitimate default
change cannot masquerade as a contract break.

**One fixture is deliberately assumed.** No `#ssid` field exists on the pinned
build at `PROGRESS-t_wz18-env-pin.md` F7 read time, and the fixture SSID is
only used inside a doubled `/api/wifi-scan` response for the STA branch. The
WAN path — the default — does not touch it. Flagged in case the green task
changes the uplink UX.

**Scope respected.** No GL-MT6000 hardware run, no flash, no video capture.
That is the separate downstream WZ1.8 task. Pushing the task branch is correct
per the brief (the evidence artifact is what is being handed on); the green
task should continue on this same branch rather than opening a new one.

**One shared-repo trap worth recording.** The git dir is shared with the
primary tree at `/home/c08r4d0r/wz16-verify`, whose `info/exclude` has
`/PROGRESS-t_*.md`. So `git add PROGRESS-t_gate8-red.md` silently no-ops; it
needs `-f`. `PROGRESS.md` at the worktree root is *tracked* and is owned by
whichever card wrote last — it already carries three other cards' sections, so
I did not touch it and used a task-scoped file instead.

## 8. REMAINING steps

None. The task is complete at the captured, referenced red run, as specified
("Red comes first; this task ends at a captured, referenced red run").

Handoff to the follow-on green task:

1. Make the wizard satisfy the contract in `wizard-flow.mjs`: 7 step
   containers with `data-testid="step-<id>"`, a per-step `Next`/`Continue`
   control, a `Skip` control on Advanced, per-step gating, and a terminal
   `Live` step exposing `routerReachable`, `tollgateApiServing` and
   `captivePortalServing` evidence nodes.
2. Re-run: `../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs` → expect exit 0.
3. Then the separate WZ1.8 hardware/video gate (GL-MT6000 run, video evidence)
   owns the device work and must consume the device lock first.
