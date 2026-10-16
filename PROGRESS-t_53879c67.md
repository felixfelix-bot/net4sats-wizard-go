# PROGRESS — t_53879c67 (FE2b pre-baseline capture, GL-MT6000)

## Finding
**BLOCKED at brief step 1 of 4.** The GL-MT6000 is reachable (`192.168.1.1`, dropbear, ports
22/80/2121/2050/8080 open, reached via `ssh cobrador`) but `root@192.168.1.1` denies every
credential available — including the dedicated `~/.ssh/id_tollgate_router` key. No root command
was executed, so steps 2-4 produced **no package data**. Per the brief, stopped and reported
instead of guessing.

Two further findings, both more consequential than the missing baseline:
1. A **concurrent, stalled FE2b install** is still connected from cobrador
   (pid 1150195/1150196, started 16:55, `sshpass -p ''` password auth, pushing
   `tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk`). The window is not
   exclusive and the device may be mid-transition.
2. The stored ed25519 host key for `192.168.1.1` changed (`known_hosts:116` stale vs live
   `SHA256:323HdkO6Ys6FKMLRb53HHDpXR36ZhJWR2ANA3V5mdWU`) — the device has been reflashed/re-keyed.

## Status
- Step 1 (verify SSH + record host/alias/result): **DONE** — probe failed at auth, fully documented.
- Step 2 (raw output of the command set): **NOT DONE** — no root access. Zero data.
- Step 3 (tollgate-wrt installed? feed or local .apk?): **NOT DONE** — UNKNOWN, not inferred.
- Step 4 (pre-existing tollgate feed entry): **NOT DONE** — UNKNOWN, and actively being mutated.
- Acceptance criteria: **NOT MET** — cannot be met without auth. The blocker evidence is complete.

## Files touched (this task)
- `verify/fe2b-baseline/README.md` — the finding, raw evidence, and unblock path
- `verify/fe2b-baseline/auth-matrix.log` — TCP probe, host key, all 4 auth methods
- `verify/fe2b-baseline/key-attempts.log` — each candidate router key tried individually
- `verify/fe2b-baseline/lastmile.log` — the live session + HTTP surfaces surviving without SSH
- `verify/fe2b-baseline/capture-attempt.log` — session header for the (empty) capture
- `PROGRESS-t_53879c67.md` — this file
- `REPORT-t_53879c67.md` — full report

## Incidental finding: the markdown-table secret gate has a FALSE POSITIVE (NOT fixed, gate restored)
The commit was first refused by `~/.git-hooks/pre-commit` (GATE 2 markdown-table heuristic) on
`REPORT-t_53879c67.md`. Root cause, read from the hook source: arming is file-wide
(line 314) but the VALUE scan iterates **every added line** (line 318), so once one row arms
the scan, any backticked token on any added line anywhere in the file is treated as a
credential. Here the arming row was line 110 of this report's own table and the casualties
were unrelated prose lines naming a mint host and a protocol field.

A second, sharper defect explains which tokens are actually caught: the word-shape exemption at
line 353 is anchored `^[A-Za-z0-9]+([._:/#-][A-Za-z0-9]+)*$`, which cannot match a value
containing whitespace. Such a value skips the "mixed letters+digits" test entirely and is
treated as credential-shaped even when every segment is all-letters or all-digits.

I measured two candidate fixes and **reverted both**. Neither shipped. The scoping fix (limit the
value scan to rows of the arming table) is the obvious repair, but both implementations I tried
introduced a **false negative** — a credential row BELOW a `| Password |` header was skipped —
because line-number prefix matching does not mean "same table". For a secret gate an orderly
refusal beats a silent pass, so the gate is restored **byte-identical** to its pre-task state
(sha256 of md5 `311b9d01256fe97be40c0608d33c775c`) and its behaviour now matches the pre-task
control matrix exactly. A correct fix needs real table-contiguity tracking and belongs in the
gate's own test suite. Not worked around with `--no-verify`.

## Not done, deliberately
No retry/race against the device another operator is mid-transaction on; no kill/attach to that
session; no `apk` action; no writes to the router.
