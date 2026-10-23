# wz1.5 — official Tollgate artifact channel: DECISION

**Card:** `t_1df48589` (channel decision note + Q5 status)
**Worker:** `t_gate8-19b1b1c6`, profile `worker-tollgate`
**Date:** 2026-10-23 (probes run 21:30–21:40 UTC)
**Decision status:** DECIDED for the interim channel · **Q5 = OPEN**
**Re-verified:** 2026-10-24 by card `t_1df48589` (second, independent pass) — all live claims re-run,
one defect found and fixed (the composed artifact-filename template 404s; see §2.3a / §3 rule 1b).
Probe timestamps in this file are otherwise the 2026-10-23 run.
**Implementation status of the decision: BLOCKED** — see §6. The chosen channel publishes
no artifact for either version `deploy.go` installs today.

---

## 1. Decision

**Chosen channel (interim, documented): `releases.tollgate.me`.**

**Base URL pattern to implement:**

```
https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<openwrt-arch>/<artifact-file>
```

Instantiated for the wizard's two package formats (the `<arch>` segment is the *OpenWrt*
`DISTRIB_ARCH` the wizard already detects, e.g. `aarch64_cortex-a53` for the GL-MT6000 on
25.12.5; it is **not** the upstream target tuple — see §4):

```
https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<arch>/<artifact-file>
```

**`<artifact-file>` is not a derivable template — and the first draft of this note got that
wrong.** The chosen host's only live version proves it: its `.ipk` is named
`tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk`, i.e. it carries a
**build-flavour suffix** that a `tollgate-wrt_v<X.Y.Z>_<arch>.ipk` template does not produce.
The templated name is **404** on the one version that exists (§2.3), while the suffixed name is
200 / 5,321,054 B. The same tag *also* publishes an unsuffixed
`tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk` (7,153,725 B) next to a suffixed apk
(5,307,464 B), so the suffix is **per-artifact, not per-tag**. Consequence: **take
`<artifact-file>` verbatim from that version's own index `Filename:` field, never compose it**
(§3 rule 1b).

Rationale for choosing `releases.tollgate.me` and not the GitHub release lane: it is the only
host the project names that is a **Tollgate-owned release channel with a channel segment in the
path**. That segment is part of the publisher's own vocabulary, not an invention of this note —
the explorer bundle at `/static/js/main.ca30081d.js` (474,543 B, 200) declares

```
kr={STABLE:"stable",BETA:"beta",ALPHA:"alpha",DEV:"dev"}
Tr={channels:[kr.STABLE],products:["tollgate-os","tollgate-wrt","tollgate-module-basic-go"]}
```

and the live static tree uses `/alpha/` today. The GitHub lane is a per-asset upload directory
owned by an account, not a release channel, and it is the thing Q5 is asking about.

**Why not a feed/index line:** no opkg/apk index exists for any current release (§2). The chosen
pattern above resolves a **single artifact file** per URL, so it works for the wizard's
download-then-install path (which is what `deploy.go` does) and does **not** work as a
`src/gz`/`packages.adb` repositories line. Any later feed work is a separate decision.

---

## 2. Evidence

### 2.0 Input provenance — one named input does not exist

The task named two parent fragments. **Only one was present:**

```
$ ls -la /home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/
-rw-rw-r-- 1 c08r4d0r c08r4d0r 17507 Oct 23 21:00 schedule.md
-rw-rw-r-- 1 c08r4d0r c08r4d0r 18604 Oct 23 22:58 wz1.5-phase2-feed-status.md
-rw-rw-r-- ...  wz1.5-url-inventory.{md,part-a.md,part-b.md,audit-*.md}
-rw-rw-r-- ...  wz1.5-validation.md

$ read_file ~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-probe-raw.md
File not found: .../wz1.5-channel-probe-raw.md

$ find /home/c08r4d0r -name '*channel-probe*' -o -name '*probe-raw*'
(no output)
```

`wz1.5-channel-probe-raw.md` was **never written**. Its expected content ("raw probe evidence +
example versioned URLs") is therefore **re-derived in full, live, in this session** — every line
in §2.1–§2.6 below was produced by a command run here, none is quoted from a missing file or from
another card's report. `schedule.md` was read as instructed: it is the FELIX Q3 / dev-split gate
file (Q3 = B resolved 2026-10-20) and contains **nothing** about artifact channels — no input from
it is used.

### 2.1 Parent verdict carried forward

`wz1.5-phase2-feed-status.md` **VERDICT: NOT PUBLISHED** — no fetchable Phase-2 feed index at any
host the project names. That verdict is re-confirmed independently in §2.2 and §2.3, and it is the
branch of the decision rule this note follows.

### 2.2 The chosen channel's live surface — one stale version only

Index probe matrix (`curl -sSL -m 25 -o /dev/null -r 0-1023 -w 'http=%{http_code} bytes=%{size_download} type=%{content_type}'`):

```
https://releases.tollgate.me/                                          http=206 bytes=697 type=text/html; charset=utf-8
https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages        http=206 bytes=1322 type=application/octet-stream
https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.gz     http=206 bytes=439  type=application/gzip
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/Packages               http=404 bytes=697  type=text/html; charset=utf-8
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/Packages.gz            http=404 bytes=697
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/APKINDEX.tar.gz        http=404 bytes=697
https://releases.tollgate.me/alpha/v0.6.1/aarch64_cortex-a53/Packages               http=404 bytes=697
https://releases.tollgate.me/alpha/v0.5.0-beta2/aarch64_cortex-a53/Packages         http=404 bytes=697
https://releases.tollgate.me/alpha/v0.7.0-alpha10/aarch64_cortex-a53/Packages       http=404 bytes=697
```

Full tag matrix over `{v0.4.0, v0.5.0-alpha1, v0.5.0-alpha2, v0.5.0, v0.6.0-rc6,
v0.6.1-post-merge, v0.7.0-alpha10} × {Packages, Packages.gz, APKINDEX.tar.gz, packages.adb}`:
**28/28 = 404.** `v0.5.0-alpha3` is the only tag on this host with any index at all, and only for
`aarch64_cortex-a53` — the same four paths on `aarch64_cortex-a72-bcm27xx-bcm2711`,
`x86_64-x86-64` and `mipsel_24kc-ramips-mt7621` are **404**.

Channel and directory paths are dead too — the channel segment is not browsable:

```
https://releases.tollgate.me/alpha/      http=404 bytes=697
https://releases.tollgate.me/stable/     http=404 bytes=697
https://releases.tollgate.me/beta/       http=404 bytes=697
https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/   http=404 bytes=697
https://releases.tollgate.me/os/57e0f2468a17…?channel=stable          http=404 bytes=697
https://releases.tollgate.me/package/v0.7.0-alpha10                   http=404 bytes=697
```

So `releases.tollgate.me` is a **React SPA** — `<title>TollGate Release Explorer</title>`, bundle
`/static/js/main.ca30081d.js` (200, 474,543 B), PWA manifest at `/manifest.json` (206, 315 B) — whose
only data lives at deep static paths. There is no directory listing and no JSON API
(`/releases.json`, `/api/releases`, `/index.json`, `/data/releases.json`, `/os.json` → all 404).

### 2.3 What the one live version actually contains — and why it cannot be used

```
$ curl -sSI https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages
HTTP/2 200
content-type: application/octet-stream
server: cloudflare
last-modified: Wed, 01 Jul 2026 19:52:40 GMT

$ curl -sSL .../Packages | head -12
Package: tollgate-wrt
Version: v0.5.0-alpha3
Depends: nodogsplash, luci, jq
Architecture: aarch64_cortex-a53
Maintainer: TollGate <tollgate@tollgate.me>
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
Size: 5321054
SHA256sum: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3
```

Payload and index are consistent and digest-verified (this is a *real* artifact, not a placeholder):

```
$ curl -sSL -o a3.ipk .../tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
http=200 bytes=5321054
$ sha256sum a3.ipk
c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3  a3.ipk     # == index SHA256sum

$ curl -sSL -o APKINDEX.tar.gz .../APKINDEX.tar.gz   -> http=200 bytes=407
$ sha256sum APKINDEX.tar.gz
ca1ccf02fcbf57b01af98f4b670864dddeafddbc78f3c7504d50932bad39e206
$ tar -xzOf APKINDEX.tar.gz APKINDEX | head -8
P:tollgate-wrt / A:aarch64_cortex-a53 / V:v0.5.0-alpha3 / S:5307464
I:5307464 / T:Q1c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f
```

But it is **one minor release behind** the published line and **unsigned**:

```
.../aarch64_cortex-a53/Packages.sig          http=404 bytes=697
.../aarch64_cortex-a53/packages.adb          http=404 bytes=697
https://releases.tollgate.me/keys/tollgate-testing.pem   http=404 bytes=697
```

Consequence, measured on real userlands by the sibling card and reproduced in
`docs/feed-install.md:265-329`: apk-tools 3.0.5 reports `UNTRUSTED signature` and refuses the
package (`tollgate-wrt (no such package)`), and opkg 24.10.2 downloads `Packages.gz` then
**deletes** the list when `Packages.sig` 404s. **An unsigned index is not consumable on either
package manager.** That is a producing-side defect, not a consumer-side one — and it is why the
repo's own doc says, at `docs/feed-install.md:40`: *"Do not tell a tester 'install from the feed'
until the producing side ships an index."*

The chosen channel's single live artifact is therefore **not a deployable version** and **not a
trustable source**, while still being the correct *channel*.

### 2.3a The artifact filename is per-artifact, not derivable — measured

Independent second pass, 2026-10-24 (`curl -sSI`, `awk` on the headers, plus one full download):

```
$ B=https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53
$ curl -sSL "$B/Packages" | grep -E '^(Package|Version|Filename|Size|SHA256sum):'
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
Size: 5321054
SHA256sum: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk
Size: 5307464
SHA256sum: c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk
Size: 7153725
SHA256sum: 5332bb5ea15edbef3a41d5a34cf0e5e66a8d0d9a152385d4e0a7f06e928e6f3d

$ for f in <the three Filenames above>; do curl -sSI "$B/$f" | awk '...'; done
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk  http=200 clen=5321054
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk  http=200 clen=5307464
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk                  http=200 clen=7153725
```

Every `content-length` equals that entry's own `Size:` — the three index entries are three real,
distinct artifacts, and the index is a faithful manifest.

**The templated name is wrong.** Composing `tollgate-wrt_v<ver>_<arch>.ipk` from the pattern gives
a URL that does not exist:

```
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.ipk"           -> http=404
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk" -> http=200
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk"           -> http=200 (7,153,725 B)
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk" -> http=200 (5,307,464 B)
```

Note the asymmetry: for **apk** both the plain and the suffixed name exist (7.15 MB vs 5.31 MB —
different builds of the same version), while for **ipk** only the suffixed name exists. So the
suffix is neither absent-by-default nor uniform. Only the index knows.

Full download + digest check on the suffixed apk (not a ranged GET):

```
$ curl -sSL -o chk.apk "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk"
http=200 bytes=5307464      on-disk=5307464
$ sha256sum chk.apk
c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f   # == index SHA256sum:
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk"
HTTP/2 200   content-type: application/vnd.android.package-archive   content-length: 5307464
             last-modified: Wed, 01 Jul 2026 19:52:40 GMT
```

### 2.4 The two real example versioned URLs (live, digest-verified)

No two *different* versions carry artifacts on `releases.tollgate.me` (§2.2: exactly one does), so
**the "two recent versions" requirement cannot be met from the chosen channel** — stating that
plainly rather than manufacturing a second URL. The two live, version-pinned, digest-verified
examples come from the only lane that has more than one published version, the `FreedomTechFeed`
per-asset release lane:

```
https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk
https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
```

Observed HTTP + digest evidence for both (`curl -sSL -o <file> …` then `sha256sum` vs the tag's own
`SHA256SUMS`):

```
v0.6.0     asset http=200 bytes=8629996
           calc c51c492a9acf504f2f9b87af028aad401ff857df2efd12bf5fe56ba2ae45d868
           sums c51c492a9acf504f2f9b87af028aad401ff857df2efd12bf5fe56ba2ae45d868   MATCH
v0.6.0-rc6 asset http=200 bytes=8630160
           calc 3b3a85cc3ff5c4612a365d54ba56117600f079a416c3f3ce60d847d47a93e26a
           sums 3b3a85cc3ff5c4612a365d54ba56117600f079a416c3f3ce60d847d47a93e26a   MATCH
```

`SHA256SUMS` is signed and the signature verifies against the published key (ed25519, namespace
`freedomtechfeed-release-manifest`), for **both** tags:

```
$ ssh-keygen -Y verify -f allowed_signers -I release-signing@freedomtechfeed \
      -n freedomtechfeed-release-manifest -s SIG-v0.6.0 < SH-v0.6.0
Good "freedomtechfeed-release-manifest" signature for release-signing@freedomtechfeed
     with ED25519 key SHA256:lQWK9yivuF2bVKLTYaKlu2X3zgRdgfbQf1ymAazgets
$ ... -s SIG-v0.6.0-rc6 < SH-v0.6.0-rc6        -> same "Good" line, rc=0
```

Release inventory (no drafts, no rolling tag):

```
$ curl -sS "https://api.github.com/repos/FreedomTechFeed/packages/releases?per_page=100"
v0.6.0      2026-10-21T02:06:45Z  assets=26  draft=False
v0.6.0-rc6  2026-10-15T19:50:37Z  assets=26  draft=False
$ curl -sS https://api.github.com/repos/FreedomTechFeed/packages/releases/latest
v0.6.0   published=2026-10-21T02:06:45Z  draft=False
```

Two further facts that bound the choice: **upstream `OpenTollGate/tollgate-module-basic-go`
publishes `v0.6.0` (2026-10-21T03:27:56Z) with `assets=0`** — it carries no binaries, so
"upstream releases" is not a usable channel; and the `FreedomTechFeed` lane has **no** asset for
either tag the wizard pins today (§5).

### 2.5 The rolling `latest` target is not addressable as an artifact (rule evidence)

```
https://github.com/FreedomTechFeed/packages/releases/download/latest/SHA256SUMS   http=404 bytes=9 text/plain
$ curl -sSI https://github.com/FreedomTechFeed/packages/releases/latest
HTTP/2 302
$ curl -sSL -o /dev/null -w 'final=%{url_effective}' .../releases/latest
final=https://github.com/FreedomTechFeed/packages/releases/tag/v0.6.0
```

`latest` resolves only to an HTML **tag page**, never to a downloadable blob; the
`/download/latest/…` form is a flat 404. This is the mechanical reason for the rule in §3.

### 2.6 Corrections to the brief's framing (measured, so the implementer argues with the right facts)

| Brief / prior framing | Measured this session |
|---|---|
| `deploy.go:24-52` points at `blossom2.orangesync.tech` and `felixfelix-bot` | **Only `felixfelix-bot`.** `deploy.go:43`, `:47`, `:50` — all three are `github.com/felixfelix-bot/...`. `grep -rn -i 'blossom2\|orangesync' --include='*.go' .` → **0 hits**. blossom2 survives only in history (`11f63e7:deploy.go:42`). |
| "replace the unofficial mirrors" | There is exactly one unofficial *owner* (`felixfelix-bot`, the project's own fork account) across three constants — one of which (`configwizURL`, `:50`) is a **different project** (`configurationwizzard`) and is out of scope for a tollgate-wrt channel decision. |
| `wz1.5-channel-probe-raw.md` supplies the probe evidence | **File absent** (§2.0). All evidence above is first-hand from this run. |

---

## 3. Version-pinning rule (binding on the implementation)

1. **Only immutable, version-pinned paths are used.** Every URL must carry the full release
   version in its path — `<base>/v<X.Y.Z>[-<pre>]/<arch>/<artifact-file>` — where `v<X.Y.Z>` is a
   tag that already exists and will never be re-pointed at different bytes.
1b. **`<artifact-file>` is copied from the version's own index, never composed.** The filename is
   **not** a function of `<version>` and `<arch>`: the chosen host's one live version publishes a
   suffixed `…-upx-ultra-brute.ipk` with no plain `.ipk` counterpart at all (§2.3a). Read
   `Filename:` out of that version's `Packages` (or the `F:` field of `APKINDEX`) and use it
   byte-for-byte. A composed name silently 404s, which is the failure mode this rule exists to
   prevent.
2. **NEVER hash-pin a moving blob.** A `sha256` is only a valid pin when it is taken from an
   artifact reachable at an *already-frozen* versioned path. Computing a digest of whatever a
   moving URL serves today and calling that a pin is forbidden: it converts a stable reference
   into a silently stale one and defeats `artifactPins`' fail-closed posture.
3. **NEVER reference a rolling `latest` in code.** `releases/download/latest/…` is a flat 404
   (§2.5) and `/releases/latest` 302s to an HTML tag page; a `latest` reference in a Go constant
   would be both unresolvable and non-reproducible. Resolve versions out-of-band and write the
   resolved tag as a literal.
4. **The pin and the URL move together.** `artifactPins` (`deploy.go:67-70`) maps URL → sha256 and
   fails closed on a missing entry (`TestArtifactPinsCoverDownloadedURLConstants`). Any URL change
   is a **two-part edit**: the constant *and* its digest, in the same commit.

---

## 4. Architecture mapping the implementer must not get wrong

The chosen pattern's `<arch>` segment is the OpenWrt `DISTRIB_ARCH` (`aarch64_cortex-a53`),
confirmed for the GL-MT6000 on 25.12.5 from upstream's `profiles.json`, not assumed. The
`FreedomTechFeed` lane instead uses a **flat asset name** carrying a *target tuple*
(`tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk`). These naming schemes are **not**
interchangeable: a wizard that hard-codes the tuple form cannot be repointed at the channel
pattern by a host swap alone — the **filename** changes too.

Two shapes that exist on the channel host and must not be confused:

| shape | probe | verdict |
|---|---|---|
| `<channel>/<tag>/<arch>/Packages[.gz]` | 206 only for `v0.5.0-alpha3` | index; **not** usable (unsigned, stale) |
| `<channel>/<tag>/<arch>/<artifact-file>` | 206 for `v0.5.0-alpha3` `.ipk`/`.apk` | the shape this decision pins |

---

## 5. Q5 status

> **Q5 — OPEN.** Interim channel vs final channel is **undecided**.

Justification (one sentence): the decision rule resolves toward the final channel **only if the
Phase-2 feed is PUBLISHED**, and it is **NOT PUBLISHED** — no index exists for any current release
(28/28 index paths 404 across seven tags, §2.2) and the only live index (`v0.5.0-alpha3`,
2026-07-01) is stale and unsigned, so `releases.tollgate.me` is adopted as the **documented interim
channel** and the interim-vs-final choice stays open.

**Q5 is not merely procedural — it is load-bearing.** Two measured facts keep it open, and neither
may be resolved by an implementer acting alone:

- The interim channel **cannot serve the versions `deploy.go` pins today** (§6) — so "adopt the
  interim channel" and "keep deploying a working router" are in direct tension right now.
- The one lane that *does* have live, signed, digest-verified artifacts
  (`FreedomTechFeed/packages`, §2.4) is a **per-asset upload directory on an account**, which is
  precisely the thing Q5 asks whether to accept as final. Deciding that here would be deciding Q5,
  which this card's rule forbids.

Therefore: **decide the channel name (done, above); do not decide its contents (open).**

---

## 6. Follow-on: what the implementer must change in `deploy.go:24-52`

Current state at HEAD (`c60cb49`), all three constants on `felixfelix-bot`:

```
:43  tollgatePkgURL    = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.7.0-alpha10/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk"
:47  tollgatePkgURLApk = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.6.1-post-merge/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk"
:50  configwizURL      = "https://github.com/felixfelix-bot/configurationwizzard/releases/download/v1.0.10-alpha/net4sats-configwiz-1.0.10-alpha.tar.gz"
```

Both current pins are **live** (206 on a ranged GET), so this is a channel-migration edit, not a
broken-link repair:

```
v0.7.0-alpha10 …/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk          http=206 bytes=1024
v0.6.1-post-merge …/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk      http=206 bytes=1024
```

### 6.1 The instruction

**Target-state literals** — rewrite `:43` and `:47` to the chosen pattern, keeping the *same
versions and the same arch* the wizard already installs:

```
:43  tollgatePkgURL    = "https://releases.tollgate.me/alpha/v0.7.0-alpha10/aarch64_cortex-a53/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk"
:47  tollgatePkgURLApk = "https://releases.tollgate.me/alpha/v0.6.1-post-merge/aarch64_cortex-a53/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk"
```

Two caveats on those literals, both measured (§2.3a): the **path** shape is right, but (a) the
**filename** segment above is a *composed* guess — on the chosen host the real filename can carry a
build-flavour suffix (`…-upx-ultra-brute.ipk`) and the plain form does not exist, so the
implementer must take `<artifact-file>` from that version's `Packages` `Filename:` field rather
than transcribe the line above; and (b) the `<arch>` segment must be the OpenWrt `DISTRIB_ARCH`
the wizard detects, not a hard-coded `aarch64_cortex-a53` (§4). Both literals are therefore
**illustrative of shape only**.

**DO NOT LAND THESE LITERALS YET.** Both target URLs were probed and both are **404**:

```
https://releases.tollgate.me/alpha/v0.7.0-alpha10/aarch64_cortex-a53/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk
  http=404 bytes=697 type=text/html; charset=utf-8
https://releases.tollgate.me/alpha/v0.6.1-post-merge/aarch64_cortex-a53/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk
  http=404 bytes=697 type=text/html; charset=utf-8
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk
  http=404 bytes=697 type=text/html; charset=utf-8
```

Landing them would turn a working deploy into a hard 404 for every fresh router — the exact
regression recorded in the `SW4a (Aug 2026)` comment at `deploy.go:37-42` (repinned *away from*
`main.53`, "asset never existed on this release — HTTP 404, broke every fresh deploy"). The
implementer must not repeat it.

### 6.2 Precondition gate the implementer must satisfy first

Execute the change **only** when all four hold, each verifiable with one command:

1. `<chosen-base>/v<X.Y.Z>/<arch>/<artifact-file>` returns **HTTP 200** for the exact tag being
   pinned (not 206-on-range alone — confirm `content-length` equals the full artifact size).
2. The artifact's `sha256sum` equals the digest published **inside that version's own index**
   (`SHA256sum:` in `Packages`, or `T:` in `APKINDEX`) — a digest taken from anywhere else is a
   moving-blob pin and violates §3 rule 2.
3. The index carries a signature the router will accept — `Packages.sig` for opkg, a signed
   `APKINDEX.tar.gz` for apk — **and** its key is fetchable. Today `Packages.sig` and
   `keys/tollgate-testing.pem` are both 404 (§2.3), so this gate is currently **closed**.
4. The published version is **not older** than what a deployed router already runs. The wizard has
   a downgrade refusal for exactly this; `v0.5.0-alpha3` would trip it.

If the channel cannot satisfy the gate, the correct action is **not** to repin to the GitHub lane
unilaterally — that is Q5's open question (§5). Escalate: record the failing probe against each of
the four gates on the card, and leave `deploy.go` as-is.

### 6.3 Scope note

`:50` `configwizURL` is **out of scope** for this decision: it is the `configurationwizzard`
project, not a tollgate-wrt artifact, and the chosen channel hosts no configwiz asset. Do not
rewrite it as part of the tollgate-wrt migration. (A separate card, `t_deabf3d6`, has already
recommended deleting that constant outright on other grounds.)

### 6.4 Doc follow-on

When the gate in §6.2 is satisfied, `docs/feed-install.md` §0 must be updated in the same PR — its
table currently records `releases.tollgate.me/alpha/v0.6.0/... → 404` and
`.../v0.5.0-alpha3/.../Packages.sig → 404`, and its standing instruction at `:40` ("do not tell a
tester to install from the feed until §0 shows a 200") is the prose form of this gate.

---

## 7. Acceptance self-check

| Criterion | Where | Status |
|---|---|---|
| Exact chosen base URL and base URL pattern | §1 | **MET** — `https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<arch>/<artifact-file>` |
| Why it was chosen | §1, §5 | **MET** — channel-with-path-segment, Tollgate-owned; feed branch of the rule not taken |
| Evidence: commands and observed output, verbatim | §2.1–§2.6 | **MET** — all first-hand this run; sources cited by command and raw line |
| Version-pinning rule (immutable paths only; no hash-pin of moving blobs; no rolling `latest` in code) | §3 | **MET** — rules 1–4, with `latest` 404/302 evidence in §2.5 |
| ≥2 real example versioned URLs, two different recent versions, with HTTP evidence | §2.4 | **MET for the release lane** (`v0.6.0`, `v0.6.0-rc6`, both 200 + digest match) · **NOT MET for the chosen channel** — only one version exists there (§2.2), stated explicitly rather than fabricated |
| Q5 status OPEN/RESOLVED + one-sentence justification | §5 | **MET** — **OPEN**, on the NOT-PUBLISHED branch of the decision rule |
| Follow-on: actionable `deploy.go:24-52` instruction | §6 | **MET** — target literals given, landing explicitly gated on four measured preconditions |
| The instantiated artifact-filename pattern resolves | §1, §2.3a | **CORRECTED** — the first draft's composed `tollgate-wrt_v<X.Y.Z>_<arch>.ipk` is **404** on the one version that exists; the filename must be read from the index (`Filename:`), and this is now stated in §1 and enforced as §3 rule 1b |
| Note re-verified by a second, independent pass | §2.3a, §2.5, header | **MET** — every live claim re-run 2026-10-24 (host/404s/206s, index `Filename:`+`Size`+`SHA256sum`, one full download+digest, `latest` 404/302, release inventories, and a 27-URL sweep confirming exactly one version carries artifacts) |
| No unsupported claims | throughout | **MET** — every claim carries a command or a `file:line`; the two claims that could not be supported (a second channel-host version; the raw-probe input) are declared in §2.4 and §2.0 |

---

## 8. Open items handed off

1. **Q5 (channel finality)** — OPEN by rule. Needs an operator/architecture decision, not a
   measurement.
2. **Producing-side gap** — the interim channel publishes **no signed index and no artifact for any
   0.6.x/0.7.x release**. Until it does, §6.2's gate stays closed and `deploy.go` keeps pointing at
   `felixfelix-bot`. This is the single blocker for wz1.5 Phase 2.
3. **Missing sibling artifact** — `wz1.5-channel-probe-raw.md` was never written (§2.0). If a later
   card is specified as consuming it, that card's dependency is unsatisfiable as written.
