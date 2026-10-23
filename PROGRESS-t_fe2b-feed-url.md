# PROGRESS — FE2b (step 1): GL-MT6000 feed base URL + exact src/gz line

Finding -> Status -> Files touched. Worktree root. 2026-10-23. Read-only task.

1. Finding: Two different "feed" surfaces are published and the brief conflates them.
   `releases.tollgate.me/alpha/<tag>/<arch>/` is a REAL opkg repo (Packages.gz + Packages
   + matching .ipk). `github.com/FreedomTechFeed/packages/releases/download/<tag>/` is the
   CURRENT release and has flat assets only, zero index.
   Status: RESOLVED — the exact line exists only for the stale tag.
   Files touched: none (probes in /tmp).

2. Finding: The exact line, verbatim, for the only published feed:
   `src/gz tollgate https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/packages`
   Tag `v0.5.0-alpha3`, channel `alpha/`, last-modified 2026-07-01.
   Status: CAPTURED. Note the real layout is `<arch>/Packages.gz` — there is NO `/packages`
   subdirectory, so the brief's `<BASE>/<arch>/packages` template needed correcting.
   Files touched: none.

3. Finding: Index is genuinely fetchable and genuine (not SPA catch-all): 200,
   `content-type: application/gzip`, gzip magic, and the asset it names downloads with a
   sha256 identical to the index stanza (c508e0b5...97fb3, 5321054 B).
   Status: VERIFIED. No Packages.sig (404) — feed is unsigned.
   Files touched: none.

4. Finding: CORRECTION to REPORT-fe2b-feed-mt6000.md (2026-10-21). Its "no feed index is
   published, in any form, at any candidate base" is too strong: it is true for the GitHub
   release prefix it probed, but a live index exists on releases.tollgate.me for alpha3.
   Its negative result still holds for the CURRENT release.
   Status: RESOLVED (supersede one claim, keep the rest).
   Files touched: none.

5. Finding: v0.6.0 (newest, 2026-10-21) publishes 26 flat assets and no index of any shape;
   the URL a src/gz line would fetch returns 404. Root cause in source: the installer
   defines "feed" as a flat-URL prefix (`arch.go:107,119,321,349`), and the index generator
   is still unmerged (`FreedomTechFeed/feed@ci/auto-feed-indices`, 1569f78, 2026-07-05).
   Status: CONFIRMED — so a src/gz line and the current release are mutually exclusive today.
   Files touched: none.

6. Finding: alpha3 index resolves `tollgate-wrt v0.5.0-alpha3` (opkg .ipk), which is a
   DOWNGRADE from the 0.6.0_rc6-r1 the prior card measured on the MT6000; repo has
   refuseDowngrade for exactly this.
   Status: RISK FLAGGED to manager. Not installed, not configured.
   Files touched: none.

7. Finding: 192.168.1.1 is no longer an MT6000 — it is a QEMU x86_64 guest (25.12.5),
   matching REPORT-t_23b8b386.md. Arch therefore not confirmed from device this run.
   Confirmed instead from upstream metadata: mediatek/filogic profiles.json ->
   arch_packages aarch64_cortex-a53, profile glinet_gl-mt6000 present.
   Status: DEVICE CAVEAT recorded; no arch substitution made.
   Files touched: none.

8. Finding: Root REPORT.md/PROGRESS.md belong to other cards; shared
   `.git/info/exclude` (outside this worktree) blocks /PROGRESS*.md and /REPORT*.md.
   Status: DEVIATION — wrote task-scoped REPORT-t_fe2b-feed-url.md +
   PROGRESS-t_fe2b-feed-url.md at the worktree root, added with `git add -f`.
   Files touched: those two files.
