# REPORT — t_be3261d3-lock-claim (claim GL-MT6000 lock for owner t_be3261d3)

Task id:  net4sats-mvp-v2:t_be3261d3 (claim executed on its behalf by this fleet offload card)
Board:    net4sats-mvp-v2
Repo:     net4sats-mvp-v2, branch worker-heavy/t_gate8-19b1b1c6
Host:     x280 (this session's host, user c08r4d0r)
Mode:     BOOKKEEPING ONLY — zero device contact (no ssh/scp/ping/curl/nc to any router),
          zero flashes, zero reboots, zero power-cycles. Zero lock writes.

## VERDICT: CLAIM NOT MADE — BLOCKED

Two independent gates failed. Per the task's own stop rule ("If the claim fails or freeness was
not confirmed upstream, STOP and report the blocker"), no lock file was written, appended to,
modified, or deleted, and no device command of any kind was issued.

### BLOCKER 1 — freeness was NOT confirmed upstream (this alone stops the claim)

The claim card's declared precondition is the freeness gate card `t_c895de5b` ("Confirm GL-MT6000
free and extract lock convention + device id"), whose body ends with: *"If freeness CANNOT be
positively confirmed ... STOP. Report the blocker ... and do not proceed."*

That gate never ran:

    $ sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db \
        'select id,status,coalesce(result,"<NULL>") from tasks where id="t_c895de5b";'
    t_c895de5b|triage|<NULL>
    $ sqlite3 ... 'select count(*) from task_comments where task_id="t_c895de5b";'
    0

`status=triage`, `result` NULL, zero comments ⇒ no verdict exists, positive or negative.
The sibling read-only inspector `t_6d8235d5` is likewise still `triage` with no result.

The earlier claim attempt by `t_d0c3e172` ("Claim GL-MT6000 lock for owner t_be3261d3") is now
`blocked` with `block_kind=completion` and `result` NULL — it claimed on 2026-10-12 and never
reached a passing gate. So the claim deliverable has never been satisfied upstream.

### BLOCKER 2 — the live lock store on this node is HELD by a different owner

Documented store on this node (`verify/t_6b512dad-lock-convention/recon-quote.txt`, Source A/B/F):

    /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock   (mode 600)

Live state at 2026-10-24T01:30:04Z, sha256 `ed8a4573e88605c6fa4c72dcb9c85da706dc9efe18ef6371ed7a65fd90f46582`:

    lock_id=gl-mt6000-192.168.1.1.lock
    device=GL-MT6000
    device_addr=192.168.1.1
    holder=t_2d0c6fa8
    owner_task=net4sats-mvp-v2:t_2d0c6fa8
    claim_ts=2026-10-14T05:38:34Z
    pid=2174746
    phase=FD4 preflight — acquire and hold hardware lock for FD4 deploy + E2E video run (task t_2d0c6fa8)
    release_condition=released ONLY after the E2E/run task (t_557bb77d child 3, Playwright + video) reports completion or explicit abort

Applying the documented release test (the only authoritative one):

    $ grep -c '^released=' <lock>
    0        (rc=1)  ⇒  STILL HELD

The holder owns a *different* task id (`t_2d0c6fa8`, whose card is now `archived`), and the
recorded pid 2174746 is no longer running — so the claim is stale in substance. But staleness is
not release under this convention. The recon note is explicit:

> "Do **not** rename the file, force-clear, or steal a live claim."

> "**Load-bearing rule:** the only fields that matter for mutual exclusion are `owner_task=` plus
> the absence of a `released=` line."

There is no documented takeover, preemption, or staleness-timeout rule for this store (a grep of
the recovered convention for `takeover|preempt|stale|orphan|force` returns only the prohibition
above plus an unrelated advisory-kit note). Writing `owner_task=net4sats-mvp-v2:t_be3261d3` over a
file that carries no `^released=` line would be exactly the "steal a live claim" the convention
forbids, and would destroy the only record of the competing holder. I did not do it.

### BLOCKER 3 — no writable, verifiable claim target exists on this node

`t_d0c3e172` recorded its claim against a **different** store, on a **different** host:

    /home/c03rad0r/physical-router-test-automation/hardware.lock      (peer node, c03rad0r/dq05)

That path does not exist on x280 (`ls` → rc=2, No such file), the peer store is not reachable from
x280 (`ssh c03rad0r@192.168.2.32` → *No route to host*; `dq05` resolves to `dq05.fips`, which fails
name resolution), and its own recorded durable copies are gone on this host:

    /home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/gl-mt6000-lock-claim-record.md   → absent

So the store the prior claim used cannot be read or written from here, and the store that IS here
is held. There is no medium on which a non-destructive claim could be recorded.

## DELIVERABLE (as requested, with honest nulls)

| Field | Value |
|---|---|
| **Exact lock handle / name** | **NOT ISSUED** — no claim operation completed. The handle the claim *would* have used, per the documented convention, is `gl-mt6000-192.168.1.1.lock` at `/home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock` (mode 600). It remains owned by `net4sats-mvp-v2:t_2d0c6fa8`. A fabricated handle is worse than none, so none is reported as claimed. |
| **Claim timestamp** | **NONE** — no claim was written. (Observation window timestamp: 2026-10-24T01:30:04Z.) |
| **Expected release time** | **NONE** — an expected release time was NOT invented. The live holder's `release_condition` gates release on the same run task `t_557bb77d` child 3; per the recon note, `expected_release` is event-driven here and a later task "must take the release-gate task id out of `release_condition`". Any ISO date I could produce would be a guess, not a claim field. |
| **Where persisted** | **No lock state persisted** (correct for a refused claim). The evidence trail IS persisted, sibling-readable, at: `verify/t_be3261d3-lock-claim/lock-state-capture.txt` (raw captures, this repo, committed+pushed), `PROGRESS.md` and this `REPORT.md` at the worktree root, and the card result. |

### Required-field checklist (all would-be fields, none written)

- owner task id `t_be3261d3` — would-be value, **not written**.
- device id/serial, exactly per recon (no substitution): device id **`GL-MT6000`**; authoritative
  board string `glinet,gl-mt6000`; serial **UNRECORDED — a documented absence** (recon §0(d): "the
  lock record's 'Device id/serial' field holds the **model string** in place of a serial"). Had the
  claim proceeded, the field would have read `GL-MT6000` — the model string — not an invented serial.
- claim timestamp UTC ISO-8601 — **not written**.
- expected release time — **not written** (see above).
- release condition — the required literal string is `t_557bb77d`; the wording supplied
  ("released only after run task t_557bb77d child 3 (Playwright + video run that consumes the
  device) reports completion, or on explicit abort") happens to be materially identical to the
  already-live `release_condition` on the held lock, so the gate is unambiguous either way.

## Device untouched — attestation

No `ssh`, `scp`, `ping`, `curl`, `nc`, `adb`, `fastboot`, `sysupgrade`, or any other router-facing
command was executed against GL-MT6000 or any other device in this run. Every command was a local
read: `ls`, `cat`, `sha256sum`, `grep`, `ps`, `sqlite3` SELECT, `git` read-only. Nothing was
flashed, rebooted, logged into, or power-cycled.

## What the parent/manager must do to unblock

1. Run the freeness gate `t_c895de5b` (or `t_6d8235d5`) to a real verdict — promote it out of
   `triage` and let it execute its documented inspect.
2. Resolve the stale hold on `/home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock`.
   Its `owner_task=net4sats-mvp-v2:t_2d0c6fa8` card is `archived` and pid 2174746 is dead, but the
   file carries no `^released=` line, so under the documented convention it is still a live claim.
   Only a documented release (append `released=<UTC> outcome=aborted` + a lane note) or an explicit
   operator decision authorises the overwrite. A worker must not do this unilaterally.
3. Decide the store, deliberately: x280's peripheral-lock store (held, readable) vs the peer
   `c03rad0r` prta kit store (unreachable from here). The recon note already flags that two live
   lock files exist for this one device on two nodes, unsynchronised.
4. Then re-dispatch this card. Nothing about the claim itself is hard — steps 1–3 are the blockers.

## Reproduce every claim above

    ls -la /home/c08r4d0r/.hermes/peripheral_locks/
    cat  /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock
    sha256sum /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock
    grep -c '^released=' /home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock   # 0 => HELD
    ps -p 2174746 -o pid,etimes,cmd                                                          # dead
    sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db \
      'select id,status,coalesce(result,"<NULL>") from tasks where id in ("t_c895de5b","t_6d8235d5","t_d0c3e172");'
    ls -la /home/c08r4d0r/src/physical-router-test-automation/hardware.lock                   # ENOENT

Raw, untrimmed capture of all of the above: `verify/t_be3261d3-lock-claim/lock-state-capture.txt`
