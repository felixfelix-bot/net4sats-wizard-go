# REPORT — t_0509c320: End-to-end verify (snapshot build, `--version` smoke, full suite, push)

Board: `net4sats-mvp-v2` · Card: `t_0509c320` · Assignee: worker-heavy · Author: worker-heavy (fleet offload)
Host: **x280** (`c08r4d0r`) · Date: 2026-10-15
Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate` — branch `worker-heavy/t_0509c320`, cut at `0b5322a`
Subject under test: **`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`**
= `feat(release): link-time version injection, --version, goreleaser config and release docs`
= `origin/ci/add-github-actions-workflow-clean` = head of **PR #35** (`net4sats/net4sats-wizard-go`)
Toolchains: go1.24.4 (`~/wz19-audit/go/bin/go`) for builds; **go1.22.12** (CI-matching, `GOTOOLCHAIN`) for lint
goreleaser 2.18.3 · golangci-lint 1.61.0 · gh 2.45.0
Raw captures: **`/home/c08r4d0r/wz-captures-t_0509c320/`**

## VERDICT: PASS — every gate green, CI green at the pushed head SHA

- `go vet ./...` → rc 0. `gofmt -l .` → no output, rc 0.
- `go test ./... -race -count=1` → **123 RUN / 123 PASS / 0 FAIL / 0 SKIP**, **0** `WARNING: DATA RACE`, exit 0, `ok net4sats-wizard 102.824s`.
- `golangci-lint run` (v1.61.0, CI-matching go1.22.12) → rc 0, no findings.
- `goreleaser check` → rc 0, `1 configuration file(s) validated`.
- `goreleaser release --snapshot --clean` → rc 0, **6 archives** (linux/windows/darwin × amd64/arm64), `SHA256SUMS` with **6 lines** (6 == 6), `sha256sum -c` → **6/6 OK**, plus 3 independent recomputes → **3/3 MATCH**.
- Snapshot binary `--version` → `tollgate-wizard 0.5.1-snapshot.0b5322a`, **rc 0** — a real injected version, not `dev`. Manually re-injecting two other values changed the printed line both times.
- **CI green at the pushed head SHA**: fork run **38677575488**, `conclusion=success`, `headSha=0b5322a…`, jobs `test` ✓ (1m46s) and `lint` ✓ (22s).
- No source, test, or workflow file was changed. The branch under test is byte-identical to `0b5322a`; PR #35's head is unchanged.
- The only commits on this task's branch are its own PROGRESS/REPORT docs, both pushed.

### Four card premises that are wrong on this host (disclosed, not silently worked around)

1. **The repo in the card does not exist.** `net4sats/tollgate-wizard-go` → 404 / "Repository not found". The real repo is `net4sats/net4sats-wizard-go`. Detail in STEP 5.
2. **`gh run list --repo <upstream>` cannot show this PR's CI at all** — upstream `main` has no `.github/`, so no workflow is registered there and PR #35 reports *no checks*. The CI workflow this PR adds only runs on the fork. Green CI was therefore obtained from the fork at the exact head SHA. Detail in STEP 5.
3. **`golangci-lint run` fails under go1.24** with `unsupported version: 2` — a toolchain/export-data mismatch, not a code defect. Re-run under the CI-matching go1.22.12 → clean. Detail in STEP 2c.
4. **`go test ./...` alone would not have been the mandated gate** — the card says `-race`, and `-race` needs `CGO_ENABLED=1` + a C toolchain (gcc 13.3.0 present here). The command actually run is named in every capture.

---

## STEP 1 — clean tree + pushed, ahead/behind 0/0

Capture: `01-clean-tree-and-push.txt` (run in the clone `~/wz16-verify`, the checkout that holds the branch under test).

```
### host: x280 c08r4d0r 2026-10-15T05:36:45+02:00
### go version: go version go1.24.4 linux/amd64
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff 0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
### branch: ci/add-github-actions-workflow-clean

$ git status --porcelain
porcelain_rc=0                     # no output = clean

$ git status -sb
## ci/add-github-actions-workflow-clean...origin/ci/add-github-actions-workflow-clean

$ git push
Everything up-to-date
push_rc=0

$ git fetch origin && git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
From https://github.com/felixfelix-bot/net4sats-wizard-go
 + 134e235...cc5daa8 wz19-video-evidence -> origin/wz19-video-evidence  (forced update)
 * [new branch]      wz19-video-evidence-pre-atomic -> origin/wz19-video-evidence-pre-atomic
0	0                              # ahead/behind = 0/0

$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff	refs/heads/ci/add-github-actions-workflow-clean
```

Clean, nothing to push, ahead/behind `0 0`, remote ref = the local head. **PASS.**

The worktree used for the gates was then created from that exact sha (`00-worktree-add.txt`):

```
$ git worktree add -b worker-heavy/t_0509c320 /home/c08r4d0r/worktrees/wz-e2e-gate 0b5322a
Preparing worktree (new branch 'worker-heavy/t_0509c320')
HEAD is now at 0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
```

---

## STEP 2 — full suite

### 2a — `go vet ./...` and module shape

Capture: `02-vet-gofmt.txt`.

```
### go: go version go1.24.4 linux/amd64
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff

$ go vet ./...
VET_RC=0

$ gofmt -l .
GOFMT_RC=0                         # no output = no drift

$ go list ./...
net4sats-wizard
LIST_RC=0

$ find . -name go.mod (module count)
./go.mod
```

One module, one package — so `./...` is genuinely the whole suite here (the multi-module false-green trap does not apply: `go list ./...` prints the single package).

### 2b — `go test ./... -race -count=1` (123/123)

Capture: `04-test-race.txt` (254 lines of raw `-v` output on disk).

```
### go: go version go1.22.12 linux/amd64  CGO_ENABLED=1
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
### command: go test ./... -race -count=1 -v   (identical to .github/workflows/ci.yml test job)

=== RUN   TestHandleScan
--- PASS: TestHandleScan (51.49s)
...
--- PASS: TestVersionStringFormat (0.00s)
=== RUN   TestVersionFlagInjectedViaLdflags
--- PASS: TestVersionFlagInjectedViaLdflags (0.64s)
PASS
ok  	net4sats-wizard	102.824s
TEST_RACE_RC=0
```

Tally over the raw capture (script `tally.sh`, plus an independent Python recount):

```
RUN(top-level)   : 123
RUN(all incl sub): 123
PASS(all)        : 123
FAIL(all)        : 0
SKIP(all)        : 0
DATA RACE warns  : 0
ok/FAIL lines    : ['--- PASS: TestVersionStringFormat (0.00s)', '--- PASS: TestVersionFlagInjectedViaLdflags (0.64s)', 'ok  \tnet4sats-wizard\t102.824s']
RC               : ['TEST_RACE_RC=0']
total lines      : 254
```

`grep -E '^--- (FAIL|SKIP)'` over the capture returns nothing. **0 skips** — note the two jq round-trip hardening tests do *not* skip here, because `jq 1.7.1` is present at `~/.local/bin/jq` (installed by the sibling card `t_10a913ec`). **PASS.**

`-race` was run with `CGO_ENABLED=1` and gcc 13.3.0 — with cgo disabled this command exits 2 before running anything, so the flag is named explicitly rather than implied.

### 2c — `golangci-lint run`

Capture: `03-golangci-lint.txt`.

First attempt, go1.24.4 in PATH (recorded for honesty — this is an environment defect):

```
### golangci-lint: golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z

$ golangci-lint run
level=warning msg="[runner] Can't run linter goanalysis_metalinter: buildir: failed to load package goarch: could not load export data: internal error in importing \"internal/goarch\" (unsupported version: 2); please report an issue"
level=error msg="Running error: can't run linter goanalysis_metalinter\nbuildir: failed to load package goarch: could not load export data: internal error in importing \"internal/goarch\" (unsupported version: 2); please report an issue"
GOLANGCI_RC=3
```

That is golangci-lint 1.61.0 (built with go1.23.1) being handed go1.24 export data it cannot parse — **not** a finding about this repo. CI does not hit it because `actions/setup-go@v5` installs from `go-version-file: go.mod` (`go 1.22`). Re-run under that same toolchain, with a fresh `GOCACHE` so no go1.24 artifacts leak in:

```
### go in PATH: go version go1.22.12 linux/amd64
### golangci-lint: golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z
### GOCACHE=/home/c08r4d0r/.cache/go-build-122 (fresh, avoids go1.24 export-data poisoning)

$ golangci-lint run
GOLANGCI_RC=0                      # no findings
```

The repo's `.golangci.yml` enables exactly the v1 default set with `max-issues-per-linter: 0` / `max-same-issues: 0`, i.e. every finding is reported — so rc 0 means genuinely clean, not capped. **PASS.**

### 2d — `goreleaser check`

Capture: `05-goreleaser-check.txt`.

```
### goreleaser: GitVersion:    2.18.3
GitCommit:     1942d44355492f44cd3774c57bd3c27fe544dfe9

$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
GORELEASER_CHECK_RC=0
```

The local goreleaser is **2.18.3**, the same version `.github/workflows/release.yml` pins (`version: v2.18.3`), and the config is v2 schema. **PASS.**

---

## STEP 3 — snapshot release, artifacts and checksums

Capture: `06-snapshot-artifacts-smoke.txt` (and the corrected count check in `07-artifact-count-corrected.txt`).

```
### goreleaser: GitVersion:    2.18.3
### git status --porcelain (must be empty):
porcelain_rc=0

$ goreleaser release --snapshot --clean
  • starting release
  • skipping announce, publish, and validate...
  • cleaning distribution directory
  • loading environment variables
  • getting and validating git state
    • using tags                                     previous=v0.3.5 current=v0.5.0
    • pipe skipped or partially skipped              reason=disabled during snapshot mode
  • parsing tag
  • setting defaults
  • snapshotting
    • building snapshot...                           version=0.5.1-snapshot.0b5322a
  • running before hooks
    • running                                        hook=go mod verify
  • ensuring distribution directory
  • setting up metadata
  • writing release metadata
  • loading go mod information
  • build prerequisites
  • building binaries
    • building                                       paths=. binaries=net4sats-wizard target=darwin_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=linux_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=linux_arm64_v8.0
    • building                                       paths=. binaries=net4sats-wizard target=windows_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=darwin_arm64_v8.0
    • building                                       paths=. binaries=net4sats-wizard target=windows_arm64_v8.0
      • took: 1m31s
  • archives
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 1m34s
  • thanks for using GoReleaser!
GORELEASER_SNAPSHOT_RC=0
```

All 6 targets built (linux/windows/darwin × amd64/arm64), 6 archives written, checksums calculated. Snapshot version = **`0.5.1-snapshot.0b5322a`** (from tag `v0.5.0` + `incpatch` + short commit, per `snapshot.version_template`).

### 3a — artifact list

```
$ ls -la dist/
total 27784
-rw-r--r--  1 c08r4d0r c08r4d0r     744 Oct 15 05:44 SHA256SUMS
-rw-r--r--  1 c08r4d0r c08r4d0r    4754 Oct 15 05:44 artifacts.json
-rw-r--r--  1 c08r4d0r c08r4d0r    4379 Oct 15 05:42 config.yaml
-rw-r--r--  1 c08r4d0r c08r4d0r     249 Oct 15 05:42 metadata.json
-rw-rw-r--  1 c08r4d0r c08r4d0r 4862526 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
-rw-rw-r--  1 c08r4d0r c08r4d0r 4699536 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
-rw-rw-r--  1 c08r4d0r c08r4d0r 4804887 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
-rw-rw-r--  1 c08r4d0r c08r4d0r 4542637 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
-rw-rw-r--  1 c08r4d0r c08r4d0r 4884019 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
-rw-rw-r--  1 c08r4d0r c08r4d0r 4579788 Oct 15 05:44 net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_darwin_amd64_v1
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_darwin_arm64_v8.0
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_linux_amd64_v1
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_linux_arm64_v8.0
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_windows_amd64_v1
drwxr-xr-x  2 c08r4d0r c08r4d0r    4096 Oct 15 05:44 net4sats-wizard_windows_arm64_v8.0
```

### 3b — checksum count == artifact count (6 == 6)

**Correction, disclosed:** the first count in `06-…` used `find dist -maxdepth 1 -type f ! -name SHA256SUMS`, which also matched goreleaser's own working files in `dist/` (`config.yaml`, `metadata.json`, `artifacts.json`) and so printed `ARTIFACT_COUNT=9` / `COUNT_MATCH=NO`. Those three are not release artifacts. Re-counted against the published set (`.tar.gz` + `.zip`) in `07-artifact-count-corrected.txt`:

```
--- release artifacts = dist/*.tar.gz + dist/*.zip ---
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
ARTIFACT_COUNT=6

--- checksum file ---
CHECKSUM_LINE_COUNT=6
CHECKSUM_FILE_NAME=SHA256SUMS

COUNT_MATCH=YES (6 checksum lines == 6 release artifacts)

--- per-target coverage: expect linux/windows/darwin x amd64/arm64 ---
  linux    amd64  -> 1 archive(s)
  linux    arm64  -> 1 archive(s)
  windows  amd64  -> 1 archive(s)
  windows  arm64  -> 1 archive(s)
  darwin   amd64  -> 1 archive(s)
  darwin   arm64  -> 1 archive(s)

--- goreleaser internal working files in dist/ (NOT release artifacts) ---
artifacts.json
config.yaml
metadata.json

--- metadata.json (version goreleaser injected) ---
{"project_name":"net4sats-wizard","tag":"v0.5.0","previous_tag":"v0.3.5","version":"0.5.1-snapshot.0b5322a","commit":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","date":"2026-10-15T05:42:44.255197998+02:00","runtime":{"goos":"linux","goarch":"amd64"}}
```

`SHA256SUMS` is the pinned name from `.goreleaser.yaml` (`checksum.name_template: SHA256SUMS`) — the name README documents for `sha256sum -c`.

### 3c — checksum verification

`sha256sum -c` over the whole file:

```
$ (cd dist && sha256sum -c SHA256SUMS)
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
SHA256SUM_C_RC=0
```

6/6 OK. Plus independent recomputes (`Get-FileHash`-style: rehash the file, compare to the digest on its `SHA256SUMS` line) for **one linux archive, one non-linux archive, and one darwin arm64 archive**:

```
  net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz  expected=4c379b59b8f084dcaf334d74e4a78715d4a8db09355909f385641f3f022bce86
                                                             actual  =4c379b59b8f084dcaf334d74e4a78715d4a8db09355909f385641f3f022bce86
  -> MATCH
  net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip   expected=3eb48cf6f572aae50361b437b1292b9aafa600da3e1eb93d9b64e5b71250a0b0
                                                             actual  =3eb48cf6f572aae50361b437b1292b9aafa600da3e1eb93d9b64e5b71250a0b0
  -> MATCH
  net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz expected=fe85b8dbf46c3b22a8748070d927eb760ce86fbf014d3dee80c07987d37b67bf
                                                             actual  =fe85b8dbf46c3b22a8748070d927eb760ce86fbf014d3dee80c07987d37b67bf
  -> MATCH
```

**PASS.** (Verification was run on this Linux host; `Get-FileHash` itself is PowerShell-only, so the equivalent operation — independent digest recompute against the published line — was performed natively for the non-linux artifact.)

Archive contents sanity (`unzip -l` on the windows zip, `tar -tzf` on the linux one) show `CHANGELOG.md`, `README.md` and the binary/`.exe`, so the non-linux artifacts are real, populated archives:

```
$ unzip -l dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
     1575  2026-10-15 05:36   CHANGELOG.md
     8264  2026-10-15 05:36   README.md
 12382720  2026-10-15 05:44   net4sats-wizard.exe
---------                     -------
 12392559                     3 files
```

---

## STEP 4 — binary smoke (`--version`) + manual injection

Capture: `06-snapshot-artifacts-smoke.txt`, GATE 4 section.

Raw built linux/amd64 binary:

```
raw built binary: dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
$ <raw linux/amd64 binary> --version
tollgate-wizard 0.5.1-snapshot.0b5322a
RAW_VERSION_RC=0
```

The binary shipped inside the linux/amd64 archive (extracted with `tar -xzf`, contents `CHANGELOG.md`, `README.md`, `net4sats-wizard`):

```
$ <extracted binary> --version
tollgate-wizard 0.5.1-snapshot.0b5322a
ARCHIVE_VERSION_RC=0
```

Printed line is exactly `tollgate-wizard <version>` with the **injected** `0.5.1-snapshot.0b5322a` — not `dev` — and exit code 0. (`tollgate-wizard` is `main.go`'s `const programName`, independent of the `net4sats-wizard` binary/project name, so the card's expected string is right.)

Manual injection, proving the value is really coming from `-X main.version` and is not hardcoded:

```
$ go build (no ldflags) -o scratch/plain && scratch/plain --version
tollgate-wizard dev
PLAIN_VERSION_RC=0

$ go build -ldflags "-X main.version=v9.9.9-manual-smoke" -o scratch/injected && scratch/injected --version
tollgate-wizard v9.9.9-manual-smoke
INJECTED_VERSION_RC=0

$ go build -ldflags "-X main.version=v0.0.0-second-different-value" -o scratch/injected2 && scratch/injected2 --version
tollgate-wizard v0.0.0-second-different-value
INJECTED2_VERSION_RC=0
```

Three distinct outputs for three different inputs (`dev`, `v9.9.9-manual-smoke`, `v0.0.0-second-different-value`) — the injected value demonstrably changes what `--version` prints. **PASS.**

Build hygiene: `git status --porcelain` after all of the above is still empty (`post_porcelain_rc=0`) — the build did not dirty any tracked file (the repo tracks `release/SHA256SUMS`, but goreleaser writes to the gitignored `dist/`, so that tracked file is untouched).

---

## STEP 5 — GitHub Actions proof

Capture: `08-gha-ci-proof.txt` (verbatim, all four probes).

### 5a — the card's exact command fails: the repo does not exist

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
GH_RUN_LIST_BRIEF_RC=1

$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)
GH_REPO_VIEW_RC=1

$ git ls-remote https://github.com/net4sats/tollgate-wizard-go
remote: Repository not found.
fatal: repository 'https://github.com/net4sats/tollgate-wizard-go/' not found
LS_REMOTE_RC=128
```

`gh` is authenticated (`felixfelix-bot`, token scopes include `repo`), so this is a genuine "no such repository", not an auth failure. The real repo is `net4sats/net4sats-wizard-go`, as the lane docs already alias ("The real PR is #32 in `net4sats/net4sats-wizard-go` (a.k.a. `tollgate-wizard-go`)").

### 5b — the real upstream repo has ZERO CI runs, and PR #35 can never report checks

```
$ gh repo view net4sats/net4sats-wizard-go --json nameWithOwner,isPrivate,defaultBranchRef
{"defaultBranchRef":{"name":"main"},"isPrivate":false,"nameWithOwner":"net4sats/net4sats-wizard-go"}

$ gh run list --repo net4sats/net4sats-wizard-go --limit 5
GH_RUN_LIST_UPSTREAM_RC=0          # empty output = no workflow run has EVER existed upstream

$ gh workflow list --repo net4sats/net4sats-wizard-go
                                   # empty = no workflow is registered upstream

$ gh api repos/net4sats/net4sats-wizard-go/contents/.github
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":404}

$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json state,headRefOid,headRefName,baseRefName
{"baseRefName":"main","headRefName":"ci/add-github-actions-workflow-clean","headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","state":"OPEN"}

$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch
GH_PR_CHECKS_RC=1

$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json statusCheckRollup --jq '.statusCheckRollup'
[]
```

**Root cause:** upstream `main` has no `.github/` at all — that is exactly what PR #34/#35 are trying to land. Until the workflow exists on the base repo's default branch, a fork PR's `pull_request` run cannot be registered, so PR #35 shows no checks and the upstream `gh run list` is empty. This is inherent to the fork-first flow for this repo, not something a re-push can fix; the fix is landing the CI workflow upstream (which is the PR under test).

### 5c — green CI, produced on the fork at the exact head SHA

The fork owns the branch and the workflow:

```
$ gh workflow list --repo felixfelix-bot/net4sats-wizard-go
CI	active	376283443
```

A **fresh** run was dispatched by this task so the proof is current and unambiguously tied to `0b5322a`:

```
$ gh workflow run ci.yml --repo felixfelix-bot/net4sats-wizard-go --ref ci/add-github-actions-workflow-clean
DISPATCH_RC=0

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540632745	2m6s	2026-10-14T09:12:20Z
completed	cancelled	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540629420	4s	2026-10-14T09:12:19Z
completed	success	CI	CI	ci/add-github-actions-workflow	workflow_dispatch	37456300358	1m58s	2026-10-06T11:25:27Z
completed	failure	CI	CI	ci/add-github-actions-workflow	workflow_dispatch	37450967243	2m3s	2026-10-06T10:37:06Z
```

The required line, with the head SHA proven to be the pushed one:

```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5 \
    --json databaseId,headSha,headBranch,event,status,conclusion,workflowName,createdAt
[{"conclusion":"success","createdAt":"2026-10-15T03:41:31Z","databaseId":38677575488,"event":"workflow_dispatch",
  "headBranch":"ci/add-github-actions-workflow-clean",
  "headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","workflowName":"CI"}, …]
```

`gh run view` on that run:

```
$ gh run view 38677575488 --repo felixfelix-bot/net4sats-wizard-go

✓ ci/add-github-actions-workflow-clean CI #29 · 38677575488
Triggered via workflow_dispatch about 4 minutes ago

JOBS
✓ test in 1m46s (ID 116234576134)
✓ lint in 22s (ID 116234576515)
```

and the per-step JSON, which shows the CI jobs really ran the mandated commands:

```
{"conclusion":"success","event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean",
 "headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff",
 "jobs":[{"name":"test","conclusion":"success","startedAt":"2026-10-15T03:41:35Z","completedAt":"2026-10-15T03:43:21Z",
          "steps":[{"name":"Set up job","conclusion":"success"},
                   {"name":"Run actions/checkout@v4","conclusion":"success"},
                   {"name":"Run actions/setup-go@v5","conclusion":"success"},
                   {"name":"Run go build ./...","conclusion":"success"},
                   {"name":"Run go vet ./...","conclusion":"success"},
                   {"name":"Run go test ./... -race -count=1","conclusion":"success"}, …]},
         {"name":"lint","conclusion":"success","startedAt":"2026-10-15T03:41:36Z","completedAt":"2026-10-15T03:41:58Z",
          "steps":[{"name":"Run actions/checkout@v4","conclusion":"success"},
                   {"name":"Run actions/setup-go@v5","conclusion":"success"},
                   {"name":"Run golangci/golangci-lint-action@v6","conclusion":"success"}, …]}]}
```

**CI is green at `0b5322a`** — `conclusion=success`, both jobs green, and the `test` job ran `go build ./...`, `go vet ./...` and `go test ./... -race -count=1` successfully. The only annotations are GitHub platform deprecation notices (Node.js 20 → 24 for the pinned actions; `ubuntu-latest` → Ubuntu 26 migration date). **No check was weakened to make this pass** — the workflow was not edited at all.

---

## Final state — pushed, clean

Capture: `09-final-state.txt` (re-run after the REPORT push).

```
### worktree: /home/c08r4d0r/worktrees/wz-e2e-gate  (branch worker-heavy/t_0509c320)

$ git status
On branch worker-heavy/t_0509c320
Your branch is up to date with 'origin/worker-heavy/t_0509c320'.

nothing to commit, working tree clean

$ git status --porcelain
porcelain_rc=0      # no output = clean

$ git log --oneline -3
418a260 docs(t_0509c320): full e2e verification REPORT (suite 123/123, snapshot 6/6, checksums, --version, CI green)
05902b2 docs(t_0509c320): PROGRESS — e2e gate findings (repo alias, fork-only CI, lint toolchain) at 0b5322a
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs

$ git push
Everything up-to-date
push_rc=0

$ git rev-list --left-right --count HEAD...origin/worker-heavy/t_0509c320
0	0

$ git ls-remote origin refs/heads/worker-heavy/t_0509c320
418a2602394a7645b40c4183b68da5c3cfc9dc02	refs/heads/worker-heavy/t_0509c320

$ git rev-parse HEAD
418a2602394a7645b40c4183b68da5c3cfc9dc02

### the branch under test must still be byte-identical to 0b5322a:
$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff	refs/heads/ci/add-github-actions-workflow-clean

$ git -C /home/c08r4d0r/wz16-verify status --porcelain   (clone of the branch under test)
clone_porcelain_rc=0   # no output = the PR head tree was never touched

$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json state,headRefOid
{"headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","state":"OPEN"}
```

Both docs commits are pushed to the fork (`felixfelix-bot/net4sats-wizard-go`):
`PROGRESS.md` at `05902b2`, `REPORT.md` at `418a260` (remote `ls-remote` confirms the same sha, ahead/behind `0 0`).
This final-state paragraph itself lands in a trailing `docs(t_0509c320): final-state capture` commit on the
same branch.

The branch under test was **not** modified: `origin/ci/add-github-actions-workflow-clean` is still
`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`, so PR #35's head is byte-identical to the reviewed sha.
Build outputs (`dist/`, `net4sats-wizard`) are gitignored, so "clean" is genuine and not an artifact of
ignoring new files.

## Reproduce

```
export PATH=$HOME/wz19-audit/go/bin:$HOME/.local/bin:$HOME/wz16-tools:$PATH
cd /home/c08r4d0r/worktrees/wz-e2e-gate          # or: git worktree add -b <b> <dir> 0b5322a

go vet ./...                                     # expect rc 0
CGO_ENABLED=1 go test ./... -race -count=1       # expect: ok net4sats-wizard ~103s, rc 0
                                                 # (use ~/wz16-tools/go122/bin/go for the lint gate)
golangci-lint run                                # expect rc 0  [go1.22.12 in PATH!]
goreleaser check                                 # expect rc 0
goreleaser release --snapshot --clean            # expect rc 0, 6 archives + SHA256SUMS
( cd dist && sha256sum -c SHA256SUMS )           # expect 6/6 OK
dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version   # expect: tollgate-wizard 0.5.1-snapshot.<sha>

gh workflow run ci.yml --repo felixfelix-bot/net4sats-wizard-go --ref ci/add-github-actions-workflow-clean
gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 3   # expect: completed success
```

## Residuals / notes for the fleet

- **Card premise correction:** the repo is `net4sats/net4sats-wizard-go` (public), not `net4sats/tollgate-wizard-go`; the fork is `felixfelix-bot/net4sats-wizard-go`. Any card that says "`gh run list --repo net4sats/tollgate-wizard-go` must show success" is unsatisfiable as written — 404.
- **CI venue:** upstream cannot report checks until the workflow lands on upstream `main`. Fork-side `workflow_dispatch` on `ci/add-github-actions-workflow-clean` is the only current source of green CI for this PR. The run dispatched here is `38677575488`.
- **Local lint toolchain:** use `GOTOOLCHAIN=go1.22.12` (shimmed at `~/wz16-tools/go122/bin/go`) with a separate `GOCACHE`; go1.24 makes golangci-lint 1.61.0 exit 3 with `unsupported version: 2`. Nothing in the repo needs changing for this.
- `-race` requires `CGO_ENABLED=1` and gcc; the box has gcc 13.3.0, so race evidence is available here.
- Worktree `/home/c08r4d0r/worktrees/wz-e2e-gate` is registered against the clone `/home/c08r4d0r/wz16-verify`; remove with `git -C /home/c08r4d0r/wz16-verify worktree remove /home/c08r4d0r/worktrees/wz-e2e-gate` when done.
- Raw captures for every step: `/home/c08r4d0r/wz-captures-t_0509c320/` (`01`…`08` plus the gate scripts).
