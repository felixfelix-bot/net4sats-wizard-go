# PROGRESS — t_1df48589 (wz1.5 channel decision note + Q5 status) — second pass

## Finding

**The deliverable exists and is correct, but was never committed — and it contained one real
defect, now fixed.**

1. **The note was written by an earlier pass of this card but never landed in git.**
   `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md` is on disk (22.7 KB, 414
   lines) and the plans directory is **not a git repo** (`git rev-parse --show-toplevel` → fatal),
   while the earlier `PROGRESS-t_1df48589.md` claims it "cannot be committed". The card's own
   report/PROGRESS files were also **not** in history — `git log --all -- PROGRESS-t_1df48589.md
   REPORT-t_1df48589.md` → **0 commits** — because `/PROGRESS-t_*.md` and `/REPORT-t_*.md` are
   gitignored in this worktree (`/home/c08r4d0r/wz16-verify/.git/info/exclude:4-5`), unlike earlier
   cards that landed with `git add -f`. So the decision was durable only as an unversioned file in
   a home directory.
2. **Defect found and fixed on re-derivation.** The note's instantiated artifact-filename pattern
   (`tollgate-wrt_v<X.Y.Z>_<arch>.ipk`) is **404** on the one version that exists on the chosen
   host — the real `.ipk` carries a build-flavour suffix
   (`tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk`) and no plain `.ipk` exists
   at all, while the **apk** exists in *both* plain (7,153,725 B) and suffixed (5,307,464 B) form.
   The filename is therefore per-artifact and must be read from that version's own index
   `Filename:` field, never composed. New §2.3a carries the evidence; §3 rule 1b enforces it; §1
   and §6.1 now say so explicitly.
3. **Every other load-bearing claim was independently re-verified this pass** (2026-10-24) and
   still holds: `releases.tollgate.me` 206; all seven probed tag versions 404 for indexes except
   `v0.5.0-alpha3`; `Packages.sig` and `keys/tollgate-testing.pem` 404; both target literals 404;
   `releases/download/latest/…` 404 and `/releases/latest` → 302 → HTML tag page; release
   inventories unchanged (`FreedomTechFeed/packages` v0.6.0 + v0.6.0-rc6, 26 assets each; upstream
   `OpenTollGate/tollgate-module-basic-go` v0.6.0 with **0 assets**; neither wizard-pinned tag
   exists on `FreedomTechFeed`, API 404 each); one full download of the suffixed apk (5,307,464 B)
   matched its index `SHA256sum`; and a 27-URL sweep across 9 tags confirmed **exactly one**
   version carries artifacts on the chosen host.

Decision unchanged: interim channel `releases.tollgate.me`, **Q5 = OPEN**.

## Status

**DONE — decision note corrected and delivered; Q5 recorded OPEN.** No Go source touched
(`deploy.go` and `pins_test.go` unmodified — the four-point landing gate stays closed), so nothing
was built or tested and no build artifact's behaviour can have changed. `git diff --stat` over the
repo's Go files is empty.

Blocked on the same two things as before: the producing side publishing a signed index and a
per-arch artifact for a current release (note §6.2 gates 1–3, all currently failing), and the
operator/architecture answer to Q5.

## Files touched

- `~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md` — **the deliverable**
  (patched: header re-verification line, §1 correction, new §2.3a, §3 rule 1b, §6.1 caveat, §7
  acceptance rows). Canonical location per the card spec; it sits outside any git repo, which is
  the durability gap this pass fixes below.
- `docs/wz1.5-channel-decision.md` — **durable, versioned mirror** of the deliverable, committed
  in this repo (the plans directory cannot be committed: `git rev-parse --show-toplevel` → fatal).
  Byte-identical to the canonical file at commit time; the canonical file remains the deliverable
  path the card names.
- `docs/evidence/2026-10-24-wz15-channel-decision-reverify.md` — verbatim captures for this
  second pass, following the sibling card's `docs/evidence/<date>-<name>.md` convention.
- `REPORT-t_1df48589.md` — full report (updated)
- `PROGRESS-t_1df48589.md` — this file (updated)
- `PROGRESS.md` — appended one section (shared file; pathspec repeated on commit)
- `REPORT.md` — **intentionally untouched**: it carries a sibling card's uncommitted edits
  (`git diff --stat REPORT.md` → 32 insertions / 22 deletions, not mine).
- `deploy.go`, `pins_test.go` — **untouched** (read-only inspection)

## Verification performed

- Presence/absence: `read_file`, `ls -la` of the plans dir, repo-wide `grep -rn`, `find`.
- Git state: `git ls-files`, `git check-ignore -v`, `git log --all -- <paths>`, `git status -sb`,
  `HEAD` vs `origin/branch` (both `7998fc3`).
- Live probes (2026-10-24): 12-URL status matrix, 27-URL tag×filename sweep, `Packages` field
  extraction, `curl -sSI` content-length cross-check against each `Size:`, one full artifact
  download with `sha256sum` vs index `SHA256sum:`, GitHub releases API inventories, `/releases/latest`
  redirect target, and the explorer bundle re-fetch (474,543 B) for the channel vocabulary.
- Environment: disk was at **935 MB free (100% used)** at the start of this pass — every download in
  this pass was ≤ 5.4 MB and the scratch file was deleted immediately after hashing.
