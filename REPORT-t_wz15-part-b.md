
# REPORT — WZ1.5 part B: `tollgateOSURL` / `configwizURL` dead-code adjudication (t_wz15-part-b)

## Deliverable

`/home/c08r4d0r/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-url-inventory.part-b.md`
(~21 KB). Standalone: commands run, per-constant reference-count table, the
hardcoded-duplicates table, and a headed Verdict section with one bullet per
constant. **Not merged with part A** — §6 of the doc records a cross-check
discrepancy instead.

Task type: read-only adjudication + a deletion experiment on a throwaway copy.
No file in this repo was modified; `git status` is byte-identical to session
start except the two untracked dirs that were already there.

## Acceptance criteria — status

| criterion | status | evidence |
|---|---|---|
| every constant has an unambiguous verdict | **MET** | §5: `tollgateOSURL` = DEAD (already deleted, no action); `configwizURL` = DEAD (delete, 4 edit sites). No "unclear" anywhere. |
| verdict backed by file:line | **MET** | both verdicts carry declaration sites, every reference site, the pickaxe commits that added/removed each consumer, and HEAD token counts (2 and 3 lines). |
| reference counts split prod / test / comment | **MET** | §2 table: prod refs 0 and 0; test refs 1 (`pins_test.go:49`, a *string* in a guard) and 2 (`pins_test.go:31`, `:77`); comment refs 1 and 1. |
| hardcoded duplicates table | **MET** | §4a exact-value table + §4b same-artifact/same-host table. Key negative result: **no shipped duplicate of either exact URL exists**, so deletion is single-source clean. |
| commands run recorded | **MET** | §1, verbatim, including the AST check, the branch sweep, the curl liveness probes and the throwaway-copy build/test run. |

## Headline findings

1. **`tollgateOSURL` is gone and should stay gone — but the guard must stay too.**
   0 identifier occurrences in any `.go` file at HEAD; its literal value appears
   0 times anywhere on disk. History is unambiguous: added `fa306d2` (2026-07-16)
   as `releases.tollgate.me/os/57e0f246…?channel=stable`, declaration-only its
   entire life (`uses=1` in all 24 intervening commits that touch `deploy.go`),
   deleted `d25c018` (2026-08-17). The URL is **404 today** (measured). The two
   surviving lines are `pins_test.go:45-46` (prose) and `:49`
   (`forbiddenDeployIdentifiers`). Keep them: the guard compares with
   `strings.Contains` over identifiers *and* string literals, so it is what would
   stop the 404 URL from coming back — and it is live-tested by
   `TestDeployGoPinRegistry` (PASS observed). Six stale remote branches still
   *declare* it (declaration-only there too; four are 82-87 commits behind HEAD,
   two are ancestors of HEAD), so the guard has real future work to do.

2. **`configwizURL` is dead, and the live-URL test is why nobody noticed.**
   3 token lines: `deploy.go:50` (declaration), `pins_test.go:31` (name string in
   `expectedPinnedURLConsts`), `pins_test.go:77`
   (`"configwizURL": configwizURL` in `liveCheckPins`). Zero production
   references. `git log --all -S'httpGetFile(configwizURL'` → **empty**: it has
   never been passed to the download helper in any commit. Go AST confirms the
   only `*ast.Ident` named `configwizURL` are the declaration and
   `pins_test.go:77`. Meanwhile `TestPinnedURLsAreLive/configwizURL` **passes**
   (real HTTP 302→206, 966,774 B) — a green test that measures the asset's
   existence, not the binary's behaviour. That is the trap the verdict had to
   survive: passing ≠ used.

3. **It was killed on purpose and early.** `d70b8be` (2026-07-16) introduced the
   constant *with* a consumer (`wget -qO /tmp/configwiz.tar.gz '…' + tar xzf`).
   `adfaffb` (same day, 2026-07-16T21:39:06+05:30, "embed admin panel in wizard
   binary … no GitHub download needed") removed that consumer. From that commit
   to HEAD the constant has been declaration-only; `pins_test.go` only started
   referencing it at `d25c018` (2026-08-17). So it has been dead for the whole
   life of the current embedded deploy path (`embed.go:11-15` → `portalFS`/
   `adminFS` → `sshDeployFS` at `deploy.go:678`, `:694`, `:762`).

4. **Deleting it is proven safe, not assumed.** On a throwaway copy
   (`/tmp/wz15-del-exp`, repo untouched): `go build ./...` OK, `go vet ./...` OK,
   `gofmt -l deploy.go pins_test.go` clean, `golangci-lint run --timeout 5m`
   exit 0 (so the `unused` linter does not depend on it either), and
   `go test ./...` green in **64.5 s vs 63.2 s** on the unmodified HEAD. The
   delete must be 3 code edits in one commit (`deploy.go:48-50`,
   `pins_test.go:31`, `pins_test.go:77`) because `TestDeployGoPinRegistry`
   asserts declared-set == registry == live-check map; a partial delete goes red.

5. **The tarball's content is the thing worth keeping, and it is already in the
   tree.** The release asset's 21 payload files are byte-identical to the tracked
   `portal/` tree (`diff -r --brief` clean; per-file sha256 equal for all 21;
   tarball sha256 `87e5b9c8…`, 966,774 B). The repo's `portal/` does carry 4
   extra stale hashed assets (`splash-BYAdP4yt.js`, `splash-C0y8YH4o.js.map`,
   `splash-DkfE4jIG.js`, `splash-DkfE4jIG.js.map`) that the current HTML no
   longer references — noted, out of scope here. Recommendation is therefore
   "delete the constant, record the version+digest provenance in
   `docs/router-setup.md`", not "repin it somewhere else".

6. **Repinning would be actively wrong twice over.** `docs/router-setup.md:193`
   is the only shipped doc row for this artifact and it is stale *and* in the
   wrong namespace: it names `net4sats-configwiz-1.0.1.tar.gz` from fork `v1.0.1`
   while the tree actually vendors `v1.0.10-alpha`. GitHub API: the fork
   (`felixfelix-bot/configurationwizzard`, `fork: true`) is **ahead** of the
   official `net4sats/configurationwizzard` (`fork: false`, newest release
   `v1.0.0`). Neither "point at the official org" nor "keep the v1.0.1 row"
   describes reality — hence the third edit of the recommended commit.

7. **Part A disagrees with this and part A is wrong on this point.** Part A's
   classification key (its lines 94-98) declares class (b) empty and calls
   `configwizURL` "referenced from `deploy.go` and asserted by
   `pins_test.go`/`artifact_pins_test.go`". It is not in `artifactPins`
   (`deploy.go:67-70`) and never reaches `httpGetFile`. Part A's 163/163 line
   accounting is unaffected; its class-(b)=0 claim is a line-count artifact, not
   a reference count. Recorded in part B §6 for the synthesis task.

## What I did NOT do (scope honesty)

- Did not modify `deploy.go`, `pins_test.go`, `docs/router-setup.md` or any other
  repo source. The deletion is *recommended and proven* on a copy, not applied —
  this card's deliverable is the adjudication doc.
- Did not run the physical-router / Playwright E2E gate (not applicable: no
  shipped behaviour changes; the delete is statically proven).
- Did not touch CI, git hooks, or kanban state.

## Verification of the deliverable itself

- File exists, **21,461 bytes**, read back from disk after the §6 addition.
- Both constants appear in the doc's Verdict section with a single unhedged
  verdict token each (`DEAD (already deleted)` / `DEAD (delete)`); grepped the
  doc for `unclear` → 0 hits.
- Every file:line cited in §2/§4/§5 was produced by a command in §1 during this
  session; none was carried over from another card's report unverified. Where an
  earlier card's claim was reused (`REPORT-t_gate8-19b1b1c6.md` §4/§5), it was
  re-derived independently and matches.

## Bottom line

Two dead constants, two different actions: `tollgateOSURL` needs **nothing**
(already deleted; keep the tripwire and never reuse the name or the 404 URL);
`configwizURL` needs a **5-line deletion in one commit plus one doc-row
correction**, with the version provenance relocated to `docs/router-setup.md`
as a see-also rather than a Go constant. No repin, no relocation to code.
