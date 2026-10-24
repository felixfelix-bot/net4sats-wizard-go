# REPORT — WZ1.9 baseline full test suite + failing-test inventory

- Task: capture baseline full-suite state (READ-ONLY, no fixes, no product commits)
- Repo: net4sats-mvp-v2 @ /home/c08r4d0r/repos/net4sats-mvp-v2
- Branch: worker-heavy/t_gate8-19b1b1c6 (HEAD b56fae5)
- Verdict: **SUITE FAILS — exactly ONE failing test, caused by an UNTRACKED
  scratch doc; the committed tree is fully green.**

---

## 1. Exact command used

```
export PATH=$PATH:/usr/local/go/bin        # `go` is not on PATH in this shell
cd /home/c08r4d0r/repos/net4sats-mvp-v2
go test ./... 2>&1 | tee /tmp/wz19-suite.log
```

Resolution of "canonical full-suite command":
- `make test` does **not** exist. The Makefile (21 lines) defines only nak/nostr
  helpers: nostr-key, nostr-publish-issue, nostr-update-status,
  nostr-query-board, nostr-relays. No build/test/lint target.
- `go.mod`: module `net4sats-wizard`, go 1.22. `go list ./...` -> exactly one
  package: `net4sats-wizard` (flat main package at repo root, tests co-located).
- Therefore `go test ./...` **is** the full suite.
- Toolchain: `go version go1.22.12 linux/amd64` (binary at /usr/local/go/bin/go).

## 2. Exit code + artifacts

| Item | Value |
|---|---|
| Raw log | `/tmp/wz19-suite.log` (exists, 31 lines) |
| Exit code of `go test ./...` | **1** (FAIL) |
| Wall time | 77.428s |
| Packages total / tested | 1 / 1 |
| Tests failed | **1** |
| `t.Errorf` leak findings | 28, all from one file |
| Control: clean HEAD export, same cmd | exit **0**, `ok net4sats-wizard 71.577s` (`/tmp/wz19-head-clean-verify.log`) |

Note on the tee pipeline: `$?` after a pipe reports `tee`'s status, so the exit
code was captured with `${PIPESTATUS[0]}`. Both a piped and a direct run were
consistent (exit 1).

## 3. Last ~30 lines of raw output (verbatim, /tmp/wz19-suite.log)

```
--- FAIL: TestNoInternalProcessLeaksOnShippedSurfaces (0.05s)
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:3: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:4: internal detail leaked (personal fleet hostname): "c08r4d0r"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:5: internal detail leaked (internal branch name): "worker-heavy/"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:22: internal detail leaked (internal board path): "hermes/kanban"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:26: internal detail leaked (internal card id): "t_afa46931"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:34: internal detail leaked (internal card id): "t_afa46931"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:49: internal detail leaked (internal card id): "t_41067114"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:60: internal detail leaked (internal card id): "t_5d6e1bae"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:86: internal detail leaked (internal card id): "t_75852fce"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:86: internal detail leaked (internal board path): "hermes/kanban"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:88: internal detail leaked (internal card id): "t_5d6e1bae"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:99: internal detail leaked (personal fleet hostname): "c08r4d0r"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:137: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:165: internal detail leaked (internal card id): "t_6b512dad"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:169: internal detail leaked (personal fleet hostname): "c08r4d0r"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:176: internal detail leaked (internal card id): "t_2d0c6fa8"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:177: internal detail leaked (internal card id): "t_2d0c6fa8"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:180: internal detail leaked (internal card id): "t_557bb77d"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:191: internal detail leaked (internal card id): "t_2d0c6fa8"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:195: internal detail leaked (internal card id): "t_557bb77d"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:198: internal detail leaked (internal card id): "t_be3261d3"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:224: internal detail leaked (internal decision-gate id): "FELIX"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:237: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:238: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:239: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:271: internal detail leaked (internal card id): "t_cc76c746"
    debrand_internal_test.go:120: docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md:272: internal detail leaked (internal card id): "t_557bb77d"
FAIL
FAIL	net4sats-wizard	77.428s
FAIL
```

## 4. Failing-test inventory

- `net4sats-wizard` — `TestNoInternalProcessLeaksOnShippedSurfaces`
  (defined at `debrand_internal_test.go:109`)

That is the **complete** list. No other test, package, build, or vet failure
appears in the log. Not "all green".

## 5. Root cause (diagnosis only — nothing was fixed)

The single failing test is failing on **one untracked scratch file**:

```
docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md
git status: ?? docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md
git ls-files --error-unmatch <path>: did not match any file(s) known to git
```

Mechanism, from the test source:

- `TestNoInternalProcessLeaksOnShippedSurfaces` scans every path returned by
  `shippedSurfaces(t)` against `internalLeakPatterns()`.
- `shippedSurfaces()` (debrand_internal_test.go:52-103) builds the docs portion
  with a **live directory read**: `os.ReadDir("docs")` (line 79), appending
  every non-dir `*.md` under `docs/`.
- It is *not* a `git ls-files` / tracked-only enumeration. Consequently any
  untracked `.md` parked in `docs/` is treated as a shipped surface and reds
  the guard.
- `docs/` also holds `feed-install.md`, `router-setup.md`,
  `wz1.5-channel-decision.md` (tracked) plus a tracked `docs/evidence/` dir.

Why this is not pre-existing breakage: the culprit is untracked, so it is not
part of the committed tree. Verified non-destructively (the other card's file
was left untouched):

```
rm -rf /tmp/wz19-head-clean && mkdir -p /tmp/wz19-head-clean
git archive HEAD | tar -x -C /tmp/wz19-head-clean
cd /tmp/wz19-head-clean && go test ./...
ok  	net4sats-wizard	71.577s      exit 0
```

So: **committed tree = green; working tree = red, solely due to one untracked doc.**

The guard itself is working as intended (it caught a real leak: internal card
ids, fleet hostname `c08r4d0r`, branch name `worker-heavy/`, board path
`hermes/kanban`, decision-gate id `FELIX`). The gap is in the *enumeration*, not
the detection.

## 6. Scope of the defect class

Only `docs/*.md` is enumerated by raw filesystem read. The walker covers
`portal/`, `admin/`, `rpcd/` (skipping `assets/`) and root `*.go`; the untracked
`verify/*` directories do contain `.json` files but are outside all scanned
roots, so they cannot red this test. At capture time the untracked `docs/`
offender is the **only** head of this class.

## 7. Disposition / recommendations (NOT executed — out of scope, read-only task)

1. Decide the fate of `docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-*.md`:
   it is a task-history artifact, but the test's own contract (comment at
   debrand_internal_test.go:43-46) says task-history artifacts belong in
   `REPORT-*/PROGRESS-*/EVIDENCE-*` at root, which are explicitly **out of
   scope**. Moving it out of `docs/` (or de-identifying it) clears the red.
2. Structurally harden `shippedSurfaces()` to enumerate docs via `git ls-files`
   so untracked scratch files cannot red the debrand guard — or, if the
   live-directory contract is deliberate, state it explicitly in the comment so
   the next worker does not re-diagnose this.
3. Re-run `go test ./...` and confirm exit 0 before any gate claims green.

## 8. Honesty notes

- The suite was NOT made green by this task; nothing was fixed, deleted, moved,
  or committed in the product tree.
- `go` is absent from PATH; every run used `/usr/local/go/bin/go` explicitly.
- Raw log at /tmp/wz19-suite.log is the unedited tee capture; the control log at
  /tmp/wz19-head-clean-verify.log is a separate clean-export run.
