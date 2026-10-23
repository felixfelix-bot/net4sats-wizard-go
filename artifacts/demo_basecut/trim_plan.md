# trim_plan.md — cut plan for the demo base-cut encode pass

**Task:** `t_ae815553` — "Assemble trim_plan.json + trim_plan.md from upstream findings"
**Board:** `net4sats-mvp-v2` · **Worker:** worker-heavy (fleet offload) · **Host:** c08r4d0r / CobradorWave
**Generated:** 2026-10-23T15:50:53Z (epoch 1792770653)
**Machine-readable companion:** `artifacts/demo_basecut/trim_plan.json`

---

## VERDICT: **BLOCKED — plan is intentionally non-executable**

All three declared upstream inputs are **absent**. There is no raw take for this chain on
this host, therefore there is no `t=0`, no EOF, no step boundaries, and no measured dead-air
regions. The task's own INPUTS clause governs:

> "If any is missing or contradictory, STOP and ask for clarification rather than filling
> gaps with invented timestamps."

So this plan **stops**. It contains zero timestamps for source, cuts, step map, and flagged
regions. Every such field is `null` or `[]` on purpose. A number in those fields would be a
fabrication, and a fabricated cut point is precisely the failure this artifact exists to
prevent — the encode pass would happily clip a step and nobody would catch it until QC.

**The encoder must not run against this plan.**

---

## 1. Why it is blocked — the three inputs, measured

Status read directly from `kanban.db` (`tasks`, `task_runs`, `task_comments`,
`task_attachments`), not inferred.

| # | Input | Card | Observed status | Runs | Completed | `result` | Comments | Attachments |
|---|---|---|---|---|---|---|---|---|
| 1 | source probe report | `t_15057569` | `triage` | **0** | 0 | NULL | 0 | 0 |
| 2 | frame-verified step map | `t_6d731c04` | `triage` | **0** | 0 | NULL | 0 | 0 |
| 3 | full-runtime scan report | `t_be3a9ae1` | `triage` | **0** | 0 | NULL | 0 | 0 |

Zero runs on all three means they never executed. There is nothing to read, from the board
or from disk.

### Root cause: the parent recording task never produced a take

| Card | Title | Status | Runs | Completed | Attachments |
|---|---|---|---|---|---|
| `t_882b3922` | Record full E2E de-branded wizard happy-path video | `triage` | 3 | **0** | **0** |
| `t_acee3752` | Prepare and validate 1080p capture toolchain | `triage` | **0** | 0 | 0 |

`t_882b3922`'s own run log records the blocker in its own words: the wizard happy path needs
a reachable OpenWrt router with SSH on 22, and none was reachable from the recording host. It
probed the LAN and the peer nodes (x240, dq05) for `192.168.1.1:22` and found nothing.

No recording run → no footage → no probe → no step map → no scan → no plan. The chain is not
missing one input; it is severed at its source.

### Independent corroboration, 12 days earlier

The downstream QC card `t_54b0d370` wrote this on 2026-10-11 and attached it to the board
(`attachments/t_54b0d370/trim_log.md`):

> "D3 | Raw footage does not exist, so no plan/cut is derivable | card `t_882b3922` —
> status `ready`, 2 runs, 0 completed, 0 attachments"

Same conclusion, same root cause, reached independently before this run. This is a persistent
chain-level stall, not a transient miss.

---

## 2. Source properties

**None.** `source.path` is `null`, `source.properties` is `null`.

No file was identified as this chain's raw take, so container, duration, resolution, fps,
pixel format, audio stream summary, and the CFR/VFR verdict are all **not measurable**. The
plan states this explicitly rather than borrowing plausible values from an unrelated file.

The encoder needs all of these before it can choose a scale/pad expression or an output fps:

`path, container, duration_s, width, height, r_frame_rate, avg_frame_rate, pix_fmt,
audio_streams, cfr_vfr_verdict, sha256` — carried in the JSON as
`source.required_before_encode_can_run`.

---

## 3. Cut points

| field | value |
|---|---|
| `head_trim_in_s` | `null` |
| `head_trim_out_s` | `null` |
| `tail_trim_in_s` | `null` |
| `tail_trim_out_s` | `null` |
| `interior_cuts_made` | `0` |
| `status` | `NOT_COMPUTABLE` |

Head trim is defined as *"from t=0 up to, but not including, the first frame where real action
begins"*. Tail trim is *"from the last meaningful frame of the final step to EOF"*. Both
definitions are relative to a known take. With no take there is no first action frame and no
final step, so neither edge can be placed — and the "~0.5s inside the dead zone" rule cannot
be applied to a dead zone that was never measured.

---

## 4. Step map

**Empty — 0 steps.** `step_map.ordered_steps = []`.

An empty list is the correct representation of "the step map was never produced." Inventing a
plausible 7-step wizard sequence from the card titles would look helpful and be worthless:
every timestamp in it would be fiction.

---

## 5. Flagged interior dead air — flag only, never remove

**0 regions inherited.** `flagged_interior_dead_air.regions = []`,
`removal_proposed = false`.

Because the scan card never ran, there are zero inherited regions to carry over verbatim.
**This must not be read as "the take has no interior dead air."** It means the take was never
scanned. Nothing was dropped and nothing was invented.

The structural guarantee holds by construction: this plan proposes **no** interior cut
anywhere, and `interior_cuts_made` is `0`.

---

## 6. Framing anomalies

**0 listed, `status = NOT_AVAILABLE`.** Not measured for this chain's take. Again, not a
clean bill of health — an absence of measurement.

---

## 7. Dead-air detection evidence — real commands, correctly scoped

The acceptance criteria require at least one dead-air detection result **and the command it
came from**. Two real detectors were executed with exit code 0. Their raw output is embedded
in `trim_plan.json` → `dead_air_detection_evidence`.

**Critical scope limit, stated up front:** these ran against **CANDIDATE-A**, which was
*rejected* as this chain's take (see §8). They prove the detectors execute and emit real
output on this host. They are **not** the plan's dead-air regions, are marked
`POSITIVE_CONTROL_ONLY`, and **must not be applied to any encode**.

Tooling: `ffmpeg 6.1.1-3ubuntu5`, `ffprobe 6.1.1-3ubuntu5` — installed and working.

### 7a. Freeze detection — real output

```
$ ffmpeg -hide_banner -nostdin -i take-attempt-blocked.mp4 \
        -vf freezedetect=n=-60dB:d=2 -an -f null -
[freezedetect] lavfi.freezedetect.freeze_start: 0.76
[freezedetect] lavfi.freezedetect.freeze_duration: 4.36
[freezedetect] lavfi.freezedetect.freeze_end: 5.12
[freezedetect] lavfi.freezedetect.freeze_start: 5.16
[freezedetect] lavfi.freezedetect.freeze_duration: 5.08
[freezedetect] lavfi.freezedetect.freeze_end: 10.24
[freezedetect] lavfi.freezedetect.freeze_start: 10.28
[freezedetect] lavfi.freezedetect.freeze_duration: 5.08
[freezedetect] lavfi.freezedetect.freeze_end: 15.36
[freezedetect] lavfi.freezedetect.freeze_start: 15.4
[freezedetect] lavfi.freezedetect.freeze_duration: 5.08
[freezedetect] lavfi.freezedetect.freeze_end: 20.48
[freezedetect] lavfi.freezedetect.freeze_start: 20.52
[freezedetect] lavfi.freezedetect.freeze_duration: 5.08
[freezedetect] lavfi.freezedetect.freeze_end: 25.6
[freezedetect] lavfi.freezedetect.freeze_start: 25.64
[freezedetect] lavfi.freezedetect.freeze_duration: 4.24
[freezedetect] lavfi.freezedetect.freeze_end: 29.88
[freezedetect] lavfi.freezedetect.freeze_start: 36.16
[freezedetect] lavfi.freezedetect.freeze_duration: 4.8
[freezedetect] lavfi.freezedetect.freeze_end: 40.96
```

7 regions detected, each ≥ 4.24s. Exit code 0.

### 7b. Silence detection — real output, and why it is empty

```
$ ffmpeg -hide_banner -nostdin -i take-attempt-blocked.mp4 \
        -af silencedetect=noise=-35dB:d=2 -f null -
(tool runs cleanly; no silencedetect events emitted)
```

No regions — **and the reason is recorded, not hidden**: this file has **zero audio streams**
(`ffprobe -select_streams a` → count `0`). There is no audio to analyse, so `silencedetect`
correctly reports nothing. The scan card explicitly asked for this reason to be recorded when
a detector produces no output.

---

## 8. Candidate footage disposition — every plausible file probed and rejected

This is what makes the BLOCKED verdict a measured finding rather than a shrug. Every video on
this host that could conceivably be mistaken for the take was probed with real `ffprobe` and
explicitly ruled out.

| Candidate | Probe | Verdict | Why rejected |
|---|---|---|---|
| `.../wz19-video-capture/artifacts/take-attempt-blocked.mp4` (CANDIDATE-A) | h264, 1920x1080, yuv420p, 42.88s, 25/1 CFR, 1072 frames, **0 audio**, sha256 `628d7ad3…` | **REJECT** | Captures a **failed** deploy. Its companion `step-log.txt` ends `STEP 7 deploy-result — ERROR VIEW SHOWN`, `sawSuccess=false sawError=true`. Chain needs the happy path. Parent card has 0 completed runs, so it cannot be its output. Filename declares a blocked attempt. |
| `wz19-atomic/repo/docs/demo/demo-trimmed.mp4` | h264, 1920x1080, 25/1, 21.84s, sha256 `83bdc911…` | **REJECT** | **A derivative, not raw footage.** Its sha256 is byte-identical to the `output.sha256` recorded by an older trim plan (`t_3e3082f1`, FD9 pass 12). Trimming an already-trimmed master would double-cut. Also a different demo chain (cashu/SMS wallet UI). |
| `attachments/t_7015c244/intro_card.mp4` | 1920x1080, 4.00s, sha256 `782766e0…` | **REJECT** | 4s synthetic intro card. No wizard steps → no head action, no final step. Bookend asset, not a take. |
| `reports/…/film5/net4sats-installer-tour.mp4` | 1280x800, 284.24s | **REJECT** | Below the ≥1080p bar; different aspect; no step log maps to it. |
| `reports/…/film/e2e-net4sats.mp4` | 1280x800, 40.88s | **REJECT** | Below ≥1080p bar; no matching step log. |
| `reports/…/film3/net4sats-branded-e2e.mp4` | 1280x800, 85.00s | **REJECT** | Below ≥1080p bar; **"branded"** — opposite of the de-branded requirement. |
| `/tmp/demo-video.webm` (FD9-era source) | — | **REJECT** | Does not exist on this host (ENOENT). |

**Accepted candidate: none.**

Two findings worth carrying forward, because both would have caused silent damage if a
greedy worker had "found some footage":

1. `demo-trimmed.mp4` looks like a perfect 1080p 25fps source. It is actually the **output of
   a previous trim pass**, byte-verified by hash. Using it would compound two cuts.
2. CANDIDATE-A is a genuinely valid 1920x1080 CFR capture with a real step log — but the log
   documents a **failure**, and the file is **branded**. It is the single most tempting false
   positive on this host.

### Candidate-A measured properties (for the record only — not this chain's source)

```
format_name = mov,mp4,m4a,3gp,3g2,mj2     duration = 42.880000   size = 647139
video: h264, 1920x1080, yuv420p, r_frame_rate 25/1, avg_frame_rate 25/1, nb_frames 1072
audio: none (0 streams)
CFR/VFR  = CFR   (r_frame_rate == avg_frame_rate == 25/1;
                  every sampled packet duration_time == 0.040000)
```

Framing control on the same file: `ffprobe … -show_entries frame=width,height | sort -u` →
`1920,1080` only, i.e. no mid-take resolution change *in that file*. Method validated; the
verdict does not transfer to the real take because the real take does not exist.

---

## 9. Validation results

| ID | Requirement | Result | Note |
|---|---|---|---|
| V1 | Both plan files exist at stated paths | **PASS** | `artifacts/demo_basecut/trim_plan.json`, `artifacts/demo_basecut/trim_plan.md` |
| V2 | Every step has start & end timestamp | **VACUOUS** | 0 steps exist. True only degenerately — not evidence a step map exists. |
| V3 | No trim point inside any step interval | **VACUOUS** | All 4 trim points are `null`, so none can fall in an interval. Deliberate. |
| V4 | Every flagged region carried verbatim | **VACUOUS** | 0 inherited regions to carry. Nothing dropped, nothing invented. |
| V5 | No interior content proposed for removal | **PASS** | Regions empty; `removal_proposed=false`; `interior_cuts_made=0`. |
| V6 | ≥1 dead-air detection with its command | **PASS (scoped)** | 2 commands, exit 0, raw output embedded — but marked `POSITIVE_CONTROL_ONLY` against a rejected candidate. |
| V7 | Inherited uncertainty carried into this .md | **PASS** | See §10. |
| V8 | No invented timestamps | **PASS** | Zero timestamps in source/cuts/step_map/regions. |

**Overall: BLOCKED.** V2, V3 and V4 pass only vacuously. Read them as *"no violation exists
because no subject exists"*, never as *"the cut plan checks out."*

---

## 10. Inherited uncertainty and unresolved questions

Carried forward verbatim rather than resolved by guessing:

1. **No take exists.** Unresolved. `t_882b3922` never completed; its 3 runs all ended
   non-productively. Which router, which host, and which capture toolchain will actually
   produce the take is unknown.
2. **Toolchain unvalidated.** `t_acee3752` ("Prepare and validate 1080p capture toolchain")
   has **0 runs**. The ≥1080p/CFR/stable-framing assumptions for the eventual take are
   therefore unproven, not merely unmeasured.
3. **Step-log path unknown.** No durable step-log path exists for this chain. The
   `t_54b0d370` QC log flagged this specifically, asking for a `STEP <n> <label> <epoch_ms>`
   log written to a **durable** path. Not done.
4. **Branding state of the intended take is unverified.** This chain's whole point is a
   **de-branded** flow. CANDIDATE-A is branded per `BRANDING-AUDIT.md`. Whether a future take
   would be de-branded is untested; a branding audit is needed before the plan is re-run.
5. **Step-count and step-name list unknown.** Card titles imply build → provisioning →
   configuration → resulting working state, but no step log exists to confirm the actual
   sequence, count, or naming. The JSON therefore contains **no** step entries.
6. **Audio presence unknown.** If the eventual take has no audio track (CANDIDATE-A has none),
   `silencedetect` will again be inapplicable and the plan should record that reason rather
   than treating silence as absent.
7. **Framing stability unknown.** Whether the eventual take has window moves, relayouts, or
   resolution changes is unmeasured. The encoder's need for one constant output size is
   therefore unaddressed.
8. **Dead-air threshold behaviour on a real take unknown.** `freezedetect` at `-60dB:d=2`
   produced dense, overlapping-looking regions on a mostly-static UI capture. On the real take
   this may flag long stretches of legitimately static UI; the threshold may need revisiting
   at re-run time. Noting it now as a known tuning risk.

---

## 11. Exact commands used in this run

```bash
# --- Board state (read-only; sqlite3 CLI is version-mismatched here, used Python sqlite3) ---
# db: /home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db
#     tables: tasks, task_runs, task_comments, task_attachments, task_events, task_links

# --- ffprobe every candidate ---
ffprobe -v error -show_streams -show_format <file>
ffprobe -v error -show_entries format=format_name,duration,size \
        -show_entries stream=index,codec_type,codec_name,width,height,\
r_frame_rate,avg_frame_rate,pix_fmt,nb_frames,sample_rate,channels -of json <file>
ffprobe -v error -select_streams a -show_entries stream=index -of csv=p=0 <file> | wc -l
ffprobe -v error -select_streams v:0 -show_entries frame=width,height -of csv=p=0 <file> | sort -u

# --- CFR/VFR verdict ---
ffprobe -v error -select_streams v:0 -show_entries packet=pts_time,duration_time \
        -of csv=p=0 <file> | head -12        # all duration_time == 0.040000 => CFR

# --- Dead-air detectors (verbatim from the scan card) ---
ffmpeg -hide_banner -nostdin -i <file> -vf freezedetect=n=-60dB:d=2 -an -f null -
ffmpeg -hide_banner -nostdin -i <file> -af silencedetect=noise=-35dB:d=2 -f null -

# --- Candidate provenance / double-cut guard ---
sha256sum <file>     # compared against the older plan's recorded output.sha256

# --- Confirm the chain's inputs are truly absent ---
ls -la /home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/attachments/{t_15057569,t_6d731c04,t_be3a9ae1}
ls -la /tmp/demo_meta /tmp/demo-video.webm          # -> NO_DIR / MISSING
find .../workspaces .../attachments -type f \( -iname '*.mp4' -o -iname '*.webm' \
     -o -iname '*.mkv' -o -iname '*.mov' \)          # -> only intro_card.mp4
```

---

## 12. What unblocks this card

In order — each is a hard prerequisite, not a suggestion:

1. **`t_acee3752`** completes: validated 1080p capture toolchain (frame size, CFR, stable framing).
2. **`t_882b3922`** completes: one continuous take of the **de-branded** happy path, plus a
   **durable** step log at an absolute path.
3. **`t_15057569`** writes the probe report with every field in §2.
4. **`t_6d731c04`** writes the step map: `{step, start, end, confidence, evidence}` per step,
   single consistent timestamp format, boundaries frame-verified, uncertain boundaries flagged.
5. **`t_be3a9ae1`** writes interior dead-air regions `{start, end, duration, reason}` plus
   framing anomalies with timestamps.
6. **Re-dispatch `t_ae815553`** against those artifacts. The plan then becomes executable and
   the encode card `t_ba89d8b4` may run.

Until step 5 lands, `t_ba89d8b4` **must not run**. Its own card says so: *"If the plan is
missing, incomplete, or its source paths do not resolve, STOP and ask rather than improvising
cut points."*

---

## 13. Scope and hygiene

- **No media was created, modified, trimmed, or re-encoded.** Read-only analysis throughout.
- **No timestamp was invented.** Every null/empty field is a deliberate representation of
  "not measured."
- **No secrets read, printed, or written.** Only video metadata, board state, and hashes.
- The rejected candidates were only *probed*, never written to.
- One genuine risk was detected and averted: `demo-trimmed.mp4` (a prior trim output) was
  identified by hash and excluded, preventing a double-cut.
