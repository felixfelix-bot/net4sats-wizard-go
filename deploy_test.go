package main

import (
	"strings"
	"testing"
)

// TestDownloadURLsPointToTollgateRepo verifies that the tollgate package
// download URLs used by the wizard target the tollgate package project.
//
// Phase 3 (feat/feed-per-arch-urls): the PRIMARY source for every known arch is
// now the FreedomTechFeed/packages release assets (the per-arch feed). The
// GitHub tollgate-module-basic-go release URLs moved into tollgateGithubFallback
// and remain the FALLBACK for arches the feed does not publish yet (or a feed
// outage).
//
// This test checks that:
//   - every feed (primary) URL points at FreedomTechFeed/packages and is a
//     GitHub release download URL;
//   - every GitHub fallback URL points at the tollgate-module-basic-go repo
//     (upstream OpenTollGate or felixfelix-bot fork) as a release download.
//
// This guards against accidentally reverting to a stale or wrong repo, or
// pointing the fallback at a non-tollgate repo.
//
// (SW4a) The nft-enforce overlay URL is gone — those rules ship inside the
// ipk — and the full set of pinned URL constants, plus a live HTTP 200 check
// for each, is enforced in pins_test.go.
//
// (feat/auto-detect-arch) The per-arch download URLs live in the
// tollgateArchAssets map in arch.go, keyed by OpenWrt tuple.
func TestDownloadURLsPointToTollgateRepo(t *testing.T) {
	for arch, asset := range tollgateArchAssets {
		for name, url := range map[string]string{"IPK": asset.IPK, "APK": asset.APK} {
			if url == "" {
				continue // no published asset yet — not checked
			}
			t.Run(arch+"_"+name, func(t *testing.T) {
				// 1. Must be a GitHub URL.
				if !strings.HasPrefix(url, "https://github.com/") {
					t.Errorf("[%s] %s = %q: must be an https://github.com URL", arch, name, url)
				}
				// 2. Primary (feed) URLs must reference the feed repo.
				if !strings.Contains(url, "FreedomTechFeed/packages") {
					t.Errorf("[%s] %s = %q: primary source must be FreedomTechFeed/packages", arch, name, url)
				}
				// 3. Must be a release download URL, not e.g. a branch archive.
				if !strings.Contains(url, "/releases/download/") {
					t.Errorf("[%s] %s = %q: must be a GitHub release download URL", arch, name, url)
				}
			})
		}
	}

	// The GitHub fallback must point at the tollgate-module-basic-go repo
	// (upstream or fork) as a release download.
	for arch, asset := range tollgateGithubFallback {
		for name, url := range map[string]string{"IPK": asset.IPK, "APK": asset.APK} {
			if url == "" {
				continue
			}
			t.Run("fallback_"+arch+"_"+name, func(t *testing.T) {
				if !strings.HasPrefix(url, "https://github.com/") {
					t.Errorf("[%s] %s = %q: must be an https://github.com URL", arch, name, url)
				}
				ownerOk := strings.Contains(url, "OpenTollGate/") ||
					strings.Contains(url, "felixfelix-bot/")
				repoOk := strings.Contains(url, "tollgate-module-basic-go")
				if !ownerOk {
					t.Errorf("[fallback] [%s] %s = %q: owner must be OpenTollGate (primary) or felixfelix-bot (fallback)", arch, name, url)
				}
				if !repoOk {
					t.Errorf("[fallback] [%s] %s = %q: must reference tollgate-module-basic-go", arch, name, url)
				}
				if !strings.Contains(url, "/releases/download/") {
					t.Errorf("[%s] %s = %q: must be a GitHub release download URL", arch, name, url)
				}
			})
		}
	}
}

// TestTollgatePkgURLAssetName verifies the aarch64 .ipk asset (the fallback
// source the wizard publishes) exists in the URL and references the
// tollgate-wrt binary. This catches typos introduced when bumping release
// tags/asset names. The architecture-correctness of asset selection is
// covered by TestArchAssetsMatchDetectedArch + TestSelectPkgURL in arch_test.go.
func TestTollgatePkgURLAssetName(t *testing.T) {
	ipk := tollgateArchAssets["aarch64_cortex-a53"].IPK
	// Must end with the .ipk extension.
	if !strings.HasSuffix(ipk, ".ipk") {
		t.Errorf("aarch64 .ipk asset = %q: must end with .ipk", ipk)
	}
	// Must target the arch tuple it is keyed under.
	if !strings.Contains(ipk, "aarch64_cortex-a53") {
		t.Errorf("aarch64 .ipk asset = %q: must target aarch64_cortex-a53", ipk)
	}
	// Must reference the tollgate-wrt binary, not some other asset.
	if !strings.Contains(ipk, "tollgate-wrt") {
		t.Errorf("aarch64 .ipk asset = %q: must reference the tollgate-wrt binary", ipk)
	}
}
