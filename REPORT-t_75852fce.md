# REPORT — t_75852fce: Run generic 7-step wizard E2E on GL-MT6000 and record full video

- **Card:** `t_75852fce` (board `net4sats-mvp-v2`, assignee `worker-tollgate`)
- **Brief repo:** `net4sats-mvp-v2` = `felixfelix-bot/net4sats-wizard-go`
- **Host:** `x280` (`c08r4d0r`), UTC run date **2026-10-23** (started 00:32Z)
- **Verdict:** **BLOCKED at the precondition gate — the task never started.**
- **Gate-9 status:** **NO VIDEO. NO GREEN RUN.** Nothing was flashed, no E2E
  was executed, no artifact was produced.
- **Lease id:** none was ever held by this session. The stale peripheral lock
  file that exists is `gl-mt6000-192.168.1.1.lock`, holder `t_2d0c6fa8`
  (claim 2026-10-14T05:38:34Z, pid 2174746 — **pid is dead**, so the file is
  stale). I did **not** take it and did **not** release it.
- **Spec commit SHA:** `ee289bf8fb9ee8c54e120ea8b205c6b659b0ca83`
  (`origin/feat/router-vm-lab`, the only branch carrying the spec)
- **Build hash:** none — no flash was performed.
- **Failure mode:** the two preconditions the brief told me to verify before
  starting are both false. I stopped instead of proceeding, per the brief's
  own instruction ("Do NOT start if fd4 or fe2b hold the lease").

---

## 1. What I was asked to do

Execute the generic 7-step wizard Playwright spec against the real GL-MT6000
running the upstream build, get a GREEN run, and capture a full end-to-end
video as the Gate-9 evidence artifact.

The brief listed exactly two preconditions to verify first:

> - The spec adaptation task is complete, pushed, and its spec file path is known.
> - The environment readiness task is complete and you hold the exclusive
>   GL-MT6000 hardware window. Do NOT start if fd4 or fe2b hold the lease.

Both fail. The gates did their job.

---

## 2. Precondition A — spec adaptation: PARTIALLY met, and NOT "pushed"

The spec is **real and competent**, which is worth stating plainly. It lives in
a different repo from the one the brief names:

```
repo:   physical-router-test-automation   (~/repos/physical-router-test-automation)
spec:   tests/browser/generic-wizard-7step.spec.mjs        (259 lines)
config: tests/browser/generic-wizard-7step.config.mjs
branch: origin/feat/router-vm-lab   tip ee289bf8   dated 2026-09-27
```

It implements precisely the seven steps the card names — Discover / Connect /
Uplink / Payouts / Advanced-optional / Deploy / Live — and its step 5 is a
genuine assertion of the Advanced-optional semantics, not a stub:

```js
await expect(page.locator('details.advanced')).not.toHaveAttribute('open', 'open');
await expect(page.locator('#deploy-btn')).toBeEnabled(); // advanced not needed
```

Step 7 tests the router for real: `http://ROUTER_IP:2121/` must respond with
`<500` and match `/kind|metric|pubkey|price_per_step/i`, then the portal on
`:2050`.

**But "complete, pushed, and its spec file path is known" is only one-third
true.**

| Claim | Measured |
|---|---|
| complete | Yes — the spec is fully written |
| pushed | **No** — exists only on `feat/router-vm-lab`; `grep` for it on `main` and on the fork sync branch returns **0** |
| path known | Yes — `tests/browser/generic-wizard-7step.spec.mjs` |

The board corroborates the incompleteness. `t_6fd69072` ("Adapt
endo-onboarding Playwright spec to generic 7-step wizard flow") is still
`todo`. `t_1eae5a7b` ("Make 7-step wizard spec green, update docs, commit and
push") is still `triage`, and its `gate-tick` comments read:

```
gate-tick: BLOCKED (completion-pending) — missing gate(s): tests_green,
pushed_or_consolidated, ci_evidence, cold_cross_family_review,
review_artifact, consolidated, secrets_clean, no_live_drift.
```

A `fleet-done rc=0 verified=True` comment sits alongside a failing gate-tick —
i.e. a worker self-reported done and the gate refused it. The gate was right.

Two further mismatches with the brief, both material to the run:

1. The brief says run the spec **headed**. The config asks for
   `headless: true`.
2. The config asks for `channel: 'chrome'`, which **cannot launch on this
   host** (see §5).

---

## 3. Precondition B — hardware: the MT6000 is real, and it is BUSY

I did not take a "not reachable" answer at face value — a probe negative is
only evidence if the probe is validated. I therefore (a) validated a known
positive, (b) swept both subnets, and (c) found the device by reading the
device's own prompt over serial rather than trusting a port or a host.

### 3.1 The device exists and I located it

```
port:      /dev/ttyUSB0
USB serial: FTDI_FT232R_USB_UART_BG02QAPG   (stable — never the tty index)
model:     GL.iNet GL-MT6000                (/proc/device-tree/model)
hostname:  tollgate-0195
firmware:  OpenWrt 24.10.2 r28739-d9340319c6
opkg:      opkg version 38eccbb1fd694d4798ac1baf88f9ba83d1eac616 (2024-10-16)
uptime:    2 days
```

This matches the board's own record: `physical-router-testing` ->
`references/spilman-channel-bench-lane.md` documents the MT6000 as a **serial-
attached, serial-only** fleet router on x280.

### 3.2 It is mid-run in ANOTHER context window's bench

The live shell prompt on that console is the evidence:

```
root@tollgate-0195:/tmp/spilman-lane-3613961#
```

and that directory contains a live Spilman-channel bench, all files created
2026-10-22 21:38–21:41:

```
-rw-r--r-- config.json       161 B
-rw------- identities.json   213 B   (root-only)
-rw-r--r-- mod.log          1862 B
-rw-r--r-- mod.pid             6 B
```

The kit documents this lane as one that "deploys its own artifacts ... on a
private port/state dir (`:2121` and `:2122` are other context windows' —
hard-abort on them)". `:2121` is exactly the port the 7-step spec's Live step
asserts against. Running now would (i) trash another CW's in-flight
measurement, and (ii) produce a meaningless result, because the box is not in
the clean pre-wizard state the brief requires.

`mod.pid` being 6 bytes means a live process is running under that lane on
this very console.

### 3.3 Recording a real lock-coverage gap

The bench mutexes are all FREE:

```
$ who-holds-tty.sh /dev/ttyUSB0     -> FREE  /dev/ttyUSB0
$ flock -n ~/.local/state/tg-bench/tollgate-A04E.lock -c 'echo FREE'  -> FREE
```

So the machinery that is *supposed* to serialize hardware access reports the
MT6000 as available **while another lane is actively using it**. The spilman
lane holds the UART, not the lease; the lease does not model the UART holder.
Anyone trusting `flock` here would walk straight into the collision I avoided.
Also note the peripheral lock file that *does* exist names
`device_addr=192.168.1.1` — an address with no path to this box at all, so it
would not have protected the console either way.

This is reported as a finding, not fixed: fixing a shared mutex used by other
lanes is out of scope for this card and is exactly the kind of change that
needs its own review.

### 3.4 There is no network path to the MT6000 (probe validated)

I validated the probe before believing its negatives. Positive control:
`192.168.2.1:80` -> OPEN, so the method works.

| Target | Result |
|---|---|
| `enp0s31f6` (kit's wired bench NIC) | admin-DOWN, and `carrier=0` after `sudo -n ip link set up` — **no cable** |
| `192.168.2.0/24` ports 22/2121/2051/8090 | only this host's own interfaces; **zero** TollGate hosts |
| `192.168.2.1` | **stock GL.iNet, unrelated firmware** — `buildInfo.js` says `project: HH63`, `version 2585f2c`, `buildTime 2023-09-25`; 2121/2051/8090/8080 all `000` |
| `10.107.73.0/24` (wifi uplink) 22/2121/2051 | 3 SSH-only hosts, none OpenWrt, none TollGate |
| `192.168.1.1` (the spec's default `ROUTER_IP`) | unroutable — routed out the default gateway and dropped; tcp 22/80/8080 closed |
| `192.168.13.231` (kit `routers.lab-backup.json`) | unreachable |

The kit's own read-only gate reaches the same conclusion independently:

```
$ bash scripts/hw-readonly-check.sh --host 192.168.1.1 --no-sweep
  G1 tollgate API full-mode : FAIL :2121/ answered 000 — is the tollgate service up?
  G2 session state          : FAIL /balance answered 000 — treating bench as NOT idle
  gates_failed : 3
  RESULT: FAIL
```

So even setting the spilman lane aside, there is no way to reach this device
over IP. Any run must go through the serial console — the same console the
other lane is holding.

---

## 4. Serial map of this host (each console identified by its own prompt)

| Port | FTDI serial | Device | Evidence |
|---|---|---|---|
| `ttyUSB0` | `BG02QAPG` | **GL-MT6000 / tollgate-0195** | `GL.iNet GL-MT6000`, prompt in `/tmp/spilman-lane-3613961` — **BUSY** |
| `ttyUSB1` | `BG046H14` | COMFAST box | prompt `COMFAST login:` |
| `ttyUSB2` | `BH00D6DL` | D-Link COVR-X1860 A1 | `D-Link COVR-X1860 A1`, OpenWrt 25.12.5 |
| `ttyUSB3` | `BH00DL3M` | silent | 10 s passive listen + newline: 0 bytes both times |

I probed read-only (one bare newline to force a prompt refresh; no command was
sent until the device was identified, and even then only `cat`/`grep`/`ls`
reads). `ttyACM0` is an unrelated FIBOCOM LTE modem.

---

## 5. Harness half: VIABLE, but the config as-written cannot launch

The brief's implied "just run it" is wrong twice over — and I checked this
properly rather than assuming, because a working harness with no device and a
broken harness with a device are different reports.

- The kit repo has **no `node_modules`**; `npm`, `npx`, `yarn`, `pnpm`,
  `corepack` are all **absent** from PATH. Node is `18.19.1`.
- Playwright is nevertheless available from an npx cache:
  `playwright 1.52.0` at
  `~/.npm/_npx/c61c9351a0dbcfa7/node_modules/.bin/playwright`.
  (The other cache holds `^1.64.1`, which refuses to run: "Playwright requires
  Node.js 20 or higher".)
- The config's `channel: 'chrome'` **fails**:
  `Chromium distribution 'chrome' is not found at /opt/google/chrome/chrome`.
- The bundled revision also fails: `chromium_headless_shell-1169` is absent.
- **But launching by explicit `executablePath` to the Chromium that IS
  installed works**, and it records video:

```
EXEPATH LAUNCH OK — page text = "harness OK via executablePath"

$ ffprobe .../vid-exepath/81af02314f2cc1040f63a8207fd9037d.webm
  codec_name=vp8   width=1280   height=900   duration=3.040000
```

So the harness gap is a **one-line config change**, not a missing capability:
drop `channel: 'chrome'`, set `executablePath` to
`~/.cache/ms-playwright/chromium-1248/chrome-linux64/chrome`, and set
`headless: false` if the brief's "headed" requirement is binding.

I recorded this so that the blocking work is only the hardware serialization,
not rediscovery of the toolchain.

---

## 6. One correction to the brief: the plan ref has no WZ1.8 entry

The brief cites
`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` as the "WZ1.8 entry".
That file exists, but it is a **FELIX-Q3-gate** document about the upstream
dev-split decision (`D-275a171c`) and the Minibits DLEQ draft. `grep WZ1.8`
over it returns **0 hits**. The authoritative WZ1.8 text is in the card bodies
on the board, not in that file. Anyone re-dispatching this should fix the ref.

---

## 7. Deviations disclosure (Gate 9 is a final gate — being explicit)

- **No green run. No video. No partial capture.** Not a partial pass — a
  non-start. Per the brief, an incomplete video is a failed task, so this is
  recorded as **FAILED/BLOCKED**, not as progress.
- **No flash was performed**, so no build hash exists to record.
- **No lease was taken or released by this session.** The stale peripheral
  lock (holder `t_2d0c6fa8`, dead pid 2174746) was deliberately left alone;
  only an operator removes an orphan lock.
- **One unintended action, disclosed:** invoking the wizard binary with
  `--help` to read its usage. This binary has **no `--help` flag** — any
  invocation starts the server. It timed out after 180 s and left no process
  behind (`pgrep -af net4sats-wizard` -> empty; nothing listening on 8099). No
  state was changed. This pitfall is already documented in the
  `wizard-e2e-quality-gate` skill; I should have honoured it.
- **No retries were used to mask anything**, because nothing was run.

---

## 8. REMAINING steps (numbered)

1. Release/verify the UART: `ttyUSB0` (FTDI `BG02QAPG`) must no longer be held
   by the spilman lane (`/tmp/spilman-lane-3613961`). Coordinate with the CW
   that owns `spilman-channel-bench-lane.md`. Do **not** force-unlock.
2. Take the serial mutex per the kit contract before any run — all three
   signals must agree: `who-holds-tty.sh /dev/ttyUSB0`, `flock -n`, `pgrep -af`.
   Separately, fix the lease so it actually covers the UART holder (§3.3).
3. Land the spec: cherry-pick `generic-wizard-7step.{spec,config}.mjs` from
   `origin/feat/router-vm-lab` onto `main`, secret-scrubbed, per the standing
   "specs must land on upstream main" rule; close `t_6fd69072` / `t_1eae5a7b`.
4. Fix `generic-wizard-7step.config.mjs` for this host: remove
   `channel: 'chrome'`, set `executablePath` to chromium-1248, honour "headed".
5. Flash/confirm the required upstream build **over the UART console** and
   record the build hash. `physical-router-testing` notes this board needs a
   **sysupgrade** image (the factory/initramfs path does not apply).
6. Serve the wizard at `localhost:8099`, then run the spec with `ROUTER_IP` and
   `ROUTER_PASSWORD` set.
7. Capture video + HTML/JSON report; verify playability with `ffprobe`
   (non-zero duration, expected resolution) and confirm the full 7-step
   traversal including the Advanced-optional step is visible.
8. Re-run once cold to confirm reproducibility; record any flake rather than
   masking it with retries.
9. Commit spec/harness fixes (atomic conventional commits) with docs updated in
   the same commit, and **PUSH**. Release the hardware lease immediately after.

---

## 9. Bottom line

The task is **not startable**. The device is real and I located it precisely,
but it is **serial-only** and **currently occupied by another context window's
live bench** on the very console needed to run anything — and no network path
to it exists at all. The spec precondition is also short of its stated bar
(not on `main`). I stopped at the gate rather than destroying another lane's
run or manufacturing a video that would not have been evidence of anything.

The single genuinely useful thing produced: the harness is proven viable with a
one-line config change, and the video recorder is verified working at
1280x900 — so once the console is free, the remaining work is mechanical.
