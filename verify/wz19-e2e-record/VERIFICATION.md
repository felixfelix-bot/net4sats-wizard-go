# VERIFICATION — t_882b3922 capture lane

Recorded 2026-10-24 on c08r4d0r/x280, tree at `639eb37` (+ the driver edit below).

## 1. Build — PASS

```
$ export PATH=$PATH:/usr/local/go/bin          # go1.22.12, not on PATH by default
$ hermes verify --json --port 8199
{"phase":"build","command":"go build ./...","exitCode":0,"ok":true}
```

## 2. Boot + readiness — PASS

`--port` only moves the readiness poll; the app still binds the manifest's 8099, which
the sibling card (`wz18-wizard` pid 3756902) legitimately holds. Bind explicitly instead:

```
$ hermes verify --json --phase start --port 8199     # -> ok:false, Errno 111
    outputTail: "…running on http://127.0.0.1:8099 … bind: address already in use"
$ WIZARD_BIND=127.0.0.1:8199 hermes verify --json --phase start --port 8199
{"ok":true,"readiness":{"url":"http://127.0.0.1:8199/","ready":true,"statusCode":200}}
```

## 3. Failure attribution — measured by elimination, one variable at a time

`hermes verify --json` run whole (`build → test → start`) reds in the test phase. The red
set is **entirely** (a) other cards' tracked docs, and (b) an environment mistake of mine
that I caught and fixed. Proven by narrowing one variable per run:

| Run | Environment | Remaining failures | Conclusion |
|---|---|---|---|
| A | `go test ./...` as-is | 3: the two docs **+** `TestMainFunctionStructure` + `TestListenAddress` asserting 8099 | — |
| B | `env -u WIZARD_BIND go test -count=1 ./...` | **only** `docs/readiness-…t_05167730….md` (22 lines) | the 8099 assertions were **my** env leak, now gone |
| C | B, then that one doc moved aside (throwaway `git archive` copy) | **`ok net4sats-wizard`** | the residual red is 100% that foreign doc |

`WIZARD_BIND` is consumed by `listenAddress()` (`main.go:53`) and the suite asserts the
8099 default (`handlers_test.go:154`, `security_hardening_test.go:394`), so exporting it
for the *boot* phase silently breaks the *test* phase. Those two phases need different
environments — hence `verify.sh` below.

```
$ PORT=8299 bash verify/wz19-e2e-record/verify.sh
== build ==                                     # go build ./...  -> exit 0
== test (clean env — no WIZARD_BIND) ==         # only foreign-doc leaks remain
== boot + readiness on :8299 ==
{"ok":true,"readiness":{"url":"http://127.0.0.1:8299/","ready":true,"statusCode":200}}
```

Both offenders are **tracked docs committed by other cards** (`22fbea4` t_05167730;
`ab33306` t_0841f821 / worker-admin), leaking internal card ids, board paths, and a
personal fleet hostname. The guard `TestNoInternalProcessLeaksOnShippedSurfaces`
enumerates `docs/` with `os.ReadDir` (`debrand_internal_test.go:79`), so any doc in
`docs/` is scanned as a shipped surface regardless of author.

Proof this lane is not implicated:

- `git diff --stat 5a9b6ae^ HEAD -- docs/` → **empty**. This lane touched no file in `docs/`.
- Guard roots are `docs/`, `portal/`, `admin/`, `rpcd/` (`debrand_internal_test.go:56,79`).
  This lane's evidence lives in `verify/wz19-e2e-record/`, which the guard **never reads**.
- Run C above: `ok net4sats-wizard` once the foreign docs are dismissed.

**The red is theirs; fixing it is their scope.** (Same finding as `t_gate8-19b1b1c6`, which
hit it for a third doc — the untracked
`docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md`. Note `docs/` currently
holds **two** offending docs, so both must go before the suite greens.)

## 4. Driver edit — verified by replay, not by assertion

`capture.cjs` used to depend on being copied next to a `node_modules` tree, because Node
resolves modules from the *file's* directory. It now resolves `@playwright/test` through
`PLAYWRIGHT_MODULE` → local require → two known install paths.

```
$ node --check verify/wz19-e2e-record/capture.cjs          -> SYNTAX_OK
$ cd <repo> && node verify/wz19-e2e-record/capture.cjs    # ran FROM the repo root
t=+0.37s START  navigating to http://127.0.0.1:8199        # browser launched = resolution OK
```

Full replay from the repo root against the same build produced the **same** take:

| | run 1 (`3feadb9`) | replay |
|---|---|---|
| routers discovered | 7 | 7 |
| deploy steps listed | 11 | 11 |
| gating | `{"disabled":false,"text":"Deploy net4sats"}` | identical |
| terminal view | `error` | `error` |
| terminal text | `Setup failed / Cannot connect to router via SSH` | identical |

Deterministic. The replay is a real second take, not a re-encode.
