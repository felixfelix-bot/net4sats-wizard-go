# REPORT — t_deabf3d6: Validate URL change (grep gate + install-step dry-run)

Card: `t_deabf3d6` — "Validate URL change: grep gate plus install-step dry-run"
Worker: fleet offload, profile `worker-tollgate`
Date: 2026-10-23
Repo: `/home/c08r4d0r/repos/net4sats-mvp-v2` (origin = `github.com/felixfelix-bot/net4sats-wizard-go`)
Branch: `worker-heavy/t_gate8-19b1b1c6` @ `681ed99`; remote `main` @ `65b9d02740211a88fa86bab05862f23a9326af8b`
Source modified: **NO**

Deliverable: `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-validation.md` (19,430 B, written
this session) — contains the raw commands and raw output behind every line item here.

---

## VERDICT: BLOCKED — there is no URL change to validate.

The card asks me to validate a pushed change. **No such push exists.** Both inputs the card names
are absent (decision note ENOENT; no implementation SHA anywhere on the remote). I did not invent
a subject to validate: steps 1–4 are reported UNVERIFIED where the precondition is missing, and I
took the surrounding measurements anyway so the next attempt has real numbers.

Two findings matter more than the missing change:

- **F1 — The planned official-channel re-pin cannot work as written.** The only ref carrying an
  official-channel attempt (`origin/feat/feed-per-arch-urls` @ `9a5f34ff`) pins
  `FreedomTechFeed/packages` at `v0.6.0-alpha1`. That tag does not exist; **8/8 of its URLs 404**.
- **F2 — The decision note is missing**, so the card's own "documented interim exception" clause
  is undecidable.

---

## What I ran, and what came back

### Precondition (step 0) — UNVERIFIED

```
$ read_file ~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md
File not found ... similar_files: [".../wz1.5-url-inventory.md"]

$ ls ~/.hermes/plans/2026-08-31-tollgate-batch/
schedule.md
wz1.5-url-inventory.audit-2026-10-14.md
wz1.5-url-inventory.audit-2026-10-15.md
wz1.5-url-inventory.md
```

`schedule.md` has zero `wz1.5`/`wz15` matches, so the decision is not recoverable from there.

Every remote branch's `deploy.go` was scanned for `felixfelix-bot`:

```
origin/main                            felixfelix=6 official=1 sha=65b9d02
origin/feat/feed-per-arch-urls         felixfelix=6 official=1 sha=9a5f34f
origin/fix/remove-minibits-dleq        felixfelix=0 official=1 sha=f43e020  <-- only zero-hit branch
```

`fix/remove-minibits-dleq` is not the change under test: `git rev-list --count` gives **0 ahead /
44 behind** main, and it pins an older unrelated set (`v0.5.0` upstream ipk, `releases.tollgate.me`
os URL with a `?channel=stable` suffix, `net4sats/configurationwizzard` v1.0.0).

```
$ git ls-remote origin | grep -iE 'wz1\.5|wz15|official|feed-per'
9a5f34ff...  refs/heads/feat/feed-per-arch-urls
9a5f34ff...  refs/tags/archive/feat-feed-per-arch-urls
```

No `wz1.5` ref exists on the remote.

### Step 1 — Grep gate — FAIL

On the pushed revision `origin/main`:

```
$ git show origin/main:deploy.go | grep -n 'orangesync\|felixfelix-bot'   # EXIT=0
28:	//   Fallback — felixfelix-bot fork releases for the v0.6.1-post-merge tag
34:	// from felixfelix-bot to OpenTollGate (keeping the same path/asset name).
43:	tollgatePkgURL = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.7.0-alpha10/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk"
47:	tollgatePkgURLApk = "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.6.1-post-merge/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk"
50:	configwizURL = "https://github.com/felixfelix-bot/configurationwizzard/releases/download/v1.0.10-alpha/net4sats-configwiz-1.0.10-alpha.tar.gz"
63:// Digests computed 2026-08-26 from the felixfelix-bot release assets:
```

5 hits, not 0. `orangesync` is 0 hits on both refs (correct — it never existed on this lineage).
Repo-wide: 39 files (listed in full in the deliverable), material ones being `deploy.go`,
`deploy_test.go`, `pins_test.go`, `CHANGELOG.md`, `docs/router-setup.md`, and the tracked binary
`net4sats-wizard`.

There is no documented interim exception to quote — the `deploy.go:25-40` comment block describes
an upstream-first/fork-as-interim *strategy* but never states that these three constants are an
accepted grep-gate exception.

### Step 2 — Version pinning — PASS on shipped code, FAIL on the official ref

Two versions of the same constant, same path shape:

```
--- d152f1e   configwizURL = ".../configurationwizzard/releases/download/v1.0.8-alpha/net4sats-configwiz-1.0.8-alpha.tar.gz"
--- 8b3220f   configwizURL = ".../configurationwizzard/releases/download/v1.0.10-alpha/net4sats-configwiz-1.0.10-alpha.tar.gz"
--- 65b9d02   configwizURL = ".../configurationwizzard/releases/download/v1.0.7-alpha/net4sats-configwiz-1.0.7-alpha.tar.gz"
```

No rolling target:

```
$ grep -rnE 'releases/latest|/latest/|channel=' --include='*.go' . --exclude-dir=.git
(empty)
```

No hash-pinning of a moving blob: `artifactPins` pins sha256 per URL and fails closed on a missing
entry. Caveat — `configwizURL` is *not* in `artifactPins` and is never passed to `httpGetFile`
(only `pins_test.go` references it), so the admin artifact is unpinned by the binary. Pre-existing,
not caused by wz1.5.

The official-channel ref fails:

```
$ curl .../releases/tags/v0.6.0-alpha1 | grep message
  "message": "Not Found",
$ curl -o /dev/null -w '%{http_code}' .../download/v0.6.0-alpha1/tollgate-wrt_0.6.0_alpha1_aarch64_cortex-a53.ipk
404
$ curl .../releases?per_page=100 | grep tag_name
    "tag_name": "v0.6.0",        (2026-10-21T02:06:45Z)
    "tag_name": "v0.6.0-rc6",    (2026-10-15T19:50:37Z)
```

Only `v0.6.0` and `v0.6.0-rc6` exist. The ref's asset naming is also structurally wrong: it uses
`tollgate-wrt_0.6.0_alpha1_<arch>.ipk`, while real assets are
`tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk` (extra target segment, no `alpha1`).

### Step 3 — Install-step dry-run — UNVERIFIED (no dry-run mode); download+verify measured

The wizard has no dry-run mode:

```
$ grep -rniE 'dry.?run' --include='*.go' . --exclude-dir=.git | wc -l
0
```

Install is an SSH sequence (`httpGetFile(selectedPkgURL)` at `deploy.go:402` → sha256 verify →
`opkg install --force-downgrade --force-reinstall --force-overwrite --force-depends` at
`deploy.go:508`). No router was contacted. I executed the download+digest half locally instead:

```
ipk  http=200 size=7475197  sha256 786565ec18f929c938c53d6b0303115ecaa8c238809a8d7e27fd79bbca100cd6
     deploy.go:67 pin 786565ec18f929c938c53d6b0303115ecaa8c238809a8d7e27fd79bbca100cd6  -> MATCH
cw   http=200 size=966774   sha256 87e5b9c8520400da6748dfb5743a7bddc474ebd4d9a6914a8b32233063b6c153
     (no pin to compare — see above)
```

Resolution matrix (network was available, so nothing is UNVERIFIED-for-network):

```
official, as the ref constructs it  -> 404  (size 9, text/plain)
official, corrected v0.6.0 name     -> 302  -> signed release-assets URL (8,591,306 B per API)
pinned fork ipk (what main ships)   -> 200  7,475,197 B  application/octet-stream
pinned admin tar.gz                 -> 200    966,774 B  application/octet-stream
```

### Step 4 — Docs in the same commit as code — UNVERIFIED

No commit changes the constants, so the property is vacuous. HEAD `681ed99` is docs-only:

```
$ git show --stat --oneline 681ed99
 PROGRESS-t_d41ef5b5.md | 70 ++++   README.md | 3 +   REPORT-t_d41ef5b5.md | 277 +++++
 docs/evidence/2026-10-23-fe2b-measured-captures.md | 286 ++++   docs/feed-install.md | 438 ++++
 5 files changed, 1074 insertions(+)
```

No `deploy.go`, no code path. And:

```
$ git log --all --oneline -G'configwizURL = "https://github.com/(OpenTollGate|FreedomTechFeed)' -- docs/ deploy.go
(empty)
```

---

## Line-item summary

| # | Item | Verdict |
|---|---|---|
| 0a | Decision note present | **UNVERIFIED** (ENOENT) |
| 0b | Implementation SHA on remote | **UNVERIFIED** (no ref drops the fork; no `wz1.5` ref) |
| 1a | `deploy.go` = 0 hits or documented exception | **FAIL** (5 hits) |
| 1b | `orangesync` absent from `deploy.go` | **PASS** (0 hits) |
| 1c | Repo-wide hits reported separately | **PASS** (39 files listed) |
| 2a | ≥2 constructed versions | **PASS** (v1.0.7 / v1.0.8 / v1.0.10) |
| 2b | No rolling `latest`/`channel=` | **PASS** (0 hits) |
| 2c | No hash-pinning of moving blob | **PASS** (`artifactPins`, fail-closed) |
| 2d | Official URLs pin a real tag | **FAIL** (8/8 → 404) |
| 3a | Install-step dry-run performed | **UNVERIFIED** (no dry-run mode; nothing installed) |
| 3b | Versioned artifact resolves | **PARTIAL PASS** (current pins 200 + digest match; official 404 as constructed) |
| 4  | Docs + code same commit | **UNVERIFIED** (no such commit) |

---

## Cross-check against prior board state

Gate-8 (`t_222a7d2b`, `done`) reached the same conclusion independently: "MERGE REFUSED … the
change was never pushed … main b436512 still pins felixfelix-bot at deploy.go:43/:47/:50 …
Validation report absent (t_deabf3d6 triage/empty)". My run confirms the substance and corrects
one detail: `b436512` is **not an ancestor** of current `origin/main` (`65b9d02`) —
`git merge-base --is-ancestor b436512 origin/main` → NO. The revision to gate against is
`65b9d02`, which still carries all three fork pins, so the verdict is unchanged.

Also worth flagging: the stale-pin claim in that verdict ("later releases exist") is now resolved
differently — the official channel's *current* release is `v0.6.0`, and its asset naming differs
from what `feat/feed-per-arch-urls` constructs. Any re-pin must re-verify names, not just bump a
version string.

---

## Route-back

1. `t_7e5033df` — actually push the change, re-pinned to a real tag (`v0.6.0` naming verified
   against the release API), with an explicit interim-exception comment if any fork URL remains.
2. `t_1df48589` — write `wz1.5-channel-decision.md` so step 1's exception clause is decidable.
3. Re-dispatch `t_deabf3d6` against the real head SHA; then `t_7da567aa` (cold review), then
   re-run Gate-8 `t_19b1b1c6`.

`/tmp/wz15-dryrun/{tol.ipk,cw.tar.gz}` is left in place so re-validation is a re-run, not a
re-download.

## Files touched

- `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-validation.md` (new — the deliverable)
- `REPORT-t_deabf3d6.md` (this file), `PROGRESS-t_deabf3d6.md`, appended section in `REPORT.md`
- `/tmp/wz15-dryrun/` (scratch artifacts)
- **No repo source file modified.**
