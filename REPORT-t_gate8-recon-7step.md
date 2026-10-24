# REPORT — 7-step wizard recon (t_gate8 recon lane)

Task: enumerate all 7 wizard steps — headings, fields, validation, transitions.
Deliverable: `docs/recon-wizard-steps.md` (read-only recon; no spec or component
file modified).

Raw identifiers are kept in this file on purpose: `docs/*.md` is a scanned
surface for the repo's internal-detail guard (`debrand_internal_test.go:20-41`,
`shippedSurfaces()` at `:52-103`), while root `REPORT-*`/`PROGRESS-*` files are
explicitly out of scope (`debrand_internal_test.go:43-46`).

---

## 0. Handoff answer

**Deliverable path (absolute):**
`/home/c08r4d0r/repos/net4sats-mvp-v2/docs/recon-wizard-steps.md`

**Committed on:** branch `worker-heavy/t_gate8-19b1b1c6`, worktree
`/home/c08r4d0r/repos/net4sats-mvp-v2`.

---

## 1. Headline finding — the 7-step flow is NOT IMPLEMENTED

This is the single most important thing for the downstream spec author, and it
is not a matter of interpretation:

- `grep -c data-testid index.html` = **0**. No step container, no field, no
  control in the wizard document carries a testid.
- **0 of 7** `[data-testid="step-*"]` containers exist.
- **2** `<h2>` elements exist, and neither is a step heading: `net4sats is live!`
  (`index.html:225`) and `Setup failed` (`index.html:234`).
- **0** buttons matching Next / Continue / Skip. The document's buttons are
  `Rescan` (`:172`), `Deploy net4sats` (`:211`), `Scan again` (`:212`),
  `Try again` (`:236`).
- No step-component files exist at all — there is no `.tsx/.jsx/.vue/.svelte`
  anywhere in the repo; `index.html` (522 lines) is the entire UI.

`index.html` is one flat form with **5 view containers**
(`:141 #scan-view`, `:149 #select-view`, `:216 #deploy-view`,
`:223 #success-view`, `:232 #error-view`) in which the Discover / Connect /
Uplink / Payouts / Advanced controls are all merged into the single
`#select-view` card (`:149-213`).

The 7-step flow exists only as a **declared contract** in a deliberately RED
Playwright spec, whose stated purpose is "the follow-on task turns it green"
(`evidence/wz1.8-red-generic-7step/RED-RUN-generic-7step-t_gate8.md:3-4`):

- `~/repos/physical-router-test-automation/tests/browser/generic-wizard-7step.spec.mjs` (401 lines) — the contract spec
- `~/repos/physical-router-test-automation/tests/browser/wizard-flow.mjs` (94 lines) — the step model / registry (`WIZARD_STEPS`)

RED status is 5 failed / 5 total, every failure on a wizard-contract assertion,
diagnosed in the artifact itself as "a **missing-generic-wizard-behaviour**
failure, not a typo, import error, or environment failure"
(`RED-RUN-generic-7step-t_gate8.md:71-72`).

---

## 2. Per-step required-field list (compressed, for handoff)

From `wizard-flow.mjs:33-41` (order/testids) and `:60-69` (`STEP_REQUIREMENTS`).
**No required field anywhere has an existing `data-testid` — all are absent.**

| # | Step | testid (declared) | Required fields (`nextRequires`) | Existing hook today | Exists? |
|---|---|---|---|---|---|
| 1 | Discover | `step-discover` | *(none)* | `#scan-view` (`index.html:141`) | Absent as step |
| 2 | Connect | `step-connect` | `router` | `#router-select` (`index.html:152`), `#password` (`:156`, optional) | Absent as step |
| 3 | Uplink | `step-uplink` | `router`, `uplinkMode` | `#mode-wan` / `#mode-sta` divs (`:161-162`), `#ssid` (`:169`), `#wifi-pass` (`:178`) | Absent as step |
| 4 | Payouts | `step-payouts` | `lightningAddress` | `#lnurl` (`index.html:183`) | Absent as step |
| 5 | Advanced (**optional**) | `step-advanced` | *(none — must not gate)* | `<details class="advanced">` (`:187`), `#devsplit` (`:191`), `#margin` (`:197`), `#mint` (`:203`) | Absent as step; fields exist |
| 6 | Deploy | `step-deploy` | `router`, `lightningAddress` | `#deploy-btn` (`index.html:211`) | Absent as step |
| 7 | Live | `step-live` | *(none; terminal)* | `#success-view` (`index.html:223`) | Absent as step |

Validation rules that gate Next:

- Discover: none; auto-satisfies when `/api/scan` returns ≥1 router.
- Connect: a router must be selected — Next **disabled** until then, and a forced
  click on the disabled control must not advance
  (`generic-wizard-7step.spec.mjs:331-347`).
- Uplink: a mode must be chosen (`wizard-flow.mjs:64`); spec checks the WAN radio
  (`generic-wizard-7step.spec.mjs:189-192`).
- Payouts: format-valid Lightning target, two forms — email-shaped
  `/^[^\s@]+@[^\s@]+\.[^\s@]+$/` or raw LNURL
  `/^lnurl1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{6,}$/`. Client
  `validLightning()` (`index.html:419-424`), server `validLightningAddress()`
  (`main.go:23-42`). Next stays disabled for `not-a-lightning-address`.
- Advanced: **none** — Next must already be enabled on arrival.
- Deploy: carries `router` + `lightningAddress`; server also rejects empty IP
  (`main.go:521-524`) and invalid LNURL (`:525-528`), and silently clamps
  `devSplit` to 0–50 / `margin` to 0–100 (`deploy.go:886-887`,`:1313`).
  Forward control is `/deploy|start/i`, not `/next|continue/i`
  (`generic-wizard-7step.spec.mjs:372`).
- Live: terminal; no forward control may remain
  (`generic-wizard-7step.spec.mjs:395-398`).

---

## 3. Advanced optionality — the brief's explicit question

**Skip control: NOT IMPLEMENTED.** There is no skip control anywhere in the repo.

- Contract: `/skip/i` accessible name (`generic-wizard-7step.spec.mjs:286-289`);
  it is a **role+name locator, not a testid** — `wizard-flow.mjs` declares no
  testid for it.
- Today: `grep -c -i '>(next|continue|skip)' index.html` = **0**.
- The only nod to optionality is the disclosure label
  `<summary>Advanced settings (optional)</summary>` (`index.html:188`) — a
  `<details>` collapse, which neither advances the flow nor records
  skipped/configured status.
- `[data-testid="advanced-summary"]` is **absent** (`grep -c advanced-summary
  index.html` = 0), so the required `data-status="skipped"|"configured"`
  distinction (`generic-wizard-7step.spec.mjs:300-302`,`:325-327`) cannot be
  reported today.
- Does skipping change Deploy state? Per contract, **yes, observably**: Deploy
  requires only `['router','lightningAddress']` (`wizard-flow.mjs:67`), so
  skipping leaves `devSplit`/`margin`/`mint` at defaults — `devSplit: 0`,
  `margin: 0`, `mint: https://8333.space/` (`index.html:191,197,204`;
  declared `main.go:508-510`).

---

## 4. Workspace discrepancy (deviation from the brief — stated, not worked around)

The brief named worktree/repo `configurationwizzard` on branch
`feature/upstream-wifi-scan`, and said the onboarding Playwright spec was in
that worktree. Verified reality:

| Brief said | Reality |
|---|---|
| worktree `configurationwizzard` | **Does not exist.** `git worktree list` for this repo lists 14 worktrees, none by that name. `configurationwizzard` is a **kanban board** (`~/.hermes/kanban/boards/configurationwizzard/board.json:2`, `"repo": "configurationwizzard"`), not a checkout. |
| branch `feature/upstream-wifi-scan` | **Does not exist** locally or on any remote. `git branch -a --list '*wifi*' '*upstream*'` → `upstream-main`, `origin/fix/upstream-release-urls`, `origin/fix/v0.3.8-wifi-scan-null`, `origin/fix/wifi-default-upstream-mode`. |
| onboarding Playwright spec in the worktree | Not in this repo at all. It lives in a **different repo**: `~/repos/physical-router-test-automation/tests/browser/endo-onboarding.spec.mjs`. |

I did **not** switch worktrees or branches (explicitly forbidden). I worked
read-only in `/home/c08r4d0r/repos/net4sats-mvp-v2` on the already-checked-out
`worker-heavy/t_gate8-19b1b1c6`, which is what the brief's repo name resolves to.

`endo-onboarding.spec.mjs` itself is the **old 5-view** walkthrough, not a
7-step flow: it asserts on `#select-view:not(.hidden)` and `#router-select`
directly (lines 33-37, 48-57) and targets the real GL-MT6000 at `192.168.1.1`
on OpenWrt 25.12.0 (lines 15-18, 27). The 7-step contract is a **separate,
newer** spec (`generic-wizard-7step.spec.mjs`).

---

## 5. Schedule file (step 1 of the brief)

Read `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` (337 lines) first,
as instructed.

**It says nothing about the wizard 7-step flow and nothing about the GL-MT6000.**
It is entirely about an operator decision gate (Q3 = do **not** change the
upstream dev-split default; net4sats may apply its own value downstream) and a
Minibits DLEQ upstream-tracking question (Q3b = don't file; tracking already
exists upstream). Zero mentions of Discover/Connect/Uplink/Payouts/Advanced/
Deploy/Live, zero mentions of the step model, zero mentions of the target
device. No constraint from it applies to this task — recorded so the next lane
does not re-read it hoping for wizard guidance.

---

## 6. Acceptance criteria — checked

| Criterion | Status | Evidence |
|---|---|---|
| All 7 steps present, correct order, none hand-waved | **PASS** | `docs/recon-wizard-steps.md` §2 — one section per step, all 7, order per `wizard-flow.mjs:33-41`; each has heading, container, fields, validation, transition, conditional-reachability, and today's-UI row |
| Every selector/testid cited exists today, or is explicitly marked absent | **PASS** | Every cited `id=` was read from `index.html` at the cited line; every absent testid is marked **Absent** in §4 of the artifact, backed by `grep -c data-testid index.html` = 0 |
| No spec or component file modified | **PASS** | `git status --porcelain` below — only `docs/recon-wizard-steps.md` is new; all other paths pre-existed this task |

### Read-only proof

Before this task, `git status --porcelain` showed 7 untracked paths, all
pre-existing (`.hermes/`, `docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-*`,
`evidence/wz1.8-red-generic-7step/VERIFICATION-t_gate8-red.md`, and four
`verify/` dirs). After: those same 7 plus `docs/recon-wizard-steps.md` and this
report. **No tracked file was modified** — `git diff --stat` is empty for
tracked paths.

Nothing under `tests/browser/` (the spec) and nothing under the wizard UI was
touched: `index.html`, `main.go`, `deploy.go` and the Playwright specs are all
unmodified.

---

## 7. Guard check (repo's own debrand test)

`docs/*.md` is a scanned surface, so `docs/recon-wizard-steps.md` was run through
`TestNoInternalProcessLeaksOnShippedSurfaces` (`debrand_internal_test.go:108-126`):

```
$ go test -run TestNoInternalProcessLeaksOnShippedSurfaces ./... | grep -c recon-wizard-steps
0
```

**Zero hits on the new file.** The test still FAILS overall, but on three
**pre-existing** docs that are not this task's:
`docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md`,
`docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md`,
`docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md`.
That failure pre-dates this lane (it is the same guard failure documented in
`evidence/wz1.8-red-generic-7step/VERIFICATION-t_gate8-red.md` §4, isolated there
to commit `22fbea4`). This lane did not introduce or fix it.

To keep the new doc clean, raw internal identifiers were kept **out of** the
artifact and placed in this report instead (which the guard does not scan).

---

## 8. Files touched

| File | Change |
|---|---|
| `docs/recon-wizard-steps.md` | **NEW** — the deliverable (426 lines) |
| `REPORT-t_gate8-recon-7step.md` | **NEW** — this report |
| `PROGRESS.md` | appended one section (never overwritten) |

No other file created, modified, or deleted.

---

## 9. What a downstream spec author should know first

1. The 7-step flow must be **built**, not documented-against. The contract is
   complete (titles, order, testids, gating, skip semantics, Live evidence) but
   there is no implementation to test.
2. Two contract details are likely to bite: Uplink's WAN/STA are **`<div onclick>`
   cards, not `<input type="radio">`** (`index.html:161-162`), so the spec's
   `getByRole('radio')` would match nothing; and the router **password** is
   required operationally for `/api/wifi-scan` and `/api/deploy` but is **not
   placed in any step's `nextRequires`**.
3. Deploy's 11 internal sub-steps (`deploy.go:234-249`) share **names** with the
   spec's stubbed status payload but not **descriptions**
   (`generic-wizard-7step.spec.mjs:92-104`) — assert on `name`, not `desc`.
