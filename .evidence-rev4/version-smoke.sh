#!/usr/bin/env bash
# Step 4: --version smoke on the built linux/amd64 binary + manual ldflags injection proof.
set -u
cd "$(dirname "$0")/.." || exit 1
OUT=.evidence-rev4/04-version-smoke.txt
exec > >(tee "$OUT") 2>&1
export PATH="/usr/local/go/bin:$HOME/wz16-tools:$PATH"

BIN=dist/net4sats-wizard_linux_amd64_v1/net4sats-wizard

echo "===== STEP4a: snapshot binary --version ====="
echo "binary: $BIN"
file "$BIN"
"$BIN" --version
echo "[version_rc=$?]"
echo
echo "===== STEP4b: -version short form ====="
"$BIN" -version
echo "[short_rc=$?]"
echo

echo "===== STEP4c: manually build with injected ldflags value #1 ====="
go build -ldflags "-s -w -X main.version=v9.9.9-evidence" -o /tmp/rev4-v1 .
echo "[build1_rc=$?]"
/tmp/rev4-v1 --version
echo "[v1_rc=$?]"

echo "===== STEP4d: manually build with injected ldflags value #2 (must differ) ====="
go build -ldflags "-s -w -X main.version=v0.0.1-changed" -o /tmp/rev4-v2 .
echo "[build2_rc=$?]"
/tmp/rev4-v2 --version
echo "[v2_rc=$?]"

echo "===== STEP4e: control build with NO ldflags (must print 'dev') ====="
go build -o /tmp/rev4-v0 .
echo "[build0_rc=$?]"
/tmp/rev4-v0 --version
echo "[v0_rc=$?]"
