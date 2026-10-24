#!/usr/bin/env python3
"""Export the draft master at FINAL quality target settings.

Name scheme: demo-<pr>-<shortsha>.mp4, both fields read from the frozen
metadata JSON -- never typed.
"""
import json
import os
import shutil
import subprocess

BUILD = "/tmp/build"
REPO = "/home/c08r4d0r/repos/net4sats-mvp-v2"
OUTDIR = f"{REPO}/artifacts/demo_master"
META = ("/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/"
        "attachments/t_0143bb0e/intro_card_meta.json")

CRF = 16          # FINAL-quality target (draft pass used 18)
PRESET = "slower"


def run(cmd):
    p = subprocess.run(cmd, capture_output=True, text=True)
    if p.returncode != 0:
        raise RuntimeError(" ".join(cmd) + "\n" + p.stderr[-2000:])
    return p


def main():
    meta = json.load(open(META))
    pr = int(meta["pr_number"])
    sha = meta["head_sha"]
    short = sha[:8]
    name = f"demo-pr{pr}-{short}.mp4"
    os.makedirs(OUTDIR, exist_ok=True)
    out = os.path.join(OUTDIR, name)

    src = f"{BUILD}/master_silent.mp4"
    # FINAL export: re-encode the assembled master at the target settings
    run(["ffmpeg", "-y", "-v", "error", "-i", src,
         "-c:v", "libx264", "-preset", PRESET, "-crf", str(CRF),
         "-pix_fmt", "yuv420p", "-profile:v", "high", "-level", "4.1",
         "-movflags", "+faststart", "-an", out])

    info = json.loads(run(["ffprobe", "-v", "error", "-show_streams", "-show_format",
                           "-of", "json", out]).stdout)
    v = [s for s in info["streams"] if s["codec_type"] == "video"][0]
    print(f"name       : {name}   (pr={pr} sha={sha} short={short})")
    print(f"path       : {out}")
    print(f"size       : {os.path.getsize(out):,} bytes")
    print(f"codec      : {v['codec_name']} profile={v.get('profile')} "
          f"level={v.get('level')} pix_fmt={v.get('pix_fmt')}")
    print(f"resolution : {v['width']}x{v['height']}")
    print(f"fps        : {v['r_frame_rate']} (avg {v['avg_frame_rate']})")
    print(f"duration   : {float(info['format']['duration']):.3f}s")
    print(f"bitrate    : {int(info['format']['bit_rate']) // 1000} kbps")

    sha_out = run(["sha256sum", out]).stdout.split()[0]
    print(f"sha256     : {sha_out}")

    # keep the builder inputs with the artifact so the build is reproducible
    for f in ("build_master.py", "make_strips.py", "verify.py", "occlusion_proof.py"):
        shutil.copy(os.path.join(BUILD, f), os.path.join(OUTDIR, f))

    json.dump({"name": name, "path": out, "size": os.path.getsize(out),
               "sha256": sha_out, "pr_number": pr, "head_sha": sha, "short_sha": short,
               "codec": v["codec_name"], "profile": v.get("profile"),
               "resolution": f"{v['width']}x{v['height']}",
               "fps": v["r_frame_rate"], "duration": float(info["format"]["duration"]),
               "audio_streams": len([s for s in info["streams"] if s["codec_type"] == "audio"]),
               "final_crf": CRF, "final_preset": PRESET},
              open(os.path.join(OUTDIR, "master.json"), "w"), indent=1)
    print("wrote", os.path.join(OUTDIR, "master.json"))


if __name__ == "__main__":
    main()
