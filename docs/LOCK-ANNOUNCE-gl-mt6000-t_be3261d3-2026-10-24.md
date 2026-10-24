# GL-MT6000 lock — public announcement (fd4 / fe2b)

**Card:** t_0841f821 (announce) · owner task **t_be3261d3**
**Posted:** 2026-10-24T02:45Z · host `x280` (`c08r4d0r`), repo `net4sats-mvp-v2`
**Mode:** single announcement. No edit, no re-claim, no lock write, no device contact.

---

## 0. Channel identification — and the honest result

**There is no shared lock channel that fd4 and fe2b read.** This is not a guess; it is
the conclusion of a dedicated fleet audit, recovered verbatim at
`verify/t_6b512dad-lock-convention/recon-quote.txt:107-119`:

| Question | Answer (audit) |
|---|---|
| A shared channel announcement? | **NO. A dedicated audit concluded FAIL** — no authoritative shared lock channel exists on any reachable medium. Public Nostr board mirror has columns backlog/inprogress/humanreview/blocked/done — **no lock column**. The kind-9 fleet bus on `wss://relay.orangesync.tech` has tags `fleet-task/fleet-claim/fleet-running/fleet-done/fleet-reply-*` — **no lock tag**. |
| The card meant to create one (`t_0841f821`) | `triage`, never ran (re-confirmed: `result` NULL). |

Two further first-hand checks this session confirm the audit still holds:

```
$ timeout 15 curl -s -o /dev/null -w 'relay_https=%{http_code}\n' https://relay.orangesync.tech/
relay_https=502                       # coordinator down, not a publishable medium
$ ls docs/ | grep -i lock
(no lock note; docs/ holds no lock channel artifact)
```

**Therefore step 1 cannot be satisfied as literally specified.** Per the recon
convention, the sibling-readable announcement medium that *actually works* is
**this repo's `docs/`, committed to the public remote**
(`recon-quote.txt:76-77` and `:231-232`). This file is that announcement. Stated
plainly rather than faked onto a channel that does not exist.

## 1. Channel identifier + message reference (deliverable)

| Item | Value |
|---|---|
| **Channel identifier** | `docs/` in git repo **`felixfelix-bot/net4sats-wizard-go`** (PUBLIC — verified `gh repo view … --json visibility` → `PUBLIC`), branch `worker-heavy/t_gate8-19b1b1c6` |
| **Channel path** | `docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` |
| **Message reference** | git blob at the commit SHA recorded in `PROGRESS.md` §A3 / the final reply's verified commit line |
| **Permalink** | `https://github.com/felixfelix-bot/net4sats-wizard-go/blob/<SHA>/docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` |
| **Medium kind** | public Git repository, commit-addressed (immutable once pushed) — the medium the recon convention names as the working one |
| **Absent alternative** | kind-9 NIP-29 fleet bus on `wss://relay.orangesync.tech` — **no lock tag exists**; relay currently 502. Not used. |

## 2. The announcement message (all five fields explicitly labeled)

> ### GL-MT6000 HARDWARE LOCK — OWNER ANNOUNCEMENT
>
> **owner task id:** `t_be3261d3`
> **device id/serial:** `GL-MT6000` — exact string from recon (§0(d): *"Device id `GL-MT6000`
> (GL.iNet "Flint 2"). Authoritative board string: `glinet,gl-mt6000`"*). **Serial is an
> UNRECORDED, documented absence** — SSH never achieved from any lane, so the lock record's
> "device id/serial" field holds the **model string** in place of a serial. No substitution
> or invention made.
> **claim timestamp (UTC ISO-8601):** `2026-10-12T07:33:17Z`
> **expected release time:** `2026-10-14T07:33:17Z` — a 48 h *estimate* (flash-prep + 4 h
> WZ1.8 budget + margin); **the actual release is event-driven** per the release condition
> below. Source: kanban card `t_d0c3e172` comment id `29831` LOCK CLAIMED, and the
> post-completion VERIFICATION CORRECTION comment id `29832` on the same card.
> **release condition:** released only after run task `t_557bb77d` child 3 (Playwright +
> video run that consumes the device) reports completion, or on explicit abort.
>
> **current lock state (read-only, measured this session):** on this node the store
> `/home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock` is still **HELD** by a
> *different* owner — `owner_task=net4sats-mvp-v2:t_2d0c6fa8`, `claim_ts=2026-10-14T05:38:34Z`,
> sha256 `ed8a4573e88605c6fa4c72dcb9c85da706dc9efe18ef6371ed7a65fd90f46582`, release test
> `grep -c '^released='` → `0` (rc=1) ⇒ **NOT released**. This announcement is bookkeeping; it
> does **not** change that state.

## 3. Field-completeness confirmation

| # | Required field | Present? | Literal in the posted message |
|---|---|---|---|
| 1 | owner task id | ✅ | `t_be3261d3` |
| 2 | device id/serial (exact string from recon) | ✅ | `GL-MT6000` (+ authoritative board string `glinet,gl-mt6000`; serial documented-absent) |
| 3 | claim timestamp (UTC ISO-8601) | ✅ | `2026-10-12T07:33:17Z` |
| 4 | expected release time | ✅ | `2026-10-14T07:33:17Z` (event-driven per condition) |
| 5 | release condition naming literal `t_557bb77d` | ✅ | `…released only after run task t_557bb77d child 3 (Playwright + video run that consumes the device) reports completion, or on explicit abort.` |

All five present. `t_557bb77d` appears verbatim in field 5.

## 4. Provenance of each value (no value invented here)

| Field | Provenance |
|---|---|
| `t_be3261d3` | Card body of t_0841f821 and t_be3261d3; board DB `tasks.id`. |
| `GL-MT6000` / `glinet,gl-mt6000` | `verify/t_6b512dad-lock-convention/recon-quote.txt:85` (§0(d) device identity table). |
| `2026-10-12T07:33:17Z` | `t_d0c3e172` comment 29831 ("Claim timestamp (UTC ISO-8601): 2026-10-12T07:33:17Z"), re-asserted in comment 29832 (13-key parse, 34/34 checks PASS). Peer node `c03rad0r` prta store is not reachable from x280 (`ssh 192.168.2.32 → No route to host`), so this is a *cited* value, not a first-hand read, and is labeled as such. |
| `2026-10-14T07:33:17Z` | Same comment 29831 — recorded as a 48 h **estimate** with release event-driven. |
| release condition + `t_557bb77d` | `t_d0c3e172` comment 29831/29832; the **live** local lock file's `release_condition=` line names `t_557bb77d` too (read this session). |

Note the sibling card's own final deliverable
(`verify/t_be3261d3-lock-claim/blockers.json`, verdict `BLOCKED_CLAIM_NOT_MADE`) records
`claim_timestamp: null` and `expected_release_time: null` because **no claim was made from
x280**. That is not a contradiction: the timestamps above come from the peer-node claim
(`t_d0c3e172`, on behalf of t_be3261d3) that actually issued them. Both are stated together
so a reader is not misled into thinking x280 issued them.

## 5. Compliance with step 4 — no edit, no re-claim

- No lock file created, appended to, modified, or deleted.
- No `released=` line appended (that would be a false release of another owner's claim).
- No SSH/SCP/ping/curl/nc to GL-MT6000 or any device; nothing flashed, rebooted, powered.
- Exactly **one** announcement artifact written. No second copy, no channel spam.
