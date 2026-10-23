# REPORT — t_e01d7b78

**Task:** Scan wizard repo for SERVER-SIDE identity/seed route definitions (inbound handlers).
**Mode:** READ-ONLY static inspection. No source file edited, no build, no test run, no state-mutating
command. The only writes are this report and `PROGRESS-t_e01d7b78.md`.

**Repo root:** `/home/c08r4d0r/repos/net4sats-mvp-v2`
**Remote:** `origin = https://github.com/felixfelix-bot/net4sats-wizard-go.git`
**Branch:** `worker-heavy/t_gate8-19b1b1c6`
**HEAD at scan start:** `681ed99fe0edb98bc8d36d5a7dead7cbd4db51f4`
**HEAD at scan end:** `7c8d032dad997c5a034e5f36d0bc4d7ae503a350`
(sibling commits `eb430b5`, `2edfa0b`, `7c8d032` landed mid-scan; none touched `main.go` — the route
table below was re-verified at the new HEAD with `git grep … HEAD -- main.go` and is byte-identical.)

**Language inventory (measured, not assumed):** Go (22 root `*.go`, the only server language),
POSIX shell (`rpcd/tollgate`), HTML + inline JS (`index.html`), minified vendored Preact/Cashu
bundles (`admin/assets/*.js`, `portal/assets/*.js`), one Python artifact validator with no server
(`artifacts/demo_basecut/validate_trim_plan.py`). **No** TypeScript, Rust, Python-web or Node source
defines a route here. No `configurationwizzard` checkout exists on this host.

---

## INPUT caveat (stated up front, per the brief)

The brief's INPUT ("task 0 / Locate wizard repo root") **has no readable result**:

- Board `net4sats-mvp-v2`, card `t_715c00e3` "Locate wizard repo root and record path, remote,
  HEAD SHA" — status `triage`, `result` column **NULL**. It never ran.
- The documented fallback, `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` item **WZ1.7**,
  **does not exist**: `grep -c "WZ1.7"` over that file and its `wz1.5-url-inventory*.md` siblings
  → **0**.

Root / HEAD / language inventory were therefore **re-derived first-hand** in this session, as shown
above. `rg` is not installed on this host (`rg: command not found`), so step 1's scan was run as
`grep -rEn` and `git grep -nE` equivalents.

---

## FINDINGS — every SERVER-SIDE (inbound) route definition in the repo

| Path | Verb | File:line | Verbatim snippet (≤3 lines) | Auth / params / body notes |
|---|---|---|---|---|
| `/api/scan` | **unconstrained** (no `r.Method` check; handler reads nothing from the request) | `main.go:628` (registration), handler `main.go:146` | `mux.HandleFunc("/api/scan", handleScan)` / `func handleScan(w http.ResponseWriter, r *http.Request) {` | No auth. No path params. No body. CORS middleware wraps it (`main.go:635`); ACAO is echoed **only** for `http://127.0.0.1:8099` and `http://localhost:8099` (`main.go:63-66`, `main.go:74-76`). |
| `/api/wifi-scan` | **POST-only** (405 otherwise) | `main.go:629` (registration), guard `main.go:294` | `mux.HandleFunc("/api/wifi-scan", handleWifiScan)` / `if r.Method != http.MethodPost {` / `writeError(w, 405, "method not allowed")` | No auth. No path params. JSON body = `wifiScanRequest` `{"ip","password"}` (`main.go:153-156`); `password` optional (empty root password on a fresh-reset router, `main.go:303-304`). CORS as above; `OPTIONS` short-circuits 204 (`main.go:79-82`). |
| `/api/deploy` | **unconstrained** (no `r.Method` check in the handler) | `main.go:630` (registration), handler `main.go:513` | `mux.HandleFunc("/api/deploy", handleDeploy)` / `func handleDeploy(w http.ResponseWriter, r *http.Request) {` / `if err := json.NewDecoder(r.Body).Decode(&req); err != nil {` | No auth, no path params. JSON body = `deployRequest` `{"ip","password","mode","ssid","wifiPass","lnurl","devSplit","margin","mint"}` (`main.go:501-511`). **Not identity/seed related** — listed for completeness because it is the one route that drives root SSH and writes `/etc/tollgate/identities.json` (see row below). |
| `/api/status/` (prefix; `jobID` = remainder) | **unconstrained** (no `r.Method` check) | `main.go:631` (registration), handler `main.go:543` | `mux.HandleFunc("/api/status/", handleStatus)` / `jobID := r.URL.Path[len("/api/status/"):]` | No auth. **Path parameter**: the job id is the path suffix (`/api/status/<jobID>`). No body. 404 `{"error":"job not found"}` on unknown id (`main.go:548-551`). |
| `/` (catch-all) | **unconstrained** | `main.go:632` (registration), handler `main.go:568` | `mux.HandleFunc("/", handleIndex)` / `func handleIndex(w http.ResponseWriter, r *http.Request) {` / `_, _ = w.Write(indexHTML)` | No auth. Serves the embedded `index.html` (`embed.go:8-9`). Its inline JS calls only `/api/wifi-scan`, `/api/scan`, `/api/deploy`, `/api/status/<id>` (`index.html:309`, `:381`, `:461`, `:477`). |
| **`ubus tollgate config_save_identities`** (router-side RPC, not HTTP) | JSON-RPC method, dispatched on `$2` | declaration `rpcd/tollgate:13`, dispatch `rpcd/tollgate:60`, exec `rpcd/tollgate:70`; ACL `rpcd/tollgate_acl.json:33`; deployed by `deploy.go:768-771` | `"config_save_identities":{"json":"string"},` / `config_save_identities)` / `echo "$JSON" | tollgate --json config save-identities 2>/dev/null` | **Identity-adjacent and the only identity-touching inbound surface this repo serves.** No auth inside the plugin — access is governed by the rpcd ACL file it ships (`/usr/share/rpcd/acl.d/tollgate.json`), which grants it in the **`write`** section (`rpcd/tollgate_acl.json:28-39`). Body shape: one string field `{"json": "..."}`; empty → `{"success":0,"error":"json is required"}` (`rpcd/tollgate:63-68`). |

### Route-definition idiom sweep (to prove the table is exhaustive)

- Go: `git grep -nE 'HandleFunc|Handle\(|ServeMux|http\.Handle' -- '*.go'` (excluding `_test.go`)
  → **only** `main.go:627-632` (the `http.NewServeMux()` block) and `main.go:72-73`
  (`corsMiddleware`'s `http.HandlerFunc` wrapper, not a route). `http.ListenAndServe` at
  `main.go:644`. No `http.ServeMux` anywhere else, no `.GET(` / `.POST(` idiom (Gin/Echo absent).
- Express/Koa/Fastify (`app.get(`, `router.post(`, `.route(`): **0 hits**.
- FastAPI/Flask (`@app.get`, `@app.post`, `@router.`, `add_url_rule`, `@app.route`): **0 hits**.
- Rust/Actix/Axum (`.route(`, `Router::new`, `#[get(`, `#[post(`): **0 hits**.
- OpenAPI / Swagger / `urls.py` / routes manifest / nginx / lighttpd config: **not found in repo**
  (`git ls-files | grep -iE 'openapi|swagger|nginx|lighttpd|\.conf$|routes'` → empty).
- Other listeners/servers: the only `http.ListenAndServe` is `main.go:644`. `uhttpd_net4sats`
  (UCI asset, `embed.go:23-24`) configures the **router's** uhttpd (`:8090` net4sats admin home
  `/www/net4sats`, `:8080` luci) — it ships no CGI/Lua handler and the repo tracks no `.cgi`/`.lua`
  file. `deploy.go:807-849` writes the same uhttpd UCI from the wizard over SSH.
- Vendored JS bundles: `grep -ciE 'identity|reveal.?seed|npub'` → **0** for all four
  `admin/assets/*.js` and all three `portal/assets/*.js`. Their only backend paths are `/whoami`
  and `/balance` (`admin/assets/balance-BNg56sOO.js`, `portal/assets/splash-BYAdP4yt.js`).
- The embedded wizard `index.html` has **0** case-insensitive hits for `identity|reveal|mnemonic|seed`.

---

## THE TWO EXPLICIT QUESTIONS

**Q1 — Does this repo SERVE `GET /identity`?**
→ **NOT FOUND IN REPO.**
Cite: `git grep -n -E '"/identity'` over the whole tree → **0 hits**; `git grep -n 'reveal-seed'` →
only prose at `AGENTS.md:30` and `AGENTS.md:47`. The repo's complete inbound HTTP surface is
`main.go:628-632` and contains no `/identity` path. The verbatim served paths are
`"/api/scan"`, `"/api/wifi-scan"`, `"/api/deploy"`, `"/api/status/"`, `"/"` (`main.go:628-632`).

**Q2 — Does this repo SERVE `POST /identity/reveal-seed`?**
→ **NOT FOUND IN REPO.**
Cite: same absence; plus `git grep -n -i -E 'reveal.?seed|mnemonic|bip39' HEAD -- '*.go' '*.sh'
'rpcd/*'` → **exit 1, zero hits** (i.e. the repo contains no seed/mnemonic/BIP39 code at all).

**Do the served path/verb differ from the asserted ones?** There is no competing definition to
compare against — nothing in this repo serves either path. Both routes are served by a **different
repo**, `tollgate-module-basic-go` (the router backend), which was **not** scanned for this card;
that is where `AGENTS.md:30`'s `GET /identity, POST /identity/reveal-seed` describes the server.
Note for the reader: `AGENTS.md:29-31` and `:46-47` are **prose describing the backend**, not
route definitions in this repo, and must not be read as inbound handlers here.

---

## OUT OF SCOPE — flagged for the caller scan (card `t_1f51893a`)

One outbound HTTP client exists and is **not** reported as a route:
`deploy.go:1219` `netClient := &http.Client{Timeout: 60 * time.Second}` used at `deploy.go:1220`
`resp, err := netClient.Get(url)` inside `func httpGetFile(url string)` (`deploy.go:1218`) — release
asset download. **Belongs to caller scan.** Likewise the router-side `wget` health probe
(`deploy.go:947`, `deploy.go:431`) and the `curl` probes in `scripts/fe2b_feed_check.sh` are
outbound callers. The sibling card `t_1f51893a` "Scan wizard repo for caller-side identity/seed
request sites" exists on the same board with the same `workspace_path`, status `triage`.

---

## Acceptance check

- Every row above carries ≥1 `file:line` citation **and** a verbatim snippet. ✔
- Unproven items are reported as **not found in repo** (OpenAPI/Swagger/urls.py manifest; both
  identity routes) rather than inferred. ✔
- No caller-side hit is presented as a route; the one that exists is flagged as caller-scan. ✔

## Commands actually issued (all read-only)

`git rev-parse HEAD`, `git branch --show-current`, `git remote -v`, `git status --short`,
`git ls-files`, `git ls-remote`, `git check-ignore -v`, `git grep -n -E …` (several), `grep -rEn`
(several, incl. `-ci` counts), `find`, `ls`, `ps -eo pid,etime,cmd`, `sqlite3 'file:…?mode=ro'`
against `~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`. No build, no test run, no write to any
source file, no mutating git command.

## REMAINING steps

None. All brief steps (1 broad content scan, 2 route-idiom sweep, 3 verbatim path/verb extraction,
CLASSIFY, DELIVERABLE table + the two explicit verdicts) are complete.
