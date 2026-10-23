# REPORT — FELIX Q3 gate: dev-split 0 vs 10 (worker-heavy, 2026-10-23)

Card: FELIX Q3 gate — "Resolve FELIX Q3 gate: get dev-split 0 vs 10 answer before dispatch"
Workspace: /home/c08r4d0r/repos/net4sats-mvp-v2
Deliverable: ~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md
Scope: no code changes (respected — no source file in the repo was modified).

## Verdict

**GATE NOT RESOLVED — but the reason it has been stuck for 9 days is now proven, and the
question is on a live channel for the first time.**

The acceptance criteria ("schedule.md contains a dated, quoted answer to Q3 and the exact
approved upstream issue body") are **NOT met**, and cannot be met by an agent: they require a
human decision. I did not mark anything resolved on a timeout or on a "sounds fine". The
`ANSWER` slot in schedule.md remains explicitly EMPTY.

What I *did* accomplish: found and fixed the actual blocker, corrected two false claims in the
plan file, made the plan file self-contained, and re-delivered the question over a channel that
is demonstrably being read.

## 1. Root cause found — the delivery channel was dead, not ignored

The gate has looked like "question delivered, operator silent" since 2026-10-14. It was never
delivered.

- `wss://relay.orangesync.tech` → **HTTP 502** on `/`. Verified from `x280` over both IPv4
  (`23.182.128.51`) and IPv6 (`64:ff9b::17b6:8033`). TCP/443 *connects* — the TLS listener is up
  but the HTTP/WS layer error-502s. A live authenticated `nak req -k 9 …` returns `rc=3` with
  `connecting to wss://relay.orangesync.tech... status 502` and
  `failed to connect to any of the given relays.`
- `~/.hermes/state/decisions_post_health.json` independently records `"probe": "down"`.
- `~/.hermes/state/decisions_responder.json` = `{"since": 1790753481, "seen": []}` — `seen` is
  **empty**. The responder that consumes operator commands from that group has never ingested a
  single event. Corroborates the relay being unreachable for the whole period.
- LAN fallback `ws://192.168.2.32:7780` (from the manager profile's `nostr_disabled.relays`) →
  `No route to host`.

⇒ The 2026-10-14 nostr post (event `13716130a9d4c9…`) reached no one. This is a **fleet-wide**
problem: any gate delivered only to `operator-decisions` since ~2026-10-12 is silently parked
for the same reason. Recorded in schedule.md; fixing the relay is fleet-ops work, out of scope
here (no code changes).

## 2. Re-delivery — done over Signal, observed to work

Signal is the only platform with working credentials. Matrix is configured in `config.yaml`
but `hermes send --list matrix` → *"no targets found"*. Signal lists **33** live targets.

| # | Target | Observed result |
|---|---|---|
| 1 | `signal:net4sats-MVP` (group `l3N6g5Ge/PPjWERCFRHW7m8uye9OW6+wZ/vjD4qJzpA=`) | `sent` |
| 2 | `signal:c08r4d0r number 3` (operator DM) | `Sent to signal home channel (chat_id: +181****0908)` |

The operator was **live on Signal during this run** — `gateway.log` shows inbound messages from
`c08r4d0r number 3` at 04:32–04:56 on other projects. So this channel is genuinely read.

Message body (also reproduced verbatim in schedule.md):

```
FELIX Q3 GATE (D-275a171c) - DECISION NEEDED

Q3: Should the UPSTREAM default dev-split be 0, with 10 applied only via net4sats downstream?

  A = YES - upstream default 0; net4sats applies 10 downstream
  B = NO  - upstream default stays 10
  C = DEFER - do not dispatch remaining WZ1.4 yet

Blocks WZ1.4 cards: t_cfb0a06a (set dev-split default 0 + gate testnut mint injection behind --test-mode), t_58a11d0d (WZ1.4 parent), t_f765f0f1.
Current repo state: main.go:508 DevSplit default 10; index.html:191 slider value=10 (0-50).

Q3b: Pre-approve filing the Minibits DLEQ upstream issue against OpenTollGate/tollgate-module-basic-go?
  (a) FILE IT - draft is at docs/WZ1.4-minibits-dleq-issue-draft.md (commit ff2290bc, author Felix 2026-10-11)
  (b) DON'T FILE - accept removal-only (upstream fix already shipped: gonuts-tollgate >= v0.10.0, tmb-go pins v0.13.0, proof.Id-keyed DLEQ verification)

REPLY: decide D-275a171c <A|B|C> [and (a)|(b)]

DELIVERY NOTE: the nostr operator-decisions relay (relay.orangesync.tech) has been DOWN (HTTP 502) since ~2026-10-12, so the original delivery path for this gate is dead. This Signal message is the replacement. The gate has been unresolved since 2026-10-14.
```

**No reply was received** after delivery. In-session `clarify` also timed out a second time
(2026-10-23), same as 2026-10-14.

## 3. Two false claims in schedule.md corrected

Both prior statements about the Minibits DLEQ draft were wrong, in opposite directions:

| Claim | Source | Reality |
|---|---|---|
| "tracked in `HEAD`" of `net4sats-mvp-v2` | schedule.md, 2026-10-14 | **FALSE for the dispatch workspace.** `/home/c08r4d0r/repos/net4sats-mvp-v2` (remote `felixfelix-bot/net4sats-wizard-go`, HEAD `417b911`) has no `docs/WZ1.4-minibits-dleq-issue-draft.md` in `HEAD` and no commit `ff2290bc`. |
| "exists nowhere" | `t_f765f0f1` report §2 | **FALSE.** Its searches never covered the `net4sats-mvp-v2` worktrees' tracked trees. |

**The truth:** the draft is real and byte-identical in three worktrees of
`felixfelix-bot/net4sats-mvp-v2`, where commit `ff2290bc` (author `Felix <felix@embedsmart.de>`,
2026-10-11) is an ancestor of `main` and the file is in `HEAD`:

```
/home/c08r4d0r/worktrees/t_be2366ea/docs/WZ1.4-minibits-dleq-issue-draft.md
/home/c08r4d0r/worktrees/t_11d0873e/wt/docs/WZ1.4-minibits-dleq-issue-draft.md
/home/c08r4d0r/worktrees/t_1eae5a7b/docs/WZ1.4-minibits-dleq-issue-draft.md
```

All three: 2,537 B, sha256 `0c9105c394f7f554a1f0c50d89136992e4b8268cb61bc1b31d7b8ca9165be76c`.

The full 62-line draft body is now quoted **verbatim inside schedule.md**, so the operator can
approve it without chasing the file across worktrees, and no future lane needs to re-derive its
location. (This was the specific thing the previous lane `t_f765f0f1` could not deliver — it
reported the approved issue body as unavailable and therefore could not file.)

## 4. What the gate actually asks (restated for the decision)

**Q3.** Should the *upstream* default dev-split be `0`, with `10` applied only via *net4sats
downstream*?

- A = YES — upstream default 0; net4sats applies 10 downstream. Slider stays `0-50` but is
  advanced-only (behind the advanced-settings toggle).
- B = NO — upstream default stays 10.
- C = DEFER — do not dispatch the remaining WZ1.4 work yet.

Current repo state, measured this run: `main.go:508` `DevSplit int // 0-50, default 10`;
`index.html:191` `<input type="range" id="devsplit" min="0" max="50" value="10">`;
`api_test.go:115-116` asserts default `DevSplit == 0`; `deploy.go:795` clamps to `0-50`.

Downstream: `t_cfb0a06a` (set dev-split default 0 + gate the testnut mint injection at
`deploy.go:528-538` behind `--test-mode`) and `t_58a11d0d` (WZ1.4 parent) cannot start until
this is answered. `t_f765f0f1` cannot file the issue without Q3b.

## 5. Acceptance criteria — honest status

| Criterion | Status |
|---|---|
| schedule.md contains a **dated, quoted answer to Q3** | **NOT MET** — no human answer exists. Slot explicitly EMPTY. |
| schedule.md contains the **exact approved upstream issue body** | **PARTIALLY MET** — the exact *draft* body is now quoted verbatim in the file, but it is **not approved**; Q3b is unanswered. |
| No code changes | **MET** — no repo source file modified. |
| Decision must be written down, not verbal | **MET as a process rule** — nothing was marked resolved on a timeout. |

The deliverable file was substantially improved but the *decision* is still owed by a human.

## 6. Files touched

- `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` — **the deliverable.** Corrected the
  draft-location claim; quoted the full draft body verbatim; added the 2026-10-23 delivery
  record with the relay-outage evidence and the verbatim Signal message; flagged the fleet-wide
  relay problem. (Note: `~/.hermes/plans/` is **not** a git repo, so this file is not committed
  by git — it is a plain file on disk. Called out explicitly rather than implying a commit.)
- `PROGRESS-t_19b1b1c6.md`, `REPORT-t_19b1b1c6.md` — this worktree's task-scoped prose, committed.

## 7. REMAINING steps (all human-dependent)

1. Operator replies `A` / `B` / `C` (optionally `(a)`/`(b)`) on Signal — the question is live in
   `signal:net4sats-MVP` and the operator DM.
2. Append that reply **verbatim** into the `ANSWER` slot in schedule.md, with date + channel.
3. Then WZ1.4 dispatch unblocks: `t_cfb0a06a`, `t_58a11d0d`, `t_f765f0f1`.
4. Fleet-ops: fix `relay.orangesync.tech` (502) or re-point `operator-decisions` at a live
   relay — otherwise every other gate delivered there stays silently parked.

I stopped at step 1 of the above: the question is delivered and verified as sent; the answer is
the human's to give.
