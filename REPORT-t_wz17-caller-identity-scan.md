# REPORT — WZ1.7: Caller-side identity / seed-reveal request sites in the wizard repo

Task:      net4sats-mvp-v2 WZ1.7 (fleet offload, worker-heavy)
Repo root: /home/c08r4d0r/repos/net4sats-mvp-v2
HEAD:      c60cb4942a39e58a7c5b227b858696651a35e961
Branch:    worker-heavy/t_gate8-19b1b1c6
Host:      x280
Mode:      READ-ONLY. No repo source edited. No test run. No network call. No kanban mutation.
Method:    static inspection only (grep/rg content scan + minified-bundle literal extraction).

## 1. VERDICT (headline)

| Asserted call | Verdict | Evidence |
|---|---|---|
| any code issues `GET /identity` | **FALSE** — not found in repo | Zero occurrences of the token pair `identity` + HTTP call in any caller idiom. The only `identity` matches in code are the uci identity *record* ops in `deploy.go:873` and a test fixture string (`security_hardening_test.go:294`). |
| any code issues `POST /identity/reveal-seed` | **not found in repo** | No `reveal-seed`, `reveal_seed`, `revealSeed`, `reveal` (any case), `mnemonic`, `seed_phrase`, or `seedPhrase` match anywhere in code. `/identity/reveal-seed` appears ONLY as prose in `AGENTS.md:30` and `AGENTS.md:47`. |

The verb+path pair asserted in `AGENTS.md:30` ("API: GET /identity, POST
/identity/reveal-seed") and `AGENTS.md:47` ("Wizard and admin UI call `GET
/identity` and `POST /identity/reveal-seed`") is **NOT implemented and NOT
called anywhere in this repository**. The documentation claims a caller that
does not exist.

## 2. INVENTORY / LANGUAGE FACTS

Root: `/home/c08r4d0r/repos/net4sats-mvp-v2`

| ext | count | note |
|---|---|---|
| .go | 20 | 8 non-test sources + test files |
| .js | 9 | 7 are minified Vite bundles under `admin/assets`, `portal/assets` |
| .html | 7 | wizard `index.html`, admin/portal pages, venue templates |
| .sh | 3 | `scripts/fe2b_feed_check.sh`, `verify/*` harnesses |
| .py | 1 | `verify/wz15-tools` entrypoint |
| Makefile | 1 | nak/Nostr helpers only |

Relevant source file list (non-test): `main.go`, `deploy.go`, `discover.go`,
`ssh.go`, `embed.go`, favourites/tests. Plus `rpcd/tollgate` (shell rpcd
plugin), `rpcd/tollgate_acl.json`, `.github/workflows/ci.yml`.

## 3. FINDING TABLE — every CALLER-SIDE outbound request that EXISTS

Every row is an outbound request actually issued by this repo. Path column is
verbatim from the source. Reminder: server-route definitions are out of scope
here (sibling server scan owns them) — none of the rows below are route defs.

| Verb | Path (verbatim) | File:line | Verbatim snippet (≤3 lines) | Base URL / auth / body notes |
|---|---|---|---|---|
| GET | `Be()+`/`` → `/` on `http://<host>:2121` | portal/assets/splash-Be2Pt9NX.js:15680 | `fetch(Be()+`/`);if(!e.ok)throw Error(`Pricing fetch failed: ${e.status}`);` | Base: `function Be(){return`http://${window.location.hostname||`192.168.1.1`}:${ze}`}` with `var ze=2121`. No auth header. No body. (portal, client-side) |
| GET | `Be()+`/whoami`` → `/whoami` on :2121 | portal/assets/splash-Be2Pt9NX.js:15911 | `fetch(Be()+`/whoami`);if(!e.ok)throw Error(`Whoami fetch failed`);` | Same base. No auth. No body. MAC parsed out of the text body. |
| POST | `Be()+`/`` + optional `?mac=…` → `/` on :2121 | portal/assets/splash-Be2Pt9NX.js:16203 | `fetch(Be()+`/`+(t?`?mac=${encodeURIComponent(t)}`:``),{method:`POST`,headers:n,body:r})` | Same base. Body = Cashu token (var `r`). No Authorization header. |
| POST | `Be()+`/ln-invoice`` | portal/assets/splash-Be2Pt9NX.js:16637 | `fetch(Be()+`/ln-invoice`,{method:`POST`,headers:{"Content-Type":`application/json`},body:JSON.stringify(r)})).json()` | Same base. Content-Type JSON, no auth header. Body = amount/mint request. |
| GET | `Be()+`/ln-invoice?quote=${…}`` | portal/assets/splash-Be2Pt9NX.js:16887 | `fetch(n)).json()}async function qe(e,t){let n=Be()+`/ln-invo…` | Same base. Query param `quote`. No auth, no body. |
| POST | `/ubus` JSON-RPC (`Fe=/ubus`) | admin/assets/index-Cv-gJynA.js:14439 | `fetch(Fe,{method:`POST`,headers:{"Content-Type":`application/json`},body:JSON.stringify({jsonrpc:`2.0`,id:++Re,method:`call`,params:[Y,e,t,n]})})).json()` | Base `Fe=`/ubus`` (same-origin). Body = ubus JSON-RPC `call` with object `Y` (var, default `00000000000000000000000000000000`), method `e`, args. No auth header — ubus session token is the `Y` param / `net4sats_session` localStorage key. |
| POST | `/ubus` JSON-RPC (second call site) | admin/assets/index-Cv-gJynA.js:15070 | `fetch(Fe,{method:`POST`,headers:{"Content-Type":`application/json`},body:JSON.stringify({jsonrpc:`2.0`,id:++Re,method:`call`,params:[`00000000000000000000000000…` | Same as above (second literal call site). |
| GET | `/whoami` (same-origin) | admin/assets/balance-xS7ctLQ6.js:17778 | `fetch(Te()+`/whoami`);if(!e.ok)throw Error(`Whoami fetch failed`);` | `function Te(){return`http://${window.location.hostname||`192.168.1.1`}:${we}`}` with `var we=2121` (same :2121 base as portal). No auth. |
| GET | `/balance` (same-origin) | admin/assets/balance-xS7ctLQ6.js:18003 | `fetch(Te()+`/balance`)).json()` | Same :2121 base. No auth, no body. |
| GET | `venue.html` (relative) | admin/assets/balance-xS7ctLQ6.js:19011 | `fetch(`venue.html`).then(e=>{if(!e.ok)throw Error(String(e.status));return e.text()})` | Relative same-origin asset fetch. No auth. |
| GET | `/api/scan` | index.html:381 | `const resp = await fetch('/api/scan');` | Wizard same-origin (backend `main.go:628`). No auth. |
| POST | `/api/wifi-scan` | index.html:309 | `const resp = await fetch('/api/wifi-scan', {` | Same-origin. JSON body. No auth. |
| POST | `/api/deploy` | index.html:461 | `const resp = await fetch('/api/deploy', {` | Same-origin. JSON body. No auth. |
| GET | `/api/status/` + jobId | index.html:477 | `const resp = await fetch('/api/status/' + jobId);` | Same-origin. No auth. |
| GET | `selectedPkgURL` (GitHub release asset, laptop side) | deploy.go:402 | `data, dlErr := httpGetFile(selectedPkgURL)` | URL constants at `deploy.go:43,47,50` (`https://github.com/felixfelix-bot/...`). Laptop-side outbound. No auth header (public releases). No body. |
| GET | `baseURL + ndsPkg` (OpenWrt package dir listing + ipk) | deploy.go:469,483 | `baseURL := "https://downloads.openwrt.org/releases/24.10.4/packages/aarch64_cortex-a53/"` / `ndsData, ndsErr := httpGetFile(routingURL + ndsPkg)` | Laptop-side via `httpGetFile`. No auth, no body. |
| GET | `baseURL + jqPkg` | deploy.go:473,492 | `jqListHTML := string(httpGetFileOrEmpty(packagesURL))` | Laptop-side via `httpGetFile`. No auth. |
| GET | `http://127.0.0.1:2121/` (router-side health probe) | deploy.go:948 | `healthOut = sshRun(client, "wget -qO- http://127.0.0.1:2121/ 2>/dev/null | head -c 100 || echo 'health check failed'")` | Executed ON THE ROUTER over SSH, not by the wizard host. No auth, no body. |
| GET | `$FEED_URL` (feed reachability probe) | scripts/fe2b_feed_check.sh:348 | `code="$(curl -s -o /dev/null -w '%{http_code}' -L -m 25 "$FEED_URL" 2>/dev/null)"` | CI/verification script. No auth header. No body. |
| GET | `http://127.0.0.1:$HEALTH_PORT/`, `:PORTAL_PORT$PORTAL_PATH`, `:PORTAL_PORT/` | scripts/fe2b_feed_check.sh:329,337,343 | `h="$(rsh "curl -s -o /dev/null -w '%{http_code}' -m 8 http://127.0.0.1:$HEALTH_PORT/" | tr -d '\r' | tail -1)"` | Executed ON THE ROUTER via `rsh`. No auth, no body. |

Single outbound primitive in Go: `httpGetFile(url)` at `deploy.go:1219`
(`netClient := &http.Client{Timeout: 60 * time.Second}` / `resp, err :=
netClient.Get(url)`). All 6 Go call sites (deploy.go:402,472,473,483,492)
pass release/package URLs only. There is **no** `http.NewRequest`,
`http.Post`, `client.Post`, or `Do(req)` anywhere in non-test code.

## 4. KEYWORD-SCAN HIT TRIAGE (why semantic matches are NOT callers)

Broad scan `rg -n -E "identity|reveal[-_]seed|revealSeed|/seed|mnemonic|seed_phrase|seedPhrase" <repo>`:

| Hit | File:line | Classification |
|---|---|---|
| `identity` docs prose | AGENTS.md:16,17,24,29,30,31,46,47; README.md:139,161 | Documentation, not code. AGENTS.md:30/47 are the unbacked claim. |
| `identity` uci record ops | deploy.go:126,200,201,873 | uci `identities.json` (owner identity / profit_share **record**), not an HTTP endpoint. |
| `identity` test JSON fixture | security_hardening_test.go:294 | `{"identity":"owner"}` inside a Go string fixture. |
| `identity` board prose | docs/feed-install.md:45,424 | Docs. |
| `seed` = time-seeded | deploy.go:568,1170 | `// Fallback: time-seeded` (SSID suffix RNG). |
| `seed` = "preseeded" prose | defaults_test.go:12,32,34,45 | Form-value prose. |
| `seed` = crypto internals | portal/assets/splash-*.js | BIP32 master seed / Cashu `seed` blinding var / NUT-09 wallet-restore prose inside minified crypto libs. Not a request. |
| `reveal` (any case) | — | **Zero matches in *.go/*.js/*.html/*.sh/*.json/Makefile.** |
| `mnemonic`, `seed_phrase`, `seedPhrase`, `revealSeed` | — | **Zero matches.** |
| `/identity` path in a sweep script | verify/grepsweep-t_1e234edf/sweep.sh:47,50 | A **grep pattern list** (`for p in ... '/identity' ...`) scanning a log. Not an HTTP call. |

## 5. WHAT THE REPO'S ENDPOINTS ACTUALLY ARE (for contrast)

| Layer | Endpoint surface | Source |
|---|---|---|
| Wizard backend (Go, local) | `/api/scan`, `/api/wifi-scan`, `/api/deploy`, `/api/status/`, `/` | main.go:628–632 |
| Admin UI → router | `/ubus` JSON-RPC (ubus objects `tollgate`, `network.*`, `dhcp`, `uci`, `system`) | admin/assets/index-Cv-gJynA.js:14439; rpcd/tollgate_acl.json |
| Portal + balance → router | `:2121` `/`, `/whoami`, `/ln-invoice`, `/balance` | portal/assets/splash-Be2Pt9NX.js:15680–16887; admin/assets/balance-xS7ctLQ6.js:17778,18003 |
| rpcd method allowlist | `config_schema`, `config_get/set/save`, `config_save_identities`, `wallet_*`, `status`, `health`, `pricing`, `sessions`, `activate`, `upstream_*`, `clients` | rpcd/tollgate_acl.json |
| CI | `go build ./...`, `go vet ./...`, `go test ./... -race -count=1` — **no curl/wget/http step** | .github/workflows/ci.yml:27–31 |

Note the closest identity-adjacent surface: rpcd exposes
`config_save_identities` (a **write** of identity config via ubus), and
`deploy.go:873` writes `lightning_address` into `identities.json` over SSH.
Neither is `GET /identity` nor `POST /identity/reveal-seed`, and neither
transfers seed material.

## 6. ACCEPTANCE CHECK

- Every row in §3 carries ≥1 `file:line` citation and a verbatim snippet. YES.
- Unproven items reported as "not found in repo": `GET /identity` → verb+path
  pair absent (FALSE); `POST /identity/reveal-seed` → absent.
- Server-side route definitions excluded per instructions ("belongs to server scan").
- Read-only honoured: no source file modified; only this report + the
  accompanying `PROGRESS-t_wz17-caller-identity-scan.md` were added; no test,
  no build, no network call, no kanban mutation.

## 7. RESIDUAL UNCERTAINTY (stated, not hidden)

1. Three stale minified bundles exist on disk but are NOT referenced by any
   HTML (`portal/assets/splash-BYAdP4yt.js`, `portal/assets/splash-DkfE4jIG.js`,
   `admin/assets/index-Dz-uU1-s.js`, `admin/assets/balance-xS7ctLQ6.js` is
   referenced by nothing — `admin/balance.html` loads `balance-BNg56sOO.js`).
   I scanned the *referenced* bundles exhaustively and keyword-scanned all
   bundles; none of the three dead ones contains `identity`/`reveal*` either.
2. `*.js.map` sourcemaps exist for the bundles. Not parsed here (the emitted
   bundles are authoritative for what is *called*). If the sibling task wants
   original source paths, they are in `admin/assets/*.js.map`,
   `portal/assets/*.js.map`.
3. `deploy.go` embeds router-side shell (`sshRun`) that itself performs `wget`;
   one such line is inventoried (deploy.go:948). A full router-side command
   audit would be a separate scope (server/device scan).
