# REPORT — t_10a913ec: Build wizard and run full go test suite green

Board: net4sats-mvp-v2 · Assignee: worker-heavy · Author: worker-heavy (fleet offload)
Date: 2026-10-14 · Host: x280 · Worktree: `/home/c08r4d0r/worktrees/wizard-fix`
Branch: `worker-heavy/t_10a913ec` (fork `felixfelix-bot/net4sats-wizard-go`)
HEAD under test: **`0b5322a`** — `feat(release): link-time version injection, --version, goreleaser config and release docs`
(= `origin/ci/add-github-actions-workflow-clean`, the newest wizard branch head on the fork)
Toolchain: go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`; `go.mod` requires 1.22)
Raw captures: `/home/c08r4d0r/wz-captures-t_10a913ec/`

## VERDICT: PASS

- `go build -buildvcs=false` → exit 0, 15,690,826-byte `net4sats-wizard`.
- Full `go test ./...` → **123 tests run, 123 PASS, 0 FAIL, 0 SKIP**, exit 0.
- Vendored portal proven embedded: all **25/25** `portal/` files (4,250,775 bytes) appear
  byte-for-byte inside the built binary.
- **No wizard source change was needed and none was made.** The hardening was not reverted;
  the vendored asset / version-bump head builds and passes as-is.
- Port `:8099` was free before the work and is free after; the only manual invocation was a
  one-shot server that was killed.
- `git status` clean, branch pushed (`a242f4b` → origin, ahead/behind 0/0).

The only commits on this branch are the task's own PROGRESS/REPORT docs. Every card step is
backed by raw terminal output below.

## Card premises that are stale on this host (disclosed, not silently worked around)

1. **`WORKTREE: /home/c03rad0r/worktrees/wizard-fix` does not exist.** This host's `$HOME` is
   `/home/c08r4d0r` (`whoami` = `c08r4d0r`), and there is no wizard clone under
   `/home/c08r4d0r/repos/`. I created the worktree at the equivalent path on this host:
   `/home/c08r4d0r/worktrees/wizard-fix`, from the wizard clone `/home/c08r4d0r/wz16-verify`:

   ```
   $ git -C /home/c08r4d0r/wz16-verify worktree add -b worker-heavy/t_10a913ec /home/c08r4d0r/worktrees/wizard-fix 0b5322a
   Preparing worktree (new branch 'worker-heavy/t_10a913ec')
   HEAD is now at 0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
   ```

   (A sibling card, `t_93d60eb5`, was blocked on this same nonexistent path — see its comment
   #29871. I did not block on it because the intent — build + full suite at the latest head —
   is unambiguously executable here.)

2. **"there is NO `--version` flag — ANY invocation starts a server" is stale at this HEAD.**
   `--version` / `-version` exists as of the version-bump commit that is the head under test
   (`main.go:637-645`, answered *before* the port bind). Proof:

   ```
   $ ./net4sats-wizard --version
   tollgate-wizard dev
   VERSION_RC=0

   $ go build -buildvcs=false -ldflags "-X main.version=v0.7.0-alpha24-verify" -o /tmp/n4s-wizard-ldflags .
   $ /tmp/n4s-wizard-ldflags --version
   tollgate-wizard v0.7.0-alpha24-verify
   VERSION_LDFLAGS_RC=0
   ```

   The port hygiene below was performed anyway, exactly as the card demands, and the server
   invocation was a one-shot that was killed.

3. **STEP 4's "pre-existing `go vet` IPv6 warning and gofmt drift" do not exist at this HEAD** —
   both are already clean (the lint commit `7ceade3` fixed them). Nothing was touched, and
   nothing was "fixed" either. See STEP 4.

4. The card calls the head "vendored/version-bumped". Both hold for this lineage: the vendored
   portal landed in `b436512 chore: vendor portal v1.0.10-alpha (pricing block) + bump
   configwizURL (#31)` and the version bump in `0b5322a` — both are ancestors of the tested HEAD
   (`git merge-base --is-ancestor` = YES).

---

## STEP 1 — Server hygiene (raw output)

### Before any local run — no wizard process, port 8099 free

```
$ pkill -x net4sats-wizard
pkill_rc=1                 # 1 = no matching process (nothing was running)
$ sleep 1
$ ss -ltnp | grep -E ':8099'
ss_8099_rc=1               # 1 = no listener = FREE
$ pgrep -a -x net4sats-wizard
pgrep_rc=1                 # 1 = no process
```

### One-shot server smoke (manual invocation, then killed)

Started the freshly built binary, verified it serves, then killed it. Capture:
`05-oneshot-server-smoke.txt`.

```
$ ./net4sats-wizard            # background, ONE-SHOT
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).

$ ss -ltnp | grep :8099        # server UP
LISTEN 0      4096            127.0.0.1:8099      0.0.0.0:*    users:(("net4sats-wizard",pid=2992050,fd=3))

$ curl -sS -o /dev/null -w '%{http_code}' http://127.0.0.1:8099/
200
$ curl -sS http://127.0.0.1:8099/api/scan
{"routers":[{"ip":"192.168.2.1","mac":"e8:8f:6f:df:9e:11", ... }, ... ]}   # live discovery JSON, 6 hosts
```

### After the work — port 8099 free again

Capture: `06-after-kill-8099-free.txt`, re-verified at the end in `10-final-8099-free.txt`.

```
$ pkill -x net4sats-wizard
pkill_rc=0                 # the one-shot server was killed (SIGTERM, exit 143)
$ sleep 1
$ pgrep -a -x net4sats-wizard
pgrep_rc=1                 # 1 = no process
$ ss -ltnp | grep :8099
ss_8099_rc=1               # 1 = no listener = FREE
```

Full `ss -ltnp` at the end (nothing on 8099; other listeners are unrelated host services):

```
State  Recv-Q Send-Q      Local Address:Port Peer Address:Port Process
LISTEN 0      5                 0.0.0.0:9100      0.0.0.0:*    users:(("python3",pid=1243,fd=3))
LISTEN 0      4096              0.0.0.0:631       0.0.0.0:*
LISTEN 0      4096              0.0.0.0:22        0.0.0.0:*
LISTEN 0      128               0.0.0.0:8443      0.0.0.0:*
LISTEN 0      4096        127.0.0.53%lo:53        0.0.0.0:*
LISTEN 0      1024            127.0.0.1:8388      0.0.0.0:*
LISTEN 0      4096           127.0.0.54:53        0.0.0.0:*
LISTEN 0      5                       *:9099            *:*    users:(("socat",pid=1244,fd=5))
LISTEN 0      4096                 [::]:631          [::]:*
LISTEN 0      1024                [::1]:8388         [::]:*
LISTEN 0      4096                 [::]:22           [::]:*
LISTEN 0      50     [::ffff:127.0.0.1]:8080            *:*    users:(("java",pid=1245,fd=66))
```

## STEP 2 — `go build -buildvcs=false` (raw output)

Capture: `02-build.txt`.

```
### host: x280 c08r4d0r 2026-10-14T14:56:56+02:00
### go version: go version go1.24.4 linux/amd64
### HEAD: 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff 0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs

$ go build -buildvcs=false
BUILD_RC=0
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 14 14:56 net4sats-wizard
```

Exit 0, no output on stdout/stderr, binary produced. `go vet ./...` also exits 0 (STEP 4).

## STEP 3 — FULL `go test ./...` (raw output, zero failures)

The module has a single package, so `./...` is the whole suite:

```
$ go list ./...
net4sats-wizard
```

### Run A — plain full suite (raw, verbatim)

Capture: `03-test-full.txt`.

```
$ go test ./...
2026-10-14T14:57:03+02:00
ok  	net4sats-wizard	67.141s
TEST_RC=0
2026-10-14T14:58:11+02:00
```

### Run B — verbose full suite, after `jq` was made available (raw tally)

Capture: `08-test-verbose-with-jq.txt` (253 lines of raw `-v` output on disk).

```
$ go test ./... -count=1 -v
ok  	net4sats-wizard	91.754s
TEST_VERBOSE_JQ_RC=0

=== tally:
RUN:  123
PASS: 123
FAIL: 0
SKIP: 0
```

Full per-test list is on disk in the capture file; `grep -E -- '--- FAIL|--- SKIP'` over it returns
**nothing**. Post-suite check: no leftover listener on `:8099` (`ss_rc=1`).

### Why an earlier run showed 2 skips (disclosed, then eliminated)

The first verbose run (`07-test-verbose.txt`, before the `jq` install) reported:

```
RUN: 120   PASS: 118   FAIL: 0   SKIP: 2
--- SKIP: TestLnIdentCommandJQRoundTrip (0.00s)
--- SKIP: TestCfgConfigCommandJQRoundTrip (0.00s)
ok  	net4sats-wizard	104.485s
```

Both skips are self-describing environment skips in `security_hardening_test.go:236` and `:291`:
`jq not in PATH — --rawfile round-trip needs a local jq`. Since those two tests are part of the
**hardening** evidence the card protects, I did not leave them skipped: I installed jq 1.7.1
user-space (`/home/c08r4d0r/.local/bin/jq`, static binary from the jqlang/jq release; no repo or
system change) and re-ran. Result: 123/123 PASS, 0 SKIP (Run B above). No test was skipped,
excluded, or filtered in the final run.

## STEP 4 — Pre-existing vet/gofmt state, hardening untouched

Capture: `04-preexisting-vet-gofmt.txt`.

```
### pre-existing go vet + gofmt state at HEAD 0b5322a (LEAVE UNTOUCHED per card STEP 4)

$ go vet ./...
VET_RC=0                   # no output — no IPv6 warning at this HEAD

$ gofmt -l .
GOFMT_RC=0                 # no output — no gofmt drift at this HEAD
```

Both symptoms the card warns about are **already resolved upstream of this HEAD** (lint commit
`7ceade3 chore(lint): pin the standard golangci-lint set and fix every finding it reports`). I
therefore had nothing to leave untouched and fixed nothing: the branch contains **zero** source
edits from this task. The only file mentioning IPv6 is an explanatory comment
(`discover.go:130`, `net.JoinHostPort`), which was neither needed nor touched.

Because no fix was required, the only commits are the task's own docs:

```
$ git log --oneline -3
a242f4b docs(t_10a913ec): wizard build + full-suite verification PROGRESS at 0b5322a
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
7ceade3 chore(lint): pin the standard golangci-lint set and fix every finding it reports
```

### Extra proof — the vendored portal is really in the built artifact

`embed.go` carries `//go:embed all:portal`. Every vendored portal file was checked to occur
byte-for-byte inside the built binary:

```
vendored portal files: 25  total bytes: 4250775
files whose exact bytes appear in built binary: 25/25
missing: []
sha256(portal/index.html) = a1fcfed53fbe57693aa22a55e53f3a78c85aa329911450f66243852f67f38819
sha256(built binary)      = 28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc
```

## Final state — `git status` clean + pushed

Capture: `09-final-git-status.txt`.

```
$ git status
On branch worker-heavy/t_10a913ec
Your branch is up to date with 'origin/worker-heavy/t_10a913ec'.

nothing to commit, working tree clean

$ git status --porcelain
porcelain_rc=0                 # no output = clean
$ git rev-list --left-right --count HEAD...origin/worker-heavy/t_10a913ec
0	0
$ git ls-remote origin refs/heads/worker-heavy/t_10a913ec
a242f4b90d3bc8a87fe73bed92ebdc65d55adfcc	refs/heads/worker-heavy/t_10a913ec
```

The block above was captured right after the `PROGRESS.md` push. This report rides the next docs
commit on the same branch, and the final state was re-captured after that push
(`09-final-git-status.txt`, overwritten by the re-run):

```
$ git status
On branch worker-heavy/t_10a913ec
Your branch is up to date with 'origin/worker-heavy/t_10a913ec'.

nothing to commit, working tree clean

$ git log --oneline -3
a645c56 docs(t_10a913ec): full verification REPORT (build + 123/123 suite green at 0b5322a)
a242f4b docs(t_10a913ec): wizard build + full-suite verification PROGRESS at 0b5322a
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs

$ git rev-list --left-right --count HEAD...origin/worker-heavy/t_10a913ec
0	0
$ git ls-remote origin refs/heads/worker-heavy/t_10a913ec
a645c56595d0ea6356f95f0928cb7c4326efebf7	refs/heads/worker-heavy/t_10a913ec
```

Both docs commits are pushed: `PROGRESS.md` at `a242f4b`, `REPORT.md` at `a645c56` (remote
`ls-remote` confirms the same SHA). This final-state paragraph itself lands in the trailing
`docs(t_10a913ec): final-state capture` commit on the same branch — the head is therefore that
docs commit, with the branch still clean and ahead/behind 0/0 against origin (checked after the
push). Build artifacts — `net4sats-wizard`, `dist/`, `coverage*.out` — are all `.gitignore`d, so
"clean" is genuine and not an artifact of ignoring new files.

## Reproduce

```
export PATH=/home/c08r4d0r/wz19-audit/go/bin:$PATH
pkill -x net4sats-wizard; sleep 1; ss -ltnp | grep :8099     # expect no listener
cd /home/c08r4d0r/worktrees/wizard-fix
go build -buildvcs=false                                     # expect exit 0
go test ./...                                                # expect: ok net4sats-wizard  ~90s, exit 0
pkill -x net4sats-wizard; sleep 1; ss -ltnp | grep :8099     # expect no listener (free again)
git status                                                   # expect: clean
```

## Residuals / notes for the fleet

- The card's `--version` premise and the `c03rad0r` path are stale; both are documented above so
  sibling cards (`t_93d60eb5`, `t_99243b9a`, `t_0509c320`) can be corrected instead of re-blocking.
- `jq` is now present at `/home/c08r4d0r/.local/bin/jq` (user-space, static binary). The wizard
  suite's two jq round-trip hardening tests therefore RUN rather than skip on this host.
- The worktree `/home/c08r4d0r/worktrees/wizard-fix` is registered against the clone
  `/home/c08r4d0r/wz16-verify`; remove with
  `git -C /home/c08r4d0r/wz16-verify worktree remove /home/c08r4d0r/worktrees/wizard-fix` when done.
- The wizard CI branch and its open PR (#35) were not modified: this task's branch was cut at
  `0b5322a` and pushed as `worker-heavy/t_10a913ec`, leaving the PR head byte-identical.
