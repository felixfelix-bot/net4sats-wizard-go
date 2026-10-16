#!/usr/bin/env bash
# wz16 pass 6 — final tip CI proof.
set -uo pipefail
export PATH=/usr/local/go/bin:$HOME/wz16-tools:$HOME/go/bin:$PATH
WT=/home/c08r4d0r/worktrees/wz16-pass6
CAP=/home/c08r4d0r/wz16-captures-v6
REPO=felixfelix-bot/net4sats-wizard-go
BR=ci/add-github-actions-workflow-clean
CODE_TIP=e85720cd39a41d4af8bdff5172752a827e25bb86
cd "$WT"
TIP=$(git rev-parse HEAD)
ID=""
for i in $(seq 1 30); do
  sleep 10
  J=$(gh run list --repo "$REPO" --branch "$BR" --limit 1 --json databaseId,status,headSha 2>/dev/null)
  S=$(echo "$J" | jq -r '.[0].status')
  SH=$(echo "$J" | jq -r '.[0].headSha')
  ID=$(echo "$J" | jq -r '.[0].databaseId')
  echo "[$((i*10))s] id=$ID status=$S headSha=${SH:0:7}"
  if [ "$S" = "completed" ] && [ "$SH" = "$TIP" ]; then break; fi
done

{
echo "### PASS-6 FINAL — CI at the final pushed tip — $(date -Iseconds)"
echo "\$ git rev-parse HEAD"; git rev-parse HEAD
echo "\$ git rev-parse origin/$BR"; git rev-parse "origin/$BR"
echo "\$ git status --porcelain"; git status --porcelain
echo "PORCELAIN_LINES=$(git status --porcelain | grep -c .)"
echo
echo "\$ gh run list --repo $REPO --branch $BR --limit 3"
gh run list --repo "$REPO" --branch "$BR" --limit 3 2>&1
echo
echo "\$ gh run view $ID --repo $REPO --json databaseId,conclusion,headSha,headBranch,event,status,jobs,workflowName"
gh run view "$ID" --repo "$REPO" --json databaseId,conclusion,headSha,headBranch,event,status,jobs,workflowName 2>&1 \
  | jq -c '{id:.databaseId,wf:.workflowName,branch:.headBranch,sha:.headSha,event:.event,status:.status,conclusion:.conclusion,jobs:[.jobs[]|{name,conclusion}]}'
echo
echo "-- step detail --"
gh run view "$ID" --repo "$REPO" --json jobs 2>&1 \
  | jq -r '.jobs[] | "JOB \(.name): \(.conclusion)\n" + (.steps[]? | "    \(.conclusion)\t\(.name)")'
echo
echo "FINAL_TIP      = $TIP"
echo "RUN_HEAD_SHA   = $(gh run view "$ID" --repo "$REPO" --json headSha --jq .headSha)"
echo "RUN_CONCLUSION = $(gh run view "$ID" --repo "$REPO" --json conclusion --jq .conclusion)"
echo
echo "\$ git diff --name-only $CODE_TIP HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'"
git diff --name-only "$CODE_TIP" HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum'
echo "CODE_DIFF_LINES=$(git diff --name-only "$CODE_TIP" HEAD -- '*.go' '.goreleaser.yaml' '.github' 'go.mod' 'go.sum' | grep -c .)"
} 2>&1 | tee "$CAP/10-final-tip-ci.txt"
echo "CAPTURED 10"
