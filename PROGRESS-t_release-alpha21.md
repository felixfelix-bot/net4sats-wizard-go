# PROGRESS — Release v0.7.0-alpha21 on fork + curl-verify asset

**Card:** `t_2f29676c` (resolved live from the board DB; escalation card `t_f9d1db3d`)
**Status: BLOCKED before Step 1. NO MUTATION PERFORMED.** Full report: `REPORT-t_release-alpha21.md` @ commit `7625acf`.

## Finding (BLOCKER 1: the brief's tag is stale, not unpublished)

`v0.7.0-alpha21` **already exists on the fork** — tag, GitHub Release, and a
correct `net4sats-wizard` asset. It is NOT an uncut version:

- Tag: `6c7f78035560437df7ba5d4cfa3d1e9ec3b59ef1` "feat: pre-flight WiFi password
  validation for STA setup" (2026-09-07 14:57 +0530)
- Tag type: **`commit`** — correctly version-named. The brief's stated fear ("a past
  finisher mistakenly tagged the raw SHA") **does not apply**: both the local and
  remote `refs/tags/v0.7.0-alpha21` peel to a commit, not a raw object.
- Release: createdAt 2026-09-07T09:27:00Z, `isDraft=false`, published
- Asset: name `net4sats-wizard`, size 14748994, `state:uploaded`, `downloadCount:5`

alpha22 / alpha23 / alpha24 are later **linear descendants**, and alpha24 is
GitHub `Latest`. So the brief's tag was **3 releases stale when written**.

Card created **2026-10-06T10:00:54Z** — a month *after* alpha21 shipped, five days
after alpha24. Upstream cards `t_93d60eb5` (bump to alpha21) and `t_4c470aa0`
(release alpha21) are **both still `blocked`**; alpha21 was cut out-of-band while
they sat blocked. Stale-brief / double-publish, not a never-cut release.

Not done unilaterally: re-pointing alpha21 is a **destructive history rewrite** —
it would overwrite a release already downloaded 5x and leave alpha22/23/24
dangling. (Unlike the safe n→n+1 merge case, this moves n *onto* newer content.)

## Finding (BLOCKER 2: no HEAD artifact can be called "tested")

`go test ./...` at HEAD **FAILS** (exit 1): `TestNoInternalProcessLeaksOnShippedSurfaces`,
105 "internal detail leaked" violations (internal card ids, fleet hostnames
`c08r4d0r`/`c03rad0r`/`CobradorWave`, board path `hermes/kanban`) all from one tracked
file — `docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md` — added by HEAD
itself (`22fbea4`) and **absent from `v0.7.0-alpha24`**. Publishing HEAD would ship
105 leaks. Go is not on `PATH`; it is at `/usr/local/go/bin/go` (go1.22.12).

## Verified THIS run

- **Step 3-equivalent, read-only, of the EXISTING asset — VERIFIED LIVE:**
  HTTP 200, magic `7f 45 4c 46` (\x7fELF), `ELF 64-bit LSB executable x86-64`,
  **14748994 bytes**, sha256
  `a5dc66c5ee894c1137a100a39875d74a0cceec5fc1ff86368cf04be3b99262fb`.
  Fetched and confirmed — not an assumed URL (9-byte "Not Found" lesson respected).
- **Step 4:** `:8099` FREE, `:8109` FREE, no stray `net4sats-wizard` process.

## Board actions

- Comment posted to `t_2f29676c` (→ "Comment added", exit 0; card now has 4 comments).
- Escalation card `t_f9d1db3d` created (`blocked`, assignee `manager`, priority 1).
- Status left `triage` — **NOT set to review**. `review` would assert a publishable
  artifact exists; it does not. The CLI also cannot transition a `triage` card from a
  session (`claim`/`block`/`unblock`/`promote` all refused), so it needs operator or
  manager intervention.

## Files touched

- `PROGRESS-t_release-alpha21.md` (new, this file)
- `REPORT-t_release-alpha21.md` (new)
- `/tmp/alpha21probe`, `/tmp/gotest-head.txt` (scratch, outside repo)

**No tag, branch, release, or asset was created or modified. No push beyond the
docs commit. Nothing destructive was executed.**
