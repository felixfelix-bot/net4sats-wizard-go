# Frame-level branding audit — OCR over the recorded take

Artifact audited: `artifacts/take-attempt-blocked.mp4`
Method: `ffmpeg -i ... -vf fps=2` -> 86 PNG frames -> `tesseract <frame> -` per frame,
then case/whitespace-normalised substring search for the token `net4sats`.

## Result

| Metric | Value |
|---|---|
| Frames OCR'd | 86 (2 fps over 42.88 s) |
| Frames containing `net4sats` | **83** |
| Frames clean of `net4sats` | 3 |
| Verdict | **BRANDING PRESENT — acceptance criterion 3 is NOT met** |

This is a measurement on rendered pixels, not on source. It is the check the card
asks for ("confirm every branded string/logo is absent from the frames"), and it
fails.

## Where the branding appears on screen

Every hit is the same element: the **page header**. Full OCR of the first and last
frames:

```
$ tesseract f0001.png -          $ tesseract f0086.png -
net4sats                         net4sats
                                 
Bitcoin WiFi — Setup your         Bitcoin WiFi — Setup your
router in 60 seconds             router in 60 seconds
                                 
d                                Setup failed

Scanning for routers on          Cannot connect to router via SSH
your network...
```

So the wordmark sits in the top-left header on **every** screen of the wizard —
start, config, deploy progress, and the terminal state. It cannot be cropped out
of a full-window recording without cutting the UI itself.

Additional branded strings that exist in the code and render during a *successful*
run (not reached in this take, but they are in the same built binary — verified by
`strings`):

- `Deploy net4sats` (primary button)
- `Deploying net4sats...` (progress heading)
- `net4sats is live!` (success headline)
- deploy step labels: `Installing net4sats package + patching backend...`,
  `Branding captive portal as net4sats...`, `Deploying net4sats captive portal...`,
  `Installing net4sats admin panel...`

## Binary-level cross-check (not just the 4 source strings)

```
$ strings /tmp/wiz-wz14 | grep -ci net4sats
120
```

The compiled artifact the recorder rendered contains 120 `net4sats` occurrences,
including all of the visible UI strings above. This rules out the possibility that
the source grep and the running binary disagree.

## Repo-wide check: no clean ref exists

```
$ for ref in $(git for-each-ref --format='%(refname:short)' refs/heads refs/remotes); do
    echo "$(git show "$ref:index.html" 2>/dev/null | grep -ci net4sats)  $ref"; done | sort -n | head
4  origin/pr32
4  origin/worker-heavy/t_c92a4003
4  origin/wz14-debrand
4  origin/wz19-video-evidence
4  origin/wz19-video-evidence-pre-atomic
4  wz19-video-evidence
5  ci/add-github-actions-workflow-clean
5  ci/github-actions-workflow-atomic
```

Minimum across the entire repository is 4. **Zero refs are clean.** There is no
branch a de-branded recording could have been taken from.

## Why this is not a tooling failure

The capture path is sound — it produced a valid 1920x1080 H.264 take of the real
UI. The branding is in the product, not in the recorder. Re-running the same
toolchain against any current ref will reproduce this table.
