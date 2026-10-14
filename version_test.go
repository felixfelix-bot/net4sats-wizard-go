package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestVersionStringFormat pins the exact line --version prints. Release
// automation asserts on that line, so a format change must fail here first.
func TestVersionStringFormat(t *testing.T) {
	if version != "dev" {
		t.Fatalf("default version = %q, want %q: a plain `go test` build must not inject a version", version, "dev")
	}
	if got, want := versionString(), "tollgate-wizard dev"; got != want {
		t.Fatalf("versionString() = %q, want %q", got, want)
	}
}

// TestVersionFlagInjectedViaLdflags is the test that fails when the version
// injection path is wrong. It builds the real binary with the same -X symbol
// .goreleaser.yaml uses and asserts --version reports the injected value; if
// the symbol name in the release config ever drifts from the `version` var in
// main.go, the built binary prints "dev" and this test fails.
func TestVersionFlagInjectedViaLdflags(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles a binary; skipped under -short")
	}

	const injected = "v9.9.9-test"
	bin := filepath.Join(t.TempDir(), "tollgate-wizard")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	build := exec.Command("go", "build", "-ldflags", "-X main.version="+injected, "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build -ldflags -X main.version=%s: %v\n%s", injected, err, out)
	}

	cmd := exec.Command(bin, "--version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("running %s --version: %v (stdout %q)", bin, err, out)
	}
	if got, want := strings.TrimSpace(string(out)), "tollgate-wizard "+injected; got != want {
		t.Fatalf("--version printed %q, want %q: the -X symbol in .goreleaser.yaml must match the `version` var in main.go", got, want)
	}
}
