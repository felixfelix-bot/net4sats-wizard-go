package main

import (
	"strings"
	"testing"
)

// The uhttpd :2051 portal docroot and the /www mirror are two deploy targets
// with different content rules. This file pins the behaviour added with the
// :2051 docroot fix:
//
//   - the :2051 docroot receives the FULL bundle. It is the only directory
//     :2051 reads (uhttpd.portal.home is set to it, index_page is
//     splash.html), so a bundle copied anywhere else never reaches splash
//     clients — they keep seeing a stale portal.
//   - the /www mirror receives the same bundle MINUS index.html, because
//     /www/index.html is OpenWrt's LuCI redirect page (served by the luci
//     uhttpd instance on :8080) and must not be replaced by the portal SPA.
//
// Two regressions this guards against, both invisible to `go build`:
// dropping the exclusion silently replaces LuCI's :8080 landing page, and
// dropping the :2051 target silently serves a stale portal again.

// wantPortalDocroot is the uhttpd :2051 document root the deploy must write.
const wantPortalDocroot = "/etc/tollgate/net4sats-captive-portal-site"

// portalPlanPaths returns the set of remote paths a deploy plan touches.
func portalPlanPaths(plan []fsDeployEntry) map[string]bool {
	out := make(map[string]bool, len(plan))
	for _, e := range plan {
		out[e.RemotePath] = true
	}
	return out
}

func TestPortalDeployTargets(t *testing.T) {
	if portalDeployDir != wantPortalDocroot {
		t.Fatalf("portalDeployDir = %q, want %q (the uhttpd :2051 docroot)", portalDeployDir, wantPortalDocroot)
	}

	t.Run("2051 docroot gets the full bundle", func(t *testing.T) {
		plan, err := planFSDeploy(portalFS, "portal", portalDeployDir)
		if err != nil {
			t.Fatalf("planFSDeploy: %v", err)
		}
		if len(plan) == 0 {
			t.Fatal("plan is empty — the embedded portal bundle did not walk")
		}
		got := portalPlanPaths(plan)
		for _, want := range []string{"index.html", "splash.html", "sw.js", "manifest.json"} {
			if !got[portalDeployDir+"/"+want] {
				t.Errorf(":2051 docroot plan is missing %s — uhttpd index_page is splash.html and the SPA needs its own assets", want)
			}
		}
	})

	t.Run("/www mirror keeps LuCI index.html but still carries the portal", func(t *testing.T) {
		plan, err := planFSDeploy(portalFS, "portal", "/www", portalMirrorExclude...)
		if err != nil {
			t.Fatalf("planFSDeploy: %v", err)
		}
		got := portalPlanPaths(plan)
		if got["/www/index.html"] {
			t.Error("/www/index.html is written — that file is OpenWrt's LuCI redirect page (uhttpd :8080); the portal SPA must not clobber it")
		}
		for _, want := range []string{"splash.html", "sw.js", "manifest.json"} {
			if !got["/www/"+want] {
				t.Errorf("/www mirror is missing %s — the mirror must stay in sync with the :2051 docroot", want)
			}
		}
	})

	t.Run("deploy path requests both targets with the exclusion", func(t *testing.T) {
		// Source-level binding, the same idiom as pins_test.go and the S-4
		// tests: the behavioural plans above only matter if deploy step 6
		// actually asks for them. deployGoSrc is the exact source that
		// compiles into this package.
		for _, want := range []string{
			`sshDeployFS(client, portalFS, "portal", portalDeployDir)`,
			`sshDeployFS(client, portalFS, "portal", "/www", portalMirrorExclude...)`,
			`"uci set uhttpd.portal.home='" + portalDeployDir + "'"`,
			`"uci set uhttpd.portal.index_page='splash.html'"`,
		} {
			if !strings.Contains(deployGoSrc, want) {
				t.Errorf("deploy.go no longer contains %q — the portal deploy path drifted from the tested targets", want)
			}
		}
	})
}
