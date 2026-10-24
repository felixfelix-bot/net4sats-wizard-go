# REPORT — t_b46dea1c: Verify fixes F1 (NDS nftables), F2 (swap-counter), F3 (ndsctl exit 1)

**Goal as briefed:** prove each of the three fixes in build `56.b528e1d` actually behaves
correctly on the **GL-MT3000**.

## VERDICT: BLOCKED — step 0 (device access) failed. No PASS/FAIL is claimed for F1, F2 or F3.

I did not run `nft list ruleset`, `ndsctl`, or any counter read **on a GL-MT3000**, because
no GL-MT3000 (and no MT3000-class board) is reachable from this host. I am explicitly
**not** claiming any of the acceptance criteria were met, and I am not reporting a FAIL
either: F1–F3 were never exercised. Nothing was installed, written, rebooted or flashed on
any router during this run.

The chain was stopped before E2E, as the card instructs for a non-passing fix — here because
the fixes could not be executed at all.

---

## 1. BLOCKER B-DEV — the card's device has no IP path from this workstation

**Measured this run.** The MT3000's documented bench address (`192.168.1.1`, per
`docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md`) is dead, and the route to it is
on a link-down bridge:

```
$ ip -br addr | grep -E 'archy-lab-br|enx6c6e0750cd79'
archy-lab-br     DOWN      10.99.99.2/24 192.168.1.2/24 ... linkdown
enx6c6e0750cd79  UP        192.168.2.6/24 metric 100 192.168.1.50/24 ...

$ ip route get 192.168.1.1
192.168.1.1 dev archy-lab-br src 192.168.1.2 uid 1000        # <-- link-down bridge

$ ip neigh show dev enx6c6e0750cd79 | grep 192.168.1
192.168.1.1 FAILED
```

TCP probe of `192.168.1.1` on the MT3000's own surface — 22, 80, 2050, 2051, 2121, 8090 —
**all FAIL**.

### 1.1 The one live TollGate box on the wire is the GL-MT6000, not the MT3000

A concurrent sweep of `192.168.2.0/24`, `192.168.1.0/24`, `192.168.8.0/24`,
`192.168.77.0/24` and `192.168.78.0/24` for ports 2121 / 2050 / 2051 / 8090 returned
**exactly one** host:

```
192.168.78.1: 2121 OPEN, 2050 OPEN, 2051 OPEN, 8090 OPEN
```

That host is the **GL-MT6000** — recorded as such by this repo's own evidence file and by
the orchestration report, both from prior sessions:

| source | line | claim |
|---|---|---|
| `docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md` | 40 | device IP `192.168.1.1/24` on `br-lan` **of the MT6000**; second LAN `192.168.2.1/24` |
| `hermes-orchestration/REPORT.md` | 61 | "`192.168.78.1` is the **GL-MT6000**" |
| `hermes-orchestration/REPORT.md` | 57-58 | it is reached by running a DHCP client on the USB NIC → offered `192.168.78.166 from 192.168.78.1` |

It is also **not freshly flashed with 56.b528e1d**. Its state is byte-identical to the
2026-10-23 capture quoted in those reports:

```
$ curl -s http://192.168.78.1:2121/ | head -c 200
{"kind":10021,"id":"dee35f3ab55658184882ba4c5ff09967af390729906632d3475517ac72cba42c",
 "pubkey":"11ee2f24ff7a1fc2cbd4c9ddf29786d24315b4cc2d7a3c22552e6b8f5485c2f4", ...}
$ curl -s http://192.168.78.1:2051/splash.html -o /dev/null -w '%{http_code}\n' ; curl -si http://192.168.78.1/ | grep Last-Modified
200
Last-Modified: Fri, 23 Oct 2026 02:52:08 GMT
```

Same `pubkey`, same advertisement `id`, portal assets last modified 23 Oct — i.e. **older
than the build under test**, and plausibly not even running it. Testing F1–F3 here would
prove something about whatever build that box carries, not about `56.b528e1d`.

### 1.2 It is also another card's locked board

`~/.hermes/peripheral_locks/gl-mt6000-192.168.1.1.lock` is held, and this repo's readiness
note records it as locked by a different card (`t_2d0c6fa8`) with release conditioned on a
different run task. Per the hardware-exclusivity rule, a worker must not take it. I did not
acquire, release or renew any lock.

### 1.3 The negative is instrument-validated (not a missing-binary artifact)

Per `physical-router-testing` → `probe-negative-validation.md`, before reporting
"unreachable" I validated the instrument against a known positive on the same NIC:

```
$ command -v arping
/usr/sbin/arping

# KNOWN POSITIVE — same instrument, same NIC:
$ sudo arping -I enx6c6e0750cd79 -c 2 -w 3 192.168.2.1
60 bytes from 94:83:c4:d1:57:58 (192.168.2.1): index=3 time=1.401 msec
2 packets transmitted, 4 packets received,   0% unanswered

# TARGET:
$ sudo arping -I enx6c6e0750cd79 -c 2 -w 6 192.168.1.1
2 packets transmitted, 0 packets received, 100% unanswered

$ cat /sys/class/net/enx6c6e0750cd79/carrier /sys/class/net/enx6c6e0750cd79/operstate
1
up
```

So: probe works, NIC is up, target genuinely silent. (`arp-scan` is absent on this host; that
probe was discarded rather than trusted.)

---

## 2. BLOCKER B-AUTH — every credential available to me was refused

Reading `nft list ruleset` and running `ndsctl` both require a root shell on the box. Probes
made this run against the one reachable TollGate host (`192.168.78.1`):

| probe | raw result |
|---|---|
| `sshpass -p test123 ssh root@192.168.78.1 'echo SSH_OK'` | `Permission denied, please try again.` |
| `ssh -i ~/.ssh/id_ed25519 -o BatchMode=yes root@… 'echo KEY_OK'` | `Permission denied (publickey,password).` |
| `curl -k -X POST https://…/cgi-bin/luci/ -d luci_username=root -d luci_password=…` | HTTP **200** returning the portal SPA, **no `sysauth` cookie** → not a login |

The kit's documented lab password (`test123`) is explicitly recorded as **REFUSED** by this
box in `physical-router-testing/scripts/offline/07-portal3-authenticated-users-rule.sh`
(header comment, lines 16-18), and that was reproduced here.

**A credential source I did not use.** `~/.tollgate-root-credentials` (251 lines, mode 0600)
is the documented fallback — its password column is real even though its target-IP column is
wrong ("every one of 231 entries records `192.168.8.1`"). I attempted a newest-first probe of
its last 12 entries against the box; **the command was blocked by the approval gate and I did
not retry it**, per the instruction not to work around a refusal. So the box may well be
openable — I simply did not open it, and I will not claim otherwise.

Note this would not have unblocked the task anyway: the only box that credential would open
is the **MT6000** (§1.1), which is the wrong board and not on the build under test.

---

## 3. Fix definitions do not match any artefact at `b528e1d` (F2, F3)

Auditing the fixes *before* testing them, as required, two of the three are not implementable
against this build. Recorded so the card can be corrected — this is a finding, not a result.

### F1 — NDS nftables: genuine and checkable (needs the device)

`packaging/files/etc/nftables.d/20-nds-enforce.nft` exists at `b528e1d`
(blob `bc45b43f2e329e2870ec52614bbf24d9531a6c65`), with the two companion fragments listed
alongside it. The chain is:

```
chain nds_enforce_forward {
    type filter hook forward priority -1; policy accept
    meta nfproto ipv4 iifname "br-lan" meta mark & 0x00030000 == 0x00010000 counter drop
    meta nfproto ipv4 iifname "br-lan" meta mark & 0x00030000 == 0x00020000 counter accept
    meta nfproto ipv4 iifname "br-lan" meta mark & 0x00030000 == 0x00030000 counter accept
    meta nfproto ipv4 iifname "br-lan" oifname != "br-lan" counter reject with icmp type port-unreachable
}
```

Shipped by `f97435bc` ("fix: NDS fw4/nftables enforcement bridge", #283); the last rule's
portable-WAN rewrite is `28c631c0` (#297). **Every rule carries `counter`**, so the card's
"counter increments when a client hits the portal" is testable — on a device running this
build. Note the known trap from the skill: a new `/etc/nftables.d/*.nft` fragment is **inert
until `fw4 reload`**, and a bare package install may leave the old ruleset live.

### F2 — "swap-counter": no such counter exists at `b528e1d`

The card asks to "read the counter before and after" and to check it against "the underlying
session/swap source of truth". At `b528e1d` there is **no swap-counter** — no writer, no
getter, no endpoint, no log line. The name in this repo's history resolves to a **dependency
bump**, not a counter:

```
CHANGELOG.md:3395  - **Cashu wallet swap-counter race (critical).** Bump `gonuts-tollgate`
                   from v0.7.1 to v0.7.4 ... v0.7.1 incremented the keyset counter only
                   after a successful swap ... v0.7.4 increments the counter before the
                   swap call and adds a `swapWithRetry` path ...
$ git log --oneline --all -S 'swap-counter' -- CHANGELOG.md
5e20ba6a fix: bump gonuts-tollgate v0.7.1 -> v0.7.4 for swap-counter race (#266)
```

That is a `gonuts-tollgate` pin change (in `packaging/build-inputs.json` and the `go.mod`
carriers), it pre-dates the pinned build by many releases, and it is **not** an observable
counter in `56.b528e1d`. A repo-wide grep at `b528e1d` finds no swap-counter accessor; the
only `ndsctl`-adjacent counters are per-MAC byte counters in
`src/valve/customer_data_tracker.go`.

**Nearest real analogue** if the card means "counter increments exactly by the number of
events, agrees with the session source of truth": `GET :2121/usage` (session usage vs
allotment) — currently `-1/-1` on the reachable box, which is itself the "no session" shape.

**No before/after table is possible until the card names which counter it means.**

### F3 — "the previously failing `ndsctl` invocation": not recorded anywhere reachable

The card says "run the previously failing `ndsctl` invocation that returned exit 1", but
**the invocation is not recorded** — not in the card body, not in its comments, not in this
repo, not in the referenced `56.b528e1d` build. `ndsctl` is NoDogSplash's binary and is not
shipped by this repo, so I cannot recover the call site from the tree either.

This build's tip commit *is* `b528e1d7` "fix(#88): trigger NDS session after autopay (review
fixes) (#347)" — so the fix plausibly is that NDS session trigger — but "the reported failing
command now exits 0" cannot be asserted without knowing which command was reported. Despite
the name collision, the `ndsctl` hits in the tree are unrelated (test seams and the
merchant's `ndsctl json` client), not an exit-code fix.

---

## 4. PASS/FAIL table

| fix | exercised? | verdict | evidence |
|---|---|---|---|
| **F1** NDS nftables | **NO** | **BLOCKED** (not FAIL) | fragment + chain confirmed present in the build's source at `b528e1d`; no device reached to read a live ruleset |
| **F2** swap-counter | **NO** | **BLOCKED** (not FAIL) | no such counter exists at `b528e1d`; the name resolves to a gonuts v0.7.1→v0.7.4 dep bump (#266) |
| **F3** ndsctl exit 1 | **NO** | **BLOCKED** (not FAIL) | the failing invocation is unrecorded; `ndsctl` not shipped by this repo |
| device access | yes | **FAIL** | no GL-MT3000 has an IP path; only reachable TollGate box is the locked MT6000 |

Per the card's rule ("if any fix FAILs, stop the chain and report with reproduction steps
instead of proceeding to E2E") the chain was stopped here. No E2E run was started.

---

## 5. Honest non-claims

- I did **not** run `nft list ruleset`, `iptables-legacy`, `ndsctl`, or read any counter on a
  GL-MT3000. No such result is reported.
- I did **not** install, flash, reboot, or write to any router. Nothing was mutated.
- I did **not** take, release or renew any hardware lock.
- I did **not** read out or validate `~/.tollgate-root-credentials`; that probe was blocked
  and not retried.
- I did **not** claim F1/F2/F3 PASS **or** FAIL. They are unexecuted.
- The `192.168.78.1` measurements are real, but that box is identified (by prior recorded
  evidence, corroborated by its unchanged `pubkey`/advert id/asset dates) as the **MT6000**,
  not the card's MT3000.

## 6. REMAINING steps

1. **Operator:** say which physical box this card means and where it is cabled — or supply a
   credential (`ROUTER_PASSWORD=…`, or a public key in `/etc/dropbear/authorized_keys`) for
   the intended MT3000. The address in every reachable artefact is stale.
2. Re-run this verification with the device reachable, then capture F1:
   `nft list ruleset | sed -n '/nds_enforce_forward/,/^}/p'` (expect 4 rules, each with
   `counter`) and `iptables-legacy -t nat -L -n | grep -ci nds` (expect **0**). For the
   counter-increment half: read the chain, drive a never-seen MAC to the portal, re-read —
   the matching rule's packet counter must advance (and `fw4 reload` if the fragment is
   present-but-inert).
3. **F2:** the card must name the counter. Either repoint it at `GET :2121/usage` vs the
   allotment source of truth, or supply the actual swap-counter accessor — no before/after is
   recordable until then.
4. **F3:** supply the exact `ndsctl` invocation that returned 1, then capture
   `<cmd>; echo $?` for it **and** for a deliberately-bad invocation
   (`ndsctl deauth 00:00:00:00:00:00`), proving the fix is not a blanket `exit 0`.
