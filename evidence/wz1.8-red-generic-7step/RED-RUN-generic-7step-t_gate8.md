# RED RUN — generic 7-step wizard spec vs pinned upstream build

Status: **RED (expected).** This is the red-first artifact for the generic 7-step
wizard flow. It is deliberately NOT green: the follow-on task turns it green.

## Exact command

```
cd ~/repos/physical-router-test-automation/tests/browser
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs
```

Exit code: **1** (5 failed / 5 total). Reproducible via
`reproduce-red-run.sh` in this directory.

## Environment / pin

| Item | Value |
|---|---|
| Repo under test | `felixfelix-bot/net4sats-wizard-go` |
| Repo HEAD at capture | `867d45604de243c41406e387ea7cf0a88e78d859` |
| Served `index.html` sha256 | `72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9` |
| Repo `HEAD:index.html` sha256 | `72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9` |
| Server | `/tmp/wz18-wizard`, sha256 `43f0ae6456a2fdc842416526835418123ee71983d6cd2c2f4c5da1478741422e` |
| Server bind | `127.0.0.1:8099` (main.go:48 `defaultListenAddr`) |
| Suite | `~/repos/physical-router-test-automation` |
| Spec | `tests/browser/generic-wizard-7step.spec.mjs` |
| Flow model (data) | `tests/browser/wizard-flow.mjs` |
| Config | `tests/browser/generic-wizard-7step.config.mjs` |
| Playwright | `@playwright/test` 1.59.1 |
| Browser | `channel: 'chromium'` -> `~/.cache/ms-playwright/chromium-1217/chrome-linux64/chrome` |

The pin is **verified by content, not by assumption**: the bytes served on
`127.0.0.1:8099` hash identically to `HEAD:index.html`, so the running server
is the pinned source.

Raw pin capture:

```
=== wizard process ===
    PID                  STARTED COMMAND
3756902 Sat Oct 24 05:04:02 2026 /tmp/wz18-wizard

=== serving bytes == /proc/exe binary ===
lrwxrwxrwx 1 c08r4d0r c08r4d0r 0 Oct 24 05:04 /proc/3756902/exe -> /tmp/wz18-wizard
43f0ae6456a2fdc842416526835418123ee71983d6cd2c2f4c5da1478741422e  /tmp/wz18-wizard

=== served index.html sha (from 8099) ===
72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9  -
=== HEAD index.html sha (git) ===
72d79dd2c4246ed18affded543a6ce64846ce7f2104333db5737216178a416c9  -
=== repo HEAD ===
867d45604de243c41406e387ea7cf0a88e78d859
```

## Expected failure reason (stated explicitly)

**The pinned build implements no step model.** `index.html` renders one flat
`#select-view` form (Discover/Connect/Uplink/Payouts/Advanced controls merged
into a single card) and shows no per-step container, no forward-navigation
control, no per-step gating, and no terminal `Live` step with post-deploy
evidence. Concretely, against the pinned build:

- `grep -c data-testid` on the served document = **0** — the spec's testid
  step hooks (`[data-testid="step-discover"]` … `step-live`) do not exist.
- The served document contains only **2** `<h2>` elements (`net4sats is live!`,
  `Setup failed`); the spec requires the 7 declared step headings.
- There is no `Next`/`Continue` control and no `Skip` control anywhere in the
  served document, so no step can be traversed or gated.

This is a **missing-generic-wizard-behaviour** failure, not a typo, import
error, or environment failure: the spec parses, the browser launches, the
server serves, and every assertion reaches the DOM before failing. Each case
fails on a wizard-contract assertion, not on a fixture/selector syntax error.

## Why the API is stubbed

`/api/scan` scans the real LAN and `/api/deploy` drives a real router. No test
router is reachable from this host (no route to `192.168.1.1`; `/api/scan`
returns no routers), so a hardware-dependent spec would fail for a *network*
reason — a wrong-reason red. The spec therefore intercepts `/api/scan`,
`/api/wifi-scan`, `/api/deploy` and `/api/status/<job>` with deterministic
doubles. The doubled `status` payload is a *completed* deployment, so the only
thing that can keep the flow from reaching `Live` is missing wizard UI — which
is exactly the property under test.

## Raw output (verbatim, unedited)

```

Running 5 tests using 1 worker

  ✘  1 [generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:2 › generic 7-step wizard flow › step order: all 7 steps appear in the declared order (5.9s)
  ✘  2 [generic-wizard-7step] › generic-wizard-7step.spec.mjs:272:2 › generic 7-step wizard flow › advanced SKIP path: skipping Advanced still reaches Deploy (891ms)
  ✘  3 [generic-wizard-7step] › generic-wizard-7step.spec.mjs:303:2 › generic 7-step wizard flow › advanced VISITED path: visiting Advanced applies and reports its values (905ms)
  ✘  4 [generic-wizard-7step] › generic-wizard-7step.spec.mjs:331:2 › generic 7-step wizard flow › gating: cannot advance past a step while its requirements are unmet (947ms)
  ✘  5 [generic-wizard-7step] › generic-wizard-7step.spec.mjs:370:2 › generic 7-step wizard flow › live: terminal Live step reached with post-deploy evidence (906ms)


  1) [generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:2 › generic 7-step wizard flow › step order: all 7 steps appear in the declared order 

    Error: step "Discover" is missing: expected a container with [data-testid="step-discover"]

    [2mexpect([22m[31mlocator[39m[2m).[22mtoBeAttached[2m([22m[2m)[22m failed

    Locator: locator('[data-testid="step-discover"]')
    Expected: attached
    Timeout: 5000ms
    Error: element(s) not found

    Call log:
    [2m  - step "Discover" is missing: expected a container with [data-testid="step-discover"] with timeout 5000ms[22m
    [2m  - waiting for locator('[data-testid="step-discover"]')[22m


      252 | 				stepLocator(page, step),
      253 | 				`step "${step.title}" is missing: expected a container with ${stepSelector(step)}`,
    > 254 | 			).toBeAttached();
          | 			  ^
      255 | 		}
      256 |
      257 | 		// DOM order must match declared order — compare the sequence of testids
        at /home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:254:6

    attachment #1: screenshot (image/png) ──────────────────────────────────────────────────────────
    ../../test-results/generic-wizard-7step-gener-e510c-ppear-in-the-declared-order-generic-wizard-7step/test-failed-1.png
    ────────────────────────────────────────────────────────────────────────────────────────────────

    Error Context: ../../test-results/generic-wizard-7step-gener-e510c-ppear-in-the-declared-order-generic-wizard-7step/error-context.md

  2) [generic-wizard-7step] › generic-wizard-7step.spec.mjs:272:2 › generic 7-step wizard flow › advanced SKIP path: skipping Advanced still reaches Deploy 

    Error: expected exactly one declared step container visible, found 0 (none) — the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live

      151 | 	}
      152 | 	if (visible.length !== 1) {
    > 153 | 		throw new Error(
          | 		      ^
      154 | 			`expected exactly one declared step container visible, found ${visible.length} ` +
      155 | 				`(${visible.join(', ') || 'none'}) — the wizard declares ${STEP_COUNT} steps: ` +
      156 | 				STEP_ORDER.join(' -> '),
        at currentStepId (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:153:9)
        at walkTo (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:210:8)
        at /home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:275:3

    attachment #1: screenshot (image/png) ──────────────────────────────────────────────────────────
    ../../test-results/generic-wizard-7step-gener-4e965-vanced-still-reaches-Deploy-generic-wizard-7step/test-failed-1.png
    ────────────────────────────────────────────────────────────────────────────────────────────────

    Error Context: ../../test-results/generic-wizard-7step-gener-4e965-vanced-still-reaches-Deploy-generic-wizard-7step/error-context.md

  3) [generic-wizard-7step] › generic-wizard-7step.spec.mjs:303:2 › generic 7-step wizard flow › advanced VISITED path: visiting Advanced applies and reports its values 

    Error: expected exactly one declared step container visible, found 0 (none) — the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live

      151 | 	}
      152 | 	if (visible.length !== 1) {
    > 153 | 		throw new Error(
          | 		      ^
      154 | 			`expected exactly one declared step container visible, found ${visible.length} ` +
      155 | 				`(${visible.join(', ') || 'none'}) — the wizard declares ${STEP_COUNT} steps: ` +
      156 | 				STEP_ORDER.join(' -> '),
        at currentStepId (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:153:9)
        at walkTo (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:210:8)
        at /home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:306:3

    attachment #1: screenshot (image/png) ──────────────────────────────────────────────────────────
    ../../test-results/generic-wizard-7step-gener-d47d9-lies-and-reports-its-values-generic-wizard-7step/test-failed-1.png
    ────────────────────────────────────────────────────────────────────────────────────────────────

    Error Context: ../../test-results/generic-wizard-7step-gener-d47d9-lies-and-reports-its-values-generic-wizard-7step/error-context.md

  4) [generic-wizard-7step] › generic-wizard-7step.spec.mjs:331:2 › generic 7-step wizard flow › gating: cannot advance past a step while its requirements are unmet 

    Error: expected exactly one declared step container visible, found 0 (none) — the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live

      151 | 	}
      152 | 	if (visible.length !== 1) {
    > 153 | 		throw new Error(
          | 		      ^
      154 | 			`expected exactly one declared step container visible, found ${visible.length} ` +
      155 | 				`(${visible.join(', ') || 'none'}) — the wizard declares ${STEP_COUNT} steps: ` +
      156 | 				STEP_ORDER.join(' -> '),
        at currentStepId (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:153:9)
        at walkTo (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:210:8)
        at /home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:334:3

    attachment #1: screenshot (image/png) ──────────────────────────────────────────────────────────
    ../../test-results/generic-wizard-7step-gener-44424--its-requirements-are-unmet-generic-wizard-7step/test-failed-1.png
    ────────────────────────────────────────────────────────────────────────────────────────────────

    Error Context: ../../test-results/generic-wizard-7step-gener-44424--its-requirements-are-unmet-generic-wizard-7step/error-context.md

  5) [generic-wizard-7step] › generic-wizard-7step.spec.mjs:370:2 › generic 7-step wizard flow › live: terminal Live step reached with post-deploy evidence 

    Error: expected exactly one declared step container visible, found 0 (none) — the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live

      151 | 	}
      152 | 	if (visible.length !== 1) {
    > 153 | 		throw new Error(
          | 		      ^
      154 | 			`expected exactly one declared step container visible, found ${visible.length} ` +
      155 | 				`(${visible.join(', ') || 'none'}) — the wizard declares ${STEP_COUNT} steps: ` +
      156 | 				STEP_ORDER.join(' -> '),
        at currentStepId (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:153:9)
        at walkTo (/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:210:8)
        at /home/c08r4d0r/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs:371:3

    attachment #1: screenshot (image/png) ──────────────────────────────────────────────────────────
    ../../test-results/generic-wizard-7step-gener-37f31-d-with-post-deploy-evidence-generic-wizard-7step/test-failed-1.png
    ────────────────────────────────────────────────────────────────────────────────────────────────

    Error Context: ../../test-results/generic-wizard-7step-gener-37f31-d-with-post-deploy-evidence-generic-wizard-7step/error-context.md

  5 failed
    [generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:2 › generic 7-step wizard flow › step order: all 7 steps appear in the declared order 
    [generic-wizard-7step] › generic-wizard-7step.spec.mjs:272:2 › generic 7-step wizard flow › advanced SKIP path: skipping Advanced still reaches Deploy 
    [generic-wizard-7step] › generic-wizard-7step.spec.mjs:303:2 › generic 7-step wizard flow › advanced VISITED path: visiting Advanced applies and reports its values 
    [generic-wizard-7step] › generic-wizard-7step.spec.mjs:331:2 › generic 7-step wizard flow › gating: cannot advance past a step while its requirements are unmet 
    [generic-wizard-7step] › generic-wizard-7step.spec.mjs:370:2 › generic 7-step wizard flow › live: terminal Live step reached with post-deploy evidence
```

## The two Advanced cases, individually runnable

Exact title filters:

```
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs --grep "advanced SKIP path"
../../node_modules/.bin/playwright test --config=generic-wizard-7step.config.mjs --grep "advanced VISITED path"
```

Both select exactly 1 test and both fail for the expected reason:

```
$ ... --grep "advanced SKIP path"
Running 1 test using 1 worker
  x  1 ... advanced SKIP path: skipping Advanced still reaches Deploy (881ms)
    Error: expected exactly one declared step container visible, found 0 (none) - the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live
  1 failed

$ ... --grep "advanced VISITED path"
Running 1 test using 1 worker
  x  1 ... advanced VISITED path: visiting Advanced applies and reports its values (832ms)
    Error: expected exactly one declared step container visible, found 0 (none) - the wizard declares 7 steps: discover -> connect -> uplink -> payouts -> advanced -> deploy -> live
  1 failed
```

(ASCII-transcribed from the raw run; the `->` above renders as the em-dash
arrow in the raw output. The full verbatim capture is in the block above.)

## Determinism compliance

- No `waitForTimeout`, no sleeps, no wall-clock arithmetic anywhere in the spec
  or its helpers. Every wait is an explicit UI-state condition:
  `expect(locator).toBeVisible()`, `toBeEnabled()`, `toBeAttached()`,
  `toHaveCount()`, `toHaveText()`, `toHaveAttribute()`.
- Selectors are `getByRole` / `getByTestId` first. The only CSS used is
  `[data-testid="..."]`, which *is* the testid hook, and a single
  `[data-testid^="step-"]` for DOM-order comparison.
- Per-case reset: `beforeEach` installs fresh API doubles and navigates;
  `afterEach` runs `page.unrouteAll()` and returns to `about:blank`. Playwright
  gives each case a fresh context and page, so no case can leak state.
- `retries: 0` and `workers: 1` in the config, so a flake cannot be mistaken
  for the red result.

## Coverage map (brief -> case)

| Brief | Case title (grep-able) |
|---|---|
| a. 7 steps in declared order | `step order: all 7 steps appear in the declared order` |
| b. Advanced SKIP path | `advanced SKIP path: skipping Advanced still reaches Deploy` |
| c. Advanced VISITED path | `advanced VISITED path: visiting Advanced applies and reports its values` |
| d. forward-navigation gating | `gating: cannot advance past a step while its requirements are unmet` |
| e. terminal Live + evidence | `live: terminal Live step reached with post-deploy evidence` |
