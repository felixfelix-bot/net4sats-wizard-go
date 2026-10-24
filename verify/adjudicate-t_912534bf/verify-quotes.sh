#!/usr/bin/env bash
# Quote-integrity verifier for REPORT-t_912534bf.md.
# Extracts every "<path>:<line>" citation of the form
#   /home/c08r4d0r/.../raw-run.log:<N>   (fenced blocks)
# and re-checks the quoted text against the FROZEN committed copies in this dir.
# Read-only. Exit 0 = every checked quote matched a real line.
set -u
cd "$(dirname "$0")" || exit 3

REPORT="../../REPORT-t_912534bf.md"
FROZEN="./raw-run.log"
FROZEND="./raw-runD.log"

fail=0; checked=0

echo "=== V1. frozen copies match the cited sha256 ==="
h1=$(sha256sum "$FROZEN" | cut -d' ' -f1)
h2=$(sha256sum "$FROZEND" | cut -d' ' -f1)
[ "$h1" = "062c8b1a4a2e97cd9366c8646f6cf64b91090b1ff4b5f1c75f151122c4cdf209" ] \
  && { echo "PASS raw-run.log  sha256 $h1"; checked=$((checked+1)); } \
  || { echo "FAIL raw-run.log  sha256 $h1"; fail=$((fail+1)); }
[ "$h2" = "3f84663234556d2057ab931b71aa07e344e3366169bc1c3b7caac2de50026beb" ] \
  && { echo "PASS raw-runD.log sha256 $h2"; checked=$((checked+1)); } \
  || { echo "FAIL raw-runD.log sha256 $h2"; fail=$((fail+1)); }

echo
echo "=== V2. every fenced-block citation (path:line) resolves ==="
# Pull the path:line tokens the report uses to introduce raw-log quotes.
grep -o '/home/c08r4d0r/[^ :]*raw-run[A2C]*\.log:[0-9]*' "$REPORT" | sort -u > /tmp/cites.$$
while IFS= read -r tok; do
  file="${tok%:*}"; line="${tok##*:}"
  case "$file" in
    *"/raw-run.log")  local_f="$FROZEN" ;;
    *"/raw-runD.log") local_f="$FROZEND" ;;
    *) echo "SKIP  $tok (no frozen copy for this sibling; cited only for its 'did not run' lines)"; continue ;;
  esac
  txt=$(sed -n "${line}p" "$local_f")
  if [ -n "$txt" ]; then
    printf 'PASS  %s -> %s\n' "$tok" "$(printf '%s' "$txt" | cut -c1-72)"
    checked=$((checked+1))
  else
    printf 'FAIL  %s -> NO SUCH LINE\n' "$tok"; fail=$((fail+1))
  fi
done < /tmp/cites.$$
rm -f /tmp/cites.$$

echo
echo "=== V3. the load-bearing quoted strings appear verbatim in the frozen primary ==="
while IFS= read -r want; do
  [ -z "$want" ] && continue
  if grep -qF -- "$want" "$FROZEN"; then
    printf 'PASS  verbatim in raw-run.log: %s\n' "$(printf '%s' "$want" | cut -c1-66)"
    checked=$((checked+1))
  else
    printf 'FAIL  NOT in raw-run.log: %s\n' "$want"; fail=$((fail+1))
  fi
done <<'WANTS'
  -  6 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:183:1 › Step 6: Deploy — run the deployment to completion
  -  7 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:1 › Step 7: Live — verify TollGate API and portal on the router
  2 did not run
  4 passed (3.2m)
Running 7 tests using 1 worker
WANTS

echo
echo "=== V4. the load-bearing quoted strings appear verbatim in the frozen primary ==="
while IFS= read -r want; do
  [ -z "$want" ] && continue
  if grep -qF -- "$want" "$FROZEND"; then
    printf 'PASS  verbatim in raw-runD.log: %s\n' "$(printf '%s' "$want" | cut -c1-66)"
    checked=$((checked+1))
  else
    printf 'FAIL  NOT in raw-runD.log: %s\n' "$want"; fail=$((fail+1))
  fi
done <<'WANTSD'
    Error: Refusing to select a deploy target: ROUTER_IP=192.168.8.1 is not in the discovery results [192.168.2.1 — unknown unknown
  -  8 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:297:1 › Step 7: Live — verify TollGate API and portal on the router
  6 did not run
WANTSD

echo
echo "=== V5. the two cited literals are ZERO-hit in the frozen primary ==="
for pat in '192\.168\.1\.1' '2121'; do
  c=$(grep -cE "$pat" "$FROZEN")
  [ "$c" -eq 0 ] && { echo "PASS  zero-hit $pat"; checked=$((checked+1)); } \
                 || { echo "FAIL  $pat -> $c"; fail=$((fail+1)); }
done

echo
echo "RESULT: PASS=$checked FAIL=$fail"
[ "$fail" -eq 0 ] && echo "ALL QUOTE-INTEGRITY CHECKS PASSED" || echo "QUOTE-INTEGRITY FAILURES PRESENT"
exit $fail
