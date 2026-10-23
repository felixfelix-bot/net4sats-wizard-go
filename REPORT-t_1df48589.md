# REPORT — t_1df48589 · wz1.5 channel decision note + Q5 status

**Worker:** `t_gate8-19b1b1c6`, profile `worker-tollgate`
**Repo/worktree:** `/home/c08r4d0r/repos/net4sats-mvp-v2`, branch `worker-heavy/t_gate8-19b1b1c6`
**Deliverable:** `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md`
**Status:** DONE — note corrected and delivered; **Q5 = OPEN**
**Repo source changed:** none (`deploy.go`, `pins_test.go` untouched — nothing built or tested)

---

## 0. Headline

A complete decision note **already existed** at the deliverable path when this pass started — written
by an earlier pass of this same card (`t_1df48589`, 2026-10-23 23:38, 414 lines, 22.7 KB). This
pass therefore did not rewrite it. It did three things:

1. **Re-derived every live claim independently** (2026-10-24) instead of trusting the prior pass.
   All load-bearing claims hold.
2. **Found and fixed one real defect** in the note: its instantiated artifact-filename template is a
   **404** on the one version that exists on the chosen host. See §3.
3. **Closed the durability gap**: the note, its captures, and this card's report files were
   **never committed anywhere** — the plans directory is not a git repo and this worktree gitignores
   `/PROGRESS-t_*.md` + `/REPORT-t_*.md`. Fixed by committing a byte-identical mirror plus
   verbatim captures. See §2.

The decision itself is unchanged: interim channel `releases.tollgate.me`, **Q5 = OPEN**.

---

## 1. Inputs

| Input named by the brief | Status |
|---|---|
| `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-phase2-feed-status.md` | **Present** (18,604 B, 2026-10-23 22:58). Verdict: **NOT PUBLISHED** |
| `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-probe-raw.md` | **ABSENT — never written.** Confirmed by `read_file`, `ls -la` of the plans dir, and a filesystem `find`. Its intended content (raw probe evidence + example versioned URLs) was re-derived first-hand instead |
| `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` | **Present** (17,507 B). It is the FELIX Q3 / dev-split gate file (Q3 = B, resolved 2026-10-20). It contains **nothing** about artifact channels — no input drawn from it |
| `deploy.go:24-52` | **Read.** At HEAD `7998fc3`: `:43` `tollgatePkgURL`, `:47` `tollgatePkgURLApk`, `:50` `configwizURL` — all three `github.com/felixfelix-bot/...`. `grep -rn -i 'blossom2\|orangesync' --include='*.go' .` → **0 hits** |

The brief's framing "replace unofficial mirrors" (plural, implying `blossom2.orangesync.tech`) is
half-wrong for `deploy.go`: there is exactly **one** unofficial *owner* (`felixfelix-bot`, the
project's own fork account), and one of the three constants (`configwizURL`, `:50`) is a
**different project** (`configurationwizzard`) and out of scope for a tollgate-wrt channel
decision. `blossom2` survives only in history (`11f63e7:deploy.go:42`).

The decision rule was applied as written: Phase-2 feed **NOT PUBLISHED** → documented interim
channel = `releases.tollgate.me`, **Q5 flagged OPEN**.

---

## 2. Durability gap found and closed

The prior pass's own `PROGRESS-t_1df48589.md` claimed its files were "tracked, committed". **They
were not.** Measured:

```
$ git log --all --oneline -- PROGRESS-t_1df48589.md REPORT-t_1df48589.md
(no output — 0 commits)

$ git log --all --oneline -S'wz1.5 -- official Tollgate artifact channel'
(no output — the note's content is in no commit)

$ cd ~/.hermes/plans/2026-08-31-tollgate-batch && git rev-parse --show-toplevel
fatal: not a git repository (or any of the parent directories): .git

$ git check-ignore -v PROGRESS-t_1df48589.md REPORT-t_1df48589.md
/home/c08r4d0r/wz16-verify/.git/info/exclude:4:/PROGRESS-t_*.md	PROGRESS-t_1df48589.md
/home/c08r4d0r/wz16-verify/.git/info/exclude:5:/REPORT-t_*.md	REPORT-t_1df48589.md
```

Earlier cards landed despite this via `git add -f` (`PROGRESS-t_855cbaac.md`,
`PROGRESS-t_gate8-19b1b1c6.md` are tracked); this card's files simply never were. So the wz1.5
decision existed only as an unversioned file in a home directory — not reproducible, not
reviewable, not diffable.

Fix applied (no change to the canonical deliverable path the card names):

- `docs/wz1.5-channel-decision.md` — **byte-identical mirror** of the deliverable, committed here.
  `sha256sum` of both paths: `409de19066fe06ba28225a627b77aac9ab40bff7eb98a4027dfe5855f66da31b`
  (501 lines after this pass's edits).
- `docs/evidence/2026-10-24-wz15-channel-decision-reverify.md` — verbatim captures for this pass,
  following the sibling card's `docs/evidence/<date>-<name>.md` convention
  (`docs/evidence/2026-10-23-fe2b-measured-captures.md`).

---

## 3. DEFECT FOUND AND FIXED — the artifact filename is not derivable

The note's §1 instantiated its pattern as:

```
https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<arch>/tollgate-wrt_v<X.Y.Z>[-<pre>]_<arch>.ipk
```

**That template 404s on the only version that exists on the chosen host.** Measured:

```
$ B=https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53
$ curl -sSL "$B/Packages" | grep -E '^(Package|Version|Filename|Size|SHA256sum):'
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
Size: 5321054
SHA256sum: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk
Size: 5307464
SHA256sum: c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk
Size: 7153725
SHA256sum: 5332bb5ea15edbef3a41d5a34cf0e5e66a8d0d9a152385d4e0a7f06e928e6f3d

$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.ipk"                    -> 404
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk"   -> 200 (5,321,054 B)
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk"                   -> 200 (7,153,725 B)
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk"   -> 200 (5,307,464 B)
```

Three facts fall out, all of which the note now records:

1. The `.ipk` carries a **build-flavour suffix** (`-upx-ultra-brute`) and there is **no plain
   `.ipk` at all** — the templated name does not exist for that version.
2. The suffix is **not uniform**: the **apk** exists in *both* plain (7,153,725 B) and suffixed
   (5,307,464 B) form — two different builds of the same version.
3. Therefore **only the index knows the filename**; it must be read from `Filename:` (opkg) /
   `F:` (apk index), never composed from version + arch.

Every index entry's `content-length` matches its own `Size:`, so the index is a faithful manifest
and the three entries are three real, distinct artifacts. One full download confirmed the bytes:

```
$ curl -sSL -o chk.apk "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk"
http=200 bytes=5307464      on-disk=5307464
$ sha256sum chk.apk
c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f   # == index SHA256sum:
```

**Why this matters:** the note is the implementer's only instruction sheet. Following its §1
template literally would have produced a 404 URL — the exact regression class recorded at
`deploy.go:37-42` (SW4a: "asset never existed on this release — HTTP 404, broke every fresh
deploy"). A note that mandates a version-pinning rule while itself handing over an unresolvable
filename is a live hazard, so this was fixed rather than noted.

### Edits made to the note

| § | Change |
|---|---|
| header | added a **Re-verified 2026-10-24** line naming the defect and its fix |
| §1 | replaced the two composed filename templates with the bare `<artifact-file>` placeholder + an explicit "this is not a derivable template" paragraph carrying the measured counter-example |
| §2.3a | **new** — full evidence block: index fields, the three `content-length` cross-checks, the 404 on the composed name, the apk asymmetry, and the full-download digest |
| §3 | **new rule 1b** — "`<artifact-file>` is copied from the version's own index, never composed", with the failure mode stated |
| §6.1 | caveat that the target literals are **illustrative of shape only**: the filename segment must come from the index and the arch segment from the router's `DISTRIB_ARCH` |
| §7 | two new acceptance rows: the filename-pattern correction, and the second-pass re-verification |

---

## 4. Independent re-verification of every live claim (2026-10-24)

Nothing below is inherited from the prior pass or from a sibling report; each line was produced by a
command in this session. Captures: `docs/evidence/2026-10-24-wz15-channel-decision-reverify.md`.

| Claim in the note | Re-measured | Verdict |
|---|---|---|
| `releases.tollgate.me/` answers | `206`, 697 B | holds |
| Only `v0.5.0-alpha3` has an index; all other tags 404 | `…/v0.5.0-alpha3/…/Packages` → `206` 1,322 B; `…/v0.6.0/…/Packages` → `404` | holds |
| Exactly one version has artifacts on the chosen host | 27-URL sweep, 9 tags × 3 filename shapes → **zero** non-404 hits | holds |
| The channel is unsigned | `Packages.sig` → `404`; `keys/tollgate-testing.pem` → `404` | holds |
| Both target literals derived from the pattern are 404 | `…/v0.7.0-alpha10/…ipk` → `404`; `…/v0.6.1-post-merge/…apk` → `404` | holds |
| Rolling `latest` is not addressable | `…/releases/download/latest/SHA256SUMS` → `404` (9 B); `/releases/latest` → `302` → `…/tag/v0.6.0` (HTML) | holds |
| Two real example URLs, two versions, digest-verified | `v0.6.0` ipk and `v0.6.0-rc6` ipk both `206` | holds |
| `FreedomTechFeed/packages` inventory | `v0.6.0` 2026-10-21T02:06:45Z assets=26; `v0.6.0-rc6` 2026-10-15T19:50:37Z assets=26; no drafts | holds |
| Upstream `OpenTollGate/tollgate-module-basic-go` publishes no binaries | `v0.6.0` 2026-10-21T03:27:56Z **assets=0**; also 0 for `v0.6.0-alpha1`, `v0.5.0-alpha3`, `v0.5.0-alpha2` (only `v0.5.0` has 16) | holds |
| Neither wizard-pinned tag exists on the `FreedomTechFeed` lane | `v0.7.0-alpha10` api=404; `v0.6.1-post-merge` api=404 | holds |
| Channel vocabulary is the publisher's, not invented | explorer bundle `main.ca30081d.js` 474,543 B → `STABLE:"stable",BETA:"beta",ALPHA:"alpha",DEV:"dev"`, products `tollgate-os`/`tollgate-wrt`/`tollgate-module-basic-go` | holds |

One transient: a probe batch failed with `curl: (6) Could not resolve host: api.github.com` while
`github.com` and `releases.tollgate.me` resolved fine in the same seconds. Retried with backoff;
results unchanged. Recorded so a future lane does not misread a DNS blip as a dead endpoint.

---

## 5. The decision, restated (unchanged)

- **Chosen channel (interim, documented): `releases.tollgate.me`.**
- **Base URL pattern:** `https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<arch>/<artifact-file>`
  where `<arch>` is the OpenWrt `DISTRIB_ARCH` the wizard detects (e.g. `aarch64_cortex-a53` for
  the GL-MT6000 on 25.12.5) and `<artifact-file>` is copied from that version's index `Filename:`.
- **Why:** it is the only host the project names that is a Tollgate-owned **release channel** with a
  channel segment in the path (`/alpha/`), and that segment is the publisher's own vocabulary.
- **Not a feed line:** no opkg/apk index exists for any current release, so the pattern resolves a
  single artifact file — usable for the wizard's download-then-install path, not as a
  `src/gz`/`packages.adb` repositories line.
- **Version-pinning rule:** immutable version-pinned paths only; never hash-pin a moving blob;
  never reference a rolling `latest` in code; URL and digest move together.
- **Two real example versioned URLs** (from the only lane with two published versions, with HTTP +
  digest evidence):
  `https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk`
  and
  `https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk`
- **Q5 = OPEN.** Interim vs final channel undecided: the decision rule resolves toward the final
  channel only if the Phase-2 feed is PUBLISHED, and it is NOT PUBLISHED.
- **Follow-on for `deploy.go:24-52`:** rewrite `:43` and `:47` onto the chosen base, taking the
  filename from the index and the arch from the router — but **do not land them yet**: both target
  URLs are 404 today, and a four-point precondition gate (200 + index-matching digest + accepted
  signature + no version downgrade) must pass first. `:50` (`configwizURL`) is out of scope.

---

## 6. Acceptance criteria

| Criterion | Where | Status |
|---|---|---|
| Exact chosen base URL + pattern | note §1 | **MET** |
| Why chosen | note §1, §5 | **MET** |
| Evidence: commands + observed output, verbatim | note §2.1–§2.6 + `docs/evidence/2026-10-24-wz15-channel-decision-reverify.md` | **MET** — first-hand this pass |
| Version-pinning rule (immutable paths, no moving-blob hash pin, no rolling `latest` in code) | note §3 rules 1–4 | **MET** |
| ≥2 real example versioned URLs, two versions, HTTP evidence | note §2.4 | **MET on the release lane** (both 206, digests matched); **explicitly NOT MET on the chosen channel** — only one version exists there, stated rather than fabricated |
| Q5 status + one-sentence justification | note §5 | **MET** — **OPEN** |
| Follow-on: actionable `deploy.go:24-52` instruction | note §6 | **MET** — target literals + four-point gate |
| No unsupported claims | throughout | **MET** — every claim carries a command or a `file:line`; the two unsatisfiable ones (second channel version; the missing raw-probe input) are declared |

---

## 7. Verification of the repo state (no source change)

```
$ git status --short
 M REPORT.md                  <- sibling card's edits, not mine
?? .hermes/  ?? verify/…      <- untracked scratch from other cards
$ git diff --stat -- deploy.go pins_test.go
(empty — no Go source modified)
```

No build or test run applies: no Go file changed, so no build artifact's behaviour can differ. The
only code-adjacent work is URL probing. `HEAD` == `origin/worker-heavy/t_gate8-19b1b1c6` ==
`7998fc3` before this pass's commit.

---

## 8. Open items handed off

1. **Q5 (channel finality)** — OPEN by rule. Needs an operator/architecture decision, not a
   measurement.
2. **Producing-side gap** — the interim channel publishes **no signed index and no artifact for any
   0.6.x/0.7.x release**. Until it does, the landing gate stays closed and `deploy.go` keeps pointing
   at `felixfelix-bot`. This is the single blocker for wz1.5 Phase 2.
3. **Missing sibling artifact** — `wz1.5-channel-probe-raw.md` was never written. Any card specified
   as consuming it has an unsatisfiable dependency as written; either write it or drop it from those
   cards' input lists.
4. **`REPORT.md` is shared and dirty** — a sibling card holds 32 insertions / 22 deletions there. This
   card's report therefore lives in `REPORT-t_1df48589.md`, and the shared `PROGRESS.md` received an
   append only.
