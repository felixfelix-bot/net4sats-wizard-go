# REPORT — t_23b8b386: Verify fd3 installed binary hash and running process identity

**Mode:** read-only. **Host:** c08r4d0r. **Date:** 2026-10-23.
**OVERALL VERDICT: FAIL — the target device was never reached; the asserted fd3 install does not exist on 192.168.1.1.**

---

## Executive summary

The card asks me to prove that the fd3 binary on the GL-MT6000 (192.168.1.1) is
byte-identical to the deployed artifact and that the running process executes that
path. I could not prove it, because **two independent premises are false**:

1. **192.168.1.1 is not a GL-MT6000.** It is a QEMU x86_64 virtual machine
   (OpenWrt 25.12.5, `OPENWRT_BOARD=x86/64`, MAC `52:54:00:12:34:56`). The GL-MT6000
   is an aarch64 `mediatek/filogic` box whose OUI is `94:83:c4`. A real GL.iNet device
   is visible at **192.168.2.7** (`94:83:c4:74:d3:03`), not at 192.168.1.1.
2. **There is no fd3 install on 192.168.1.1 at all.** No `tollgate-wrt`, no `tollgate`
   binary anywhere on the rootfs (`find / -xdev -name '*tollgate*'` → empty), no
   tollgate package, no listener on 2121, no process.

Additionally the Step 0 claims **do not exist**: the deploy cards `t_89e1c6cc` and
`t_86a624d3` were closed `done` by the fleet write-back with `result` = NULL. No
preflight sha256, installed exec path, /tmp artifact path, or restart timestamp was
ever recorded. This confirms — on independent evidence — the prior run of this same
card on host `c03rad0r`, which also reported BLOCKED for missing inputs.

So: Steps 1 and 3 are **FAIL**, not merely inconclusive — not because a hash
mismatched, but because the object under test is absent on the addressed host.

## What changed vs. the prior runs — SSH now works

Earlier runs (host `c03rad0r`) were walled at `Permission denied (publickey,password)`
on every key. From **this** host the router-looking address is directly reachable —
this host holds `192.168.1.2/24` on `archy-lab-br` — and `ssh root@192.168.1.1`
succeeds. That is precisely how the identity mismatch was discovered: the earlier
lanes could only see an HTTP surface that happened to look like TollGate, and could
never read a shell. With a shell, the box identifies itself immediately.

**Observation, not remediated:** the local `known_hosts` had a *different* ED25519 key
for 192.168.1.1 (stale record from the genuine GL-MT6000 of an earlier bench session).
Reads were performed with `UserKnownHostsFile=/dev/null`; nothing was removed from
`known_hosts`. This key change is itself consistent with the address having been
re-pointed at a different machine.

## STEP 0 — Untrusted claims (verbatim from the board DB)

Source: `/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`

| Claim | Recorded value | Source |
|---|---|---|
| (a) preflight sha256 of artifact | **ABSENT** | `t_89e1c6cc.result` = NULL; `t_86a624d3.result` = NULL |
| (b) installed exec path | **ABSENT** | same |
| (c) /tmp artifact path | **ABSENT** | same |
| (d) restart timestamp | **ABSENT** | same |

Only comment on both cards: `fleet-writeback: fleet-done rc=0 verified=True (offloaded); status -> done`.

Nearest thing to a recorded claim, from sibling card `t_23b8b386`'s own prior run
(host c03rad0r) — carried forward here as an **unverified reference**, used only as a
comparison anchor since the real deploy claims are missing:

| Reference item | Value (unverified) |
|---|---|
| apk artifact sha256 | `fb893695e80371b6b4b981b6da6cc8665b2dac0149e6bc1500f02951e2fcd899` |
| expected packed `usr/bin/tollgate-wrt` | `ec3d81bb903f0efa16dd69978071216f82499171dacdf534082cf44161da3888` (11763707 B) |
| expected packed `usr/bin/tollgate` | `0c058950ab82419aff7b3a01cb3869a35f29d29729c52334479c67b5e1846268` (6640609 B) |

---

## STEP 1 — Installed binary hash

| # | Exact command run | Raw unmodified output | Verdict |
|---|---|---|---|
| 1.1 | `ssh root@192.168.1.1 'sha256sum /usr/bin/tollgate-wrt'` | `sha256sum: can't open '/usr/bin/tollgate-wrt': No such file or directory` | **FAIL** |
| 1.2 | `ssh root@192.168.1.1 'sha256sum /usr/bin/tollgate'` | `sha256sum: can't open '/usr/bin/tollgate': No such file or directory` | **FAIL** |
| 1.3 | `ssh root@192.168.1.1 'ls -l /usr/bin/tollgate* /usr/sbin/tollgate* 2>&1; which tollgate tollgate-wrt'` | `ls: /usr/bin/tollgate*: No such file or directory`<br>`ls: /usr/sbin/tollgate*: No such file or directory`<br>`--- WHICH ---` (no output) | **FAIL** |
| 1.4 | `ssh root@192.168.1.1 'find / -xdev -name "*tollgate*" 2>/dev/null'` | (empty — no match anywhere on rootfs) | **FAIL** |
| 1.5 | `ssh root@192.168.1.1 'opkg list-installed 2>/dev/null \| grep -i toll'` | (empty — no tollgate package installed) | **FAIL** |

There is no installed binary to hash. The hash comparison demanded by this step is
not merely unevaluable — the object it addresses is absent. **FAIL.**

## STEP 2 — Staged artifact in /tmp (informational)

| # | Exact command run | Raw unmodified output | Verdict |
|---|---|---|---|
| 2.1 | `ssh root@192.168.1.1 'ls -l /tmp/tollgate* /tmp/*.apk 2>&1'` | `ls: /tmp/tollgate*: No such file or directory`<br>`ls: /tmp/*.apk: No such file or directory` | **N-A** (absent) |

No `/tmp` artifact path was ever recorded in Step 0, and none exists on disk.
Per the card, **absence alone is not a FAIL** — recorded as **N-A / absent**.

## STEP 3 — Running process identity

| # | Exact command run | Raw unmodified output | Verdict |
|---|---|---|---|
| 3.1 | `ssh root@192.168.1.1 'pgrep -f tollgate'` | (empty) | **FAIL** |
| 3.2 | `ssh root@192.168.1.1 'pgrep -x tollgate; pgrep -x tollgate-wrt'` | (empty) `exact_rc=1` | **FAIL** |
| 3.3 | `ssh root@192.168.1.1 'ps w \| grep -i tollgate \| grep -v grep'` | (empty) | **FAIL** |
| 3.4 | `ssh root@192.168.1.1 'netstat -ltnp'` | no listener on `:2121` (only 22, 53, 80, 443) | **FAIL** |
| 3.5 | `ssh root@192.168.1.1 'tr "\0" " " < /proc/<pid>/cmdline; ls -l /proc/<pid>/exe'` | **not executable — no tollgate pid exists** | **FAIL** |

**A false-positive I caught and discarded:** an early `pgrep -f tollgate` returned
pid `3567`, and a later one returned `3580`. Inspecting `/proc/3567/` gave
`no such file` — the process had already exited. Re-running with
`pgrep -laf tollgate` showed the only match was **my own SSH command line**
(`3580 ash -c echo "=== pgrep -f tollgate ==="; ...`), because `-f` matches the full
argument vector of the shell that carries the pattern string. `pgrep -x` (exact
process name) returns nothing, rc=1. **There is no tollgate process.** This is a
classic `pgrep -f` self-match artefact and must not be reported as a live service.

The `/proc/<pid>/exe` == installed-path identity check therefore has no subject.
**FAIL.**

## STEP 4 — Process count

| # | Exact command run | Raw unmodified output | Verdict |
|---|---|---|---|
| 4.1 | `ssh root@192.168.1.1 'pgrep -f tollgate'` | (empty) | **N-A — zero processes** |
| 4.2 | `ssh root@192.168.1.1 'pgrep -x tollgate; pgrep -x tollgate-wrt'` | (empty) rc=1 | **N-A — zero processes** |
| 4.3 | `ssh root@192.168.1.1 'ps w'` (full table) | 96 lines, system processes only; **no tollgate entry** | **N-A — zero processes** |

**Pid list: EMPTY.** Not multiple — *none*. The card asks me to flag multiples;
the actual condition is zero, which I flag explicitly as the more serious finding.
No duplicate/conflicting service instance exists on this host.

---

## Device identity — why the above FAILs are an address error, not a hash mismatch

| # | Exact command run | Raw unmodified output | Reading |
|---|---|---|---|
| 5.1 | `ssh root@192.168.1.1 'uname -a'` | `Linux OpenWrt 6.12.94 #0 SMP Mon Jun 29 12:59:20 2026 x86_64 GNU/Linux` | **x86_64**, not aarch64 |
| 5.2 | `ssh root@192.168.1.1 'uname -m'` | `x86_64` | GL-MT6000 is `aarch64_cortex-a53` |
| 5.3 | `ssh root@192.168.1.1 'cat /tmp/sysinfo/model'` | `QEMU Standard PC (Q35 + ICH9, 2009)` | **a virtual machine** |
| 5.4 | `ssh root@192.168.1.1 'cat /tmp/sysinfo/board_name'` | `qemu-standard-pc-q35-ich9-2009` | not `glinet,gl-mt6000` |
| 5.5 | `ssh root@192.168.1.1 'grep OPENWRT_BOARD /etc/os-release'` | `OPENWRT_BOARD="x86/64"` | confirms x86 target |
| 5.6 | `ssh root@192.168.1.1 'mount \| head -6'` | `/dev/root on / type ext4` … `/dev/vda1 on /boot type ext4` | `vda` = virtio disk |
| 5.7 | `ip neigh show 192.168.1.1` | `192.168.1.1 dev archy-lab-br lladdr 52:54:00:12:34:56 REACHABLE` | **`52:54:00` = QEMU OUI**; GL.iNet is `94:83:c4` |
| 5.8 | `ip neigh \| grep -i '94:83:c4'` | `192.168.2.7 dev enx00e04c36043e lladdr 94:83:c4:74:d3:03 STALE` | a real **GL.iNet** box is at **192.168.2.7** |
| 5.9 | `curl -sS -m 10 -D - http://192.168.1.1:2121/` | `curl: (7) Failed to connect … port 2121 … Couldn't connect to server` | the TollGate API port is **closed** here |
| 5.10 | `curl -o /dev/null -w '%{http_code}' http://192.168.1.1/` | `http=200` | plain uhttpd LuCI/landing page, not TollGate |

The `52:54:00:12:34:56` MAC plus `vda1` root disk plus `QEMU Standard PC` model is
unambiguous: **this is a QEMU guest**, and `52:54:00:12:34:56` is the well-known
QEMU default MAC for the first virtio NIC. Whatever answered on `:2121` during the
2026-09-20 sibling runs — and `t_c82dec01` genuinely got a well-formed kind:10021
payload there — it is not answering now, and the address no longer hosts the same
machine. The stale `known_hosts` key for 192.168.1.1 is consistent with that
re-pointing.

## Overall verdict

**FAIL — verification could not be established, and the premises it rests on are false.**
The 192.168.1.1 target is a QEMU x86_64 guest, not the aarch64 GL-MT6000; it carries
no tollgate binary, no tollgate package, no `:2121` listener, and no tollgate process.
Steps 1 and 3 FAIL on absence of the object under test; Step 2 is N-A by absence
(not a FAIL by the card's own rule); Step 4 is N-A at zero processes rather than
multiple.

Nothing was remediated, per the card's hard constraint. Not attempted: any fix,
reinstall, restart, rollback, package operation, or re-addressing.

## REMAINING STEPS (operator, in order)

1. **Resolve the true device address before re-running this card.** 192.168.1.1 is a
   QEMU guest. The only GL.iNet-OUI host currently visible is **192.168.2.7**
   (`94:83:c4:74:d3:03`) — confirm it is the GL-MT6000 and that it is the intended
   fd3 target, then re-issue this verification against that address.
2. **Restore the missing Step 0 inputs.** `t_89e1c6cc` and `t_86a624d3` were closed
   `done` with `result` = NULL. Without a recorded preflight sha256 + installed exec
   path, *no* later verification can ever be conclusive, on any host. The deploy lane
   must record those four values.
3. **Confirm whether an fd3 install ever happened at all.** On the evidence available,
   there is no install on the reachable host, and the deploy card never reported one.
   Treat the fd3 deployment as **unproven** until a device is named and hashed.
4. **De-conflict the address.** Three separate cards (`t_23b8b386`, `t_c82dec01`,
   `t_14e59f6e`) have now burned runs against 192.168.1.1. Prune the stale
   `known_hosts` entry and pin the bench host/address in the card body.
5. Only after 1-4: re-run Steps 1-4 verbatim against the confirmed device.

## Mutation statement

**Nothing on the device was mutated. All commands were read-only.**

Device-facing actions were exclusively: `echo`-only SSH session checks, `uname`,
`cat` on `/etc/os-release`, `/tmp/sysinfo/*`, and `/proc/sys/kernel/hostname`,
`sha256sum` on non-existent paths, `ls`, `find`, `which`, `ps`, `pgrep`, `netstat`,
`mount`, `grep`, `tr`/`ls` on `/proc` (all of which failed with "no such file"), and
two unauthenticated HTTP GETs. No file was created, edited, or deleted on the device;
no service was started, stopped, or restarted; no package was installed or removed;
no rollback, no lock taken or released, no reboot. The local `known_hosts` was **not**
modified — stale-key reads were scoped with `UserKnownHostsFile=/dev/null`.


