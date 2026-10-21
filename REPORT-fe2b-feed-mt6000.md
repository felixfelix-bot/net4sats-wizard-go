# REPORT — FE2b: install tollgate-wrt on GL-MT6000 strictly from the published FEED

**Date:** 2026-10-21
**Repo under work:** `felixfelix-bot/net4sats-wizard-go` (worktree `/home/c08r4d0r/repos/net4sats-mvp-v2`,
detached HEAD `8a8d006`)
**Target:** GL-MT6000 at `192.168.1.1` (board `glinet,gl-mt6000`, hostname `tollgate-0195`)
**Card:** no board record exists for this task (searched `~/.hermes/kanban.db` for MT6000 /
feed / `src/gz` / install titles → 0 rows). Worked from the brief text alone.

---

## 0. Verdict

**PARTIAL — 1 of 5 steps completed. The card's core deliverable CANNOT be met: it is
blocked upstream, not on this device.**

| # | Step | Result |
|---|---|---|
| 1 | Add feed line idempotently | **DONE** — added once, re-run appends nothing, persists |
| 2 | `apk update`, capture full output, confirm feed index fetched | **FAILED** — `apk: not found` (device is opkg); the opkg equivalent fails because **no index is published** |
| 3 | `apk add tollgate-wrt`, resolved from the feed URL | **FAILED** — `apk: not found`; opkg cannot resolve it from any feed because no feed exists |
| 4 | Prove origin is the feed, not a local file | **ANSWERED, NEGATIVE** — origin is a local file `/tmp/tollgate-wrt.ipk`, byte-identical to the release asset. No tollgate feed list exists on the device |
| 5 | Feed line persists across reboot | **DONE** — verified present in the overlay upper dir + read back |

Two independent root causes, both outside this worker's authority:

- **A. No feed index is published** (structural, upstream). Not the `src/gz` form, not
  the URL: there is nothing to serve an index.
- **B. The device has no working WAN** (environmental). Even a real feed could not be
  fetched from this router right now.

Per the card's own escape hatch — *"If the artifacts do not contain a [feed], STOP and
report to the manager rather than substituting an architecture"* — I did not substitute
an architecture, a publisher, or an index form. I executed the card's steps literally
and captured the failures verbatim.

---

## 1. Premise correction: the parent did NOT record a FEED_BASE_URL

The brief states: *"parent task preflight confirmed FE2a is merged and the release is
tagged, and recorded the exact feed base URL and src/gz line. Read that parent's result
first."*

I read it. It says the opposite, three times over
(`REPORT-t_1b1cf42f.md`, 922 lines, committed as `8a8d006`/`1fd6778`/`7260034`):

```
| asked for | answer |
| exact `src/gz tollgate <FEED_BASE_URL>/<arch>/packages` line | DOES NOT EXIST — no feed index is published (evidence §3) |
| raw curl output proving the index is fetchable | no index is fetchable — 404 at every candidate path (§3.3) |
```

Its run-3 conclusion: *"signature files exist but sign flat assets, no feed index."* And
run 2 found the authoritative definition in the publisher's own installer code:

```
feedRepoSlug = "FreedomTechFeed/packages"
feedReleaseTagDefault = "v0.6.0-rc6"
func feedReleaseURLPrefix() string { return "https://github.com/"+feedRepoSlug+"/releases/download/"+feedReleaseTag+"/" }
func feedAssetURL(arch, target, ext string) string { return feedReleaseURLPrefix()+"tollgate-wrt_"+feedPkgVersion()+"_"+arch+"-"+target+ext }
```

i.e. in this project **"feed" denotes a per-asset URL prefix, not an index directory.**
There is no FEED_BASE_URL in the parent's result to use. Rather than refuse on a
technicality, I adopted the parent's own candidate base
(`https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6`) and ran
every step against it.

---

## 2. The device, and the parent's open conflict resolved

Prior cards recorded the MT6000 as ping-dead and credential-denied. This run it answers:

```
$ ssh -o BatchMode=yes root@192.168.1.1 'cat /tmp/sysinfo/board_name; cat /tmp/sysinfo/model; uname -a'
glinet,gl-mt6000
GL.iNet GL-MT6000
Linux tollgate-0195 6.6.93 #0 SMP Mon Jun 23 20:40:36 2025 aarch64 GNU/Linux

$ ... 'cat /etc/openwrt_release | head -5'
DISTRIB_ID='OpenWrt'
DISTRIB_RELEASE='24.10.2'
DISTRIB_REVISION='r28739-d9340319c6'
DISTRIB_TARGET='mediatek/filogic'
DISTRIB_ARCH='aarch64_cortex-a53'

$ ... 'which opkg apk; opkg --version | head -2'
/bin/opkg
opkg version 38eccbb1fd694d4798ac1baf88f9ba83d1eac616 (2024-10-16)
```

**This settles the parent's unresolved §10.3 conflict.** `TARGETS.md` said MT6000 =
opkg/24.10.2; the prior card `t_877ad4ef` claimed the FE2b target was apk/25.12.5 but
admitted it had measured a **MT3000** because no MT6000 was on the bench. Measured
directly, this MT6000 is **OpenWrt 24.10.2 / opkg**. So:

- the `src/gz` form is **correct** for this device (the parent's §10.2 "src/gz is the
  wrong form" correction applied to the MT3000, not here);
- **every `apk …` step in the brief is unrunnable on this device** — it is transcribed
  from the 25.12/apk lane.

---

## 3. Step 1 — feed line added idempotently (DONE)

Target file: `/etc/opkg/customfeeds.conf` (24.10 opkg layout; no `feeds.d/` here —
`/etc/opkg/` holds `customfeeds.conf`, `distfeeds.conf`, `keys/`).

```
$ cp /etc/opkg/customfeeds.conf /tmp/customfeeds.conf.pre-fe2b.bak && echo 'backup: ok'
backup: ok
$ grep -qF 'src/gz tollgate ' /etc/opkg/customfeeds.conf && echo 'ALREADY PRESENT - no append' \
    || { printf '\n# FE2b trial: tollgate feed\nsrc/gz tollgate %s/aarch64_cortex-a53/packages\n' "$FEED_BASE" >> /etc/opkg/customfeeds.conf; echo 'appended'; }
appended
```

Read back:

```
$ cat /etc/opkg/customfeeds.conf
# add your custom package feeds here
#
# src/gz example_feed_name http://www.example.com/path/to/files

# FE2b trial: tollgate feed
src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages
```

Idempotency, second run:

```
$ grep -qF 'src/gz tollgate ' /etc/opkg/customfeeds.conf && echo 'ALREADY PRESENT - no append (2nd run)'
ALREADY PRESENT - no append (2nd run)
$ grep -c '^src/gz tollgate ' /etc/opkg/customfeeds.conf
1
```

Exactly one line, no duplicate.

---

## 4. Step 2 — "apk update" is impossible here; the opkg equivalent fails

### 4.1 The literal command (as the card words it)

```
$ apk update
ash: apk: not found
EXIT=127
```

### 4.2 The device's real package manager, `opkg update` — FULL OUTPUT

```
$ opkg update
Downloading https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages/Packages.gz
Failed to send request: Operation not permitted
Failed to send request: Operation not permitted
Failed to send request: Operation not permitted
Failed to send request: Operation not permitted
*** Failed to download the package list from https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/packages/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/packages/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/base/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/base/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/kmods/6.6.93-1-6a9e125268c43e0bae8cecb014c8ab03/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/kmods/6.6.93-1-6a9e125268c43e0bae8cecb014c8ab03/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/luci/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/luci/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/packages/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/packages/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/routing/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/routing/Packages.gz

Downloading https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/telephony/Packages.gz
*** Failed to download the package list from https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/telephony/Packages.gz

Collected errors:
 * opkg_download: Failed to download https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/packages/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/base/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/targets/mediatek/filogic/kmods/6.6.93-1-6a9e125268c43e0bae8cecb014c8ab03/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/luci/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/packages/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/routing/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.
 * opkg_download: Failed to download https://downloads.openwrt.org/releases/24.10.2/packages/aarch64_cortex-a53/telephony/Packages.gz, wget returned 4.
 * opkg_download: Check your network settings and connectivity.

EXIT=8
```

**The tollgate feed index was NOT fetched.** Note carefully: OpenWrt's own seven feeds
fail identically. The failure is therefore not specific to the tollgate line — it is the
device's egress (cause B). But even with egress, the tollgate line would fail for a
second, independent reason (cause A, §6): the URL 404s.

---

## 5. Step 3 — "apk add" impossible; opkg resolves nothing from the feed

```
$ apk add tollgate-wrt
ash: apk: not found
EXIT=127

$ opkg upgrade tollgate-wrt
EXIT=0            <- no output: nothing to upgrade, no candidate from any feed

$ opkg list | grep -i tollgate
tollgate-wrt - 0.6.0_rc6-r1
```

`opkg upgrade` exiting 0 with zero output means opkg had **no feed candidate** for
`tollgate-wrt`; the line `opkg list` prints comes from the installed-package status
database (`/usr/lib/opkg/status`), not from an index. The decisive proof:

```
$ ls /var/opkg-lists/ | grep -i toll || echo "NO tollgate list file"
NO tollgate list file
```

Only the seven OpenWrt lists are present, all dated Oct 20 22:09 (pre-existing, intact
after the failed update — no regression).

---

## 6. Step 4 — provenance: the installed package came from a LOCAL FILE

This is the answer to the card's provenance question, and it is negative.

### 6.1 The package is already installed — and predates this run

```
$ opkg list-installed | grep -i tollgate
tollgate-wrt - 0.6.0_rc6-r1

$ grep -A20 '^Package: tollgate-wrt' /usr/lib/opkg/status
Package: tollgate-wrt
Version: 0.6.0_rc6-r1
Depends: libc, nodogsplash, jq, ca-bundle
Status: install user installed
Architecture: aarch64_cortex-a53
Conffiles:
 /etc/nftables.d/20-nds-enforce.nft 8eb1987d96fdde60ac5f8a3e31514c932813bb7da62dd0b6fa1326a911c02641
Installed-Time: 1792514447

$ date -d @1792514447
Tue Oct 20 16:40:47 UTC 2026          <- ~9.5 h before this run (device date: Wed Oct 21 02:13:25 UTC 2026)
```

### 6.2 The local artifact is still on the device, and is byte-identical to the release

```
$ ls -la /tmp/tollgate-wrt.ipk /tmp/tollgate-setup.log /tmp/jq_1.8.1-r1_aarch64_cortex-a53.ipk /tmp/nodogsplash_5.0.2-r1_aarch64_cortex-a53.ipk
-rw-r--r-- 1 root root   423212 Oct 20 16:40 /tmp/jq_1.8.1-r1_aarch64_cortex-a53.ipk
-rw-r--r-- 1 root root    50340 Oct 20 16:40 /tmp/nodogsplash_5.0.2-r1_aarch64_cortex-a53.ipk
-rw------- 1 root root     2798 Oct 20 16:40 /tmp/tollgate-setup.log
-rw-r--r-- 1 root root  8630160 Oct 20 16:40 /tmp/tollgate-wrt.ipk

$ sha256sum /tmp/tollgate-wrt.ipk
3b3a85cc3ff5c4612a365d54ba56117600f079a416c3f3ce60d847d47a93e26a  /tmp/tollgate-wrt.ipk

# same asset downloaded from the release, on my host:
$ curl -sSL -o /tmp/rc6asset.ipk .../tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
$ sha256sum /tmp/rc6asset.ipk
3b3a85cc3ff5c4612a365d54ba56117600f079a416c3f3ce60d847d47a93e26a  /tmp/rc6asset.ipk
```

Same digest, same 8630160 bytes. The installed package is the release's flat `.ipk`,
installed from a local file in `/tmp` — **exactly the mechanism the card says must not
be used.** Also present: `/tmp/tollgate-installer.*/installer.log` (the installer's
scratch dir) and its own provenance line:

```
$ grep -iE 'feed|repositor|src/gz|opkg (update|install)|\.ipk' /tmp/tollgate-installer.pjNefI/installer.log | head
Feed: FreedomTechFeed/packages release v0.6.0-rc6 (package 0.6.0_rc6)
```

That is the installer *naming* its source. It is not a feed index and it is not an
`opkg install` — it is the flat-asset download path from §1.

### 6.3 The `apk info -a` / `apk list --installed` steps in the brief

```
$ apk info -a tollgate-wrt
ash: apk: not found
$ apk list --installed tollgate-wrt
ash: apk: not found
```

Unrunnable on an opkg device. The opkg equivalents (`opkg info`, `opkg list-installed`)
are quoted above and both show the status-database origin, not a feed.

---

## 7. Step 5 — persistence (DONE)

`/etc/opkg/customfeeds.conf` is on the overlayfs upper (f2fs `/dev/loop0` on `/overlay`),
so it survives reboot — verified directly rather than assumed:

```
$ mount | grep -E 'overlay|/overlay'
/dev/loop0 on /overlay type f2fs (rw,lazytime,noatime,...)
overlayfs:/overlay on / type overlay (rw,noatime,lowerdir=/,upperdir=/overlay/upper,workdir=/overlay/work,...)

$ ls -la /overlay/upper/etc/opkg/customfeeds.conf
-rw-r--r--    1 root     root           249 Oct 21 02:06 /overlay/upper/etc/opkg/customfeeds.conf
```

The file exists in the **upper** (writable, persistent) layer, not in the read-only
lower image. Read back, and the persistent backup:

```
$ cat /etc/opkg/customfeeds.conf
# add your custom package feeds here
#
# src/gz example_feed_name http://www.example.com/path/to/files

# FE2b trial: tollgate feed
src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages

$ grep -c '^src/gz tollgate ' /etc/opkg/customfeeds.conf
1

$ ls -la /root/customfeeds.conf.pre-fe2b.bak
-rw-r--r--    1 root     root           249 Oct 21 02:14 /root/customfeeds.conf.pre-fe2b.bak
```

Because the line is known-inert-but-noisy (§8), I left a one-command revert on the
device: `/root/revert-fe2b-feed.sh` (restores the backup, or strips the two trial lines
if the backup is gone).

---

## 8. Root cause A — there is no published feed index (the structural blocker)

Re-derived this run, not inherited from the parent. The release surface is unchanged:

```
$ curl -sS -o /dev/null -w '%{http_code} %{url_effective}\n' -L .../FreedomTechFeed/packages/releases/latest
200 https://github.com/FreedomTechFeed/packages/releases/tag/v0.6.0-rc6
$ git ls-remote --tags .../FreedomTechFeed/packages.git | tail
v0.6.0-rc6            <- still newest; no new tag in the ~3 days since the parent ran
```

26 assets, all flat; the aarch64_cortex-a53/mediatek-filogic tuple exists as
`.apk`, `.ipk`, and `-offline.tar.gz`. **No directory entry, no index, of any form:**

```
$ B=.../releases/download/v0.6.0-rc6
$ for p in Packages Packages.gz Packages.sig aarch64_cortex-a53/Packages \
           aarch64_cortex-a53/Packages.gz aarch64_cortex-a53/packages.adb \
           packages.adb aarch64_cortex-a53/packages/Packages \
           aarch64_cortex-a53/packages/Packages.gz aarch64_cortex-a53/packages/Packages.sig \
           mediatek-filogic/Packages.gz aarch64_cortex-a53-mediatek-filogic/Packages.gz \
           index.json SHA256SUMS tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk; do
    printf '%s  %s\n' "$(curl -sS -o /dev/null -w '%{http_code}' -L "$B/$p")" "$p"; done
404  Packages
404  Packages.gz
404  Packages.sig
404  aarch64_cortex-a53/Packages
404  aarch64_cortex-a53/Packages.gz
404  aarch64_cortex-a53/packages.adb
404  packages.adb
404  aarch64_cortex-a53/packages/Packages
404  aarch64_cortex-a53/packages/Packages.gz
404  aarch64_cortex-a53/packages/Packages.sig
404  mediatek-filogic/Packages.gz
404  aarch64_cortex-a53-mediatek-filogic/Packages.gz
404  index.json
200  SHA256SUMS
200  tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk
```

**15 index shapes probed → every one 404. Flat assets → 200.** The exact URL opkg
requested in §4.2 is the first line: `.../aarch64_cortex-a53/packages/Packages.gz` → 404.

### 8.1 I also checked the one thing the parent did not: a real feed repo

`FreedomTechFeed/feed` and `OpenTollGate/tollgate-feed` both exist (HTTP 200). I
inspected both rather than assuming.

**`FreedomTechFeed/feed` — a `src-git` BUILD-TIME feed, not a binary repository.**

```
$ head -20 README.md
# TollGate OpenWrt Feed
An OpenWrt package feed that builds `tollgate-wrt` from the upstream
`OpenTollGate/tollgate-module-basic-go` source using OpenWrt's `golang-package.mk`.
...
## Using the feed
echo "src-git tollgate https://github.com/FreedomTechFeed/feed.git" >> feeds.conf
./scripts/feeds update tollgate
./scripts/feeds install tollgate-wrt
```

```
$ ls main branch
.github  AGENTS.md  PLAN.md  README.md  net  scripts
$ net/tollgate-wrt/Makefile -> 200
$ git ls-remote --tags .../FreedomTechFeed/feed.git -> (none)
$ releases page -> no release assets
```

It contains a **Makefile**, not packages or an index. `src-git` is consumed by the
OpenWrt *build system* on a workstation; it is not a runtime `src/gz` opkg feed and the
router's opkg cannot use it.

The index generator exists only on an **unmerged branch**, `ci/auto-feed-indices`:

```
blob .github/workflows/build-feed.yml
blob FEED-MANIFEST.conf
blob docs/feed-strategy.md
blob scripts/generate-apk-index.sh
blob scripts/generate-packages-index.sh
blob scripts/index-feed.sh
```

And its strategy doc is a **plan**, not a shipped state — the indices are labelled as
things to add:

```
### Option A: Hosted Package Index on releases.tollgate.me (Recommended)
Adding `Packages.gz` (opkg) and `packages.adb`/`APKINDEX.tar.gz` (apk) index files turns it into a proper package repository.
...
      aarch64_cortex-a53/
        1fcc1635...71d.apk
        Packages.gz         ← NEW: opkg index
        packages.adb        ← NEW: apk index
```

**"← NEW" is the answer.** The feed index is a planned deliverable. It does not exist
yet. `FreedomTechFeed/feed` main is source-only; the index generator is unreleased.

So cause A is upstream and structural: **nothing in the org publishes a runtime package
index for `tollgate-wrt`.** A `src/gz` line cannot work until one is hosted — which is
exactly what the parent card concluded three times, now with the feed repo itself
checked and the *reason* identified (the generator is on an unmerged branch).

---

## 9. Root cause B — the device has no working WAN (environmental)

Independent of cause A: this router cannot reach the internet at all right now.

```
$ ip route
default via 10.2.0.1 dev phy0-sta0
10.2.0.0/16 dev phy0-sta0 proto kernel scope link src 10.2.1.122
192.168.1.0/24 dev br-lan proto kernel scope link src 192.168.1.1
192.168.2.0/24 dev br-private proto kernel scope link src 192.168.2.1

$ iw dev phy0-sta0 link
Connected to b2:36:4a:2a:03:32 (on phy0-sta0)
        SSID: StarGate
        signal: -69 dBm

$ ping -c 2 -W 3 10.2.0.1        -> 2 received, 0% loss   (its own gateway answers)
$ ping -c 2 -W 3 1.1.1.1         -> 2 transmitted, 0 received, 100% packet loss
$ nslookup github.com 1.1.1.1    -> connection timed out; no servers could be reached
$ nslookup github.com 10.2.0.1   -> ** server can't find github.com: NXDOMAIN / SERVFAIL
$ curl --resolve github.com:443:140.82.121.3 https://github.com -> (28) timed out after 12001 ms
$ cat /etc/resolv.conf           -> nameserver 127.0.0.1 ; nameserver ::1
```

The device's only uplink is a STA association on `phy0-sta0`; it reaches `10.2.0.1` but
nothing beyond, and its local dnsmasq returns NXDOMAIN/SERVFAIL for everything.

**This is not a firewall problem** — I checked, because the first hypothesis (STA
interface in no zone) would have been fixable:

```
$ nft list chain inet fw4 output | head -12
chain output { type filter hook output priority filter; policy accept;
    oif "lo" accept
    ct state vmap { established : accept, related : accept }
    oifname "br-lan" jump output_lan
    oifname { "eth1", "phy0-sta0" } jump output_wan      <- STA IS in the wan zone
    oifname "br-private" jump output_private }

$ uci show network | grep -E 'proto=|ifname|device='
network.wan.device='phy0-sta0'
network.wan.proto='dhcp'
```

Output policy is `accept` and `phy0-sta0` is explicitly in `output_wan`. The traffic is
not being dropped locally — **the upstream (SSID "StarGate", 10.2.0.1) does not route to
the internet.** That is outside this device and outside this worker's authority.

Consequence: even if a feed index existed, `opkg update` could not fetch it from this
router until its WAN is fixed. Both blockers must clear.

---

## 10. Acceptance criteria — scored honestly

| Criterion | Met? |
|---|---|
| `apk update` output captured verbatim | **N/A** — `apk` does not exist on this device; `opkg update` captured in full (§4.2), tollgate index NOT fetched |
| install output showing package resolved from the feed URL | **NO** — no install from a feed occurred; `opkg upgrade` found no candidate (§5) |
| installed-package origin is the feed | **NO — it is a local file** (§6): `/tmp/tollgate-wrt.ipk`, sha256 `3b3a85cc…`, == the release asset |
| feed line persists | **YES** (§7) — overlay upper + read back + backup + revert script |
| no local package file involved | **NO** — a local `/tmp/*.ipk` is exactly what was involved, pre-dating this run |

Overall: **PARTIAL / BLOCKED.** The card's central requirement — install *from the feed*,
not from a scp'd file — cannot be satisfied because there is no feed to install from, and
the device cannot reach the internet to try.

---

## 11. Repo change

**No source change is needed or made.** This is a router-state + upstream-publication
finding; there is no code in `net4sats-wizard-go` that should change as a result. In
particular I did **not** add a `src/gz`/feed URL constant to this repo — the repo is a
thin UI wrapper (AGENTS.md) and `deploy.go` already pins per-asset URLs, which §1 shows
is the actual published design. Adding a feed URL constant would encode a URL that 404s.

Artifacts committed to this repo are **records only**:
`PROGRESS-fe2b-feed-mt6000.md`, `REPORT-fe2b-feed-mt6000.md` (this file).

### 11.1 Deviation: I did not overwrite the root REPORT.md / PROGRESS.md

The brief says to write `PROGRESS.md` at the worktree root and the full report to
`REPORT.md`. In this worktree those two files belong to **other cards** — root `REPORT.md`
is "WZ1.9 history atomicity (QGATES v3.2)" and root `PROGRESS.md` is "testid gap analysis".
Overwriting them would destroy another card's record. I wrote card-scoped files at the
worktree root instead (`PROGRESS-fe2b-feed-mt6000.md`, `REPORT-fe2b-feed-mt6000.md`).
The deliverable is at the root as asked, just namespaced.

---

## 12. Device state after this run

| Item | State |
|---|---|
| `tollgate-wrt` | 0.6.0_rc6-r1, installed (unchanged, `Installed-Time` 1792514447) |
| service | running (PIDs 12996/16305). NB pre-existing: `init.d status` reports *"API :2121 is listening but /var/run/tollgate.sock is absent (CLI commands will fail)"* — not caused by this run |
| `/etc/opkg/customfeeds.conf` | +2 lines (comment + `src/gz tollgate …`), backup at `/root/customfeeds.conf.pre-fe2b.bak` |
| `/var/opkg-lists/` | unchanged — 7 OpenWrt lists only, no tollgate list |
| revert | `/root/revert-fe2b-feed.sh` |
| **recommended** | run the revert, or fix the WAN, until a real index is published — the line currently makes every `opkg update` exit 8 |

---

## 13. What would unblock this (for the manager)

1. **Publish the index.** Merge `FreedomTechFeed/feed`'s `ci/auto-feed-indices`
   (`scripts/index-feed.sh`, `generate-packages-index.sh`) and publish `Packages.gz` +
   `Packages.sig` (opkg, for 24.10 devices) and `packages.adb` (apk, for 25.12 devices)
   next to the flat assets. **Until then no `src/gz`/`packages.adb` line can work** — this
   is the single blocking dependency.
2. **Fix the MT6000's WAN.** Its STA uplink ("StarGate" → 10.2.0.1) does not route; the
   router cannot reach the internet. Both blockers must clear before any feed install.
3. **Correct the card text.** It is transcribed from the apk/25.12 lane and asks for
   `apk update`/`apk add`/`apk info` on a device that has only `opkg`; and its premise
   ("parent recorded the exact feed base URL") contradicts the parent's actual result.
4. **Decide the local-file question.** The device is currently in the forbidden state
   (installed from `/tmp/tollgate-wrt.ipk`). Either accept the flat-asset installer as the
   supported path (which is what `deploy.go` and the installer actually implement), or
   gate the feed work on (1)+(2) — but the two are not equivalent and the card assumes the
   latter.
