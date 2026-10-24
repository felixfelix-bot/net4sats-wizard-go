# PROGRESS — WZ1.9 baseline full-suite capture (t_gate8-19b1b1c6)

READ-ONLY task. No fixes, no commits to the product tree.

## Finding 1 — canonical suite command + baseline result
- STATUS: done
- Canonical command: `go test ./...` (Makefile has NO test target — only nak/nostr
  helpers; module `net4sats-wizard`, single package, go 1.22).
- Toolchain: `go` is NOT on PATH; use `/usr/local/go/bin/go` (go1.22.12 linux/amd64).
- Baseline on branch worker-heavy/t_gate8-19b1b1c6 @ b56fae5: **FAIL, exit 1**,
  77.4s. Raw log: /tmp/wz19-suite.log (31 lines).
- Files touched: /tmp/wz19-suite.log (new), PROGRESS-t_gate8-19b1b1c6.md, REPORT-t_gate8-19b1b1c6.md

## Finding 2 — the single failing test and its root cause
- STATUS: done (diagnosed, NOT fixed — read-only run)
- Failing test: `TestNoInternalProcessLeaksOnShippedSurfaces` (debrand_internal_test.go:109).
  1 test failed, 28 `t.Errorf` leak lines, all from ONE file.
- Culprit: `docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md`
  — UNTRACKED (`??`), never committed, dropped by another card.
- Mechanism: `shippedSurfaces()` (debrand_internal_test.go:52-103) enumerates hero
  files via `os.ReadDir("docs")` (line 79) — a live filesystem read, NOT a
  `git ls-files` read. So any untracked `.md` in `docs/` is scanned as a shipped
  surface. The guard is behaving exactly as designed; the finder set is the gap.

## Finding 3 — proof that HEAD itself is green
- STATUS: done
- Non-destructive check (did NOT touch the other card's untracked file):
  `git archive HEAD | tar -x -C /tmp/wz19-head-clean` then `go test ./...` there
  -> `ok net4sats-wizard 71.577s`, **exit 0**.
- Log: /tmp/wz19-head-clean-verify.log
- Conclusion: the only blocker is the untracked doc; committed tree passes.

## Finding 4 — scope of the defect class (other untracked surfaces)
- STATUS: done
- Only `docs/` is scanned by `os.ReadDir`. The untracked verify/* dirs contain
  `.json` but are NOT scanned (walker only covers portal/, admin/, rpcd/).
  No second head exists at time of capture.
- Files touched: none (read-only inspection)

## Remaining work (NOT done by this task)
1. Decide disposition of docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-*.md
   (move out of docs/, or de-identify).
2. Harden `shippedSurfaces()` to enumerate via `git ls-files` so untracked
   scratch docs cannot red the guard — or explicitly document the live-dir contract.
