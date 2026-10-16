## Pass-6 — CI went RED once at the final tip, then GREEN on re-run (flake, diagnosed)

Recorded because a red run happened during this pass and hiding it would be dishonest.

At `2026-10-16T15:27:22Z`, run **`38924071342`** on the final tip `e46f58f9dd6a23eb51e7ab70447e939a2122dc8e`
concluded **`failure`**: job `test` succeeded, job `lint` failed at step `Run golangci/golangci-lint-action@v6`.
The verbatim error (capture `11-lint-failure-log.txt`):

```
##[error]Failed to run: Error: Command failed: /home/runner/golangci-lint-1.61.0-linux-amd64/golangci-lint config verify
Failed executing command with error: [.golangci.yml] validate: compile schema: jsonschema
  https://golangci-lint.run/jsonschema/golangci.v1.61.jsonschema.json compilation failed:
  Get "https://golangci-lint.run/jsonschema/golangci.v1.61.jsonschema.json":
  context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

**This is not a lint finding and not a defect in the code.** Three independent reasons:

1. The linter never ran — the action died in its `config verify` pre-step, which validates `.golangci.yml`
   against a JSON schema it downloads over HTTP from `golangci-lint.run`. The failure is an outbound HTTP timeout.
2. The **identical** code tree passed minutes earlier: run `38923652846` @ `ec84984` = `success`, and
   `git diff --name-only e85720c e46f58f -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'` is empty
   (`CODE_DIFF_LINES=0`). Nothing in the diff can change a lint outcome.
3. The same pinned toolchain and config pass locally in this session: `golangci-lint 1.61.0` → `GOLANGCI_RC=0`,
   zero findings (§2).

**Response — re-run, not a weakened check:**

```
$ gh run rerun 38924071342 --failed --repo felixfelix-bot/net4sats-wizard-go
RERUN_RC=0
[10s] status=completed  [{"n":"lint","c":"success"},{"n":"test","c":"success"}]
FINAL: {"status":"completed","conclusion":"success"}      # headSha still e46f58f (the pushed tip)
```

**No workflow file was edited, no check was weakened or removed, no re-push was needed.** The gate is green at the
pushed tip by a faithful run of the unmodified `.github/workflows/ci.yml`.

**Recommendation to the manager (deliberately NOT applied here).** The `lint` job carries an avoidable network
dependency: `golangci-lint-action@v6` runs `golangci-lint config verify`, which fetches the schema from
`golangci-lint.run`. It failed once in ~8 runs on this branch today. A reliability fix — *not* a weakening — would be
to run the same pinned linter directly (`golangci-lint run` at v1.61.0) instead of through the action's
config-verify path, or to pin an action version that does not fetch the schema, or to add a bounded retry. I did not
apply it because `.github/workflows/ci.yml` is the very artifact this gate exists to verify; editing it here would
change the subject under test and invalidate this pass's evidence. Escalated instead.

---

## Pass-6 close-out — CI green at the pushed evidence tip

The evidence commit of this pass is **`ec849846d4894a5e6610ba63187f6aaa678a788e`**, pushed to
`ci/add-github-actions-workflow-clean`. CI was dispatched on it and is green **at that exact sha**:

```
$ gh workflow run CI --repo felixfelix-bot/net4sats-wizard-go --ref ci/add-github-actions-workflow-clean
DISPATCH_RC=0

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38923652846	1m59s	2026-10-16T15:24:17Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852530642	2m0s	2026-10-16T04:48:45Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852255804	2m2s	2026-10-16T04:45:27Z

$ gh run view 38923652846 --repo felixfelix-bot/net4sats-wizard-go --json …
{"id":38923652846,"wf":"CI","branch":"ci/add-github-actions-workflow-clean","sha":"ec849846d4894a5e6610ba63187f6aaa678a788e",
 "event":"workflow_dispatch","status":"completed","conclusion":"success",
 "jobs":[{"name":"lint","conclusion":"success"},{"name":"test","conclusion":"success"}]}

$ gh run view 38923652846 --json jobs   # step detail — the real commands CI ran
JOB lint: success
    success	Run actions/checkout@v4
    success	Run golangci/golangci-lint-action@v6
JOB test: success
    success	Run go build ./...
    success	Run go vet ./...
    success	Run go test ./... -race -count=1

local HEAD = ec849846d4894a5e6610ba63187f6aaa678a788e
CI run headSha = ec849846d4894a5e6610ba63187f6aaa678a788e   -> SHA_MATCH = OK

=== code-tree identity: evidence tip vs the code tip the local suite/snapshot ran on ===
$ git diff --name-only e85720cd39a41d4af8bdff5172752a827e25bb86 HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'
CODE_DIFF_LINES=0
HEAD^{tree} = c0fb58376089bd2ce2fb298726e5dec1869865ec     main.go blob = 30b35e02055417e35adc4695a08283ea9d0ad92c
.goreleaser.yaml blob = 0741e31079072bb26135fb605e689c66233d7434
```

The diff between the code tip whose local suite/snapshot/`--version` evidence was produced (`e85720c`) and the pushed
tip `ec84984` is **docs/evidence only** (`PROGRESS.md`, `REPORT.md`, `verify/wz16-e2e-gate/pass6/*`), so the green CI
run certifies the same buildable source that was actually exercised locally.

**Definitive per-tip check (use this for any later tip, it never needs re-derivation):**
`gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 1 --json headSha,conclusion` must equal
`git rev-parse HEAD`. Every commit added after this one is docs/evidence only, proven by re-running the
`git diff --name-only <code-tip> HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'` command above.

---

# REPORT — wz16 e2e integration gate (t_0509c320) — PASS 6, independent re-run

**Verdict: PASS — all 6 card steps were executed live this pass, on the pushed tip; the tree is clean and
committed; CI is green at the exact tip SHA.**

- Author: worker-heavy (fleet offload) · host **x280** (`c08r4d0r`) · run start `2026-10-16T17:14+02:00`
- Subject under test: **`e85720cd39a41d4af8bdff5172752a827e25bb86`** = `origin/ci/add-github-actions-workflow-clean`
  on fork `felixfelix-bot/net4sats-wizard-go` (upstream `net4sats/net4sats-wizard-go`)
- Worktree for this pass: `/home/c08r4d0r/worktrees/wz16-pass6` (detached at the tip)
- Toolchain: `go1.22.12 linux/amd64`, `golangci-lint 1.61.0` (== the CI pin), `goreleaser 2.18.3`,
  `gh` authenticated as `felixfelix-bot`
- Raw captures (verbatim): `verify/wz16-e2e-gate/pass6/*.txt` → mirrored at `/home/c08r4d0r/wz16-captures-v6/`
- Re-runnable harness: `verify/wz16-e2e-gate/pass6/wz16-0*.sh`

This is pass 6. Passes 1–5 exist on this branch; the pass-5 section below is preserved. Every command below was
re-executed live in this session — nothing was replayed from an earlier capture.

## 0. Why this pass ran in `~/worktrees/wz16-pass6` and not in the assigned workspace

The assigned workspace `/home/c08r4d0r/repos/net4sats-mvp-v2` was found **detached at `13de851`** — the tip of the
*old* `ci/github-actions-workflow-atomic` branch, which carries `ci.yml` **but no `.goreleaser.yaml` at all**:

```
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 rev-parse HEAD
13de851f8b56256524dade8f1fecad4edbf76292
$ git -C /home/c08r4d0r/repos/net4sats-mvp-v2 ls-tree -r --name-only HEAD | grep -i goreleaser
(empty)
```

The card says this gate "runs after the CI workflow, **goreleaser config**, and docs are committed", so steps 2 and 3
are undefined at that commit. That workspace also holds another card's untracked `REPORT.md`/`PROGRESS.md`, which
collide with the *tracked* copies at the deliverable tip. Verification was therefore run in a dedicated worktree
checked out at the deliverable tip, leaving the board workspace untouched:

```
$ git worktree add --detach /home/c08r4d0r/worktrees/wz16-pass6 origin/ci/add-github-actions-workflow-clean
HEAD is now at e85720c docs(wz16): add pass-5 card comment to the evidence bundle
$ git status --porcelain
PORCELAIN_LINES=0
$ ls .goreleaser.yaml .github/workflows/
.goreleaser.yaml    .github/workflows/: ci.yml  release.yml
```

## 1. Clean tree check — PASS (raw)

```
$ git status --porcelain
PORCELAIN_RC=0  LINES=0

$ git push          # literal, in the branch's own worktree /home/c08r4d0r/wz16-verify
Everything up-to-date
GIT_PUSH_RC=0

$ git push origin HEAD:refs/heads/ci/add-github-actions-workflow-clean   # same op, from this detached worktree
Everything up-to-date
PUSH_EXPLICIT_RC=0

$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0
local_HEAD            = e85720cd39a41d4af8bdff5172752a827e25bb86
origin/ci/add-github-actions-workflow-clean = e85720cd39a41d4af8bdff5172752a827e25bb86

$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
e85720cd39a41d4af8bdff5172752a827e25bb86	refs/heads/ci/add-github-actions-workflow-clean

### deliverable files present at this tip (CI workflow + goreleaser config + docs):
.github/workflows/ci.yml
.github/workflows/release.yml
.goreleaser.yaml
docs/endo-runbook.md
```

Porcelain empty, both push paths `rc=0`, ahead/behind `0	0`, and the remote's own reading of the ref
(`ls-remote`, independent of local tracking refs) equals local HEAD.

## 2. Full suite — PASS (raw)

```
$ go vet ./...
GO_VET_RC=0                       # no output = clean

$ golangci-lint run
GOLANGCI_RC=0                     # v1.61.0, no findings

$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
GORELEASER_CHECK_RC=0

$ go test ./... -race
ok  	net4sats-wizard	67.918s
GO_TEST_RACE_RC=0

$ go test ./... -race -v -count=1
ok  	net4sats-wizard	82.248s
GO_TEST_RACE_VERBOSE_RC=0
--- per-test tally ---
RUN_COUNT=123
PASS_COUNT=52
FAIL_COUNT=0
SKIP_COUNT=0
DATA_RACE_COUNT=0
--- PASS: TestVersionFlagInjectedViaLdflags (0.85s)
ok  	net4sats-wizard	82.248s
```

Two independent race-enabled runs, both green: 123 `=== RUN`, 52 `--- PASS`, **0 fail, 0 skip, 0 data race**.
The suite includes `TestVersionFlagInjectedViaLdflags`, which builds a binary with the same `-X main.version`
symbol `.goreleaser.yaml` uses and fails if that symbol drifts — the injection path used in step 4 is guarded by
the suite itself.

## 3. Snapshot release — PASS (raw; 6/6 artifacts, 3×2 matrix)

```
$ goreleaser release --snapshot --clean
   • snapshotting   version=0.5.1-snapshot.e85720c
   • building       6 targets: linux_amd64_v1 linux_arm64_v8.0 darwin_amd64_v1 darwin_arm64_v8.0 windows_amd64_v1 windows_arm64_v8.0
   • archiving      6 archives (4 .tar.gz + 2 .zip)
   • calculating checksums
   • release succeeded
GORELEASER_SNAPSHOT_RC=0

artifact list (name  size):
net4sats-wizard_0.5.1-snapshot.e85720c_darwin_amd64.tar.gz	4543771
net4sats-wizard_0.5.1-snapshot.e85720c_darwin_arm64.tar.gz	4411830
net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz	4491442
net4sats-wizard_0.5.1-snapshot.e85720c_linux_arm64.tar.gz	4267818
net4sats-wizard_0.5.1-snapshot.e85720c_windows_amd64.zip	4567009
net4sats-wizard_0.5.1-snapshot.e85720c_windows_arm64.zip	4309017

matrix breakdown: linux = 2   darwin = 2   windows = 2     (amd64 + arm64 each)
per-target built binaries:
net4sats-wizard_darwin_amd64_v1/net4sats-wizard	11395152
net4sats-wizard_darwin_arm64_v8.0/net4sats-wizard	11162146
net4sats-wizard_linux_amd64_v1/net4sats-wizard	11215000
net4sats-wizard_linux_arm64_v8.0/net4sats-wizard	10944664
net4sats-wizard_windows_amd64_v1/net4sats-wizard.exe	11427328
net4sats-wizard_windows_arm64_v8.0/net4sats-wizard.exe	11023872

$ wc -l dist/SHA256SUMS
6 dist/SHA256SUMS
ARCHIVE_COUNT=6  BINARY_COUNT=6  CHECKSUM_LINES=6
PARITY: OK  (6 checksum lines == 6 archives)

$ cd dist && sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.e85720c_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.e85720c_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.e85720c_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.e85720c_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.e85720c_windows_arm64.zip: OK
SHA256SUM_C_RC=0   OK_COUNT=6  FAILED_COUNT=0

-- required: one archive + one non-linux artifact --
A (archive,   linux/amd64)  = net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz
B (non-linux, darwin/arm64) = net4sats-wizard_0.5.1-snapshot.e85720c_darwin_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.e85720c_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz: OK
TWO_ARTIFACT_VERIFY_RC=0

-- independent recompute (not trusting the sums file), vs the published line --
71c6250fd8b7c7f799dc79d017e905b835232fb13fdd9504a4d2b0907e9d56c6  net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz
3f4182b93d9b2a55f487698abf7330501b5a3d4c76b63b024ad6b93d53c4ccfe  net4sats-wizard_0.5.1-snapshot.e85720c_darwin_arm64.tar.gz
published lines:  (identical to the recompute above — see capture 05-06)

-- negative control: a corrupted copy MUST fail, proving the check is non-vacuous --
net4sats-wizard_0.5.1-snapshot.e85720c_linux_amd64.tar.gz: FAILED
NEGATIVE_CONTROL_RC=1  (non-zero = control is real)
```

Checksum file name is `dist/SHA256SUMS` (pinned by `checksum.name_template` in `.goreleaser.yaml`, and the README
documents `sha256sum -c SHA256SUMS`, so the name is part of the public interface).

**Honest gap (not hidden):** Windows `Get-FileHash` was **not** executed — there is no Windows host on this machine.
The equivalent integrity proof (`sha256sum -c`, the same FIPS 180-4 SHA-256 over the same file bytes) passed for all
six artifacts, plus an independent recompute and a corruption control. Not claimed as a Windows run.

## 4. Binary smoke + version injection — PASS (raw)

```
$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, ... stripped

$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
STDOUT: tollgate-wizard 0.5.1-snapshot.e85720c
STDERR: (empty)
EXIT_CODE=0
LINE=[tollgate-wizard 0.5.1-snapshot.e85720c]
FORMAT:   OK   (prints 'tollgate-wizard <version>')
EXIT:     OK   (0)
INJECTION: OK   (non-dev version, not a fallback constant)

-- manual injection control: two different -X values, plus an uninjected build --
CTL_OUT=[tollgate-wizard v9.9.9-pass6-control]
CTL2_OUT=[tollgate-wizard v0.0.1-changed]
NOLD_OUT=[tollgate-wizard dev]   (no ldflags -> 'dev' is the fallback)
INJECTION_CHANGES_WITH_FLAG: OK
SNAPSHOT_VALUE_IS_STRIPPED_BINARY: OK
```

Required line `tollgate-wizard <version>` confirmed with the goreleaser-injected value `0.5.1-snapshot.e85720c`
(not `dev`), exit code 0. The injected value demonstrably changes with the ldflag (two different values), and the
no-ldflags control prints `dev` — so the snapshot value is really injected, not a coincidence.

## 5. GHA proof — PASS (raw)

The card's literal repo does not exist (measured again this pass, not cited):

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
RC=1
$ gh api repos/net4sats/tollgate-wizard-go --jq .message
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/repos#get-a-repository","status":"404"}
RC=1
```

`gh` is authenticated as `felixfelix-bot` and the real repo is public, so this is a genuine 404 — the card's
`--repo` is wrong, not a privacy/auth problem. Run against the real repo:

```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	worker-heavy/t_10a913ec	workflow_dispatch	38899734040	2m1s	2026-10-16T12:19:42Z
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38870340847	1m53s	2026-10-16T07:55:58Z
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38869481487	2m0s	2026-10-16T07:47:50Z
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38867038864	2m0s	2026-10-16T07:24:41Z
completed	success	CI	CI	worker-heavy/t_0509c320-r3	workflow_dispatch	38866690995	1m55s	2026-10-16T07:21:27Z
RC=0

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852530642	2m0s	2026-10-16T04:48:45Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852255804	2m2s	2026-10-16T04:45:27Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38852038477	1m56s	2026-10-16T04:43:00Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851802391	1m56s	2026-10-16T04:40:12Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851577790	1m59s	2026-10-16T04:37:28Z

$ gh run view 38852530642 --repo felixfelix-bot/net4sats-wizard-go --json databaseId,conclusion,headSha,headBranch,event,status,jobs,workflowName
{"id":38852530642,"wf":"CI","branch":"ci/add-github-actions-workflow-clean",
 "sha":"e85720cd39a41d4af8bdff5172752a827e25bb86","event":"workflow_dispatch","status":"completed",
 "conclusion":"success","jobs":[{"name":"lint","conclusion":"success"},{"name":"test","conclusion":"success"}]}

local HEAD          = e85720cd39a41d4af8bdff5172752a827e25bb86
CI run headSha      = e85720cd39a41d4af8bdff5172752a827e25bb86
CI run conclusion   = success
SHA_MATCH_CONCLUSION = OK  (the green run is exactly this pushed tip)
```

The load-bearing check is the SHA match, because "some green run exists on the branch" is a weaker claim than
"the pushed commit is green". Here the newest run on the branch is *at* the pushed tip, both jobs green.

**Cross-head identity (why the sibling branch's green CI covers the same code).** The other candidate head
`worker-heavy/t_0509c320-r3` @ `6345e2b` is also CI-green (`run 38870340847`, `conclusion=success`). The delta
between the two heads is exactly two test files:

```
$ git diff --name-only 6345e2b3f794494d5c66ab9e30050bd0b9ff5163 e85720cd39a41d4af8bdff5172752a827e25bb86 \
      -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'
handlers_test.go
security_hardening_test.go

blob identity for the build/CI-relevant paths (r3 vs this tip):
  IDENTICAL  .goreleaser.yaml
  IDENTICAL  .github/workflows/ci.yml
  IDENTICAL  .github/workflows/release.yml
  IDENTICAL  go.mod
  IDENTICAL  go.sum
  IDENTICAL  main.go
  IDENTICAL  deploy.go
  IDENTICAL  discover.go
  IDENTICAL  ssh.go
  IDENTICAL  embed.go
```

So the snapshot/`--version`/checksum evidence and both heads' CI runs describe the same buildable source; only the
test files differ, and both heads' suites pass.

**Structural caveat (recorded, not papered over).** These runs are `event=workflow_dispatch`, not automatic
`pull_request` runs. That is structural to the fork: a fork PR executes the **base** repo's workflows, and upstream
`net4sats/net4sats-wizard-go:main` has no `.github/` at all, so PR #29/#36 report an empty `statusCheckRollup`.
The dispatched run, however, executes the branch's *own* `.github/workflows/ci.yml` at the branch head — the same
file GitHub would load for the PR — so it is a faithful run of the real definition. **No workflow file was edited
and no check was weakened to make anything pass.** Making the check merge-blocking on the upstream PR requires
escalating the workflow to upstream `main` (manager decision, outside this gate).

## Gate table (this pass)

| # | Gate | Result | Capture |
|---|---|---|---|
| 1 | `git status --porcelain` empty | PASS (0 lines) | `pass6/01-02-steps1-2.txt` |
| 1 | pushed, ahead/behind 0 | PASS (`0	0`; ls-remote == HEAD; both `git push` rc 0) | `pass6/01-02-steps1-2.txt` |
| 2 | `go vet ./...` | PASS (rc 0) | `pass6/01-02-steps1-2.txt` |
| 2 | `go test ./... -race` | PASS (`ok 67.918s`, rc 0) | `pass6/03-test-race.txt` |
| 2 | `go test ./... -race -v` tally | PASS (123 RUN / 52 PASS / 0 FAIL / 0 SKIP / 0 DATA RACE) | `pass6/04-*` |
| 2 | `golangci-lint run` | PASS (rc 0, v1.61.0 == CI pin) | `pass6/01-02-steps1-2.txt` |
| 2 | `goreleaser check` | PASS (1 config validated) | `pass6/01-02-steps1-2.txt` |
| 3 | `goreleaser release --snapshot --clean` | PASS (6 targets / 6 archives, 2 per OS) | `pass6/05-06-steps3-4.txt` |
| 3 | checksum count == artifact count | PASS (6 == 6, parity OK) | `pass6/05-06-steps3-4.txt` |
| 3 | `sha256sum -c` ≥2 (one archive + one non-linux) | PASS (all 6 OK, rc 0; recompute matches; corruption control fails) | `pass6/05-06-steps3-4.txt` |
| 4 | linux/amd64 `--version` | PASS (`tollgate-wizard 0.5.1-snapshot.e85720c`, rc 0, not `dev`) | `pass6/05-06-steps3-4.txt` |
| 4 | injected value changes manually | PASS (`v9.9.9-pass6-control`, `v0.0.1-changed`, uninjected → `dev`) | `pass6/05-06-steps3-4.txt` |
| 5 | CI green via `gh run list` at the tip | PASS (run 38852530642, headSha == HEAD, lint+test success) | `pass6/07-step5-gha-proof.txt` |
| 6 | report + push + card comment | this report; commits on the branch; comment on `t_0509c320` | — |

**No step was skipped. Every command above was actually executed in this session and its output captured.**

## Findings / honest limits

1. **The card's repo slug does not exist** (`net4sats/tollgate-wizard-go` → HTTP 404, both `gh run list` and
   `gh api`). Real repo `net4sats/net4sats-wizard-go`, push target the fork `felixfelix-bot/net4sats-wizard-go`.
   Reported verbatim; nothing silently substituted.
2. **Windows `Get-FileHash` not run** (no Windows host). Substituted with the equivalent `sha256sum -c` over all
   six artifacts + independent recompute + a corruption control. Stated as a substitution.
3. **CI runs are `workflow_dispatch`, not `pull_request`** — structural to the fork/base relationship, explained
   above. The green run is at the pushed head SHA and its jobs are the card's `test` + `lint`.
4. **The card is in status `triage`**, which worker sessions cannot transition (`complete`/`promote`/`claim` all
   refuse on a triage card). This worker therefore attaches this report and comments; the **manager** must promote
   the card. The work itself is pushed and CI-green at the tip.
5. **Passes 1–5 already claimed PASS on this card.** This re-run re-executed everything from scratch rather than
   trusting them; the results agree. The pass-5 sections below are preserved unmodified.

**Honesty statement:** I did not create the CI workflow, goreleaser config, docs, or any earlier evidence commit —
those pre-date this run. What this run *did* create: the pass-6 scripts and captures, this report section, the
`PROGRESS.md` section, and their commits/pushes on `ci/add-github-actions-workflow-clean`. No source, test, or
workflow file was modified, no force-push was used, and no failing command was suppressed.

---

# REPORT — wz16 e2e integration gate (pass 5, independent re-run)

**Verdict: PASS — all 6 steps re-executed live; pushed; tree clean; CI green at the exact pushed tip.**

- Author: worker-heavy (fleet offload) · host **x280** · 2026-10-16T06:27:34+02:00
- Branch `ci/add-github-actions-workflow-clean` on fork `felixfelix-bot/net4sats-wizard-go`
- Verified code tip `64008dc2cfa83f9e3679f1f40d21f186c765642d`; pushed evidence tip **`c2af939bfe52068ab48d2eaefc550391f6cd0bb8`**
- Toolchain: go1.22.12, golangci-lint v1.61.0, goreleaser 2.x, gh authenticated as `felixfelix-bot`
- Raw captures: `verify/wz16-e2e-gate/pass5/` (12 files) — full report `verify/wz16-e2e-gate/REPORT.wz16-e2e-gate-pass5.md`

Summary of the 6 steps: (1) porcelain 0 lines, push up-to-date, ahead/behind `0	0`, ls-remote == HEAD;
(2) `go vet` rc 0, `golangci-lint` rc 0, `goreleaser check` validated, `go test ./... -race` 52 pass/0 fail/0 race;
(3) snapshot 6 archives (linux|windows|darwin × amd64|arm64), checksum parity 6==6, `sha256sum -c` 6/6 OK;
(4) `--version` → `tollgate-wizard 0.5.1-snapshot.64008dc` exit 0, ldflag control changes value, uninjected → `dev`;
(5) `gh run list` run `38844256545` headSha == HEAD, lint + test success; code tree identical across CI range.

Two card-string discrepancies recorded, not silently substituted: repo `net4sats/tollgate-wizard-go` is a 404
(real: `felixfelix-bot/net4sats-wizard-go`); the card's `--version` expectation `tollgate-wizard <version>`
is CORRECT (`main.go:55`). No source or workflow file was changed and no check was weakened.

---

# REPORT — End-to-end integration gate (t_0509c320, pass 4)

**Verdict: PASS — every gate re-executed live; branch pushed; CI green at the pushed head SHA.**

- Author: worker-heavy (fleet offload) · host **x280** · 2026-10-16T04:47–04:55Z
- Worktree root: `/home/c08r4d0r/wz16-verify` · branch `ci/add-github-actions-workflow-clean`
- **Subject under test: `9c37d2e04fe5bfafeb8d210202375ff7f34b773d`** = `origin/ci/add-github-actions-workflow-clean`
  = **PR #35 head** (OPEN) on fork `felixfelix-bot/net4sats-wizard-go`, base upstream `net4sats/net4sats-wizard-go:main`
- Toolchain: go1.22.12 (GOTOOLCHAIN=go1.22.12), gcc 13.3.0, golangci-lint **v1.61.0** (== CI pin),
  goreleaser 2.x, gh 2.45.0 authenticated as `felixfelix-bot`
- Raw captures (verbatim, 13 files): `/home/c08r4d0r/wz16-captures-t_0509c320-p4/`
- Re-runnable harness: `/home/c08r4d0r/wz16-tools/wz16-0*.sh`, `wz16-9*.py/.sh`

---

## Step 1 — clean tree check + push + ahead/behind

Required: `git status --porcelain` empty; `git push`; ahead/behind 0 vs remote.

```
$ git status --porcelain
PORCELAIN_RC=0   PORCELAIN_LINES=0   # 0 lines = clean tree

$ git push
Everything up-to-date
PUSH_RC=0

$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0
AB_RC=0   # format: <ahead>	<behind>

local_HEAD        = 9c37d2e04fe5bfafeb8d210202375ff7f34b773d
origin/ci/add-github-actions-workflow-clean = 9c37d2e04fe5bfafeb8d210202375ff7f34b773d

$ git ls-remote origin refs/heads/ci/add-github-actions-workflow-clean
9c37d2e04fe5bfafeb8d210202375ff7f34b773d	refs/heads/ci/add-github-actions-workflow-clean
LSREMOTE_RC=0
```

**PASS.** Clean tree, push a no-op, ahead/behind `0	0`, and remote ground truth (`ls-remote`) agrees with
local HEAD. (`git rev-list --left-right --count origin/main...HEAD` → `16	11` is the *base-vs-tip* distance,
i.e. 11 commits ahead of `main` and the base 16 behind — expected for an open PR branch, not a push gap.)

## Step 2 — full suite

### 2a `go vet ./...` (+ gofmt)

```
$ go list ./...
net4sats-wizard
GOLIST_RC=0

$ go vet ./...
GO_VET_RC=0

$ gofmt -l .
GOFMT_RC=0   # no output above = formatted
```

Single module (`net4sats-wizard`), so the mandated `./...` gates cover the whole test surface.

### 2b `go test ./... -race`

```
$ CGO_ENABLED=1 go test ./... -race -count=1
ok  	net4sats-wizard	108.034s

real	1m49.242s
TEST_RACE_RC=0
```

Verbose re-run counters:

```
RUN   = 123
PASS  = 52
FAIL  = 0
SKIP  = 0
DATA_RACE_WARNINGS = 0
final ok line:
--- PASS: TestVersionStringFormat (0.00s)
--- PASS: TestVersionFlagInjectedViaLdflags (0.43s)
ok  	net4sats-wizard	90.007s
```

### 2c `golangci-lint run`

```
$ golangci-lint run   # golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09T17:44:42Z
GOLANGCI_RC=0
```

### 2d `goreleaser check`

```
$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
GORELEASER_CHECK_RC=0
```

**PASS** on all four.

## Step 3 — snapshot release, artifacts, checksums

```
$ goreleaser release --snapshot --clean
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 3s
SNAPSHOT_RC=0
```

Artifacts — **6**, exactly the required target matrix:

```
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_arm64.tar.gz
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_amd64.zip
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip
ARCHIVE_COUNT=6
distinct_targets=6
```

`file(1)` confirms real formats (4 × gzip tar, 2 × Zip deflate).

Checksum file + count:

```
$ ls -la SHA256SUMS
-rw-r--r-- 1 c08r4d0r c08r4d0r 744 Oct 16 04:52 SHA256SUMS
$ wc -l < SHA256SUMS
6
6  ==  6  (checksum_count == artifact_count)
```

`sha256sum -c` (real exit code, no pipe):

```
$ sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip: OK
SHA256SUM_C_REAL_RC=0
OK_LINES=6  FAILED_LINES=0
```

Two mandated spot checks, independently recomputed (one archive, one non-linux):

```
artifact=net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz type=tar.gz recomputed=ce227c9c80e726ff2277f6bb8413a72639666d558d1030093a28122241ac4419 expected=ce227c9c80e726ff2277f6bb8413a72639666d558d1030093a28122241ac4419 MATCH
artifact=net4sats-wizard_0.5.1-snapshot.9c37d2e_darwin_arm64.tar.gz type=tar.gz recomputed=f2da1a6e81fc70283de80f182fb95bb923a46fc8fa41d00bfaf54eb3a33fd188 expected=f2da1a6e81fc70283de80f182fb95bb923a46fc8fa41d00bfaf54eb3a33fd188 MATCH
artifact=net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_amd64.zip type=zip recomputed=97544c41ef17f89d9e5c9dd4860e22ef66db8aa60c58a9829a02295dd318f02e expected=97544c41ef17f89d9e5c9dd4860e22ef66db8aa60c58a9829a02295dd318f02e MATCH
POWERSHELL_PRESENT=no
```

`PowerShell` is not installed on this host, so the Windows `Get-FileHash` equivalent is `sha256sum` on the
zip — `windows_amd64.zip` → `MATCH`.

### 3b Non-vacuity controls (prove `sha256sum -c` can fail)

```
--- CONTROL A (positive): pristine copies -- expect 6/6 OK, REAL_RC=0 ---
REAL_RC=0
OK=6  FAILED=0

### tamper: flip one byte of the EXPLICIT path net4sats-wizard_..._linux_amd64.tar.gz
tamper_cmd_rc=0   size_after=4492301
--- CONTROL B (negative): 1 byte flipped -- expect that line FAILED, REAL_RC=1 ---
net4sats-wizard_0.5.1-snapshot.9c37d2e_linux_amd64.tar.gz: FAILED
sha256sum: WARNING: 1 computed checksum did NOT match
REAL_RC=1
OK=5  FAILED=1

--- CONTROL C (negative): windows_arm64.zip missing -- expect FAILED open or read, REAL_RC=1 ---
sha256sum: net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip: No such file or directory
net4sats-wizard_0.5.1-snapshot.9c37d2e_windows_arm64.zip: FAILED open or read
REAL_RC=1
OK=5  FAILED=1
```

(Controls use **explicit** filenames — a `dd of=<glob>` never expands and would have produced a vacuous
control — and read the **real** `sha256sum` rc via a temp file, never through a pipe.)

**PASS.**

## Step 4 — binary smoke (`--version`) + injection change

Binary extracted from the **built snapshot archive** (not a rebuild):

```
$ file ./net4sats-wizard
./net4sats-wizard: ELF 64-bit LSB executable, x86-64, ..., statically linked, ..., stripped

$ ./net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.9c37d2e
exit_code=0
LINE=[tollgate-wizard 0.5.1-snapshot.9c37d2e]
FOUND_DEV=no (good)   shape_ok=yes
injected_version_token=0.5.1-snapshot.9c37d2e

$ ./net4sats-wizard -version
tollgate-wizard 0.5.1-snapshot.9c37d2e
dash_version_rc=0
```

Prints `tollgate-wizard <version>` with the **injected** version — not `dev` — exit code 0.

Manual injection, proving the value demonstrably changes:

```
$ go build -ldflags "-X main.version=v9.9.9-manual-p4" -o .../wiz-manual .
$ .../wiz-manual --version
tollgate-wizard v9.9.9-manual-p4
run_rc=0

$ go build -ldflags "-X main.version=v0.0.0-second-different-value-p4" -o .../wiz-manual .
$ .../wiz-manual --version
tollgate-wizard v0.0.0-second-different-value-p4
run_rc=0

# control: plain build, NO ldflags -> dev default
$ go build -o .../wiz-nodev .
$ .../wiz-nodev --version
tollgate-wizard dev
run_rc=0

$ go tool nm .../wiz-manual | grep -F ' main.version'
  eb5c20 D main.version
  ca6b20 R main.version.str

# source anchors
52:var version = "dev"
59:func versionString() string {
641:		case "--version", "-version":
```

**PASS.**

## Step 5 — GHA proof

### The literal card command, verbatim (the card names `net4sats/tollgate-wizard-go`)

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found (https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
literal_real_rc=1

$ gh repo view net4sats/tollgate-wizard-go
GraphQL: Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'. (repository)
repoview_real_rc=1

$ git ls-remote https://github.com/net4sats/tollgate-wizard-go
remote: Repository not found.
fatal: repository 'https://github.com/net4sats/tollgate-wizard-go/' not found
lsremote_real_rc=128
```

`gh` is authenticated as `felixfelix-bot` over HTTPS, so this is a **real 404** — that repo does not exist.
Per the card, I report the exact error instead of asserting success. The real repo is
**`net4sats/net4sats-wizard-go`**.

### Upstream cannot report CI on this PR yet (the CI venue is the fork)

```
$ gh api repos/net4sats/net4sats-wizard-go/contents/.github --jq .message
{"message":"Not Found", ... "status":404}gh: Not Found (HTTP 404)
upstream_ghdir_rc=1

$ gh pr checks 35 --repo net4sats/net4sats-wizard-go
no checks reported on the 'ci/add-github-actions-workflow-clean' branch
prchecks_rc=1

$ gh pr view 35 --repo net4sats/net4sats-wizard-go --json number,state,headRefOid,statusCheckRollup
{"headRefOid":"9c37d2e04fe5bfafeb8d210202375ff7f34b773d","number":35,"state":"OPEN","statusCheckRollup":[]}
```

Upstream `main` has no `.github/` — that is precisely what PR #35 adds — so no check can exist upstream
until the workflow lands. **No check was weakened and no workflow was edited.**

### CI green at the pushed head SHA (fork venue)

```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838430526	2m1s	2026-10-16T02:03:24Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838000347	2m1s	2026-10-16T01:58:39Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38834487839	1m57s	2026-10-16T01:19:29Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
fork_runs_rc=0

# a run at THE PUSHED HEAD SHA with conclusion success
successful_runs_at_9c37d2e04fe5bfafeb8d210202375ff7f34b773d = 1
38838430526 completed success 9c37d2e04fe5bfafeb8d210202375ff7f34b773d https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38838430526
```

Run detail + the gate steps that actually executed:

```
$ gh run view 38838430526 --json conclusion,headSha,headBranch,status,workflowName,event,createdAt,url
{"conclusion":"success","createdAt":"2026-10-16T02:03:24Z","event":"workflow_dispatch","headBranch":"ci/add-github-actions-workflow-clean","headSha":"9c37d2e04fe5bfafeb8d210202375ff7f34b773d","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38838430526","workflowName":"CI"}

✓ ci/add-github-actions-workflow-clean CI #29 · 38838430526
✓ lint in 20s (ID 116772683647)
✓ test in 1m57s (ID 116772683899)

test	Run go build ./...	    2026-10-16T02:03:41Z go build ./...
test	Run go vet ./...	    2026-10-16T02:03:42Z go vet ./...
test	Run go test ./... -race -count=1	2026-10-16T02:03:43Z go test ./... -race -count=1
test	Run go test ./... -race -count=1	2026-10-16T02:05:21Z ok  	net4sats-wizard	82.366s
lint	Run golangci/golangci-lint-action@v6	version: v1.61.0
lint	Run golangci/golangci-lint-action@v6	2026-10-16T02:03:43Z golangci-lint found no issues
```

**PASS** — the CI workflow ran and concluded `success` at the exact pushed SHA, with the same three gates
this pass ran locally (`go build` / `go vet` / `go test ./... -race -count=1`) plus the `golangci-lint` job.

## Self-verification (independently re-derived from the captures)

```
$ python3 wz16-90-adhoc-verify.py <captures>
TOTAL CHECKS = 62   FAILURES = 0
ALL CHECKS PASS

$ wz16-91-adhoc-control.sh   # 4 facts corrupted in a copy
TOTAL CHECKS = 62   FAILURES = 4
FAILED CHECKS:
  - 02b: verbose 0 FAIL
  - 03: sha256sum -c real rc 0
  - 04: --version output line exactly matches the injected version
  - 05: gh run view confirms conclusion success at that SHA
CONTROL_VERIFIER_EXIT=1
```

62/62 real checks pass; corrupting four load-bearing facts fails exactly those four — the checks are not vacuous.

Secret scan: clean (`SCAN_RC=1` = no match) over the captures and over `PROGRESS.md`/`REPORT.md`.

## Honest caveats

1. **The card's repo name is wrong** (F1). Every gate was run against the real repo/fork; the literal
   command's 404 is captured verbatim rather than papered over.
2. **`gh pr checks 35` on upstream is empty** because upstream has no workflow yet — a structural fact, not
   a red CI. Green CI evidence is from the fork run at the identical head SHA.
3. This pass **changed no source, test, workflow, or goreleaser file**; the only commits are documentation
   (`PROGRESS.md`, `REPORT.md`) plus the re-runnable harness. The reviewed branch content is unchanged
   (`git diff 0b5322a..9c37d2e` = documentation only).
4. `PROGRESS.md`/`REPORT.md` were rewritten from the 27-line doc-only versions that were already on the
   docs branch head — i.e. the previous pass's dispatch notes are replaced by this pass's full report, not
   appended to. The tip therefore advances with this commit; CI was re-run at the new tip and is green
   (below).

## Addendum — CI re-verified at the post-report tip

The pass's own documentation commit moved the tip from `9c37d2e` to `649d90f`. CI was re-dispatched and
re-verified at that new tip:

```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38842920848	2m1s	2026-10-16T02:55:37Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838430526	2m1s	2026-10-16T02:03:24Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838000347	2m1s	2026-10-16T01:58:39Z

$ gh run view 38842920848 --json conclusion,headSha,status,event,url
{"conclusion":"success","event":"workflow_dispatch","headSha":"649d90ff94bea71f7f98a9767374df5e9a176223","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38842920848"}

$ gh run view 38842920848 --json jobs --jq '[.jobs[] | "\(.name) \(.status) \(.conclusion)"]'
["lint completed success","test completed success"]
```

So CI is green at **both** the reviewed code tip `9c37d2e` (run 38838430526) and the documentation tip
`649d90f` (run 38842920848). The addendum commit moved the tip once more to `5efd948`; CI was dispatched
again and is green there too:

```
$ gh run view 38843145298 --json conclusion,headSha,status,url
{"conclusion":"success","headSha":"5efd9483cc11286639b132f3e0264364902097b7","status":"completed","url":"https://github.com/felixfelix-bot/net4sats-wizard-go/actions/runs/38843145298"}

$ gh run list --repo felixfelix-bot/net4sats-wizard-go --branch ci/add-github-actions-workflow-clean --limit 4
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38843145298	2m0s	2026-10-16T02:58:21Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38842920848	2m1s	2026-10-16T02:55:37Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838430526	2m1s	2026-10-16T02:03:24Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838000347	2m1s	2026-10-16T01:58:39Z
```

Every documentation commit is treated the same way: pushed, CI re-dispatched, its own run id recorded. The
reviewed *code* tree is unchanged across all of these — `git diff 9c37d2e <tip> -- ':(exclude)*.md' -- ':(exclude)verify/'` is empty.

Final state: `git status --porcelain` → 0 lines; `git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean` → `0	0`.



---

## Pass-5 close-out: CI green at the final pushed tip

Per-tip confirmation loop (each line = a `workflow_dispatch` CI run whose `headSha` equals that tip):

| tip | CI run | conclusion |
|---|---|---|
| `64008dc2` (pass-4 tip, code tree) | 38844256545 | success |
| `95743c3` (pass-5 evidence) | 38851358337 | success |
| `21ba636` | 38851577790 | success |
| `c71b2f8` (**FINAL tip = HEAD = remote**) | 38851802391 | success |

Final state, verbatim:

```
### FINAL tip = HEAD = remote = c71b2f860bb51c951c4cd243586ec6cda8dcc947
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 3
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851802391	1m56s	2026-10-16T04:40:12Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851577790	1m59s	2026-10-16T04:37:28Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38851358337	1m58s	2026-10-16T04:34:52Z

$ gh run view 38851802391 --json conclusion,headSha,status
{"conclusion":"success","headBranch":"ci/add-github-actions-workflow-clean","headSha":"c71b2f860bb51c951c4cd243586ec6cda8dcc947","status":"completed"}

porcelain=[]
ahead/behind=0	0
```

Capture on disk: `verify/wz16-e2e-gate/pass5/11-ci-at-c71b2f8.txt`.
**All six card steps pass on a pushed, clean, CI-green tip.**

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
