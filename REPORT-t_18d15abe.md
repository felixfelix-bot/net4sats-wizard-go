# REPORT — Preflight FE2b (t_18d15abe): verify FE2a merged, release tagged, capture feed URL

Worker pass: 2026-10-21 (started ~02:05Z) · node x280 · board net4sats-mvp-v2
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2 (origin felixfelix-bot/net4sats-wizard-go, detached HEAD 8a8d006 at session start)

## VERDICT: STOP — FE2b's precondition DOES NOT hold. Report to the manager, do not proceed to feed install.

The card's step-1 escape hatch fired verbatim: "If either is missing, STOP and report back to the
manager — do not proceed." The release IS tagged (v0.6.0, published 2026-10-21T01:17:07Z) but the
FE2a PR is NOT merged and the merged artifacts still contain NO feed index — the bridge FE2a was
to build does not exist. Only the release-tag leg holds; the FE2a leg does not.

Despite the STOP, the card also asks for a *failing* baseline check (TDD / QGATES v3.2). That part
IS deliverable without mutating anything: `scripts/fe2b_feed_check.sh` was written, run, and its
raw failing output captured (5 passed / 3 failed). Everything is read-only against the world
except that script and these records.

## 1. FE2b definition of done (step 1)

- Plan file `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` (read in full): contains NO
  FE2b definition of done — it is the unrelated "FELIX Q3 gate" (dev-split decision D-275a171c)
  document. No "FE2b"/"FE2a"/"feed" occurrence anywhere in it.
- Authoritative DoD from the board cards (bodies read read-only via sqlite3 this pass):
  - t_53879c67 (FE2b consumption, status blocked): "after FE2a merged + release published: on
    GL-MT6000 add feeds.conf src/gz line -> apk update -> install tollgate-wrt from FEED (not
    local scp) -> service up, portal serves. Proves feed consumable on real hardware. Serialize
    with fd4/wz18. QGATES v3.2: TDD (failing test first where code) | full suite green | docs
    same-commit | atomic conventional commits | PUSHED (done=pushed) | cold-review by
    reviewer-kimi before complete"
  - t_df4304d5 (FE2a, status todo): "rebase ci/auto-feed-indices onto main (take PR side
    wholesale for build-feed.yml, redo docs), ADD missing bridge job (workflow artifacts ->
    tagged release assets), OpenWrt arch tuple naming (aarch64_cortex-a53 etc.), gh-pages
    publishing (Q3 default), tag first release, SDK-container consumption test (feeds.conf src/gz
    -> apk update -> install -> deps resolve). Push branch, open/update PR -> status review.
    Manager merges."

## 2. Precondition check — the STOP (step 1, decisive)

| Precondition leg | Evidence (live, this pass) | Holds? |
|---|---|---|
| Release tagged/published | gh release list -R FreedomTechFeed/packages: **v0.6.0** published 2026-10-21T01:17:07Z, isLatest, not draft | YES |
| FE2a PR merged | gh pr list (30 most recent, all states): no merged PR adds a feed index / gh-pages / bridge. Branch ci/auto-feed-indices: API 404 Branch not found. Card t_df4304d5 board status **todo**. build-feed.yml workflow on master: contents API 404. gh-pages branch: 404. GitHub Pages site: 404 | NO |
| Feed index published | Latest release v0.6.0 = 26 flat assets (apk/ipk/offline + SHA256SUMS), no Packages index, no arch dir. Raw curl 404 at .../v0.6.0/Packages, .../v0.6.0/aarch64_cortex-a53/Packages, .../v0.6.0/Packages.gz. Code search src/gz in org: 0 results | NO |

Prior-lane context (read first, not re-derived): predecessor lane t_1b1cf42f (REPORT-t_1b1cf42f.md,
runs 1-3) reached the same negative on rc6 artifacts and recorded the structural explanation
(FE2a todo => no bridge => no index). This pass re-verified on the NEW v0.6.0 release published
today: the negative persists. PR #93 (release/v0.6.0, OPEN) touches only net/tollgate-wrt/Makefile —
a version-bump release PR, not the feed bridge.

## 3. Exact feed src/gz line (step 2)

Device: GL-MT6000, OpenWrt 24.10.2 (r28739), target mediatek/filogic, arch aarch64_cortex-a53,
**opkg-based** (no apk binary on the box; opkg version 38eccbb1 (2024-10-16)).

Exact src/gz line on the device today (line 6 of /etc/opkg/customfeeds.conf, read first-hand,
staged by the in-flight FE2b consumption lane — see §5):

    src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages

Live status of that URL: **HTTP 404** — the index does not exist. The line matches the card
template `src/gz tollgate <FEED_BASE_URL>/<arch>/packages` with
FEED_BASE_URL = https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6
and arch = aarch64_cortex-a53. No working feed base URL can be resolved from the release
artifacts because NO index is published under any path shape — the documented resolution
mechanism is a per-asset direct URL (scripts/tollgate-release-url.sh), not a src/gz index
(confirmed in the predecessor report §4; re-confirmed by today's flat asset list).

## 4. Serialization check (step 3) — CONFLICT DETECTED: hold and wait

- fd4/wz18 hardware/feed work: no fd4/wz18-named process is running (pgrep over fd4|wz18|fe2b:
  none). The GL-MT6000 bench lock (~/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock, 721 B,
  mode 600) is still held by t_2d0c6fa8 (FD4 preflight holder, claim 2026-10-14T05:38:34Z,
  release_condition ties release to the E2E run task t_557bb77d child 3) — an FD4-family hold.
- **HOWEVER a live competing lane IS running FE2b feed work right now**: process pid 2118663,
  profile worker-tollgate, brief "Install tollgate-wrt from the published FEED on GL-MT6000 ...
  Steps: add src/gz line, apk update, apk add tollgate-wrt ...". It was observed running at
  02:17Z this pass (ps output captured in PROGRESS). Its brief presumes the FE2a precondition
  holds — this preflight's finding says it does NOT (no index published; the device's own feed
  line 404s).
- The device evidence corroborates that lane's activity: /tmp holds freshly-staged local .ipk
  files dated 2026-10-20 16:40 (tollgate-wrt.ipk 8.6 MB, nodogsplash, jq) — i.e. a LOCAL scp-based
  install path was staged on the box, which is exactly what FE2b forbids ("install from FEED, not
  local scp").

Conclusion for step 3: serialization does NOT hold for this preflight to hand off to a feed
install — wait for the worker-tollgate lane to finish, and note its local-scp staging as a
finding for the manager.

## 5. GL-MT6000 reachability + package state (step 4) — all captured first-hand

Host x280 reaches the box on wlp59s0 (192.168.1.158/24). TCP 22/2121/2051 all open. SSH works
with the documented lab credential (sshpass, StrictHostKeyChecking=accept-new). Read-only state
captured verbatim this pass:

- /etc/openwrt_release: DISTRIB_RELEASE='24.10.2', DISTRIB_TARGET='mediatek/filogic',
  DISTRIB_ARCH='aarch64_cortex-a53'
- /etc/apk/repositories: **empty file** (this is an opkg box, not apk)
- /etc/apk/feeds.d: does not exist
- /etc/opkg/customfeeds.conf: contains the tollgate src/gz line (§3) plus the stock comment block
- opkg version 38eccbb1fd694d4798ac1baf88f9ba83d1eac616 (2024-10-16); no apk binary
- opkg list-installed | grep -i tollgate: **tollgate-wrt - 0.6.0_rc6-r1**
- /usr/lib/opkg/status for tollgate-wrt: Version 0.6.0_rc6-r1, Depends libc, nodogsplash, jq,
  ca-bundle; Status: install user installed; **no Source/origin field pointing at a feed**
- Service: /etc/init.d/tollgate-wrt present; process `/usr/bin/tollgate-wrt` (pid 12996) RUNNING;
  uhttpd instances on :8090/:8443 (admin), :8080/:443 (portal mirror), :2051 (captive portal)
- /tmp local staged artifacts (2026-10-20 16:40): tollgate-wrt.ipk (8,630,160 B), jq ipk,
  nodogsplash ipk — local-scp install path staged
- HTTP probes from the host: :2121 => 200 (API healthy); :2051/ => **403** (portal root refuses
  the bench host — vantage/expected-403 note: the portal serves splash.html, not /, to a
  pre-auth client; the check script records the 403 as the failing baseline)

## 6. The failing check script + raw baseline output (step 5)

Script: **scripts/fe2b_feed_check.sh** (chmod +x, sh -n clean). Host-side, read-only; encodes the
FE2b acceptance: feed index fetchable (A2), device carries feed line (A3), tollgate-wrt installed
(A4), origin is feed-not-local (A5), service running (A6), :2121 2xx (A7), portal 200 (A8), plus
SSH reachability (A1). Raw output of the 2026-10-21T02:23:16Z baseline run (exit 1):

    === fe2b_feed_check.sh — FE2b acceptance (TDD baseline; failing is the expected starting state) ===
    router   : 192.168.1.1
    feed_url : https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages
    arch     : aarch64_cortex-a53
    host     : x280  date: 2026-10-21T02:23:16Z

    --- A1: host can reach the router over SSH ---
    PASS  A1 ssh reachability to 192.168.1.1
    --- A2: the published FEED INDEX is fetchable (the feed must exist to be consumable) ---
    FAIL  A2 feed index fetchable: HTTP 404 at https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages (no published index -> nothing to install from)
    --- A3: device carries the feed line ---
    PASS  A3 device has a tollgate feed line
          src/gz tollgate https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/aarch64_cortex-a53/packages
    --- A4: tollgate-wrt is installed ---
    PASS  A4 tollgate-wrt installed
          tollgate-wrt - 0.6.0_rc6-r1
    --- A5: tollgate-wrt ORIGIN is the feed, not a local file ---
    FAIL  A5 origin is the feed: a LOCAL staged package file is present on the device (install was scp-based, not feed-based)
          local artifact(s): /tmp/tollgate-wrt.ipk
          opkg status origin fields:
          Status: install user installed
    --- A6: tollgate-wrt service is running ---
    PASS  A6 tollgate-wrt process is running
    --- A7: :2121 API health is 2xx ---
    PASS  A7 :2121 health HTTP 200
    --- A8: captive portal returns HTTP 200 ---
    FAIL  A8 portal HTTP 403 (expected 200)

    === RESULT: 5 passed, 3 failed ===
    FE2b NOT SATISFIED — this is the expected TDD baseline before the feed install.
    Assertions A2 (feed index) and A5 (feed origin) are the ones the feed install must flip.
    EXIT=1

Interpretation: A2 (feed index 404) is the STRUCTURAL failure — the FE2a bridge does not exist, so
no feed install can succeed yet. A5 (local staged .ipk in /tmp) is the provenance failure — the
box's current tollgate-wrt came via local scp, not a feed. A8 (portal 403 on / from the bench
host) is the vantage-shaped failure; the FE2b consumption lane should probe a portal path a
pre-auth client is actually served (e.g. /splash.html), and the assertion can be re-tuned then.
A2 and A5 are the two the feed install must flip.

## 7. What the manager should decide

1. FE2a (t_df4304d5, todo) must run and merge BEFORE FE2b consumption retries: build-feed.yml,
   the workflow-artifacts -> release-assets bridge, arch-tuple naming, gh-pages publishing, and
   the SDK-container consumption test. Without it there is no index for a src/gz line to address.
2. The live worker-tollgate lane (pid 2118663) is attempting a feed install whose premise this
   preflight refutes: its feed line 404s and /tmp shows a staged local .ipk path. The manager
   should stop or redirect that lane before it reports a fabricated "feed install".
3. If a feed-index shape is finally published, the target-qualified naming decision from the
   predecessor report (arch-only paths are ambiguous for mips_24kc) still applies: the index path
   must carry the target slot.

## 8. Honesty ledger

- STOP fired per the card's own instruction; no feed install was attempted, no device mutation
  made (SSH reads only; the /tmp .ipk files were NOT removed — they belong to the competing lane).
- The v0.6.0 release facts, PR list, branch/workflow/pages 404s, and the device state were all
  read live this pass. The predecessor reports were read before working; §2/§3 build on them.
- The baseline run output above is real command output (exit 1, 5 pass / 3 fail, captured
  2026-10-21T02:23:16Z). No fabricated URLs: the feed line documented is the one physically on
  the device, and its 404 is quoted.
- The plan file does not contain the FE2b DoD; the board cards are the source used, quoted above.
- The A8 403 is reported as the baseline, with the vantage caveat named rather than tuned away.
## 9. Verification addendum (2026-10-21, post-ship)

The shipped script was re-verified after commit:

- `sh -n scripts/fe2b_feed_check.sh` → clean (shellcheck not installed on this host).
- Committed bytes == on-disk bytes: `git show HEAD:scripts/fe2b_feed_check.sh | sha256sum` and
  `sha256sum scripts/fe2b_feed_check.sh` both `e6394ed5f1939e22c3a866dc82a75c7c9696830c3118f9c35c8b45c2991fae79`
  (at commit d5d2373).
- **Determinism**: 3 runs (02:23Z, and two 10 s apart post-commit) — identical verdict every time:
  5 PASS / 3 FAIL, same failing assertions (A2, A5, A8). True script exit code captured directly
  (not through a pipe): **exit 1**.
- `go vet ./...` → clean. `go test ./...` → 1 failure: `TestSshRunFunction`
  (coverage_test.go:18) — **pre-existing environment-coupled debt, unrelated to this change**:
  the test asserts `sshIdentify("192.168.1.1", "")` returns empty strings (assumes no reachable
  router), but the live lab router is up and answers, returning
  `firmware="OpenWrt 24.10.2 r28739-d9340319c6", vendor="OpenWrt", model="glinet,gl-mt6000"`.
  This change touched zero Go files (new shell script + docs only), so it cannot affect Go tests;
  the failure is the test's no-router assumption violated by the lab bench being powered — i.e.
  the exact condition FE2b depends on.
- Canonical recipe `hermes verify --json` (venv binary; the bare `hermes` on PATH dies in
  `import yaml` — different interpreter): `{"ok": false, "build": {"go build ./...": exit 0,
  PASS}, "test": {"go test ./...": exit 1, FAIL — only TestSshRunFunction}}`.
  Scoped isolation: `go test -run 'Test[^S]|TestS[^s]' ./...` → **`ok net4sats-wizard 45.152s`**
  — every test except the one env-coupled case is green. The sole red test fails because the lab
  router is powered on; it passes nowhere a reachable MT6000 exists. Concrete blocker for a fully
  green `hermes verify`: that pre-existing test, owned by other in-flight lanes (Gate-8 t_19b1b1c6
  adjudication explicitly cites "suite red"). Not repaired here: rewriting its no-router premise
  is a product decision for the owning lane, and this preflight is read-only for Go sources.
- Note: this shared worktree's detached HEAD advanced past d5d2373 (a competing lane committed
  t_31a002af report files on top). The FE2b deliverables remain committed at d5d2373 and pushed
  to `worker-heavy/t_18d15abe` (ls-remote verified full sha d5d2373682c3777b94356001a5987cd7a0cafc24).
