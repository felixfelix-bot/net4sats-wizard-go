# Changelog

All notable changes to the net4sats wizard are documented here.
This project loosely follows [Keep a Changelog](https://keepachangelog.com/)
and [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed

- **Payout-share default is now `0%`, and no share is routed unless the
  operator sets one.** The advanced payout-share control was preseeded with
  `10`, so an operator who never opened Advanced settings silently routed a
  tenth of every payment to the developer payout address — they had to notice
  the control and move it to opt *out*. The slider now starts at `0` (label,
  readout and hint all say `0%`), so the default sends the full share to the
  operator and any other value is a deliberate choice. The control itself is
  unchanged: it still spans `0–50` and is still the only thing that populates
  `config.json` `profit_share`.

  This is the wizard's own downstream default. It does not touch any upstream
  (unbranded) default, and the `devSplit` wire field, the `devsplit` /
  `devsplit-val` DOM ids and the persisted `developer` identity name are all
  kept, so existing clients and already-deployed `config.json` files keep
  working. Guard: `TestDevSplitDefaultIsZero` plus
  `TestDevSplitSliderIsAdvancedOnly` (`defaults_test.go`), both verified RED
  against the old markup before the fix.

- **Test-network mints are no longer injected into every deploy.** The deploy
  step used to merge two `testnut` (Cashu test-network) mints into
  `config.json` `accepted_mints` unconditionally. Those mints issue worthless
  test ecash, so a production router accepting them lets a client pay with
  value the operator cannot redeem. The production set is now unconditional and
  the test set is added **only** when the wizard is started with
  `--test-mode`:

  ```sh
  ./net4sats-wizard --test-mode   # test router: testnut mints included
  ./net4sats-wizard               # default: production mints only
  ```

  Test mode is a **process-level** switch, not a request field: the wizard's
  HTTP API is reachable by anything running on the host, and a request-supplied
  flag would let a single `POST /api/deploy` opt itself into test mints.
  `TestDeployRequestHasNoTestModeField` enforces that. The streamed deploy log
  names the mints it wrote and calls out the test set when it is active.

  Guards: `TestNoTestnutMintsWithoutTestMode`,
  `TestDefaultMintsIncludeTestnutOnlyUnderFlag`, `TestConfigCommandGatesTestnut`
  and `TestParseOptionsTestMode` (`defaults_test.go`), verified RED before the
  fix; the jq round-trip in `security_hardening_test.go` now asserts both the
  off and on cases against a real `jq`.

- **Internal operator context removed from every shipped surface.** The wizard
  is intended to be donated/handed over, so the operator's private working
  context has no business in it. Removed: a personal laptop model and first name
  used as machine labels in `discover.go`'s discovery comments; an internal
  kanban card id + branch name in `ssh.go`'s retained-helper comment; an
  internal handover process name in `deploy.go` and `deploy_test.go`; an
  internal decision-gate id in this file; and the public `README.md` link text.

  The operator runbook is renamed from its old handover-derived filename to
  **`docs/router-setup.md`** (with the README link updated in the same commit),
  because the filename itself carried an internal process name onto a public doc
  path.

  No behaviour changed — every edit on the Go side is a comment; the discovery
  candidate IP list, the pinned download URLs and the deploy logic are byte-identical.

  The regression guard is `TestNoInternalProcessLeaksOnShippedSurfaces`
  (`debrand_internal_test.go`). It scans the embedded UI, the packaged
  `portal/`/`admin/`/`rpcd/` trees, the non-test Go sources and the docs for the
  internal-detail family, and was verified RED (8 distinct leaks reported) before
  any fix and GREEN after, plus a negative control that re-injected the retired
  wording and confirmed the guard fails.

  Deliberately retained, and recorded in the audit note rather than silently
  "fixed": the `10.47.41.1` probe candidate (functional discovery data, not
  branding) and the `felixfelix-bot/...` release URLs (live artifact pins —
  renaming them breaks the install; a pin migration is separate upstream work).

  A rename has two halves, so a second guard covers the half a content scan
  cannot see: **`TestNoInternalProcessNamedFiles`** fails if any published file
  *name* still carries internal detail. It exists because the first attempt at
  the runbook rename committed only the new path and left the old filename
  published; the guard was verified to fail on that exact state and pass once
  the stale file was removed.

- **Retired control name removed from the remaining operator-visible copy.**
  The previous de-branding pass neutralised the advanced control's *hint* but
  left the control's own **label** reading `Dev split`, the deploy-log line
  echoing the raw `devSplit` field name, and `README.md` naming the setting as
  "dev split" in two places. All three now read **payout share** — the wording
  the hint already used. No behaviour changed: the slider, its `0–50` range, its
  `10%` default and the `profit_share` factors it writes are untouched, so the
  already-decided upstream-default question (`Q3=B`, upstream default stays
  as-is) is not re-litigated here.

  The regression guard now covers the surfaces it previously missed:
  `TestNoSplitWordingInVisibleCopy` strips `<script>`/`<style>` bodies and all
  tags before matching, so it polices genuinely operator-visible text without
  tripping on the retained `devsplit` / `devsplit-val` DOM ids;
  `TestReadmeUsesNeutralPayoutWording` keeps the README in step; and
  `TestDeployLogUsesNeutralPayoutWording` covers the streamed log line. All
  three were verified to fail against the old copy before the fix.

  Still deliberately unchanged, as protocol/config identifiers rather than copy:
  the `devSplit` JSON wire field of `POST /api/deploy`, the `devsplit` /
  `devsplit-val` DOM ids, the `devSplit` Go local, and the `developer` identity
  name inside the persisted `config.json` `profit_share` array. Renaming any of
  those is a wire/config migration with a compatibility cost, not a copy change.

- **Revenue-share branding removed from the wizard UI.** The advanced
  dev-split control now reads as what it is — an optional share of payments
  routed to the developer payout address (`config.json` `profit_share`), with
  the documented 0 opt-out. The retired marketing phrasing is gone from the
  embedded UI, the packaged portal/admin HTML, the code comments describing the
  `devSplit` field, and this changelog; a regression test
  (`debrand_copy_test.go`) fails the build if it comes back.

  Unchanged on purpose: the `devSplit` JSON wire field of `POST /api/deploy`,
  the `devsplit` / `devsplit-val` DOM ids, the `devSplit` local in `deploy.go`,
  and the `developer` identity name inside the persisted `config.json`
  `profit_share` array (also used in `security_hardening_test.go` fixtures).
  Renaming any of those is a protocol/config migration, not a copy change, and
  is tracked separately rather than slipped into this commit.

### Fixed

- **Root shell-injection sinks closed.** All five deploy-command sinks
  (password, SSID, WPA key, LNURL callback, mint URL) now pass values
  through base64 encoding and stdin redirection instead of shell-string
  interpolation, eliminating the injection vector that allowed LAN drive-by
  RCE via the `/api/deploy` endpoint. The wizard HTTP server is now bound
  to loopback only (`127.0.0.1:8099`) and CORS is restricted to
  `localhost` origins, preventing cross-origin abuse from other LAN
  hosts.

- **Release artifacts now verified with sha256 pins.** Downloaded
  tollgate-wrt packages are checked against pinned sha256 digests before
  any SSH push to the router. Both the laptop-side download path and the
  router-side `wget` fallback verify the digest; a mismatch or missing pin
  aborts the deploy with an explicit error. Unpinned URLs fail closed by
  default.

- **IPK filename extraction hardened.** `extractIPKFilename` now
  rejects metacharacters and requires a strict `^<pkg>_[0-9][A-Za-z0-9._-]*\.ipk$`
 pattern, preventing traversal or injection via crafted package names in
  the nodogsplash download scrape.

- **Router-side wget fallback size-check bracket typo fixed.** A doubled
  bracket `]]` was glued to the filename, making the `WGET_OK` path
  unreachable. The fallback arm of the sha256 gate is now reachable.
