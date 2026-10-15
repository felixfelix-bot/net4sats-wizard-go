# PROGRESS — t_0509c320: end-to-end verify (snapshot build, --version smoke, full suite, push)

Board `net4sats-mvp-v2` · card `t_0509c320` · assignee worker-heavy · host **x280** (`c08r4d0r`) · 2026-10-15
Worktree (this pass): `/home/c08r4d0r/worktrees/wz-e2e-gate-r2` — branch `worker-heavy/t_0509c320-r2`,
cut at `0b5322a` (unchanged source; docs-only commits on top).
Subject under test: **`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`** =
`origin/ci/add-github-actions-workflow-clean` (checkout `~/wz16-verify`) = head of **PR #35**.
Raw captures: `/home/c08r4d0r/wz-captures-t_0509c320-reverify/` · scripts `/home/c08r4d0r/wz16-tools/rev2-*.sh`

This is pass 2 of the same card (pass 1 = `PROGRESS.t_0509c320.md`). Pass 2 re-runs **every** gate
live on this host rather than re-citing pass-1 captures; the subject SHA is identical.

## Finding 1 — the repo named in the card does not exist (unchanged from pass 1)
- `gh run list --repo net4sats/tollgate-wizard-go --limit 5`
  → `failed to get runs: HTTP 404: Not Found (…/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5…)`
- `gh repo view net4sats/tollgate-wizard-go` → `GraphQL: Could not resolve to a Repository…`
- `git ls-remote https://github.com/net4sats/tollgate-wizard-go` → `remote: Repository not found.` (rc 128)
- `gh` IS authenticated (account `felixfelix-bot`), so this is a real 404, not an auth failure.
- Real repo: **`net4sats/net4sats-wizard-go`**; push target (upstream is READ-only): fork
  **`felixfelix-bot/net4sats-wizard-go`**.
- Status: RESOLVED (use the real repo + fork). Files touched: none.

## Finding 2 — upstream cannot report CI for this PR at all
- `gh workflow list --repo net4sats/net4sats-wizard-go` → empty.
- `gh api repos/net4sats/net4sats-wizard-go/contents/.github` → `{"message":"Not Found",…"status":"404"}`.
- `gh pr checks 35 --repo net4sats/net4sats-wizard-go` → `no checks reported on the 'ci/add-github-actions-workflow-clean' branch`.
- `gh pr view 35 --json statusCheckRollup` → `[]`.
- Cause: `pull_request` runs come from the BASE repo's workflows; upstream `main` has no `.github/`,
  which is exactly what PR #34/#35 add. So the fork is the only venue where this workflow can run.
- Status: RESOLVED — fresh `workflow_dispatch` on the fork at the exact head SHA (no check weakened,
  no workflow edited). Files touched: none.

## Finding 3 — golangci-lint 1.61.0 cannot read go1.24 export data
- Under the host toolchain go1.24.4, `golangci-lint run` exits 3 with
  `buildir: failed to load package goarch: could not load export data: internal error in importing
  "internal/goarch" (unsupported version: 2)`.
- That is a toolchain/export-data mismatch, **not** a code defect. CI pins `go-version-file: go.mod`
  (= go 1.22) + `golangci-lint v1.61.0`, so the CI-matching gate is go1.22.12.
- Fix: run every gate with `/home/c08r4d0r/wz16-tools/go122/bin/go` (go1.22.12) first on PATH.
- Status: RESOLVED. Files touched: none (toolchain shim lives outside the repo).

## Finding 4 — no source change was required; every gate is green at 0b5322a
- `go vet ./...` rc 0 · `gofmt -l .` no output rc 0 · `golangci-lint run` rc 0 · `goreleaser check` rc 0.
- `go test ./... -race -count=1` (CGO_ENABLED=1, gcc 13.3.0) → `ok net4sats-wizard 133.926s` rc 0;
  verbose re-run: 123 `=== RUN`, 0 `--- FAIL`, 0 `--- SKIP`, **0** `WARNING: DATA RACE`.
- `goreleaser release --snapshot --clean` rc 0 → version `0.5.1-snapshot.0b5322a`, 6 archives
  (linux/windows/darwin × amd64/arm64; 4 tar.gz + 2 zip), `SHA256SUMS` 6 lines == 6 artifacts,
  `sha256sum -c SHA256SUMS` 6/6 OK, 3 independent recomputes + 1 openssl recompute MATCH.
- Binary smoke: `--version` → `tollgate-wizard 0.5.1-snapshot.0b5322a` rc 0 (not `dev`);
  two manual `-X main.version=…` builds changed the printed line; a no-ldflags control printed `dev`.
- Tree hygiene: `git status --porcelain` empty after the full suite (gates left no tracked churn;
  `dist/` is gitignored).
- Status: DONE. Files touched: `PROGRESS.md`, `REPORT.md` (worktree root, this task's branch only).

## Gate status (pass 2)

| # | Gate | Result |
|---|------|--------|
| 1 | clean tree + pushed, ahead/behind 0/0 | PASS |
| 2a | `go vet ./...` (go1.22.12) | PASS (rc 0) |
| 2b | `CGO_ENABLED=1 go test ./... -race -count=1` | PASS (123 RUN / 0 FAIL / 0 SKIP / 0 data races) |
| 2c | `golangci-lint run` v1.61.0 | PASS (rc 0) |
| 2d | `goreleaser check` (v2.18.3) | PASS (rc 0, 1 config validated) |
| 3 | `goreleaser release --snapshot --clean` | PASS (6/6 archives, SHA256SUMS 6/6) |
| 3b | `sha256sum -c` + independent recomputes | PASS (6/6 OK; 4/4 recomputed MATCH) |
| 4 | binary `--version` + manual injection | PASS (`tollgate-wizard 0.5.1-snapshot.0b5322a` rc 0) |
| 5 | CI green on the fork at the pushed head SHA | PASS (run recorded in REPORT.md STEP 5) |

## Card-level note
The card sits in status `triage`, which is worker-immutable: `complete` / `promote` / `claim` /
`request-review` all refuse from a worker session, and the latest summary is
`reclaim backoff until 1791730323 (dead worker pid, strike 1)`. Evidence is therefore attached +
commented; **the manager must promote the card** for it to leave `triage`.

