package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// internalLeakPatterns returns the patterns that must never appear on a shipped
// surface. These are not product branding (net4sats is the product and stays) —
// they are the operator's private context that has no business in a repository
// being donated: a personal laptop model, a first name used as a machine label,
// internal kanban card/branch identifiers, internal process names, and internal
// decision-gate ids.
//
// Needles are assembled from fragments so that a repo-wide grep for the retired
// wording stays clean: this file must not become the last place it survives.
func internalLeakPatterns() []struct {
	name string
	re   *regexp.Regexp
} {
	laptop := "T14" + "Gen" + "5"
	handover := "En" + "do"
	return []struct {
		name string
		re   *regexp.Regexp
	}{
		{"personal laptop model", regexp.MustCompile(`\b` + laptop + `\b`)},
		{"operator first name as machine label", regexp.MustCompile(`\bFelix'?s\b`)},
		{"internal decision-gate id", regexp.MustCompile(`\bFELIX\b`)},
		{"internal handover process name", regexp.MustCompile(`(?i)\b` + handover + `[- ](handover|onboarding|runbook)\b`)},
		{"internal card id", regexp.MustCompile(`\bt_[0-9a-f]{8}\b`)},
		{"internal branch name", regexp.MustCompile(`worker-(heavy|tollgate)/`)},
		{"personal fleet hostname", regexp.MustCompile(`\b(c08r4d0r|c03rad0r|CobradorWave)\b`)},
		{"internal board path", regexp.MustCompile(`(hermes/kanban|kanban\.db|workspaces/t_)`)},
		// 10.47.41.1 is intentionally NOT here: it is a probe candidate IP
		// (functional behaviour), not branding. See AUDIT-t_c92a4003.md.
	}
}

// shippedSurfaces lists the files that actually reach an operator: the embedded
// wizard UI, the packaged portal/admin trees, the Go sources, and the
// operator-facing docs. Task-history artifacts (REPORT-*/PROGRESS-*/EVIDENCE-*)
// and vendored minified bundles are deliberately out of scope.
//
// The test file itself is excluded from the scan: the needles are assembled
// from fragments precisely so the retired wording is absent, but this file's
// own explanatory prose and card references are not shipped copy, and scanning
// it would make the guard fail on its own documentation.
func shippedSurfaces(t *testing.T) []string {
	t.Helper()
	paths := []string{"index.html", "README.md", "CHANGELOG.md"}

	for _, dir := range []string{"portal", "admin", "rpcd"} {
		err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// Vendored minified bundles: third-party code, not our copy.
				if filepath.Base(p) == "assets" {
					return filepath.SkipDir
				}
				return nil
			}
			switch strings.ToLower(filepath.Ext(p)) {
			case ".html", ".json", ".js":
				paths = append(paths, p)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}

	entries, err := os.ReadDir("docs")
	if err != nil {
		t.Fatalf("read docs/: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			paths = append(paths, filepath.Join("docs", e.Name()))
		}
	}

	gos, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob *.go: %v", err)
	}
	for _, g := range gos {
		if strings.HasSuffix(g, "_test.go") {
			// Test sources are not shipped copy and legitimately cite the
			// card ids an audit is about; only non-test Go ships.
			continue
		}
		paths = append(paths, g)
	}

	return paths
}

// TestNoInternalProcessLeaksOnShippedSurfaces is the regression guard for the
// WZ1.9 audit's leak findings (AUDIT-t_c92a4003.md L1–L8). It is a black-box
// file scan on purpose: the guard must fail when the retired words reappear
// anywhere on a shipped surface, not merely when some helper changes.
func TestNoInternalProcessLeaksOnShippedSurfaces(t *testing.T) {
	patterns := internalLeakPatterns()

	for _, path := range shippedSurfaces(t) {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for i, line := range strings.Split(string(body), "\n") {
			for _, p := range patterns {
				if m := p.re.FindString(line); m != "" {
					t.Errorf("%s:%d: internal detail leaked (%s): %q",
						path, i+1, p.name, m)
				}
			}
		}
	}
}
