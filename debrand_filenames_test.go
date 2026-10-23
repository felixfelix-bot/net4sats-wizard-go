package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoInternalProcessNamedFiles asserts that no file published in the repo is
// named after an internal process or person. This is the gap that a filename
// rename exposes: a rename is two halves (add the new path, delete the old one)
// and it is entirely possible to commit only the first, leaving the retired
// filename visible in the public tree while every content scan passes.
//
// Found the hard way while closing the WZ1.9 audit: `docs/endo-runbook.md` was
// re-added alongside `docs/router-setup.md` in one commit with a scoped
// pathspec, so the leaky filename survived the "rename".
func TestNoInternalProcessNamedFiles(t *testing.T) {
	handover := "en" + "do"
	bannedSubstrings := []string{
		handover,      // internal handover process name
		"t14" + "gen", // personal laptop model
		"felix",       // operator first name used as a label
	}

	skipDirs := map[string]bool{
		".git": true, "admin": true, "portal": true, "verify": true, "artifacts": true,
	}

	var violations []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skipDirs[info.Name()] || info.Name() == "assets" {
				return filepath.SkipDir
			}
			return nil
		}
		// Task-history artifacts are not shipped copy: they legitimately
		// record which internal card produced them.
		base := filepath.Base(path)
		lower := strings.ToLower(base)
		if strings.HasPrefix(lower, "report-") || strings.HasPrefix(lower, "progress-") ||
			strings.HasPrefix(lower, "report.") || strings.HasPrefix(lower, "progress.") ||
			strings.HasPrefix(lower, "audit-") || strings.HasPrefix(lower, "evidence-") {
			return nil
		}
		for _, bad := range bannedSubstrings {
			if strings.Contains(lower, bad) {
				violations = append(violations, path)
				break
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	for _, v := range violations {
		t.Errorf("file name still carries internal detail: %s", v)
	}
}
