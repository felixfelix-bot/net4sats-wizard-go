#!/usr/bin/env bash
# wz16 STEP 4 — binary smoke: --version from the built linux/amd64 artifact + injection-change proof.
set -u
CAP=/home/c08r4d0r/wz16-captures-t_0509c320-p4
WT=/home/c08r4d0r/wz16-verify
export PATH="/home/c08r4d0r/wz16-tools/go122/bin:$PATH"
export GOTOOLCHAIN=go1.22.12
export GOCACHE="${GOCACHE:-/home/c08r4d0r/.cache/go-build-wz122}"
SMOKE=/home/c08r4d0r/wz16-tests-t_0509c320-p4/smoke
rm -rf "$SMOKE"; mkdir -p "$SMOKE"

{
  echo "### host: $(hostname) $(id -un) $(date -Is)"
  echo "### scratch: $SMOKE"
  echo
  echo '# ---- extract the BUILT linux/amd64 binary from the snapshot archive (not a fresh build) ----'
  cd "$SMOKE" || exit 9
  echo "\$ tar -xzf $WT/dist/$(ls -1 "$WT"/dist | grep '_linux_amd64.tar.gz$')"
  tar -xzf "$WT"/dist/*_linux_amd64.tar.gz
  echo "tar_rc=$?"
  echo '$ ls -la'
  ls -la
  echo
  echo '$ file ./net4sats-wizard'
  file ./net4sats-wizard
  echo
  echo '# ---- REQUIRED: execute it with --version; expect "tollgate-wizard <version>", rc 0, not dev ----'
  echo '$ ./net4sats-wizard --version'
  ./net4sats-wizard --version > out.txt 2>&1
  ver_rc=$?
  cat out.txt
  echo "exit_code=$ver_rc"
  line=$(cat out.txt)
  echo "LINE=[$line]"
  case "$line" in
    "tollgate-wizard dev") echo "FOUND_DEV=YES  -> BAD (injection did not reach the binary)";;
    "tollgate-wizard "*)   echo "FOUND_DEV=no (good)   shape_ok=yes";;
    *)                     echo "FOUND_DEV=unknown  -> shape unexpected";;
  esac
  echo "injected_version_token=$(echo "$line" | awk '{print $2}')"
  echo
  echo '# ---- also -version and the version baked into the SNAPSHOT version string ----'
  echo '$ ./net4sats-wizard -version'
  ./net4sats-wizard -version; echo "dash_version_rc=$?"
} 2>&1 | tee "$CAP/06-version-smoke.txt"

# manual injection change proof: build twice with different -X main.version and with none
{
  echo "### host: $(hostname) $(date -Is)   worktree: $WT"
  echo '# ---- REQUIRED: build the injected value manually and confirm it CHANGES ----'
  cd "$WT" || exit 9
  for v in "v9.9.9-manual-p4" "v0.0.0-second-different-value-p4"; do
    echo
    echo "\$ go build -ldflags \"-X main.version=$v\" -o $SMOKE/wiz-manual ."
    go build -ldflags "-X main.version=$v" -o "$SMOKE/wiz-manual" . 2>&1; echo "build_rc=$?"
    echo "\$ $SMOKE/wiz-manual --version"
    "$SMOKE/wiz-manual" --version 2>&1; echo "run_rc=$?"
  done
  echo
  echo '# ---- control: plain build with NO ldflags must fall back to the dev default ----'
  echo "\$ go build -o $SMOKE/wiz-nodev ."
  go build -o "$SMOKE/wiz-nodev" . 2>&1; echo "build_rc=$?"
  echo "\$ $SMOKE/wiz-nodev --version"
  "$SMOKE/wiz-nodev" --version 2>&1; echo "run_rc=$?"
  echo
  echo '# ---- locate the injected symbol so the mechanism is proven, not asserted ----'
  echo "\$ go tool nm $SMOKE/wiz-manual | grep -F ' main.version'"
  go tool nm "$SMOKE/wiz-manual" 2>/dev/null | grep -F ' main.version'
  echo
  echo '# ---- source anchors ----'
  grep -n 'var version = ' main.go
  grep -n 'func versionString' main.go
  grep -n '"--version", "-version"' main.go
} 2>&1 | tee "$CAP/06b-version-injection-change.txt"
echo "WROTE 06-* and 06b-*"
