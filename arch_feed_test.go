package main

import (
	"strings"
	"testing"
)

// TestSelectPkgURLFallbackPath verifies the GitHub-release fallback branch of
// selectPkgURL: when the feed (tollgateArchAssets) has no asset for an
// arch/format, selectPkgURL must fall back to tollgateGithubFallback instead
// of returning ok=false.
//
// The real maps cover this implicitly (every known arch is in the feed), so
// the test temporarily empties the aarch64 feed entries to force the fallback
// path — the same branch that fires during a feed outage or before Phase 2
// publishes a new arch.
func TestSelectPkgURLFallbackPath(t *testing.T) {
	// Save + restore the aarch64 feed entries so we don't poison the live
	// maps for sibling tests (these do not run in parallel).
	saved := tollgateArchAssets["aarch64_cortex-a53"]
	defer func() {
		tollgateArchAssets["aarch64_cortex-a53"] = saved
	}()

	// Force the fallback: empty both feed entries for aarch64.
	tollgateArchAssets["aarch64_cortex-a53"] = struct{ IPK, APK string }{}

	// opkg path falls back to the GitHub .ipk.
	url, ext, ok := selectPkgURL("aarch64_cortex-a53", "opkg")
	if !ok {
		t.Fatalf("selectPkgURL(aarch64, opkg) w/o feed = ok=false, want GitHub fallback")
	}
	if ext != ".ipk" {
		t.Errorf("fallback opkg ext = %q, want .ipk", ext)
	}
	if url != tollgateGithubFallback["aarch64_cortex-a53"].IPK {
		t.Errorf("fallback opkg url = %q, want GitHub .ipk %q", url, tollgateGithubFallback["aarch64_cortex-a53"].IPK)
	}

	// apk path falls back to the GitHub .apk.
	url, ext, ok = selectPkgURL("aarch64_cortex-a53", "apk")
	if !ok {
		t.Fatalf("selectPkgURL(aarch64, apk) w/o feed = ok=false, want GitHub fallback")
	}
	if ext != ".apk" {
		t.Errorf("fallback apk ext = %q, want .apk", ext)
	}
	if url != tollgateGithubFallback["aarch64_cortex-a53"].APK {
		t.Errorf("fallback apk url = %q, want GitHub .apk %q", url, tollgateGithubFallback["aarch64_cortex-a53"].APK)
	}

	// Unknown pkgMgr with the feed-empty arch must still fail (no fallback
	// branch matches an unrecognized package manager).
	if _, _, ok := selectPkgURL("aarch64_cortex-a53", "rpm"); ok {
		t.Errorf("selectPkgURL(aarch64, rpm) = ok=true, want false (unknown pkgMgr)")
	}
}

// TestSelectPkgURLFeedPrimary is the RED→GREEN contract for Phase 3: the
// wizard's arch→URL selection must point at the FEED's per-arch stable URLs
// (FreedomTechFeed/packages release assets) for every known arch, NOT the
// GitHub tollgate-module-basic-go release URLs. The feed publishes all four
// canonical tuples (aarch64_cortex-a53, mipsel_24kc, mips_24kc, x86_64) in
// both .apk and .ipk, so selectPkgURL must return a feed URL for each.
//
// The GitHub-release URLs remain as a FALLBACK for arches the feed does not
// publish yet — selectPkgURL must consult them only when the feed has no
// asset for the requested arch/format.
func TestSelectPkgURLFeedPrimary(t *testing.T) {
	// Every known arch must resolve to a FEED URL (FreedomTechFeed/packages),
	// not the GitHub tollgate-module-basic-go fallback.
	for _, arch := range []string{"aarch64_cortex-a53", "mipsel_24kc", "mips_24kc", "x86_64"} {
		for _, pkgMgr := range []string{"apk", "opkg"} {
			url, ext, ok := selectPkgURL(arch, pkgMgr)
			if !ok {
				t.Errorf("selectPkgURL(%q, %q) = ok=false, want true (feed publishes this arch)", arch, pkgMgr)
				continue
			}
			if !strings.Contains(url, "FreedomTechFeed/packages") {
				t.Errorf("selectPkgURL(%q, %q) = %q, want a FEED URL (FreedomTechFeed/packages)", arch, pkgMgr, url)
			}
			wantExt := ".apk"
			if pkgMgr == "opkg" {
				wantExt = ".ipk"
			}
			if ext != wantExt {
				t.Errorf("selectPkgURL(%q, %q) ext = %q, want %q", arch, pkgMgr, ext, wantExt)
			}
		}
	}
}

// TestSelectPkgURLFeedFallback verifies the GitHub-release URLs are preserved
// as a fallback for arches the feed does not publish yet. The fallback map
// (tollgateGithubFallback) must still carry the aarch64_cortex-a53 GitHub
// assets so a feed outage or a not-yet-published arch can fall back.
func TestSelectPkgURLFeedFallback(t *testing.T) {
	// The fallback map must exist and carry the aarch64 GitHub assets.
	if len(tollgateGithubFallback) == 0 {
		t.Fatal("tollgateGithubFallback is empty — GitHub-release fallback must be preserved")
	}
	gh := tollgateGithubFallback["aarch64_cortex-a53"]
	if gh.IPK == "" || !strings.HasPrefix(gh.IPK, "https://github.com/") || !strings.HasSuffix(gh.IPK, ".ipk") {
		t.Errorf("aarch64_cortex-a53 GitHub .ipk fallback missing/malformed: %q", gh.IPK)
	}
	if gh.APK == "" || !strings.HasPrefix(gh.APK, "https://github.com/") || !strings.HasSuffix(gh.APK, ".apk") {
		t.Errorf("aarch64_cortex-a53 GitHub .apk fallback missing/malformed: %q", gh.APK)
	}
	// The fallback URLs must NOT be the feed URLs (they are distinct sources).
	if strings.Contains(gh.IPK, "FreedomTechFeed") {
		t.Errorf("GitHub fallback .ipk must not point at the feed: %q", gh.IPK)
	}
}
