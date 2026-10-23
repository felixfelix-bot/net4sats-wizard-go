# PROGRESS — t_138dc54e: Confirm FE2a merged and release tagged (precondition gate)

Session 2026-10-23 (offload worker). Read-only against the world; only this task's
records written. Finding -> Status -> Files touched.

1. Finding: **Step 1 of the brief cannot be satisfied as worded.** The named plan
   file `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` contains **zero**
   occurrences of FE2a/FE2b/FE2/feed. It is the unrelated FELIX Q3 dev-split gate.
   `grep -r` for `FE2` across the whole plans tree = 0 files.
   Status: VERIFIED. The FE2b definition of done and precondition wording were
   instead recovered **verbatim from the kanban cards themselves**
   (`t_df4304d5` FE2a, `t_53879c67` FE2b, `t_18d15abe` preflight), which is where
   they actually live. This matches the earlier independent finding recorded at
   `REPORT-t_1b1cf42f.md:278-280`.
   Files touched: none.

2. Finding: **FE2a PR is NOT merged.** `FreedomTechFeed/feed` PR #1
   (`ci/auto-feed-indices` -> main, "feat(ci): auto-generate opkg+apk indices on
   every release (all arches)") = state **CLOSED**, `merged:false`,
   `mergedAt:null`, `mergeCommit:null`, closed 2026-09-08T11:16:51Z. Confirmed by
   both `gh pr view --json` and GraphQL `pullRequest.merged`. The branch
   `ci/auto-feed-indices` still exists at `1569f78` — it was abandoned, not merged.
   Status: VERIFIED.
   Files touched: none.

3. Finding: **No release is tagged or published.** `FreedomTechFeed/feed` has
   `releases` length **0**, `releases/latest` -> HTTP **404** via API, tags
   endpoint -> **404** (no tags namespace). No tag, no release, no assets.
   Status: VERIFIED.
   Files touched: none.

4. Finding: The `https://github.com/FreedomTechFeed/feed/releases/latest` URL
   returns HTTP 200 — but that is a **redirect to the empty `/releases` page**
   (`final=https://github.com/FreedomTechFeed/feed/releases`), not a release.
   Recorded so a future lane does not mistake the 200 for a published release.
   Status: VERIFIED (curl -w url_effective).
   Files touched: none.

5. Finding: No feed index is reachable anywhere, so no `src/gz` line could work
   even if a tag existed: `releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/packages.adb`
   -> 404 and `releases.tollgate.me/v0.6.0/...` -> 404.
   Status: VERIFIED (live probes).
   Files touched: none.

6. Status: **DONE — verdict PRECONDITION NOT MET.** Stopped per the brief: no
   device work, no repo/source modification, no FE2b work started. Both halves of
   the precondition are missing. Full evidence in `REPORT-t_138dc54e.md`.
   Files touched: `PROGRESS-t_138dc54e.md`, `REPORT-t_138dc54e.md`.

## Deviation

The brief said "keep PROGRESS.md at the worktree root". Root `PROGRESS.md` /
`REPORT.md` in this repo are gitignored and **shared across cards**, so appending
to them would collide with other lanes' in-flight work. Written as task-scoped
`PROGRESS-t_138dc54e.md` / `REPORT-t_138dc54e.md` at the worktree root instead
(the tracked convention in this repo), and committed with an explicit pathspec.
