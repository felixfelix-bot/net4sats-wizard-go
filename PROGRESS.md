# PROGRESS — t_0509c320: End-to-end verify (snapshot build, --version smoke, full suite, push)

Board: `net4sats-mvp-v2` · Card: `t_0509c320` · Assignee: worker-heavy · Author: worker-heavy (fleet offload)
Host: x280 (`c08r4d0r`) · Date: 2026-10-15
Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate` (branch `worker-heavy/t_0509c320`, cut at `0b5322a`)
Subject under test: `0b5322a` = `origin/ci/add-github-actions-workflow-clean` = PR #35 head
Raw captures: `/home/c08r4d0r/wz-captures-t_0509c320/`

## Finding 1 — the repo named in the card does not exist
- `gh repo view net4sats/tollgate-wizard-go` → `GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'`.
- `git ls-remote https://github.com/net4sats/tollgate-wizard-go` → `remote: Repository not found.` (rc 128).
- `gh run list --repo net4sats/tollgate-wizard-go --limit 5` → `failed to get runs: HTTP 404: Not Found`.
- Real repo = **`net4sats/net4sats-wizard-go`** (public, default branch `main`) — the name the lane docs already alias.
- Status: RESOLVED. Files touched: none.

## Finding 2 — upstream cannot report CI checks on this PR; the fork is the only CI venue
- `gh workflow list --repo net4sats/net4sats-wizard-go` → **empty**.
- `gh api repos/net4sats/net4sats-wizard-go/contents/.github` → 404: upstream `main` has no `.github/` at all.
- `gh pr checks 35 --repo net4sats/net4sats-wizard-go` → `no checks reported on the 'ci/add-github-actions-workflow-clean' branch`.
- `gh pr view 35 … --json statusCheckRollup` → `[]`.
- So the CI workflow this PR adds cannot run on the PR itself yet (a fork PR's `pull_request` run needs the workflow on the base repo's default branch, which is the very thing PR #35/#34 are trying to land).
- Green CI therefore has to be produced and read on the fork: `felixfelix-bot/net4sats-wizard-go`, workflow `CI` (id 376283443, active).
- Status: RESOLVED (worked around with an explicit fork-side `workflow_dispatch`; no upstream change made, no check weakened).

## Finding 3 — golangci-lint 1.61.0 cannot read go1.24 export data
- `golangci-lint run` with go1.24.4 in PATH → rc 3:
  `buildir: failed to load package goarch: could not load export data: internal error in importing "internal/goarch" (unsupported version: 2)`.
- Environment/toolchain mismatch, NOT a code defect. CI pins `go-version-file: go.mod` = `go 1.22`, and `golangci-lint-action@v6` pins v1.61.0.
- Fix: run the lint gate under the CI-matching toolchain `GOTOOLCHAIN=go1.22.12` with a fresh `GOCACHE`. Then `golangci-lint run` → **rc 0, no findings**.
- Status: RESOLVED. Files touched: none (toolchain shim at `~/wz16-tools/go122/bin/go`, outside the repo).

## Finding 4 — no source change was required
- Every gate is green at `0b5322a` as-is. Zero source edits, zero test edits, no workflow weakening.
- The only commits on `worker-heavy/t_0509c320` are this task's own PROGRESS/REPORT docs.
- The branch under test (`ci/add-github-actions-workflow-clean`) is left byte-identical: PR #35 head stays `0b5322a`.
- Status: DONE. Files touched: `PROGRESS.md`, `REPORT.md` (worktree root, task branch only).

## Gate status

| # | Gate | Result |
|---|------|--------|
| 1 | clean tree + pushed, ahead/behind 0/0 | PASS |
| 2a | `go vet ./...` | PASS (rc 0) |
| 2b | `go test ./... -race -count=1` | PASS (123/123, 0 FAIL, 0 SKIP, 0 data races) |
| 2c | `golangci-lint run` (v1.61.0 under go1.22.12) | PASS (rc 0) |
| 2d | `goreleaser check` | PASS (rc 0, 1 config validated) |
| 3 | `goreleaser release --snapshot --clean` | PASS (6/6 archives, SHA256SUMS 6/6 lines) |
| 3b | checksum verification (`sha256sum -c` + independent recompute) | PASS (6/6 OK; 3/3 recomputed MATCH) |
| 4 | binary `--version` smoke + manual injection | PASS (`tollgate-wizard 0.5.1-snapshot.0b5322a`, rc 0) |
| 5 | CI green at the pushed head SHA | PASS (fork run 38677575488, conclusion `success`, headSha `0b5322a…`) |

## Files touched
- `PROGRESS.md` (this file) — worktree root, branch `worker-heavy/t_0509c320`.
- `REPORT.md` — worktree root, same branch.
- No file in the repo under test was modified.
