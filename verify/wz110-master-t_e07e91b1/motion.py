#!/usr/bin/env python3
"""Measure per-second visual motion (dead-air detector) by PIL frame differencing.

No ffmpeg log parsing: extract frames at a fixed fps, scale to a small grey
canvas, then mean absolute difference between consecutive frames.
"""
import os, subprocess, sys, json
from PIL import Image, ImageChops, ImageStat

FILMS = {
    "inst":   "/home/c08r4d0r/reports/net4sats-mvp/film5/01-installer-deploy.webm",
    "portal": "/home/c08r4d0r/reports/net4sats-mvp/film5/02-portal.webm",
    "admin":  "/home/c08r4d0r/reports/net4sats-mvp/film5/03-admin-tour.webm",
}
FPS = 2
W = 320

def frames(src, outdir):
    os.makedirs(outdir, exist_ok=True)
    for f in os.listdir(outdir):
        os.remove(os.path.join(outdir, f))
    subprocess.run(
        ["ffmpeg", "-y", "-v", "error", "-i", src,
         "-vf", f"fps={FPS},scale={W}:-1,format=gray",
         os.path.join(outdir, "f_%05d.png")],
        check=True)
    return sorted(os.path.join(outdir, f) for f in os.listdir(outdir))

def main():
    result = {}
    for tag, src in FILMS.items():
        fl = frames(src, f"/tmp/build/mot_{tag}")
        vals = []
        prev = None
        for p in fl:
            im = Image.open(p).convert("L")
            if prev is not None:
                d = ImageChops.difference(im, prev)
                vals.append(ImageStat.Stat(d).mean[0])
            prev = im
        result[tag] = vals
        print(f"### {tag}  frames={len(fl)}  step={1.0/FPS:.2f}s  dur={len(fl)/FPS:.2f}s")
        # 2-second windows
        win = FPS * 2
        for i in range(0, len(vals), win):
            c = vals[i:i+win]
            if not c:
                continue
            print(f"  t={i/FPS:6.1f}..{(i+len(c))/FPS:6.1f}s  max={max(c):6.3f} mean={sum(c)/len(c):6.3f} min={min(c):6.3f}")
    json.dump({k: v for k, v in result.items()}, open("/tmp/build/motion.json", "w"))

main()
