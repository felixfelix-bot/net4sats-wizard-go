# PROGRESS — Gate-8 review/merge, wizard official-URL change (t_19b1b1c6 re-run)

Session 2026-10-21 (fleet offload, worker-heavy). Finding -> Status -> Files touched.
Prior sections below belong to earlier tasks sharing this worktree root (preserved whole).

1. Finding: Gate-8 step-1 precondition FALSE — cold review is BLOCKED, not APPROVE.
   Status: VERIFIED from primary sources (board DB + live comment read, not the prior card prose).
   Evidence: t_7da567aa comment 2026-10-11 (worker-reviewer-kimi): "VERDICT: BLOCKED — nothing to
   review" with 7 findings; card still status=blocked, result empty; never re-dispatched since.
   Files touched: PROGRESS.md, REPORT.md, PROGRESS-t_19b1b1c6.md, REPORT-t_19b1b1c6.md.

2. Finding: No implementation exists to merge. No wz1.5 branch/PR on either remote.
   Status: VERIFIED live. Upstream net4sats/net4sats-wizard-go main=b436512 still pins
   felixfelix-bot URLs in deploy.go (6 hits; :43/:47/:50). Fork open PRs: #33/#34 only.
   Nearest-candidate ref origin/feat/feed-per-arch-urls@9a5f34f carries the FreedomTechFeed
   pin map, but has NO PR to upstream (fork PR #6 targets fork main; upstream #1 CLOSED).
   Files touched: none in repo source.

3. Finding: UNVERIFIED network check re-verified live TODAY by this gate (step-1 decision).
   Decision RECORDED: feed asset availability+integrity = VERIFIED (HTTP 200, 7,937,246 B,
   sha256 5134f29dd64f9985f9b993ea6a06b76823af221a88a6cb319d59d8448edb3b90 == pin).
   Install-step dry-run = NOT ACCEPTED as satisfied — remains an open blocking gap owned by
   t_deabf3d6 (impossible to run: no change exists to dry-run against).
   Files touched: as in (1).

4. Finding: Full suite NOT green on the would-be merged revision (step 2).
   Status: MEASURED in probe worktrees. Candidate 9a5f34f: build OK, vet OK, go test
   -count=1 -> 2 FAIL. TestSshRunFunction fails at base b436512 too (pre-existing, live lab
   router answers empty-cred SSH identify). TestRenewOSCommandOnCurrentOS is INTRODUCED by the
   branch (absent at base) and FAILS on this host -> branch-introduced red, not green.
   Files touched: transient probe worktrees, removed after measurement.

5. Finding: NEW since the 2026-10-14 refusal — feed published v0.6.0 FINAL (2026-10-21T02:06Z,
   26 assets incl. per-target ipk/apk). Candidate's pinned v0.6.0-alpha1 (2026-09-09) is now
   ~6 weeks / 7 tags stale. Re-pin required before any merge.
   Status: VERIFIED live (GitHub API + asset download).
   Files touched: as in (1).

6. Finding: Q5 (final vs interim channel) is explicitly tracked — NOT silently dropped.
   Status: VERIFIED. Card t_f6c8aa8e (status=triage, assignee=worker-tollgate, created
   2026-10-14) exists with full brief; v0.6.0-final publication strengthens the FINAL-channel
   path but the decision remains that card's deliverable.
   Files touched: as in (1).

7. Gate outcome: MERGE REFUSED — merge commit SHA: NONE (nothing to merge; merging would
   require creating the implementation myself, which is route-back work owned by t_7e5033df ->
   t_deabf3d6 -> t_7da567aa). Matches and re-confirms verdict card t_222a7d2b (2026-10-14).
   Acceptance-criteria deviation stated honestly: "merged commit SHA recorded" is NOT
   satisfiable; recorded value is "none — merge refused".
   Status: DONE. Committed + pushed on branch worker-heavy/t_gate8-19b1b1c6.

