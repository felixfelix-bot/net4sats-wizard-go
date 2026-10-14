# net4sats wizard

[![CI](https://github.com/net4sats/net4sats-wizard-go/actions/workflows/ci.yml/badge.svg)](https://github.com/net4sats/net4sats-wizard-go/actions/workflows/ci.yml)

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

## Verify binaries

Every release asset is covered by a `SHA256SUMS` file published with the
[release](https://github.com/net4sats/net4sats-wizard-go/releases).

```sh
# from the directory holding SHA256SUMS and the archive you downloaded
sha256sum -c --ignore-missing SHA256SUMS
```

On Windows, `Get-FileHash` is the equivalent check:

```powershell
Get-FileHash .\net4sats-wizard_<version>_windows_amd64.zip -Algorithm SHA256
# the printed hash must match that file's line in SHA256SUMS
```

## Building

```sh
git clone https://github.com/net4sats/net4sats-wizard-go.git
cd net4sats-wizard-go
go build -o net4sats-wizard .
./net4sats-wizard --version
```

`--version` prints `tollgate-wizard <version>` and exits 0 without starting the
server. A plain `go build` leaves the version at `dev`; inject a local one with
the same ldflags symbol the release config uses (`.goreleaser.yaml`):

```sh
go build -ldflags "-X main.version=v0.4.0" -o net4sats-wizard .
./net4sats-wizard --version   # tollgate-wizard v0.4.0
```

Cross-compile for all platforms (no cgo — the release build sets
`CGO_ENABLED=0`):

```sh
GOOS=darwin  GOARCH=arm64 go build -o dist/net4sats-wizard-darwin-arm64 .
GOOS=darwin  GOARCH=amd64 go build -o dist/net4sats-wizard-darwin-amd64 .
GOOS=linux   GOARCH=amd64 go build -o dist/net4sats-wizard-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/net4sats-wizard-windows-amd64.exe .
```

## Testing

These are the exact commands CI runs on every pull request and on every push to
`main` (`.github/workflows/ci.yml`, two parallel jobs):

```sh
go build ./...
go vet ./...
go test ./... -race -count=1
golangci-lint run
```

`golangci-lint` is pinned to **v1.61.0** in the workflow; its configuration is
[`.golangci.yml`](.golangci.yml) at the repo root.

## Releases

Releases are cut from a tag. `.github/workflows/release.yml` then runs
[goreleaser](https://goreleaser.com) with [`.goreleaser.yaml`](.goreleaser.yaml):

```sh
git tag vX.Y.Z
git push origin vX.Y.Z
```

Each release publishes six archives — linux, windows and darwin, each on amd64
and arm64 — named `net4sats-wizard_<version>_<os>_<arch>.tar.gz` (`.zip` on
Windows), plus the `SHA256SUMS` checksum file that covers all six. The binary
inside is `net4sats-wizard` (`net4sats-wizard.exe` on Windows) and reports
itself as `tollgate-wizard <version>` from `--version`.

To reproduce a release build locally, without tagging or publishing:

```sh
goreleaser check
goreleaser release --snapshot --clean   # artifacts land in dist/
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
