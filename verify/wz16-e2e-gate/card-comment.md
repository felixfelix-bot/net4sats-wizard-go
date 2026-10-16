## t_0509c320 — End-to-end verify (pass 4): VERDICT **PASS** — all 5 steps re-executed live, CI green at the pushed tip

Author: worker-heavy (fleet offload) · Host x280 · 2026-10-16T04:47–05:02Z
Worktree root: `/home/c08r4d0r/wz16-verify` · branch `ci/add-github-actions-workflow-clean`
Subject under test: **`9c37d2e04fe5bfafeb8d210202375ff7f34b773d`** = PR #35 head (OPEN) · fork `felixfelix-bot/net4sats-wizard-go`
Final tip after this pass's documentation commits: **`b71af0b5d3a09849434a56aca8e491a32d3c4404`**
Full report: attachment `REPORT.wz16-e2e-gate-pass4.md` · progress: `PROGRESS.wz16-e2e-gate-pass4.md`
Raw captures + re-runnable harness: `verify/wz16-e2e-gate/` on the branch (13 captures + 9 scripts)

### Gate results (raw output in the attachment and the captures)

**1. Clean tree + pushed** — `git status --porcelain` → 0 lines; `git push` → `Everything up-to-date`;
`git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean` → `0	0`;
remote ground truth `git ls-remote` → `9c37d2e… refs/heads/ci/add-github-actions-workflow-clean`. PASS

**2a `go vet ./...`** → rc 0; `gofmt -l .` → no output; single module (`net4sats-wizard`), so `./...` covers everything. PASS

**2b `CGO_ENABLED=1 go test ./... -race -count=1`** (go1.22.12, gcc 13.3.0) → `ok net4sats-wizard 108.034s`, rc 0;
verbose re-run **123 RUN / 0 FAIL / 0 SKIP / 0 `WARNING: DATA RACE`**. PASS

**2c `golangci-lint run`** (v1.61.0 — exactly the CI pin) → rc 0, no findings. PASS

**2d `goreleaser check`** → rc 0, `1 configuration file(s) validated`. PASS

**3. Snapshot release** `goreleaser release --snapshot --clean` → rc 0, version `0.5.1-snapshot.9c37d2e`.
Artifacts (6, 6 distinct targets):

```
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_amd64.zip
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip
```

`wc -l < SHA256SUMS` → `6` (**== artifact count**); `sha256sum -c SHA256SUMS`:

```
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip: OK
SHA256SUM_C_REAL_RC=0
```

The two mandated spot checks are inside that run (one archive = `linux_amd64.tar.gz` → OK; one non-linux =
`darwin_arm64.tar.gz` → OK), and 3 digests were independently recomputed → `MATCH`. PowerShell is not on
this host, so the Windows `Get-FileHash` equivalent is `sha256sum` on the zips (`windows_amd64.zip` → MATCH). PASS

**3b. Non-vacuity controls** — positive rc 0 (6/6 OK); flipping 1 byte of an **explicit** path copy flips
`linux_amd64.tar.gz` to `FAILED` (rc 1); removing an artifact gives `FAILED open or read` (rc 1). PASS

**4. Binary smoke** — the linux/amd64 binary extracted from the snapshot archive:

```
$ ./net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.9c37d2e
exit_code=0
FOUND_DEV=no (good)
```

Injected value demonstrably changes: `-X main.version=v9.9.9-manual-p4` → `tollgate-wizard v9.9.9-manual-p4`,
`-X main.version=v0.0.0-second-different-value-p4` → that value, no-ldflags control → `tollgate-wizard dev`.
`go tool nm` locates `main.version`; source anchors `main.go:52` (`var version = "dev"`), `main.go:59`
(`versionString()`), `main.go:641` (the `--version`/`-version` case). PASS

**5. CI green.** The literal card command **fails because that repo does not exist** — reported verbatim,
not papered over:

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
literal_real_rc=1

$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)

$ git ls-remote https://github.com/net4sats/tollgate-wizard-go
remote: Repository not found.        lsremote_real_rc=128
```

`gh` is authenticated as `felixfelix-bot` over HTTPS, so this is a real 404. The real repo is
`net4sats/net4sats-wizard-go`; upstream `main` has **no `.github/`** (`gh api …/contents/.github` → 404;
`gh pr checks 35` → "no checks reported"; `statusCheckRollup` → `[]`), which is exactly what this PR adds,
so the **fork is the CI venue** until it lands. CI at the pushed tip:

```
$ gh run view 38843393341 --json conclusion,headSha,status,url
{"conclusion":"success","headSha":"b71af0b5d3a09849434a56aca8e491a32d3c4404","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38843393341"}

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843393341	2m0s	2026-10-16T03:02:00Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843145298	2m0s	2026-10-16T02:58:21Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838430526	2m1s	2026-10-16T02:03:24Z
```

Run 38838430526 is green at the reviewed **code** tip `9c37d2e`; runs 38842920848 / 38843145298 / 38843393341
are green at the subsequent documentation tips. Gate steps observed in the job log:
`✓ lint in 20s` (golangci-lint v1.61.0 → "found no issues") and `✓ test in 1m57s`
(`go build ./...` / `go vet ./...` / `go test ./... -race -count=1` → `ok net4sats-wizard 82.366s`). PASS

### Card premises that are wrong on this host (unchanged from passes 1–3)

1. **`net4sats/tollgate-wizard-go` does not exist** (404 / "Could not resolve to a Repository" / ls-remote rc 128).
   The repo is `net4sats/net4sats-wizard-go`; the push + CI venue is the fork `felixfelix-bot/net4sats-wizard-go`.
2. **Upstream cannot report CI on this PR at all** until the workflow lands (no `.github/` upstream).
   **No check was weakened and no workflow was edited** to make anything pass.
3. `golangci-lint run` needs the CI-matching go1.22.12 (under host go1.24 it reports `unsupported version: 2`).

### Self-verification

`wz16-90-adhoc-verify.py` re-derives every claim above from the captures — **62 checks / 0 failures**; the
tamper control corrupts 4 load-bearing facts and gets **4 failures / exit 1**, so the checks are not
vacuous. Secret scan over captures and the committed docs: **clean**.

### State / hygiene

- **No source, test, workflow, module, or goreleaser file was changed.** All gates were green as received.
  This pass's commits are documentation + harness only; the reviewed code tree is unchanged
  (`git diff 9c37d2e <tip> -- ':(exclude)*.md' -- ':(exclude)verify/'` is empty).
- Final: `git status --porcelain` → 0 lines; ahead/behind → `0	0`; tip pushed to the fork and CI green there.

### Card status

The card is in status **`triage`**, so a worker session cannot close it: `complete`/`promote --force` are
refused on triage cards. Evidence is attached + commented above — **the manager must promote the card**
out of `triage`. Every gate in the card's acceptance criteria is green on the pushed tip; nothing is left
for a worker to do on this card.
