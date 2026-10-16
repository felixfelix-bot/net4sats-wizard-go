#!/usr/bin/env bash
# wz16 e2e gate — pass 6. Step 5: GHA proof (gh run list) + cross-head identity check.
set -uo pipefail
export PATH=/usr/local/go/bin:$HOME/wz16-tools:$HOME/go/bin:$PATH
WT=/home/c08r4d0r/worktrees/wz16-pass6
CAP=/home/c08r4d0r/wz16-captures-v6
REPO=felixfelix-bot/net4sats-wizard-go
BR=ci/add-github-actions-workflow-clean
cd "$WT"
HEAD_SHA=$(git rev-parse HEAD)

{
echo "### STEP 5 — GHA proof — $(date -Iseconds)"
echo "local HEAD = $HEAD_SHA"
echo
echo "=== 5a. THE CARD'S LITERAL REPO DOES NOT EXIST (measured, not cited) ==="
echo "\$ gh run list --repo net4sats/tollgate-wizard-go --limit 5"
gh run list --repo net4sats/tollgate-wizard-go --limit 5 2>&1
echo "RC=$?"
echo
echo "\$ gh api repos/net4sats/tollgate-wizard-go --jq .message"
gh api repos/net4sats/tollgate-wizard-go --jq .message 2>&1
echo "RC=$?"
echo "-> fall back to the real repo: $REPO (public fork; upstream net4sats/net4sats-wizard-go)"
echo
echo "=== 5b. THE CARD'S LITERAL COMMAND, on the real repo ==="
echo "\$ gh run list --repo $REPO --limit 5"
gh run list --repo $REPO --limit 5 2>&1
echo "RC=$?"
echo
echo "=== 5c. runs for THIS branch ==="
echo "\$ gh run list --repo $REPO --branch $BR --limit 5"
gh run list --repo $REPO --branch "$BR" --limit 5 2>&1
echo "RC=$?"
echo
echo "=== 5d. the green run that IS the pushed tip (load-bearing check) ==="
NEWEST=$(gh run list --repo "$REPO" --branch "$BR" --limit 1 --json databaseId --jq '.[0].databaseId')
echo "newest run id on this branch = $NEWEST"
echo "\$ gh run view $NEWEST --repo $REPO --json databaseId,conclusion,headSha,headBranch,event,status,jobs,workflowName"
gh run view "$NEWEST" --repo "$REPO" --json databaseId,conclusion,headSha,headBranch,event,status,jobs,workflowName 2>&1 | jq -c '{id:.databaseId,wf:.workflowName,branch:.headBranch,sha:.headSha,event:.event,status:.status,conclusion:.conclusion,jobs:[.jobs[]|{name,conclusion}]}'
RUN_SHA=$(gh run view "$NEWEST" --repo "$REPO" --json headSha --jq .headSha)
RUN_CONC=$(gh run view "$NEWEST" --repo "$REPO" --json conclusion --jq .conclusion)
echo
echo "local HEAD          = $HEAD_SHA"
echo "CI run headSha      = $RUN_SHA"
echo "CI run conclusion   = $RUN_CONC"
if [ "$RUN_SHA" = "$HEAD_SHA" ] && [ "$RUN_CONC" = "success" ]; then
  echo "SHA_MATCH_CONCLUSION = OK  (the green run is exactly this pushed tip)"
else
  echo "SHA_MATCH_CONCLUSION = MISMATCH/other -> a fresh dispatch is required at the pushed tip"
fi
echo
echo "=== 5e. the run's step log (proves which commands CI actually ran) ==="
gh run view "$NEWEST" --repo "$REPO" --log 2>&1 | grep -E 'go build|go vet|go test|golangci|Run |##\[error\]' | sed 's/^[^ ]* //' | sort -u | head -20
echo
echo "=== 5f. CROSS-HEAD identity: is the OTHER candidate head's green CI covering the same code? ==="
echo "\$ gh run list --branch worker-heavy/t_0509c320-r3 --limit 3"
gh run list --repo "$REPO" --branch worker-heavy/t_0509c320-r3 --limit 3 2>&1
R3=$(git rev-parse origin/worker-heavy/t_0509c320-r3)
echo "r3 tip = $R3"
echo
echo "\$ git diff --name-only $R3 $HEAD_SHA -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'"
git diff --name-only "$R3" "$HEAD_SHA" -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'
echo
echo "blob identity for the build/CI-relevant paths (r3 vs this tip):"
for f in .goreleaser.yaml .github/workflows/ci.yml .github/workflows/release.yml go.mod go.sum main.go deploy.go discover.go ssh.go embed.go; do
  a=$(git rev-parse "$R3:$f" 2>/dev/null || echo MISSING)
  b=$(git rev-parse "$HEAD_SHA:$f" 2>/dev/null || echo MISSING)
  [ "$a" = "$b" ] && echo "  IDENTICAL  $f" || echo "  DIFFERS    $f"
done
} 2>&1 | tee "$CAP/07-step5-gha-proof.txt"
echo "PASS6 STEP 5 CAPTURED"
