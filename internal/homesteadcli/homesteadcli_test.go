package homesteadcli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
)

func TestRunShellInitBashFish(t *testing.T) {
	var stderr bytes.Buffer
	if c := RunShellInit([]string{"bash"}, &stderr); c != 0 {
		t.Fatalf("bash exit %d: %s", c, stderr.String())
	}
	if c := RunShellInit([]string{"fish"}, &stderr); c != 0 {
		t.Fatalf("fish exit %d: %s", c, stderr.String())
	}
	if c := RunShellInit([]string{}, &stderr); c != 2 {
		t.Fatalf("want 2, got %d", c)
	}
	if c := RunShellInit([]string{"csh"}, &stderr); c != 2 {
		t.Fatalf("want 2, got %d", c)
	}
}

func TestRunExportProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	prefsPath := filepath.Join(dir, "homestead", "preferences.yaml")
	if err := os.MkdirAll(filepath.Dir(prefsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := preferences.Save(prefsPath, preferences.DefaultPreferences()); err != nil {
		t.Fatal(err)
	}

	profPath := filepath.Join(dir, "homestead", "profile.yaml")
	st := profilestate.State{}
	profilestate.RecordInstalled(&st, "demo-pkg")
	if err := profilestate.Save(profPath, st); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	code := RunExportProfile([]string{"-format", "json"}, "unit", &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	var raw map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatalf("%s: %v", stdout.String(), err)
	}
}
