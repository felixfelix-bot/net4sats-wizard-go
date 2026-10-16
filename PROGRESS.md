# PROGRESS — wz16 e2e integration gate (t_0509c320, pass 4)

Worktree root: `/home/c08r4d0r/wz16-verify` · branch `ci/add-github-actions-workflow-clean`
Verified tip: **`9c37d2e04fe5bfafeb8d210202375ff7f34b773d`** (= PR #35 head, fork `felixfelix-bot/net4sats-wizard-go`)
Pass start: 2026-10-16T04:47Z · host x280

## Findings

- **F1 — the card names a repo that does not exist.** `gh run list --repo net4sats/tollgate-wizard-go`
  → `HTTP 404`, `gh repo view` → "Could not resolve to a Repository", `git ls-remote` → rc 128.
  The real repo is **`net4sats/net4sats-wizard-go`**; the push/CI venue (upstream is read-only and its
  `main` has no `.github/` — which is what this PR adds) is the fork **`felixfelix-bot/net4sats-wizard-go`**.
  Status: reported, not "fixed" (nothing to fix — the literal command is captured verbatim per the card).
- **F2 — upstream PR #35 cannot report checks yet.** `gh pr checks 35 --repo net4sats/net4sats-wizard-go`
  → "no checks reported"; `statusCheckRollup` → `[]`; `repos/.../contents/.github` → 404. The workflow
  is *not yet* on upstream main, so the fork is the only venue where this workflow can run. No check was
  weakened and no workflow file was touched to make anything pass.
- **F3 — tip advanced beyond the previously verified SHA.** The rev3 pass verified `0b5322a`; the branch
  has since gained 4 documentation-only commits and the head is now `9c37d2e`. Source/workflow/goreleaser
  are unchanged (`git diff 0b5322a 9c37d2e` = only `PROGRESS.md` + `REPORT.md`), and CI is green at the
  new head. This pass re-executed every gate live at `9c37d2e` rather than inheriting the old result.
- **F4 — no source change was needed.** All gates were green as received; commits from this pass are
  documentation only.

## Steps (all executed; raw captures under `/home/c08r4d0r/wz16-captures-t_0509c320-p4/`)

| # | Gate | Result | Capture |
|---|------|--------|---------|
| 1 | clean tree + push + ahead/behind | PASS — porcelain 0 lines, push up-to-date, 0/0 vs remote | `01-clean-tree-and-push.txt` |
| 2a | `go vet ./...` | PASS — rc 0 (+ `gofmt -l` clean) | `02-vet-lint-check.txt` |
| 2b | `go test ./... -race` | PASS — `ok net4sats-wizard 108.034s`, 123 RUN / 0 FAIL / 0 SKIP / 0 DATA RACE | `03-test-race.txt`, `04-*` |
| 2c | `golangci-lint run` | PASS — rc 0, v1.61.0 (== CI pin) | `02-vet-lint-check.txt` |
| 2d | `goreleaser check` | PASS — rc 0, "1 configuration file(s) validated" | `02-vet-lint-check.txt` |
| 3 | `goreleaser release --snapshot --clean` | PASS — rc 0, 6 artifacts / 6 distinct targets, checksums 6 == 6, `sha256sum -c` rc 0, 2 spot checks MATCH | `05-goreleaser-snapshot.txt`, `05b-artifacts-and-checksums.txt` |
| 3b | checksum negative controls | PASS — positive rc 0/6 OK; 1-byte flip rc 1/1 FAILED; missing file rc 1 | `08-checksum-controls.txt` |
| 4 | binary `--version` + injection change | PASS — `tollgate-wizard 0.5.1-snapshot.9c37d2e`, exit 0, not `dev`; injection changes; no-ldflags → `dev` | `06-version-smoke.txt`, `06b-version-injection-change.txt` |
| 5 | `gh run list` + CI green | PASS — literal card command 404s (captured); **CI run 38838430526 = success at `9c37d2e`** in the fork | `07-gh-run-list-and-detail.txt`, `09-run-detail-and-log.txt` |

Self-verification: `wz16-90-adhoc-verify.py` re-derives every claim from the captures —
**62 checks / 0 failures**; the tamper control (4 facts corrupted) yields **4 failures / exit 1**,
so the checks are not vacuous (`11-adhoc-verifier-run.txt`, `12-tamper-control.txt`).

## Files touched

- `PROGRESS.md`, `REPORT.md` (this worktree root) — documentation only.
- `verify/wz16-e2e-gate/` — the re-runnable harness (`wz16-0*.sh`, `wz16-9*.py`, captures manifest).
- No `.go`, no `.github/`, no `.goreleaser.yaml`, no test, no module file was modified in the reviewed branch.

## Status: COMPLETE (gate green, CI green, pushed)

Nothing was skipped. The only premise not satisfiable as literally written is F1 (repo name); its exact
error output is captured rather than asserted around.
