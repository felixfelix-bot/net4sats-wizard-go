# Measured captures — wz1.5 channel decision, second-pass re-verification (2026-10-24)

Verbatim captures behind
`~/.hermes/plans/2026-08-31-tollgate-batch/wz1.5-channel-decision.md` (the wz1.5 channel
decision note, card `t_1df48589`). Everything below was run on **2026-10-24** from this build
host; nothing is reconstructed from an older report. The first pass of this card ran 2026-10-23.

Purpose of this pass: the decision note was written by an earlier pass but **never committed**
(the plans directory is not a git repo), and its claims had not been independently re-derived.
This file is the re-derivation. **One defect was found** (§3) and is fixed in the note.

Environment note: the root filesystem was at **935 MB free / 100% used** when this pass started.
Every download here is ≤ 5.4 MB and each scratch file was deleted immediately after hashing.

---

## 1. The decision-relevant status matrix

```
$ w() { printf '%-118s ' "$1"; curl -sSL -m 20 -o /dev/null -r 0-1023 -w 'http=%{http_code} bytes=%{size_download}\n' "$1"; }

https://releases.tollgate.me/                                                                        http=206 bytes=697
https://releases.tollgate.me/alpha/v0.7.0-alpha10/aarch64_cortex-a53/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk http=404 bytes=697
https://releases.tollgate.me/alpha/v0.6.1-post-merge/aarch64_cortex-a53/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk http=404 bytes=697
https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages                         http=206 bytes=1322
https://releases.tollgate.me/alpha/v0.6.0/aarch64_cortex-a53/Packages                                http=404 bytes=697
https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53/Packages.sig                     http=404 bytes=697
https://releases.tollgate.me/keys/tollgate-testing.pem                                               http=404 bytes=697
https://github.com/FreedomTechFeed/packages/releases/download/latest/SHA256SUMS                      http=404 bytes=9
https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0/tollgate-wrt_0.6.0_aarch64_cortex-a53-mediatek-filogic.ipk http=206 bytes=1024
https://github.com/FreedomTechFeed/packages/releases/download/v0.6.0-rc6/tollgate-wrt_0.6.0_rc6_aarch64_cortex-a53-mediatek-filogic.ipk http=206 bytes=1024
```

Reading: the two literals `deploy.go:43/47` would become under the chosen pattern are **404**;
the only live index on the chosen host is the stale unsigned `v0.5.0-alpha3`; `Packages.sig` and
the signing key are **404**; and `releases/download/latest/…` is a flat **404**.

## 2. Exactly one version carries artifacts on the chosen host

27 URLs = 9 tags × 3 filename shapes, `curl -sSI` status only, non-404s printed:

```
$ for tag in v0.4.0 v0.5.0-alpha1 v0.5.0-alpha2 v0.5.0 v0.5.0-beta2 v0.6.0-rc6 v0.6.0 v0.6.1-post-merge v0.7.0-alpha10; do
    for f in "tollgate-wrt_${tag#v}_aarch64_cortex-a53.ipk" \
             "tollgate-wrt_${tag#v}_aarch64_cortex-a53.apk" \
             "tollgate-wrt_${tag#v}_aarch64_cortex-a53-upx-ultra-brute.ipk"; do
      code=$(curl -sSI -m 15 -o /dev/null -w '%{http_code}' "https://releases.tollgate.me/alpha/$tag/aarch64_cortex-a53/$f")
      [ "$code" != "404" ] && echo "HIT $code  $f"
    done
  done
=== done ===
```

**Zero hits.** Combined with §3 (where the `v0.5.0-alpha3` names that *do* exist are read out of
the index), this confirms the note's claim that the "two real example versioned URLs from two
different versions" acceptance criterion **cannot be met from the chosen channel** — exactly one
version has artifacts there. The two examples in the note are therefore carried from the
`FreedomTechFeed` per-asset lane, which is the only lane with more than one published version.

## 3. DEFECT FOUND — the artifact filename is per-artifact, not derivable

The note's first draft instantiated its pattern as `tollgate-wrt_v<X.Y.Z>[-<pre>]_<arch>.ipk`.
That is **wrong for the only version that exists**:

```
$ B=https://releases.tollgate.me/alpha/v0.5.0-alpha3/aarch64_cortex-a53
$ curl -sSL -m 25 "$B/Packages" | grep -E '^(Package|Version|Filename|Size|SHA256sum):'
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk
Size: 5321054
SHA256sum: c508e0b51613a93731e56e4f9c2fd44134cb99a6d0bae3dcc39e1cc6d8a97fb3
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk
Size: 5307464
SHA256sum: c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f
Package: tollgate-wrt
Version: v0.5.0-alpha3
Filename: tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk
Size: 7153725
SHA256sum: 5332bb5ea15edbef3a41d5a34cf0e5e66a8d0d9a152385d4e0a7f06e928e6f3d
```

Each index entry's `content-length` equals its own `Size:`:

```
$ for f in <the three Filenames above>; do printf '%-70s ' "$f"; curl -sSI -m 25 "$B/$f" | awk 'BEGIN{IGNORECASE=1} /^HTTP\//{c=$2} /^content-length:/{l=$2} END{printf "http=%s clen=%s\n", c, l}'; done
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.ipk      http=200 clen=5321054
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk      http=200 clen=5307464
tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.apk                      http=200 clen=7153725
```

The composed name is a **404**:

```
$ curl -sSI "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53.ipk"
(non-200)   <-- the note's first-draft template would have produced this URL
```

Asymmetry worth keeping: for **apk** the plain (7,153,725 B) and suffixed (5,307,464 B) names
**both** exist — two different builds of the same version — while for **ipk** only the suffixed
name exists. So the suffix is neither absent-by-default nor uniform, and only the index knows.

Full download + digest (not a ranged GET), and the header proof it is a real 5.3 MB artifact:

```
$ curl -sSL -m 120 -o chk.apk "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk" -w 'http=%{http_code} bytes=%{size_download}\n'
http=200 bytes=5307464
$ stat -c 'on-disk=%s' chk.apk
on-disk=5307464
$ sha256sum chk.apk
c9e09282b8e0ca8b0af45f388a317b995454bfcca71d1f1e6852b4b99dbbfe9f  chk.apk
   # == the index's own SHA256sum: for that entry

$ curl -sSI -m 30 "$B/tollgate-wrt_v0.5.0-alpha3_aarch64_cortex-a53-upx-ultra-brute.apk"
HTTP/2 200
content-type: application/vnd.android.package-archive
content-length: 5307464
last-modified: Wed, 01 Jul 2026 19:52:40 GMT
```

## 4. The rolling `latest` target is not addressable as an artifact

```
$ curl -sSI -m 30 https://github.com/FreedomTechFeed/packages/releases/latest
HTTP/2 302
$ curl -sSL -m 30 -o /dev/null -w 'final=%{url_effective} code=%{http_code}\n' https://github.com/FreedomTechFeed/packages/releases/latest
final=https://github.com/FreedomTechFeed/packages/releases/tag/v0.6.0 code=200
$ curl -sSL -m 20 -o /dev/null -r 0-1023 -w 'http=%{http_code} bytes=%{size_download}\n' https://github.com/FreedomTechFeed/packages/releases/download/latest/SHA256SUMS
http=404 bytes=9
```

`latest` resolves only to an HTML **tag page**; the `/download/latest/…` form is a flat 404. This
is the mechanical basis for the note's rule "never reference a rolling `latest` in code".

## 5. Release inventories (unchanged from the first pass)

```
$ curl -sS -m 30 "https://api.github.com/repos/FreedomTechFeed/packages/releases?per_page=100" | jq -r '.[] | "\(.tag_name)\t\(.published_at)\tassets=\(.assets|length)\tdraft=\(.draft)"'
v0.6.0      2026-10-21T02:06:45Z  assets=26  draft=false
v0.6.0-rc6  2026-10-15T19:50:37Z  assets=26  draft=false

$ curl -sS -m 30 "https://api.github.com/repos/OpenTollGate/tollgate-module-basic-go/releases?per_page=5" | jq -r '.[] | "\(.tag_name)\t\(.published_at)\tassets=\(.assets|length)\tdraft=\(.draft)"'
v0.6.0        2026-10-21T03:27:56Z  assets=0   draft=false
v0.6.0-alpha1 2026-09-08T11:42:45Z  assets=0   draft=false
v0.5.0        2026-07-03T11:39:09Z  assets=16  draft=false
v0.5.0-alpha3 2026-06-27T14:50:35Z  assets=0   draft=false
v0.5.0-alpha2 2026-06-17T13:53:17Z  assets=0   draft=false

$ for t in v0.7.0-alpha10 v0.6.1-post-merge; do curl -sS -m 30 -o /dev/null -w "$t api=%{http_code}\n" "https://api.github.com/repos/FreedomTechFeed/packages/releases/tags/$t"; done
v0.7.0-alpha10       api=404
v0.6.1-post-merge    api=404
```

Two load-bearing facts confirmed: **upstream `OpenTollGate/tollgate-module-basic-go` `v0.6.0`
carries 0 assets** (so "upstream releases" is not a usable artifact channel), and **neither
wizard-pinned tag exists on the `FreedomTechFeed` lane** (so that lane cannot serve today's
`deploy.go` pins either).

## 6. Publisher channel vocabulary read out of the live bundle

```
$ curl -sSL -m 25 https://releases.tollgate.me/static/js/main.ca30081d.js -o exp.js; stat -c 'bytes=%s' exp.js
bytes=474543
$ grep -o 'STABLE:"[a-z]*",BETA:"[a-z]*",ALPHA:"[a-z]*",DEV:"[a-z]*"' exp.js
STABLE:"stable",BETA:"beta",ALPHA:"alpha",DEV:"dev"
$ grep -oE '(Cr|Ar|Fr)="[^"]{0,40}"' exp.js | sort -u
Ar="tollgate-wrt"
Cr="tollgate-os"
Fr="tollgate-module-basic-go"
```

The `alpha` path segment used by the chosen base URL is the publisher's own channel vocabulary,
not an invention of the decision note.

## 7. Transient-failure note (not a finding)

One probe batch in this pass failed with `curl: (6) Could not resolve host: api.github.com`
while `releases.tollgate.me` and `github.com` resolved normally in the same seconds. Retried with
backoff and succeeded unchanged; the §5 inventories above are from the successful retry. Recorded
so a future lane does not read a single DNS blip as a dead endpoint.

---

## Result

- **Decision unchanged:** interim channel `releases.tollgate.me`, pattern
  `https://releases.tollgate.me/alpha/v<X.Y.Z>[-<pre>]/<arch>/<artifact-file>`; **Q5 = OPEN**.
- **One correction:** `<artifact-file>` must be read from the version's own index `Filename:`
  field, never composed — the note now says so in §1, §2.3a and §3 rule 1b.
- **No repo source modified:** `deploy.go` and `pins_test.go` are untouched, so no build or test
  run could have changed behaviour.
