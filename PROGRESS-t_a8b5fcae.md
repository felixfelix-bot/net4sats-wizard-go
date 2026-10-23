# PROGRESS — t_a8b5fcae (GL-MT6000 smoke: reachable on required upstream build)

Finding -> Status -> Files touched. Written 2026-10-23.

## 1. Required build is MISSING — task 0 never ran
Finding: the requirements task `t_b009e963` ("Extract WZ1.8 GL-MT6000 requirements")
is `status=triage`, `result=NULL`. The plan ref it names
(`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`) has **0 hits** for `WZ1.8`
— it is a FELIX-Q3-gate document, not the WZ1.8 entry. No upstream commit SHA for
the WZ1.8 flash exists in the plan tree, the board DB, or the repo.
Status: **BLOCKER B1** recorded per the card PREREQ; current build captured anyway
rather than refusing. Required hash was NOT guessed.
Files touched: none.

## 2. Device reached — serial only, identified by its own prompt
Finding: GL-MT6000 is powered and booted: `root@tollgate-0195`, `glinet,gl-mt6000`,
GL.iNet GL-MT6000, bridge MAC `94:83:c4:d1:57:58`. Reachable on `/dev/ttyUSB0`
(FTDI `BG02QAPG`, 115200). The other three UARTs are different boards and were
ruled out by their own prompts: `ttyUSB1`=COMFAST login, `ttyUSB2`=`dlink,covr-x1860-a1`
(`root@tollgate-AF92`), `ttyUSB3`=`glinet,gl-e750` (`root@OpenWrt`).
Status: RESOLVED (identification). Nothing flashed, nothing rebooted.
Files touched: none.

## 3. Running build captured
Finding: OpenWrt `24.10.2`, revision **`r28739-d9340319c6`**, `BUILD_ID` same,
kernel `6.6.93` built 2025-06-23, target `mediatek/filogic`; package
`tollgate-wrt 0.6.0_rc6-r1`. LAN `192.168.1.1/24` (br-lan) + `192.168.2.1/24`
(br-private); STA uplink `10.2.1.122/16` gw `10.2.0.1`.
Status: DONE (capture). **Match/mismatch: unanswerable — required side MISSING.**
Files touched: none.

## 4. Web UI serves — smoke passed
Finding: `curl -sS -o /dev/null -w '%{http_code}\n' -k https://192.168.1.1/` -> `200`.
Body is the live TollGate Admin SPA (`<title>TollGate Admin</title>`,
`/assets/index-C9aCgOwi.js`). Latency ~0.112 s, three consecutive clean probes,
stable ETag `204-2fe-6ad10037`, `Content-Length 766`. Also: `:8090` -> `307` ->
`https://192.168.1.1:8443/` -> `200` LuCI; `:2051/splash.html` -> `200`; `:2121`
is `tollgate-wrt`. uhttpd sends **no Cache-Control/Expires** — stated as a caveat
rather than glossed.
Status: DONE. Executed from the device's own shell (see finding 5), which proves
the UI serves but not host reachability.
Files touched: none.

## 5. No IP path from the bench host
Finding: all device addresses time out from `x280` (`10.2.1.122` 22/80/443/8090/8099,
`192.168.1.1` 22/80/443); host has no `192.168.1.0/24` route. `192.168.2.7` answers
ConnectionRefused but is a **different** GL.iNet host (`94:83:c4:74:d3:03`), not the
testbed. Same environmental fault as `REPORT-fe2b-feed-mt6000.md` §9.
Status: **BLOCKER B2** — blocks the WZ1.8 E2E card's `ROUTER_IP` requirement.
Files touched: none.

## 6. Credentials source recorded (no secret)
Finding: admin web login = `/etc/shadow` root entry, mode `0600` root:root, set at
deploy via `docs/endo-runbook.md` Step 5; SSH = `/etc/dropbear/authorized_keys`
(mode `0600`, 102 B); E2E uses env vars `ROUTER_IP`/`ROUTER_PASSWORD` per
`REPORT-t_75852fce.md:295`; `/etc/config/tollgate` holds only non-secret identity
(`device.code='0195'`, `device.nym='c08r4d0r'`).
Status: DONE. No secret value read or written anywhere.
Files touched: none.

## 7. Bench mutex honoured; one stale lane noted
Finding: `tg_bench_lock_is_held` -> FREE, `fuser -v /dev/ttyUSB0` -> no holders.
Every UART batch wrapped in `tg_bench_lock_acquire`/`_release` from
`/home/c08r4d0r/worktrees/bench-mutex/scripts/tg-bench-lock.sh`. The device shell's
cwd was `/tmp/spilman-lane-3613961` — a prior lane left an open shell on the router
(not a host path). No process held the UART so it did not block, but a spilman lane
may still claim the device.
Status: DONE (released after every batch).
Files touched: none.

## 8. Deliverable namespacing (shared worktree)
Finding: root `PROGRESS.md`/`REPORT.md` in this worktree belong to other cards;
overwriting them would destroy another card's record. The shared
`.git/info/exclude` also blocks `PROGRESS-t_*/REPORT-t_*.md`.
Status: DEVIATION (namespaced files at worktree root, `git add -f` needed).
Files touched: `REPORT-t_a8b5fcae.md`, `PROGRESS-t_a8b5fcae.md`.

## Net status
Verdict **BLOCKER** (B1 required build MISSING, B2 no IP path). UI-serves half of
the GOAL is PASSED with raw evidence. No flash, no config write, no reboot —
correct per the card's BLOCKER RULE. Full detail: `REPORT-t_a8b5fcae.md`.
