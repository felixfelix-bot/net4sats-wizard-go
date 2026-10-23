# REPORT — FE2b: resolve GL-MT6000 feed base URL and exact src/gz line

**Date:** 2026-10-23
**Repo:** `felixfelix-bot/net4sats-wizard-go` (worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`,
branch `worker-heavy/t_gate8-19b1b1c6`)
**Card.** FE2b, step 1 of 3. Read-only. No device was modified.
**Predecessor.** `REPORT-fe2b-feed-mt6000.md` (2026-10-21) — its conclusion is **superseded**
in one important respect (see §2). Its negative finding for the *current* release stands.

---

## 0. Verdict

**The exact `src/gz` line exists, but it belongs to a superseded release, and it would
not install the package the router is supposed to run.**

Two different things are published, and the brief's premise conflates them:

| surface | what it is | state for MT6000 `aarch64_cortex-a53` |
|---|---|---|
| `releases.tollgate.me/alpha/<tag>/` | a real **opkg repository**: `Packages.gz` + `Packages` + the `.ipk` assets they name | only for **`v0.5.0-alpha3`** |
| `github.com/FreedomTechFeed/packages/releases/download/<tag>/` | the **current** feed (v0.6.0): flat per-arch assets, **no index** | flat `.ipk`/`.apk` only, 26 assets, zero index |

So there is **no `src/gz` line that installs the current release**, and the one line that
does exist installs a stale build. Per the card's STOP clause I did not substitute an
architecture — and I am not proposing to substitute a *release* either. Reporting to the
manager, with the exact line FE2b would use if the stale pin is acceptable.

### The exact line (verbatim), for the only published feed

```
src/gz tollgate https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/packages
```

**Release/tag:** `v0.5.0-alpha3` (channel path `alpha/`, published 2026-07-01).
**Package it resolves:** `tollgate-wrt v0.5.0-alpha3` (Depends: nodogsplash, luci, jq).

The arch path is `<ARCH>/` and the index lives **directly** in it (`.../aarch64_cortex-a53/Packages.gz`).
There is no `/packages/` subdirectory — the literal template in the brief
(`<BASE>/<arch>/packages`) does **not** match the real artifact layout, so the trailing
`/packages` segment is wrong here. That is a layout correction, not an arch substitution.

---

## 1. Raw curl proof the index is fetchable

```
$ curl -sI -L https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.gz
HTTP/2 200
date: Fri, 23 Oct 2026 16:46:22 GMT
content-type: application/gzip
content-length: 439
server: cloudflare
last-modified: Wed, 01 Jul 2026 19:52:40 GMT
access-control-allow-origin: *
etag: "6a457008-1b7"
```

Uncompressed index body (complete, 3 stanzas):

```
Package: tollgate-wrt
Version: v0.5.0-alpha3
Depends: nodogsplash, luci, jq
Provides: nodogsplash-files
Replaces: nodogsplash, base-files
Architecture: aarch64_cortex-a53
Maintainer: TollGate <tollgate@tollgate.me>
License: CC0-1.0
Description: TollGate Basic Module for OpenWrt
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
Size: 5321054
SHA256sum: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3
... (same package repeated for -upx-ultra-brute.apk and .apk)
```

**The index is genuine, not an SPA catch-all.** `releases.tollgate.me` is a React SPA
whose router answers 200 with `text/html` for arbitrary deep paths, so a bare 200 there
proves nothing. Three independent checks rule that out for this URL:

1. `content-type: application/gzip` (SPA fallback returns `text/html`).
2. `file` → `gzip compressed data, was "Packages"`; first bytes `037 213` (magic).
3. The asset the index names is downloadable and its **SHA256 matches the index stanza**:

```
$ curl -sSL -o idx.ipk https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
http=200 bytes=5321054
$ sha256sum idx.ipk
c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3  idx.ipk
   index declares: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3   <- match
```

A plain (uncompressed) `Packages` also exists beside it, `content-type: application/octet-stream`,
200 — so both `src/gz` (gz) and an uncompressed `src` variant are servable.

**No signature file:** `Packages.sig` → 404. The feed is unsigned; opkg on a stock image
will need `--force-signature`-style handling or a `check_signature` toggle.

---

## 2. Correction to the 2026-10-21 report

`REPORT-fe2b-feed-mt6000.md` §8 concluded "no feed index is published, in any form, at
any candidate base" after probing **15 path shapes**. That conclusion is **too strong**.
It is true for `github.com/FreedomTechFeed/packages/releases/download/<tag>/` — which is
where it probed — but the org also serves a **real opkg repository** on
`releases.tollgate.me`, and the prior card's candidate list did not include that host's
`alpha/<tag>/<arch>/` layout. Its own §8.1 quoted the feed-strategy doc describing
exactly this location as "the recommended option" and treated the `← NEW` labels as
proof it was unbuilt. For `v0.5.0-alpha3` it **was** built and is live (last-modified
2026-07-01). I verified rather than inherited, which is the only reason the line in §0 exists.

The prior card's *other* findings still hold and are re-confirmed below.

---

## 3. The current release publishes no index — re-derived today

`v0.6.0` is the newest feed release (published 2026-10-21T02:06:45Z, i.e. **after** the prior card ran).

```
$ gh release view v0.6.0 -R FreedomTechFeed/packages --json assets -q '.assets|length'
26
$ ... --json assets -q '.assets[].name' | grep -iE 'index|packages|adb|Packages\.gz'
NONE
```

Every index shape at the current tag → 404:

```
404  Packages
404  Packages.gz
404  Packages.sig
404  index.json
404  aarch64_cortex-a53/Packages
404  aarch64_cortex-a53/Packages.gz
404  aarch64_cortex-a53/packages/Packages
404  aarch64_cortex-a53/packages/Packages.gz
404  aarch64_cortex-a53/packages.adb
404  aarch64_cortex-a53/APKINDEX.tar.gz
404  packages.adb
200  SHA256SUMS
```

The literal URL a `src/gz` line against the current tag would make opkg fetch:

```
$ curl -o /dev/null -w '%{http_code}' -L .../releases/download/v0.6.0/aarch64_cortex-a53/packages/Packages.gz
404
```

The arch is still published — as a **flat asset**, which is how the installer consumes it:

```
$ curl -sI -L .../releases/download/v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk
HTTP/2 200
content-length: 8629996
```

**Why v0.6.0 has no index (root cause, confirmed in source):** the installer defines
"feed" as a flat-URL prefix, not an index directory —
`FreedomTechFeed/packages` is the slug and the asset URL is derived, so no index is ever
needed by its path:

```go
// tollgate-installer internal/app/arch.go:107,119,321-323,349-351
feedRepoSlug          = "FreedomTechFeed/packages"
feedReleaseTagDefault = "v0.6.0"
func feedReleaseURLPrefix() string { return "https://github.com/"+feedRepoSlug+"/releases/download/"+feedReleaseTag+"/" }
func feedAssetURL(arch, target, ext string) string { return feedReleaseURLPrefix()+... }
```

The index generator is still **unmerged**: `FreedomTechFeed/feed` branch
`ci/auto-feed-indices`, single commit `1569f78 2026-07-05 feat(ci): auto-generate opkg+apk
indices on every release`. `main` carries no `scripts/index-feed.sh` and no `docs/`.
Its strategy doc proposes `releases.tollgate.me/<channel>/<tag>/<arch>/Packages.gz` — the
layout that **is** populated for alpha3 and never again (alpha1/alpha2/alpha4/alpha5/
v0.5.0/alpha/0.6.0 tag dirs all 404).

I also probed the alternative hosts named in that doc — `net4sats.github.io/feed/...` → 404,
`feed.tollgate.me` / `downloads.tollgate.me` → DNS/connect failure (000). Legacy feed repos
(`OpenTollGate/tollgate-feed`, `custom-nostr-feed-decomissioned`) have no `gh-pages` and no
release assets; `FreedomTechFeed/feed` has **zero releases**.

---

## 4. Why the alpha3 line is dangerous for this target

The only fetchable index resolves a package **older than what is on the router**, and the
repo's own code treats that as a downgrade hazard:

- index version: `tollgate-wrt v0.5.0-alpha3`
- prior card measured installed on the MT6000: `tollgate-wrt 0.6.0_rc6-r1`
- index asset is a **`.ipk` (opkg/ar)**; the `.apk` stanzas name
  `-upx-ultra-brute.apk` variants that my probe found only as index text, not checked as fetchable
- `internal/app/arch.go:368+` documents the fallback dump `v0.5.0` as an *older* release
  and `refuseDowngrade` (deploy.go) exists specifically to stop a pinned feed pushing a
  lower version over an installed one.

So installing from this feed on a box already running 0.6.0_rc6 is a **downgrade**, and
its `Depends: nodogsplash, luci, jq` are resolved from OpenWrt's own feeds, not from this one.

Also note the alpha3 index predates the target-qualified asset naming
(`_<arch>-<target>`); its `Filename:` is arch-only.

---

## 5. Device-state caveat: the MT6000 is not at 192.168.1.1 anymore

The prior card reached a real MT6000 at `192.168.1.1`. It is **not there now**:

```
$ ssh root@192.168.1.1 'cat /tmp/sysinfo/board_name; . /etc/openwrt_release; echo "$DISTRIB_RELEASE $DISTRIB_ARCH"'
qemu-standard-pc-q35-ich9-2009
25.12.5 x86_64
```

That is a **QEMU x86_64 guest** (OpenWrt 25.12.5), matching the finding already recorded
in this repo's `REPORT-t_23b8b386.md`. `aarch64_cortex-a53` is therefore **not** confirmed
from the device this run. It is confirmed from upstream's own target metadata for the exact
board, which is the strongest available substitute:

```
$ curl -sS https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/profiles.json
target: mediatek/filogic | arch_packages: aarch64_cortex-a53
profiles.glinet_gl-mt6000.titles: [{'model': 'GL-MT6000', 'vendor': 'GL.iNet'}]  -> exists
```

MT6000 on 24.10.2 is **opkg**, so the `src/gz` form is the right *form* for that image
(consistent with the prior card's measurement). But nothing in this task touched a real
MT6000, so FE2b must not assume the box in front of it is one.

---

## 6. What FE2b should actually do

1. **Do not** add the alpha3 line expecting the current package. It is a downgrade and its
   index is older than the installed build.
2. If a `src/gz` line is required for the card's literal acceptance criteria, the line in
   §0 is the only real one; pin it with the tag, and treat it as a **v0.5.0-alpha3
   repository**, not as "the tollgate feed".
3. The supported path for the current release is the installer's flat-asset lane
   (`feedAssetURL`), i.e. download
   `.../v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk` and push it —
   which is exactly the mechanism a previous card's brief forbade. **That tension is a
   manager decision, not a worker one.** Until an index is published for `v0.6.0`, a
   `src/gz` line and the current release are mutually exclusive.
4. The durable fix is upstream: merge `FreedomTechFeed/feed@ci/auto-feed-indices` so CI
   emits `Packages.gz`/`packages.adb` per release.

---

## 7. Hygiene / deviations

- **Read-only honoured.** No repo file, no device config, no release was modified. Only
  `/tmp` scratch (`fe2bprobe/`, `inst/`, `fr/`, `fmain/`, `fidx/`, `tl_pkgs.gz`, `idx.ipk`, `prof.json`).
- **Files written:** `REPORT-t_fe2b-feed-url.md`, `PROGRESS-t_fe2b-feed-url.md` (task-scoped,
  because root `REPORT.md`/`PROGRESS.md` in this worktree belong to other cards — the shared
  `.git/info/exclude` outside this worktree blocks `/PROGRESS.md`, `/REPORT.md`,
  `/PROGRESS-t_*.md`, `/REPORT-t_*.md`, so these were added with `git add -f`).
- **No device mutation**, no creds read, no `--no-verify`, no history rewrite.
- **Not inherited:** every URL above was probed this run (2026-10-23); the prior card's
  verdict was deliberately re-tested rather than cited.

## 8. Acceptance criteria

| criterion | status |
|---|---|
| exact `src/gz` line, verbatim | **MET** for v0.5.0-alpha3 (§0); **does not exist** for v0.6.0 (§3) |
| release tag it belongs to | **MET** — `v0.5.0-alpha3`, `alpha/` channel |
| raw curl output proving index fetchable | **MET** — 200 `application/gzip`, §1, plus sha256 cross-check |
| do not substitute an architecture | **HONOURED** — aarch64_cortex-a53 confirmed from upstream profiles.json; no substitution |
| STOP and report if no aarch64 feed | **REPORTED** — no aarch64_cortex-a53 feed for the current release; stale one exists |
