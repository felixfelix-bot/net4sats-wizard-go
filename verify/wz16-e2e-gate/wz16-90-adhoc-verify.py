#!/usr/bin/env python3
"""wz16 ad-hoc verifier — re-derives THIS pass's claims from the raw capture files.

Run:  wz16-90-adhoc-verify.py <captures-dir>
Exit 0 iff every claim about this pass's evidence holds. A companion control
(wz16-91-adhoc-control.sh) tampers a copy of the captures and requires this
script to FAIL on it, proving the checks are not vacuous.
"""
import re
import sys
from pathlib import Path

CAP = Path(sys.argv[1] if len(sys.argv) > 1 else "/home/c08r4d0r/wz16-captures-t_0509c320-p4")
HEAD = "9c37d2e04fe5bfafeb8d210202375ff7f34b773d"
VERSION = "0.5.1-snapshot.9c37d2e"
RUN = "38838430526"

fails = []
checks = 0


def text(name):
    p = CAP / name
    if not p.is_file():
        fails.append(f"MISSING capture file: {name}")
        return ""
    t = p.read_text(errors="replace")
    if not t.strip():
        fails.append(f"EMPTY capture file: {name}")
    return t


def need(label, cond, detail=""):
    global checks
    checks += 1
    if cond:
        print(f"  PASS  {label}")
    else:
        print(f"  FAIL  {label}   {detail}")
        fails.append(label)


def has(blob, needle, label):
    need(label, needle in blob, f"missing {needle!r}")


def rc(blob, line, value, label):
    m = re.search(rf"^{re.escape(line)}=(\d+)", blob, re.M)
    need(label, bool(m) and m.group(1) == str(value),
         f"{line} = {m.group(1) if m else 'ABSENT'}, want {value}")


# ---------------- STEP 1 ----------------
t01 = text("01-clean-tree-and-push.txt")
print("== STEP 1: clean tree + push + ahead/behind ==")
has(t01, "PORCELAIN_LINES=0", "01: git status --porcelain EMPTY (0 lines)")
rc(t01, "PORCELAIN_RC", 0, "01: porcelain rc 0")
has(t01, "Everything up-to-date", "01: git push was a no-op (branch already pushed)")
rc(t01, "PUSH_RC", 0, "01: git push rc 0")
need("01: rev-list count is 0\\t0 (ahead/behind == 0)",
     re.search(r"^0\t0$", t01, re.M) is not None)
has(t01, f"local_HEAD        = {HEAD}", "01: local HEAD == verified SHA")
has(t01, f"{HEAD}\trefs/heads/ci/add-github-actions-workflow-clean",
    "01: ls-remote remote ref == verified SHA (remote ground truth)")

# ---------------- STEP 2 ----------------
t02 = text("02-vet-lint-check.txt")
t03 = text("03-test-race.txt")
t04 = text("04-test-race-summary.txt")
print("== STEP 2: vet / test -race / golangci-lint / goreleaser check ==")
has(t02, "go version go1.22.12 linux/amd64", "02: go1.22.12 (CI-matching toolchain)")
has(t02, "net4sats-wizard", "02: single module net4sats-wizard (so ./... covers all)")
rc(t02, "GO_VET_RC", 0, "02a: go vet ./... rc 0")
rc(t02, "GOFMT_RC", 0, "02: gofmt -l rc 0 (no unformatted files)")
has(t02, "golangci-lint has version 1.61.0", "02c: golangci-lint v1.61.0 == CI pin")
rc(t02, "GOLANGCI_RC", 0, "02c: golangci-lint run rc 0")
rc(t02, "GORELEASER_CHECK_RC", 0, "02d: goreleaser check rc 0")
has(t02, "1 configuration file(s) validated", "02d: goreleaser validated the config")
rc(t03, "TEST_RACE_RC", 0, "02b: go test ./... -race rc 0")
has(t03, "ok  \tnet4sats-wizard", "02b: suite reported ok")
has(t04, "RUN   = 123", "02b: verbose 123 tests RUN")
has(t04, "FAIL  = 0", "02b: verbose 0 FAIL")
has(t04, "SKIP  = 0", "02b: verbose 0 SKIP")
has(t04, "DATA_RACE_WARNINGS = 0", "02b: zero 'WARNING: DATA RACE'")

# ---------------- STEP 3 ----------------
t05 = text("05-goreleaser-snapshot.txt")
t05b = text("05b-artifacts-and-checksums.txt")
t08 = text("08-checksum-controls.txt")
print("== STEP 3: snapshot release + artifacts + checksums ==")
rc(t05, "SNAPSHOT_RC", 0, "03: goreleaser release --snapshot --clean rc 0")
has(t05b, "ARCHIVE_COUNT=6", "03: exactly 6 archives")
has(t05b, "distinct_targets=6", "03: 6 distinct OS/arch targets")
for tgt in ["darwin_amd64", "darwin_arm64", "linux_amd64", "linux_arm64",
            "windows_amd64", "windows_arm64"]:
    has(t05b, tgt, f"03: target {tgt} present")
need("03: checksum count == artifact count (6 == 6)",
     re.search(r"^6  ==  6  \(checksum_count == artifact_count\)$", t05b, re.M) is not None)
rc(t05b, "SHA256SUM_C_REAL_RC", 0, "03: sha256sum -c real rc 0")
has(t05b, "OK_LINES=6  FAILED_LINES=0", "03: 6/6 checksum lines OK, 0 failed")
# two mandated spot checks, independently recomputed
need("03: spot-check#1 linux archive recomputed MATCH",
     re.search(rf"{re.escape('net4sats-wizard_'+VERSION+'_linux_amd64.tar.gz')} type=tar\.gz recomputed=([0-9a-f]{{64}}) expected=\1 MATCH", t05b) is not None)
need("03: spot-check#2 non-linux (darwin arm64) recomputed MATCH",
     re.search(rf"{re.escape('net4sats-wizard_'+VERSION+'_darwin_arm64.tar.gz')} type=tar\.gz recomputed=([0-9a-f]{{64}}) expected=\1 MATCH", t05b) is not None)
need("03: windows zip (Get-FileHash equivalent) recomputed MATCH",
     re.search(rf"{re.escape('net4sats-wizard_'+VERSION+'_windows_amd64.zip')} type=zip recomputed=([0-9a-f]{{64}}) expected=\1 MATCH", t05b) is not None)
# non-vacuity controls
print("== STEP 3b: checksum non-vacuity controls ==")
need("03b: CONTROL A positive rc 0, OK=6",
     re.search(r"CONTROL A[^\n]*\n(?:(?!CONTROL B).)*?REAL_RC=0\s*\nOK=6  FAILED=0", t08, re.S) is not None)
need("03b: CONTROL B (1 byte flipped) rc 1, exactly 1 FAILED",
     re.search(r"CONTROL B[^\n]*\n(?:(?!CONTROL C).)*?REAL_RC=1\s*\nOK=5  FAILED=1", t08, re.S) is not None)
has(t08, f"{VERSION}_linux_amd64.tar.gz: FAILED",
    "03b: the tampered artifact is the one reported FAILED")
need("03b: CONTROL C (artifact removed) rc 1, FAILED open or read",
     re.search(r"CONTROL C[^\n]*\n(?:(?!\Z).)*?REAL_RC=1\s*\nOK=5  FAILED=1", t08, re.S) is not None)
has(t08, "FAILED open or read", "03b: CONTROL C reports 'FAILED open or read'")

# ---------------- STEP 4 ----------------
t06 = text("06-version-smoke.txt")
t06b = text("06b-version-injection-change.txt")
print("== STEP 4: binary smoke --version + injection change ==")
has(t06, f"tollgate-wizard {VERSION}", "04: built binary prints 'tollgate-wizard <version>'")
need("04: --version output line exactly matches the injected version",
     re.search(rf"^LINE=\[tollgate-wizard {re.escape(VERSION)}\]$", t06, re.M) is not None)
has(t06, "exit_code=0", "04: --version exit code 0")
has(t06, "FOUND_DEV=no (good)", "04: version is NOT 'dev'")
has(t06b, "tollgate-wizard v9.9.9-manual-p4", "04: manual injection #1 visible")
has(t06b, "tollgate-wizard v0.0.0-second-different-value-p4", "04: manual injection #2 (changed)")
has(t06b, "tollgate-wizard dev", "04: no-ldflags control falls back to dev")
has(t06b, "main.version", "04: go tool nm locates the injected symbol")
has(t06b, 'var version = "dev"', "04: source anchor for the default value")

# ---------------- STEP 5 ----------------
t07 = text("07-gh-run-list-and-detail.txt")
t09 = text("09-run-detail-and-log.txt")
print("== STEP 5: GHA proof ==")
has(t07, "gh run list --repo net4sats/tollgate-wizard-go --limit 5",
    "05: the LITERAL card command was executed")
has(t07, "HTTP 404: Not Found", "05: literal command 404s (repo does not exist)")
rc(t07, "literal_real_rc", 1, "05: literal command real rc 1")
has(t07, "Could not resolve to a Repository with the name 'net4sats/tollgate-wizard-go'",
    "05: gh repo view cannot resolve the named repo")
has(t07, "statusCheckRollup\":[]", "05: upstream PR #35 has no checks (no .github/ upstream)")
need(f"05: >=1 successful run at the pushed head SHA {HEAD}",
     re.search(rf"successful_runs_at_{HEAD} = (\d+)", t07) is not None
     and int(re.search(rf"successful_runs_at_{HEAD} = (\d+)", t07).group(1)) >= 1)
has(t07, f"{RUN} completed success {HEAD}", "05: the run is completed + success at that SHA")
need("05: gh run view confirms conclusion success at that SHA",
     f'"conclusion":"success"' in t09 and f'"headSha":"{HEAD}"' in t09)
has(t09, "ok  \tnet4sats-wizard", "05: CI test job ran the suite (ok net4sats-wizard)")
has(t09, "go test ./... -race -count=1", "05: CI ran go test -race")
has(t09, "go vet ./...", "05: CI ran go vet")
has(t09, "golangci-lint found no issues", "05: CI lint job clean")

# ---------------- VERDICT ----------------
print()
print(f"TOTAL CHECKS = {checks}   FAILURES = {len(fails)}")
if fails:
    print("FAILED CHECKS:")
    for f in fails:
        print(f"  - {f}")
    sys.exit(1)
print("ALL CHECKS PASS")
sys.exit(0)
