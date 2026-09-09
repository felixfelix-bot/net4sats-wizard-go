package main

import (
	"strings"
	"testing"
)

// TestNormalizeBareArch verifies normalizeBareArch maps bare CPU arch strings
// (as reported by `apk --print-arch` or `uname -m`) to canonical OpenWrt arch
// tuples, and passes canonical tuples through unchanged. A bare "aarch64" must
// become "aarch64_cortex-a53" — this is the exact bug the wizard previously
// papered over by hardcoding aarch64_cortex-a53 in every download URL.
func TestNormalizeBareArch(t *testing.T) {
	cases := map[string]string{
		// bare names from apk --print-arch / uname -m
		"aarch64": "aarch64_cortex-a53",
		"mipsel":  "mipsel_24kc",
		"mips":    "mips_24kc",
		"x86_64":  "x86_64",
		"amd64":   "x86_64",
		// whitespace + case are tolerated
		"  AARCH64  ": "aarch64_cortex-a53",
		"Mipsel":      "mipsel_24kc",
		// canonical tuples pass through unchanged
		"aarch64_cortex-a53": "aarch64_cortex-a53",
		"mipsel_24kc":        "mipsel_24kc",
		"mips_24kc":          "mips_24kc",
		// unknowns normalize to ""
		"":        "",
		"sparc":   "",
		"armv7l":  "",
		"riscv64": "",
	}

	for in, want := range cases {
		if got := normalizeBareArch(in); got != want {
			t.Errorf("normalizeBareArch(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestSelectPkgURL verifies selectPkgURL returns the correct download URL and
// extension for the requested package manager, and FAILS (ok=false) for an
// unknown arch or an arch without a published asset in the requested format.
func TestSelectPkgURL(t *testing.T) {
	// aarch64_cortex-a53 carries the only currently-published fallback assets.
	if url, ext, ok := selectPkgURL("aarch64_cortex-a53", "opkg"); !ok {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, opkg): ok=%v, want true", ok)
	} else if ext != ".ipk" {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, opkg): ext=%q, want .ipk", ext)
	} else if url != tollgateArchAssets["aarch64_cortex-a53"].IPK {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, opkg): url=%q, want %q", url, tollgateArchAssets["aarch64_cortex-a53"].IPK)
	}

	if url, ext, ok := selectPkgURL("aarch64_cortex-a53", "apk"); !ok {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, apk): ok=%v, want true", ok)
	} else if ext != ".apk" {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, apk): ext=%q, want .apk", ext)
	} else if url != tollgateArchAssets["aarch64_cortex-a53"].APK {
		t.Errorf("selectPkgURL(aarch64_cortex-a53, apk): url=%q, want %q", url, tollgateArchAssets["aarch64_cortex-a53"].APK)
	}

	// Arch entries in the map but with no published asset must FAIL loudly
	// (never silently fall back to aarch64_cortex-a53).
	for _, arch := range []string{"mipsel_24kc", "mips_24kc", "x86_64"} {
		if _, _, ok := selectPkgURL(arch, "opkg"); ok {
			t.Errorf("selectPkgURL(%q, opkg) = ok=true, want false (no published asset; must not default)", arch)
		}
		if _, _, ok := selectPkgURL(arch, "apk"); ok {
			t.Errorf("selectPkgURL(%q, apk) = ok=true, want false (no published asset; must not default)", arch)
		}
	}

	// Unknown archs must fail, including the empty string.
	for _, arch := range []string{"sparc", "riscv64", ""} {
		if _, _, ok := selectPkgURL(arch, "opkg"); ok {
			t.Errorf("selectPkgURL(%q, opkg) = ok=true, want false (unknown arch)", arch)
		}
	}
}

// TestArchAssetsMatchDetectedArch ties arch detection to asset selection:
// every key in tollgateArchAssets must be a canonical tuple (so a value
// returned by detectArch's normalizeBareArch is always selectable), and the
// aarch64_cortex-a53 entry — the fallback source — must carry BOTH the .ipk
// and .apk GitHub release URLs for tollgate-wrt.
func TestArchAssetsMatchDetectedArch(t *testing.T) {
	// Every map key must be a canonical tuple so selectPkgURL can look up any
	// arch that detectArch might return.
	for tuple := range tollgateArchAssets {
		if normalizeBareArch(tuple) != tuple {
			t.Errorf("tollgateArchAssets key %q is not a canonical tuple (normalizeBareArch(%q)=%q)", tuple, tuple, normalizeBareArch(tuple))
		}
	}

	ipk := tollgateArchAssets["aarch64_cortex-a53"].IPK
	apk := tollgateArchAssets["aarch64_cortex-a53"].APK

	if ipk == "" || !strings.HasPrefix(ipk, "https://github.com/") || !strings.HasSuffix(ipk, ".ipk") {
		t.Errorf("aarch64_cortex-a53 .ipk fallback missing or malformed: %q", ipk)
	}
	if !strings.Contains(ipk, "tollgate-wrt") {
		t.Errorf("aarch64_cortex-a53 .ipk must reference tollgate-wrt: %q", ipk)
	}
	if apk == "" || !strings.HasPrefix(apk, "https://github.com/") || !strings.HasSuffix(apk, ".apk") {
		t.Errorf("aarch64_cortex-a53 .apk fallback missing or malformed: %q", apk)
	}
	if !strings.Contains(apk, "tollgate-wrt") {
		t.Errorf("aarch64_cortex-a53 .apk must reference tollgate-wrt: %q", apk)
	}
}
