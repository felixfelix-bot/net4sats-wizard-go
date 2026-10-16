# REPORT — End-to-end verify: snapshot build, `--version` smoke, full suite, push

Task: t_0509c320 "End-to-end verify: snapshot build, --version smoke, full suite, push" (re-dispatch)
Board: net4sats-mvp-v2 · lane repo (this worktree): /home/c08r4d0r/wz16-verify
Repo under test: **felixfelix-bot/net4sats-wizard-go** (fork; upstream net4sats/net4sats-wizard-go)
Branch: `ci/add-github-actions-workflow-clean` (head of OPEN PR #35)
**Verified tip: `4ca0b381ab3aa40f00b74f133f7b9024740277d7`**
Environment: git 2.43.0, go1.22.12 linux/amd64, golangci-lint 1.61.0, goreleaser 2.18.3
Raw captures: `/home/c08r4d0r/wz16-redispatch/captures/` (one file per step, verbatim)
Run window: 2026-10-16T01:54Z → 02:05Z

VERDICT: **all 5 steps executed, none skipped; tree clean and pushed; CI `success` at the shipped
tip.** Two things this report deliberately does NOT do: it does not relay the previous run's
green, and it does not claim a step it did not re-execute.

---

## Step 0 — task's repo literal

The task names `net4sats/tollgate-wizard-go`. That repository does not exist:

```
$ gh run list --repo net4sats/tollgate-wizard-go --limit 5
failed to get runs: HTTP 404: Not Found
(https://api.github.com/repos/net4sats/tollgate-wizard-go/actions/runs?per_page=5&exclude_pull_requests=true)
rc=1
```

The wizard repo is `net4sats/net4sats-wizard-go`. Upstream is read-only for this account, so the
branch and its CI live on the fork `felixfelix-bot/net4sats-wizard-go` — the same arrangement every
prior wizard PR used, and this worktree's `origin`.

---

## Step 1 — clean tree check + push

capture: `captures/01-clean-tree-and-push.txt`

```
$ git status --porcelain
(empty output)
rc=0

$ git fetch origin
(no output)
$ git push origin ci/add-github-actions-workflow-clean
Everything up-to-date
push rc=0

$ git rev-list --left-right --count HEAD...origin/ci/add-github-actions-workflow-clean
0	0
(ahead<TAB>behind) — want 0  0

$ git rev-parse HEAD
4ca0b381ab3aa40f00b74f133f7b9024740277d7

$ git branch -vv | grep '^\*'
* ci/add-github-actions-workflow-clean 4ca0b38 [origin/ci/add-github-actions-workflow-clean] docs(ci-verify): mark gate COMPLETE — CI re-confirmed green at the shipped tip 11c661b
```

PASS: tree clean, remote up to date, ahead/behind 0/0.

## Step 2 — full suite

captures: `captures/02-vet-and-test-race.txt`, `captures/02c-lint-and-goreleaser-check.txt`

```
$ go vet ./...
(no output)
go vet rc=0

$ go test ./... -race
ok  	net4sats-wizard	121.327s
go test rc=0

$ golangci-lint run
(no output)
golangci-lint rc=0
golangci-lint found no issues

$ goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
goreleaser check rc=0

$ gofmt -l .
(empty)
gofmt rc=0
```

PASS: vet clean, race suite green, lint 0 issues, goreleaser config valid, formatted.

## Step 3 — snapshot release: 6 artifacts + checksum verification

captures: `captures/03a-snapshot-release.txt`, `captures/03b-inventory-and-checksums.txt`,
`captures/03e-checksum-verification.txt`

```
$ goreleaser release --snapshot --clean
  • skipping announce, publish, and validate...
  • cleaning distribution directory
  • getting and validating git state
    • using tags                                     previous=v0.3.5 current=v0.5.0
  • snapshotting
    • building snapshot...                           version=0.5.1-snapshot.4ca0b38
  • running before hooks
    • running                                        hook=go mod verify
  • building binaries
    • building   binaries=net4sats-wizard target=darwin_arm64_v8.0
    • building   binaries=net4sats-wizard target=darwin_amd64_v1
    • building   binaries=net4sats-wizard target=windows_arm64_v8.0
    • building   binaries=net4sats-wizard target=windows_amd64_v1
    • building   binaries=net4sats-wizard target=linux_arm64_v8.0
    • building   binaries=net4sats-wizard target=linux_amd64_v1
  • archives
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_arm64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_amd64.zip
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_arm64.zip
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_amd64.tar.gz
    • archiving  name=dist/net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 12s
goreleaser rc=0
```

**Artifact list** (`dist/`):

| # | file | size (bytes) |
|---|------|--------------|
| 1 | net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz   | 4492303 |
| 2 | net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_arm64.tar.gz   | 4268921 |
| 3 | net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_amd64.tar.gz  | 4544366 |
| 4 | net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz  | 4413920 |
| 5 | net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_amd64.zip    | 4567564 |
| 6 | net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_arm64.zip    | 4310151 |

```
archives: 6
binaries: 6
checksum files: dist/SHA256SUMS

executable binaries built:
dist/net4sats-wizard_darwin_amd64_v1/net4sats-wizard
dist/net4sats-wizard_darwin_arm64_v8.0/net4sats-wizard
dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
dist/net4sats-wizard_linux_arm64_v8.0/net4sats-wizard
dist/net4sats-wizard_windows_amd64_v1/net4sats-wizard.exe
dist/net4sats-wizard_windows_arm64_v8.0/net4sats-wizard.exe
```

**Checksum count == artifact count**: `dist/SHA256SUMS` has **6 lines** for **6 archives**.

```
$ cat dist/SHA256SUMS
9b2aefbc046af4cfefb5f4f3a64731d16a3e114cb7dd43c1dfd56f10a4f96f07  net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_amd64.tar.gz
864f63d75e48fa129159dfadc4c70609bd63dfaf0b4c5e7cb2448c988163bfa4  net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz
119b50f7bcbeac275219ddfd3155bb20caf954f3ce43b24809095ac223d0f5df  net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz
6f847aeb9e94915b8ded054de38ef0e4e438aa543c9b672d7fd3714b6e962ba8  net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_arm64.tar.gz
742e7c9694bd0a7f2ea001048c126aa94b3bf18521d62992856ad5218e140c80  net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_amd64.zip
eb3d2f22c298e6f97307ff076887924aa6855cfc28b80fadc764623b6b792a66  net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_arm64.zip
```

**Full verification (`sha256sum -c`) and strict mode:**

```
$ cd dist && sha256sum -c SHA256SUMS
net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_arm64.tar.gz: OK
net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_amd64.zip: OK
net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_arm64.zip: OK
sha256sum -c rc=0

$ sha256sum -c SHA256SUMS --ignore-missing --strict
(all six : OK)
rc=0
```

**Two-artifact spot check — one archive, one NON-linux:**
```
$ grep linux_amd64.tar.gz SHA256SUMS | sha256sum -c -
net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz: OK

$ grep darwin_arm64.tar.gz SHA256SUMS | sha256sum -c -      # NON-linux
net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz: OK

$ grep windows_amd64.zip SHA256SUMS | sha256sum -c -        # NON-linux
net4sats-wizard_0.5.1-snapshot.4ca0b38_windows_amd64.zip: OK
```

**Independent recomputation vs the recorded hash:**
```
$ sha256sum net4sats-wizard_0.5.1-snapshot.4ca0b38_linux_amd64.tar.gz
119b50f7bcbeac275219ddfd3155bb20caf954f3ce43b24809095ac223d0f5df  ...linux_amd64.tar.gz
   (== line 3 of SHA256SUMS)
$ sha256sum net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_amd64.tar.gz
9b2aefbc046af4cfefb5f4f3a64731d16a3e114cb7dd43c1dfd56f10a4f96f07  ...darwin_amd64.tar.gz
   (== line 1 of SHA256SUMS)
```

**Negative control** (the check must be able to fail — appended one byte to a copy):
```
$ printf 'x' >> /tmp/tamper.tar.gz
$ sha256sum /tmp/tamper.tar.gz | awk '{print $1}'
2a0d54b44bdd0ecf8fd096149691e1a1a8fe70c8b73de552cc8e9cc2a66feffa
$ grep 'darwin_arm64.tar.gz' SHA256SUMS | awk '{print $1}'
864f63d75e48fa129159dfadc4c70609bd63dfaf0b4c5e7cb2448c988163bfa4
$ sha256sum -c /tmp/tamper.sum
net4sats-wizard_0.5.1-snapshot.4ca0b38_darwin_arm64.tar.gz: FAILED
sha256sum: WARNING: 1 computed checksum did NOT match
negative-control rc=1
```

PASS: 6/6 artifacts checksummed, count matches, two named artifacts verified individually
(one linux archive + two non-linux), and a tampered copy is correctly rejected.

*(PowerShell `Get-FileHash` equivalent is documented in the repo README, but this host is Linux —
the achieved check is `sha256sum -c`, which is the same SHA-256 verification. Stated plainly rather
than asserting an untested Windows command.)*

## Step 4 — binary smoke: `--version` on the built linux/amd64 binary

captures: `captures/04a-binary-smoke.txt`, `captures/04d-ldflags-injection.txt`

```
$ file dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard
... ELF 64-bit LSB executable, x86-64, statically linked, Go BuildID=D142HRxzy..., stripped

$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard --version
tollgate-wizard 0.5.1-snapshot.4ca0b38
exit code = 0

assertions:
ASSERT format OK (prefix 'tollgate-wizard ')
ASSERT OK: version is NOT dev
ASSERT OK: exit code 0

$ dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard -version
tollgate-wizard 0.5.1-snapshot.4ca0b38
rc=0

5-run determinism:
run 1: tollgate-wizard 0.5.1-snapshot.4ca0b38
run 2: tollgate-wizard 0.5.1-snapshot.4ca0b38
run 3: tollgate-wizard 0.5.1-snapshot.4ca0b38
run 4: tollgate-wizard 0.5.1-snapshot.4ca0b38
run 5: tollgate-wizard 0.5.1-snapshot.4ca0b38
```

**Injected value changes when built manually** (required by the task):
```
$ go build -ldflags "-s -w -X main.version=v9.9.9-manual-test" -o /tmp/wiz-injected ./
build rc=0
$ /tmp/wiz-injected --version
tollgate-wizard v9.9.9-manual-test
exit code = 0
ASSERT OK: injected value appears in output

$ go build -o /tmp/wiz-nodev ./
build rc=0
$ /tmp/wiz-nodev --version
tollgate-wizard dev
ASSERT OK: un-injected build prints 'dev' (proves ldflags is what changes it)
```

PASS: `tollgate-wizard <version>` with the injected release version, rc 0, not `dev`; injection
demonstrably changes the value.

## Step 5 — GHA proof

captures: `captures/05a-gha-list.txt`, `05e-gha-dispatch.txt`, `05f-gha-watch.txt`, `05g-gha-final.txt`

On arrival, **to the current tip was not covered**:
```
$ git rev-parse HEAD
4ca0b381ab3aa40f00b74f133f7b9024740277d7
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 20 --json headSha,conclusion,...
runs matching tip: 0
```

The newest green run was one commit behind:
```
databaseId 38834487839  conclusion success  headSha 11c661b253ea29cfc051ae801d6bd99238c5201f
```

So CI was re-dispatched on the exact tip and watched to completion:
```
$ gh workflow list --repo felixfelix-bot/net4sats-wizard-go
CI	active	376283443
$ gh workflow run .github/workflows/ci.yml --repo felixfelix-bot/net4sats-wizard-go --ref ci/add-github-actions-workflow-clean
dispatch rc=0
$ gh run watch 38838000347 --repo felixfelix-bot/net4sats-wizard-go --exit-status
✓ lint in 22s                        (checkout, setup-go, golangci-lint-action@v6 — all ✓)
✓ test in 1m56s                      (go build ✓, go vet ✓, go test ./... -race -count=1 ✓)
watch rc=0

$ gh run view 38838000347 --json headSha,status,conclusion,jobs
{"conclusion":"success","headSha":"4ca0b381ab3aa40f00b74f133f7b9024740277d7",
 "status":"completed","jobs":[{"name":"lint","conclusion":"success"},{"name":"test","conclusion":"success"}]}
```

Final `gh run list`:
```
$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38838000347	2m1s	2026-10-16T01:58:39Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38834487839	1m57s	2026-10-16T01:19:29Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38761913548	1m54s	2026-10-15T15:37:06Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38738030458	2m3s	2026-10-15T12:52:40Z
completed	success	CI	CI	ci/add-github-actions-workflow-clean	workflow_dispatch	38677575488	1m50s	2026-10-15T03:41:31Z
```

Assertion over the run list:
```
runs at tip 4ca0b381ab3a : 1
AT LEAST ONE SUCCESS AT TIP: True
ASSERT OK
```

CI was **not red** at any point, so the task's "fix the workflow, do not weaken the checks"
instruction never triggered; no workflow or linter configuration was touched.

Honest caveat: `gh pr checks 35 --repo net4sats/net4sats-wizard-go` still answers
"no checks reported on the 'ci/add-github-actions-workflow-clean' branch" — GitHub Actions do not
run on PRs in this fork→upstream arrangement. The green evidence is the branch-head run, and the
branch head IS the PR head (`headRefOid` 4ca0b38…, PR #35 OPEN). This is a limitation of the
fork arrangement, not a missing verification.

## Step 6 — report / comment

Summary comment posted on PR #35 (raw output above, verbatim, plus the artifact and checksum
captures):
https://github.com/net4sats/net4sats-wizard-go/pull/35#issuecomment-6206347289

---

## Acceptance checklist

| # | requirement | state |
|---|-------------|-------|
| 1 | clean tree (`git status --porcelain` empty) + pushed, 0/0 | PASS |
| 2 | `go vet`, `go test -race`, `golangci-lint run`, `goreleaser check` | PASS (all rc=0) |
| 3 | snapshot: 6 artifacts + archives + checksum file; counts match; ≥2 verified (archive + non-linux) | PASS |
| 4 | linux/amd64 `--version` prints `tollgate-wizard <version>`, not `dev`, rc 0; injection changes it | PASS |
| 5 | `gh run list` shows the CI run for this branch with conclusion `success` at the shipped tip | PASS (run 38838000347 @ 4ca0b38) |
| 6 | report with raw outputs, artifact list, checksum output, `--version` output, `gh run list` line | this document + PR comment |

Nothing skipped. No source, workflow, or release config was modified: every step was already green
locally; the only action needed was re-running CI on the commit that had never been tested.

## Files touched by this re-dispatch
- `PROGRESS.md` (worktree root) — this run's findings
- `REPORT.md` (worktree root) — this document
- No repo source/config changes. Evidence lives outside the repo in
  `/home/c08r4d0r/wz16-redispatch/captures/`.
