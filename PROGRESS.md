# PROGRESS — WZ1.8 env pin: wizard E2E environment (upstream commit, env vars, fixtures)

Card: fleet offload worker, repo net4sats-mvp-v2, branch worker-heavy/t_gate8-19b1b1c6
Host: x280 / CobradorWave (c08r4d0r), Ubuntu 24.04.3 LTS
Deliverable: written record only. NO spec edits, no spec writes, no hardware.


Card: fleet offload worker, repo net4sats-mvp-v2
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Host: CobradorWave / x280 (c08r4d0r), Ubuntu 24.04.3 LTS
Deliverable: written record only. NO spec edits.

---
# PROGRESS — t_0841f821: Announce GL-MT6000 lock on shared channel (fd4/fe2b)

Card: fleet offload worker, repo net4sats-mvp-v2
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Host: x280 / CobradorWave (c08r4d0r), Ubuntu 24.04.3 LTS
Mode: single announcement. No edit, no re-claim, no lock write, zero device contact.

## A1 — Channel identification: NO shared lock channel exists (step 1 as written is unsatisfiable)

Finding: the documented lock convention is a lock FILE, not a channel
(`/home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock`), and a prior dedicated
fleet audit concluded **FAIL** on the existence of a shared lock channel: the public Nostr board
mirror has no lock column, and the kind-9 fleet bus on `wss://relay.orangesync.tech` carries tags
`fleet-task/fleet-claim/fleet-running/fleet-done/fleet-reply-*` — **no lock tag**. The card meant
to create one is this card. Re-verified live this session: relay returns **HTTP 502**; `docs/`
carries no lock-channel artifact.
Status: VERIFIED (recon-quote.txt:107-119 read directly, not summarised).
Files touched: none.

## A2 — Substituted the medium the convention names as the working one

Finding: `recon-quote.txt:76-77` and `:231-232` state verbatim — "Sibling-readable announcement
medium that actually works: this repo's `docs/`." Remote `felixfelix-bot/net4sats-wizard-go` is
**PUBLIC** (`gh repo view --json visibility` → PUBLIC), so a committed `docs/` note is
anonymously readable by any peer lane — the substantive requirement ("readable by fd4 and fe2b").
Posted one announcement there rather than faking a channel that does not exist.
Status: DONE.
Files touched: `docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md` (new).

## A3 — Posted + pushed + publicly verified

Finding: commit `ab333063dcde01aa9288bf760b855c1fc0b8f78c` pushed to
`origin/worker-heavy/t_gate8-19b1b1c6` (`git ls-remote` confirms the ref equals that SHA).
Permalink fetch verified **anonymously** (no auth):
`https://github.com/felixfelix-bot/net4sats-wizard-go/blob/ab333063dcde01aa9288bf760b855c1fc0b8f78c/docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md`
→ **http=200**; raw.githubusercontent.com same SHA → **http=200**.
Status: VERIFIED.
Files touched: `PROGRESS.md`, `REPORT.md`, the announce note.

## A4 — All five required fields confirmed present in the posted message

Finding: literal counts in the posted file — `t_be3261d3` ×8, `GL-MT6000` ×6,
`glinet,gl-mt6000` ×3, `2026-10-12T07:33:17Z` ×3, `2026-10-14T07:33:17Z` ×3, `t_557bb77d` ×4.
Field 5 carries the release condition verbatim at line 61: "released only after run task
`t_557bb77d` child 3 (Playwright + video run that consumes the device) reports completion, or on
explicit abort." `secret-scan: clean`.
Status: VERIFIED.
Files touched: (record only).

## A5 — Two honest qualifiers carried into the announcement

Finding: (a) the device **serial is a documented absence** (SSH never achieved from any lane);
recon says the model string stands in for it — reported as `GL-MT6000`, no substitution invented.
(b) The claim timestamp + expected release come from the peer-node claim t_d0c3e172 on behalf of
t_be3261d3; the sibling x280 card t_be3261d3 recorded both as **null** because no claim was made
from x280. Both facts are stated side by side in the note so no reader is misled.
Status: VERIFIED (blockers.json + t_d0c3e172 comments 29831/29832 read first-hand).
Files touched: none.

---

## F1 — Upstream commit/branch pinned

Finding: HEAD = `b56fae5e3f0eb1850153b43792825b1b8f5184b3`, branch
`worker-heavy/t_gate8-19b1b1c6`, remote `felixfelix-bot/net4sats-wizard-go`.
Status: DONE
Files touched: (record only)

## F2 — The plan ref has NO WZ1.8 entry (same precondition gap prior lanes hit)

Finding: `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` is a FELIX-Q3-gate
document; `grep WZ1.8` -> 0 hits. Consistent with `t_75852fce` / `t_a8b5fcae` /
`t_cc76c746`. Authoritative WZ1.8 text lives in card bodies, not the plan file.
Status: DONE
Files touched: (record only)

## F3 — The 7-step spec is NOT in this repo; it lives in the test kit

Finding: the Playwright suite lives in `~/repos/physical-router-test-automation`
(HEAD `5344e354`, branch `fork-main-sync`). The 7-step spec
`tests/browser/generic-wizard-7step.spec.mjs` exists ONLY on
`origin/feat/router-vm-lab` (tip `ee289bf`), NOT on main / fork-main-sync.
The brief's "existing endo-onboarding spec" = `tests/browser/endo-onboarding.spec.mjs`
on the checked-out branch.
Status: DONE
Files touched: (record only)

## F4 — Server start command resolved + verified booting

Finding: wizard binds `127.0.0.1:8099` (main.go:48 `defaultListenAddr`), override
via `WIZARD_BIND`. Build via `/usr/local/go/bin/go` (go1.22.12) — `go` is NOT on
the default PATH of non-login shells. Boot verified: LISTEN on 127.0.0.1:8099,
HTTP 200, 19982 bytes.
Status: VERIFIED
Files touched: (record only)

## F5 — BLOCKER: no browser runtime on this host

Finding: `node` v18.19.1 present at /usr/bin/node, but **npx and npm are absent**;
no system chrome/chromium at any of the 5 paths the kit's launcher probes; no
`~/.cache/ms-playwright` browser bundle. A reusable playwright 1.59.1 exists at
`/data/worktrees/prta/node_modules`. Playwright configs pin `channel: 'chrome'`.
Consequence: the endo spec CANNOT be executed as-is on this host.
Status: BLOCKER (execution), recorded with evidence
Files touched: (record only)

## F6 — DOM contract: endo spec selectors EXIST in HEAD

Finding: all 12 ids the endo spec uses (`#select-view`, `#router-select`,
`#password`, `#mode-wan`, `#sta-fields`, `#lnurl`, `#deploy-btn`, `#deploy-view`,
`#success-view`, `#error-view`, `#error-detail`, `#steps-list`) are present in
HEAD's served index.html. `details.advanced` also present.
Status: VERIFIED
Files touched: (record only)

## F7 — CORRECTED: the selector contract is COMPLETE in HEAD; one value drift

Finding: my earlier F7 note was WRONG and is corrected here (I mis-grepped).
Re-verified by reading the served HTML directly:
- `#scan-view` EXISTS at line 141 — the 7-step spec's `discoverAndWait()` gate
  (`waitForSelector('#scan-view')` then `#select-view:not(.hidden)`) is satisfied.
- The wizard IS a scan -> select -> deploy -> success flow. "Discover" = the
  scan view. Steps 1-4 (Discover/Connect/Uplink/Payouts) all map onto
  select-view controls; step 6 = deploy-view; step 7 = post-deploy router check.
- So the 7-step spec is NOT structurally incompatible with HEAD, as I first said.

Real drift found instead — Step 5 asserts a value HEAD does not have:
  spec:  await expect(page.locator('#devsplit')).toHaveValue('10');
  HEAD:  <input type="range" id="devsplit" min="0" max="50" value="0"  (line 191)
  HEAD:  <span id="devsplit-val">0%</span>                              (line 190)
This is exactly the WZ1.4/Q3 dev-split question (schedule.md: net4sats fork
main.go:508 DevSplit default 10 -> now 0 at this commit; slider 0-50). The spec
was written against the `10` state. Step 5 ALSO asserts `#margin` == '0' which
HEAD satisfies. Whether this is a spec bug or a wizard regression depends on the
Q3 gate answer, which is recorded RESOLVED = B (upstream stays as upstream has
it; net4sats may differ downstream). Record it as a fixture value to pin, not a
silent failure.
Status: VERIFIED (contract present), DRIFT recorded (devsplit 0 vs 10)
Files touched: (record only)

---

1|# PROGRESS — t_be3261d3: Claim GL-MT6000 lock for owner t_be3261d3
2|
3|Finding -> status -> files touched. Worktree root: /home/c08r4d0r/repos/net4sats-mvp-v2
4|Mode: bookkeeping only. Zero device contact, zero lock writes.
5|
6|## STATUS: BLOCKED — claim NOT made. Stopped at step 1 (claim) per the task's own stop rule.
7|
8|1. Finding: the documented lock convention is a lock FILE, not a channel/column/script, at
9|   `/home/c08r4d0r/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock` (mode 600); claim is a
10|   read-modify-write with 7 fields; release is a line-anchored `^released=` line or file deletion.
11|   Recovered verbatim from `verify/t_6b512dad-lock-convention/recon-quote.txt` (Sources A/B/C/F).
12|   Status: VERIFIED (re-read the recon capture directly, did not trust a summary).
13|   Files touched: none.
14|
15|2. Finding: the freeness precondition was NEVER satisfied upstream. The gate card `t_c895de5b`
16|   is `status=triage`, `result` NULL, 0 comments; sibling inspector `t_6d8235d5` likewise `triage`
17|   with no result. The prior claim attempt `t_d0c3e172` sits `blocked` (block_kind=completion,
18|   result NULL) since 2026-10-12. There is no positive freeness verdict to consume.
19|   Status: VERIFIED (`sqlite3` SELECT on the board DB, reproduced in the capture file).
20|   Files touched: none.
21|
22|3. Finding: the live store on this node is HELD by a different owner — `owner_task=
23|   net4sats-mvp-v2:t_2d0c6fa8`, `holder=t_2d0c6fa8`, and `grep -c '^released='` returns 0 (rc=1),
24|   which the convention defines as STILL HELD. Recorded pid 2174746 is dead and that owner's card
25|   is `archived`, so the hold is substantively stale — but staleness is not release, and the
26|   convention explicitly forbids "rename the file, force-clear, or steal a live claim". No
27|   documented takeover/staleness rule exists for this store. I did NOT overwrite it.
28|   Status: VERIFIED (file read + sha256 + documented release test + `ps`).
29|   Files touched: none.
30|
31|4. Finding: no writable claim target exists from this node. The store the prior claim used is the
32|   peer-node prta kit at `/home/c03rad0r/physical-router-test-automation/hardware.lock` — absent on
33|   x280 (ENOENT), peer host unreachable (`192.168.2.32` no route to host; `dq05` -> `dq05.fips`
34|   unresolvable), and its recorded durable claim-record copy is absent in the plan dir here.
35|   Status: VERIFIED.
36|   Files touched: none.
37|
38|5. Finding: no lock state was persisted because no claim was made. What IS persisted is the
39|   evidence trail the sibling announcement (`t_0841f821`) and verification tasks can read back:
40|   raw captures + a plain-language blocker set with the exact unblocking actions.
41|   Status: DONE.
42|   Files touched: `verify/t_be3261d3-lock-claim/lock-state-capture.txt`, `PROGRESS.md`,
43|   `REPORT.md` (full report, written before the final reply).
44|
45|## Explicitly NOT done
46|- No lock file created, appended to, modified, or deleted.
47|- No `released=` line appended (that would be a false release of someone else's claim).
48|- No lock handle, claim timestamp, or expected release time invented — all three are reported as
49|  NOT ISSUED / NONE, because a fabricated handle is worse than an honest null.
50|- No ssh/scp/ping/curl/nc to GL-MT6000 or any device; nothing flashed, rebooted, or power-cycled.
51|
52|## Unblock path for the manager (detail in REPORT.md)
53|1. Run `t_c895de5b` / `t_6d8235d5` to a real freeness verdict.
54|2. Resolve the stale `t_2d0c6fa8` hold by documented release or explicit operator decision.
55|3. Decide the store (x280 peripheral_locks vs peer c03rad0r prta kit) — the two are unsynchronised.
56|4. Re-dispatch this card.
57|
---
# PROGRESS — t_05167730: Write GL-MT6000 readiness note with evidence and blockers

Card: fleet offload worker, repo net4sats-mvp-v2
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Host: x280 / CobradorWave (c08r4d0r), Ubuntu 24.04.3 LTS
Mode: documentation only. Zero device contact, zero code edits, zero lock writes.

## R1 — VERDICT = BLOCKED (all four PREREQ tasks failed to complete)

Finding: PREREQ required requirements + lease booking + device smoke + Playwright/runner to
have COMPLETED. All four are `status=triage`, `result=NULL` in
`~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`:
t_b009e963 (requirements, 0 comments), t_77c59d23 (lease, 0 comments),
t_a8b5fcae (smoke, 1 comment = a manager urgency note only), t_63ec0ed2 (runner, closed
2026-10-11 with verdict NOT READY, then blocked).
Status: BLOCKED. 10 BLOCKER lines in the note.

## R2 — Three of the four "missing" fields are missing because their card never ran

Finding: required upstream build hash, booked hardware window, and the runner host name all
flow from cards that never produced a result. The plan ref `schedule.md` has NO WZ1.8 entry
(`grep -n WZ1.8` -> 0 hits), so the requirement must come from card bodies or the manager.
Status: RECORDED as BLOCKER B-REQ / B-LEASE.

## R3 — Stale-claim re-measurement: the 2026-10-23 "no IP path" finding still holds

Finding: this host has SINCE gained `192.168.1.50/24` on `enx6c6e0750cd79`, which invalidates a
naive carry of the smoke report's reachability claim. Re-probed directly this run:
192.168.1.1:{22,443,2121} TimeoutError; :{80,8090} Errno 113 No route to host; same when the
source is pinned to the .50 address; `ip neigh` -> 192.168.1.1 FAILED on archy-lab-br (linkdown).
`ping` is not installed on this box (rc=127) so TCP probes were used instead. Blockers survive.
Status: VERIFIED (fresh measurement, not inherited).
Files touched: none.

## R4 — Lock re-read this run: still held, no `released=` line, digest unchanged

Finding: `gl-mt6000-192.168.1.1.lock` sha256 `ed8a4573...f46582` (identical to the
t_be3261d3 capture), `grep -c '^released='` -> 0 (rc=1), holder `t_2d0c6fa8` (archived),
release conditioned on run task `t_557bb77d` child 3. No staleness/takeover rule exists in the
recovered convention. A worker must not take it.
Status: BLOCKED (fd4 serialization unresolved).

## R5 — Runner toolchain: this host cannot run the 7-step spec at all

Finding: node v18.19.1 but npm AND npx absent; no system chrome/chromium at all 6 probed paths;
no ~/.cache/ms-playwright; no node_modules/@playwright/test in the checked-out kit; the 7-step
spec + config exist ONLY on `origin/feat/router-vm-lab` (tip ee289bf8…), NOT on the checked-out
`fork-main-sync` (tip 5344e35). Carried runner figures from t_63ec0ed2 name user `c03rad0r`
paths that do not exist on this host -> could not be re-verified here.
Status: BLOCKED (ffmpeg-1011 + chromium-1223 + spec-on-branch).
Files touched: none.

## R6 — Deliverable written to both required places, byte-identical

Finding: note written at
`/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/readiness-gl-mt6000-wz1.8.md` and a
tracked copy at `docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md`.
sha256 both `1dd9ad0f328234723fa35d4275a617b99eaaee1c3508e7af68ed04043cd2472d`.
Top-level `VERDICT: BLOCKED` on line 3, 10 `BLOCKER: ` lines.
Status: DONE.
Files touched: the two note copies, PROGRESS.md (append), REPORT.md (append).

## R7 — Honest non-claims (round-3 pass)

Finding: this card touched no application code and no device. The smoke evidence in the note is
pasted verbatim from `REPORT-t_a8b5fcae.md` §4 (committed `303c76e`) and is NOT a fresh
device measurement by this card. Flash/re-image is OUTSTANDING and owned by the tollgate flow
owner — the smoke card requires a sysupgrade image (factory/initramfs does not apply).
Status: STATED in note §5.

## t_gate8 recon — 7-step wizard flow enumeration (2026-10-24)

Finding: the 7-step wizard flow (Discover -> Connect -> Uplink -> Payouts ->
Advanced -> Deploy -> Live) is **NOT IMPLEMENTED** in this repo. `index.html`
is a flat 5-view form (`#scan-view`, `#select-view`, `#deploy-view`,
`#success-view`, `#error-view`) with Discover/Connect/Uplink/Payouts/Advanced
controls merged into one card. `grep -c data-testid index.html` = 0; 0 of 7
`step-*` containers; 2 `<h2>`s neither of which is a step heading; 0
Next/Continue/Skip buttons; no step-component files exist anywhere.

The 7 steps exist only as a declared contract in a deliberately RED spec in a
DIFFERENT repo: `physical-router-test-automation/tests/browser/
generic-wizard-7step.spec.mjs` + `wizard-flow.mjs` (the step registry).

Brief deviations (stated, not worked around): worktree `configurationwizzard`
and branch `feature/upstream-wifi-scan` do not exist on this host
(`configurationwizzard` is a kanban board name); the onboarding Playwright spec
is not in this repo. Did not switch worktrees/branches.

Advanced optionality: marked optional only by `optional: true` in the step
model and the `<summary>Advanced settings (optional)</summary>` disclosure
label. **No skip control is implemented** (`/skip/i` matches nothing).

Status: COMPLETE. Deliverable `docs/recon-wizard-steps.md` (426 lines) with
per-step tables and file:line citations; full report
`REPORT-t_gate8-recon-7step.md`.

Files touched: `docs/recon-wizard-steps.md` (new), `REPORT-t_gate8-recon-7step.md`
(new), `PROGRESS.md` (this appended section). No spec and no component file
modified — read-only recon, as briefed.
