# PROGRESS — t_n4s_c4_reverify: C4 re-verify real payment on re-deployed bench

**VERDICT: BLOCKED.** The C3 precondition has not passed and no re-deployed bench exists
on the wire. No payment was attempted; C4's purchase claim is asserted neither PASS nor FAIL.
The one reachable box **does serve a working balance page** (HTTP-layer PASS), but it is the
**pre-reflash build** and its accounting is `0` / no-session until a real purchase happens.

Card: `net4sats-mvp-v2/t_n4s_c4_reverify` (status `ready`, assignee `manager`)
Body (verbatim): "After C3 passes: Cashu purchase end-to-end, balance page renders
(142.7MB-class accounting). Was verified pre-reflash."

Worker: fleet offload, host x280 (c08r4d0r). Worktree /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6

---

## Finding 1 — the C4 precondition (C3) has NOT passed
- STATUS: **BLOCKED** (precondition unmet; no payment was attempted, so C4 is NOT claimed)
- C4's body is explicitly conditional ("After C3 passes"). C3 is
  `t_n4s_c3_redeploy` "C3 re-deploy + 4-marker brand read-back on vanilla MT6000",
  board status **`triage`** (never dispatched, `result` NULL, 0 comments).
- C3's own body: "BLOCKED on operator merges: installer #154/#150/#151/#153".
  Measured this run with `gh pr view` on `felixfelix-bot/tollgate-installer`:
  all four are **OPEN, `mergedAt: null`, `mergeable: MERGEABLE`** — i.e. still awaiting the
  operator merge, unchanged from when C3 was parked.
- Files touched: none (read-only).

## Finding 2 — the reachable bench is the GL-MT6000, NOT the "re-deployed"/reflashed bench
- STATUS: MEASURED
- Only one TollGate box answers on the wire: `192.168.78.1`, on USB NIC `enx9c69d3c70689`.
  Ports 22/80/2050/2051/2121/8090 all OPEN (TCP, not ICMP — the box drops ICMP).
- Identity read back from it is **unchanged from the 2026-10-23 record**: `pubkey`
  `11ee2f24ff7a1fc2cbd4c9ddf29786d24315b4cc2d7a3c22552e6b8f5485c2f4`, same as
  `PROGRESS-t_b46dea1c.md` P1 (`192.168.78.1 = the GL-MT6000`, locked board).
- The card's premise is "re-deployed bench ... was verified pre-reflash". No data shows this
  box was re-deployed or reflashed since 2026-10-23; on the contrary it is byte-identical on
  every marker I can read without credentials.
- Files touched: none (read-only measurement).

## Finding 3 — the box's OWN local accounting: never used, and the allotment unit is bytes
- STATUS: MEASURED (device responses captured verbatim below)
- `GET :2121/balance` -> `{"status":1,"session_active":false,"usage":0,"allotment":0,"remaining":0,"mac":"9c:69:d3:c7:06:89"}`
- `GET :2121/` (kind:10021 advertisement) -> `metric=bytes`, `step_size=22020096`
  (21 MiB), 6 live mainnet mints, `price_per_step=cashu 1 sat`.
  462422016 bytes = 22 steps × 22020096 = 21 MiB/step → "142.7MB-class" = 7 steps.
- `GET :2121/usage` -> `200, 5 bytes`.
- Files touched: `evidence/wz-c4-n4s-c4-reverify/*` (new, this run).

## Finding 4 — the balance page DOES render, but on the :2051 portal — and it is at the pre-reflash build
- STATUS: **PASS** (the balance page surface is present and functional on this bench)
- Supersedes my first reading (recorded here for honesty): "admin balance page not served" was
  too strong. The corrected picture, measured:
  - `http://192.168.78.1:2051/balance.html` -> **200, 935 bytes**, and it is a **fully-formed Vite
    bundle load**, not a stub: the page's own JS/CSS all return 200 from the same instance —
    `assets/balance-CJspyFkd.js` (3973 B), `assets/index-YkGiQMp2.js` (361923 B),
    `assets/balance-5KSP87Gz.css` (954 B). So the balance page **loads**.
  - Its accounting render is **real and byte-based**, exactly the "142.7MB-class" shape the card
    names. From the deployed `balance-CJspyFkd.js`: it polls `` fetch(`${N()}/balance`) `` every
    `C=5e3` ms, and formats `metric==='bytes'` as `GB / MB / KB / B` (>=1048576 → `x/1048576` MB).
    `142.7 MB` = an allotment-class value rendered by this exact formatter.
- What that page reports on this box right now is however **not** a paid session:
  `:2121/balance` -> `usage:0, allotment:0, remaining:0, session_active:false` and
  `:2121/usage` -> `-1/-1` (the documented "no session" value). A `142.7MB-class` render therefore
  **cannot be reproduced without a real purchase** — which is exactly what C4 was to verify.
- The **:8090 admin instance is broken as a surface**: every path 307s to
  `https://192.168.78.1:8443/...` for a vhost that `404`s (`<h1>Not Found</h1>`, 79 B) — including
  the build's own `assets/balance-BNg56sOO.js`, `assets/balance-D7WaP2L-.css` and `manifest.json`.
  That is device state (a TLS/redirect_https vhost misconfigured on this box), not a build defect:
  this repo's `deploy.go:847` deliberately sets `uhttpd.luci.redirect_https='0'`.
- **Build-identity finding (decisive):** the live `:2051` balance page serves asset hashes
  `balance-CJspyFkd.js` / `balance-5KSP87Gz.css` / `index-YkGiQMp2.js`, while the build in this
  worktree (`admin/balance.html`) references `balance-BNg56sOO.js` / `balance-D7WaP2L-.css` /
  `index-Cv-gJynA.js`. **The deployed bundle is a DIFFERENT, older build.** Combined with the
  unchanged `pubkey` (Finding 2), this confirms the card's phrase "the re-deployed bench" does not
  describe this box: it is still running a **pre-reflash-class build**.
- Files touched: `evidence/wz-c4-n4s-c4-reverify/*` (new, this run).
- Related defect, same family: `http://192.168.78.1:2050/splash.html` -> 307 to
  `http://192.168.77.1:2050/splash.html` — a **dead** redirect target (192.168.77.1 does not
  answer); the sibling card `verify/wz-e2e-captive-t_gate8` recorded the same defect on this branch.
- Files touched: `evidence/wz-c4-n4s-c4-reverify/*` (new, this run).

## Finding 4b — the balance page's HTTP surface, and what I could NOT measure about it
- STATUS: PARTIAL (surface PASS; interactive render NOT measured)
- Measured: the SPA is served (200) with all four of its bundles (see Finding 4), and
  `:2121/balance` (the API it polls) answers with well-formed JSON. Those are HTTP-layer facts.
- **NOT measured:** an actual browser render (`:2051` -> paints "142.7 MB" or similar). I have
  Playwright + chromium on this host, but driving it was **deliberately deprioritised**: without a
  paid session the page can only render `0 B` / `—`, so a screenshot would prove nothing about the
  C4 claim, and the "142.7MB-class accounting" render is *definitionally* post-purchase. A visual
  PASS is therefore contingent on the purchase step, not on the plumbing.
- Files touched: none (decision recorded).

## Finding 5 — no credentialed route to the on-device wallet/mint, and the purchase client is self-exempt
- STATUS: BLOCKED (two independent blockers, both measured)
- **No shell.** `sshpass -p test123 ssh root@192.168.78.1` -> `Permission denied, please try
  again.` (same refusal as `PROGRESS-t_b46dea1c.md` P2, `REPORT-t_b46dea1c.md:117`,
  `verify/wz-e2e-captive-t_gate8/PROGRESS.md:148`). The operator public key is not installed
  on the box. Without a shell I cannot read/normalise the on-device wallet, cannot clear a
  poisoned `wallet.db`, and cannot clear a `trustedmac` exemption — all three are documented
  prerequisites of an honest purchase (`physical-router-testing` →
  `references/cashu-payment-session-flip-proof.md`, Blockers 1 and 2).
- **The only reachable client is the bench NIC itself.** `:2121/balance` reports
  `"mac":"9c:69:d3:c7:06:89"` — that is this workstation's own NIC
  (`enx9c69d3c70689`). The module authorises the MAC of the REQUESTING SOCKET, so running the
  purchase from the host would test an already-seen/possibly-exempt client, not a customer.
  The card's own adversarial run on this branch
  (`verify/wz-e2e-captive-t_gate8/PROGRESS.md` A4) proved a never-seen MAC
  (`02:e2:01:0e:aa:07`, `wzfresh0` 192.168.78.222) reaches the internet **unfiltered, without
  paying** — so a purchase "succeeding" from that client could not distinguish a real gate
  open from the pre-existing pass-through defect.
- Files touched: none (read-only).

## Finding 6 — the payable-hardware half cannot even be denominated from here
- STATUS: MEASURED (partial failure)
- The box advertises **6 mainnet mints**. `nofee.testnut.cashu.space` (the documented test mint
  for purchase lanes) is **unreachable from this host**: `curl https://nofee.testnut.cashu.space/v1/info`
  -> `000` (no connection), while the 6 advertised mints all return `200` and `example.com`
  returns `200`. So there is no test-credit path to buy the box's advertised price without
  spending real mainnet sats from the operator's wallet.
- Files touched: none (read-only).

---

## Status of the two C4 claims

| C4 assertion | status | why |
|---|---|---|
| Cashu purchase end-to-end | **NOT RUN** | precondition C3 unmet; no shell for wallet/trustedmac hygiene; only reachable client is self-exempt; only mainnet mints advertised; test mint unreachable |
| balance page renders (142.7MB-class accounting) | **FAIL (measured)** | admin balance page not served on this box (:8090 307→:8443 404); box's own `/balance` shows `usage:0, allotment:0` |

## Remaining steps
See `REPORT-t_n4s_c4_reverify.md` → "REMAINING STEPS" (numbered, hand-off ready).
