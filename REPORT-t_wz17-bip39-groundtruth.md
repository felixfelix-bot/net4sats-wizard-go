# REPORT — WZ1.7: BIP39 12-vs-24 ground-truth table + AGENTS.md candidate statements

Card:        net4sats-mvp-v2 WZ1.7 (fleet offload, worker-tollgate / worker-heavy)
Repo:        /home/c08r4d0r/repos/net4sats-mvp-v2  (remote felixfelix-bot/net4sats-wizard-go)
Wizard HEAD: ee330514f41819bce28ac48bb4713b81d3fabae0 (branch worker-heavy/t_gate8-19b1b1c6)
Wizard AGENTS.md blob: f2ff1775a57b434ff966e7c8eaf6287e30be9fc5 (unchanged at HEAD)
Backend:     /home/c08r4d0r/repos/tollgate-module-basic-go
             working tree HEAD eb12d7f8e7c1c79e46264c1b705de4e74544975b (branch pr-1037)
             origin/main  95835413 (2026-10-14, "Merge docs: dashboard repo at cvm-services")
Mode:        READ-ONLY. No code, doc, config or AGENTS.md edited. No source committed.
             Only this report + its PROGRESS file were added.

Citation convention: wizard paths are relative to the wizard repo root; backend paths are
written `tollgate-module-basic-go/src/...` and resolve under
`/home/c08r4d0r/repos/tollgate-module-basic-go/`; module-cache paths are written
`go-bip39@v1.1.0/<file>:NNN` and resolve under
`/home/c08r4d0r/go/pkg/mod/github.com/tyler-smith/go-bip39@v1.1.0/`. A ref prefix
(`origin/main:`, `7a87fb4d^:`) means that git ref was read, not the working tree.

---

## Ground-truth table

| Claim | Verdict (TRUE/FALSE) | Evidence (path:line) | Verbatim snippet |
|---|---|---|---|
| **CORE — the mnemonic is 12 words / 128-bit entropy** | **TRUE** | `tollgate-module-basic-go/src/identity/identity.go:137` (generator body; function at `:136`), doc comment at `:133` | `entropy, err := bip39.NewEntropy(128)` — and at `:133`: `// GenerateMnemonic generates a 12-word BIP39 mnemonic (128 bits of entropy).` |
| The mnemonic library is `github.com/tyler-smith/go-bip39`, exact version **v1.1.0** | **TRUE** | `tollgate-module-basic-go/src/identity/go.mod:8` (direct require, identity module); `tollgate-module-basic-go/src/go.mod:115` (indirect require, root module); checksum `src/identity/go.sum:132` = `src/go.sum:413` | `github.com/tyler-smith/go-bip39 v1.1.0` — root module: `github.com/tyler-smith/go-bip39 v1.1.0 // indirect` |
| That library has **no default strength**: the size argument is mandatory and range-checked, so nothing can silently produce 24 words | **TRUE** | `go-bip39@v1.1.0/bip39.go:99` (signature), `:315-316` (range check) | `func NewEntropy(bitSize int) ([]byte, error) {` — `if (bitSize%32) != 0 \|\| bitSize < 128 \|\| bitSize > 256 {` — `return ErrEntropyLengthInvalid` |
| 128-bit entropy maps to exactly 12 words, and the mapping is computed rather than tabulated | **TRUE** | `go-bip39@v1.1.0/bip39.go:165-167`; executed probe against the pinned v1.1.0 (see §Probes) | `entropyBitLength := len(entropy) * 8` / `checksumBitLength := entropyBitLength / 32` / `sentenceLength := (entropyBitLength + checksumBitLength) / 11` |
| The generator's 12-word output is asserted by the backend's own tests, and those tests pass at this HEAD | **TRUE** | `tollgate-module-basic-go/src/identity/identity_test.go:57`, `:88`, `:238`; `tollgate-module-basic-go/src/identity_handler_test.go:78`; first-hand run quoted in §Probes | `assert.Len(t, words, 12, "expected 12 BIP39 words")` — `assert.Len(t, words, 12, "mnemonic must be exactly 12 words")` |
| The **parse/recovery path does NOT pin 12** — it accepts any valid BIP39 length (12/15/18/21/24), so the 12-word figure is a generation-time pin, not an input restriction | **TRUE** | `tollgate-module-basic-go/src/identity/identity.go:151` → `go-bip39@v1.1.0/bip39.go:275-277` → `:356`; accepted lengths enumerated at `go-bip39@v1.1.0/bip39.go:29-35` | `if !nip06.ValidateWords(mnemonic) {` — `12: big.NewInt(15),` … `24: big.NewInt(255),` — `if numOfWords%3 != 0 \|\| numOfWords < 12 \|\| numOfWords > 24 {` |
| **INCORRECT 24-WORD DOC #1 (the claim under adjudication): `AGENTS.md` says "BIP39 mnemonic (24 words)"** — false for this repo's code | **FALSE** (the assertion) | `AGENTS.md:43`; contradicted by `tollgate-module-basic-go/src/identity/identity.go:137` | `└─ BIP39 mnemonic (24 words) → seed phrase (PrivateKeyToMnemonic)` |
| **INCORRECT 24-WORD DOC #2 (same sentence, same line): it names `PrivateKeyToMnemonic` as the function.** No such symbol exists in the backend at its working HEAD or on `origin/main` | **FALSE** | absence: `grep -rn 'PrivateKeyToMnemonic' --include='*.go' src/` over the backend working tree → 0 hits; `git log origin/main -S 'PrivateKeyToMnemonic' -- '*.go'` → 0 commits. Implemented symbols are `GenerateMnemonic` (`identity.go:136`) and `MnemonicToPrivateKey` (`identity.go:150`) | `AGENTS.md:43` again: `└─ BIP39 mnemonic (24 words) → seed phrase (PrivateKeyToMnemonic)`; the real definition: `func MnemonicToPrivateKey(mnemonic string) (string, error) {` — `not found in repo` for `PrivateKeyToMnemonic` |
| **INCORRECT 24-WORD DOC #3 (parallel line): `AGENTS.md:29` lists `PrivateKeyToMnemonic` among backend identity functions** | **FALSE** | `AGENTS.md:29`; absence as above | `- Identity: `src/identity/` — DeriveIPv4, DeriveMAC, PrivateKeyToMnemonic, MnemonicToPrivateKey` |
| A **24-word code path exists** in the pinned dependency stack (`nip06.GenerateSeedWords` → `NewEntropy(256)`), but it is dead code — zero callers in the backend | **TRUE** (the path exists) / **FALSE** (that anything calls it) | `github.com/nbd-wtf/go-nostr@v0.51.12/nip06/nip06.go:10-11`; caller count: `grep -rn 'GenerateSeedWords' --include='*.go' src/` → 0 | `func GenerateSeedWords() (string, error) {` / `entropy, err := bip39.NewEntropy(256)` — callers: `not found in repo` |
| The **second** mnemonic generator in the backend (wallet side, `cdk-ffi`) also produces 12 words; strength is hardcoded one layer below the call | **TRUE** (behaviour, measured 6/6) / version literal **not found in repo** (see residual uncertainty) | `tollgate-module-basic-go/src/tollwallet/cdk_wallet.go:55`; binding pinned at `src/go.mod:68`; shipped library `cdk-go@v0.17.3/bindings/cdkffi/native/linux_amd64/libcdk_ffi.so` | `mnemonic, err := cdk_ffi.GenerateMnemonic()` — `github.com/cashubtc/cdk-go v0.17.3 // indirect` — measured: `[(12, 'ugly'), (12, 'prize'), (12, 'toast'), (12, 'license'), (12, 'mutual'), (12, 'ribbon')]` |
| **Vendored wordlist / hardcoded length constant in THIS repo (the wizard)**: none. The wizard neither vendors a BIP39 wordlist nor hardcodes a word-count constant | **not found in repo** | wizard `go.sum` → 0 hits for `bip39\|mnemonic\|hdwallet`; `git ls-files \| grep -i -E 'wordlist\|english.txt'` → 0; all 7 shipped bundles `grep -c -i bip39` → 0; `strings -a net4sats-wizard \| grep -c 'tyler-smith/go-bip39'` → `0` | `not found in repo` |
| **Wordlist + literal values where they actually live (library, not vendored here)**: the 2048-word English list, CRC32-verified at init | **TRUE** | `go-bip39@v1.1.0/wordlists/english.go:13`, `:15`, `:23`; count measured on the pinned file = **2048** words, first `abandon`, last `zoo` | `// c1dbd296` — `checksum := crc32.ChecksumIEEE([]byte(english))` — `var english = ` + backtick string beginning `abandon` |
| **Hardcoded length constant that exists in the backend and is easy to mistake for a mnemonic length**: `6` (word count of the derived root/WiFi passwords), not a BIP39 word count | **TRUE** (constant exists, literal `6`) | `tollgate-module-basic-go/src/identity/identity.go:161`, `:172` | `parts := make([]string, 6)` |
| The wizard repo contains **no** BIP39/mnemonic code of its own, so `AGENTS.md` is the only place in this repo where the 24-word figure lives | **TRUE** | `git grep -n -i -E 'bip39\|mnemonic' -- '*.go' '*.sh' '*.html' Makefile` → 0 hits; the only non-report matches repo-wide are `AGENTS.md:17`, `:29`, `:43`, `:53` | `AGENTS.md:17:- Does NOT implement: DeriveIPv4, DeriveMAC, BIP39 mnemonic, kind:0 publishing` |
| **History attribution (parent conflict, resolved)**: "the backend switched from 24 to 12 words in commit `7a87fb4d` (2026-07-08)" is true **only on `origin/feat/identity-v2-carried`**; `7a87fb4d` is **NOT** an ancestor of `origin/main`. On `origin/main`, `src/identity/identity.go` has exactly two commits and was **introduced** already 12-word by `d491a4c4` (2026-08-17, #331) | **FALSE as stated for origin/main** / **TRUE for the carried branch** | `git merge-base --is-ancestor 7a87fb4d origin/main` → NO (`7a87fb4d` appears only on `origin/feat/identity-v2-carried`); `git log --oneline origin/main -- src/identity/identity.go` → 2 commits (`c60ea071`, `d491a4c4`); 24-word predecessor exists only at `origin/feat/identity-v2-carried` (`406b3865`, 2026-07-06) | `identity: switch from 24-word raw-entropy to 12-word NIP-06` (`7a87fb4d`, carried branch) vs `feat(identity): unified NIP-06 + HKDF + RevealSeed (v2, rebased) (#331)` (`d491a4c4`, the commit that put this file on `origin/main`) |
| **The 24-word era really did exist** — so the doc's figure was not invented, merely frozen | **TRUE** | `origin/feat/identity-v2-carried` history, `7a87fb4d^:src/identity/identity.go:124-133` (and same text at `406b3865:src/identity/identity.go:121-130`) | `// PrivateKeyToMnemonic converts a hex private key into a 24-word BIP39` / `func PrivateKeyToMnemonic(hexPrivKey string) (string, error) {` / `mnemonic, err := bip39.NewMnemonic(b)` |
| **24-word assertion in the third-party admin UI** (`configurationwizzard`), quoted by a parent: `identity-panel.tsx:174` "Write down these 24 words" | **not found in repo** (unverifiable on this host) | no `configurationwizzard` checkout exists: `ls -d /home/c08r4d0r/repos/configurationwizzard` → No such file or directory; `find /home/c08r4d0r -maxdepth 4 -name 'identity-panel.tsx'` → no results; the wizard repo tracks 0 `.tsx` files | `not found in repo` — the claim exists only in `EVIDENCE-NOTE-t_93620520.md` §105 and `attachments/t_e7c505b4/AUDIT-evidence-note-t_e7c505b4.md:8`, both quoting a ref not present here |
| The `v0.5.0-alpha3` tag that `AGENTS.md:33` claims compatibility with contains **no** identity package at all, so the 12-word pin cannot be asserted about it | **TRUE** (absence in that tag) | `git ls-tree -r --name-only v0.5.0-alpha3 \| grep -i identity` → 0; `git grep -n -i -E 'mnemonic\|bip39' v0.5.0-alpha3 -- '*.go'` → 0; tag resolves to `d8e8288918820524fcb52376dcba8e3767ef7d8f` (2026-06-20) | `not found in repo` (for that tag) |

No row above is speculative. Positive rows are direct file/ref reads; negative rows carry the
exact grep whose emptiness is the evidence, and are labelled `not found in repo` rather than dropped.

---

## Parent reconciliation — where the three inputs conflicted, and who is right

Re-opened every contested citation myself; no averaging, no hedging.

1. **`identity_test.go` line for the round-trip assertion.** One parent (the audit, quoted in
   `EVIDENCE-NOTE-t_93620520.md`) said `:56`; another (`REPORT-t_e75c5cd6.md`) said `:57`.
   **Truth: `:57`.** `assert.Len(t, words, 12, "expected 12 BIP39 words")` is on line 57; line 56
   is `words := strings.Fields(mnemonic)`.
2. **`reveal-seed` handler line numbers.** Parents cite `main.go:2189`/`:2195` (working tree) and
   `:1862-1891` (an older HEAD). **Truth at working HEAD `eb12d7f8`:** comment at `:2175`, `func`
   at `:2179`, 405 guard at `:2186`, 1 KB cap at `:2189`, `DeriveFromMnemonic` at `:2195`.
   Route registration is at `:1976` and `:1982` (not `:1935`/`:1941`, which is `origin/main`
   numbering — `origin/main` registers at `:1876`/`:1882`). Both parents' numbers are *correct for
   the ref they read*; they are not interchangeable, and this report uses the working-tree numbers.
3. **The 24→12 switch commit (`7a87fb4d`) — this is the real conflict, and two parents are wrong.**
   `REPORT-t_d60caf8a.md` and `REPORT-t_e75c5cd6.md` both present `7a87fb4d` (2026-07-08) as the
   backend's switch, implying the 24-word doc was true until then. `git merge-base --is-ancestor
   7a87fb4d origin/main` → **NO**; that commit lives only on `origin/feat/identity-v2-carried`.
   On `origin/main`, `src/identity/identity.go` has **two** commits and was introduced whole as
   12-word by `d491a4c4` (2026-08-17). So for `origin/main` the doc's 24-word figure was **never**
   true — the 24-word code (`PrivateKeyToMnemonic`, `bip39.NewMnemonic(b)` over the 32-byte key)
   existed only on the carried feature branch. The correction matters for the candidate text: the
   defensible statement is "the wizard's AGENTS.md froze a figure from the pre-NIP-06 design that
   never landed on the backend's default branch", not "main switched on 2026-07-08".
   The 24-word era itself is real and citable — at `7a87fb4d^:src/identity/identity.go:124-133`.
4. **The `cdk-ffi` version literal.** `REPORT-t_e75c5cd6.md` cites
   `/tmp/cdkffi173/cdk-ffi-0.17.3/src/wallet.rs:873-879` for `Mnemonic::generate(12)`. **That path
   no longer exists** (`ls -d /tmp/cdkffi173` → No such file or directory), and the cargo registry
   on this host holds `cdk-ffi-0.18.1` only (`wallet.rs:1119-1120`, same literal). The *behaviour*
   at the pinned version is nonetheless confirmed first-hand by calling the shipped
   `libcdk_ffi.so` from `cdk-go@v0.17.3` — 6/6 runs returned 12 words (§Probes). Verdict: behaviour
   TRUE, source literal at 0.17.3 unverified → residual uncertainty.
5. **`AGENTS.md` provenance.** All three parents agree it is a single commit; confirmed:
   `git log --oneline -- AGENTS.md` → `475a9ef docs: add AGENTS.md ...` only, blob
   `f2ff1775...` identical to HEAD. `475a9ef` is dated **2026-07-06**, i.e. before both the
   carried-branch switch (2026-07-08) and the `origin/main` introduction (2026-08-17) — the doc
   was written against the 24-word design and never revisited.

---

## Probes (executed, not inferred)

Go probe built against **only** the pinned `github.com/tyler-smith/go-bip39 v1.1.0`
(`GOPROXY=off`, module cache only, go1.22.12 at `/usr/local/go/bin/go`):

```
bits=128 bytes=16 words=12 valid=true
bits=160 bytes=20 words=15 valid=true
bits=192 bytes=24 words=18 valid=true
bits=224 bytes=28 words=21 valid=true
bits=256 bytes=32 words=24 valid=true
24-word phrase accepted by IsMnemonicValid: true (words=24)
```

The last line is the reason the wrong figure survived review: **24 words is valid input**, so
nothing rejects the doc's number at the parse layer — only the generator pins 12.

Second generator, called directly through `ctypes` against the library the backend links
(`cdk-go@v0.17.3/bindings/cdkffi/native/linux_amd64/libcdk_ffi.so`,
`uniffi_cdk_ffi_fn_func_generate_mnemonic`, 6 invocations):

```
[(12, 'ugly'), (12, 'prize'), (12, 'toast'), (12, 'license'), (12, 'mutual'), (12, 'ribbon')]
```

Backend tests, first-hand, in the identity module (`GOWORK=off` required — `src/go.work`
otherwise resolves the wrong main module):

```
=== RUN   TestMnemonic_RoundTrip
--- PASS: TestMnemonic_RoundTrip (0.41s)
=== RUN   TestMnemonicToPrivateKey_RejectsInvalid
--- PASS: TestMnemonicToPrivateKey_RejectsInvalid (0.00s)
=== RUN   TestGenerateMnemonic_TwelveWords
--- PASS: TestGenerateMnemonic_TwelveWords (0.00s)
PASS
ok  	github.com/OpenTollGate/tollgate-module-basic-go/src/identity	0.427s
```

---

## Candidate statements for AGENTS.md (Part B)

One line each, present tense, factual, every line carrying a `path:line` citation and naming the
library + pinned version that fixes the behaviour. **PROPOSED text only — nothing was applied;
`AGENTS.md` is untouched at blob `f2ff1775...`.**

1. **The device recovery mnemonic is BIP39 12 words / 128-bit entropy, generated by
   `tollgate-module-basic-go`, never by the wizard** — `tollgate-module-basic-go/src/identity/identity.go:137` (`entropy, err := bip39.NewEntropy(128)`), documented at `identity.go:133`.
2. **The generator is pinned to `github.com/tyler-smith/go-bip39 v1.1.0`** — `tollgate-module-basic-go/src/identity/go.mod:8` (direct), `tollgate-module-basic-go/src/go.mod:115` (indirect in the root module).
3. **`go-bip39` v1.1.0 has no default strength — `128` is passed literally, and 128 bits encodes to exactly 12 words** — `go-bip39@v1.1.0/bip39.go:99` (mandatory `bitSize`), `:165-167` (word-count computation).
4. **The 12-word output is enforced by tests, not by the recovery endpoint** — `tollgate-module-basic-go/src/identity/identity_test.go:57`, `:88`, `:238`; `identity_handler_test.go:78`.
5. **The recovery endpoint accepts any valid BIP39 length (12/15/18/21/24 words); 12 words is a generation-time pin, not an input restriction** — `tollgate-module-basic-go/src/identity/identity.go:151` with `go-bip39@v1.1.0/bip39.go:275-277` and `:356`.
6. **The correct backend symbols are `GenerateMnemonic` and `MnemonicToPrivateKey`; `PrivateKeyToMnemonic` does not exist** — `tollgate-module-basic-go/src/identity/identity.go:136`, `:150` (grep for `PrivateKeyToMnemonic` over `*.go` → 0 hits).
7. **A 24-word generator exists in the pinned dependency but is dead code: `nip06.GenerateSeedWords` calls `bip39.NewEntropy(256)` and has zero callers** — `github.com/nbd-wtf/go-nostr@v0.51.12/nip06/nip06.go:10-11`; caller grep → 0.
8. **The wallet-side mnemonic generator (`cdk-ffi`, via `github.com/cashubtc/cdk-go v0.17.3`) also emits 12 words** — `tollgate-module-basic-go/src/tollwallet/cdk_wallet.go:55`, `tollgate-module-basic-go/src/go.mod:68`; measured 6/6 runs at 12 words against the shipped `libcdk_ffi.so`.
9. **This repo vendors no BIP39 wordlist and hardcodes no word count; the only wordlist in play is the library's 2048-word English list** — wizard `go.sum` (0 hits), `net4sats-wizard` binary (0 hits for `tyler-smith/go-bip39`), `go-bip39@v1.1.0/wordlists/english.go:23` (CRC32 `c1dbd296` at `:13-15`).
10. **The "24 words" figure predates NIP-06 and never matched `origin/main`: the 24-word code exists only on `origin/feat/identity-v2-carried` (`7a87fb4d^`), while `origin/main` introduced `src/identity` already 12-word in `d491a4c4`** — `git merge-base --is-ancestor 7a87fb4d origin/main` → NO; `git log --oneline origin/main -- src/identity/identity.go` → `c60ea071`, `d491a4c4`.

**Do not write** (each FALSE per the table, or unverifiable here):
`BIP39 mnemonic (24 words) → seed phrase (PrivateKeyToMnemonic)` (`AGENTS.md:43`);
`PrivateKeyToMnemonic` as a backend function (`AGENTS.md:29`);
any statement that the 128-bit pin applies to `v0.5.0-alpha3` (that tag has no identity package);
any statement sourced from `configurationwizzard/src/components/identity-panel.tsx` — no checkout
exists on this host, so it is unverifiable here, not false.

---

## Residual uncertainty

**Residual uncertainty:** the `cdk-ffi` **0.17.3** Rust source is not on this host (only the 0.18.1
crate is in the cargo registry, and the parent's `/tmp/cdkffi173` extraction is gone), so the
literal `Mnemonic::generate(12)` is verified at the pinned **behaviour** (6/6 calls into
`cdk-go@v0.17.3`'s shipped `libcdk_ffi.so`) but not at the pinned **source line**; fetching the
`cdk-ffi 0.17.3` crate tarball (sha256 `bb89b448…`) or reading the `CDK_COMMIT`
(`34882e8a…`) tree would settle it. Secondary, non-blocking: `configurationwizzard`'s
`identity-panel.tsx:174` "24 words" cannot be confirmed without a checkout of that repo at a named
ref; `go-bip39` is marked `// indirect` in `src/go.mod:115` while imported directly at
`identity.go:137`; and the `v0.5.0-alpha3` compatibility claim at `AGENTS.md:33` is unadjudicated
beyond the fact that the tag ships no identity package, so what a deployed alpha3 device emits is
outside this evidence.

---

## Repo-state attestation

- `git status --porcelain` at end: only this report + `PROGRESS-t_wz17-bip39-groundtruth.md`
  added, plus the pre-existing untracked files of other cards (`.hermes/`, `docs/WZ1.8-...`,
  `verify/...`) which were already untracked at session start and were not touched.
- No edit to `AGENTS.md`, any `.go`, `.js`, `.html`, `go.mod`, `go.sum`, config, or hook.
- No build of the wizard, no write inside the backend repo. The only writes outside this repo
  were throwaway probes under `/tmp`.
