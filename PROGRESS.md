# PROGRESS — WZ1.9 / t_8b5cf200: Record E2E demo video of de-branded wizard flow

Board: net4sats-mvp-v2 | Host: x280 (c08r4d0r) | Started: 2026-10-23T05:15Z
Worktree: /home/c08r4d0r/worktrees/wz19-video (branch wz19-video-capture)

## Verdict up front

**BLOCKED — the card's central premise is false, and the acceptance criteria are
self-contradictory.** A real, unbroken, 1080p E2E happy-path video of a
*de-branded* wizard **cannot be produced today**. Three independent facts:

1. **No de-branded wizard exists in any branch.** Every ref in the repo still
   renders `net4sats` in the wizard UI. The card asks for a video with "NO
   net4sats branding", which no recorded run of the current code can satisfy.
2. **No happy path is reachable.** No TollGate/OpenWrt router is addressable from
   this host (LAN sweep + 4 serial consoles + SSH all negative), so every deploy
   fails at step 1 and the wizard shows the error view. Step 2 says keep the
   failure visible -> the video would be a footage of an error, not a happy path.
3. **The PR-branch instruction does not resolve.** `net4sats/tollgate-wizard-go`
   is 404. The PR is #32 on `net4sats/net4sats-wizard-go`, branch `wz14-debrand`.

I produced a genuine 1080p take of the real flow anyway (so the toolchain and the
failure are both on the record) and verified its branding by OCR.

## Finding -> status -> files touched

### F1 — No de-branded wizard exists on ANY branch (card premise false)
- status: CONFIRMED (measured, not inherited)
- Measured: for every `refs/heads` + `refs/remotes`, `git show <ref>:index.html |
  grep -ci net4sats`. Minimum across the whole repo is **4**; most refs are 5.
  **Zero refs return 0.** No exception exists.
- PR head `b64555a` (`wz14-debrand`, PR #32) still renders:
  - `index.html:6` `<title>net4sats — Router Setup</title>`
  - `index.html:212` `Deploy net4sats` (the primary button)
  - `index.html:218` `Deploying net4sats...`
  - `index.html:226` `net4sats is live!` (the success headline)
- Also confirmed on the **built binary**, not just source:
  `strings /tmp/wiz-wz14 | grep -ci net4sats` -> **120** hits, including all four
  visible UI strings above. So the branding is in the compiled artifact the
  recording would actually render.
- The deploy step list is itself branded and unavoidable: "Installing net4sats
  package + patching backend...", "Branding captive portal as net4sats...",
  "Deploying net4sats captive portal...", "Installing net4sats admin panel..."
  (captured live, see F3).
- PR #32's own body concedes this: "Residual `net4sats` strings — apk package
  name, admin-panel HTML titles/icon filenames and release URLs still carry
  branded text in this tree."
- files touched: none

### F2 — The instruction "check out the PR branch of net4sats/tollgate-wizard-go" is unresolvable
- status: CORRECTED (repo does not exist)
- `gh api repos/net4sats/tollgate-wizard-go` -> **404 Not Found**.
- Real repo: **`net4sats/net4sats-wizard-go`**; PR **#32**, state OPEN,
  `mergeable: MERGEABLE`, head `wz14-debrand` @ `b64555a`.
- Independently corroborated by the board's own Gate-8 comment on `t_0a3a77bb`
  ("this repo does not exist — use `net4sats/net4sats-wizard-go`").
- files touched: none

### F3 — No reachable device: the happy path is unreachable (BLOCKING)
- status: **BLOCKING** — this is what makes the deliverable impossible, not slow
- Measured live, this host, this session:
  - **LAN sweep** 192.168.2.0/24 on 22/2121/2051 (concurrent, positive control
    `192.168.2.1:80` -> OPEN, so the probe is valid): 6 hosts answer :22, and
    **zero answer :2121** (the TollGate backend port the deploy must reach).
  - **SSH to every candidate**: `.30` -> `Permission denied (publickey,password)`;
    `.130` -> `Permission denied (publickey)`; `.121` -> same. `.130` is not a
    router at all — it serves `Wifistr` / "Open wifi maps powered by nostr" on :80.
  - **All four serial consoles are dead**: `/dev/ttyUSB0..3`, 115200 raw, newline
    sent, 3-4s read -> all silent (no prompt). Nothing to drive.
  - Token-scan fallback (`/dev/tcp`) timed out at 180s; the Python sweep above
    replaced it and is the measurement of record.
- Consequence, observed for real: the deploy run fails at **step 1** and the
  wizard renders the error view `Cannot connect to router via SSH`. 11 steps were
  listed; step 1 never passed.
- Step 2 of the card says "if a step fails, keep the failure visible and report it
  as a blocker." I have done exactly that — but it means the artifact is a
  recording of a failure, which is the opposite of the "COMPLETE happy path" the
  same sentence demands. The card cannot be satisfied as written.
- files touched: none (probes only)

### F4 — Toolchain is VIABLE and was proven by a real recording
- status: WORKING (the one thing that is not blocked)
- Playwright 1.52.0 from the npx cache
  (`~/.npm/_npx/c61c9351a0dbcfa7/node_modules/.bin/playwright`); Node here is
  18.19.1, so the newer npx cache (requires Node >= 20) is NOT usable.
- `channel: 'chrome'` fails (no `/opt/google/chrome/chrome`); bundled chromium at
  the pinned revision also fails. **Working combination:** explicit
  `executablePath = ~/.cache/ms-playwright/chromium-1248/chrome-linux64/chrome`,
  `headless: false`, Xvfb `:97` at 1920x1080x24.
- Produced a genuine take: `/tmp/wz19-take/*.webm`, with a timestamped
  frame-level step log at `/tmp/wz19-take/step-log.txt`.
- files touched: `verify/wz19-video-capture/capture.cjs` (new, in worktree)

## What was actually recorded (honest scope)

A real, continuous, un-cut 1080p take of the wizard driven through the UI:
start -> discovery (5 routers listed) -> router select -> password -> upstream
mode -> Lightning address -> Advanced open/close -> deploy click -> 11-step
progress list -> **error view "Cannot connect to router via SSH"**. That is the
true end state of the flow on this host. It is NOT a happy path and NOT
de-branded; I am not going to label it as either.

## REMAINING steps (numbered, for whoever unblocks this)

1. **Fix the branding** (owner: the WZ1.4/de-brand lane). Strip the 4 visible
   `net4sats` strings in `index.html` (title, deploy button, deploying heading,
   success headline) plus the 4 branded deploy step labels in `deploy.go`, and
   re-verify with `strings <binary> | grep -ci net4sats` on the BUILT artifact.
   Until a ref returns 0 there, acceptance criterion 3 is unsatisfiable.
2. **Provide a reachable aarch64 device** (operator action). The pins are
   arch-locked to `aarch64_cortex-a53` (`deploy.go:43,47`) and
   `verifyArtifactPin` fails closed, so only a real aarch64 router yields green.
   Needs: bench GL-MT6000 powered + cabled (its NIC `enp0s31f6` is
   `NO-CARRIER`/DOWN today) with SSH enabled, or a documented aarch64 VM/QEMU
   lab with a working qemu-system-aarch64 (only the user-mode
   `qemu-aarch64-static` is installed here).
3. **Resolve the card's self-contradiction** before re-dispatch: "NO net4sats
   branding in any frame" and "record the COMPLETE happy path" cannot both hold
   until step 1 lands. Split into (a) de-brand commit, (b) video capture against
   the de-branded head.
4. Re-run `verify/wz19-video-capture/capture.cjs` against the de-branded,
   reachable head; convert webm -> mp4 H.264 and gif; OCR-verify absence of
   branding on extracted frames.
5. Do not post the video anywhere (card: a later task owns PR attachment).

## Files touched in the repo (worktree wz19-video)

- `verify/wz19-video-capture/capture.cjs` — new, the real capture driver
- `verify/wz19-video-capture/PROGRESS.md` — this file
- `verify/wz19-video-capture/REPORT.md` — full report
- No source file was modified. No router was flashed. No happy path was faked.
