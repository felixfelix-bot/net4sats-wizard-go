# PROGRESS — wz16 e2e integration gate (pass 5, independent re-run)

## Finding

Card task: end-to-end verify (clean tree → full suite → snapshot → `--version` smoke → push → CI green)
for the wizard Go repo, after the CI workflow / goreleaser config / docs are committed.

Two card-string discrepancies, both checked against live GitHub rather than assumed:

1. Repo `net4sats/tollgate-wizard-go` **does not exist** (HTTP 404 / GraphQL unresolved). The real target
   is `felixfelix-bot/net4sats-wizard-go` (public fork; upstream `net4sats/net4sats-wizard-go`).
2. The card's expected `--version` line `tollgate-wizard <version>` **is correct** — pinned by
   `main.go:55 const programName = "tollgate-wizard"`. My first smoke script wrongly asserted the
   repo-derived name `net4sats-wizard` and mis-reported `FORMAT: FAIL`; that was my script's bug.
   Corrected, re-run, capture replaced.

Prior state found: a complete pass-4 bundle already existed at `/home/c08r4d0r/wz16-verify` on
`ci/add-github-actions-workflow-clean` @ `64008dc2…`. This pass re-executed all five step groups live
rather than replaying those captures.

## Status — ALL STEPS GREEN, pushed, CI green at the exact tip

| Step | Result | Evidence (this pass) |
|---|---|---|
| 1 clean tree + push + 0/0 | PASS — porcelain 0 lines, up-to-date, `0	0`, ls-remote == HEAD | `pass5/01-*.txt` |
| 2 `go vet` / `golangci-lint` / `goreleaser check` | PASS — 0 / 0 / validated | `pass5/02-*.txt` |
| 2 `go test ./... -race` | PASS — 52 pass, 0 fail, 0 data race, rc 0 | `pass5/03`,`04*` |
| 3 snapshot 6 artifacts + archives + checksums | PASS — 6 archives (3 OS × 2 arch), parity 6==6 | `pass5/05*` |
| 3 checksum `sha256sum -c` | PASS — 6/6 OK, rc 0 (all, not just 2) | `pass5/05b-*.txt` |
| 4 `--version` smoke | PASS — `tollgate-wizard 0.5.1-snapshot.64008dc`, exit 0 | `pass5/06-*.txt` |
| 4 injection control + negative control | PASS — ldflag changes value; uninjected → `dev` | `pass5/06-*.txt` |
| 5 CI green via `gh run list` | PASS — run 38844256545 headSha == HEAD, lint+test success | `pass5/07`,`07b` |
| 5 code-tree identity across CI range | PASS — 0 `*.go`/`.goreleaser.yaml`/`.github/**` diffs | `pass5/08-*.txt` |

Commands actually executed this pass (host x280, go1.22.12, golangci-lint 1.61.0, goreleaser):
`git status --porcelain` · `git push` · `git rev-list --left-right --count` · `git ls-remote` ·
`go vet ./...` · `golangci-lint run` · `goreleaser check` · `go test ./... -race` ·
`go test ./... -race -v` · `goreleaser release --snapshot --clean` · `sha256sum -c SHA256SUMS` ·
`./dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version` · two control builds ·
`gh run list` (×3 repos) · `gh run view` · `gh repo view` · `git ls-tree` · `git diff --name-only`.
Nothing was skipped and no workflow check was weakened.

## Honest caveats (recorded, not hidden)

- PR #29 (`ci/add-github-actions-workflow-clean → main`, `felixfelix-bot/…`) is OPEN at
  `headRefOid 64008dc2…` but has an **empty `statusCheckRollup`** and is `CONFLICTING`. CI green is
  therefore established from the workflow runs (`workflow_dispatch`), **not** PR check-rollup. Both are
  true statements; they are not the same claim.
- No Windows host here, so `Get-FileHash` was not literally run. The equivalent integrity proof is
  `sha256sum -c` over the full `SHA256SUMS` (same SHA-256 comparison).
- The CI run at the tip was produced before this pass; `git push` returned "Everything up-to-date", so
  the tip CI green is unchanged by this pass. This pass adds only docs/evidence commits.

## Files touched (this pass)

```
verify/wz16-e2e-gate/REPORT.wz16-e2e-gate-pass5.md   (new — full report)
verify/wz16-e2e-gate/PROGRESS.wz16-e2e-gate-pass5.md (new — this file)
verify/wz16-e2e-gate/pass5/                          (new — 12 raw captures + 8 scripts)
PROGRESS.md                                          (updated)
REPORT.md                                            (updated)
```

## Card status

Card is `triage` and worker-immutable (`complete`/`promote --force` refuse). Evidence attached + comment;
the manager must promote the card out of `triage`. All gates are green on the pushed tip `64008dc2`.

---

## Hand-off note on the evidence-tip offset

The ledger above ends at `c71b2f8`; the next pushed tip `db00314b4cde3b4f939fea653e1f50212955ceed`
was also confirmed CI-green (`gh run list` run `38852038477`, `conclusion: success`,
`headSha == db00314b…` == HEAD == remote).

Each further evidence/docs commit necessarily shifts the tip by one commit. That does not weaken the
result, because the **code tree is unchanged** across every one of these tips:

```
git diff --name-only <code-tip> <any-later-tip> -- '*.go' '.goreleaser.yaml' '.github/**'   # -> empty
```

So the CI-green verdict transfers: the tip is green, and the tip contains the identical `.go`,
`.goreleaser.yaml`, and `.github/**` blobs that produced the snapshot, checksum, and `--version`
evidence. The definitive per-tip check is always
`gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 1 --json headSha,conclusion`
against `git rev-parse HEAD`.
