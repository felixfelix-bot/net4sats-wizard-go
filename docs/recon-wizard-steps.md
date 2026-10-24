# Recon — 7-step wizard contract (Discover → Connect → Uplink → Payouts → Advanced → Deploy → Live)

Task: enumerate all 7 wizard steps (headings, fields, validation, transitions).
Scope: READ + RECORD only — no spec and no component file is modified.
Target device: GL-MT6000.

---

## 0. BOTTOM LINE FIRST (read this before the per-step tables)

**The 7-step wizard flow is NOT IMPLEMENTED in this repository.** Every one of the
7 steps below is a **declared contract**, sourced from the Playwright contract
spec and its flow-definition data module — not from a rendered UI. There is no
step container, no step registry array consumed by the app, no per-step
component, no `Next`/`Continue` control and no `Skip` control in the wizard
document. The spec that declares the 7 steps is a **deliberately RED**,
red-first artifact whose stated purpose is "the follow-on task turns it green"
(`evidence/wz1.8-red-generic-7step/RED-RUN-generic-7step-t_gate8.md:3-4`).

Consequence for the downstream spec author: the per-step tables below give you
the **intended contract** (titles, order, optionality, gating, evidence items)
with exact citations, plus, for each step, the **closest existing UI today** in
`index.html` so you can see what must be built rather than re-reading the
codebase. Where the brief asks for "the current CSS selector or existing
`data-testid`", the honest answer for **every** field is: **no `data-testid`
exists anywhere in the wizard document** — `grep -c data-testid index.html` = **0**
(`index.html`, whole file; corroborated by
`evidence/wz1.8-red-generic-7step/served-dom-audit.txt:2-3`). Existing hooks are
plain `id=` attributes on a single flat form.

### 0.1 Workspace discrepancy (reported, not worked around)

The dispatch brief named workspace `configurationwizzard` on branch
`feature/upstream-wifi-scan`. **Neither exists on this host:**

| Brief said | Reality on disk (verified) |
|---|---|
| worktree `configurationwizzard` | No such directory. `git worktree list` for this repo lists 14 worktrees, none named `configurationwizzard`. `configurationwizzard` is the name of a **task board**, not a checkout — see §6 for why the board path is not reproduced here. |
| branch `feature/upstream-wifi-scan` | No local or remote branch matches. `git branch -a --list '*wifi*' '*upstream*'` returns `upstream-main`, `origin/fix/upstream-release-urls`, `origin/fix/v0.3.8-wifi-scan-null`, `origin/fix/wifi-default-upstream-mode`. |

I did **not** switch worktrees or branches. All work was done read-only in the
repo root the brief's repo name resolves to, on the branch that was already
checked out there (both are named verbatim in the task report at the worktree
root, not here — see §6). The brief also stated the operator-onboarding
walkthrough spec is "in the worktree"; it is not in this repo — see §1.

---

## 1. Where the flow is declared (the artifacts I traced)

Order of tracing, exactly as the brief asked (spec → outward to the step model):

| # | Artifact | Role |
|---|---|---|
| 1 | `~/repos/physical-router-test-automation/tests/browser/` — the operator-walkthrough spec (177 lines) | The "existing onboarding Playwright spec" named in the brief. **Not in this repo.** Walks the *old* 5-view form, not 7 steps. |
| 2 | `.../tests/browser/generic-wizard-7step.spec.mjs` (401 lines) | The **7-step contract spec** — the authoritative enumeration of the 7 steps. |
| 3 | `.../tests/browser/wizard-flow.mjs` (94 lines) | The **step model / registry**: `WIZARD_STEPS` array, titles, testids, optionality, per-step requirements. This is the "step registry" in the brief's sense. |
| 4 | `.../tests/browser/generic-wizard-7step.config.mjs` (37 lines) | Runner config (`retries: 0`, `workers: 1`). |
| 5 | `index.html` (522 lines, repo root) | The **actual wizard document served on :8099**. Flat 5-view form. |
| 6 | `main.go`, `deploy.go` | Server-side validation and the deploy job's 11 sub-steps. |

**There is no step-container component, no step-registry module consumed by the
application, and no per-step component files.** The only step array that exists
in the repo is the *deploy job's* 11 internal sub-steps
(`deploy.go:234-249`, `deploySteps()`), which are **not** the 7 wizard steps —
they are what runs *inside* Deploy and are what the Deploy view lists.

### 1.1 The declared step model, verbatim

`wizard-flow.mjs:33-41`:

```
id          title       optional   testId
discover    Discover    false      step-discover
connect     Connect     false      step-connect
uplink      Uplink      false      step-uplink
payouts     Payouts     false      step-payouts
advanced    Advanced    true       step-advanced
deploy      Deploy      false      step-deploy
live        Live        false      step-live
```

`STEP_COUNT = 7` (`wizard-flow.mjs:44`), frozen via `Object.freeze`
(`wizard-flow.mjs:34-42`), so order and count have a single source of truth.

### 1.2 The testid selector convention

`wizard-flow.mjs:82-84` — `stepSelector(step)` returns
`[data-testid="<step.testId>"]`. This is the **only** CSS the spec uses
(`generic-wizard-7step.spec.mjs:20-22`), plus `[data-testid^="step-"]` for DOM-order
comparison (`generic-wizard-7step.spec.mjs:282`).

---

## 2. Per-step contract

Legend: **Exists?** = does the UI described exist in `index.html` today.
Every "Absent" means *absent from the repo*, not merely from the served copy.

### Global facts that apply to all 7 steps

- **Headings**: each step must render an `<h2>` whose trimmed text equals the
  step title, in declared order —
  `expect(headings).toEqual([...STEP_TITLES])`
  (`generic-wizard-7step.spec.mjs:288-290`). Today the document has exactly
  **2** `<h2>` elements, `net4sats is live!` (`index.html:225`) and
  `Setup failed` (`index.html:234`) — neither is a step heading.
  The only other heading is `<h3>Deploying net4sats...</h3>` (`index.html:217`)
  and the brand `<h1>` (`index.html:136`).
- **Forward control**: a single button matching `getByRole('button', { name:
  /next|continue/i })` (`generic-wizard-7step.spec.mjs:134-136`). **Absent** —
  the document's buttons are `Deploy net4sats` (`index.html:211`), `Scan again`
  (`index.html:212`), `Rescan` (`index.html:172`), `Try again` (`index.html:236`).
  Today's advance is implicit: the text field's `oninput` re-runs `checkReady()`
  which enables `#deploy-btn` (`index.html:426-437`).
- **No i18n**: the document is English-only, hardcoded literals; there is no
  i18n layer, no `t()`/`$t()` calls and no locale files anywhere in the repo
  (`index.html` is plain HTML with inline strings; no `i18n`/`locale` module
  exists in the tree). **Every heading/step-indicator string below is a literal
  with no i18n key.**
- **Exactly one step visible at a time** is the spec's core invariant —
  `currentStepId()` throws unless exactly one declared step container is visible
  (`generic-wizard-7step.spec.mjs:148-158`).

---

### Step 1 — Discover

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Discover` (literal, no i18n key) | `wizard-flow.mjs:35`; asserted `generic-wizard-7step.spec.mjs:288-290` |
| Step container | `[data-testid="step-discover"]` | `wizard-flow.mjs:35,82-84` |
| Required fields | **none** — `nextRequires: []` | `wizard-flow.mjs:62` |
| Validation gating Next | Resolves on its own once the scan returns; the spec satisfies it by waiting for the first `<option>` to be visible, then advances | `generic-wizard-7step.spec.mjs:178-181` |
| Server-side check | `GET /api/scan` → `{routers:[{ip,vendor,model}]}`; the stub used is `{ip:'192.0.2.1', vendor:'TestVendor', model:'FixtureRouter'}` | `main.go:146` (`handleScan`); fixture `generic-wizard-7step.spec.mjs:42-46,60-66` |
| Transition on success | → `connect` | `wizard-flow.mjs:36`; order assertion `generic-wizard-7step.spec.mjs:237-245` |
| Conditional reachability | **No** — always the first step, unconditionally | `wizard-flow.mjs:35`; `TERMINAL_STEP`/order model |
| **Exists in UI today?** | **Absent as a step.** The equivalent is `#scan-view` (`index.html:141-146`) with spinner + literal text `Scanning for routers on your network...` (`index.html:144`), driven by `scan()` (`index.html:376-396`). No testid, no `<h2>`. | `index.html:141-146,376-396` |

Discover's failure path today rewrites the spinner's `<p>` rather than advancing:
`No routers found. Make sure your router is powered on and connected.`
(`index.html:386-388`).

---

### Step 2 — Connect

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Connect` (literal, no i18n key) | `wizard-flow.mjs:36` |
| Step container | `[data-testid="step-connect"]` | `wizard-flow.mjs:36` |
| Required fields | **`router`** (logical id) — `nextRequires: ['router']` | `wizard-flow.mjs:63` |
| Field selector used by spec | `getByRole('combobox', { name: /router/i })`, then `selectOption({index: 0})` | `generic-wizard-7step.spec.mjs:184-187` |
| Validation gating Next | Router must be selected. Asserted explicitly: Next is **disabled** until a router is chosen, and a forced click on the disabled control must **not** advance | `generic-wizard-7step.spec.mjs:331-347` |
| Format / regex | none declared | `wizard-flow.mjs:63` |
| Async / server check | Discovery must have returned at least one router (Discover is a precondition) | `wizard-flow.mjs:62-63` |
| Transition on success | → `uplink` | `wizard-flow.mjs:37`; asserted `generic-wizard-7step.spec.mjs:351-354` |
| Conditional reachability | **No** — always second | `wizard-flow.mjs:33-41` |
| **Exists in UI today?** | **Absent as a step.** Equivalent controls, all in the single `#select-view` card: `<select id="router-select">` (`index.html:152`) with `<label>Select your router</label>` (`index.html:151`); password `<input type="password" id="password">` (`index.html:156`) labelled `Router password` (`index.html:155`). Options are populated `ip — vendor model` (`index.html:406-410`). | `index.html:149-157,398-413` |

**Note on the password field**: the contract's Connect step requires only
`router` (`wizard-flow.mjs:63`). The existing form additionally carries a router
password, which is **optional** — a fresh-reset OpenWrt router ships with an
empty root password (`main.go:519-520`; UI placeholder
`Router password (leave blank on fresh reset)`, `index.html:156`). The spec never
fills it. A downstream spec author should decide whether password belongs to
Connect or Uplink, because the declared model does not place it.

---

### Step 3 — Uplink

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Uplink` (literal, no i18n key) | `wizard-flow.mjs:37` |
| Step container | `[data-testid="step-uplink"]` | `wizard-flow.mjs:37` |
| Required fields | **`router`, `uplinkMode`** — `nextRequires: ['router','uplinkMode']` | `wizard-flow.mjs:64` |
| Field selector used by spec | `getByRole('radio', { name: /wan|ethernet/i })` then `.check()` | `generic-wizard-7step.spec.mjs:189-192` |
| Validation gating Next | A mode must be chosen; the spec satisfies it by checking the WAN radio | `generic-wizard-7step.spec.mjs:189-192`; requirement `wizard-flow.mjs:64` |
| Additional (not in contract) | In the existing form, STA mode additionally requires an SSID or Deploy stays disabled | `index.html:431-435` |
| Async / server check | SSID list via `POST /api/wifi-scan {ip,password}` → `{ssids:[{name,encryption}], debug?}`; stub returns `[{name:'FixtureUplink', encryption:'wpa2'}]` | `main.go:293` (`handleWifiScan`); stub `generic-wizard-7step.spec.mjs:68-74` |
| Transition on success | → `payouts` | `wizard-flow.mjs:38` |
| Conditional reachability | **Step is unconditional**; but its *sub-fields* are conditional in the existing UI. The declared model says `uplinkMode` is always required — so the mode **choice itself** is mandatory even though WAN is the default. | `wizard-flow.mjs:64`; existing UI `index.html:250-259` |
| **Exists in UI today?** | **Absent as a step.** Equivalent controls: `Upstream connection` label (`index.html:159`) over two `div.radio-card` elements — `#mode-wan` text `WAN (Ethernet)` (`index.html:161`) and `#mode-sta` text `WiFi Repeater` (`index.html:162`) — toggled by `selectMode()` (`index.html:250-259`). **These are `<div onclick>`, not `<input type="radio">`**, so `getByRole('radio')` would match nothing: the a11y role is absent. STA panel `#sta-fields` (`index.html:165-180`) holds `<select id="ssid">` (`index.html:169`), `#rescan-btn` (`index.html:172`), `<input type="password" id="wifi-pass">` (`index.html:178`). | `index.html:158-180,250-259` |

**Nested `data-encryption`**: SSID options carry `opt.dataset.encryption`
(`index.html:347`), purely to build the hint text (`index.html:354-365`).

---

### Step 4 — Payouts

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Payouts` (literal, no i18n key) | `wizard-flow.mjs:38` |
| Step container | `[data-testid="step-payouts"]` | `wizard-flow.mjs:38` |
| Required fields | **`lightningAddress`** — `nextRequires: ['lightningAddress']` | `wizard-flow.mjs:65` |
| Field selector used by spec | `getByRole('textbox', { name: /lightning/i })` | `generic-wizard-7step.spec.mjs:196-199,378` |
| Validation gating Next | Must be a **valid** Lightning address, not merely non-empty: Next stays disabled for `not-a-lightning-address` and enables only for `operator@wallet.example` | `generic-wizard-7step.spec.mjs:374-385` |
| Format / regex (client) | Two accepted forms: email-shaped `/^[^\s@]+@[^\s@]+\.[^\s@]+$/` **or** raw LNURL `/^lnurl1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{6,}$/` — `validLightning()` | `index.html:419-424` |
| Format / regex (server) | Same two patterns, identical: `lnAddrRe` (`main.go:24`) and `lnurlRe` (`main.go:28`), combined in `validLightningAddress()` (`main.go:39-42`) | `main.go:23-42` |
| Server-side check | `POST /api/deploy` rejects with `400 "a valid Lightning address is required"` | `main.go:525-528` |
| Async check | none (pattern only; "actual resolution happens at payout time on the router" — `main.go:36-38`) | `main.go:31-38` |
| Transition on success | → `advanced` | `wizard-flow.mjs:39` |
| Conditional reachability | **No** — always fourth, and always required (`optional: false`) | `wizard-flow.mjs:38` |
| **Exists in UI today?** | **Absent as a step.** Equivalent field: `<label>Lightning address <span class="req">*</span></label>` (`index.html:182`) + `<input type="text" id="lnurl" placeholder="you@wallet.app  or  lnurl1..." oninput="checkReady()">` (`index.html:183`), hint `Lightning address (you@wallet.app) or raw LNURL (lnurl1...). Required.` (`index.html:184`). The `*` is the only "required" marker in the document; class `.req` is styled orange (`index.html:35`). | `index.html:181-185` |

**Test fixture value**: `FIXTURE_LNURL = 'operator@wallet.example'`
(`generic-wizard-7step.spec.mjs:49`).

---

### Step 5 — Advanced  *(the only optional step — see §3)*

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Advanced` (literal, no i18n key) | `wizard-flow.mjs:39` |
| Step container | `[data-testid="step-advanced"]` | `wizard-flow.mjs:39` |
| Required fields | **none** — `nextRequires: []` (deliberate: optionality must not gate) | `wizard-flow.mjs:66` |
| Field with a declared selector | Payout share slider: `getByRole('slider', { name: /payout share/i })`, value reported in `[data-testid="devsplit-value"]` as `10%` | `generic-wizard-7step.spec.mjs:314-316` |
| Validation gating Next | **Next must be ENABLED without touching anything** — an optional step cannot gate forward navigation | `generic-wizard-7step.spec.mjs:311-312`; requirement `wizard-flow.mjs:66` |
| Transition on success | → `deploy` (by Next **or** by Skip) | `generic-wizard-7step.spec.mjs:318-323,296-299` |
| Conditional reachability | **Always reachable**, but skippable. `optional: true` is the flag; `OPTIONAL_STEP` resolves to it and the spec throws if no optional step is declared | `wizard-flow.mjs:39,51-52`; `generic-wizard-7step.spec.mjs:273-274` |
| Status reporting | `[data-testid="advanced-summary"]` must carry `data-status="skipped"` or `"configured"`, and in the configured case contain the applied value (`10%`) | `generic-wizard-7step.spec.mjs:300-302,325-327` |
| **Exists in UI today?** | **Absent as a step, but the fields exist** as a collapsed `<details class="advanced">` (`index.html:187-209`) with `<summary>Advanced settings (optional)</summary>` (`index.html:188`). Fields: `Payout share` range `#devsplit` min 0 max 50 value 0 (`index.html:190-192`) with readout `#devsplit-val` default `0%`; `Margin` range `#margin` min 0 max 100 value 0 (`index.html:196-198`) with readout `#margin-val` `0%`; `Cashu mint preference` `<select id="mint">` (`index.html:202-206`) with options `8333.space (default)` and `LNbits Legend`. | `index.html:187-209` |

Existing hints (useful for the spec author's copy): dev-split is "Optional share
of payments routed to the developer payout address (config.json
`profit_share`). Set to 0 to keep the full share. Default 0%."
(`index.html:193`); margin is "Operator markup added on top of data costs.
Default 0%." (`index.html:199`); mint is "Which Cashu mint the router should
prefer for ecash payments." (`index.html:207`).

---

### Step 6 — Deploy

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Deploy` (literal, no i18n key) | `wizard-flow.mjs:40` |
| Step container | `[data-testid="step-deploy"]` | `wizard-flow.mjs:40` |
| Required fields | **`router`, `lightningAddress`** — `nextRequires: ['router','lightningAddress']` | `wizard-flow.mjs:67` |
| Forward control | **NOT** the generic Next button — a distinct control matching `getByRole('button', { name: /deploy|start/i })` | `generic-wizard-7step.spec.mjs:370-372` |
| Validation gating | The two carried requirements must still hold; the deploy button then triggers the job | `wizard-flow.mjs:67` |
| Server-side checks | `400 "IP required"` when `ip==""` (`main.go:521-524`); `400 "a valid Lightning address is required"` (`main.go:525-528`); `400 "invalid JSON"` (`main.go:515-518`) | `main.go:513-528` |
| Request body | `{ip, password, mode, ssid, wifiPass, lnurl, devSplit, margin, mint}` | `deployRequest` `main.go:501-511`; built `index.html:449-459` |
| Server clamps | `devSplit = clamp(req.DevSplit, 0, 50)`, `margin = clamp(req.Margin, 0, 100)` — silent, no error | `deploy.go:886-887`; `clamp` `deploy.go:1313` |
| Transition on success | → `live` (terminal). Trigger = the Live container becoming visible, i.e. **UI state, not elapsed time** | `generic-wizard-7step.spec.mjs:374-381` |
| Job polling | `POST /api/deploy` → `{job_id}`, then `GET /api/status/<job_id>` polled every 1000 ms | `main.go:530-541`, `index.html:461-473`; `main.go:543` (`handleStatus`) |
| Conditional reachability | **No** — always sixth | `wizard-flow.mjs:40` |
| **Exists in UI today?** | **Absent as a step.** Equivalent: `<button class="btn" id="deploy-btn" onclick="startDeploy()" disabled>Deploy net4sats</button>` (`index.html:211`), and the progress view `#deploy-view` (`index.html:216-220`) titled `Deploying net4sats...` as an `<h3>` (`index.html:217`) listing `#steps-list` (`index.html:218`) plus `#deploy-log` (`index.html:219`). | `index.html:211,216-220` |

The **11 sub-steps rendered inside Deploy** (from `deploySteps()`,
`deploy.go:234-249`, in order) are: `Verifying SSH access to router...`,
`Checking firmware version...`, `Setting root password...`,
`Configuring upstream connection...`,
`Installing net4sats package + patching backend...`,
`Branding captive portal as net4sats...`,
`Deploying net4sats captive portal...`,
`Installing net4sats admin panel...`,
`Configuring Lightning address...`, `Restarting services...`,
`Running health check...`. Note the spec's stubbed status payload uses
*different* descriptions for the same 11 names
(`generic-wizard-7step.spec.mjs:92-104`) — names match, copy does not. A
downstream spec should assert on `name`, not `desc`.

---

### Step 7 — Live (terminal)

| Aspect | Contract | Citation |
|---|---|---|
| Heading / indicator | `Live` (literal, no i18n key) | `wizard-flow.mjs:41` |
| Step container | `[data-testid="step-live"]` | `wizard-flow.mjs:41` |
| Required fields | **none** — `nextRequires: []` | `wizard-flow.mjs:68` |
| Post-deploy evidence (mandatory) | **Three** testids must each be visible: `routerReachable`, `tollgateApiServing`, `captivePortalServing` | `wizard-flow.mjs:71-75`; asserted `generic-wizard-7step.spec.mjs:386-392` |
| Forward control | **Must be absent** — `getByRole('button', {name:/next|continue/i})` must have count 0 | `generic-wizard-7step.spec.mjs:395-398` |
| Transition | None — terminal. `TERMINAL_STEP` is the last element and `optional` must be `false` | `wizard-flow.mjs:57`; `generic-wizard-7step.spec.mjs:384` |
| Conditional reachability | **No** — always the terminal step | `wizard-flow.mjs:41,57` |
| **Exists in UI today?** | **Absent as a step.** Equivalent: `#success-view` (`index.html:223-229`) with `<h2>net4sats is live!</h2>` (`index.html:225`) and two `<p>`s. **No evidence testids, no reachability/API/portal indicators.** | `index.html:223-229` |

Health evidence does exist server-side but is **not surfaced as UI**: step 10
`health` reports `API healthy on :2121` (`deploy.go:972`) and services restart at
`deploy.go:935`. Those strings land in the deploy log/step detail, not in a Live
step.

---

## 3. Advanced step — optionality, in full detail

This is the step the brief singles out, so it gets its own section.

**How it is marked optional in the contract**

1. **In data**: `optional: true` on the Advanced entry —
   `wizard-flow.mjs:39`. It is the **only** entry with `true`; all six others are
   `false` (`wizard-flow.mjs:35-41`). `OPTIONAL_STEP = WIZARD_STEPS.find(s => s.optional)`
   (`wizard-flow.mjs:51-52`).
2. **By the absence of requirements**: `advanced: []` in `STEP_REQUIREMENTS`
   (`wizard-flow.mjs:66`) — so its Next control must be enabled on arrival.
3. **By the header comment** stating the invariant: "Advanced (index 4) is the
   ONLY optional step and is skippable" (`wizard-flow.mjs:17-18`).

**The exact skip control**

- **Button label**: matched by `getByRole('button', { name: /skip/i })` — i.e. any
  button whose accessible name contains "Skip" (`generic-wizard-7step.spec.mjs:286`).
  The spec's assertion message calls it "an explicit Skip control"
  (`generic-wizard-7step.spec.mjs:287-289`).
- **Selector / testid**: the spec uses a **role+name** locator, **not** a
  `data-testid`. There is **no declared testid for the Skip control** in
  `wizard-flow.mjs` — the only testids in the contract are the 7 step
  containers plus `advanced-summary`, `devsplit-value` and the 3 `LIVE_EVIDENCE`
  items.

**→ ANSWER TO THE BRIEF'S QUESTION: the skip control is NOT IMPLEMENTED.**
`grep -c -i '>(next|continue|skip)' index.html` = **0**; the document has no
button whose text is Next, Continue or Skip (`index.html`, buttons at
`:172, :211, :212, :236`). The existing UI's only nod to optionality is the
`<summary>Advanced settings (optional)</summary>` disclosure label
(`index.html:188`) — a `<details>` collapse, which is **not** a skip control: it
neither advances the flow nor records a skipped/configured status.

**Does skipping change the state passed to Deploy?** Yes, per contract —
`deploy` requires only `['router','lightningAddress']` (`wizard-flow.mjs:67`);
advanced values (`devSplit`, `margin`, `mint`) are **not** in Deploy's
`nextRequires`, so skipping leaves them at defaults. The difference is
observable, not silent:

| Path | `[data-testid="advanced-summary"]` | `data-status` | Advanced values |
|---|---|---|---|
| Skip | present | `skipped` | untouched — not applied | `generic-wizard-7step.spec.mjs:300-302` |
| Visited | present | `configured` | applied, summary contains `10%` | `generic-wizard-7step.spec.mjs:325-327` |

On the wire, skipping means the defaults are sent: `devSplit: 0`, `margin: 0`,
`mint: https://8333.space/` (the first `<option>`, `index.html:204`; defaults
declared `index.html:191,197` and `main.go:508-510`), and the server clamps
regardless (`deploy.go:886-887`). **Today there is no summary element at all**,
so a skipped-vs-configured distinction cannot be reported.

---

## 4. Steps not implemented / stubs — explicit list

Not guessing: the following are **absent** from the repo today.

| Item | Status | Evidence |
|---|---|---|
| Any `data-testid` in the wizard document | **Absent (0 occurrences)** | `grep -c data-testid index.html` = 0; `served-dom-audit.txt:2-3` |
| All 7 `[data-testid="step-*"]` containers | **Absent (0 of 7)** | same grep; spec RED assertions `RED-RUN-...md:102-113` |
| Step container / per-step component | **Absent** — no component files exist (no `.tsx/.jsx/.vue/.svelte` in repo; `index.html` is the whole UI) | repo tree; `index.html` |
| Step registry consumed by the app | **Absent** — the only step array is the deploy job's 11 sub-steps | `deploy.go:234-249` |
| Per-step `<h2>` headings | **Absent** — 2 `<h2>`s exist, neither a step heading | `index.html:225,234` |
| `Next` / `Continue` control | **Absent** | `index.html` (buttons `:172,211,212,236`) |
| `Skip` control | **Absent** | same |
| `advanced-summary` element / `data-status` | **Absent** | `grep -c advanced-summary index.html` = 0 |
| `devsplit-value` testid | **Absent** as a testid — the readout is `id="devsplit-val"` | `index.html:190` |
| `routerReachable` / `tollgateApiServing` / `captivePortalServing` | **Absent** | `index.html:223-229` |
| Per-step forward gating | **Absent** — one global `checkReady()` gates one global Deploy button | `index.html:426-437` |
| i18n | **Absent** — no keys, no locale files | repo tree |
| WAN/STA as `<input type="radio">` | **Absent** — they are `<div onclick>` cards, so `getByRole('radio')` matches nothing | `index.html:161-162` |

**Nothing here is a "stub" or placeholder in the sense of an empty component** —
the 7-step flow simply has not been built. The RED run's own diagnosis:
"the pinned build implements no step model … one flat `#select-view` form
(Discover/Connect/Uplink/Payouts/Advanced controls merged into a single card)"
(`RED-RUN-generic-7step-t_gate8.md:58-61`), and the 5 failures are
"a missing-generic-wizard-behaviour failure, not a typo, import error, or
environment failure" (`RED-RUN-generic-7step-t_gate8.md:71-72`).

---

## 5. Open questions for the downstream spec author

1. **Password placement.** Connect requires only `router` (`wizard-flow.mjs:63`),
   yet the router root password is needed for both `POST /api/wifi-scan`
   (`index.html:312`) and `POST /api/deploy` (`main.go:503`). Where does it live
   in the 7-step model? Not declared.
2. **`uplinkMode` semantics.** `uplink` requires `uplinkMode` (`wizard-flow.mjs:64`)
   but the existing UI defaults to WAN and only forces a choice implicitly
   (`index.html:246`). Does Uplink require an *explicit* selection (Next disabled
   until a mode is touched), or is the WAN default sufficient? The spec satisfies
   it by checking the WAN radio (`generic-wizard-7step.spec.mjs:189-192`).
3. **STA SSID/password are not in the requirement model.** If `uplinkMode=sta`,
   the existing UI cannot deploy without an SSID (`index.html:431-435`), but
   `STEP_REQUIREMENTS.uplink` lists only `['router','uplinkMode']`. Should the
   model become conditional?
4. **Deploy uses a different forward control.** `/deploy|start/i`
   (`generic-wizard-7step.spec.mjs:372`) instead of `/next|continue/i`. Confirm the
   Deploy step intentionally has no generic Next.
5. **Advanced testid fixture.** Should the Skip control get a testid? Today it is
   role+name only, so any copy change breaks the spec.
6. **Deploy sub-step copy divergence.** `deploySteps()` (`deploy.go:234-249`) vs the
   spec's stubbed descriptions (`generic-wizard-7step.spec.mjs:92-104`) differ.
   Assert on `name`, not `desc`.
7. **GL-MT6000 is never referenced in the 7-step contract.** The only device
   identity in the flow is the stubbed `{ip:'192.0.2.1', vendor:'TestVendor',
   model:'FixtureRouter'}` (`generic-wizard-7step.spec.mjs:42-46`); the endo spec
   is the one that targets the real GL-MT6000 at `192.168.1.1` on OpenWrt 25.12.0
   (the operator-walkthrough spec, lines 15-18 and 27). If the downstream spec must be
   device-specific, that linkage does not exist yet.
8. **Schedule file says nothing about the wizard flow.** As instructed I read
   `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` first; it concerns an
   operator decision gate on the upstream dev-split default plus a Minibits DLEQ
   upstream issue (both resolved; the answer was to leave the upstream default
   alone and not to file). It contains **zero** mentions of the 7 steps, the wizard step
   model, or the target router model. No constraint from it applies to this task.
   (The exact plan path and the gate id are recorded in the task report at the
   worktree root, not here — see §6.)

---

## 6. Verification of the "read-only" constraint

`git status --porcelain` in the worktree shows this artifact plus pre-existing
untracked paths that were **not** created by this task — a `.hermes/` directory,
one pre-existing `docs/WZ1.8-*review-verdict-*.md` file, one pre-existing
`evidence/wz1.8-red-generic-7step/VERIFICATION-*.md` file, and several `verify/`
subdirectories. No spec file and no component file was modified; the exact
before/after porcelain capture is in the task report at the worktree root.

(Those pre-existing paths are left untouched and uncommitted by this task. Their
filenames are deliberately not reproduced here, because `docs/*.md` is a scanned
surface for the repo's internal-detail guard — see
`debrand_internal_test.go:108-126` and `shippedSurfaces()` at `:52-103`.)
