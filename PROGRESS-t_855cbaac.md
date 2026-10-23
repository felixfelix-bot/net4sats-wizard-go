# PROGRESS — t_855cbaac (Phase-2 feed publication status)

## Finding

**Phase-2 feed: NOT PUBLISHED.** No fetchable feed index at any host this project names. Every
opkg/apk index path 404s except a stale, unsigned `v0.5.0-alpha3` alpha index on
`releases.tollgate.me` (`last-modified: Wed, 01 Jul 2026 19:52:40 GMT`, `Packages.sig` + key 404).

What exists is a per-asset release lane: `FreedomTechFeed/packages` `v0.6.0`
(2026-10-21T02:06:45Z, 26 flat assets, `SHA256SUMS` + `SHA256SUMS.sig`, signature + digest both
re-verified this run). Candidate base for a downstream decision:
`https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/` — asset-only, no index.

Two brief corrections, both measured:
1. `schedule.md` does **not** mention the feed (`grep -ciE 'phase.?2|feed|fe2b'` → `0` of 337
   lines). The term is defined in `feed-freedomtech/docs/per-arch-release-assets.md:4`.
2. `deploy.go` at HEAD has **zero** `blossom2`/`orangesync` references (all three constants are
   `felixfelix-bot`); `blossom2.orangesync.tech` survives only in history (`11f63e7:deploy.go:42`).

## Status

**DONE** — read-only probe + evidence record. No repo source file modified. Nothing built or
tested: the task changes no code, and Go was not invoked.

Blocked-on-nothing. Handoff target is the downstream channel-decision card, which consumes
`~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-phase2-feed-status.md`.

## Files touched

- `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-phase2-feed-status.md` — deliverable (new;
  outside any git repo, so it cannot be committed)
- `REPORT-t_855cbaac.md` — full report (new, tracked)
- `PROGRESS-t_855cbaac.md` — this file (new, tracked)
- `PROGRESS.md` — appended one section (shared; was clean when appended)
- `REPORT.md` — **intentionally untouched**: sibling card `t_wz15-part-b` had uncommitted edits
  there during this run, and appending would have swept their work into this commit

## Verification performed

- `curl` index probe matrix (23 URLs × status/bytes/content-type) — full transcript in the
  deliverable, re-runnable via `/tmp/phase2-probe.sh`
- `FreedomTechFeed/packages` v0.6.0 asset downloaded in full (`8629996` B), sha256 recomputed and
  matched against the published `SHA256SUMS` row
- `SHA256SUMS.sig` verified with `ssh-keygen -Y verify -n freedomtechfeed-release-manifest` →
  `Good ... signature`
- GitHub API state for both repos: releases/tags/pages/branches/PR states
- DNS resolution check for six candidate hostnames
