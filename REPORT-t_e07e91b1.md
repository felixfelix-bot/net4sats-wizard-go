# REPORT — t_e07e91b1

**Card:** `t_e07e91b1` — "Assemble master with intro card and phase captions, export draft"
**Board:** net4sats-mvp-v2 · **Repo:** net4sats-wizard-go (working tree `net4sats-mvp-v2`)
**Role:** fleet offload worker · **Branch:** `worker-heavy/t_gate8-19b1b1c6`
**Commit:** `0a5ea88` — pushed; remote blob sha256 verified identical to local
**Status:** DONE — draft master assembled, exported at FINAL settings, verified 19/19 against the exported file

---

## 1. Deliverable — the one exported MP4

```
/home/c08r4d0r/repos/net4sats-mvp-v2/artifacts/demo_master/demo-pr32-2d770966.mp4
```

| property | value |
|---|---|
| name | `demo-pr32-2d770966.mp4` (scheme `demo-<pr>-<shortsha>.mp4`) |
| size | 5,284,235 bytes |
| sha256 | `d14bf9c6753dcb3dce5de86a5e3cfb4f9dcd8848b38f7a2828eb8a6dbf856b3d` |
| container | mp4, `+faststart` (moov at front) |
| video | h264, profile High, level 4.1, yuv420p |
| resolution | **1920x1192** (>= 1080p) |
| frame rate | 25/1 CFR (r_frame_rate == avg_frame_rate) |
| duration | 54.000 s |
| bitrate | 782 kbps |
| audio | **0 streams — silent by design** (see §5) |

It plays end-to-end: `ffmpeg -v error -i <file> -f null -` decoded every frame with zero errors.

## 2. Acceptance criteria — measured scorecard

| criterion | result | evidence |
|---|---|---|
| One exported MP4, path stated explicitly, plays end-to-end without errors | **PASS** | path above; full decode, `stderr` empty |
| Phase order build -> provisioning -> config -> working state visible **and** explained by captions | **PASS** | captions OCR-read off the exported file in timeline order: `1/4 BUILD` -> `2/4 PROVISIONING` -> `3/4 CONFIGURATION` -> `4/4 WORKING STATE` |
| Captions legible and non-occluding | **PASS** | 38px bold (= 34.4px normalised to 1080p, floor 24px); occlusion proven at pixel level — see §4 |
| Framing stable | **PASS** | constant frame size across sampled frames; source pillarboxed with symmetric 96px bars, no letterbox jump (motion deviates steadily) |
| Resolution >= 1080p | **PASS** | 1920x1192 |
| SHA on the intro card matches the prep-task head SHA character-for-character | **PASS** | see §3 |

**Total: 19/19 PASS** on the exported artifact (`artifacts/demo_master/final-verify.json`).

## 3. The SHA — both values stated

Prep-task head SHA, read programmatically from the frozen metadata JSON
(`attachments/t_0143bb0e/intro_card_meta.json`, corroborated by `t_f53c5650`'s recovered artifact
and by `t_0143bb0e` run 1409's live `gh pr view 32 --json headRefOid`):

```
2d770966b163cd2deafb1697bfd892e0a7a8d9e2
```

SHA read back off the rendered intro-card frames **inside the exported master**, by OCR
(`tesseract --psm 6`, 8 frames across the 4s card):

```
2d770966b163cd2deafb1697bfd892e0a7a8d9e2
```

Character-for-character identical: 40 chars, all hex, zero divergence. `short_sha` = `2d770966`,
computed programmatically, never typed.

I did **not** re-render or re-type the SHA. The intro-card clip was copied byte-for-byte from
`attachments/t_7015c244/intro_card.mp4` (sha256 `782766e049c927905257fdbc4e27062351bd1ed3700f15e533fbcd4db3c9ac4c`)
and scaled 1920x1080 -> 1920x1080, i.e. an **identity transform**, so its SHA-bearing pixels are
preserved 1:1.

## 4. Non-occlusion — the criterion I treated as the hard one

Burning captions into the frame would have violated "must NOT cover the UI element being
demonstrated": the admin UI uses the vertical centre and the portal puts its purchase control
mid-frame. So the canvas is a **video belt plus caption belt**, not an overlay:

```
1920x1192 = rows 0..1079     video belt    (a true 1920x1080 frame, untouched)
          + rows 1080..1191  caption belt  (112 px; ALL caption ink lives here)
```

Because "the caption doesn't overlap the UI" is easy to assert and hard to prove, I proved it at
pixel level rather than by eye. I rendered the **same** frames twice — once with the real caption
strip, once with a blank strip of identical height — extracted them as raw rgb24 PNGs with **no
video codec in the diff path** (encoder rate-control noise otherwise swamps the signal; my first
attempt measured 0.30 mean-diff purely from x264 noise at CRF 18), then differenced:

```
frames compared : 20 (5 fps over a 4.0s window of the BUILD segment)
VIDEO BELT      rows 0..1079       worst mean |diff| = 0.000000  (differing pixels = 0)
CAPTION BELT    rows 1080..1191     worst mean |diff| = 29.310812
VERDICT: PASS
```

Zero differing pixels in the video belt across 20 frames. The caption ink provably cannot touch a
single video pixel. Machine-readable: `verify/wz110-master-t_e07e91b1/occlusion.json`.

Note on vision: `vision_analyze` returned HTTP 503 (`no lane declares this model (flat router)`),
so I used OCR readback and pixel differencing throughout. Per the intro-card skill these are the
stronger evidence anyway for "exact characters present" and "nothing moved"; no conclusion in this
report rests on a vision model.

## 5. Phase map — derived from frames, not from prose

Boundaries were chosen by OCR-scanning frames at 1-3 s intervals **and** by measuring per-second
visual motion (PIL frame differencing) so each window satisfies the dead-air bar.

| seg | phase | source | window | content verified by OCR |
|---|---|---|---|---|
| 00 | INTRO | intro card | whole 4.0 s | `PR #32`, `fix(wz1.4): de-brand wizard business logic…`, head SHA |
| 01 | BUILD | installer | 8–20 s | "Select your router", "192.168.1.1 Use this address" |
| 02 | PROVISIONING | installer | 22–36 s | "Verifying SSH access", "DISTRIB_ID='OpenWrt' … 24.10.2", "Flashing OpenWrt", "installing tollgate-wrt package + patching backend" |
| 03 | CONFIGURATION | admin-tour | 36–52 s | "WiFi / radio0 up / Channel 1", "Devices – 0 connected", "Wallet / Balance" |
| 04 | WORKING STATE | portal | 20–28 s | "Purchase Internet Access", "Pay with Cashu to access the internet", "Pay with BTC Lightning to access the internet" |

Dead air per window (longest consecutive visual freeze): **0.5 s, 1.5 s, 0.0 s, 0.0 s** — all <= 2 s.
Whole-master footage belt: longest freeze **1.5 s**.

The 3.5 s motionless stretch in the master is the **intro card** (t=0..3.5 s). A static title card
is motionless by design — its own independent verifier `t_069b6842` ruled it a PASS on exactly
that property — so counting it as "dead air" would be a category error. I report both numbers
rather than silently excluding it.

**Audio:** 0 audio streams, deliberately. The recording chain (t_882b3922) has 3 runs and
**0 completed**, produced no narration take, and the audio stream of the source webm files does
not exist (ffprobe: `vp8,video` only). There is therefore no narration to encode. The task asks
for "narration **or** on-screen captions"; I shipped captions, and the default-delivery contract
for this demo is captioned/silent (prepass card t_74b12192). I state this as measured absence,
**not** as "audio present and intelligible" — that clause is satisfied by the caption branch, not
by narration.

## 6. Input honesty — two of the three declared inputs never existed

This is the part a reader should not skim. The card's INPUTS list three upstream artifacts.
Measured against the board DB and the filesystem:

| declared input | from card | reality | usable |
|---|---|---|---|
| intro card (SHA-bearing frames) | `t_7015c244` | present — 80,817 B, h264 1920x1080 30fps 4.000 s | **YES** |
| cleaned base cut (trim task) | `t_ba89d8b4` | **never produced.** Card `todo`; run 1356 `blocked`: "trim_plan.json (t_ae815553) has never completed … the raw take it must reference does not exist on any reachable host". `artifacts/demo_basecut/` holds only `trim_plan.json`/`.md`/validator, all `status: BLOCKED`, `executable: false` | **NO** |
| step log / phase descriptions | `t_882b3922` | **never produced.** 3 runs, 0 completed, 0 attachments — the recorder needs a reachable OpenWrt router over SSH and never got one | **NO** |

So **step 1 as literally written — "prepend the intro card to the base cut" — was not performed,
because there is no base cut to prepend to.** I did not invent one, and I did not silently
substitute a lookalike.

What I did instead: searched the host for real net4sats footage, checked every candidate by OCR,
and derived the phase map from the accepted footage's frames.

**Decoys found and rejected (named, because rejecting them is the substantive work here):**

* `demo-trimmed.mp4` (1080p, 21.8 s, in `worktrees/t_11d0873e` and `t_1eae5a7b`) — **REJECTED.**
  OCR reads *"Install PWA"*, *"Pay in sats"*, *"Continue to payment"*: it is a **Minibits wallet**
  recording, not the net4sats router flow. The prior audit `t_f53c5650` flagged this same file as
  belonging to a different chain. It would have satisfied "1080p" mechanically while showing the
  wrong product.
* `take-attempt-blocked.mp4` (1080p, 42.9 s, `worktrees/wz19-video`) — **REJECTED.** OCR reads
  *"Scanning for routers on your network"* then **"Setup failed"**. It is an explicitly-named
  blocked attempt from a different card (`t_8b5cf200`) and it ends in failure.
* **ACCEPTED:** `/home/c08r4d0r/reports/net4sats-mvp/film5/` — durable on-host net4sats captures:
  `01-installer-deploy.webm` (192 s), `02-portal.webm` (30 s), `03-admin-tour.webm` (62 s).
  OCR confirms the net4sats wizard, the captive portal, and the router admin UI.

## 7. Honest gap — read this before calling the artifact "1080p native"

The footage above is **1280x800**. It is fitted with identity AR to 1728x1080 and pillarboxed
96 px L/R into the 1920x1080 video belt, so the **canvas** is genuinely 1920x1080+ and no pixels
are cropped — but the underlying image is a **1.35x upscale**, not a natively-captured 1080p take.
UI text is legible and survives OCR readback; that is a real, measured property. It is not the
same claim as "recorded at 1080p".

The durable fix is unchanged and is the blocker the chain has been stuck on for weeks: record a
real 1080p take on the operator-gated GL-MT6000 with a reachable SSH endpoint, then re-run the
trim task and re-assemble against that base cut. Everything in this build is scripted and reruns
against a real take by swapping the source paths in `build_master.py`.

## 8. Reproduction — all inputs committed

Everything needed to rebuild and re-verify is in the repo:

* `artifacts/demo_master/` — the MP4, `master.json`, `final-verify.json`, and the four builder scripts
* `verify/wz110-master-t_e07e91b1/` — `build_master.py`, `make_strips.py`, `verify.py`,
  `occlusion_proof.py`, `motion.py`, `deadair.py`, `locate.py`, plus the JSON result files
  (`verify.json`, `occlusion.json`, `motion.json`, `deadair.json`)

## 9. Commit / push

* **commit** `0a5ea88` — "feat(demo): t_e07e91b1 — assemble draft master (intro card + 4 phase
  captions), export 1080p"
* **push** to `origin` (`github.com/felixfelix-bot/net4sats-wizard-go`),
  `worker-heavy/t_gate8-19b1b1c6`: `95ae8a6..0a5ea88`
* **push verified, not assumed:** `git ls-remote` shows the remote ref at `0a5ea88e97d63f3a16ebbd02add0848f1b7e3fa0`
  == local HEAD, and hashing the blob retrieved from the remote ref
  (`git cat-file blob origin/...:artifacts/demo_master/demo-pr32-2d770966.mp4 | sha256sum`) returns
  `d14bf9c6…6b3d`, byte-identical to the local artifact. The remote git blob hash
  (`2aa8e15ba4dfd5b82d893b7b49cc8f3c49a274ad`) matches on both sides.
* pre-commit credential gate: clean (30 literals checked, mode=staged)

## 10. Criteria I am NOT claiming

* **Not claiming** narration: there is no audio stream, and I did not synthesise one. The caption
  branch of the "narration or captions" requirement is what is satisfied.
* **Not claiming** a natively-1080p take: the canvas is 1080p+, the footage is a 1.35x upscale (§7).
* **Not claiming** the base cut was used: it does not exist (§6).
* **Not claiming** the master is a final publishable demo: the card says *draft*, and the upscale
  and the missing real take are open.
