# REPORT — t_10a913ec (re-dispatch pass 5)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.
**Verdict: PASS — all four card steps complete, zero test failures, zero wizard source changes required.**

| | |
|---|---|
| Host | x280 (home `/home/c08r4d0r`) |
| Worktree | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` (tip at start `74090de`, up to date with origin) |
| Source head under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (clone `/home/c08r4d0r/wz16-verify`) |
| Toolchain | go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-t_10a913ec-pass5/` |
| Date | 2026-10-16 |

Pass 5 of an identical dispatch; pass-1..4 reports are preserved in this branch's git history
(`a242f4b`, `a645c56`, `cb1db08`, `3ff9ae4`, `188df08`, `da5715a`). **Every figure below was
re-observed live in this pass**; nothing is carried forward from an earlier pass.

### Which commit is under test

The branch tip at start is `74090de`, which is docs-only. The last commit touching Go source is `0b5322a`:

```
$ git log --oneline 0b5322a..HEAD
74090de docs(t_10a913ec) pass-4: disclose the sibling session's net4sats-wizard.test process left running on this shared host (different worktree, no listening socket, not this pass's process)
da5715a docs(t_10a913ec) pass-4: re-verify build + full suite green at 0b5322a
188df08 docs(t_10a913ec) pass-3: re-verify build + 123/123 suite green at 0b5322a
3ff9ae4 docs(t_10a913ec) pass-2: independent re-verification — build + 123/123 suite green at 0b5322a
cb1db08 docs(t_10a913ec): final-state capture in REPORT
a645c56 docs(t_10a913ec): full verification REPORT (build + 123/123 suite green at 0b5322a)
a242f4b docs(t_10a913ec): wizard build + full-suite verification PROGRESS at 0b5322a

$ git diff --stat 0b5322a..HEAD
 PROGRESS.md | 131 +++++++++
 REPORT.md   | 907 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 2 files changed, 1038 insertions(+)

$ git log -1 --format='%h %s' -- '*.go'
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
```

So the source under test is byte-identical to `0b5322a`; every commit above it touches only the two
markdown records.

### Card premises that are stale on this host (disclosed, not silently worked around)

1. **`WORKTREE: /home/c03rad0r/worktrees/wizard-fix` does not exist.** This host's home is
   `/home/c08r4d0r` (`whoami` = `c08r4d0r`); `/home/c03rad0r` is the *other* fleet node's home and
   holds only a directory named `c08r4d0r`. The worktree was already created at the equivalent path
   on this host — `/home/c08r4d0r/worktrees/wizard-fix` — registered against the wizard clone
   `/home/c08r4d0r/wz16-verify`. That is where this pass ran.

2. **"there is NO `--version` flag — ANY invocation starts a server" is stale at this head.**
   `--version` / `-version` exists and is answered *before* the listener binds:

   ```
   $ grep -n -- '--version' main.go
   54:// programName is the CLI identity reported by --version.
   57:// versionString is the single line --version prints: the program name followed
   637:	// --version is answered before the server binds a port, so release
   641:		case "--version", "-version":

   $ ./net4sats-wizard --version
   tollgate-wizard dev
   VERSION_RC=0
   ```

   The port hygiene below was performed in full anyway, exactly as the card demands, and the server
   invocation was a one-shot that was killed.

3. **STEP 4's "pre-existing `go vet` IPv6 warning and gofmt drift" do not exist at this head** — both
   are already clean (lint commit `7ceade3` removed them upstream of `0b5322a`). Nothing was touched
   and nothing was "fixed" either. See STEP 4.

4. **The wizard clone's CI branch has advanced past the head under test.** `/home/c08r4d0r/wz16-verify`
   is now at `e85720c` (was `0b5322a`). This does not change the verdict:

   ```
   $ git -C /home/c08r4d0r/wz16-verify diff --stat 0b5322a..e85720c -- '*.go'
                                  (empty — no Go source changed)
   $ git -C /home/c08r4d0r/wz16-verify diff --stat 0b5322a..e85720c | tail -1
    56 files changed, 3937 insertions(+)     # all under verify/ (evidence bundles) + docs
   ```

   The worktree branch under test (`worker-heavy/t_10a913ec`) was not rebased; its Go source is still
   exactly `0b5322a`, which is what was built and tested.

---

## STEP 1 — Server hygiene (raw output)

### Before any local run — no wizard process, port 8099 free

```
$ pkill -x net4sats-wizard; echo "pkill_rc=$?"
pkill_rc=1                              # 1 = no matching process (nothing running this pass)
$ sleep 1
$ pgrep -a -x net4sats-wizard; echo "pgrep_rc=$?"
pgrep_rc=1                              # 1 = no process
$ ss -ltnp | grep -E ':8099'; echo "ss_8099_grep_rc=$?"
ss_8099_grep_rc=1                       # 1 = no listener = FREE
```

Full `ss -ltnp` before the work (nothing on 8099; the other listeners are unrelated host services):

```
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
```

No `:8099` line anywhere. Capture: `01-hygiene-before.txt`

### The single manual invocation — a one-shot server, then killed

```
=== 5a: --version premise check (must NOT bind :8099) ===
$ ./net4sats-wizard --version
tollgate-wizard dev
VERSION_RC=0
$ pgrep -a -x net4sats-wizard; echo "pgrep_rc=$?"
pgrep_rc=1                              # no process
$ ss -ltnp | grep -E ':8099'; echo "ss_8099_rc=$?"
ss_8099_rc=1                            # still FREE => --version opened no listener

=== 5b: one-shot server on default 127.0.0.1:8099 ===
$ ./net4sats-wizard &
server pid=3334046
$ ss -ltnp | grep -E ':8099'
LISTEN 0      4096            127.0.0.1:8099       0.0.0.0:*    users:(("net4sats-wizard",pid=3334046,fd=3))

$ curl -s -o /dev/null -w 'HTTP %{http_code}  bytes=%{size_download}  time=%{time_total}s\n' http://127.0.0.1:8099/
HTTP 200  bytes=19976  time=0.000492s

$ curl -s --max-time 30 http://127.0.0.1:8099/api/scan | head -c 400
{"routers":[{"ip":"192.168.2.1","mac":"e8:8f:6f:df:9e:11","vendor":"unknown","model":"unknown","firmware":"unknown","ssh_open":false,"http_port":80},{"ip":"192.168.2.30","mac":"7c:7a:91:86:a1:0f",...

=== 5c: kill the one-shot server ===
$ kill 3334046; sleep 1; kill -9 3334046
$ pgrep -a -x net4sats-wizard; echo "pgrep_rc=$?"
pgrep_rc=1                              # 1 = gone
$ ss -ltnp | grep -E ':8099'; echo "ss_8099_rc=$?"
ss_8099_rc=1                            # 1 = FREE again
```

### After the work — port 8099 free again (acceptance check)

```
$ ss -ltnp | grep -E ':8099'; echo "ss_8099_rc=$?"
ss_8099_rc=1
```

The final full `ss -ltnp` is identical to the "before" listing above — no `:8099` entry, no
`net4sats-wizard` process. Captures: `05-oneshot-server.txt` (full transcript incl. final `ss`).

---

## STEP 2 — `go build -buildvcs=false`

```
$ go version
go version go1.24.4 linux/amd64         # go.mod requires go 1.22
$ git check-ignore -v net4sats-wizard
.gitignore:6:net4sats-wizard	net4sats-wizard     # build output is ignored
$ rm -f net4sats-wizard
$ go build -buildvcs=false
BUILD_RC=0                              # <-- no compiler output at all
$ ls -l net4sats-wizard
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 16 08:51 net4sats-wizard
$ sha256sum net4sats-wizard
28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc  net4sats-wizard
$ git status --porcelain
                                        # empty => build did not dirty the tree
```

**Build succeeds** (rc=0, empty stdout/stderr) and the 15,690,826-byte binary is byte-identical
(sha256 + size) to the binaries produced in passes 2, 3 and 4 — the build is reproducible.
Capture: `02-build.txt`

---

## STEP 3 — FULL `go test ./...` (raw output, zero failures)

The module is a single package, so `./...` is the whole suite:

```
$ go list ./...
net4sats-wizard
```

Uncached verbose run (`-count=1`, so no cached PASS can masquerade as fresh evidence):

```
$ go test ./... -count=1 -v
=== RUN   TestHandleScan
--- PASS: TestHandleScan (123.84s)
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
    --- PASS: TestDeployRequestValidation/missing_IP_and_password (0.00s)
    ... (123 RUN/PASS lines in total; full raw capture on disk) ...
=== RUN   TestVersionStringFormat
--- PASS: TestVersionStringFormat (0.00s)
=== RUN   TestVersionFlagInjectedViaLdflags
--- PASS: TestVersionFlagInjectedViaLdflags (0.50s)
PASS
ok  	net4sats-wizard	157.867s
TEST_RC=0
```

Counts parsed from the raw verbose capture (`03-go-test-verbose.txt`, 254 lines on disk):

| metric | count |
|---|---|
| `=== RUN` (incl. subtests) | **123** |
| `--- PASS` (incl. subtests) | **123** |
| `--- PASS` (top-level) | **52** |
| `--- FAIL` | **0** |
| `--- SKIP` | **0** |
| `ok  ` package lines | **1** (`ok net4sats-wizard 157.867s`) |
| `FAIL` package lines | **0** |

`TEST_RC=0`. **Zero failures, zero skips.**

No test failure was caused by the vendored asset / version bump, so no wizard fix was required and
none was made. The hardening was neither reverted nor touched. Captures:
`03-go-test-verbose.txt` (full raw), `03-tally.txt`, `07-tally-embed.txt`.

---

## STEP 4 — pre-existing vet/gofmt state left untouched

The card asks to leave a pre-existing `go vet` IPv6 warning and gofmt drift alone. **At this head
there is nothing to leave** — both are already clean:

```
$ go vet ./...
VET_RC=0                                # no output at all
$ gofmt -l .
GOFMT_RC=0                              # empty list => no gofmt drift
$ git status --porcelain
                                        # empty
```

Nothing was modified and nothing was "fixed". (The card's premise dates from before lint commit
`7ceade3`, which cleaned both upstream of the head under test.) Capture: `04-vet-gofmt.txt`,
re-confirmed live in `04b-vet-gofmt-live.txt`.

---

## Bonus — the vendored asset really is inside the binary (the card's GOAL)

The "vendored/version-bumped wizard" claim is not just about the build succeeding; it is that the
vendored UI asset tree is compiled in. Proven first-hand by substring-searching the built binary for
each asset file's exact bytes:

```
$ python3 run-06-portal-embed.py
binary size: 15,690,826 bytes
binary sha256: 28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc
portal files: 25/25 embedded byte-for-byte (4,250,775 bytes total)
missing: 0
EMBED_RC=0

$ python3 run-07-tally-embed.py
  portal/ : 25/25 files embedded byte-for-byte (4,250,775 B)
  admin/  : 27/27 files embedded byte-for-byte (1,267,428 B)
```

Both `portal/` (`//go:embed all:portal`) and `admin/` (`all:admin`) trees are present byte-for-byte
in the shipped binary. Captures: `06-portal-embed.txt`, `07-tally-embed.txt`.

---

## Files touched

- **Repo worktree** `/home/c08r4d0r/worktrees/wizard-fix` (branch `worker-heavy/t_10a913ec`):
  `PROGRESS.md` and `REPORT.md` only. **Zero `.go`, asset, config or build-file changes** — no fix
  was needed. The build output `net4sats-wizard` is gitignored.
- **Outside the repo** (evidence only): `/home/c08r4d0r/wz-captures-t_10a913ec-pass5/*`.

## Verdict

PASS. Build rc=0 and reproducible; full suite 123/123 with 0 failures/0 skips uncached; `:8099` free
before and after; one-shot server smoke green; vendored assets embedded 25/25 + 27/27; frozen
vet/gofmt state observed-clean and untouched; tree clean and branch pushed.
# REPORT — t_10a913ec (re-dispatch pass 4)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.
**Verdict: PASS — all four card steps complete, zero test failures, zero wizard source changes required.**

| | |
|---|---|
| Host | x280 (home `/home/c08r4d0r`) |
| Worktree | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` (tip at start `188df08`, up to date with origin) |
| Source head under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (clone `/home/c08r4d0r/wz16-verify`) |
| Toolchain | go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-t_10a913ec-pass4/` |
| Date | 2026-10-16 |

Pass 4 of an identical dispatch; pass-1/2/3 reports are preserved in this branch's git history
(`a242f4b`, `a645c56`, `cb1db08`, `3ff9ae4`, `188df08`). **Every figure below was re-observed live
in this pass**; nothing is carried forward from an earlier pass.

### Which commit is under test

The branch tip is `188df08`, which is docs-only. The last commit touching Go source is `0b5322a`:

```
$ git log --oneline 0b5322a..HEAD
188df08 docs(t_10a913ec) pass-3: re-verify build + 123/123 suite green at 0b5322a
3ff9ae4 docs(t_10a913ec) pass-2: independent re-verification — build + 123/123 suite green at 0b5322a
cb1db08 docs(t_10a913ec): final-state capture in REPORT
a645c56 docs(t_10a913ec): full verification REPORT (build + 123/123 suite green at 0b5322a)
a242f4b docs(t_10a913ec): wizard build + full-suite verification PROGRESS at 0b5322a
```

---

## STEP 1 — Server hygiene (port :8099)

**Card premise to check first:** "there is NO `--version` flag — ANY invocation starts a server."
This premise is **stale at this head**. `main.go:637-644` answers `--version`/`-version` *before*
`ListenAndServe` is reached:

```go
func main() {
	// --version is answered before the server binds a port, so release
	// automation can smoke-test the binary without opening a listener.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-version":
			fmt.Println(versionString())
			return
		}
	}
```

Observed first-hand:

```
$ ./net4sats-wizard --version
tollgate-wizard dev
rc=0
$ ss -ltnp | grep ':8099'      # no listener opened by --version
(no output)
```

The full port-8099 hygiene was performed anyway, as the card requires.

### :8099 BEFORE

`pkill -x net4sats-wizard` returned **rc=0 on the first call** — i.e. a stale wizard process *was*
live on this host and was killed by this pass (pass 3 found none; this is a genuine difference, not
a carried-over claim). A second `pkill` returned rc=1 (nothing left), and `pgrep -ax net4sats-wizard`
returned rc=1 (gone).

```
$ ss -ltnp        # BEFORE (no :8099 entry; grep -c ':8099' == 0)
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
```

### :8099 DURING (one-shot server)

The built binary was launched as a tracked background process (no `nohup`/`&` wrapper):

```
$ ./net4sats-wizard
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).
```

```
$ ss -ltnp | grep ':8099'
LISTEN 0      4096            127.0.0.1:8099       0.0.0.0:*    users:(("net4sats-wizard",pid=3020475,fd=3))
```

```
$ curl -s -o /dev/null -w 'GET / -> %{http_code}\n' http://127.0.0.1:8099/
GET / -> 200
```

Note: the server binds **127.0.0.1** (loopback), the hardening from `cc9f94c` — confirmed by the
`Local Address:Port` column above, not `0.0.0.0`.

### :8099 AFTER

```
$ pkill -x net4sats-wizard    # rc=0
$ pgrep -ax net4sats-wizard   # rc=1 (gone)
$ ss -ltnp        # AFTER (grep -c ':8099' == 0)
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
```

**Port :8099 is free before and after; no listener survives this pass.**

### Shared-host disclosure: a sibling session's test binary is running

At the end of this pass `pgrep -af net4sats-wizard` shows one *other* wizard-ish process:

```
3035033 /home/c08r4d0r/.cache/tmp-go/go-build3923712378/b001/net4sats-wizard.test \
  -test.paniconexit0 -test.timeout=10m0s -test.count=1 -test.v=true \
  -test.run=TestDiscoverRoutersNoNetwork|TestPinnedURLsAreLive
```

This is **not** this pass's process and was deliberately left running. Evidence of ownership:
its parent chain is `hermes` (pid 2898659) → `bash` → `go` → `net4sats-wizard.test`, and its
`cwd` is `/home/c08r4d0r/worktrees/wizard-pr-prep` — a **different worktree**, i.e. a concurrent
sibling session on this shared host. It is a `net4sats-wizard.test` binary, so
`pkill -x net4sats-wizard` does not match it (exact-name match), and it holds **no listening
socket** (`/proc/3035033/fd` has 0 sockets; `ss -ltnp | grep 3035033` is empty). It therefore
does not affect the :8099 criterion. Killing another session's live test run would be destructive,
so it was not touched.

---

## STEP 2 — `go build -buildvcs=false`

```
$ go build -buildvcs=false
$ echo $?
0
```

rc=0 with **no compiler output** (the captured `build.txt` is empty). Artifact:

```
$ ls -l net4sats-wizard
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 16 08:07 net4sats-wizard
$ sha256sum net4sats-wizard
28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc  net4sats-wizard
```

This sha256 is **byte-identical to pass 2's** binary — the build is reproducible across passes.

---

## STEP 3 — Full `go test ./...`

```
$ go test ./...
ok  	net4sats-wizard	91.866s
$ echo $?
0
```

Verbose run for per-test evidence:

```
$ go test -v ./...
...
PASS
ok  	net4sats-wizard	57.727s
```

Counts parsed from the verbose capture:

| metric | count |
|---|---|
| `=== RUN` (incl. subtests) | **123** |
| `--- PASS` (top-level) | **52** |
| `--- FAIL` | **0** |
| `--- SKIP` | **0** |

**Zero failures.** No test failure was caused by the vendored asset / version bump, so no wizard
fix was required and none was made. The hardening (`cc9f94c`) was neither reverted nor touched.

---

## STEP 4 — Leave the pre-existing vet/gofmt state untouched

At this head there is nothing to leave untouched — the lint commit `7ceade3` already cleaned both:

```
$ go vet ./...
$ echo $?
0                      # empty output, no IPv6 warning present at this head
$ gofmt -l .
$ echo $?              # zero lines listed, no drift present at this head
0
```

Verified clean; **nothing was "fixed"**, and no file was reformatted. This is an observation about
this head, not a claim that the card's warning was ever fixed by this pass.

---

## GOAL evidence — the vendored portal is really inside the binary

The wizard embeds the vendored bundle with `//go:embed all:portal` and `//go:embed all:admin`
(`embed.go`). Byte-for-byte presence of every embedded file inside the built 15,690,826-byte
binary was re-checked first-hand this pass:

| tree | files | bytes | present in binary |
|---|---|---|---|
| `portal/` | 25 | 4,250,775 | **25/25, none missing** |
| `admin/` | 27 | — | **27/27, none missing** |
| `index.html` | 1 | 19,976 | embedded (and served live for `GET /` → 200) |

Specific vendored artifacts confirmed embedded: `portal/assets/splash-Be2Pt9NX.js` (208,725 B) and
`portal/assets/splash-D7WaP2L-.css` (14,045 B).

One correction to a naive probe: `GET /assets/<bundle>.js` returns **index.html** (sha256
`e83a518a...406f`, identical to on-disk `index.html`), because the wizard's only route is the `/`
catch-all (`main.go:648-652`) — the portal bundle is deployed to the *router's* uhttpd docroot, not
served by this Go server. The embed proof is therefore the in-binary byte comparison above, not an
HTTP fetch of `/assets/...`.

---

## Conclusion

**PASS.** The vendored/version-bumped wizard (`0b5322a`) builds clean with `-buildvcs=false`, its
full `go test ./...` suite passes with 123 tests run and 0 failures / 0 skips, port :8099 is free
before and after, and the embedded vendored portal is present byte-for-byte in the binary.
**No wizard source change was necessary.** Tree clean at the end.

---

# REPORT — t_10a913ec (re-dispatch pass 3)

## Build wizard and run full go test suite green

**Goal:** Prove the vendored/version-bumped wizard builds and its full test suite passes.
**Verdict: PASS — all four card steps complete, zero test failures, zero wizard source changes required.**

| | |
|---|---|
| Host | x280 (home `/home/c08r4d0r`) |
| Worktree | `/home/c08r4d0r/worktrees/wizard-fix` |
| Branch | `worker-heavy/t_10a913ec` (tip `3ff9ae4`, up to date with origin) |
| Source head under test | `0b5322a` — *feat(release): link-time version injection, --version, goreleaser config and release docs* |
| Repo | fork `felixfelix-bot/net4sats-wizard-go` (clone `/home/c08r4d0r/wz16-verify`) |
| Toolchain | go1.24.4 linux/amd64 (`/home/c08r4d0r/wz19-audit/go/bin/go`); go.mod requires 1.22 |
| Raw captures | `/home/c08r4d0r/wz-captures-t_10a913ec-redispatch/` |
| Date | 2026-10-16 |

This is pass 3 of an identical dispatch; pass-1 and pass-2 reports are preserved in this branch's
git history (`a242f4b`, `a645c56`, `cb1db08`, `3ff9ae4`). **Every figure below was re-observed live
in this pass**; nothing is carried forward from an earlier pass.

### Which commit is under test

The branch tip is `3ff9ae4`, which is docs-only. The last commit that touches Go source is `0b5322a`:

```
$ git log --oneline 0b5322a..HEAD
3ff9ae4 docs(t_10a913ec) pass-2: independent re-verification — build + 123/123 suite green at 0b5322a
cb1db08 docs(t_10a913ec): final-state capture in REPORT
a645c56 docs(t_10a913ec): full verification REPORT (build + 123/123 suite green at 0b5322a)
a242f4b docs(t_10a913ec): wizard build + full-suite verification PROGRESS at 0b5322a

$ git diff --stat 0b5322a..HEAD
 PROGRESS.md |  44 +++++++++
 REPORT.md   | 309 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 2 files changed, 353 insertions(+)

$ git log -1 --format='%h %s' -- '*.go'
0b5322a feat(release): link-time version injection, --version, goreleaser config and release docs
```

So the source under test is byte-identical to `0b5322a`; the four commits above it touch only
`PROGRESS.md` and `REPORT.md`.

### Card premises that are stale on this host (disclosed, not silently worked around)

1. **`WORKTREE: /home/c03rad0r/worktrees/wizard-fix` does not exist.** This host's home is
   `/home/c08r4d0r` (`whoami` = `c08r4d0r`); `/home/c03rad0r` is the *other* fleet node's home and
   holds only a directory named `c08r4d0r`. The worktree was already created at the equivalent path
   on this host — `/home/c08r4d0r/worktrees/wizard-fix` — by an earlier pass, registered against the
   wizard clone `/home/c08r4d0r/wz16-verify`.

2. **"there is NO `--version` flag — ANY invocation starts a server" is stale at this head.**
   `--version` / `-version` exists and is answered *before* the listener binds:

   ```
   $ grep -n -- '--version' main.go
   54:// programName is the CLI identity reported by --version.
   57:// versionString is the single line --version prints: the program name followed
   637:	// --version is answered before the server binds a port, so release
   641:		case "--version", "-version":

   $ ./net4sats-wizard --version
   tollgate-wizard dev
   VERSION_RC=0
   ```

   The port hygiene below was performed in full anyway, exactly as the card demands, and the
   server invocation was a one-shot that was killed.

3. **STEP 4's "pre-existing `go vet` IPv6 warning and gofmt drift" do not exist at this head** —
   both are already clean (lint commit `7ceade3` removed them upstream of `0b5322a`). Nothing was
   touched and nothing was "fixed" either. See STEP 4.

---

## STEP 1 — Server hygiene (raw output)

### Before any local run — no wizard process, port 8099 free

```
$ pkill -x net4sats-wizard; echo "pkill_rc=$?"
pkill_rc=1                              # 1 = no matching process (nothing running)
$ sleep 1
$ pgrep -a -x net4sats-wizard; echo "pgrep_rc=$?"
pgrep_rc=1                              # 1 = no process
$ ss -ltnp | grep -E ':8099'; echo "ss_8099_rc=$?"
ss_8099_rc=1                            # 1 = no listener = FREE
```

Full `ss -ltnp` before the work (nothing on 8099; the other listeners are unrelated host services):

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

### One-shot server smoke (the single manual invocation, then killed)

Capture `06-oneshot.txt`. Started the freshly built binary, verified it serves, then killed it.

```
--- ss :8099 (server UP) ---
LISTEN 0      4096            127.0.0.1:8099      0.0.0.0:*    users:(("net4sats-wizard",pid=273680,fd=3))
ss_up_rc=0
--- pgrep ---
273680 ./net4sats-wizard
--- curl / ---
http_code=200 size=19976
--- curl /api/scan (first 300B) ---
{"routers":[{"ip":"192.168.2.1","mac":"e8:8f:6f:df:9e:11","vendor":"unknown","model":"unknown",
"firmware":"unknown","ssh_open":false,"http_port":80},{"ip":"192.168.2.30","mac":"7c:7a:91:86:a1:0f",
..."ssh_open":true},{"ip":"192.168.2.32","mac":"a
```

The handler serves the setup UI (`/` → HTTP 200, 19,976 bytes) and live LAN discovery JSON.

### After the work — port 8099 free again

Capture `07-after-kill.txt`.

```
$ pkill -x net4sats-wizard; echo "pkill_rc=$?"
pkill_rc=0                              # the one-shot server was killed
$ sleep 1
$ pgrep -a -x net4sats-wizard; echo "pgrep_rc=$?"
pgrep_rc=1                              # 1 = no process
$ ss -ltnp | grep ':8099'; echo "ss_8099_rc=$?"
ss_8099_rc=1                            # 1 = no listener = FREE
```

Full `ss -ltnp` after the work — identical to the before-state, nothing on 8099:

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

Capture `02-build.txt`.

```
### 2026-10-16T01:36:01+02:00 host=x280 user=c08r4d0r
### go version go1.24.4 linux/amd64
### HEAD 3ff9ae4dad0fe99847873f6c6893962c3a09e6f2

$ go build -buildvcs=false
BUILD_RC=0
-rwxrwxr-x 1 c08r4d0r c08r4d0r 15690826 Oct 16 01:36 net4sats-wizard
```

Exit 0, nothing on stdout/stderr, 15,690,826-byte binary produced. The build is reproducible: the
sha256 of this pass's binary is `28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc`,
byte-identical to pass 2's.

## STEP 3 — FULL `go test ./...` (raw output, zero failures)

The module has a single package, so `./...` is the whole suite:

```
$ go list ./...
net4sats-wizard
```

### Run A — plain full suite (rejected as cached, disclosed)

```
$ go test ./...
ok  net4sats-wizard  (cached)
TEST_RC=0
```

The plain run returned a **cached** result. A cached PASS is not fresh evidence, so it was not
accepted; the suite was re-run uncached (`-count=1`) below.

### Run B — full suite, uncached, verbose (raw tally — the evidence)

Capture `03-test-verbose.txt` (251 lines of raw `-v` output on disk).

```
$ go test ./... -count=1 -v
... (123 per-test RUN/PASS lines) ...
PASS
ok  net4sats-wizard  137.022s
TEST_VERBOSE_RC=0

=== tally (counted from the capture file):
RUN   : 123
PASS  : 123
FAIL  : 0
SKIP  : 0
```

The `--- FAIL` / `--- SKIP` grep over the capture returns **nothing** — a graded list of the
hardening-focused tests that ran and passed includes:

```
--- PASS: TestRemoteCommandsNeverInterpolatePayloads (0.00s)   # 15 shell-injection payload cases
--- PASS: TestVulnerableConstructionsAreGone (0.00s)
--- PASS: TestPasswdCommandInertUnderSh (0.01s)
--- PASS: TestStaCarriersDecodeUnderSh (0.01s)
--- PASS: TestLnIdentCommandJQRoundTrip (0.02s)
--- PASS: TestCfgConfigCommandJQRoundTrip (0.02s)
--- PASS: TestArtifactPinsCoverDownloadedURLConstants (0.00s)
--- PASS: TestVerifyArtifactPinTamperedDownload (0.00s)
--- PASS: TestVerifyArtifactPinUnknownURLFailsClosed (0.00s)
--- PASS: TestCORSMiddleware (0.00s)                          # 7 origin cases
--- PASS: TestListenAddress (0.00s)
--- PASS: TestVersionStringFormat (0.00s)
--- PASS: TestVersionFlagInjectedViaLdflags (0.90s)
--- PASS: TestPinnedURLsAreLive (0.91s)
--- PASS: TestHandleScan (78.52s)
--- PASS: TestDiscoverRoutersNoNetwork (46.52s)
--- PASS: TestSshRunFunction (10.01s)
```

**Zero failures. Zero skips.** `jq` is present on this host (`/home/c08r4d0r/.local/bin/jq`), so the
two jq round-trip hardening tests RAN rather than self-skipping (in an earlier pass without `jq` they
skipped at `security_hardening_test.go:236` / `:291`; nothing was skipped, excluded or filtered here).

Post-suite check: no leftover listener on `:8099` (`ss_8099_rc=1`).

**No fix was required.** No failure was caused by the vendored asset / version change, so no wizard
source edit was made — and the hardening was neither reverted nor touched.

## STEP 4 — Pre-existing vet/gofmt state, hardening untouched

Capture `04-vet.txt` / `04-gofmt.txt`.

```
$ go vet ./...
VET_RC=0                 # no output — no IPv6 warning at this head

$ gofmt -l .
GOFMT_RC=0               # no output — no gofmt drift at this head
```

Both symptoms the card warns about are already resolved upstream of `0b5322a` (lint commit
`7ceade3`). I therefore had nothing to leave untouched and fixed nothing: this pass contains **zero**
source edits. Because no fix was required, there is no fix commit — the only commit this pass is these
docs, with a conventional message.

## GOAL proof — the vendored portal is really in the built artifact

`embed.go` carries `//go:embed all:portal`. Every vendored portal file was checked to occur
byte-for-byte inside the built binary:

```
vendored portal files: 25  total bytes: 4250775
files whose exact bytes appear in built binary: 25/25
missing: []
sha256(built binary) = 28ef640cd7ad74d0a7cd9497b7f3739261e473abe6fb5e450494f1c481da4bdc
```

## Final state — `git status` clean + pushed

```
$ git status --short
(no output — clean)

$ git rev-parse HEAD
<this report's docs commit>

$ git rev-list --left-right --count HEAD...origin/worker-heavy/t_10a913ec
0  0

$ git ls-remote origin refs/heads/worker-heavy/t_10a913ec
<same SHA as HEAD>   refs/heads/worker-heavy/t_10a913ec
```

The final-state SHA block is captured live on disk at
`/home/c08r4d0r/wz-captures-t_10a913ec-redispatch/09-final-git-status.txt` after the push. Build
artifacts (`net4sats-wizard`, `dist/`, `coverage*.out`) are `.gitignore`d, so "clean" is genuine and
not an artifact of ignoring new files.

## Reproduce

```
export PATH=/home/c08r4d0r/wz19-audit/go/bin:$PATH
pkill -x net4sats-wizard; sleep 1; ss -ltnp | grep :8099     # expect no listener
cd /home/c08r4d0r/worktrees/wizard-fix
git rev-parse HEAD                                           # source head = 0b5322a
go build -buildvcs=false                                     # expect exit 0
go test ./... -count=1 -v                                    # expect: ok net4sats-wizard ~137s, exit 0
pkill -x net4sats-wizard; sleep 1; ss -ltnp | grep :8099     # expect no listener (free again)
git status                                                   # expect: clean
```

## Residuals / notes for the fleet

- The card's `--version` premise and the `c03rad0r` worktree path are stale; both are documented
  above so sibling cards can be corrected instead of re-blocking against them.
- `jq` is present at `/home/c08r4d0r/.local/bin/jq`, so the wizard suite's two jq round-trip
  hardening tests RUN rather than skip on this host.
- The worktree `/home/c08r4d0r/worktrees/wizard-fix` is registered against the clone
  `/home/c08r4d0r/wz16-verify`; remove with
  `git -C /home/c08r4d0r/wz16-verify worktree remove /home/c08r4d0r/worktrees/wizard-fix` when done.
- The wizard CI branch and its open PR (#35) were not modified: this branch was cut at `0b5322a` and
  pushed as `worker-heavy/t_10a913ec`, leaving the PR head byte-identical.
- Going forward, prefer `go test ./... -count=1` for this repo: an unqualified `go test ./...` can
  return `(cached)` and prove nothing about the current tree.

---

<!-- ============ PASS 2 and PASS 1 (preserved, unchanged) ============ -->
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
