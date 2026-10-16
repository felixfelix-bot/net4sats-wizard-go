# REPORT — t_53879c67: FE2b pre-baseline capture, GL-MT6000

**Date:** 2026-10-16
**Task:** Capture GL-MT6000 current package state over SSH; record the pre-FE2b baseline.
**Result: BLOCKED — device reachable, root access denied. Steps 2-4 captured NO data.**
**Evidence:** `verify/fe2b-baseline/` (`README.md`, `auth-matrix.log`, `key-attempts.log`, `lastmile.log`)

---

## 1. One-line summary requested by the brief

**Installed? UNKNOWN. Origin? UNKNOWN. Package-manager version? UNKNOWN.** — not measured,
because `root@192.168.1.1` denied every available credential. Nothing below is inferred; each
statement is tied to raw output in the logs.

The brief's escape hatch was: *"If SSH is unreachable, STOP and report to the manager instead of
guessing."* The router is **up and reachable** — the thing that is unreachable is
*authenticated* access to it. I applied the same instruction and stopped rather than guessing.

---

## 2. Step 1 — SSH reachability: PASS (raw result)

Access path: `c08r4d0r` -> `ssh cobrador` -> `ssh root@192.168.1.1`.

cobrador is the only fleet node with a route to the bench segment. dq05 was probed
(`192.168.1.1`, `192.168.8.1`, `192.168.10.1`, `192.168.21.1` — all closed) and has no bench
LAN; the whole `192.168.2.0/24` was swept from this host and contains **no dropbear host at
all** — only Ubuntu OpenSSH fleet nodes and the home gateway. Had I trusted this host alone I
would have reported "not on this LAN"; the device was found only by probing from cobrador.

```
OPEN   192.168.1.1:22     banner: SSH-2.0-dropbear
OPEN   192.168.1.1:80
OPEN   192.168.1.1:2121   (tollgate backend API)
OPEN   192.168.1.1:2050   (captive portal)
OPEN   192.168.1.1:8080
```

## 3. Step 2 — the read-only command set: NOT RUN

`/etc/apk/repositories`, `apk list --installed | grep -i tollgate`, the `opkg` equivalent,
`apk --version` / `opkg --version`, and `uname -a` were **never executed**. There is no device
output to attach because no session was ever established.

Authentication, all methods available on cobrador, each tried individually:

```
A  -i ~/.ssh/id_tollgate_router   ->  root@192.168.1.1: Permission denied (publickey,password)
B  -i ~/.ssh/id_ed25519           ->  root@192.168.1.1: Permission denied (publickey,password)
C  -i ~/.ssh/tollgate             ->  root@192.168.1.1: Permission denied (publickey,password)
D  -i ~/.ssh/id_rsa               ->  root@192.168.1.1: Permission denied (publickey,password)
E  sshpass -p ''      (password only)  -> no working session
F  sshpass -p tollgate (password only) -> no working session
```

Method A failing is the load-bearing result: `~/.ssh/id_tollgate_router` exists *specifically*
for this device (created 2026-09-23) and it is refused, so this is not a wrong-key selection
problem. The server offers `publickey,password`; neither category is satisfied by anything on
cobrador.

### 3a. The host key on record is stale — a trap that disguises itself as "bad credentials"

```
live  192.168.1.1 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHVQW4eljQFAE6fvIU/fzFZxX+JDz6xW9Pad5y2X8+hl
      SHA256:323HdkO6Ys6FKMLRb53HHDpXR36ZhJWR2ANA3V5mdWU
stored  cobrador:~/.ssh/known_hosts line 116 -> a DIFFERENT ed25519 key
```

A plain `ssh root@192.168.1.1` on cobrador aborts on `REMOTE HOST IDENTIFICATION HAS CHANGED`
*before* offering any key. Anyone reading that output would conclude "SSH is unreachable" or
"credentials blocked" — the device was in fact up on all five ports. Re-keying means it was
reflashed or its keys regenerated since that line was written, which is also consistent with
`AuthorizedKeysFile` pointing at the `/tmp` tmpfs: any previously authorized key would be
discarded by a reboot, which alone explains every `Permission denied` above.

## 4. Steps 3 and 4 — cannot be answered without access

- **Is `tollgate-wrt` installed, and is its origin a feed or a local `.apk`?** UNKNOWN. Requires
  `apk info -a tollgate-wrt` (the `origin:` field) or the raw
  `/lib/apk/db/installed` entry — neither was run.
- **Is there a pre-existing tollgate entry in `/etc/apk/repositories`?** UNKNOWN, and see §5:
  that file may be under active mutation right now, which makes reading it *before* the window
  is quiet unsafe as a baseline.

## 5. The finding that matters more than the missing baseline: **FE2b is already in flight**

While checking for an existing authenticated session, I found one — and it is doing this task's
own work, for a version this task was never told about:

```
ESTAB  192.168.1.106:49798 -> 192.168.1.1:22   users:(("ssh",pid=1150196,fd=3))

sshpass -p '' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  -o PreferredAuthentications=password root@192.168.1.1 \
  'cd /tmp
   FILE=tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk
   apk update 2>&1 | tail -2
   echo === installing ===
   apk add --allow-untrusted --force-overwrite /tmp/$FILE 2>&1 | tail -8
   echo === installed? ===
   apk list -I 2>/dev/null | grep tollgate-wrt'
```

| fact | value |
|---|---|
| started | Fri Oct 16 16:55:01 2026 |
| still established at | 15:39Z (17:39 local) — ~45 min idle |
| parent | pid 1150195, `ppid=1` (bare `bash -c` string, reparented to init) |
| last hop auth | sshpass with an empty value, plus PreferredAuthentications=password |
| target artifact | `tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk` |

Three consequences:

1. **The window is not exclusive.** The card mandates serialization against fd4/wz18; a third
   actor is already on the device.
2. **The state to be recorded as "pre-FE2b" may not be pre-FE2b.** `/tmp/tollgate-wrt_..._rc6....apk`
   is staged and `apk update` may already have run. A baseline taken now risks recording a
   half-applied rc6 — worse than no baseline, because it would be *authoritative-looking*.
3. **The stalled auth is the same one that fails for us** (`sshpass -p ''`, password-only),
   which is consistent with the session being **stuck at the password prompt**, not mid-write.
   That is the reassuring reading — but it is a reading, not a measurement, so I did not treat
   it as permission to proceed.

Respecting the delegated-destructive-ops rule, I did **not** attach to, kill, share, or race
that session, and I ran no `apk` command and made no write to the device. Racing another
operator's `apk` transaction is how a package database is left dirty.

## 6. Evidence that survives without SSH (explicitly NOT the baseline)

`:2121/` answers **200** with a well-formed `kind:10021` advert — FULL mode, 7 mints
(`mint.coinos.io`, `mint.minibits.cash`, `mint.macadamia.cash`, `mint.westernbtc.com`,
`kashu.me`, `mint.cubabitcoin.org`), `created_at: 0`, `step_size 22020096`. `:2050/splash.html`
answers 200; `:8080/` 307-redirects to `https://192.168.1.1/`.

So a tollgate backend **is running and serving traffic**. This does **not** establish which
version it is, nor whether it arrived from a feed or a local file — which is precisely the
question the capture was meant to answer. Anyone reading this report should not mistake it for
the baseline.

## 7. Card precondition: also unmet

`t_53879c67` states *"Precondition: manager merged FE2a PR + release tagged."*
Card `t_df4304d5` (**FE2a**: "Feed PR#1 rebase + artifact-to-release bridge + arch tuples") is
still **todo**, not merged. The most recent published `FreedomTechFeed/packages` release is
**v0.6.0-rc6 (2026-10-15T19:50Z)**. This task is therefore gated on FE2a independently of the
access problem.

## 8. To unblock (operator/card-owner, in order)

1. **Clear the stalled cobrador session** (pid 1150195/1150196) or confirm it has exited, so the
   FE2b window is exclusive as the card requires.
2. **Repair router SSH auth:** install a known public key into the dropbear authorized-keys file
   at a path that survives reboot, then hand the worker the key path — or supply the credential
   through a channel that is not a shell argument.
3. **Re-run this capture** with a quiet window and working auth. Steps 2-4 are four read-only
   commands and take under a minute; the probe script is already written and needs only the
   auth to succeed.

Optional, and worth doing while the window is quiet: confirm whether the dropbear authorized-keys
file is a tmpfs symlink, since that would turn every future reboot into a silent lockout.

---

## 9. Honesty ledger

- Steps 2, 3, 4: **not performed**. No device output exists for them and none is invented here.
- No commit, push, or test is claimed in this report that was not observed to succeed.
- The `:2121`/`:2050`/`:8080` responses are real observations and are labelled as *not* the
  baseline.
- The "stuck at password prompt" reading of the stalled session is labelled as a reading.
