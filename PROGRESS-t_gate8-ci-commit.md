# PROGRESS — t_gate8-19b1b1c6 CI commit/push card

Task: "Commit and push CI workflow as atomic conventional commit. Done == pushed."
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6 (== origin)

Note: task id collides with the wz1.5 URL-inventory card that already owns
REPORT-t_gate8-19b1b1c6.md / PROGRESS-t_gate8-19b1b1c6.md in this tree, so this card's
prose lives under -t_gate8-ci-commit.

## Finding -> status -> files touched

1. Finding: the brief's premise was stale. The CI work was ALREADY committed and pushed by
   parents before this session: commit 13de851 "ci: add gha workflow for test/vet/lint on pr
   and main", 10 files (.github/workflows/ci.yml new, .golangci.yml new, + lint fixes in
   deploy.go discover.go embed.go main.go ssh.go and 3 test files).
   Status: VERIFIED (ancestor of HEAD; on origin; local HEAD 1a0c0ec == ls-remote tip).
   Therefore NO new commit was made and NO push was performed (no delta to push).

2. Finding: git status is NOT clean but the dirty set is unrelated and pre-existing --
   M REPORT.md (sibling card, mtime 22:56), ?? .hermes/ (07:59), ?? verify/wz16-e2e-gate/
   (2026-10-16). Change set of CI files = empty.
   Status: EXCLUDED deliberately. No git add -A, no stash, no discard. Committing them would
   sweep another card's work; deleting them would destroy it.

3. Finding: the workflow had NEVER executed on the landed commit -- actions/runs?head_sha=
   13de851 -> total_count 0; the repo's 5 green CI runs are all on a different branch.
   "Local commands pass" did not prove the workflow runs.
   Status: FIXED IN EVIDENCE TERMS -- ran the full matrix locally (build, vet, test,
   test -race, golangci-lint 1.61.0 == pinned version: all exit 0) and triggered a real run.
   Run 40025804189 on ci/github-actions-workflow-atomic @ 13de851: conclusion=success,
   jobs test=success lint=success.
   Files touched: none (verification only).

4. Finding: PR #34 is mergeable=CONFLICTING. 19 behind / 5 ahead of main; 9 conflicting
   paths, all SOURCE files (CHANGELOG.md, admin/index.html, artifact_pins_test.go, deploy.go,
   handlers_test.go, main.go, portal/index.html, portal/splash.html,
   security_hardening_test.go). ci.yml/.golangci.yml do not conflict.
   Status: REPORTED, not resolved. Resolving rewrites a pushed shared branch and is outside
   "commit and push".

5. Finding (most important): origin/main contains NEITHER .github/workflows/ci.yml NOR
   .golangci.yml. The only registered workflow object (376283443, active) points at a path
   that does not exist on the default branch -> the repo has NO CI gate on main today.
   Status: ESCALATED in report §4. Blocked on PR #34 landing.

6. Mistake made and corrected: I wrote this card's report to REPORT-t_gate8-19b1b1c6.md,
   which was already the wz1.5 inventory card's tracked file, overwriting it.
   Status: RESTORED byte-exact (git checkout HEAD --, blob 4b7787b7 == HEAD:...), verified
   identical, then re-issued under REPORT-t_gate8-ci-commit.md. No sibling data lost.

## Status

DONE for everything the brief can honestly mean here: the commit and push already existed and
are verified on origin, and the workflow is now proven green in real GitHub Actions rather than
assumed. The card's literal acceptance lines "git log -1 shows the conventional commit" and
"git status clean" are NOT met and cannot be met without faking a commit or discarding another
card's work; both are documented instead of papered over.

## FILES TOUCHED
- PROGRESS-t_gate8-ci-commit.md   (new, this file)
- REPORT-t_gate8-ci-commit.md     (new, full report)
- REPORT-t_gate8-19b1b1c6.md      (restored byte-exact after my overwrite; net unchanged)

## REMAINING
1. Resolve PR #34's 9 conflicts (rebase onto main; keep main's files, re-apply CI-only hunks).
2. Land #34 so main finally gets a workflow file.
3. Decide if CI should also land directly on main.
4. Confirm owner of REPORT.md / .hermes/ / verify/wz16-e2e-gate/ before anyone commits them.
5. Optional: bump actions pins off the Node-20-deprecated majors.
