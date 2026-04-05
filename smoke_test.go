package main

import (
	"os"
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

	runExport := exec.Command(bin, "export-profile", "-format", "json")
	runExport.Env = append(os.Environ(), "XDG_CONFIG_HOME="+t.TempDir())
	outEx, err := runExport.CombinedOutput()
	if err != nil {
		t.Fatalf("homestead export-profile: %v\n%s", err, outEx)
	}
	if !strings.Contains(string(outEx), `"version"`) {
		t.Fatalf("expected JSON export, got: %s", outEx)
	}

	runFish := exec.Command(bin, "shell-init", "fish")
	outFish, err := runFish.CombinedOutput()
	if err != nil {
		t.Fatalf("homestead shell-init fish: %v\n%s", err, outFish)
	}
	if !strings.Contains(string(outFish), "HOMESTEAD_CONFIG_DIR") {
		t.Fatalf("expected HOMESTEAD_CONFIG_DIR in fish snippet: %s", outFish)
	}
}
