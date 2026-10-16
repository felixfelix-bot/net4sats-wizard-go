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

## Not done, deliberately
No retry/race against the device another operator is mid-transaction on; no kill/attach to that
session; no `apk` action; no writes to the router.
