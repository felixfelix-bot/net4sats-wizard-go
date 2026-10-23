# PROGRESS — t_deabf3d6 (Validate URL change: grep gate + install-step dry-run)

Session 2026-10-23 (fleet offload, worker-tollgate). Finding -> Status -> Files touched.
Prior sections below belong to earlier tasks sharing this worktree root (preserved whole).

Scope: read-only validation. No source file was modified.

1. Finding: The card's two named inputs are both absent. The decision note
   `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md` is ENOENT (the plans dir
   holds only `schedule.md` and the three `wz1.5-url-inventory*` files; `schedule.md` has zero
   `wz1.5` matches). No implementation commit SHA is reachable on the remote.
   Status: BLOCKED at precondition, reported rather than passed.
   Files touched: none.

2. Finding: No ref on `origin` carries the official-URL change. Scanned every remote branch's
   `deploy.go` for `felixfelix-bot` hits: `origin/main` = 6, and the only zero-hit branch
   (`origin/fix/remove-minibits-dleq`) is 44 commits behind main and pins unrelated older URLs.
   `git ls-remote origin | grep -iE 'wz1\.5|wz15'` returns nothing.
   Status: step 0b UNVERIFIED.
   Files touched: none.

3. Finding: Grep gate measured on the pushed revision anyway. `origin/main:deploy.go` has 5
   `felixfelix-bot` hits at lines 28, 34, 43, 47, 50, 63; `orangesync` has 0 hits. So "0 hits, or
   only the documented interim exception" is FAIL, not the expected pass.
   Status: step 1 FAIL (with `orangesync` sub-check PASS).
   Files touched: none.

4. Finding: Version pinning is correct on the shipped code but wrong on the official-channel ref.
   URL shape is `releases/download/<TAG>/<ASSET>` with explicit tags (v1.0.7-alpha / v1.0.8-alpha /
   v1.0.10-alpha shown), and `grep -rnE 'releases/latest|/latest/|channel=' --include='*.go'`
   returns zero hits. But `origin/feat/feed-per-arch-urls` @ 9a5f34ff hard-codes `v0.6.0-alpha1`
   for 8 feed URLs, and that tag does not exist: API reports only `v0.6.0` (2026-10-21) and
   `v0.6.0-rc6` (2026-10-15). Pinned aarch64 ipk URL = HTTP 404.
   Status: step 2a/2b/2c PASS, step 2d FAIL.
   Files touched: none.

5. Finding: The wizard has no dry-run mode at all — `grep -rniE 'dry.?run' --include='*.go'`
   returns 0, and install is SSH `opkg install` at deploy.go:508. So a true install-step dry-run
   is impossible as shipped. Executed the download+digest-verify half locally in /tmp/wz15-dryrun
   instead: ipk HTTP 200, 7,475,197 B, sha256 786565ec...a100cd6 == the deploy.go:67 pin (match);
   configwiz tar.gz HTTP 200, 966,774 B (no pin exists for it — configwizURL is not passed to
   httpGetFile). Network was available, so nothing is UNVERIFIED-on-network.
   Status: step 3a UNVERIFIED (no dry-run mode), step 3b PARTIAL PASS.
   Files touched: /tmp/wz15-dryrun/tol.ipk, /tmp/wz15-dryrun/cw.tar.gz (scratch only).

6. Finding: Step 4 cannot hold — no commit anywhere changes the URL constants, so docs+code
   coupling is vacuous. HEAD commit 681ed99 is docs-only (5 files, 0 under deploy.go).
   `git log --all -G'configwizURL = "https://github.com/(OpenTollGate|FreedomTechFeed)' -- docs/ deploy.go`
   is empty.
   Status: step 4 UNVERIFIED.
   Files touched: none.

7. Finding: Deliverable written with raw commands + raw output, every line item labelled
   PASS/FAIL/UNVERIFIED.
   Status: DONE.
   Files touched: ~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-validation.md,
   REPORT-t_deabf3d6.md (this worktree root), PROGRESS-t_deabf3d6.md.

VERDICT: BLOCKED — no URL change exists to validate. Downstream cards t_7da567aa (cold review)
and t_19b1b1c6 (Gate-8) cannot proceed until t_7e5033df actually pushes, re-pinned to a real tag.
