#!/usr/bin/env bash
# wz16v2 STEP 5 — GHA proof: gh run list + run detail + HEAD-SHA match
set -u
CAP=/home/c08r4d0r/wz16-captures-v2
REPO=/home/c08r4d0r/wz16-verify
BR=ci/add-github-actions-workflow-clean
OUT="$CAP/07-gh-run-list-and-detail.txt"
cd "$REPO" || exit 9
HEAD_SHA=$(git rev-parse HEAD)

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### local HEAD (branch $BR): $HEAD_SHA"
  echo "### remote SHA (ls-remote): $(git ls-remote origin refs/heads/$BR | awk '{print $1}')"
  echo
  echo '# ---- STEP 5a: gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5 ----'
  echo '$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5'
  gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5
  echo
  echo '# ---- STEP 5b: the card names net4sats/tollgate-wizard-go — that repo does NOT exist ----'
  echo '$ gh run list --repo net4sats/tollgate-wizard-go --limit 5'
  gh run list --repo net4sats/tollgate-wizard-go --limit 5 2>&1
  echo "NET4SATS_TGW_RC=$?"
  echo
  echo '$ gh repo view net4sats/tollgate-wizard-go'
  gh repo view net4sats/tollgate-wizard-go 2>&1
  echo "REPO_VIEW_RC=$?"
  echo
  echo '# ---- STEP 5c: gh run list --repo net4sats/net4sats-wizard-go (the real upstream) --limit 5 ----'
  echo '$ gh run list --repo net4sats/net4sats-wizard-go --limit 5'
  gh run list --repo net4sats/net4sats-wizard-go --limit 5 2>&1
  echo "NET4SATS_NWG_RC=$?"
  echo
  echo '# ---- STEP 5d: JSON view of the 5 runs (conclusion + headSha + event) ----'
  echo '$ gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5 --json databaseId,conclusion,status,headBranch,headSha,event,workflowName,createdAt'
  gh run list --repo felixfelix-bot/net4sats-wizard-go --limit 5 \
     --json databaseId,conclusion,status,headBranch,headSha,event,workflowName,createdAt
} > "$OUT" 2>&1
cat "$OUT"
