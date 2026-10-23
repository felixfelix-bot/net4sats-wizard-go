# REPORT — t_74b12192: Check fd4/wz18 serialization before any hardware or feed work

Task id: net4sats-mvp-v2:t_74b12192
Board:   net4sats-mvp-v2
Repo:    net4sats-mvp-v2 (branch worker-heavy/t_gate8-19b1b1c6)
Host:    x280 (this session's host)
Mode:    READ-ONLY — no device touched, no kanban state mutated, no repo source edited.

---

## 1. VERDICT

**FREE** — neither fd4 nor wz18 holds hardware/feed work in flight.
**No hold recommendation; FE2b may proceed without collision.**

Scope note: this verdict covers *fd4/wz18 serialization only*. The
separate FE2a precondition gate is a different check and its outcome is
recorded elsewhere (commit 6d75819: "PRECONDITION NOT MET"). This report
does not supersede that gate.

---

## 2. WHAT fd4 AND wz18 ACTUALLY ARE (the key correction)

fd4 and wz18 are **workstream labels**, not worker/agent node ids.

Evidence A — fleet roster has no such node:

    $ cat ~/.hermes/bot/nodes.json | jq -r '.nodes[].id'
    cobrador
    dq05
    hermes-nvme
    testserver2
    t440
    x280
    x240

Neither `fd4` nor `wz18` appears. (`wz18` does not appear as a bare token
anywhere except the body of THIS task card.)

Evidence B — the only bare-word occurrences of fd4/wz18 in all of
`~/.hermes` are inside task bodies in `~/.hermes/bot/fleet_queue_state.json`,
where they name *competing hardware workstreams*, not agents:

    net4sats-mvp-v2:t_75069ebe "Ready GL-MT6000 upstream-build environment and lock hardware slot"
      ...reserve exclusive hardware access so this run does not collide with the
      other hardware workstreams fd4 and fe2b.

    net4sats-mvp-v2:t_75852fce "Run generic 7-step wizard E2E on GL-MT6000 and record full video"
      ...Do NOT start if fd4 or fe2b hold the lease.

    net4sats-mvp-v2:t_6d8235d5 "Execute documented inspect ... report GL-MT6000 lock state"
      ...Never preempt fd4, fe2b, or any other holder.

    net4sats-mvp-v2:t_0841f821 "Announce GL-MT6000 lock on shared channel (fd4/fe2b)"

    net4sats-mvp-v2:t_31119c9a "Fetch FD4 branch, record remote HEAD SHA, enumerate kimi blockers"
    net4sats-mvp-v2:t_e7309ecc "Run blocker-named failing tests against pushed FD4 HEAD"

So the serialization question reduces to: **is any hardware/feed work in
flight under either label right now?** Answer below: no.

---

## 3. BOARD STATE (raw evidence)

### 3a. net4sats-mvp-v2 — the repo under test

    $ sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db \
        "SELECT status, count(*) FROM tasks GROUP BY status;"
    archived|222
    blocked|62
    done|126
    ready|14
    scheduled|2
    todo|26
    triage|150

Rows that carry a live claim_lock or worker_pid, or an in-flight status:

    $ sqlite3 ... "SELECT id,title,status,assignee,claim_lock,worker_pid FROM tasks
                   WHERE claim_lock IS NOT NULL OR worker_pid IS NOT NULL
                   OR status IN ('running','in_progress','claimed','active','doing');"
    0 rows

**No claimed/in-flight task on the net4sats-mvp-v2 board.**

### 3b. ALL 189 board DBs swept for live hardware/feed work

Query per board: rows with claim_lock NOT NULL OR worker_pid NOT NULL OR
status in (running,in_progress,claimed,active,doing).

Result: **14 rows total across all boards — 0 of them hardware/feed.**
The 14 are stale items from Sep 2026 with no live process behind them:

    e80-bench:t_427b43fd        blocked  worker-reviewer-kimi   hb 2026-09-16  "HARM-T6 cold review"
    e80-bench:t_fbc3d2a1        blocked  worker-balloon         hb 2026-09-16  "HARM-T4 ESP32 bench firmware"
    house-keeping:t_467c973e    ready    worker-base            hb 2026-09-14  "Read-only survey of state.db WAL"
    house-keeping:t_3536afa8    blocked  worker-inspector       hb 2026-09-14  "Snapshot live hermes processes"
    pae:t_d7653106              done     worker-pae             hb 2026-09-06  "[PAE-4] conformance vectors"
    plebeian-my-prs:t_dfd4a214  completed worker-reviewer-kimi  hb 2026-09-11  "Review PR #1299"
    plebeian-my-prs:t_97786ee1  done     worker-reviewer-kimi   hb 2026-09-10  "Review PR #1252"
    plebeian-my-prs:t_29be37e1  completed worker-reviewer-kimi  hb 2026-09-10  "Review PR #1232"
    plebeian-pr-reviews:t_2fa49a05 ready  worker-reviewer-kimi-code hb 2026-09-28 "Review PR #1401"
    plebeian-pr-reviews:t_c76fef86 ready  worker-reviewer-kimi   hb 2026-09-28  "Unblock CI on PR #1401"
    plebeian-pr-reviews:t_8de8485a ready  worker-reviewer-kimi-code hb 2026-09-25 "Review PR #1391"
    plebeian-pr-reviews:t_1564de2e done   worker-reviewer-kimi   hb 2026-09-20  "Review PR #1333"
    setupwizzard:4c9cb562...    done     worker-tollgate        hb 2026-09-08  "verify zero net4sats references"
    setupwizzard:451efdc7...    done     worker-tollgate        hb 2026-09-08  "De-brand net4sats references"

Hardware/feed-keyword match among those live rows: **0 of 14**.

---

## 4. LIVE DISPATCH LOCKS AND WHAT THEY HOLD

Snapshot time: 2026-10-23 22:43 local (epoch ~1792788197).

    $ ~/.hermes/bot/.locks/*.json  (live dispatch locks)
    net4sats-mvp-v2:t_2a104c92.json | node x280 | pid 2455141 | ts 22:42:13 | age 72s
    net4sats-mvp-v2:t_6b512dad.json | node x280 | pid 2446938 | ts 22:38:12 | age 313s
    net4sats-mvp-v2:t_74b12192.json | node x280 | pid 2455140 | ts 22:42:13 | age 72s  <= THIS task
    router-maintenance:t_d3a82cc0.json | node x280 | pid 2426256 | ts 22:28:10 | age 916s

    $ ~/.hermes/bot/.fleet_slots/*.slot
    2426256.slot  22:28:10
    2446938.slot  22:38:12
    2455140.slot  22:42:13   <= THIS task
    2455141.slot  22:42:13

These are **harness/dispatch process locks**, not device claims. The
kanban rows for t_2a104c92 / t_6b512dad / t_74b12192 show
status=`triage`, started_at 2026-10-06, last_heartbeat NULL — the
`triage` status is a stale board field; the row is not a hardware claim.

The actual **live worker processes** and the tasks they were dispatched
with (read verbatim from `ps -o args=`), all from 2026-10-23 22:38–22:42:

    pid 2455141 (worker-tollgate) — t_2a104c92
      Task: "Determine whether tollgateOSURL and configwizURL are dead code."
      Nature: Go source grep / reference-count analysis. NO DEVICE. NO FEED.
              Writes a planning .md under ~/.hermes/plans/. NOT hardware/feed.

    pid 2446938 (worker-admin) — t_6b512dad
      Task: "Recover documented GL-MT6000 lock convention (read-only, no device commands)."
      Nature: pure search/extraction over board+workspace notes.
              HARD CONSTRAINT in the card: "runs ZERO lock commands and touches NO
              device." It is researching the lock *convention*; it is not holding
              the device. NOT hardware/feed work in flight.

    pid 2455140 (worker-admin) — t_74b12192  <= THIS task (read-only precheck)
    pid 2426256 (worker-proxy)  — router-maintenance:t_d3a82cc0
      Row: status=blocked, started 2026-09-11, last_heartbeat 2026-09-20.
      Title: ollama pool capacity audit. Unrelated to tollgate hardware/feed.

Conclusion: **no live worker is performing hardware or feed work.**

---

## 5. LOCAL HARDWARE / DEVICE EVIDENCE (x280)

    $ ps -eo args | grep -Ei 'adb|fastboot|minicom|picocom|screen /dev|cu -l|/dev/tty(USB|ACM)|espflash|esptool|openocd|opkg|apk |tollgate|flash|playwright|chromium.*--headless'
    (none)

    $ ps -eo args | grep -Ei 'serial_run|probe_ports|ubus|telnet|socat.*tty|python.*ttyUSB'
    (none)

Serial devices exist but are **unheld** (fuser reports no holder):

    /dev/ttyUSB0  (188,0)  no holder
    /dev/ttyUSB1  (188,1)  no holder
    /dev/ttyUSB2  (188,2)  no holder
    /dev/ttyUSB3  (188,3)  no holder
    /dev/ttyACM0  (166,0)  no holder

USB serial adapters present (FTDI FT232 x4, bus 001 dev 007/014/015/017)
but no process has any of them open.

Device-lock files — both **stale and unheld** (not a live device claim):

    ~/.local/state/tg-bench/tollgate-A04E.lock        mtime 2026-10-23 16:48  (age ~5.9h)  fuser: no holder
    ~/.local/state/tg-bench/tollgate-COVR-X1860.lock  mtime 2026-10-23 14:57  (age ~7.8h)  fuser: no holder

`find ~/.local/state/tg-bench -newermt '-24 hours'` returns only the
stale lock + the static bench scripts; no fresh run artifacts.

Last device-ish activity on this host was therefore ~16:48, ~5.9 h before
this check, and it left no process behind. Nothing is running now.

---

## 6. HOST SCOPE CAVEAT (honest limitation)

This check ran on host **x280**. The fleet is multi-node
(cobrador, dq05, hermes-nvme, testserver2, t440, x280, x240). Board state
and `~/.hermes/bot/` locks were read on x280's copy of the shared state.
- The board sweep covers all 189 board DBs visible from this host, and no
  hardware/feed task is claimed on any of them — that is host-independent
  (kanban state is the shared source of truth).
- Live *process* evidence (ps/fuser/lsusb) is authoritative for x280 only.
- The live GL-MT6000 camera/E2E workstreams in the historical cards
  (t_75069ebe / t_75852fce) ran on a device-adjacent host, not necessarily
  x280. Those cards are **done** (done=True in fleet_queue_state.json) and
  hold no lock, so they impose no constraint now.

No cross-node device lock is currently held: the `router-maintenance`
lock and both tg-bench locks are the only device-adjacent lock artifacts,
and none is held by a live hardware process.

---

## 7. COMMANDS USED (reproducible)

    git status --short; git branch --show-current; git log --oneline -3
    ls -d ~/.hermes/kanban/boards/*/                       # 189 boards
    sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db ".schema tasks"
    sqlite3 ... "SELECT status,count(*) FROM tasks GROUP BY status;"
    sqlite3 <each board> "SELECT id,title,status,assignee,claim_lock,worker_pid,
      COALESCE(last_heartbeat_at,started_at,created_at) FROM tasks
      WHERE claim_lock IS NOT NULL OR worker_pid IS NOT NULL
      OR status IN ('running','in_progress','claimed','active','doing');"
    jq -r '.nodes[].id' ~/.hermes/bot/nodes.json
    grep -rIlE '\b(fd4|wz18)\b' ~/.hermes/bot/            # token provenance
    ls -la ~/.hermes/bot/.locks/ ~/.hermes/bot/.fleet_slots/
    cat ~/.hermes/bot/.locks/*.json
    ps -o pid,lstart,etime,args -p <pid>                  # live worker tasks
    ps -eo args | grep -Ei 'adb|fastboot|esptool|...'      # device processes
    fuser -v /dev/ttyUSB{0,1,2,3} /dev/ttyACM0
    lsusb
    find ~/.local/state/tg-bench -newermt '-24 hours'

---

## 8. BOTTOM LINE FOR THE MANAGER

- fd4/wz18 are workstream labels, not nodes; neither is currently active.
- 0 hardware/feed tasks are claimed or in flight on any of 189 boards.
- 0 device-holding processes; all serial devices and every device lock are unheld.
- The only live workers on x280 are doing Go dead-code analysis and a
  read-only lock-convention search — both explicitly device-free.
- **No hold. FE2b hardware/feed work will not collide with fd4/wz18.**
  (Separately, obey the FE2a precondition gate — see commit 6d75819.)

---

# REPORT — WZ1.5 part A: exhaustive inventory of unofficial artifact-URL tokens

Task: grep the ENTIRE repo (not just `deploy.go`) for `orangesync`, `blossom2`,
`felixfelix-bot`, `releases.tollgate.me` (case-insensitive) and tabulate every hit with a
fixed classification. Deliverable:
`~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.part-a.md`.
Part A only — the URL-constant analysis (part B) and the synthesis are other tasks.

## 1. What was measured

Worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`, branch `worker-heavy/t_gate8-19b1b1c6`,
`HEAD = ed3e548`.

| command | result |
|---|---|
| `grep -rniE 'orangesync\|blossom2\|felixfelix-bot\|releases\.tollgate\.me' . --exclude-dir=.git` | 163 text lines / 45 files + 1 binary notice |
| same with `--exclude=net4sats-wizard` (the ledger's input) | 163 |
| `git grep -nIi -e orangesync -e blossom2 -e felixfelix-bot -e releases.tollgate.me` | 134 |
| `strings net4sats-wizard \| grep -cEi '<tokens>'` | 2 (`felixfelix-bot`, compiled `deploy.go` string table) |

Per token (matching lines / raw occurrences / tracked-only lines):

- `orangesync` — 17 / 17 / 16
- `blossom2` — 3 / 3 / 2
- `felixfelix-bot` — 89 / 92 / 63
- `releases.tollgate.me` — 59 / 59 / 58
- total — **163** lines / 171 occurrences / 139 tracked token-lines

The columns exceed the line total because 5 lines carry two tokens
(`PROGRESS-t_gate8-19b1b1c6.md:13`, `PROGRESS-t_deabf3d6.md:23`, `REPORT-t_deabf3d6.md:74`,
`REPORT-t_gate8-19b1b1c6.md:8`, `REPORT-t_gate8-19b1b1c6.md:31`); `git grep`'s 134 = 139 − 5.

## 2. Reconciliation (the brief's step 2)

`grep -rniE` is strictly broader than `git grep`, and the difference is fully explained:

- 163 − 134 = **29 lines**, all in `.git/info/exclude`'d task logs. The exclude file is
  `/home/c08r4d0r/wz16-verify/.git/info/exclude` (a *shared* git dir — this is a linked
  worktree), and it hides `/PROGRESS-t_*.md` and `/REPORT-t_*.md`. That is why 10 files with
  hits are invisible to `git grep` but present on disk. Nothing else was hidden by
  `.gitignore` except the build binary.
- The tracked subset reproduced per token by my sweep (16 / 2 / 63 / 58) is **identical** to
  `git grep`'s per-token counts, so there is no unexplained residual in either direction.
- No docs live outside the worktree: `docs/` is inside it, and no hit references a path
  outside it.
- Sibling worktrees of the same git dir were counted but deliberately not inventoried
  (wz16-pass6 158, wz16-verify 116, wz-e2e-gate-r3 60, wz-e2e-gate 41, wizard-fix 37,
  wz-e2e-gate-r2 39, wz19-atomicity 18, wz19-base 16, wz19-video 18, wz-ci-atomic/repo 18).
  The brief scopes this to "the ENTIRE repo from the git root", i.e. this worktree.

## 3. Classification outcome

| class | lines | where |
|---|---|---|
| (a) URL to replace with official channel | **15** | `deploy.go` 28,34,43,47,50,63 · `deploy_test.go` 10,14,38,40,43 · `pins_test.go` 54,59 · `docs/router-setup.md` 192,193 |
| (b) unused constant to delete | **0** | — |
| (c) see-also/documentation link to move | **4** | `docs/feed-install.md` 27,28,29,34 |
| (d) unrelated / leave alone | **144** | task logs 128 · evidence captures 12 · CHANGELOG 1 · feed-install transcript 1 · pins_test guard 1 · verify recon quote 1 |

Why (b) is empty: none of the four tokens appears in a declaration of an unused constant.
`tollgateOSURL` — the one dead constant this area ever had — is gone; it survives only as a
forbidden-identifier regression guard (`pins_test.go:45-49`), so that line is (d). The three
live constants carrying the token (`tollgatePkgURL`, `tollgatePkgURLApk`, `configwizURL`) are
all referenced from `deploy.go` and asserted by the pin tests.

Why so much (d): `felixfelix-bot` is *also* the name of this repo's git remote, so every
report header and push transcript matches; `orangesync` is always the fleet Nostr relay host
(`relay.orangesync.tech`, HTTP 502 since ~2026-10-12), never an artifact; `blossom2` appears
only in branding-census lines; `releases.tollgate.me` in logs and evidence files is recorded
probe output, not a shipped reference.

## 4. Constraints handed to the synthesis/implementation task

- `pins_test.go:59` mirrors `deploy.go:43` — repinning one without the other turns
  `TestTollgatePkgURLPinsExistingAsset` RED.
- `deploy.go:67` `artifactPins` is keyed by the URL string, so a repin without a refreshed
  digest makes `verifyArtifactPin` fail closed (by design).
- `deploy_test.go:40` accepts *either* owner, so an owner swap does not fail that guard.
- `pins_test.go:120` only registers constants whose name ends in `URL`, which is why
  `tollgatePkgURLApk` is outside `expectedPinnedURLConsts`/`liveCheckPins`. Any new pin must
  keep the `…URL` suffix or it escapes the live-200 registry.
- `docs/feed-install.md:272` (class d) is a copy-pasteable feed line inside a measured-failure
  transcript; if the host is retired, annotate the fence rather than editing the capture.

## 5. Verification of the deliverable (not a self-report)

The generator (`~/.hermes/plans/2026-08-31-tollgate-batch/wz15-tools/gen_report.py`) refuses to
emit a document if any hit is unclassified. After writing, the following were run against the
file itself:

- 163 table rows parsed, 7 cells each, **0 malformed**.
- Row set vs a **fresh** live grep: **0 missing, 0 extra, 0 duplicates**.
- Every `exact line` cell reversed (`&#124;` -> `|`) compared to the actual source bytes at
  that path/line: **163/163 identical, 0 mismatches**.
- Generator kept **outside** the repo so the document cannot match its own token list (the
  first draft, placed in `verify/`, pushed the count from 163 to 190 — caught and removed).

## 6. Bottom line

Every occurrence of the four tokens in this worktree is in the table; the actionable set is
19 lines (15 (a) + 4 (c)) across 6 files; the rest is internal logs and captured evidence that
must not be rewritten.
