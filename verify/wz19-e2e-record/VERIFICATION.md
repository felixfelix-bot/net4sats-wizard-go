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

## 3. Full suite — RED, and NOT caused by this lane

```
$ go test -count=1 ./... | grep -oE 'docs/[A-Za-z0-9._-]+\.md' | sort | uniq -c
  22 docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md   (then, after quarantine:)
  22 docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md
```

Both offenders are **tracked docs committed by other cards** (`22fbea4` t_05167730;
`ab33306` t_0841f821 / worker-admin). Both leak internal card ids and board paths. The
guard `TestNoInternalProcessLeaksOnShippedSurfaces` enumerates `docs/` with
`os.ReadDir` (`debrand_internal_test.go:79`), so any doc in `docs/` is scanned as a
shipped surface regardless of author.

Proof this lane is not implicated:

- `git diff --stat 5a9b6ae^ HEAD -- docs/` → **empty**. This lane touched no file in `docs/`.
- The guard's scanned roots are `docs/`, `portal/`, `admin/`, `rpcd/`
  (`debrand_internal_test.go:56,79`). This lane's evidence lives in
  `verify/wz19-e2e-record/`, which the guard **never reads**.
- Removing only those two foreign docs, in a throwaway `git archive HEAD` copy:

```
$ mv docs/readiness-….md <out of dir>; mv docs/LOCK-ANNOUNCE-….md docs/evidence/
$ go test -count=1 ./...
ok  net4sats-wizard  56.545s
```

**Conclusion: the committed tree passes once two other cards' docs are dismissed. The
red is theirs; fixing it is their scope, not this lane's.** (Same finding as
`t_gate8-19b1b1c6`, which hit it for a third doc — the untracked
`docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md`.)

Pitfall for the next lane: exporting `WIZARD_BIND` in the shell **breaks
`TestListenAddress` and `TestSecurityHardening`** (they assert the 8099 default). Always
run the suite with `env -u WIZARD_BIND`.

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
