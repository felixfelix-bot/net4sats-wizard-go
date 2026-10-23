# REPORT — t_855cbaac: Phase-2 feed publication status

Worker lane `t_gate8-19b1b1c6`. Worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`.
Probed 2026-10-23 20:46–21:00 UTC. Read-only: no repo source file was modified.

**Deliverable:** `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-phase2-feed-status.md`
**Probe script / raw output (this host, ephemeral):** `/tmp/phase2-probe.sh`, `/tmp/phase2-probe.out`

---

## Verdict

**Phase-2 feed: NOT PUBLISHED.** There is no fetchable feed index at any host this project names.
Every opkg/apk index path returns 404 except one stale, **unsigned** `v0.5.0-alpha3` alpha index on
`releases.tollgate.me` (`last-modified: Wed, 01 Jul 2026 19:52:40 GMT`).

What *is* published is a per-asset release lane — `FreedomTechFeed/packages` `v0.6.0`
(2026-10-21T02:06:45Z, 26 flat assets, `SHA256SUMS` + `SHA256SUMS.sig`). That is release
artifacts, not a feed. Its base is
`https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/` and it serves
`<tag>/<asset-file>` URLs only — it cannot express a `src/gz`/`packages.adb` repositories line.

---

## Two corrections to the task brief (both measured, neither speculative)

### 1. `schedule.md` does not describe the Phase-2 feed

The brief says the Phase-2 feed is "described in the plan at
`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`". It is not.

```
$ cd ~/.hermes/plans/2026-08-31-tollgate-batch && grep -ciE 'phase.?2|feed|fe2b' schedule.md
0
$ wc -l schedule.md
337 schedule.md
```

`schedule.md` is the FELIX Q3 / dev-split gate file (Q3 = B, Q3b = DON'T FILE). The term "Phase 2"
in this batch appears in `wz1.5-url-inventory.md:414` (wz1.5's own *implementation* phase — the
downstream repin work, blocked on the channel decision) and, authoritatively for the feed, in
`/data/repos/feed-freedomtech/docs/per-arch-release-assets.md:4`:

> The feed publishes tollgate-wrt packages for every router architecture the wizard can select, at
> stable, deterministic URLs. This is the Phase 2 prerequisite for "test the feed via the installer"

The two sibling kanban cards (`t_66e66ae0`, `t_1df48589`) use the same phrase, and `t_1df48589`
names this deliverable path as its required input — so the reading used here (per-arch published
`tollgate-wrt` assets) is the one the batch itself supports.

### 2. `deploy.go` does NOT currently reference `blossom2.orangesync.tech`

The brief says `deploy.go:24-52` "points at unofficial mirrors (blossom2.orangesync.tech and
felixfelix-bot release URLs)". At HEAD the file has **zero** blossom/orangesync references:

```
$ grep -rn -i -E 'blossom2|orangesync' --include='*.go' . | wc -l
0
```

All three constants are `github.com/felixfelix-bot/...`:

```
deploy.go:43  tollgatePkgURL     = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.7.0-alpha10/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk"
deploy.go:47  tollgatePkgURLApk  = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.6.1-post-merge/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk"
deploy.go:50  configwizURL       = "https://github.com/felixfelix-bot/configurationwizzard/releases/download/v1.0.10-alpha/net4sats-configwiz-1.0.10-alpha.tar.gz"
```

`blossom2.orangesync.tech` exists **only in history**:

```
$ git log --all --oneline -S'blossom2' -- deploy.go
11f63e7 fix(deploy): point tollgate-wrt URL to portal-fix e705430 blossom
$ git show 11f63e7:deploy.go | grep -n blossom
42:	tollgatePkgURL = "https://blossom2.orangesync.tech/c6b1484e...c2bba8.apk"
```

The host itself is alive (`http=200`, `23.182.128.219`) but publishes no feed: `/feed`, `/feed/`,
`/manifest.json`, `/index.json`, `/tollgate/`, `/aarch64_cortex-a53/Packages.gz` are all 404.
This does not change the verdict; it changes who the downstream task argues with — the live
mirrors are `felixfelix-bot` GitHub releases (a fork of upstream), not blossom.

---

## What was probed (verbatim, all in the deliverable)

`wz1.5-phase2-feed-status.md` §1 carries the full 56-line transcript. The load-bearing results:

| target | result |
|---|---|
| `releases.tollgate.me/` | `http=200 text/html` — a "TollGate Release Explorer" SPA, not a repo server |
| `releases.tollgate.me/alpha/v0.6.0/.../{Packages,Packages.gz,APKINDEX.tar.gz,packages.adb,Packages.sig}` | **404** (all five) |
| `releases.tollgate.me/alpha/v0.5.0-alpha3/.../{Packages,Packages.gz,APKINDEX.tar.gz}` | 206 / 206 / 206 — 2026-07-01, **`Packages.sig` 404**, `keys/tollgate-testing.pem` 404 |
| `FreedomTechFeed/packages/releases/download/v0.6.0/{Packages,Packages.gz,aarch64_cortex-a53/Packages.gz,packages.adb}` | **404** (all four) |
| `FreedomTechFeed/packages/.../v0.6.0/{SHA256SUMS,SHA256SUMS.sig}` | 206 / 206 |
| `FreedomTechFeed/packages/.../v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.{ipk,apk}` | 206 / 206 |
| `FreedomTechFeed/packages/.../v0.6.0-rc6/..._0.6.0_rc6_...ipk` | 206 |
| `FreedomTechFeed/packages/.../v0.6.0-alpha1/...` | **404** (the tag `feat/feed-per-arch-urls` pinned) |
| `FreedomTechFeed/feed` releases / tags | **0 / 0** |
| pages API + `gh-pages` branch (both repos) | **404 / 404** |
| `freedomtechfeed.github.io/{feed,packages}` | **404** |
| `blossom.net4sats.com`, `feed.net4sats.com`, `feed/packages/repo.tollgate.me` | **NO-DNS** |

Positive proof the asset lane is real (re-derived this run, not quoted):

```
$ curl -sSL -o ipk ".../v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk"
ipk http=200 bytes=8629996
$ sha256sum ipk
c51c492a9acf504f2f9b87af028aad401ff857df2efd12bf5fe56ba2ae45d868  ipk
$ grep 'tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk$' SHA256SUMS
c51c492a9acf504f2f9b87af028aad401ff857df2efd12bf5fe56ba2ae45d868  tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk

$ ssh-keygen -Y verify -f allowed_signers -I release-signing@freedomtechfeed \
      -n freedomtechfeed-release-manifest -s SHA256SUMS.sig < SHA256SUMS
Good "freedomtechfeed-release-manifest" signature for release-signing@freedomtechfeed
     with ED25519 key SHA256:lQWK9yivuF2bVKLTYaKlu2X3zgRdgfbQf1ymAazgets
```

(The namespace had to be discovered by trying candidates; `ssh-keygen` reports
`namespace does not match` for wrong ones. Worth knowing: the signing key is ssh-ed25519, not
`openssl`-compatible — `openssl dgst -verify` fails with "Could not read public key".)

---

## Why there is no index — absence is by design

- The consumer resolves a single asset: `feedReleaseURLPrefix()` + `feedAssetURL(arch,target,ext)`
  in `FreedomTechFeed/packages` — no index path, no `<arch>/packages`, no `src/gz`.
- The producer agrees: index generation was proposed once and **closed unmerged** —
  `FreedomTechFeed/feed` PR #1, `state=closed merged_at=null` (2026-09-08).
- The Gen-1 feed repo is stale and not the publisher: `main` @ `24fb8216` (2026-09-08) pins
  `PKG_VERSION:=0.5.0` while `v0.6.0` ships from `FreedomTechFeed/packages`; PR #11 ("point the
  standalone feed at the canonical recipe") is still **open** and its body says the sync is
  "**deliberately not executed here**".

So "feed" in this project = a naming convention for stable per-asset URLs. Any downstream card
that expects to write a `src/gz` line is working against a mechanism that has never existed.

---

## Answer set for the downstream decision card

1. **Is the Phase-2 feed live?** No.
2. **Is anything published?** Yes — `FreedomTechFeed/packages` `v0.6.0`, 2026-10-21, 26 assets,
   signed manifest verified this run.
3. **Candidate base URL:** `https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/`
   (asset-only, immutable per tag, digest-verified). Two minors *behind* the wizard's current
   `felixfelix-bot` pins (`v0.7.0-alpha10` / `v0.6.1-post-merge`), so adopting it is a downgrade
   unless the fork's assets are promoted upstream first.
4. **Rejected candidates:** `releases.tollgate.me` (stale unsigned index, key 404),
   `freedomtechfeed.github.io` (404, no Pages), `blossom2.orangesync.tech` (up, no feed paths),
   `blossom.net4sats.com` / `feed.net4sats.com` / `*.tollgate.me` (NO-DNS).

---

## Evidence gaps (stated, not papered over)

- **No authoritative definition of "Phase-2 feed" exists in any repo.** The reading used here is
  the one the feed project's own doc and the two sibling cards support; a different reading is
  possible but undocumented anywhere in the fleet.
- **No CI history was checked for a removed/renamed index publisher.** What was checked: 0
  releases/0 tags on `FreedomTechFeed/feed`, index PR #1 closed unmerged, no `gh-pages` branch,
  Pages API 404 on both repos. A publisher that once existed and was deleted would not show up on
  the live surface.
- `scripts/tollgate-release-url.sh --check` was **not** re-run; the equivalent proof was obtained
  directly (full download + digest compare + signature verify) and is stronger than the script's
  own check.
- The probe matrix covers the path *shapes* the project documents (flat, arch-only, arch-target)
  on the hosts it names. It does not enumerate arbitrary unknown hosts.

---

## Files touched

- `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-phase2-feed-status.md` — **the deliverable**
  (new, 18 KB). Outside any git repo; cannot be committed.
- `REPORT-t_855cbaac.md` (this file) — new, repo-tracked.
- `PROGRESS-t_855cbaac.md` — new, repo-tracked.
- `PROGRESS.md` — appended one section (shared file, was clean at append time).
- **`REPORT.md` deliberately NOT touched** — sibling card `t_wz15-part-b` had uncommitted edits in
  it during this run (`git status` showed `M REPORT.md`); appending and committing would have
  swept their work into this commit.
- No repo source file modified; no build/test state changed (`go` not run — nothing to run).
