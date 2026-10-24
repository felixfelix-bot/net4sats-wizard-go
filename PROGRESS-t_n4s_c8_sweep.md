# PROGRESS — C8 archive sweep (Phase 5, LAST) — card t_n4s_c8_sweep

Checklist: /home/c08r4d0r/reports/net4sats-mvp/ARCHIVE-SWEEP-READY.md
Worktree: /home/c08r4d0r/repos/net4sats-mvp-v2 @ worker-heavy/t_gate8-19b1b1c6

## F1 — The checklist is STALE; the state has already moved (measured, not assumed)

| Checklist says | Live reality (2026-10-24, gh api) |
|---|---|
| net4sats-feed "0 open PRs, ✅ READY" | **PR #1 is OPEN** (`pr/archive-pointer`, MERGEABLE, ad59d3a) |
| tollgate-captive-portal-site #1+#2 open | **both CLOSED already** |
| net4sats-captive-portal-site #4 open | **CLOSED already** |
| net4sats-wizard-go #32/#33/#34 open, #35 to close | #32/#33/#34 OPEN; **#35 already CLOSED** |
| configurationwizzard #13 open | OPEN, authored by **c03rad0r** (not us), **CONFLICTING/DIRTY** |

Status → recorded. No work done off the stale table; every action below is
against the measured state.

## F2 — THE BLOCKER: this account cannot archive anything

`gh auth status` → logged in as **felixfelix-bot**. Measured permission on
**every** net4sats/* repo in scope:

    net4sats-feed                  perm=READ  admin=false  push=false
    tollgate-module-basic-go       perm=READ  admin=false  push=false
    gonuts-tollgate                perm=READ  admin=false  push=false
    configurationwizzard           perm=READ  admin=false  push=false
    tollgate-captive-portal-site   perm=READ  admin=false  push=false
    net4sats-captive-portal-site   perm=READ  admin=false  push=false
    net4sats-wizard-go             perm=READ  admin=false  push=false
    net4sats.github.io             perm=READ  admin=false  push=false
    net4sats/installer             perm=ADMIN admin=true            <-- only admin repo

Empirically confirmed, not just inferred from the flag:

    $ git push --dry-run https://github.com/net4sats/net4sats-feed.git HEAD:refs/heads/worker/probe-perm
    remote: Permission to net4sats/net4sats-feed.git denied to felixfelix-bot.
    fatal: ... The requested URL returned error: 403

`gh repo archive net4sats/<repo>` therefore cannot succeed — archiving needs
admin, and the pointer-README commit that MUST land before archiving needs push.
Both are absent. The already-staged `OPERATOR-ARCHIVE-RUNBOOK.sh` in
`reports/net4sats-mvp/archive-pointers/` opens with exactly this note
("felixfelix-bot lacks push+admin on these — run from an account with admin"),
so this is a known, previously-identified operator gate, not a new discovery.

Status → BLOCKED on operator credentials for every archive step (checklist
steps 1-5). Escalated rather than retried: `block_recurrences: 2` on this card
already records two prior `fleet-offload:x280` blocks, i.e. re-dispatching to
the same read-only identity is the loop that keeps failing.

## F3 — What the token CAN do: PR closes (author-scoped, not repo-scoped)

Closing your OWN PR does not need push/admin. Measured per-PR:

    net4sats-wizard-go #32  canClose=true  canUpdate=true  (cross-repo, head felixfelix-bot)
    net4sats-wizard-go #33  canClose=true  canUpdate=true
    net4sats-wizard-go #34  canClose=true  canUpdate=true
    net4sats-feed      #1   canClose=true  canUpdate=true
    configurationwizzard #13 canClose=FALSE canUpdate=FALSE  <-- c03rad0r's PR

Fork `felixfelix-bot/net4sats-wizard-go` is `admin=true push=true`, so branch
preservation (cut `archive/*` tags) IS available.

## F4 — Branch preservation audit (checklist: "preserve branch")

Fork already carries an `archive/*` tag convention (20+ tags) plus:

    refs/heads/wz14-debrand                  b64555a5  = PR #32 head   ✓
    refs/heads/pr/cors-write-path            fbac751d  = PR #33 head   ✓
    refs/heads/ci/add-github-actions-workflow 87938089 = PR #34 head   ✓
    refs/heads/ci/add-github-actions-workflow-clean 9b272d6c = PR #35 head ✓
    refs/heads/pr/archive-pointer (feed)     ad59d3a7  = feed PR #1 head ✓
    refs/tags/archive/wz14-debrand           2d770966  (ANCESTOR of PR#32 head
                                                       b64555a5 — tag is STALE,
                                                       branch carries 2 more commits)
    refs/tags/archive/pr-cors-write-path     fbac751d  = PR #33 head   ✓
    (no tag for #34, none for #35, none for feed #1)

Status → every PR head survives as a fork BRANCH today, so the "preserve
branch" requirement is satisfied by the branches themselves. Two tag gaps
(#34, feed #1) plus the stale `archive/wz14-debrand` tag were found and FIXED
(additive only) — see F5.

## F5 — Branch preservation HARDENED (done, verified on remote)

Annotated tags cut and pushed. The PRE-EXISTING `archive/wz14-debrand`
(`2d770966`) was deliberately NOT re-pointed — it is an ancestor of the PR head
and `UPSTREAM-PINS.md` forbids moving a tag; a new name was cut instead.

    felixfelix-bot/net4sats-wizard-go
      archive/ci-github-actions-workflow        -> 8793808950c0  (= PR #34 head)
      archive/ci-github-actions-workflow-clean  -> 9b272d6c0071  (= PR #35 head)
      archive/wz14-debrand-final                -> b64555a5efff  (= PR #32 head, FULL tip)
      (untouched: archive/wz14-debrand          -> 2d770966      (ancestor, stale))
    felixfelix-bot/net4sats-feed
      archive/feed-archive-pointer              -> ad59d3a76a32  (= feed PR #1 head)

Status → VERIFIED by `git ls-remote <fork> 'refs/tags/<t>^{}'` peeled readback
(annotated-tag object SHAs differ from commit SHAs by design; the `^{}` peeled
values above are the readback, not a local read).

Net effect: after the closes in F6, EVERY branch/commit from every PR in this
sweep is reachable from a durable, already-pushed tag on a fork where this
account holds admin. Nothing depends on GitHub retaining `refs/pull/*`.

