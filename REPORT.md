# REPORT — t_10a913ec pass-7 (fleet offload: heavy/deepseek-v4-flash)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.

**Verdict: PASS — all four card steps executed and re-verified live in THIS pass.
Zero test failures. Zero wizard source changes required. Build succeeds.**

| | |
|---|---|
| Host | x280 (home `/home/c08r4d0r`, `whoami` = `c08r4d0r`) |
| Worktree (actual) | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` |
| Branch tip at start | `a6805a4` (a sibling session advanced it to `b939329` mid-pass; tree stayed clean) |
| Source head under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (module `net4sats-wizard`, PUBLIC) |
| Toolchain | go1.22.12 linux/amd64 (`/usr/local/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-heavypass/` |
| Date | 2026-10-16 |

Every figure below was re-observed live in THIS pass. Prior passes' reports are preserved in this
branch's git history; nothing is carried forward unverified.

### Which commit is under test

The branch tip at start was `a6805a4`, docs-only. No Go source changed above `0b5322a`:

```
$ git log -1 --format='%h %s' -- '*.go'
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs

$ git log --oneline 0b5322a..HEAD -- '*.go'
(empty)
```

So the source under test is byte-identical to `0b5322a` — the vendored/version-bumped head with the
hardening intact.

### Card premises that are stale on this host (disclosed, not silently worked around)

1. **`WORKTREE: /home/c03rad0r/worktrees/wizard-fix` does not exist.** This host's home is
   `/home/c08r4d0r`; `/home/c03rad0r` is the *other* fleet node's home. The worktree is registered
   at `/home/c08r4d0r/worktrees/wizard-fix` (`git worktree list` confirms). That is where this ran.

2. **"there is NO `--version` flag — ANY invocation starts a server" is stale at this head.**
   `--version`/`-version` exists (`main.go:637-641`) and is answered *before* the listener binds;
   `TestVersionFlagInjectedViaLdflags` in the suite asserts it. Port-8099 hygiene was performed
   regardless, and the smoke test used a **bare (no-arg) invocation** so the real server path was
   exercised — exactly what the card asks for.

3. **STEP 4's "pre-existing `go vet` IPv6 warning + gofmt drift" do not exist at this head.**
   `go vet ./...` exits 0 with empty output; `gofmt -l .` prints zero lines. Nothing was "fixed" —
   and there was nothing pre-existing to leave untouched.

---

## STEP 1 — Server hygiene (`ss :8099` BEFORE)

```
### STEP1: pkill -x net4sats-wizard
pkill rc=1
--- pgrep -a -x net4sats-wizard ---
pgrep rc=1
--- ss -ltnp (all) ---
State  Recv-Q Send-Q      Local Address:Port  Peer Address:PortProcess
LISTEN 0      5                 0.0.0.0:9100       0.0.0.0:*    users:(("python3",pid=1243,fd=3))
LISTEN 0      4096              0.0.0.0:631        0.0.0.0:*
LISTEN 0      5                 0.0.0.0:8899       0.0.0.0:*    users:(("python3",pid=807504,fd=3))
LISTEN 0      4096              0.0.0.0:22         0.0.0.0:*
LISTEN 0      128               0.0.0.0:8443       0.0.0.0:*
LISTEN 0      4096        127.0.0.53%lo:53         0.0.0.0:*
LISTEN 0      4096              0.0.0.0:18081      0.0.0.0:*
LISTEN 0      511             127.0.0.1:8898       0.0.0.0:*    users:(("node-MainThread",pid=820004,fd=21))
LISTEN 0      1024            127.0.0.1:8388       0.0.0.0:*
LISTEN 0      4096           127.0.0.54:53         0.0.0.0:*
LISTEN 0      5                       *:9099             *:*    users:(("socat",pid=1244,fd=5))
LISTEN 0      4096                 [::]:631           [::]:*
LISTEN 0      1024                [::1]:8388          [::]:*
LISTEN 0      4096                 [::]:22            [::]:*
LISTEN 0      50     [::ffff:127.0.0.1]:8080             *:*    users:(("java",pid=1245,fd=66))
--- ss -ltnp | grep 8099 ---
grep8099 rc=1
```

`pkill -x net4sats-wizard` rc=1 and `pgrep -a -x` rc=1 — no wizard process existed to kill.
`:8099` had **no listener** (grep rc=1 = FREE). Raw: `wz-captures-heavypass/01-server-hygiene-before.txt`.

---

## STEP 2 — `go build -buildvcs=false`

```
### go version
go version go1.22.12 linux/amd64
### go.mod
module net4sats-wizard

go 1.22

require golang.org/x/crypto v0.24.0

require golang.org/x/sys v0.21.0 // indirect
### STEP2: go build -buildvcs=false
BUILD rc=0
### artifact
-rwxrwxr-x 1 c08r4d0r c08r4d0r 14257829 Oct 16 14:22 net4sats-wizard
445934e6ec8ef5ed69b6d0e419eb948c822aaf127ad78dee7e59d1c962b891a9  net4sats-wizard
### git status --porcelain after build
 M PROGRESS.md
```

**rc=0, no compiler output.** Binary `net4sats-wizard`, 14,257,829 bytes, sha256
`445934e6ec8ef5ed69b6d0e419eb948c822aaf127ad78dee7e59d1c962b891a9`.

Raw: `wz-captures-heavypass/02-build.txt`.

Note on sha provenance: this pass built with go1.22.12 (the toolchain on `PATH` at
`/usr/local/go/bin`), whereas earlier passes' reports cite go1.24.4 from
`/home/c08r4d0r/wz19-audit/go/bin/go` — hence a different binary sha for the same source. Both
toolchains satisfy `go.mod`'s `go 1.22`. The source is unchanged; only the compiler differs.

---

## STEP 3 — FULL `go test ./...` (zero failures)

Non-vacuity check first — this repo is a single module, so `./...` covers everything:

```
$ find . -name go.mod -not -path '*/.git/*' | wc -l
1
$ go list ./...
net4sats-wizard
```

Command and tail:

```
### STEP3: go test -count=1 -v ./...
...
=== RUN   TestVersionFlagInjectedViaLdflags
--- PASS: TestVersionFlagInjectedViaLdflags (0.54s)
PASS
ok  	net4sats-wizard	66.226s
TEST rc=0
```

Head of the same raw run (proves real tests executed, not a no-op):

```
=== RUN   TestHandleScan
--- PASS: TestHandleScan (28.99s)
=== RUN   TestHandleIndex
--- PASS: TestHandleIndex (0.00s)
=== RUN   TestDeployRequestValidation
=== RUN   TestDeployRequestValidation/missing_IP_and_password
=== RUN   TestDeployRequestValidation/empty_password_(fresh_reset)
=== RUN   TestDeployRequestValidation/missing_IP
=== RUN   TestDeployRequestValidation/invalid_lightning_address
=== RUN   TestDeployRequestValidation/valid_lightning_address
=== RUN   TestDeployRequestValidation/valid_lnurl
--- PASS: TestDeployRequestValidation (0.00s)
...
=== RUN   TestArtifactPinsCoverDownloadedURLConstants
--- PASS: TestArtifactPinsCoverDownloadedURLConstants (0.00s)
=== RUN   TestVerifyArtifactPinBenignDownload
```

Tally computed from the full raw file (`wz-captures-heavypass/03-test-verbose.txt`):

| metric | value |
|---|---|
| `=== RUN   Test…` (incl. subtests) | **123** |
| `--- PASS` (top-level) | **52** |
| `--- FAIL` | **0** |
| `--- SKIP` | **0** |
| package result | `ok  net4sats-wizard  66.226s` |
| exit code | **0** |

Grep for failure markers returned nothing:

```
$ grep -nE '^(FAIL|--- FAIL|panic)' 03-test-verbose.txt
(no output)
```

**Zero failures, zero skips.** No test failed, so there was no vendored-asset/version-change
breakage to fix. The hardening was **not** reverted (no source file was touched at all).

### Vendored assets really are compiled in

`embed.go` carries `//go:embed all:portal` and `all:admin`. The vendored-asset pins are asserted by
the suite itself — `TestArtifactPinsCoverDownloadedURLConstants`, `TestVerifyArtifactPinBenignDownload`
and `TestVersionFlagInjectedViaLdflags` all PASS — and the version-injection path is exercised by
`TestVersionStringFormat` + `TestVersionFlagInjectedViaLdflags`.

---

## One-shot server smoke — and `ss :8099` AFTER

Bare (no-arg) invocation, so the real server path is exercised:

```
### one-shot server (BARE invocation, no --version)
--- pgrep -a -x net4sats-wizard ---
1345858 ./net4sats-wizard
pgrep rc=0
--- ss -ltnp | grep 8099 ---
LISTEN 0      4096            127.0.0.1:8099       0.0.0.0:*    users:(("net4sats-wizard",pid=1345858,fd=3))
grep rc=0
recorded PID=1345858
--- GET / ---
http_code=200 size=19976
--- server log ---
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).
--- kill by recorded PID only ---
kill rc=0
--- ss -ltnp | grep 8099 (after kill) ---
grep rc=1
--- pkill -x net4sats-wizard final sweep ---
pkill rc=1
--- FINAL ss -ltnp | grep 8099 ---
final grep rc=1  (1 = FREE)
```

Bound `127.0.0.1:8099` (pid 1345858; `ss` shows that exact PID owns the port), `GET /` → HTTP 200
(19,976 B), then killed **by recorded PID only** and re-proved free.

Raw: `wz-captures-heavypass/04-smoke-after.txt`.

**`:8099` is free after the work** (final grep rc=1, and the post-kill sweep found no wizard).

---

## STEP 4 — vet / gofmt left untouched, tree clean

```
### go vet ./...
VET rc=0
### gofmt -l . (empty = clean)
GOFMT_LINES=0
### git status --porcelain
(end)
### git status -sb
## worker-heavy/t_10a913ec...origin/worker-heavy/t_10a913ec
### HEAD
b93932986a8775a7362594e845390c2a8ed05567
### go source above HEAD?
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
```

`go vet ./...` rc=0 with **empty** output; `gofmt -l .` **0 lines**; `git status --porcelain`
**empty**; branch in sync with its remote.

Raw: `wz-captures-heavypass/05-gates.txt`.

Note: the worktree is **shared with a concurrent sibling session** (`hermes -p worker-heavy` running
the "squash history into atomic conventional commits" task, visible in `ps`). That session committed
`b939329` (docs-only) **during this pass**, taking HEAD from `a6805a4` → `b939329` and occupying the
pathspec I had modified in PROGRESS.md. No Go file was touched by it. The build/test evidence above
is unaffected in substance: `git log -1 --format='%h %s' -- '*.go'` still resolves to `0b5322a`, so
the compiled+tested source is exactly what earlier passes verified.

---

## Commit / tree state

- **No fix commit was required** (zero Go source changes, zero test failures).
- `git status --porcelain` empty at end of the verification (STEP 4 block).
- Sibling session committed docs-only `b939329` mid-pass (PROGRESS.md + REPORT.md).
- This pass appends its own PROGRESS.md/REPORT.md update and pushes (see final section).

## REMAINING steps

None — all of STEP 1..4 are complete and verified in this pass.
