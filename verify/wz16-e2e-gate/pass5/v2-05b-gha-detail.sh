#!/usr/bin/env bash
# wz16v2 STEP 5e — run detail for the tip run + tree/CI assertions
set -u
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
OUT="$CAP/07b-tip-run-detail.txt"
RUN=38844256545
cd "$REPO" || exit 9
HEAD_SHA=$(git rev-parse HEAD)
RUN_SHA=$(gh run view "$RUN" --repo felixfelix-bot/net4sats-wizard-go --json headSha -q .headSha)

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### run id under scrutiny: $RUN"
  echo "local HEAD SHA        = $HEAD_SHA"
  echo "CI run headSha        = $RUN_SHA"
  if [ "$HEAD_SHA" = "$RUN_SHA" ]; then
    echo "SHA_MATCH             = OK  (the green run IS this exact tip)"
  else
    echo "SHA_MATCH             = MISMATCH"
  fi
  echo
  echo '# ---- gh run view 38844256545 (detail) ----'
  echo '$ gh run view 38844256545 --repo felixfelix-bot/net4sats-wizard-go'
  gh run view "$RUN" --repo felixfelix-bot/net4sats-wizard-go 2>&1
  echo
  echo '# ---- gh run view --json (conclusion/status/workflow) ----'
  gh run view "$RUN" --repo felixfelix-bot/net4sats-wizard-go \
     --json databaseId,conclusion,status,headBranch,headSha,event,workflowName,jobs 2>&1
  echo
  echo '# ---- CI workflow + goreleaser config present at this tip? ----'
  echo '$ git ls-tree -r --name-only HEAD | grep -E "workflows|\\.goreleaser"'
  git ls-tree -r --name-only HEAD | grep -E 'workflows|\.goreleaser' | sort
  echo
  echo '# ---- last 6 commits on the branch (docs/ci/goreleaser must be committed) ----'
  echo '$ git log --oneline -6'
  git log --oneline -6
} > "$OUT" 2>&1
cat "$OUT"
