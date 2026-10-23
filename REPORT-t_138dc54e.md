# REPORT — t_138dc54e: Confirm FE2a merged and release tagged (precondition gate)

**VERDICT: PRECONDITION NOT MET.**

Both halves of FE2b's precondition fail: the FE2a PR is **closed unmerged**, and
**no release is tagged or published** — the feed repo has zero tags and zero
releases. Per the brief, work stopped here: no device work, no repo source
modified, no FE2b step begun.

- Date: 2026-10-23
- Card: `t_138dc54e` — "Confirm FE2a merged and release tagged (precondition gate)"
- Board: `net4sats-mvp-v2`
- Workspace: `/home/c08r4d0r/repos/net4sats-mvp-v2` (branch `worker-heavy/t_gate8-19b1b1c6`)
- Mode: read-only verification; only this task's two records written

---

## 1. Step 1 — the definition of done and the precondition, verbatim

**The plan file named in the brief does not contain FE2 content.** First-hand:

```
$ grep -rn "FE2" --exclude-dir=.git --exclude-dir=node_modules .   # in-repo hits are prior reports only
$ grep -rln "FE2" /home/c08r4d0r/.hermes/plans/                    # 0 files
```

`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` (337 lines, 17,507 B) is
the **FELIX Q3 dev-split gate**. It never mentions FE2, FE2a, FE2b, or a feed.
This independently reproduces the note already on disk at
`REPORT-t_1b1cf42f.md:278-280`:

> (The plan file at `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` does not mention FE2a,
> FE2b, or a feed at all; it is an unrelated Q3 dev-split gate. Do not look there for FE2b's
> definition of done.)

The wording exists in the **kanban cards**, which is where the definition of done
and the precondition actually live. Quoted verbatim from the board DB
(`.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`, read via the recovered JSON
snapshots in `.hermes/state/kanban_recovery/…/boards/net4sats-mvp-v2/tasks/`):

### FE2b — `t_53879c67` (status: `blocked`)

> **Title:** `FE2b: Feed consumption install on GL-MT6000`
>
> **Body:** `Plan ref: ~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md | Objective: after FE2a merged + release published: on GL-MT6000 add feeds.conf src/gz line -> apk update -> install tollgate-wrt from FEED (not local scp) -> service up, portal serves. Proves feed consumable on real hardware. Serialize with fd4/wz18. QGATES v3.2: TDD (failing test first where code) | full suite green | docs same-commit | atomic conventional commits | PUSHED (done=pushed) | cold-review by reviewer-kimi before complete | manager Gate-8 review/merge within 24h.`
>
> **Validation:** `apk output + curl :2121 health + portal 200 in result.`
>
> **Precondition:** `manager merged FE2a PR + release tagged.`

The **definition of done** is therefore: on a GL-MT6000 (Flint 2), add a
`feeds.conf` `src/gz` line, `apk update`, install `tollgate-wrt` **from the feed
(not local scp)**, service up, portal serves — validated by `apk output + curl
:2121 health + portal 200` captured **in the card result**.

The **precondition** is the single sentence: **`manager merged FE2a PR + release
tagged.`** Note it is a conjunction — both terms required.

### FE2a — `t_df4304d5` (status: `todo`)

> **Title:** `FE2a: Feed PR#1 rebase + artifact-to-release bridge + arch tuples`
>
> **Body:** `Plan ref: ~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md | Objective: rebase ci/auto-feed-indices onto main (take PR side wholesale for build-feed.yml, redo docs), ADD missing bridge job (workflow artifacts -> tagged release assets), OpenWrt arch tuple naming (aarch64_cortex-a53 etc.), gh-pages publishing (Q3 default), tag first release, SDK-container consumption test (feeds.conf src/gz -> apk update -> install -> deps resolve). Push branch, open/update PR -> status review. Manager merges. QGATES v3.2: … | PUSHED (done=pushed) | cold-review by reviewer-kimi before complete | manager Gate-8 review/merge within 24h.`
>
> **Validation:** `consumption test log in result; PR mergeable.`
>
> **Result:** `None`

FE2a's own objective includes "tag first release" — i.e. the tag half of the
precondition is FE2a's deliverable, and FE2a has not started.

### Preflight — `t_18d15abe` (status: `triage`)

> `Goal: establish that FE2b's precondition holds and that we have a *failing* check that the feed install will satisfy (TDD / QGATES v3.2).`
>
> `1. Read the plan at ~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md and confirm the FE2b definition of done. Confirm with the manager (or the merged PR record) that the FE2a PR is merged AND the release is tagged/published. If either is missing, STOP and report back to the manager — do not proceed.`

This card (`t_138dc54e`) is the precondition gate the brief describes, and the
"if either is missing, STOP" instruction is the brief's step 4.

---

## 2. Step 2 — is FE2a merged? **NO.**

FE2a's PR is **`FreedomTechFeed/feed` PR #1**, branch `ci/auto-feed-indices` ->
`main`, exactly as the FE2a card title names it ("Feed PR#1 rebase").

```
$ gh pr view 1 -R FreedomTechFeed/feed --json number,state,mergedAt,mergeCommit,title,url,closedAt
{"closedAt":"2026-09-08T11:16:51Z","mergeCommit":null,"mergedAt":null,"number":1,
 "state":"CLOSED",
 "title":"feat(ci): auto-generate opkg+apk indices on every release (all arches)",
 "url":"https://github.com/FreedomTechFeed/feed/pull/1"}
```

Cross-checked against GraphQL, which reports the `merged` boolean directly:

```
$ gh api graphql -f query='{repository(owner:"FreedomTechFeed",name:"feed"){pullRequest(number:1){number state merged mergedAt mergeCommit{oid} closedAt}}}'
{"data":{"repository":{"pullRequest":{"number":1,"state":"CLOSED","merged":false,
 "mergedAt":null,"mergeCommit":null,"closedAt":"2026-09-08T11:16:51Z"}}}}
```

| Field | Value |
|---|---|
| PR number | **#1** (`FreedomTechFeed/feed`) |
| State | **CLOSED** — not merged |
| `merged` | **false** |
| Merge commit SHA | **none** (`null`) |
| Merge timestamp | **none** (`mergedAt: null`) |
| Closed at | 2026-09-08T11:16:51Z |
| Head branch | `ci/auto-feed-indices` @ `1569f78d81e4e7017e7304b8add93bafaef52496` (still exists — abandoned, not merged) |
| Base | `main` @ `24fb821` |

Full PR inventory confirms nothing else merged the feed-index work. Of 11 PRs in
the repo, only #3, #4 and #7 are MERGED, and none of those is the index
generator (#3 = archive-from-file fix, #4 = SDK build targets, #7 = cache
namespacing). PR #1 remains CLOSED and the four newest PRs (#8–#11) are OPEN.

**Feed vector on the FE2a definition: not merged, no merge commit exists to
record.**

---

## 3. Step 3 — is the release tagged AND published, with assets? **NO.**

```
$ gh release list -R FreedomTechFeed/feed
(no output — zero releases)

$ gh api repos/FreedomTechFeed/feed/releases --jq 'length'
0

$ gh api repos/FreedomTechFeed/feed/releases/latest
{"message":"Not Found","status":"404"}

$ gh api repos/FreedomTechFeed/feed/git/refs/tags
{"message":"Not Found","status":"404"}     # no tags namespace at all
```

| Field | Value |
|---|---|
| Tag name | **none exist** |
| Release URL | **none exist** |
| Artifacts attached | **none** — no release to attach to |
| Release count | **0** |
| Tag count | **0** |

### Trap worth recording: `releases/latest` returns HTTP 200

```
$ curl -s -o /dev/null -w 'final=%{url_effective} code=%{http_code}\n' -L \
    https://github.com/FreedomTechFeed/feed/releases/latest
final=https://github.com/FreedomTechFeed/feed/releases code=200
```

The 200 is a **redirect to the empty `/releases` page**, not a release. The API
endpoint for the same thing is a hard 404. A future lane probing with
`curl -o /dev/null -w '%{http_code}'` would wrongly read this as "a release
exists". Use `gh api .../releases/latest` or `gh release list`, not the HTML URL.

### No index is reachable anywhere either

Even if a tag existed, the artifact the FE2b `src/gz` line needs does not resolve:

```
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/packages.adb -> 404
https://releases.tollgate.me/v0.6.0/aarch64_cortex-a53/packages.adb       -> 404
```

The only release host with wizard artifacts is `felixfelix-bot/net4sats-wizard-go`,
whose releases (`v0.7.0-alpha18` … `v0.7.0-alpha24`, latest
`v0.7.0-alpha24` 2026-10-01T11:37:34Z) are the **wizard's own binaries, not a
package feed** — they carry no `packages.adb` / `Packages.gz` index, so they
cannot satisfy the FE2b precondition's "release tagged" term in the sense the
card means (a feed release).

---

## 4. Step 4 — STOP condition hit

Both precondition terms are missing:

1. **FE2a PR not merged.** PR #1 CLOSED 2026-09-08, `merged:false`, no merge
   commit, no merge timestamp.
2. **Release not tagged or published.** 0 tags, 0 releases, 0 assets in
   `FreedomTechFeed/feed`.

Per the brief and per `t_18d15abe` step 1 ("If either is missing, STOP and report
back to the manager — do not proceed"), **no device work was performed and no
repository source file was modified.** Specifically I did **not**: touch the
GL-MT6000, write a `feeds.conf` line, run `apk update`/`apk add`, create
`scripts/fe2b_feed_check.sh`, or begin any FE2b/FE2a implementation.

### What is missing, for the manager

| # | Missing | Owner card | Evidence |
|---|---|---|---|
| 1 | FE2a feed-index PR merged to `FreedomTechFeed/feed@main` | `t_df4304d5` (FE2a, `todo`) | PR #1 CLOSED unmerged 2026-09-08 |
| 2 | A tagged + published feed release with index artifacts attached | `t_df4304d5` (FE2a objective: "tag first release") | 0 releases, 0 tags |

Note that #2 cannot be satisfied before #1 in practice: the artifact-to-release
bridge and index generation are precisely what PR #1 + FE2a's "missing bridge
job" would create. FE2b stays blocked until FE2a is done and a real feed release
exists.

---

## 5. Honesty ledger

- Every command output quoted above was executed in this session, from
  `/tmp`, read-only. No output is reconstructed or inferred.
- I did **not** obtain the FE2b definition of done from `schedule.md`, because it
  is not there. I said so rather than quoting the Q3 gate text as if it were the
  FE2b definition. The brief's step 1 premise is the thing that is wrong, and it
  is flagged as such.
- I did **not** claim a merge commit or merge timestamp: none exists, so the
  fields are recorded as `null` rather than left blank or filled with the PR's
  close time (which is a different event — 2026-09-08T11:16:51Z is when it was
  *closed*, not merged).
- I verified tag/release absence two ways (API count and `releases/latest` 404)
  and explicitly probed the misleading HTML-200 case, so the negative is not a
  single-endpoint artifact.
- Read-only respected: `git status` at session start showed only the two
  pre-existing untracked dirs (`.hermes/`, `verify/wz16-e2e-gate/`); this task
  added only its own two record files.
- Prior lanes' reports were read before working (`REPORT-t_1b1cf42f.md`,
  `PROGRESS-t_d41ef5b5.md`); this report reproduces their negative findings
  independently rather than citing them as evidence.
