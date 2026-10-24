# REPORT — Release v0.7.0-alpha21 on fork + curl-verify asset

**Card:** offload task "Release v0.7.0-alpha21 on fork and curl-verify asset"
**Run:** 2026-10-24, host `x280` (`c08r4d0r`)
**Workspace:** `/home/c08r4d0r/repos/net4sats-mvp-v2`
**Branch:** `worker-heavy/t_gate8-19b1b1c6` (HEAD `22fbea4`)
**Fork (only writable target):** `felixfelix-bot/net4sats-wizard-go` (`origin`, HTTPS, gh auth OK as `felixfelix-bot`)

## VERDICT: **BLOCKED BEFORE STEP 1. NO MUTATION PERFORMED.**

The task cannot be executed as written, and executing it literally would be
destructive. Two independent hard blockers, both measured first-hand.

---

## Blocker 1 — `v0.7.0-alpha21` is not an unpublished version; it is 3 releases stale

The task says "Publish **the tested** wizard as a release ... TAG NAME must be
EXACTLY `v0.7.0-alpha21`". That tag is already published, and not recently:

| Fact | Measured value |
|---|---|
| Local tag object | `git rev-parse v0.7.0-alpha21` → `6c7f78035560437df7ba5d4cfa3d1e9ec3b59ef1` |
| Tag type | `git cat-file -t` → **`commit`** (an annotated/lightweight tag pointing at a commit) |
| Tag commit | `6c7f780 feat: pre-flight WiFi password validation for STA setup` — 2026-09-07 14:57:00 +0530 |
| Remote tag | `refs/tags/v0.7.0-alpha21` = `6c7f780…` (peeled → commit) |
| Release | `gh release view v0.7.0-alpha21` → createdAt **2026-09-07T09:27:00Z**, `isDraft:false`, `isPrerelease:false` |
| Asset | name **`net4sats-wizard`** (exactly as specified, no suffix), size **14748994**, `state:uploaded`, `downloadCount:5`, createdAt 09:29:29Z / updatedAt 09:29:33Z |

**The task's stated failure mode does not apply here.** The brief warns "a past
finisher mistakenly tagged the raw SHA; do not repeat that". Both the local and
remote `v0.7.0-alpha21` refs point at a **commit**, not a raw/unpeeled object —
this tag is correctly version-named. There is nothing to repair.

**It is also not the tip of the release line.** alpha22, alpha23 and alpha24
were all published *after* alpha21, and each is a linear descendant:

```
alpha21 -> alpha22 : YES (linear ancestor)
alpha22 -> alpha23 : YES
alpha23 -> alpha24 : YES
```

| Tag | Published | Asset size |
|---|---|---|
| v0.7.0-alpha21 | 2026-09-07T09:27Z | 14,748,994 |
| v0.7.0-alpha22 | 2026-09-07T11:58Z | — |
| v0.7.0-alpha23 | 2026-09-07T12:59Z | — |
| **v0.7.0-alpha24** | **2026-10-01T11:35Z** | 16,946,325 (GitHub `Latest`) |

Consequence: there is no coherent reading in which "publish *the tested* wizard
as v0.7.0-alpha21" is a forward action.

- If it means *publish the alpha21 artifact* — it was published 2026-09-07 and is
  live right now (verified below). Nothing to do.
- If it means *re-point alpha21 at newer content* — that is a **history rewrite
  of a published release**. `v0.7.0-alpha21..v0.7.0-alpha24` is a 4-commit chain
  (alpha22 → alpha24); moving alpha21 forward would leave alpha22/23/24 as
  dangling, unreachable tags, and would overwrite a release whose asset has
  already been downloaded 5 times. My brief forbids unrecoverable operations.
- If it means *publish HEAD as alpha21* — HEAD is not in that lineage at all:
  alpha21, alpha23 and alpha24 are **all** non-ancestors of HEAD
  (`git merge-base --is-ancestor` → NO for each; both alpha21 and alpha24 are
  83 commits away from HEAD). Publishing HEAD under a tag that already holds
  different content would misrepresent both.

I did **not** pick one of these unilaterally. Choosing wrong is either a no-op
or a destructive rewrite.

## Blocker 2 — no artifact published from this tree could honestly be called "tested"

The task's acceptance list requires "the full go test output". At HEAD, that
output is a **FAILURE**:

```
$ export PATH=$PATH:/usr/local/go/bin   # go1.22.12 linux/amd64
$ go test ./...
--- FAIL: TestNoInternalProcessLeaksOnShippedSurfaces (0.11s)
FAIL
FAIL	net4sats-wizard	74.297s
FAIL
EXIT=1
```

One test fails, with **105** `internal detail leaked` violations, all sourced
from a single file:

```
docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:193: internal detail leaked (internal card id): "t_a8b5fcae"
docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:218: internal detail leaked (personal fleet hostname): "c08r4d0r"
docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:231: internal detail leaked (personal fleet hostname): "c03rad0r"
docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:411: internal detail leaked (internal board path): "hermes/kanban"
... (105 total)
```

That file is tracked, and it was added by **HEAD itself**
(`22fbea4 docs(wz1.8): readiness note for GL-MT6000 wizard E2E — VERDICT BLOCKED (t_05167730)`).
It is **absent from `v0.7.0-alpha24`**, so the failure is new work on this
branch, not a pre-existing release-branch bug. A release cut from HEAD would
ship 105 internal-detail leaks and a red suite. I will not label that "tested".

(Note: Go is NOT on `PATH` in this session — it lives at `/usr/local/go/bin/go`.
Anyone re-running this must export that path first.)

---

## Steps NOT performed

1. **Push branch + create tag `v0.7.0-alpha21` — NOT DONE.** Branch is already
   pushed and in sync (`git log origin/worker-heavy/t_gate8-19b1b1c6..HEAD` →
   empty). No tag created: the tag already exists with different content, and
   re-pointing it is destructive (Blocker 1).
2. **Build/attach asset `net4sats-wizard` — NOT DONE.** No build was run; no
   upload attempted. Uploading HEAD's binary to alpha21 would overwrite a
   downloaded-5x asset; uploading it anywhere else needs a version decision.
3. **Curl-verify — DONE, but against the EXISTING asset** (read-only, no
   mutation). Output below. This is a verification of the artifact that is live
   today, not of a new upload by this run.
4. **`:8099` free + no stray wizard process — DONE.**
5. **Status = review — NOT DONE.** Reporting BLOCKED instead. Setting `review`
   would imply a publishable artifact exists; it does not.

## Step 3 (read-only) — curl-verify of the EXISTING v0.7.0-alpha21 asset

```
$ curl -sIL https://github.com/felixfelix-bot/net4sats-wizard-go/releases/download/v0.7.0-alpha21/net4sats-wizard
final_http=200
last-modified: Mon, 07 Sep 2026 09:29:33 GMT
etag: "0x8DF0CC286DE948C"
server: Windows-Azure-Blob/1.0 Microsoft-HTTPAPI/2.0
```

```
$ curl -sL https://github.com/felixfelix-bot/net4sats-wizard-go/releases/download/v0.7.0-alpha21/net4sats-wizard -o /tmp/alpha21probe
HTTP=200 bytes=14748994

$ head -c 4 /tmp/alpha21probe | xxd
00000000: 7f45 4c46                                .ELF

$ file /tmp/alpha21probe
ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked,
interpreter /lib64/ld-linux-x86-64.so.2, Go BuildID=Zst4d9vlaI7sp7X-aTT-/…,
BuildID[sha1]=38f03857d88abcf576c5e89486492e6ab8535795, with debug_info, not stripped

$ stat -c '%s %n' /tmp/alpha21probe
14748994 /tmp/alpha21probe

$ sha256sum /tmp/alpha21probe
a5dc66c5ee894c1137a100a39875d74a0cceec5fc1ff86368cf04be3b99262fb
```

HTTP 200 + `\x7fELF` magic + 14,748,994 bytes, matching the API-reported asset
size exactly. **This is a verified-live URL.** Contrast with the "9-byte 'Not
Found'" lesson: that failure was reporting a URL without fetching it. This URL
was fetched, and the fetch returned real ELF content. Had I uploaded blind and
reported, I would have risked exactly that error.

## Step 4 — port and process hygiene

```
$ ss -ltnp | grep -E ':8099|:8109'
8099/8109 FREE

$ pgrep -af net4sats-wizard
(no stray wizard process — only this Hermes session's own command line matched)
```

`:8099` free, `:8109` (the wizard's documented fallback port) free, no zombie
wizard. Clean.

## Full `go test ./...` output (the acceptance item)

```
$ export PATH=$PATH:/usr/local/go/bin && go version
go version go1.22.12 linux/amd64

$ go test ./...
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:193: internal detail leaked (internal card id): "t_a8b5fcae"
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:194: internal detail leaked (internal card id): "t_a8b5fcae"
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:202: internal detail leaked (internal card id): "t_77c59d23"
    ... [105 leak lines, all from the same doc] ...
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:411: internal detail leaked (internal board path): "hermes/kanban"
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:415: internal detail leaked (personal fleet hostname): "CobradorWave"
    debrand_internal_test.go:120: docs/readiness-gl-mt6000-wz1.8-t_05167730-2026-10-24.md:419: internal detail leaked (internal card id): "t_63ec0ed2"
--- FAIL: TestNoInternalProcessLeaksOnShippedSurfaces (0.11s)
FAIL
FAIL	net4sats-wizard	74.297s
EXIT=1
```

Full raw log preserved at `/tmp/gotest-head.txt`.

## `git log --oneline -5`

```
22fbea4 docs(wz1.8): readiness note for GL-MT6000 wizard E2E — VERDICT BLOCKED (t_05167730)
a380e10 docs(lock): t_0841f821 — progress + full report for the GL-MT6000 lock announcement
ab33306 docs(lock): t_0841f821 — public GL-MT6000 lock announcement (owner t_be3261d3, release gate t_557bb77d)
fda8cea docs(wz1.9): baseline full-suite capture + failing-test inventory (t_gate8-19b1b1c6)
b56fae5 docs(demo): t_e07e91b1 — full report for the assembled draft master
```

## What is needed to proceed

**Operator decision required** — pick the intended version identity:

- **(A)** alpha21 is already done; the real intent is to release HEAD content
  under a **new** tag (alpha25) — needs the debrand failure fixed first.
- **(B)** alpha21 must genuinely be re-published with current content — this is
  an explicit, confirmed history rewrite of a 5-download published release,
  leaving alpha22/23/24 dangling. I will not do this without an explicit order.
- **(C)** The target was a **different** version and `v0.7.0-alpha21` in the brief
  is a stale copy-paste (alpha24 is `Latest`).
- **(D)** Ship alpha21 as-is and simply report the already-verified asset —
  no work needed, but then the task's "publish" framing is wrong.

Regardless of choice: **`TestNoInternalProcessLeaksOnShippedSurfaces` must go
green before any HEAD-derived release**, otherwise the release ships 105
internal-detail leaks.

## Files touched by this run

- `PROGRESS-t_release-alpha21.md` (new)
- `REPORT-t_release-alpha21.md` (new, this file)
- `/tmp/alpha21probe` (re-downloaded asset probe, outside repo)
- `/tmp/gotest-head.txt` (raw test log, outside repo)

**No source, test, tag, branch, or release was modified. No push. No tag. No
upload.** The only writes are the two documentation files above and two `/tmp`
scratch files.
