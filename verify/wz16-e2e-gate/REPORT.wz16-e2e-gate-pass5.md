# REPORT — wz16 e2e integration gate (pass 5, independent re-run)

**Card task:** End-to-end verify: snapshot build, --version smoke, full suite, push. Runs after the
CI workflow, goreleaser config, and docs are committed. Integration gate before cold review.

**Verdict:** ALL 6 STEPS PASS. Pushed, tree clean (0/0 vs remote), CI green at the exact tip.

**Verifier this pass:** worker-heavy / deepseek-v4-flash, host `x280`, run start 2026-10-16T06:27:34+02:00.
**Method:** every command below was re-executed live this pass; raw output captured under
`/home/c08r4d0r/wz16-captures-v2/`. This is an *independent re-run*, not a replay of the pass-4 captures.

---

## 0. Two repo-name corrections (verified, not assumed)

The card's repo/step strings do not resolve as written. Both checked against live GitHub:

| Card says | Reality (verified) | Evidence |
|---|---|---|
| repo `net4sats/tollgate-wizard-go` | **does not exist** — HTTP 404 / GraphQL "Could not resolve to a Repository" | `08` §5b |
| binary prints `tollgate-wizard <version>` | **CORRECT** — pinned by `main.go:55 const programName = "tollgate-wizard"` | `06-version-smoke.txt` |

The real target is `felixfelix-bot/net4sats-wizard-go` (public fork; upstream `net4sats/net4sats-wizard-go`).
Working clone/branch: `/home/c08r4d0r/wz16-verify` @ `ci/add-github-actions-workflow-clean`.
So the card's `--version` expectation is right, but its `--repo` is wrong. Both are stated plainly rather
than silently substituted.

Note on my own first attempt: my initial smoke script asserted the line starts with `net4sats-wizard`
(inferred from the repo name) and therefore printed `FORMAT: FAIL`. That was **my script's bug**, not a
product bug — `main.go:55` pins `tollgate-wizard`. Corrected and re-run; the corrected capture is what
`06-version-smoke.txt` contains. The failed first capture was not retained as evidence of a pass.

---

## 1. Clean tree + push + ahead/behind  → PASS

Capture: `01-clean-tree-and-push.txt`

```
### branch: ci/add-github-actions-workflow-clean
### HEAD: 64008dc2cfa83f9e3679f1f40d21f186c765642d
$ git status --porcelain
PORCELAIN_RC=0  LINES=0
$ git push
Everything up-to-date
PUSH_RC=0
$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0
local_HEAD          = 64008dc2cfa83f9e3679f1f40d21f186c765642d
origin/ci/add-...   = 64008dc2cfa83f9e3679f1f40d21f186c765642d
$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
64008dc2cfa83f9e3679f1f40d21f186c765642d	refs/heads/ci/add-github-actions-workflow-clean
```

`git status --porcelain` empty, push up-to-date, ahead/behind = `0 0`, and remote ground truth
(`ls-remote`, independent of local tracking refs) equals local HEAD.

---

## 2. Full suite  → PASS

Captures: `02-vet-lint-check.txt`, `03-test-race.txt`, `04-test-race-summary-verbose.txt`

```
$ go vet ./...                → VET_RC=0
$ golangci-lint run           → LINT_RC=0        (golangci-lint 1.61.0)
$ goreleaser check            → GRCHECK_RC=0     ("1 configuration file(s) validated")
```

```
$ go test ./... -race
ok  	net4sats-wizard	78.885s
TEST_RACE_RC=0
```

Test-level counts (`go test ./... -race -v`, go1.22.12 linux/amd64):

```
--- PASS = 52    --- FAIL = 0    --- SKIP = 0
ok_pkgs  = 1     DATA_RACE = 0   panic = 0
TEST_RACE_VERBOSE_RC=0
```

Zero failures, zero data races, zero panics. Included and passing is
`TestVersionFlagInjectedViaLdflags`, which itself builds a binary with the same `-X` symbol
`.goreleaser.yaml` uses and fails if the symbol drifts — i.e. the injection path is guarded by the suite.

---

## 3. Snapshot release  → PASS (6 artifacts, 3×2 matrix)

Captures: `05-goreleaser-snapshot.txt`, `05b-artifacts-and-checksums.txt`

```
$ goreleaser release --snapshot --clean
  • building  6 targets: linux_amd64_v1 linux_arm64_v8.0 darwin_amd64_v1 darwin_arm64_v8.0 windows_amd64_v1 windows_arm64_v8.0
  • archiving 6 archives (4 .tar.gz + 2 .zip)
  • calculating checksums
  • release succeeded after 5s
GORELEASER_SNAPSHOT_RC=0
```

Artifact inventory (all `0.5.1-snapshot.64008dc`):

```
net4sats-wizard_..._linux_amd64.tar.gz     4492303
net4sats-wizard_..._linux_arm64.tar.gz     4268919
net4sats-wizard_..._darwin_amd64.tar.gz    4544369
net4sats-wizard_..._darwin_arm64.tar.gz    4413921
net4sats-wizard_..._windows_amd64.zip      4567564
net4sats-wizard_..._windows_arm64.zip      4310153
```

```
ARCHIVE_COUNT=6   (linux=2  windows=2  darwin=2)
BINARY_COUNT=6    (dist/<target>/net4sats-wizard[.exe], one per target)
CHECKSUM_LINES=6  ARCHIVE_COUNT=6
PARITY: OK  (6 checksum lines == 6 archives)
```

Checksum file is `dist/SHA256SUMS` (that is the `name_template` goreleaser emitted; not `checksums.txt`).

```
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.64008dc_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.64008dc_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.64008dc_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.64008dc_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.64008dc_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.64008dc_windows_arm64.zip: OK
SHA256SUM_C_ALL_RC=0     OK_COUNT=6  FAILED_COUNT=0
```

All six verify, not just the two the card minimally required — including one archive
(`linux_amd64.tar.gz`) and one non-linux artifact (`darwin_arm64.tar.gz`, `windows_amd64.zip`).
Cross-platform `Get-FileHash` was not executed (no Windows host here); the equivalent integrity proof is
`sha256sum -c` over the whole `SHA256SUMS`, which is the same SHA-256 comparison `Get-FileHash` performs.

---

## 4. Binary smoke + version injection  → PASS (with controls)

Capture: `06-version-smoke.txt`

```
$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
ELF 64-bit LSB executable, x86-64, statically linked, stripped

$ ./dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
STDOUT: tollgate-wizard 0.5.1-snapshot.64008dc
STDERR: (empty)
VERSION_EXIT_CODE=0
LINE=[tollgate-wizard 0.5.1-snapshot.64008dc]
FORMAT: OK    (prints 'tollgate-wizard <version>')
EXIT_CODE: OK (0)
INJECTION: OK (non-dev version printed: 0.5.1-snapshot.64008dc)
```

Controls proving the injected value is real and not a coincidental constant:

```
$ go build -ldflags "-X main.version=v9.9.9-injection-control" -o /tmp/wz16v2-ctl .
CTL_LINE=[tollgate-wizard v9.9.9-injection-control]
CONTROL: OK  (value CHANGES with the ldflag)

$ go build -o /tmp/wz16v2-nold .            # no ldflags
NOLD_LINE=[tollgate-wizard dev]
NEGATIVE: OK (no injection -> 'dev')
```

The injected snapshot value is `0.5.1-snapshot.64008dc` (not `dev`), exit code 0, and the negative
control proves an uninjected build would have printed `dev`.

---

## 5. GHA proof  → PASS (CI green at the exact tip)

Captures: `07-gh-run-list-and-detail.txt`, `07b-tip-run-detail.txt`, `08-code-tree-identity.txt`

```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38844256545	1m54s	2026-10-16T03:11:19Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38844056356	1m57s	2026-10-16T03:08:56Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843855025	1m58s	2026-10-16T03:06:32Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843632861	1m48s	2026-10-16T03:04:03Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843393341	1m52s	2026-10-16T03:01:14Z
```

The newest run is **at this exact tip** — this is the load-bearing check, because "a green run exists
somewhere on the branch" is not the same claim as "the pushed commit is green":

```
local HEAD SHA  = 64008dc2cfa83f9e3679f1f40d21f186c765642d
CI run headSha  = 64008dc2cfa83f9e3679f1f40d21f186c765642d   (run 38844256545)
SHA_MATCH       = OK  (the green run IS this exact tip)

$ gh run view 38844256545
✓ ci/add-github-actions-workflow-clean CI #29 · 38844256545
JOBS
✓ lint in 15s
✓ test in 1m49s
  steps: go build ./... ✓ | go vet ./... ✓ | go test ./... -race -count=1 ✓
```

CI/goreleaser/docs are committed at this tip:

```
$ git ls-tree -r --name-only HEAD | grep -E 'workflows|\.goreleaser'
.github/workflows/ci.yml
.github/workflows/release.yml
.goreleaser.yaml
```

**CI is green via `gh run list` for this branch.** The card correctly anticipated the repo/privacy
contingency; the failure it flagged (`gh` unauthenticated / private repo) is not what happened — `gh` is
authenticated as `felixfelix-bot` and the repo is public. The actual obstacle was that
`net4sats/tollgate-wizard-go` does not exist (§0). No workflow was weakened or edited this pass;
nothing needed fixing, so nothing was touched.

**Honest caveat (recorded, not hidden):** PR #29
(`consolidation: ci/add-github-actions-workflow-clean -> main`) is OPEN at `headRefOid 64008dc2…` but its
`statusCheckRollup` is **empty** and `mergeable: CONFLICTING`. So CI green is established from the
workflow runs (triggered by `workflow_dispatch`), **not** from PR check-rollup. The green evidence is real;
it just does not appear as a PR check. The PR's conflicting/mergeable state is a consolidation concern for
the manager, outside this card's steps.

**Code identity (why the green CI covers the verified snapshot):** across the whole CI-covered range
(`b71af0b` → `64008dc`) there are **zero** changes to any `*.go`, `.goreleaser.yaml`, or `.github/**`
file — the diff is docs only:

```
$ git diff --name-only b71af0b 64008dc
PROGRESS.md
verify/wz16-e2e-gate/10-ci-rerun-at-final-tip.txt
verify/wz16-e2e-gate/PROGRESS.wz16-e2e-gate-pass4.md
verify/wz16-e2e-gate/REPORT.wz16-e2e-gate-pass4.md
verify/wz16-e2e-gate/card-comment.md

GO_CODE_DIFF_LINES=0  (expect 0)
HEAD^{tree} = 74d3597fd4377a704d903e02bfaf7eb163b61ef6
main.go blob = 30b35e02055417e35adc4695a08283ea9d0ad92c
.goreleaser.yaml blob = 0741e31079072bb26135fb605e689c66233d7434
```

Therefore the CI-green SHA and the SHA the snapshot/`--version` evidence was produced from carry an
identical code tree. The green CI genuinely validates the code behind this evidence.

---

## Acceptance checklist

| Requirement | Status | Evidence |
|---|---|---|
| `git status --porcelain` empty | PASS (0 lines) | `01` |
| branch pushed, ahead/behind 0 | PASS (`0 0`, ls-remote == HEAD) | `01` |
| `go vet ./...` | PASS (rc 0) | `02` |
| `go test ./... -race` | PASS (52 pass / 0 fail / 0 race) | `03`,`04` |
| `golangci-lint run` | PASS (rc 0) | `02` |
| `goreleaser check` | PASS (validated) | `02` |
| snapshot: 6 artifacts + archives + checksums | PASS (6/6, parity OK) | `05`,`05b` |
| checksum verification (≥2, one archive one non-linux) | PASS (all 6 OK) | `05b` |
| binary `--version` = `tollgate-wizard <version>`, exit 0 | PASS | `06` |
| injected value manually changed | PASS (control + negative control) | `06` |
| CI green at pushed tip via `gh run list` | PASS (run 38844256545 @ HEAD) | `07`,`07b` |
| verbatim output for every command | PASS (captures listed) | all |

**Nothing was skipped. Nothing was weakened to pass. All commands actually executed this pass.**

## Raw captures (this pass)

```
/home/c08r4d0r/wz16-captures-v2/01-clean-tree-and-push.txt
/home/c08r4d0r/wz16-captures-v2/02-vet-lint-check.txt
/home/c08r4d0r/wz16-captures-v2/03-test-race.txt
/home/c08r4d0r/wz16-captures-v2/04-test-race-summary.txt
/home/c08r4d0r/wz16-captures-v2/04-test-race-summary-verbose.txt
/home/c08r4d0r/wz16-captures-v2/04-test-race-verbose.log
/home/c08r4d0r/wz16-captures-v2/05-goreleaser-snapshot.txt
/home/c08r4d0r/wz16-captures-v2/05b-artifacts-and-checksums.txt
/home/c08r4d0r/wz16-captures-v2/06-version-smoke.txt
/home/c08r4d0r/wz16-captures-v2/07-gh-run-list-and-detail.txt
/home/c08r4d0r/wz16-captures-v2/07b-tip-run-detail.txt
/home/c08r4d0r/wz16-captures-v2/08-code-tree-identity.txt
```

---

## Pass-5 close-out: CI green at the final pushed tip

Per-tip confirmation loop (each line = a `workflow_dispatch` CI run whose `headSha` equals that tip):

| tip | CI run | conclusion |
|---|---|---|
| `64008dc2` (pass-4 tip, code tree) | 38844256545 | success |
| `95743c3` (pass-5 evidence) | 38851358337 | success |
| `21ba636` | 38851577790 | success |
| `c71b2f8` (**FINAL tip = HEAD = remote**) | 38851802391 | success |

Final state, verbatim:

```
### FINAL tip = HEAD = remote = c71b2f860bb51c951c4cd243586ec6cda8dcc947
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851802391	1m56s	2026-10-16T04:40:12Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851577790	1m59s	2026-10-16T04:37:28Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851358337	1m58s	2026-10-16T04:34:52Z

$ gh run view 38851802391 --json conclusion,headSha,status
{"conclusion":"success","headBranch":"ci/add-github-actions-workflow-clean","headSha":"c71b2f860bb51c951c4cd243586ec6cda8dcc947","status":"completed"}

porcelain=[]
ahead/behind=0	0
```

Capture on disk: `verify/wz16-e2e-gate/pass5/11-ci-at-c71b2f8.txt`.
**All six card steps pass on a pushed, clean, CI-green tip.**
