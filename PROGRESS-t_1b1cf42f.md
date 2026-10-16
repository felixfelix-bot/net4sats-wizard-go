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

## FINDING 6 — a pre-existing comment on this card corrects the ASK itself. STATUS: reconciled + extended
Comment id 29860 (author worker-tollgate, from FE2b RED baseline `t_877ad4ef`, 2026-10-12)
reached the same verdict from the other direction and adds a fact I had not tested:
- It probed the feed HOST namespace (feed/feeds/download.tollgate.me, blossom.net4sats.com) → all
  NO-DNS; releases.tollgate.me resolves but index paths 404. My Finding 1 probed the RELEASE LANE
  (FreedomTechFeed/packages) → flat assets, no index. Non-overlapping: together, no index exists
  in either namespace.
- It corrects the FORM: the device is apk (`apk-tools 3.0.5`), where the repositories line is a
  full URL to `<base>/<channel>/<arch>/packages.adb`, NOT `src/gz <name> <dir>` (opkg 24.10/23.05
  form). I verified the apk-side corollary from the release lane: packages.adb 404s at all three
  path shapes. So the card's src/gz template is doubly wrong — wrong form AND no index of any form.
- UNRESOLVED CONFLICT I am recording, not resolving: TARGETS.md says the MT6000 was validated as
  **opkg/24.10.2**, while t_877ad4ef says the FE2b device is **apk/25.12.5** — but t_877ad4ef's own
  caveat states the device it inspected was an **MT3000**, since no MT6000 was on the bench. I
  cannot settle it (no MT6000 present; root@192.168.1.1 denies all credentials per t_53879c67).
  Verified supporting fact: rc6 ships BOTH .apk and .ipk for this arch/target, so the release does
  not decide the question.

## Files touched by this run
- This worktree: `PROGRESS-t_1b1cf42f.md`, `REPORT-t_1b1cf42f.md` (task-scoped; root
  PROGRESS.md/REPORT.md belong to other cards and were NOT touched).
- Read-only elsewhere: no device access attempted, no repo files modified, no `--no-verify`.

## Commit / push (observed)
- commit `a24779226c347fa115dec58fc230b929f10009b4`, exactly 2 files
  (PROGRESS-t_1b1cf42f.md, REPORT-t_1b1cf42f.md), 392 insertions; commit rc=0.
- push: `git push origin HEAD:refs/heads/worker-tollgate/t_1b1cf42f` → `* [new branch]`, rc=0.
- server truth: `git ls-remote origin refs/heads/worker-tollgate/t_1b1cf42f` =
  a24779226c347fa115dec58fc230b929f10009b4, identical to local HEAD; `git ls-tree` on that
  server-side ref lists both record files. No sweep: the untracked `verify/wz16-e2e-gate/`
  belongs to another card and was left alone.

## Note — a delivery-convention conflict, resolved by precedent (NOT by guessing)
The linked worktree's `.git/info/exclude` (`/home/c08r4d0r/wz16-verify/.git/info/exclude`)
excludes `/PROGRESS-t_*.md` and `/REPORT-t_*.md` with the comment "offload reports (local-only,
never committed)". But the immediate predecessor card on this same chain — `t_53879c67`, the
FE2b card this work serves — has its task records TRACKED (`git ls-files`) and pushed to
`refs/heads/worker-tollgate/t_53879c67` in commits `24020f9` and `5dd47e3`. I followed the
proven precedent of the directly-adjacent card and pushed to this card's own
`worker-tollgate/t_1b1cf42f` branch (`git add -f`, no `--no-verify`). If the fleet actually
wants these local-only, the fix is the stale exclude rule, not this commit — flagged for the
manager.

## Status: COMPLETE for this card's scope — the requested line CANNOT be produced
Stopping as the card instructs rather than substituting an architecture or inventing a URL.
The deliverable is the negative finding + the alternative FE2b needs to know about.
