# PROGRESS — t_23b8b386: Verify fd3 installed binary hash and running process identity (read-only)

## Finding -> status -> files touched

### F1 — STEP 0 (gather untrusted claims): DONE
- Deploy card `t_c1d989ab` (per card body) is itself the parent verification umbrella; its
  children are t_23b8b386 / t_c82dec01 / t_14e59f6e / t_ce52427f. `result` = NULL.
- The actual deploy/install cards: `t_89e1c6cc` (deploy) and `t_86a624d3` (install) —
  BOTH status=done with **`result` = NULL (empty)**; only comment is
  `fleet-writeback: fleet-done rc=0 verified=True (offloaded); status -> done`.
- Consequence: (a) preflight sha256, (b) installed exec path, (c) /tmp artifact path,
  (d) restart timestamp were **never recorded** in the board DB. Searched board DB
  (tasks, task_comments), plus prior sibling reports.
- Partial claims DO exist from sibling cards (also unverified):
  - t_23b8b386 (prior run, host c03rad0r): artifact apk sha256 fb893695e80371b6b4b981b6da6cc8665b2dac0149e6bc1500f02951e2fcd899;
    expected packed payloads usr/bin/tollgate-wrt = ec3d81bb903f0efa16dd69978071216f82499171dacdf534082cf44161da3888 (11763707 B),
    usr/bin/tollgate = 0c058950ab82419aff7b3a01cb3869a35f29d29729c52334479c67b5e1846268 (6640609 B).
  - t_86a624d3: reports STOPPED at preflight gate, nothing installed.
- STATUS: claims are INCOMPLETE/MISSING — proceed to Steps 1-4 and report what the
  device actually holds.

### F2 — Device reachable + SSH works from THIS host (c08r4d0r), unlike prior runs
- Prior runs on host `c03rad0r` were walled: `Permission denied (publickey,password)`.
- This host has a direct interface on 192.168.1.0/24 (`archy-lab-br` = 192.168.1.2/24).
- `ssh root@192.168.1.1` → OK. Device is **OpenWrt 25.12.5 r33051-f5dae5ece4, SMP, 6.12.94**.
- Note: existing known_hosts entry for 192.168.1.1 had a DIFFERENT ED25519 key (stale
  record from the earlier box). Used `UserKnownHostsFile=/dev/null` for reads; this is
  flagged as an observation, not a target of this card.
- STATUS: unblocked. Proceeding.

### Steps 1-4: see REPORT.md for the full evidence table.

### F3 — STEP 1 (installed binary hash): **FAIL**
- `sha256sum /usr/bin/tollgate-wrt` → `No such file or directory`
- `sha256sum /usr/bin/tollgate` → `No such file or directory`
- `find / -xdev -name '*tollgate*'` → empty. `opkg list-installed | grep -i toll` → empty.
- No installed binary exists, so no hash comparison is possible.

### F4 — STEP 2 (/tmp artifact, informational): **N-A (absent)**
- `/tmp/tollgate*` and `/tmp/*.apk` → both No such file. No path was ever recorded either.
- Card states absence alone is not a FAIL → N-A.

### F5 — STEP 3 (process identity): **FAIL**
- `pgrep -f tollgate` → empty. `pgrep -x tollgate`/`tollgate-wrt` → empty, rc=1.
- `ps w` (96 lines) → no tollgate entry. `netstat -ltnp` → no :2121 listener.
- /proc/<pid>/exe check not executable: no tollgate pid.
- **Trap caught:** an early `pgrep -f tollgate` returned 3567 then 3580; /proc/3567
  was already gone, and `pgrep -laf` showed 3580 was MY OWN ssh ash command line.
  `pgrep -f` self-matches the wrapper shell. Not a live service. Discarded.

### F6 — STEP 4 (process count): **N-A — ZERO processes** (not multiple)
- Pid list EMPTY. Flagged explicitly: the risk here is zero instances, not duplicates.

### F7 — ROOT CAUSE: 192.168.1.1 is not the GL-MT6000
- `uname -m` → `x86_64`; `/tmp/sysinfo/model` → `QEMU Standard PC (Q35 + ICH9, 2009)`;
  `board_name` → `qemu-standard-pc-q35-ich9-2009`; `OPENWRT_BOARD="x86/64"`;
  root on `/dev/root`, `/dev/vda1` on /boot (virtio); MAC `52:54:00:12:34:56` (QEMU OUI).
- GL-MT6000 is aarch64 cortex-a53; GL.iNet OUI is `94:83:c4`.
- A real GL.iNet host is at **192.168.2.7** (`94:83:c4:74:d3:03`) — candidate target.
- :2121 refused; :80 = plain 200 (uhttpd), so the earlier TollGate-looking HTTP surface
  was on a different occupant of this address.
- Local known_hosts had a stale/different ED25519 key for 192.168.1.1 → consistent with
  the address having been re-pointed at a different machine. NOT modified.

### F8 — Nothing mutated
- Read-only: uname/cat/ls/find/which/ps/pgrep/netstat/mount/grep/tr-on-proc + 2 unauth GETs.
- No writes, no restarts, no package ops, no rollback, no lock taken/released, no reboot.

## Files touched
- `PROGRESS-t_23b8b386.md` (this file)
- `REPORT-t_23b8b386.md` (full evidence table + verdict)

## Status: COMPLETE — verdict FAIL, reported, nothing remediated.

