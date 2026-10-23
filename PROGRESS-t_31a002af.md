# PROGRESS — t_31a002af: Adjudicate identity API claims and emit AGENTS.md evidence table

Worktree root: `/home/c08r4d0r/repos/net4sats-mvp-v2` (linked worktree → `/home/c08r4d0r/wz16-verify/.git/worktrees/net4sats-mvp-v2`)
Host: c08r4d0r · Board: net4sats-mvp-v2 · Repo: `felixfelix-bot/net4sats-wizard-go`
HEAD at session start: `e5fedf0` · HEAD at inspection and at commit: `d5d2373`
Mode: read-only against code/docs; only report files written.

## Finding 1 — task 1 and task 2 produced nothing to consume — STATUS: resolved by re-deriving everything
- Board read-only: `t_1f51893a` (caller-side scan) and `t_e01d7b78` (server-side scan) both have `result` = NULL.
- Their card bodies were used as the method spec only. Every citation in the report was re-opened at that file:line in this session; no summary line was trusted.

## Finding 2 — the two AGENTS.md claims: server TRUE, caller FALSE — STATUS: resolved
- **(a) "API: GET /identity, POST /identity/reveal-seed"** (AGENTS.md:30) → **TRUE**, both routes verified at `tollgate-module-basic-go/src/main.go:1935` and `:1941`; POST-only guard at `:2144`.
- The asserted path and verb are exactly right — `/identity` and `/identity/reveal-seed` on `:2121`, no `/api/v1` prefix, no path params. The falsehood is the caller attribution below, not the path.
- **(b) "Wizard and admin UI call GET /identity and POST /identity/reveal-seed"** (AGENTS.md:47) → **FALSE, not found in repo**, in both halves:
  - Wizard: its whole HTTP surface is `main.go:619-623`; only outbound HTTP is a release-asset download (`deploy.go:1128-1129`); `git grep "/identity" -- '*.go'` → 0 hits; `git grep "reveal-seed"` → only AGENTS.md:30 and :47.
  - Admin UI: all seven vendored bundles carry 0 case-insensitive matches for `identity`/`reveal-seed`/`npub`; they call only `/whoami` and `/balance` (`admin/assets/balance-BNg56sOO.js:1`).

## Finding 3 — real contract differs from what the claim implies — STATUS: verified
- `/identity` documents GET but never checks the verb (`src/main.go:2121-2132` has no `r.Method`).
- `/identity/reveal-seed` is a mnemonic **recovery** endpoint, not a seed reveal: body is the raw 12-word mnemonic as **text** (`:2148`, `:2153`), response is the identity derived from it (`:2154`), loopback-only 403 (`:2140`), 1 KB cap, no auth headers.

## Finding 4 — the card's prior result came from another node and diverges on the admin half — STATUS: recorded, not imported
- Prior result (node c03rad0r, HEAD `f840385`) marks the admin UI caller TRUE citing `configurationwizzard/src/components/identity-panel.tsx:38`.
- No `configurationwizzard` checkout exists on this host (`find` → nothing); this repo tracks 0 `.tsx` files. Recorded as "not found in repo" rather than borrowing an unopenable citation.

## Finding 5 — HEAD moved mid-session, AGENTS.md did not — STATUS: verified
- `e5fedf0` → `d5d2373` (another card's fe2b commit landed in this shared tree).
- `AGENTS.md` blob is `f2ff177` at both, and `git log -- AGENTS.md` returns only the introducing commit `475a9ef`. Claims audited are current, not stale-vs-HEAD.

## Finding 6 — four adjacent defects in the same AGENTS.md block — STATUS: recorded with cites
- `:41`/`:42`: derivation is HKDF-SHA256 (salt `tollgate-v1`, seps `tollgate-ipv4-v1:` / `tollgate-mac-v1:<iface>:`), not `SHA256(npub + "...")` — `src/identity/identity.go:260`, `:111`, `:125`.
- `:41`: CGNAT `100.64.0.0/10`, not `10.X.Y.1/24` — `:58`, `:114`.
- `:43`: mnemonic is **12 words**, not 24; and `PrivateKeyToMnemonic` does not exist in the backend (`grep` → 0 hits) — `:133`, `:137`.
- `:22`: shipped uhttpd listens on 8090 and 8080, not "port 80 for captive portal" — `uhttpd_net4sats:2`, `:12`.

## Files touched by this run
- CREATED `REPORT-t_31a002af.md` — full citable report (adjudication table, header block, distinguishing details, candidate statements).
- CREATED `PROGRESS-t_31a002af.md` — this file.
- No source, config, AGENTS.md, or test file modified. No tests run.

## Delivery — STATUS: pushed, observed
- Local commit `841538e` in the shared worktree (its parent `d5d2373` is unpushed, so not pushed from there).
- Shipped branch `worker-tollgate/t_31a002af` = `origin/main` `8fcc30e` + one cherry-picked commit `42043e5`.
- Verified: `git ls-remote origin refs/heads/worker-tollgate/t_31a002af` → `42043e579132d1868b214d54f57392700190e797`, equal to local HEAD.
- No PR opened (card asks for a report, not a merge).

## Remaining work
- None for this card's acceptance. If the author wants the **admin-UI half** asserted rather than dropped, a `configurationwizzard` checkout must be made available — it is not on this host.
