#!/usr/bin/env bash
# wz16v2 STEP 5f — code-tree identity: prove the green CI tip and the snapshot tip
# differ only in non-code files, so the 6/6 + --version evidence applies to the CI-green SHA.
set -u
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
OUT="$CAP/08-code-tree-identity.txt"
cd "$REPO" || exit 9
{
  echo "### host: $(hostname) $(date -Is)"
  echo "### HEAD = $(git rev-parse HEAD)"
  echo
  echo '# ---- source tree hash at HEAD (the SHA the snapshot + smoke ran against) ----'
  echo '$ git rev-parse HEAD^{tree}'
  git rev-parse 'HEAD^{tree}'
  echo
  echo '# ---- every commit in the range that CI covered, with the files each touched ----'
  echo '$ for c in 64008dc cbbb092 4cee1d2 fadd82c b71af0b; do git show --stat --oneline $c | head -6; done'
  for c in 64008dc cbbb092 4cee1d2 fadd82c b71af0b; do
    git show --stat --oneline "$c" | head -8
    echo "  ---"
  done
  echo
  echo '# ---- are any .go / .goreleaser.yaml / workflow files modified between the oldest'
  echo '#      CI-green tip (b71af0b) and HEAD (64008dc)? (empty = code identical) ----'
  echo '$ git diff --name-only b71af0b 64008dc'
  git diff --name-only b71af0b 64008dc
  echo "CODE_TREE_DIFF_RC=$?"
  echo
  echo '$ git diff --name-only b71af0b 64008dc -- "*.go" ".goreleaser.yaml" ".github/**"'
  git diff --name-only b71af0b 64008dc -- '*.go' '.goreleaser.yaml' '.github/**'
  echo "GO_CODE_DIFF_LINES=$(git diff --name-only b71af0b 64008dc -- '*.go' '.goreleaser.yaml' '.github/**' | wc -l)  (expect 0)"
  echo
  echo '# ---- blob identity: main.go + .goreleaser.yaml blob SHAs at HEAD ----'
  git ls-tree HEAD main.go .goreleaser.yaml
} > "$OUT" 2>&1
cat "$OUT"
