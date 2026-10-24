# REPORT-t_912534bf — Adjudication of every device-contact claim; evidence level 1–4

Card: `net4sats-mvp-v2:t_912534bf` — “Adjudicate every device-contact claim; set evidence level 1-4”
Workspace: `/home/c08r4d0r/repos/net4sats-mvp-v2` (branch `worker-heavy/t_gate8-19b1b1c6`)
Host: x280 (`c08r4d0r`). Pass: first execution of this card on this host.
Mode: read-only over every subject artefact. No device command was run. No network probe was made.
Local-only exclude rules match this filename (`/REPORT-t_*.md`, `/PROGRESS-t_*.md` in the shared
`info/exclude`), so both are force-added and tracked — the repo's existing convention: 33
`PROGRESS-t_*` and 40+ `REPORT-t_*` files are already tracked the same way.

---

## 0. VERDICT, UP FRONT

1. **No input gate fired.** The raw log exists and is readable; therefore no WITHHELD form is
   emitted. The upstream locate task (`t_47795719`) emitted nothing at all — no path, no WITHHELD
   string, zero comments/events/runs (`PROGRESS`/`REPORT.t_1e234edf.md` §0; board DB read
   transcribed at `verify/grepsweep-t_1e234edf/../fd4-device-contact-sweep-t_1e234edf/` in the
   sibling tree). Its absence does not gate this card: a readable primary log was still designated,
   and the card's stop condition is “if either emitted WITHHELD”. Neither did.
2. **The report under test is a set of negative assertions, not device-contact claims.** I extracted
   the six device-contact assertion-shapes it actually tests. Every one of them is a *negative*
   assertion (“no line exists showing…”). **Zero positive device-contact claims are made by any
   version of the evidence report, and therefore zero could be confirmed.** I did not invent
   positive claims to adjudicate — the card body's illustrative examples are labelled as such.
3. **Evidence level for the adjudicated subject: 1** — and level 1 here is *below* the card's level-1
   description, because `192.168.1.1`/`2121` do not even appear in the raw log. They appear only in
   the spec source in the repository, which is a different artefact from the log. Level 1 alone =
   **NEGATIVE / ABSENT for device contact**.
4. **The actual endpoint is named: `http://localhost:8099` on `127.0.0.1`, port 8099** — 5 evidence
   records, all HTTP 200, all loopback. It is the host-side wizard web app, not the device.
5. **A competing, second "device-contact evidence" assertion set exists on this host and it is
   POSITIVE at level 4** — for `192.168.1.1:2121`, with a real response body, from a *different*
   task family (`fe2b-*`, 2026-10-16, via jumphost `cobrador`). So the true answer to “did anything
   on this host receive a device UI body from `192.168.1.1:2121`?” is **yes** — just not from the
   audited run. Both levels are stated, separately, below. Neither is hidden behind the other.

Confidence: **high** for the adjudicated subject. Pivotal missing evidence named in §6.

---

## 1. Inputs, and the WITHHELD gate

| input | card | status observed this pass | bearing |
|---|---|---|---|
| raw log path | `t_47795719` (worker-base) | **no output of any kind.** Board DB: `result` NULL body_len 1381; 0 comments, 0 events, 0 runs | not gating — see below |
| grep sweep inventory | `t_1e234edf` (worker-heavy) | **produced.** 5-file sweep + transcript, 84 commands, ad-hoc verifier 12/12 PASS | satisfied |

Neither emitted the WITHHELD form (“VERDICT withheld — BLOCKER: missing artifact …”), so the
adjudication proceeds. The raw log was located and designated by first-hand enumeration, not by
inheritance: `raw-run.log` is one of five `raw-run*.log` files in the FD4 evidence tree and is the
unsuffixed canonical member of the run family.

### 1.1 Primary raw log — identity, re-measured this pass

```
$ wc -lc raw-run.log
  59 4585 raw-run.log
$ sha256sum raw-run.log
062c8b1a4a2e97cd9366c8646f6cf64b91090b1ff4b5f1c75f151122c4cdf209  raw-run.log
```

Absolute path of the copy I read and cite:
`/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log`

Two independent on-host copies were re-measured and are **byte-identical** (same 64-hex sha256):

- `/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log`
- `/home/c08r4d0r/worktrees/t_1eae5a7b/verify/fd4-t_be2366ea-device-traffic/raw-run.log`

A frozen copy of the two cited logs is committed with this work at
`verify/adjudicate-t_912534bf/raw-run.log` and `verify/adjudicate-t_912534bf/raw-runD.log`
(hashes 062c8b1a…c4cdf209 and 3f846632…0026beb), so every quote below re-verifies without the
other trees.

Five-file family, measured this pass (`raw-run.log` = PRIMARY; A/A2/C/D = siblings):

```
   59  4585 raw-run.log        sha256 062c8b1a...c4cdf209
   48  4049 raw-runA.log
   58  4356 raw-runA2.log
    7   381 raw-runC.log
   49  4766 raw-runD.log       sha256 3f846632...0026beb
  221 18137 total
```

---

## 2. STEP 1 — every distinct device-contact claim the evidence report makes, numbered

The evidence report under test exists in four revisions on this host; all four carry the **same six
device-contact assertion-shapes**. The nearest thing to a card-produced report (`t_be2366ea`) is the
uncommitted copy at
`/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/attachments/t_be2366ea/FD4-device-traffic-verbatim-t_be2366ea.md`
(14 770 B; sha256 `c45aaa1b9b6115c5b3b56f50d33825cd7c72b7ac417afcd505c347d2db7e9964`), byte-identical
to `REPORT.fd4-device-traffic-t_be2366ea.md` in the sibling tree. The revisions disclose their own
status: one says outright “**No fabricated device-contact claim exists in the repo either — every
artifact consistently records the run as not executed.** Trigger words `client connected to
192.168.1.1:2121`, `portal login page retrieved`, and `TCP handshake succeeded` **appear nowhere in
any version.** Of the four revisions, two declare **no log exists at all** (including the
t_be2366ea attachment above) and two report against a designated log. Both postures reach the same
device-contact conclusion. Claim list, therefore, is the set of *falsifiable* device-contact
assertions the family actually makes — all negative:

| # | device-contact claim (as asserted, verbatim intent) | asserted by |
|---|---|---|
| 1 | “client connected to `192.168.1.1:2121`” — asserted to be absent / a real connection exists | t_912534bf body (example); t_be2366ea task STEP 4 |
| 2 | “portal login page retrieved” — asserted absent | t_912534bf body (example) |
| 3 | “TCP handshake succeeded” to the device — asserted absent | t_912534bf body (example) |
| 4 | “at least one line showing a real request AND a real response to `192.168.1.1:2121`” — asserted: **no** | t_be2366ea STEP 4(a) |
| 5 | “one [request AND response] to the endo-onboarding UI on the device IP” — asserted: **no** | t_be2366ea STEP 4(b) |
| 6 | “this run touched the real device” — asserted: **no** | t_be2366ea STEP 4(c) |

**Step 1 is the step at which this card was previously mistried**: the card seen in `t_912534bf`'s
own log was archived with zero extractable claims because the report was then thought not to exist.
It does exist, in four revisions, and I extracted them above.

---

## 3. STEP 2 — claim-by-claim adjudication

Rule applied: quote the exact line verbatim with absolute path + line number, or mark UNVERIFIED.
No claim is restated that I did not personally locate. All quotes below were produced this pass by
`verify/adjudicate-t_912534bf/adjudicate-sweep.sh` (transcript `adjudicate-sweep.out`, committed).

Negative control proving the grep surface is live (so the zeros are real absences, not a broken
sweep) — verbatim:

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:2:Running 7 tests using 1 worker
line count: 59        ESC bytes: 7
```

| # | claim | adjudication | quote (absolute path : line) | tag |
|---|---|---|---|---|
| 1 | client connected to `192.168.1.1:2121` | **Not restatable — no log line asserts it, and I located no such line.** `192.168.1.1` = 0 hits and `2121` = 0 hits across all five logs. The nearest located text is a **refusal**, not a connection (see #6). | **UNVERIFIED** as a connection. Disconfirming quote: `/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:5` — `  ✓  2 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:111:1 › Step 2: Connect — select router and enter password (37.2s)` (that is a test *title*; no URL, no host, no transport follows it) | blocker |
| 1b | (sibling) the only connection-shaped **transport datum** anywhere in the family: an SSH banner captured on the wire from the device | **Located, but in a command log, not in the Playwright run log.** | `/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log:33` — `SSH-2.0-dropbear` (produced by the `head -c 40 <&3 after connect to :22` command echoed at `:32`, immediately after `192.168.1.1:22 OPEN`) | major |
| 2 | portal login page retrieved | **UNVERIFIED** — no response body from any device endpoint exists in the run log; `endo` = 0 hits, `onboarding` = 0 hits, `HTTP/1.` = 0 hits, `Host:` = 0 hits in all five logs. | **UNVERIFIED** (quote could not be located). | blocker |
| 3 | TCP handshake succeeded (to the device, in the run) | **UNVERIFIED** in the run log — `connected` = 0 hits, `SSH-2.0` = 0, `220 ` = 0, `curl` = 0. A device TCP success *is* located elsewhere in the evidence base (see §4b). | **UNVERIFIED** (run log). Located counterpart: `…/fe2b-device-baseline-pass5-2026-10-16/evidence.log:30` — `192.168.1.1:2121 OPEN` | blocker (run) |
| 4 | a real request AND response to `192.168.1.1:2121` in the run | **No such line.** 0 hits for both literals in all five logs; raw protocol-line shapes = none. | Disconfirming quote: `/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-runD.log:16` (opening) — `    Error: Refusing to select a deploy target: ROUTER_IP=192.168.8.1 is not in the discovery results [192.168.2.1 — unknown unknown \| …]` | blocker |
| 5 | a request AND response to the endo-onboarding UI on the device IP | **No such line.** `endo`/`onboarding` = 0 across 221 lines; the run's own captured HAR contains only loopback URLs (§4). | Disconfirming quote: `…/raw-run.log:10` — `  -  7 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:1 › Step 7: Live — verify TollGate API and portal on the router` (leading `  -  ` = did not run; this is the only step whose code navigates to the device) | blocker |
| 6 | did this run touch the real device | **No.** The run's device-facing legs never executed and its only peer was loopback:8099. | `…/raw-run.log:9` — `  -  6 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:183:1 › Step 6: Deploy — run the deployment to completion`; `…/raw-run.log:10` — `  -  7 … Step 7: Live …`; `…/raw-run.log:58` — `  2 did not run`; `…/raw-run.log:59` — `  4 passed (3.2m)` | blocker |

### 3.1 The mechanism, quoted (why no device line can exist in this log)

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:9
  -  6 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:183:1 › Step 6: Deploy — run the deployment to completion
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:10
  -  7 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:237:1 › Step 7: Live — verify TollGate API and portal on the router
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:58
  2 did not run
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-run.log:59
  4 passed (3.2m)
```

Corroboration in the sibling run whose line numbers pin it to the current spec revision:

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-runD.log:10
  -  7 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:243:1 › Step 6: Deploy — run the deployment to completion
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-runD.log:11
  -  8 [desktop-generic-wizard-7step] › generic-wizard-7step.spec.mjs:297:1 › Step 7: Live — verify TollGate API and portal on the router
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-runD.log:48
  6 did not run
```

Re-measured this pass against the current spec: `e2e/generic-wizard-7step.spec.mjs:243:test('Step 6:
Deploy — run the deployment to completion'` and `:297:test('Step 7: Live — verify TollGate API and
portal on the router'` — the runD-reported coordinates are exact. Those two tests are the only code
that navigates to the device (`spec.mjs:299` `page.goto(\`http://${ROUTER_IP}:${TOLLGATE_PORT}/\``,
`spec.mjs:311` `page.goto(\`http://${ROUTER_IP}:${PORTAL_PORT}/\``).

### 3.2 The only non-zero task term in the family, quoted in full

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/raw-runD.log:16
    Error: Refusing to select a deploy target: ROUTER_IP=192.168.8.1 is not in the discovery results [192.168.2.1 — unknown unknown | 192.168.2.13 — unknown unknown | 192.168.2.43 — unknown unknown | 192.168.2.21 — unknown unknown | 192.168.2.34 — unknown unknown | 192.168.2.33 — unknown unknown | 192.168.2.14 — unknown unknown | 192.168.2.25 — unknown unknown]. Deploying to the first discovered router would hit an unrelated device. Set ALLOW_FIRST_ROUTER=1 to override, or fix ROUTER_IP / the device.
```

`raw-runD.log` is the only file in the family carrying an IP literal; `raw-run.log` itself carries
none. The two `device` hits in the whole 221-line family are this line and
`raw-runD.log:24` (`      98 | \t\t\t\t\`discovered router would hit an unrelated device. Set \` +`),
re-measured: counts `device → 0 0 0 0 2` for raw-run / A / A2 / C / D.

---

## 4. STEP 3 — explicit evidence level

The card's scale: (1) the address string appears only in config/argv/env; (2) a connection ATTEMPT
is made; (3) a connection SUCCEEDS; (4) a device UI response BODY is actually received.

### 4a. For the adjudicated subject (the audited wizard run, `192.168.1.1:2121`)

**LEVEL 1 — and the level-1 literals are not even in the log.**

Strictly: `192.168.1.1` = 0 hits, `2121` = 0 hits in all five run logs. The literals exist only in
repository **source**, which is the level-1 description:

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/e2e/generic-wizard-7step.spec.mjs:49:const ROUTER_IP = process.env.ROUTER_IP || '192.168.1.1';
/home/c08r4d0r/worktrees/t_11d0873e/wt/e2e/generic-wizard-7step.spec.mjs:63:const TOLLGATE_PORT = process.env.TOLLGATE_PORT || '2121';
/home/c08r4d0r/worktrees/t_11d0873e/wt/e2e/generic-wizard-7step.spec.mjs:18: *   ROUTER_IP        default 192.168.1.1          (MUST appear in scan results —
```

By the card's own rule — “Level 1 alone means NEGATIVE/ABSENT for device contact” — the adjudicated
subject is **NEGATIVE / ABSENT**.

### 4b. For the strongest device-contact evidence located **anywhere in this host's evidence base**

**LEVEL 4** — a device API response body was actually received from `192.168.1.1:2121`. It belongs to
the `fe2b` task family, **not** to the audited wizard run:

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log:104
$ curl -s http://192.168.1.1:2121/
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log:105
{"kind":10021,"id":"7551ca8554a8c536eeaaf7b2b00203b1675d1a7ff5edc4611713fcb92d90a39a",
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log:115
$ curl -s http://192.168.1.1:2121/balance
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fe2b-device-baseline-pass5-2026-10-16/evidence.log:116
{"status":1,"session_active":false,"usage":0,"allotment":0,"remaining":0,"mac":"6c:6e:07:50:c3:3b"}
```

That file is 126 lines / 6 949 B, sha256 `562b1ab832e3438877464fafb9ac38f18130a4ea10cd5ee3247176f809f6575e`,
**tracked in git** (`git ls-files --error-unmatch` → success), self-describing as
`### PASS5 FE2b pre-baseline — GL-MT6000 package state (raw transcript)` /
`### Executed 2026-10-16T12:23:58Z .. 12:29Z by worker host x280 (c08r4d0r)`.

Convergent, independently-captured level-4 bodies for the same address exist in the repo report
`REPORT.tinstall-fd3-dispatch4.md:111` (`{"kind":10021,"id":"06cff9d704ae3b8d…`, captured
`2026-10-15T18:53:48Z`) and `REPORT.fe2b-device-baseline.md:87-92` (distinct `id`, same device MAC).
Note the two dispatched bodies are **different events** (`id` 06cff9d7… vs 7551ca85…) while the
advertising `pubkey` is identical (`6a8f7e6cbf0f47a2…`) — one device, two advertisements.

**Both levels are true at once and neither cancels the other:** level 1/NEGATIVE for the audited
Playwright run; level 4 for the fleet's direct-curl baseline of the same physical address, a day
later, through a different access path.

---

## 5. STEP 4 — host/port mismatch, named plainly

**A mismatch is found.**

The audited run's only network peer, taken from the run's own extracted HAR resource snapshots
(`traces-network-dump.txt`, 18 lines, 9 549 B), re-measured this pass:

```
resource-snapshot lines: 4
"url":"http://localhost:8099/"        "status":200  "serverIPAddress":"127.0.0.1"  "_serverPort":8099
"url":"http://localhost:8099/api/scan" "status":200 "serverIPAddress":"127.0.0.1"  "_serverPort":8099
"url":"http://localhost:8099/"        "status":200  "serverIPAddress":"127.0.0.1"  "_serverPort":8099
"url":"http://localhost:8099/api/scan" "status":200 "serverIPAddress":"127.0.0.1"  "_serverPort":8099
```

Every recorded request header on all four snapshots: `"name":"Host","value":"localhost:8099"`.
Content of those four responses: `"size":19976, "mimeType":"text/html; charset=utf-8"` (the wizard
page), `"size":1021, "mimeType":"application/json"` (scan #1),
`"mimeType":"application/json", "size":1143` (scan #2) — two scan runs, different lengths.

And the harness's own verdict on the device endpoint, verbatim from the dump:

```
/home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic/traces-network-dump.txt:13
(NONE)
```

**Named actual endpoint: `http://localhost:8099`, resolving to `127.0.0.1:8099`** — the host-side
wizard Go binary, not `192.168.1.1:2121`. The device address the run was pointed at was, moreover, a
**third** network address: `ROUTER_IP=192.168.8.1` (quoted §3.2), whose discovery set was entirely
`192.168.2.0/24` (re-measured from the `/api/scan` body at `traces-network-dump.txt:10`:
192.168.2.1, .13, .14, .20, .21, .25, .33, .34, .43 — `192.168.1.1` and `192.168.8.1` both absent).

Additional mismatch traps, named so they cannot be mistaken for traffic:

- The literal `192.168.1.1:2121` **does** occur in the evidence directory — only inside synthetic
  anti-vacuity fixtures: `verify-scripts-adhoc.sh:93`, `.repass-verify.sh:34`,
  `repass2026-10-16-vacuity.sh:7`. Fixtures, not run traffic.
- `192.168.1.1` also occurs inside the captured trace archive as the **text of the spec source**
  (`traces-inventory.txt:57-60,67-68` — `resources/src@7ac9bdc0….txt:18` and `:49`). Source text,
  not a request.
- This host **can** route to `192.168.1.1` (on `archy-lab-br`, src `192.168.1.2`) and to
  `192.168.8.1` (on `enx9c69d3c70689`, src `192.168.8.5`). A reachable route on the adjudicating
  host is not evidence about a run that happened on host `c03rad0r` on 2026-10-13.

---

## 6. STEP 5 — confidence, and the single most pivotal missing evidence

**Confidence: HIGH** for the adjudicated subject (level 1 / NEGATIVE).

Basis: two independently-hashed log copies agree byte-for-byte; the grep surface is proven live by a
positive control on the same file; the null survives case-insensitive widening across 20 patterns
(`ftp`, `sftp`, `endo`, `onboarding`, `Host:`, `GET `, `POST `, `HTTP/1.`, `connected`, `220 `,
`SSH-2.0`, `curl`, `localhost`, `127.0.0.1`, `8099`, `mt6000`, `192.168.1.1`, `2121`, `2050`);
and the failure mechanism is quoted, not inferred (`  -  ` = did not run on both device-facing steps,
in the primary and in the current-revision sibling). The repo also contains **no** `test-results/`,
`playwright-report/`, `generic-wizard-report/` or `blob-report/` anywhere (searched); `e2e/` holds
exactly README + spec + config.

**Pivotal missing evidence (one item):** the original `trace.zip` for the audited run — or its
`0-trace.network` member. The run was configured `trace: 'on'` and `video: 'on'`, and the reporter
output names its path (`…/raw-run.log:49` → `test-results/generic-wizard-7step-Step--9a5f9-…/trace.zip`).
That archive would settle device contact from the run's own packet record rather than from greps of
a `list`-reporter stdout that structurally cannot contain URLs. It is **absent from this host**: the
runner directory `/home/c03rad0r/.tmp/gate9-e2e-t_5d6e1bae` does not exist here (`/home/c03rad0r`
absent), the five logs' own attachment paths resolve to nothing in the repo, and the only `trace.zip`
files found on the host are three 2026-10-24 archives from a **different, later** run
(`/home/c08r4d0r/wz18-gate8-red/…` ×2, `/home/c08r4d0r/wz18-gate8-run/…` ×1, the last 9 255 B with a
**0-byte** `0-trace.network`). Those three are named here and are **not** evidence for the audited
run: one is a later run of the same spec, and a zero-byte network member carries no traffic.

Secondary item that would materially change the *whole-system* verdict (not the run verdict): a
first-hand re-run of the `fe2b`-style direct probe, to re-derive §4b's level-4 body from the raw
device on demand rather than from a committed transcript.

---

## 7. Contradictions and honesty section

1. **The family contradicts itself on whether a log exists.** Two revisions (2026-09-25, 2026-09-27)
   assert “no Playwright run log exists anywhere on this host”; two later revisions (2026-10-15,
   2026-10-16) report verbatim against a designated log. The larger of the two later ones, however,
   reports the task workspace `/home/c08r4d0r/repos/net4sats-mvp-v2` as non-existent (dangling
   worktree pointer) and read from `/home/c08r4d0r/worktrees/t_be2366ea/` — a path that does **not**
   exist on this host this pass. I therefore re-derived everything from
   `/home/c08r4d0r/worktrees/t_11d0873e/wt/` and `/home/c08r4d0r/worktrees/t_1eae5a7b/`, both of
   which exist, and whose `raw-run.log` is byte-identical to the hash those reports cite.
2. **Claim-shape correction to the earlier adjudication.** An earlier pass of this same card asserted
   that the strongest located device-contact evidence on this host is level 3 against
   `10.17.171.1:2121` (GL-MT3000), a different device from `192.168.1.1`. I could not locate that
   `10.17.171.1` artefact in the subject trees; I therefore do not adopt the claim. What I did
   locate, first-hand and quoted, is **level 4 against `192.168.1.1:2121` itself** (§4b). The card
   instruction was to name the *strongest* evidence, and level 4 supersedes level 3.
3. **Level 4 does not transfer to the run.** §4b concerns a direct `curl` run on 2026-10-16 from
   worker host x280 via jumphost `cobrador`. The audited Playwright run happened 2026-10-13T22:38–39Z
   on host `c03rad0r` and never emitted a device request. Citing §4b as evidence for the run would
   be exactly the misattribution this card exists to prevent.
4. **An unrelated 2026-10-27 session log carries live device curl traffic to `192.168.1.1:2121`**
   (board `fips`, `logs/t_9dc240a3.log:1097-1099`). Different board, different task, different day;
   it is not evidence about this repo's run. Noted so a future sweep cannot quietly promote it.
5. **Two reports on this host assert a device TCP/HTTP reachability state that the repo's own spec
   conflicts with.** `e2e/README.md:86-88` records the device exposing **0/1024 TCP ports** from the
   fleet network as of 2026-10-14; `REPORT.fe2b-device-baseline.md` (2026-10-15) and the pass5
   transcript (2026-10-16) record 22/80/2050/2051/2121/8090 **OPEN** via `cobrador`. Both are
   reports; I re-ran no probe, so I state neither as fact — they are in conflict and the conflict is
   unadjudicated here.
6. **What I did not verify.** The artefact decomposition for §4b is by surrounded command echo and
   body id. I matched the body id `7551ca85…` against the whole on-host evidence base and found no
   second occurrence; the public `pubkey` `6a8f7e6cbf0f47a2…` recurs across `tinstall-fd3` and
   `fe2b` captures. I did not re-derive the body from the raw device.

---

## 8. Method, commands, and reproduction

One script produces every number in §3–§5:

```
verify/adjudicate-t_912534bf/adjudicate-sweep.sh   # committed
verify/adjudicate-t_912534bf/adjudicate-sweep.out  # committed transcript (4 333 B)
verify/adjudicate-t_912534bf/verify-quotes.sh      # committed quote-integrity verifier
verify/adjudicate-t_912534bf/verify-quotes.out     # committed transcript, 18/18 PASS
```

```
bash verify/adjudicate-t_912534bf/adjudicate-sweep.sh
# defaults: evidence dir /home/c08r4d0r/worktrees/t_11d0873e/wt/verify/fd4-t_be2366ea-device-traffic
#           spec dir     /home/c08r4d0r/worktrees/t_11d0873e/wt/e2e
# both overridable as argv 1 and 2
```

Read-only throughout. Two `git add -f` additions (this report and `PROGRESS-t_912534bf.md`) sit
alongside the frozen raw logs and the sweep + verifier scripts and transcripts, in commits
`00a8b96` (evidence bundle) and `b18e6f7`/`HEAD` (verifier + report) of
`worker-heavy/t_gate8-19b1b1c6`, pushed to
`origin` = `https://github.com/felixfelix-bot/net4sats-wizard-go.git` (remote sha verified equal to
local HEAD by `git rev-parse HEAD` vs `git ls-remote origin`).

**Verification status: ad-hoc, not a project suite.** This repo exposes no canonical test command for
a log-evidence audit. What was actually run and observed: `bash -n` on the sweep script (clean);
`bash adjudicate-sweep.sh` (exit 0, transcript captured); the negative control in §3 (hit at line 2,
59 lines, 7 ESC bytes); two-copy sha256 agreement; `git ls-files --error-unmatch` on the §4b
transcript; `git ls-remote` after push. No test-suite green is claimed.

**Quote-integrity gate actually executed.** `verify-quotes.sh` extracts every `path:line` citation
from this report, resolves it against the freshly `sha256sum`-matched frozen copies, and prints each
resolved line. Observed result, verbatim from the committed transcript's last line:

```
RESULT: PASS=18 FAIL=0
ALL QUOTE-INTEGRITY CHECKS PASSED
```

Every citation in §3, §3.1, §4a, §4b and §5 that points into `raw-run.log` or `raw-runD.log`
re-resolved. The five quotes that point into *other* artefacts (`evidence.log`, the spec, the
trace dump, `dispatch4:111`, `README.md:86-88`) cannot resolve from the frozen copies and were
re-checked by direct `sed -n` against the live files this pass; all five matched byte-for-byte,
including the CRLF endings in `evidence.log` (`cat -A` → `OPEN^M$`) and the `(NONE)` at
`traces-network-dump.txt:13`.

### 8.1 Canonical project recipe — run this pass, and the one pre-existing failure

`hermes verify --detect-only --json` reports the repo's recipe as a Go project:
`build: ["go build ./..."]`, `test: ["go test ./..."]`, `start: "go run ."`, port 8099,
`evidence: ["Detected go.mod"]`. `/usr/local/go/bin` is not on this shell's PATH; the toolchain
(go1.22.12) was exported for the run.

```
$ /usr/local/go/bin/go build ./...      # BUILD_EXIT=0   (silent)
$ /usr/local/go/bin/go vet ./...        # VET_EXIT=0     (silent)
$ /usr/local/go/bin/go test -v ./...
--- PASS: 65 tests
--- FAIL: 1 test
    --- FAIL: TestNoInternalProcessLeaksOnShippedSurfaces (0.11s)
```

`TestNoInternalProcessLeaksOnShippedSurfaces` is the WZ1.9 debrand guard. It is the **only**
failure. Its `shippedSurfaces` set is `index.html`, `README.md`, `CHANGELOG.md`, `docs/*.md`,
`portal/|admin/|rpcd/**`, and `*.go` (`debrand_internal_test.go:52-102`). Three files trip it:

```
     56 hits  docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md     tracked, added by 22fbea4
     27 hits  docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md  untracked, foreign card t_cc76c746
     22 hits  docs/LOCK-ANNOUNCE-gl-mt6000-t_be3261d3-2026-10-24.md       tracked, added by ab33306
```

**This failure is not mine, and I did not touch it.** Four independent checks: (i) no line of the
failure output cites any file I wrote — re-measured, `NONE OF MY FILES CITED`; (ii) my four commits
touch ten paths (`git diff --name-only e8e6274..HEAD`) and none is in the test's scan set — the
guard file, `docs/`, and `index.html` are untouched (`git diff --stat` empty); (iii) the two tracked
offenders were added by commits that are **ancestors of my base** `e8e6274`
(`git merge-base --is-ancestor` → true for both); (iv) deterministic replay at the base commit in a
disposable worktree: `git worktree add --detach /tmp/base-e8e6274 e8e6274` → the same test FAILS
there with **78** `internal detail leaked` lines (mine, with the third file present, prints more
because the untracked `docs/WZ1.8-gate8-gate9-review-verdict-…md` is in scope on this working tree;
that file was authored by card `t_cc76c746`, not by this card). The worktree was then removed and
`git worktree prune` run.

The guard's own header states task-history artefacts are deliberately out of scope for it
(`debrand_internal_test.go:43-46`), which is why this card's `REPORT-t_*`/`verify/` output cannot and
does not participate. This report is therefore delivered with a **known pre-existing red** in the
repo's own suite, reported rather than hidden, and not repaired here — repairing it means editing
three `docs/*.md` files owned by three other cards, which is scope this card was not given.
