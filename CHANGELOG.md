# Changelog

All notable changes to the net4sats wizard are documented here.
This project loosely follows [Keep a Changelog](https://keepachangelog.com/)
and [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed

- **Retired control name removed from the remaining operator-visible copy.**
  The previous de-branding pass neutralised the advanced control's *hint* but
  left the control's own **label** reading `Dev split`, the deploy-log line
  echoing the raw `devSplit` field name, and `README.md` naming the setting as
  "dev split" in two places. All three now read **payout share** — the wording
  the hint already used. No behaviour changed: the slider, its `0–50` range, its
  `10%` default and the `profit_share` factors it writes are untouched, so the
  resolved FELIX Q3 gate (`Q3=B`, upstream default stays as-is) is not
  re-litigated here.

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
