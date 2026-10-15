# PROGRESS — t_0509c320 (pass 3): end-to-end verify — snapshot build, --version smoke, full suite, push

Board `net4sats-mvp-v2` · card `t_0509c320` · assignee worker-heavy · host **x280** (`c08r4d0r`) · 2026-10-15
Worktree (this pass): `/home/c08r4d0r/worktrees/wz-e2e-gate-r3` — branch `worker-heavy/t_0509c320-r3`,
cut at `0b5322a` (no source change; docs-only commits on top).
Subject under test: **`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`** =
`origin/ci/add-github-actions-workflow-clean` (checkout `~/wz16-verify`) = head of **PR #35** (OPEN).
Raw captures: `/home/c08r4d0r/wz-captures-t_0509c320-p3/` (`00`…`07`) · scripts `~/wz16-tools/rev3-*.sh`
Toolchain: go1.22.12 (CI-matching; `ci.yml` uses `go-version-file: go.mod` = `go 1.22`), golangci-lint
1.61.0, goreleaser 2.18.3, gcc 13.3.0, gh 2.45.0 (account `felixfelix-bot`).

Pass 3 of the same card (pass 1 = `PROGRESS.t_0509c320.md`, pass 2 = `PROGRESS.t_0509c320-pass2.md` in the
lane repo). Every gate is re-executed live on this host; no pass-1/2 capture is re-cited as this pass's
evidence. The subject SHA is unchanged since pass 1.

## Finding 1 — the repo named in the card does not exist (unchanged, re-measured)
- `gh run list --repo net4sats/tollgate-wizard-go --limit 5`
  → `failed to get runs: HTTP 404: Not Found (…/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5…)`
  (rc 1). `gh repo view net4sats/tollgate-wizard-go` → GraphQL "Could not resolve to a Repository".
  `git ls-remote https://github.com/net4sats/tollgate-wizard-go` → `remote: Repository not found.` (rc 128).
- `gh` IS authenticated (account `felixfelix-bot`), so this is a real 404, not an auth failure.
- Real repo: **`net4sats/net4sats-wizard-go`**; push target (upstream is READ-only): fork
  **`felixfelix-bot/net4sats-wizard-go`**.
- Status: RESOLVED (use the real repo + fork). Files touched: none.

## Finding 2 — upstream cannot report CI for this PR at all
- `gh run list --repo net4sats/net4sats-wizard-go --limit 5` → empty (0 runs).
- `gh workflow list --repo net4sats/net4sats-wizard-go` → empty.
- `gh api repos/net4sats/net4sats-wizard-go/contents/.github` → 404 Not Found.
- `gh pr checks 35 --repo net4sats/net4sats-wizard-go` → `no checks reported on the
  'ci/add-github-actions-workflow-clean' branch`; `statusCheckRollup` → `[]`.
- Cause: `pull_request` runs come from the BASE repo's workflows; upstream `main` has no `.github/`,
  which is exactly what PR #34/#35 add. The fork is the only venue where this workflow can run.
- Status: RESOLVED — fresh `workflow_dispatch` on the fork at the exact head SHA. No check weakened,
  no workflow edited. Files touched: none.

## Finding 3 — golangci-lint 1.61.0 needs the CI-matching toolchain (go1.22.12), not host go1.24
- Under go1.24.x, `golangci-lint run` exits 3 with
  `buildir: failed to load package goarch: could not load export data: internal error in importing
  "internal/goarch" (unsupported version: 2)`.
- Toolchain/export-data mismatch, **not** a code defect. CI pins `go-version-file: go.mod` (go 1.22) +
  golangci-lint v1.61.0, so the CI-matching gate is go1.22.12 (`~/wz16-tools/go122/bin/go`, first on PATH).
- Status: RESOLVED. Files touched: none (toolchain shim lives outside the repo).

## Finding 4 — no source change required; every gate green at 0b5322a (see REPORT.md for raw output)
- `go vet ./...` rc 0 · `gofmt -l .` no output rc 0 · `golangci-lint run` rc 0 · `goreleaser check` rc 0.
- `go test ./... -race -count=1` (CGO_ENABLED=1, gcc 13.3.0) → green; verbose counts in REPORT.md §2b.
- `goreleaser release --snapshot --clean` rc 0 → version `0.5.1-snapshot.0b5322a`, 6 archives
  (linux/windows/darwin × amd64/arm64; 4 tar.gz + 2 zip), `SHA256SUMS` 6 lines == 6 artifacts,
  `sha256sum -c SHA256SUMS` 6/6 OK, independent recomputes MATCH, 2 negative controls FAIL as expected.
- Binary smoke: `--version` → `tollgate-wizard 0.5.1-snapshot.0b5322a` rc 0 (not `dev`); two manual
  `-X main.version=…` builds changed the printed line; a no-ldflags control printed `dev`.
- Tree hygiene: `git status --porcelain` empty after the full suite (`dist/` is gitignored).
- Status: DONE. Files touched: `PROGRESS.md`, `REPORT.md` (this worktree root, this task's branch only).

## Gate status (pass 3) — measured

| # | Gate | Result |
|---|------|--------|
| 1 | clean tree + pushed, ahead/behind 0/0 | PASS — porcelain empty; `git push` = `Everything up-to-date`; `0	0`; remote ref == HEAD |
| 2a | `go vet ./...` (go1.22.12) | PASS — rc 0; `gofmt -l .` no output; single module, so `./...` covers the suite |
| 2b | `CGO_ENABLED=1 go test ./... -race -count=1` | PASS — `ok net4sats-wizard 92.390s`; 123 RUN / 0 FAIL / 0 SKIP / 0 data races |
| 2c | `golangci-lint run` v1.61.0 | PASS — rc 0, no findings |
| 2d | `goreleaser check` (v2.18.3) | PASS — rc 0, 1 config validated |
| 3 | `goreleaser release --snapshot --clean` | PASS — `0.5.1-snapshot.0b5322a`, 6 archives (l/w/d × amd64/arm64) |
| 3b | `sha256sum -c` + independent recomputes + controls | PASS — SHA256SUMS 6 lines == 6 artifacts, 6/6 OK rc 0, 3 recomputes MATCH, 1 positive + 2 negative controls |
| 4 | binary `--version` + manual injection | PASS — `tollgate-wizard 0.5.1-snapshot.0b5322a` rc 0 (not `dev`); two manual values changed it; control = `dev` |
| 5 | CI green on the fork at the pushed head SHA | PASS — run `38761913548` `conclusion=success`, `headSha=0b5322a…`, jobs `test` + `lint` |
| 6 | report + card comment | PASS — comment on the card at 2026-10-15 17:45; attachments 92 (REPORT) and 93 (PROGRESS) |

## Ad-hoc verification of this pass's own harness
`rev3-90-adhoc-verify.py` re-derives every claim in REPORT.md from the capture files: **62 checks /
0 failures / `ADHOC_VERDICT=PASS`**. Its control (`rev3-91-adhoc-control.sh`) tampers a COPY of the
captures (flip a gate rc, rename a MATCH, flip the corrupted-artifact control, change the run conclusion,
flip the card-refusal rc) and the same verifier then reports **5 failures** — so the checks are not
vacuous. `bash -n` is clean on every `rev3-*.sh`. This is harness/report verification, NOT the repo suite;
the repo suite is gate 2b. Capture: `11-adhoc-verifier-run.txt`.

## Self-corrections made during this pass (both caught by the verifier above)
1. The first checksum negative control was **vacuous**: `dd of=<glob>` did not expand the glob (it wrote
   a literal `net4sats-wizard_*_linux_amd64.tar.gz` file), so the "corrupted" copy was intact and the
   control printed OK. Fixed with the explicit path; the flawed attempt is preserved as
   `05b-FIRST-ATTEMPT-broken-controls.txt` and the corrected controls are `08-checksum-controls.txt`.
   The checksum run itself was never affected.
2. Two status labels reported a **pipe's** exit code instead of the command's (`sha256sum -c | grep`,
   `gh repo view | head`). Fixed by capturing the real codes (`SHA256SUM_C_REAL_RC=0`, `REPOVIEW_RC=1`,
   `LSREMOTE_RC=128`); corrected values are in `09-gh-run-detail-and-real-rc.txt`.

## Files touched by this pass
- Worktree (branch `worker-heavy/t_0509c320-r3`, pushed to the fork): `PROGRESS.md`, `REPORT.md` —
  commits `141d57f` then `ae49ed0`; remote read-back `0 0`.
- Lane repo `main`: `PROGRESS.t_0509c320-pass3.md`, `REPORT.t_0509c320-pass3.md`,
  `verify/t_0509c320-pass3/` (35 files) — commits `9acadc59` then `d2f6159e`, both in `origin/main`.
- **No source, test, or workflow file was changed**: the branch is green at `0b5322a` as received.

## Card-level note
The card sits in status `triage`, so this worker session cannot close it — both transitions refused
verbatim (capture `10-card-transition-attempts.txt`):
`hermes kanban complete t_0509c320` → `cannot complete t_0509c320 (unknown id or terminal state)`;
`hermes kanban promote t_0509c320 --force` → `cannot promote t_0509c320: task t_0509c320 is 'triage';
promote only applies to 'todo' or 'blocked'`. Evidence is attached (92, 93) and commented, so **the
manager must promote the card** out of `triage`. Every acceptance criterion in the card is green on the
pushed SHA; there is no residual worker action.

