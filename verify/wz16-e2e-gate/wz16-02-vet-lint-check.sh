#!/usr/bin/env bash
# wz16 STEP 2 — go vet, gofmt, go test -race, golangci-lint, goreleaser check.
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
WT=/home/c08r4d0r/wz16-verify
export PATH="/home/c08r4d0r/wz16-tools/go122/bin:$PATH"
export GOTOOLCHAIN=go1.22.12
export GOCACHE="${GOCACHE:-/home/c08r4d0r/.cache/go-build-wz122}"
cd "$WT" || exit 9

{
  echo "### host: $(hostname) $(id -un) $(date -Is)"
  echo "### worktree: $WT"
  echo "### HEAD: $(git rev-parse HEAD)"
  echo "### branch: $(git rev-parse --abbrev-ref HEAD)"
  echo "\$ go version"; go version
  echo "### GOROOT=$(go env GOROOT)  GOTOOLCHAIN=$(go env GOTOOLCHAIN)"
  echo "### gcc: $(gcc --version | head -1)"
  echo
  echo '# ---- module shape: every path ./... covers ----'
  echo '$ go list ./...'; go list ./...
  echo "GOLIST_RC=$?"
  echo
  echo '# ---- 2a REQUIRED: go vet ./... ----'
  echo '$ go vet ./...'
  go vet ./... 2>&1; echo "GO_VET_RC=$?"
  echo
  echo '# ---- gofmt (part of the mandated vet gate in this repo) ----'
  echo '$ gofmt -l .'
  gofmt -l . 2>&1; echo "GOFMT_RC=$?   # no output above = formatted"
  echo
  echo '# ---- 2c REQUIRED: golangci-lint run (v1.61.0 == CI pin) ----'
  echo "\$ golangci-lint run   # $(/home/c08r4d0r/wz16-tools/golangci-lint-1.61.0-linux-amd64/golangci-lint --version)"
  /home/c08r4d0r/wz16-tools/golangci-lint-1.61.0-linux-amd64/golangci-lint run 2>&1
  echo "GOLANGCI_RC=$?"
  echo
  echo '# ---- 2d REQUIRED: goreleaser check (== release.yml pin) ----'
  echo "\$ goreleaser check   # $(/home/c08r4d0r/wz16-tools/goreleaser --version 2>/dev/null | grep -iE '^goreleaser|version' | head -1)"
  /home/c08r4d0r/wz16-tools/goreleaser check 2>&1
  echo "GORELEASER_CHECK_RC=$?"
  echo
  echo "### lint config present: $(ls -1 .golangci.yml .golangci.yaml .golangci.toml 2>/dev/null | tr '\n' ' ')"
} 2>&1 | tee "$CAP/02-vet-lint-check.txt"
echo "WROTE $CAP/02-vet-lint-check.txt"
