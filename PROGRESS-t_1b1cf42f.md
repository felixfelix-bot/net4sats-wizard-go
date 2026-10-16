# PROGRESS — t_1b1cf42f: Resolve GL-MT6000 feed base URL and exact src/gz line

Card: `t_1b1cf42f` · Board: net4sats-mvp-v2 · Repo: felixfelix-bot/net4sats-wizard-go
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2 (detached HEAD 5dd47e3) · Host c08r4d0r
Date: 2026-10-16

## FINDING 1 — no `src/gz` feed exists in the published artifacts. STATUS: resolved (negative, evidence-backed)
The card's escape hatch applies: "If the artifacts do not contain a aarch64_cortex-a53
feed, STOP and report to the manager rather than substituting an architecture."
The aarch64_cortex-a53 **package** exists, but there is no feed index — flat release
assets only, no `Packages`/`Packages.gz` at any path.

Evidence (raw in REPORT.md §3):
- `gh release view v0.6.0-rc6 -R FreedomTechFeed/packages --json assets` → 26 assets, all
  flat. For aarch64_cortex-a53/mediatek-filogic: `.apk`, `.ipk`, `-offline.tar.gz`.
  No `Packages`, no `Packages.gz`, no per-arch directory entry.
- HTTP probes all 404: `<release>/Packages`, `<release>/Packages.gz`,
  `<release>/aarch64_cortex-a53/Packages`.
- Master `release-publish.yml` header is explicit: "uploads the resulting .apk / .ipk as
  release assets with deterministic, TARGET-QUALIFIED names".
- `docs/per-arch-release-assets.md`: "The feed publishes tollgate-wrt packages ... at
  stable, deterministic URLs ... the wizard's `tollgateArchAssets` map ... points at these
  URLs". i.e. per-asset URLs, not an index feed.
- `TARGETS.md` documents `scripts/tollgate-release-url.sh` as the resolution path and
  `src/gz` is absent from it.
- `gh search code 'src/gz tollgate'` and `'Packages.gz' --owner FreedomTechFeed` → 0 hits.

## FINDING 2 — the real artifact mechanism is a direct package URL. STATUS: verified live
`sh scripts/tollgate-release-url.sh --check v0.6.0-rc6 aarch64_cortex-a53 mediatek-filogic apk`
→ exit 0, digest verified against the release's signed SHA256SUMS:

  OK  https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
      sha256 f59311148130ad483d1287854104dd4e9916a306e365d44eac35ccbb0b38c99f

## FINDING 3 — release pin. STATUS: verified
Latest `FreedomTechFeed/packages` release = **v0.6.0-rc6** (published 2026-10-15T19:50:37Z,
`isLatest:true`, no drafts outstanding). This is the release any FE2b feed line would be
pinned to, and it is the release the survived-baseline report (`t_53879c67`) also named.

## FINDING 4 — card precondition independently unmet. STATUS: confirmed
`t_df4304d5` (FE2a: "Feed PR#1 rebase + artifact-to-release bridge + arch tuples") still
reads status **todo**. FE2a is the card that would create the release-bridge this card asks
about. Confirms the FE2b chain is gated.

## FINDING 5 — brief workspace is a different checkout than this board's repo. STATUS: flagged
Brief said "Repo: net4sats-mvp-v2"; this worktree's origin is
`felixfelix-bot/net4sats-wizard-go`. Predecessor cards t_00ee52a3/t_53879c67 found the same
and reported it. Loaded prior FE2b reports before working, per that finding.

## Files touched by this run
- This worktree: `PROGRESS-t_1b1cf42f.md`, `REPORT-t_1b1cf42f.md` (task-scoped; root
  PROGRESS.md/REPORT.md belong to other cards and were NOT touched).
- Read-only elsewhere: no device access attempted, no repo files modified, no `--no-verify`.

## Commit / push (observed)
- Pending at time of writing; see the Closing section appended below after the push.

## Status: COMPLETE for this card's scope — the requested line CANNOT be produced
Stopping as the card instructs rather than substituting an architecture or inventing a URL.
The deliverable is the negative finding + the alternative FE2b needs to know about.
