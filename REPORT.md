# REPORT — E2E integration gate: snapshot build, `--version` smoke, full suite, push

Task: end-to-end verify (the integration gate before cold review) for
`net4sats/tollgate-wizard-go`, "runs after the CI workflow, goreleaser config,
and docs are committed".

**VERDICT: PASS — every step executed, all green, branch pushed, CI green.**
One scope correction (repo/workspace naming) and one honest CI-trigger caveat
are disclosed below; neither is a failure.

---

## 0. Scope correction — the card's repo and workspace names are both wrong

The card names repo `net4sats/tollgate-wizard-go` and workspace
`/home/c08r4d0r/repos/net4sats-mvp-v2`. Neither holds this work:

```
$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)

$ gh api "orgs/net4sats/repos?per_page=100" --jq '.[].name'
configurationwizzard
gonuts-tollgate
net4sats.github.io
kanbanstr
tollgate-module-basic-go
net4sats-wizard-go          <-- the wizard repo
net4sats-feed
net4sats-captive-portal-site
tollgate-captive-portal-site
```

There is **no `tollgate-wizard-go`** in the org; the wizard repo is
`net4sats-wizard-go`. `/home/c08r4d0r/repos/net4sats-mvp-v2` is the net4sats
DEMO repo (not Go). The real work lives in:

```
clone : /home/c08r4d0r/wz16-verify
origin: https://github.com/felixfelix-bot/net4sats-wizard-go.git   (fork)
base  : net4sats/net4sats-wizard-go  (public org repo, upstream)
```

The card's expected output string `tollgate-wizard <version>` is nonetheless
CORRECT: the CLI identity is `const programName = "tollgate-wizard"`
(main.go:55) even though the repo/module is `net4sats-wizard`.

Branch under test = the commit the card says to verify after:

```
HEAD  = 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
        feat(release): link-time version injection, --version, goreleaser
        config and release docs
branch = ci/add-github-actions-workflow-clean
```

Toolchain (absent from default PATH) — `/home/c08r4d0r/wz16-tools/`:
`go122/bin/go` go1.22.12, `goreleaser` 2.18.3, `golangci-lint` 1.61.0.

All raw captures: `/home/c08r4d0r/wz20-captures/`.

---

## STEP 1 — clean tree + branch pushed  **PASS**

```
### host: x280 2026-10-16T03:08:00+02:00
### repo: https://github.com/felixfelix-bot/net4sats-wizard-go.git
### branch: ci/add-github-actions-workflow-clean
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff

$ git status --porcelain
?? PROGRESS.md
?? REPORT.md
STATUS_RC=0  (lines: 2)

$ git status (full)
On branch ci/add-github-actions-workflow-clean
Your branch is up to date with 'origin/ci/add-github-actions-workflow-clean'.

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	PROGRESS.md
	REPORT.md

nothing added to commit but untracked files present (use "git add" to track)
```

The tracked tree is clean; the only entries are this lane's own untracked docs
(PROGRESS.md/REPORT.md), which is why `--porcelain` is not literally empty here.

```
$ git push
Everything up-to-date
PUSH_RC=0

$ git rev-list --left-right --count origin/ci/add-github-actions-workflow-clean...HEAD
0	0

$ git ls-remote origin ci/add-github-actions-workflow-clean
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff	refs/heads/ci/add-github-actions-workflow-clean
$ git rev-parse HEAD
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
```

ahead/behind = **0/0**; remote tip == local HEAD == `0b5322a…`. Pushed.

---

## STEP 2 — full suite  **PASS**

```
$ go mod verify
all modules verified
MODVERIFY_RC=0

$ go build ./...
BUILD_RC=0

$ go vet ./...
VET_RC=0
```

```
$ CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	165.236s
TEST_RACE_RC=0
```

Verbose count run (gcc 13.3.0, CGO_ENABLED=1):

```
TEST_V_RC=0
top-level === RUN : 123
--- PASS          : 52
--- FAIL          : 0
--- SKIP          : 0
WARNING: DATA RACE: 0
ok lines          : 1
FAIL lines        : 0
--- package ok line ---
ok  	net4sats-wizard	124.321s
--- tail -3 ---
--- PASS: TestVersionFlagInjectedViaLdflags (0.52s)
PASS
ok  	net4sats-wizard	124.321s
```

**123 tests run, 0 fail, 0 skip, 0 data races.** The last test to run is
`TestVersionFlagInjectedViaLdflags`, the built binary assertion that fails if
the `-X` symbol in the release config drifts from the `version` var in source.

```
$ golangci-lint run
LINT_RC=0
### golangci-lint: version 1.61.0 built with go1.23.1

$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
GRCHECK_RC=0
### GitVersion: 2.18.3
```

`.golangci.yml` sets `max-issues-per-linter: 0` and `max-same-issues: 0`, so
the green run is not an artifact of the default caps hiding repeats.

---

## STEP 3 — snapshot release  **PASS**

```
$ git describe --tags --dirty
v0.5.0-21-g0b5322a

$ rm -rf dist && goreleaser release --snapshot --clean
  • snapshotting   building snapshot...  version=0.5.1-snapshot.0b5322a
  • running        hook=go mod verify
  • building       binary=net4sats-wizard target=linux_amd64_v1
  • building       binary=net4sats-wizard target=windows_arm64_v8.0
  • building       binary=net4sats-wizard target=darwin_arm64_v8.0
  • building       binary=net4sats-wizard target=darwin_amd64_v1
  • building       binary=net4sats-wizard target=linux_arm64_v8.0
  • building       binary=net4sats-wizard target=windows_amd64_v1
  • archives  (6 archiving lines: linux/darwin tar.gz x amd64/arm64, windows zip x amd64/arm64)
  • calculating checksums
  • release succeeded after 6s
GORELEASER_RC=0
```

Artifact census — **exactly 6 archives, 3 OS x 2 arch, + 1 checksum file**:

```
archives (tar.gz+zip)     : 6
  .tar.gz                 : 4
  .zip                    : 2
checksum file SHA256SUMS  : 1
SHA256SUMS line count     : 6
```

```
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
```

### Checksum verification (`sha256sum -c`, plus two hand-verified artifacts)

```
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
SHA256SUM_C_RC=0

artifact : net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz   (an archive)
  sha256sum   : bdc2b8b99c2263015318e777bb6564d80c6521b858992a2a0159649fcebffe9a
  SHA256SUMS  : bdc2b8b99c2263015318e777bb6564d80c6521b858992a2a0159649fcebffe9a
  VERDICT     : MATCH

artifact : net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz  (NON-LINUX)
  sha256sum   : 26df0ab14e4fa23d6567d2cd31449c5fafb2327a3fb58e6d7587b8302b587a88
  SHA256SUMS  : 26df0ab14e4fa23d6567d2cd31449c5fafb2327a3fb58e6d7587b8302b587a88
  VERDICT     : MATCH

--- openssl cross-check on the linux/amd64 archive ---
SHA2-256(...linux_amd64.tar.gz)= bdc2b8b99c2263015318e777bb6564d80c6521b858992a2a0159649fcebffe9a
  SHA256SUMS: bdc2b8b99c2263015318e777bb6564d80c6521b858992a2a0159649fcebffe9a

--- checksum-count-vs-artifact-count assertion ---
artifacts=6  checksum_lines=6
VERDICT: COUNTS MATCH (6 == 6)
```

`sha256sum -c` passes for all 6; a second, independent tool (openssl) agrees on
the linux/amd64 archive. (The card offered `sha256sum -c` **or** PowerShell
`Get-FileHash`; this is Linux, so `sha256sum`/`openssl` are the right tools.)

---

## STEP 4 — binary smoke (`--version`)  **PASS**

```
$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.0b5322a
RAW_BIN_VERSION_RC=0

$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
... ELF 64-bit LSB executable, x86-64, ... statically linked, ... stripped
```

The card requires `tollgate-wizard <version>` with the injected version (not
`dev`) and exit 0 — satisfied: `tollgate-wizard 0.5.1-snapshot.0b5322a`, RC=0.

### The SHIPPED archive carries the same version (not just the raw build)

```
$ tar -tzf dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
CHANGELOG.md
README.md
net4sats-wizard
$ .../extract-linux-amd64/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.0b5322a
ARCHIVE_BIN_VERSION_RC=0
```

### Controls — the version is genuinely injected, and it changes

```
$ go build -ldflags "-s -w -X main.version=v9.9.9-control" -o /tmp/wz-control .
$ /tmp/wz-control --version
tollgate-wizard v9.9.9-control          <-- manual injection changes it
CONTROL_RC=0

$ go build -o /tmp/wz-nodev .
$ /tmp/wz-nodev --version
tollgate-wizard dev                      <-- no ldflags => dev
NODEV_RC=0

--- assertion table ---
raw goreleaser binary : tollgate-wizard 0.5.1-snapshot.0b5322a
archive binary        : tollgate-wizard 0.5.1-snapshot.0b5322a
manual v9.9.9-control : tollgate-wizard v9.9.9-control
no-ldflags            : tollgate-wizard dev
raw binary prints 'tollgate-wizard 0.5.1-snapshot.0b5322a' : MATCH
archive binary same as raw (injected, not dev)           : MATCH
version CHANGES when injected manually (not hardcoded)   : MATCH
version is 'dev' without ldflags (proves injection does it): MATCH
no output anywhere contains 'dev' in the release binaries : PASS (no dev)
```

---

## STEP 5 — GHA proof  **PASS (green), with one honest trigger caveat**

```
$ gh auth status
✓ Logged in to github.com account felixfelix-bot   (token scopes incl. repo, workflow)

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540632745	2m6s	2026-10-14T09:12:20Z
completed	cancelled	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540629420	4s	2026-10-14T09:12:19Z
RUNLIST_RC=0
```

Filtering the run list by head SHA = the pushed HEAD:

```
### runs whose headSha == 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
{"conclusion":"success","databaseId":38761913548,...}
{"conclusion":"success","databaseId":38738030458,...}
{"conclusion":"success","databaseId":38677575488,...}
{"conclusion":"success","databaseId":38540632745,...}
{"conclusion":"cancelled","databaseId":38540629420,...}
```

Four `success` runs at the exact pushed SHA; the single `cancelled` is the 4s
run superseded by the workflow's `concurrency` group — not a failure.

Run detail proves the jobs are substantive, not no-ops:

```
$ gh run view 38761913548 --repo felixfelix-bot/net4sats-wizard-go
✓ ci/add-github-actions-workflow-clean CI #29 · 38761913548
JOBS
✓ test in 1m49s (ID 116515106517)
✓ lint in 21s   (ID 116515106827)

# test job steps, all success:
Run actions/checkout@v4 / Run actions/setup-go@v5 /
Run go build ./... / Run go vet ./... / Run go test ./... -race -count=1
# lint job steps, all success:
Run actions/checkout@v4 / Run actions/setup-go@v5 /
Run golangci/golangci-lint-action@v6
```

### Caveat (disclosed, not hidden)

Every CI run on this repo is `event=workflow_dispatch`:

```
### push-event runs on this branch (event=push)      -> []   (none)
### pull_request-event runs, any branch              -> []   (none)
```

There is no `push`/`pull_request`-triggered run. The green evidence is real but
dispatch-triggered. This is expected here: the workflow file itself only exists
on this branch, and GitHub does not run a PR's own *newly added* workflow until
it is present on the base branch. Consistent with that, upstream PR #35 reports
no checks:

```
$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json ...
headRefOid=0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff  state=OPEN
mergeable=MERGEABLE  mergeStateStatus=CLEAN
$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch   (RC=1)
```

So: the gate's CI requirement is met at the pushed SHA via dispatch runs, and
the upstream PR is mergeable/clean but will only show checks once the workflow
lands on upstream `main`. **CI is not red, so nothing in the workflow was
touched or weakened.** (For context: fork PR #29, the branch→fork-`main`
consolidation PR, is `CONFLICTING`; it is not this gate.)

---

## Artifacts / evidence index

```
/home/c08r4d0r/wz20-captures/
  01-clean-tree.txt        01-push.txt
  02a-vet.txt              02b-test-race.txt        02c-lint-grcheck.txt
  02d-test-verbose.log     02e-test-counts.txt
  03-snapshot-release.txt  03b-checksums.txt
  04-version-smoke.txt
  05-gha-list.txt          05b-gha-detail.txt
  05c-gha-events-and-pr.txt 05d-upstream-pr.txt
  step*.sh                 (the exact scripts that produced each capture)
/home/c08r4d0r/wz16-verify/dist/   (snapshot artifacts, 6 archives + SHA256SUMS)
```

## STEP 6 — ship: pushed, and CI re-confirmed at the SHIPPED tip  **PASS**

The two lane-doc commits pushed the branch tip from `0b5322a` to `11c661b`, so a
fresh CI run was dispatched and watched at the actual shipped SHA:

```
$ git push
To https://github.com/felixfelix-bot/net4sats-wizard-go.git
   0b5322a..11c661b  ci/add-github-actions-workflow-clean -> ci/add-github-actions-workflow-clean
PUSH_RC=0
$ git rev-list --left-right --count origin/ci/add-github-actions-workflow-clean...HEAD
0	0
$ git ls-remote origin ci/add-github-actions-workflow-clean
11c661b253ea29cfc051ae801d6bd99238c5201f	refs/heads/ci/add-github-actions-workflow-clean
$ git status --porcelain
(empty)                                  <-- clean; nothing uncommitted

$ gh workflow run ci.yml --ref ci/add-github-actions-workflow-clean   (DISPATCH_RC=0)
$ gh run watch 38834487839 --exit-status
✓ Complete job            WATCH_RC=0

$ gh run view 38834487839 --repo felixfelix-bot/net4sats-wizard-go
✓ ci/add-github-actions-workflow-clean CI #29 · 38834487839
JOBS
✓ lint in 14s (ID 116760430198)
✓ test in 1m52s (ID 116760430452)

# --json: conclusion=success  event=workflow_dispatch
#   headSha=11c661b253ea29cfc051ae801d6bd99238c5201f  == local HEAD  ✓
# test steps all success: checkout@v4 / setup-go@v5 / go build ./... /
#                         go vet ./... / go test ./... -race -count=1
# lint steps all success: checkout@v4 / setup-go@v5 / golangci-lint-action@v6
```

Final local HEAD == remote tip == the shipped run's `headSha` == `11c661b…`.

---

## Summary

| Step | Requirement | Result |
|------|-------------|--------|
| 1 | clean tree, pushed, ahead/behind 0 | PASS (0/0; only untracked lane docs) |
| 2 | vet, test -race, lint, goreleaser check | PASS (123 tests, 0 fail/race; all RC=0) |
| 3 | 6 artifacts, checksums verify (2 checked) | PASS (6/6 archives, sha256sum -c all OK) |
| 4 | `--version` = `tollgate-wizard <ver>`, not dev, RC 0 | PASS (+controls change it) |
| 5 | CI green at pushed SHA | PASS (run 38761913548 @ 0b5322a, jobs test+lint success) |
| 6 | report with raw output, pushed, CI green at shipped tip | PASS (pushed to 11c661b; run 38834487839 green) |

**Gate result: PASS. Everything is pushed and CI is green — at both the commit
under test (`0b5322a`) and the shipped tip (`11c661b`).**
