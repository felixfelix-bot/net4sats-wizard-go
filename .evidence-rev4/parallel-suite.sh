#!/usr/bin/env bash
# Does the suite fail under concurrent load (contended go build cache / ports)?
set -u
cd "$(dirname "$0")/.." || exit 1
export PATH="/usr/local/go/bin:$HOME/wz16-tools:$PATH"

echo "nproc=$(nproc)"; uptime

pids=()
for i in 1 2 3; do
  ( go test ./... -count=1 > "/tmp/rev4-par-$i.txt" 2>&1; echo "p$i rc=$? :: $(tail -1 "/tmp/rev4-par-$i.txt")" ) &
  pids+=($!)
done
for p in "${pids[@]}"; do wait "$p"; done
echo "done"
