#!/usr/bin/env python3
"""Assemble the draft master: intro card + 4 phase segments + caption strips.

Canvas contract:
    video area  1920x1080  (a true 1080p frame)
    caption strip 1920x112 appended BELOW it
    output      1920x1192  (>= 1080p)

The strip is a separate raster row-belt stacked under the video, so caption
pixels can never overwrite demonstrated UI pixels.

The intro card is scaled 1920x1080 -> 1920x1080 (identity), so its SHA-bearing
pixels are preserved 1:1.
"""
import json
import os
import subprocess

BUILD = "/tmp/build"
FPS = 25
CRF = 18
VW, VH = 1920, 1080
STRIP_H = 112
OUT_H = VH + STRIP_H

FILMS = {
    "inst":   "/home/c08r4d0r/reports/net4sats-mvp/film5/01-installer-deploy.webm",
    "portal": "/home/c08r4d0r/reports/net4sats-mvp/film5/02-portal.webm",
    "admin":  "/home/c08r4d0r/reports/net4sats-mvp/film5/03-admin-tour.webm",
}

INTRO = f"{BUILD}/intro_card.mp4"

# (segment_id, kind, source_key, start, end, strip_png)
SEGMENTS = [
    ("00_intro", "card",  None,     None,  None, f"{BUILD}/strip_0.png"),
    ("01_build", "film", "inst",     8.0,  20.0, f"{BUILD}/strip_1.png"),
    ("02_prov",  "film", "inst",    22.0,  36.0, f"{BUILD}/strip_2.png"),
    ("03_cfg",   "film", "admin",   36.0,  52.0, f"{BUILD}/strip_3.png"),
    ("04_work",  "film", "portal",  20.0,  28.0, f"{BUILD}/strip_4.png"),
]


def run(cmd):
    p = subprocess.run(cmd, capture_output=True, text=True)
    if p.returncode != 0:
        raise RuntimeError(f"cmd failed: {' '.join(cmd)}\n{p.stderr[-3000:]}")
    return p


def build_segment(seg):
    sid, kind, key, a, b, strip = seg
    out = f"{BUILD}/seg_{sid}.mp4"
    if kind == "card":
        # identity scale: card is already 1920x1080
        vf = (f"[0:v]scale={VW}:{VH},setsar=1,fps={FPS}[bg];"
              f"[bg][1:v]vstack=inputs=2[v]")
        cmd = ["ffmpeg", "-y", "-v", "error", "-i", INTRO, "-i", strip,
               "-filter_complex", vf, "-map", "[v]",
               "-c:v", "libx264", "-preset", "slow", "-crf", str(CRF),
               "-pix_fmt", "yuv420p", "-movflags", "+faststart", out]
    else:
        src = FILMS[key]
        # 1280x800 (1.6) -> 1728x1080 exactly, pillarboxed 96px each side
        vf = (f"[0:v]scale=1728:{VH},pad={VW}:{VH}:96:0:color=0x101014,setsar=1,"
              f"fps={FPS}[bg];[bg][1:v]vstack=inputs=2[v]")
        cmd = ["ffmpeg", "-y", "-v", "error", "-ss", str(a), "-to", str(b),
               "-i", src, "-i", strip, "-filter_complex", vf, "-map", "[v]",
               "-an", "-c:v", "libx264", "-preset", "slow", "-crf", str(CRF),
               "-pix_fmt", "yuv420p", "-movflags", "+faststart", out]
    run(cmd)
    return out


def probe(path):
    p = run(["ffprobe", "-v", "error", "-select_streams", "v:0",
             "-show_entries", "stream=width,height,r_frame_rate,nb_frames,codec_name",
             "-show_entries", "format=duration,size", "-of", "json", path])
    return json.loads(p.stdout)


def main():
    parts = []
    for seg in SEGMENTS:
        p = build_segment(seg)
        info = probe(p)
        st = info["streams"][0]
        parts.append(p)
        print(f"{os.path.basename(p):22s} {st['codec_name']} {st['width']}x{st['height']} "
              f"{st['r_frame_rate']} frames={st['nb_frames']} "
              f"dur={float(info['format']['duration']):.3f}s")

    lst = f"{BUILD}/concat.txt"
    with open(lst, "w") as fh:
        for p in parts:
            fh.write(f"file '{p}'\n")

    master = f"{BUILD}/master_silent.mp4"
    run(["ffmpeg", "-y", "-v", "error", "-f", "concat", "-safe", "0", "-i", lst,
         "-c", "copy", master])
    print("master:", master, probe(master)["format"]["duration"])


if __name__ == "__main__":
    main()
