## t_0509c320 — pass 4 (rev4) re-verification: Steps 1–4 GREEN, **Step 5 FAILS as literally specified (blocker)**

Author: worker-heavy (fleet offload) · Host x280 · 2026-10-16
Worktree `/home/c08r4d0r/worktrees/wz-e2e-gate-r3` · branch `worker-heavy/t_0509c320-r3` · HEAD `a412b61`
Raw captures: `.evidence-rev4/` (this worktree) · subject under test: HEAD `a412b61` (no source change; docs-only on top of the release commits)

### Finding A — the repo named by the card does not exist (re-measured, not re-cited)
```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
[rc=1]
$ gh api repos/net4sats/tollgate-wizard-go
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/repos#get-a-repository","status":404}
gh: Not Found (HTTP 404)
```
`gh` is authenticated (account `felixfelix-bot`), so this is a real 404. The real repo is
`net4sats/net4sats-wizard-go` (upstream, read-only); push target is the fork `felixfelix-bot/net4sats-wizard-go`.

### STEPS 1–4 — ALL GREEN (verbatim output in the attachment)

**Step 1 — clean tree + pushed**
```
$ git status --porcelain
[porcelain_rc=0] [lines=0]
$ git rev-parse HEAD; git rev-parse @{u}; git rev-list --left-right --count @{u}...HEAD
a412b616c019e1950b523e2340b54a04e2656eb9
a412b616c019e1950b523e2340b54a04e2656eb9
0	0
$ git push
Everything up-to-date
[push_rc=0]
```

**Step 2 — full suite**
```
$ go vet ./...
[go_vet_rc=0]                # no output = clean
$ gofmt -l .
[gofmt_rc=0]                 # no output = clean
$ CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	62.623s
[test_race_rc=0]
# verbose pass: 123 === RUN · 52 --- PASS · 0 --- FAIL · 0 --- SKIP · 0 DATA RACE → ok net4sats-wizard 59.831s
$ golangci-lint run
[golangci_rc=0]              # v1.61.0, no findings
$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
[goreleaser_check_rc=0]
```

**Step 3 — snapshot release**
```
$ goreleaser release --snapshot --clean
  • snapshotting  • building snapshot...  version=0.5.1-snapshot.a412b61
  • building  windows_arm64_v8.0, darwin_arm64_v8.0, linux_amd64_v1, linux_arm64_v8.0, darwin_amd64_v1, windows_amd64_v1
  • archiving 6 archives (4 tar.gz + 2 zip)
  • calculating checksums
  • release succeeded after 4s
[goreleaser_snapshot_rc=0]

artifacts: 6   checksums in SHA256SUMS: 6   (== matches)
$ sha256sum -c SHA256SUMS
<all 6 lines>: OK    [sha256sum_c_rc=0]
independent recompute linux_amd64.tar.gz   → INDEPENDENT_RECOMPUTE_A=MATCH
independent recompute windows_arm64.zip    → INDEPENDENT_RECOMPUTE_B=MATCH
negative control (truncated copy)          → NEGATIVE_CONTROL=FAIL_AS_EXPECTED   (control is non-vacuous)
```

**Step 4 — binary smoke**
```
$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.a412b61        [rc=0]   (NOT 'dev')
$ go build -ldflags "-X main.version=v9.9.9-evidence" … → tollgate-wizard v9.9.9-evidence
$ go build -ldflags "-X main.version=v0.0.1-changed"  … → tollgate-wizard v0.0.1-changed
$ go build (no ldflags)                               … → tollgate-wizard dev
```
(The injected value changes with the flag — proven by two differing values plus a `dev` control.)

### Finding B — STEP 5 **FAILS** as literally specified; pass 3's "CI green" was on a DIFFERENT branch
```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3 --limit 10
(no output — 0 runs)   [rc=0]

$ gh pr view 36 --repo felixfelix-bot/net4sats-wizard-go --json number,headRefName,headRefOid,mergeStateStatus,statusCheckRollup
{"headRefName":"worker-heavy/t_0509c320-r3","headRefOid":"a412b616c019e1950b523e2340b54a04e2656eb9",
 "mergeStateStatus":"DIRTY","number":36,"statusCheckRollup":[]}    ← NO CHECKS

$ gh api …/actions/runs?per_page=100 | head_sha | grep -c a412b61
0        ← not one of the 100 most recent runs is at this branch's head SHA
```
Every recent `success` run is `event=workflow_dispatch` on branch **`ci/add-github-actions-workflow-clean`** — including the
`0b5322a` runs pass 3 cited. So pass 3's Step-5 green was a MANUAL dispatch on a different branch, NOT a CI run on this
branch's pushed head. Measured on the pushed branch, Step 5 is **RED: zero CI runs exist**.

**Root cause** (matches pass-3 Finding 2, now confirmed against the pushed head): a fork PR's `pull_request` runs execute
the BASE repo's workflows. Upstream `main` has no `.github/` at all —
```
$ gh api repos/net4sats/net4sats-wizard-go/contents/.github/workflows
{"message":"Not Found",…,"status":404}      (and: git ls-tree -r --name-only origin/main .github/workflows/ → empty)
```
— so PR #36's check set is empty. The branch itself DOES carry `.github/workflows/{ci,release}.yml` at its head.
Additionally PR #36 is `CONFLICTING`/`DIRTY` (branch-only 12 vs main-only 16 since merge-base `c747aae`), so it cannot be
merged to inherit anything either.

### Status
**Steps 1, 2, 3, 4: PASS.** **Step 5: FAIL (no CI run at the pushed head).** The card's completion gate — "no completion
until pushed and CI verified green" — is therefore **NOT met**. No check was weakened and no workflow file was edited.
Escalating to the manager with a proposed fix (escalate to upstream so `pull_request` checks have a definition; or add
`branches:` to the fork's `push` trigger + push a tie commit so the run is genuinely at the branch head).
The card is in status `triage`, so this worker session cannot transition it.
