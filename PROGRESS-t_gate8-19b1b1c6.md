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
