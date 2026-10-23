# REPORT — t_gate8-19b1b1c6 / wz1.5 URL inventory (fleet offload)

Deliverable: `/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.md`
(replaces the stale 2026-10-15 edition of that path, which described a base tree the wizard branch no
longer carries and predated `tollgatePkgURLApk` entirely).

Read-only inventory: no repo source file modified, no test run. Full detail in the per-card
`REPORT-t_gate8-19b1b1c6.md` at the worktree root.

ACCEPTANCE: all six tokens accounted for — 172 raw lines at the frozen HEAD, 26 actionable (listed)
+ 146 leave-alone (registered by file); implementer needs no re-grep (file:line + exact line + class
(a/b/c/d) + concrete action per row, sequenced in 3 phases). Verdicts: `tollgateOSURL` = already
deleted, keep the tripwire and never reuse the name; `configwizURL` = genuinely unused AND
unofficial -> delete (4 edits, one commit), do not re-pin.

HEADLINE RISK: the official-channel repin is **not** a drop-in. The current pins live in
`felixfelix-bot` forks that are AHEAD of their upstreams (upstream newest `tollgate-wrt` .ipk is
`v0.5.0`; the wizard pins `v0.7.0-alpha10`; upstream `v0.6.0` publishes no asset at all). An
owner-only swap 404s on every fresh deploy — exactly the SW4a outage. `releases.tollgate.me` serves
the alpha3 index but its signing key is 404. So Phase 2 is blocked on a channel decision this
document cannot make for you; Phase 1 (the `configwizURL` deletion) can land now.

GIT: no source commit. The two task-scoped report files are committed and pushed to
`worker-heavy/t_gate8-19b1b1c6`; root `PROGRESS.md`/`REPORT.md` are worktree-local (gitignored) and
were prepended additively. Three tracked files (`deploy.go`, `defaults_test.go`, `index.html`) carry
a SIBLING worker's uncommitted mint-filter edit — inspected, disjoint from this task's lines, not
touched.


---

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


---

# REPORT — t_ba89d8b4: Encode base cut (head/tail trim only, constant 1080p+ frame size)

Session 2026-10-23 (fleet offload worker, worker-heavy). Full detail: REPORT-t_ba89d8b4.md.
Prior sections below belong to earlier tasks sharing this worktree root (preserved whole).

## VERDICT: BLOCKED — no encode performed, and none is possible from this chain's inputs.

**Deliverable path: NONE. No base-cut MP4 was produced.**

The trim plan exists and is deliberately non-executable; its handoff_to_encoder names this exact
card (t_ba89d8b4) with the instruction "DO NOT RUN", and the card's own INPUTS clause says stop
when the plan's source paths do not resolve. They do not. Stopped at STEP 1.

- Plan: status=BLOCKED, executable=false, source.path=null (DOES_NOT_EXIST), all four cut points
  null, cuts.encode_must_not_proceed=true, step_map=[] , regions=[]. Validator: 21/21 PASS.
- Upstream, live from kanban.db: t_882b3922 (record take) 3 runs / 0 completed / 0 attachments;
  t_15057569, t_6d731c04, t_be3a9ae1 all triage, 0 runs. Chain severed at the recording step.
- Source genuinely absent: /tmp/demo-video.webm ENOENT; whole-disk find by exact size
  (1353893c) and by name both return nothing; /tmp is tmpfs, so it was wiped between runs.
- Decoy rejected: artifacts/demo_basecut/demo-trimmed.mp4 in two sibling worktrees
  (sha 83bdc911...) is an already-trimmed derivative of a DIFFERENT chain (Cashu SMS wallet),
  and its take never resolves past "Paying..." to a working state. Adopting it would have
  satisfied the mechanical acceptance rows while failing the four-phase criterion.
- No ffmpeg encode was run; no trim point was invented; no interior edit was made.
  ffmpeg/ffprobe 6.1.1 present and exercised — tooling was not the blocker.
- Vision unavailable this pass (vision_analyze 503), so no frame was visually inspected; all
  claims rest on ffprobe/ffmpeg pixel statistics and the board DB.

Unblocks when: t_acee3752 -> t_882b3922 (absolute + durable take path) -> t_15057569 ->
t_6d731c04 -> t_be3a9ae1 -> re-dispatched t_ae815553 (executable plan) -> re-run of this card.

---

# REPORT — Gate-8 review & merge, wizard official-URL change (card t_19b1b1c6, re-run)

Board: net4sats-mvp-v2. Session: 2026-10-21 (fleet offload worker, worker-heavy profile).
Repo under gate: net4sats/net4sats-wizard-go (upstream) + felixfelix-bot fork.
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2 (detached HEAD 8a8d006).

## VERDICT: GATE CANNOT PROCEED — MERGE REFUSED. Merge commit SHA: NONE.

This is the second execution of this gate card. The first (2026-10-14, verdict card
t_222a7d2b, FULL_REVIEW: /home/c08r4d0r/reports/reviews/net4sats-mvp-v2-t_19b1b1c6-deepseek-v4-flash.md)
refused the merge. Every blocker it cited still holds today, re-verified live from primary
sources (not from the prior card's prose). One materially NEW fact strengthens the refusal.

## Step 1 — cold review + validation report

Cold review t_7da567aa ("Cold review the wizard official-URL change before completion"):
VERDICT BLOCKED, not APPROVE, not REQUEST-CHANGES. Read directly from the board DB
(~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db, task_comments 2026-10-11,
worker-reviewer-kimi): "nothing to review" — producer card t_7e5033df never produced output
(still status=triage, empty result), plan inputs absent. Seven findings (F1-F7), three
blockers. Card remains status=blocked; no re-dispatch since.

Validation report t_deabf3d6 ("grep gate + install-step dry-run"): status=triage, empty
result. The validation report does not exist.

UNVERIFIED-network-check decision (explicitly recorded, per step 1):
- ACCEPTED as verified (by re-measuring today, not by trusting): feed channel reality.
  FreedomTechFeed/packages v0.6.0-alpha1 asset
  tollgate-wrt_0.6.0_alpha1_aarch64_cortex-a53.ipk -> HTTP 200, 7,937,246 B, sha256
  5134f29dd64f9985f9b993ea6a06b76823af221a88a6cb319d59d8448edb3b90 == the pin in the
  candidate branch. Measured by download + sha256sum this session.
- NOT ACCEPTED: the install-step dry-run. It remains structurally impossible (no change
  exists to dry-run against) and stays an open blocking gap owned by t_deabf3d6. No
  vacuous pass recorded.

## Step 2 — pushed? docs same commit? suite green?

NOT PUSHED. Live evidence:
- Upstream net4sats/net4sats-wizard-go: no branch and no PR for the wz1.5 official-URL
  change. main = b43651279b41903c6ba76f73af1d46f7823182cf (2026-08-27) still pins
  felixfelix-bot at deploy.go:43/:47/:50 — grep 'felixfelix-bot' count = 6 on main's
  deploy.go TODAY. Open PRs: #33 (CORS), #34 (CI). #35 CLOSED.
- Fork: only ref carrying the FreedomTechFeed pin map is feat/feed-per-arch-urls@9a5f34f
  (23 commits ahead of upstream main, 0 behind, merge-tree CLEAN). Its PRs: fork #6
  ("consolidation ... -> main", OPEN, targets FORK main, not upstream) and upstream #1
  (CLOSED 2026-09-09). No upstream PR exists.
- Docs-in-same-commit: moot — there is no implementation commit at all.

Suite on the candidate revision (probe worktree, since upstream main is its ancestor
and "the merged revision" would be its tip): build OK (go build -buildvcs=false), vet OK,
`go test -count=1 ./...` = 2 FAIL:
- TestSshRunFunction — ALSO FAILS at base b436512 (pre-existing; a live lab router answers
  the empty-credential identify probe on this network). Environmental, not a regression.
- TestRenewOSCommandOnCurrentOS — test ABSENT at base, INTRODUCED by the branch, FAILS on
  this host ("expected a known renew command on host-OS, got cmd='"). Branch-introduced red.
Gate-8's "full suite green on the merged revision" is therefore false on the only ref that
carries the change. Probe worktrees removed after measurement.

## Step 3 — Gate-8 checklist and merge

Checklist outcome:
1. Cold review APPROVE? NO — BLOCKED since 2026-10-11.
2. Change pushed + docs same commit? NO — no change exists on any remote.
3. Full suite green on merged revision? NO — see Step 2.
4. Official channel correct? The only pin map in existence points at FreedomTechFeed
   v0.6.0-alpha1 — real and byte-verified, but STALE: the feed published v0.6.0 FINAL
   today (2026-10-21T02:06:45Z, 26 assets incl. per-target ipk/apk for 8 targets).
   Merging today would pin a superseded tag.
MERGE: REFUSED. No merge commit SHA exists to record — recording "none" rather than
fabricating one. A merge would have required me to author the implementation mid-gate,
which is route-back work (t_7e5033df), not gate work.

## Step 4 — Q5 (final vs interim channel)

EXPLICITLY TRACKED, not silently dropped: card t_f6c8aa8e ("Q5 follow-up: record
final-vs-interim Tollgate artifact channel decision for the wizard (WZ1.5)"), status=triage,
assignee=worker-tollgate, created 2026-10-14 by the first gate run, brief intact, and the
board confirms no interim URL is baked into any main — the risk remains latent. New evidence
for that card: v0.6.0 FINAL exists on the feed, which strengthens the FINAL-channel path;
the decision itself remains that card's deliverable.

## Route-back (unchanged from t_222a7d2b, plus one addition)

1. Implement + push the official-URL change on a branch/PR against net4sats:main, pinned at
   FreedomTechFeed v0.6.0 (NOT alpha1 — stale) — t_7e5033df.
2. Produce the validation report (grep gate + install-step dry-run) — t_deabf3d6.
3. Produce the channel decision note — t_1df48589; close Q5 decision — t_f6c8aa8e.
4. Fix the branch-introduced TestRenewOSCommandOnCurrentOS failure (host-OS renew-command
   coverage) before re-running the suite gate. NEW item this run.
5. Re-dispatch cold review with producer->review link (t_7da567aa F6) — t_7da567aa.
6. Re-run this Gate-8 card against the real head SHA.

## Acceptance criteria status (honest accounting)

- "Merged commit SHA recorded": NOT SATISFIABLE — recorded as NONE, merge refused. This is
  the gate doing its job, not a failure to execute it.
- "Q5 follow-up either closed or explicitly tracked": SATISFIED — explicitly tracked
  (t_f6c8aa8e, open, owner assigned) and re-confirmed this session.
- "Task marked complete with the gate outcome stated": gate outcome = MERGE REFUSED, stated
  above and in the commit.

secret-scan note: no credentials in this report; all checks were read-only plus one
public-asset download to /tmp.

---

# REPORT ADDENDUM — Preflight FE2b t_18d15abe (2026-10-21): STOP, FE2a not merged

FE2b preflight re-run on v0.6.0 (published 2026-10-21T01:17:07Z). Verdict: STOP — the release-tag
leg holds, the FE2a-merged leg does NOT: no feed index exists in the published artifacts (26 flat
assets; Packages/arch-dir paths all 404), the FE2a card t_df4304d5 is todo, and the feed-bridge
branch/workflow/gh-pages surfaces are all absent. The TDD failing baseline IS delivered:
scripts/fe2b_feed_check.sh (5 PASS / 3 FAIL: A2 feed-index 404, A5 local-scp origin with
/tmp/tollgate-wrt.ipk staged on the box, A8 portal 403). Device: opkg 24.10.2 aarch64_cortex-a53,
tollgate-wrt 0.6.0_rc6-r1 installed (local origin), service up, :2121 200. Serialization: a
competing worker-tollgate lane (pid 2118663) is mid FE2b feed-install — HOLD before any device work.
Full report: REPORT-t_18d15abe.md.

---

# REPORT — net4sats-mvp-v2 / WZ1.9 history atomicity (QGATES v3.2)

Worker pass, 2026-10-16. Subject repo: `felixfelix-bot/net4sats-wizard-go` (fork).
Worktree used: `/home/c08r4d0r/worktrees/wz19-atomicity`.

**Verdict: PASS — no rewrite was owed. The atomic history already exists and was already
published. I executed zero history rewrites and zero force-pushes.**

---

## 0. Card-prose corrections (the prompt named things that do not exist)

| Prompt says | Reality |
|---|---|
| repo `net4sats/tollgate-wizard-go` | **404** (`gh api repos/net4sats/tollgate-wizard-go` → Not Found). |
| worktree on branch `WZ1.9` | No `WZ1.9` branch and no such worktree anywhere in `git worktree list` or `git log --all --grep=WZ1`. |
| repo `net4sats-mvp-v2` | The local clone at that path is a **linked worktree whose remote is `felixfelix-bot/net4sats-wizard-go`**. |

Resolved subject by evidence:

- Real upstream (PR base repo): **`net4sats/net4sats-wizard-go`**, main = `b436512`.
- Fork (where the PRs live): `felixfelix-bot/net4sats-wizard-go`, main = `8fcc30e`.
- The "WZ1.9" label attaches to the lineage `origin/wz19-video-evidence`
  (its commits carry the string `WZ1.9`, e.g. `cc5daa8 docs(demo): host WZ1.9 E2E demo artifacts…`).
- `net4sats/tollgate-wizard-go` is a **different, unrelated repo** (AGENTS.md's "net4sats-wizard-go"
  naming; a "tollgate-wizard-go" surfaces in `boards/conwrt`) — it is not this subject.

## 1. Base resolution (skill Step 0) — the decisive measurement

| Base ref | sha | PR #32 behind/ahead | conflict markers |
|---|---|---|---|
| `upstream-main` (= `net4sats/net4sats-wizard-go` main) | `b436512` | **0 behind / 4 ahead** | **0** |
| `origin/main` (= fork main, `felixfelix-bot`) | `8fcc30e` | 16 behind / 8 ahead | **20** |

`git merge-base upstream-main origin/wz19-video-evidence` == `b436512` == `upstream-main` itself,
so the branch is rooted exactly at the real base and is **0 behind** it. No rebase owed.

PR #32 reads `mergeable=CONFLICTING / mergeStateStatus=DIRTY` — **this is pure fork-main base
pollution**, not a content conflict. The PR is declared against the fork's `main`
(`baseRefName: main` → `felixfelix-bot/net4sats-wizard-go@8fcc30e`), which is 12 commits ahead of
upstream with unrelated unreleased work. Fix is a **PR base retarget to `net4sats:main`**, not a rebase.
Measured tool output (`git merge-tree <mb> <base> <head> | grep -cE '<<<<<<<|CONFLICT'`):
20 vs fork main, **0** vs upstream main.

## 2. Branch is / is not the PR head (skill Step 1)

Live read, `gh pr view 32`:

```
number=32 state=OPEN headRefName=wz19-video-evidence
headRefOid=cc5daa8a22df01758965ba5a9c777df0b42c29d2 baseRefName=main mergeable=CONFLICTING
url=https://github.com/felixfelix-bot/net4sats-wizard-go/pull/32
```

`git ls-remote origin 'refs/heads/wz19*'` → `wz19-video-evidence = cc5daa8…` (matches the PR head,
byte for byte). **The branch I audited IS the PR head.** Sibling PR #33 points at
`wz19-video-evidence-pre-atomic @ 134e235` — the pre-atomic ref, kept as a sibling, not the head.

## 3. DELIVERABLE — `git log --oneline`

There is **no local ref named `main`** in this clone (`entrypoint % git remote -v` → only `origin`;
`git log main..HEAD` → `fatal: ambiguous argument`). So the literal `main..HEAD` range cannot be
produced here; both meaningful readings follow.

```
$ git log --oneline upstream-main..HEAD     # ← the real PR base (net4sats/net4sats-wizard-go main = b436512)
cc5daa8 docs(demo): host WZ1.9 E2E demo artifacts on a scratch evidence branch
75b9c88 fix(ui): restore Minibits mint option, drop stale DLEQ-workaround note
bdaf560 fix(ui): de-brand wizard — dev split defaults to 0, developer-fund copy removed
6cdbf65 fix(deploy): stop force-injecting testnut mints into router config

$ git log --oneline origin/main..HEAD       # fork main = 8fcc30e — NOT the PR base (base pollution)
cc5daa8 docs(demo): host WZ1.9 E2E demo artifacts on a scratch evidence branch
75b9c88 fix(ui): restore Minibits mint option, drop stale DLEQ-workaround note
bdaf560 fix(ui): de-brand wizard — dev split defaults to 0, developer-fund copy removed
6cdbf65 fix(deploy): stop force-injecting testnut mints into router config
b436512 chore: vendor portal v1.0.10-alpha (pricing block) + bump configwizURL (#31)
cc9f94c fix(deploy): close root shell-injection sinks; bind loopback; localhost-only CORS (#24)
8b10897 fix: deploy portal bundle to uhttpd :2051 docroot (+ /www mirror) (#30)
de8c36a fix: embed portal with balance redirect fix + bump configwizURL to v1.0.9-alpha (#23)

$ git log --oneline main..HEAD
fatal: ambiguous argument 'main..HEAD': unknown revision or path not in the working tree.
```

**Final branch name: `wz19-video-evidence`** (tip `cc5daa8a22df01758965ba5a9c777df0b42c29d2`).

### Atomic + conventional: every entry

| sha | subject | atomic? | conventional? |
|---|---|---|---|
| `cc5daa8` | `docs(demo): host WZ1.9 E2E demo artifacts on a scratch evidence branch` | yes — single concern: evidence hosting | yes `docs(demo)` |
| `75b9c88` | `fix(ui): restore Minibits mint option, drop stale DLEQ-workaround note` | yes — one concern, one fix + its doc | yes `fix(ui)` |
| `bdaf560` | `fix(ui): de-brand wizard — dev split defaults to 0, developer-fund copy removed` | yes — one concern | yes `fix(ui)` |
| `6cdbf65` | `fix(deploy): stop force-injecting testnut mints into router config` | yes — one concern | yes `fix(deploy)` |

Census (`upstream-main..HEAD`): `commits=4  merges=0  non_conventional=0  wip_words=0  behind_real_base=0`.

## 4. Docs-same-commit (skill Step 3) — PASS

Per-commit, with the doc delta itself as evidence:

| commit | files | CHANGELOG added lines |
|---|---|---|
| `6cdbf65` | `CHANGELOG.md`, `deploy.go`, `docs/endo-runbook.md`, `security_hardening_test.go` | 6 |
| `bdaf560` | `CHANGELOG.md`, `api_test.go`, `index.html`, `main.go`, `ui_defaults_test.go` | 5 |
| `75b9c88` | `CHANGELOG.md`, `deploy.go`, `index.html`, `ui_defaults_test.go` | 6 |
| `cc5daa8` | `docs/demo/README.md`, `docs/demo/demo-trimmed.mp4`, `docs/demo/wz19-demo.gif` | 0 |

- The three behaviour-change commits each carry their own CHANGELOG delta **in the same commit**. Not stragglers.
- `cc5daa8` is a pure docs/evidence commit (README + two demo artifacts) — it owes no CHANGELOG. Stated explicitly rather than left blank.
- `CHANGELOG.md` appearing in all three code commits is **not** a mixed concern; single-concern is justified from behaviour, not file overlap.

**No straightener commits exist on the head.** A mechanical scan for CHANGELOG-only commits:

```
origin/wz19-video-evidence-pre-atomic: changelog_only_commits=1   <- 9098c84 docs(changelog): …
origin/wz19-video-evidence:            changelog_only_commits=0
```

**Zero** CHANGELOG-only commits on the head. A prior pass already squashed them.

## 5. Content-preservation proof (skill Step 4)

```
atomic     tree: e64ab1b4fac3ef8a0a67bc3a37160d6e2c151079
pre-atomic tree: e64ab1b4fac3ef8a0a67bc3a37160d6e2c151079     <- identical
$ git diff --name-status origin/wz19-video-evidence-pre-atomic origin/wz19-video-evidence
   (empty)
```

Shared interior boundary is also tree-identical: `eb5a95d` (pre) tree `a2b53e5f…` == `6cdbf65` (atomic) tree `a2b53e5f…`.

Combined content change is byte-identical across the whole range:

```
diff eb5a95d..134e235  sha256: aeaad7a2429260e8bbae6a581901c818543c893a0285f9b45d5e2be6a53ed302
diff 6cdbf65..cc5daa8  sha256: aeaad7a2429260e8bbae6a581901c818543c893a0285f9b45d5e2be6a53ed302   <- equal
```

**The rewrite was shape-only.** Every prior gate verdict on this change set remains valid by
construction.

### Where the two stragglers went

The pre-atomic ref carried exactly two non-atomic artifacts:

- `9098c84 docs(changelog): drop internal board path from Minibits note` — a **docs-only commit repairing the previous commit's doc text**. A straggler by construction.
- `2d77096 chore(debrand): strip remaining "developer fund" wording from docs and code` — a **follow-up changing a CHANGELOG line an earlier commit had added**, plus one `main.go` line.

Both are absorbed into the atomic commits that own the concern (`75b9c88` and `bdaf560`
respectively), which is why 6 pre-atomic commits → 4 atomic ones.

## 6. Scope discipline — nothing was rewritten

The atomic form is **already published**. Local tip == remote tip:

```
local  HEAD            = cc5daa8a22df01758965ba5a9c777df0b42c29d2
remote wz19-video-evidence = cc5daa8a22df01758965ba5a9c777df0b42c29d2
porcelain = 0
```

Consequently:

- **No rewrite command was run. No force-push was issued. No commit was created.** A force-push here
  would move an open PR's head and invalidate SHA-anchored review — a real cost with zero benefit,
  since the content is already correct.
- I checked out the existing branch locally as `wz19-video-evidence` tracking `origin/wz19-video-evidence`
  in the isolated worktree `/home/c08r4d0r/worktrees/wz19-atomicity` — this is a **local ref creation
  only, tip already equal to the remote; the remote was not touched.**

## 7. Suite re-run at the published tip (skill Step 7)

Ran in the **isolated** worktree at `cc5daa8`, so the primary checkout was never mutated.

```
$ go build ./...                 -> exit 0
$ go vet ./...                   -> exit 0
$ go test ./... -race -count=1 -v
  ok  net4sats-wizard  101.957s
  RUN=127  PASS=127  FAIL=0  SKIP=0
  (counts include subtests: `grep -cE '^ *--- PASS'`)
```

Post-run assertion — the suite did not disturb the tree:

```
HEAD=cc5daa8a22df01758965ba5a9c777df0b42c29d2   (unmoved)
porcelain=0 lines
```

`golangci-lint v1.61.0` (CI's pinned version) is **pre-existing-red and unchanged**:

| | findings | errcheck | ineffassign | unused |
|---|---|---|---|---|
| base `b436512` | 14 | 10 | 2 | 2 |
| tip `cc5daa8` | 14 | 10 | 2 | 2 |

The only difference between the two finding sets is 4 **line-number shifts** in `deploy.go`
(`620→619`, `636→635`, `681→680`, `1078→1076`), caused by `eb5a95d`'s deletion of one line. Finding
identities are otherwise identical → **the 4 commits introduced no new lint finding.** The `lint` CI
job was already failing on `main` before these commits.

## 8. Residual operator actions (NOT executed — deliberate)

1. **Fix PR metadata, not the branch.** Retarget PR #32's base from `felixfelix-bot:main` to
   `net4sats:main`. That clears the misleading `CONFLICTING / DIRTY` flag. No rebase. No force-push.
2. **Optional:** retire the non-atomic sibling refs `wz19-video-evidence-pre-atomic` (PR #33) and/or
   close PR #33 — they are superseded by, and content-identical to, the head.
   (Deliberately left alone: they are evidence history, and deleting published refs is the operator's call.)

Not a residual action: the `lint` CI job's pre-existing red — it predates this change set and is out of scope.

## 9. Honest statement of what was NOT observed

- No kanban card matching this card-text was found in any board (searched `tasks.title`/`tasks.body`
  across `/home/c08r4d0r/.hermes/kanban.db` and all `boards/*/kanban.db` for `atomic conventional`,
  `atomic commits`, `squash`, `docs in the same commit`, `QGATES`). **No card evidence was attached,
  because I could not resolve a real card id** — a fabricated `t_xxxxxxxxx` would be a silent lie.
- I did **not** run any history rewrite. The acceptance clause "only conventional atomic commits
  between main and HEAD" is satisfied by measurement of the **already-published** branch, with the
  rewrite performed by an earlier pass.
- `origin/main` was fetched as present locally; I did not re-fetch. `git ls-remote` was used to confirm
  the live remote tips, which matched.

---

**Final state:** branch `wz19-video-evidence`, tip `cc5daa8a22df01758965ba5a9c777df0b42c29d2`,
published and unchanged by this pass; 4 atomic conventional commits vs the real base `b436512`;
suite green (127/127) at the tip; lint unchanged (pre-existing 14); PR #32's CONFLICTING flag is
fork-base pollution fixable by a base retarget.

---

# REPORT — t_deabf3d6: Validate URL change (grep gate + install-step dry-run)

Session 2026-10-23 (fleet offload worker, worker-tollgate). Full detail:
`REPORT-t_deabf3d6.md`; validation deliverable at
`~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-validation.md`.

## VERDICT: BLOCKED — there is no URL change to validate.

Both inputs named by the card are absent. The decision note
`wz1.5-channel-decision.md` is ENOENT, and no pushed revision implements the official-URL change:
`origin/main` @ `65b9d02` still pins `felixfelix-bot` at `deploy.go:43,:47,:50`, and
`git ls-remote origin` shows no `wz1.5` ref. Steps 1–4 are therefore UNVERIFIED (precondition
absent), with the surrounding measurements still taken and quoted.

Highest-value finding: the only ref carrying an official-channel attempt
(`origin/feat/feed-per-arch-urls` @ `9a5f34ff`) pins `FreedomTechFeed/packages` at tag
`v0.6.0-alpha1`, and **8/8 of those URLs 404** — that tag does not exist. The channel publishes
only `v0.6.0` (2026-10-21) and `v0.6.0-rc6` (2026-10-15), with asset names carrying an extra
target segment (`tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk`). The planned
re-pin cannot be shipped as written.

Also measured, on what actually ships: version pinning is correct
(`releases/download/<TAG>/`, zero `latest`/`channel=` hits in `*.go`), digests are content-pinned
fail-closed, and the current fork pins resolve 200 with a byte-exact sha256 match
(7,475,197 B, `786565ec…a100cd6`). No router was contacted; nothing was installed; no source
file was modified.


---

# REPORT — WZ1.4: remove the Minibits DLEQ filter, track the fix upstream

Card `t_gate8-19b1b1c6` · Repo `net4sats-mvp-v2` · Branch
`worker-heavy/t_gate8-19b1b1c6` · Remote `felixfelix-bot/net4sats-wizard-go`.
Section above belongs to the prior lane sharing this worktree; preserved.

## Outcome

| Acceptance criterion | Result |
|---|---|
| `go test ./...` PASS | **PASS** — `ok net4sats-wizard 49.209s`; `go build ./...` rc=0 |
| No Minibits-specific DLEQ filtering remaining | **MET** — filter construct absent (it was already gone); all brand-specific DLEQ copy removed; two guards prevent reintroduction |
| Upstream issue filed and URL recorded in plan + commit msg | **MET, via the existing issue** — `#1057`; recorded in `schedule.md` and in commit `docs(wz1.4)`. No duplicate filed (see below) |
| Docs updated same-commit if user-facing docs referenced the filter | **MET** — `CHANGELOG.md` entry in the code commit; no other user-facing doc referenced it |
| Atomic conventional commit | **MET** — two commits: one `fix(ui)` (code+tests+changelog), one `docs(wz1.4)` (plan/issue record) |
| PUSHED | **MET** — `c56ce38` pushed (`fea18ff..c56ce38`); second commit pushed (see §Push) |

## The central judgement call: I did not file a second issue

The brief's step 2 says file the upstream issue. **An upstream issue already exists
and covers exactly this scope**: `OpenTollGate/tollgate-module-basic-go#1057`,
opened 2026-10-23T14:45:48Z by `felixfelix-bot` — the operator's own account — and
its body explicitly states it **supersedes** the approved draft
(`docs/WZ1.4-minibits-dleq-issue-draft.md`, commit `ff2290bc`).

Filing again would have produced a duplicate of the operator's own issue, against an
explicit operator instruction that says **don't file** (Q3b = (b)). I treated the
acceptance criterion's intent — *upstream tracking exists and is recorded* — as
satisfied by recording the existing issue, and said so plainly rather than
manufacturing a new issue to satisfy a literal reading.

I verified the premise of (b) myself instead of trusting the brief, two ways:

1. **The fix is present.** `~/go/pkg/mod/.../gonuts-tollgate@v0.13.0/cashu/nuts/nut12/nut12.go:43`
   defines `VerifyProofsDLEQWithKeysets`, which groups proofs by `proof.Id` and
   fetches the matching historical keyset. Pinned by the backend at `src/go.mod:44`
   and `src/merchant/go.mod:6`. This is exactly what the draft requested, and
   ≥ the v0.10.0 threshold in the gate text.
2. **The issue is accurate, not just present.** `grep -rli dleq --include='*.go' src/`
   in `tollgate-module-basic-go` returns zero matches — independently reproducing
   #1057's central claim. So the issue tracks a real gap (the keyed verify is not
   *called*), which is nonetheless a backend concern that does not block removing
   wizard-side copy.

## What actually needed removing

Worth stating precisely, because the brief's framing ("remove the filter") does not
match the tree: **the jq filter was already gone.** It was added in `f43e020`
(2026-07-17) as
`.accepted_mints = (.accepted_mints | map(select(.url | test("minibits") | not)))`,
and already reversed — `git log -S'test("minibits")' -- deploy.go` returns nothing,
and `deploy.go`'s jq program contains only URL de-duplication.

What survived was the workaround's *narrative*, which is the real de-branding
liability:

- `index.html:207` — operator-visible hint: "Minibits temporarily removed (DLEQ
  keyset rotation bug)." Contradicts the tree (Minibits is a listed production mint
  at `deploy.go:153`) and names one mint's bug in a front-end meant to carry no
  per-mint special-casing.
- `deploy.go:882` — comment narrating the exception ("Does NOT strip minibits …")
  instead of stating the general rule.

Both removed. `defaults_test.go` gained the two guards described below.

## Guards, and evidence they actually test the thing

- `TestNoMintIsFilteredFromConfigCommand` — rejects the filter *construct* (a regex
  `test()` against a mint URL) in the jq program, and asserts every production mint
  — Minibits included — reaches the emitted `config.json` command, while testnut
  mints do not leak into a default deploy.
- `TestNoMinibitsDleqCopyInVisibleUI` — rejects brand-specific mint copy in
  operator-visible UI.

**RED/GREEN, observed.** Against a scratch copy with the old filter and old hint
re-injected, both fail:
`config jq program still filters mints by URL pattern ".url | test("` and
`operator-visible UI still carries brand-specific mint copy "minibits"/"dleq"`.
Against the fixed tree, both pass.

**A wrong test was written first and discarded.** My initial version banned the bare
words `select(` and `minibits`, which produced false failures: the jq program
legitimately uses `select(...)` for `profit_share` and URL de-duplication, and
Minibits is a legitimate listed mint whose name correctly appears in
`productionMints`. Narrowed to the filtering construct. Recorded because the naive
version would have "passed" only by accident or forced a wrong fix.

## Toolchain correction (affects earlier lanes' reports)

Prior reports on this repo state the Go toolchain is absent. It is not — it is at
`/usr/local/go/bin/go` (go1.22.12) and merely not on the default PATH. Both build and
test gates ran successfully after
`export PATH=$PATH:/usr/local/go/bin`.

## Incident: I overwrote a gitignored worktree-root `PROGRESS.md`

**Disclosed rather than quietly repaired.** My first write targeted
`/home/c08r4d0r/repos/net4sats-mvp-v2/PROGRESS.md` with `write_file`, which
overwrites. That file is gitignored (via
`/home/c08r4d0r/wz16-verify/.git/info/exclude`), so it is not recoverable from git
objects; a `git fsck` dangling-blob scan found nothing matching it.

It held the prior lane's progress text for this same worktree root. I restored it
from `PROGRESS-t_gate8-19b1b1c6.md` (4,699 B) — that lane's own task-scoped sibling,
very likely the same text — but I **cannot prove byte identity**, and I did not
record the original size before overwriting. The file now holds the restored
content followed by this card's section.

The repo's own convention (task-scoped `PROGRESS-t_<task>.md`, root files shared
across cards) exists precisely to prevent this, and I should have followed it on the
first write instead of after discovering the collision.

## Commits

| SHA | Subject | Contents |
|---|---|---|
| `c56ce38` | `fix(ui): drop the Minibits DLEQ workaround copy from the wizard` | `index.html`, `deploy.go`, `defaults_test.go`, `CHANGELOG.md` |
| _(final)_ | `docs(wz1.4): record upstream issue #1057 and close out the DLEQ workaround` | `PROGRESS.md`, `PROGRESS-t_gate8-*.md`, `REPORT.md`, `REPORT-t_gate8-*.md` |

The first commit's body carries the #1057 URL, per the acceptance criterion.

## Push

`git push origin HEAD` → `fea18ff..c56ce38  HEAD -> worker-heavy/t_gate8-19b1b1c6`
(observed). The docs commit is pushed as the final step of this run; the exact SHA
and ref update are recorded in the session summary, and the tree is left with no
uncommitted source changes.

## Remaining

1. **Upstream `#1057`** stays open, owned by the backend lane: which path verifies
   DLEQ on receive (CDK vs gonuts), and whether it routes through the keyed verify.
   Nothing here blocks on it.
2. **`t_cfb0a06a`** ("set dev-split default to 0") is REJECTED per Q3 = B — it needs
   closing/re-scoping on the board. Not a code change in this repo.
3. **`t_58a11d0d`** (WZ1.4 parent) unblocks for its other scope (de-brand business
   logic), not a default change.
4. **Board bookkeeping for this card** — mark complete with the two SHAs.
