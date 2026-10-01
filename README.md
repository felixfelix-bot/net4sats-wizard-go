# net4sats wizard

Cross-platform onboarding wizard that turns an OpenWrt router into a Bitcoin
WiFi access point. Single Go binary — serves a web UI that auto-discovers
routers on your LAN and deploys net4sats over SSH.

```
┌─ Your Laptop ─────────────────────────────┐
│                                           │
│  net4sats-wizard binary                   │
│  └─ web UI at http://localhost:8099       │
│     └─ scans LAN for routers              │
│     └─ deploys via SSH                    │
│                                           │
└──────────┬────────────────────────────────┘
           │ SSH
           ▼
┌─ Router (OpenWrt) ────────────────────────┐
│  tollgate-wrt   (:2121)  payment backend  │
│  nodogsplash    (:2050)  captive portal    │
│  configurationwizzard (:8090) admin panel  │
└───────────────────────────────────────────┘
```

## Quick start

1. **Flash your router to OpenWrt** (wizard does NOT flash firmware). For
   GL.iNet routers, use the web UI at `http://192.168.8.1` → Advanced →
   Upload Firmware → untick "Keep settings".

   Alternatively, SSH in and sysupgrade:
   ```sh
   # Download the sysupgrade image for GL-MT3000 (OpenWrt 25.12.5)
   curl -O https://downloads.openwrt.org/releases/25.12.5/targets/mediatek/filogic/openwrt-25.12.5-mediatek-filogic-glinet_gl-mt3000-squashfs-sysupgrade.bin

   # Upload to router and flash (clean, no config kept)
   cat openwrt-25.12.5-mediatek-filogic-glinet_gl-mt3000-squashfs-sysupgrade.bin | \
     ssh root@192.168.1.1 'cat > /tmp/sysupgrade.bin && sysupgrade -n /tmp/sysupgrade.bin'
   ```
   See [OpenWrt firmware downloads](https://downloads.openwrt.org/releases/25.12.5/targets/mediatek/filogic/)
   for other devices. Use `-n` for clean install (no config kept), omit for
   in-place upgrade preserving settings.

2. **Set a root password:**
   ```sh
   ssh root@192.168.1.1
   passwd
   ```

3. **Download the wizard** from the
   [latest release](https://github.com/net4sats/net4sats-wizard-go/releases/latest):

   | OS | File |
   |---|---|
   | macOS (Intel) | `net4sats-wizard-darwin-amd64` |
   | macOS (Apple Silicon) | `net4sats-wizard-darwin-arm64` |
   | Linux (x86_64) | `net4sats-wizard-linux-amd64` |
   | Windows | `net4sats-wizard-windows-amd64.exe` |

4. **Run it:**
   ```sh
   chmod +x net4sats-wizard-*
   ./net4sats-wizard-darwin-arm64   # replace with your OS
   ```

5. **Browser opens at** `http://localhost:8099` — follow the wizard:
   - It scans your network and lists detected routers
   - Select your router, enter the root password
   - Choose upstream connection (Ethernet WAN or WiFi repeater)
   - Enter your Lightning address (where payouts go)
   - Click "Deploy net4sats"

6. After ~30 seconds: connect to the `net4sats-portal` WiFi and open any
   website — the captive portal appears with payment options.

## What the wizard does

The deployment runs 9 steps over SSH:

| # | Step | Action |
|---|------|--------|
| 1 | Verify SSH | Connects, reads `/etc/openwrt_release` |
| 2 | Check firmware | Parses OpenWrt version |
| 3 | Set password | Sets the root password you entered |
| 4 | Configure upstream | WiFi STA mode or WAN passthrough |
| 5 | Install net4sats | `apk add net4sats` (pulls tollgate + nodogsplash) |
| 6 | Brand portal | Sets gateway name to `net4sats` |
| 7 | Configure Lightning | Sets your Lightning address, dev split, margin, mint |
| 8 | Restart services | Restarts `tollgate-wrt` + `nodogsplash` |
| 9 | Health check | Verifies TollGate API responding on `:2121` |

## Request security (origin gate)

The wizard drives a **root SSH session** on the router, so its HTTP API is
gated:

- **Loopback bind by default** — it listens on `127.0.0.1:8099`. Opt into LAN
  exposure with `WIZARD_BIND=0.0.0.0:8099` (accepting that any LAN host can
  then drive deploys).
- **Origin allowlist, enforced on every request** — a request carrying an
  `Origin` header outside the allowlist is refused with `403` *before* any
  handler runs, reads and writes alike. The default allowlist is
  `http://localhost:8099` and `http://127.0.0.1:8099`; a concrete
  (non-wildcard) `WIZARD_BIND` address is added to it, so the browser UI
  served from that address keeps working.
- **Non-browser clients keep working** — requests with no `Origin` header
  (`curl`, the CLI, E2E harnesses) are unaffected.

Why the write path matters: a cross-origin "simple request" (for example a
`POST` with `Content-Type: text/plain`) needs no preflight, so merely
withholding `Access-Control-Allow-Origin` hides the *response* while the
request still reaches the server. Refusing foreign origins outright is what
stops a malicious page — or a DNS-rebinding host that spoofs the `Host`
header — from driving `/api/deploy` from the operator's browser. The `Host`
header is never consulted when building the allowlist.

## Verify binaries

```sh
sha256sum -c SHA256SUMS
```

Checksums are published with each [release](https://github.com/net4sats/net4sats-wizard-go/releases).

## Build from source

```sh
git clone https://github.com/net4sats/net4sats-wizard-go.git
cd net4sats-wizard-go
go build -o net4sats-wizard .
./net4sats-wizard
```

Cross-compile for all platforms:

```sh
GOOS=darwin  GOARCH=arm64 go build -o dist/net4sats-wizard-darwin-arm64 .
GOOS=darwin  GOARCH=amd64 go build -o dist/net4sats-wizard-darwin-amd64 .
GOOS=linux   GOARCH=amd64 go build -o dist/net4sats-wizard-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/net4sats-wizard-windows-amd64.exe .
```

## Documentation

- **[Setup guide (step by step)](https://github.com/net4sats/net4sats.github.io/pull/1)** —
  full walkthrough with screenshots, troubleshooting, and advanced settings
  (dev split, margin, mint selection)
- **[Admin panel guide](https://github.com/net4sats/net4sats.github.io/pull/1)** —
  using the configurationwizzard admin dashboard after deployment (dashboard,
  WiFi, devices, settings, wallet, identity)
- **[Endo onboarding runbook](docs/endo-runbook.md)** — operator-facing setup
  guide tested on a live GL-MT6000

## Prerequisites

- **Router** running OpenWrt (24.10.x or 25.x). The wizard does NOT flash
  firmware — see GL.iNet or OpenWrt docs for flashing.
- **SSH access** — port 22 open, root password set.
- **Upstream internet** — either Ethernet cable into WAN port, or WiFi
  credentials for the router to join an existing network.
- **Lightning address** — where Bitcoin payments from customers route
  (e.g. `you@walletofsatoshi.com` or a raw LNURL).

## Architecture

This is a **thin UI wrapper**. All business logic lives in
[tollgate-module-basic-go](https://github.com/net4sats/tollgate-module-basic-go).
The wizard only discovers routers, renders a deployment UI, and runs SSH
commands — no payment processing, identity derivation, or Nostr logic.

| Repo | Role |
|------|------|
| **net4sats-wizard-go** (this repo) | Laptop-side onboarding wizard |
| [configurationwizzard](https://github.com/net4sats/configurationwizzard) | Router-side admin panel + captive portal |
| [tollgate-module-basic-go](https://github.com/net4sats/tollgate-module-basic-go) | Payment backend (Cashu + Lightning) |

## License

MIT
