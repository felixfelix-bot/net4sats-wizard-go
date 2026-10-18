# REPORT — t_1b1cf42f: Resolve GL-MT6000 feed base URL and exact `src/gz` line

**Date:** 2026-10-16
**Card:** `t_1b1cf42f` (board net4sats-mvp-v2, status triage)
**Repo under work:** `felixfelix-bot/net4sats-wizard-go` (worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`, detached HEAD `5dd47e3`)

**Result: STOPPED — the published artifacts DO NOT CONTAIN A FEED.**

The card's own escape hatch applies verbatim:

> Acceptance criteria: … If the artifacts do not contain a aarch64_cortex-a53 feed,
> STOP and report to the manager rather than substituting an architecture.

The aarch64_cortex-a53 **package** exists. There is no **feed** (no package index, no
`src/gz`-addressable directory) in the published artifacts. So there is no feed BASE URL and
no `src/gz` line to capture — and I have not invented one, nor substituted an architecture.

---

## 1. Direct answer to the card's deliverable

| asked for | answer |
|---|---|
| exact `src/gz tollgate <FEED_BASE_URL>/<arch>/packages` line | **DOES NOT EXIST** — no feed index is published (evidence §3) |
| release tag the feed is pinned to | n/a for a feed line. The relevant release is **`v0.6.0-rc6`** (2026-10-15T19:50:37Z, `isLatest`) |
| raw curl output proving the index is fetchable | **no index is fetchable — 404 at every candidate path** (§3.3). Raw 2xx + digest evidence for the actual artifact mechanism is in §4.3 |

A feed line cannot be produced from these artifacts because no index file exists to serve it.
Substituting one would be a fabricated URL, which the card forbids in two places.

---

## 2. What "the published release artifacts" actually are

The project's release lane is **`FreedomTechFeed/packages`** (confirmed as the source of
truth by the predecessor card `t_53879c67` and by the repo's own docs). Latest release:

```
$ gh release list -R FreedomTechFeed/packages -L 5 --json tagName,createdAt,isLatest,isDraft
[{"createdAt":"2026-10-15T18:45:10Z","isDraft":false,"isLatest":true, ... "tagName":"v0.6.0-rc6"},
 {"createdAt":"2026-10-14T11:22:08Z","isDraft":false,"isLatest":false,"tagName":"v0.6.0-rc5"}, ...]

$ gh api repos/FreedomTechFeed/packages/releases --paginate -q '.[] | select(.draft==true) | .tag_name'
(no output — no drafts)
```

**Release tag this work is pinned to: `v0.6.0-rc6`.**
Release URL: https://github.com/FreedomTechFeed/packages/releases/tag/v0.6.0-rc6

---

## 3. The decisive evidence: there is no feed index

### 3.1 The full asset list is FLAT — 26 assets, no index, no arch directory

```
$ gh release view v0.6.0-rc6 -R FreedomTechFeed/packages --json assets -q '.assets[].name'
SHA256SUMS
SHA256SUMS.sig
tollgate-wrt-0.6.0_rc6-aarch64_cortex-a53-mediatek-filogic-offline.tar.gz
tollgate-wrt-0.6.0_rc6-aarch64_cortex-a72-bcm27xx-bcm2711-offline.tar.gz
tollgate-wrt-0.6.0_rc6-arm_cortex-a7_neon-vfpv4-bcm27xx-bcm2709-offline.tar.gz
tollgate-wrt-0.6.0_rc6-mips64_octeonplus-octeon-generic-offline.tar.gz
tollgate-wrt-0.6.0_rc6-mipsel_24kc-ramips-mt7621-offline.tar.gz
tollgate-wrt-0.6.0_rc6-mips_24kc-ath79-generic-offline.tar.gz
tollgate-wrt-0.6.0_rc6-mips_24kc-ath79-nand-offline.tar.gz
tollgate-wrt-0.6.0_rc6-x86_64-x86-64-offline.tar.gz
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a72-bcm27xx-bcm2711.apk
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a72-bcm27xx-bcm2711.ipk
tollgate-wrt_0.6.0_rc6_arm_cortex-a7_neon-vfpv4-bcm27xx-bcm2709.apk
tollgate-wrt_0.6.0_rc6_arm_cortex-a7_neon-vfpv4-bcm27xx-bcm2709.ipk
tollgate-wrt_0.6.0_rc6_mips64_octeonplus-octeon-generic.apk
tollgate-wrt_0.6.0_rc6_mips64_octeonplus-octeon-generic.ipk
tollgate-wrt_0.6.0_rc6_mipsel_24kc-ramips-mt7621.apk
tollgate-wrt_0.6.0_rc6_mipsel_24kc-ramips-mt7621.ipk
tollgate-wrt_0.6.0_rc6_mips_24kc-ath79-generic.apk
tollgate-wrt_0.6.0_rc6_mips_24kc-ath79-generic.ipk
tollgate-wrt_0.6.0_rc6_mips_24kc-ath79-nand.apk
tollgate-wrt_0.6.0_rc6_mips_24kc-ath79-nand.ipk
tollgate-wrt_0.6.0_rc6_x86_64-x86-64.apk
tollgate-wrt_0.6.0_rc6_x86_64-x86-64.ipk

$ gh release view v0.6.0-rc6 -R FreedomTechFeed/packages --json assets -q '.assets | length'
26
```

Every name is a flat file. There is no `Packages`, no `Packages.gz`, no `Packages.sig`, and
no `<arch>/` prefix anywhere in the asset set. The `aarch64_cortex-a53` variant is present —
as two **direct package files**, not as an arch directory containing an index.

Note the asset name is `<arch>-<target>` qualified (`aarch64_cortex-a53-mediatek-filogic`),
which matters for §5.

### 3.2 The publish workflow says outright that it publishes flat per-arch assets

`.github/workflows/release-publish.yml` (master), header:

```
name: Publish per-arch tollgate-wrt release assets

# Builds tollgate-wrt for every arch the wizard can select (see
# net4sats-wizard-go arch.go) and uploads the resulting .apk / .ipk as
# release assets with deterministic, TARGET-QUALIFIED names:
#
#   tollgate-wrt_<PKG_VERSION>_<arch>-<target>.apk
#   tollgate-wrt_<PKG_VERSION>_<arch>-<target>.ipk
...
# `.release.target` and builds the URL with feedAssetURL(arch, target, ext),
# so the wizard's tollgateArchAssets map (arch.go) points at these stable URLs,
# so the correct binary is fetched for the detected router arch AND target.
```

The phrase "**feedAssetURL**" is the tell: the "feed" in this project's vocabulary resolves to
a **per-asset URL**, not an opkg/apk index. No index is built or uploaded.

### 3.3 The decisive negative: every candidate index path returns 404

```
$ for u in \
    ".../releases/download/v0.6.0-rc6/Packages" \
    ".../releases/download/v0.6.0-rc6/Packages.gz" \
    ".../releases/download/v0.6.0-rc6/aarch64_cortex-a53/Packages"; do
    curl -s -o /dev/null -w "%{http_code}  $u\n" -L "$u"; done

404  https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/Packages
404  https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/Packages
404  https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/Packages.gz
```

This is the raw curl output the card asked for, and it is negative. There is nothing at any
path shaped like `<FEED_BASE_URL>/<arch>/packages`, which is exactly the shape the card's
`src/gz` template presumes. **There is no `curl -sI <url>/Packages` that returns 2xx**, so the
card's step-3 sanity check cannot be satisfied — by anything, not just by me.

### 3.4 Nothing in the org publishes a `src/gz` feed either

```
$ gh search repos 'freedomtech feed'
(no output)

$ gh search code 'src/gz tollgate'          → 0 results
$ gh search code 'Packages.gz' --owner FreedomTechFeed → 0 results

$ gh api repos/FreedomTechFeed/packages/branches --paginate -q '.[].name' | grep -i -E 'feed|gz|index|publish|static'
feat/feed-publish-per-arch
feat/release-sha256sums-manifest
feed60-resolve
pr/feed-gate-naming
pr/feed-gates-ci
pr/publish-without-offline-bundle
pr/release-matrix-club-archs
```

The branch names confirm the *lane* is per-arch asset publication / SHA256SUMS manifests.
`pr/publish-without-offline-bundle` is the strongest existing signal: a branch exists about
publishing *without* the bundle — i.e. the per-arch-asset lane, still not an index feed.

### 3.5 The arch is NOT substituted — it is confirmed present, then correctly rejected as a feed

Per the card, I checked that aarch64_cortex-a53 exists before concluding, so this is not a
missing-architecture case:

- arch `aarch64_cortex-a53`, target `mediatek-filogic` — **present** as `.apk` and `.ipk` and
  as an offline bundle.
- `TARGETS.md` lists it as a **validated** pair: *"GL-MT3000, GL-MT6000, Cudy WR3000 |
  c03rad0r | 2026-10-15 | rc6 | MT3000 (apk/25.12.5) + MT6000 (opkg/24.10.2) upgraded in
  place to `0.6.0_rc6-r1`"*.

So the correct arch is embedded in the artifact names. The blocker is purely that no index
artifact exists to be addressed by a `src/gz` line.

---

## 4. What FE2b / the manager must know: the real resolution mechanism

### 4.1 The project's own resolver, and why it is not a feed

`TARGETS.md` documents the resolution path:

```
scripts/tollgate-release-url.sh v0.6.0-rc6 <arch> <target>          # apk lane
scripts/tollgate-release-url.sh 0.6.0_rc6 <arch> <target> ipk       # opkg lane
```

That script's own header states the design intent and the failure modes it exists to close:

```
# tollgate-release-url.sh -- resolve the tollgate-wrt download URL for a release.
# The release TAG and the package asset name are NOT the same string:
#     tag   : v0.6.0-rc6            (leading 'v', hyphens)
#     asset : tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
...
# This script needs a POSIX shell and nothing else. It never prints a URL it
# cannot construct: empty or malformed input is a hard error, not a warning.
#
# Exit codes: 0 ok | 2 usage | 3 empty/malformed input | 4 --check failed.
```

The script is the *documented, supported* way to resolve the artifact. It emits a **direct
download URL**, never an index. There is no `src/gz` in it because there is no feed.

### 4.2 The A-side of the bridge is not in the published artifacts

Per `docs/per-arch-release-assets.md`:

> The feed publishes tollgate-wrt packages for every router architecture the wizard can
> select, at stable, deterministic URLs. This is the Phase 2 prerequisite for "test the feed
> via the installer": the wizard's `tollgateArchAssets` map (in `net4sats-wizard-go`
> `arch.go`) points at these URLs so the correct binary is fetched for the detected router
> arch.

So the intended consumer is **the wizard fetching an individual package URL** — this repo's
own deploy path. The "feed" is a *naming convention for stable per-asset URLs*, not an apk/opkg
repository. Anything expecting `src/gz` is off the published design.

### 4.3 Live proof the direct-URL mechanism works (the alternative deliverable)

I ran the project's own self-gating check, which requires the asset to be listed in the
release's signed `SHA256SUMS` and to match on digest:

```
$ sh scripts/tollgate-release-url.sh --check v0.6.0-rc6 aarch64_cortex-a53 mediatek-filogic apk
OK  https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
    sha256 f59311148130ad483d1287854104dd4e9916a306e365d44eac35ccbb0b38c99f
    (matches SHA256SUMS for v0.6.0-rc6; authenticate the manifest with SHA256SUMS.sig before trusting it)
```

**exit code 0.** The ipk lane resolves to:

```
https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
```

A negative control confirms the script refuses rather than guesses: pointing it at a
near-miss arch/target (`mediateki-filogic`) produced

```
tollgate-release-url: asset is not listed in SHA256SUMS for v0.6.0-rc6: tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediateki-filogic.apk
```

— a hard failure, not a 404 URL. So the 2xx-equivalent proof the card asked for exists for the
**artifact**, and the corresponding proof for an **index** does not, because there is no index.

### 4.4 The target is load-bearing — a `src/gz` line could not express it anyway

From `docs/per-arch-release-assets.md`:

> The target slot is **load-bearing, not decoration**: one CPU arch can serve several OpenWrt
> targets. `mips_24kc` is both `ath79/generic` … and `ath79/nand` … so an arch-only name
> (`tollgate-wrt_<ver>_mips_24kc.apk`) made those two release rows **collide on a single asset
> name** … the installer can resolve the asset for the target the router actually reports
> (`internal/app/arch.go`, `feedAssetURL(arch, target, ext)`).

The card's template (`<FEED_BASE_URL>/<arch>/packages`) carries **no target slot**. Even if an
index feed were built later, an arch-only path would be ambiguous for `mips_24kc` and violates
the invariant this project explicitly hardened against. Any real feed path would have to be
`<BASE>/<arch>-<target>/packages`, which is a design decision beyond this card's remit.

---

## 5. Precondition check (context, not this card's deliverable)

`t_1b1cf42f`'s sibling `t_18d15abe` requires FE2a merged before FE2b proceeds. Observed:

```
sqlite3 .../net4sats-mvp-v2/kanban.db \
  "SELECT id,status,title FROM tasks WHERE title LIKE '%FE2%'"
t_df4304d5|todo   |FE2a: Feed PR#1 rebase + artifact-to-release bridge + arch tuples
t_1b1cf42f|triage |Resolve GL-MT6000 feed base URL and exact src/gz line
t_53879c67|blocked|FE2b: Feed consumption install on GL-MT6000
```

**FE2a (`t_df4304d5`) is still `todo`, not merged.** Its title — "artifact-to-release bridge" —
is precisely the bridge this card was sent to find. That is the structural explanation for the
negative finding: the bridge does not exist yet because the card that builds it has not run.
(The plan file at `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` does not mention FE2a,
FE2b, or a feed at all; it is an unrelated Q3 dev-split gate. Do not look there for FE2b's
definition of done.)

---

## 6. Brief / environment premises that were wrong (flagged, per fleet guidance)

1. **"Repo: net4sats-mvp-v2" is the board name, not the repo.** This worktree's origin is
   `felixfelix-bot/net4sats-wizard-go` (detached HEAD `5dd47e3`), which is why `arch.go` — the
   `tollgateArchAssets` file the feed docs reference — is not present here. Predecessor cards
   `t_00ee52a3` and `t_53879c67` hit the same mismatch. `arch.go`, `feedAssetURL`, and
   `tollgateArchAssets` are therefore **not verifiable from this worktree**; I did not claim to
   read them. Everything above is sourced from `FreedomTechFeed/packages` (live) and this repo's
   `deploy.go`, which *is* present.
2. **Card `t_1b1cf42f`'s `workspace_path` names workspace `t_53879c67`**, i.e. the blocked
   FE2b baseline card's workspace.
3. The surviving pre-baseline card `t_53879c67` is **blocked on device auth** and recorded that
   a *third actor* was mid-`apk` on the GL-MT6000. Nothing here required device access, so that
   contention did not affect this task — but FE2b's consumption step is still gated on it.

---

## 7. Honesty ledger

- I did **not** produce a `src/gz` line. No feed exists to derive one from; producing one would
  be fabrication, which the card forbids. This is the card's own STOP condition, not a partial
  result.
- I did **not** substitute an architecture. `aarch64_cortex-a53` was confirmed present in the
  artifacts (§3.5) before being rejected as a *feed*, which is a different question.
- I did **not** touch the device, `deploy.go`, or any other repo file. This run is read-only
  against the world; the only writes are this task's two records.
- The `--check` output in §4.3 is real command output, exit 0. The 404s in §3.3 are real.
- I could not verify `internal/app/arch.go` / `feedAssetURL` first-hand (wrong checkout, §6.1);
  those are quoted from the feed's own docs and labelled as such.
- Predecessor reports were read before working, so this report does not re-derive their
  findings as if new.

## 8. What the manager should decide

FE2b cannot receive a `src/gz` line, because the published artifacts contain no feed index.
Three coherent paths, in the manager's court:

1. **Drop the feed line from FE2b** and let FE2b consume the direct, digest-verified package
   URL that §4.3 proves works (the mechanism `docs/per-arch-release-assets.md` already names as
   intended). Cheapest; matches the shipped design.
2. **Have FE2a (`t_df4304d5`) publish a real index** — and decide the arch/target path shape,
   since an arch-only path is ambiguous for `mips_24kc` (§4.4). This makes `src/gz` possible
   later but is new work, not a lookup.
3. **Re-scope this card** if a feed index was expected from a different publisher than
   `FreedomTechFeed/packages` — I found no such publisher (`gh search repos`/`code` empty, §3.4),
   so I could not have found it on my own.

Whichever is chosen, the feed BASE URL cannot be resolved today because no index exists to
resolve it from.

---

## 9. Delivery — commit, push, and a convention conflict

**Commit:** `a24779226c347fa115dec58fc230b929f10009b4` — exactly 2 files
(`PROGRESS-t_1b1cf42f.md`, `REPORT-t_1b1cf42f.md`), 392 insertions, rc=0.
Committed with an explicit pathspec (`git commit -F <msg> -- <paths>`) so the shared tree's
sibling-card files could not be swept; the untracked `verify/wz16-e2e-gate/` belongs to another
card and was left untouched.

**Push:** `git push origin HEAD:refs/heads/worker-tollgate/t_1b1cf42f` → `* [new branch]`, rc=0.

**Server-side verification** (the exit code alone is not treated as proof):

```
$ git rev-parse HEAD
a24779226c347fa115dec58fc230b929f10009b4
$ git ls-remote origin refs/heads/worker-tollgate/t_1b1cf42f
a24779226c347fa115dec58fc230b929f10009b4	refs/heads/worker-tollgate/t_1b1cf42f
$ git ls-tree --name-only origin/worker-tollgate/t_1b1cf42f -- PROGRESS-t_1b1cf42f.md REPORT-t_1b1cf42f.md
PROGRESS-t_1b1cf42f.md
REPORT-t_1b1cf42f.md
```

Remote tip equals local HEAD and the tree on the server-side ref contains both records.

### The conflict I had to resolve to deliver at all

`git add` was initially **refused** — both files are ignored by the linked worktree's
`.git/info/exclude` (`/home/c08r4d0r/wz16-verify/.git/info/exclude`):

```
# t_0f285464 offload reports (local-only, never committed)
/PROGRESS.md
/REPORT.md
/PROGRESS-t_*.md
/REPORT-t_*.md
```

Taken at face value that rule forbids exactly the artifact the brief requires
("write your FULL report to REPORT.md … commit AND push after each verified milestone").
I did not resolve it by assumption. I checked the precedent on **this** worktree:

```
$ git ls-files | grep -E "PROGRESS-t_|REPORT-t_"
PROGRESS-t_53879c67.md
REPORT-t_53879c67.md
$ git log --all --oneline --source -- 'REPORT-t_53879c67.md'
24020f9  refs/remotes/origin/worker-tollgate/t_53879c67  docs(fe2b): record GL-MT6000 pre-baseline capture as BLOCKED
5dd47e3  refs/remotes/origin/worker-tollgate/t_53879c67  docs(fe2b): record markdown-table secret-gate false positive
```

The **immediately adjacent card on this exact chain** (the GL-MT6000 FE2b card `t_53879c67`)
has its task records tracked and pushed to its own `worker-tollgate/<task>` branch. I followed
that proven precedent: `git add -f`, **no `--no-verify`**, pushed to this card's own
`worker-tollgate/t_1b1cf42f` branch. Flagged for the manager: if these reports are genuinely
meant to be local-only, the stale exclude rule from `t_0f285464` is the thing to fix — not this
Fail, which matches the two most recent pushes on this branch family.

---

## 10. Reconciliation with the prior finding on this card — and a correction to the ASK itself

After I committed §1–§9 I found a **pre-existing comment on this same card** (`t_1b1cf42f`,
comment id 29860, author `worker-tollgate`) referencing an earlier FE2b RED baseline
(`t_877ad4ef`, 2026-10-12). It reached the *same* verdict by a *different* route and adds a
fact I had not tested. Both are recorded here.

### 10.1 Agreement — two independent negative proofs, from two different directions

| route | examined | result |
|---|---|---|
| the prior card `t_877ad4ef` | feed **host namespace**: `feed.tollgate.me`, `feeds.tollgate.me`, `download.tollgate.me`, `blossom.net4sats.com` | all **NO-DNS**; `releases.tollgate.me` resolves but `/feed/.../packages.adb` and `/alpha/...` return **404** |
| **this run** | the **release-artifact lane** `FreedomTechFeed/packages` (a different publisher, examined directly) | 26 flat assets, **no index**; every candidate index path 404 (§3.3) |

These do not overlap, so together they are a stronger result than either alone: **no feed index
exists in the host namespace or in the release-asset lane.** The prior comment's closing line —
*"This card's own acceptance criteria ('curl -sI <url> must return 2xx') cannot be satisfied
until one is hosted"* — is exactly what §3.3 measures from the release side.

### 10.2 Correction I am adopting: `src/gz` is the wrong FORM for this device

The prior card's live device evidence:

- The tollgate device on the bench runs **OpenWrt 25.12.5, `apk-tools 3.0.5`**, and declares
  `opkg: not found`.
- Its `/etc/apk/repositories.d/distfeeds.list` contains only **full `.../packages.adb` URLs** —
  never a bare directory base.
- Its packages line is therefore `<base>/<channel>/<arch>/packages.adb` in
  `/etc/apk/repositories.d/<name>.list`, **not** `src/gz <name> <dir>` (which is the
  opkg 24.10/23.05 form in `/etc/opkg/customfeeds.conf`).

I verified the apk-side corollary **from the release lane myself**, which the prior card could
not (it had no index to look for at that publisher):

```
$ for f in packages.adb aarch64_cortex-a53/packages.adb aarch64_cortex-a53-mediatek-filogic/packages.adb; do
    curl -s -o /dev/null -w "%{http_code}  $f\n" -L ".../releases/download/v0.6.0-rc6/$f"; done
404  packages.adb
404  aarch64_cortex-a53/packages.adb
404  aarch64_cortex-a53-mediatek-filogic/packages.adb
```

So **neither** index form exists — not the opkg `Packages`/`Packages.gz` (§3.3) and not the apk
`packages.adb` (above). The card's `src/gz` template is doubly wrong here: the form is opkg-only,
*and* no index of either form is published.

### 10.3 Unresolved conflict — which package manager does the GL-MT6000 actually use?

This matters because it decides which of the two forms FE2b would eventually use, and the two
sources disagree:

- **`TARGETS.md`** (validated-by table, dated 2026-10-15, validator `c03rad0r`) says:
  *"MT3000 (apk/25.12.5) + MT6000 (**opkg/24.10.2**) upgraded in place to `0.6.0_rc6-r1`"*.
- **The prior card `t_877ad4ef`** says the FE2b target device is **apk** (25.12.5,
  `apk-tools 3.0.5`) — but its own honest caveat states the device it actually inspected was a
  **GL-MT3000** (`BOARD=glinet,gl-mt3000`, hostname `tollgate-FDE7`), because **no MT6000 was on
  the bench** (`192.168.1.1` ping-dead, 22/80 closed). It then generalizes from the MT3000 to
  "the FE2b target device".

So the apk finding is real but was measured on the **MT3000**, not the MT6000. `TARGETS.md`
independently claims the MT6000 is **opkg/24.10.2**. **I cannot resolve this from here** — the
MT6000 is not on the bench and `root@192.168.1.1` denies every available credential
(predecessor `t_53879c67`). I am recording the conflict rather than picking a side.

Relevant supporting fact I did verify: **rc6 ships BOTH package formats** for this exact
arch/target, so the release itself does not decide the question and FE2b must:

```
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk    <- apk lane (25.12+)
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk    <- opkg lane (24.10)
```

### 10.4 Net effect on this card's ask

Both §3 and §10 agree on the deliverable: **there is no `src/gz` line to capture**, and no
`packages.adb` line either. The ask is unanswerable for two compounding reasons — the form is
wrong for an apk device (§10.2), and no index is published under any form (§3.3, §10.1). The
manager decision in §8 stands, now with a fourth option: settle the MT6000's package manager
first (§10.3), since that determines the form a resolved feed line would take.

---

# RUN 2 — 2026-10-17: re-derivation from the publisher's own code

**Result: STOPPED — the card's STOP condition applies.** No package index is published, in any
form, at any of the three bases the project's own repos document. There is therefore no
`src/gz tollgate <BASE>/<arch>/packages` line to capture, and no `curl -sI <url>/Packages` that
returns 2xx.

This run is not a restatement of RUN 1. RUN 1 established the absence empirically from the
release lane and the host namespace. This run found the **authoritative definition of the feed
in the publisher's own code**, which settles *why* the line cannot exist: in this project the
word "feed" denotes a per-asset URL prefix, not an index directory.

---

## 1. Why the requested line cannot exist — the decisive evidence

`OpenTollGate/tollgate-installer` is the canonical installer (cloned fresh this run,
`internal/app/arch.go`):

```
L106-107:  // feedRepoSlug is the repo whose CI builds and publishes the feed packages.
           feedRepoSlug = "FreedomTechFeed/packages"
L119:      feedReleaseTagDefault = "v0.6.0-rc6"
L320-324:  // feedReleaseURLPrefix is the fixed prefix of every feed asset URL, derived
           // from the effective tag (and therefore from the same single literal).
           func feedReleaseURLPrefix() string {
               return "https://github.com/" + feedRepoSlug + "/releases/download/" + feedReleaseTag + "/"
           }
L326-352:  // feedAssetURL builds the deterministic tollgate-wrt download URL for a
           // canonical OpenWrt arch tuple and file extension. The feed publishes every
           // arch it builds at a stable, predictable URL ...
           func feedAssetURL(arch, target, ext string) string {
               return feedReleaseURLPrefix() + "tollgate-wrt_" + feedPkgVersion() + "_" + arch + "-" + target + ext
           }
```

The project's single source of truth for "where is the feed" resolves to
`https://github.com/FreedomTechFeed/packages/releases/download/<tag>/` + an **individual
package file name**. There is no index path in it, no `<arch>/packages` directory, and no
`src/gz`. `arch_test.go` asserts exactly these strings
(`TestFeedReleasePublishesEachDerivedAssetName` checks the derived asset names against the live
release), so this is tested, current behaviour — not stale docs.

Consequence: the card's template `src/gz tollgate <FEED_BASE_URL>/<arch>/packages` presumes an
opkg index repository that this project has never published. The absence is **by design**, and
no amount of searching the artifacts can produce the line.

## 2. The card's step-3 sanity check, run in every form — all negative

Tag pinned: **v0.6.0-rc6** (`isLatest:true`, published 2026-10-15T18:45:10Z; no drafts).
Base probed: `https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6`

```
### B. the card's sanity check - index paths (all forms)
Packages                                       404
Packages.gz                                    404
aarch64_cortex-a53/Packages                    404
aarch64_cortex-a53/Packages.gz                 404
aarch64_cortex-a53/packages.adb                404
aarch64_cortex-a53/APKINDEX.tar.gz             404
aarch64_cortex-a53-mediatek-filogic/Packages   404
aarch64_cortex-a53-mediatek-filogic/packages.adb 404
```

Both index families are covered (opkg `Packages`/`Packages.gz`; apk `packages.adb` /
`APKINDEX.tar.gz`), and all three path shapes (flat, arch-only, arch-target). **Nothing returns
2xx, so the acceptance criterion "raw curl output proving the index is fetchable" cannot be
met — vacuously, because no index is published.** (The `src/gz` form is also opkg-only; an apk
device would use a `packages.adb` repositories line, which is equally 404 — see RUN 1 §10.2.)

## 3. The artifacts DO contain aarch64_cortex-a53 — as direct packages, 200 + digest-matched

This is the positive control proving I rejected a *feed*, not an *architecture*:

```
$ curl -sI -L .../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk | head -6
HTTP/2 302
location: https://release-assets.githubusercontent.com/...filename=tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk...
HTTP/2 200
content-length: 8591280

$ curl -sI -L .../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk | head -6
HTTP/2 302
location: https://release-assets.githubusercontent.com/...filename=tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk...
HTTP/2 200
content-length: 8630160
```

`signed` manifest `SHA256SUMS` (HTTP 200, 3010 B) lists both, so the fetch is verifiable:

```
f59311148130ad483d1287854104dd4e9916a306e365d44eac35ccbb0b38c99f  tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
3b3a85cc3ff5c4612a365d54ba56117600f079a416c3f3ce60d847d47a93e26a  tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
```

The release also ships the `-offline.tar.gz` bundle for the same arch/target. rc6 ships **both**
`.apk` and `.ipk`, so the release itself does not settle the MT6000's package-manager question
(RUN 1 §10.3 conflict stands unresolved).

## 4. Three documented feed bases exist in project repos — all dead or empty

RUN 1 searched the host namespace and the release lane. This run found three bases named in the
project's **own repositories**, which is a broader net than either:

| # | documented base | documented in | live result |
|---|---|---|---|
| 1 | `https://blossom.net4sats.com/feed/net4sats` | `net4sats/net4sats-feed` README — as BOTH the apk line (`>> /etc/apk/repositories`) and the opkg line (`src/gz net4sats ...`) | **NO-DNS** (`socket.gaierror`) — host does not resolve |
| 2 | `https://releases.tollgate.me/alpha/v0.5.0-beta2/{arch}` | that repo's `MANIFEST` (`src/gz net4sats <base>/{arch}`) | host resolves (Cloudflare 104.21.18.69) but every index path **404** |
| 3 | `https://freedomtechfeed.github.io/feed/<arch>/…` | `FreedomTechFeed/feed` PR #1 (gh-pages target) | **404**; repo has no `gh-pages` branch, Pages API 404 |

Base 1 is the most important finding of this run: the net4sats project **documents an exact
`src/gz net4sats https://blossom.net4sats.com/feed/net4sats` line** in its own README. That is
the closest thing in existence to the card's deliverable — and the host is dead. I record it as
the line the project *intends*, explicitly **not** as the card's deliverable, because it cannot
pass step 3 (it does not even resolve).

Raw:
```
$ python3 -c 'socket.getaddrinfo(...)'
blossom.net4sats.com     NO-DNS (gaierror)
feed.tollgate.me         NO-DNS (gaierror)
feeds.tollgate.me        NO-DNS (gaierror)
download.tollgate.me     NO-DNS (gaierror)
releases.tollgate.me     RESOLVES ['104.21.18.69', '172.67.180.198']
freedomtechfeed.github.io RESOLVES ['185.199.108.153', ...]
```

## 5. Root cause: both index lanes were built and then abandoned

- **`FreedomTechFeed/feed` PR #1** — *"feat(ci): auto-generate opkg+apk indices on every release
  (all arches)"*. Body describes exactly the missing artifact (`Packages.gz` + `APKINDEX.tar.gz`
  per arch, flat `<arch>/` layout, gh-pages publication). State: **CLOSED**, `mergedAt: null`,
  closed 2026-09-08. Its files are not in `main`: `scripts/` contains only
  `sync-from-upstream.sh`, and there is no `build-feed.yml` in `main`'s workflows.
- **`net4sats/net4sats-feed`** — `.github/workflows/build-package.yml` is a stub:
  ```
  # TODO: Add blossom upload + NIP-94 kind:1063 event
  ...
  - name: Update feed index
    run: |
      echo "Feed index update would go here"
  ```
- The release explorer at `releases.tollgate.me` is a React SPA whose own bundled install
  instructions are `opkg install /tmp/<file>` and `apk add --allow-untrusted /tmp/<file>` —
  i.e. per-file install, consistent with an asset lane and no feed.

So FE2a's "artifact-to-release bridge" (still `todo`, `t_df4304d5`) is not merely unmerged — its
index half was **written and closed unmerged twice**. That is the structural reason this card
has no answer, and it is the same conclusion RUN 1 reached by a different route.

## 6. The one LIVE feed — `src-git`, not `src/gz`, and not usable by FE2b

`FreedomTechFeed/feed` (main, live) README §"Using the feed":

```
echo "src-git tollgate https://github.com/FreedomTechFeed/feed.git" >> feeds.conf
./scripts/feeds update tollgate
./scripts/feeds install tollgate-wrt
```

This is a **firmware build-time source feed**: it compiles `tollgate-wrt` from the upstream
tarball via `golang-package.mk` during an image build. It cannot be consumed on-device by
opkg/apk, so it cannot install onto an already-flashed GL-MT6000, and it pins the stale
`PKG_VERSION:=0.5.0` / `PKG_RELEASE:=1`. It is a real feed, but the wrong *kind*, and it does not
satisfy FE2b. Recorded so the manager is not surprised to find a "live feed" that still doesn't
help.

## 7. What FE2b should actually consume (verified mechanism)

Derived from `feedAssetURL` (§1) and proven fetchable (§3), for arch `aarch64_cortex-a53`,
target `mediatek-filogic`, tag `v0.6.0-rc6`, both formats:

```
apk : https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk   (8,591,280 B)
ipk : https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk   (8,630,160 B)
sums: https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/SHA256SUMS        (3,010 B)
```
On-device: `apk add --allow-untrusted /tmp/<file>` or `opkg install /tmp/<file>` — the exact
commands the project's own release-explorer UI emits for this arch.

## 8. Premise corrections (carried forward, re-confirmed)

1. **"Repo: net4sats-mvp-v2" is the board name, not the repo.** This worktree's origin is
   `felixfelix-bot/net4sats-wizard-go`. `arch.go` is not in this worktree; I read it from a fresh
   clone of `OpenTollGate/tollgate-installer` instead, and cite the upstream path.
2. The card id lives in board DB `/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`
   (not the `boards/*.db` flat files, which are 0-byte, and not the manager profile DB). Card is
   `triage`; FE2a `t_df4304d5` is still `todo`; FE2b `t_53879c67` still `blocked`.
3. RUN 1's unresolved MT6000 package-manager conflict (opkg/24.10.2 per `TARGETS.md` vs
   apk/25.12.5 measured on an MT3000) is **not** resolved by this run; rc6 shipping both formats
   keeps it open. It decides which §7 command FE2b runs.

## 9. Honesty ledger

- I did **not** produce a `src/gz` line as the deliverable. The nearest documented line
  (`src/gz net4sats https://blossom.net4sats.com/feed/net4sats`, §4) is reported explicitly as
  **non-deliverable**: its host is NO-DNS.
- I did **not** substitute an architecture. `aarch64_cortex-a53` was confirmed present and
  fetchable (§3) *before* being rejected as an index.
- I did **not** invent a URL. Every URL in this report was either curl'd in this run or is a
  verbatim quote from a cloned repo / published release manifest.
- The 404/200/NO-DNS results are real output from this run; the sha256 values are quoted from
  the published SHA256SUMS, which I fetched (200) but did **not** re-verify against the
  downloaded bodies this run (RUN 1's `tollgate-release-url.sh --check` did verify the apk digest).
- Read-only: no device access, no repo files modified, no `--no-verify`. Scratch clones in /tmp.

## 10. Recommended manager decision (updated from RUN 1)

1. **Drop the feed line from FE2b** and have it consume the §7 asset URL + SHA256SUMS. This is
   the shipped design (§1) and is the only path with a working, digest-verifiable artifact today.
2. If a real on-device feed is genuinely required, the work is **FE2a**: revive
   `FreedomTechFeed/feed` PR #1 (closed unmerged) or the `net4sats-feed` stub, and **point it at a
   live host** — `blossom.net4sats.com` is NO-DNS, so the currently documented base cannot serve
   as-is. Note the path must carry the target slot (`<arch>-<target>`) or `mips_24kc`
   ath79-generic/nand collide (RUN 1 §4.4).
3. Settle the MT6000's package manager before writing any feed line, since that decides
   `src/gz` (opkg) vs a `packages.adb` repositories entry (apk).

---

# RUN 3 — 2026-10-18: re-derivation #3 (live), with the signature sub-clause closed

## 11.0 Verdict (unchanged, now signature-complete)

**No `src/gz` line can be produced.** The acceptance criteria's STOP clause applies verbatim:
the artifacts do not contain a `aarch64_cortex-a53` **feed**. I did not substitute an
architecture, and I did not substitute a publisher.

What this run adds over RUN 1/RUN 2 is the one acceptance-criteria sub-clause neither had tested:
*"and the matching package index/**signature** files."* Result: signature files **do** exist — and
their existence makes the negative *stronger*, not weaker, because they sign the flat assets
rather than an index (§11.3).

Everything below is raw output from this run, 2026-10-18.

## 11.1 Re-derived, not inherited: the release surface is unchanged

RUN 1 and RUN 2 both concluded NEGATIVE. I re-verified from scratch rather than citing them.

```
$ curl -sS -o /dev/null -w '%{http_code} %{url_effective}\n' -L https://github.com/FreedomTechFeed/packages/releases/latest
200 https://github.com/FreedomTechFeed/packages/releases/tag/v0.6.0-rc6

$ gh release list -R FreedomTechFeed/packages --limit 20
v0.6.0-rc6	Latest	v0.6.0-rc6	2026-10-15T19:50:37Z
v0.6.0-rc5		v0.6.0-rc5	2026-10-14T11:58:14Z
v0.6.0-rc4		v0.6.0-rc4	2026-10-13T10:51:53Z
... (18 more, all older)

$ gh release view v0.6.0-rc6 -R FreedomTechFeed/packages --json tagName,publishedAt,isDraft,isPrerelease,name -q .
{"isDraft":false,"isPrerelease":false,"name":"v0.6.0-rc6","publishedAt":"2026-10-15T19:50:37Z","tagName":"v0.6.0-rc6"}
```

**Release pin for any FE2b feed line: `v0.6.0-rc6`, published 2026-10-15T19:50:37Z.** No new tag
has appeared in the ~3 days since RUN 2, which is *why* the answer is stable rather than
merely repeated. (Note: `isLatest` is not a supported `--json` field on this gh version; the
`Latest` column of `gh release list` and the `releases/latest` redirect both establish it.)

## 11.2 The arch IS present — as packages, not as a feed

`gh release view ... --json assets -q '.assets[].name' | sort` → 26 assets. Every one is a flat
file. The aarch64_cortex-a53/mediatek-filogic row exists in three forms:

```
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
tollgate-wrt-0.6.0_rc6-aarch64_cortex-a53-mediatek-filogic-offline.tar.gz
```

plus `SHA256SUMS`, `SHA256SUMS.sig`, and the same apk/ipk/offline triple for 7 other
arch-target tuples. **There is no directory entry, no `Packages`, no `Packages.gz`, no
`Packages.sig`, no `packages.adb`.** This is the STOP condition: the arch is present as a
package and absent as a feed. I did not substitute an architecture.

## 11.3 The acceptance-criteria check, run verbatim — all 404

The card specifies `curl -sI <url>/Packages`. Raw headers, this run:

```
--- curl -sI .../releases/download/v0.6.0-rc6/Packages
HTTP/2 404
date: Sun, 18 Oct 2026 00:05:03 GMT
content-type: text/plain; charset=utf-8

--- curl -sI .../releases/download/v0.6.0-rc6/aarch64_cortex-a53/Packages
HTTP/2 404

--- curl -sI .../releases/download/v0.6.0-rc6/packages.adb
HTTP/2 404

--- curl -sI .../releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages.adb
HTTP/2 404

--- curl -sI .../releases/download/v0.6.0-rc6/Packages.sig
HTTP/2 404
```

Ten path shapes were probed in total this run (adding `Packages.gz`, `Packages.sig`,
`aarch64_cortex-a53/mediatek-filogic/{Packages,packages.adb}`, `aarch64_cortex-a53/packages/Packages`).
**Every index path in both index forms returned 404.** No `curl -sI` returned 2xx, so the
acceptance criterion "must return 2xx" is **unmet — and unmet for a structural reason**, not a
transient one (§11.6).

## 11.4 NEW: the signature files exist — and that closes the sub-clause *against* the card

The card asks to locate the feed "and the matching package index/signature files". The signature
files are real and fetchable:

```
200 8591280B  .../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
200 8630160B  .../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
200    3010B  .../SHA256SUMS
200     334B  .../SHA256SUMS.sig
200 9490009B  .../tollgate-wrt-0.6.0_rc6-aarch64_cortex-a53-mediatek-filogic-offline.tar.gz
```

`SHA256SUMS.sig` (334 B) is a genuine detached signature over `SHA256SUMS`. So **the release is
signed** — but what is signed is the checksum list of the **flat per-asset files**. There is no
`Packages.sig`, no `Packages.gz.sig`, no `packages.adb.sig`, because **there is no index to
sign** (all three 404, §11.3).

This is the decisive point for the sub-clause: a `src/gz` feed requires a *signed index*. A
signature over loose assets is not a feed signature and cannot be consumed by opkg or apk as a
repository. **The signature files exist and are correctly rejected as feed signatures.**

## 11.5 NEW: the offline bundle is the real consumption path, and it contains no index

`tollgate-wrt-0.6.0_rc6-aarch64_cortex-a53-mediatek-filogic-offline.tar.gz` → HTTP 200, 9490009 B.
Extracted:

```
$ find rc6b -type f | wc -l
44
$ find rc6b -iname '*packages*' -o -iname '*APKINDEX*' -o -iname '*.adb'
(no output)
```

44 files, **zero index files**. Its README states the bundle is *"38 packages: the tollgate-wrt
package + its closed dependency set"* fetched *"from the PUBLISHED apk indexes of OpenWrt 25.12.5
for `aarch64_cortex-a53` / `mediatek-filogic`"*, with `MANIFEST.sha256` for integrity and the
release's signed `SHA256SUMS` as the provenance chain.

That sentence is the answer to *why* no feed index is published: **the publisher resolves the
dependency closure at build time and ships the bytes**, so the device never needs an index. It
also states the limitation explicitly:

> The `.ipk` / OpenWrt <= 24.x (`opkg`) lane is NOT covered: this bundle is the apk-tools 3 lane
> only (OpenWrt 25.x).

and, in the runbook, that `--allow-untrusted` is required on the apk lane because the packages
carry no repository signature. Consistent: signed *asset sums*, unsigned *packages*, no *index*.

## 11.6 The only live index at this arch/target is OpenWrt's own — explicitly NOT substituted

```
https://downloads.openwrt.org/releases/25.12.5/targets/mediatek/filogic/packages/packages.adb
  → HTTP 200, 14417 B, magic bytes "ADBd" (a valid apk index)
  → grep -c -i tollgate → 0

https://downloads.openwrt.org/releases/25.12.5/targets/mediatek/filogic/packages/Packages
  → 404
https://downloads.openwrt.org/releases/25.12.5/targets/mediatek/filogic/packages/Packages.sig
  → 404
```

This is the one URL in the whole investigation that satisfies `curl -sI <url>/...` with a 2xx and
is a real index at the correct arch/target. It is **not** the feed, and I am not offering it:
wrong publisher, and its contents do not include `tollgate-wrt` (0 matches). Substituting it
would produce a `src/gz` line that installs nothing — exactly the failure the card's STOP clause
exists to prevent. It is recorded so the manager can see the alternative was tested and refused
on evidence rather than overlooked.

## 11.7 Why the negative is structural, not transient — both index lanes are still abandoned

```
$ gh api repos/FreedomTechFeed/packages/contents/.github/workflows -q '.[].name'
build-tollgate-sdk-image.yml  feed-gates.yml  llm-review.yml  multi-arch-test-build.yml
pkg-tarball-parity.yml  release-manifest-backfill.yml  release-publish.yml
release-tooling.yml  vendor-drift.yml  scripts

$ gh api 'search/code?q=repo:FreedomTechFeed/packages+packages.adb'   -q .total_count   → 0
$ gh api 'search/code?q=repo:FreedomTechFeed/packages+APKINDEX'       -q .total_count   → 0

$ gh pr view 1 -R FreedomTechFeed/feed --json state,mergedAt,closedAt,title -q .
{"closedAt":"2026-09-08T11:16:51Z","mergedAt":null,"state":"CLOSED",
 "title":"feat(ci): auto-generate opkg+apk indices on every release (all arches)"}

$ gh api repos/FreedomTechFeed/feed/branches/ci/auto-feed-indices -q '.commit.commit.author.date'
2026-07-05T16:27:53Z
$ gh api repos/FreedomTechFeed/feed/commits/main -q '.commit.author.date'
2026-09-08T10:46:33Z
```

`release-publish.yml`'s own header describes its output as *"release assets with deterministic,
TARGET-QUALIFIED names"* — flat assets by design. The only index-generating PR was closed
unmerged; its branch has not moved since 2026-07-05. **Nothing in the publisher's own CI can
produce a feed index**, so re-running this card will keep returning the same verdict until the
index lane is built.

Host namespace re-probed this run, all still dead:

```
blossom.net4sats.com      → NO-DNS (curl exit 6: Could not resolve host)
feed.tollgate.me          → NO-DNS
feeds.tollgate.me         → NO-DNS
download.tollgate.me      → NO-DNS
feed.net4sats.com         → NO-DNS
releases.tollgate.me      → A=172.67.180.198 (Cloudflare); serves only a React SPA
                            ("TollGate Release Explorer"); 12 candidate index paths → all 404
freedomtechfeed.github.io/feed/aarch64_cortex-a53/packages.adb → 404
                            (FreedomTechFeed/feed has no gh-pages branch; Pages API → 404)
```

Also re-confirmed from the docs of record: `net4sats/net4sats-feed` README documents
`blossom.net4sats.com/feed/net4sats` as **both** the apk and opkg base, and that host does not
resolve; its `MANIFEST` documents `releases.tollgate.me/alpha/v0.5.0-beta2/{arch}` as the feed
base, and every path under it 404s.

## 11.8 What FE2b should do instead (unchanged recommendation, now with the bundle)

1. **Drop the feed line.** Consume the digest-verified flat asset, or the offline bundle:
   `.../releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk`
   (or `.ipk` for the opkg lane), verifying against `SHA256SUMS` / `SHA256SUMS.sig`.
2. **For a no-uplink install**, use the offline bundle (44 files, closed closure, `MANIFEST.sha256`).
3. **If a real `src/gz` feed is required**, the index lane must be built first (the closed
   `FreedomTechFeed/feed` PR #1 is the work item), and the path must be arch **and target**
   qualified — `mips_24kc` covers both `ath79-generic` and `ath79-nand`, so an arch-only path
   collides. A `<FEED_BASE_URL>/<arch>/packages` template cannot express this project's tuples.
4. **Note the `src/gz` form is opkg-only.** The FE2b bench device is apk/25.12.5 (per the prior
   card's live evidence), where the repositories line is a full `.../packages.adb` URL in
   `/etc/apk/repositories.d/`, not `src/gz <name> <dir>`.

## 11.9 Honesty notes for this run

- I did **not** produce a `src/gz` line, because no index exists to derive one from.
- I did **not** substitute an architecture: aarch64_cortex-a53 was confirmed present as a
  *package* (§11.2) before being rejected as a *feed*.
- I did **not** substitute a publisher: OpenWrt's live index was found, then explicitly refused
  on evidence (§11.6).
- I did **not** substitute a *form* either: both the opkg (`Packages`/`Packages.gz`) and apk
  (`packages.adb`) index forms were probed at every candidate path, in the release lane and in
  all four documented host bases.
- The only 2xx responses this run were for flat assets, signature files, the offline bundle, and
  OpenWrt's own index. **No tollgate feed index returned 2xx anywhere.**
- Read-only: no repo files modified outside this card's two record files; no device access; no
  `--no-verify`. Scratch clones and downloads confined to `/tmp`.
