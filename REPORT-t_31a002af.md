# REPORT — t_31a002af: Adjudicate identity API claims and emit AGENTS.md evidence table

**Date:** 2026-10-21
**Card:** `t_31a002af` — "Adjudicate identity API claims and emit AGENTS.md evidence table" (board `net4sats-mvp-v2`, status `blocked`)
**Mode:** static inspection only — no code, docs, config, or AGENTS.md edited; no tests run.
**Full report file (this file):** `REPORT-t_31a002af.md` at the worktree root.

---

## HEADER BLOCK (required)

| field | value |
|---|---|
| Repo root inspected | `/home/c08r4d0r/repos/net4sats-mvp-v2` — a **linked worktree**; `.git` file → `gitdir: /home/c08r4d0r/wz16-verify/.git/worktrees/net4sats-mvp-v2` |
| Remote | `origin https://github.com/felixfelix-bot/net4sats-wizard-go.git` (fetch + push; sole remote) |
| HEAD SHA at inspection | `d5d2373682c3777b94356001a5987cd7a0cafc24` (detached; subject `docs(fe2b): t_18d15abe preflight — STOP, FE2a not merged; TDD baseline captured`) |
| HEAD SHA at session start | `e5fedf025091ee82e4f8843b7014e03c0bdec67a`. HEAD moved to `d5d2373` mid-session — another card's commit landed in this shared tree. `AGENTS.md` is **byte-identical** at both: blob `f2ff1775a57b434ff966e7c8eaf6287e30be9fc5` at `e5fedf0:AGENTS.md`, `d5d2373:AGENTS.md` and `HEAD:AGENTS.md`. |
| AGENTS.md path examined | `AGENTS.md` at repo root (65 lines, 2976 B), blob `f2ff177` |
| AGENTS.md revision | Introduced in `475a9ef docs: add AGENTS.md — architecture: thin UI wrapper over tollgate-module-basic-go` (2026-07-06). `git log --oneline -- AGENTS.md` returns **that commit only** — the file has never been modified since. |
| Server-side repo inspected | `/home/c08r4d0r/repos/tollgate-module-basic-go`, HEAD `b11db6af658ce4fd896185d2918bf14c330cf339`, remotes `origin`/`fork` = `felixfelix-bot/tollgate-module-basic-go`, `upstream` = `OpenTollGate/tollgate-module-basic-go` |

### Note on the dispatch premises and on task 1/task 2 inputs

- The brief's `Repo: net4sats-mvp-v2` is the **board** name. The workspace path does hold a real checkout, but of `net4sats-wizard-go` — that is the repo under audit, and all caller-side citations below are against it.
- **Task 1 and task 2 produced no results to consume.** On the board, `t_1f51893a` (caller-side scan) and `t_e01d7b78` (server-side scan) both have `result` = NULL (queried read-only via `sqlite3 -readonly` on `boards/net4sats-mvp-v2/kanban.db`). Their bodies were read as the method spec only. **Every citation below was re-derived first-hand** by opening the file at that line in this session.

---

## ADJUDICATION TABLE

Claim (a) is the AGENTS.md `tollgate-module-basic-go` section's route inventory — a claim about **what is served**, so it is satisfied by SERVER evidence.
Claim (b) is the "Wizard and admin UI call …" sentence — a claim about **code calls**, so it is satisfied only by a CALLER site.

| # | Claim | Verdict | Evidence (path:line) | Verbatim snippet (≤3 lines) |
|---|---|---|---|---|
| 1 | **(a)** "`API: GET /identity, POST /identity/reveal-seed`" — backend serves `GET /identity` | **TRUE** | `tollgate-module-basic-go/src/main.go:1935` (registration); handler `src/main.go:2121`; route doc `src/main.go:2120` | `mux.HandleFunc("/identity", func(w http.ResponseWriter, r *http.Request) {` — doc line: `// merchant private key. Registered at GET /identity.` |
| 2 | **(a)** — backend serves `POST /identity/reveal-seed`, and the verb is POST | **TRUE** | `tollgate-module-basic-go/src/main.go:1941` (registration); POST-only guard `src/main.go:2144-2146`; route doc `src/main.go:2137` | `mux.HandleFunc("/identity/reveal-seed", func(w http.ResponseWriter, r *http.Request) {` and `if r.Method != http.MethodPost {` / `http.Error(w, "method not allowed: use POST", http.StatusMethodNotAllowed)` |
| 3 | **(a)** — the two paths are exactly `/identity` and `/identity/reveal-seed`: no `/api/v1` prefix, no path parameters | **TRUE** | `src/main.go:1935`, `src/main.go:1941`, served on `:2121` at `src/main.go:426`; both registered on the same mux as `/whoami` (`src/main.go:1865`) | `apiListenAddr = ":2121"` — and `mux.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {` sits on the same mux as the identity routes (`registerAPIHandlers`, called at `src/main.go:677`) |
| 4 | **(b)** — "Wizard … call `GET /identity`" | **FALSE — not found in repo** | Absence: `net4sats-mvp-v2/main.go:619-623` (the wizard's entire HTTP surface); only outbound call is a release-asset download, `deploy.go:1128-1129`; repo-wide `git grep -n "/identity" -- '*.go'` → **zero hits** | `mux.HandleFunc("/api/scan", handleScan)` / `mux.HandleFunc("/api/wifi-scan", handleWifiScan)` / `mux.HandleFunc("/api/deploy", handleDeploy)` — and `func httpGetFile(url string) ([]byte, error) {` / `resp, err := netClient.Get(url)` |
| 5 | **(b)** — "Wizard … call `POST /identity/reveal-seed`" | **FALSE — not found in repo** | `git grep -n "reveal-seed" -- .` returns **only** `AGENTS.md:30` and `AGENTS.md:47` — no code hit anywhere in the tree | `AGENTS.md:30:- API: GET /identity, POST /identity/reveal-seed` (the sole non-prose match is the claim itself) |
| 6 | **(b)** — "**admin UI** call[s] `GET /identity` and `POST /identity/reveal-seed`" | **FALSE for the admin UI shipped by this repo — not found in repo** | Vendored admin dist: `admin/assets/balance-BNg56sOO.js:1`, `admin/assets/balance-xS7ctLQ6.js:1`, `admin/assets/index-Cv-gJynA.js`, `admin/assets/index-Dz-uU1-s.js` — zero case-insensitive matches for `identity`/`reveal-seed`/`npub`. Their only backend calls are `/whoami` and `/balance` against `http://<host>:2121` | `var we=2121;function Te(){return` + `` `http://${window.location.hostname||`192.168.1.1`}:${we}` `` and `` async function Ee(){let e=await fetch(Te()+`/whoami`); `` |
| 7 | **(b)** — "… to get the results. Zero computation in the UI layer." (the *purpose* clause) | **FALSE — no such call site exists in this repo to serve that purpose** | `deploy.go:126-127` (identities.json written over SSH, not via API); `rpcd/tollgate:60` (router-side ubus method is `config_save_identities`, an exec of the `tollgate` CLI) | `// lnJQProgram updates the owner identity's lightning_address in` / `// /etc/tollgate/identities.json ($la supplied via --rawfile).` and `config_save_identities)` → `echo "$JSON" \| tollgate --json config save-identities 2>/dev/null` |
| 8 | Related `AGENTS.md:16` — wizard "Calls the Go backend's REST API on the router for identity, status, etc." (identity half) | **FALSE — not found in repo** | `deploy.go:857` is the sole interaction with the backend's HTTP surface, and it is a **wget executed on the router over SSH**, not a call from wizard code | `healthOut = sshRun(client, "wget -qO- http://127.0.0.1:2121/ 2>/dev/null \| head -c 100 \|\| echo 'health check failed'")` |

**Answers to step 2's "state the true path/verb":** for the *server*, the asserted strings are **already exact** — `/identity` on GET and `/identity/reveal-seed` on POST, on port **2121**, with no `/api/v1` prefix and no path parameters. There is no `/identity/seed` alias and no GET reveal route anywhere in the tree. The falsehood is **not** a wrong path; it is the **caller attribution**: no wizard or vendored-admin code in this repo requests either route.

---

## Step 3 — what materially distinguishes the real routes from the asserted ones

**1. `/identity` — the verb is documented but NOT enforced.**
`handleIdentityDerive` contains no `r.Method` check at all (`src/main.go:2121-2132`). The only `r.Method` tests in the identity block are at `src/main.go:2144` (reveal-seed). So a `POST` or `DELETE` to `/identity` is served identically to `GET`. The AGENTS.md wording is correct as a description of the handler's intent but overstates enforcement.

**2. `/identity/reveal-seed` is a mnemonic *recovery* endpoint, not a seed-reveal.** This is the most misleading part of the name:
- **Request body is the raw 12-word BIP39 mnemonic as plain TEXT** — not JSON. `body, err := io.ReadAll(io.LimitReader(r.Body, 1024))` then `mnemonic := strings.TrimSpace(string(body))` (`src/main.go:2148`, `:2153`).
- **Response is the full identity derived FROM that mnemonic.** `full, err := identity.DeriveFromMnemonic(mnemonic)` (`src/main.go:2154`; function at `src/identity/identity.go:208`). The server does not hand out its own stored seed.
- **Auth is loopback-only, not header-based.** `if !isLocalRequest(r) { http.Error(w, "forbidden: loopback only", http.StatusForbidden)` (`src/main.go:2140-2141`; `isLocalRequest` = `ip.IsLoopback()` at `src/main.go:2020`). There is **no** `Authorization`, `X-API-Key`, or cookie check on either route; `"Content-Type, Authorization"` at `src/main.go:1044` is CORS preflight advertisement only.
- **Body cap 1024 bytes** (`src/main.go:2148`); non-POST → 405 (`:2144`).
- **Response shape** `FullIdentity` carries secrets: `mnemonic`, `privatekey`, `root_password`, `wifi_password` (`src/identity/identity.go:75-81`).

**3. `/identity` response shape** is `DerivedIdentity{npub, ipv4, macs}` only — explicitly "safe to return from GET /identity" (`src/identity/identity.go:64-70`). No secrets.

**4. Both routes are conditionally registered.** If `identities.json` has no usable merchant key, neither route exists (404), and the process logs a warning and boots normally: `mainLogger.Warn("identity: no merchant private key in identities.json — /identity routes disabled")` (`src/main.go:1929`; gate at `:1924-1934`). Any report of these routes must state the precondition.

---

## Reconciliation with the prior run recorded on this card

The card's `result` column already holds an adjudication written from **another node** (`/home/c03rad0r/repos/net4sats-wizard-go`, HEAD `f840385`, branch `wz14-debrand`). It largely agrees with rows 1-5 above, but it differs on the admin half, and the difference is a real tree difference, not a disagreement about code:

| point | prior run (node c03rad0r, `f840385`) | this run (`d5d2373`) |
|---|---|---|
| admin UI calls `/identity` | TRUE, cited `configurationwizzard/src/components/identity-panel.tsx:38` | **not found in repo** — no `configurationwizzard` checkout exists on this host (`find /home/c08r4d0r -maxdepth 6 -name identity-panel.tsx` → nothing), this repo tracks **0** `.tsx` files (`git ls-files` = 104 files, 0 matching `\.tsx?$`), and every vendored admin bundle has **0** identity matches |
| backend line cites | `main.go:1683`, `:1846-1848`, `:1671-1694` | `main.go:1935`, `:1941`, `:1924-1934` — different backend commit; same routes |

Per the constraint, I do not import a citation I cannot re-open on this tree. I record the admin claim as **not found in repo** rather than borrowing the other node's TRUE. If the author wants the admin half asserted, it must be cited against a `configurationwizzard` checkout that is not present here.

---

## Adjacent defects observed in the same AGENTS.md block (not part of the two claims, cited anyway)

| AGENTS.md line | assertion | observed | evidence |
|---|---|---|---|
| `:41` | `SHA256(npub + "tollgate-lan-ipv4") → 10.X.Y.1/24` | Derivation is **HKDF-SHA256** with salt `tollgate-v1` and domain separator **`tollgate-ipv4-v1:`**, over the **hex pubkey**; the network is **`100.64.0.0/10` (CGNAT)**, not `10.0.0.0/8` | `src/identity/identity.go:111` `h := deriveHash("tollgate-ipv4-v1:", pubKeyHex)`; `:260` `okm, _ := hkdf.Key(sha256.New, []byte(keyHex), []byte("tollgate-v1"), domainSep, 32)`; `:114` `return fmt.Sprintf("%d.%d.%d.%d", CGNATPrefix, octet2, octet3, 1)` |
| `:42` | `SHA256(npub + "tollgate-lan-mac")` | domain separator is **`tollgate-mac-v1:<iface>:`**, keyed on the hex pubkey, per-interface | `src/identity/identity.go:125` `h := deriveHash("tollgate-mac-v1:"+iface+":", pubKeyHex)` |
| `:43` | `BIP39 mnemonic (**24 words**) → seed phrase (PrivateKeyToMnemonic)` | the mnemonic is **12 words** (128-bit entropy), and `PrivateKeyToMnemonic` **does not exist in the backend** — the implemented direction is `MnemonicToPrivateKey` (NIP-06) and `GenerateMnemonic` | `src/identity/identity.go:133` `// GenerateMnemonic generates a 12-word BIP39 mnemonic (128 bits of entropy).`; `:137` `entropy, err := bip39.NewEntropy(128)`; `grep -rn "PrivateKeyToMnemonic" --include=*.go tollgate-module-basic-go/` → **no hits** |
| `:22` | "Served by uhttpd on the router (port 8090 for admin, **port 80 for captive portal**)" | the shipped uhttpd config listens on **8090** (net4sats) and **8080** (luci). No `listen_http` on port 80 in this file. | `uhttpd_net4sats:2` `list listen_http '0.0.0.0:8090'`; `:12` `list listen_http '0.0.0.0:8080'` |

---

## Candidate statements for AGENTS.md (Part A)

One line each, present tense, each with its citation.

1. `tollgate-module-basic-go` serves `GET /identity` on port 2121, returning the public derived identity (`npub`, `ipv4`, `macs`) as JSON — `tollgate-module-basic-go/src/main.go:1935`, handler `src/main.go:2121`, shape `src/identity/identity.go:66-70`.
2. `tollgate-module-basic-go` serves `POST /identity/reveal-seed` on port 2121; the request body is a **12-word BIP39 mnemonic as plain text** (1 KB cap) and the response is the identity derived from it — `src/main.go:1941`, `:2148`, `:2153-2154`.
3. `POST /identity/reveal-seed` is POST-only (other verbs → 405) and loopback-only (remote callers → 403), with no auth headers — `src/main.go:2144-2146`, `:2140-2141`, `:2020`.
4. Neither identity route enforces the GET verb: `handleIdentityDerive` performs no method check — `src/main.go:2121-2132`.
5. Both identity routes are registered only when `/etc/tollgate/identities.json` holds a usable merchant private key, and are absent (404) otherwise — `src/main.go:1924-1934`, `:1929`.
6. The identity routes carry no `/api/v1` prefix and no path parameters; they sit on the same mux as `/whoami`, `/ln-invoice`, `/balance`, `/usage` and `/session-state` — `src/main.go:1859-1891`, `apiListenAddr = ":2121"` at `src/main.go:426`.
7. The onboarding wizard makes **no** identity-API calls; its entire HTTP surface is its own local UI server — `net4sats-mvp-v2/main.go:619-623`.
8. The wizard's only outbound HTTP request downloads release assets on the laptop — `deploy.go:1128-1129`.
9. The wizard's only contact with the router's backend surface is a health check executed over SSH on the router — `deploy.go:857`.
10. The wizard writes the merchant's Lightning address by editing `identities.json` over SSH with `jq`, not through the API — `deploy.go:134-140`, `:126-127`.
11. The admin dashboard vendored into this repo calls only `/whoami` and `/balance` on `http://<router>:2121`; it contains no identity or seed-reveal caller — `admin/assets/balance-BNg56sOO.js:1`.
12. The router-side ubus surface exposes identity configuration as `config_save_identities`, which shells out to `tollgate --json config save-identities` — `rpcd/tollgate:60`, `:70`.
13. Identity attributes are derived with HKDF-SHA256 (salt `tollgate-v1`, per-attribute domain separator) over the hex public key — `src/identity/identity.go:260`, `:111`, `:125`.
14. The derived IPv4 address lives in the CGNAT range `100.64.0.0/10`, not `10.0.0.0/8` — `src/identity/identity.go:58`, `:114`.
15. The recovery mnemonic is 12 words (128-bit entropy), and the backend implements `MnemonicToPrivateKey` (NIP-06), not `PrivateKeyToMnemonic` — `src/identity/identity.go:133-137`, `:150-156`.

---

## Delivery

| item | observed |
|---|---|
| Local commit in the shared worktree | `841538e4050d59b6d5f1e9da9d995dba45300c54` — `docs(wz1.7): adjudicate AGENTS.md identity-API claims — server TRUE, wizard caller FALSE` (2 files, +160) |
| Why a separate branch | That commit's parent (`d5d2373`) is **not** on `origin/main` — another card's unpushed commit. Publishing from there would carry foreign commits, so the report was cherry-picked onto `origin/main` in a throwaway worktree (`/tmp/wt-t31a`). |
| Shipped branch | `worker-tollgate/t_31a002af`, cut from `origin/main` `8fcc30e`, single commit `42043e579132d1868b214d54f57392700190e797` |
| Push observed | `git ls-remote origin refs/heads/worker-tollgate/t_31a002af` → `42043e579132d1868b214d54f57392700190e797` — **identical to local HEAD**. Not a self-report. |
| Credential gate | `cred-gate: clean (30 literals checked, mode=staged)` on the local commit; pass-with-no-PR at push |
| No PR opened | Correct for this card — it asks for a report, not a merge. PR URL offered by the remote is unused. |

## Method / reproducibility

Read-only throughout. Commands used, all against the two checkouts: `git rev-parse HEAD`, `git remote -v`, `git log -- AGENTS.md`, `git grep -c -i -E "identity|reveal-seed"`, `git grep -n "reveal-seed"`, `grep -rn` over `*.go`/`*.js`/`*.html`, `read_file` at explicit offsets, and `sqlite3 -readonly` against the board DB. No build, no test run, no write to any source file.

Bundle-level absence claims (`admin/`, `portal/`) were measured, not assumed: per-file counts of case-insensitive `identity|reveal-seed|npub` are `0` for all seven shipped JS bundles and for `index.html`; the only `identity` strings in the tree are `AGENTS.md`, `README.md:125,144`, `deploy.go:126,160,161,785` (all `profit_share` JSON), `security_hardening_test.go:294` (same JSON), and three `portal/assets/splash-*.js.map` source-map hits that are noble-curves doc comments ("point identity", "Size-1 FFT is the identity") — cryptographically unrelated.
