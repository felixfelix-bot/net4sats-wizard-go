# REPORT — t_10a913ec (re-dispatch pass 2)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.
**Verdict: PASS — all four card steps complete, zero test failures, zero wizard source changes required.**

| | |
|---|---|
| Host | x280 (`/home/c08r4d0r`) |
| Worktree | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` |
| HEAD under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* (= `origin/ci/add-github-actions-workflow-clean`, the newest wizard branch head; the commits above it on this branch are docs only) |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (clone `/home/c08r4d0r/wz16-verify`) |
| Toolchain | go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-redispatch/` |
| Date | 2026-10-15 |

This is pass 2 of an identical earlier dispatch. Pass-1's report is preserved in this branch's git
history (`a242f4b`, `a645c56`, `cb1db08`). **Every figure below was re-observed live in this pass**;
nothing is copied forward.

---

## STEP 1 — Server hygiene (WIZARD QUIRKS)

### 1a. `pkill -x net4sats-wizard` BEFORE any local run

```
$ pkill -x net4sats-wizard; echo "pkill_rc=$?"
pkill_rc=1          # rc=1 == no matching process was running
```

### 1b. `ss -ltnp` BEFORE (full output — no :8099 listener)

```
$ ss -ltnp
State  Recv-Q Send-Q      Local Address:Port  Peer Address:PortProcess
LISTEN 0      5                 0.0.0.0:9100       0.0.0.0:*    users:(("python3",pid=1243,fd=3))
LISTEN 0      4096              0.0.0.0:631        0.0.0.0:*
LISTEN 0      4096            127.0.0.1:18445      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=27))
LISTEN 0      4096              0.0.0.0:22         0.0.0.0:*
LISTEN 0      128               0.0.0.0:8443       0.0.0.0:*
LISTEN 0      4096        127.0.0.53%lo:53         0.0.0.0:*
LISTEN 0      1024            127.0.0.1:8388       0.0.0.0:*
LISTEN 0      5                 0.0.0.0:18543      0.0.0.0:*    users:(("python3",pid=4040268,fd=3))
LISTEN 0      128               0.0.0.0:18443      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=12))
LISTEN 0      4096              0.0.0.0:18444      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=28))
LISTEN 0      4096           127.0.0.54:53         0.0.0.0:*
LISTEN 0      5                       *:9099             *:*    users:(("socat",pid=1244,fd=5))
LISTEN 0      4096                 [::]:631           [::]:*
LISTEN 0      1024                [::1]:8388          [::]:*
LISTEN 0      4096                 [::]:22            [::]:*
LISTEN 0      50     [::ffff:127.0.0.1]:8080             *:*    users:(("java",pid=1245,fd=66))
LISTEN 0      4096                 [::]:18444         [::]:*    users:(("bitcoind",pid=3498175,fd=26))

$ ss -ltnp | grep -E ':(8099|8109)' || echo "NO LISTENER ON :8099 or :8109 -> PORT FREE"
NO LISTENER ON :8099 or :8109 -> PORT FREE
```

### 1c. One-shot server smoke (bare invocation), then killed

The card says "there is NO --version flag — ANY invocation starts a server". **That premise is stale
at this HEAD** — see STEP 1d. Both behaviours were exercised:

```
$ ./net4sats-wizard --version
tollgate-wizard dev
version_rc=0                       # exits immediately, binds nothing (main.go:637-645)

$ ./net4sats-wizard                 # bare invocation = one-shot server (backgrounded)
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).

$ ss -ltnp | grep -E ':8099'
LISTEN 0      4096            127.0.0.1:8099       0.0.0.0:*    users:(("net4sats-wizard",pid=123270,fd=3))

$ curl -s -o /dev/null -w 'GET / -> HTTP %{http_code}\n' http://127.0.0.1:8099/
GET / -> HTTP 200
$ curl -s -o /dev/null -w 'GET /api/status -> HTTP %{http_code}\n' http://127.0.0.1:8099/api/status
GET /api/status -> HTTP 301

# killed by recorded PID (process kill -> pid 123270 gone)
```

### 1d. STEP 1d — `pkill -x net4sats-wizard` AFTER work, `ss -ltnp` AFTER

```
$ pkill -x net4sats-wizard; echo "pkill_rc=$?"
pkill_rc=1
$ ps aux | grep -E 'net4sats-wizard' | grep -v grep || echo "no net4sats-wizard process"
no net4sats-wizard process

$ ss -ltnp
State  Recv-Q Send-Q      Local Address:Port  Peer Address:PortProcess
LISTEN 0      5                 0.0.0.0:9100       0.0.0.0:*    users:(("python3",pid=1243,fd=3))
LISTEN 0      4096              0.0.0.0:631        0.0.0.0:*
LISTEN 0      4096            127.0.0.1:18445      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=27))
LISTEN 0      4096              0.0.0.0:22         0.0.0.0:*
LISTEN 0      128               0.0.0.0:8443       0.0.0.0:*
LISTEN 0      4096        127.0.0.53%lo:53         0.0.0.0:*
LISTEN 0      1024            127.0.0.1:8388       0.0.0.0:*
LISTEN 0      5                 0.0.0.0:18543      0.0.0.0:*    users:(("python3",pid=4040268,fd=3))
LISTEN 0      128               0.0.0.0:18443      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=12))
LISTEN 0      4096              0.0.0.0:18444      0.0.0.0:*    users:(("bitcoind",pid=3498175,fd=28))
LISTEN 0      4096           127.0.0.54:53         0.0.0.0:*
LISTEN 0      5                       *:9099             *:*    users:(("socat",pid=1244,fd=5))
LISTEN 0      4096                 [::]:631           [::]:*
LISTEN 0      1024                [::1]:8388          [::]:*
LISTEN 0      4096                 [::]:22            [::]:*
LISTEN 0      50     [::ffff:127.0.0.1]:8080             *:*    users:(("java",pid=1245,fd=66))
LISTEN 0      4096                 [::]:18444         [::]:*    users:(("bitcoind",pid=3498175,fd=26))

$ ss -ltnp | grep -E ':(8099|8109)' || echo "NO LISTENER ON :8099 or :8109 -> PORT FREE"
NO LISTENER ON :8099 or :8109 -> PORT FREE
```

**:8099 free before AND after. No wizard process left behind.**

---

## STEP 2 — `go build -buildvcs=false`

```
$ go build -buildvcs=false
$ echo "BUILD_RC=$?"
BUILD_RC=0
```

No compiler output at all (clean build). Artifact:

```
$ ls -l net4sats-wizard
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 15 05:52 net4sats-wizard

$ sha256sum net4sats-wizard
28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc  net4sats-wizard
```

`net4sats-wizard` is listed in `.gitignore`, so the build does not dirty the tree.

---

## STEP 3 — Full `go test ./...`

Package set is a single package:

```
$ go list ./...
net4sats-wizard
```

### Raw output, uncached (`-count=1`)

```
$ go test ./... -count=1
ok  	net4sats-wizard	91.116s
$ echo $?
0
```

(The first run in this pass reported `ok net4sats-wizard (cached)`; it was re-run with `-count=1`
so the evidence is a real execution, not a cache hit.)

### Verbose run — pass/fail accounting

```
$ go test ./... -count=1 -v > step3-gotest-verbose.txt 2>&1; echo $?
0

$ grep -c '^=== RUN'   step3-gotest-verbose.txt
123
$ grep -c '^--- PASS'  step3-gotest-verbose.txt
52
$ grep -c '^--- FAIL'  step3-gotest-verbose.txt
0
$ grep -c '^--- SKIP'  step3-gotest-verbose.txt
0
$ grep -E '^--- (FAIL|SKIP)|^(FAIL|panic)' step3-gotest-verbose.txt || echo NONE
NONE

$ tail -3 step3-gotest-verbose.txt
--- PASS: TestVersionFlagInjectedViaLdflags (0.50s)
PASS
ok  	net4sats-wizard	84.275s
```

**123 RUN / 52 top-level PASS / 0 FAIL / 0 SKIP / rc=0.** (52 `--- PASS` lines are the top-level
test functions; the 123 `=== RUN` lines include their subtests.)

Full list of the 52 passing top-level tests:

```
TestHandleScan                                  TestJobStructFields
TestHandleIndex                                 TestDeployStepStruct
TestDeployRequestValidation                     TestLogEntryStruct
TestDeployRequestDefaults                       TestConcurrentJobOperations
TestStepInitialization                          TestTruncate
TestArtifactPinsCoverDownloadedURLConstants     TestNewJob
TestVerifyArtifactPinBenignDownload             TestAddLog
TestVerifyArtifactPinTamperedDownload           TestSetStep
TestVerifyArtifactPinUnknownURLFailsClosed      TestWriteError
TestVerifyArtifactPinHex                        TestTcpProbe
TestExtractIPKFilenameStrict                    TestJobMutexSafety
TestClamp                                       TestTollgatePkgURLPinsExistingAsset
TestSshRunFunction                              TestDeployGoPinRegistry
TestDiscoverRoutersNoNetwork                    TestPinnedURLsAreLive
TestEmptyJobOperations                          TestRemoteCommandsNeverInterpolatePayloads
TestDownloadURLsPointToTollgateRepo             TestVulnerableConstructionsAreGone
TestTollgatePkgURLAssetName                     TestBenignValuesRoundTripThroughCarriers
TestScanFailedHeuristic                         TestPasswdCommandInertUnderSh
TestAllRadiosUp                                 TestStaCarriersDecodeUnderSh
TestIfaceUp                                     TestLnIdentCommandJQRoundTrip
TestHTTPGetFile                                 TestCfgConfigCommandJQRoundTrip
TestStaSetupScript                              TestListenAddress
TestHandleStatus                                TestCORSMiddleware
TestHandleDeploy                                TestValidLightningAddressUnchangedOnBenign
TestMainFunctionStructure                       TestValidLightningAddress
                                                TestVersionStringFormat
                                                TestVersionFlagInjectedViaLdflags
```

### Failures caused by the vendored asset / version change

**None.** Zero failures, zero skips — so no wizard fix was required and nothing had to be
re-run. The hardening from `cc9f94c` (root shell-injection sinks closed, loopback bind,
localhost-only CORS) is intact and its guards are exactly what these tests assert
(`TestVulnerableConstructionsAreGone`, `TestRemoteCommandsNeverInterpolatePayloads`,
`TestPasswdCommandInertUnderSh`, `TestListenAddress`, `TestCORSMiddleware`).

### The vendored asset really is in the binary (the point of the task)

`embed.go` declares `//go:embed all:portal`. Independent check — every file under `portal/` was
read from disk and searched for byte-for-byte inside the freshly built `net4sats-wizard`:

```
portal files: 25/25 present byte-for-byte in net4sats-wizard (15690826 byte binary)
missing: none
```

The bundle the embedded HTML references (`portal/index.html` → `assets/splash-Be2Pt9NX.js`) is
embedded, and that file is the one vendored by commit `b436512`
(*chore: vendor portal v1.0.10-alpha (pricing block) + bump configwizURL (#31)*) — i.e. the
vendored asset bump and the version bump both build and pass.

---

## STEP 4 — Pre-existing `go vet` IPv6 warning + gofmt drift: left untouched

There is **nothing to leave untouched at this HEAD** — both are already clean, because the lint
commit `7ceade3` (*chore(lint): pin the standard golangci-lint set and fix every finding it
reports*) removed them. Verified, not assumed, and no file was modified:

```
$ go vet ./...
$ echo "GOVET_RC=$?"
GOVET_RC=0                 # empty output, rc=0

$ gofmt -l .
$ echo "GOFMT_RC=$?"
GOFMT_RC=0                 # empty output == no unformatted file
```

No `.go`, asset, or config file was edited in this pass. **No fix commit was needed**, so the only
commits added are the two documentation files below.

---

## Commits and final state

Docs only (no code change):

```
$ git log --oneline -3
<this pass> docs(t_10a913ec) pass-2: independent re-verification — build + 123/123 suite green at 0b5322a
0b5322a    feat(release): link-time version injection, --version, goreleaser config and release docs
```

```
$ git status --porcelain
                       # empty == clean tree
$ git rev-list --left-right --count HEAD...origin/worker-heavy/t_10a913ec
0	0                  # HEAD == pushed remote head
```

Pushed to `origin` = fork `felixfelix-bot/net4sats-wizard-go`, branch `worker-heavy/t_10a913ec`.

---

## Acceptance checklist

- [x] `ss -ltnp` before (no :8099) **and** after (no :8099) pasted
- [x] `go build -buildvcs=false` output pasted — rc=0, clean
- [x] Full `go test ./...` output pasted — `ok net4sats-wizard 91.116s`, 123 RUN / 0 FAIL / 0 SKIP
- [x] `git status` clean
- [x] Vendored asset (`portal/`, 25/25 files byte-for-byte) proven embedded in the built binary
- [x] Hardening not reverted, not touched; `go vet` + `gofmt` left as found

## Findings for the operator

1. **Card worktree path is for the other fleet node.** `/home/c03rad0r/worktrees/wizard-fix` does not
   exist on x280; the worktree is `/home/c08r4d0r/worktrees/wizard-fix`. Same class as the
   `/home/c03rad0r` vs `/home/c08r4d0r` mix-up recorded in pass 1.
2. **STEP 1's premise is outdated.** `--version`/`-version` now exists (`main.go:641`) and is answered
   *before* the listener binds, so a version check is safe at this HEAD. The port-8099 hygiene
   guidance is still correct for bare invocations and was followed in full.
3. **STEP 4's premise is outdated too.** No `go vet` IPv6 warning and no gofmt drift remain at
   `0b5322a`; the lint commit already cleared them. Nothing was "fixed" because nothing was there.
