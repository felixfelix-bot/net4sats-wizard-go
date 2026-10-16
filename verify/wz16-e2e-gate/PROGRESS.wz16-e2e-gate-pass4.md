# PROGRESS — wz16 e2e integration gate (t_0509c320, pass 4)

Worktree root: `/home/c08r4d0r/wz16-verify` · branch `ci/add-github-actions-workflow-clean`
Verified code tip: **`9c37d2e04fe5bfafeb8d210202375ff7f34b773d`** (= PR #35 head at pass start)
Final tip after this pass's docs: **`fadd82c3423d1764670685592c75873a5e3e4d19`** (= current PR #35 head)
Pass start: 2026-10-16T04:47Z · host x280 · **STATUS: COMPLETE (all gates green, pushed, CI green)**

## Findings

- **F1 — the card names a repo that does not exist.** `gh run list --repo net4sats/tollgate-wizard-go`
  → `HTTP 404`, `gh repo view` → "Could not resolve to a Repository", `git ls-remote` → rc 128.
  The real repo is **`net4sats/net4sats-wizard-go`**; the push/CI venue (upstream is read-only and its
  `main` has no `.github/` — which is what this PR adds) is the fork **`felixfelix-bot/net4sats-wizard-go`**.
  Status: reported verbatim, not "fixed" (nothing to fix — the literal command is captured per the card).
- **F2 — upstream PR #35 cannot report checks yet.** `gh pr checks 35` → "no checks reported";
  `statusCheckRollup` → `[]`; `contents/.github` → 404. The workflow is not yet on upstream main, so the fork
  is the only venue. No check was weakened and no workflow file was touched to make anything pass.
- **F3 — tip advanced beyond the previously verified SHA.** rev3 verified `0b5322a`; the branch has since
  gained documentation commits. Source/workflow/goreleaser are unchanged (`git diff 0b5322a 9c37d2e` = docs
  only). This pass re-executed every gate live at `9c37d2e` rather than inheriting the old result.
- **F4 — no source change was needed.** All gates green as received; commits from this pass are docs only.

## Steps (all executed; raw captures under `/home/c08r4d0r/wz16-captures-t_0509c320-p4/`)

| # | Gate | Result | Capture |
|---|------|--------|---------|
| 1 | clean tree + push + ahead/behind | PASS — porcelain 0 lines, push up-to-date, 0/0 vs remote | `01-clean-tree-and-push.txt` |
| 2a | `go vet ./...` | PASS — rc 0 (+ `gofmt -l` clean) | `02-vet-lint-check.txt` |
| 2b | `go test ./... -race` | PASS — `ok net4sats-wizard 108.034s`, 123 RUN / 0 FAIL / 0 SKIP / 0 DATA RACE | `03-test-race.txt`, `04-*` |
| 2c | `golangci-lint run` | PASS — rc 0, v1.61.0 (== CI pin) | `02-vet-lint-check.txt` |
| 2d | `goreleaser check` | PASS — rc 0, "1 configuration file(s) validated" | `02-vet-lint-check.txt` |
| 3 | `goreleaser release --snapshot --clean` | PASS — rc 0, 6 artifacts / 6 distinct targets, checksums 6 == 6, `sha256sum -c` rc 0, 2 spot checks MATCH | `05-*`, `05b-*` |
| 3b | checksum negative controls | PASS — positive rc 0/6 OK; 1-byte flip rc 1/1 FAILED; missing file rc 1 | `08-checksum-controls.txt` |
| 4 | binary `--version` + injection change | PASS — `tollgate-wizard 0.5.1-snapshot.9c37d2e`, exit 0, not `dev`; injection changes; no-ldflags → `dev` | `06-*` |
| 5 | `gh run list` + CI green | PASS — literal card command 404s (captured); CI green at every pushed tip | `07-*`, `09-*`, `10-*` |

Self-verification: `wz16-90-adhoc-verify.py` re-derives every claim from the captures —
**62 checks / 0 failures**; tamper control (4 facts corrupted) → **4 failures / exit 1** (`11-*`, `12-*`).

## CI at each pushed tip (fork `felixfelix-bot/net4sats-wizard-go`)

| tip | meaning | run | conclusion |
|-----|---------|-----|-----------|
| `9c37d2e` | reviewed **code** tip (= PR #35 head at pass start) | 38838430526 | success |
| `649d90f` | + pass-4 report/progress | 38842920848 | success |
| `5efd948` | + CI-re-run note | 38843145298 | success |
| `b71af0b` | + final-state addendum | 38843393341 | success |
| `fadd82c` | + card evidence bundle (= PR #35 head now) | 38843632861 | success |

The reviewed code tree is byte-identical across all of these:
`git diff 9c37d2e fadd82c -- ':(exclude)*.md' ':(exclude)verify/'` is empty.

## Files touched

- `PROGRESS.md`, `REPORT.md` (worktree root) + `verify/wz16-e2e-gate/` (13 captures, 9 harness scripts,
  report/progress/card-comment copies).
- No `.go`, no `.github/`, no `.goreleaser.yaml`, no test, no module file modified anywhere.

## Card status

Card is in **`triage`** and worker-immutable: `complete` → "unknown id or terminal state";
`promote --force` → "promote only applies to 'todo' or 'blocked'". Evidence attached (3 files) + commented;
**the manager must promote the card** out of `triage`. All gates are green on the pushed tip.
