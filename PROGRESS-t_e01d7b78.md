# PROGRESS — t_e01d7b78 (scan wizard repo for SERVER-SIDE identity/seed route definitions)

Session 2026-10-23 (fleet offload, worker-heavy). Finding -> Status -> Files touched.
Read-only static inspection: no source file edited, no test run, no build, no state-mutating command.

## INPUT resolution (task 0 result NOT visible)

1. Finding: The named INPUT (task 0, "Locate wizard repo root") has **no recorded result** to read.
   Board `net4sats-mvp-v2` (kanban.db, read-only): `t_715c00e3` "Locate wizard repo root and record
   path, remote, HEAD SHA" is status=`triage` with `result` = NULL — never ran.
   Status: VERIFIED (read-only sqlite3 against the board DB).
   Files touched: none.

2. Finding: The documented fallback — `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`
   item **WZ1.7** — **does not exist**. `grep -c "WZ1.7"` over that schedule.md and its sibling
   `wz1.5-url-inventory*.md` → **0 hits**. The only WZ1.x items present are WZ1.4 / WZ1.5.
   Status: VERIFIED (grep, exit 0 / count 0). Re-derived root/HEAD/languages first-hand instead.
   Files touched: none.

3. Finding: Re-derived inventory, first-hand, in this session:
   - root `/home/c08r4d0r/repos/net4sats-mvp-v2`
   - remote `origin = https://github.com/felixfelix-bot/net4sats-wizard-go.git`
   - branch `worker-heavy/t_gate8-19b1b1c6`
   - HEAD at session start `681ed99fe0edb98bc8d36d5a7dead7cbd4db51f4`;
     HEAD at end of scan `7c8d032dad997c5a034e5f36d0bc4d7ae503a350`
     (3 sibling commits landed mid-scan: `eb430b5`, `2edfa0b`, `7c8d032` — none touched main.go,
     so the route table was re-verified at the new HEAD; identical).
   - languages actually present: **Go** (22 `.go` files, root package), **POSIX shell**
     (`rpcd/tollgate`), **HTML+inline JS** (`index.html`), **minified vendored JS** (Preact/Cashu
     bundles under `admin/assets/`, `portal/assets/`), **Python** (one artifact validator,
     `artifacts/demo_basecut/validate_trim_plan.py`, no server). No TS/JS/RS/Python source that
     defines a route. No `configurationwizzard` checkout on this host.
   Status: VERIFIED.
   Files touched: PROGRESS-t_e01d7b78.md, REPORT-t_e01d7b78.md.

## Scan results

4. Finding: `rg` is **not installed** on this host (`rg: command not found`). The task's step 1
   command was run as `grep -rEn` / `git grep -nE` equivalents. The `*.go` glob adaptation was
   valid — Go is the only server language here.
   Status: VERIFIED (command error observed, workaround used).
   Files touched: none.

5. Finding: The repo defines exactly **five** inbound HTTP routes, all in one `http.ServeMux`
   at `main.go:627-632`, plus **one** router-side ubus method table (`rpcd/tollgate`, deployed by
   `deploy.go:768-771`). Full table + verbatim snippets in REPORT-t_e01d7b78.md.
   Status: VERIFIED (`git grep -n -E 'mux\.HandleFunc' HEAD -- main.go`, `HandleFunc|Handle\(|ServeMux`
   repo-wide, `\.GET\(|\.POST\(|Router::new|app\.(get|post)|add_url_rule|@app\.route` → 0 hits).
   Files touched: none.

6. Finding: **No identity or seed-reveal route is served by this repo.** `git grep -nE '"/identity'`
   → 0 hits; `git grep -n 'reveal-seed'` → only prose at `AGENTS.md:30` and `AGENTS.md:47` (plus the
   prose of the pre-existing WZ1.7 evidence note). Zero hits in `*.go`, `rpcd/*`, `uhttpd_net4sats`,
   or the vendored JS bundles. The only identity-adjacent *served* thing is the rpcd method
   `tollgate.config_save_identities` (`rpcd/tollgate:13`, dispatch `:60`, exec `:70`).
   Status: VERIFIED.
   Files touched: none.

7. Finding: Two claims to answer explicitly: repo SERVE `GET /identity` → **NOT FOUND IN REPO**;
   repo SERVE `POST /identity/reveal-seed` → **NOT FOUND IN REPO**. Both are served by the *backend*
   (`tollgate-module-basic-go`), which is a different repo — not scanned here.
   Status: VERIFIED.
   Files touched: none.

8. Finding: Sibling card `t_1f51893a` "Scan wizard repo for caller-side identity/seed request sites"
   exists on the same board, same `workspace_path`, status `triage`, assignee `worker-heavy`.
   Caller-side hits are therefore deliberately NOT reported here; the one outbound HTTP client in the
   repo (`deploy.go:1219`) is flagged as "belongs to caller scan" only.
   Status: VERIFIED (board DB read-only).
   Files touched: none.

## Deliverable location

- Full report (table + verdicts + citations): `REPORT-t_e01d7b78.md` at this worktree root.
- Both files are **gitignored** by this worktree's shared exclude file
  (`/home/c08r4d0r/wz16-verify/.git/info/exclude` lines 4-5: `/PROGRESS-t_*.md`, `/REPORT-t_*.md`),
  so they need `git add -f` to be committed. They were committed with explicit pathspecs.
- The shared root `PROGRESS.md` / `REPORT.md` were **also** prepended (additive-only, verified
  0 deletions), because the brief names those exact filenames; they are local-only in this worktree.
