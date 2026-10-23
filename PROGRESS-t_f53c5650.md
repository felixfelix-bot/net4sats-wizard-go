# PROGRESS — t_f53c5650 (Write /tmp/demo-meta.json from resolved metadata + probed footage)

## Finding

**BLOCKED — /tmp/demo-meta.json deliberately NOT written.** Parent task B (`t_dd851d4e`,
"Locate recording footage, ffprobe each file, record step-log path") has **0 runs, 0 attachments,
result NULL** — it has never executed on this host. Parent A (`t_85fda20d`) is likewise 0 runs /
result NULL.

Recovered A's four fields from a genuine chain artifact attached to this very card
(`attachments/t_f53c5650/resolved_pr_meta_partial-t_f53c5650.json`, 1251 B), which self-declares
PARTIAL and carries `footage: []` / `step_log: null`. head_sha `2d770966b163cd2deafb1697bfd892e0a7a8d9e2`
is corroborated by two further independent board records (t_0143bb0e run 1409 live `gh pr view 32`
headRefOid; t_069b6842 frame-OCR cross-check) — 3/3 agree, len 40, all hex, short `2d770966`.

Root cause upstream: recorder `t_882b3922` has 3 runs / 0 completed / 0 attachments;
`/tmp/demo-video.webm` ENOENT; board holds no take attachment. Two decoy MP4s that *do* ffprobe
(take-attempt-blocked.mp4 647139 B, sibling demo-trimmed.mp4 1205700 B) were examined and
rejected — neither is B's output, and adopting one would put a false path in a machine-readable
artifact. Prior QA on this chain (t_8bb677fc) reached the same BLOCKED verdict.

## Status

BLOCKED — neither blockable criterion satisfiable without fabrication. Read-only honoured: no
repo source edited, staged or committed to any shared branch.

## Files touched

- `REPORT-t_f53c5650.md` — full report (new)
- `PROGRESS-t_f53c5650.md` — this file (new)

## REMAINING steps (unblock order)

1. Operator-gated hardware available (GL-MT6000 + real SSH) → recorder produces a real take on a
   durable path (not /tmp).
2. Parent B `t_dd851d4e` runs: ffprobe each footage file, emit absolute footage paths + step_log.
3. Parent A `t_85fda20d` runs: emit verbatim pr_number/pr_title/branch/40-char head_sha.
4. Re-run this card: write /tmp/demo-meta.json with exactly the 7 keys, re-read + validate, print.
