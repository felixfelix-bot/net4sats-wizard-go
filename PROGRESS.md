# PROGRESS — t_6b512dad: Recover documented GL-MT6000 lock convention (read-only)

Finding -> status -> files touched. Worktree root: /home/c08r4d0r/repos/net4sats-mvp-v2

## Status: COMPLETE — all three items (a)(b)(c) recovered verbatim, with source notes. ZERO lock commands run, zero device contact.

## Findings

1. Searched kanban board `net4sats-mvp-v2` (titles/bodies/comments/attachments) and the
   workspace with each term: GL-MT6000, lock, recon, flint2, mt6000 (individually + combined).
   - `flint2` -> only the search-term boilerplate inside the task bodies themselves
     (t_6b512dad, t_bc00634f). No artifact by that name.
   - The convention is NOT documented in the net4sats-mvp-v2 repo itself; it lives in the
     recon note attached to card `t_cc28693c` ("Recon plan, upstream build artifacts, and
     hardware-lock convention") and in companion RECON docs committed under docs/.

2. Source note (authoritative) identified:
   - Card: `t_cc28693c` — "Recon plan, upstream build artifacts, and hardware-lock convention"
   - Attachment (id 99): `recon-note-2026-10-16.md` (pass 4) and (pass 5)
     `RECON-tollgate-batch-plan-artifacts-hwlock-device-2026-10-16-pass5.md`
   - On-disk (pass5): /home/c08r4d0r/worktrees/t_11d0873e/wt/docs/RECON-tollgate-batch-plan-artifacts-hwlock-device-2026-10-16-pass5.md
   - Attachment path (pass4 note the claim task cited):
     /home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/attachments/t_cc28693c/recon-note-2026-10-16.md
   - Re-verified on the filesystem this pass: lock file EXISTS
     (ls -la ~/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock -> 721 B, mode 600).
     I did NOT open/execute any lock command. Reads of public docs only.

3. (a) lock name/path + inspect command — recovered verbatim (see REPORT.md §Verbatim).
4. (b) device id/serial — recovered verbatim: device id `GL-MT6000`; serial is a
   DOCUMENTED ABSENCE (UNRECORDED; model string stands in for serial). Not a gap I introduced.
5. (c) shared channel / board column — recovered verbatim: NOT a channel, NOT a column.
   The recon note answers (c) with an explicit NO on both, naming the audit that concluded FAIL.

## Files touched (this task)
- PROGRESS.md (this file, worktree root)
- REPORT.md (worktree root)
- verify/t_6b512dad-lock-convention/recon-quote.txt (raw extracted note text, provenance)

## Constraint check
- ZERO lock commands executed. No claim/unlock/inspect command was run.
- No device contact (no ssh/scp/ping/curl/nc to any router). Reads were local files only.

═══════════════════════════════════════════════════════════════════════
## PRIOR CONTENT (preserved — previous card on this shared file)
═══════════════════════════════════════════════════════════════════════

# PROGRESS — t_gate8-19b1b1c6 → wz1.5 URL inventory (fleet offload, worker-heavy)

Session 2026-10-23. Finding -> Status -> Files touched. Read-only inventory: no source file in
this repo was edited, no test run, no build, no state-mutating git command issued.

1. Finding: the brief's six tokens produce **172 matching lines** repo-wide (case-insensitive,
   text files, `.git` excluded) — but only **26 are shipped/tracked repo content**. The other 146
   are worker scratch (`PROGRESS*/REPORT*/AUDIT*/EVIDENCE*` at the tree root) plus 12 rows of
   measured capture transcript under `docs/evidence/`.
   Status: VERIFIED (grep counts reproduced token by token).
   Files touched: none.

2. Finding: `orangesync` and `blossom2` have **ZERO hits in any shipped artifact** — 12 and 1
   lines respectively, all inside worker reports (the 12 are `relay.orangesync.tech` outage
   evidence in REPORT/PROGRESS-t_19b1b1c6 + two census reports). Nothing to replace or delete.
   Status: VERIFIED (`grep -rniI` per token + per-file tracked/ignored split).
   Files touched: none.

3. Finding: `tollgateOSURL` is **not present in `deploy.go`** (`grep -c` -> 0). It was added at
   `fa306d2` (2026-07-16, "constant for future firmware flash") with a
   `releases.tollgate.me/os/57e0f2468a17…?channel=stable` value and deleted at `d25c018`. Its only
   two remaining references are the tripwire in `pins_test.go:45` (comment) and `:49`
   (`forbiddenDeployIdentifiers`). Verdict: **keep the guard**; the name is a forbidden identifier
   in `deploy.go`, so a future firmware-image constant must not reuse it.
   Status: VERIFIED (`git log --all -S tollgateOSURL`, `git show d25c018^:deploy.go`).
   Files touched: none.

4. Finding: `configwizURL` (`deploy.go:50`) has **zero production references** — only its
   declaration plus two test-registry entries (`pins_test.go:31` `expectedPinnedURLConsts`, `:77`
   `liveCheckPins`). The admin panel ships **embedded** (`embed.go:14` `//go:embed all:admin`,
   deployed at `deploy.go:762` step 7a), superseding the "download from a GitHub release" design
   (`d70b8be` -> superseded by `adfaffb`). Verdict: **delete, do not re-pin** — 4 edits in one
   commit or `TestDeployGoPinRegistry` fails.
   Status: VERIFIED (`grep -rnIw configwizURL` over go/md/sh/Makefile = 3 lines; `artifact_pins_test.go:24`).
   Files touched: none.

5. Finding: **the official-channel decision is genuinely blocked, and a pure owner swap would
   break deploys.** `felixfelix-bot/tollgate-module-basic-go` is a GitHub **fork of
   `OpenTollGate/tollgate-module-basic-go`**, but the fork is *ahead*: upstream's newest
   `tollgate-wrt` .ipk is `v0.5.0` (HTTP 206) while the wizard pins `v0.7.0-alpha10`; upstream
   `v0.6.0`/`v0.6.0-alpha1` publish **no** `tollgate-wrt` asset at all. Swapping owner-only ->
   404 on every fresh deploy. `net4sats/configurationwizzard` (official, parent of the fork) does
   publish `net4sats-configwiz-1.0.0.tar.gz` (206) but that const is dead anyway.
   `releases.tollgate.me` alpha3 index is 206 while `keys/tollgate-testing.pem` is **404** (unsigned
   -> uninstallable). `FreedomTechFeed/feed` has 0 releases/tags.
   Status: VERIFIED (GitHub releases API + `curl -r 0-0 -L` status codes; fork lineage via repos API).
   Files touched: none.

6. Finding: all three current pins are still **live** (HTTP 206 this run), so nothing is broken
   today — the issue is provenance, not availability.
   Status: VERIFIED (`curl -sS -o /dev/null -w '%{http_code}' -r 0-0 -L` on each pin).
   Files touched: none.

7. Finding: the repo's `go` toolchain is **absent on this host** (`go: command not found`), so no
   build/test gate could be run. This task is read-only, so nothing was left unverified that the
   task required — but Phase 3 gates in the deliverable are explicitly marked UNVERIFIED.
   Status: VERIFIED (command error observed).
   Files touched: none.

8. Deliverable written: `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.md`
   — 26-row merged ledger (file:line, exact line, class a/b/c/d, action), a leave-alone register
   accounting for the remaining 146 lines by file, per-constant verdicts for `tollgateOSURL` /
   `configwizURL` / `releases.tollgate.me`, the measured official-channel census, and a 3-phase
   implementation checklist.
   Status: DONE (written; nothing to commit — the repo is untouched by design).
   Files touched: `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.md` (overwrote the
   stale 2026-10-15 edition), `PROGRESS-t_gate8-19b1b1c6.md`, `REPORT-t_gate8-19b1b1c6.md`
   (both worktree-only, untracked and gitignored via `/home/c08r4d0r/wz16-verify/.git/info/exclude`).


---

# PROGRESS — WZ1.4: remove the Minibits DLEQ filter, track the fix upstream

Card: `t_gate8-19b1b1c6` (worker-heavy) · Repo: `net4sats-mvp-v2` ·
Remote: `felixfelix-bot/net4sats-wizard-go` · Branch: `worker-heavy/t_gate8-19b1b1c6`

Finding -> Status -> Files touched. Section above belongs to the prior lane sharing
this worktree; preserved unmodified.

---

## 1. Finding: the "DLEQ filter" was already half-removed; what remained was
##    brand-specific copy, not a filter

Status: **VERIFIED first-hand.**

`git show HEAD:deploy.go` has **no** jq mint filter. The filter lived in `f43e020`
("fix: remove minibits from default mints due to DLEQ keyset rotation bug",
2026-07-17), which injected:

```
.accepted_mints = (.accepted_mints | map(select(.url | test(\"minibits\") | not)))
```

That commit **is** an ancestor of HEAD, but its filter line was already reversed:
`git log -S'test("minibits")' -- deploy.go` returns no commit, and the jq program
at `deploy.go:199-211` contains only URL de-duplication.

What actually remained was the workaround's *narrative*, stale and brand-specific
in two operator-visible places:

| Location | Text | Why it is a liability |
|---|---|---|
| `index.html:207` | mint hint: "... Minibits temporarily removed (DLEQ keyset rotation bug)." | Says a mint is removed when it is a listed production mint; names one mint's bug in a front-end meant to carry no per-mint special-casing |
| `deploy.go:882` | "Does NOT strip minibits (DLEQ keyset rotation bug fixed in gonuts v0.11.1)." | Narrates the exception instead of stating the general rule |

Files touched: read-only at this stage.

## 2. Finding: Q3b is (b) DON'T FILE, and its premise checks out

Status: **VERIFIED two independent ways.**

The brief carries Q3b = **(b) DON'T FILE — accept removal-only (upstream fix already
shipped)**. The plan file at `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`
did **not** reflect it — it still read `NOT RESOLVED` with an empty `ANSWER` slot.
Verified rather than trusted:

1. `tollgate-module-basic-go/src/go.mod:44` and `src/merchant/go.mod:6` both pin
   `github.com/OpenTollGate/gonuts-tollgate v0.13.0`.
2. In the module cache at
   `~/go/pkg/mod/github.com/!open!toll!gate/gonuts-tollgate@v0.13.0/cashu/nuts/nut12/nut12.go:43`,
   `VerifyProofsDLEQWithKeysets` groups proofs by `proof.Id` and fetches the
   matching historical keyset — precisely the fix the draft asked for. The v0.10.0
   threshold in the gate text is satisfied.

Files touched: none (read-only verification).

## 3. Finding: upstream issue #1057 already exists — filing would duplicate it

Status: **VERIFIED.** `gh api repos/OpenTollGate/tollgate-module-basic-go/issues/1057`

- Title: "DLEQ on the receive path: the proof.Id-keyed verify exists in gonuts
  v0.13.0 but is unreachable from this module"
- State: open · opened 2026-10-23T14:45:48Z · author `felixfelix-bot`
- Its body explicitly states it **supersedes**
  `docs/WZ1.4-minibits-dleq-issue-draft.md` (commit `ff2290bc`) — the approved
  draft's content is already carried upstream.

**Decision taken:** do **not** file a second issue. The acceptance criterion is
"upstream issue filed and URL recorded"; upstream tracking exists and is recorded.
Filing a duplicate would be the wrong act, and Q3b already says don't file. I also
reproduced #1057's central claim myself (`grep -rli dleq --include='*.go' src/` →
zero matches), so the recorded issue is accurate, not merely present.

Files touched: `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` (URL +
verification recorded; header and gate-body status corrected from `NOT RESOLVED`
to resolved, since three downstream cards grep that heading).

## 4. Finding: removal + guards

Status: **DONE, tests GREEN and verified RED before the fix.**

| File | Change |
|---|---|
| `index.html` | Mint hint reduced to the neutral one-liner; brand-specific note gone |
| `deploy.go` | Comment states the general rule ("no per-mint special-casing") instead of the Minibits exception |
| `defaults_test.go` | `TestNoMintIsFilteredFromConfigCommand` + `TestNoMinibitsDleqCopyInVisibleUI` |
| `CHANGELOG.md` | Unreleased entry, no internal board paths (public repo) |

The first guard rejects the filter *construct* (a regex `test()` against a mint URL)
and asserts every production mint — Minibits included — reaches the emitted config
command. The second rejects brand-specific mint copy in operator-visible UI.

**RED/GREEN evidence.** Both tests run against a scratch copy with the old filter and
old hint re-injected: both FAIL (`config jq program still filters mints by URL
pattern ".url | test("`; `operator-visible UI still carries brand-specific mint copy
"minibits"/"dleq"`). Against the fixed tree both PASS.

An earlier draft of the first guard falsely failed on the words `select(` and
`minibits` — the jq program legitimately uses `select(...)` for `profit_share` and
URL de-duplication, and Minibits is a legitimate listed mint. Corrected to target
the filtering construct; recorded because the naive version would have pushed a
wrong "fix".

## 5. Toolchain note

`go` is **not on PATH** on this host; it is at `/usr/local/go/bin/go` (go1.22.12).
Prior lanes recorded "go toolchain absent" — that is wrong, it is merely unexported
(`export PATH=$PATH:/usr/local/go/bin`). Builds and tests ran fine.

## 6. Incident (disclosed): I overwrote an ignored root `PROGRESS.md`

Status: **PARTIALLY RECOVERED — reporting it rather than hiding it.**

My first write used `write_file` on the worktree-root `PROGRESS.md`, which is
gitignored (via `/home/c08r4d0r/wz16-verify/.git/info/exclude`) and therefore not
recoverable from git objects (`git fsck` dangling-blob scan found nothing). It held
the prior lane's progress text for this same worktree root. Recovered content is a
copy of `PROGRESS-t_gate8-19b1b1c6.md` (4,699 B), which is that lane's own
task-scoped sibling and very likely the same text — but I cannot prove byte
identity, and I did not observe the original size before overwriting it.

The correct pattern (per this repo's own convention, which I should have followed
first) is a task-scoped `PROGRESS-t_<task>.md`. Both files now exist: the restored
root file plus this section, and the task-scoped pair for this card.

## 7. Commits and push

| SHA | Subject | Pushed |
|---|---|---|
| `c56ce38` | `fix(ui): drop the Minibits DLEQ workaround copy from the wizard` | yes — `fea18ff..c56ce38` |
| _(this commit)_ | `docs(wz1.4): record upstream issue #1057 and close out the DLEQ workaround` | see REPORT.md |

---

## Status

**COMPLETE.** `go test ./...` PASS on the changed tree; no Minibits-specific DLEQ
filtering or copy remains; upstream tracking recorded (#1057, not duplicated); plan
file reflects the real decision. Both milestones committed and pushed.

## Files touched

- `index.html`, `deploy.go`, `defaults_test.go`, `CHANGELOG.md`
- `PROGRESS.md` (this root file, appended), `PROGRESS-t_gate8-19b1b1c6.md`,
  `REPORT.md` (appended), `REPORT-t_gate8-19b1b1c6.md`
- `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` (plan ref, outside repo)

## Remaining (none blocking this card)

1. Upstream `#1057` is open and owned by the backend lane: it asks which code path
   verifies DLEQ on the receive path, since the keyed verify is not called from the
   module. Nothing in this repo blocks on it.
2. `t_cfb0a06a` ("set dev-split default to 0") is REJECTED per Q3 = B — needs
   closing/re-scoping on the board, not by code in this repo.
3. `t_58a11d0d` (WZ1.4 parent) unblocks for its other scope (de-brand business
   logic), not for a default change.
