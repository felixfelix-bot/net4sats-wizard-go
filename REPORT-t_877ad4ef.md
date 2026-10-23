# REPORT — t_877ad4ef: failing feed check script + RED baseline (FE2b)

Worker: fleet offload lane · node x280 · board net4sats-mvp-v2
Repo/worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Commit: **80c8d56** on new branch **worker-heavy/t_877ad4ef** (pushed, remote sha verified)
Script: **scripts/fe2b_feed_check.sh**

## VERDICT

The RED baseline is produced, reproducible, and falsifiable: **5 passed / 3 failed, exit 1**,
3 consecutive identical runs. The RED set is A2[O1], A2[O2] and A6 — the two provenance
conditions and the feed index that the FE2b feed install must turn green.

## 1. What changed and why the predecessor script could not be used as-is

A `scripts/fe2b_feed_check.sh` already existed at commit d5d2373 (lane t_18d15abe). It is
**ssh-only**, and on this bench the GL-MT6000 has **no IP path**: TCP 22, 2121 and 2051 on
192.168.1.1 are all closed, and `ping` is not installed on this host. Run as-is it reports
"cannot reach the router" — that is an *unrunnable* test, not a *RED* test. A gate whose
failure means "the harness is blind" cannot certify anything.

So the script was rewritten with a **serial-first transport**:

- auto-detects: TCP 22 open -> ssh; else USB UART present -> serial.
- the device is identified by **its own prompt**, never by port name. `/dev/ttyUSB0`
  answered `root@tollgate-0195` / `glinet,gl-mt6000`; `/dev/ttyUSB1` answered a
  **COMFAST login**, i.e. ports on this bench are not interchangeable.
- all device reads are read-only; nothing on the router or in the repo is mutated.

`FE2B_TRANSPORT=ssh` still works if the box ever returns to the LAN.

## 2. Device state used to make the checks concrete (re-verified first-hand, 2026-10-23)

Read live over /dev/ttyUSB0, matching the package-state task:

| Fact | Value |
|---|---|
| board | `glinet,gl-mt6000`, GL.iNet GL-MT6000, hostname `tollgate-0195` |
| firmware | OpenWrt 24.10.2 r28739-d9340319c6, mediatek/filogic, aarch64_cortex-a53 |
| package manager | **opkg only** — `apk: not found` (opkg 38eccbb1, 2024-10-16) |
| installed | `tollgate-wrt - 0.6.0_rc6-r1`, `Status: install user installed`, Installed-Time 1792514447 |
| feed line | `/etc/opkg/customfeeds.conf:6` -> `src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages` |
| local staged artifact | `/tmp/tollgate-wrt.ipk` 8,630,160 B (plus jq + nodogsplash .ipk, Oct 20 16:40) |
| feed index on box | **none** — `/var/opkg-lists/` holds only openwrt_{base,core,kmods,luci,packages,routing,telephony} |
| service | `/etc/init.d/tollgate-wrt` RUNNING, `S95tollgate-wrt` enabled, binary `/usr/bin/tollgate-wrt` 13,557,944 B |
| :2121 | HTTP **200**, body is a `kind:10021` advert (JSON) |
| :2051 | `/splash.html` -> **200**; `/` -> **403**; `/index.html` -> 404 |
| host-side feed index | `src/gz` target -> **HTTP 404** (also 404 for `packages.gz`, `/Packages`, release-root `Packages`) |

The exact `src/gz` line quoted above is the one from the feed-URL task (t_1b1cf42f) as it
physically sits on the device today; the checks use it verbatim.

## 3. Assertions the script enforces

Each prints a `PASS`/`FAIL` line; any FAIL makes the exit code non-zero.

| ID | Assertion | Failing today? |
|---|---|---|
| A0 | a transport answers (board_name read back) | no |
| A1 | `tollgate-wrt` is installed | no |
| A2[O1] | origin is the feed — **no locally staged tollgate package** | **YES** |
| A2[O2] | origin is the feed — **tollgate index present in /var/opkg-lists/** | **YES** |
| A2[O3] | installed version equals the version the index advertises | not reached (no index) |
| A3 | service running (authoritative `/etc/init.d/tollgate-wrt running`) and rc.d-enabled | no |
| A4 | `:2121` health returns 2xx | no |
| A5 | portal returns HTTP 200 (`/splash.html`) | no |
| A6 | the `src/gz` feed index is fetchable | **YES** |

Two deliberate design decisions, both documented in-file:

1. **A2 is decomposed into three falsifiable conditions.** opkg 1.x on OpenWrt 24.10 writes
   **no** `Source:`/`Filename:` field into `/usr/lib/opkg/status`, so "origin is the feed"
   is not directly readable and asserting it from status alone would be unfalsifiable.
   Instead: no staged local file, a tollgate index exists, and the installed version is the
   version that index advertises. All three must hold.
2. **Portal 200 is asserted on `/splash.html`, not `/`.** `:2051/` returns 403 *by design*
   (the document root has no `index.html` and directory listing is disabled). The
   predecessor script asserted 200 on `/`, which no install could ever satisfy — a
   permanently-red check is not a gate. The 403 is still printed as a note.

## 4. The baseline — full raw failing output

Command: `sh scripts/fe2b_feed_check.sh`  (from /home/c08r4d0r/repos/net4sats-mvp-v2)

```
=== fe2b_feed_check.sh — FE2b acceptance gate (TDD baseline; RED is the expected start) ===
router    : 192.168.1.1
feed_line : src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages
feed_url  : https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages
transport : auto (requested)
host      : x280  date: 2026-10-23T05:21:33Z

--- A0: a transport to the device works ---
      using transport: serial
PASS  A0 transport=serial answered from device (glinet,gl-mt6000)
      board: glinet,gl-mt6000
--- A1: tollgate-wrt is installed ---
PASS  A1 tollgate-wrt installed
      tollgate-wrt - 0.6.0_rc6-r1
--- A2: tollgate-wrt ORIGIN is the feed, not a local file path ---
FAIL  A2[O1] origin is the feed: a LOCAL staged tollgate package is present (install was file-based, not feed-based)
      local artifact: /tmp/tollgate-wrt.ipk
FAIL  A2[O2] origin is the feed: no tollgate index in /var/opkg-lists/ (nothing to have installed from)
--- A3: the tollgate-wrt service is running ---
      binary: -rwxr-xr-x    1 root     root      13557944 Oct 21 23:32 /usr/bin/tollgate-wrt
PASS  A3 tollgate-wrt running (/etc/init.d/tollgate-wrt reports RUNNING; S95tollgate-wrt enabled)
--- A4: :2121 health returns 2xx ---
PASS  A4 :2121 health HTTP 200
--- A5: portal returns HTTP 200 ---
PASS  A5 portal HTTP 200 at :2051/splash.html
      note: :2051/ returns 403 by design (no index.html; listing disabled) — not asserted
--- A6: the src/gz feed index is fetchable (precondition of A2) ---
FAIL  A6 feed index fetchable: HTTP 404 at https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages (no index -> nothing to install from)

=== RESULT: 5 passed, 3 failed ===
FE2b NOT SATISFIED — this is the expected TDD RED baseline before the feed install.
```

True process exit code captured directly (not through a pipe): **1**.
Raw copies on disk: `verify/fe2b-baseline/red-run-2026-10-23.txt` (committed).

## 5. THE RED SET (what the feed install must flip)

1. **A2[O1]** — `/tmp/tollgate-wrt.ipk` is still on the box: the installed package came from
   a local file, which is exactly what FE2b forbids.
2. **A2[O2]** — `/var/opkg-lists/` contains no tollgate index: nothing was ever installed
   from a feed on this device.
3. **A6** — the `src/gz` target returns **HTTP 404**: the feed index does not exist, so a
   feed install is currently impossible no matter what is run on the router.

A2[O3] is conditional on O2: it becomes reachable only once an index exists, then requires
the index to advertise the same version that is installed.

## 6. Why this RED baseline is trustworthy (not just red)

- **Reproducible**: 3 consecutive runs, `5 passed, 3 failed`, exit 1, and the PASS/FAIL line
  set is **byte-identical** across all three (`diff` empty). Raw: `verify/fe2b-baseline/red-run-2026-10-23.txt`.
- **Falsifiable**: a control run with a deliberately wrong health port and portal path
  (`FE2B_HEALTH_PORT=2122 FE2B_PORTAL_PATH=/nonexistent.html`) flips A4 and A5 to FAIL,
  giving `3 passed, 5 failed` (exit 1). The green checks therefore respond to real device
  behaviour instead of passing unconditionally.
  Raw: `verify/fe2b-baseline/control-run-2026-10-23.txt` (committed).
- **Not a transport artefact**: A0 proves the console answered and returned the board name,
  so every later reading is device output, not a timeout.

Bugs found and fixed while building the serial transport (all would have produced silent,
wrong verdicts — they are recorded because they are the trap for the next lane):

1. the shell's **echo of the wrapped command contains the framing markers**, so recording
   started on the echo line and stopped immediately -> markers are now emitted by `printf`
   with a quote-split (`__TGST""ART__`) so the echo cannot match them;
2. `curl -w` output has **no trailing newline**, so the end marker glued onto the code and
   recording overran into the prompt -> the end marker is preceded by a newline;
3. that added newline then left a **trailing blank line**, which made every `rsh | tail -1`
   read an empty string (this presented as "device did not answer") -> trailing blanks are
   trimmed in the helper;
4. a single lost write on the shared UART made A0 flap -> bounded retries (3) on the
   liveness probe and in the helper.

## 7. Where the script lives / branch statement

- Path: **`scripts/fe2b_feed_check.sh`** (the path the card names; nothing in
  `~/.hermes/plans/` designates a different one — the plan does not mention this script).
- Committed in **80c8d56** and pushed to a **new** branch **`worker-heavy/t_877ad4ef`**
  (remote sha verified: `git ls-remote` returns 80c8d56c0e57b3edc900f94f0775d7dd0270e49e).
- **Explicit branch statement:** the card asks for one. The script is on a feature branch
  (`worker-heavy/t_877ad4ef`), **not** on `main`. It was pushed to its own branch because
  the shared worktree's HEAD (`worker-heavy/t_gate8-19b1b1c6`) belongs to the Gate-8 lane,
  and pushing there would have mixed this commit into another workstream's branch.
- Nothing belonging to another workstream was committed. `verify/wz16-e2e-gate/` is
  untracked and untouched — not mine.

## 8. Honesty ledger / deviations

- **Root `PROGRESS.md` and `REPORT.md` are not written by me.** They carry other cards'
  records (Gate-8, FE2b preflight, WZ) and HEAD is a Gate-8 lane; overwriting them would
  destroy another workstream's evidence. I wrote task-scoped `PROGRESS-t_877ad4ef.md` and
  `REPORT-t_877ad4ef.md` at the worktree root instead, as this repo's convention already
  does (16 such files exist). This is a deviation from the literal instruction, named here.
- All device readings, the HTTP codes and the feed 404 were observed this run, first-hand.
  Nothing was carried over from the predecessor report unverified.
- I did **not** modify the router: no package install/remove, no `opkg update`, no file
  deletion. `/tmp/tollgate-wrt.ipk` was deliberately left in place — removing it would
  destroy the very evidence A2[O1] reports.
- `shellcheck` is not installed on this host; `sh -n` is clean.
- The board's FE2b consumption card (t_53879c67) remains **blocked**, and A6 shows why: there
  is no feed index to consume. This script does not by itself unblock it — it defines what
  "unblocked" must look like.
- **Card state**: t_877ad4ef reads `blocked` with reason `fleet-offload:x280`, set by the
  manager at 05:12:42 — i.e. the offload harness parks its own card while the work runs
  off-node. That is not a failure signal from this task, and I did not fight it: the
  deliverable summary is recorded as a card comment instead. `request-review` is refused
  while the card is in that state ("task is not in running/ready").
- **Brief-vs-reality on paths**: the card's `workspace_path` points at
  `.../workspaces/t_53879c67` — a *different* card's directory, and `branch_name` is empty.
  The real worktree is `/home/c08r4d0r/repos/net4sats-mvp-v2`. I worked the repo named in
  the brief and the branch I created, not the card's stale workspace pointer.
- A prior completion of this same card exists (2026-10-12) with its own
  `fe2b_feed_check.sh` + evidence attachments, against host `10.17.171.1`. This run is the
  re-run against the current bench (`192.168.1.1`, serial-only); the script here supersedes
  that version.

## 9. Post-ship verification (after the final commit)

- Committed script bytes == on-disk bytes:
  `git show HEAD:scripts/fe2b_feed_check.sh | sha256sum` and `sha256sum scripts/fe2b_feed_check.sh`
  both `44f2148d7caf9cf1c4e45aa0180e9bb2e292f6825a33fbd6d35ea403d18f42cb`.
- The **committed** script was re-run against the live device: exit **1**,
  `5 passed, 3 failed`, and the PASS/FAIL lines are `diff`-identical to the committed
  baseline file `verify/fe2b-baseline/red-run-2026-10-23.txt` ("MATCHES COMMITTED BASELINE").
  So the shipped artefact, not just a working copy, produces the documented RED baseline.
- `sh -n scripts/fe2b_feed_check.sh` clean. `shellcheck` not installed on this host.

## 10. Post-ship verification of the repo (`hermes verify`)

The canonical recipe was run and now passes, build -> test -> boot/readiness:

| Phase | Command | Result |
|---|---|---|
| build | `go build ./...` | exit 0 |
| test | `go test ./...` | exit 0 — `ok net4sats-wizard 59.8s` (**full suite green**, incl. `TestSshRunFunction`, which the FE2b preflight lane had recorded as red) |
| boot + readiness | wizard started, polled | HTTP **200** |

Two harness notes, because both nearly produced a false verdict:

1. **`go` is not on PATH** in this session (`/usr/local/go/bin/go`); every `go` command
   needs `export PATH=$PATH:/usr/local/go/bin`. Without it the recipe's phases cannot run
   at all.
2. **The recipe's readiness poll defaulted to `:8000` while the wizard binds `:8099`**, so
   the first run reported `ready: false` with connection-refused even though build and test
   passed. `:8099` was also **held by a sibling offload lane** (worker-heavy's WZ1.9 wizard
   server, pid 583866), so a `--port 8099` run could have been satisfied by *their* process
   rather than mine — that result was discarded as ambiguous, not reported as a pass.
   Readiness was instead proven on **`:8098`, confirmed free beforehand**, with
   `WIZARD_BIND=127.0.0.1:8098` for the start phase only:
   `{"ready": true, "statusCode": 200, "outputTail": "net4sats wizard running on http://127.0.0.1:8098"}`.
3. `WIZARD_BIND` must **not** be exported for the test phase: two tests assert the default
   bind (`TestMainFunctionStructure`, `TestListenAddress` expect `127.0.0.1:8099`), so a
   global override makes them fail. That failure was my contamination, reproduced and then
   removed by scoping the override to the start phase — hence `--phase build --phase test`
   separately from `--phase start`.

Repo-level change in this card is only `scripts/fe2b_feed_check.sh` (+ the two evidence
captures and the two task-scoped docs). The other files touched this session are host
scratch, outside the repo: `~/.local/state/tg-bench/probe_ports.py`,
`~/.local/state/tg-bench/serial_run.py`, `/tmp/*.sh`, `/tmp/t_helper.py`.
