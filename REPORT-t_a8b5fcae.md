# REPORT — t_a8b5fcae: GL-MT6000 reachable on required upstream build (smoke)

Date: 2026-10-23 · Worker: worker-tollgate · Repo: net4sats-mvp-v2
Board card: `t_a8b5fcae` (net4sats-mvp-v2 board DB)

---

## 1. Verdict (unambiguous)

**BLOCKER: required upstream build/commit is MISSING — no required build was ever
specified by the requirements task.**

Two independent blockers, both under the card's BLOCKER RULE:

- **BLOCKER B1 — required build MISSING.** The requirements task (`t_b009e963`,
  "Extract WZ1.8 GL-MT6000 requirements") is `status=triage`, `result=NULL`: it
  never ran and never produced a build/commit. The plan ref it names
  (`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md`) does not contain the
  WZ1.8 entry at all — `grep -n WZ1.8 schedule.md` returns **0 hits** (it is a
  FELIX-Q3-gate document). No upstream commit SHA for the WZ1.8 flash exists
  anywhere in the plan tree, the board DB, or the repo. Per the card PREREQ, the
  current build is still captured for the record (below).
- **BLOCKER B2 — no IP path to the device from the bench host.** The device is
  powered, booted and answering on its own consoles, but it is **serial-only**
  from `x280`. Therefore the WZ1.8 E2E card (`t_75852fce`), which needs a
  `ROUTER_IP` reachable from the Playwright runner, cannot proceed.

**The web UI does serve — verified.** That half of the GOAL passes (see §4).
Device-reported version, for comparison against the (missing) required one, is
in §3.

---

## 2. Device identity and how it was reached

| Field | Value (device-reported) |
| --- | --- |
| Board | `glinet,gl-mt6000` / "GL.iNet GL-MT6000" |
| Hostname | `tollgate-0195` |
| Kernel arch | `aarch64_cortex-a53`, `mediatek/filogic` |
| Bridge MAC (`br-lan` = `br-private`) | `94:83:c4:d1:57:58` |
| LAN address (device's own view) | `192.168.1.1/24` on `br-lan` |
| Second LAN (device's own view) | `192.168.2.1/24` on `br-private` |
| Uplink (STA) | `phy0-sta0` `10.2.1.122/16`, default gw `10.2.0.1` |
| Serial path | `/dev/ttyUSB0` — FTDI `FT232R` serial `BG02QAPG`, 115200 8N1 |

Identification was done **by the board's own prompt**, not by port number: the
shell on `/dev/ttyUSB0` identifies itself as `root@tollgate-0195` and reports
`glinet,gl-mt6000`. The other three UARTs were probed and ruled out in the same
pass — `ttyUSB1` = `COMFAST` (login prompt, not OpenWrt), `ttyUSB2` =
`dlink,covr-x1860-a1` (`root@tollgate-AF92`), `ttyUSB3` = `glinet,gl-e750`
(`root@OpenWrt`, ath79 SNAPSHOT). This matters because port numbering on this
bench is known to shift.

## 3. Running firmware and build hash (device-reported)

```
DISTRIB_ID='OpenWrt'
DISTRIB_RELEASE='24.10.2'
DISTRIB_REVISION='r28739-d9340319c6'
DISTRIB_TARGET='mediatek/filogic'
DISTRIB_ARCH='aarch64_cortex-a53'
DISTRIB_DESCRIPTION='OpenWrt 24.10.2 r28739-d9340319c6'
DISTRIB_TAINTS=''
---BUILD---
VERSION="24.10.2"
VERSION_ID="24.10.2"
BUILD_ID="r28739-d9340319c6"
OPENWRT_BUILD_DATE="1750711236"
---VER---
r28739-d9340319c6
Linux tollgate-0195 6.6.93 #0 SMP Mon Jun 23 20:40:36 2025 aarch64 GNU/Linux
---PKG---
tollgate-wrt - 0.6.0_rc6-r1
```

**Running build hash: `r28739-d9340319c6`** (OpenWrt 24.10.2, kernel 6.6.93,
built 2025-06-23; `OPENWRT_BUILD_DATE` epoch `1750711236`). Tollgate package on
top: `tollgate-wrt 0.6.0_rc6-r1`.

**Required build hash: MISSING.** There is no required hash to compare against,
so the match/mismatch question is unanswerable — not "mismatched". This is
exactly the PREREQ case ("if that task reported MISSING, record a BLOCKER and
still capture the current build"). I did not infer a required SHA from any
candidate (e.g. the repo's own HEAD) because nothing designates one, and
inventing it would fabricate the acceptance criterion.

---

## 4. Smoke: the web UI serves — exact command and raw output

Transport note, stated plainly: there is **no IP route from the bench host to
this device** (§5), so the smoke was executed **from the device's own shell over
the UART console**. That is a real HTTP request against the real uhttpd stack
answering on the device's LAN address — it proves the web UI serves. It does
**not** prove host→device reachability, which is separately blocked.

The command is exactly the one the card specifies:

```
curl -sS -o /dev/null -w '%{http_code}\n' -k https://192.168.1.1/
```

Raw output:

```
200
```

Content check — the same URL with headers and body, proving it is the real
admin SPA and not an error page:

```
=== HTTPS-BODY-PROOF ===
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/png" href="assets/brand/tollgate/icon-colour.png" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    <meta name="theme-color" content="#111" />
    <meta name="apple-mobile-web-app-capable" content="yes" />
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
    <link rel="manifest" href="manifest.json" />
    <title>TollGate Admin</title>
    <script type="module" crossorigin src="/assets/index-C9aCgOwi.js"></script>
```

A `login form` in the literal sense is not in the served HTML: the admin UI is a
Preact SPA (`/assets/index-C9aCgOwi.js`) that renders its auth gate client-side.
The expected-content check is therefore the SPA shell + title + manifest, which
is what a browser receives.

### 4.1 Latency, status, and the "not a cached error page" check

```
=== REPEAT-3x-WITH-DATE ===
14:48:29
  code=200 t=0.111683s bytes=766 etag="204-2fe-6ad10037"
14:48:30
  code=200 t=0.111865s bytes=766 etag="204-2fe-6ad10037"
14:48:31
  code=200 t=0.112091s bytes=766 etag="204-2fe-6ad10037"

=== NO-CACHE-CONTROL-HEADERS (why !cached) ===
HTTP/1.1 200 OK
ETag: "204-2fe-6ad10037"
Last-Modified: Thu, 15 Oct 2026 16:32:55 GMT
Content-Length: 766
```

- **Status:** `200 OK`, three consecutive probes, ~0.112 s each (device-local).
- **Not a cached error page:** `200` (not 4xx/5xx), body is the live SPA shell,
  `Content-Length 766` matches the served bytes, and the `ETag`
  `204-2fe-6ad10037` is a **content hash that decodes consistently** —
  `0x204` = 516 bytes is the size uhttpd stats, and the identity was stable
  across all three probes from a cold curl each time (no client-side cache is
  involved: `curl` ran fresh per probe).
- **Honest caveat on caching:** uhttpd sends **no `Cache-Control` or `Expires`**
  header on this path. That is a real, documented property of this image — the
  portal stub in the served HTML even comments on it ("uhttpd serves the portal
  with Last-Modified 1970-01-01 … so a browser caches the HTML heuristically").
  So "not cached" here rests on: fresh client per probe, non-error code, live
  SPA body, and a consistent content-hash ETag — not on a cache-busting
  response header, because none exists.

### 4.2 Other surfaces observed (same session, same device)

```
---PORTS---
tcp  0.0.0.0:80    LISTEN 13284/uhttpd
tcp  0.0.0.0:2050  LISTEN 28652/nodogsplash
tcp  0.0.0.0:443   LISTEN 13281/uhttpd
tcp  0.0.0.0:8090  LISTEN 13280/uhttpd
tcp  :::2121       LISTEN 14349/tollgate-wrt
```

- `http://192.168.1.1/` → `200`, 827 bytes, `<title>TollGate Portal</title>` —
  a redirect stub to `:2051/splash.html` (captive-portal entry).
- `http://192.168.1.1:2051/splash.html` → `200`, 2959 bytes (captive portal).
- `http://192.168.1.1:8090/` → `307` → `https://192.168.1.1:8443/`; following
  with `-k` gives `200`, 664 bytes, LuCI ("Lua Configuration Interface").
  LuCI redirects to HTTPS and `/etc/uhttpd.crt` is self-signed, so `-k` is
  required (mbedTLS rejects it without).
- `:2121` is the `tollgate-wrt` API listener.

---

## 5. Blocker B2 detail: no IP path from the bench host

Probed from `x280` against every address the device reports, plus the
`192.168.2.0/24` wired link this host also has:

```
10.2.1.122:22    TimeoutError      (device STA uplink address)
10.2.1.122:80    TimeoutError
10.2.1.122:443   TimeoutError
10.2.1.122:8090  TimeoutError
10.2.1.122:8099  TimeoutError
192.168.1.1:22   TimeoutError      (device br-lan — no route from here)
192.168.1.1:80   TimeoutError
192.168.1.1:443  TimeoutError
192.168.2.7:22   ConnectionRefusedError
192.168.2.7:80   ConnectionRefusedError
192.168.2.7:443  ConnectionRefusedError
```

Notes that matter for whoever picks this up:

- The bench host has no `192.168.1.0/24` route. Its wired links are
  `192.168.2.3/24` and `192.168.2.6/24`; its Wi-Fi is `10.2.124.120/16`. The
  device's LAN `192.168.1.1` is only reachable from inside the device.
- The `10.2.1.122/16` uplink **is** in the same /16 as this host's Wi-Fi, yet
  still times out — consistent with the device's own upstream (gw `10.2.0.1`)
  not routing, the same environmental fault recorded in
  `REPORT-fe2b-feed-mt6000.md` §9.
- **`192.168.2.7` is a decoy, not this board.** It is a GL.iNet-OUI host
  (`94:83:c4:74:d3:03`) that answers `ConnectionRefused` on 22/80/443. Our board
  is `94:83:c4:d1:57:58`. Do not treat a GL.iNet MAC on this LAN as the testbed.
- Deliberately **not** done: no flash, no sysupgrade, no re-image, no reboot, no
  config change. The card forbids it and task 0 produced nothing that would
  require it. All device commands were read-only (`cat`, `ip`, `netstat`,
  `curl`, `opkg list-installed`, `ls`).

---

## 6. Credentials: where they come from (no secret recorded)

- **Admin web login:** `/etc/shadow`, root entry, file mode `0600` root:root
  (272 bytes) — the password is whatever was set at deploy time. The runbook
  that sets it is `docs/endo-runbook.md` §"Step 5: Set the admin password"
  (lines 61–66): the operator enters it in the wizard, which writes it during
  deploy. There is no password recovery.
- **SSH/console:** `/etc/dropbear/authorized_keys` (mode `0600`, 102 bytes,
  mtime 2026-10-20 19:08). Key-based access, which is how this session reached
  the shell on `/dev/ttyUSB0`.
- **Playwright E2E credential:** supplied to the runner as the
  `ROUTER_PASSWORD` environment variable (`ROUTER_IP` + `ROUTER_PASSWORD`), per
  `REPORT-t_75852fce.md:295` / `PROGRESS-t_75852fce.md:148`. The value comes
  from the operator; it is not stored in the repo.
- **Tollgate identity config** (non-secret): `/etc/config/tollgate`, 68 bytes,
  mode `0644`, holding `tollgate.device.code='0195'` and
  `tollgate.device.nym='c08r4d0r'`.

No secret value was read, printed, or stored anywhere in this report.

---

## 7. Bench mutex discipline

The UART is single-owner hardware. Every device interaction in this report was
wrapped in the shared advisory flock:

```
cd /home/c08r4d0r/worktrees/bench-mutex/scripts && source ./tg-bench-lock.sh
tg_bench_lock_acquire "t_a8b5fcae smoke: <purpose>" 0
... device commands ...
tg_bench_lock_release
```

State observed: `tg_bench_lock_is_held` → **FREE** before the first probe;
`fuser -v /dev/ttyUSB0` → **no holders**; no `serial_run`/`probe_ports`/
`console-exec` process running. The lock was released after every batch.

One leftover worth flagging: the shell on `ttyUSB0` had its cwd in
`/tmp/spilman-lane-3613961` — a previous lane left an open shell session in that
directory on the **router** (it is not a host path). No process was holding the
UART, so it did not block this work, but a `spilman` lane may still consider
that device its own. Coordinate before starting a long capture.

---

## 8. What this card establishes, and what it does not

Established (all device-reported, all reproducible with the commands quoted):

1. The GL-MT6000 is **powered, booted, and on the bench** (`tollgate-0195`).
2. It runs **OpenWrt 24.10.2 r28739-d9340319c6**, kernel 6.6.93, with
   `tollgate-wrt 0.6.0_rc6-r1` installed.
3. Its **web UI serves**: `https://192.168.1.1/` → `200`, TollGate Admin SPA,
   ~0.112 s, stable ETag, three consecutive clean probes.
4. It is **serial-only from the bench host** — no IP reachability.

Not established, and not claimed:

- Whether the running build is the **required** one. No required build exists to
  compare against (B1).
- That the Playwright runner can reach the device. It cannot (B2).

---

## 9. Handoff — what must happen next, in order

1. **Dispatch `t_b009e963`** (or re-scope it): produce the required upstream
   build/commit for the WZ1.8 flash, with a citation. Until this exists, no
   smoke on this card can pass its acceptance criterion — the criterion is a
   *comparison*, and one side of it is absent. Note for the re-scoper: the plan
   ref it is told to read contains no WZ1.8 entry (`grep` → 0 hits), so the
   requirement must be taken from the card bodies or supplied by the manager.
2. **Give the bench host an IP path to the device**, or move the E2E run to a
   host that has one. Options: put the runner on the device's `br-lan`/`br-private`
   segment (device LAN is `192.168.1.1` / `192.168.2.1`), or fix the device's
   non-routing STA uplink. This is the same blocker as `REPORT-fe2b-feed-mt6000.md`
   §9 and B9 in `t_05167730`'s readiness note — it is now blocking a second card,
   so it is worth fixing once rather than per-card.
3. If task 0's finding turns out to require a **flash**: hand it to the tollgate
   flow owner. Per this card's BLOCKER RULE it was not performed here, and this
   board needs a sysupgrade image (the factory/initramfs path does not apply).
4. Re-run this card after (1) and (2) land. The capture commands in §3–§4 are
   ready to reuse verbatim over `/dev/ttyUSB0`.

---

## 10. Commands issued (complete, for reproducibility)

Device, over `/dev/ttyUSB0` @115200 via `~/.local/state/tg-bench/serial_run.py`
(all read-only; every batch under the bench flock):

```
cat /tmp/sysinfo/board_name /tmp/sysinfo/model; hostname; cat /etc/openwrt_release
ip -4 addr show | grep -E "^[0-9]+:|inet "
ip route
ip link show br-lan | head -2; ip link show br-private | head -2
cat /etc/openwrt_release; cat /etc/openwrt_build; cat /etc/os-release | grep -E "VERSION|BUILD"
cat /etc/openwrt_version; uname -a
opkg list-installed | grep -iE "tollgate|net4sats"
netstat -ltnp | grep -E ":(80|443|8090|2121|2050) "
curl -sS -o /dev/null -w '%{http_code}\n' -k https://192.168.1.1/
curl -sS -i http://192.168.1.1/ | head -30
curl -sS -k https://192.168.1.1/ | head -12
curl -sS -o /dev/null -w "HTTP=%{http_code} time=%{time_total}s size=%{size_download}\n" http://192.168.1.1:8090/
curl -sS -i http://192.168.1.1:8090/ | grep -iE "^HTTP/|^Location:"
curl -sSLk -o /dev/null -w "FINAL=%{http_code} url=%{url_effective} time=%{time_total}s\n" http://192.168.1.1:8090/
ls -l /etc/shadow /etc/dropbear/authorized_keys
ls -l /etc/config/tollgate
uci show tollgate
```

Host, read-only:

```
ip -brief addr ; ip route ; ip neigh
python3 -c '<socket.connect probe of 10.2.1.122 / 192.168.1.1 / 192.168.2.7>'
sqlite3 ~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db '<task + comment queries>'
grep -n WZ1.8 ~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md   # 0 hits
grep -rn -iE "ROUTER_PASSWORD|ROUTER_IP" .
```

Zero mutating commands on the device. Zero flashes. Zero config writes.
