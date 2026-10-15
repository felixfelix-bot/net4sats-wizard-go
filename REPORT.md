# REPORT — t_0509c320 (pass 3): end-to-end verify — snapshot build, `--version` smoke, full suite, push

Board `net4sats-mvp-v2` · card `t_0509c320` · assignee worker-heavy · host **x280** (`c08r4d0r`) · 2026-10-15
Worktree (this pass): `/home/c08r4d0r/worktrees/wz-e2e-gate-r3` — branch `worker-heavy/t_0509c320-r3`, cut at `0b5322a`
Subject under test: **`0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff`** =
`feat(release): link-time version injection, --version, goreleaser config and release docs`
= `origin/ci/add-github-actions-workflow-clean` (checkout `~/wz16-verify`) = **head of PR #35** (OPEN)
Toolchain: **go1.22.12** (CI-matching — `ci.yml` uses `go-version-file: go.mod` = `go 1.22`), golangci-lint
**1.61.0**, goreleaser **2.18.3**, gcc 13.3.0, gh 2.45.0 (account `felixfelix-bot`)
Raw captures: **`/home/c08r4d0r/wz-captures-t_0509c320-p3/`** (`00`…`11`) · scripts `~/wz16-tools/rev3-*.sh`

**VERDICT: PASS** — every mandated gate executed live in this pass; 6/6 artifacts verified with
independent recomputes and negative controls; `--version` prints the injected version (not `dev`); the CI
workflow is green on the fork **at the pushed head SHA**.

This is **pass 3** of the card (pass 1 = `PROGRESS.t_0509c320.md` / `REPORT.t_0509c320.md`, pass 2 =
`…-pass2.md` in the lane repo). Nothing from pass 1/2 is re-cited as this pass's evidence: every command
below was re-executed on this host and its raw output captured. The subject SHA is unchanged since pass 1.

Three card premises are wrong on this host, disclosed rather than silently worked around:
(1) the repo `net4sats/tollgate-wizard-go` does not exist — it is `net4sats/net4sats-wizard-go`;
(2) the mandated `gh run list --repo <the card's repo>` returns **HTTP 404**, and even the real upstream
repo has **no workflows at all** (its `main` lacks `.github/`, which is what the PR adds), so green CI can
only be read on the fork;
(3) `golangci-lint run` fails under the host go1.24 with `unsupported version: 2` — a toolchain/export-data
mismatch, not a code defect; under the CI-matching go1.22.12 it is clean.

---

## STEP 0 — worktree registration (unchanged source)

```
### host: x280 c08r4d0r 2026-10-15T17:44:28+02:00
### NOTE: the `git worktree add` stdout was emitted in the session transcript and NOT tee'd at the time;
###       the registration it produced is proven below from git's own state (authoritative).

$ git worktree list
/home/c08r4d0r/wz16-verify               0b5322a [ci/add-github-actions-workflow-clean]
/home/c08r4d0r/worktrees/wizard-fix      3ff9ae4 [worker-heavy/t_10a913ec]
/home/c08r4d0r/worktrees/wz-e2e-gate     eb8a1a6 [worker-heavy/t_0509c320]
/home/c08r4d0r/worktrees/wz-e2e-gate-r2  8a84dc8 [worker-heavy/t_0509c320-r2]
/home/c08r4d0r/worktrees/wz-e2e-gate-r3  0b5322a [worker-heavy/t_0509c320-r3]
/home/c08r4d0r/wz-ci-atomic/repo         13de851 [ci/github-actions-workflow-atomic]

$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 rev-parse HEAD
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 rev-parse --abbrev-ref HEAD
worker-heavy/t_0509c320-r3

$ git worktree list | grep -c wz-e2e-gate-r3   # exactly one registration
1

$ git -C /home/c08r4d0r/wz16-verify status --porcelain   # the reviewed checkout is untouched
porcelain_rc=0
```

---

## STEP 1 — clean tree check + push, ahead/behind 0/0

```
### host: x280 c08r4d0r 2026-10-15T17:35:52+02:00
### repo: /home/c08r4d0r/wz16-verify
### remote: origin	https://github.com/felixfelix-bot/net4sats-wizard-go.git (fetch)
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff  0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
### branch: ci/add-github-actions-workflow-clean

$ git status --porcelain
porcelain_rc=0   # no lines printed above = clean tree

$ git status -sb
## ci/add-github-actions-workflow-clean...origin/ci/add-github-actions-workflow-clean

$ git push
Everything up-to-date
push_rc=0

$ git fetch origin
fetch_rc=0

$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0
count_rc=0   # <ahead>\t<behind>; want 0\t0

$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff	refs/heads/ci/add-github-actions-workflow-clean

$ git diff --stat origin/ci/add-github-actions-workflow-clean..HEAD
diff_rc=0   # empty = byte-identical to the pushed head
```

Clean tree, `git push` = `Everything up-to-date`, ahead/behind **0/0**, remote ref == local HEAD.
**PASS.**

---

## STEP 2 — full suite

### 2a — `go vet ./...`, `gofmt`, module shape

```
### host: x280 c08r4d0r 2026-10-15T17:36:05+02:00
### worktree: /home/c08r4d0r/worktrees/wz-e2e-gate-r3
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ go version
go version go1.22.12 linux/amd64
### GOROOT: /home/c08r4d0r/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.22.12.linux-amd64
### GOFLAGS: 
### GOTOOLCHAIN: go1.22.12
### CGO_ENABLED: 1

$ find . -name go.mod -not -path '*/.git/*' | wc -l
1

$ go list ./...
net4sats-wizard
list_rc=0

$ go vet ./...
VET_RC=0

$ gofmt -l .
GOFMT_RC=0   # no filenames above = no formatting drift

$ git status --porcelain   # did vet/gofmt dirty the tree?
porcelain_rc=0
```

**A single module** (`go list ./...` → one package `net4sats-wizard`), so the mandated `./...` gates cover
the whole test surface — no nested-module vacuity. `git status --porcelain` empty after vet/gofmt.
**PASS.**

### 2b — `go test ./... -race -count=1`


```
### host: x280 c08r4d0r 2026-10-15T17:36:40+02:00
### go: go version go1.22.12 linux/amd64   cpus: 8   load: 7.33 3.64 3.61
### gcc: gcc (Ubuntu 13.3.0-6ubuntu2~24.04.1) 13.3.0
### CGO_ENABLED=1
### worktree: /home/c08r4d0r/worktrees/wz-e2e-gate-r3
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff

$ cd /home/c08r4d0r/worktrees/wz-e2e-gate-r3 && CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	92.390s
TEST_RC=0
```

counts derived from the verbose re-run (`04-test-race-verbose.log`, 123 `=== RUN` lines):


```
### counts derived from 04-test-race-verbose.log
RUN  (=== RUN)   : 123
PASS (--- PASS)  : 52
FAIL (--- FAIL)  : 0
SKIP (--- SKIP)  : 0
DATA RACE warns  : 0
ok lines         : 1
ok  	net4sats-wizard	94.214s
```

`123 RUN / 0 FAIL / 0 SKIP / 0 data races`, exit 0. **PASS.**

### 2c — `golangci-lint run` (v1.61.0)

```
### golangci-lint  (PATH: /home/c08r4d0r/wz16-tools/golangci-lint)
golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z

$ golangci-lint run
LINT_RC=0   # no findings printed above = clean
```

### 2d — `goreleaser check` (v2.18.3)

```
$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
GRELEASER_CHECK_RC=0
```

---

## STEP 3 — snapshot release: 6 artifacts, checksums

```
### host: x280 c08r4d0r 2026-10-15T17:40:15+02:00
### worktree: /home/c08r4d0r/worktrees/wz-e2e-gate-r3
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
### goreleaser:   ____       ____      _

$ cd /home/c08r4d0r/worktrees/wz-e2e-gate-r3 && goreleaser release --snapshot --clean
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
    • building                                       paths=. binaries=net4sats-wizard target=windows_arm64_v8.0
    • building                                       paths=. binaries=net4sats-wizard target=darwin_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=linux_arm64_v8.0
    • building                                       paths=. binaries=net4sats-wizard target=linux_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=windows_amd64_v1
    • building                                       paths=. binaries=net4sats-wizard target=darwin_arm64_v8.0
      • took: 1m4s
  • archives
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
    • archiving                                      name=dist/net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 1m9s
  • thanks for using GoReleaser!
SNAPSHOT_RC=0
```

Artifact inventory + checksum verification (this file deliberately contains **only** the real run;
corruption controls live in §3b):

```
### host: x280 c08r4d0r 2026-10-15T17:44:07+02:00
### cwd: /home/c08r4d0r/worktrees/wz-e2e-gate-r3/dist   (dist/ of the at-head snapshot release)

$ ls -1 *.tar.gz *.zip | sort
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip
archive_count=6

$ ls -1 *.tar.gz *.zip | sed 's/.*_\(linux\|windows\|darwin\)_.*/\1/' | sort | uniq -c
      2 darwin
      2 linux
      2 windows

$ for a in *.tar.gz *.zip; do echo "$a -> $(file -b $a)"; done
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz -> gzip compressed data, max compression, original size modulo 2^32 11445760
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz -> gzip compressed data, max compression, original size modulo 2^32 11225088
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz -> gzip compressed data, max compression, original size modulo 2^32 11265536
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz -> gzip compressed data, max compression, original size modulo 2^32 11023872
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip -> Zip archive data, at least v2.0 to extract, compression method=deflate
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip -> Zip archive data, at least v2.0 to extract, compression method=deflate

$ ls -1 *.tar.gz *.zip | sed 's/.*\.\(tar\.gz\|zip\)$/\1/' | sort | uniq -c   # archive formats
      4 tar.gz
      2 zip

$ wc -l < SHA256SUMS     # checksum lines must equal the artifact count
6

$ cat SHA256SUMS
868a8845c8825b2b73723deb9ca6581de6f64ba9c1a43607c83c7f038c686601  net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz
627deeecaa79eb5c5404a99b1a79180279926a978d1a99769e2979839eafcaf0  net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz
77cefac426e777045af6ce34b9b3aabb1ca242801982bcfc421a324d225b7507  net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz
1b9171ac3aadaab24adf166b70c614ad348af65aa13ffc62912feba2ab46b7ad  net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz
6a646141c4ccc62de3414c27f4e08ca82d67c929732a9bfa212aa64b6c8a16ba  net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip
6c6afa43c3126f03246d8cdfd6eecf7e8eee3ffda63c2593914437c75e0a17a4  net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip

$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
SHA256SUM_C_REAL_RC=0   # rc of sha256sum itself (no pipe)

# independent recompute (not via -c): fresh digest vs the SHA256SUMS entry
MATCH    net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz  77cefac426e777045af6ce34b9b3aabb1ca242801982bcfc421a324d225b7507
MATCH    net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz  627deeecaa79eb5c5404a99b1a79180279926a978d1a99769e2979839eafcaf0
MATCH    net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip  6a646141c4ccc62de3414c27f4e08ca82d67c929732a9bfa212aa64b6c8a16ba

# Windows equivalent of Get-FileHash: PowerShell is not installed on this host, so the
# zip digests are computed with sha256sum and compared against SHA256SUMS above (MATCH).
pwsh/powershell ABSENT -> sha256sum is the equivalent used
```

6 archives = linux / windows / darwin × amd64 / arm64 (4 `tar.gz` + 2 `zip`); `SHA256SUMS` **6 lines == 6
artifacts**; `sha256sum -c` → **6/6 OK, rc 0**; 3 independent digest recomputes → **3/3 MATCH** (includes
the linux `tar.gz` archive and the non-linux `darwin_arm64.tar.gz`). **PASS.**

### 3b — non-vacuity: checksum controls

```
### host: x280 2026-10-15T17:42:30+02:00   cwd: /home/c08r4d0r/wz-tests-t_0509c320-p3/negctl-run
### files copied from /home/c08r4d0r/worktrees/wz-e2e-gate-r3/dist

### CONTROL A (positive): unmodified copies must verify OK, rc 0
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
REAL_RC=0   # want 0

### CONTROL B (negative): corrupt 1 byte at offset 1024 of net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz — EXPLICIT path
$ printf 'X' | dd of=net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz bs=1 seek=1024 count=1 conv=notrunc
1+0 records in
1+0 records out
1 byte copied, 2.806e-05 s, 35.6 kB/s
dd_rc=0   new_size=4492115 bytes (unchanged size = in-place write, not truncating)

$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: FAILED
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: OK
sha256sum: WARNING: 1 computed checksum did NOT match
REAL_RC=1   # want non-zero

### CONTROL C (negative): a missing artifact must be reported FAILED open or read
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz: FAILED
net4sats-wizard_0.5.1-snapshot.0b5322a_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip: OK
sha256sum: net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: No such file or directory
net4sats-wizard_0.5.1-snapshot.0b5322a_windows_arm64.zip: FAILED open or read
sha256sum: WARNING: 1 listed file could not be read
sha256sum: WARNING: 1 computed checksum did NOT match
REAL_RC=1   # want non-zero

### CONCLUSION: -c reports OK on unmodified (A, rc 0) and FAILED on corrupted (B, rc 1) and
###             missing (C, rc 1), so the 6/6 OK run in 05b-artifacts-and-checksums.txt is not vacuous.
```

Control A (positive, rc 0) reports 6/6 OK; control B (1 byte corrupted in an explicit-path copy) flips
`linux_amd64.tar.gz` to `FAILED`, rc 1; control C (artifact removed) reports `FAILED open or read`, rc 1.
So the §3 `6/6 OK` is not vacuous. **PASS.**

---

## STEP 4 — binary smoke (`--version`) and manual injection

```
### host: x280 c08r4d0r 2026-10-15T17:42:34+02:00
### scratch: /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke

# extract the built linux/amd64 binary from the snapshot archive (not a fresh build)
$ tar -xzf /home/c08r4d0r/worktrees/wz-e2e-gate-r3/dist/net4sats-wizard_*_linux_amd64.tar.gz
tar_rc=0
$ ls -la
total 11148
drwxrwxr-x    2 c08r4d0r c08r4d0r     4096 Oct 15 17:42 .
drwx------ 3367 c08r4d0r c08r4d0r   135168 Oct 15 17:42 ..
-rw-rw-r--    1 c08r4d0r c08r4d0r     1575 Oct 15 17:35 CHANGELOG.md
-rw-rw-r--    1 c08r4d0r c08r4d0r     8264 Oct 15 17:35 README.md
-rwxr-xr-x    1 c08r4d0r c08r4d0r 11251864 Oct 15 17:40 net4sats-wizard

$ file ./net4sats-wizard
./net4sats-wizard: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=MmKS5_eawR9iar_ccpT8/sGCA5M6nz2bZqVA7rUcg/Dyl7ftR2QVYukp4NS1s0/ItFfTY-2v-I6FP3GmoIW, stripped

$ ./net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.0b5322a
exit_code=0

$ ./net4sats-wizard --version | sed -n 's/^tollgate-wizard //p'   # injected version extraction
0.5.1-snapshot.0b5322a

# is the printed version 'dev'? (want: no)
FOUND_DEV=no (good)
```

```
### cwd: /home/c08r4d0r/worktrees/wz-e2e-gate-r3   HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
# manual injection: two DIFFERENT values must change the printed line

$ go build -buildvcs=false -ldflags "-X main.version=v9.9.9-manual-p3" -o /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/inject-v9.9.9-manual-p3 .
build_rc=0
$ /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/inject-v9.9.9-manual-p3 --version
tollgate-wizard v9.9.9-manual-p3
exit_code=0

$ go build -buildvcs=false -ldflags "-X main.version=v0.0.0-second-different-value-p3" -o /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/inject-v0.0.0-second-different-value-p3 .
build_rc=0
$ /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/inject-v0.0.0-second-different-value-p3 --version
tollgate-wizard v0.0.0-second-different-value-p3
exit_code=0

# control: no -ldflags -> the zero value must be 'dev'
$ go build -buildvcs=false -o /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/plain-build .
build_rc=0
$ /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/plain-build --version
tollgate-wizard dev
exit_code=0

# the -X target must exist in the source (stripped release binary has no symbol table)
$ go tool nm /home/c08r4d0r/.cache/hermes/tmp/wz-p3-smoke/plain-build | grep -E 'main\.version|main\.versionString'
  eb5b60 D main.version

# where the version var and the --version flag live in source
$ grep -rn 'version' --include=*.go . | grep -v _test.go | head -20
./main.go:44:// version is the wizard build version. Release builds inject it at link time:
./main.go:46://	go build -ldflags "-X main.version=v0.4.0" .
./main.go:52:var version = "dev"
./main.go:54:// programName is the CLI identity reported by --version.
./main.go:57:// versionString is the single line --version prints: the program name followed
./main.go:58:// by the (possibly injected) build version.
./main.go:59:func versionString() string {
./main.go:60:	return programName + " " + version
./main.go:192:// Some versions prefix with "Cell 01 - Address: ..." instead of the interface name.
./main.go:212:		// "Cell" prefix (some iwinfo versions) or a new ESSID line indicates a new block
./main.go:350:	//      interface. This is what worked in the July version.
./main.go:405:			// on AP interfaces on many OpenWrt versions (e.g. GL-MT3000
./main.go:637:	// --version is answered before the server binds a port, so release
./main.go:641:		case "--version", "-version":
./main.go:642:			fmt.Println(versionString())
./deploy.go:194:		{Name: "firmware", Desc: "Checking firmware version...", Status: "pending"},
./deploy.go:228:	fwOut := sshRun(client, "cat /etc/openwrt_release 2>/dev/null || cat /etc/openwrt_version 2>/dev/null || echo 'not openwrt'")
./deploy.go:244:	versionLine := ""
./deploy.go:249:				versionLine = strings.Trim(parts[1], "'")
./deploy.go:253:	job.addLog("Firmware: " + versionLine)

# git tree still clean after all builds?
$ git status --porcelain
?? PROGRESS.md
porcelain_rc=0
```

```
$ go build -ldflags "-X main.version=v9.9.9-plain-no-buildvcs-flag" -o /tmp/plainbg .
PLAIN_BUILD_RC=0
tollgate-wizard v9.9.9-plain-no-buildvcs-flag
exit_code=0
```

The linux/amd64 binary **extracted from the snapshot archive** prints `tollgate-wizard 0.5.1-snapshot.0b5322a`
(exit 0) — the injected version, not `dev`. Two manual `-X main.version=…` builds printed two *different*
injected values, and a no-ldflags control printed `tollgate-wizard dev`, so the injection demonstrably
changes. `go tool nm` locates the `-X` target (`D main.version`) in `main.go:52`, consumed by
`versionString()` at `main.go:59` and printed for `--version`/`-version` at `main.go:641`. A plain build
without `-buildvcs=false` also succeeds in this worktree (rc 0), so the worktree VCS-stamp pitfall does
not apply here. **PASS.**

---

## STEP 5 — GitHub Actions proof

```
### host: x280 c08r4d0r 2026-10-15T17:42:49+02:00
### gh: gh version 2.45.0 (2025-07-18 Ubuntu 2.45.0-1ubuntu0.3)   account: felixfelix-bot

# ---- the LITERAL command from the task, against the repo the task names ----
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
literal_rc=1

$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)
repoview_rc=0

$ git ls-remote https://github.com/net4sats/tollgate-wizard-go
remote: Repository not found.
fatal: repository 'https://github.com/net4sats/tollgate-wizard-go/' not found
lsremote_rc=0

# ---- the real repo (alias) and its upstream workflows ----
$ gh run list --repo net4sats/net4sats-wizard-go --limit 5
upstream_runs_rc=0   # empty output = no runs

$ gh workflow list --repo net4sats/net4sats-wizard-go
workflow_list_rc=0

$ gh api repos/net4sats/net4sats-wizard-go/contents/.github
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":"404"}gh: Not Found (HTTP 404)
api_github_dir_rc=0

$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch
prchecks_rc=1

$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json statusCheckRollup,headRefOid
{"headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","statusCheckRollup":[]}

# ---- the fork: where this workflow CAN run (upstream main has no .github/) ----
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540632745	2m6s	2026-10-14T09:12:20Z
completed	cancelled	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38540629420	4s	2026-10-14T09:12:19Z
fork_runs_rc=0

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 5 --json databaseId,conclusion,status,headSha,headBranch,event,workflowName,createdAt,url
[{"conclusion":"success","createdAt":"2026-10-15T15:37:06Z","databaseId":38761913548,"event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548","workflowName":"CI"},{"conclusion":"success","createdAt":"2026-10-15T12:52:40Z","databaseId":38738030458,"event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38738030458","workflowName":"CI"},{"conclusion":"success","createdAt":"2026-10-15T03:41:31Z","databaseId":38677575488,"event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38677575488","workflowName":"CI"},{"conclusion":"success","createdAt":"2026-10-14T09:12:20Z","databaseId":38540632745,"event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38540632745","workflowName":"CI"},{"conclusion":"cancelled","createdAt":"2026-10-14T09:12:19Z","databaseId":38540629420,"event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38540629420","workflowName":"CI"}]
fork_runs_json_rc=0
```

Real exit codes (the labels above for `repoview`/`lsremote`/`api` were the rc of a downstream `head` in a
pipe; the values below are the commands' own):

```
### host: x280 2026-10-15T17:43:05+02:00
# CORRECTION: in 07-gh-run-list-and-detail.txt the rc labels for repoview/lsremote/api were the
# rc of the downstream `head` in the pipe, not of gh/git. Real exit codes below (no pipes).

$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
LITERAL_RC=1

$ gh repo view net4sats/tollgate-wizard-go >/dev/null
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)
REPOVIEW_RC=1

$ git ls-remote https://github.com/net4sats/tollgate-wizard-go >/dev/null
remote: Repository not found.
fatal: repository 'https://github.com/net4sats/tollgate-wizard-go/' not found
LSREMOTE_RC=128

$ gh api repos/net4sats/net4sats-wizard-go/contents/.github >/dev/null
gh: Not Found (HTTP 404)
API_GITHUB_DIR_RC=1

$ gh run list --repo net4sats/net4sats-wizard-go --limit 5   (upstream: no runs)
UPSTREAM_RUNS_RC=0

$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch
PRCHECKS_RC=1

=========== run 38761913548 (dispatched by THIS pass) ===========
$ gh run view 38761913548 --repo felixfelix-bot/net4sats-wizard-go --json status,conclusion,headSha,headBranch,event,url,workflowName
{"conclusion":"success","event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548","workflowName":"CI"}
RUNVIEW_RC=0

$ gh run view 38761913548 --repo felixfelix-bot/net4sats-wizard-go   # per-job

✓ ci/add-github-actions-workflow-clean CI #29 · 38761913548
Triggered via workflow_dispatch about 6 minutes ago

JOBS
✓ test in 1m49s (ID 116515106517)
✓ lint in 21s (ID 116515106827)

ANNOTATIONS
! Node.js 20 is deprecated. The following actions target Node.js 20 but are being forced to run on Node.js 24: actions/checkout@v4, actions/setup-go@v5. For more information see: https://github.blog/changelog/2025-09-19-deprecation-of-node-20-on-github-actions-runners/
test: .github#2

- "The ubuntu-latest label will migrate to Ubuntu 26 beginning October 19, 2026. For more information, see https://github.com/actions/runner-images/issues/14748"
test: .github#1

! Node.js 20 is deprecated. The following actions target Node.js 20 but are being forced to run on Node.js 24: actions/checkout@v4, actions/setup-go@v5, golangci/golangci-lint-action@v6. For more information see: https://github.blog/changelog/2025-09-19-deprecation-of-node-20-on-github-actions-runners/
lint: .github#2

- "The ubuntu-latest label will migrate to Ubuntu 26 beginning October 19, 2026. For more information, see https://github.com/actions/runner-images/issues/14748"
lint: .github#1


For more information about a job, try: gh run view --job=<job-id>
View this run on GitHub: https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548
RUNVIEW_TABLE_RC=0

$ gh api repos/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548/jobs --jq '.jobs[] | {name,conclusion,status,started_at,completed_at}'
{"completed_at":"2026-10-15T15:38:59Z","conclusion":"success","name":"test","started_at":"2026-10-15T15:37:10Z","status":"completed"}
{"completed_at":"2026-10-15T15:37:32Z","conclusion":"success","name":"lint","started_at":"2026-10-15T15:37:11Z","status":"completed"}
JOBS_RC=0
```

Job-log proof that the workflow really ran its gates (not skipped, not `continue-on-error`):

```
### job-log proof for run 38761913548 — the gate commands the workflow actually ran
1:test	Set up job	﻿2026-10-15T15:37:11.2885103Z Current runner version: '2.337.0'
183:test	Run go build ./...	﻿2026-10-15T15:37:20.8907796Z ##[group]Run go build ./...
184:test	Run go build ./...	2026-10-15T15:37:20.8908282Z [36;1mgo build ./...[0m
185:test	Run go build ./...	2026-10-15T15:37:20.9103527Z shell: /usr/bin/bash -e {0}
186:test	Run go build ./...	2026-10-15T15:37:20.9103839Z ##[endgroup]
187:test	Run go vet ./...	﻿2026-10-15T15:37:21.4463001Z ##[group]Run go vet ./...
188:test	Run go vet ./...	2026-10-15T15:37:21.4463751Z [36;1mgo vet ./...[0m
189:test	Run go vet ./...	2026-10-15T15:37:21.4499400Z shell: /usr/bin/bash -e {0}
190:test	Run go vet ./...	2026-10-15T15:37:21.4499717Z ##[endgroup]
191:test	Run go test ./... -race -count=1	﻿2026-10-15T15:37:22.6639986Z ##[group]Run go test ./... -race -count=1
192:test	Run go test ./... -race -count=1	2026-10-15T15:37:22.6640390Z [36;1mgo test ./... -race -count=1[0m
193:test	Run go test ./... -race -count=1	2026-10-15T15:37:22.6676224Z shell: /usr/bin/bash -e {0}
194:test	Run go test ./... -race -count=1	2026-10-15T15:37:22.6676569Z ##[endgroup]
195:test	Run go test ./... -race -count=1	2026-10-15T15:38:56.3285998Z ok  	net4sats-wizard	82.110s
220:lint	Set up job	﻿2026-10-15T15:37:11.9666877Z Current runner version: '2.337.0'
252:lint	Set up job	2026-10-15T15:37:13.0421027Z Download action repository 'golangci/golangci-lint-action@v6' (SHA:55c2c1448f86e01eaae002a5a3a9624417608d84)
404:lint	Run golangci/golangci-lint-action@v6	﻿2026-10-15T15:37:26.3647471Z ##[group]Run golangci/golangci-lint-action@v6
405:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3647841Z with:
406:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3648048Z   version: v1.61.0
407:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3648276Z   install-mode: binary
408:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3650932Z   github-token: ***
409:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3651154Z   verify: true
410:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3651376Z   only-new-issues: false
411:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3651612Z   skip-cache: false
412:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3651827Z   skip-save-cache: false
413:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3652064Z   problem-matchers: false
414:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3652324Z   cache-invalidation-interval: 7
415:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.3652584Z ##[endgroup]
416:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.5040848Z ##[group]prepare environment
417:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.5046041Z Checking for go.mod: go.mod
418:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.5051660Z (node:2376) [DEP0040] DeprecationWarning: The `punycode` module is deprecated. Please use a userland alternative instead.
419:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.5052627Z (Use `node --trace-deprecation ...` to show where the warning was created)
420:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.7753929Z Cache hit for: golangci-lint.cache-Linux-2963-ec4850e71b02b9858515893ba8635d17a34933de
421:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:26.7876551Z (node:2376) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
422:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.3376772Z Received 169855 of 169855 (100.0%), 0.6 MBs/sec
423:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.3377460Z Cache Size: ~0 MB (169855 B)
424:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.3406072Z [command]/usr/bin/tar -xf /home/runner/work/_temp/4dcdbf3b-dc16-48c2-a6c0-6ff09892feda/cache.tzst -P -C /home/runner/work/net4sats-wizard-go/net4sats-wizard-go --use-compress-program unzstd
425:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.3992958Z Cache restored successfully
426:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.3997852Z Restored cache for golangci-lint from key 'golangci-lint.cache-Linux-2963-ec4850e71b02b9858515893ba8635d17a34933de' in 895ms
427:lint	Run golangci/golangci-lint-action@v6	2026-10-15T15:37:27.4000069Z Finding needed golangci-lint version...
LOG_RC=0
```

`gh run list --repo net4sats/tollgate-wizard-go --limit 5` — the literal command — fails with
`HTTP 404: Not Found`: **that repo does not exist**. The real repo `net4sats/net4sats-wizard-go` has **no
workflows and no runs at all** (`main` has no `.github/`), so `gh pr checks 35` reports nothing and
`statusCheckRollup` is `[]` — a property of the *base* repo, not of this branch. The green run is
therefore **38761913548** on the fork: `conclusion=success`, `status=completed`,
`headSha=0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff` (exactly the pushed head), `headBranch=ci/add-github-actions-workflow-clean`, jobs
`✓ test in 1m49s` and `✓ lint in 21s`. Its log shows `go build ./...`, `go vet ./...`,
`go test ./... -race -count=1` → `ok  net4sats-wizard 82.110s`, and golangci-lint `version: v1.61.0`.

No workflow was edited, no check weakened, nothing was force-pushed: CI is green on the pushed SHA.
**PASS.**

---

## STEP 6 — card transition + this pass's own report

```
### card transition attempts (verbatim, 2026-10-15T17:43:31+02:00)
$ HERMES_KANBAN_BOARD=net4sats-mvp-v2 hermes kanban show t_0509c320 | head -4
Task t_0509c320: End-to-end verify: snapshot build, --version smoke, full suite, push
  status:    triage
  assignee:  worker-heavy
  workspace: scratch @ /home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/workspaces/t_bd4576f8

$ hermes kanban complete t_0509c320
cannot complete t_0509c320 (unknown id or terminal state)
complete_rc=1

$ hermes kanban promote t_0509c320 --force
cannot promote t_0509c320: task t_0509c320 is 'triage'; promote only applies to 'todo' or 'blocked'
promote_rc=1
```

The card is in status `triage`; a worker session cannot complete or promote it (both attempts refused
verbatim above). Evidence is therefore attached + commented on the card, and **the manager must promote
it** out of `triage`.

### Ad-hoc verification of this pass's own harness

`rev3-90-adhoc-verify.py` re-derives every claim above from the capture files (not from this prose):
**62 checks / 0 failures / `ADHOC_VERDICT=PASS`**. Its negative control (`rev3-91-adhoc-control.sh`)
tampers a **copy** of the captures — flips `TEST_RC=0→1`, renames one `MATCH` to `NOMATCH`, flips the
corrupted-artifact control back to rc 0, changes the run conclusion `success→failure`, and flips
`promote_rc=1→0` — and the same verifier then reports **5 failures**, so the checks are not vacuous. The
originals are re-checked afterwards and still pass. `bash -n` is clean on every `rev3-*.sh` script.

```
=== 1) verifier on the REAL captures (must PASS) ===
== STEP 1: clean tree + push ==
  PASS  01: rev-list count is 0\t0 (ahead/behind)
  PASS  01: push was a no-op (branch already pushed)
  PASS  01: remote ref == verified head SHA
  PASS  01: porcelain rc 0
== STEP 2a: vet / gofmt / module shape ==
  PASS  02a: CI-matching toolchain go1.22.12 in use
  PASS  02a: single module listed (./... covers the suite)
  PASS  02a: go vet rc 0
  PASS  02a: gofmt rc 0 (no formatting drift)
  PASS  02a: go list rc 0
== STEP 2b: go test -race ==
  PASS  03: race suite rc 0
  PASS  03: package reported ok
  PASS  04: 123 tests RUN
  PASS  04: 0 FAIL
  PASS  04: 0 SKIP
  PASS  04: 0 data races
== STEP 2c/2d: lint + goreleaser check ==
  PASS  02c: golangci-lint rc 0 (v1.61.0)
  PASS  02c: the version CI pins
  PASS  02d: goreleaser check rc 0
== STEP 3: snapshot + checksums ==
  PASS  05: goreleaser release --snapshot --clean rc 0
  PASS  05: snapshot version derived from the head SHA
  PASS  05b: 6 archives
  PASS  05b: 2 linux archives
  PASS  05b: 2 windows archives
  PASS  05b: 2 darwin archives
  PASS  05b: SHA256SUMS has 6 lines (== 6 artifacts)
  PASS  05b: sha256sum -c real rc 0 (captured without a pipe)
  PASS  05b: sha256sum -c reports 6 OK and no FAILED
  PASS  05b: contains no control lines (those live in 08)
  PASS  05b: 3 independent recomputes MATCH (>=1 archive, >=1 non-linux)
== STEP 3b: checksum controls (non-vacuity) ==
  PASS  08: CONTROL A positive rc 0
  PASS  08: CONTROL B corrupted -> rc 1
  PASS  08: B flags the corrupted artifact as FAILED
  PASS  08: C flags the missing artifact
  PASS  08: no 'CONTROL_B_RC' piped-status artifact remains
== STEP 4: --version smoke + injection ==
  PASS  06: --version prints 'tollgate-wizard 0.5.1-snapshot.0b5322a'
  PASS  06: exit code 0
  PASS  06: printed version is NOT 'dev'
  PASS  06: binary came from the linux/amd64 archive (ELF x86-64)
  PASS  06b: manual injection value #1 changed the line
  PASS  06b: manual injection value #2 changed the line
  PASS  06b: no-ldflags control prints dev
  PASS  06b: the -X target symbol exists in the source
  PASS  06c: plain go build (no -buildvcs=false) rc 0
== STEP 5: GHA proof ==
  PASS  07: literal card repo returns HTTP 404
  PASS  09: literal gh run list rc 1 (real code, no pipe)
  PASS  09: gh repo view unreadable rc 1
  PASS  09: git ls-remote rc 128
  PASS  09: upstream .github 404
  PASS  09: gh pr checks 35 reports nothing
  PASS  09: pr checks empty (upstream runs no workflows)
  PASS  09: own dispatched run concluded success
  PASS  09: run's headSha == the pushed head SHA
  PASS  09: both jobs (test, lint) success
  PASS  09: job table shows test and lint green
  PASS  09b: CI really ran the race suite
  PASS  09b: CI really ran go vet
  PASS  09b: CI really ran go build
  PASS  09b: CI race suite reported ok
  PASS  09b: CI lint used the pinned golangci-lint v1.61.0
== STEP 6: card transitions ==
  PASS  10: complete refused (triage)
  PASS  10: promote --force refused (triage)
  PASS  10: refusal reason recorded verbatim

checks=62  failures=0
ADHOC_VERDICT=PASS
REAL_RUN_RC=0

=== 2) verifier on the TAMPERED copy (must FAIL) ===
== STEP 1: clean tree + push ==
  PASS  01: rev-list count is 0\t0 (ahead/behind)
  PASS  01: push was a no-op (branch already pushed)
  PASS  01: remote ref == verified head SHA
  PASS  01: porcelain rc 0
== STEP 2a: vet / gofmt / module shape ==
  PASS  02a: CI-matching toolchain go1.22.12 in use
  PASS  02a: single module listed (./... covers the suite)
  PASS  02a: go vet rc 0
  PASS  02a: gofmt rc 0 (no formatting drift)
  PASS  02a: go list rc 0
== STEP 2b: go test -race ==
  FAIL  03: race suite rc 0   TEST_RC = 1, want 0
  PASS  03: package reported ok
  PASS  04: 123 tests RUN
  PASS  04: 0 FAIL
  PASS  04: 0 SKIP
  PASS  04: 0 data races
== STEP 2c/2d: lint + goreleaser check ==
  PASS  02c: golangci-lint rc 0 (v1.61.0)
  PASS  02c: the version CI pins
  PASS  02d: goreleaser check rc 0
== STEP 3: snapshot + checksums ==
  PASS  05: goreleaser release --snapshot --clean rc 0
  PASS  05: snapshot version derived from the head SHA
  PASS  05b: 6 archives
  PASS  05b: 2 linux archives
  PASS  05b: 2 windows archives
  PASS  05b: 2 darwin archives
  PASS  05b: SHA256SUMS has 6 lines (== 6 artifacts)
  PASS  05b: sha256sum -c real rc 0 (captured without a pipe)
  PASS  05b: sha256sum -c reports 6 OK and no FAILED
  PASS  05b: contains no control lines (those live in 08)
  FAIL  05b: 3 independent recomputes MATCH (>=1 archive, >=1 non-linux)   ['net4sats-wizard_0.5.1-snapshot.0b5322a_darwin_arm64.tar.gz', 'net4sats-wizard_0.5.1-snapshot.0b5322a_windows_amd64.zip']
== STEP 3b: checksum controls (non-vacuity) ==
  PASS  08: CONTROL A positive rc 0
  FAIL  08: CONTROL B corrupted -> rc 1   
  PASS  08: B flags the corrupted artifact as FAILED
  PASS  08: C flags the missing artifact
  PASS  08: no 'CONTROL_B_RC' piped-status artifact remains
== STEP 4: --version smoke + injection ==
  PASS  06: --version prints 'tollgate-wizard 0.5.1-snapshot.0b5322a'
  PASS  06: exit code 0
  PASS  06: printed version is NOT 'dev'
  PASS  06: binary came from the linux/amd64 archive (ELF x86-64)
  PASS  06b: manual injection value #1 changed the line
  PASS  06b: manual injection value #2 changed the line
  PASS  06b: no-ldflags control prints dev
  PASS  06b: the -X target symbol exists in the source
  PASS  06c: plain go build (no -buildvcs=false) rc 0
== STEP 5: GHA proof ==
  PASS  07: literal card repo returns HTTP 404
  PASS  09: literal gh run list rc 1 (real code, no pipe)
  PASS  09: gh repo view unreadable rc 1
  PASS  09: git ls-remote rc 128
  PASS  09: upstream .github 404
  PASS  09: gh pr checks 35 reports nothing
  PASS  09: pr checks empty (upstream runs no workflows)
  FAIL  09: own dispatched run concluded success   
  PASS  09: run's headSha == the pushed head SHA
  PASS  09: both jobs (test, lint) success
  PASS  09: job table shows test and lint green
  PASS  09b: CI really ran the race suite
  PASS  09b: CI really ran go vet
  PASS  09b: CI really ran go build
  PASS  09b: CI race suite reported ok
  PASS  09b: CI lint used the pinned golangci-lint v1.61.0
== STEP 6: card transitions ==
  PASS  10: complete refused (triage)
  FAIL  10: promote --force refused (triage)   promote_rc = 0, want 1
  PASS  10: refusal reason recorded verbatim

checks=62  failures=5
  - 03: race suite rc 0
  - 05b: 3 independent recomputes MATCH (>=1 archive, >=1 non-linux)
  - 08: CONTROL B corrupted -> rc 1
  - 09: own dispatched run concluded success
  - 10: promote --force refused (triage)
TAMPERED_RUN_RC=1   # want non-zero

=== diff of the tampering (real -> tampered) ===
--- 03-test-race.txt
10c10
< TEST_RC=0
---
> TEST_RC=1
--- 05b-artifacts-and-checksums.txt
51c51
< MATCH    net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz  77cefac426e777045af6ce34b9b3aabb1ca242801982bcfc421a324d225b7507
---
> NOMATCH    net4sats-wizard_0.5.1-snapshot.0b5322a_linux_amd64.tar.gz  77cefac426e777045af6ce34b9b3aabb1ca242801982bcfc421a324d225b7507
--- 08-checksum-controls.txt
29c29
< REAL_RC=1   # want non-zero
---
> REAL_RC=0   # want non-zero
42c42
< REAL_RC=1   # want non-zero
---
> REAL_RC=0   # want non-zero
--- 09-gh-run-detail-and-real-rc.txt
31c31
< {"conclusion":"success","event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548","workflowName":"CI"}
---
> {"conclusion":"failure","event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38761913548","workflowName":"CI"}
62,63c62,63
< {"completed_at":"2026-10-15T15:38:59Z","conclusion":"success","name":"test","started_at":"2026-10-15T15:37:10Z","status":"completed"}
< {"completed_at":"2026-10-15T15:37:32Z","conclusion":"success","name":"lint","started_at":"2026-10-15T15:37:11Z","status":"completed"}
---
> {"completed_at":"2026-10-15T15:38:59Z","conclusion":"failure","name":"test","started_at":"2026-10-15T15:37:10Z","status":"completed"}
> {"completed_at":"2026-10-15T15:37:32Z","conclusion":"failure","name":"lint","started_at":"2026-10-15T15:37:11Z","status":"completed"}
--- 10-card-transition-attempts.txt
14c14
< promote_rc=1
---
> promote_rc=0

=== 3) bash -n on every rev3 script ===
OK   /home/c08r4d0r/wz16-tools/rev3-01-clean-tree-and-push.sh
OK   /home/c08r4d0r/wz16-tools/rev3-02-vet-lint-check.sh
OK   /home/c08r4d0r/wz16-tools/rev3-02b-test-race.sh
OK   /home/c08r4d0r/wz16-tools/rev3-03-snapshot.sh
OK   /home/c08r4d0r/wz16-tools/rev3-03b-negative-controls.sh
OK   /home/c08r4d0r/wz16-tools/rev3-03c-inventory-and-checksums.sh
OK   /home/c08r4d0r/wz16-tools/rev3-04-binary-smoke.sh
OK   /home/c08r4d0r/wz16-tools/rev3-05-gha-proof.sh
OK   /home/c08r4d0r/wz16-tools/rev3-05b-gha-detail.sh
OK   /home/c08r4d0r/wz16-tools/rev3-91-adhoc-control.sh

=== 4) the originals are unmodified: real run again (must PASS) ===

checks=62  failures=0
ADHOC_VERDICT=PASS
REAL_RUN_2_RC=0
```

---

## Corrections made during this pass (disclosed, not hidden)

Both were caught by this pass's own verifier and fixed in the artifacts, not papered over:

1. **The first negative control was vacuous.** In the first version of the artifact capture, the
   corrupting `dd of=<glob>` did **not** expand the glob — `dd` created a literal file named
   `net4sats-wizard_*_linux_amd64.tar.gz` and left the real archive intact, so the control printed
   `OK` and `NEGCTL_RC=0`. The checksum itself was never in doubt (the `6/6 OK` run and the independent
   recomputes were always against the real artifacts). Fixed by using the explicit filename; the flawed
   first attempt is preserved verbatim as
   `05b-FIRST-ATTEMPT-broken-controls.txt`, and the corrected controls are `08-checksum-controls.txt`.
2. **Two status labels were false markers.** `sha256sum -c … | grep …; echo $?` and
   `gh repo view … | head; echo $?` report the *pipe's* rc, not the command's. Fixed by redirecting to a
   file / `/dev/null` and reading the real rc (`SHA256SUM_C_REAL_RC`, `REPOVIEW_RC=1`, `LSREMOTE_RC=128`).
   The `head`/`grep`-derived labels remain in the earlier capture files next to their corrected
   counterparts, and this note is why.
3. **`05b` used to mix the real run with the controls**, which made it ambiguous for a reviewer. It was
   regenerated to contain only the real run (`SHA256SUM_C_REAL_RC=0`); the verifier now asserts that.

## What this pass did NOT change

No source, test, workflow, or release-config file was touched: every gate is green at `0b5322a` as
received. The commits on this pass's branch are documentation only (`PROGRESS.md`, `REPORT.md`).

---

## STEP 9 — post-push final state (all ends), observed SHAs

### post-push final state   host: x280 2026-10-15T17:46:07+02:00

=== A) the docs branch on the fork ===
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 status --porcelain
porcelain_rc=0   # empty = nothing uncommitted
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 rev-parse HEAD  /  rev-parse origin/worker-heavy/t_0509c320-r3
141d57fafdf884c64ce23923bb541dd181918312
141d57fafdf884c64ce23923bb541dd181918312
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 rev-list --left-right --count HEAD...origin/worker-heavy/t_0509c320-r3
0	0

=== B) the branch UNDER TEST is unchanged ===
$ git -C /home/c08r4d0r/wz16-verify status --porcelain
porcelain_rc=0
$ git -C /home/c08r4d0r/wz16-verify rev-parse HEAD
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ git -C /home/c08r4d0r/wz16-verify ls-remote origin refs/heads/ci/add-github-actions-workflow-clean | cut -f1
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json state,headRefOid
{"headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","state":"OPEN"}

=== C) lane repo: is this pass's commit on the remote main? ===
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 merge-base --is-ancestor 9acadc59 origin/main
ancestor_rc=0   # 0 = my commit is in origin/main
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 ls-tree -r --name-only origin/main | grep -c verify/t_0509c320-pass3/
32
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 log --oneline -3 origin/main
22ab6da9 docs(demo-meta): record run-3 git outcome (push verified both ends; shared-tree sibling lanes noted)
9acadc59 verify(t_0509c320): pass 3 — full integration gate re-executed; evidence + reports landed
f9b752bd docs(demo-meta): run-3 intro_card_meta assemble+validate - both validators pass, parent-gate null reconciled, renderer consumes artifact

=== D) CI at the pushed head SHA, re-read now ===
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z

=== E) the card comment, read back from the board ===
27:Comments (5):
29:  [2026-10-15 05:49] worker-heavy: ## t_0509c320 — End-to-end verify: VERDICT **PASS** (all gates green, CI green at the pushed head SHA)
65:  [2026-10-15 14:57] worker-heavy: ## t_0509c320 — pass 2: end-to-end verify — VERDICT **PASS** (every gate re-executed live; CI green at the pushed head SHA)
140:  [2026-10-15 17:45] worker-heavy: ## t_0509c320 — End-to-end verify (pass 3): VERDICT **PASS** — all gates re-executed live, CI green at the pushed head SHA
$ hermes kanban show t_0509c320 | grep -c 'pass 3'
1

### Observed push SHAs (quote these, do not re-derive)

| What | SHA | Remote read-back |
|---|---|---|
| Docs branch for this pass (`PROGRESS.md` + `REPORT.md`) | `141d57fafdf884c64ce23923bb541dd181918312` | `origin/worker-heavy/t_0509c320-r3` = `141d57f…`, ahead/behind `0 0` |
| Lane board repo evidence commit | `9acadc59` (`verify(t_0509c320): pass 3 — …`) | ancestor of `origin/main` (rc 0); 32 files of `verify/t_0509c320-pass3/` present |
| Branch under test (UNCHANGED) | `0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff` | `origin/ci/add-github-actions-workflow-clean` = same; PR #35 `state=OPEN`, `headRefOid` = same |
| CI run for that SHA | `38761913548` | `conclusion=success`, `headSha=0b5322a…` |

A peer session on this shared lane tree committed `22ab6da9` on top of `9acadc59` immediately after the
push; `9acadc59` is an ancestor of `origin/main`, so this pass's evidence is intact. Card comment for
this pass read back from the board at `2026-10-15 17:45` (5 comments total on the card). Attachments 92
(`REPORT-t_0509c320-pass3.md`, 47708 bytes) and 93 (`PROGRESS-t_0509c320-pass3.md`).

Capture: `12-final-state.txt`.

---

## STEP 10 — gate-evidence supplement for the fleet gate engine

### 10a — `secrets_clean`: attested with the gate's own scanner

### scan_secrets from /home/c08r4d0r/.hermes/scripts/gate_engine.py
### targets: 40 files (this pass's evidence only; not the whole repo)
  clean lane-evidence:00-worktree-registration.txt
  clean lane-evidence:01-clean-tree-and-push.txt
  clean lane-evidence:02a-vet-gofmt-module-shape.txt
  clean lane-evidence:02c-golangci-lint.txt
  clean lane-evidence:02d-goreleaser-check.txt
  clean lane-evidence:03-test-race.txt
  clean lane-evidence:04-test-race-summary.txt
  clean lane-evidence:04-test-race-verbose.log
  clean lane-evidence:05-goreleaser-snapshot.txt
  clean lane-evidence:05b-FIRST-ATTEMPT-broken-controls.txt
  clean lane-evidence:05b-artifacts-and-checksums.txt
  clean lane-evidence:06-version-smoke.txt
  clean lane-evidence:06b-version-injection-change.txt
  clean lane-evidence:06c-plain-build-no-flag.txt
  clean lane-evidence:07-gh-run-list-and-detail.txt
  clean lane-evidence:08-checksum-controls.txt
  clean lane-evidence:09-gh-run-detail-and-real-rc.txt
  clean lane-evidence:09b-run-log-gate-lines.txt
  clean lane-evidence:10-card-transition-attempts.txt
  clean lane-evidence:11-adhoc-verifier-run.txt
  clean lane-evidence:12-final-state.txt
  clean lane-evidence:rev3-01-clean-tree-and-push.sh
  clean lane-evidence:rev3-02-vet-lint-check.sh
  clean lane-evidence:rev3-02b-test-race.sh
  clean lane-evidence:rev3-03-snapshot.sh
  clean lane-evidence:rev3-03b-negative-controls.sh
  clean lane-evidence:rev3-03c-inventory-and-checksums.sh
  clean lane-evidence:rev3-04-binary-smoke.sh
  clean lane-evidence:rev3-05-gha-proof.sh
  clean lane-evidence:rev3-05b-gha-detail.sh
  clean lane-evidence:rev3-12-final-state.sh
  clean lane-evidence:rev3-90-adhoc-verify.py
  clean lane-evidence:rev3-91-adhoc-control.sh
  clean lane-evidence:rev3-92-make-report.py
  clean lane-evidence:rev3-card-comment.md
  clean lane:PROGRESS.t_0509c320-pass3.md
  clean lane:REPORT.t_0509c320-pass3.md
  clean worktree docs (this pass):PROGRESS.md
  clean worktree docs (this pass):REPORT.md
  clean card-comment:rev3-card-comment.md

files_scanned=40  total_hits=0
secret-scan: clean (regex fallback)

gitleaks is not installed on x280 (`which gitleaks` → `command not found`), so the attestation uses the
`(regex fallback)` spelling, which `gate_engine.RE_SECRET_TOOL` accepts. The scanner is
`gate_engine.scan_secrets` **imported from the engine itself** (not a re-implementation), so the gate's
independent re-scan runs the same detectors. 40 evidence files, 0 hits, `ATTEST_RC=0`.

### 10b — `ci_evidence`: the engine's verifier reads ngit CI, not GitHub Actions

```
$ python3 ~/.hermes/scripts/ngit_ci_evidence.py felixfelix-bot/net4sats-wizard-go \
    --commit 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff --json --require-conclusion --timeout 60
warning: nak req exited 2 but returned events; continuing
no kind-9842 CI results found for repo=felixfelix-bot/net4sats-wizard-go commit=0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff (408 event(s) scanned from wss://relay.ngit.dev, wss://gitnostr.com)
[]
CI_EVIDENCE_RC=0

$ python3 ~/.hermes/scripts/gate_engine.py evaluate --board net4sats-mvp-v2 --task t_0509c320
net4sats-mvp-v2/t_0509c320: tier=code verdict=block missing=['ci_evidence', 'cold_cross_family_review', 'review_published', 'consolidated', 'secrets_clean', 'pr_branch_naming', 'no_live_drift']
gate_engine_rc=0
```

`ngit_ci_evidence.py` reads **kind-9842 CI-result events from Nostr relays** ("no kind-9842 CI results
found … 408 event(s) scanned"); this repository's CI runs on **GitHub Actions**, which publishes no such
events. The gate's `ci_evidence` predicate is therefore unsatisfiable for this card as configured — that
is a **predicate/tooling gap, not missing CI**, since the GitHub run is green at the exact head SHA (see
§5 and the run re-read in §9). I deliberately did **not** post the canonical
`ci_evidence repo=… head=… ref=…` line as satisfied: the engine parses that line and then invokes the
verifier, which returns `no_results`, so claiming it would assert something the tool cannot confirm.

**Decision for the manager** (not a worker action): either apply an audited scope override / GitHub-CI
evidence form for this card, or require CI results to be published as kind-9842 events (ngit CI) before
`ci_evidence` can pass here.

### 10c — what remains genuinely outside a verification pass

From the engine's own evaluation, `missing` = `ci_evidence, cold_cross_family_review, review_published,
consolidated, secrets_clean, pr_branch_naming, no_live_drift`. `secrets_clean` is now attested (§10a);
`cold_cross_family_review` / `review_published` / `consolidated` belong to the review pipeline;
`no_live_drift` and `pr_branch_naming` are properties of the reviewed PR (its branch is
`ci/add-github-actions-workflow-clean`, not a `pr/<slug>` name). None of these is producible by a
verification pass, and the card still needs a manager promotion out of `triage`.

Capture: `13-secrets-attestation.txt`. Posted as a second comment on t_0509c320.

---

## FINAL — all four ends re-checked after every commit of this pass

### FINAL state   host: x280 c08r4d0r 2026-10-15T17:49:02+02:00

=== END 1: docs branch on the fork (this pass) ===
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 log --oneline -4
f992c9d docs(t_0509c320): pass 3 — append §10 gate-evidence supplement (secrets attestation + ci_evidence gap)
abc5bab docs(t_0509c320): pass 3 — record the measured gate verdicts in PROGRESS.md
ae49ed0 docs(t_0509c320): pass 3 — append §9 post-push final state (observed SHAs)
141d57f verify(t_0509c320): pass 3 — integration gate re-executed at 0b5322a (all gates green)
$ local HEAD vs origin/worker-heavy/t_0509c320-r3
f992c9d9eeb129b8397435dad884ea4ae806de1a
f992c9d9eeb129b8397435dad884ea4ae806de1a
0	0
$ git -C /home/c08r4d0r/worktrees/wz-e2e-gate-r3 status --porcelain   # want empty
porcelain_rc=0

=== END 2: lane repo main (this pass's evidence) ===
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 log --oneline -3 origin/main
623fc3ec docs(t_0509c320): pass 3 — append §10 gate-evidence supplement (secrets attestation + ci_evidence gap)
2553cd76 docs(t_0509c320): pass 3 — record the measured gate verdicts in PROGRESS.md
d2f6159e docs(t_0509c320): pass 3 — append §9 post-push final state (observed SHAs)
  IN origin/main: 9acadc59
  IN origin/main: d2f6159e
  IN origin/main: 2553cd76
  IN origin/main: 623fc3ec
$ evidence files on origin/main: 38

=== END 3: branch UNDER TEST (must be unchanged) ===
$ git -C /home/c08r4d0r/wz16-verify status --porcelain
porcelain_rc=0
$ git -C /home/c08r4d0r/wz16-verify rev-parse HEAD
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ remote ref
0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
$ gh pr view 35 ... --json state,headRefOid
{"headRefOid":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","state":"OPEN"}

=== END 4: CI for that exact SHA + the card ===
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 2
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
$ gh run view 38761913548 --json conclusion,headSha,status
{"conclusion":"success","headSha":"0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff","status":"completed"}
$ hermes kanban show t_0509c320 | head -3
Task t_0509c320: End-to-end verify: snapshot build, --version smoke, full suite, push
  status:    triage
  assignee:  worker-heavy
$ hermes kanban show t_0509c320 | grep -c 'pass 3'   # pass-3 comments on the card
1
$ gate_engine evaluate
net4sats-mvp-v2/t_0509c320: tier=code verdict=block missing=['ci_evidence', 'cold_cross_family_review', 'review_published', 'consolidated', 'pr_branch_naming', 'no_live_drift']

### SHA ledger for this pass (final, observed)

| End | Content | SHA | Read-back |
|---|---|---|---|
| Fork branch `worker-heavy/t_0509c320-r3` | `PROGRESS.md` (docs) | `141d57f` | pushed |
| | `REPORT.md` §9 post-push state | `ae49ed0` | pushed |
| | `PROGRESS.md` measured verdicts | `abc5bab` | pushed |
| | `REPORT.md` §10 gate supplement | `f992c9d` | **local HEAD == remote, ahead/behind `0 0`** |
| Lane repo `main` | evidence + reports + captures | `9acadc59` → `d2f6159e` → `2553cd76` → `623fc3ec` | all four verified ancestors of `origin/main`; 38 files under `verify/t_0509c320-pass3/` |
| Branch under test (UNCHANGED) | `ci/add-github-actions-workflow-clean` | `0b5322a` | checkout clean; PR #35 `OPEN`, `headRefOid` identical |
| CI for that SHA | run `38761913548` | — | `conclusion=success`, `headSha=0b5322a`, `status=completed` |

**The one gate item that moved:** `gate_engine evaluate` now reports
`missing=['ci_evidence', 'cold_cross_family_review', 'review_published', 'consolidated', 'pr_branch_naming', 'no_live_drift']`
— **`secrets_clean` is no longer missing**, i.e. the engine independently re-scanned this pass's evidence
and accepted the attestation in §10a. The remaining items are the review pipeline's (`cold_cross_family_review`,
`review_published`, `consolidated`), the PR's branch naming, live-drift, and the ngit-only `ci_evidence`
predicate explained in §10b. The card is still `triage` and needs a manager promotion.

Capture: `14-final-state-after-commits.txt`.
