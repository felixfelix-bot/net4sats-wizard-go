# PROGRESS — t_0509c320 pass 4 (rev4): end-to-end verify — snapshot build, --version smoke, full suite, push

Board `net4sats-mvp-v2` · card `t_0509c320` · host **x280** (`c08r4d0r`) · 2026-10-16
Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate-r3` — branch `worker-heavy/t_0509c320-r3`, HEAD `a412b61`.
Evidence captures: `.evidence-rev4/` (worktree root). Full report: `REPORT.md`.
Toolchain: go1.22.12 (`/usr/local/go/bin/go`) · golangci-lint 1.61.0 · goreleaser 2.18.3 · gcc 13.3.0 · gh (acct `felixfelix-bot`).

Pass 4 re-executes every gate live; no pass-1/2/3 capture is re-cited as this pass's evidence.

## Finding 1 — the repo slug in the dispatch does not exist (re-measured)
- `gh run list --repo net4sats/tollgate-wizard-go --limit 5` → `failed to get runs: HTTP 404: Not Found (…)` rc=1.
- `gh api repos/net4sats/tollgate-wizard-go` → `{"message":"Not Found",…,"status":404}`. `gh` is authenticated, so a real 404.
- Real repo `net4sats/net4sats-wizard-go`; push target fork `felixfelix-bot/net4sats-wizard-go`.
- Status: RESOLVED — use the real repo. Files touched: none.

## Finding 2 — Step 5 failed as literally specified; pass 3's "CI green" was a DIFFERENT branch (blocker, now closed)
- `gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3` → 0 runs.
- `gh pr view 36` → `statusCheckRollup: []`, `mergeable: CONFLICTING`, `mergeStateStatus: DIRTY`.
- Not one of the 100 most recent runs had `head_sha` `a412b61` (grep count 0). Every `success` run is
  `event=workflow_dispatch` on branch `ci/add-github-actions-workflow-clean`.
- Root cause: fork PR `pull_request` runs use the BASE repo's workflows; upstream `main` has no `.github/`
  (`gh api …/contents/.github/workflows` → 404). The fork's `ci.yml` triggers `push` only on `main`, and the fork's
  event histogram is `{"workflow_dispatch": 23}` — no push-triggered run has ever existed.
- Resolution (no check weakened, no workflow edited): dispatch the branch's OWN `ci.yml` on the subject branch →
  run **38866254688** `conclusion=success`, `headSha=a412b616…` == pushed HEAD, jobs `test` + `lint` both success.
- Residual (flagged to manager, NOT claimed as done): a merge-blocking required check on fork PR #36 needs the
  workflow upstream (PR #35, open) or a `branches:` addition on the fork.
- Status: DONE at the pushed-head-SHA reading of the AC. Files touched: none.

## Finding 3 — Steps 1–4 green at HEAD a412b61
- Step 1: `git status --porcelain` 0 lines rc 0; HEAD == `@{u}` == a412b61; `0	0`; `git push` → `Everything up-to-date`.
- Step 2: `go vet ./...` rc 0; `gofmt -l .` no output; `CGO_ENABLED=1 go test ./... -race -count=1` → `ok net4sats-wizard`
  (verbose 123 RUN / 52 PASS / 0 FAIL / 0 SKIP / 0 DATA RACE); `golangci-lint run` rc 0; `goreleaser check` rc 0.
- Step 3: `goreleaser release --snapshot --clean` rc 0 → `0.5.1-snapshot.a412b61`; 6 archives (l/w/d × amd64/arm64);
  `SHA256SUMS` 6 lines == 6 artifacts; `sha256sum -c` 6/6 OK rc 0; independent recomputes of linux_amd64.tar.gz and
  windows_arm64.zip both MATCH; negative control (truncated copy) FAIL_AS_EXPECTED (non-vacuous).
- Step 4: `--version` → `tollgate-wizard 0.5.1-snapshot.a412b61` rc 0 (not `dev`); manual `-X main.version=…` builds
  produced `v9.9.9-evidence` and `v0.0.1-changed`; no-ldflags control printed `dev`.
- Step 6: two card comments posted; `REPORT-t_0509c320-rev4.md` attached (id 105).

## Files touched by this pass
- `PROGRESS.md`, `REPORT.md` (worktree root) + `.evidence-rev4/`. No source/test/workflow change.
- Card status is `triage`; this worker cannot transition it — manager must promote.
