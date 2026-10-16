#!/usr/bin/env bash
# wz16 e2e gate — pass 6. Step 5b: run detail for the evidence tip + code-tree identity.
set -uo pipefail
export PATH=/usr/local/go/bin:$HOME/wz16-tools:$HOME/go/bin:$PATH
WT=/home/c08r4d0r/worktrees/wz16-pass6
CAP=/home/c08r4d0r/wz16-captures-v6
REPO=felixfelix-bot/net4sats-wizard-go
BR=ci/add-github-actions-workflow-clean
ID="$1"
CODE_TIP=e85720cd39a41d4af8bdff5172752a827e25bb86
cd "$WT"
{
echo "### STEP 5b — CI run detail at the evidence tip — $(date -Iseconds)"
echo "\$ gh run view $ID --repo $REPO --json databaseId,conclusion,headSha,headBranch,event,status,workflowName,jobs"
gh run view "$ID" --repo "$REPO" --json databaseId,conclusion,headSha,headBranch,event,status,workflowName,jobs 2>&1 \
  | jq -c '{id:.databaseId,wf:.workflowName,branch:.headBranch,sha:.headSha,event:.event,status:.status,conclusion:.conclusion}'
echo
echo "-- job + step detail (proves which commands CI executed) --"
gh run view "$ID" --repo "$REPO" --json jobs 2>&1 \
  | jq -r '.jobs[] | "JOB \(.name): \(.conclusion)\n" + (.steps[]? | "    \(.conclusion)\t\(.name)")' 
echo
echo "local HEAD                  = $(git rev-parse HEAD)"
echo "origin/$BR = $(git rev-parse origin/$BR)"
RUN_SHA=$(gh run view "$ID" --repo "$REPO" --json headSha --jq .headSha)
echo "CI run $ID headSha   = $RUN_SHA"
[ "$RUN_SHA" = "$(git rev-parse HEAD)" ] && echo "SHA_MATCH = OK (green run is exactly the pushed tip)" || echo "SHA_MATCH = NO"
echo
echo "=== code-tree identity: this tip vs the code tip the snapshot/suite evidence came from ==="
echo "\$ git diff --name-only $CODE_TIP HEAD"
git diff --name-only "$CODE_TIP" HEAD
echo "\$ git diff --name-only $CODE_TIP HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'"
N=$(git diff --name-only "$CODE_TIP" HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum' | tee /dev/stderr | grep -c .)
echo "CODE_DIFF_LINES=$N   (0 => identical buildable source, so the CI-green run covers the verified snapshot)"
echo "HEAD^{tree}                = $(git rev-parse HEAD^{tree})"
echo "main.go blob               = $(git rev-parse HEAD:main.go)"
echo ".goreleaser.yaml blob      = $(git rev-parse HEAD:.goreleaser.yaml)"
echo
echo "\$ gh run list --repo $REPO --branch $BR --limit 3"
gh run list --repo "$REPO" --branch "$BR" --limit 3 2>&1
} 2>&1 | tee "$CAP/09-run-at-ec84984.txt"
echo "CAPTURED 09"
