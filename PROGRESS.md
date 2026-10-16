# PROGRESS — E2E integration gate: snapshot build, --version smoke, full suite, push

Repo (actual): net4sats-wizard-go  — clone /home/c08r4d0r/wz16-verify
  origin = https://github.com/felixfelix-bot/net4sats-wizard-go.git  (fork)
  upstream = net4sats/net4sats-wizard-go (public org repo)

CARD REPO-NAME MISMATCH (finding, not a blocker):
  Card says repo "net4sats/tollgate-wizard-go", workspace
  /home/c08r4d0r/repos/net4sats-mvp-v2. Neither exists:
    - `gh repo view net4sats/tollgate-wizard-go` -> GraphQL: Could not resolve.
    - `gh api orgs/net4sats/repos` lists 10 repos; NO tollgate-wizard-go.
      The wizard repo is `net4sats-wizard-go`.
    - /home/c08r4d0r/repos/net4sats-mvp-v2 is the net4sats DEMO repo, not Go.
  Work genuinely lives in /home/c08r4d0r/wz16-verify (felixfelix-bot fork).
  CLI identity is `tollgate-wizard` (const programName, main.go:55) even though
  the repo/module is net4sats-wizard — so the card's expected
  "tollgate-wizard <version>" output is CORRECT and matches source.

Branch under test: ci/add-github-actions-workflow-clean
  HEAD = 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff
  "feat(release): link-time version injection, --version, goreleaser config and release docs"
  This is the commit the card says to verify AFTER (CI workflow + goreleaser + docs).

Toolchain (not on default PATH): /home/c08r4d0r/wz16-tools/
  go122/bin/go = go1.22.12 ; goreleaser = 2.18.3 ; golangci-lint = 1.61.0

STATUS: IN PROGRESS — steps 1-4 complete and green; 5-6 pending.

## Findings -> status -> files touched

1. FINDING: card repo name + workspace are both wrong (see above).
   STATUS: routed to correct clone/ref. FILES: none (triage only).

2. FINDING: `go`/`goreleaser`/`golangci-lint` absent from default PATH.
   STATUS: fixed by exporting /home/c08r4d0r/wz16-tools paths. FILES: none.

3. STEP 1 (clean tree + pushed): `git status --porcelain` -> only 2 UNTRACKED
   lane docs (PROGRESS.md, REPORT.md); tracked tree clean. `git push` ->
   "Everything up-to-date"; ahead/behind 0/0; ls-remote tip == local HEAD
   == 0b5322a070ca00538d0d4d9be3fc3b2df41fc9ff.
   STATUS: PASS. Raw: wz20-captures/01-clean-tree.txt, 01-push.txt.

4. STEP 2 (full suite): go mod verify "all modules verified"; go build ./... RC=0;
   `go vet ./...` RC=0 clean; `go test ./... -race -count=1` -> ok 165.236s;
   verbose -> 123 tests run, 52 PASS, 0 FAIL, 0 SKIP, 0 DATA RACE;
   `golangci-lint run` RC=0 (v1.61.0); `goreleaser check` 1 config validated RC=0.
   STATUS: PASS. Raw: 02a-vet.txt, 02b-test-race.txt, 02c-lint-grcheck.txt,
   02d-test-verbose.log, 02e-test-counts.txt.

5. STEP 3 (snapshot): `goreleaser release --snapshot --clean` RC=0, version
   0.5.1-snapshot.0b5322a. EXACTLY 6 archives: linux/windows/darwin x amd64/arm64
   (4 tar.gz + 2 zip) + SHA256SUMS with 6 lines == 6 artifacts.
   `sha256sum -c SHA256SUMS` -> all 6 OK, RC=0. Hand-verified MATCH for
   linux_amd64.tar.gz (archive) AND darwin_arm64.tar.gz (non-linux); openssl
   cross-check agrees. Counts 6==6.
   STATUS: PASS. Raw: 03-snapshot-release.txt, 03b-checksums.txt.

6. STEP 4 (--version smoke): raw goreleaser linux/amd64 binary prints
   `tollgate-wizard 0.5.1-snapshot.0b5322a`, RC=0 (proves the card's exact
   expected shape, version injected not dev). The SHIPPED archive's binary
   prints the same -> the artifact itself is correct. Positive control
   (-X main.version=v9.9.9-control) -> `tollgate-wizard v9.9.9-control`;
   negative control (no ldflags) -> `tollgate-wizard dev`. Binary is ELF
   x86-64, statically linked, stripped.
   STATUS: PASS. Raw: 04-version-smoke.txt.

## Capture directory
/home/c08r4d0r/wz20-captures/   (all raw command output for this pass)
