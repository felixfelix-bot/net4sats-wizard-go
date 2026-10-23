# Installing `tollgate-wrt` from a feed (GL-MT6000)

Operator procedure for installing `tollgate-wrt` on a GL-MT6000 from a package
feed instead of pushing a local `.apk`/`.ipk`. It exists because the wizard's
normal deploy path `scp`s a package and installs it by path — an install whose
`origin` is a local file, not a feed.

**Read §0 before handing anyone a command from this page.** The feed this
document describes is not published yet, so the procedure is reproducible end to
end but the feed URL in it is not yet resolvable.

Everything marked *measured* below was run on real OpenWrt userlands with the
real package managers on 2026-10-23; the transcripts are in
[`evidence/`](evidence/). Anything not marked measured is marked *planned*.

## 0. Precondition: is the feed consumable yet?

The feed is produced by the feed repository (index generation + signing) and
consumed here. As of **2026-10-23** the producing side is **not merged**, so no
public index exists for the current release:

| probe | result (measured 2026-10-23) |
|---|---|
| `FreedomTechFeed/feed` PR "auto-generate opkg+apk indices" | **CLOSED**, `mergedAt: null` — never merged |
| `FreedomTechFeed/feed` releases / tags | **0 releases, 0 tags** |
| `.../releases/download/v0.6.0/aarch64_cortex-a53/packages.adb` | **404** |
| `releases.tollgate.me/alpha/v0.6.0/...` | **404** |
| `releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/APKINDEX.tar.gz` | **200**, but **unsigned** |
| `releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.sig` | **404** |

So today there are exactly two feed-shaped artifacts, and **neither yields an
install**:

- the `v0.5.0-alpha3` index on `releases.tollgate.me` — real, but stale
  (a downgrade from any 0.6.x build) and **unsigned**, so apk refuses it
  (§4) and opkg discards it for the missing signature (§5);
- `v0.6.0` on GitHub — flat assets with **no index at all**, so there is
  nothing for a feed line to point at.

**Do not tell a tester "install from the feed" until the producing side ships an
index and this page's §0 table shows a 200 for the arch they are on.** The rest
of this page is the procedure that is ready to run the moment that happens, plus
the measured failure modes you will hit if you run it early.

Arch/board identity, which the line depends on:

```sh
. /etc/openwrt_release && echo "$DISTRIB_DESCRIPTION $DISTRIB_ARCH"
```

GL-MT6000 on OpenWrt 25.12.5 is `aarch64_cortex-a53` (`mediatek/filogic`),
confirmed from upstream's own `profiles.json` for the exact board rather than
assumed.

## 1. Which package manager the router has decides the whole procedure

The two OpenWrt lines are not interchangeable, and picking the wrong one is the
most common way this procedure fails:

| OpenWrt | package manager | index format | repository file |
|---|---|---|---|
| 24.10.x | `opkg` | `Packages.gz` | `/etc/opkg/customfeeds.conf`, `src/gz <name> <url>` |
| 25.12.x | `apk` (apk-tools 3.x) | `packages.adb` (or `APKINDEX.tar.gz`) | `/etc/apk/repositories.d/*.list` |

Check first:

```sh
which apk opkg
. /etc/openwrt_release && echo "$DISTRIB_DESCRIPTION $DISTRIB_ARCH"
```

Measured: on 24.10.2 `apk` does not exist (`apk: not found`, exit 127), and on
25.12.5 `apk --version` reports `apk-tools 3.0.5`.

## 2. Precondition the feed itself must satisfy

Before any router-side command, the feed must publish, for `aarch64_cortex-a53`:

1. a **signed index** — `packages.adb` (apk v3) or `APKINDEX.tar.gz` for apk,
   `Packages.gz` **plus `Packages.sig`** for opkg;
2. the **signing public key**, published out of band, to install under
   `/etc/apk/keys/` (apk) or `/etc/opkg/keys/` (opkg);
3. the package payload(s) the index names, at the path the index's `Filename:`
   / `o:` field refers to.

An index without its signature is not a feed — it is a file that every package
manager on a stock image will refuse (§4, §5).

## 3. Install, apk (OpenWrt 25.12.x, the GL-MT6000 line)

This is the target procedure. Steps 1–2 are *planned* (no key or index is
published yet); steps 3–7 are *measured* shapes, taken from a real
apk-tools 3.0.5 userland.

**Step 1 — trust the feed's signing key** *(planned)*

```sh
wget -O /etc/apk/keys/tollgate.pem https://<feed-host>/keys/tollgate.pem
```

The file name under `/etc/apk/keys/` is not significant; apk matches the key by
content.

**Step 2 — point apk at the feed** *(planned URL form, measured behaviour)*

```sh
echo "https://<feed-host>/<channel>/<line>/aarch64_cortex-a53/packages.adb" \
  > /etc/apk/repositories.d/tollgate.list
```

The line must be the **full URL of the index file**. apk accepts a bare
directory on the command line but, as a *repository line*, it appends
`/<arch>/APKINDEX.tar.gz` and the update fails. Measured on apk-tools 3.0.5:

```
$ echo <base>/aarch64_cortex-a53 > /etc/apk/repositories.d/tollgate.list
$ apk update
ERROR: wget: exited with error 8
WARNING: updating and opening <base>/aarch64_cortex-a53/x86_64/APKINDEX.tar.gz: unexpected end of file
1 unavailable, 0 stale; 11247 distinct packages available
```

Note the doubled arch path in that message: that is apk telling you it treated
your directory as a repository root.

**Step 3 — refresh the index** *(measured)*

```sh
apk update
```

Success looks like the tollgate line appearing in the fetched list and the
package count being non-zero with nothing unavailable. Measured shape on a
healthy feed:

```
 [https://downloads.openwrt.org/releases/25.12.5/packages/x86_64/video/packages.adb]
OK: 11247 distinct packages available
```

A feed line that is reachable but **unsigned** instead gives:

```
WARNING: updating and opening <url>: UNTRUSTED signature
1 unavailable, 0 stale; 11247 distinct packages available
```

`1 unavailable` with the feed's own URL named is the tell. It is not a warning to
be waved through: the feed contributes **nothing** to the package database, and
`--allow-untrusted` is not the fix — it disables verification for everything,
which is exactly the property the signature exists to provide.

**Step 4 — confirm the feed offers the package** *(measured)*

```sh
apk policy tollgate-wrt
```

Measured, on a working feed, this prints the version and the URL it would come
from:

```
tollgate-wrt policy:
  0.6.0_alpha2-r1:
    <feed-url>
```

If instead you see `WARNING: opening from cache <url>: No such file or
directory`, the index never made it into the cache — the feed was refused at
Step 3 and there is no candidate to install.

**Step 5 — install from the feed** *(measured)*

```sh
apk add tollgate-wrt
```

Note there is no `--allow-untrusted` and no local path: the whole point of a feed
install is that the package is selected by the index, fetched by apk, and
verified by signature. Measured shape of a success:

```
(1/2) Installing jq (1.8.1-r2)
(2/2) Installing tollgate-wrt (0.6.0_alpha2-r1)
OK: 39.1 MiB in 138 packages
```

Measured shape when the feed was refused (today's state):

```
WARNING: updating and opening <url>: UNTRUSTED signature
ERROR: unable to select packages:
  tollgate-wrt (no such package):
    required by: world[tollgate-wrt]
```

**Step 6 — prove the package came from the feed** *(measured)*

This is the check that distinguishes a feed install from the wizard's normal
local-file install, so it is not optional:

```sh
apk info -a tollgate-wrt | head -20
grep -A8 '^P:tollgate-wrt' /lib/apk/db/installed
```

The `origin` / `o:` field is the verdict:

| `o:` value | meaning |
|---|---|
| `o:/feed/net/tollgate-wrt` | installed from a **local file** — not a feed install |
| `o:tollgate-wrt` with a feed URL in the policy output | installed **from the feed** |

Measured example of the local-file case (the wizard's path, captured on the real
v0.6.0 x86_64 package):

```
P:tollgate-wrt
V:0.6.0-r1
A:x86_64
o:/feed/net/tollgate-wrt
```

`apk info -a` on that same install prints description, webpage, installed size
and `depends on:`, but no feed URL — which is the distinction in one screen.

**Step 7 — verify the service, not just the package** *(measured)*

Installing the package is not the acceptance criterion; a backend answering on
`:2121` is.

```sh
/etc/init.d/tollgate-wrt enable
/etc/init.d/tollgate-wrt start
sleep 5
wget -qO- http://127.0.0.1:2121/
```

Measured, from the real v0.6.0 package started directly (procd unavailable in
the container):

```
tcp        0      0 :::2121                 :::*                    LISTEN      922/tollgate-wrt
{"kind":10021,"id":"36fb76...","pubkey":"d3e7c936...","created_at":0,
 "tags":[["metric","bytes"],["step_size","22020096"],["tips","1","2"],
         ["price_per_step","cashu","1","sat","https://mint.coinos.io","1"], ...]}
```

A `kind:10021` advert from `:2121` is the health signal: the backend is up and
advertising pricing.

Portal is a separate service (nodogsplash on `:2050`) and must be checked
separately:

```sh
wget -qO /dev/null http://127.0.0.1:2050/splash.html && echo "portal 200"
```

Measured caveat: in a container without procd/ubus, nodogsplash refuses to start
(`Failed to connect to ubus`, `Interface br-lan not detected`) even though the
package ships `/etc/tollgate/tollgate-captive-portal-site/splash.html` correctly.
So a portal probe can only be run on a real router, not in a rehearsal
container — do not record a container portal result as a pass or as a failure.

## 4. Why the one published index will not install (measured)

For the record, because "there is an index, why doesn't it work" is the question
this page will get. On a real apk-tools 3.0.5 userland, pointed at the only
published index:

```
$ echo https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/APKINDEX.tar.gz \
    > /etc/apk/repositories.d/tollgate.list
$ apk update
WARNING: updating and opening .../APKINDEX.tar.gz: UNTRUSTED signature
1 unavailable, 0 stale; 11247 distinct packages available
$ apk add tollgate-wrt
ERROR: unable to select packages:
  tollgate-wrt (no such package):
    required by: world[tollgate-wrt]
```

Two independent blockers, both measured:

1. **Unsigned.** `Packages.sig` → 404, and the `APKINDEX.tar.gz` carries no
   signature apk trusts. apk requires the index to be signed; there is no
   supported router-side flag that makes an unsigned index usable.
2. **Stale.** The index resolves `tollgate-wrt v0.5.0-alpha3`, which is older
   than the 0.6.x already on a deployed router. The wizard's own deploy path has
   a downgrade refusal for exactly this, so installing it would be a regression
   even if it were signed.

Also measured: the index advertises `A:aarch64_cortex-a53` while an x86_64 host
is asking, so on a mismatched host even a signed index would stop at
`apk add` with `error: uninstallable / arch: aarch64_cortex-a53` — the update
succeeds and only the install fails, which is a confusing failure worth
recognising.

## 5. The 24.10.x / opkg path (measured, and it fails the same way)

If the board is on 24.10.x the line goes in `/etc/opkg/customfeeds.conf` as
`src/gz <name> <base>`, where `<base>` is the **directory containing**
`Packages.gz` — opkg appends `/Packages.gz`, exactly as apk appends
`/<arch>/APKINDEX.tar.gz`. That asymmetry is the trap: the same URL is wrong for
one manager and right for the other.

Measured on real OpenWrt 24.10.2 (`opkg` 38eccbb):

```
# with the base INCLUDING a trailing /packages segment
Downloading .../aarch64_cortex-a53/packages/Packages.gz
*** Failed to download the package list from .../aarch64_cortex-a53/packages/Packages.gz
 * opkg_download: ... wget returned 8.

# with the base being the arch directory itself
Downloading .../aarch64_cortex-a53/Packages.gz
Updated list of available packages in /var/opkg-lists/tollgate
Downloading .../aarch64_cortex-a53/Packages.sig
Signature file download failed.
Remove wrong Signature file.
 * opkg_download: ... Packages.sig, wget returned 8.
-> /var/opkg-lists/tollgate is REMOVED; opkg list tollgate-wrt -> nothing
```

So on opkg the index downloads and then is **discarded** because `Packages.sig`
404s. `option check_signature 0` does not rescue it in 24.10.2's opkg — the
signature fetch is still attempted and the list is still removed. The measured
conclusion is the same as §4: an unsigned index is not consumable on either
package manager, and this is a producing-side defect, not a consumer-side one.

## 6. Upgrade path

Moving to a later feed release is the same three commands, because the index URL
is stable across versions and version selection is the package manager's job:

```sh
apk update            # fetch the new index
apk policy tollgate-wrt   # confirm a newer version is now offered
apk add -u tollgate-wrt   # upgrade (apk)
# opkg: opkg update && opkg upgrade tollgate-wrt
```

Two things to pin deliberately:

- **Pin by URL, not by habit.** The channel/line segment in the repository URL
  (`<channel>/<line>/<arch>/`) is what decides which release stream the router
  follows. A device on a `stable` line upgrades within stable; a device on
  `testing` follows testing. Do not point a production router at a testing line
  and expect it not to move.
- **Pin a version when reproducibility matters.** apk will otherwise take the
  highest version the index offers. To hold a specific build:
  `apk add 'tollgate-wrt=0.6.0-r1'`.

Verify after upgrading with the same Step 6 origin check plus Step 7 health
probe — an upgrade that changes the origin from the feed to a local file has
silently replaced the feed install.

## 7. Rollback and removal

Removal is the supported rollback; there is no "downgrade in place" to rely on,
and pinning a lower version than the index offers is a fight with the resolver
rather than a procedure.

```sh
# stop and disable first, so procd does not restart something you are removing
/etc/init.d/tollgate-wrt stop
/etc/init.d/tollgate-wrt disable

# remove the package (apk)
apk del tollgate-wrt
# opkg: opkg remove tollgate-wrt --force-depends

# confirm it is gone
ls /usr/bin/tollgate /usr/bin/tollgate-wrt   # expect: No such file or directory
```

Measured on a real install (`apk del` after a feed install):

```
(1/2) Purging tollgate-wrt (0.6.0_alpha2-r1)
(2/2) Purging jq (1.8.1-r2)
OK: 11.4 MiB in 136 packages
ls: /usr/bin/tollgate: No such file or directory
ls: /usr/bin/tollgate-wrt: No such file or directory
```

Note `jq` was removed with it in that transcript because it had been pulled in as
a dependency and nothing else needed it. On a real router `jq` is usually
required by other packages and will stay; do not read its removal here as
expected behaviour.

To **stop following a feed** without removing the package, delete the repository
line and refresh, so no future `apk update` can move the package again:

```sh
rm /etc/apk/repositories.d/tollgate.list
apk update
```

Config under `/etc/tollgate/` and `/etc/config/` is not package-owned state and
survives removal; clear it deliberately if you want a clean re-onboard.

## 8. Troubleshooting, keyed on the message you actually see

All messages below are verbatim from the 2026-10-23 captures.

| message | cause | fix |
|---|---|---|
| `apk: not found` | board is 24.10.x (opkg) | use §5, not §3 |
| `ERROR: wget: exited with error 8` + doubled `/<arch>/APKINDEX.tar.gz` | repository line is a directory, not the index file | use the full index URL (§3 Step 2) |
| `WARNING: ... UNTRUSTED signature` | index is unsigned, or the signing key is not in `/etc/apk/keys/` | install the feed's public key; do not use `--allow-untrusted` |
| `apk add`: `no such package` after a clean-looking update | feed contributed nothing — usually the unsigned case above | fix the key/signature, not the package name |
| `error: uninstallable` / `arch: aarch64_cortex-a53` on a different host | wrong arch package for this host | use the arch that matches `DISTRIB_ARCH` |
| `apk add /path/file.apk` fails on `jq`/`nodogsplash` with `Operation not permitted` | no network, so dependencies cannot be fetched | fix egress before installing |
| opkg: `Signature file download failed` then list REMOVED | feed publishes no `Packages.sig` | producing-side defect; not fixable on the router |
| `Failed to connect to ubus`, `Interface br-lan not detected` from nodogsplash | running in a container, not a router | probe the portal on real hardware only |

## 9. What this procedure does not claim

Stated plainly so a reader does not over-read the evidence:

- **No GL-MT6000 was used for the transcripts.** They come from real OpenWrt
  userlands (25.12.5 / apk-tools 3.0.5 and 24.10.2 / opkg 38eccbb) on x86_64,
  and from the real published `v0.6.0` x86_64 package. The arch/board identity
  for the MT6000 is confirmed from upstream `profiles.json`, not from a device.
- **No aarch64 package was executed.** The aarch64 `.apk` was downloaded and its
  digest checked against `SHA256SUMS`, but it cannot be run on an x86_64 host.
- **No install from a public feed succeeded**, because no public feed index for
  the current release exists. That is the §0 finding, not a gap in the
  procedure.
- **The portal was never seen answering 200** in these captures; nodogsplash
  cannot start without procd/ubus. The `:2121` backend was seen answering.

## 10. Evidence

| file | what it is |
|---|---|
| [`evidence/2026-10-23-fe2b-measured-captures.md`](evidence/2026-10-23-fe2b-measured-captures.md) | the captures this page quotes: feed probes, apk refusal, real local-file install + origin, `:2121` health, opkg 24.10.2 src/gz runs, arch negative |
