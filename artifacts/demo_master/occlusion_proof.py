#!/usr/bin/env python3
"""Independent, pixel-level proof that captions do not occlude the demo UI.

Method: build the SAME composite frame twice -- once with the real caption
strip, once with a blank strip of identical height -- and write the frames
straight out as PNG (NO video codec involved, so encoder rate-control noise
cannot contaminate the comparison). Then diff ONLY the video belt
(rows 0..1079) of the two renders.

If the captions were overlaying the UI, the video belt would differ.
"""
import json
import os
import subprocess

from PIL import Image, ImageChops, ImageStat

BUILD = "/tmp/build"
VW, VH, STRIP_H = 1920, 1080, 112
OUT_H = VH + STRIP_H
FPS = 25

SRC = "/home/c08r4d0r/reports/net4sats-mvp/film5/01-installer-deploy.webm"
A, B = 8.0, 20.0
STRIP = f"{BUILD}/strip_1.png"
BLANK = f"{BUILD}/strip_blank.png"
WINDOW = 4.0
PICK = 5.0   # fps for the sampled frames


def run(cmd):
    p = subprocess.run(cmd, capture_output=True, text=True)
    if p.returncode != 0:
        raise RuntimeError(" ".join(cmd) + "\n" + p.stderr[-2000:])


def render_frames(strip_png, outdir):
    os.makedirs(outdir, exist_ok=True)
    for f in os.listdir(outdir):
        os.remove(os.path.join(outdir, f))
    vf = (f"[0:v]scale=1728:{VH},pad={VW}:{VH}:96:0:color=0x101014,setsar=1,"
          f"fps={FPS},trim=duration={WINDOW},setpts=PTS-STARTPTS[bg];"
          f"[1:v]format=rgb24[strip];[bg][strip]vstack=inputs=2,"
          f"fps={PICK}[v]")
    run(["ffmpeg", "-y", "-v", "error", "-ss", str(A), "-to", str(B), "-i", SRC,
         "-loop", "1", "-i", strip_png, "-filter_complex", vf, "-map", "[v]",
         "-frames:v", str(int(PICK * WINDOW)), "-pix_fmt", "rgb24",
         os.path.join(outdir, "f_%03d.png")])
    return sorted(os.path.join(outdir, f) for f in os.listdir(outdir))


def main():
    Image.new("RGB", (VW, STRIP_H), (16, 16, 20)).save(BLANK)

    fa = render_frames(STRIP, f"{BUILD}/occl_a")
    fb = render_frames(BLANK, f"{BUILD}/occl_b")
    n = min(len(fa), len(fb))

    worst_video = 0.0
    worst_strip = 0.0
    vpx = 0
    for pa, pb in zip(fa[:n], fb[:n]):
        ia = Image.open(pa).convert("RGB")
        ib = Image.open(pb).convert("RGB")
        assert ia.size == (VW, OUT_H), ia.size

        va, vb = ia.crop((0, 0, VW, VH)), ib.crop((0, 0, VW, VH))
        d = ImageChops.difference(va, vb)
        m = ImageStat.Stat(d).mean[0]
        worst_video = max(worst_video, m)
        if d.getbbox() is not None:
            vpx += sum(1 for p in d.convert("L").getdata() if p != 0)

        sa, sb = ia.crop((0, VH, VW, OUT_H)), ib.crop((0, VH, VW, OUT_H))
        worst_strip = max(worst_strip,
                          ImageStat.Stat(ImageChops.difference(sa, sb)).mean[0])

    print(f"frames compared : {n} ({PICK:.0f} fps over a {WINDOW}s window of the BUILD segment)")
    print(f"VIDEO BELT      rows 0..{VH - 1}       worst mean |diff| = {worst_video:.6f} "
          f"(differing pixels total = {vpx})")
    print(f"CAPTION BELT    rows {VH}..{OUT_H - 1}     worst mean |diff| = {worst_strip:.6f}")
    verdict = worst_video == 0.0 and worst_strip > 1.0
    print()
    print("VERDICT:", "PASS - the video belt is BYTE-IDENTICAL with and without the caption "
                     "strip, while the caption belt changes; caption ink is provably confined "
                     "to rows >= 1080" if verdict else "FAIL")
    json.dump({"frames": n, "video_belt_worst_diff": worst_video,
               "video_belt_differing_pixels": vpx,
               "caption_belt_worst_diff": worst_strip,
               "video_belt_rows": [0, VH - 1], "caption_belt_rows": [VH, OUT_H - 1],
               "method": "direct rgb24 PNG extraction, no video codec",
               "verdict": "PASS" if verdict else "FAIL"},
              open(f"{BUILD}/occlusion.json", "w"), indent=1)


if __name__ == "__main__":
    main()
