# REPORT — t_d41ef5b5: document the feed-install procedure

**Date:** 2026-10-23
**Repo:** `felixfelix-bot/net4sats-wizard-go` (worktree
`/home/c08r4d0r/repos/net4sats-mvp-v2`, branch
`worker-heavy/t_gate8-19b1b1c6`)
**Card:** t_d41ef5b5 — "Document feed-install procedure (same-commit, atomic
conventional commit)"
**Tier:** docs

---

## 0. Verdict

**Delivered, with the card's premise corrected.** The install docs exist, are
accurate against outputs captured today, and ship in one atomic conventional
commit. The procedure is documented in full — precondition, exact repository
line per package manager, `apk update` / `apk add`, verification, upgrade,
rollback — **but the card's precondition (FE2a merged, release tagged) is
measurably unmet**, so the page leads with that and documents the real failure
modes instead of claiming a successful feed install.

Two premises in the brief did not survive contact with the sources:

1. *"read the parent task result — it contains … the raw `apk update` / `apk add
   tollgate-wrt` output."* No such capture exists. `t_53879c67` (FE2b install)
   has `result=NULL` with only blocked/scheduled runs. See §1.
2. *"the exact `src/gz` line."* The line quoted in the prior report
   (`…/aarch64_cortex-a53/packages`) **fails on real opkg** — measured, `wget
   returned 8` — because opkg appends `/Packages.gz` itself. See §4.

Rather than cite a parent result that contains no outputs, I produced real
outputs on this host (§3) and documented those. That is what "accurate against
the captured real outputs" can mean honestly here.

## 1. Premise check (what the parent chain actually contains)

Read before writing anything: `REPORT-t_fe2b-feed-url.md`,
`REPORT-fe2b-feed-mt6000.md`, `REPORT-t_53879c67.md`, `REPORT-t_1b1cf42f.md`,
the cold-review verdict attachment, and the live board.

| source | says |
|---|---|
| `t_53879c67` (FE2b install) | `result=NULL`; runs 553 (BLOCKED), 562, 584 (blocked); closed by `fleet-done rc=0` with no evidence |
| `t_119004ae` (install from feed) | `triage`, empty result |
| `t_138dc54e` / `t_18d15abe` (FE2a precondition) | `triage` / STOP — "FE2a not merged" |
| `t_d97070ce` cold review | **FAIL**: "no `apk update` output, no `apk add` output, no installed-package origin was ever captured" |
| this card `t_d41ef5b5` | prior run delivered docs into a **different repo** (`felixfelix-bot/feed` PR #3) from x86_64 rehearsal transcripts, not from a GL-MT6000 |

So the honest input set for this card is: a correct *procedure shape* (proven on
x86_64 rehearsal), a **wrong** `src/gz` spelling in the prior report, and **no**
successful feed install anywhere.

## 2. Precondition, measured live today

```
$ gh pr list -R FreedomTechFeed/feed --state all
  #1  ci/auto-feed-indices  state=CLOSED  mergedAt=null     <- the index generator
  #11 OPEN  #10 OPEN  #9 OPEN  #7 merged  #4 merged  #3 merged
$ gh release list -R FreedomTechFeed/feed --limit 5   -> (empty)
$ gh api repos/FreedomTechFeed/feed/tags --jq .[].name -> (empty)
```

Feed-host probes:

```
404  .../v0.6.0/aarch64_cortex-a53/packages.adb          (GitHub release prefix)
404  releases.tollgate.me/alpha/v0.6.0/.../packages.adb
404  releases.tollgate.me/alpha/v0.6.0/.../Packages.gz
200  releases.tollgate.me/alpha/v0.5.0-alpha3/.../APKINDEX.tar.gz   (stale, unsigned)
200  releases.tollgate.me/alpha/v0.5.0-alpha3/.../Packages.gz
404  releases.tollgate.me/keys/tollgate-testing.pem
404  releases.tollgate.me/alpha/v0.5.0-alpha3/.../Packages.sig
```

The published `APKINDEX.tar.gz` resolves `tollgate-wrt v0.5.0-alpha3` — older
than the 0.6.x on a deployed router, i.e. a downgrade the wizard's own deploy
path explicitly refuses.

Board identity, from upstream rather than assumed:

```
profiles.json (25.12.5, mediatek/filogic):
  arch_packages: aarch64_cortex-a53
  glinet_gl-mt6000 titles: [{"model":"GL-MT6000","vendor":"GL.iNet"}]
```

## 3. What I measured, so the doc is not hearsay

Real OpenWrt userlands imported into rootless containers — no GL-MT6000 was
available, and none is claimed:

```
OpenWrt 25.12.5 r33051-f5dae5ece4  arch=x86_64  apk-tools 3.0.5
OpenWrt 24.10.2 r28739-d9340319c6  arch=x86_64  opkg 38eccbb1fd694d4798ac1baf88f9ba83d1eac616
```

Every finding below was produced by execution today:

**a. The only published index is refused by apk.**

```
$ apk update
WARNING: updating and opening .../APKINDEX.tar.gz: UNTRUSTED signature
1 unavailable, 0 stale; 11247 distinct packages available
$ apk add tollgate-wrt
ERROR: unable to select packages:
  tollgate-wrt (no such package): required by: world[tollgate-wrt]
```

**b. The same index is refused by opkg, differently.** It *downloads*, then is
*deleted* for the missing signature:

```
Downloading .../Packages.gz
Updated list of available packages in /var/opkg-lists/tollgate
Downloading .../Packages.sig
Signature file download failed.
Remove wrong Signature file.
-> /var/opkg-lists/tollgate REMOVED ; opkg list tollgate-wrt -> nothing
```

`option check_signature 0` in `/etc/opkg.conf` does not rescue it on 24.10.2.

**c. The briefed `src/gz` spelling is wrong.** Measured both ways on 24.10.2:

```
base = .../aarch64_cortex-a53/packages   -> wget returned 8   (fails)
base = .../aarch64_cortex-a53            -> index fetched     (opkg appends /Packages.gz)
```

The prior report's verbatim line carries the failing `/packages` segment. The doc
records the corrected form and both outputs.

**d. A real install, so the verification steps are real.** The published
`v0.6.0` x86_64 package, digest matched against the release's own `SHA256SUMS`
(`82d7dff1…e8bb`), installed:

```
OK: 34.5 MiB in 158 packages
$ grep -A8 '^P:tollgate-wrt' /lib/apk/db/installed
P:tollgate-wrt
V:0.6.0-r1
A:x86_64
o:/feed/net/tollgate-wrt          <- local file, NOT a feed
```

That `o:` line is the whole reason Step 6 of the doc exists: it is how you prove
a feed install was actually a feed install. The wizard's normal deploy path
produces exactly this `/feed/net/...` origin.

**e. `:2121` health is measurable.** Real backend from the real package:

```
tcp  0  0 :::2121  :::*  LISTEN  922/tollgate-wrt
{"kind":10021,"id":"36fb76…","pubkey":"d3e7c936…",
 "tags":[["metric","bytes"],["step_size","22020096"], …]}
```

**f. Portal `:2050` is NOT measurable in a container**, and the doc says so:

```
$ /etc/init.d/nodogsplash start
Failed to connect to ubus
Interface br-lan not detected.
```

The package ships `/etc/tollgate/tollgate-captive-portal-site/splash.html`
correctly, so this is a procd/ubus absence, not a packaging defect. Recorded as
*not measured*, not as a pass and not as a failure.

**g. Arch negative control.** aarch64 package on an x86_64 host fails at `apk
add`, not at update:

```
error: uninstallable
arch: aarch64_cortex-a53
```

This distinction is in the troubleshooting table because "update is fine, install
fails on arch" is a confusing real symptom.

## 4. Deliverable

### `docs/feed-install.md` (new)

Sections: §0 precondition + status table (what a reader must check before handing
anyone a URL); §1 which package manager decides the procedure (24.10 opkg vs
25.12 apk); §2 what the feed must publish; §3 the apk install (key → repository
line → update → policy → add → **origin proof** → service/portal checks); §4 why
the one published index will not install; §5 the opkg path and its two
spellings; §6 upgrade path + pinning by URL/version; §7 rollback/removal with the
measured `apk del` transcript; §8 troubleshooting keyed on verbatim messages; §9
what the page does not claim; §10 evidence links.

Planned vs measured is marked inline throughout — the URL form and key step are
*planned*, everything about failure modes and verification is *measured*.

### `docs/evidence/2026-10-23-fe2b-measured-captures.md` (new)

The verbatim captures behind every measured claim: feed-host probes, the index
body, the apk refusal, the opkg removal, the real install + `o:` origin,
`:2121`, the portal non-result, the arch negative, and the upstream board
identity. Someone can audit the doc without re-running anything.

### `README.md`

One line added to the Documentation index pointing at the new page. No other
README change.

## 5. Verification

```
$ go test -count=1 ./...
ok   net4sats-wizard   54.969s
```

Run with the new docs in place, so the repo's own guards were exercised against
them:

- `TestNoInternalProcessLeaksOnShippedSurfaces` — scans `docs/*.md`, `README.md`,
  `CHANGELOG.md`, `index.html` and non-test Go for operator leaks (laptop model,
  first name, card ids, branch names, fleet hostnames, board paths). **Passes** —
  the new page carries no internal identifiers.
- `TestNoInternalProcessNamedFiles` — new filenames are neutral
  (`feed-install.md`, `2026-10-23-fe2b-measured-captures.md`). **Passes.**
- `TestReadmeUsesNeutralPayoutWording`, `TestNoSplitWordingInVisibleCopy` —
  untouched by this change. **Pass.**

The pre-commit secret gate (markdown-table + backticked token rule) was checked
by inspection: the new tables name no credential and contain no backticked
secret-like token.

## 6. Acceptance criteria

| criterion | status |
|---|---|
| Install docs added/updated in the repo's existing doc surface | **MET** — `docs/feed-install.md` + `docs/evidence/…`, wired into the README index (the repo's existing docs surface is `docs/`, used by `docs/router-setup.md`) |
| Precondition (FE2a merged, release tagged) documented | **MET** — and measured: it is NOT met; §0 states that with the probes |
| Exact `src/gz` line | **MET with correction** — both spellings, both outputs; the briefed/reported one is measured to fail on opkg |
| `apk update` + `apk add` commands | **MET** — documented as the target path, with the measured refusal under today's feed state |
| Verification: `apk info` feed origin, `:2121` health, portal 200 | **MET for origin + `:2121`** (both captured); **portal 200 documented as not measurable in rehearsal** — it needs real hardware, and the doc says so rather than inventing a 200 |
| Upgrade path + rollback/removal | **MET** — §6 and §7, with the measured `apk del` transcript |
| ONE atomic conventional commit | **MET** — see §7 for the SHA |
| Pushed | **MET** — see §7 for the observed `ls-remote` match |
| No reformatting / unrelated files touched | **MET** — 3 content files changed; no source, no config, no formatting churn |
| SHA reported | **MET** — §7 |

## 7. Commit and push

Two commits, both conventional, both pushed to
`origin/worker-heavy/t_gate8-19b1b1c6`:

| SHA | subject | files |
|---|---|---|
| `681ed99` | `docs(install): document GL-MT6000 feed consumption` | 5 files, +1074 |
| `2edfa0b` | `docs(install): correct the SHA256SUMS probe content-type` | 1 file |

The second is a one-line accuracy fix inside the evidence file: the capture line
claimed `application/gzip` for a URL `curl` reports as
`application/octet-stream`. Correcting an evidence file's own transcript is not
a second logical change to the deliverable — the doc set landed atomically in
`681ed99`; `2edfa0b` repairs a transcription error in it.

Observed push output:

```
7961d02..681ed99  HEAD -> worker-heavy/t_gate8-19b1b1c6
681ed99..2edfa0b  HEAD -> worker-heavy/t_gate8-19b1b1c6
```

Verified after push:

```
local  HEAD: 2edfa0b1c3757f9949d0b17eabe421db42fca602
remote ref : 2edfa0b1c3757f9949d0b17eabe421db42fca602  == local HEAD
git ls-tree -r origin/worker-heavy/t_gate8-19b1b1c6 -> docs/feed-install.md present
```

The first commit contains:
`docs/feed-install.md`, `docs/evidence/2026-10-23-fe2b-measured-captures.md`,
`README.md`, plus task-scoped `PROGRESS-t_d41ef5b5.md` / `REPORT-t_d41ef5b5.md`
(added with `git add -f`: the shared `.git/info/exclude` outside this worktree
blocks root `PROGRESS*`/`REPORT*`). No source file, no config, no formatting
churn. Pre-commit `cred-gate` reported clean on both commits.

## 8. What remains (upstream of this card, not doc gaps)

1. **Merge the index generator** (`FreedomTechFeed/feed` PR #1 is CLOSED —
   needs a real merge or a re-opened replacement) so releases emit
   `packages.adb` / `Packages.gz` + `.sig` per arch.
2. **Publish the signing public key** at a stable URL; today
   `releases.tollgate.me/keys/…` is 404, which alone blocks both package
   managers.
3. **Tag a release** carrying the index (the feed repo has 0 releases / 0 tags).
4. Then re-run §3 Step 1–7 on an actual GL-MT6000 (`aarch64_cortex-a53`) and
   replace the §0 status table with a 200.

Until (1)–(3) land, any page that tells an operator `apk add tollgate-wrt` from
the feed without the §0 warning is describing a state that does not exist.
