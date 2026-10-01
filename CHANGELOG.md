# Changelog

All notable changes to the net4sats wizard are documented here.
This project loosely follows [Keep a Changelog](https://keepachangelog.com/)
and [Semantic Versioning](https://semver.org/).

## [Unreleased]

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

- **CORS write-path gap closed.** The origin allowlist introduced by #24 only
  withheld `Access-Control-Allow-Origin`; it did not refuse the request, so a
  cross-origin "simple request" `POST` (no preflight, e.g.
  `Content-Type: text/plain`) still reached `/api/deploy` — the endpoint that
  drives the root SSH session. The middleware now rejects any request whose
  `Origin` is present and not allowlisted with `403` before the handler runs,
  for reads and writes alike. Non-browser clients that send no `Origin`
  (`curl`, the CLI, E2E harnesses) and same-origin browser callers are
  unaffected. A concrete (non-wildcard) `WIZARD_BIND` address is added to the
  allowlist so the UI served from it still writes; a wildcard bind (or a
  spoofed `Host` header) never widens it.
