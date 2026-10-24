#!/usr/bin/env python3
"""Render the caption strip PNGs: one title strip + four phase strips.

Layout contract (chosen so captions can NEVER occlude the UI):
  canvas 1920x1192 = video area 1920x1080 (top) + caption strip 1920x112 (bottom)
The strip is a separate raster row-belt stacked under the video, so caption
pixels can never overwrite demonstrated UI pixels.

All strips use ONE uniform main font size = the largest size that fits every
caption, so type size never jumps between phases.

The PR number in the title strip is read from the intro-card metadata JSON at
render time -- never retyped by hand.
"""
import json
import os

from PIL import Image, ImageDraw, ImageFont

FONT_BOLD = "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"

META = ("/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/"
        "attachments/t_7015c244/intro_card.mp4.meta.json")
META_FALLBACK = ("/home/c08r4d0r/.hermes/kanban/boards/net4sats-mvp-v2/"
                 "attachments/t_0143bb0e/intro_card_meta.json")

W = 1920
STRIP_H = 112
BG = (16, 16, 20)
FG = (245, 245, 248)
ACCENT = (247, 147, 26)   # bitcoin orange
MUTED = (150, 155, 165)
PAD_X = 40
TAG_GAP = 28

SIZE_TAG = 30
SIZE_MAIN_START = 52
SIZE_MAIN_FLOOR = 26


def load_pr_number():
    """Read the PR number from the frozen metadata; never type it."""
    for p in (META, META_FALLBACK):
        if os.path.exists(p):
            with open(p) as fh:
                d = json.load(fh)
            if "pr_number" in d:
                return int(d["pr_number"]), p
    raise RuntimeError("no metadata JSON with pr_number found")


PR_NUMBER, PR_SOURCE = load_pr_number()

TITLE_TAG = "INTRO"
TITLE_TEXT = f"net4sats wizard -- de-branded E2E walkthrough (PR #{PR_NUMBER})"

PHASES = [
    (1, "BUILD", "Finds the router on the LAN, takes the root password, then Deploy"),
    (2, "PROVISIONING", "Verifies SSH, flashes OpenWrt, installs the tollgate-wrt package"),
    (3, "CONFIGURATION", "Router admin UI: WiFi radio, connected devices, wallet balance"),
    (4, "WORKING STATE", "Captive portal sells access -- pay with Cashu or Lightning"),
]

# strips 0..4 -> (tag, text); strip 0 is the intro title strip
STRIPS = [(TITLE_TAG, TITLE_TEXT)] + [(f"{i}/4  {t}", x) for i, t, x in PHASES]


def uniform_main_size():
    """Largest main size that fits every strip's text in its own available width."""
    f_tag = ImageFont.truetype(FONT_BOLD, SIZE_TAG)
    widest_tag = max(f_tag.getlength(tag) for tag, _ in STRIPS)
    for size in range(SIZE_MAIN_START, SIZE_MAIN_FLOOR - 1, -1):
        f = ImageFont.truetype(FONT_BOLD, size)
        if all(f.getlength(text) <= (W - (PAD_X + widest_tag + TAG_GAP) - PAD_X)
               for _, text in STRIPS):
            return size
    raise RuntimeError("no uniform size fits")


def render(idx, tag, text, size):
    im = Image.new("RGB", (W, STRIP_H), BG)
    d = ImageDraw.Draw(im)
    d.rectangle([0, 0, W, 4], fill=ACCENT)

    f_tag = ImageFont.truetype(FONT_BOLD, SIZE_TAG)
    f_main = ImageFont.truetype(FONT_BOLD, size)

    x = PAD_X
    asc, desc = f_tag.getmetrics()
    d.text((x, (STRIP_H - (asc + desc)) // 2 - 3), tag, font=f_tag, fill=ACCENT)
    x += f_tag.getlength(tag) + TAG_GAP

    d.rectangle([x - 14, 26, x - 13, STRIP_H - 26], fill=MUTED)

    asc, desc = f_main.getmetrics()
    ink = f_main.getlength(text)
    d.text((x, (STRIP_H - (asc + desc)) // 2 - 3), text, font=f_main, fill=FG)

    out = f"/tmp/build/strip_{idx}.png"
    im.save(out)
    print(f"strip_{idx}.png {W}x{STRIP_H} tag={tag!r} main={size}px "
          f"ink_w={ink:.0f} avail={W - x - PAD_X:.0f} "
          f"left={x:.0f} right_margin={W - x - ink:.0f}")
    return out


if __name__ == "__main__":
    print(f"pr_number={PR_NUMBER} read from {PR_SOURCE}")
    size = uniform_main_size()
    print(f"uniform main size = {size}px (floor {SIZE_MAIN_FLOOR}px; requirement >=24px)")
    for idx, (tag, text) in enumerate(STRIPS):
        render(idx, tag, text, size)
