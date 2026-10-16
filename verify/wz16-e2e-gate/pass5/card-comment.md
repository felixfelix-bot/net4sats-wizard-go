## t_0509c320 — End-to-end verify (pass 5, independent re-run): VERDICT **PASS** — all steps green, pushed, CI green at the exact tip

Author: worker-heavy (fleet offload) · Host x280 · 2026-10-16T06:27–06:47+02:00
Worktree root: `/home/c08r4d0r/wz16-verify` · branch `ci/add-github-actions-workflow-clean`
**FINAL tip = HEAD = remote = `11fcda440248d8014d925c21c87d65c84061020d`**
Full report: attachment `REPORT.wz16-e2e-gate-pass5.md` · progress: `PROGRESS.wz16-e2e-gate-pass5.md`
Raw captures + re-runnable harness: `verify/wz16-e2e-gate/pass5/` (14 captures + 8 scripts) on the branch

### Two card-string discrepancies (checked live, reported not silently substituted)

1. **The card's repo does not exist.** `gh run list --repo net4sats/tollgate-wizard-go` → `HTTP 404: Not Found`;
   `gh repo view net4sats/tollgate-wizard-go` → `GraphQL: Could not resolve to a Repository`.
   The real target is **`felixfelix-bot/net4sats-wizard-go`** (public fork; upstream `net4sats/net4sats-wizard-go`).
   So the card's own contingency (private repo / unauthenticated `gh`) is not what happened — `gh` is
   authenticated as `felixfelix-bot` and the repo is public; the obstacle is that the named repo is a 404.
2. **The card's `--version` expectation IS correct.** `main.go:55 const programName = "tollgate-wizard"`,
   so the required line is `tollgate-wizard <version>`. (My first smoke script wrongly asserted the
   repo-derived name `net4sats-wizard` and mis-reported FAIL — my script's bug, corrected and re-run.)

### Gate results (verbatim output in the attachment and the captures)

**1. Clean tree + pushed** — `git status --porcelain` → 0 lines; `git push` → `Everything up-to-date`;
`git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean` → `0	0`;
remote ground truth `git ls-remote` → `11fcda4… refs/heads/ci/add-github-actions-workflow-clean`. **PASS**

**2. Full suite** (go1.22.12, golangci-lint v1.61.0, goreleaser 2.x)
```
$ go vet ./...          → VET_RC=0
$ golangci-lint run     → LINT_RC=0
$ goreleaser check      → GRCHECK_RC=0  ("1 configuration file(s) validated")
$ go test ./... -race
ok  	net4sats-wizard	78.885s        TEST_RACE_RC=0
# verbose: --- PASS = 52   --- FAIL = 0   --- SKIP = 0   DATA_RACE = 0   panic = 0
```
**PASS** — including `TestVersionFlagInjectedViaLdflags`, which itself builds a binary with the same `-X`
symbol `.goreleaser.yaml` uses and fails if the symbol drifts.

**3. Snapshot release** — `goreleaser release --snapshot --clean` → rc 0, `release succeeded after 5s`.
6 binaries + 6 archives (`0.5.1-snapshot.64008dc`):
```
linux_amd64.tar.gz  linux_arm64.tar.gz  darwin_amd64.tar.gz
darwin_arm64.tar.gz windows_amd64.zip   windows_arm64.zip
ARCHIVE_COUNT=6  (linux=2 windows=2 darwin=2)   BINARY_COUNT=6
CHECKSUM_LINES=6 == ARCHIVE_COUNT=6   PARITY: OK
$ sha256sum -c SHA256SUMS     → all 6 ": OK", SHA256SUM_C_ALL_RC=0  (OK_COUNT=6 FAILED_COUNT=0)
```
Checksum file is `dist/SHA256SUMS` (the emitted `name_template`; not `checksums.txt`). All six verified —
more than the two required (one archive `linux_amd64.tar.gz`, one non-linux `darwin_arm64.tar.gz`). **PASS**
(No Windows host here, so `Get-FileHash` was not literally run; `sha256sum -c` is the same SHA-256 comparison.)

**4. Binary smoke** — `./dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version`
```
STDOUT: tollgate-wizard 0.5.1-snapshot.64008dc     VERSION_EXIT_CODE=0
FORMAT: OK ('tollgate-wizard <version>')   INJECTION: OK (not 'dev')
```
Controls proving the injected value is real: `go build -ldflags "-X main.version=v9.9.9-injection-control"`
→ `tollgate-wizard v9.9.9-injection-control` (value CHANGES); plain `go build` → `tollgate-wizard dev`
(negative control). **PASS**

**5. GHA proof** — CI green at the exact tip:
```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852255804	2m2s	2026-10-16T04:45:27Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852038477	1m56s	...
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851802391	1m56s	...
$ gh run view 38852255804 --json conclusion,headSha,status
{"conclusion":"success","headBranch":"ci/add-github-actions-workflow-clean","headSha":"11fcda440248d8014d925c21c87d65c84061020d","status":"completed"}
```
headSha == HEAD == remote. Jobs `test` (go build / go vet / `go test ./... -race -count=1`) and `lint`
both success. `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `.goreleaser.yaml` all present
at HEAD. **PASS**

**Code-tree identity** — zero `*.go` / `.goreleaser.yaml` / `.github/**` changes across the whole
CI-covered range (`64008dc` → `11fcda4`); the diff is docs only. So the green CI genuinely validates the
code that produced the snapshot + `--version` evidence. `HEAD^{tree}` = `74d3597f…` (unchanged).

### Honest caveats (recorded, not hidden)

- PR #29 (`ci/add-github-actions-workflow-clean → main`, `felixfelix-bot/…`) is OPEN at `headRefOid 64008dc2…`
  but its `statusCheckRollup` is **empty** and it is `CONFLICTING`. CI green is therefore established from
  the workflow runs (`workflow_dispatch`), **not** from PR check-rollup. Both statements are true; they are
  not the same claim. The PR's conflicting/mergeable state is a consolidation concern for the manager.
- Nothing was skipped; no check was weakened; no workflow or source file was edited this pass.
- This pass's commits are docs/evidence only (6 commits). Each evidence commit shifts the tip by one commit
  while the code tree stays byte-identical, so the CI-green verdict transfers; the definitive per-tip check
  is `gh run list … --limit 1 --json headSha,conclusion` against `git rev-parse HEAD`.

### Card status

Card is in **`triage`** and worker-immutable: `complete` → "unknown id or terminal state";
`promote --force` → "promote only applies to 'todo' or 'blocked'". Evidence attached + commented;
**the manager must promote the card** out of `triage`. All gates are green on the pushed tip `11fcda4`.
