# PROGRESS — WZ1.7 caller-side identity/seed request scan (net4sats-mvp-v2)

Mode: READ-ONLY static inspection. No edits to repo source, no test runs, no network calls.

| # | Finding | Status | Files touched |
|---|---------|--------|---------------|
| 0 | Repo root confirmed `/home/c08r4d0r/repos/net4sats-mvp-v2`, HEAD `c60cb4942a39e58a7c5b227b858696651a35e961`. Languages: 20 Go, 9 JS (7 minified bundles), 7 HTML, 3 sh, 1 py, 1 Makefile. | DONE | — |
| 1 | Broad keyword scan (`identity\|reveal[-_]seed\|revealSeed\|/seed\|mnemonic\|seed_phrase\|seedPhrase`) → 20 hits, ALL non-caller: AGENTS.md/README prose, `deploy.go` owner-identity uci ops, `defaults_test.go` "preseeded" prose, test JSON fixture `{"identity":"owner"}`, minified crypto/bip32 `seed` vars. Zero endpoint-path hits. | DONE | — |
| 2 | Client/outbound idiom sweep (Go `http.Client/Get/NewRequest`, JS `fetch(`/XHR, sh/CI `curl`/`wget`). Go has exactly ONE outbound helper `httpGetFile` (deploy.go:1219); every `Be()`/`Te()`/`Fe` JS base URL enumerated and resolved. | DONE | — |
| 3 | Server-route definitions deliberately excluded (belongs to sibling server scan). Repo serves only `/api/scan`, `/api/wifi-scan`, `/api/deploy`, `/api/status/`, `/` + `/ubus` (rpcd). None identity/seed. | DONE | — |
| 4 | VERDICT: `GET /identity` = FALSE (not found). `POST /identity/reveal-seed` = not found in repo. | DONE | — |
| 5 | Deliverable written to `REPORT-t_wz17-caller-identity-scan.md`; committed + pushed in a bounded path-scoped commit. | DONE | REPORT-t_wz17-caller-identity-scan.md, PROGRESS-t_wz17-caller-identity-scan.md |

Umbrella `PROGRESS.md` intentionally NOT touched (shared, owned by last writer).
