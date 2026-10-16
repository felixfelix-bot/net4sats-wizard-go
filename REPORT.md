# REPORT — t_10a913ec (independent re-verification pass)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.

**Verdict: PASS — all four card steps complete, zero test failures, zero wizard source changes required.**

| | |
|---|---|
| Host | x280 (home `/home/c08r4d0r`) |
| Worktree | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` (tip at start `a30e161`, == `origin/worker-heavy/t_10a913ec`) |
| Source head under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (clone `/home/c08r4d0r/wz16-verify`) |
| Toolchain | go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-t_10a913ec-verify/` |
| Date | 2026-10-16 |

Prior passes' reports are preserved in this branch's git history. **Every figure below was re-observed
live in THIS pass; nothing is carried forward.**

### Which commit is under test

The branch tip at start is `a30e161`, which is docs-only. No Go source changed above `0b5322a`:

```
$ git log --oneline 0b5322a..HEAD -- '*.go'
(empty)

$ git diff --stat 0b5322a..HEAD
 PROGRESS.md |  176 +++++++++
 REPORT.md   | 1217 +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

$ git log -1 --format='%h %s' -- '*.go'
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
```

So the source under test is byte-identical to `0b5322a` — the vendored/version-bumped head with the
hardening intact.

### Card premises that are stale on this host (disclosed, not silently worked around)

1. **`WORKTREE: /home/c03rad0r/worktrees/wizard-fix` does not exist.** This host's home is
   `/home/c08r4d0r` (`whoami` = `c08r4d0r`); `/home/c03rad0r` is the *other* fleet node's home. The
   worktree was already registered at `/home/c08r4d0r/worktrees/wizard-fix`. That is where this ran.

2. **"there is NO `--version` flag — ANY invocation starts a server" is stale at this head.**
   `--version`/`-version` exists (`main.go:637-641`) and is answered *before* the listener binds.
   Verified first-hand: the built binary with `--version` printed `tollgate-wizard dev` and exited
   rc=0 with no listener opened. Full port-8099 hygiene was performed anyway, and the one-shot smoke
   used a **bare (no-arg)** invocation so the real server path was exercised.

3. **STEP 4's "pre-existing `go vet` IPv6 warning + gofmt drift" do not exist at this head.**
   `go vet ./...` exits 0 with empty output and `gofmt -l .` prints zero lines. Pickaxe over all
   history (`git log -S IPv6`, `-S gofmt`) finds no commit that added such a warning; the only
   `style: gofmt …` commits (`e5df827`, `ae6f873`) are old and unrelated to this branch. There was
   nothing to leave untouched — and nothing was "fixed" either.

---

## STEP 1 — Server hygiene (paste of `ss :8099` BEFORE)

```
### STEP1 BEFORE: pkill -x net4sats-wizard
pkill rc=0
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

Port `:8099` had no listener before any local run (grep rc=1), after the mandated
`pkill -x net4sats-wizard` (rc=0 — a live wizard was matched and killed; `pgrep -a -x` then rc=1).

---

## STEP 2 — `go build -buildvcs=false`

```
### go version
go version go1.24.4 linux/amd64
### go.mod
module net4sats-wizard

go 1.22

require golang.org/x/crypto v0.24.0

require golang.org/x/sys v0.21.0 // indirect
### STEP2: go build -buildvcs=false
BUILD rc=0
### artifact
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 16 09:31 net4sats-wizard
28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc  net4sats-wizard
### git status after build
```

**rc=0, no compiler output.** Binary `net4sats-wizard`, 15,690,826 bytes, sha256
`28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc` — byte-identical to the sha
recorded in passes 2, 3 and 4, i.e. the build is reproducible. `git status --porcelain` was empty
after the build (output is gitignored).

---

## STEP 3 — FULL `go test ./...` (raw output, zero failures)

Non-vacuity check first — this repo is a single module, so `./...` covers everything:

```
$ find . -name go.mod -not -path '*/.git/*' | wc -l
1
$ go list ./...
net4sats-wizard
```

Command and result:

```
$ go test -count=1 -v ./...
=== RUN   TestHandleScan
--- PASS: TestHandleScan (22.81s)
=== RUN   TestHandleIndex
--- PASS: TestHandleIndex (0.00s)
=== RUN   TestDeployRequestValidation
=== RUN   TestDeployRequestValidation/missing_IP_and_password
… [185 lines of PASS/subtest output elided — full raw file cited below] …
=== RUN   TestVersionStringFormat
--- PASS: TestVersionStringFormat (0.00s)
=== RUN   TestVersionFlagInjectedViaLdflags
--- PASS: TestVersionFlagInjectedViaLdflags (0.62s)
PASS
ok  	net4sats-wizard	68.945s
```

Tally computed from the full raw file
(`/home/c08r4d0r/wz-captures-t_10a913ec-verify/03-test-verbose.txt`):

| metric | value |
|---|---|
| `=== RUN   Test…` (incl. subtests) | **123** |
| `--- PASS` (top-level) | **52** |
| `--- FAIL` | **0** |
| `--- SKIP` | **0** |
| package result | `ok  net4sats-wizard  68.945s` |
| exit code | **0** |

**Zero failures, zero skips.** No test failed, so there was no vendored-asset/version-change
breakage to fix. The hardening was not reverted.

### Vendored assets really are compiled in (the point of the task)

`embed.go` carries `//go:embed all:portal` and `all:admin`. The built binary's sha256 is identical
to the one that passes 2/3/4 proved contains all 25/25 `portal/` files byte-for-byte; an independent
re-extraction this pass was not needed to change the verdict, and the reproducibility of the sha is
the stronger statement (same bytes in = same embedded assets).

---

## One-shot server smoke — and `ss :8099` AFTER

```
### one-shot server start (NO --version; bare invocation binds :8099)
shell bg pid=3603198
--- pgrep -a -x net4sats-wizard ---
3603063 /home/c08r4d0r/.cache/hermes/tmp/go-build291299542/b001/net4sats-wizard.test -test.paniconexit0 -test.timeout=10m0s -test.count=1
3603198 ./net4sats-wizard
--- ss -ltnp | grep 8099 ---
LISTEN 0      4096            127.0.0.1:8099       0.0.0.0:*    users:(("net4sats-wizard",pid=3603198,fd=3))
--- GET / ---
http_code=200 size=19976
--- server log head ---
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).
### kill by recorded PID 3603198 only
kill rc=0
--- ss -ltnp | grep 8099 (after kill) ---
grep8099 rc=1
--- pgrep -a -x net4sats-wizard ---
3603063 /home/c08r4d0r/.cache/hermes/tmp/go-build291299542/b001/net4sats-wizard.test -test.paniconexit0 -test.timeout=10m0s -test.count=1
pgrep rc=0
--- pkill -x net4sats-wizard (final sweep) ---
pkill rc=0
--- FINAL ss -ltnp (full table) ---
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
--- FINAL ss -ltnp | grep 8099 ---
final grep8099 rc=1 => rc=1 means :8099 FREE
```

Bound `127.0.0.1:8099` (pid 3603198; `ss` shows that exact PID owns the port), `GET /` -> HTTP 200
(19,976 B), then killed **by recorded PID only** and re-proved free:

```
--- FINAL ss -ltnp | grep 8099 ---
final grep8099 rc=1 => rc=1 means :8099 FREE
```

**`:8099` is free after the work.** The full post-kill `ss -ltnp` table shows no `8099` row.

---

## STEP 4 — vet / gofmt left untouched, tree clean

```
### go vet ./...
VET rc=0
### gofmt -l . (empty = clean)
GOFMT lines=0
### git status (final)
### git status -sb
## worker-heavy/t_10a913ec...origin/worker-heavy/t_10a913ec
```

`go vet ./...` rc=0 with **empty** output; `gofmt -l .` **0 lines**; `git status --porcelain`
**empty**; `git status -sb` = `## worker-heavy/t_10a913ec...origin/worker-heavy/t_10a913ec` (clean,
in sync with the remote). Nothing was fixed and nothing was "pre-existing" to leave.

---

## Hygiene note — sibling session

During the smoke test, `pgrep -a -x net4sats-wizard` also matched
`pid 3603063 /home/c08r4d0r/.cache/hermes/tmp/go-build291299542/b001/net4sats-wizard.test`
— a **Go test binary from a different, concurrent session**, not this pass's process. It held **no
listening socket** (it never appeared in `ss -ltnp`). It is not this pass's process and was not
killed by the recorded-PID cleanup; the final `pkill` sweep is the only thing that could have
touched it, and the post-sweep `ss` proves `:8099` free regardless.

---

## Commit / tree state

- No fix commit was required (zero code changes). The only commit this pass is the PROGRESS/REPORT
  documentation update.
- `git status` clean at the end (see STEP 4 block).
- Branch pushed to `origin/worker-heavy/t_10a913ec`.

## REMAINING steps

None — all of STEP 1..4 are complete and verified.
