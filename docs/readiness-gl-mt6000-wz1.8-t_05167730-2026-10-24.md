# READINESS — GL-MT6000 / WZ1.8 generic 7-step wizard E2E

VERDICT: BLOCKED

**Card:** `t_05167730` — Write GL-MT6000 readiness note with evidence and blockers
**Repo:** `net4sats-mvp-v2` (worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`)
**Produced by:** fleet offload worker (`worker-tollgate`), host `x280` (CobradorWave)
**Authoring date:** 2026-10-24 (UTC)
**Gates:** the WZ1.8 generic 7-step wizard E2E run (card `t_75852fce`)

This note exists at:

```
/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/readiness-gl-mt6000-wz1.8.md
```

A durable, tracked copy is committed into the repo at:

```
docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md
```

> VERDICT is BLOCKED because the BLOCKERS section below is non-empty. Per the
> card's ACCEPTANCE rule: *any open blocker => BLOCKED*. No blocker has been
> softened, worked around, or omitted.

---

## 1. Required fields (card STEPS 1)

Each field is either **measured on this host during this run**, or carried from
the named upstream artifact, or marked **MISSING / BLOCKER**. Freshness of every
carried value is stated, because two upstream artifacts were produced on a
different host/user context than this one and were re-checked here.

| # | Field | Value | Status |
| --- | --- | --- | --- |
| 1 | **Upstream build hash (running)** | `r28739-d9340319c6` — OpenWrt 24.10.2, kernel 6.6.93, built 2025-06-23 (`OPENWRT_BUILD_DATE` epoch `1750711236`); on top `tollgate-wrt 0.6.0_rc6-r1` | MEASURED (device-reported, `REPORT-t_a8b5fcae.md` §3) |
| 2 | **Upstream build hash (REQUIRED)** | `MISSING` — **BLOCKER B-REQ** | **BLOCKER** |
| 3 | **Device IP** | `192.168.1.1/24` on `br-lan` (device's own view); second LAN `192.168.2.1/24` on `br-private`; STA uplink `phy0-sta0` `10.2.1.122/16`, gw `10.2.0.1` | MEASURED (device-reported) |
| 4 | **Device identity** | board `glinet,gl-mt6000` / "GL.iNet GL-MT6000"; hostname `tollgate-0195`; `br-lan` MAC `94:83:c4:d1:57:58`; `tollgate.device.code='0195'` | MEASURED |
| 5 | **Host→device reachability** | `NO IP PATH` from `x280` to `192.168.1.1` — **BLOCKER B-IP** | **BLOCKER** |
| 6 | **Serial path (only working path)** | `/dev/ttyUSB0` — FTDI `FT232R` serial `BG02QAPG`, 115200 8N1, prompt `root@tollgate-0195` | MEASURED |
| 7 | **Runner host** | `CobradorWave` (card named it; `t_b009e963` never named one — see BLOCKER B-REQ) | CARRIED (`t_63ec0ed2`, 2026-10-11) |
| 8 | **Runner host == this host?** | **NO — different host.** `t_63ec0ed2` recorded runner paths under `/home/c03rad0r/...`, which do **not exist** on this host (`/home/c03rad0r` → `No such file or directory`). Local user is `c08r4d0r@x280`. | MEASURED (this run) |
| 9 | **Runner: node version** | runner host `v22.22.1` (carried). **On this host:** `v18.19.1` | CARRIED / differs locally |
| 10 | **Runner: npm / pnpm** | runner `npm 9.2.0`, `pnpm 10.8.0` (carried). **On this host: `npm: command not found`, `npx: command not found`** | CARRIED / **BLOCKER locally** |
| 11 | **Runner: @playwright/test** | runner `1.60.0` (harness `~/physical-router-test-automation`, `package.json` declares `^1.52.0`). **On this host:** no `node_modules/@playwright/test` in the checked-out kit | CARRIED / MISSING locally |
| 12 | **Runner: browser build** | `channel: 'chrome'` → system Google Chrome `151.0.7922.137` (runner). Bundled `chromium-1223` / `chromium_headless_shell-1223` **absent** — **BLOCKER B-BROWSER** | **BLOCKER** |
| 13 | **Runner: video encoder (ffmpeg)** | **BLOCKER B-FFMPEG** — `~/.cache/ms-playwright/ffmpeg-1011/ffmpeg-linux` absent; only `ffmpeg-1013` present → no `.webm` can be produced | **BLOCKER** |
| 14 | **Runner: display / headed-capable** | `xvfb-run` + `Xvfb` present; `$DISPLAY` empty (expected for headless server). Video path uses headed-capable Chrome under Xvfb. | MEASURED |
| 15 | **Video output dir** | `~/.cache/ms-playwright/...`-independent: the 7-step config sets `video: 'on'` with no explicit `outputDir`, so Playwright's default `test-results/` under the config's `testDir` applies. Previously-verified concrete dir: `/home/c03rad0r/physical-router-test-automation/tests/browser/test-results` (runner host) | CARRIED (config read from `origin/feat/router-vm-lab`) |
| 16 | **Video output dir on this host** | `MISSING` — `/home/c08r4d0r/repos/physical-router-test-automation/tests/browser/test-results` does not exist | **MISSING** |
| 17 | **Free disk (this host)** | `149G` avail on `/` (`/dev/mapper/ubuntu--vg-ubuntu--lv`, 247G total, 38% used) | MEASURED |
| 18 | **Free disk (runner host)** | `7.4G` in the video dir's filesystem (was `7.3G` at first read) | CARRIED (`t_63ec0ed2`, 2026-10-11) |
| 19 | **7-step spec availability** | `tests/browser/generic-wizard-7step.spec.mjs` **exists ONLY on `origin/feat/router-vm-lab`** (tip `ee289bf8fb9ee8c54e120ea8b205c6b659b0ca83`), **NOT** on `main` / `fork-main-sync` (checked-out tip `5344e35`) | MEASURED (this run) |
| 20 | **Booked hardware window** | `MISSING` — **BLOCKER B-LEASE** | **BLOCKER** |
| 21 | **fd4 serialization** | UNRESOLVED — the live lock store is held by the fd4 side (`owner_task=net4sats-mvp-v2:t_2d0c6fa8`, no `released=` line). See §4. | **BLOCKER** |
| 22 | **fe2b serialization** | UNRESOLVED — no window/lease record exists naming either workstream. See §4. | **BLOCKER** |

Config fact worth recording (read verbatim from the spec header on
`origin/feat/router-vm-lab`): the spec takes `WIZARD_URL` (default
`http://localhost:8099`), `ROUTER_IP` (default `192.168.1.1`, **must appear in
scan results**), `ROUTER_PASSWORD`, `LNURL`, `DEPLOY_MODE`, `DEV_SPLIT`
(default: leave wizard default 10), `MARGIN`, `STRICT_DEPLOY`, `DEPLOY_TIMEOUT`,
`TOLLGATE_PORT` (`2121`), `PORTAL_PORT` (`2050`).

---

## 2. Raw smoke command and captured output (card STEPS 2)

Pasted **verbatim** from `REPORT-t_a8b5fcae.md` §4 (committed `303c76e`). No
reformatting, no paraphrase.

Transport note, quoted from that report: *"there is no IP route from the bench
host to this device, so the smoke was executed from the device's own shell over
the UART console. That is a real HTTP request against the real uhttpd stack
answering on the device's LAN address — it proves the web UI serves. It does not
prove host→device reachability, which is separately blocked."*

The command is exactly the one the card specifies:

```
curl -sS -o /dev/null -w '%{http_code}\n' -k https://192.168.1.1/
```

Raw output:

```
200
```

Content check — the same URL with headers and body, proving it is the real admin
SPA and not an error page:

```
=== HTTPS-BODY-PROOF ===
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/png" href="assets/brand/tollgate/icon-colour.png" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    <meta name="theme-color" content="#111" />
    <meta name="apple-mobile-web-app-capable" content="yes" />
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
    <link rel="manifest" href="manifest.json" />
    <title>TollGate Admin</title>
    <script type="module" crossorigin src="/assets/index-C9aCgOwi.js"></script>
```

Latency / status / "not a cached error page" check:

```
=== REPEAT-3x-WITH-DATE ===
14:48:29
  code=200 t=0.111683s bytes=766 etag="204-2fe-6ad10037"
14:48:30
  code=200 t=0.111865s bytes=766 etag="204-2fe-6ad10037"
14:48:31
  code=200 t=0.112091s bytes=766 etag="204-2fe-6ad10037"

=== NO-CACHE-CONTROL-HEADERS (why !cached) ===
HTTP/1.1 200 OK
ETag: "204-2fe-6ad10037"
Last-Modified: Thu, 15 Oct 2026 16:32:55 GMT
Content-Length: 766
```

Listening surfaces observed in the same session:

```
---PORTS---
tcp  0.0.0.0:80    LISTEN 13284/uhttpd
tcp  0.0.0.0:2050  LISTEN 28652/nodogsplash
tcp  0.0.0.0:443   LISTEN 13281/uhttpd
tcp  0.0.0.0:8090  LISTEN 13280/uhttpd
tcp  :::2121       LISTEN 14349/tollgate-wrt
```

### 2.1 Re-measurement of the reachability claim, this run

The smoke report is dated 2026-10-23. Because this host has since gained a
`192.168.1.0/24` address (`enx6c6e0750cd79` = `192.168.1.50/24`), the
"no IP path" claim was **re-tested first-hand** rather than carried. It still
holds:

```
$ ip -brief addr | grep -E "192\.168|10\.2\."
enx6c6e0750cd79  UP  192.168.2.6/24 metric 100 192.168.1.50/24 ...
enx00e04c36043e  UP  192.168.2.3/24 metric 100
archy-lab-br     DOWN 10.99.99.2/24 192.168.1.2/24 ... linkdown

$ python3 -c '<socket probe results, this run>'
192.168.1.1:22    TimeoutError: timed out
192.168.1.1:80    OSError: [Errno 113] No route to host
192.168.1.1:443   TimeoutError: timed out
192.168.1.1:8090  OSError: [Errno 113] No route to host
192.168.1.1:2121  TimeoutError: timed out
# same result when the source address is pinned to the 192.168.1.0/24 iface:
from .50: TimeoutError timed out

$ ip neigh show | grep '192.168.1.1'
192.168.1.1 dev archy-lab-br FAILED
```

So **B-IP is confirmed open as of this run**, not inherited stale evidence. The
only path to the device remains the UART at `/dev/ttyUSB0`.

---

## 3. BLOCKERS (card STEPS 3 + 4)

Every blocker raised by an upstream task, verbatim in substance, each with the
task that raised it. The PREREQ for this card required *"all upstream
verification tasks must have completed (requirements, lease booking, device
smoke, Playwright/runner)"* — **all four failed to complete**, which is itself
the root of most of the list below.

```
BLOCKER: required upstream build/commit for the WZ1.8 flash — no required build was ever specified
   Raised by: t_b009e963 ("Extract WZ1.8 GL-MT6000 requirements from schedule plan")
   State: status=triage, result=NULL, 0 comments — the task never ran and never produced a build.
   Evidence: the plan ref it names (~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md) contains
     NO WZ1.8 entry at all — grep -n WZ1.8 schedule.md -> 0 hits (it is a FELIX-Q3-gate document).
     No upstream commit SHA for the WZ1.8 flash exists in the plan tree, the board DB, or the repo.
   Consequence: the smoke acceptance criterion is a COMPARISON, and one side is absent. The
     running build is captured (§1 #1); the required build cannot be compared against.
   Also blocks: the runner host was to be named by this task. It never was (t_63ec0ed2 had to
     fall back to "the host the plan implies").

BLOCKER: no reachable IP path from the runner/bench host to the GL-MT6000
   Raised by: t_a8b5fcae ("Verify GL-MT6000 reachable on required upstream build (smoke)")
   State: status=triage, result=NULL; evidence committed as REPORT-t_a8b5fcae.md (303c76e).
   Evidence: 10.2.1.122:{22,80,443,8090,8099} TimeoutError; 192.168.1.1:{22,80,443} TimeoutError;
     192.168.2.7:{22,80,443} ConnectionRefused (decoy GL.iNet host 94:83:c4:74:d3:03, NOT our board
     94:83:c4:d1:57:58). Re-verified this run (§2.1): still unreachable.
   Consequence: the Playwright runner needs ROUTER_IP reachable from the runner host. It is not.
     The device is serial-only (/dev/ttyUSB0) from the bench.

BLOCKER: no lease/serialization mechanism found, and no hardware window was booked
   Raised by: t_77c59d23 ("Book exclusive GL-MT6000 window; note fd4/fe2b serialization")
   State: status=triage, result=NULL, 0 comments — step 4 ("Update the registry entry itself")
     never happened.
   Evidence: no lease/window/booking record exists in ~/.hermes/plans/2026-08-31-tollgate-batch/
     (ls | grep -iE 'lease|window|booking|lock' -> empty). The only mechanism that exists is the
     fleet peripheral-lock store; it has no free-window/booking semantics.
   Consequence: the card STEP 1 field "the booked hardware window with explicit serialization
     notes for fd4 and fe2b" CANNOT be filled. No window start/end, no owner, no release condition.

BLOCKER: GL-MT6000 freeness was never confirmed — the gate card never ran
   Raised by: t_c895de5b ("Confirm GL-MT6000 free and extract lock convention + device id")
   State: status=triage, result=NULL, 0 comments. Sibling inspector t_6d8235d5 also triage/NULL.
   Consequence: every downstream claim that depends on a positive freeness verdict is unbacked.

BLOCKER: the live lock is STILL HELD by a different owner and cannot be taken
   Raised by: t_be3261d3 ("Claim GL-MT6000 lock for owner t_be3261d3") -> verdict BLOCKED_CLAIM_NOT_MADE
   Evidence (re-read this run, unchanged): /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock
     mode 600, 721 bytes, mtime 2026-10-14 07:38, sha256 ed8a4573e88605c6fa4c72dcb9c85da706dc9efe18ef6371ed7a65fd90f46582,
     holder=t_2d0c6fa8, owner_task=net4sats-mvp-v2:t_2d0c6fa8, claim_ts=2026-10-14T05:38:34Z,
     release_condition="released ONLY after the E2E/run task (t_557bb77d child 3, Playwright + video)
     reports completion or explicit abort".
     grep -c '^released=' -> 0 (rc=1) => STILL HELD. Recorded holder pid 2174746 is NOT running;
     holder card t_2d0c6fa8 is archived.
   Convention prohibition (recon-quote.txt:228): "Do not rename the file, force-clear, or steal a
     live claim." No staleness-timeout or takeover rule exists for this store.
   Consequence: staleness in substance is not release in form. A worker must not take this lock.

BLOCKER: no writable, verifiable claim target reachable from this node
   Raised by: t_be3261d3 (blocker B3)
   Evidence: prior claim store /home/c03rad0r/physical-router-test-automation/hardware.lock -> ENOENT here;
     peer ssh c03rad0r@192.168.2.32 -> "No route to host" (rc=255); ssh dq05 -> name resolution failure.
   Consequence: even with a freeness verdict, there is no medium here for a non-destructive claim.

BLOCKER: missing ffmpeg binary for @playwright/test 1.60.0 — no .webm can be produced
   Raised by: t_63ec0ed2 ("Verify Playwright runner toolchain and video capture path")
   Verbatim: "missing ffmpeg binary for @playwright/test 1.60.0 — expected
     ~/.cache/ms-playwright/ffmpeg-1011/ffmpeg-linux (revision 1011); only ffmpeg-1013 present.
     Raw: Error: browserContext.newPage: Executable doesn't exist at
     /home/c03rad0r/.cache/ms-playwright/ffmpeg-1011/ffmpeg-linux -> no .webm can be produced."
   Re-confirmed by that card's own verifier: "video:'on' fails closed on the missing ffmpeg-1011 and
     produces no .webm"; "no new .webm (before=0 after=0)".
   Consequence: the WZ1.8 acceptance is a FULL VIDEO of the 7-step run. Video cannot be recorded.

BLOCKER: missing bundled browser for @playwright/test 1.60.0
   Raised by: t_63ec0ed2
   Verbatim: "missing bundled browser for @playwright/test 1.60.0 — chromium-1223 /
     chromium_headless_shell-1223 absent. Raw: browserType.launch: Executable doesn't exist at
     .../chromium_headless_shell-1223/chrome-headless-shell-linux64/chrome-headless-shell
     -> only channel:'chrome' (system Google Chrome 151.0.7922.137) launches."
   Consequence: the only launch path is system Chrome. The 7-step config pins channel:'chrome',
     so any host without system Chrome cannot launch at all.

BLOCKER: no browser runtime on the bench host
   Raised by: PROGRESS-t_wz18-env-pin.md F5 (and independently observed this run)
   Evidence (this run): node v18.19.1 present, but npm and npx are BOTH absent; no system
     chrome/chromium at any of the 6 probed paths (/usr/bin/google-chrome, google-chrome-stable,
     chromium, chromium-browser, /opt/google/chrome/chrome, /snap/bin/chromium) -> none found;
     no ~/.cache/ms-playwright bundle; no node_modules/@playwright/test in the checked-out kit
     (only /data/worktrees/prta/node_modules/@playwright/test exists, version unknown/pinned).
   Consequence: the 7-step spec CANNOT be executed as-is on this host.

BLOCKER: the 7-step spec does not exist on the checked-out branch
   Raised by: PROGRESS-t_wz18-env-pin.md F3
   Evidence (re-verified this run): tests/browser/generic-wizard-7step.spec.mjs and
     .config.mjs exist ONLY on origin/feat/router-vm-lab (tip ee289bf8fb9ee8c54e120ea8b205c6b659b0ca83).
     Checked-out branch fork-main-sync tips at 5344e35; neither spec file is present there.
   Consequence: the run command is not runnable from the checked-out tree as-is.
```

### 3.1 Note on scope of "upstream blockers"

The card says *"listing every blocker from the upstream tasks"*. The list above
is every blocker from the four PREREQ tasks plus the lock-claim task that
serialization depends on. The separate **fd4 verification-blocker set**
(`verify/fd4_blockers.json`, card `t_31119c9a`, revision v2) was also read for
completeness: it enumerates **9 findings (F1–F9; 6 BLOCKER, 2 MAJOR, 1 MINOR)**
against the *fd4 review lane* (e.g. F1 — "QGATES gate 5 PUSHED FAILS: the lane
is not on the remote"). Those are blockers on fd4's own lane, not blockers on
this WZ1.8 readiness gate, so they are recorded here as context rather than
folded into the list above. Their load-bearing effect on this note is the one
already stated: fd4 is an open, unfinished workstream competing for the same
device.

---

## 4. Booked hardware window + fd4 / fe2b serialization (card STEP 1, last field)

**Booked window: MISSING.** There is none. I did not invent one, per the lease
card's own rule (*"If no lease/serialization mechanism exists, do NOT invent one
silently"*).

What was searched, and what was found:

* `~/.hermes/plans/2026-08-31-tollgate-batch/` — no file matching
  `lease|window|booking|lock` (empty result).
* The fleet peripheral-lock store `~/.hermes/peripheral_locks/` — 5 files:
  `balloon-rx.lock` (0 B), `balloon-tx.lock` (0 B), `tg-bench.lock` (0 B),
  `e80-rx-x280.lock` (237 B, unrelated peripheral), and
  `gl-mt6000-192.168.1.1.lock` (721 B, **held**, §3).
  This store expresses *mutual exclusion right now*; it has **no** window
  start/end, no booking, and no scheduling semantics.
* The bench mutex kit `/home/c08r4d0r/worktrees/bench-mutex/scripts/`
  (`tg-bench-lock.sh`, `bench_lock.py`) — an advisory flock over the UART,
  acquired per-batch and released immediately. Also not a booking mechanism.

**fd4 serialization — UNRESOLVED.** The live lock is *held by the fd4 side*:
`owner_task=net4sats-mvp-v2:t_2d0c6fa8`, `phase=FD4 preflight — acquire and hold
hardware lock for FD4 deploy + E2E video run (task t_2d0c6fa8)`,
`release_condition=released ONLY after the E2E/run task (t_557bb77d child 3,
Playwright + video) reports completion or explicit abort`. The holder's pid is
dead and its card is archived, but **no `released=` line was ever appended**, so
under the documented rule (*"the only fields that matter for mutual exclusion
are owner_task= plus the absence of a released= line"*) it is still held. fd4
must, meanwhile: **release the hold by documented means** (append
`released=<UTC> outcome=aborted` + lane note) **or** obtain an explicit operator
decision. A worker must not do this unilaterally. Until that happens the
serialization note for fd4 is: **fd4 currently owns the device lock and blocks
this run; no window has been agreed.**

**fe2b serialization — UNRESOLVED.** No lease or window record names fe2b at
all, so there is nothing to serialize against. The one fe2b artifact that exists
is `REPORT-fe2b-feed-mt6000.md`, whose own verdict is PARTIAL/BLOCKED upstream
(1 of 5 steps done; no published feed index; device has no working WAN). fe2b is
therefore not *competing* for the device at this moment, but that is an
observation about fe2b's state — it is **not** a booked serialization agreement,
and it must not be read as one. Recorded as: **fe2b has no booked slot; it is
not currently device-active; no agreement exists.**

---

## 5. Code and device mutation disclosure (card STEPS 5)

**No application code was modified by card `t_05167730`.** The only files this
task created or touched are documentation and evidence:

* `docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md` (this note, tracked copy, new)
* `PROGRESS.md` (append only)
* `REPORT.md` (append only)
* `/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/readiness-gl-mt6000-wz1.8.md` (this note, plan-dir copy)

No Go, JS/TS, HTML, CSS, config, or packaging file was edited. No test file was
edited. The wizard binary was not rebuilt.

**Device mutation: NONE.** Zero device contact occurred during this task — no
ssh, scp, serial write, curl, ping, nc, opkg, sysupgrade, reboot, flash, or
power-cycle against the GL-MT6000 or any other device. All device facts in this
note are carried from the cited committed reports, plus local read-only host
measurements.

**Flash / re-image: OUTSTANDING, and owned by the tollgate flow owner.**
Disposition is unchanged on this evidence:

* No flash or re-image was performed, and none is performed here.
* The card `t_b009e963` (requirements) was the task mandated to say *"whether
  the plan requires flashing/re-imaging the device in THIS task or not"*. It
  never ran, so **the flash decision has never been made** (BLOCKER B-REQ).
* The smoke card's own handoff (`REPORT-t_a8b5fcae.md` §9.3) states that if a
  flash turns out to be required it must be handed to the tollgate flow owner,
  and that this board needs a **sysupgrade image** (the factory/initramfs path
  does not apply).
* Therefore: **flash/re-image is outstanding and is owned by the tollgate flow
  owner — not by any worker on this card.** No worker should flash this board on
  the strength of this note.

---

## 6. What must happen for this note to become READY

In dependency order. Every item is a human/manager decision or another card —
none of them is actionable by this card.

1. Dispatch or re-scope `t_b009e963` to produce the **required upstream
   build/commit** with a citation. Note for the re-scoper: the plan ref it is
   told to read contains no WZ1.8 entry (`grep` → 0 hits), so the requirement
   must come from card bodies or from the manager.
2. Give the runner host an **IP path to the device**, or move the run to a host
   that has one (device LAN is `192.168.1.1` / `192.168.2.1`), or repair the
   device's non-routing STA uplink. This is now blocking at least two cards
   (`t_a8b5fcae`, `t_75852fce`) — worth fixing once, not per-card.
3. Repair the runner toolchain: install the **matching `ffmpeg-1011`** bundle and
   the **`chromium-1223` / `chromium_headless_shell-1223`** browser for
   `@playwright/test 1.60.0`, or pin the harness to a Playwright version whose
   bundles are present. Do not silently substitute a browser or disable video.
4. Make the **7-step spec available on the branch the run will use**
   (`origin/feat/router-vm-lab` tip `ee289bf8fb…`, or cherry-pick it forward).
5. Resolve the **fd4 hold** by documented release or explicit operator decision,
   then **book an actual window** with start/end, owner, purpose, and release
   condition recorded in the peripheral-lock store *and* in
   `~/.hermes/plans/2026-08-31-tollgate-batch/`.
6. Record the **fe2b serialization agreement** in the same booking.
7. Re-run the smoke card against the now-known required build, then re-run this
   note. If all of the above land, the VERDICT flips to READY.

---

## 7. Provenance of every value in this note

| Value | Source | Freshness |
| --- | --- | --- |
| Device identity, board, MAC, LAN IPs, running build, tollgate pkg | `REPORT-t_a8b5fcae.md` §2–§3 (committed `303c76e`) | device-reported 2026-10-23 |
| Raw smoke command + output + ports + latency | `REPORT-t_a8b5fcae.md` §4 | 2026-10-23 |
| Reachability re-test | measured this run on `x280` | 2026-10-24 |
| Lock state, sha256, `released=` count | re-read this run; cross-checked `verify/t_be3261d3-lock-claim/lock-state-capture.txt` | 2026-10-24 |
| Lock-claim blockers B1–B3 | `verify/t_be3261d3-lock-claim/blockers.json` | 2026-10-24 |
| Runner host, node/npm/pnpm, Playwright 1.60.0, Chrome 151, ffmpeg/browser blockers, 7.4G | `t_63ec0ed2` task comment (2026-10-11) | 13 days old — flagged |
| Local toolchain (node 18.19.1, no npm/npx, no Chrome, no playwright bundle, ffmpeg 6.1.1, xvfb present) | measured this run on `x280` | 2026-10-24 |
| 7-step spec location + config contract | `git show` on `origin/feat/router-vm-lab`, this run | 2026-10-24 |
| Free disk 149G | `df -h /`, this run | 2026-10-24 |
| fd4 9-finding set (context only) | `verify/fd4_blockers.json` (card `t_31119c9a`, rev v2) | 2026-10-24 |
| Card states (triage / result=NULL) | board DB `sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`, this run | 2026-10-24 |
| fe2b upstream state | `REPORT-fe2b-feed-mt6000.md` | 2026-10-21 |

**Freshness caveat, stated plainly:** the runner-host figures (#9–#12, #16, #18)
were measured 2026-10-11 on `CobradorWave` under user `c03rad0r`. That user and
those paths **do not exist on this host**, so they could not be re-verified here
— only the *this-host* counterparts (#9, #10, #11, #13, #14, #17) could be
measured first-hand. Treat the carried runner numbers as 13 days old until
someone re-runs `t_63ec0ed2` on the runner host itself.

No secret value appears anywhere in this note. No credential was read, printed,
or stored.
