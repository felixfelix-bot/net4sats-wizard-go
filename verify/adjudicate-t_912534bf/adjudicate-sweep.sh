#!/usr/bin/env bash
# Adjudication sweep for card t_912534bf — per-term hit counts across the five
# raw run logs of the t_be2366ea evidence family, plus the trace-dump checks.
# READ-ONLY. Prints one "term | counts" row per term.
#
# Usage: bash adjudicate-sweep.sh [EVIDENCE_DIR]
set -u
EVD="${1:-/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic}"
SPEC="${2:-/home/c08r4d0r/worktrees/t_11d0873e/wt/e2e}"
cd "$EVD" || { echo "FATAL: evidence dir not readable: $EVD"; exit 3; }

LOG="raw-run.log"
FAM="raw-run.log raw-runA.log raw-runA2.log raw-runC.log raw-runD.log"

echo "EVIDENCE_DIR: $PWD"
echo "SPEC_DIR:     $SPEC"
echo

echo "== A. family identity =="
wc -lc $FAM
sha256sum $LOG
echo

echo "== B. per-term counts (order: raw-run, A, A2, C, D) =="
while IFS='|' read -r label pat; do
  [ -z "$label" ] && continue
  counts=$(grep -ciE "$pat" $FAM | cut -d: -f2 | tr '\n' ' ')
  printf '%-22s %s\n' "$label" "$counts"
done <<'TERMS'
192.168.1.1|192\.168\.1\.1
2121|2121
2050|2050
ftp|ftp
sftp|sftp
endo|endo
onboarding|onboarding
Host:|Host:
GET-space|GET 
POST-space|POST 
HTTP/1-dot|HTTP/1\.
connected|connected
220-space|(^|[^0-9])220 
SSH-2-dot0|SSH-2\.0
curl|curl
localhost|localhost
127.0.0.1|127\.0\.0\.1
8099|8099
mt6000|mt6000
device|device
TERMS
echo

echo "== C. primary log: negative control (grep surface is live) =="
grep -n -i -F 'Running 7 tests' "$LOG"
printf 'line count: '; grep -c '' "$LOG"
printf 'ESC bytes  : '; grep -c $'\x1b' "$LOG"
echo

echo "== D. primary log: device-facing steps and totals (verbatim) =="
grep -n -E 'Step 6: Deploy|Step 7: Live|did not run|passed \(' "$LOG"
echo

echo "== E. trace dump: snapshot inventory =="
printf 'resource-snapshot lines: '; grep -c '^{"type":"resource-snapshot"' traces-network-dump.txt
grep '^{"type":"resource-snapshot"' traces-network-dump.txt \
  | grep -o '"url":"[^"]*"\|"status":[0-9]*\|"serverIPAddress":"[^"]*"\|"_serverPort":[0-9]*' \
  | paste - - - -
echo

echo "== F. trace dump: device endpoint in .network files (verbatim) =="
sed -n '12,13p' traces-network-dump.txt
echo

echo "== G. trace dump: discovery body IPs =="
grep -o '"ip":"[0-9.]*"' traces-network-dump.txt | sort -u
echo

echo "== H. spec literals (repo source, NOT a log) =="
grep -n "192\.168\.1\.1\|2121\|2050\|8099" \
  "$SPEC/generic-wizard-7step.spec.mjs" \
  "$SPEC/generic-wizard-7step.config.mjs" 2>&1
echo

echo "== I. synthetic-fixture locations of the literal 192.168.1.1:2121 in this evidence dir =="
grep -n "192\.168\.1\.1:2121" verify-scripts-adhoc.sh .repass-verify.sh repass2026-10-16-vacuity.sh 2>&1 | cut -c1-160
echo

echo "== J. END =="
