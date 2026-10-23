# PROGRESS — WZ1.8 / t_75852fce: generic 7-step wizard E2E on GL-MT6000 + video

Board: net4sats-mvp-v2 | Repo claimed by brief: net4sats-mvp-v2
(`felixfelix-bot/net4sats-wizard-go`) | Host: x280 (c08r4d0r)
Card: t_75852fce (status `triage`, assignee `worker-tollgate`)
Started: 2026-10-23T00:32Z

## Verdict up front

**BLOCKED — precondition gate failed. The task never started as specified.**
Two independent hard preconditions in the brief are false, and one of them is
held by another context window's live work. See REPORT-t_75852fce.md for the
full evidence. Nothing was flashed, no E2E was run, no video was recorded.

## Finding -> status -> files touched

### F1 — Plan ref exists but contains NO WZ1.8 entry (precondition vacuous)
- status: CORRECTED (brief assumption false)
- The plan ref `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` is a
  FELIX-Q3-gate document about the upstream dev-split decision (D-275a171c).
  It contains no WZ1.8 section at all — `grep WZ1.8` -> 0 hits. The brief's
  "PLAN REF ... (WZ1.8 entry)" points at a document that does not describe
  this task. The real WZ1.8 text lives in the card bodies on the board.
- files touched: none

### F2 — Spec-adaptation precondition: spec EXISTS, but NOT pushed to main
- status: PARTIALLY MET (weaker than the brief claims)
- The adapted spec is real and good:
  `tests/browser/generic-wizard-7step.spec.mjs` (259 lines) +
  `generic-wizard-7step.config.mjs`, in repo `physical-router-test-automation`.
  It implements exactly the 7 steps the card names
  (Discover/Connect/Uplink/Payouts/Advanced-optional/Deploy/Live), asserts the
  Advanced `<details>` is closed-by-default and NOT required for deploy
  readiness, and tests the `Live` step against `:2121` + portal.
- BUT it exists ONLY on branch `origin/feat/router-vm-lab`
  (tip ee289bf8, dated 2026-09-27). It is absent from `main` and from the
  fork's sync branch. Its `config.mjs` also requests `channel: 'chrome'` and
  `headless: true`, which contradicts the brief's "run headed".
- The upstream precondition card (`t_6fd69072`) is still `todo`; the
  spec-authoring card `t_1eae5a7b` is still `triage` with a `fleet-done`
  comment but failing gate-ticks (missing tests_green, pushed_or_consolidated,
  ci_evidence, cold_cross_family_review, review_artifact, secrets_clean).
  So "complete, pushed, path known" is only 1/3 true.
- files touched: none

### F3 — Hardware precondition: the MT6000 is REAL and LOCATED, but BUSY
- status: **BLOCKING** — do not start
- Located it positively. It is `/dev/ttyUSB0` on this host (USB serial
  `FTDI_FT232R_USB_UART_BG02QAPG`), identified by its OWN prompt, not by port:
  - `/proc/device-tree/model` -> `GL.iNet GL-MT6000`
  - hostname `tollgate-0195`, `DISTRIB_DESCRIPTION='OpenWrt 24.10.2 r28739-d9340319c6'`
  - `opkg version 38eccbb1...` (opkg, not apk)
- **It is mid-run in another lane's work.** Its live shell prompt is
  `root@tollgate-0195:/tmp/spilman-lane-3613961#` and that directory holds a
  live Spilman-channel bench: `config.json` (161 B), `identities.json` (213 B,
  mode 600), `mod.log` (1862 B), `mod.pid` (6 B), all created 2026-10-22
  21:38-21:41. Uptime 2 days.
- The board's own record agrees: `physical-router-testing` ->
  `references/spilman-channel-bench-lane.md` states the MT6000 is serial-only
  on x280, reachable ONLY through a UART console — i.e. through the very
  console this live lane is occupying. There is NO IP path to it
  (see F5).
- files touched: `evidence/mt6000-state.txt` (outside the repo)

### F4 — The GL-MT6000 has NO network path; it is serial-only
- status: CONFIRMED (probe validated against a known positive)
- The brief implies network automation. There is none. Measured:
  - `enp0s31f6` (the kit's documented wired bench NIC): admin-down AND
    `carrier=0` — no cable. Brought up via `sudo -n ip link set up`; still
    `carrier: 0`.
  - Full probe of `192.168.2.0/24` on ports 22/2121/2051/8090: only this
    host's own interfaces. Positive control first: `192.168.2.1:80` OPEN
    (the probe is valid). Zero TollGate-answering hosts.
  - `192.168.2.1` is a **stock GL.iNet router with unrelated firmware** —
    `buildInfo.js` reports `project: HH63`, `version 2585f2c`,
    `buildTime 2023-09-25`; ports 2121/2051/8090/8080 all `000`. Not the lab
    unit, not TollGate.
  - Sweep of `10.107.73.0/24` (wifi uplink): 3 hits, all SSH-only, none
    OpenWrt, none TollGate.
  - `192.168.1.1` (the spec's default `ROUTER_IP`) is unroutable — packets go
    out the default gateway and die. tcp 22/80/8080 all closed.
  - The kit's `config/routers.lab-backup.json` names an unrelated COVR-X1860
    at `192.168.13.231` — unreachable.
  - The kit's own read-only gate agrees: `scripts/hw-readonly-check.sh
    --host 192.168.1.1` -> `RESULT: FAIL`, gates_failed 3
    (`:2121` answered `000`).
- files touched: none

### F5 — Serial map on this host (verified by reading each console's prompt)
- `ttyUSB0` FTDI `BG02QAPG` -> **GL-MT6000, hostname tollgate-0195** (BUSY,
  other CW's spilman lane) — this is the target device
- `ttyUSB1` FTDI `BG046H14` -> COMFAST box (`COMFAST login:`)
- `ttyUSB2` FTDI `BH00D6DL` -> D-Link COVR-X1860 A1, OpenWrt 25.12.5
- `ttyUSB3` FTDI `BH00DL3M` -> silent (no console responding)
- All four mutexes confirm FREE (`who-holds-tty.sh` -> FREE  /dev/ttyUSB0;
  `flock -n ~/.local/state/tg-bench/tollgate-A04E.lock` -> FREE). The bench
  lease therefore does NOT reflect the occupying spilman lane — the lane holds
  the UART, not the lease. **Recording this as a real lock-coverage gap.**
- files touched: `tmp-wz18/serial_probe*.py` (scratch, outside the repo)

### F6 — Harness half is VIABLE (proven, with a one-line config change)
- status: MITIGATION IDENTIFIED
- The kit repo has no `node_modules` and npm/npx are not on PATH. But the
  harness IS viable, proven by launching a real browser and recording a real
  video:
  - `playwright` 1.52.0 is available from the npx cache
    (`~/.npm/_npx/c61c9351a0dbcfa7/node_modules/.bin/playwright`).
  - The config's `channel: 'chrome'` **FAILS**: `Chromium distribution
    'chrome' is not found at /opt/google/chrome/chrome`.
  - Bundled chromium also fails at the pinned revision
    (`chromium_headless_shell-1169` missing).
  - **But** the present build works by explicit `executablePath`:
    `~/.cache/ms-playwright/chromium-1248/chrome-linux64/chrome` launched
    successfully and the video recorder produced a valid file.
  - Verified artifact: `vid-exepath/81af0231....webm` ->
    `codec vp8, 1280x900, duration 3.04s`. Real, playable, correct resolution.
- Conclusion: the fix is a config change (drop `channel: 'chrome'`, add
  `executablePath` for chromium-1248), NOT a missing capability. Node here is
  18.19.1; the newer playwright in the other npx cache requires Node >= 20,
  which is why 1.52.0 + chromium-1248 is the working combination.
- files touched: `tmp-wz18/harness_check.cjs`, `harness_check2.cjs` (scratch)

## Files touched in the repo

- `PROGRESS-t_75852fce.md` (this file) — new
- `REPORT-t_75852fce.md` — new
- No source, spec, or config file was modified. No flash, no E2E run.

## REMAINING steps (numbered, for whoever unblocks this)

1. Release/verify the UART: `ttyUSB0` (FTDI `BG02QAPG`) must stop being held by
   the spilman lane (`/tmp/spilman-lane-3613961`). Coordinate with the CW
   owning `spilman-channel-bench-lane.md`. Do NOT steal it — the lane's own
   rule is "never force-unlock".
2. Take the serial mutex per the kit contract before any run:
   `who-holds-tty.sh /dev/ttyUSB0` + `flock -n` + `pgrep -af` must agree.
   Fix the lease so it actually covers the UART holder (see F5 gap).
3. Land the spec: cherry-pick `generic-wizard-7step.{spec,config}.mjs` from
   `origin/feat/router-vm-lab` onto `main` (secret-scrub first) per the
   standing "specs must land on upstream main" rule.
4. Fix `generic-wizard-7step.config.mjs` for this host: remove
   `channel: 'chrome'`, set `executablePath` to chromium-1248, and set
   `headless: false` if the brief's "headed" requirement is binding.
5. Flash/confirm the required upstream build over the UART console; record the
   build hash. `physical-router-testing` documents the sysupgrade path for this
   board (it needs a sysupgrade image; the factory path does not apply).
6. Serve the wizard at `localhost:8099` (`net4sats-wizard` binary), then run
   the spec with `ROUTER_IP`/`ROUTER_PASSWORD` set.
7. Capture the video + HTML/JSON report; verify playability with `ffprobe`
   (non-zero duration, expected resolution) and that the full 7-step traversal
   including the Advanced-optional step is visible.
8. Re-run once cold to confirm reproducibility; record any flake rather than
   masking it with retries.
9. Commit the spec/harness fixes (atomic conventional commits), update docs in
   the same commit, and PUSH. Release the hardware lease immediately after.
