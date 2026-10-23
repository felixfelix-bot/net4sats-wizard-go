# REPORT — t_f53c5650 (Write /tmp/demo-meta.json from resolved metadata + probed footage)

## VERDICT: BLOCKED — parent task B absent; /tmp/demo-meta.json NOT written

One of the two required parents has **never run on this host**. Per the task's own IF BLOCKED
clause ("Report which field(s) are missing and what you need from the parent tasks"), the
artifact was **not** written, because writing it would require inventing footage paths.

## Parent status — measured live against the board DB, not assumed

Source: `/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`

| parent | card | role | status | task_runs | task_attachments | result |
|---|---|---|---|---|---|---|
| A | `t_85fda20d` | "Resolve prep artifact metadata: PR number, PR title, branch, 40-char head SHA" | triage | **0** | 0 | **NULL** (len 0) |
| B | `t_dd851d4e` | "Locate recording footage, ffprobe each file, and record the step-log path" | triage | **0** | 0 | **NULL** (len 0) |

A board-wide sweep of every `kanban.db` on the host returned these two ids in exactly one
database (this board) — there is no synced copy elsewhere on `/home/c08r4d0r` or `/tmp`, and the
prior sibling evidence path referenced on the card (`/home/c03rad0r/...`) does not exist on this
host (`ls: cannot access '/home/c03rad0r'`).

## What I recovered anyway (provenance-stamped, NOT used to fabricate)

A genuine chain artifact — `attachments/t_f53c5650/resolved_pr_meta_partial-t_f53c5650.json`
(1251 B) — is attached to **this same card** from a prior run. It supplies A's four fields and
**explicitly declares itself PARTIAL with `footage: []` and `step_log: null`**:

```
{"_status": "PARTIAL — t_f53c5650 blocked on missing task-B inputs; /tmp/demo-meta.json deliberately NOT written",
 "_missing": ["footage", "step_log"],
 "_needed_from": "t_882b3922 (E2E recording, 0 completed runs) -> t_dd851d4e (ffprobe probe + step-log path)",
 "pr_number": 32,
 "pr_title": "fix(wz1.4): de-brand wizard business logic — dev split 0, production-only mints, Minibits restored",
 "branch": "wz14-debrand",
 "head_sha": "2d770966b163cd2deafb1697bfd892e0a7a8d9e2",
 "footage": [], "step_log": null}
```

That head_sha was independently corroborated a third time by this board's own completed runs —
`t_0143bb0e` run 1409 (`outcome=completed`) records `head_sha=2d770966b163cd2deafb1697bfd892e0a7a8d9e2`
from a live `gh pr view 32 --repo net4sats/net4sats-wizard-go --json headRefOid`, and `t_069b6842`
cross-checked the same value from rendered intro-card frames. All three agree character-for-character.
Verified programmatically: length **40**, all hex, `short_sha` = `2d770966`.

Recovering A is **not** sufficient: the task says "Use ONLY those outputs", and B's outputs —
the footage list and step_log — do not exist anywhere on this host.

## Why the chain is severed (verified, not inferred)

- `t_882b3922` "Record full E2E de-branded wizard happy-path video" — the recorder B depends on:
  3 runs, **0 completed**, 0 attachments, result NULL. The board has **zero** video attachments
  except `t_7015c244/intro_card.mp4` (80,817 B — an intro card, not a take).
- `/tmp/demo-video.webm` (the path the recorder chain names): **ENOENT**. `/tmp` is tmpfs.
- Prior independent QA on this same chain (`t_8bb677fc`, `t_069b6842`) reached the same verdict:
  "QA VERDICT: BLOCKED — missing input, not a media defect", "0/10 physically checkable".
- **Decoys present, deliberately rejected:** `.../wz19-video/verify/wz19-video-capture/artifacts/take-attempt-blocked.mp4`
  (647,139 B) does probe cleanly — h264 1920x1080 25/1, 42.88 s — and so does the sibling lane's
  `artifacts/demo_basecut/demo-trimmed.mp4` (1,205,700 B, 21.84 s). Neither is B's output:
  B never ran, the first is explicitly named a **blocked attempt** sourced from a *different*
  card (`t_8b5cf200`), and the second is the pre-existing sibling "WZ1.9" lane the wz19 REPORT
  already flags as a derivative of a different demo chain. Substituting either would satisfy the
  mechanical acceptance rows while putting a false path in a machine-readable artifact.

## ACCEPTANCE CRITERIA — honest scorecard (1 of 4 satisfiable)

| criterion | result | evidence |
|---|---|---|
| `/tmp/demo-meta.json` exists, valid JSON, exactly 7 keys | **FAIL — not written** | writing it requires a footage list that does not exist |
| head_sha verbatim, 40 chars, no ellipsis | **PASS (recovered)** | `2d770966b163cd2deafb1697bfd892e0a7a8d9e2`; len 40, all hex; 3 independent sources agree |
| `short_sha == head_sha[:8]` and is a prefix | **PASS (computable)** | `2d770966` — computed programmatically, never typed |
| footage non-empty, entries = parent probe task's paths | **FAIL — parent B never ran** | A's own artifact declares `footage: []` |

Because criterion 4 is unreachable and criterion 2's SHA has no live parent to be compared
against, the 7-key artifact cannot be produced without fabrication. **Not written.**

## Blockers — exactly what is needed from the parents

1. **Parent B (`t_dd851d4e`) must run.** It needs a real recording first.
2. **The recorder (`t_882b3922`, or successor) must produce a take on a host-durable path**
   (not `/tmp`) — currently 0 completed runs, 0 attachments, and `/tmp/demo-video.webm` is ENOENT.
   The demo happy path also needs operator-gated hardware (GL-MT6000 + real SSH) — the one take
   ever captured on this host failed at STEP 7 with `"Cannot connect to router via SSH"`.
3. Only then can B emit the footage list + absolute step_log path, and A/B together unblock
   `/tmp/demo-meta.json`.

## Files touched

- `REPORT-t_f53c5650.md` — this report (new)
- `PROGRESS-t_f53c5650.md` — progress note (new)

No repo source file was modified. No artifact was written to `/tmp/demo-meta.json` —
that absence is the deliverable's honest state, not an oversight.
