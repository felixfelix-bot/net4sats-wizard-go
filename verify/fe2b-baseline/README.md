# FE2b pre-baseline: GL-MT6000 package state — NOT CAPTURED (task blocked)

**Card:** `t_53879c67` — "FE2b: Feed consumption install on GL-MT6000" (status: blocked)
**Captured:** 2026-10-16T15:38Z
**Result: STOPPED at step 1 of the brief — no root command on the device was executed.**
The brief's own escape hatch applies: "If SSH is unreachable, STOP and report to the manager
instead of guessing." The router is reachable; *authenticated access to it is not*, and it is
reachable-but-contended by another concurrent FE2b attempt. No package state was recorded.

## Files here

| file | content |
|---|---|
| `auth-matrix.log` | TCP reachability + host key + all 4 documented auth methods |
| `key-attempts.log` | every candidate router key on cobrador, tried individually |
| `lastmile.log` | running ssh sessions, and the router's HTTP surfaces (evidence that survives without SSH) |

## The honest one-line summary the brief asked for

**Installed?** UNKNOWN — not measured. **Origin?** UNKNOWN. **Manager version?** UNKNOWN.
Everything else in this directory is the *reason* it is unknown, which is itself real evidence.

## What WAS measured (raw output in the logs)

**Step 1 — reachability: PASS.** Host `192.168.1.1`, reached as `ssh cobrador` -> `ssh root@192.168.1.1`
(cobrador is the only fleet node on the bench segment; dq05 has no route to it).

```
192.168.1.1:22   OPEN   banner SSH-2.0-dropbear
192.168.1.1:80   OPEN
192.168.1.1:2121 OPEN   <- tollgate backend API
192.168.1.1:2050 OPEN   <- captive portal
192.168.1.1:8080 OPEN
```

**Step 2 — the read-only command set: NOT RUN.** `root@192.168.1.1` denied every credential:

```
A dedicated key ~/.ssh/id_tollgate_router  -> Permission denied (publickey,password)
B ~/.ssh/id_ed25519                        -> Permission denied (publickey,password)
C ~/.ssh/tollgate                          -> Permission denied (publickey,password)
D ~/.ssh/id_rsa                            -> Permission denied (publickey,password)
```

The dedicated key that exists *specifically* for this device fails, so this is not a
"wrong key selected" problem. `sshpass` is present and the documented empty-password and
lab-default (`tollgate`) paths do not produce a working session either. Note that
`AuthorizedKeysFile` may live on the `/tmp` tmpfs, in which case every passwordless
authorized key would be discarded by a reboot.

**Additional finding — the stored host key no longer matches.** The key currently served is
`SHA256:323HdkO6Ys6FKMLRb53HHDpXR36ZhJWR2ANA3V5mdWU`, while cobrador's `known_hosts:116`
holds a different ed25519 key. The device has been reflashed (or otherwise re-keyed) since
that line was written. This is why a plain `ssh root@192.168.1.1` on cobrador dies on the
host-key warning *before* it ever offers a key — a trap that reads as "credentials blocked"
when the real question is authentication.

## The finding that matters more than the baseline: FE2b is already running

`lastmile.log` caught a still-connected session on cobrador attempting exactly the FE2b
install, for a package version this task was never told about:

```
ESTAB  192.168.1.106:49798 -> 192.168.1.1:22   users:(("ssh",pid=1150196,fd=3))

sshpass -p '' ssh ... PreferredAuthentications=password root@192.168.1.1 \
  'cd /tmp; FILE=tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.apk;
   apk update; apk add --allow-untrusted --force-overwrite /tmp/$FILE;
   apk list -I | grep tollgate-wrt'
```

- Started **16:55:01**; observed still established at **15:39Z (17:39 local)** — ~45 min idle.
- The last hop is `sshpass -p ''` + `PreferredAuthentications=password`, i.e. exactly the
  method that fails for us, which is consistent with it being **stalled at the password
  prompt** rather than mid-transaction.
- It is not a tty someone is watching: the parent is a bare `bash -c` string reparented to
  init (pid 1150195, ppid 1).

This changes the task's premise. The brief frames FE2b as a future change whose *pre*-state
must be recorded. But an attempt to mutate this device is already open, so:
1. **the window is not exclusive** despite the card requiring serialization, and
2. **the current device state may already be mid-transition**, which is precisely the state
   that is unsafe to record as "the baseline".

Per the delegated-destructive-ops rule, I did not attach to, kill, or share that session, and
I did not retry auth against a device another operator is mid-transaction on. Racing it on
`apk` is how a package DB is left dirty.

## Why this matters for FE2b's correctness

The card's premise is that the device is sitting at a known pre-FE2b state. The rc6 artifact
in that stalled command means a **partial rc6 footprint may already be staged** (`/tmp/...apk`
staged, `apk update` possibly run against a feed list that may or may not include the FE2b
tollgate feed line). Two consequences:

- A "baseline" captured now could record a **half-applied rc6** state and be wrong.
- The step-4 answer the brief wanted (*is there already a tollgate feed entry in
  `/etc/apk/repositories`?*) is **itself the thing under active mutation**, so it must be read
  after the window is confirmed quiet, not before.

## Evidence that survived without SSH

`http://192.168.1.1:2121/` answers **200** with a valid `kind:10021` advert, FULL mode, 7 mints
(`mint.coinos.io`, `mint.minibits.cash`, `mint.macadamia.cash`, `mint.westernbtc.com`,
`kashu.me`, `mint.cubabitcoin.org`). `:2050/splash.html` answers 200; `:8080` 307-redirects to
`https://192.168.1.1/`. So a tollgate backend **is** running and serving — but this says nothing
about *which version* or whether it came from a feed or a local file, which is the whole point
of the requested capture. Do not read this as the baseline.

## To unblock (operator/card-owner actions, in order)

1. Confirm or clear the stalled session on cobrador (pid 1150195/1150196) so the window is
   exclusive — the card already requires FE2b to be serialized against fd4/wz18.
2. Rotate/repair router SSH auth: install a known public key in
   `/etc/dropbear/authorized_keys` (or a persistent `/etc/dropbear/authorized_keys` path that
   survives reboot), then hand the key path to the worker.
3. Only then re-run this capture: with a quiet window and working auth, steps 2-4 are a
   few read-only commands and take under a minute.

## Card precondition also unmet

`t_53879c67` states "Precondition: manager merged FE2a PR + release tagged."
Card `t_df4304d5` (FE2a) is still **todo** — not merged. The last published
`FreedomTechFeed/packages` release is **v0.6.0-rc6 (2026-10-15)**. So this task is also gated
on FE2a, independent of the access problem above.
