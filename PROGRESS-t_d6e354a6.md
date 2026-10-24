# PROGRESS — t_d6e354a6 · Add .golangci.yml and fix pre-existing lint failures

Workspace: /home/c08r4d0r/repos/net4sats-mvp-v2
Branch: worker-heavy/t_gate8-19b1b1c6
Date: 2026-10-24

## Acceptance state (summary)

| # | Requirement | Status |
|---|---|---|
| 1 | `.golangci.yml` exists | PASS — pre-existed on disk (untracked before my commit), v1 schema matching the installed v1.61.0 |
| 2 | `golangci-lint run` exits 0 | PASS — exit 0, verified on a cold cache |
| 3 | `go build ./...` passes | PASS — exit 0 |
| 4 | `go vet ./...` passes | PASS — exit 0 |
| 5 | `go test ./...` passes | PASS on the committed tree (see finding 2) |

## Finding 1 — the lint config already existed, so step 1 was a no-op

`.golangci.yml` was already present and correct. It is v1 schema:

    linters:
      enable: [errcheck, gosimple, govet, ineffassign, staticcheck, unused]

The card asked for a v2-schema config (`version: "2"`, `default: standard`) but
explicitly said to adapt to the installed version. Installed tool:

    $ golangci-lint --version
    golangci-lint has version 1.61.0 built with go1.23.1

v1.61.0 parses the v1 schema only. I therefore **kept the v1 form** — this is the
version CI itself pins (`.github/workflows/ci.yml`: `golangci/golangci-lint-action@v6`
with `version: v1.61.0`), so config and CI agree. A v2 file would have made CI red,
which is the exact opposite of this card's goal.

The six enabled linters are exactly the v1 default set written out explicitly, so
the intended coverage cannot drift with upstream defaults. It also sets
`issues.max-issues-per-linter: 0` and `max-same-issues: 0`, so no finding is
silently capped — a green run means genuinely clean.

Status: NO CHANGE NEEDED. I did not rewrite the file to v2.

## Finding 2 — `go test ./...` was already red before I touched anything

`go build ./...` and `golangci-lint run` are green, but `go test ./...` failed on
`TestNoInternalProcessLeaksOnShippedSurfaces` (debrand_internal_test.go). I proved
this was **pre-existing** by extracting a pristine `git archive HEAD` tree into
/tmp/pristine-head and running the guard there — it failed identically at HEAD,
with no untracked files present.

Root cause: commit `9bc9898` added a tracked mirror at
`docs/wz1.5-channel-decision.md`, a path the guard scans (all `docs/*.md` are
shipped surfaces). That document is an internal work note; it carried card ids,
the fleet hostname, the operator's decision-gate label, and a raw `/home/...`
listing. The guard landed earlier in `4b9dfb0`; the mirror was committed after it
and the guard was never re-run against it.

Status: FIXED in source (the doc), not waived. I did **not** add an `//nolint`, and
I did not widen the guard's exclusion list — the repo's own convention
(debrand_internal_test.go header, debrand_filenames_test.go history) is to fix the
leak, not exempt it.

De-identification applied, with every technical claim preserved:

| line | before | after |
|---|---|---|
| 3-4 | ``**Card:** `t_1df48589` `` / worker card + profile | "internal task id, withheld from this shipped surface" |
| 7 | "by card `t_1df48589`" | "in a second, independent pass on the same internal task" |
| 71-73 | `ls -la /home/<user>/...`, `c08r4d0r c08r4d0r` ownership | `~/.hermes/...`, `<user> <user>` |
| 80 | `find /home/c08r4d0r` | `find ~` |
| 87 | "the FELIX Q3 / dev-split gate" | "the operator's Q3 / dev-split gate" |
| 464 | "A separate card, `t_deabf3d6`," | "A separate internal task has already" |

No numbers, URLs, 404 results, file sizes or decision logic were altered; the
remaining 8 guard hits are now 0, confirmed by an independent regex sweep as well
as by the guard itself.

## Finding 3 — an untracked scratch file also trips the guard, but is not mine

`docs/WZ1.8-gate8-gate9-review-verdict-t_cc76c746-2026-10-24.md` (untracked, written
by the `t_cc76c746` review card) is a task-history artifact sitting in `docs/`.
Because the guard walks the *directory*, not the git index, it also fails on that
file locally. It is **not committed**, so CI (which checks out HEAD) never sees it.

I deliberately did NOT delete, edit or commit another card's file. Note for whoever
finishes this lane: if that file is ever `git add`-ed into `docs/`, the guard goes
red for everyone. It belongs at the worktree root (like the other REPORT-*/PROGRESS-*
artifacts) or outside the repo.

I verified `go test ./...` on a tree built from `git archive HEAD` + my one doc fix,
which is exactly what CI will check out.

## Verification commands and exit codes

    $ golangci-lint --version
    golangci-lint has version 1.61.0 built with go1.23.1 from a1d6c560 on 2024-09-09

    $ golangci-lint cache clean && golangci-lint run ./...     # cold cache
    (no output)  -> exit 0

    $ go build ./...                                          -> exit 0
    $ go vet ./...                                            -> exit 0
    $ go test ./... -count=1     (pristine HEAD + doc fix)     -> see REPORT.md
