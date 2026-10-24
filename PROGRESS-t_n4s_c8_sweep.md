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

## F6 — Criterion #3 read-back PASSES ⇒ closes are honest, not premature

The checklist gates the whole sweep on "criterion #3 read-back passes —
archiving before replacements proven removes the only working install path".
Measured today, all three legs:

    net4sats/installer exists, main, pushed_at 2026-10-24, not archived
    https://raw.githubusercontent.com/net4sats/installer/main/install.sh  -> HTTP 200
    install.sh header: selects a brand, forwards argv to the UPSTREAM launcher,
      verifies vendor/ against scripts/upstream-launcher.sha256 before executing
      any byte, and states it contains "no discovery, SSH/trust, deploy, feed
      resolution, flashing, identity derivation or pricing"  (AGENTS.md-compliant)
    pins.tsv (the pin manifest, DATA) + every pin resolves on the remote:
      net4sats-pin/installer-comfast        b042a2396425  == recorded  ✓
      net4sats-pin/installer-gateway-domain a72a4f6698f7  == recorded  ✓
      net4sats-pin/installer-rebrand-gutter 049b68918684  == recorded  ✓
      net4sats-pin/module-ssid-bang         41ada96054bd  == recorded  ✓
      launcher sha256 pin present: 6595640e29e4...

Status → the replacement install path is proven live and pinned BEFORE anything
is closed/archived, which is exactly the ordering the checklist demands.

## F7 — PR closes performed (the part this identity CAN do)

Closing your OWN PR is author-scoped and needs no repo write, so these went
through. Closes verified by reading `state`/`closedAt` back:

    net4sats-wizard-go #32  -> CLOSED 2026-10-24T05:54:09Z
    net4sats-wizard-go #33  -> CLOSED 2026-10-24T05:53:59Z
    net4sats-wizard-go #34  -> CLOSED 2026-10-24T05:53:45Z
    (net4sats-wizard-go #35 was ALREADY closed 2026-10-16 — no action owed)

Each carry a full rationale on the PR: withdrawal-with-the-repo, the measured
branch/tag preservation SHAs, and — importantly — an HONEST statement that the
diff did NOT already land on main (verified: `cors_write_path_test.go` absent
from main; main differs from each head in ~197-198 files). "Superseded" here
means superseded-by-architecture, not already-merged. A follow-up comment on #32
corrects my own mislabeling of the work as "business logic" (AGENTS.md makes this
repo a zero-logic thin wrapper).

Status → DONE. Commented, not merged, branch preserved.

## F8 — net4sats-feed PR #1: CANNOT merge (repo write), escalated on the PR

PR #1 is OPEN/MERGEABLE/CLEAN but merging needs repo write, which the account
lacks. A blocking comment was posted to the PR (comment id 6380613573) carrying
the evidence, and the branch is tag-preserved (F5).

Re-measured the harm it fixes, on `main` today, rather than trusting the PR body:

    files/etc/uci-defaults/95-net4sats-feed:6  FEED_URL="https://blossom.net4sats.com/feed"
    install-net4sats.sh:10                     NET4SATS_FEED=".../feed/net4sats"
    README.md:16,20                            same host, two snippets
    getent hosts blossom.net4sats.com          -> NXDOMAIN (curl: could not resolve)

So the uci-defaults script appends an **unresolvable** origin to
`/etc/apk/repositories` / `/etc/opkg/customfeeds.conf` on first boot and breaks
every later `apk update` / `opkg update`. A live hardware bug, not cosmetics.

**ORDERING TRAP recorded for the operator:** merging #1 as written means main's
README never becomes the archive pointer (the pointer text is in the PR body),
and the repo would then be archived with a README still telling operators to add
a dead feed. Intended order: merge the fix and put the pointer text into main's
README → then archive.


