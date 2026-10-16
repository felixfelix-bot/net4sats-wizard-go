#!/usr/bin/env bash
# wz16 tamper control — copy the captures, corrupt a few load-bearing facts, and require the
# verifier to FAIL on the copy. If it still passes, the checks are vacuous and the whole pass is void.
set -u
SRC=/home/c08r4d0r/wz16-captures-t_0509c320-p4
TMP=/home/c08r4d0r/wz16-tests-t_0509c320-p4/tamper
CAP="$TMP/tampered"
rm -rf "$TMP"; mkdir -p "$CAP"
cp "$SRC"/*.txt "$SRC"/*.log "$SRC"/*.out "$CAP"/ 2>/dev/null

# 1. pretend --version printed the default (dev) instead of the injected version
sed -i 's/^LINE=\[tollgate-wizard 0\.5\.1-snapshot\.9c37d2e\]$/LINE=[tollgate-wizard dev]/' "$CAP/06-version-smoke.txt"
# 2. pretend the checksum run failed
sed -i 's/^SHA256SUM_C_REAL_RC=0$/SHA256SUM_C_REAL_RC=1/' "$CAP/05b-artifacts-and-checksums.txt"
# 3. pretend CI concluded failure
sed -i 's/"conclusion":"success"/"conclusion":"failure"/' "$CAP/09-run-detail-and-log.txt"
# 4. pretend the verbose suite had a failure
sed -i 's/^FAIL  = 0$/FAIL  = 1/' "$CAP/04-test-race-summary.txt"

{
  echo "### host: $(hostname) $(date -Is)"
  echo "### tampered copy: $CAP   (4 facts corrupted)"
  echo
  echo '$ python3 wz16-90-adhoc-verify.py '"$CAP"'   # MUST exit non-zero'
  python3 /home/c08r4d0r/wz16-tools/wz16-90-adhoc-verify.py "$CAP" 2>&1
  echo "CONTROL_VERIFIER_EXIT=$?   # non-zero == the checks are NOT vacuous (required)"
} 2>&1 | tee "$SRC/12-tamper-control.txt"
echo "WROTE $SRC/12-tamper-control.txt"
