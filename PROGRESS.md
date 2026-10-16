# PROGRESS — WZ1.6 e2e integration gate, re-dispatch (t_0509c320)

Worktree: /home/c08r4d0r/wz16-verify · branch ci/add-github-actions-workflow-clean
Tip verified: 4ca0b381ab3aa40f00b74f133f7b9024740277d7
Evidence: /home/c08r4d0r/wz16-redispatch/captures/

## Finding 1 — task's repo literal does not exist → RESOLVED (no file changes)
- `gh run list --repo net4sats/tollgate-wizard-go` → HTTP 404 Not Found; GraphQL "could not
  resolve to a Repository". Real repo = net4sats/net4sats-wizard-go (public); work lands on the
  fork felixfelix-bot/net4sats-wizard-go, which is this worktree's `origin`.
- Status: DONE. Files: none.

## Finding 2 — the recorded "CI green" tip was STALE → re-dispatched on the true tip
- On arrival the newest green CI run (38834487839) was at headSha 11c661b — one docs commit
  BEHIND the current tip 4ca0b38. Query over 20 runs: `runs matching tip: 0`. So the earlier
  report's green was not evidence about the shipped head.
- Action: re-dispatched CI on the exact tip (`gh workflow run .github/workflows/ci.yml --ref
  ci/add-github-actions-workflow-clean`), watched to completion.
- Result: run 38838000347 → conclusion success, headSha 4ca0b38…, jobs lint=success, test=success.
  Assertion: "runs at tip 4ca0b38 : 1 · AT LEAST ONE SUCCESS AT TIP: True".
- Status: DONE. Files: none (no workflow change needed — CI was already correct; it had just
  never been run on this commit).

## Finding 3 — all 5 gate steps re-executed green against the current tip
- Step 1 clean tree: `git status --porcelain` empty, push "Everything up-to-date", ahead/behind 0 0.
- Step 2 suite: `go vet` rc=0, `go test ./... -race` rc=0 (ok 121.327s), `golangci-lint run` rc=0
  (no issues), `goreleaser check` rc=0, `gofmt -l .` empty.
- Step 3 snapshot: `goreleaser release --snapshot --clean` rc=0, version 0.5.1-snapshot.4ca0b38,
  6 binaries + 6 archives; SHA256SUMS has 6 lines == 6 artifacts; `sha256sum -c SHA256SUMS` all OK
  rc=0; two-artifact spot check (linux_amd64 archive + darwin_arm64 NON-linux) OK; negative control
  on a byte-tampered copy FAILED with rc=1 (check is real, not a rubber stamp).
- Step 4 smoke: `--version` → `tollgate-wizard 0.5.1-snapshot.4ca0b38`, rc=0, not `dev`; manual
  `-ldflags -X main.version=v9.9.9-manual-test` build prints the injected value; un-injected build
  prints `dev` (proves ldflags is what changes it).
- Status: DONE. Files: none — this is a verification task; no source or config was modified.

## Finding 4 — summary comment posted
- https://github.com/net4sats/net4sats-wizard-go/pull/35#issuecomment-6206347289
- Note: `gh pr checks 35` still reports "no checks reported" — Actions do not fire on PRs in this
  fork/upstream arrangement, so the gate's CI evidence is the branch-head run, at the same commit.
- Status: DONE.

## Finding 5 — no workflow fix was required, and none was made
- CI was NOT red anywhere in this run, so the "fix the workflow, do not weaken the checks" branch of
  the task never triggered. Zero changes to .github/ or .goreleaser.yaml.

## Milestones
- (this commit) docs(wz16): e2e gate re-dispatch — all 5 steps green at tip 4ca0b38, CI re-run green
- Report: REPORT.md (worktree root)
