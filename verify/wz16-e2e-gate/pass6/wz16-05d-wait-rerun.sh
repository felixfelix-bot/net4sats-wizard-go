#!/usr/bin/env bash
# wait for the re-run of run 38924071342 to finish
set -uo pipefail
REPO=felixfelix-bot/net4sats-wizard-go
ID=38924071342
for i in $(seq 1 30); do
  sleep 10
  J=$(gh run view "$ID" --repo "$REPO" --json status,conclusion,jobs 2>/dev/null)
  S=$(echo "$J" | jq -r .status)
  echo "[$((i*10))s] status=$S  $(echo "$J" | jq -c '[.jobs[]|{n:.name,c:.conclusion}]')"
  if [ "$S" = "completed" ]; then
    echo "FINAL: $(echo "$J" | jq -c '{status,conclusion}')"
    break
  fi
done
