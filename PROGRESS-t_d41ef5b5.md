# PROGRESS — t_d41ef5b5: document feed-install procedure (atomic conventional commit)

Session 2026-10-23 (worker-heavy, workspace /home/c08r4d0r/repos/net4sats-mvp-v2,
branch `worker-heavy/t_gate8-19b1b1c6`). Finding -> Status -> Files touched.

1. Finding: The card's premise ("read the parent task result — it contains the
   exact feed base URL, the src/gz line, and the raw `apk update` / `apk add`
   output") is FALSE. Parent is t_d41ef5b5's own FE2b chain: `t_53879c67`
   (FE2b install) has result=NULL and only blocked/scheduled runs; the
   design-named parent `t_1b1cf42f` resolves a URL but reports NO feed index
   exists for the current release. A sweep of the board DB for
   `src/gz|apk update|apk add|tollgate-wrt` found no captured install
   transcript for this card to document. There is no successful
   `apk update` / `apk add tollgate-wrt` against a real GL-MT6000 anywhere.
   Status: VERIFIED — read parent reports, board DB (net4sats-mvp-v2), and the
   cold-review verdict attachment before writing anything.
   Files touched: none.

2. Finding: The FE2a precondition is still UNMET, measured live today.
   `gh pr list -R FreedomTechFeed/feed --state all`: index-generation PR #1
   (`ci/auto-feed-indices`) = CLOSED, mergedAt=null; 0 releases; 0 tags.
   `.../v0.6.0/aarch64_cortex-a53/packages.adb` -> 404. `releases.tollgate.me/
   alpha/v0.6.0/...` -> 404. So no `src/gz` line can install the current
   release.
   Status: VERIFIED (live probes, 2026-10-23).
   Files touched: none.

3. Finding: The only published index is `v0.5.0-alpha3` on
   releases.tollgate.me, and it is BOTH stale and UNSIGNED. Measured on a real
   apk-tools 3.0.5 userland: `apk update` -> `UNTRUSTED signature`,
   `1 unavailable`; `apk add tollgate-wrt` -> `no such package`. Measured on
   real 24.10.2 opkg: index downloads, then `Packages.sig` 404 -> list REMOVED,
   `opkg list` shows nothing; `check_signature 0` does not rescue it.
   Status: VERIFIED by execution, not citation. This is why the doc documents a
   precondition + real failure modes rather than a copy-paste "it works" run.
   Files touched: none.

4. Finding: The card's briefed `src/gz` spelling is WRONG for opkg. Measured on
   real 24.10.2: base including a trailing `/packages` -> `wget returned 8`
   (opkg appends `/Packages.gz` itself); base = arch directory -> index fetched.
   The parent report's verbatim line carries the bad `/packages` segment.
   Status: CORRECTED in the doc with both spellings and their outputs.
   Files touched: `docs/feed-install.md`.

5. Finding: Real evidence was capturable here without a GL-MT6000, so the doc
   cites today's own captures rather than the 2026-10-11 rehearsal transcripts.
   Real OpenWrt 25.12.5 (apk-tools 3.0.5) and 24.10.2 (opkg 38eccbb) userlands
   imported into containers; the real published v0.6.0 x86_64 package
   digest-checked against the release SHA256SUMS, installed, `o:` origin read
   back as `/feed/net/tollgate-wrt` (proving a local-file, NOT feed, install),
   and `:2121` probed for a `kind:10021` advert.
   Status: VERIFIED. Portal `:2050` NOT measurable (nodogsplash needs
   procd/ubus + br-lan) and is documented as not-measured, not as a failure.
   Files touched: `docs/evidence/2026-10-23-fe2b-measured-captures.md`.

6. Status: DONE. Docs written, README index wired, repo guards green
   (`go test -count=1 ./...` PASS, 54.97s incl. the shipped-surface leak guard
   and the filename guard). Committed as ONE atomic conventional commit and
   pushed.
   Files touched: `docs/feed-install.md`, `docs/evidence/
   2026-10-23-fe2b-measured-captures.md`, `README.md` (docs index line),
   `PROGRESS-t_d41ef5b5.md`, `REPORT-t_d41ef5b5.md`.

## Deviation

Per the card's acceptance criteria ("docs are accurate against the captured real
outputs"), accuracy was satisfied by producing real outputs today rather than by
citing a parent result that contains none. The card's precondition (FE2a merged,
release tagged) is NOT met, so the doc documents the procedure plus the measured
blockers instead of claiming a successful feed install. Reported, not hidden.
