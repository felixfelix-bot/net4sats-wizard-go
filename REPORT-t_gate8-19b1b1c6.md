# REPORT — t_gate8-19b1b1c6 / wz1.5 URL inventory

## Deliverable

`/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.md` (33 KB).
This **replaces** the stale 2026-10-15 edition of the same path, which described a base tree
(`2d77096…`, PR #32 head) that is no longer what the wizard branch carries: it listed 15
`felixfelix-bot` hits and 1 `releases.tollgate.me` hit, versus 11 and 5 today, and it predates
`tollgatePkgURLApk` (the OpenWrt-25 apk pin) entirely.

Task type: **read-only inventory.** No file in this repo was modified, no test was run, no git
state-changing command was issued. The repo working tree is byte-identical to session start.

## Acceptance criteria — status

| criterion | status | evidence |
|---|---|---|
| repo-wide grep for the six tokens fully accounted for | **MET** | 172 lines total; 26 in the merged ledger (§3), 146 in the leave-alone register (§4). Per-file counts in §4 sum back to the totals in §1. |
| implementer can act without re-grepping | **MET** | every ledger row carries `file:line`, the exact line, a class (a/b/c/d) and the concrete action; §7 sequences the edits into 3 phases; §5.2 lists all four required edits for the one constant that must be deleted. |
| per-constant verdict for `tollgateOSURL` and `configwizURL` | **MET** | §5.1 and §5.2, each with its own verdict, reference census and recommendation. |

## Headline findings

1. **The token sweep is mostly noise, and saying so matters.** 172 lines match, but **only 26 are
   shipped/tracked repo content**. The rest is worker scratch (`PROGRESS*`/`REPORT*`/`AUDIT*`/
   `EVIDENCE*` at the tree root; 12 of those files are untracked *and* gitignored via
   `/home/c08r4d0r/wz16-verify/.git/info/exclude`) plus 12 rows of measured capture transcript in
   `docs/evidence/2026-10-23-fe2b-measured-captures.md`. An implementer who acts on the raw grep
   count would rewrite other cards' evidence.

2. **`orangesync` and `blossom2` are absent from every shipped artifact — 0 hits.** Their 12 and 1
   occurrences are all inside worker reports, and the 12 are outage evidence for
   `relay.orangesync.tech` (HTTP 502), not artifact references. Nothing to replace or delete.

3. **Zero hits in CI, git hooks, Makefile, install scripts, `release/`, `artifacts/`, `admin/`,
   `portal/`, `rpcd/`, or embedded UI assets.** Measured with an explicit sweep (M7), not assumed.

4. **`tollgateOSURL` — already dead and already deleted; only the tripwire remains.** Not in
   `deploy.go` (`grep -c` → 0). Added `fa306d2`, deleted `d25c018`; the two surviving references are
   `pins_test.go:45` (prose) and `:49` (`forbiddenDeployIdentifiers`). **Verdict: keep both lines.**
   Consequence: the name is a *forbidden identifier* inside `deploy.go`, so if the official channel
   turns out to be `releases.tollgate.me`, the new firmware-image constant must be named something
   else or the forbidden list must be amended.

5. **`configwizURL` — genuinely unused *and* unofficial. Verdict: delete, do not re-pin.**
   Reference census: `deploy.go:50` (declaration), `pins_test.go:31`, `pins_test.go:77`. **Zero
   production references.** The admin panel is embedded (`embed.go:14` → `adminFS`, deployed at
   `deploy.go:762`), and `artifact_pins_test.go:24` already asserts the install-path download set is
   exactly `{tollgatePkgURL, tollgatePkgURLApk}`. The `pins_test.go:77` live-check entry is what
   makes the constant *look* load-bearing — `TestPinnedURLsAreLive` issues a real request to it —
   while nothing downloads it. Deletion is **4 edits in one commit** or `TestDeployGoPinRegistry`
   goes red (declared set == registry == live-check map).

6. **The official-channel replacement is not a drop-in, and this is the one thing that can break
   every deploy.** `felixfelix-bot/tollgate-module-basic-go` is a GitHub **fork of
   `OpenTollGate/tollgate-module-basic-go`** (fork lineage confirmed via the repos API), so
   `felixfelix-bot` *is* the upstream project in a fork namespace. But the fork is **ahead**:
   upstream's newest `tollgate-wrt` .ipk is `v0.5.0` (206), upstream `v0.6.0` / `v0.6.0-alpha1`
   publish **no** `tollgate-wrt` asset, and upstream `v0.5.0` also carries an apk. So:
   owner-swap-only → **404** (exactly the SW4a outage `TestPinnedURLsAreLive` exists to catch);
   swap to upstream's real newest → **silent two-release downgrade**; and any genuinely different
   asset needs a **newly computed** digest, never the fork digest carried over.
   `net4sats/configurationwizzard` (official, parent of the fork) does publish
   `net4sats-configwiz-1.0.0.tar.gz` (206) — irrelevant, since that const is dead.
   `releases.tollgate.me` serves the alpha3 index (206) but `keys/tollgate-testing.pem` is **404**
   (unsigned → uninstallable). `FreedomTechFeed/feed`: 0 releases, 0 tags.

7. **All three current pins are live right now (HTTP 206), so nothing is broken today** — the
   problem is provenance, not availability. That is why wz1.5 is a *change*, and also why it can
   wait for the channel decision without an outage.

## Test evidence

**No test suite was run.** Two reasons, both observed rather than assumed:

* This task is a read-only inventory; it edits no source, so there is no changed behaviour to
  verify.
* The Go toolchain is **absent on this host**: `go version` → `/usr/bin/bash: line 4: go: command
  not found`. `go build ./...`, `go vet ./...`, `go test ./...` and the repo's `make` targets are
  all unrunnable here.

Consequently the three verification gates the implementer must run (`go test ./...` without
`-short`, `TestDeployGoPinRegistry`, and a `strings ./net4sats-wizard | grep -c felixfelix-bot` → 0
rebuild check, since the ignored binary reports `2` today) are marked **UNVERIFIED** in §7/§8 of the
deliverable.

## Git state

**No source commit** (read-only task), but the two task-scoped report files **are committed and
pushed** — observed, not assumed:

```
commit 7c5cbdc  docs(wz1.5): t_gate8-19b1b1c6 — exhaustive URL inventory (26 actionable rows, 146 leave-alone)
        PROGRESS-t_gate8-19b1b1c6.md | 68 +++++++
        REPORT-t_gate8-19b1b1c6.md   | 134 +++++++++
push    c56ce38..7c5cbdc  HEAD -> worker-heavy/t_gate8-19b1b1c6   (origin)
verify  git ls-remote origin refs/heads/worker-heavy/t_gate8-19b1b1c6 -> 7c5cbdc7cf88ae11f65f6bdd25402295f8aab8b5
        git rev-parse HEAD                                             -> 7c5cbdc7cf88ae11f65f6bdd25402295f8aab8b5   (equal)
```

Staged explicitly by path (`git add -f PROGRESS-t_gate8-19b1b1c6.md REPORT-t_gate8-19b1b1c6.md`);
`git diff --cached --name-only` confirmed exactly those two files before committing — the sibling's
in-flight `deploy.go`/`defaults_test.go`/`index.html` edits were never staged. (That sibling committed
them itself as `c56ce38` mid-session, so they no longer appear in `git status`.) The pre-commit
cred-gate passed: `clean (30 literals checked, mode=staged)`.

Remaining `git status --short`: `?? .hermes/` and `?? verify/wz16-e2e-gate/` only — both pre-existing,
neither mine.

Observed facts:

* `git branch --show-current` → `worker-heavy/t_gate8-19b1b1c6`; `HEAD = 15f5d29` at inventory time
  and `fea18ff` at report time — the branch moved **twice** under me (`15f5d29` → `fea18ff`) as
  sibling cards committed. Both transitions were checked mechanically: `git grep -nE <tokens>` at
  `15f5d29` and at `HEAD` yields an **identical** 26-line actionable set (`added: NONE, lost: NONE`),
  and `git diff --stat 15f5d29..HEAD` touches only `PROGRESS-t_e75c5cd6.md`,
  `PROGRESS-t_f53c5650.md`, `REPORT-t_e75c5cd6.md`, `REPORT-t_f53c5650.md` — four sibling worker
  reports, zero source files.
* `git status --short` at report time shows `M defaults_test.go`, `M deploy.go`, `M index.html` —
  **not mine.** Inspected: they are a sibling worker's Cashu mint-filter change (`deploy.go:878`
  jq/copy comment, the `productionMintsJSON` block ~`:895-905`, the `index.html:207` hint line, and
  a new `TestNoMintIsFilteredFromConfigCommand` in `defaults_test.go`) — a live edit belonging to a
  different card. It overlaps this task's file `deploy.go` but **none of its hunks are within the
  six token lines**, so the inventory's line numbers remain valid. Not touched, not reverted,
  not staged.
* The sibling hunks were located exactly (`git diff -U0 | grep '^@@'`) to prove they do not collide
  with this task: `deploy.go` hunk starts at `:881` (this inventory's rows are `:28`–`:63`,
  `:762`); `defaults_test.go` inserts 45 lines at `:260`; `index.html` touches `:207`.
* The two report files this task wrote (`PROGRESS-t_gate8-19b1b1c6.md`,
  `REPORT-t_gate8-19b1b1c6.md`) are **untracked and gitignored** (`git check-ignore -v` →
  `/home/c08r4d0r/wz16-verify/.git/info/exclude`). They are worktree-only, per repo convention —
  they do **not** exist on any remote.
* The plan deliverable lives outside the repo entirely
  (`~/.hermes/plans/2026-08-31-tollgate-batch/`, not a git repository), so it cannot be committed
  or pushed from here by construction.

## What could NOT be established

* **Which channel is official.** No document, ticket or constant in this repo names the replacement
  host/owner. §6 of the deliverable therefore reports what each candidate *actually publishes*
  rather than picking one on the implementer's behalf.
* **Whether the chosen replacement asset's digest matches the current pin.** Must be computed from
  real downloaded bytes at implementation time.
* **Build/test status under any change** — no Go toolchain on this host.

## Recommended next step

Phase 1 of the deliverable (§7) is channel-independent and can land immediately: delete
`configwizURL` and its three registry lines, drop the `net4sats-configwiz-1.0.1.tar.gz` row from
`docs/router-setup.md`, and mark the `CHANGELOG.md` retained-pin note superseded. Phase 2 (the
repins) is blocked on the channel decision; §6 shows why it cannot be guessed.
