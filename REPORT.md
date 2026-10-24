# REPORT — t_0841f821: Announce GL-MT6000 lock on shared channel (fd4/fe2b)

Card:      fleet offload worker, repo `net4sats-mvp-v2`
Worktree:  `/home/c08r4d0r/repos/net4sats-mvp-v2`
Branch:    `worker-heavy/t_gate8-19b1b1c6` (remote `felixfelix-bot/net4sats-wizard-go`)
Host:      x280 / CobradorWave (`c08r4d0r`), Ubuntu 24.04.3 LTS, UTC run 2026-10-24
Scope honoured: single announcement. No lock edit, no re-claim, no device contact.

## VERDICT

**ANNOUNCED — but not on a "shared lock channel", because none exists.** Step 1 as
written is unsatisfiable, and I say so rather than fabricate a channel. The announcement
was published on the medium the recovered lock convention itself names as the working
one (this repo's `docs/`, **PUBLIC** remote, anonymously readable), it carries all five
required fields verbatim including the literal `t_557bb77d`, and its permalink was
verified 200 with **no auth**.

### Deliverable (as asked: channel identifier + message reference/permalink)

| Item | Value |
|---|---|
| **Channel identifier** | `docs/` in git repo **`felixfelix-bot/net4sats-wizard-go`** — visibility **PUBLIC** (verified via `gh repo view --json visibility`) |
| **Channel path** | `docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` |
| **Message reference** | commit `ab333063dcde01aa9288bf760b855c1fc0b8f78c` on `refs/heads/worker-heavy/t_gate8-19b1b1c6` |
| **Permalink** | `https://github.com/felixfelix-bot/net4sats-wizard-go/blob/ab333063dcde01aa9288bf760b855c1fc0b8f78c/docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` |
| **Fetch evidence** | anonymous `curl -L` → **http=200**; anonymous `raw.githubusercontent.com` same SHA → **http=200** |
| **Absent alternative** | kind-9 NIP-29 fleet bus on `wss://relay.orangesync.tech` — **has no lock tag**, and the relay itself returns **HTTP 502**. Not used. |

### Field-completeness confirmation (all five, with `t_557bb77d`)

| # | Field | Literal posted | Occurrences |
|---|---|---|---|
| 1 | owner task id | `t_be3261d3` | 8 |
| 2 | device id/serial | `GL-MT6000` (+ `glinet,gl-mt6000`) | 6 (+3) |
| 3 | claim timestamp (UTC ISO-8601) | `2026-10-12T07:33:17Z` | 3 |
| 4 | expected release time | `2026-10-14T07:33:17Z` | 3 |
| 5 | release condition ↔ `t_557bb77d` | `released only after run task t_557bb77d child 3 (Playwright + video run that consumes the device) reports completion, or on explicit abort.` | 4 |

`secret-scan: clean`. `FULL_REVIEW: t_0841f821 announcement — channel identification,
five-field completeness and public reachability each verified first-hand.`

---

## 1. Why there is no shared channel to post to

The recon convention was recovered first-hand from
`verify/t_6b512dad-lock-convention/recon-quote.txt` (a raw multi-source extraction, read
directly rather than trusting the summary). It answers the dispatch question explicitly:

> **A shared channel announcement?** **NO.** A dedicated audit concluded **FAIL**: no
> authoritative shared lock channel exists on any reachable medium (public Nostr board
> mirror has columns backlog/inprogress/humanreview/blocked/done — **no lock column**;
> the kind-9 fleet bus on `wss://relay.orangesync.tech` has tags
> `fleet-task/fleet-claim/fleet-running/fleet-done/fleet-reply-*` — **no lock tag**).
> The card meant to create one, `t_0841f821`, remains **`triage`, never ran**.
> — `recon-quote.txt:107-119`

The same file, twice, names the substitute medium:

> **Sibling-readable announcement medium that actually works: this repo's `docs/`.** Every
> lock claim/renewal/verdict in this batch is a committed `docs/` note.
> — `recon-quote.txt:76-77`, repeated at `:231-232`

Two further first-hand checks this session, so the audit is not being trusted on faith:

```
$ timeout 15 curl -s -o /dev/null -w 'relay_https=%{http_code}\n' https://relay.orangesync.tech/
relay_https=502
$ ls docs/ | grep -i lock
(no lock note — docs/ carries no lock-channel artifact)
```

**Note the circularity that step 1 walks into:** the only card that would have created a
shared lock channel *is this card* (`t_0841f821`), and it never ran. A step instructing me
to "identify the shared lock channel that fd4 and fe2b read" therefore has no
identification to make. I resolved it by using the medium the convention documents as the
one that works, on a **public** remote so peer lanes genuinely can read it — and by
recording the gap instead of papering over it.

## 2. What was actually posted

One artifact, committed and pushed once:

`docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` (105 lines, 6,882 B)

It contains, explicitly labeled: the owner task id, the device id/serial with the
authoritative board string and the serial's documented absence, the claim timestamp, the
expected release time (flagged as a 48 h estimate with the real release event-driven), and
the release condition naming `t_557bb77d` verbatim — plus the current lock state measured
read-only from the local store so a sibling can see the live situation without opening it.

## 3. Provenance of every value — nothing invented

| Field | Source, read first-hand this session |
|---|---|
| `t_be3261d3` | `tasks.id` in the board DB; card bodies of t_0841f821 and t_be3261d3. |
| `GL-MT6000` / `glinet,gl-mt6000` | `recon-quote.txt:85` — §0(d) device identity table. |
| `2026-10-12T07:33:17Z` | `t_d0c3e172` comment id `29831` ("Claim timestamp (UTC ISO-8601): 2026-10-12T07:33:17Z"), re-asserted in comment `29832` (13-key parse, 34/34 checks PASS). |
| `2026-10-14T07:33:17Z` | Same comment — recorded there as a **48 h estimate**, with release event-driven. |
| release condition + `t_557bb77d` | Comments `29831`/`29832`; independently corroborated by the **live local lock file's** own `release_condition=` line (read this session). |

## 4. Two honest qualifiers, stated rather than smoothed over

**(a) The serial does not exist.** Recon §0(d) records serial/case-label as
**"UNRECORDED — a documented absence"**: SSH has never been achieved from any lane, so the
lock record's "Device id/serial" field holds the **model string** in place of a serial. The
announcement reports `GL-MT6000` and labels the absence; it does not invent a serial.

**(b) The timestamps come from the peer node, not from x280.** The sibling card
`t_be3261d3` on this host `BLOCKED` with **no claim made**, and its deliverable records
`claim_timestamp: null` and `expected_release_time: null`
(`verify/t_be3261d3-lock-claim/blockers.json`). The timestamps published here are the ones
issued by claim card `t_d0c3e172` acting **on behalf of** t_be3261d3 on the peer node
`c03rad0r`, whose store is not reachable from x280 (`ssh 192.168.2.32` → *No route to
host*). Both facts sit side by side in the announcement so a reader cannot mistake one for
the other.

## 5. Live lock state (read-only, measured — recorded because it matters to the reader)

```
$ ls -la /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock
-rw------- 1 c08r4d0r c08r4d0r 721 Oct 14 07:38 …gl-mt6000-192.168.1.1.lock
$ sha256sum …gl-mt6000-192.168.1.1.lock
ed8a4573e88605c6fa4c72dcb9c85da706dc9efe18ef6371ed7a65fd90f46582
$ grep -c '^released=' …gl-mt6000-192.168.1.1.lock
0        # rc=1 ⇒ STILL HELD (the documented release test)
```

`owner_task=net4sats-mvp-v2:t_2d0c6fa8`, `claim_ts=2026-10-14T05:38:34Z` — a **different**
owner from t_be3261d3. This announcement is bookkeeping only and does not alter that state.
Flagging it because a sibling reading "t_be3261d3 owns the lock" could otherwise conclude
the local store is free, and it is not.

## 6. Compliance with step 4 (no edit, no re-claim, a single announcement)

- No lock file created, appended to, modified, or deleted.
- No `released=` line appended — that would be the false release of another owner's claim.
- No ssh/scp/ping/curl/nc/adb/fastboot against GL-MT6000 or any device; nothing flashed,
  rebooted, logged into, or power-cycled. Every command was a local read (`ls`, `cat`,
  `sha256sum`, `grep`, `curl` to public hosts, `sqlite3` SELECT, read-only `git`).
- Exactly **one** announcement artifact; no duplicate copy, no channel spam.
- No commit/push claimed that was not observed: the SHA above is the one `git ls-remote`
  returned for the branch and the one the anonymous permalink fetch resolved to 200.

## 7. Residual gap handed to the manager

The lock convention's own audit says a shared lock channel should exist and does not.
`t_0841f821` was the card meant to create one; it has now run and produced an announcement
on the documented fallback medium — **not** a channel. Creating an actual shared lock
channel (a lock tag on the kind-9 bus, or a lock column on the board mirror, on a relay
that is up) remains open and is an architecture decision, not a bookkeeping one.

---



Card:      fleet offload worker, repo `net4sats-mvp-v2`
Worktree:  `/home/c08r4d0r/repos/net4sats-mvp-v2`
Branch:    `worker-heavy/t_gate8-19b1b1c6` (remote `felixfelix-bot/net4sats-wizard-go`)
Host:      x280 / CobradorWave (`c08r4d0r`), Ubuntu 24.04.3 LTS, UTC run 2026-10-24
Scope honoured: written record only. No spec written or modified, no hardware
touched, no GL-MT6000 capture, no repo file modified except this record.
Deliverable files: `PROGRESS.md` (appended section), `REPORT.md` (this file).

## VERDICT

The environment is **fully pinned and reproducible up to the browser boundary**.
Every fact a later worker needs is resolved and verified except one, which is a
hard environmental gap on this host:

- **The wizard app-under-test boots and serves — VERIFIED first-hand.**
- **The Playwright spec cannot execute here — VERIFIED cause: no browser binary
  and no package manager on this host.** Not a regression; pre-existing.

The 7-step generic spec also has one assertion that does not match this commit
(dev-split default), recorded as a fixture value to pin.

---

## 1. Upstream build/commit the run is validated against

| Item | Value (measured) |
|---|---|
| `git rev-parse HEAD` | `b56fae5e3f0eb1850153b43792825b1b8f5184b3` |
| Branch | `worker-heavy/t_gate8-19b1b1c6` |
| Remote | `https://github.com/felixfelix-bot/net4sats-wizard-go.git` |
| HEAD commit | `docs(demo): t_e07e91b1 — full report for the assembled draft master` |
| App under test | Go binary built from this worktree (`go build -o <out> .`) |

Note: `docs(card)` precedent aside, no upstream release tag is pinned by the
E2E — the spec drives whatever binary is serving `WIZARD_URL`. The *commit*
above is therefore the pin, not a release tag.

### How the app is built and served for Playwright

Neither a dev server nor a preview server — the app is a **pre-built Go binary
that serves an embedded SPA**. Verified:

- Build: `cd /home/c08r4d0r/repos/net4sats-mvp-v2 && /usr/local/go/bin/go build -o /tmp/wizard-pin-test .`
  → rc=0, 14,264,590 bytes ELF x86-64.
- **`go` is NOT on the default PATH** for a non-login shell (`go: command not
  found`). It must be invoked as `/usr/local/go/bin/go` or with
  `export PATH=/usr/local/go/bin:$PATH`. **Go version: go1.22.12 linux/amd64.**
- Serve: the binary listens on `127.0.0.1:8099` — hardcoded default
  `defaultListenAddr = "127.0.0.1:8099"` at `main.go:48`, overridable via
  `WIZARD_BIND` (e.g. `WIZARD_BIND=0.0.0.0:8099`).
- A prebuilt binary already sits in the worktree: `./net4sats-wizard`
  (14264590 bytes, mtime Oct 24 03:35) — same size as a fresh build.
- Repo also carries `.hermes/environment.json`: recipe `Go project`,
  `start: "go run ."`, `port: 8099`, `readinessPath: "/"`.

### Exact start command (verified live)

```bash
cd /home/c08r4d0r/repos/net4sats-mvp-v2
pkill -f net4sats-wizard 2>/dev/null            # kill strays FIRST (port-8099 trap)
ss -ltn | grep 8099 && echo "PORT HELD — resolve before proceeding"
/tmp/wizard-pin-test > /tmp/wizard-boot.log 2>&1 &   # or: ./net4sats-wizard
sleep 3
ss -ltnp | grep 8099        # must show our PID owning 127.0.0.1:8099
curl -s -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8099/   # expect 200
```

Observed boot log (`/tmp/wizard-boot.log`, verbatim):

```
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).
```

Observed readiness: `LISTEN 127.0.0.1:8099 users:(("wizard-pin-test",pid=3405276))`,
`HTTP status=200 bytes=19982`.

---

## 2. Where the Playwright suite actually lives (NOT in this repo)

This repo has **no** `package.json`, no Playwright config, and no spec. The suite
is a separate repo:

| Item | Value |
|---|---|
| Kit repo | `/home/c08r4d0r/repos/physical-router-test-automation` |
| Kit HEAD | `5344e3545f70cab74cbb23716dcc428a71f16381` on `fork-main-sync` (2026-10-23) |
| Kit package.json | `"@playwright/test": "^1.52.0"` (devDependency) |
| Kit scripts | `test` / `test:list` / `test:ui` → `npx playwright test --config=tests/playwright.config.mjs` |
| Existing endo spec | `tests/browser/endo-onboarding.spec.mjs` (+ `.config.mjs`) — present on the checked-out branch |
| Generic 7-step spec | `tests/browser/generic-wizard-7step.spec.mjs` (+ `.config.mjs`) — **ONLY on `origin/feat/router-vm-lab`, tip `ee289bf8fb9ee8c54e120ea8b205c6b659b0ca83` (2026-09-27). NOT on main / fork-main-sync.** |

This confirms the earlier lane's finding (`REPORT-t_75852fce.md` §2): the 7-step
spec exists on `feat/router-vm-lab` only, and is unpushed to main.

---

## 3. Required env vars

### 3a. Existing endo spec (`endo-onboarding.spec.mjs`) — the one to run first

Explicit (`process.env`):

| Var | Default in spec | Feeds |
|---|---|---|
| `WIZARD_URL` | `http://localhost:8099` | all navigation |
| `ROUTER_IP` | `192.168.1.1` | router select + step 8 API/portal checks |
| `ROUTER_PASSWORD` | `''` (empty) | `#password`; note spec falls back to `'default-pass'` if falsy |
| `TOLLGATE_VIEWPORT` | `desktop` | config viewport (`desktop` 1280x900 / `mobile` 375x812) |

**Implicit env the endo spec relies on but does NOT read** — these bite silently:

- `channel: 'chrome'` is hardcoded in `endo-onboarding.config.mjs:31`. It is an
  *implicit* requirement on a system Chrome at `/opt/google/chrome/chrome`.
  There is no env override for it. See §5.
- `endo-onboarding.spec.mjs` has **no `.env` loading** — unlike
  `tests/playwright.config.mjs`, which hand-parses `<repo>/.env`. The endo config
  does not, so `.env` values will NOT reach the endo spec. Vars must be exported.
- `ROUTER_PASSWORD || 'default-pass'` fabrication at two call sites: on a fresh
  reset router the correct value is empty, so leaving the var unset injects the
  literal string `default-pass` into `#password`. Pin `ROUTER_PASSWORD=` (empty)
  deliberately rather than omitting it.
- Hardcoded `'endo@coinos.io'` as the Lightning address (no env var) — the spec
  is NOT fully parameterized; that is what the 7-step rewrite fixes.

### 3b. Generic 7-step spec (for the later worker)

All optional, defaults in-source. Full inventory (from the spec's own header and
`process.env` sweep):

| Var | Default | Purpose |
|---|---|---|
| `WIZARD_URL` | `http://localhost:8099` | wizard base |
| `ROUTER_IP` | `192.168.1.1` | must appear in scan results |
| `ROUTER_PASSWORD` | `''` | fresh-reset routers have none |
| `LNURL` | `generic@wallet.app` | payouts Lightning address |
| `DEPLOY_MODE` | `wan` | `wan` (ethernet) or `sta` (repeater) |
| `WIFI_SSID` | `''` | STA only |
| `WIFI_PASSWORD` | `''` | STA only |
| `DEV_SPLIT` | unset → leave wizard default | advanced; asserts `#devsplit-val` text. **See §4 drift.** |
| `MARGIN` | unset → leave wizard default | advanced |
| `STRICT_DEPLOY` | `1` | `1` fails on wizard error view; `0` lenient like endo |
| `DEPLOY_TIMEOUT` | `180000` | ms to wait for deploy result |
| `TOLLGATE_PORT` | `2121` | step 7 API check |
| `PORTAL_PORT` | `2050` | step 7 portal check |

There is **no feature flag** toggling generic-vs-endo flow — they are two
different spec files with two different configs. Selection is by `--config`.

There is **no DB/reset hook** in either spec: state reset is done out-of-band by
re-pointing the wizard at a router (the wizard itself has no DB). The only
reset-adjacent behaviour is `ssid` re-scan and the router's own state.

### 3c. Example invocation (values that match this host)

```bash
# Precondition: wizard serving on 8099, no router required only for DOM contract
export WIZARD_URL=http://localhost:8099
export ROUTER_IP=192.168.1.1
export ROUTER_PASSWORD=            # deliberately empty (fresh reset)
export TOLLGATE_VIEWPORT=desktop
```

---

## 4. Fixtures / setup dependencies

Neither config declares Playwright `projects` with `dependencies` (setup
projects) beyond a single viewport project; there is **no globalSetup / fixture
file**. Fixtures are environmental, not code:

1. **Wizard process** serving `WIZARD_URL`, owning `:8099`. Must be proven owned
   by PID before browser traffic (zombie-wizard trap — a stale server on 8099
   with a replaced binary causes failures that masquerade as router auth errors).
2. **Target router** at `ROUTER_IP` reachable by SSH with `ROUTER_PASSWORD`.
   Serial mode (`test.describe.configure({ mode: 'serial' })`) in both specs, so
   router state persists across the 5 / 7 tests.
3. **Post-deploy ports** on the router: TollGate API `:2121`, portal `:2050`
   (7-step defaults) or `:2050` (endo step 8).
4. **Viewport project** `desktop-endo-onboarding` /
   `desktop-generic-wizard-7step` — selected by `TOLLGATE_VIEWPORT`.

### Fixture drift found (pin this — do not let it fail silently)

The 7-step spec's **Step 5** asserts:

```js
await expect(page.locator('#devsplit')).toHaveValue('10');
await expect(page.locator('#margin')).toHaveValue('0');
```

Measured against HEAD's served HTML:

- `#margin` → `value="0"` ✅ matches.
- `#devsplit` → `value="0"`, `#devsplit-val` = `0%` (index.html:190-191) ❌
  spec expects **10**.

The wizard was changed from default 10 → 0 (the net4sats WZ1.4 dev-split
question; `schedule.md` records the Q3 gate RESOLVED = **B**, upstream keeps its
own default and net4sats may differ downstream). **The later 7-step worker must
decide whether to pin the expected value to `0` or to set `DEV_SPLIT=10`** — as
written, Step 5 fails against this commit for a value reason, independent of any
browser problem.

---

## 5. VERIFIED STATUS OF THE EXISTING ENDO SPEC — run attempted, result observed

I ran it rather than reasoning about it.

**Run command used (exact):**

```bash
cd /home/c08r4d0r/repos/physical-router-test-automation
node node_modules/playwright/cli.js test \
  --config=tests/browser/endo-onboarding.config.mjs \
  --project='desktop-endo-onboarding' --grep 'Step 1-2'
```

**Result: FAILED at browser launch — 1 failed (6ms), 0 passed.**

```
Error: browserType.launch: Chromium distribution 'chrome' is not found at /opt/google/chrome/chrome
Run "npx playwright install chrome"
```

**Classification: PRE-EXISTING ENVIRONMENTAL FAILURE — NOT a spec regression.**
The spec never reached a single assertion. Cause verified independently:

| Probe | Result |
|---|---|
| `node --version` | `v18.19.1` present at `/usr/bin/node` |
| `npm` | **absent** (`npm: command not found`) |
| `npx` | **absent** (`npx: command not found`) |
| system chrome/chromium at the 5 probed paths | **none exist** (`/usr/bin/google-chrome`, `-stable`, `chromium`, `chromium-browser`, `/snap/bin/chromium`) |
| `find / -maxdepth 8 -name chrome -o -name chromium -o -name headless_shell` | only `/etc/apparmor.d/chrome` (a profile, not a binary) |
| `~/.cache/ms-playwright` | **absent** — no bundled browser |
| kit `node_modules` | absent in the kit; a reusable `@playwright/test` **1.59.1** exists at `/data/worktrees/prta/node_modules` |

I made the spec *reachable* to prove the failure is the browser, not module
resolution: symlinked that `node_modules` into the kit, then test discovery
worked and `--list` enumerated all 5 endo tests. The run then failed at launch.
(The symlink and the `test-results/`/`endo-report/` output my attempt created were
**removed afterwards**; the kit is left exactly as found — `git status` shows only
the pre-existing untracked `.venv/`.)

**Consequence:** the endo spec **cannot be executed on this host as-is**. Any
worker must first provision a browser. `playwright install chrome` needs npm and
network; the better fix on this host is a system Chromium plus overriding the
launch path — but note `endo-onboarding.config.mjs` hardcodes
`channel: 'chrome'` with **no env override**, unlike `tests/playwright.config.mjs`
which probes `PLAYWRIGHT_CHROME` / `CHROME_PATH` and 5 system paths. Using the
kit's *root* config is the only browser-overridable entry point today.

---

## 6. Run command + reporter config + grep/filter syntax

### Reporter config (measured)

`endo-onboarding.config.mjs`:

```js
reporter: [['html', { outputFolder: 'endo-report', open: 'never' }], ['list']]
timeout: 180000   retries: 0   workers: 1
use: { baseURL: 'http://localhost:8099', headless: true, channel: 'chrome',
       screenshot: 'on', trace: 'on', video: 'on', actionTimeout: 30000 }
```

`generic-wizard-7step.config.mjs` is the same shape with
`timeout: 240000`, `outputFolder: 'generic-wizard-report'`,
`testMatch: 'generic-wizard-7step.spec.mjs'`.

Kit root `tests/playwright.config.mjs` additionally emits
`['json', { outputFile: 'report/report.json' }]` and is the only config with the
`chromiumLaunchOptions()` browser-path override and `.env` loading.

### Exact run commands

```bash
cd /home/c08r4d0r/repos/physical-router-test-automation/tests/browser

# canonical (requires npm/npx — ABSENT on this host)
npx playwright test --config=endo-onboarding.config.mjs
npx playwright test --config=generic-wizard-7step.config.mjs

# the resolution-free form that works with a borrowed playwright install
node <path>/node_modules/playwright/cli.js test --config=<config>.mjs
```

### Grep / filter syntax — VERIFIED working

Playwright `--grep` filters by test **title**. Proven on the endo spec:

```bash
# list only (no browser needed) — VERIFIED output: 2 tests matched
node node_modules/playwright/cli.js test --config=tests/browser/endo-onboarding.config.mjs \
  --list --grep 'Deploy'

# exact endo titles (measured via --list, line numbers included)
#   Step 1-2: Auto-detect routers on network
#   Step 3-6: Configure deployment (select router, password, upstream, LNURL)
#   Step 7: Deploy and watch deployment progress
#   Step 8: Verify captive portal and TollGate API on router
#   Full experience screenshot summary
```

**7-step spec — exact test titles** (for the Advanced cases the task calls out):

```
Step 1: Discover — auto-detect routers on the network
Step 2: Connect — select router and enter password
Step 3: Uplink — choose upstream mode
Step 4: Payouts — enter Lightning address
Step 5: Advanced (optional) — open advanced settings, defaults intact
Step 6: Deploy — run the deployment to completion
Step 7: Live — verify TollGate API and portal on the router
```

```bash
# Advanced-skip case        (Step 5 not visited)
--grep 'Step 5' --grep-invert
# or run the whole file and exclude:
--grep-invert 'Advanced'

# Advanced-visited case     (Step 5 only)
--grep 'Step 5: Advanced'

# a single title verbatim
--grep 'Step 5: Advanced \(optional\)'
```

`--grep-invert` is the correct switch for the "Advanced-skip" case; `--grep`
alone for "Advanced-visited". Note `--grep` uses a regex, so escape `(`/`)`.

---

## 7. Selector / DOM contract (both specs) — VERIFIED against HEAD

All ids the endo spec uses exist in HEAD's served HTML (measured counts = 1 each):
`#select-view`, `#router-select`, `#password`, `#mode-wan`, `#sta-fields`,
`#lnurl`, `#deploy-btn`, `#deploy-view`, `#success-view`, `#error-view`,
`#error-detail`, `#steps-list`.

The 7-step spec's additional gates also exist: `#scan-view` (index.html:141),
`details.advanced` (index.html:187), `#devsplit` / `#devsplit-val` (:190-191),
`#margin` / `#margin-val` (:196-197), `#mint` (:203), `#wifi-pass` (:178),
`#ssid` (:169).

Flow shape at HEAD: `scan-view` → `select-view` → `deploy-view` → `success-view`,
with an `error-view`. The 7-step spec's `discoverAndWait()` waits on `#scan-view`
then `#select-view:not(.hidden)` — **satisfied**. (Correction recorded: an
intermediate note in `PROGRESS.md` claiming HEAD had "no 7-step flow" was wrong
and is corrected there; the selectors are all present. The only real mismatch is
the `#devsplit` **value** in §4.)

---

## 8. Pre-existing failures vs regressions

- **Pre-existing:** no browser / no npm / no npx on this host → every
  browser-launching spec fails at launch. Independent of any spec content.
- **Pre-existing:** the 7-step spec is unpushed (only on
  `origin/feat/router-vm-lab`) and its Step 5 dev-split expectation predates the
  10 → 0 wizard change.
- **Regressions found: none.** Nothing I observed was caused by the current
  commit; the DOM contract at HEAD is fully intact.

## 9. What a later worker still must do (handoff)

1. Provision a browser on the target host (npm + `playwright install chrome`, or
   a system Chromium plus a config that honours `PLAYWRIGHT_CHROME`/`CHROME_PATH`
   — `endo-onboarding.config.mjs` hardcodes `channel: 'chrome'` today).
2. Decide the `#devsplit` expectation for Step 5 (`0` vs `DEV_SPLIT=10`).
3. Get `generic-wizard-7step.spec.mjs` + `.config.mjs` onto the branch under test
   (currently only on `origin/feat/router-vm-lab`).
4. Provide the hardware fixture (`ROUTER_IP` reachable, SSH password) for steps 6
   and 7 — out of this task's scope (no hardware contact performed).

## 10. Process notes

- The plan ref named in the brief
  (`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`, "entry WZ1.8") has
  **no WZ1.8 entry** — `grep WZ1.8` → 0 hits; it is a FELIX-Q3-gate document.
  Same precondition gap prior lanes recorded (`t_75852fce`, `t_a8b5fcae`,
  `t_cc76c746`).
- The untracked `PROGRESS-t_*.md` / `REPORT-t_*.md` naming used by sibling cards
  is blocked for the `t_wz18-*` prefix here: a **different worktree's**
  `info/exclude` (`/home/c08r4d0r/wz16-verify/.git/info/exclude:4`) carries
  `/PROGRESS-t_*.md` and is applied to this shared linked worktree. Hence this
  record is written into the tracked, non-ignored `PROGRESS.md` / `REPORT.md`.
