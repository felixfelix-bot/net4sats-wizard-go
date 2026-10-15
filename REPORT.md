# REPORT — t_0509c320: end-to-end verify (snapshot build, `--version` smoke, full suite, push)

Board `net4sats-mvp-v2` · card `t_0509c320` · assignee worker-heavy · host **x280** (`c08r4d0r`) · 2026-10-15
Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate-r2` — branch `worker-heavy/t_0509c320-r2`, cut at `0b5322a`
Subject under test: **`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`** =
`feat(release): link-time version injection, --version, goreleaser config and release docs`
= `origin/ci/add-github-actions-workflow-clean` (checkout `~/wz16-verify`) = **head of PR #35**
Toolchain: **go1.22.12** (CI-matching — `ci.yml` uses `go-version-file: go.mod` = `go 1.22`),
golangci-lint **1.61.0**, goreleaser **2.18.3**, gcc 13.3.0, gh 2.45.0 (account `felixfelix-bot`)
Raw captures: **`/home/c08r4d0r/wz-captures-t_0509c320-reverify/`** (`00`…`09`) · scripts `~/wz16-tools/rev2-*.sh`

**VERDICT: PASS** — every mandated gate re-executed live, 6/6 artifacts verified, `--version` prints the
injected version, and the CI workflow is green on the fork **at the pushed head SHA**.

Three card premises are wrong on this host and are disclosed rather than silently worked around:
(1) the repo name `net4sats/tollgate-wizard-go` does not exist — it is `net4sats/net4sats-wizard-go`;
(2) the mandated `gh run list --repo <the card's repo>` returns **HTTP 404**, and even the real upstream
repo has **no workflows at all** (its `main` lacks `.github/`), so green CI can only be read on the fork;
(3) `golangci-lint run` fails under go1.24 with `unsupported version: 2` — a toolchain mismatch, not a
code defect; re-run under the CI-matching go1.22.12 it is clean.

---

## STEP 1 — clean tree check + push, ahead/behind 0/0

Capture: `01-clean-tree-and-push.txt` (run in `~/wz16-verify`, the checkout that holds the branch).

```
### host: x280 c08r4d0r 2026-10-15T14:45:34+02:00
### repo: /home/c08r4d0r/wz16-verify
### remote: origin  https://github.com/felixfelix-bot/net4sats-wizard-go.git
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff  0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
### branch: ci/add-github-actions-workflow-clean

$ git status --porcelain
porcelain_rc=0   # no output = clean

$ git status -sb
## ci/add-github-actions-workflow-clean...origin/ci/add-github-actions-workflow-clean

$ git push
Everything up-to-date
push_rc=0

$ git fetch origin

$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0

$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff	refs/heads/ci/add-github-actions-workflow-clean

$ git diff --stat origin/ci/add-github-actions-workflow-clean..HEAD
(empty = byte-identical to the pushed head)
```

Clean tree, `git push` = `Everything up-to-date`, ahead/behind **0/0**, remote ref == local HEAD.
**PASS.**

The gate worktree was created from that exact sha (`00-worktree-add.txt`):

```
$ git worktree add -b worker-heavy/t_0509c320-r2 /home/c08r4d0r/worktrees/wz-e2e-gate-r2 0b5322a
Preparing worktree (new branch 'worker-heavy/t_0509c320-r2')
HEAD is now at 0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
# one registration only (guards the two-registrations-one-path trap):
$ git worktree list | grep -c wz-e2e-gate-r2
1
```

---

## STEP 2 — full suite

### 2a — `go vet ./...`, `gofmt`, module shape

Capture: `02-vet-gofmt-lint.txt`.

```
### go: go version go1.22.12 linux/amd64   GOROOT=/home/c08r4d0r/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.12.linux-amd64
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff

$ find . -name go.mod -not -path "*/.git/*" | wc -l
1
$ go list ./...
net4sats-wizard
list_rc=0

$ go vet ./...
VET_RC=0

$ gofmt -l .
GOFMT_RC=0   # no filenames above = no formatting drift

$ git status --porcelain      # did vet/gofmt dirty the tree?
(empty)
```

**A single module** (`go list ./...` = one package), so the mandated `./...` gates cover the whole test
surface — no nested-module vacuity here. **PASS.**

### 2b — `go test ./... -race -count=1`

Captures: `03-test-race-and-goreleaser-check.txt` (named run) and `04-test-race-summary.txt` +
`04-test-race-verbose.log` (verbose counts).

```
### go: go version go1.22.12 linux/amd64  |  cpus: 8  load: 3.97 2.75 2.46
### gcc: gcc (Ubuntu 13.3.0-6ubuntu2~24.04.1) 13.3.0

$ CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	133.926s
TEST_RC=0

$ CGO_ENABLED=1 go test ./... -race -count=1 -v          # verbose re-run for counts
TEST_V_RC=0
### counts from 04-test-race-verbose.log
RUN  (=== RUN)      : 123
PASS (--- PASS)     : 52          # top-level test funcs; the rest of the 123 are subtests
FAIL (--- FAIL)     : 0
SKIP (--- SKIP)     : 0
data races          : 0
package ok lines    : 1  ->  ok  	net4sats-wizard	126.621s
package FAIL lines  : 0
### tail -3
--- PASS: TestVersionFlagInjectedViaLdflags (0.71s)
PASS
ok  	net4sats-wizard	126.621s
```

Two independent full runs, both rc 0, **123 test entries, 0 FAIL, 0 SKIP, 0 `WARNING: DATA RACE`**.
`-race` needs cgo, so `CGO_ENABLED=1` (gcc present) is part of the command that was run. **PASS.**

### 2c — `golangci-lint run` (v1.61.0, exactly what `ci.yml` pins)

```
$ golangci-lint --version
golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z
$ golangci-lint run
LINT_RC=0
```

`.golangci.yml` is v1-schema (matching `golangci-lint-action@v6` pinned to v1.61.0) and enables exactly
the v1 default linter set with no `exclude-rules`. Note the toolchain dependency: under the host's
go1.24.4 this same command exits **3** with
`buildir: failed to load package goarch: … internal error in importing "internal/goarch" (unsupported
version: 2)` — an export-data/toolchain mismatch, not a finding. Under the CI-matching go1.22.12 it is
rc 0 with no output. **PASS.**

### 2d — `goreleaser check`

```
$ goreleaser --version
GitVersion:    2.18.3
Platform:      linux/amd64
$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
GRCHECK_RC=0
```

goreleaser 2.18.3 matches the `goreleaser-action@v6` / `version: v2.18.3` pin in
`.github/workflows/release.yml`, and the config is v2-schema. **PASS.**

---

## STEP 3 — snapshot release: 6 artifacts, checksum file, verification

Capture: `05-snapshot-release.txt` (full log) and `06-checksum-verification.txt`.

```
$ git describe --tags --dirty
v0.5.0-21-g0b5322a
$ rm -rf dist && goreleaser release --snapshot --clean
  • getting and validating git state
    • using tags                                     previous=v0.3.5 current=v0.5.0
  • snapshotting
    • building snapshot...                           version=0.5.1-snapshot.0b5322a
  • running before hooks
    • running                                        hook=go mod verify
  • building binaries
    • building   binaries=net4sats-wizard target=darwin_amd64_v1
    • building   binaries=net4sats-wizard target=windows_arm64_v8.0
    • building   binaries=net4sats-wizard target=linux_amd64_v1
    • building   binaries=net4sats-wizard target=darwin_arm64_v8.0
    • building   binaries=net4sats-wizard target=linux_arm64_v8.0
    • building   binaries=net4sats-wizard target=windows_amd64_v1
  • calculating checksums
  • release succeeded after 8s
GORELEASER_RC=0

### ARTIFACT CENSUS
archives (tar.gz + zip)      : 6
  tar.gz                     : 4
  zip                        : 2
checksum file                : 1
SHA256SUMS lines             : 6

$ cat dist/SHA256SUMS
b1f1ba21751c2b139d2f89b188b3755ad94c4ac52527876a0e30096fda62c0eb  net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
0889b70653189ef79594e27fae7c827f9d3934f7d25801b336e229c9c02834c6  net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
53c715551a4fb957db25aff15e91c9ce0a73344eaf692a95559cc0c435e68432  net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
5319594b92c158ffdb189cd08e79aabfe33d9373df3c77ce64ba011e3539d9c4  net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
f4a6bf0aedad376d7604fdee9430f5fb79674fd58578c9e0f837ad49b8dff1d9  net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
5a0074510c72952367b55f7bbe59e2fa9599e5d7e64f804cc3078c4e38044e8c  net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
```

The 6 archives are exactly linux/windows/darwin × amd64/arm64 (`.tar.gz`, with `.zip` overridden for
windows per `format_overrides`), and the checksum file is named `SHA256SUMS` — the name README documents
for `sha256sum -c`. **checksum count 6 == artifact count 6.**

Checksum verification (`06-checksum-verification.txt`):

```
$ sha256sum -c SHA256SUMS          # run from dist/
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
SHA256SUM_C_RC=0

### independent recompute (NOT the checksum file), one linux archive + non-linux both:
$ sha256sum dist/…_linux_amd64.tar.gz
53c715551a4fb957db25aff15e91c9ce0a73344eaf692a95559cc0c435e68432  dist/…_linux_amd64.tar.gz
  => MATCH
$ sha256sum dist/…_windows_amd64.zip
f4a6bf0aedad376d7604fdee9430f5fb79674fd58578c9e0f837ad49b8dff1d9  dist/…_windows_amd64.zip
  => MATCH
$ sha256sum dist/…_darwin_arm64.tar.gz
0889b70653189ef79594e27fae7c827f9d3934f7d25801b336e229c9c02834c6  dist/…_darwin_arm64.tar.gz
  => MATCH

### PowerShell-style equivalent (Get-FileHash -Algorithm SHA256) reproduced with openssl:
$ openssl dgst -sha256 dist/…_windows_amd64.zip
SHA2-256(dist/…_windows_amd64.zip)= f4a6bf0aedad376d7604fdee9430f5fb79674fd58578c9e0f837ad49b8dff1d9

### archive contents sanity (each archive carries the binary)
$ tar tzf dist/…_linux_amd64.tar.gz
CHANGELOG.md
README.md
net4sats-wizard
$ unzip -l dist/…_windows_amd64.zip
     1575  2026-10-15 14:45   CHANGELOG.md
     8264  2026-10-15 14:45   README.md
 11462144  2026-10-15 14:51   net4sats-wizard.exe
---------                     -------
 11471983                     3 files
```

`sha256sum -c` = **6/6 OK**; three independent recomputes (linux_amd64 archive, windows_amd64 zip,
darwin_arm64 archive) + one `openssl dgst` recompute all MATCH the same digests. Requirement "at least
two artifacts, one archive and one non-linux" is exceeded (3 recomputes, both non-linux targets).
**PASS.**

---

## STEP 4 — binary smoke: `--version` prints the injected version

Capture: `07-version-smoke.txt`.

```
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
### go: go version go1.22.12 linux/amd64

### A. the linux/amd64 binary goreleaser built (dist/net4sats-wizard_linux_amd64_v1/)
$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
… ELF 64-bit LSB executable, x86-64, statically linked, Go BuildID=…, stripped
$ …/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.0b5322a
VERSION_RC=0
$ …/net4sats-wizard -version
tollgate-wizard 0.5.1-snapshot.0b5322a
SHORT_FLAG_RC=0

### B. the SAME binary extracted fresh from the linux archive
$ ~/wz-captures-t_0509c320-reverify/scratch/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.0b5322a
ARCHIVE_BINARY_RC=0

### C. manual ldflags injection — the printed value must CHANGE with the injected one
$ go build -ldflags "-X main.version=v9.9.9-manual-smoke" -o /tmp/wz-manual-1 .
BUILD1_RC=0
$ /tmp/wz-manual-1 --version
tollgate-wizard v9.9.9-manual-smoke
MANUAL1_RC=0

$ go build -ldflags "-X main.version=0.0.0-second-different-value" -o /tmp/wz-manual-2 .
BUILD2_RC=0
$ /tmp/wz-manual-2 --version
tollgate-wizard 0.0.0-second-different-value
MANUAL2_RC=0

### D. control: a build with NO ldflags must print the source default
$ go build -o /tmp/wz-nodev .
BUILD3_RC=0
$ /tmp/wz-nodev --version
tollgate-wizard dev
NODEV_RC=0

### E. source of the version var / flag
$ grep -n "version" main.go | head -20
44:// version is the wizard build version. Release builds inject it at link time:
46://	go build -ldflags "-X main.version=v0.4.0" .
52:var version = "dev"
54:// programName is the CLI identity reported by --version.
57:// versionString is the single line --version prints: the program name followed
58:// by the (possibly injected) build version.
59:func versionString() string {
60:	return programName + " " + version
637:// --version is answered before the server binds a port, so release
641:		case "--version", "-version":
642:			fmt.Println(versionString())

### F. exit code hygiene: --version returns before opening any listener
$ timeout 5 $B --version >/dev/null 2>&1; echo RC=$?
TIMEOUT_GUARD_RC=0
```

Format is `tollgate-wizard <version>` (program name `tollgate-wizard`, not the goreleaser
`project_name`), rc **0**, and the value is the injected `0.5.1-snapshot.0b5322a` — **not `dev`**.
Two different manual injections changed the printed line, and a no-`-ldflags` control printed `dev`,
proving the printed value tracks the injected one and is not hardcoded. `main.go:641` answers
`--version`/`-version` before any listener binds. **PASS.**

---

## STEP 5 — GHA proof

### 5a — the card's own command, and why its premise fails

Captures: `08-gha-proof.txt`, `09-gha-fresh-run.txt`.

```
### gh auth: account felixfelix-bot (gh 2.45.0, token scopes include 'repo','workflow')

=== (1) THE COMMAND THE CARD MANDATES
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
RC=1
$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)
RC=1
$ git ls-remote https://github.com/net4sats/tollgate-wizard-go
remote: Repository not found.
fatal: repository 'https://github.com/net4sats/tollgate-wizard-go/' not found
RC=128
```

`gh` is authenticated, so this is a genuine 404 — the repo in the card does not exist, and the
mandated command therefore cannot show any run. Reported as the exact error, not as success.

### 5b — the real repo has no CI lane for this PR either

```
=== (2) the REAL repo: net4sats/net4sats-wizard-go
$ gh run list --repo net4sats/net4sats-wizard-go --limit 5
(no output)  RC=0
$ gh workflow list --repo net4sats/net4sats-wizard-go
(no output)  RC=0
$ gh api repos/net4sats/net4sats-wizard-go/contents/.github
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":"404"}
$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch
$ gh pr view 35 --json number,state,headRefOid,baseRefName,statusCheckRollup
{"baseRefName":"main","headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","number":35,"state":"OPEN","statusCheckRollup":[]}
```

Upstream `main` has **no `.github/` at all** — which is exactly the directory PR #34/#35 add. A
`pull_request` run is taken from the base repo's workflows, so this PR can never report a check until
that workflow lands upstream. That is a repo-state fact, not "CI pending". No upstream change was made.

### 5c — CI green at the pushed head SHA (fork, where the workflow CAN run)

Capture: `09-gha-fresh-run.txt`. A **fresh** `workflow_dispatch` was triggered at the pushed head
(not a reused older run), then polled to conclusion:

```
$ gh workflow list --repo felixfelix-bot/net4sats-wizard-go
CI	active	376283443

$ gh workflow run CI --repo felixfelix-bot/net4sats-wizard-go --ref ci/add-github-actions-workflow-clean
DISPATCH_RC=0
RUN_ID=38738030458
poll 1: in_progress
…
poll 9: completed success

$ gh run view 38738030458 --json databaseId,name,status,conclusion,headSha,headBranch,event,createdAt,updatedAt,url
{"conclusion":"success","createdAt":"2026-10-15T12:52:40Z","databaseId":38738030458,"event":"workflow_dispatch",
 "headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff",
 "name":"CI","status":"completed","updatedAt":"2026-10-15T12:54:43Z",
 "url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38738030458"}

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540632745	2m6s	2026-10-14T09:12:20Z
completed	cancelled	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540629420	4s	2026-10-14T09:12:19Z
completed	success	CI	CI	ci/add-github-actions-workflow	workflow_dispatch	37456300358	1m58s	2026-10-06T11:25:27Z

### head-SHA tie
run headSha = 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
pushed head = 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
```

Job level (both jobs of the workflow):

```
$ gh api repos/felixfelix-bot/net4sats-wizard-go/actions/runs/38738030458/jobs
lint -> status=completed conclusion=success started=2026-10-15T12:52:44Z id=116432231856
test -> status=completed conclusion=success started=2026-10-15T12:52:45Z id=116432231946
```

Step level — a green conclusion alone is not proof the gates ran, so the steps are read too:

```
--- job id 116432231946  (test)
  1. Set up job -> completed/success
  2. Run actions/checkout@v4 -> completed/success
  3. Run actions/setup-go@v5 -> completed/success
  4. Run go build ./... -> completed/success
  5. Run go vet ./... -> completed/success
  6. Run go test ./... -race -count=1 -> completed/success
--- job id 116432231856  (lint)
  4. Run golangci/golangci-lint-action@v6 -> completed/success
```

And the logs, proving those steps produced real output rather than short-circuiting:

```
$ gh run view --job 116432231856 --log | grep -iE 'golangci|level='
Running [/home/runner/golangci-lint-1.61.0-linux-amd64/golangci-lint config path] …
Running [/home/runner/golangci-lint-1.61.0-linux-amd64/golangci-lint config verify] …
Running [/home/runner/golangci-lint-1.61.0-linux-amd64/golangci-lint run] …
golangci-lint found no issues
Ran golangci-lint in 317ms

$ gh run view --job 116432231946 --log | grep -E 'Run go vet|ok[[:space:]]+net4sats'
test  Run go vet ./...            2026-10-15T12:52:59.0436757Z ##[group]Run go vet ./...
test  Run go vet ./...            2026-10-15T12:52:59.0437126Z go vet ./...
test  Run go test ./... -race -count=1  2026-10-15T12:54:39.1498371Z ok  	net4sats-wizard	82.644s
```

The runner used **golangci-lint 1.61.0** — the version `ci.yml` pins — and it found no issues, matching
the local gate exactly. (The step-number list skips `5` for the lint job: GitHub's jobs API returns no
step with that number for the composite action; every step it *does* return is `success`, and the log
above shows the linter itself ran.)

One incidental warning in both jobs, not a failure:
`Node.js 20 is deprecated. The following actions target Node.js 20 but are being forced to run on
Node.js 24: actions/checkout@v4, actions/setup-go@v5, golangci/golangci-lint-action@v6.`

**PASS — CI green at the pushed head SHA.** No workflow file was edited and no check was weakened.

---

## STEP 6 — report + evidence

- This report and `PROGRESS.md` are committed and pushed on `worker-heavy/t_0509c320-r2` (worktree
  `/home/c08r4d0r/worktrees/wz-e2e-gate-r2`).
- The full summary + raw outputs are posted as a comment on card `t_0509c320`, with this report and
  `PROGRESS.md` attached.
- **The card is in status `triage` and could not be closed from this worker session.** Both attempted
  transitions refused, verbatim (capture `11-card-comment-and-state.txt`):
  ```
  $ hermes kanban complete t_0509c320 --result "pass-2 end-to-end gate PASS: …"
  cannot complete t_0509c320 (unknown id or terminal state)

  $ hermes kanban promote t_0509c320 --force
  cannot promote t_0509c320: task t_0509c320 is 'triage'; promote only applies to 'todo' or 'blocked'
  ```
  `hermes kanban show t_0509c320` still reports `status: triage`
  (last summary: `reclaim backoff until 1791730323 (dead worker pid, strike 1)`). So the deliverable is
  complete and pushed, but **a manager must promote the card out of `triage`** — that is stated rather
  than worked around, and no other transition was claimed.
- The branch under test was never modified: `origin/ci/add-github-actions-workflow-clean` is still
  `0b5322a…`, PR #35's `headRefOid` is unchanged, and `~/wz16-verify` is still clean.



