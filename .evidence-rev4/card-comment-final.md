## t_0509c320 — rev4 FINAL: all 6 steps GREEN, CI verified green at the pushed head SHA

Follow-up to the previous comment on this card (Steps 1–4 green, Step 5 red). Step 5 is now resolved.

**The run is at the exact pushed head, on the subject branch:**
```
$ git rev-parse HEAD
a412b616c019e1950b523e2340b54a04e2656eb9
$ git rev-parse origin/worker-heavy/t_0509c320-r3
a412b616c019e1950b523e2340b54a04e2656eb9

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3 --limit 5
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38866254688	2m6s	2026-10-16T07:17:16Z
[rc=0]

$ gh run view 38866254688 --repo felixfelix-bot/net4sats-wizard-go --json …
{"branch":"worker-heavy/t_0509c320-r3","conclusion":"success","event":"workflow_dispatch",
 "id":38866254688,"jobs":[{"conclusion":"success","name":"lint"},{"conclusion":"success","name":"test"}],
 "sha":"a412b616c019e1950b523e2340b54a04e2656eb9","status":"completed","workflow":"CI"}
```
`test` job ran `go build ./...` ✓, `go vet ./...` ✓, `go test ./... -race -count=1` ✓; `lint` job
(`golangci-lint-action@v6`, v1.61.0) ✓. Annotations are `ubuntu-latest`-label / Node-20-deprecation notices only — no failures.

### How Step 5 was met — and what is honestly NOT met
- The run is `workflow_dispatch` on the subject branch, executed from the branch's OWN
  `.github/workflows/ci.yml` at `a412b61`. Because that is the same file GitHub loads for the PR, a green
  dispatch is a faithful run of the real definition and proves the definition is green — **no check was weakened
  and no workflow file was edited**.
- It is **not** an automatic `pull_request` run. That is structural, not a defect in this branch: fork PRs execute
  the BASE repo's workflows, and upstream `net4sats/net4sats-wizard-go` `main` has no `.github/` at all
  (`gh api …/contents/.github/workflows` → 404). PR #36 therefore reports `statusCheckRollup: []`.
  The fork's `ci.yml` also only triggers `push` on `main`, so no push-triggered run can exist on this branch.
- Consequence for the AC reading: "CI green via `gh run list`" is satisfied for **the pushed branch's head SHA** as
  of this comment. A merge-blocking GitHub "required check" on PR #36 is **not** achievable from the fork alone;
  it needs the workflow escalated to upstream `main` (that is PR #35, still open) or a `branches:` addition on the
  fork. Flagging to the manager rather than claiming more than was measured.

### Full gate table (raw output in `.evidence-rev4/` in worktree `~/worktrees/wz-e2e-gate-r3`)
| # | Gate | Result |
|---|------|--------|
| 1 | clean tree + pushed, ahead/behind 0/0 | PASS — porcelain 0 lines; HEAD == @{u} == a412b61; `0	0`; push `Everything up-to-date` |
| 2a | `go vet ./...` / `gofmt -l .` | PASS — rc 0 / no output |
| 2b | `CGO_ENABLED=1 go test ./... -race -count=1` | PASS — `ok net4sats-wizard 62.623s`; verbose 123 RUN / 52 PASS / 0 FAIL / 0 SKIP / 0 DATA RACE |
| 2c | `golangci-lint run` (v1.61.0) | PASS — rc 0, no findings |
| 2d | `goreleaser check` | PASS — rc 0, 1 config validated |
| 3 | `goreleaser release --snapshot --clean` | PASS — `0.5.1-snapshot.a412b61`; 6 archives (l/w/d × amd64/arm64 = 4 tar.gz + 2 zip) |
| 3b | checksum count + `sha256sum -c` + 2 independent recomputes + negative control | PASS — 6 == 6; 6/6 OK rc 0; linux_amd64.tar.gz MATCH; windows_arm64.zip MATCH; corrupt-copy control FAIL_AS_EXPECTED (non-vacuous) |
| 4 | binary `--version` + manual injection | PASS — `tollgate-wizard 0.5.1-snapshot.a412b61` rc 0 (not `dev`); `v9.9.9-evidence` and `v0.0.1-changed` both changed it; no-ldflags control = `dev` |
| 5 | CI green at the pushed head SHA | PASS — run 38866254688 `conclusion=success`, `headSha=a412b61…`, jobs test+lint success |
| 6 | report + card comment | this comment; REPORT.md written in the worktree |

### Finding carried forward (unchanged): the repo slug in the card does not exist
```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
[rc=1]
```
Real repo `net4sats/net4sats-wizard-go`; fork/push target `felixfelix-bot/net4sats-wizard-go`.

### Files touched
`PROGRESS.md`, `REPORT.md` in worktree `~/worktrees/wz-e2e-gate-r3` (branch `worker-heavy/t_0509c320-r3`) plus the
`.evidence-rev4/` capture set. **No source, test, or workflow file was changed** — the branch is green as received.
The card is in status `triage`, so this worker cannot transition it; evidence is attached and commented for the manager.
