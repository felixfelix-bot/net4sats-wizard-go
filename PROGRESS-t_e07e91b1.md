# PROGRESS — t_e07e91b1 (Assemble master with intro card and phase captions, export draft)

card: t_e07e91b1 (board net4sats-mvp-v2)
role: fleet offload worker
status: **DONE — draft master assembled, exported, verified 17/17**

## Inputs — what actually exists on this host (measured, not inherited)

The three declared inputs are not all present. Two of three are real; one was never produced.

| # | declared input | from | reality on this host | usable? |
|---|---|---|---|---|
| 1 | intro card (SHA-bearing frames) | t_7015c244 | `attachments/t_7015c244/intro_card.mp4`, 80,817 B, sha256 `782766e0…c9ac4c`, h264 1920x1080 30fps 4.000s | **YES** |
| 2 | cleaned base cut from the trim task | t_ba89d8b4 | **never produced**. Card is `todo`; run 1356 = `blocked` — "trim_plan.json (t_ae815553) has never completed… raw take does not exist". `artifacts/demo_basecut/` in this repo holds only `trim_plan.json` + `trim_plan.md` + a validator — all `status: BLOCKED`, `executable: false`. | **NO** |
| 3 | step log / phase descriptions | t_882b3922 | **never produced**. 3 runs, **0 completed**, 0 attachments — the recorder needs a reachable OpenWrt router over SSH and never got one. | **NO** |

So step 1 ("prepend the intro card to the base cut") has no base cut to prepend to. I did not
invent one. Instead I located real, on-host net4sats footage and derived the phase map from
frames by OCR, which is the honest substitute for inputs 2 and 3.

## Footage search — executed, with the decoys named

* `demo-trimmed.mp4` (1080p, 21.8s) in `worktrees/t_11d0873e` and `t_1eae5a7b` — **REJECTED**.
  OCR of its frames reads *"Install PWA"*, *"Pay in sats"*, *"Continue to payment"* — it is a
  **Minibits wallet** recording, not the net4sats router flow. This is the decoy the prior
  audit t_f53c5650 also flagged as belonging to a different chain.
* `take-attempt-blocked.mp4` (1080p, 42.9s) in `worktrees/wz19-video` — **REJECTED**. OCR reads
  *"Scanning for routers"* then **"Setup failed"**. It is an explicitly-named blocked attempt
  from a different card (t_8b5cf200), and it ends in failure.
* **ACCEPTED**: `/home/c08r4d0r/reports/net4sats-mvp/film5/` — durable, real net4sats footage,
  `01-installer-deploy.webm` (192s), `02-portal.webm` (30s), `03-admin-tour.webm` (62s).
  OCR confirms the net4sats wizard, the captive portal (*"Purchase Internet Access / Pay with
  Cashu to access the internet"*) and the router admin UI. These are 1280x800, so they do not
  natively meet the 1080p bar — handled explicitly below.

## Phase map — derived from frames, not from prose

Every boundary was chosen by OCR-scanning frames at 1-3s intervals and by measuring per-second
visual motion (PIL frame differencing) to satisfy the "no dead air > ~2s" bar.

| seg | phase | source | window | content verified by OCR |
|---|---|---|---|---|
| 00 | INTRO | intro card | whole 4.0s | PR #32 + head SHA (readback, see below) |
| 01 | BUILD | installer | 8–20s | "Select your router", "192.168.1.1 Use this address" |
| 02 | PROVISIONING | installer | 22–36s | "Verifying SSH access", "DISTRIB_ID='OpenWrt'…24.10.2", "Flashing OpenWrt" |
| 03 | CONFIGURATION | admin-tour | 36–52s | "WiFi / radio0 up / Channel 1", "Devices", "Wallet / Balance" |
| 04 | WORKING STATE | portal | 20–28s | "Purchase Internet Access", "Pay with Cashu", "Pay with BTC Lightning" |

Dead-air check per window (longest consecutive frame freeze): 0.5s, 1.5s, 0.0s, 0.0s — all <= 2s.

## The build — canvas contract

Naive full-frame caption burn-in would have violated the "must NOT cover the UI element being
demonstrated" rule (the admin UI uses the vertical centre; the portal puts the purchase button
mid-frame). So the canvas is a **video belt + caption belt**:

```
1920x1192 = rows 0..1079    video belt   (a true 1920x1080 frame, untouched)
          + rows 1080..1191 caption belt (112 px, caption ink lives here ONLY)
```

1280x800 (AR 1.600) footage is fitted to 1728x1080 (identity AR) and pillarboxed 96px L/R into
the video belt — full source width, no crop, constant frame size. Caption type: **38px** DejaVu
Sans Bold, uniform across all five strips (one fitted size, so type never jumps); tag 30px.
Both exceed the >=24px-at-1080p floor (38px on a 1192 frame = **34.4px** normalised to 1080p).

## Verification — 17/17 PASS, all measured from the FILE

`verify/wz110-master-t_e07e91b1/verify.json` and `occlusion.json` hold the machine-readable results.

* resolution 1920x1192 (>= 1080p), single video stream, CFR 25/25, duration 54.000s
* `ffmpeg -f null -` decoded every frame, zero errors -> plays end to end
* **intro SHA char-for-char**: prep `head_sha` read from JSON =
  `2d770966b163cd2deafb1697bfd892e0a7a8d9e2` (40 ch); OCR of the rendered intro frames returned
  exactly that string. IDENTICAL.
* phase captions read, in timeline order: INTRO -> 1/4 BUILD -> 2/4 PROVISIONING ->
  3/4 CONFIGURATION -> 4/4 WORKING STATE
* **non-occlusion proven at pixel level** (not by eye, not by a vision model): the same frames
  rendered with and without the caption strip, extracted as raw rgb24 PNGs with **no video
  codec in the path**, then differenced. Video belt (rows 0..1079): **0 differing pixels across
  20 frames**. Caption belt: mean |diff| 29.3. Caption ink is provably confined to rows >= 1080.
* dead air: footage belt 1.5s longest freeze (<= 2s PASS). The 3.5s motionless stretch is the
  **intro card**, which is static by design and was independently ruled PASS by t_069b6842;
  counting it as dead air would be a category error.
* audio: **0 audio streams — this master is silent on purpose.** The recording chain produced no
  narration take (t_882b3922, 0 completed runs), so there is no audio to encode, and the
  default-delivery contract for this demo is captioned/silent (prepass card t_74b12192).
  This is a measured property, not a claim of intelligible narration.

## Deliverable

`demo-pr32-2d770966.mp4` — H.264/mp4, yuv420p, 1920x1192, 25fps CFR, 54.0s.
See REPORT-t_e07e91b1.md for the exact path, size and sha256.

## Honest gap

Step 1 as literally written ("prepend the intro card to the cleaned base cut") could not be
performed, because the cleaned base cut does not exist on this host and never did. I substituted
the closest real, verified net4sats footage and said so. Captions are 1080p-native in
resolution but the underlying footage is 1280x800 upscaled 1.35x — the UI text is legible and
OCR-readable, but this is an upscale and is not the same as a natively-captured 1080p take.
The durable fix is still the one the chain has been blocked on: record a real 1080p take on the
operator-gated GL-MT6000.
