# Verification evidence — WZ1.9 capture lane (t_8b5cf200)

Worktree: `/home/c08r4d0r/worktrees/wz19-video`, branch `wz19-video-capture`.
Run 2026-10-23. Go toolchain: `/usr/local/go/bin` (go1.22.12 linux/amd64).

## 1. Repo verification recipe — `hermes verify` PASS

The project has no canonical `make test`, so the repo's runnable recipe was used:
`hermes verify` (detect -> build -> test -> start -> readiness poll).

### First run — `ok: false`, and the cause was the RECIPE, not the code

```json
{"recipe":"Go project","ok":false,
 "phases":[{"phase":"build","command":"go build ./...","exitCode":0,"ok":true,"duration":11.623},
           {"phase":"test","command":"go test ./...","exitCode":0,"ok":true,"duration":62.517,
            "outputTail":"ok  \tnet4sats-wizard\t59.592s\n"}],
 "readiness":{"url":"http://127.0.0.1:8000/","ready":false,
              "error":"<urlopen error [Errno 111] Connection refused>",
              "outputTail":"net4sats wizard running on http://127.0.0.1:8099\n..."},
 "source":"detected"}
```

Build exit 0, tests exit 0 (`ok net4sats-wizard 59.592s`) — but readiness failed
because the detector guessed port **8000** while the wizard binds **8099**
(`main.go: defaultListenAddr = "127.0.0.1:8099"`). The app's own stdout in the same
record says `running on http://127.0.0.1:8099`. The app was up; the poll knocked on
the wrong door.

### Proof the app was fine — `--phase start --port 8099`

```json
{"recipe":"Go project","ok":true,"phases":[],
 "readiness":{"url":"http://127.0.0.1:8099/","ready":true,"statusCode":200,
              "duration":1.006,"error":null}}
```

### Recipe corrected, then FULL verify — `ok: true`

Added `.hermes/environment.json` pinning `port: 8099` and the reason, so this does
not recur for the next lane.

```json
{"recipe":"Go project (net4sats wizard)","ok":true,
 "phases":[{"phase":"build","command":"go build ./...","exitCode":0,"ok":true},
           {"phase":"test","command":"go test ./...","exitCode":0,"ok":true,
            "outputTail":"ok  \tnet4sats-wizard\t(cached)\n"}],
 "readiness":{"url":"http://127.0.0.1:8099/","ready":true,"statusCode":200},
 "source":"manifest"}
```

### Cache defeated — fresh full suite (the regression-capable proof)

The recipe's `go test ./...` can return `(cached)`, and a cached result cannot detect
a regression — so a separate uncached run is the evidence that actually counts:

```
$ go test -count=1 ./...
ok  	net4sats-wizard	60.934s
UNCACHED_EXIT=0
```

I attempted to pin `-count=1` into `.hermes/environment.json` so the recorded run
would always be uncached. **That write was refused by the protected-file guard**
("write to protected agent-instruction file(s) … approval prompt timed out"), and I
did not retry it via another path. The manifest therefore still reads
`go test ./...`; the uncached run above stands as the real evidence, and pinning the
flag remains a follow-up (see REMAINING in REPORT.md).

## 2. `capture.cjs` syntax gate

```
$ node --check verify/wz19-video-capture/capture.cjs
capture.cjs: SYNTAX OK
```

## 3. Reproducibility — the take was re-run and matched

The same driver was re-run against the same binary on a **second, clean Xvfb display**
(`:98`, 1920x1080x24) to prove the result is not a one-off.

| | Run 1 (committed) | Run 2 (replay) |
|---|---|---|
| Routers discovered | 5 (`.1`, `.2`, `.30`, `.121`, `.130`) | **same 5, same order** |
| Deploy button ready | `{"disabled":false,"text":"Deploy net4sats"}` | **identical** |
| Steps rendered | 11, all pending | **identical 11, same labels** |
| Terminal state | error view `Cannot connect to router via SSH` | **identical** |
| Duration | 42.88 s | 40.28 s |
| Resolution | 1920x1080 H.264 | 1920x1080 H.264 |
| Frames OCR'd | 86 | 81 |
| Frames with `net4sats` | 83 (96.5 %) | 78 (**96.3 %**) |
| Last-frame OCR | `net4sats / … / Setup failed / Cannot connect…` | **byte-identical text** |

Branding ratio, failure mode, step list, and final OCR text all reproduce. The
blocker is a property of the product and the environment, not of one flaky run.

## 4. What this does and does not prove

**Proven:** the build and full test suite pass at this head; the app boots and serves
200 on :8099; the capture toolchain is syntactically valid, runs, and is
deterministic across independent displays and runs.

**NOT proven — and not claimed:** that the de-branded happy path works. It cannot be
shown here: no ref is de-branded (branding is in 96 % of rendered frames) and no
aarch64 router is reachable (deploy dies at step 1). A green build/test/readiness run
says nothing about either. The gate stays open.

## Files

- `.hermes/environment.json` — recipe with `port: 8099` pinned
- `artifacts/gotest-count1.log` — uncached full suite
- `artifacts/verify-full.json` — `hermes verify --json`, `ok: true`
- `artifacts/verify-first-run-port-mismatch.json` — the original failure, kept as-is
- `artifacts/replay-step-log.txt` — run-2 frame-timestamped step log
