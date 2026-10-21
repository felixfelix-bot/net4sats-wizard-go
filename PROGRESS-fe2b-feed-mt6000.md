# PROGRESS — FE2b: install tollgate-wrt from the published FEED on GL-MT6000

Finding -> Status -> Files touched. Written 2026-10-21.

1. Finding: The brief's premise ("parent preflight recorded the exact feed base URL
   and `src/gz` line") is FALSE. The parent card `t_1b1cf42f` recorded the opposite:
   no feed index exists, so no `src/gz` line could be derived (3 runs, 3 independent
   negative proofs). No FEED_BASE_URL was ever recorded.
   Status: CONFIRMED independently this run, then proceeded anyway with the canonical
   candidate base so the card's literal steps were executed rather than refused.
   Files touched: none.

2. Finding: The device IS reachable this run (prior cards recorded it ping-dead and
   credential-denied). `root@192.168.1.1` accepts the default key.
   Status: RESOLVED a long-standing blocker.
   Evidence: `ssh root@192.168.1.1` -> OK; `cat /tmp/sysinfo/board_name` -> glinet,gl-mt6000.

3. Finding: The parent's open conflict (§10.3 "which package manager?") is SETTLED:
   GL-MT6000 runs OpenWrt 24.10.2, `opkg` only, `apk: not found`.
   Status: RESOLVED. The `apk update` / `apk add` / `apk info` steps in the brief are
   unrunnable on this device — they are transcribed from the 25.12/apk lane.
   Files touched: none.

4. Finding: tollgate-wrt 0.6.0_rc6-r1 is ALREADY installed, and it was installed from
   a local flat file — `/tmp/tollgate-wrt.ipk`, sha256 `3b3a85cc...e26a`, byte-identical
   to the release's flat `.ipk` asset (8630160 B). `Installed-Time` 2026-10-20 16:40:47,
   predating this run. No tollgate index exists in `/var/opkg-lists/`.
   Status: CONFIRMED — i.e. the device is currently in exactly the state the card
   forbids ("NOT via local scp of a .apk/.ipk").
   Files touched: none.

5. Finding: No feed index is published, in any form, at any candidate base. 15 index
   path shapes probed -> all 404. Flat assets 200. `FreedomTechFeed/feed` is a
   `src-git` BUILD-TIME feed (`net/tollgate-wrt/Makefile`, zero tags/releases), and its
   index generator lives on an unmerged `ci/auto-feed-indices` branch; the
   `Packages.gz`/`packages.adb` in its strategy doc are labelled "NEW" = proposals.
   Status: BLOCKED (structural, upstream). This is the root cause; step 3 cannot pass.
   Files touched: none.

6. Finding: Device has no usable WAN. Only uplink is STA `phy0-sta0` -> 10.2.1.122
   (SSID "StarGate"); it reaches 10.2.0.1 but nothing beyond (ping 1.1.1.1 100% loss,
   DNS NXDOMAIN, direct-IP HTTPS timeout). This is NOT the firewall (phy0-sta0 is in
   the wan zone, output policy accept) — the upstream itself does not route.
   Status: BLOCKED (environmental). Even a real feed could not be fetched from here.
   Files touched: none.

7. Finding: Feed line added idempotently to /etc/opkg/customfeeds.conf, backup taken,
   persistence verified in overlay upper, one-command revert left on device.
   Status: DONE (the only step of the brief that could be completed as specified).
   Files touched: device `/etc/opkg/customfeeds.conf`,
   `/root/customfeeds.conf.pre-fe2b.bak`, `/root/revert-fe2b-feed.sh`;
   repo `PROGRESS-fe2b-feed-mt6000.md`, `REPORT-fe2b-feed-mt6000.md`.

8. Finding: Root `REPORT.md`/`PROGRESS.md` in this worktree belong to other cards
   (WZ1.9 atomicity, testid gap analysis). Overwriting them would destroy another
   card's record.
   Status: DEVIATION — wrote card-scoped files at the worktree root instead of
   overwriting. Deliverable is at the root as asked, just namespaced.
   Files touched: as in (7).
