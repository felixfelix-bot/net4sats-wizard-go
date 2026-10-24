# PROGRESS-t_912534bf — Adjudicate every device-contact claim; evidence level 1–4

Card: `net4sats-mvp-v2:t_912534bf` · assignee `worker-heavy` · host x280 (`c08r4d0r`)
Workspace: `/home/c08r4d0r/repos/net4sats-mvp-v2` · branch `worker-heavy/t_gate8-19b1b1c6`
Mode: read-only over all subject artefacts. No device command, no network probe.

## finding → status → files touched

| # | finding | status |
|---|---|---|
| F1 | Neither input card emitted WITHHELD, so the adjudication proceeds. `t_47795719` (locate) emitted nothing (result NULL, 0 comments/events/runs); `t_1e234edf` (sweep) produced a full 5-file inventory + transcript. | DONE |
| F2 | The evidence report under test EXISTS in four revisions (incl. the `t_be2366ea` attachment, sha256 `c45aaa1b…`), and it makes **six negative** device-contact assertions — zero positive ones. Earlier passes that reported it missing were examining a different task instance. | DONE |
| F3 | PRIMARY raw log bound and re-measured: `raw-run.log`, 59 L / 4585 B, sha256 `062c8b1a…c4cdf209`; byte-identical in two on-host trees; family 221 L / 18137 B across 5 logs. | DONE |
| F4 | Every one of the six assertions re-tested against the log: none confirmable as device contact; each carries a disconfirming verbatim quote with absolute path + line. | DONE |
| F5 | **Level = 1 for the audited subject** (`192.168.1.1`/`2121` = 0 hits in all five logs; literals exist only in spec source) ⇒ NEGATIVE / ABSENT. | DONE |
| F6 | Actual endpoint named: `http://localhost:8099` → `127.0.0.1:8099`, 4 resource-snapshots, `Host: localhost:8099`, all status 200. Harness verdict on the device endpoint in its own HAR: `(NONE)`. | DONE |
| F7 | Strongest device-contact evidence located anywhere on the host is **level 4** — a real `{"kind":10021,…}` and `/balance` body from `192.168.1.1:2121`, from the `fe2b` family (2026-10-16, via `cobrador`), NOT from the audited run. Level 3 claim (`10.17.171.1`) from an earlier pass could not be located; level 4 supersedes it. | DONE |
| F8 | Confidence HIGH. Pivotal missing evidence: the audited run's `trace.zip` / `0-trace.network` (absent from this host; runner dir `/home/c03rad0r` absent; only 2026-10-24 archives from a different run exist). | DONE |

## deliverables

- `/home/c08r4d0r/repos/net4sats-mvp-v2/REPORT-t_912534bf.md` — full report (force-added; tracked)
- `verify/adjudicate-t_912534bf/adjudicate-sweep.sh` + `.out` — committed sweep + transcript (4 333 B)
- `verify/adjudicate-t_912534bf/verify-quotes.sh` + `.out` — committed quote-integrity verifier + 18/18 PASS transcript
- `verify/adjudicate-t_912534bf/raw-run.log`, `raw-runD.log` — committed frozen copies of the two cited logs
- this file (force-added; tracked)

## commits / pushes actually observed

- `00a8b96` — `evidence(adjudicate): t_912534bf — bind raw log + device-contact sweep for claim adjudication`
- push: `e8e6274..00a8b96  HEAD -> worker-heavy/t_gate8-19b1b1c6` to `github.com/felixfelix-bot/net4sats-wizard-go.git`
- `git ls-remote origin refs/heads/worker-heavy/t_gate8-19b1b1c6` → `00a8b96…` == local `HEAD` (verified)

## verification actually run (ad-hoc; no project suite exists for this)

- `bash -n adjudicate-sweep.sh` → clean
- `bash adjudicate-sweep.sh` → exit 0, transcript `adjudicate-sweep.out`
- `bash -n verify-quotes.sh` → clean; `bash verify-quotes.sh` → exit 0, `RESULT: PASS=18 FAIL=0`
- negative control: `raw-run.log:2:Running 7 tests using 1 worker`; 59 lines; 7 ESC bytes
- two-copy sha256 agreement on `raw-run.log`
- live re-check of the 5 external quotes (`evidence.log:30/32/33/104-105/115-116`, spec `:18/49/63/243/297/299/311`, `traces-network-dump.txt:13` → `(NONE)`, `dispatch4:111`, `e2e/README.md:86-88`) → all matched
- `grep -o` re-count of the HAR dump: 4 `resource-snapshot`, 5 `"status":200`, 4 `Host: localhost:8099`, 4 `serverIPAddress 127.0.0.1`, 4 `_serverPort 8099`, sizes 19976×2 / 1021 / 1143
- `git ls-files --error-unmatch verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log` → success
- no test-suite green claimed.

## open items handed forward

1. Nothing on this card required a device probe; none was made. The `fe2b` level-4 body rests on a
   committed transcript, not a live re-derivation.
2. `10.17.171.1:2121` (earlier pass's level-3 claim) was not locatable in the subject trees — if a
   future card needs it, it must be re-located and re-quoted from scratch.
