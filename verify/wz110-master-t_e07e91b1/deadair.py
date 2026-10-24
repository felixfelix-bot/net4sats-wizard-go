#!/usr/bin/env python3
"""Dead-air analysis: for each candidate segment, report the longest run of
consecutive near-identical frames (max consecutive visual freeze in seconds).

Threshold: mean abs diff on a 320px grey downscale < 0.010 = frozen frame.
"""
import os, subprocess, json, sys
from PIL import Image, ImageChops, ImageStat

FPS = 2
W = 320
FILMS = {
    "inst":   "/home/c08r4d0r/reports/net4sats-mvp/film5/01-installer-deploy.webm",
    "portal": "/home/c08r4d0r/reports/net4sats-mvp/film5/02-portal.webm",
    "admin":  "/home/c08r4d0r/reports/net4sats-mvp/film5/03-admin-tour.webm",
}
FROZEN = 0.010

# candidate windows: (src, start, end, phase)
CANDIDATES = [
    ("inst",   0.0,  12.0),
    ("inst",   8.0,  14.0),
    ("inst",  18.0,  36.0),
    ("inst",  24.0,  36.0),
    ("inst", 126.0, 134.0),
    ("admin", 16.0,  22.0),
    ("admin", 36.0,  40.0),
    ("admin", 44.0,  50.0),
    ("admin", 52.0,  58.0),
    ("portal",20.0,  28.0),
]

def frames(src, start, end, outdir):
    os.makedirs(outdir, exist_ok=True)
    for f in os.listdir(outdir):
        os.remove(os.path.join(outdir, f))
    subprocess.run(["ffmpeg", "-y", "-v", "error", "-ss", str(start), "-to", str(end),
                    "-i", src, "-vf", f"fps={FPS},scale={W}:-1,format=gray",
                    os.path.join(outdir, "f_%05d.png")], check=True)
    return sorted(os.path.join(outdir, f) for f in os.listdir(outdir))

def longest_frozen_run(diffs):
    best = cur = 0
    for d in diffs:
        if d < FROZEN:
            cur += 1
            best = max(best, cur)
        else:
            cur = 0
    return best

rows = []
for src_tag, a, b, in [(c[0], c[1], c[2]) for c in CANDIDATES]:
    src = FILMS[src_tag]
    fl = frames(src, a, b, "/tmp/build/da")
    diffs, prev = [], None
    for p in fl:
        im = Image.open(p).convert("L")
        if prev is not None:
            diffs.append(ImageStat.Stat(ImageChops.difference(im, prev)).mean[0])
        prev = im
    run = longest_frozen_run(diffs)
    rows.append({
        "src": src_tag, "start": a, "end": b, "dur": round(b - a, 2),
        "frames": len(fl),
        "max_diff": round(max(diffs), 3) if diffs else 0.0,
        "mean_diff": round(sum(diffs)/len(diffs), 4) if diffs else 0.0,
        "longest_frozen_run_s": round(run / FPS, 2),
        "verdict": "PASS" if run / FPS <= 2.0 else "FAIL-deadair",
    })

for r in rows:
    print(f"{r['src']:7s} {r['start']:6.1f}-{r['end']:6.1f}s dur={r['dur']:5.1f} "
          f"max={r['max_diff']:6.3f} mean={r['mean_diff']:6.4f} "
          f"longest_freeze={r['longest_frozen_run_s']:4.1f}s  {r['verdict']}")
json.dump(rows, open("/tmp/build/deadair.json", "w"), indent=1)
