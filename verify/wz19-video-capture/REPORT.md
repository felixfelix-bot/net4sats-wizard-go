# REPORT — WZ1.9 / t_8b5cf200: Record E2E demo video of de-branded wizard flow

Worker pass, 2026-10-23. Host x280 (c08r4d0r).
Worktree: `/home/c08r4d0r/worktrees/wz19-video`, branch `wz19-video-capture` @ `0bb28d5` (pushed).
Evidence dir: `verify/wz19-video-capture/`.

---

## VERDICT

**BLOCKED. The requested deliverable cannot be produced, and the reason is not
tooling — two of the card's own premises are false, and its acceptance criteria
contradict each other.**

I did not produce a de-branded happy-path video, because no such thing is
producible today. I did produce a **real, un-cut, 1920x1080 H.264 recording of the
actual wizard flow** so that the failure is on the record rather than narrated, and
I OCR-verified its branding frame by frame. That audit **fails**: 83 of 86 frames
contain `net4sats`.

Three independent blockers, each measured first-hand:

| # | Blocker | Measured how |
|---|---|---|
| 1 | **No de-branded wizard exists on any branch** | Branding grep across every ref + `strings` on the built binary |
| 2 | **No reachable device, so no happy path** | LAN sweep w/ positive control, SSH to all candidates, 4 serial consoles |
| 3 | **The "PR branch" in step 1 does not exist** | `gh api repos/net4sats/tollgate-wizard-go` -> 404 |

### The card is self-contradictory as written

Step 2: *"Record the COMPLETE happy path ... No cuts that hide a failure; if a step
fails, keep the failure visible and report it as a blocker."*
Step 4 + criterion 3: *"confirm every branded string/logo is absent from the frames."*

Given blocker 1, a recording of the current code cannot satisfy step 4. Given
blocker 2, the only recording obtainable is of a failure, which cannot be "the
COMPLETE happy path". **Both halves cannot be satisfied simultaneously by any
artifact, so this card needs re-scoping, not another worker.** I have honoured the
"keep the failure visible and report the blocker" half — the take contains the real
failure, uncut.

---

## ACCEPTANCE CRITERIA — reported as demanded, including the ones that fail

### 1. Absolute path(s) + file size + duration

| Artifact | Path | Size | Duration | Format |
|---|---|---|---|---|
| Primary take (real, un-cut) | `/home/c08r4d0r/worktrees/wz19-video/verify/wz19-video-capture/artifacts/take-attempt-blocked.mp4` | 647,139 B (632 KiB) | **42.88 s** | H.264, 1920x1080, 25 fps, 1072 frames, yuv420p |
| Raw capture (pre-convert) | `/tmp/wz19-probe/f78874eab2b35a98f9d26b57737478f8.webm` | 2,943,416 B | 42.88 s | VP8, 1920x1080 |
| Step log (frame-timestamped) | `.../artifacts/step-log.txt` | 1,784 B | — | text |
| Branding audit | `.../BRANDING-AUDIT.md` | 3,179 B | — | text |
| First-frame OCR | `.../artifacts/ocr-f1.txt` | 103 B | — | text |
| Last-frame OCR | `.../artifacts/ocr-f86.txt` | 107 B | — | text |
| Capture driver | `.../capture.cjs` | 6,799 B | — | Node/Playwright |

The git-committed mp4 is 632 KiB because git-committed evidence is kept small; the
full-size master is the webm above. Only the webm->mp4 conversion was applied
(`libx264`, `-crf 20`); no cut, trim, or frame was dropped.

Note the filename: **`take-attempt-blocked`**, deliberately. It is not named as a
gate-satisfying artifact because it does not satisfy the gate.

### 2. One-paragraph, step-by-step description of what the video shows

The take opens on the wizard's start screen (branded header `net4sats` /
"Bitcoin WiFi — Setup your router in 60 seconds") with the status line "Scanning
for routers on your network...". After ~22 s of discovery the router dropdown
populates with five hosts found on 192.168.2.0/24 (`.1`, `.2`, `.30`, `.121`,
`.130`) — all reported as "unknown unknown" vendor/model, because none is an
OpenWrt target. The operator selects `192.168.2.30`, leaves the password field
empty, keeps the default "WAN (Ethernet)" upstream mode, types a Lightning address
(`demo@getalby.com`), and opens then closes the "Advanced settings (optional)"
`<details>` panel to show it is not required — the deploy button enables
(`{"disabled":false,"text":"Deploy net4sats"}`). Clicking **Deploy net4sats**
switches to the deploy view and renders the full 11-step progress list (Verifying
SSH access, Checking firmware version, Setting root password, Configuring upstream
connection, Installing net4sats package + patching backend, Branding captive portal
as net4sats, Deploying net4sats captive portal, Installing net4sats admin panel,
Configuring Lightning address, Restarting services, Running health check) — all
still in the pending state. Within ~1 s the flow terminates at **step 1** and the
error view is shown: **"Setup failed / Cannot connect to router via SSH"**. That
error frame is held for the final ~4 s. The video ends there. No success state is
reached; the recording is of the true end-to-end behaviour on this host.

### 3. Explicit confirmation re: net4sats branding — **CANNOT CONFIRM. BRANDING IS PRESENT.**

**I cannot give the requested confirmation, because the check fails.** Stating
otherwise would be false.

How I checked (this is the check the card asks for, performed on rendered pixels,
not on source):

```
$ ffmpeg -v error -i take-attempt-blocked.mp4 -vf fps=2 /tmp/wz19-ocr/f%04d.png
$ ls /tmp/wz19-ocr/*.png | wc -l
86
$ cd /tmp/wz19-ocr && for f in *.png; do tesseract "$f" - ; done
   -> normalise whitespace/case, search for token "net4sats"
```

| Metric | Result |
|---|---|
| Frames OCR'd (2 fps over 42.88 s) | 86 |
| Frames containing `net4sats` | **83** |
| Frames clean | 3 |
| **Verdict** | **FAIL — acceptance criterion not met** |

The wordmark renders in the page header on **every** screen (start, config,
deploy, error), so it cannot be cropped out of a full-window recording without
cutting the UI. Confirmed in the raw OCR:

```
f0001: net4sats / Bitcoin WiFi — Setup your router in 60 seconds / Scanning for routers...
f0086: net4sats / Bitcoin WiFi — Setup your router in 60 seconds / Setup failed
       / Cannot connect to router via SSH
```

Binary cross-check (rules out source/binary disagreement):
`strings /tmp/wiz-wz14 | grep -ci net4sats` -> **120**.

Repo-wide check — **no clean ref exists**, so no de-branded recording could have
been taken from anywhere:

```
$ for ref in $(git for-each-ref --format='%(refname:short)' refs/heads refs/remotes); do
    echo "$(git show "$ref:index.html" | grep -ci net4sats)  $ref"; done | sort -n | head
4  origin/pr32
4  origin/worker-heavy/t_c92a4003
4  origin/wz14-debrand
4  origin/wz19-video-evidence
4  origin/wz19-video-evidence-pre-atomic
4  wz19-video-evidence
5  ci/add-github-actions-workflow-clean
5  ci/github-actions-workflow-atomic
```

Minimum across the entire repository is 4; **zero refs return 0.**

### 4. Not posted anywhere — honoured

The video was **not** uploaded, embedded, attached, or linked in any PR, comment,
release, or external host. It exists only on local disk and on the
`wz19-video-capture` branch of the fork. A later task owns PR attachment.

---

## Evidence for each blocker

### Blocker 1 — no de-branded ref exists

On PR head `b64555a` (`wz14-debrand`, PR #32), `index.html` still contains four
visible branded strings:

```
6:  <title>net4sats — Router Setup</title>
212: <button ... >Deploy net4sats</button>
218: <h3 ...>Deploying net4sats...</h3>
226: <h2>net4sats is live!</h2>
```

The deploy step list is branded too and is rendered during any deploy — captured
live in this take, and present in `deploy.go` as:
`Installing net4sats package + patching backend...`,
`Branding captive portal as net4sats...`,
`Deploying net4sats captive portal...`,
`Installing net4sats admin panel...`.

PR #32's own body already concedes this: *"Residual `net4sats` strings — apk
package name, admin-panel HTML titles/icon filenames and release URLs still carry
branded text in this tree."*

A sibling lane's older work (`66756b1`, "debrand: strip internal names...") is
**not** what this needs — it strips *internal* tokens (gate labels, personal
subnets, board paths), and it explicitly left `index.html` branded; its own diff
touches `index.html` only to repin a Minibits mint URL.

### Blocker 2 — no happy path is reachable

- **No TollGate backend anywhere.** Concurrent sweep of 192.168.2.0/24 on ports
  22/2121/2051: 6 hosts answer :22; **zero answer :2121** (the backend port the
  deploy's health check and install path require). Positive control
  `192.168.2.1:80` -> OPEN, so the probe is valid.
- **Every SSH candidate refuses.** `.30` -> `Permission denied (publickey,password)`;
  `.130` and `.121` -> `Permission denied (publickey)` with `BatchMode` as `root`.
- **`.130` is not a router at all** — it serves `Wifistr`, "Open wifi maps powered
  by nostr", on :80 with a LuCI-ish :8080 that 303s to `/login`.
- **All four serial consoles are dead.** `/dev/ttyUSB0..3`, 115200 raw, newline
  sent, 3-4 s read -> silent, no prompt. Nothing to drive.
- **Arch pin makes a stand-in insufficient.** `deploy.go:43,47` pin the artifact
  to `aarch64_cortex-a53` and `verifyArtifactPin` fails closed (deploy.go:76-93:
  a URL with no pin is refused exactly like a hash mismatch). Only a real aarch64
  router yields a green run. Only user-mode `qemu-aarch64-static` is installed —
  there is no `qemu-system-aarch64` for a full-system router VM.
- **Observed end state, twice, for real:** 11 steps listed, step 1 fails, wizard
  shows `Cannot connect to router via SSH`.

### Blocker 3 — the named repo does not exist

`gh api repos/net4sats/tollgate-wizard-go` -> `404 Not Found`. The real target is
**`net4sats/net4sats-wizard-go`**, PR **#32**, head `wz14-debrand`, state OPEN,
`mergeable: MERGEABLE`. This matches the board's own Gate-8 comment on
`t_0a3a77bb` ("this repo does not exist — use `net4sats/net4sats-wizard-go`"), which
I re-verified live rather than inheriting.

---

## Toolchain status — this part WORKS (and is reusable)

The capture path is not the problem and should not be rebuilt from scratch:

- Playwright **1.52.0** from the npx cache:
  `~/.npm/_npx/c61c9351a0dbcfa7/node_modules/.bin/playwright`. Node here is
  **18.19.1**, so the newer npx cache (`e41f203b7505f1fb`) is unusable — it needs
  Node >= 20.
- `channel: 'chrome'` **fails** (no `/opt/google/chrome/chrome`); the pinned
  bundled chromium revision is **missing**. Do not use either.
- **Working combination:** explicit
  `executablePath: ~/.cache/ms-playwright/chromium-1248/chrome-linux64/chrome`,
  `headless: false`, under `Xvfb :97 -screen 0 1920x1080x24`, with
  `recordVideo.size = 1920x1080`.
- Produced a valid H.264 1080p take first try. `capture.cjs` drives the wizard
  **through the UI** (real clicks/fills, not API calls) and writes a
  frame-timestamped `step-log.txt`, which is what makes the take auditable.

Re-running it against a de-branded, reachable head is a config change, not new work.

---

## REMAINING steps (numbered — who must do what)

1. **Strip the branding** (owner: WZ1.4 / de-brand lane). Remove the 4 visible
   `net4sats` strings in `index.html` (title, deploy button, deploying heading,
   success headline) **and** the 4 branded deploy step labels in `deploy.go`, then
   verify on the **built binary**: `strings <binary> | grep -ci net4sats` must be 0
   (it is 120 today). Acceptance criterion 3 is unsatisfiable until then.
2. **Provide a reachable aarch64 device** (operator action; cannot be done from any
   agent lane). Either power on + cable the bench GL-MT6000 (its NIC `enp0s31f6` is
   `NO-CARRIER`/DOWN right now) with SSH enabled, or stand up an aarch64 full-system
   QEMU lab and document it.
3. **Re-scope this card** before re-dispatch. "NO branding in any frame" and
   "record the COMPLETE happy path" cannot both hold while (1) and (2) are false.
   Split into: (a) de-brand commit, (b) video capture against the de-branded head.
4. Re-run `verify/wz19-video-capture/capture.cjs` against the de-branded, reachable
   head; convert webm -> mp4 (H.264) + gif fallback; re-run the OCR audit
   (`ffmpeg -vf fps=2` -> `tesseract` per frame) and require **0** branding hits.
5. Evidence naming must stay honest (`take-*` not `gate-*`) until step 4 passes.
6. Still do **not** post the video anywhere — attachment remains a later task.

---

## Files touched

Worktree `/home/c08r4d0r/worktrees/wz19-video`, branch `wz19-video-capture`,
commit `0bb28d5` (pushed; remote sha verified equal to local):

- `verify/wz19-video-capture/capture.cjs` — new: real Playwright UI driver
- `verify/wz19-video-capture/artifacts/take-attempt-blocked.mp4` — new: the real take
- `verify/wz19-video-capture/artifacts/step-log.txt` — new
- `verify/wz19-video-capture/artifacts/ocr-f1.txt`, `ocr-f86.txt` — new
- `verify/wz19-video-capture/BRANDING-AUDIT.md` — new
- `verify/wz19-video-capture/PROGRESS.md` — new
- `verify/wz19-video-capture/REPORT.md` — this file

**No source file was modified. No router was flashed. Nothing was uploaded. No
happy path was faked, and no synthetic or mockup footage was produced.**

## Verification evidence (this lane's own work)

Full detail in `verify/wz19-video-capture/VERIFICATION.md`; raw JSON in `artifacts/`.

- **`hermes verify` -> `ok: true`** (build exit 0, `go test ./...` exit 0, readiness
  HTTP 200 on `http://127.0.0.1:8099/`). The first run reported `ok: false` purely
  because the detector guessed port **8000** while the wizard binds **8099**; the
  app's own stdout in that same record says `running on http://127.0.0.1:8099`.
  Confirmed by re-running `--phase start --port 8099` -> `ok: true`, HTTP 200, and
  then fixing the recipe. The failing run is kept at
  `artifacts/verify-first-run-port-mismatch.json` rather than deleted.
- **Recipe pinned:** `.hermes/environment.json` now sets `port: 8099` with the reason,
  so the next lane does not rediscover it.
- **Cache defeated:** `go test -count=1 ./...` -> `ok net4sats-wizard 57.082s`, exit 0
  (`artifacts/gotest-count1.log`).
- **`node --check capture.cjs` -> SYNTAX OK.**
- **Reproducibility:** the take was re-run on a second clean Xvfb display (`:98`).
  Same 5 routers in the same order, same 11-step list, same `Deploy net4sats` ready
  state, same terminal error view, and branding in **78/81 (96.3 %)** frames vs
  **83/86 (96.5 %)** in run 1 — last-frame OCR text byte-identical.

**What this proves:** build + full suite green, app boots and serves, capture
toolchain valid and deterministic.
**What it does NOT prove, and I am not claiming:** that the de-branded happy path
works. A green build/test/readiness run says nothing about branding or about a
deploy that dies at step 1. The gate stays open.

## Fleet writeback

- **Comment posted** to card `t_8b5cf200` on board `net4sats-mvp-v2`
  (`hermes kanban --board net4sats-mvp-v2 comment t_8b5cf200 ...`) — confirmed by
  the CLI: `Comment added to t_8b5cf200`.
- **`hermes kanban block` was REFUSED**: `cannot block t_8b5cf200`. The card sits in
  `triage` and is not claimed by this lane, so the block transition is not permitted
  from here. I did **not** force it. The blocker therefore lives in the card comment
  and in this report; a card owner or the orchestrator must apply the status change.
  Card status at hand-off remains `triage` (verified by `hermes kanban show`).
- No card status, assignee, or gate field was altered by this lane.

## Explicitly NOT claimed

- No gate PASS. Criterion 3 fails and criterion 2 is unmet.
- No PR, no comment, no video attachment, no upload.
- No merge, no flash, no deploy success.
- No card status transition — the `block` call was refused and reported as such.
- A prior sibling lane's "WZ1.9 demo video" (`docs/demo/demo-trimmed.mp4` on
  `wz19-video-evidence`) is **not** the wizard: OCR of it shows the **Cashu.me**
  wallet onboarding and an SMS mockup containing the literal text *"mockup invoice /
  placeholder - no real"*. It must not be counted as satisfying this gate.
