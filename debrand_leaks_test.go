package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoInternalBrandLeaks enforces the de-branded, donation-pitch state of
// the user-facing wizard tree: no internal names, personal device references,
// internal board paths, or internal-process labels may appear in the
// non-vendored sources. Vendored third-party assets (portal/, admin/,
// release/, rpcd/, uhttpd_net4sats/) are excluded: they contain unrelated
// library internals (e.g. the noble-curves "endomorphism" API) that are not
// brand leaks.
//
// Each token below leaked into this public repo during WZ1.9 review and was
// removed — this test pins the removal so the leaks cannot silently return.
func TestNoInternalBrandLeaks(t *testing.T) {
	// Vendored/embedded third-party trees: excluded by directory prefix.
	vendoredDirs := []string{
		"portal/", "admin/", "release/", "rpcd/", "uhttpd_net4sats/",
	}

	// token -> why it must never appear in the wizard's own sources.
	//
	// "FELIX" is matched case-SENSITIVELY: lowercase "felixfelix-bot" is the
	// legitimate public fork org used in release-asset URLs and must not
	// trigger. The leak form is the uppercase internal gate label.
	tokens := map[string]string{
		"FELIX":             "internal individual/gate label",
		"T14Gen5":           "personal laptop model",
		"10.47.41.1":        "personal subnet probe target",
		"kanban.db":         "internal kanban board database",
		"kanban/boards":     "internal kanban board path",
		"net4sats-mvp-v2":   "internal kanban board name",
		".hermes":           "internal agent-tooling path",
		"c03rad0r":          "personal host/user name",
		"dq05":              "internal host name",
		"tailscale":         "internal network tooling",
		"endo-runbook":      "internal codename in doc paths",
		"github.io/pull/":   "internal PR links used as user documentation",
		"issuecomment-":     "internal PR-comment deep links",
		"#issuecomment":     "internal PR-comment deep links",
	}
	caseSensitive := map[string]bool{"FELIX": true}

	// Walk the working tree (source files, tests, docs, HTML, hooks).
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "dist":
				return filepath.SkipDir
			}
			// Skip vendored directories entirely.
			for _, vd := range vendoredDirs {
				if strings.HasPrefix(path+"/", vd) {
					return filepath.SkipDir
			}
			}
			return nil
		}
		// Only scan text files that can carry user-facing copy.
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go", ".md", ".html", ".css", ".js", ".json", ".sh", ".yaml", ".yml", ".mod", ".sum":
		default:
			// Extension-less files: only small ones can be hook scripts;
			// anything large is a build artifact (e.g. the root wizard
			// binary) and must not be read into memory.
			if strings.Contains(filepath.Base(path), ".") {
				return nil
			}
			if fi, statErr := d.Info(); statErr == nil && fi.Size() > 64*1024 {
				return nil
			}
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// This test file itself names every token in its map, so skip it.
		if path == "debrand_leaks_test.go" {
			return nil
		}
		lower := strings.ToLower(string(data))
		for token, why := range tokens {
			needle := strings.ToLower(token)
			if caseSensitive[token] {
				// Case-sensitive match: the exact internal-label spelling.
				if strings.Contains(string(data), token) {
					t.Errorf("%s: contains %q (%s)", path, token, why)
				}
				continue
			}
			if strings.Contains(lower, needle) {
				t.Errorf("%s: contains %q (%s)", path, token, why)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}