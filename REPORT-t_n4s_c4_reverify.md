# REPORT — t_n4s_c4_reverify: C4 re-verify real payment on re-deployed bench

**Card:** `net4sats-mvp-v2/t_n4s_c4_reverify` (status `ready`, `assignee=manager`, `urgency=high`)
**Card body (verbatim, complete):** "After C3 passes: Cashu purchase end-to-end, balance page
renders (142.7MB-class accounting). Was verified pre-reflash."
**Task:** fleet offload worker, repo `net4sats-mvp-v2`, worktree
`/home/c08r4d0r/repos/net4sats-mvp-v2`, branch `worker-heavy/t_gate8-19b1b1c6`.
**Host:** x280 (`c08r4d0r`), Ubuntu, 2026-10-24.
**Commits (both pushed, verified):** `b2f56b7` then `4adae89`.

---

## VERDICT

**BLOCKED — the C4 precondition has not passed, and the bench the card names does not exist on
the wire.** No payment was attempted, so C4's first claim is asserted **neither PASS nor FAIL**.
The second claim (balance page renders) splits cleanly:

| C4 assertion | verdict | one-line basis |
|---|---|---|
| **Cashu purchase end-to-end** | **NOT RUN** | C3 unmet; no shell for the documented pre-purchase wallet/`trustedmac` hygiene; only reachable client is this host's own NIC (already known to the box); only mainnet mints advertised and the test mint is unreachable |
| **balance page renders (142.7MB-class accounting)** | **HTTP-surface PASS / interactive render NOT MEASURED** | page + all 4 bundles serve 200 on `:2051`; but accounting is `usage:0, allotment:0, session_active:false`, so a `142.7MB` render is post-purchase by definition |

I did **not** turn this into a pass. Two independent findings below (the C3 precondition, the
build-identity discriminator) show that running the purchase here would have produced a result
about the **wrong box and the wrong build**.

---

## 1. The precondition C4 depends on has not passed (decisive, measured)

C4's body is conditional: *"After C3 passes…"*. C3 is `t_n4s_c3_redeploy`, *"C3 re-deploy +
4-marker brand read-back on vanilla MT6000"*:

- board status **`triage`** — never dispatched, `result` NULL, 0 comments;
- its own body: *"BLOCKED on operator merges: installer #154/#150/#151/#153"*.

Measured this run with `gh pr view` against `felixfelix-bot/tollgate-installer`:

| PR | title | state | mergedAt | mergeable |
|---|---|---|---|---|
| #154 | fix(launcher): busybox `tr` has no character classes — asset URLs 404 on OpenWrt | **OPEN** | `null` | MERGEABLE |
| #150 | fix(branding): stop re-arming the nodogsplash `gatewaydomainname` redirect loop (#428) | **OPEN** | `null` | MERGEABLE |
| #151 | test(shell): port the whitelabel rebrand-literal gutter from the module | **OPEN** | `null` | MERGEABLE |
| #153 | feat(brand): accept a brand token and an asset slot, for white-labelling | **OPEN** | `null` | MERGEABLE |

All four are still awaiting the **operator** merge. C3 cannot pass without them, so C4's gate is
closed by an operator action, not by anything a worker can do. *This is a hand-off, not a task.*

## 2. There is no "re-deployed bench" — the reachable box is the locked GL-MT6000, pre-reflash build

- Exactly one TollGate box answers on the wire: **`192.168.78.1`** (USB NIC `enx9c69d3c70689`),
  TCP ports 22/80/2050/2051/2121/8090/8443 OPEN. (It drops ICMP — probe TCP, never ping.)
- Its identity is **byte-identical to the 2026-10-23 record**: `pubkey`
  `11ee2f24ff7a1fc2cbd4c9ddf29786d24315b4cc2d7a3c22552e6b8f5485c2f4` — the same value recorded in
  `PROGRESS-t_b46dea1c.md` P1, where this address is identified as the **GL-MT6000** (a locked
  board), *not* a GL-MT3000 and *not* freshly reflashed.
- **Build-identity discriminator (new, decisive):** the box's live balance page serves
  `balance-CJspyFkd.js` / `balance-5KSP87Gz.css` / `index-YkGiQMp2.js`, while this worktree's build
  (`admin/balance.html`) references `balance-BNg56sOO.js` / `balance-D7WaP2L-.css`. Two disjoint
  hash sets ⇒ **the deployed bundle is a different, older build than the tree under test.**
- So the card's premise — a *re-deployed* bench whose purchase flow "was verified pre-reflash" — is
  contradicted: this box is still on a **pre-reflash-class** build. Whoever issued C4 is describing
  a state that this bench is not in.

## 3. Balance page: the surface is real and loads (HTTP-layer PASS), the accounting is not paid

Correcting my own first reading (I initially wrote "not served"; that was too strong and is
corrected in `PROGRESS-t_n4s_c4_reverify.md`):

- `http://192.168.78.1:2051/balance.html` → **200, 935 B**, a fully-formed Vite entry, and all four
  of its bundles serve **200** from the same instance:
  `assets/balance-CJspyFkd.js` (3973 B), `assets/index-YkGiQMp2.js` (361923 B),
  `assets/balance-5KSP87Gz.css` (954 B).
- The page's accounting is **byte-based and real**: the deployed `balance-CJspyFkd.js` polls
  `` fetch(`${N()}/balance`) `` every 5000 ms and formats `metric === 'bytes'` as
  `GB / MB / KB / B` (`>= 1048576` → `x/1048576` MB). A `142.7 MB` render is exactly the shape this
  formatter produces — consistent with the card's "142.7MB-class accounting".
- **But** on this box the API it polls reports no session: `:2121/balance` →
  `{"status":1,"session_active":false,"usage":0,"allotment":0,"remaining":0,"mac":"9c:69:d3:c7:06:89"}`
  and `:2121/usage` → **`-1/-1`** (the documented "no session" value). A `142.7MB-class` render
  therefore **cannot be produced without a real purchase** — which is the very thing blocked.
- The **`:8090` admin instance is broken as a surface** on this box: every path 307s to
  `https://192.168.78.1:8443/…`, and that vhost 404s (`<h1>Not Found</h1>`, 79 B) — including the
  build's own `assets/balance-BNg56sOO.js` and `manifest.json`. This is **device state, not a build
  defect**: this repo's `deploy.go:847` deliberately sets `uhttpd.luci.redirect_https='0'`.
- Related defect, same family: `:2050/splash.html` 307s to a **dead** host
  (`http://192.168.77.1:2050/…`); the sibling card `verify/wz-e2e-captive-t_gate8` records the same.

**Browser render NOT measured.** I have Playwright + chromium on this host. I deliberately did not
drive them: with no paid session the page can only render `0 B` / `—`, so a screenshot would prove
nothing about the C4 claim. I state that as a decision, not as an oversight — do not read the
absence of a screenshot as a plumbing failure.

## 4. Why the purchase could not honestly be run (four measured blockers)

1. **No shell.** `sshpass -p test123 ssh root@192.168.78.1 'echo SSH_OK'` →
   `Permission denied, please try again.` — the fourth consecutive session to hit this
   (`PROGRESS-t_b46dea1c.md` P2, `REPORT-t_b46dea1c.md:117`,
   `verify/wz-e2e-captive-t_gate8/PROGRESS.md:148`). The operator public key is not installed on the
   box. The documented purchase prerequisites that
   `physical-router-testing` → `references/cashu-payment-session-flip-proof.md` calls **Blockers 1
   and 2** — clearing a `trustedmac` exemption and normalising a poisoned `wallet.db` — both need a
   shell. I could not have cleared either, so I could not have produced a *clean* purchase even if
   the money had been available.
2. **The only reachable client is this host's own NIC.** `:2121/balance` reports
   `"mac":"9c:69:d3:c7:06:89"` = `enx9c69d3c70689`. The module authorises the MAC of the
   **requesting socket**, so a purchase from here tests an already-seen client, not a customer.
   Worse, the card's own adversarial run on this branch
   (`verify/wz-e2e-captive-t_gate8/PROGRESS.md` A4) proved a **never-seen MAC**
   (`02:e2:01:0e:aa:07`, `wzfresh0` 192.168.78.222) already reaches the internet **unfiltered,
   without paying**. So a purchase "succeeding" from that vantage could not distinguish a real gate
   open from that pre-existing pass-through defect.
3. **No test-credit path.** The box advertises **6 mainnet mints** (`price_per_step=cashu,1,sat`
   each). The documented test mint `nofee.testnut.cashu.space` is **unreachable from this host**
   (`curl …/v1/info` → `000`), while the 6 mainnet mints and `example.com` all return `200`. Buying
   the box's advertised price therefore means **spending the operator's real sats** — a spend
   decision that is not mine to make, and one I did not make.
4. **Unit check (so the number in the card is meaningful).** The advertisement sets
   `metric=bytes`, `step_size=22020096` (21 MiB). `142.7 MB` ≈ **7 steps**; one step = 1 sat, so a
   `142.7MB` session is a **7-sat** purchase at the advertised rate.

## 5. What a genuinely re-deployed bench would need (for whoever re-issues this)

The card cannot be satisfied until all of the following are true; each is operator-side:

1. Installer **#154/#150/#151/#153 merged** (or C3 explicitly waived) so C3 can pass.
2. A bench that is **provably freshly deployed**: report it with its own
   `pubkey` **and** its served asset hashes, so a verifier can tell it apart from the
   pre-reflash build still live at 192.168.78.1 (this is the check that would have caught the gap).
3. **Working root access** on that bench (password or key), so the documented `trustedmac` and
   `wallet.db` hygiene can be performed before the purchase.
4. A **non-self client vantage** (macvlan / real station on the merchant LAN) that is *not* this
   host's NIC, with the policy table pinned so internet-bound probes do not leak out of the house
   NIC.
5. Either a reachable **test mint** (`nofee.testnut.cashu.space` — dead from here) or an explicit
   operator go-ahead to spend **real mainnet sats** (7 sat per 142.7 MB session at current price).

---

## REMAINING STEPS (numbered, hand-off ready)

1. **Operator:** merge installer #154 / #150 / #151 / #153 (all OPEN, MERGEABLE, `mergedAt=null`) —
   this is the single blocker holding C3, and therefore C4.
2. **Operator:** unblock `t_n4s_c3_redeploy` (board status `triage`) and run it to completion on a
   genuinely stock bench, capturing the 4-marker read-back.
3. **Operator/worker:** identify the re-deployed bench by `pubkey` **and** served asset hashes and
   record them in the card, so the pre-reflash box at `192.168.78.1` cannot be mistaken for it.
4. **Worker (needs step 3 + root):** acquire shell on that bench; clear any `trustedmac` exemption
   and normalise `wallet.db` if poisoned (the two documented pre-purchase hygiene steps).
5. **Worker:** establish a non-self client vantage (macvlan `wzfresh0`-style MAC, policy table
   pinned to the bench, verified via `ip route get … table <n>`), then confirm
   `ndsctl … state:"Preauthenticated"` before paying.
6. **Operator decision:** either restore a reachable test mint or authorise a real mainnet spend of
   ~7 sat for a `142.7MB` session.
7. **Worker:** capture `:2121/balance` **before** (expect `session_active:false`, `allotment:0`) and
   **after** (expect `session_active:true`, `allotment` in the `462422016`/`142.7MB` class) as the
   paired assertion; screenshot the balance page render with the non-zero allotment.
8. **Worker:** verify the merchant-side receipt on-device — `wallet.db` sha256 must change across
   the purchase, `tollgate --json status` → `wallet_ok: true`, and the `kind:1022` session event
   captured from `logread` before the ring rotates.
9. **Verifier:** only then assert C4 PASS, citing the before/after `balance` JSON pair, the
   screenshot, and the on-device wallet delta.

**Stopped at:** step 0 — precondition verification. I completed the measurement/attribution work
(steps that need no payment and no shell) and stopped exactly where the honest options ran out:
no re-deployed bench, no shell, no non-self client, no test credit.

---

## Evidence and artefacts

- `PROGRESS-t_n4s_c4_reverify.md` (worktree root) — findings, corrected in place where my first
  reading was wrong.
- `evidence/wz-c4-n4s-c4-reverify/capture-20261024T055255Z.txt` — raw capture (port sweep, the
  `kind:10021` advertisement, `/balance`, `/usage`, balance-page surface probes, mint reachability,
  SSH refusal, build-identity discriminator).
- Commits: `b2f56b7`, `4adae89` — both pushed to
  `origin/worker-heavy/t_gate8-19b1b1c6`; `git rev-parse HEAD` == `git rev-parse origin/…` ==
  `4adae894aff9dbfbc8d3737ecc79bf673cdebf5f` at time of writing.

### Honesty notes

- I did **not** run a payment; the claim table above says `NOT RUN`, not FAIL.
- I did **not** fabricate a `142.7MB` render; I proved the formatter exists and reported the box's
  actual `0` / no-session state.
- My first balance-page reading was wrong and I corrected it in the same branch rather than leaving
  it; the superseded text is visible in the `4adae89` diff.
- Where a capture command wrote an empty worktree-hash line (a relative-path slip), I re-recorded it
  correctly and appended the amended block rather than editing the original capture.
