# PROGRESS — t_74b12192 (Check fd4/wz18 serialization)

## Finding
VERDICT: **FREE** — fd4 and wz18 hold no hardware/feed work in flight.

- fd4/wz18 are workstream labels, NOT worker node ids. The fleet roster
  (`~/.hermes/bot/nodes.json`) has no node named fd4 or wz18; the only
  occurrences of those tokens as bare words in the whole of `~/.hermes`
  are inside `~/.hermes/bot/fleet_queue_state.json` task bodies, where
  they denote *competing hardware workstreams* ("so this run does not
  collide with the other hardware workstreams fd4 and fe2b",
  "Never preempt fd4, fe2b, or any other holder", "the FD4 branch").
  So the serialization question is: is any hardware/feed work in flight
  under either label right now?
- Board sweep across ALL 189 board DBs: 14 rows carry a live
  claim_lock/worker_pid/in-flight status — **0 of them** are
  hardware/feed work (all are stale reviews/inspections from
  Sep 2026 + blocked e80 bench rows).
- The net4sats-mvp-v2 board has **0** rows with claim_lock/worker_pid.
  Its 4 live kanban locks are dispatch-process locks, not device claims.
- Live worker processes on this host (x280) at check time held
  NON-hardware tasks (see REPORT.md for the full task text).
- No local hardware/feed work in flight: 0 processes touching
  adb/fastboot/esptool/flash/opkg/apk/playwright; no holder (fuser) on
  /dev/ttyUSB0-3 or /dev/ttyACM0; no live device lock (both
  `~/.local/state/tg-bench/*.lock` are stale and unheld).

## Status
DONE — read-only investigation complete, nothing was modified or
claimed. No hold required; FE2b may proceed (subject to the separate
FE2a precondition gate, which is a different gate).

## Files touched
- PROGRESS.md (this file) — root of worktree
- REPORT.md — full evidence report
- No repo source files touched. No device touched. No kanban state mutated.

---

# WZ1.5 part A — exhaustive artifact-URL-token inventory (t_wz15-part-a)

## Finding
Repo-wide case-insensitive sweep for the four tokens yields **163 text lines in 45 files**
plus 1 binary notice (`./net4sats-wizard`, 2 strings). Tracked-only (`git grep`) is 134
lines; the 29-line delta is entirely `.git/info/exclude`'d `PROGRESS-t_*`/`REPORT-t_*` task
logs in this shared worktree root. Per token: orangesync 17, blossom2 3, felixfelix-bot 89,
releases.tollgate.me 59 (lines; 5 lines carry two tokens).

Classification: **(a) 15** URL-to-replace, **(b) 0** unused constant, **(c) 4** doc link to
move, **(d) 144** leave alone. Class (b) is empty — the only dead constant in this area
(`tollgateOSURL`) was already deleted and survives only as a regression-guard comment.

## Status
DONE. Deliverable written and verified row-by-row against a fresh live grep: 163/163 rows,
0 missing, 0 extra, 0 duplicates, and all 163 `exact line` cells match the source file bytes
after reversing the table's `|` -> `&#124;` escape.

## Files touched
- `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.part-a.md` — the deliverable
- `~/.hermes/plans/2026-08-31-tollgate-batch/wz15-tools/gen_report.py` — generator (kept out of
  the repo so the doc has no self-hits; classified rows are explicit, unclassified rows abort)
- `PROGRESS.md` (this file) + `REPORT.md` — appended, not overwritten (shared with other cards)
- No repo source file touched; no test/build state changed; nothing merged with part B.
