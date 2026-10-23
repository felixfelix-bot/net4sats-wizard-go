# PROGRESS — t_877ad4ef: RED baseline feed check (FE2b)

Finding -> Status -> Files touched. Written 2026-10-23.

1. Finding: Card t_877ad4ef is mine ("Write and run failing feed check script,
   capture baseline", board net4sats-mvp-v2, status ready). The board also holds
   t_53879c67 (FE2b consumption) BLOCKED and t_18d15abe (preflight) which already
   wrote a `scripts/fe2b_feed_check.sh` in commit d5d2373.
   Status: CONFIRMED against the board DB (not the brief).
   Files touched: none.

2. Finding: The GL-MT6000 has NO IP path from this host — TCP 22, 2121 and 2051
   are all closed on 192.168.1.1, and `ping` is not even installed here. The
   prior lane's ssh-only script therefore CANNOT run on this bench; running it
   yields "cannot reach the router", which is not a RED baseline.
   Status: ROOT CAUSE of why the script had to be rewritten, not re-run.
   Files touched: none.

3. Finding: The device IS alive on USB UART. All four ttyUSB ports were free
   (no flock holder, no fd-holding process), and I identified the board by its
   OWN prompt rather than by port: /dev/ttyUSB0 -> root@tollgate-0195,
   /tmp/sysinfo/board_name = glinet,gl-mt6000, OpenWrt 24.10.2 r28739,
   mediatek/filogic, aarch64_cortex-a53. ttyUSB1 answered as a COMFAST login,
   so ports are NOT interchangeable — the prompt is the identifier.
   Status: RESOLVED. Serial-first transport implemented.
   Files touched: scripts/fe2b_feed_check.sh.

4. Finding: Device state re-verified first-hand (read-only) and it matches the
   package-state task: opkg (no apk at all), tollgate-wrt 0.6.0_rc6-r1
   installed, /etc/opkg/customfeeds.conf:6 carries the exact src/gz line, a
   local /tmp/tollgate-wrt.ipk is still staged, /var/opkg-lists/ has no tollgate
   index, service RUNNING and rc.d-enabled, :2121 -> 200, :2051/splash.html ->
   200 while :2051/ -> 403 by design.
   Status: CONFIRMED — this is what makes A2 decidable.
   Files touched: none.

5. Finding: A2 (origin == feed) needs care: opkg 1.x on OpenWrt 24.10 writes no
   Source/Filename field to /usr/lib/opkg/status, so status cannot prove origin.
   Implemented as three falsifiable conditions: no staged local package, a
   tollgate index exists in /var/opkg-lists/, and the installed version equals
   the version that index advertises.
   Status: DONE.
   Files touched: scripts/fe2b_feed_check.sh.

6. Finding: The prior script asserted portal 200 on `:2051/`, which is 403 BY
   DESIGN (doc root has no index.html, listing disabled). That assertion could
   never turn green, so it was not a valid gate. Now asserts /splash.html, with
   the 403 recorded as a note.
   Status: CORRECTED (deviation from the predecessor script, documented in-file).
   Files touched: scripts/fe2b_feed_check.sh.

7. Finding: RED baseline captured and reproducible. 3 consecutive runs: 5 PASS /
   3 FAIL, exit 1, PASS/FAIL set byte-identical. Control run with a deliberately
   wrong health port + portal path flips A4/A5 to FAIL (3 PASS / 5 FAIL), proving
   the green checks are falsifiable rather than vacuous.
   Status: DONE — this is the deliverable.
   Files touched: verify/fe2b-baseline/red-run-2026-10-23.txt,
   verify/fe2b-baseline/control-run-2026-10-23.txt.

8. Finding: Root PROGRESS.md and REPORT.md in this shared worktree belong to
   other cards (they carry Gate-8 / FE2b-preflight / WZ content and HEAD is a
   Gate-8 lane). Overwriting them would destroy another card's record.
   Status: DEVIATION — wrote task-scoped PROGRESS-t_877ad4ef.md /
   REPORT-t_877ad4ef.md at the worktree root instead of overwriting.
   Files touched: as named.

9. Finding: Committed 80c8d56 and pushed to a NEW branch
   worker-heavy/t_877ad4ef (remote sha verified via ls-remote). The shared HEAD
   branch belongs to the Gate-8 lane, so pushing there would have mixed my commit
   into another workstream.
   Status: DONE.
   Files touched: none (branch-side).
