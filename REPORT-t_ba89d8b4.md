# REPORT — t_ba89d8b4: Encode base cut (head/tail trim only, constant 1080p+ frame size)

**Card:** `t_ba89d8b4` · **Board:** `net4sats-mvp-v2` · **Worker:** worker-heavy (fleet offload)
**Host:** c08r4d0r / CobradorWave · **Date:** 2026-10-23
**Worktree:** `/home/c08r4d0r/repos/net4sats-mvp-v2` (linked worktree; git dir
`/home/c08r4d0r/wz16-verify/.git/worktrees/net4sats-mvp-v2`)

---

## VERDICT: **BLOCKED — no encode performed, and none is possible from this chain's inputs**

**Deliverable path: NONE. No base-cut MP4 was produced.** I am stating that explicitly because
this card's deliverable is a single video file, and there is no honest way to point at one.

The trim plan exists, is committed, and is **deliberately non-executable**. Its own
`handoff_to_encoder` block names *this very card* as the downstream consumer and instructs:
**"DO NOT RUN."** This card's own INPUTS clause independently governs the same way: *"If the
plan is missing, incomplete, or its source paths do not resolve, STOP and ask rather than
improvising cut points."* The plan's source path does not resolve, so I stopped at STEP 1.

No `ffmpeg` encode was run. No trim point was invented. No interior edit was made. Nothing was
fabricated to fill the deliverable slot.

---

## 1. What the plan actually says (read, not assumed)

`artifacts/demo_basecut/trim_plan.json` — authored by card `t_ae815553`, committed **a38acbf**
(`docs(demo): t_ae815553 — fail-closed trim_plan; upstream inputs missing`).

| Field | Value |
|---|---|
| `status` | `BLOCKED` |
| `executable` | `false` |
| `stop_reason.code` | `UPSTREAM_INPUTS_MISSING` |
| `source.path` | `null` |
| `source.path_status` | `DOES_NOT_EXIST` |
| `cuts.head_trim_in_s` | `null` |
| `cuts.tail_trim_out_s` | `null` |
| `cuts.status` | `NOT_COMPUTABLE` |
| `cuts.encode_must_not_proceed` | `true` |
| `step_map.ordered_steps` | `[]` (0 steps) |
| `flagged_interior_dead_air.regions` | `[]` (0) |
| `handoff_to_encoder.downstream_card` | **`t_ba89d8b4`** (this card) |
| `handoff_to_encoder.instruction` | "DO NOT RUN. This plan is not executable." |

I re-ran the plan's own validator as independent confirmation:

```
$ python3 artifacts/demo_basecut/validate_trim_plan.py
status=BLOCKED  executable=False  steps=0  regions=0  detections=2
  ... 21 assertions ...
21/21 assertions passed
RESULT: PASS — plan is coherent and provably non-executable (fail-closed).
EXIT=0
```

So the plan is not corrupt or half-written; it is a correct, fail-closed artifact. The absence
of cut points is its intended content, not a defect I should route around.

---

## 2. Upstream state, verified live (not from the plan's prose)

Read from a read-only copy of the live board DB
(`~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db` → `/tmp/kb_ro.db`), joined across
`tasks` / `task_runs` / `task_attachments` / `task_comments`:

| Card | Title | status | runs | completed | attachments |
|---|---|---|---|---|---|
| `t_acee3752` | Prepare and validate 1080p capture toolchain | `triage` | 0 | 0 | 0 |
| `t_882b3922` | Record full E2E de-branded wizard happy-path video | `triage` | 3 | **0** | **0** |
| `t_15057569` | Probe raw demo footage: streams, duration, CFR/VFR | `triage` | 0 | 0 | 0 |
| `t_6d731c04` | Build timestamped step map, frame-verified | `triage` | 0 | 0 | 0 |
| `t_be3a9ae1` | Scan full runtime for dead air / framing shifts | `triage` | 0 | 0 | 0 |
| `t_ae815553` | Assemble trim_plan.json + trim_plan.md | `todo` | 2 | 0 | 0 |
| `t_54b0d370` | QC the base cut | `todo` | 3 | 0 | 0 |
| `t_ba89d8b4` | **Encode base cut (this card)** | `todo` | 3 | 0 | 0 |

The chain is severed at its source: **no recording run completed → no footage → no probe → no
step map → no scan → no cut points.** Every downstream card, including the QC card, is stuck on
the same missing take. This card is the 4th run of `t_ba89d8b4` and the 2nd to reach this
conclusion; the earlier run's comment (2026-10-12) said the same thing, and the QC card
independently logged it on 2026-10-11. Persistent chain-level stall, not a transient miss.

---

## 3. The source really is gone — measured, not assumed

The FD9-era plan that *does* carry real cut points references `/tmp/demo-video.webm`. I probed
for it exhaustively before concluding anything:

| Probe | Result |
|---|---|
| `ls /tmp/demo-video.webm` | ENOENT (also `/tmp/demo-meta.json`) |
| `find / -xdev -type f -size 1353893c` (exact byte size of that source) | **no results** |
| `find / -xdev -iname 'demo-video*' -o -iname 'demo*.webm'` | **no results** |
| Video attachments on the entire `net4sats-mvp-v2` board | exactly **one**: `t_7015c244/intro_card.mp4`, 1920x1080, **4.0s** — a synthetic bookend with no wizard steps |
| Every ≥1920-wide video under `/data` + `~/worktrees` | only `neode-ui/public/assets/video/video-intro.mp4` promo clips (unrelated project) |

And `/tmp` is tmpfs — the source was demonstrably wiped between runs, not merely misplaced.
Even if I had found it, it is a *different* demo chain (Cashu SMS wallet UI), not the
build/provisioning/configuration/working-state wizard happy path this card's acceptance criteria
name.

---

## 4. The decoy file, and why I did not ship it

Two sibling worktrees contain `artifacts/demo_basecut/demo-trimmed.mp4` — superficially exactly
the deliverable. It **must not** be adopted, and I did not adopt it. Measured facts:

- **It is a derivative, not raw footage.** sha256 `83bdc9117c96467008e6d7ee1abd662f5c6750c56a90f251a5264f5368f91aa7`,
  byte-identical across all four copies on disk and identical to the `output.sha256` declared by
  the FD9 plan. Re-trimming it head/tail would **double-cut** an already-trimmed master.
- **Wrong chain.** Source was `/tmp/demo-video.webm` (Cashu SMS wallet), not the wizard.
- **It fails the four-phase criterion regardless.** Its own take flips to "Paying..." around
  source 16.5s and **never resolves through EOF at 139.48s**. There is no "resulting working
  state" phase in it at all. Its plan records this as defect `take-incomplete` (severity high):
  *"Master preserves all take content but is NOT a usable demo."*
- For completeness I re-verified its technical shape: `ffprobe` → h264 High, 1920x1080,
  `r_frame_rate == avg_frame_rate == 25/1`, yuv420p, 546 frames, 21.840s; per-frame
  width/height across all 546 frames → **one distinct value, 1920x1080** (constant frame size
  confirmed, 0 audio streams). First output frame mean luminance 3.2 (dark pre-boot, matches its
  plan's `head_trim_out` 0.68s landing on first non-black frame); last frame 101.08 (content).

Shippping this file would have satisfied the *mechanical* acceptance rows (constant ≥1080p,
CFR, duration) while silently failing the substantive ones (four phases in order; it is not this
chain's footage at all). That is exactly the failure mode a base-cut pass exists to prevent, so
it is reported here as a decoy rather than presented as the deliverable.

---

## 5. Steps 2–5 of the card: not reached

| Step | State |
|---|---|
| 1. Read plan, take trim points verbatim | **DONE → STOPPED.** Points are `null`; nothing to take. |
| 2. Determine scale/crop, CFR output | NOT REACHED — depends on a measured source |
| 3. Single-pass re-encode | **NOT RUN** — deliberately; zero ffmpeg encode invoked |
| 4. Carry interior dead air forward to QC | NOT APPLICABLE — 0 regions inherited (never scanned ≠ no dead air) |
| 5. Confirm with ffprobe | NOT APPLICABLE — no output file exists |

Tooling was **not** the blocker: `ffmpeg` and `ffprobe` **6.1.1-3ubuntu5** are installed and
were exercised for real (candidate probing, per-frame dimension enumeration, frame luminance
measurement). Vision was unavailable (`vision_analyze` → 503, *"no lane declares this model
(flat model router), capacity_exhausted"*), so **no frame was visually inspected** in this pass;
every claim above is from ffprobe/ffmpeg pixel statistics and the board DB, not from eyes on
pixels. Flagging that limit rather than papering over it.

---

## 6. What would unblock this card

Numbered, in dependency order — all upstream of me:

1. `t_acee3752` — validate a 1080p capture toolchain.
2. `t_882b3922` — record the de-branded happy-path take; publish an **absolute** take path and a
   **durable** step-log path (not `/tmp` — tmpfs is why the last source evaporated).
3. `t_15057569` — write the probe report with every field in
   `source.required_before_encode_can_run` (`path, container, duration_s, width, height,
   r_frame_rate, avg_frame_rate, pix_fmt, audio_streams, cfr_vfr_verdict, sha256`).
4. `t_6d731c04` — write the frame-verified step map (start, end, confidence, evidence per step),
   covering **all four phases**: build, provisioning, configuration, resulting working state.
5. `t_be3a9ae1` — write the interior dead-air region list plus framing anomalies.
6. Re-dispatch `t_ae815553` to regenerate an **executable** plan with real cut points.
7. Re-run this card (`t_ba89d8b4`) against that executable plan.

Until step 6 lands, every re-run of this card terminates at STEP 1 by design. Re-running it
without new upstream inputs will keep producing BLOCKED — that is the correct behaviour, not a
dispatch failure to work around.

---

## 7. Files touched

| Path | Nature |
|---|---|
| `artifacts/demo_basecut/trim_plan.json` | **read only** (pre-existing, committed a38acbf) |
| `artifacts/demo_basecut/validate_trim_plan.py` | **executed** read-only (21/21 PASS) |
| `PROGRESS-t_ba89d8b4.md` | new (this run) |
| `REPORT-t_ba89d8b4.md` | new (this run) |
| `PROGRESS.md`, `REPORT.md` | prepend-only; both are **local-only** per `.git/info/exclude` |

No media file created. No existing file modified or deleted. Pre-existing dirty files from
sibling sessions (`CHANGELOG.md`, `index.html`, `main.go`, `debrand_copy_test.go`,
`verify/wz16-e2e-gate/`, `.hermes/`) were **not** staged and not touched.
