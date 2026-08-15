package main

import (
	"strings"
	"testing"
)

// TestDownloadURLsPointToTollgateRepo verifies that the deployment download
// URLs used by the wizard target the tollgate-module-basic-go project
// (either the upstream OpenTollGate org or its felixfelix-bot fork).
//
// This guards against accidentally reverting to a stale or wrong repo, and
// documents the Endo-handover fallback strategy: upstream OpenTollGate is the
// primary source; the felixfelix-bot fork is the fallback for the
// v0.6.1-post-merge tag until an equivalent upstream release exists.
func TestDownloadURLsPointToTollgateRepo(t *testing.T) {
	cases := map[string]string{
		"tollgatePkgURL":        tollgatePkgURL,
		"tollgateNftEnforceURL": tollgateNftEnforceURL,
	}

	for name, url := range cases {
		t.Run(name, func(t *testing.T) {
			// 1. Must not be empty.
			if url == "" {
				t.Fatalf("%s is empty", name)
			}

			// 2. Must be a GitHub URL.
			if !strings.HasPrefix(url, "https://github.com/") {
				t.Errorf("%s = %q: must be an https://github.com URL", name, url)
			}

			// 3. Must reference the tollgate-module-basic-go repo, either via the
			//    upstream OpenTollGate org or the felixfelix-bot fork (fallback).
			ownerOk := strings.Contains(url, "OpenTollGate/") ||
				strings.Contains(url, "felixfelix-bot/")
			repoOk := strings.Contains(url, "tollgate-module-basic-go")
			if !ownerOk {
				t.Errorf("%s = %q: owner must be OpenTollGate (primary) or felixfelix-bot (fallback)", name, url)
			}
			if !repoOk {
				t.Errorf("%s = %q: must reference tollgate-module-basic-go", name, url)
			}

			// 4. Must be a release download URL, not e.g. a branch archive.
			if !strings.Contains(url, "/releases/download/") {
				t.Errorf("%s = %q: must be a GitHub release download URL", name, url)
			}
		})
	}
}

// TestTollgatePkgURLAssetName verifies the .ipk asset exists in the URL and
// targets the aarch64_cortex-a53 architecture the routers use. This catches
// typos introduced when bumping release tags/asset names (the v0.6.1-post-merge
// asset is named differently from the v0.5.0-e2e-test one).
func TestTollgatePkgURLAssetName(t *testing.T) {
	// Must end with the .ipk extension.
	if !strings.HasSuffix(tollgatePkgURL, ".ipk") {
		t.Errorf("tollgatePkgURL = %q: must end with .ipk", tollgatePkgURL)
	}
	// Must target the router's architecture.
	if !strings.Contains(tollgatePkgURL, "aarch64_cortex-a53") {
		t.Errorf("tollgatePkgURL = %q: must target aarch64_cortex-a53", tollgatePkgURL)
	}
	// Must reference the tollgate-wrt binary, not some other asset.
	if !strings.Contains(tollgatePkgURL, "tollgate-wrt") {
		t.Errorf("tollgatePkgURL = %q: must reference the tollgate-wrt binary", tollgatePkgURL)
	}
}

// TestConfigwizURLPointsToOrgRelease verifies the admin-panel tarball is
// pulled from the net4sats org's configurationwizzard releases, not from a
// personal fork. The org v1.0.1 release carries the fork v1.0.1 content plus
// the PR #22 balance-redirect fix and the admin-port (8090) fix.
func TestConfigwizURLPointsToOrgRelease(t *testing.T) {
	// 1. Must be a GitHub release download URL on the net4sats org repo.
	if !strings.HasPrefix(configwizURL, "https://github.com/net4sats/configurationwizzard/releases/download/") {
		t.Errorf("configwizURL = %q: must point at net4sats/configurationwizzard releases", configwizURL)
	}
	// 2. Must pull the v1.0.1 tarball asset.
	if !strings.HasSuffix(configwizURL, "/v1.0.1/net4sats-configwiz-1.0.1.tar.gz") {
		t.Errorf("configwizURL = %q: must reference the v1.0.1 net4sats-configwiz tarball", configwizURL)
	}
}
