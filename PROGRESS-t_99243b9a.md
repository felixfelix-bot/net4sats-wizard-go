# PROGRESS — Stage hardened configwiz portal dist (t_99243b9a)

Card: t_99243b9a "Stage hardened configwiz portal dist in v2keyset-fix worktree"
Repo: net4sats-mvp-v2 / upstream configwizzard (fork felixfelix-bot)
Worker host: x280 (c08r4d0r)

## Finding 1 — the brief's worktree path does not exist on this host
Brief said work in `/home/c03rad0r/worktrees/v2keyset-fix`.
`c03rad0r` is a DIFFERENT fleet user (dq05). On this host only `c08r4d0r`
exists. Verified:
- `ls /home/` → `c08r4d0r  ghrunner` (no c03rad0r)
- `git worktree list` in net4sats-mvp-v2 → no v2keyset-fix entry
- dq05 unreachable: ssh 192.168.2.32 → "No route to host";
  *.fips aliases → "Temporary failure in name resolution" (mesh DNS down)

→ Could not work in the literal path. Recreated an equivalent worktree at
`/home/c08r4d0r/worktrees/v2keyset-fix` from the real hardened branch.

## Finding 2 — V2K-1H produced NO dist; the hardened SOURCE exists on the fork
- Card t_b167c169 (V2K-1H) status=blocked, workspace empty, 0 comments.
- Hardened branch found on fork: `fix/v2keyset-gettokenmetadata`
  @ b127d488fafeed0a4cea3123390fa55b553f75ee
- Branch commits: b127d48 mint normalizePath parity / 093c9b8 V3 dual-encoding
  tests / ad91f84 getTokenMetadata v2-keyset decode / db6f580 vitest harness
- `src/lib/cashu-validate.ts` on branch: 278 lines, sha256
  9f9dfd10bb5a2661c5dd51b7df3a4dc39db63729e8576880d51dc9132a912c1e
- No prebuilt dist of it existed anywhere on this host (searched all
  worktrees + kanban workspaces). Other dist dirs found were UNHARDENED
  (tollgate-admin-ui scaffold, 154-line cashu-validate = pre-V2K-1H).

## Status — build DONE, dist verified

Worktree: /home/c08r4d0r/worktrees/v2keyset-fix (detached HEAD b127d48)
No source file edited. dist/ is gitignored → git status stays clean.

Build blockers hit and resolved (all environment, zero source edits):
1. System node 18.19.1 → Vite 8 needs >=20.19. Fixed by downloading portable
   node v22.14.0 to /tmp (no system change).
2. npm ci under node18 skipped optional native dep. Fixed by installing the
   exact lock-pinned `@rolldown/binding-linux-x64-gnu@1.0.0` with --no-save.

Build command actually used (documented repo command):
    PATH=/tmp/node-v22.14.0-linux-x64/bin:$PATH npm run build
  → runs `node scripts/build-app.mjs` = admin + portal + balance, plus the
    splash.html→index.html copy that uhttpd requires.

## Files touched
- Created (not in repo): /home/c08r4d0r/worktrees/v2keyset-fix (git worktree)
- Created (gitignored): that worktree's dist/
- Created: PROGRESS-t_99243b9a.md, REPORT-t_99243b9a.md (wizard repo root)
- NO source modification in the worktree.
