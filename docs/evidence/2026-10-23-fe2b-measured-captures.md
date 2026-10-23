# Measured captures — feed-install procedure, 2026-10-23

Verbatim captures behind [`../feed-install.md`](../feed-install.md). Everything
here was run on 2026-10-23 from the build host; nothing is reconstructed from an
older report.

Environments, both real OpenWrt userlands imported into containers:

```
OpenWrt 25.12.5 r33051-f5dae5ece4   arch=x86_64   apk-tools 3.0.5
OpenWrt 24.10.2 r28739-d9340319c6   arch=x86_64   opkg 38eccbb1fd694d4798ac1baf88f9ba83d1eac616
```

Rootfs images (public, digest-pinned at capture time):

```
2dae4250b1b742062ed5787adcdc9dd6d26c7f69f94ced46548eb512d5957e8a  openwrt-25.12.5-x86-64-rootfs.tar.gz
```

## 1. Is there a published feed index? (feed-host probes)

```
404 text/html; charset=utf-8                  https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/packages
200 application/octet-stream                  https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages
200 application/gzip                          https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.gz
404 text/html; charset=utf-8                  https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/packages.adb
200 application/gzip                          https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/APKINDEX.tar.gz
404 text/html; charset=utf-8                  https://releases.tollgate.me/keys/tollgate-testing.pem
404 text/html; charset=utf-8                  https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/Packages.gz
404 text/html; charset=utf-8                  https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/packages.adb
404                                            https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/aarch64_cortex-a53/packages.adb
200 application/gzip (via -L, octet-stream)    https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/SHA256SUMS
404                                            https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.sig
```

Index contents of the one published index (`APKINDEX.tar.gz`, 407 bytes,
gzip, mtime 2026-07-01):

```
P:tollgate-wrt
A:aarch64_cortex-a53
V:v0.5.0-alpha3
S:5307464
I:5307464
T:Q1c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f
D:nodogsplash luci jq
p:nodogsplash-files
m:TollGate <tollgate@tollgate.me>
c:TollGate Basic Module for OpenWrt
o:tollgate-wrt
u:https://tollgate.me
(second stanza, same version, S:7153725)
```

`DESCRIPTION`: `The TollGate alpha index`

Producing side:

```
$ gh pr list -R FreedomTechFeed/feed --state all
  #1  ci/auto-feed-indices  state=CLOSED  mergedAt=null
  #11 wrapper/point-to-canonical       OPEN
  #10 worker-base/FEED-SCAFFOLD-...    OPEN
  #9  ci/upstream-prep                 OPEN
  ... (#7, #4, #3 merged)
$ gh release list -R FreedomTechFeed/feed --limit 5   -> (empty)
$ gh api repos/FreedomTechFeed/feed/tags --jq .[].name -> (empty)
```

Current release assets (`FreedomTechFeed/packages` `v0.6.0`): 26 assets —
`.apk` + `.ipk` + `-offline.tar.gz` per arch/target, plus `SHA256SUMS` and
`SHA256SUMS.sig`. **No index of any shape.**

## 2. apk (25.12.5): the published index is refused

```
$ echo $U/APKINDEX.tar.gz > /etc/apk/repositories.d/tollgate.list
$ apk update
WARNING: updating and opening https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/APKINDEX.tar.gz: UNTRUSTED signature
1 unavailable, 0 stale; 11247 distinct packages available

$ apk policy tollgate-wrt
WARNING: opening from cache https://.../APKINDEX.tar.gz: No such file or directory

$ apk add tollgate-wrt
WARNING: updating and opening https://.../APKINDEX.tar.gz: UNTRUSTED signature
ERROR: unable to select packages:
  tollgate-wrt (no such package):
    required by: world[tollgate-wrt]
```

Directory-form repository line (the `src/gz`-shaped mistake on apk):

```
$ echo $U > /etc/apk/repositories.d/tollgate.list
$ apk update
ERROR: wget: exited with error 8
WARNING: updating and opening https://.../aarch64_cortex-a53/x86_64/APKINDEX.tar.gz: unexpected end of file
1 unavailable, 0 stale; 11247 distinct packages available
```

Index URL that is not an apk index at all (`Packages.gz`):

```
ERROR: wget: exited with error 8
WARNING: updating and opening https://.../Packages.gz/x86_64/APKINDEX.tar.gz: unexpected end of file
```

Healthy `apk update` on the same userland (OpenWrt's own feeds reachable):

```
 [https://downloads.openwrt.org/releases/25.12.5/packages/x86_64/video/packages.adb]
OK: 11247 distinct packages available
```

## 3. Real install of the real v0.6.0 x86_64 package (local file, for shape)

Package, digest-checked against the release's own `SHA256SUMS`:

```
82d7dff1957ea6e17b1c5e352162c13eccab7678ba7e1be28eb43c2d04ffe8bb  tollgate-wrt_0.6.0_x86_64-x86-64.apk
82d7dff1957ea6e17b1c5e352162c13eccab7678ba7e1be28eb43c2d04ffe8bb   <- SHA256SUMS entry for that name: match
```

```
$ apk add --allow-untrusted --force-overwrite /pkgs/tg-x86.apk
  * uci: Entry not found
  * /etc/rc.common: ... can't create /var/lock/procd_tollgate-wrt.lock: nonexistent directory
  * Restarting network...
  * Warning: br-lan did not appear within 15 s
  * Note: br-private not found — this is expected on devices without a private network configured.
  * Reloading remaining services...
OK: 34.5 MiB in 158 packages

$ apk info -a tollgate-wrt
tollgate-wrt-0.6.0-r1 description:
TollGate turns an OpenWrt router into a Cashu-powered payment gateway for internet access. ...
tollgate-wrt-0.6.0-r1 webpage:
https://tollgate.me
tollgate-wrt-0.6.0-r1 installed size:
21 MiB
tollgate-wrt-0.6.0-r1 depends on:
ca-bundle
jq
libc
nodogsplash

$ grep -A8 '^P:tollgate-wrt' /lib/apk/db/installed
P:tollgate-wrt
V:0.6.0-r1
A:x86_64
T:TollGate turns an OpenWrt router into a Cashu-powered payment gateway ...
L:GPL-3.0-only
o:/feed/net/tollgate-wrt      <- LOCAL FILE, not a feed
```

Dependency behaviour with no network (negative control):

```
$ apk add --allow-untrusted --force-overwrite /pkgs/tg-x86.apk   # --network=none
...
ERROR: unable to select packages:
  jq (no such package):          required by: tollgate-wrt-0.6.0-r1[jq]
  nodogsplash (no such package): required by: tollgate-wrt-0.6.0-r1[nodogsplash]
```

## 4. `:2121` health from the real installed package

Started directly (`procd` absent in the container, so the init script cannot be
used):

```
$ mkdir -p /var/lock /var/run; /usr/bin/tollgate-wrt &
$ netstat -ltnp | grep 2121
tcp        0      0 :::2121                 :::*                    LISTEN      922/tollgate-wrt

$ wget -qO- http://127.0.0.1:2121/
{"kind":10021,"id":"36fb76f3774db5936af4ae8c3a42cda762aac68fdec6b46a70442da91d7de1bc",
 "pubkey":"d3e7c936d630be38582974e22ae69693c6f1375a44173ee2922bc7e4a3846d35",
 "created_at":0,
 "tags":[["metric","bytes"],["step_size","22020096"],["tips","1","2"],
         ["price_per_step","cashu","1","sat","https://mint.coinos.io","1"],
         ["price_per_step","cashu","1","sat","https://mint.minibits.cash/Bitcoin","1"],
         ["price_per_step","cashu","1","sat","https://mint.macadamia.cash","1"],
         ["price_per_step","cashu","1","sat","https://mint.westernbtc.com","1"],
         ["price_per_step","cashu","1","sat","https://kashu.me","1"], ...]}

service log:
INFO[...] Starting HTTP server on all interfaces...   module=main
INFO[...] UpstreamDetector event loop started         module=upstream_detector
```

## 5. Portal (2050) — cannot be measured in a container

```
$ wget -qO /dev/null http://127.0.0.1:2050/splash.html ; echo rc=$?
Failed to send request: Operation not permitted
rc=4

$ /etc/init.d/nodogsplash start
Failed to connect to ubus
Interface br-lan not detected.
Can not generate uci config. Will not start instance cfg015847.

$ ls /etc/tollgate/tollgate-captive-portal-site/
404.html asset-manifest.json assets balance.html favicon.ico locales
logo192.png logo512.png manifest.json splash.html
```

The package ships the portal tree correctly; the daemon needs procd/ubus and a
real `br-lan`. Recorded as **not measurable here**, not as a failure.

## 6. opkg (24.10.2) — the `src/gz` line, two spellings, both fail

Base **including** a trailing `/packages` segment:

```
$ echo "src/gz tollgate $U/packages" >> /etc/opkg/customfeeds.conf
$ opkg update
Downloading https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/packages/Packages.gz
*** Failed to download the package list from https://.../packages/Packages.gz
 * opkg_download: ... wget returned 8.
```

Base = the arch directory itself (index downloads, then is discarded):

```
$ echo "src/gz tollgate $U" >> /etc/opkg/customfeeds.conf
$ opkg update
Downloading https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.gz
Updated list of available packages in /var/opkg-lists/tollgate
Downloading https://.../Packages.sig
Signature file download failed.
Remove wrong Signature file.
 * opkg_download: ... Packages.sig, wget returned 8.

$ ls /var/opkg-lists/ | grep -i toll   ->  (nothing: list REMOVED)
$ opkg list | grep -i tollgate         ->  (nothing)
```

`option check_signature 0` in `/etc/opkg.conf` does **not** help on 24.10.2's
opkg — the `.sig` fetch is still attempted and the list is still removed:

```
$ printf "option check_signature 0\n" > /etc/opkg.conf
$ opkg update
Updated list of available packages in /usr/lib/opkg/lists/tollgate
Signature file download failed.
Remove wrong Signature file.
```

Note: `mkdir -p /var/lock /var/run` was required in the container before any
`opkg` run, else `opkg_conf_load: Could not create lock file /var/lock/opkg.lock`.
That is a container artifact, not a router behaviour.

## 7. Arch negative control

aarch64 package offered to an x86_64 host:

```
host arch=x86_64
$ apk add --allow-untrusted --force-overwrite /pkgs/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.apk
  tollgate-wrt-0.6.0-r1:
    error: uninstallable
    arch: aarch64_cortex-a53
    satisfies: world[tollgate-wrt><Q1LN9Gtu15lLSFH/1LVNfIqs+DiJs=]
```

Note the failure appears at `apk add`, not at update — and the `satisfies:` hash
is what a resolver prints when the only candidate is excluded by arch.

aarch64 asset, digest checked but not executable here:

```
128a5846be81ed97990981285a320725fd0bdadc8e8aff13af3de2935590c4c1  tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.apk
8591278 bytes
```

## 8. Board identity (upstream, not device)

```
$ curl -s .../releases/25.12.5/targets/mediatek/filogic/profiles.json
arch_packages: aarch64_cortex-a53
mt6000 profile: [{"model": "GL-MT6000", "vendor": "GL.iNet"}]
version_number: 25.12.5
```
