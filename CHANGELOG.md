# Changelog

All notable changes to the net4sats wizard are documented here.
This project loosely follows [Keep a Changelog](https://keepachangelog.com/)
and [Semantic Versioning](https://semver.org/).

## [v0.7.0-alpha24] — 2026-10-01

Security release: PR #24 (Amperstrand audit) merged onto the alpha23
field-fix lineage of the felixfelix-bot fork. Every alpha14–alpha23 fork
release predating this one carries the injection sinks described below —
upgrade before deploying routers.

### Added

- `--version` flag: prints the release tag and exits without binding a
  port, so `curl | bash` install paths can verify what they downloaded.

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
