# PROGRESS — t_e75c5cd6 (pinned mnemonic library default strength & word-count semantics)

Session 2026-10-23 (fleet offload, worker-heavy). Finding -> Status -> Files touched.
READ-ONLY against the worktree: no repo code/doc/config/AGENTS.md edited. All probes ran
in /tmp. Full report: `REPORT-t_e75c5cd6.md` at this worktree root.

1. Finding: The brief's INPUT ("the parent task's report") DOES NOT EXIST.
   Board `net4sats-mvp-v2` card `t_48cbf9fe` is `status=triage`, `result=NULL`,
   `session_id=NULL` — it never ran. Its own INPUT pointer
   (`~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` WZ1.7) was already proven
   missing by earlier cards (`REPORT.md:29-30`).
   Status: VERIFIED (read-only sqlite3 on the board DB; prior report cross-check).
   Consequence: library, version and call sites re-derived first-hand in this session.
   Files touched: none.

2. Finding: `net4sats-mvp-v2` contains ZERO mnemonic code. `git grep -n -i -E
   'bip39|mnemonic|NewEntropy|wordlist'` over tracked non-report sources →
   NO_HITS_IN_TRACKED_SOURCES; the shipped `net4sats-wizard` binary has 0 occurrences of
   `tyler-smith/go-bip39`; the bundled JS maps' only `bip39` strings are `bip39seed`
   option keys in the vendored cashu-ts wallet. The mnemonic work is entirely in the
   backend `tollgate-module-basic-go`, consistent with `AGENTS.md`'s thin-wrapper rule.
   Status: VERIFIED (git grep, `strings`, per-file `grep -o` counts).
   Files touched: none.

3. Finding: TWO pinned libraries generate mnemonics, both at 12 words, and NEITHER has a
   "default strength" because neither API accepts an omitted strength.
   (a) `go-bip39 v1.1.0` — `NewEntropy(bitSize int)` is mandatory; validated into the
   closed set {128,160,192,224,256} (`bip39.go:99`, `:315-320`). Pinned at
   `tollgate-module-basic-go/src/go.mod:115`; go.sum hash matches the module-cache
   ziphash (source read == source built against).
   (b) `cdk-ffi 0.17.3` — the wrapper takes no argument and hardcodes
   `Mnemonic::generate(12)` (`wallet.rs:876`); underlying rust `bip39` 2.2.2 has no
   `Default` impl (`lib.rs:314`, `:682-684`).
   Status: VERIFIED (source read + executed probes).
   Files touched: none.

4. Finding: Executed confirmation, not inference.
   Go probe at `/tmp/bip39probe` (GOPROXY=off, pinned module cache, go1.24.2 at
   /usr/local/go/bin/go): 128→12, 160→15, 192→18, 224→21, 256→24 words; 0/64/96/127/288
   bits all rejected. Parse side: `IsMnemonicValid` true for 12/15/18/21/24, false for
   "", "abandon", "not a valid mnemonic at all"; a 24-word phrase is ACCEPTED.
   CDK probe: `uniffi_cdk_ffi_fn_func_generate_mnemonic` called via ctypes against the
   shipped `libcdk_ffi.so` (cdk-go@v0.17.3) → 6/6 runs returned 12 words / 128 bits.
   Status: VERIFIED (both probes ran; outputs captured verbatim in the report).
   Files touched: /tmp/bip39probe/{go.mod,go.sum,main.go}, /tmp/cdkprobe.py (outside repo).

5. Finding: Effective strength per call site. Generation: `identity.go:137` literal 128 →
   12 words; `cdk_wallet.go:55` omits the argument but the wrapper hardcodes 12 → 128 bits
   → 12 words. Parse/validate (`identity.go:208`, `:151`, `main.go:2195`): no strength is
   chosen, and `IsMnemonicValid` accepts 12/15/18/21/24 — so the recovery endpoint does
   NOT pin 12 words; only the generators do. Enforced by tests at
   `identity_test.go:57,88,238` and `identity_handler_test.go:78`.
   Status: VERIFIED (source + probe).
   Files touched: none.

6. Finding: `AGENTS.md:43` ("BIP39 mnemonic (24 words)") is WRONG for the code at
   backend HEAD `c77ef204`. The only 256-bit path is `nip06.GenerateSeedWords`
   (`nip06.go:10-11`), which has NO callers. No dynamic strength exists either:
   grep for `--words/--strength/--entropy/word_count/WordCount` and for
   `flag.(Int|String)` over backend `src/` → zero hits.
   Status: VERIFIED. NOT fixed — read-only constraint; the report states the correction
   only.
   Files touched: none.

7. Finding: Environment limits hit, disclosed rather than papered over.
   `go` is not on PATH (used /usr/local/go/bin/go). `github.com` DNS failed mid-session
   (`Could not resolve host`), so (i) the cdk-ffi git commit could not be byte-verified
   against tag 34882e8a and (ii) the final `git ls-remote` / push verification could not
   be repeated. crates.io DID resolve, so the cdk-ffi 0.17.3 crate was fetched from
   static.crates.io and its behaviour confirmed against the shipped .so.
   Status: VERIFIED (errors observed).
   Files touched: none.

8. Deliverables written: `REPORT-t_e75c5cd6.md` (full report) and this file.
   Committed with `git add -f` (both are covered by `.git/info/exclude`, which lists
   `/PROGRESS-t_*.md` and `/REPORT-t_*.md` as local-only). Commit `02598dd`, 2 files,
   +328 lines. The repo pre-commit hook ran and reported
   `cred-gate: clean (30 literals checked, mode=staged)`.
   Pushed and verified: `15f5d29..02598dd -> origin/worker-heavy/t_gate8-19b1b1c6`,
   and `git ls-remote` returns `02598dd5f96d94685047d0fd35c2c00c92b61f8c`, matching local
   HEAD exactly. (The first push attempt failed on a transient `Could not resolve host:
   github.com`; DNS recovered and the retry succeeded.)
