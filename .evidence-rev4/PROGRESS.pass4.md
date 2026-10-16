# PROGRESS — t_0509c320 (pass 4 / rev4): end-to-end verify — snapshot build, --version smoke, full suite, push

Board `net4sats-mvp-v2` · card `t_0509c320` · host **x280** (`c08r4d0r`) · 2026-10-16
Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate-r3` — branch `worker-heavy/t_0509c320-r3`, HEAD `a412b61`.
Evidence captures: `.evidence-rev4/` (this worktree root).
Toolchain: go1.22.12 (`/usr/local/go/bin/go`) · golangci-lint 1.61.0 · goreleaser 2.18.3 · gcc 13.3.0 · gh (acct `felixfelix-bot`).

Pass 4 re-executes EVERY gate live; no pass-1/2/3 capture is re-cited as this pass's evidence.

## Finding 1 — the repo slug in the dispatch does not exist (re-measured this pass)
- `gh run list --repo net4sats/tollgate-wizard-go --limit 5`
  → `failed to get runs: HTTP 404: Not Found (…/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5…)` rc=1.
- `gh api repos/net4sats/tollgate-wizard-go` → `{"message":"Not Found",…,"status":404}` / `gh: Not Found (HTTP 404)`.
- `gh` IS authenticated (`felixfelix-bot`), so this is a real 404, not an auth failure.
- Real repo: `net4sats/net4sats-wizard-go` (upstream, read-only) → fork `felixfelix-bot/net4sats-wizard-go` (push target).
- Status: RESOLVED — use the real repo. Capture: `.evidence-rev4/05-gha-proof.txt`.

## Finding 2 — Step 5 FAILS as literally specified, and pass 3's green was on a DIFFERENT branch (blocker)
Measured this pass on the pushed branch:
- `gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3 --limit 10` → **empty (0 runs)** rc=0.
- `gh pr view 36 --repo felixfelix-bot/net4sats-wizard-go --json …` →
  `{"headRefOid":"a412b61…","mergeable":"CONFLICTING","mergeStateStatus":"DIRTY","statusCheckRollup":[]}` → **no checks**.
- Sweep of all 100 most recent runs: **not one** has `head_sha` starting `a412b61` (grep count = 0).
- Every `success` run belongs to branch **`ci/add-github-actions-workflow-clean`** and every one is
  `event=workflow_dispatch` — i.e. the pass-3 green at `0b5322a` was a MANUAL dispatch on
  `~/wz16-verify`, not on this branch's pushed head.
- Root cause (matches pass-3 Finding 2): a fork PR's `pull_request` runs execute the **base repo's**
  workflows. Upstream `net4sats/net4sats-wizard-go` `main` has NO `.github/` (`gh api …/contents/.github` → 404),
  so PR #36 gets an empty check set. The branch DOES carry `.github/workflows/{ci,release}.yml` at its head.
- Also: PR #36 is `CONFLICTING`/`DIRTY` (branch-only 12, main-only 16 since the `c747aae` merge-base), so it
  cannot be merged as-is to inherit anything.
- Status: **OPEN / BLOCKING for the literal AC.** No check was weakened and no workflow file was edited.
  Capture: `.evidence-rev4/05-gha-proof.txt`.

## Finding 3 — Steps 1–4 all green at HEAD a412b61 (raw output in .evidence-rev4/)
- Step 1 clean tree: `git status --porcelain` empty (0 lines) rc=0; `git rev-parse HEAD` == `@{u}` == `a412b61…`;
  `git rev-list --left-right --count @{u}...HEAD` → `0	0`; `git push` → `Everything up-to-date` rc=0.
- Step 2a `go vet ./...` rc=0; `gofmt -l .` no output rc=0.
- Step 2b `CGO_ENABLED=1 go test ./... -race -count=1` → `ok net4sats-wizard 62.623s` rc=0.
- Step 2c `golangci-lint run` rc=0 (no findings). Step 2d `goreleaser check` rc=0 (1 config validated).
- Step 3 `goreleaser release --snapshot --clean` rc=0 → `0.5.1-snapshot.a412b61`; 6 archives
  (linux/windows/darwin × amd64/arm64 = 4 tar.gz + 2 zip); `SHA256SUMS` 6 lines == 6 artifacts;
  `sha256sum -c SHA256SUMS` 6/6 OK rc=0; independent recomputes of one archive (linux_amd64.tar.gz) and one
  non-linux artifact (windows_arm64.zip) both MATCH; negative control (truncated copy) FAILS as expected
  (non-vacuous).
- Step 4 `dist/…_linux_amd64_v1/net4sats-wizard --version` → `tollgate-wizard 0.5.1-snapshot.a412b61` rc=0
  (NOT `dev`); manual `-X main.version=v9.9.9-evidence` → `tollgate-wizard v9.9.9-evidence`;
  `-X main.version=v0.0.1-changed` → `tollgate-wizard v0.0.1-changed`; no-ldflags control → `tollgate-wizard dev`.
- Status: PASS. Files touched: none (dist/ is gitignored; no source/workflow edit).

## Files touched by this pass
- `PROGRESS.md`, `REPORT.md` (this worktree root) + `.evidence-rev4/` captures. No source/test/workflow change.
