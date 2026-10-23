# REPORT — t_e75c5cd6: pinned mnemonic library default strength & word-count semantics

**Date:** 2026-10-23
**Repo under inspection:** `net4sats-mvp-v2` (worktree
`/home/c08r4d0r/repos/net4sats-mvp-v2`, branch `worker-heavy/t_gate8-19b1b1c6`)
**Card:** `t_e75c5cd6` — "Derive the pinned mnemonic library's default strength and
word-count semantics" (parent: `t_48cbf9fe`)
**Mode:** READ-ONLY. No repo code, doc, config or `AGENTS.md` edited. All probes ran
outside the worktree in `/tmp`. Verified library behaviour was executed against the
pinned versions, not inferred.

---

## 0. INPUT caveat — the parent task has no report

The brief says "Use the parent task's report". **No such report exists.**

- `~/.hermes/kanban/boards/net4sats-mvp-v2/kanban.db`, card `t_48cbf9fe`
  ("Locate mnemonic generation/parse sites and pin the mnemonic library version"):
  `status = 'triage'`, `result = NULL`, `session_id = NULL`. It never ran.
- Its body points at `~/.hermes/plans/2026-08-31-tollgate-batch/schedule.md` item
  WZ1.7. That file was already established as missing by prior cards in this repo
  (`REPORT.md:29-30`: "`grep -c "WZ1.7"` → does not exist").

So the library name, version and call sites below were **re-derived first-hand** in
this session. Every claim carries a `path:line`.

**Structural finding that changes the shape of the answer:** the wizard repo
(`net4sats-mvp-v2`) contains **zero** mnemonic code. `git grep` over tracked
non-report sources for `bip39|mnemonic|NewEntropy|wordlist` returns
`NO_HITS_IN_TRACKED_SOURCES`; the built `net4sats-wizard` binary contains **0**
occurrences of `tyler-smith/go-bip39`; the bundled JS/JS-map assets contain **0**
mnemonic generators (their only `bip39` strings are `bip39seed` option keys inside the
vendored cashu-ts wallet, `portal/assets/splash-C0y8YH4o.js.map`). The mnemonic work
lives entirely in the backend, `tollgate-module-basic-go` — which matches this repo's
`AGENTS.md` "thin UI wrapper" principle.

**Two libraries are in play, both pinned, both generating 12 words.** Neither one has
a "default strength", because neither API accepts an omitted strength. Details below.

---

## 1. Authoritative source

| # | Library | Exact version inspected | Where read |
|---|---|---|---|
| A | `github.com/tyler-smith/go-bip39` | **v1.1.0** | `$GOMODCACHE/github.com/tyler-smith/go-bip39@v1.1.0/bip39.go` — `/home/c08r4d0r/go/pkg/mod/github.com/tyler-smith/go-bip39@v1.1.0/bip39.go:99`, `:163`, `:275`, `:315`, `:348` |
| B | `cdk-ffi` (Rust, via Go binding `github.com/cashubtc/cdk-go` v0.17.3) | **cdk-ffi 0.17.3**, crate `bip39` **2.2.2** | `/tmp/cdkffi173/cdk-ffi-0.17.3/src/wallet.rs:873-879` (crate fetched from `https://static.crates.io/crates/cdk-ffi/cdk-ffi-0.17.3.crate`, sha256 `bb89b4483ccf2fca31baba23227a343bc34b82b6e6ed46fa4d6d91e7ddcabc4c`); rust-bip39 at `~/.cargo/registry/src/index.crates.io-1949cf8c6b5b557f/bip39-2.2.2/src/lib.rs:270`, `:682` |
| C | `github.com/nbd-wtf/go-nostr/nip06` (wrapper over A) | **v0.51.12** | `$GOMODCACHE/github.com/nbd-wtf/go-nostr@v0.51.12/nip06/nip06.go:10-22`, `:53-55` |

Version provenance — why these are the pinned versions and not guesses:

- `tollgate-module-basic-go/src/go.mod:115` → `github.com/tyler-smith/go-bip39 v1.1.0 // indirect`
- `tollgate-module-basic-go/src/go.sum:413` → `h1:5eUemwrMargf3BSLRRCalXT93Ns6pQJIjYQN2nyfOP8=`,
  byte-identical to the module-cache ziphash `.../@v/v1.1.0.ziphash` — **the source I
  read is the source the backend builds against.**
- `tollgate-module-basic-go/src/go.mod:17` → `github.com/nbd-wtf/go-nostr v0.51.12`
- `tollgate-module-basic-go/src/go.mod:68` → `github.com/cashubtc/cdk-go v0.17.3 // indirect`
- `cdk-go@v0.17.3/rust/Cargo.toml` pins `cdk-ffi = { version = "=0.17.3", ... }`, and
  `cdk-go@v0.17.3/bindings/cdkffi/CDK_COMMIT` = `34882e8a90419d15fa2053b067776adccd957a2f`.
- Backend tree read at: `tollgate-module-basic-go` HEAD **`c77ef2044ae6b593f27b83ce4f721856d7c83603`**.
  `src/identity/identity.go` is clean at HEAD (last touched by `d491a4c4`, PR #331).
  `src/main.go` is staged-modified in that tree (`M src/main.go`, +42/−1); both HEAD and
  worktree line numbers are cited where they differ.

---

## 2. Default behavior — verbatim, from source

**A. `go-bip39` v1.1.0 has no default strength. The parameter is mandatory and
validated into a closed range.** `NewEntropy` is the only generator entry point, and it
takes `bitSize int` with no variadic/optional form:

```go
// bip39.go:98-107
// bitSize has to be a multiple 32 and be within the inclusive range of {128, 256}
func NewEntropy(bitSize int) ([]byte, error) {
	err := validateEntropyBitSize(bitSize)
```

```go
// bip39.go:315-320
func validateEntropyBitSize(bitSize int) error {
	if (bitSize%32) != 0 || bitSize < 128 || bitSize > 256 {
		return ErrEntropyLengthInvalid
```

Accepted values are therefore exactly **128, 160, 192, 224, 256** bits. The
word-count mapping is computed, not tabulated (`bip39.go:165-167`):
`entropyBitLength = len(entropy)*8`, `checksumBitLength = entropyBitLength/32`,
`sentenceLength = (entropyBitLength + checksumBitLength)/11`. It is echoed as a
closed enumeration on the parse side (`bip39.go:29-35`,
`12/15/18/21/24` keys of `wordLengthChecksumMasksMapping`; `:355-357`
`numOfWords%3 != 0 || numOfWords < 12 || numOfWords > 24`).

**B. `cdk-ffi` 0.17.3 also has no default — it hardcodes the argument.** The wrapper's
own signature takes no strength parameter at all, and supplies `12` literally:

```rust
// /tmp/cdkffi173/cdk-ffi-0.17.3/src/wallet.rs:873-877
/// Generates a new random mnemonic phrase
#[uniffi::export]
pub fn generate_mnemonic() -> Result<String, FfiError> {
    let mnemonic = Mnemonic::generate(12)
```

The underlying `bip39` 2.2.2 crate likewise has **no** `Default` impl and no default
word count — `generate(word_count: usize)` is mandatory
(`bip39-2.2.2/src/lib.rs:314-316`), validated by
`is_invalid_word_count` → `word_count < 12 || word_count % 3 != 0 || word_count > 24`
(`:682-684`, `MIN_NB_WORDS=12` at `:83`, `MAX_NB_WORDS=24` at `:86`).

**C. `nip06` v0.51.12 does have a fixed 256-bit generator — and it is dead code here.**
`GenerateSeedWords()` calls `bip39.NewEntropy(256)` (`nip06.go:10-11`), i.e. 24 words.
`grep -rn "GenerateSeedWords" --include=*.go` over the backend → **NO_CALLERS**.

### Executed confirmation (not inferred)

Go probe built against the pinned v1.1.0 in `/tmp/bip39probe` (`GOPROXY=off`, module
cache only), `go run .`:

```
NewEntropy(128) -> entropy bytes=16 words=12
NewEntropy(160) -> entropy bytes=20 words=15
NewEntropy(192) -> entropy bytes=24 words=18
NewEntropy(224) -> entropy bytes=28 words=21
NewEntropy(256) -> entropy bytes=32 words=24
NewEntropy(  0) -> err=Entropy length must be [128, 256] and a multiple of 32
NewEntropy(127) -> err=Entropy length must be [128, 256] and a multiple of 32
NewEntropy(288) -> err=Entropy length must be [128, 256] and a multiple of 32
256-bit roundtrip: words=24 entropy_back_bytes=32 err=<nil> valid=true
```

Parse-side probe (same binary): `IsMnemonicValid` returns `true` for 12/15/18/21/24-word
phrases and `false` for `""`, `"abandon"`, `"not a valid mnemonic at all"`; a freshly
built 24-word phrase is **accepted**.

CDK probe — `uniffi_cdk_ffi_fn_func_generate_mnemonic` called directly through
`ctypes` against the shipped
`cdk-go@v0.17.3/bindings/cdkffi/native/linux_amd64/libcdk_ffi.so`, 6 invocations:

```
run 0: words=12 entropy_bits=128 ... run 5: words=12 entropy_bits=128
word-count set observed: [12]
```

So the shipped native library that `cdk_ffi.GenerateMnemonic()` calls returns 12 words
every time — confirming the `Mnemonic::generate(12)` literal in the 0.17.3 source.

---

## 3. Effective strength per call site

Every call site that touches the mnemonic, mapped to the literal argument or the
wrapper's hardcoded value. **No call site anywhere relies on a library default, because
no library in play has one.**

| # | path:line | argument passed / how strength is determined | effective entropy bits | expected word count |
|---|---|---|---|---|
| 1 | `tollgate-module-basic-go/src/identity/identity.go:137` | **literal `128`** — `entropy, err := bip39.NewEntropy(128)` | **128** | **12** |
| 2 | `tollgate-module-basic-go/src/identity/identity.go:141` | `bip39.NewMnemonic(entropy)` — entropy flows from #1, so 16 bytes | **128** | **12** |
| 3 | `tollgate-module-basic-go/src/identity/identity.go:208` (`DeriveFromMnemonic`) | parse path only; no strength is chosen | n/a (depends on caller's phrase) | accepts 12/15/18/21/24 |
| 4 | `tollgate-module-basic-go/src/identity/identity.go:151` (`MnemonicToPrivateKey` → `nip06.ValidateWords` → `bip39.IsMnemonicValid`) | validation only; **no** strength argument exists | n/a — accepts all five BIP39 lengths | 12/15/18/21/24 all valid |
| 5 | `tollgate-module-basic-go/src/tollwallet/cdk_wallet.go:55` (`loadOrCreateMnemonic`) | **omitted at the call** — `cdk_ffi.GenerateMnemonic()` takes no argument; the strength is hardcoded one layer down at `cdk-ffi-0.17.3/src/wallet.rs:876` as `Mnemonic::generate(12)` | **128** | **12** |
| 6 | `tollgate-module-basic-go/src/identity/identity.go:160` | `bip39.GetWordList()` — wordlist access, not generation | n/a | n/a |
| 7 | `tollgate-module-basic-go/src/identity/identity.go:171` | `bip39.GetWordList()` — wordlist access, not generation | n/a | n/a |
| 8 | `tollgate-module-basic-go/src/main.go:2195` (`handleIdentityRevealSeed`) | `identity.DeriveFromMnemonic(mnemonic)` on the POST body — parse of caller-supplied text, 1 KB cap at `:2192` | caller-determined | caller-determined |
| — | **`net4sats-mvp-v2` (this repo)** | **no call site exists** — 0 hits in tracked sources, 0 symbols in the shipped binary | — | — |

Row 5 is the only place the phrase "library default" could even arise, and it is a
**wrapper default, not a library default**: the effective strength is fixed at 12 words
by `cdk-ffi` itself, so a caller who passes nothing still gets 128 bits. Confirmed
empirically (6/6 runs, all 12 words).

Rows 3/4/8 are **parse** sites. They do **not** pin 12 words: `IsMnemonicValid` accepts
12/15/18/21/24 (`bip39.go:275-278` → `EntropyFromMnemonic` → `splitMnemonicWords`
`:348-359`, masks at `:29-35`). The 12-word figure comes from the **generators**
(rows 1-2, 5), never from the recovery endpoint.

Hardened by tests in the backend:
`src/identity/identity_test.go:57` `assert.Len(t, words, 12, "expected 12 BIP39 words")`,
`:88` `assert.Len(t, words, 12, "mnemonic must be exactly 12 words")`,
`:238` `require.Len(t, words, 12, ...)`, and
`src/identity_handler_test.go:78` `require.Len(t, words, 12, "mnemonic must be 12 words")`.

---

## 4. Verdict

> **The repo's mnemonic pipeline generates 12-word / 128-bit BIP39 mnemonics. The
> strongest single citation is `tollgate-module-basic-go/src/identity/identity.go:137`
> — `entropy, err := bip39.NewEntropy(128)` — pinned to
> `github.com/tyler-smith/go-bip39 v1.1.0` at `src/go.mod:115`.**

The second generator agrees independently:
`src/tollwallet/cdk_wallet.go:55` → `cdk-ffi-0.17.3/src/wallet.rs:876`
`Mnemonic::generate(12)`.

**`AGENTS.md:43` is wrong.** It states `BIP39 mnemonic (24 words) → seed phrase
(PrivateKeyToMnemonic)`. 24 words would require `NewEntropy(256)`; the only code path
that does that is `nip06.GenerateSeedWords` (`nip06.go:10-11`), which has **zero
callers**. (`PrivateKeyToMnemonic` also does not exist — prior cards established this;
the implemented direction is `MnemonicToPrivateKey` at `identity.go:150`.) This report
does not edit `AGENTS.md`, per the read-only constraint.

**No dynamic strength exists.** `grep -rn '"words"|--words|"strength"|--strength|
"entropy"|--entropy|wordCount|word_count|WordCount'` over the backend Go sources →
zero hits; `grep -rn 'flag\.\(Int\|String\)'` over `src/` → zero hits. There is no CLI
flag, config key, or API parameter that changes the generated word count.

---

## 5. Unresolved / explicitly out of scope

1. **Parent card `t_48cbf9fe` never ran** (`result = NULL`). Its deliverable — the
   call-site inventory this card was told to consume — does not exist. The inventory in
   §3 is this session's own; it was produced by `git grep` over the backend at
   `c77ef204`, not inherited.
2. **Deployed `v0.5.0-alpha3` not inspected.** `AGENTS.md:33` claims compatibility with
   that older build. Whether that revision also pinned 128 bits was not checked — the
   backend repo's tags were not enumerated (no network during this session). Statements
   here are about HEAD `c77ef204` only.
3. **`cdk-ffi` 0.17.3 was read from the published crate, not from the tagged git tree.**
   `github.com` DNS resolution failed in this session (`Could not resolve host:
   github.com`), so `raw.githubusercontent.com/.../34882e8a.../crates/cdk-ffi/src/lib.rs`
   could not be fetched. The crates.io tarball for 0.17.3 is the same release and its
   `Cargo.toml:16` version matches, but the git commit is not byte-verified. The
   *behaviour* is independently confirmed by calling the shipped `.so` (§2C), so the
   12-word conclusion does not rest on the tarball alone.
4. **`go-bip39` is marked `// indirect`** in `src/go.mod:115` while being imported
   directly at `identity.go:137`. Whether that marker is stale is unresolved; it does
   not affect the pinned version, which is explicit.
5. **Pushes could not be verified.** `git ls-remote origin` succeeded early in the
   session (`HEAD` = `65b9d02740211a88fa86bab05862f23a9326af8b`) but failed on retry with
   `Could not resolve host: github.com`. The commit for this report is created locally;
   see §6 for the observed push outcome.

---

## 6. Method / reproducibility

Read-only against the worktree. Commands: `git grep -n -i -E 'bip39|mnemonic|NewEntropy|
wordlist'` over both trees; `git ls-files`, `strings -a net4sats-wizard | grep -c
'tyler-smith/go-bip39'` (→ 0); `grep -o -i 'bip39[a-zA-Z_]*'` over
`portal/assets/*.js.map` and `admin/assets/*.js` (→ only `bip39seed`); `read_file` at
explicit line ranges on `bip39.go`, `nip06.go`, `identity.go`, `main.go`, `wallet.rs`,
`bip39-2.2.2/src/lib.rs`; module-cache ziphash vs `go.sum` comparison; `ctypes` FFI call
into `libcdk_ffi.so`; and a throwaway Go module at `/tmp/bip39probe` (`GOPROXY=off`)
executed with `/usr/local/go/bin/go` (go1.24.2 — `go` is not on `PATH` by default on this
host). No build, test or `git` write inside either inspected repo, other than the
single docs-only commit described below.
