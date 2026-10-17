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

---

# RUN 2 — 2026-10-17 (re-derivation, independent of RUN 1)

Trigger: same card re-dispatched. RUN 1 (2026-10-16) concluded STOPPED / no feed index.
This run did NOT simply inherit that verdict: it re-verified every live surface, and found
the AUTHORITATIVE definition of "the feed" in code, which RUN 1 could not see.

## FINDING 7 — the feed's definition exists in code, and it is an ASSET-URL PREFIX, not an index. STATUS: resolved (authoritative)
`OpenTollGate/tollgate-installer` `internal/app/arch.go` (canonical installer, cloned fresh):
```
L107:  feedRepoSlug = "FreedomTechFeed/packages"
L119:  feedReleaseTagDefault = "v0.6.0-rc6"
L322:  func feedReleaseURLPrefix() string {
L323:      return "https://github.com/" + feedRepoSlug + "/releases/download/" + feedReleaseTag + "/"
L350:  func feedAssetURL(arch, target, ext string) string {
L351:      return feedReleaseURLPrefix() + "tollgate-wrt_" + feedPkgVersion() + "_" + arch + "-" + target + ext
```
This is the project's own single source of truth for feed resolution. It is a per-asset URL
prefix. There is no index path, no `<arch>/packages`, no `src/gz` anywhere in it. So the card's
premise (that a `<FEED_BASE_URL>/<arch>/packages` index exists to be pointed at) is false by
the publisher's own design — not merely unpopulated.

## FINDING 8 — three DOCUMENTED feed bases exist in project repos; all are dead or empty. STATUS: verified, all negative
- `https://blossom.net4sats.com/feed/net4sats` — documented in `net4sats/net4sats-feed` README as
  BOTH the apk (`>> /etc/apk/repositories`) and opkg (`src/gz net4sats ...`) base.
  Live: **NO-DNS** (gaierror). Host does not resolve at all.
- `https://releases.tollgate.me/alpha/v0.5.0-beta2/{arch}` — documented in that repo's `MANIFEST`
  as the feed base. Host resolves (Cloudflare) but every index path **404**.
- `https://freedomtechfeed.github.io/feed/<arch>/...` — the target of `FreedomTechFeed/feed` PR #1
  (gh-pages). **404**; repo has no `gh-pages` branch and Pages is not enabled (API 404).

## FINDING 9 — the index-generating work exists but was never landed. STATUS: verified
- `FreedomTechFeed/feed` PR #1 "auto-generate opkg+apk indices on every release (all arches)" →
  state **CLOSED**, `mergedAt: null` (closed 2026-09-08). It is NOT in `main`: main's
  `scripts/` has only `sync-from-upstream.sh`, and no `build-feed.yml`.
- `net4sats/net4sats-feed` `.github/workflows/build-package.yml` is a stub:
  `# TODO: Add blossom upload + NIP-94 kind:1063 event` and
  `echo "Feed index update would go here"`.
So both index lanes were designed and abandoned. This is the root cause of the absence.

## FINDING 10 — a real, LIVE feed exists, but it is `src-git` (build-time), not `src/gz`. STATUS: verified
`FreedomTechFeed/feed` main README:
```
echo "src-git tollgate https://github.com/FreedomTechFeed/feed.git" >> feeds.conf
./scripts/feeds update tollgate && ./scripts/feeds install tollgate-wrt
```
That is a firmware-BUILD-TIME source feed (compiles tollgate-wrt via golang-package.mk). It is
not an on-device binary feed, cannot be consumed by opkg/apk, and its `PKG_VERSION:=0.5.0`
(PKG_RELEASE 1) is stale vs rc6. It does not satisfy FE2b.

## FINDING 11 — the artifacts DO contain aarch64_cortex-a53, as direct packages only. STATUS: verified live (200, digest-matched)
Tag **v0.6.0-rc6** (isLatest, 2026-10-15). Raw curl this run:
- `.../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk` → 302 → **HTTP 200, 8591280 B**
- `.../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk` → 302 → **HTTP 200, 8630160 B**
- SHA256SUMS (200, 3010 B) lists both: apk `f5931114...c99f`, ipk `3b3a85cc...e26a`.
- ALL index paths 404 in both forms and all three path shapes (8 paths probed, see REPORT §2).

## Status: COMPLETE — deliverable is NEGATIVE. The card's STOP condition applies.
No `src/gz` line can be captured because no index is published in any form, at any of the three
documented bases, and the publisher's own code defines the feed as a per-asset URL prefix.
I did NOT substitute an architecture: aarch64_cortex-a53 was confirmed present (Finding 11)
before being rejected as an *index*. I did NOT invent a URL.

## Files touched by this run
- `PROGRESS-t_1b1cf42f.md`, `REPORT-t_1b1cf42f.md` (task-scoped; root PROGRESS.md/REPORT.md
  belong to other cards and were NOT touched). Read-only against all repos/devices.
- Scratch clones under /tmp only (`ftf-feed-probe`, `ns-feed-probe`, `ns-installer`, `og-installer`).
