#!/usr/bin/env bash
# wz16 STEP 1 — clean tree check + push + ahead/behind vs remote (real exit codes, no pipes).
set -u
REPO=/home/c08r4d0r/wz16-verify
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
BR=ci/add-github-actions-workflow-clean
mkdir -p "$CAP"
OUT="$CAP/01-clean-tree-and-push.txt"
{
  echo "### host: $(hostname) $(id -un) $(date -Is)"
  echo "### repo: $REPO"
  echo "### branch: $(git -C "$REPO" rev-parse --abbrev-ref HEAD)"
  echo "### HEAD: $(git -C "$REPO" rev-parse HEAD)"
  echo
  echo '# ---- required check: `git status --porcelain` must be EMPTY ----'
  echo '$ git status --porcelain'
  git -C "$REPO" status --porcelain > /tmp/wz16-s1-porcelain.out 2>&1
  rc=$?
  cat /tmp/wz16-s1-porcelain.out
  echo "PORCELAIN_RC=$rc   PORCELAIN_LINES=$(wc -l < /tmp/wz16-s1-porcelain.out)   # 0 lines = clean tree"
  echo
  echo '# ---- required check: branch is pushed (`git push`) ----'
  echo '$ git push'
  git -C "$REPO" push > /tmp/wz16-s1-push.out 2>&1
  rc=$?
  cat /tmp/wz16-s1-push.out
  echo "PUSH_RC=$rc"
  echo
  echo '# ---- required check: ahead/behind is 0 vs the branch itself ----'
  echo "\$ git rev-list --left-right --count HEAD...origin/$BR"
  git -C "$REPO" rev-list --left-right --count "HEAD...origin/$BR" > /tmp/wz16-s1-ab.out 2>&1
  rc=$?
  cat /tmp/wz16-s1-ab.out
  echo "AB_RC=$rc   # format: <ahead>\\t<behind>"
  echo
  echo '# ---- required check: ahead/behind is 0 vs the PR base the card names ----'
  echo '$ git rev-list --left-right --count origin/main...HEAD   # remote-tracking base'
  git -C "$REPO" rev-list --left-right --count origin/main...HEAD > /tmp/wz16-s1-ab2.out 2>&1
  echo "AB2_RC=$?   base_behind_tip_ahead"
  cat /tmp/wz16-s1-ab2.out
  echo
  echo '# ---- local vs remote SHA identity ----'
  echo "local_HEAD        = $(git -C "$REPO" rev-parse HEAD)"
  echo "origin/$BR = $(git -C "$REPO" rev-parse origin/$BR)"
  echo
  echo '# ---- remote ground truth (ls-remote, independent of local tracking refs) ----'
  echo "\$ git ls-remote origin refs/heads/$BR"
  git -C "$REPO" ls-remote origin "refs/heads/$BR" 2>&1
  echo "LSREMOTE_RC=$?"
} 2>&1 | tee "$OUT"
echo "WROTE $OUT"
