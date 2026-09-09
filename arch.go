package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

// ─── Router CPU architecture auto-detection ─────────────────────────────
//
// The wizard used to hardcode "aarch64_cortex-a53" in every package download
// URL. That is wrong for any other router — a mipsel_24kc binary, for example,
// will not execute on an ARM router, and on an x86_64 gl-gate it would be the
// same class of silent failure. detectArch reads the target's real CPU
// architecture at deploy time and the deployment then selects the matching
// tollgate-wrt asset. See docs/arch-detect.md for the full reasoning and the
// precedence ladder.

// bareArchToTuple normalizes a BARE CPU architecture name (as reported by
// `apk --print-arch` or `uname -m`) to the canonical OpenWrt arch tuple that
// appears in package feed path names and in tollgate-wrt release assets.
//
// `apk --print-arch` on a GL-MT3000 prints the bare "aarch64", not the
// "aarch64_cortex-a53" tuple the download URLs (and the OpenWrt package
// repository layout) use. Mapping that is the whole point of this function — a
// missing entry here is the bug this code deletes.
var bareArchToTuple = map[string]string{
	"aarch64": "aarch64_cortex-a53",
	"arm64":   "aarch64_cortex-a53",
	"mipsel":  "mipsel_24kc",
	"mips":    "mips_24kc",
	"x86_64":  "x86_64",
	"amd64":   "x86_64",
}

// canonicalArchTuples is the set of canonical OpenWrt arch tuples that
// normalizes to themselves. This lets canonical tuples pass through unchanged
// and lets TestArchAssetsMatchDetectedArch assert every tollgateArchAssets key
// is a selectable tuple. Add a tuple here only when it is a real OpenWrt feed
// target; an unknown bare name maps to "" (never a guessed default).
var canonicalArchTuples = map[string]bool{
	"aarch64_cortex-a53": true,
	"mipsel_24kc":        true,
	"mips_24kc":          true,
	"x86_64":             true,
}

// normalizeBareArch maps a BARE architecture string to a canonical OpenWrt
// tuple, or returns "" if it cannot be mapped. It tolerates surrounding
// whitespace and case; canonical tuples (e.g. "aarch64_cortex-a53") pass
// through unchanged. It NEVER invents a default.
func normalizeBareArch(b string) string {
	b = strings.ToLower(strings.TrimSpace(b))
	if b == "" {
		return ""
	}
	if canonicalArchTuples[b] {
		return b
	}
	if tuple, ok := bareArchToTuple[b]; ok {
		return tuple
	}
	return ""
}

// tollgateArchAssets is the per-arch map of tollgate-wrt download assets, keyed
// by the canonical OpenWrt arch tuple. IPK is used on OpenWrt <=24.x (opkg), APK
// on OpenWrt 25+ (apk).
//
// Only aarch64_cortex-a53 has published release assets today; the other entries
// are placeholders recording the naming convention so a future release (e.g.
// mipsel_24kc binaries for MT3000-class routers) can be wired by adding the URL
// here AND in arch_test.go (TestArchAssetsMatchDetectedArch) without touching
// deploy.go's selection logic. selectPkgURL refuses to return a URL for an arch
// that has no published asset rather than silently substituting aarch64.
var tollgateArchAssets = map[string]struct{ IPK, APK string }{
	// Fallback source — the CURRENT GitHub release URLs, kept as-is so existing
	// deploys keep working. Both formats published on the aarch64 tuple.
	"aarch64_cortex-a53": {
		IPK: "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.7.0-alpha10/tollgate-wrt_v0.7.0-alpha10_aarch64_cortex-a53.ipk",
		APK: "https://github.com/felixfelix-bot/tollgate-module-basic-go/releases/download/v0.6.1-post-merge/tollgate-wrt_main.56.b528e1d_aarch64_cortex-a53.apk",
	},
	// No published assets yet — selectPkgURL returns ok=false for these.
	"mipsel_24kc": {},
	"mips_24kc":   {},
	"x86_64":      {},
}

// selectPkgURL returns the tollgate-wrt download URL and file extension for the
// given canonical arch tuple and package manager. pkgMgr must be "opkg" or
// "apk".
//
// It returns ok=false when the arch is unknown OR that arch has no published
// asset in the requested format. The caller must treat ok=false as a hard
// deploy failure — NEVER silently substitute the aarch64 fallback.
func selectPkgURL(arch, pkgMgr string) (url, ext string, ok bool) {
	asset, exists := tollgateArchAssets[arch]
	if !exists {
		return "", "", false
	}
	switch pkgMgr {
	case "apk":
		if asset.APK == "" {
			return "", "", false
		}
		return asset.APK, ".apk", true
	case "opkg":
		if asset.IPK == "" {
			return "", "", false
		}
		return asset.IPK, ".ipk", true
	default:
		return "", "", false
	}
}

// distArchRe pulls DISTRIB_ARCH out of /etc/openwrt_release output. The file
// format is shell: DISTRIB_ARCH="aarch64_cortex-a53".
var distArchRe = regexp.MustCompile(`(?m)^DISTRIB_ARCH=["']?([^"'\s]+)`)

// opkgArchRe matches the first "arch <name> <pri>" line in `opkg print-architecture`.
var opkgArchRe = regexp.MustCompile(`(?m)^\s*arch\s+(\S+)\s+\d+`)

// detectArch returns the canonical OpenWrt CPU architecture tuple for the
// target router, probing sources in precedence order. It returns "" (never a
// hardcoded default) if the architecture cannot be determined.
//
// Precedence:
//  1. DISTRIB_ARCH from /etc/openwrt_release  — authoritative, exact-tuple today.
//  2. `opkg print-architecture`               — first "arch <name> <pri>" line.
//  3. `ubus call system board`                — JSON .architecture.
//  4. `apk --print-arch`                      — bare; normalized via normalizeBareArch.
//  5. `uname -m`                              — coarse; last resort, always normalized.
func detectArch(client *ssh.Client) string {
	if client == nil {
		return ""
	}
	return detectArchFrom(func(cmd string) string {
		return sshRun(client, cmd)
	})
}

// detectArchFrom is the pure, testable core of detectArch. get runs a remote
// command and returns its combined output; the precedence ladder below is
// exercised directly by TestDetectArchPrecedence without a live SSH session.
func detectArchFrom(get func(cmd string) string) string {
	if get == nil {
		return ""
	}

	// 1. Authoritative: DISTRIB_ARCH from /etc/openwrt_release.
	if out := get("grep DISTRIB_ARCH /etc/openwrt_release 2>/dev/null; true"); out != "" {
		if m := distArchRe.FindStringSubmatch(out); m != nil {
			arch := strings.TrimSpace(m[1])
			if arch != "" {
				return arch
			}
		}
	}

	// 2. opkg print-architecture — first real "arch <name> <pri>" line.
	// The command lists pseudo-architectures first ("all", "noarch"); the
	// usable CPU tuple is the first line whose name isn't a pseudo-arch.
	if out := get("opkg print-architecture 2>/dev/null; true"); out != "" {
		for _, m := range opkgArchRe.FindAllStringSubmatch(out, -1) {
			arch := strings.TrimSpace(m[1])
			if arch == "" || arch == "all" || arch == "noarch" {
				continue
			}
			return arch
		}
	}

	// 3. ubus call system board — JSON .architecture field.
	if out := get("ubus call system board 2>/dev/null; true"); out != "" {
		var board struct {
			Architecture string `json:"architecture"`
		}
		if err := json.Unmarshal([]byte(out), &board); err == nil && board.Architecture != "" {
			return board.Architecture
		}
	}

	// 4. apk --print-arch — bare arch, must be normalized.
	if out := strings.TrimSpace(get("apk --print-arch 2>/dev/null; true")); out != "" {
		if arch := normalizeBareArch(out); arch != "" {
			return arch
		}
	}

	// 5. uname -m — coarse, last resort, always normalized.
	if out := strings.TrimSpace(get("uname -m 2>/dev/null; true")); out != "" {
		if arch := normalizeBareArch(out); arch != "" {
			return arch
		}
	}

	return ""
}

// downloadBaseURL builds the OpenWrt package-repository base URL for an arch
// tuple. The wizard pins release 24.10.4 (the release it was validated on); the
// arch tuple is threaded in so nodogsplash/jq installs resolve for ANY router.
func downloadBaseURL(arch string) string {
	return fmt.Sprintf("https://downloads.openwrt.org/releases/24.10.4/packages/%s/", arch)
}
