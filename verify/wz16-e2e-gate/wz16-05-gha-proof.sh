#!/usr/bin/env bash
# wz16 STEP 5 — GHA proof.
# The card names net4sats/tollgate-wizard-go; that repo does not exist. We run the LITERAL command
# and report its exact error, then prove CI green at the pushed head SHA in the real repo + the fork
# (upstream main has no .github/, which is exactly what this PR adds -> the fork is the CI venue).
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
REPO=/home/c08r4d0r/wz16-verify
BR=ci/add-github-actions-workflow-clean
HEAD_SHA=$(git -C "$REPO" rev-parse HEAD)
FORK=felixfelix-bot/net4sats-wizard-go
UP=net4sats/net4sats-wizard-go

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### gh: $(gh --version | head -1)   account: $(gh api user --jq .login 2>/dev/null)"
  echo "### pushed head SHA: $HEAD_SHA   branch: $BR"
  echo
  echo '# ==== the LITERAL command from the card, against the repo the card names ===='
  echo '$ gh run list --repo net4sats/tollgate-wizard-go --limit 5'
  gh run list --repo net4sats/tollgate-wizard-go --limit 5 > /tmp/wz16-s5-lit.out 2>&1
  lit_rc=$?
  cat /tmp/wz16-s5-lit.out
  echo "literal_real_rc=$lit_rc"
  echo
  echo '$ gh repo view net4sats/tollgate-wizard-go'
  gh repo view net4sats/tollgate-wizard-go > /tmp/wz16-s5-rv.out 2>&1
  rv_rc=$?
  cat /tmp/wz16-s5-rv.out
  echo "repoview_real_rc=$rv_rc"
  echo
  echo '$ git ls-remote https://github.com/net4sats/tollgate-wizard-go'
  git ls-remote https://github.com/net4sats/tollgate-wizard-go > /tmp/wz16-s5-ls.out 2>&1
  ls_rc=$?
  cat /tmp/wz16-s5-ls.out
  echo "lsremote_real_rc=$ls_rc"
  echo
  echo '# ==== upstream: is there CI to read? (main has no .github/ -- added by this PR) ===='
  echo "\$ gh run list --repo $UP --limit 5"
  gh run list --repo "$UP" --limit 5 2>&1; echo "upstream_runs_rc=$?"
  echo "\$ gh api repos/$UP/contents/.github --jq .message"
  gh api "repos/$UP/contents/.github" --jq .message 2>&1; echo "upstream_ghdir_rc=$?"
  echo "\$ gh pr checks 35 --repo $UP"
  gh pr checks 35 --repo "$UP" 2>&1; echo "prchecks_rc=$?"
  echo "\$ gh pr view 35 --repo $UP --json number,state,headRefOid,statusCheckRollup"
  gh pr view 35 --repo "$UP" --json number,state,headRefOid,statusCheckRollup 2>&1
  echo
  echo '# ==== the fork: the CI venue while .github/ is not yet upstream ===='
  echo "\$ gh run list --repo $FORK --limit 5"
  gh run list --repo "$FORK" --limit 5 > /tmp/wz16-s5-fork.out 2>&1
  fork_rc=$?
  cat /tmp/wz16-s5-fork.out
  echo "fork_runs_rc=$fork_rc"
  echo
  echo "\$ gh run list --repo $FORK --branch $BR --limit 5 --json databaseId,conclusion,status,headSha,headBranch,event,workflowName,url,createdAt"
  gh run list --repo "$FORK" --branch "$BR" --limit 5 \
     --json databaseId,conclusion,status,headSha,headBranch,event,workflowName,url,createdAt 2>&1 \
     | tee /tmp/wz16-s5-forkjson.out
  echo
  echo '# ---- REQUIRED: a run at THE PUSHED HEAD SHA with conclusion success ----'
  MATCH=$(gh run list --repo "$FORK" --branch "$BR" --limit 20 \
     --json databaseId,conclusion,status,headSha \
     --jq "[.[] | select(.headSha==\"$HEAD_SHA\" and .conclusion==\"success\")] | length")
  echo "successful_runs_at_$HEAD_SHA = $MATCH"
  gh run list --repo "$FORK" --branch "$BR" --limit 20 \
     --json databaseId,conclusion,status,headSha,url \
     --jq ".[] | select(.headSha==\"$HEAD_SHA\") | \"\(.databaseId) \(.status) \(.conclusion) \(.headSha) \(.url)\""
  echo "match_rc=$?"
} 2>&1 | tee "$CAP/07-gh-run-list-and-detail.txt"
echo "WROTE $CAP/07-gh-run-list-and-detail.txt"
