# PROGRESS — t_c3-redeploy: C3 re-deploy + 4-marker brand read-back on vanilla MT6000

Finding -> Status -> Files touched. Task-scoped file (root PROGRESS.md is shared; not touched).

## Status: BLOCKED (3 independent blockers, one of them survives the required merges)

The card's plan was: wait for installer #154/#150/#151/#153 -> re-run the one-liner on a
stock bench -> wizard deploy with brand -> read back 4 UCI markers.

Three of those four preconditions are false, and one of them is false in a way the
merge does NOT repair. Detail and evidence in `REPORT-t_c3-redeploy.md`.

---

## 1. The four PRs are unmerged — and NOT in `origin/main`

Finding: all four are open; two are DRAFTS. None is an ancestor of main.
Status: VERIFIED this run.

```
PR 154: state=OPEN draft=false mergedAt=null   NOT in main   (busybox tr)
PR 150: state=OPEN draft=true  mergedAt=null   NOT in main   (gatewaydomainname)
PR 151: state=OPEN draft=true  mergedAt=null   NOT in main   (rebrand gutter)
PR 153: state=OPEN draft=false mergedAt=null   NOT in main   (brand token)
```

`origin/main` head is `4be836f` (release v0.6.0, #156). Checked with
`git merge-base --is-ancestor origin/<branch> origin/main` per PR, not by eyeball.

Note: #150 and #151 are **drafts**, so "operator merges them" needs an explicit
un-draft first; a merge cannot happen on a draft PR at all.

Files touched: none.

## 2. Marker 1 (`tollgate.brand=net4sats`) — NO writer exists, and #153 does not add one

Finding: the card asks to read back `uci get tollgate.brand`. There is **no such UCI
option anywhere** — not in the module, not in the installer, not in any open PR.

What actually exists:

| what | where | shape |
|---|---|---|
| brand selector | `/etc/tollgate/brand` (a FILE) | one alphanumeric token, newline-terminated |
| device code store | `/etc/config/tollgate` | `config device 'device'` + `option code`/`option nym` |

`grep -rn 'tollgate\.brand|tollgate\.@brand|config brand'` across net4sats-mvp-v2,
tollgate-installer and tollgate-module-basic-go returns **zero hits**. `/etc/config/tollgate`
holds a `device` section only (`99-tollgate-setup:500-502`);

PR #153 writes the brand as a **file**, deliberately, and with a comment saying so:

```go
// It writes ONLY the token file. Every derived name belongs to the module's
// load_brand(), so this must not contain an SSID, hostname or gatewayname
"printf '%s' " + shellQuoteSingle(t) + " > /etc/tollgate/brand",
...
// TestBrandTokenWriterIsEmittedBeforePackaging
for _, forbidden := range []string{"hostname", "ssid", "gatewayname"} {
```

Status: **the marker as written can never pass.** The correct read-back is
`cat /etc/tollgate/brand` -> `net4sats`. This is a card defect, not a code defect.

Files touched: none.

## 3. Marker 4 (`gatewayname = SSID + ' Portal'`) — SURVIVES THE MERGES AND STILL FAILS

Finding: this is the load-bearing one and it is NOT fixed by #154/#150/#151/#153.

`gatewayname` has **three writers**; the module documents the conflict itself in
`tollgate-module-basic-go/docs/architecture/one-device-code.md:19`:

| writer | value | when |
|---|---|---|
| module `setup_nodogsplash` (`99-tollgate-setup:1885`) | `"${GATEWAY_NAME} Portal"` | package postinst |
| module repair `recheck_private_subnet` (`:2716`) | `"${GATEWAY_NAME} Portal"` | package postinst, later |
| installer `brandingCommands` (`deploy.go:858`) | the SSID, **no suffix** | deploy step 7 |

`deploySteps()` in the installer is explicit that the installer's branding runs AFTER the
package install:

```
install  -> Installing tollgate-wrt package + patching backend...
brand    -> Branding captive portal as TollGate...     <-- deploy.go:858 wins
```

So the last writer is the installer's suffix-less value. I checked every open installer PR:
**not one** touches the `gatewayname` writer (#150 only deletes the `gatewaydomainname`
redirect-loop line; verified by diffing for `gatewayname` across all open PRs). So after all
four merges, a successful deploy still leaves:

```
uci get nodogsplash.@nodogsplash[0].gatewayname   ->  !Net4sats-XXXX      (marker FAILS)
expected by the card                             ->  !Net4sats-XXXX Portal
```

The module's repair path converges it back to `+ Portal`, but that runs in the package
postinst — i.e. BEFORE the installer's branding step — so it cannot outlive the installer.

Status: **card marker is correct per the module's contract; the installer is the defect.**
A 5th PR (one line at `deploy.go:858`, `id.SSID + " Portal"`) is required; it does not exist.

Files touched: none.

## 4. No stock/vanilla MT6000 is reachable; the box on the wire is already TollGate-deployed

Finding: the card says "re-run one-liner on stock bench". There is no stock bench box.

The only TollGate-answering host on every local subnet is `192.168.78.1`:

```
:2121 -> 200  {"kind":10021,...}   merchant API live
:2051 -> 403                        portal vhost up
:8090 -> 307 -> https://...:8443/   admin board + TLS
```

A vanilla OpenWrt has none of those ports. So `/etc/openwrt_release` is not stock here and
the "re-run the one-liner on stock" step has no target. Prior lanes reached the same box and
identified it as the GL-MT6000 (`PROGRESS-t_b46dea1c.md`, `hermes-orchestration/REPORT.md:61`).

Also measured: the box refuses SSH with every credential class I am permitted to use
(`Permission denied (publickey,password)`), and its documented credential file
(`~/.tollgate-root-credentials`, 251 entries) contains **only** host `192.168.8.1`, which is
dead on this host. I did not brute-force it — the probe was blocked by the approval gate and
I did not retry it or route around it.

Status: **BLOCKED** — no shell, so no `uci get` read-back of any marker is possible today.

Files touched: none on-device. Nothing written to any router.

## 5. PR #154's premise is real (checked, because it gates everything)

Finding: I could NOT reproduce the busybox `tr` bug with this host's busybox
(`/usr/bin/busybox` v1.36.1 maps `Linux` -> `linux` correctly), so the PR's claim needed a
second source. It has one, independent, in the module repo:

```
99-tollgate-setup:398-401
# Do NOT use `tr -d '[:space:]'` here: BusyBox tr treats the character class
# as a literal set and deletes every s/p/a/c/e from the value — a nine-character
# brand collapsed to a four-character fragment and coerced to the default.
```

That is a *measured on-device* failure of the same class, recorded by a different repo. So
#154 is a genuine upstream-stock blocker, and the one-liner on `main` today still carries the
bug (`install-and-test.sh:371` fetched live from raw.githubusercontent.com this run).

Status: **VERIFIED as a real blocker** (with the caveat that my host's busybox is not the
affected build — I did not claim a local repro I do not have).

## 6. What I deliberately did NOT do

- Did not merge the four PRs. They are upstream (`OpenTollGate/tollgate-installer`) and the
  card says the operator merges them; two are drafts, which is a decision I cannot make.
- Did not flash or reconfigure the MT6000. It is a working Merchant install with a live
  wallet; the card's target is a *vanilla* box, so touching this one would be the wrong box.
- Did not fabricate a read-back. No `uci get` was run, so no marker value is claimed.

---

## FILES TOUCHED

- `PROGRESS-t_c3-redeploy.md` (this file)
- `REPORT-t_c3-redeploy.md` (full report + verbatim evidence)

## REMAINING steps (numbered, in dependency order)

1. **Operator:** un-draft #150 and #151, then merge all four (#154, #150, #151, #153).
2. **New PR needed (does not exist):** fix `internal/app/deploy.go:858` to write
   `id.SSID + " Portal"` so the installer stops emitting the 4th spelling of the portal
   banner. Without it marker 4 fails after every other merge. Should ship with a test that
   pins the value against the module's `<captive SSID> Portal` contract.
3. **Card correction:** marker 1 must be `cat /etc/tollgate/brand` (a file), not
   `uci get tollgate.brand`. No UCI brand option exists or is planned.
4. **Provide a target:** either a genuinely vanilla OpenWrt MT6000 plus root access (password
   or a pubkey in `/etc/dropbear/authorized_keys`), or the working credential for
   `192.168.78.1` if a re-deploy over the existing install is acceptable instead of a stock one.
5. **Then** run the one-liner on stock, deploy via wizard with the brand token, and read back:
   `cat /etc/tollgate/brand`, `uci get system.@system[0].hostname`,
   `uci get wireless.@wifi-iface[0].ssid`, `uci get nodogsplash.@nodogsplash[0].gatewayname`.
