# PROGRESS — Release v0.7.0-alpha21 on fork + curl-verify asset

## Finding (BLOCKER: tag name is stale, not destination-consistent)

`v0.7.0-alpha21` **already exists on the fork** — tag, GitHub Release, and a
correct `net4sats-wizard` asset. It is NOT an unpublished version:

- Tag: `6c7f78035560437df7ba5d4cfa3d1e9ec3b59ef1` "feat: pre-flight WiFi password
  validation for STA setup" (2026-09-07 14:57 +0530)
- Remote tag `refs/tags/v0.7.0-alpha21` = `6c7f780` (peeled → commit, version-named,
  NOT a raw-SHA tag — the task's stated fear does not apply to this tag)
- Release: createdAt 2026-09-07T09:27:00Z, name `v0.7.0-alpha21`, published
  (isDraft=false)
- Asset: name `net4sats-wizard`, size 14748994, state uploaded,
  downloadCount=5, `createdAt` 09:29:29Z / `updatedAt` 09:29:33Z

`v0.7.0-alpha21` is **NOT an ancestor of HEAD** — it is 83 commits out of tree, and
so is alpha24. HEAD (`22fbea4`) sits on the docs/worker branch `worker-heavy/t_gate8-19b1b1c6`.

## The conflict

Task says "Publish the **tested** wizard as a release ... v0.7.0-alpha21".
But alpha22, alpha23 and alpha24 are **already published after alpha21**
(alpha22 2026-09-07T11:58Z, alpha23 12:59Z, alpha24 2026-10-01T11:35Z), and
alpha24 is currently GitHub `Latest`. The task's notion of "the tested release
version" is **2 weeks and 3 releases stale**.

Consequences:
- Publishing alpha21's content as "the tested release" is not possible — that
  exact artifact was published 2026-09-07.
- Re-tagging alpha21 is **destructive**: the 83 commits since, including
  alpha22/23/24 (security fixes) would no longer be reachable from the head of
  that line. Prior content is only recoverable from the reflog.
- No HEAD artifact may be published under a tag that already carries other
  content.
- alpha24's asset is a **smaller** binary (16,946,325) than alpha21's
  (14,748,994). Neither is an ancestor of HEAD, so the "tested" lineage is
  ambiguous from the repo alone.

## Status: BLOCKED on operator decision (Step 1)

Cannot start Step 1 without knowing which version identity is correct.

## Steps verified THIS far (no mutation performed)

- Step 3-equivalent verification of the EXISTING alpha21 asset (read-only):
  - HTTP 200, `bytes=14748994`, magic `7f 45 4c 46` (\x7fELF),
    `ELF 64-bit LSB executable, x86-64`, sha256
    `a5dc66c5ee894c1137a100a39875d74a0cceec5fc1ff86368cf04be3b99262fb`
- Step 4 (pre-flight): `:8099` and `:8109` FREE; no stray `net4sats-wizard`
  process (only this Hermes session matched pgrep).

## Files touched

- `PROGRESS-t_release-alpha21.md` (new, this file)
- `REPORT-t_release-alpha21.md` (new)
