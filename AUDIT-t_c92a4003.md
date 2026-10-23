# AUDIT (task 0 substitute) — WZ1.9 de-branding leaks + test/doc gaps

Repo: /home/c08r4d0r/repos/net4sats-mvp-v2 (origin = felixfelix-bot/net4sats-wizard-go)
Branch: worker-heavy/t_gate8-19b1b1c6
Base for this audit: 520b5df

## Why this file exists

Tasks 0 (`t_a8cbb51e`, audit) and 1 (`t_b2d894bd`, baseline) are BOTH `triage` with
**zero runs** on the board (`/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`).
Their inputs were never produced, so the card (`t_c92a4003`) that depends on their findings
had nothing to read. This audit is performed here, from primary sources, so the TDD work rests
on measured evidence rather than an absent upstream.

## Baseline (task 1 substitute) — measured, not asserted

```
$ go test -count=1 ./...            # at a38acbf/520b5df era
ok  net4sats-wizard  37.400s        exit 0
```

Suite is GREEN at 520b5df. The one RED guard present earlier in the session
(`TestNoSplitWordingInVisibleCopy`, `FAIL net4sats-wizard 41.098s`, exit 1) was
closed by commit `520b5df` — a sibling worker's TDD cycle that landed in this
shared tree. Logs: /tmp/wz19-baseline.log (RED), /tmp/wz19-baseline2.log (green).

## Scope of the sweep

Every tracked, wizard-shipped surface: the `go:embed`-ed `index.html`, `portal/`,
`admin/` HTML/JS, plus `*.go`, `*.md`, `*.yml`. Minified vendored bundles
(`admin/assets/*`, `portal/assets/*`) excluded — third-party bundles, not our copy.
Reports/PROGRESS/EVIDENCE and `artifacts/` excluded (task-history, not shipped; also
git-excluded in this repo).

## (a) De-branding leaks in the user-facing wizard demo path

| # | File:line | Leak | Class | Verdict |
|---|---|---|---|---|
| L1 | `discover.go:33-34` | `// Include 192.168.21.1 (GL.iNet default WiFi-AP subnet seen on Felix's T14Gen5)` | personal laptop model + operator first name in a shipped source comment | **REAL — fix** |
| L2 | `discover.go:41` | `// On some laptops (e.g. Felix's T14Gen5) the router shows up in the ARP` | same | **REAL — fix** |
| L3 | `ssh.go:125` | `// ... must not diverge further; see worker-heavy/t_d6e354a6-lint notes.` | internal kanban branch/card name in a shipped source comment | **REAL — fix** |
| L4 | `README.md:126` | `- **[Endo onboarding runbook](docs/endo-runbook.md)** — operator-facing setup guide` | internal handover person's name, in the public README | **REAL — fix** |
| L5 | `docs/endo-runbook.md` (filename) | the file itself is named after the internal handover | internal naming on a public doc path | **REAL — fix (rename)** |
| L6 | `deploy.go:25` | `// Fallback strategy (Endo handover):` | internal process name in source comment | **REAL — fix** |
| L7 | `CHANGELOG.md:18` | `` resolved FELIX Q3 gate (`Q3=B`, upstream default stays as-is) `` | internal gate id + operator name in a public changelog | **REAL — fix** |
| L8 | `deploy_test.go:10,13,14` | `felixfelix-bot fork` / `Endo-handover fallback strategy` in comments | internal process naming; owner string is a real remote | **REAL — narrow (comment wording only)** |

### Explicitly NOT leaks (documented so the section is unambiguous)

- `discover.go:37` `"10.47.41.1"` — a **probe candidate IP**, i.e. functional behaviour.
  It is a plausible router-gateway address (a carrier/CGN-style LAN gateway). It is data the
  code dials, not branding. **Retained.** (Removing it would be a behaviour change with no
  de-branding benefit, and would need its own card.)
- `deploy.go:43,47,50` + `pins_test.go:59` `felixfelix-bot/...` release URLs — these are the
  **live artifact pins** the wizard downloads from. Renaming them would break the install.
  A pin migration is owned by the upstream-switch work, not a copy pass. **Retained**, and
  already covered by `artifact_pins_test.go` / `pins_test.go`.
- `index.html:406` `r.vendor` / `discover.go:15` `Vendor` — matched only by a careless grep for
  the substring "endo" inside "v**endo**r". Not a leak. **No change.**
- The `devSplit` wire field, `devsplit` DOM ids, and the persisted `config.json`
  `developer` identity — protocol/config identifiers, already documented as deliberately
  retained in CHANGELOG + the prior report.
- `net4sats` brand itself — the wizard IS the net4sats product; branding is intended.

## (b) Behaviour changes lacking a same-commit doc update or failing-first test

| # | Gap | Evidence | Verdict |
|---|---|---|---|
| G1 | **No test guards ANY of L1–L8.** The existing debrand guards (`debrand_copy_test.go`) cover only revenue-share wording, the `Dev split` control name, README payout wording and the deploy log line. Nothing pins the internal-process/personal-name family. | `grep -rn 'T14Gen5\|10\.47\.41\.1\|Endo\|FELIX' --include='*_test.go'` → **0 hits** (the only match was a prose comment outside any assertion) | **REAL — new failing-first test** |
| G2 | `docs/endo-runbook.md` renamed ⇒ every inbound reference must move in the same commit (README). Also the doc's own self-reference. | `README.md:126` link target | **REAL — fix in same commit as rename** |
| G3 | The `Dev split` control-name cycle (**already closed by 520b5df**) landed *without* a worktree-root PROGRESS/REPORT update for this card. | `git show --stat 520b5df` — code+test+CHANGELOG+README only | noted; this task's PROGRESS/REPORT covers it |

## Sections explicitly clean

- **Internal endpoints/hostnames:** none in the shipped path. No `*.internal`, no
  hardcoded staging host, no telemetry URL. The only hosts are GitHub release URLs (real pins).
- **Internal logos/assets:** none. Every logo under `admin/assets/logo/` and
  `portal/assets/` is the net4sats mark, which is the product's own brand.
- **Board/workspace paths:** 0 hits in the shipped path (`hermes/kanban`, `kanban.db`,
  `workspaces/t_`).
- **Personal hostnames:** 0 hits (`c08r4d0r`, `c03rad0r`, `CobradorWave`) in the shipped path.
- **Locale/i18n catalogs:** the repo has none, so no locale variants to sweep.
