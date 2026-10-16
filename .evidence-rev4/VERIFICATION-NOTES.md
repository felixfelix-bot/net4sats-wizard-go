# VERIFICATION NOTES — hermes verify on branch worker-heavy/t_0509c320-r3

Added after the turn that wrote `.evidence-rev4/checksums.sh` and `.evidence-rev4/version-smoke.sh`.
Those two files are evidence-capture shell scripts, not product code; `bash -n` is clean on both.
This note records what the workspace verification recipe actually reports and why.

## Result

`hermes verify --json` (recipe: "Go project", detected):

| phase | command | exitCode | ok |
|-------|---------|----------|-----|
| build | `go build ./...` | 0 | true |
| test  | `go test ./...` | 0 | true |
| readiness | GET `http://127.0.0.1:8000/` | — | **false** |

Three consecutive runs, identical shape:
```
ok= False | build rc= 0 | test rc= 0 'ok  \tnet4sats-wizard\t90.917s' | readiness ready= False url= http://127.0.0.1:8000/ err= <urlopen error [Errno 111] Connection refused>
ok= False | build rc= 0 | test rc= 0 'ok  \tnet4sats-wizard\t(cached)'   | readiness ready= False url= http://127.0.0.1:8000/ err= <urlopen error [Errno 111] Connection refused>
ok= False | build rc= 0 | test rc= 0 'ok  \tnet4sats-wizard\t(cached)'   | readiness ready= False url= http://127.0.0.1:8000/ err= <urlopen error [Errno 111] Connection refused>
```
Both real phases (build, test) pass on every run. `ok:false` is entirely attributable to `readiness`.

## Why the readiness phase is red — misdetected port, not a defect

The harness probes **`http://127.0.0.1:8000/`**. The wizard binds **`127.0.0.1:8099`**:

```
$ grep -n 'defaultListenAddr\|8099' main.go
67:const defaultListenAddr = "127.0.0.1:8099"
83:	"http://127.0.0.1:8099": true,
84:	"http://localhost:8099": true,
658:	fmt.Printf("net4sats wizard running on http://%s\n", addr)
660:	fmt.Println("To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).")

$ grep -rn '\b8000\b' --include='*.go' --include='*.yaml' --include='*.yml' . | grep -v .evidence-rev4
(empty)                 # port 8000 appears nowhere in the project
```

Direct proof against the freshly built binary:
```
$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
net4sats wizard running on http://127.0.0.1:8099
Open this URL in your browser to set up a router.
To serve on other interfaces, set WIZARD_BIND (e.g. WIZARD_BIND=0.0.0.0:8099).

$ curl -s -o /dev/null -w 'http_code=%{http_code}\n' http://127.0.0.1:8099/
http_code=200                    # the app really is up and serving
$ curl -s -o /dev/null -w 'http_code=%{http_code}\n' --max-time 3 http://127.0.0.1:8000/
http_code=000    curl_rc=7       # nothing on 8000, by design
```
No foreground process was left behind after the probes (`ss -ltnp` shows nothing on 8000/8099).

## The one failing run that was NOT the port mismatch — cause identified

The first clean-env re-run produced a different failure: `exitCode 1`, `"signal: terminated\nFAIL\n"`, and
`readiness: null`. I initially attributed this to my own `WIZARD_BIND=127.0.0.1:8000` leak; that was wrong — it
reproduced without the variable set. Investigation:

- `version_test.go:38` runs a nested `go build` inside `go test`, so the suite compiles under itself.
- The failure is nondeterministic and external: a sibling agent session is concurrently active on this host
  running an unrelated `wizard-fix` task whose brief instructs it to `pkill -x net4sats-wizard` and to start/kill
  the same binary. `signal: terminated` is that interference (SIGTERM delivered to a contended process).
- Controlled tests, all green:
  - isolated `go test -run TestVersionFlagInjectedViaLdflags -count=1` → `ok net4sats-wizard 0.464s` three times;
  - four sequential full `go test ./... -count=1 -v` runs → PASS, PASS, PASS, PASS (`ok` 91.6s / 57.8s / 69.3s / 74.9s);
  - three *parallel* full suite runs (`.evidence-rev4/parallel-suite.sh`, nproc=8) → all three `rc=0`, all `ok`;
  - `hermes verify` three times → build rc 0 and test rc 0 every time.
- 0 `DATA RACE`, 0 OOM kills in reachable logs.

**Conclusion:** the transient `signal: terminated` runs are cross-task interference on a shared host, not a repo
defect. The claimed blocker (env leak) was retracted; the port-mismatch explanation covers the reproducible red.

## Reproduce
```
export PATH="/usr/local/go/bin:$HOME/wz16-tools:$PATH"
cd /home/c08r4d0r/worktrees/wz-e2e-gate-r3
go build ./... && go test ./... -count=1        # both green
hermes verify --json                            # build+test green; readiness red on port 8000 only
```

## Orientation note
`main.go:637-642` answers `--version` **before** binding a port. The sibling task's brief asserting "there is NO
`--version` flag — ANY invocation starts a server" describes an older revision; at `a412b61`/`18a1764` the flag
exists and exits 0 with `tollgate-wizard <version>`.
