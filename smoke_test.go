package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSmoke_VersionFlag builds the CLI and checks -version exits successfully (black-box smoke).
func TestSmoke_VersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("smoke test skipped with go test -short")
	}

	bin := filepath.Join(t.TempDir(), "homestead")
	build := exec.Command("go", "build", "-o", bin, "./cmd/homestead")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	run := exec.Command(bin, "-version")
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("homestead -version: %v\n%s", err, out)
	}
	if strings.TrimSpace(string(out)) == "" {
		t.Fatal("expected non-empty version output")
	}
}
