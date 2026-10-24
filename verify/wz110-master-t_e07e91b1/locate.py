#!/usr/bin/env python3
"""Locate where the longest visual freeze falls in the master."""
import os
import subprocess

from PIL import Image, ImageChops, ImageStat

MASTER = "/tmp/build/master_silent.mp4"
FPS = 2
W = 320
FROZEN = 0.010
SEGMENTS = [("00_intro", 0.0, 4.0), ("01_build", 4.0, 16.0), ("02_prov", 16.0, 30.0),
            ("03_cfg", 30.0, 46.0), ("04_work", 46.0, 54.0)]

d = "/tmp/build/d_loc"
os.makedirs(d, exist_ok=True)
for f in os.listdir(d):
    os.remove(os.path.join(d, f))
subprocess.run(["ffmpeg", "-y", "-v", "error", "-i", MASTER,
                "-vf", f"fps={FPS},scale={W}:-1,format=gray",
                os.path.join(d, "f_%05d.png")], check=True)
fl = sorted(os.path.join(d, f) for f in os.listdir(d))
diffs, prev = [], None
for p in fl:
    im = Image.open(p).convert("L")
    if prev is not None:
        diffs.append(ImageStat.Stat(ImageChops.difference(im, prev)).mean[0])
    prev = im

# find all frozen runs
runs = []
start = None
for i, x in enumerate(diffs):
    if x < FROZEN:
        if start is None:
            start = i
    else:
        if start is not None:
            runs.append((start, i, i - start))
            start = None
if start is not None:
    runs.append((start, len(diffs), len(diffs) - start))
runs.sort(key=lambda r: -r[2])
print("longest frozen runs (master time):")
for s, e, n in runs[:6]:
    st, en = s / FPS, e / FPS
    seg = next((nm for nm, a, b in SEGMENTS if a <= st < b), "?")
    print(f"  {st:6.2f}..{en:6.2f}s  len={n / FPS:4.1f}s  seg={seg}")
print()
print("frozen frames in the intro window (0-4s):")
for i, x in enumerate(diffs):
    t = i / FPS
    if t < 4.5:
        print(f"  t={t:4.2f}s diff={x:.4f} {'FROZEN' if x < FROZEN else ''}")
