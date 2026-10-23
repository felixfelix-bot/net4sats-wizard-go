# PROGRESS — t_93620520 consolidated identity-API + BIP39 evidence note (worker-heavy, 2026-10-23)

Task: fuse the two verification strands into the single citable evidence note with verdict table.
Finding -> Status -> Files touched. Prior sections below belong to earlier tasks sharing this
worktree root (preserved whole).

1. Finding: NEITHER parent task arrived with a result — the brief's premise is unmet.
   Status: VERIFIED from the board DB, read-only.
   Evidence: `boards/net4sats-mvp-v2/kanban.db` → `t_f047773e` (identity-API, worker-heavy) and
   `t_fd85fd04` (BIP39 ground truth, worker-tollgate) both have `length(result)=0`; the BIP39 card
   is `archived` with 3 failed runs (`reclaim backoff … dead worker pid` ×2, then
   `block_loop_detected`, loop guard). The consumer card `t_93620520` ("Write consolidated evidence
   note with verdict table") is also `triage`/`result=NULL` — which is precisely the card this
   session's brief describes.
   Consequence: I re-derived both strands from disk rather than consuming 0-byte inputs.
   Files touched: none (read-only).

2. Finding: the two strands' substance DOES exist on disk; the parents' failure was `result`
   writeback, not missing work.
   Status: VERIFIED.
   Evidence: Part A = `REPORT-t_31a002af.md` (128 lines, this worktree root) — the identity-API
   adjudication, written against wizard `d5d2373` / backend `b11db6a`. Part B = the adversarial
   audit on card `t_e7c505b4`
   (`attachments/t_e7c505b4/AUDIT-evidence-note-t_e7c505b4.md`, finding 2), which re-verified the
   BIP39 ground truth first-hand after the BIP39 strand died, and which FAILED a previous attempt
   at this very note.
   Files touched: none (read-only).

3. Finding: every citation carried forward was re-opened at CURRENT HEADs, because both repos moved.
   Status: VERIFIED first-hand.
   Evidence: wizard HEAD `18de0779` (AGENTS.md blob `f2ff1775` — unchanged since the single commit
   `475a9ef`, so the contested claims are byte-identical). Backend working tree `ffc7c433` (branch
   `pr/cdk-mint-issuance`) has moved from the parent's `b11db6a`; route lines are now
   `src/main.go:1935` / `:1941`. The parent's `CHANGELOG.md:2515-2518` cite does NOT resolve on the
   current branch — the identity entry is at `origin/main:CHANGELOG.md:3169-3174`. I used my own
   line numbers throughout.
   Files touched: none (read-only).

4. Finding: headline verdicts — GET /identity called by the code = FALSE; POST
   /identity/reveal-seed called by the code = FALSE; mnemonic 12 words (128-bit) = TRUE while the
   wizard's `AGENTS.md:43` "24 words" = FALSE.
   Status: VERIFIED. Both routes are SERVED by the backend (`src/main.go:1935`, `:1941`) and
   CALLED by nothing in the wizard repo: repo-wide `reveal-seed|/identity` hits only
   `AGENTS.md:29,30,46,47`; the wizard's whole HTTP surface is `main.go:619-623`, its only outbound
   request is `deploy.go:1130`. The 12-word pin is `bip39.NewEntropy(128)` at
   `src/identity/identity.go:137`, asserted at `src/identity/identity_test.go:57`/`:88`, library
   `go-bip39 v1.1.0` (`src/go.mod:115`).
   Files touched: EVIDENCE-NOTE-t_93620520.md.

5. Finding: the note's citations are machine-verified, not asserted.
   Status: DONE, observed. A parser extracted every ``token:line[-line]`` from the note and resolved
   each against the working trees, `origin/main` (via `git show`) and the Go module cache:
   195 resolved, 0 out of range, 0 missing files. A further 24 hand spot-checks confirmed the
   verbatim snippets match the lines they are attached to. Three claims are labelled
   `unverified - no citation available` (upstream configurationwizzard source; `AGENTS.md:22` port
   claim; v0.5.0-alpha3 backward-compat scope) and one parent citation is recorded as unresolvable
   rather than laundered.
   Files touched: EVIDENCE-NOTE-t_93620520.md.

6. Finding: an earlier attempt at this note was audited FAIL (t_e7c505b4) for an uncited 12-word
   assertion and for candidate statements that would inject new falsehoods.
   Status: fixed by construction here. The 12-word claim now carries five citations plus an explicit
   caveat that the ENDPOINT accepts any valid BIP39 length (`IsMnemonicValid`,
   `go-bip39@v1.1.0/bip39.go:275-278`, lengths at `:29-35`); the ubus attribution is corrected
   (`deploy.go:134-139`, `:868`); the `GET /identity` verb-enforcement overstatement is called out
   (`src/main.go:2121-2132` has no `r.Method` check).
   Files touched: EVIDENCE-NOTE-t_93620520.md.

## Status

DONE. Deliverable committed and pushed.

FILES TOUCHED
- EVIDENCE-NOTE-t_93620520.md   (the deliverable note; committed `824ca17`, pushed to
  `origin/worker-heavy/t_gate8-19b1b1c6` — first push was rejected non-fast-forward because a
  sibling lane advanced the remote, the remote tip `c9adde4e` was verified to be an ancestor of
  my HEAD, and the re-push was observed at `824ca17`)
- PROGRESS-t_93620520.md, REPORT-t_93620520.md   (this worktree's task-scoped prose)

NOT DONE (out of scope by the brief)
- The note was NOT pasted into card `t_93620520`'s result field — that card lives on the board and
  belongs to the dispatcher; the note is on disk and in git for the parent to lift.
- `AGENTS.md` was NOT edited. The brief is read-only and the fix is a separate card
  (`t_b44625ff`, `t_72f4489d`).
