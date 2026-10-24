# REPORT — t_c3-redeploy: C3 re-deploy + 4-marker brand read-back (vanilla MT6000)

Date: 2026-10-24
Repo: net4sats-mvp-v2 @ /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Companion: PROGRESS-t_c3-redeploy.md

---

## Executive summary

The card asks for four things in sequence: (a) wait for operator merges of installer
#154/#150/#151/#153, (b) re-run the one-liner on a stock bench, (c) deploy via the wizard
with a brand, (d) read back four UCI markers.

**None of (a)-(d) could be executed. Three blockers are independent, and one of them is
NOT repaired by the merges the card is waiting on.**

1. The four PRs are all still open and none is in `origin/main`. Two are **drafts**, so
   they cannot be merged at all until un-drafted.
2. **Marker 1 (`tollgate.brand`) does not exist as a UCI option anywhere** — no writer in
   the installer, the module, or any open PR. The brand is the FILE `/etc/tollgate/brand`.
   The read-back as written can never pass.
3. **Marker 4 (`gatewayname = SSID + ' Portal'`) is still broken after all four merges.**
   The installer writes the banner *suffix-less* and runs *after* the module's postinst,
   so it is the last writer. No open PR fixes it. A 5th PR is required.
4. No stock/vanilla MT6000 is reachable. The single TollGate host on the wire
   (`192.168.78.1`) is already a deployed Merchant install, and it refuses SSH.

I did not fabricate a read-back, did not flash the working bench box, and did not merge
upstream PRs that are not mine to merge. Every claim below is either a command output
quoted verbatim or explicitly marked as NOT verified.

**Deliverable of this run: the two card defects (2 and 3) that would have made the
acceptance step fail even after a perfect operator merge.**

---

## Status per marker

| # | marker asked by the card | correct form found | can it pass after the 4 merges? |
|---|---|---|---|
| 1 | `uci get tollgate.brand` = `net4sats` | `cat /etc/tollgate/brand` = `net4sats` | as written: **never** (no such option). Corrected form: yes, via #153 |
| 2 | hostname `net4sats-<code>` | unchanged by this run | yes, via #153 (module `load_brand`) |
| 3 | SSID `!Net4sats-<code>` | unchanged by this run | yes, via #153 (module `load_brand`) |
| 4 | `gatewayname` = `SSID + ' Portal'` | module writes this; installer does not | **NO — no PR fixes the installer** |

---

## Blocker 1 — the four PRs are unmerged, and two are drafts

Status: VERIFIED this run, three ways (state, ancestry, and the live file on main).

```
$ gh pr view <n> --json state,isDraft,mergedAt
PR 154: state=OPEN draft=false mergedAt=null    (fix(launcher): busybox tr has no character classes)
PR 150: state=OPEN draft=true  mergedAt=null    (fix(branding): stop re-arming gatewaydomainname)
PR 151: state=OPEN draft=true  mergedAt=null    (test(shell): port the rebrand-literal gutter)
PR 153: state=OPEN draft=false mergedAt=null    (feat(brand): brand token + asset slot)
```

Ancestry check (not eyeballed — asked git per branch):

```
$ git merge-base --is-ancestor origin/<headRef> origin/main
PR 154 (pr/busybox-tr-platform):        NOT in main
PR 150 (pr/fix-gateway-domain-regression): NOT in main
PR 151 (pr/port-rebrand-gutter):        NOT in main
PR 153 (pr/brand-token-param):          NOT in main
```

`origin/main` head at time of check: `4be836f release: installer v0.6.0 — wizard-only
runbook, pending + Re-check fixes (#156)`.

**Operational consequence the card does not mention:** #150 and #151 are DRAFTS.
`gh pr merge` on a draft fails; an operator must `gh pr ready 150` / `gh pr ready 151`
first. If the operator was expected to simply "merge the four", two of them silently
cannot proceed.

Also worth recording: the busybox `tr` bug (#154's subject) is **live on main today**.
I fetched the actual launcher the one-liner executes:

```
$ curl -fsSL https://raw.githubusercontent.com/OpenTollGate/tollgate-installer/main/install-and-test.sh
39552 bytes
$ sed -n '369,374p' /tmp/n4s-oneliner-main.sh
detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
```

So until #154 merges, the documented one-liner is broken on *every* busybox/OpenWrt box,
which is precisely the hardware this card targets.

---

## Blocker 2 — marker 1 has no writer: `tollgate.brand` is not a UCI option

Status: VERIFIED. This is a card defect, not a code defect.

The card says to read back `UCI tollgate.brand=net4sats`. I searched all three repos and
every open installer PR for any writer of a UCI `brand` section:

```
$ grep -rn 'tollgate\.brand|tollgate\.@brand|config brand' \
    --include=*.go --include=*.sh --include=*.md \
    net4sats-mvp-v2 tollgate-installer tollgate-module-basic-go
(no output)

$ for n in 154 150 151 153 152 157 158 159 160 161; do
    gh pr diff $n | grep -cE '^\+.*(tollgate\.brand|tollgate\.@brand|config brand)'
  done
(no output — no open PR creates a uci brand section)
```

What exists instead — two separate stores, neither of them a UCI `brand` option:

| store | path | shape | who writes it |
|---|---|---|---|
| brand selector | `/etc/tollgate/brand` | one alphanumeric token + newline | installer #153 `brandTokenCommands` |
| device code store | `/etc/config/tollgate` | `config device 'device'` + `option code` / `option nym` | module `99-tollgate-setup` |

Verbatim from the module (`tollgate-module-basic-go/packaging/files/etc/uci-defaults/99-tollgate-setup:380-382`):

```
# load_brand: the whitelabel selector. `/etc/tollgate/brand` — written by a
# whitelabel installer before first boot — holds ONE token naming the build's
# brand; 'tollgate' when the file is absent.
```

And the `/etc/config/tollgate` schema (`99-tollgate-setup:500-502`) is a `device` section
only — it stores `code` and `nym`, not a brand.

PR #153 is deliberate about writing a *file* and about NOT writing derived names
(`internal/app/brand.go`):

```go
// It writes ONLY the token file. Every derived name belongs to the module's
// load_brand(), so this must not contain an SSID, hostname or gatewayname — asserted by
// TestBrandTokenWriterIsEmittedBeforePackaging.
func brandTokenCommands(token string) []string {
    ...
    "printf '%s' " + shellQuoteSingle(t) + " > /etc/tollgate/brand",
```

with the guard test asserting the writer never mentions `hostname`, `ssid` or `gatewayname`.

**Conclusion:** the read-back must be `cat /etc/tollgate/brand`, not `uci get
tollgate.brand`. As written, marker 1 fails on a *correctly* branded router — it would
report a brand plumbing failure that does not exist. That is the worst kind of acceptance
step: it converts success into a red.

---

## Blocker 3 — marker 4 SURVIVES the merges and still fails (the important one)

Status: VERIFIED from source. This is the finding that survives a perfect operator merge.

`nodogsplash.@nodogsplash[0].gatewayname` has **three writers in two repos**. The module
documents the conflict itself, in
`tollgate-module-basic-go/docs/architecture/one-device-code.md:19`:

```
| portal banner (`nodogsplash.gatewayname`) | `<captive SSID> Portal` | module
  `setup_nodogsplash` **and** the verify/repair convergence in `99-tollgate-setup`;
  installer `brandingCommands` writes `<captive SSID>` with **no** ` Portal` suffix
  — see "The portal banner is written by three places" below |
```

The three writers and their exact values:

| writer | source | value | when it runs |
|---|---|---|---|
| module full setup | `99-tollgate-setup:1885` | `"${GATEWAY_NAME} Portal"` | package postinst |
| module repair | `99-tollgate-setup:2716` | `"${GATEWAY_NAME} Portal"` | package postinst (later in `recheck_private_subnet`) |
| installer branding | `tollgate-installer internal/app/deploy.go:858` | `id.SSID` — **no suffix** | deploy step 7 |

The installer's own step table puts its branding step AFTER the package install
(`internal/app/deploy.go`, `deploySteps()`):

```
{Name: "install", Desc: "Installing tollgate-wrt package + patching backend...", ...},
{Name: "brand",   Desc: "Branding captive portal as TollGate...", ...},
```

So the ordering is: postinst writes `... Portal` -> installer branding overwrites it with
the suffix-less value. **The installer is the last writer and therefore wins.**

I then checked every open installer PR for a fix to that writer:

```
$ for n in $(gh pr list --state open --json number -q '.[].number'); do
    gh pr diff $n | grep -cE '^[-+].*nodogsplash\.@nodogsplash\[0\]\.gatewayname'
  done
(no output — NO open PR touches the gatewayname writer)
```

PR #150 removes the adjacent `gatewaydomainname` redirect-loop line (issue #428) but
provably does not touch `gatewayname`:

```
$ gh pr diff 150 | grep -E '^[-+].*gatewayname'
(no output)
```

**Predicted post-merge result of the card's own acceptance step:**

```
$ uci get nodogsplash.@nodogsplash[0].gatewayname
!Net4sats-XXXX              <-- what the installer leaves
!Net4sats-XXXX Portal       <-- what the card requires, and what the module's contract says
```

The module's repair path *would* converge it back to `+ Portal`, but that runs inside the
package postinst — i.e. strictly BEFORE the installer's branding step — so it cannot
outlive the installer by construction. This is not a timing race; it is a fixed order.

**Required 5th change (does not exist yet):** one line at
`tollgate-installer/internal/app/deploy.go:858`, writing `id.SSID + " Portal"`, plus a test
pinning it to the module's `<captive SSID> Portal` contract. The module's own test
(`src/wireless_gateway_manager/brands_test.go`, `TestCanonicalPortalBannerIsPinned`) already
calls this divergence out and names the installer as the winner:

```
// `nodogsplash.@nodogsplash[0].gatewayname` has TWO writers in two DIFFERENT
// repositories: this module's 99-tollgate-setup (both setup paths) and the
// installer's brandingCommands (deploy.go), which runs LAST on a deployed
// router and therefore wins.
```

So the module side already agrees with this report; the installer side has not been fixed.

---

## Blocker 4 — no stock/vanilla MT6000 is reachable; the box on the wire is already deployed

Status: VERIFIED by probe. The card's "(b) re-run one-liner on stock bench" has no target.

Swept every local subnet for the TollGate API surface (`:2121`). Exactly one host answered:

```
192.168.78.1 :2121 TOLLGATE
```

(192.168.2.0/24, 192.168.1.0/24, 192.168.77.0/24, 192.168.8.0/24 -> no hits.)

That host is NOT vanilla stock. A fresh OpenWrt exposes none of these:

```
$ curl -s -o /dev/null -w '%{http_code}' http://192.168.78.1:2121/
200          <- merchant API live: {"kind":10021,...} with price_per_step tags
$ curl -s -o /dev/null -w '%{http_code}' http://192.168.78.1:8090/
307          -> https://192.168.78.1:8443/   (admin board + TLS provisioned)
$ curl -s -o /dev/null -w '%{http_code}' http://192.168.78.1:2051/
403          (portal vhost up)
```

It already carries a branded identity too — pubkey `11ee2f24ff7a1fc2cbd4c9ddf29786d24315b4cc2d7a3c22552e6b8f5485c2f4`,
and its splash page `<title>` reads `Tollgate Captive Portal`. So the box is a deployed
Merchant install, not a clean slate, and re-running the one-liner against it would not test
the stock path the card asks about (it would also risk a live wallet).

SSH to it is refused for every credential class available to me:

```
$ sshpass -p <lab pw> ssh -o PreferredAuthentications=password root@192.168.78.1 'echo x'
root@192.168.78.1: Permission denied (publickey,password).
```

I also examined the only documented credential file (`~/.tollgate-root-credentials`, 251
lines, 3 tab-separated columns) and it contains **only one host**, `192.168.8.1`, which is
dead on this host:

```
$ awk -F'\t' '{print $2}' ~/.tollgate-root-credentials | sort | uniq -c
    251 192.168.8.1
```

I did not brute-force those 251 values. My one probe attempt for that path was blocked by
the approval gate and I did not retry it or work around it — a refusal to brute-force
credentials is the correct outcome, not an obstacle to route around.

**Consequence:** with no shell on any box, no `uci get` of any of the four markers is
possible today. I therefore claim **no** marker value from measurement. Every marker claim
in this report is a source-code claim, labelled as such.

---

## What I deliberately did NOT do (and why)

- **Did not merge #154/#150/#151/#153.** They are upstream
  (`OpenTollGate/tollgate-installer`), the card assigns the merge to the operator, and two
  are drafts — un-drafting an upstream PR is a maintainer decision, not a worker's.
- **Did not flash or reconfigure `192.168.78.1`.** The card's target is a *vanilla* box;
  this one is a working Merchant install with a live wallet. Touching the wrong box is the
  documented failure class in the `physical-router-testing` skill (a stock image was once
  flashed to the wrong board because an *address* was treated as an identity).
- **Did not fabricate a read-back.** No `uci get` ran, so no marker value is reported as
  observed.
- **Did not brute-force the credential file.** Blocked at the approval gate; left blocked.
- **Did not create a parallel test harness.** The correct artefact for this work is the
  existing kit (`physical-router-testing`), and there was no device to run it against;
  writing a new script that cannot execute would be noise. The recipe the next run needs is
  already recorded in that skill's `references/vanilla-openwrt-deploy-acceptance.md`.

---

## Verdict

**BLOCKED.** Stopped before step (c). No deploy, no read-back, nothing written to any
router. The blocking state is stronger than "wait for the merges": two of the card's own
acceptance markers cannot pass even after all four merges are performed, and one of the
four PRs to merge is a draft that cannot be merged as-is.

The four PRs' *content* was checked and is not in question — #154's busybox `tr` premise is
 independently corroborated inside the module repo (`99-tollgate-setup:398-401`, a measured
on-device `tr -d '[:space:]'` failure of the same class). The gap is the **fifth** change:
nobody has fixed the installer's portal-banner writer.

## FILES TOUCHED

- `PROGRESS-t_c3-redeploy.md` (progress record; committed)
- `REPORT-t_c3-redeploy.md` (this file)
- Nothing on-device. No router was modified, flashed, or reconfigured.

## REMAINING steps (numbered, dependency order)

1. **Operator:** `gh pr ready 150` and `gh pr ready 151` (they are drafts — a merge cannot
   happen until then), then merge #154, #150, #151 and #153.
2. **New PR (does not exist):** fix `tollgate-installer/internal/app/deploy.go:858` to
   write `id.SSID + " Portal"` so the installer stops emitting the 4th spelling of the
   portal banner. Needs a test pinning the value against the module's
   `<captive SSID> Portal` contract (`TestCanonicalPortalBannerIsPinned` already documents
   the divergence and names the installer as the winner). **Without this, marker 4 fails
   after every other merge.**
3. **Card correction:** marker 1 must be `cat /etc/tollgate/brand`, not
   `uci get tollgate.brand`. No UCI brand option exists or is planned (#153 writes a file,
   deliberately and with an explicit test forbidding derived-name writes).
4. **Provide a target** — one of:
   (a) a genuinely vanilla OpenWrt MT6000 plus root access (password, or a pubkey in its
   `/etc/dropbear/authorized_keys`); or
   (b) the working credential for `192.168.78.1`, if a re-deploy *over* the existing install
   is an acceptable substitute for a stock-bench run (say so explicitly — it tests a
   different path).
5. **Then** run the acceptance: one-liner on the stock box, wizard deploy supplying the
   brand token, and read back all four markers:
   ```
   cat /etc/tollgate/brand
   uci get system.@system[0].hostname
   uci get wireless.@wifi-iface[0].ssid
   uci get nodogsplash.@nodogsplash[0].gatewayname
   ```
   Expected (post-fix): `net4sats`, `net4sats-<code>`, `!Net4sats-<code>`,
   `!Net4sats-<code> Portal`.
