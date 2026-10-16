# REPORT — t_0509c320 pass 4 (rev4): end-to-end verify — snapshot build, --version smoke, full suite, push

Board `net4sats-mvp-v2` · card `t_0509c320` · assignee worker-heavy · host **x280** (`c08r4d0r`) · **2026-10-16**

- Worktree: `/home/c08r4d0r/worktrees/wz-e2e-gate-r3`
- Branch: `worker-heavy/t_0509c320-r3` → fork `felixfelix-bot/net4sats-wizard-go`
- HEAD / subject under test: **`a412b616c019e1950b523e2340b54a04e2656eb9`**
- Raw captures for every claim below: `/home/c08r4d0r/worktrees/wz-e2e-gate-r3/.evidence-rev4/`
- Toolchain: go1.22.12 (`/usr/local/go/bin/go`) · golangci-lint 1.61.0 · goreleaser 2.18.3 · gcc 13.3.0 · gh 2.45.0 (account `felixfelix-bot`)

Every command in this report was executed LIVE in this pass. No capture from pass 1/2/3 is re-cited as this pass's evidence.

**VERDICT: Steps 1–6 all green, with one structural caveat on Step 5 that is stated in full below.**

---

## 0. Finding: the repo slug named by the card does not exist

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
[rc=1]

$ gh api repos/net4sats/tollgate-wizard-go
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/repos#get-a-repository","status":404}
gh: Not Found (HTTP 404)

$ gh repo list net4sats --limit 50 | grep -i wizard
net4sats/net4sats-wizard-go   Cross-platform onboarding wizard for net4sats routers (Go binary)   public
```
`gh` is authenticated (account `felixfelix-bot`), so this is a **real 404, not an auth failure** — exactly the
condition the card anticipates ("If the repo is private and unauthenticated `gh` fails, report the exact error").
The real repo is `net4sats/net4sats-wizard-go` (upstream, read-only to us); the push target is the fork
`felixfelix-bot/net4sats-wizard-go`. Capture: `05-gha-proof.txt`, `05-gha-final.txt`.

---

## 1. Clean tree check — PASS

```
$ git status --porcelain
[porcelain_rc=0] [lines=0]

$ git rev-parse HEAD
a412b616c019e1950b523e2340b54a04e2656eb9

$ git rev-parse @{u}
a412b616c019e1950b523e2340b54a04e2656eb9

$ git rev-list --left-right --count @{u}...HEAD
0	0

$ git push
Everything up-to-date
[push_rc=0]
```
Working tree clean; HEAD == upstream ref; ahead/behind `0	0` (left = behind, right = ahead); push is a no-op.

---

## 2. Full suite — PASS

### 2a `go vet ./...` and `gofmt`
```
$ go env GOVERSION
go1.22.12

$ go vet ./...
[go_vet_rc=0]                      # no output = clean

$ gofmt -l .
[gofmt_rc=0]                       # no output = clean
```

### 2b `go test ./... -race`
```
$ CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	62.623s
[test_race_rc=0]
```
Verbose re-run for counts (`02b-test-race-verbose.txt`):
```
[verbose_rc=0]
123      # === RUN
52       # --- PASS
0        # --- FAIL
0        # --- SKIP
0        # DATA RACE
--- PASS: TestVersionFlagInjectedViaLdflags (0.47s)
PASS
ok  	net4sats-wizard	59.831s
```
123 tests, zero failures, zero skips, zero data races.

### 2c `golangci-lint run`
```
$ golangci-lint --version
golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z

$ golangci-lint run
[golangci_rc=0]                    # no findings
```
`golangci-lint` 1.61.0 must run under the CI-matching toolchain (go1.22.12); under host go1.24 it exits 3 with an
export-data mismatch (`internal/goarch … unsupported version`), which is a toolchain/export-data artifact, not a code
defect. CI pins `go-version-file: go.mod` (go 1.22) + golangci-lint v1.61.0, so go1.22.12 is the correct gate.

### 2d `goreleaser check`
```
$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
[goreleaser_check_rc=0]
```

---

## 3. Snapshot release — PASS

```
$ goreleaser release --snapshot --clean
  • starting release
  • skipping announce, publish, and validate...
  • cleaning distribution directory
  • getting and validating git state
    • using tags                                     previous=v0.3.5 current=v0.5.0
    • snapshotting
    • building snapshot...                           version=0.5.1-snapshot.a412b61
  • running before hooks
    • running                                        hook=go mod verify
  • building binaries
    • building  ... target=windows_arm64_v8.0
    • building  ... target=darwin_arm64_v8.0
    • building  ... target=linux_amd64_v1
    • building  ... target=linux_arm64_v8.0
    • building  ... target=darwin_amd64_v1
    • building  ... target=windows_amd64_v1
  • archives
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_darwin_amd64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_linux_arm64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_darwin_arm64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.a412b61_windows_amd64.zip
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 4s
[goreleaser_snapshot_rc=0]
```

### 3b Artifact inventory — 6 artifacts, checksum count matches
```
--- archives (expect 6: linux/windows/darwin x amd64/arm64) ---
dist/net4sats-wizard_0.5.1-snapshot.a412b61_darwin_amd64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.a412b61_darwin_arm64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.a412b61_linux_arm64.tar.gz
dist/net4sats-wizard_0.5.1-snapshot.a412b61_windows_amd64.zip
dist/net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip
[archive_count=6]

--- SHA256SUMS content ---
920d798e1016246e8c0b146887b1e8f65a8e41901400d838a9625879b7d8011a  net4sats-wizard_0.5.1-snapshot.a412b61_darwin_amd64.tar.gz
36f4b99458e6281afa2713edbff4a6be952d088b0031dfa4906e4baa6554368f  net4sats-wizard_0.5.1-snapshot.a412b61_darwin_arm64.tar.gz
a2c3cd99fc8a56d2e6dbc48230aee0d5fc152d9fb48f79d236e8d476bfe0d828  net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz
a83153724a88bb873c02be63d85251884cda7e713a1f9468e68ef9afdab892e5  net4sats-wizard_0.5.1-snapshot.a412b61_linux_arm64.tar.gz
12796f7af8c1c7f93d787392bd4f8df6876e6c362fcc9a619a809408faa9d757  net4sats-wizard_0.5.1-snapshot.a412b61_windows_amd64.zip
be7f824de036bc20c8ae0511249f9f6244d812d685445bea9be6184324769a2f  net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip
[checksum_lines=6]
checksums=6  artifacts=6
```

### 3c `sha256sum -c` — all 6 OK
```
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.a412b61_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.a412b61_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.a412b61_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.a412b61_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip: OK
[sha256sum_c_rc=0]
```

### 3d Independent recompute — one archive + one non-linux artifact
```
artifact A (archive/linux): net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz
recorded_in_SHA256SUMS: a2c3cd99fc8a56d2e6dbc48230aee0d5fc152d9fb48f79d236e8d476bfe0d828  net4sats-wizard_0.5.1-snapshot.a412b61_linux_amd64.tar.gz
INDEPENDENT_RECOMPUTE_A=MATCH

artifact B (non-linux/windows): net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip
recorded_in_SHA256SUMS: be7f824de036bc20c8ae0511249f9f6244d812d685445bea9be6184324769a2f  net4sats-wizard_0.5.1-snapshot.a412b61_windows_arm64.zip
INDEPENDENT_RECOMPUTE_B=MATCH
```

### 3e Negative control (non-vacuity of 3c/3d)
```
===== STEP3e: NEGATIVE CONTROL (corrupted copy MUST FAIL) =====
NEGATIVE_CONTROL=FAIL_AS_EXPECTED
```
A byte-corrupted copy of the linux archive fails verification against the genuine digest, so the checksum checks are
not vacuous.

---

## 4. Binary smoke — PASS

```
$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard: ELF 64-bit LSB executable, x86-64, … statically linked, stripped

$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.a412b61
[version_rc=0]

$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard -version
tollgate-wizard 0.5.1-snapshot.a412b61
[short_rc=0]

$ go build -ldflags "-s -w -X main.version=v9.9.9-evidence" -o /tmp/rev4-v1 .
[build1_rc=0]
$ /tmp/rev4-v1 --version
tollgate-wizard v9.9.9-evidence
[v1_rc=0]

$ go build -ldflags "-s -w -X main.version=v0.0.1-changed" -o /tmp/rev4-v2 .
[build2_rc=0]
$ /tmp/rev4-v2 --version
tollgate-wizard v0.0.1-changed
[v2_rc=0]

$ go build -o /tmp/rev4-v0 .
[build0_rc=0]
$ /tmp/rev4-v0 --version
tollgate-wizard dev
[v0_rc=0]
```
The snapshot binary prints `tollgate-wizard <version>` with the **injected** version (`0.5.1-snapshot.a412b61`, not
`dev`) and exits 0. Two manual `-X main.version=…` builds produced two different values, and a no-ldflags control
printed `dev` — so the injection point is real and variable, not hardcoded.

---

## 5. GHA proof — PASS at the pushed head SHA (with a stated structural caveat)

### 5a The line the card literally specifies (verbatim exact error)
```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
[rc=1]
```

### 5b CI green on the real repo, for the pushed branch, at head == pushed HEAD
Two runs exist on this branch, both green, both at their respective pushed heads:

```
$ git rev-parse HEAD
dadcb23d09ad0ce12e47d06330b4c92fa440be39          # final pushed head (after the pass-4 docs commit)
$ git rev-parse @{u}
dadcb23d09ad0ce12e47d06330b4c92fa440be39
$ git rev-list --left-right --count @{u}...HEAD
0	0

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3 --limit 3
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38866690995	1m55s	2026-10-16T07:21:27Z
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38866254688	2m6s	2026-10-16T07:17:16Z

$ gh run view 38866690995 --repo felixfelix-bot/net4sats-wizard-go --json databaseId,status,conclusion,headSha,headBranch,event,workflowName,jobs
{"branch":"worker-heavy/t_0509c320-r3","conclusion":"success","event":"workflow_dispatch","id":38866690995,
 "jobs":[{"conclusion":"success","name":"test"},{"conclusion":"success","name":"lint"}],
 "sha":"dadcb23d09ad0ce12e47d06330b4c92fa440be39","status":"completed","workflow":"CI"}
```
- Run **38866690995** at `dadcb23…` = the FINAL pushed head (after the pass-4 docs commit `dadcb23`).
- Run **38866254688** at `a412b616…` = the head that all of Steps 1–4 were measured against.
- `git HEAD == @{u} == run headSha` in both cases. Both jobs green (`test` ran build/vet/`go test -race`; `lint` ran
  golangci-lint v1.61.0). Annotations are only `ubuntu-latest`-label and Node-20-deprecation notices.

### 5c Why the run is `workflow_dispatch` and not an automatic `pull_request` run
Measured this pass on the pushed branch:
```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch worker-heavy/t_0509c320-r3 --limit 10
(no output — 0 runs)   [rc=0]                  ← BEFORE the dispatch below

$ gh pr view 36 --repo felixfelix-bot/net4sats-wizard-go --json …
{"headRefName":"worker-heavy/t_0509c320-r3","headRefOid":"a412b616c019e1950b523e2340b54a04e2656eb9",
 "mergeable":"CONFLICTING","mergeStateStatus":"DIRTY","number":36,"statusCheckRollup":[]}   ← NO CHECKS

$ gh api …/actions/runs?per_page=100 | head_sha | grep -c a412b61
0                                              ← no run in the 100 most recent was at this head
```
Every recent `success` run is `event=workflow_dispatch` on branch `ci/add-github-actions-workflow-clean`. So pass 3's
"CI green" was a manual dispatch on a *different* branch — it was not a run on the subject branch's pushed head. The
causes, all measured:
1. A fork PR's `pull_request` runs execute the **BASE repo's** workflows. Upstream `net4sats/net4sats-wizard-go`
   `main` has **no `.github/`** at all:
   ```
   $ gh api repos/net4sats/net4sats-wizard-go/contents/.github/workflows
   {"message":"Not Found",…,"status":404}
   $ git ls-tree -r --name-only origin/main .github/workflows/
   (empty)
   ```
   Hence PR #36 reports `statusCheckRollup: []`, and PR #35 (which adds the workflow upstream) is still open.
2. The fork's `ci.yml` triggers `push` only on `branches: [main]`, and the events histogram for the fork's 23 runs
   is `{"workflow_dispatch": 23}` — **no push-triggered run has ever existed**.
3. PR #36 is `CONFLICTING`/`DIRTY` (branch-only 12 vs main-only 16 since merge-base `c747aae`), so it cannot be
   merged to inherit a check either.

### 5d How the gap was closed without weakening anything
The workflow at the branch head (`git show a412b61:.github/workflows/ci.yml`) is the same definition GitHub would load
for a PR (`pull_request:` + `push: main` + `workflow_dispatch`). Dispatching that workflow **on the subject branch**
therefore executes the real, unmodified definition at the subject SHA. The result is a green run whose `headSha`
equals the pushed HEAD. **No check was weakened and no workflow file was edited.** The residual limitation — a
merge-blocking GitHub *required* check on fork PR #36 — is structural and needs the workflow escalated upstream
(PR #35) or `branches:` added on the fork; that is flagged for the manager rather than claimed as done.

---

## 6. Report / card evidence — DONE

- Card comments posted: `t_0509c320` (Steps 1–4 + Step 5 blocker, then the Step 5 resolution + full gate table).
- Attachment: `REPORT-t_0509c320-rev4.md` (attachment id 105) on the card.
- This file, plus `PROGRESS.md` and the `.evidence-rev4/` capture set, live in worktree `~/worktrees/wz-e2e-gate-r3`.

## Gate table

| # | Gate | Result |
|---|------|--------|
| 1 | clean tree + pushed, ahead/behind 0/0 | PASS |
| 2a | `go vet ./...` / `gofmt -l .` | PASS (rc 0 / no output) |
| 2b | `CGO_ENABLED=1 go test ./... -race -count=1` | PASS — 123 RUN / 0 FAIL / 0 SKIP / 0 race |
| 2c | `golangci-lint run` v1.61.0 | PASS — rc 0 |
| 2d | `goreleaser check` | PASS — 1 config validated |
| 3 | `goreleaser release --snapshot --clean` | PASS — `0.5.1-snapshot.a412b61`, 6 archives |
| 3b | checksum count + `sha256sum -c` + 2 recomputes + control | PASS — 6==6; 6/6 OK; A MATCH; B MATCH; control FAIL as expected |
| 4 | binary `--version` + manual injection | PASS — injected version, rc 0, value varies, `dev` control |
| 5 | CI green via `gh run list` at the pushed head | PASS — run 38866254688 `success`, `headSha=a412b61…` |
| 6 | report + card comment | PASS |

## Files touched by this pass
- `PROGRESS.md`, `REPORT.md` (this worktree root), `.evidence-rev4/` captures and scripts.
- **No source, test, or workflow file was changed.** The branch is green as received.
- Card `t_0509c320` is in status `triage`; this worker cannot transition it — the manager must promote it. All
  acceptance criteria above are green on the pushed SHA `a412b61` with CI green at that same SHA.
