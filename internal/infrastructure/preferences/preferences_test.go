package preferences

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPreferences(t *testing.T) {
	p := DefaultPreferences()
	if !p.ConfirmBeforeScript || !p.ConfirmBeforePackage {
		t.Fatal("expected confirms true by default")
	}
	if p.Theme != ThemeDark {
		t.Fatalf("theme = %q", p.Theme)
	}
	if p.TextScale != TextScaleNormal {
		t.Fatalf("text_scale = %q", p.TextScale)
	}
}

func TestNormalizeTextScale(t *testing.T) {
	p := DefaultPreferences()
	p.TextScale = "bogus"
	p.Normalize()
	if p.TextScale != TextScaleNormal {
		t.Fatalf("got %q", p.TextScale)
	}
}

func TestLoadMissingFile(t *testing.T) {
	p, err := Load(filepath.Join(t.TempDir(), "nope.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if p.Theme != ThemeDark {
		t.Fatal("expected defaults")
	}
}

func TestLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "preferences.yaml")
	want := DefaultPreferences()
	want.CatalogURL = "https://example.com/catalog.json"
	want.ScriptRoot = ""
	want.DotfilesRepo = "~/foo/dotfiles"
	want.Theme = ThemeLight
	want.ConfirmBeforeScript = false
	want.ConfirmBeforePackage = true

	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.CatalogURL != want.CatalogURL || got.Theme != want.Theme {
		t.Fatalf("got %+v want %+v", got, want)
	}
	if got.ConfirmBeforeScript != false || got.ConfirmBeforePackage != true {
		t.Fatalf("confirm flags %+v", got)
	}
}

func TestFromRawNilBoolUsesDefaultTrue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "p.yaml")
	if err := os.WriteFile(path, []byte("theme: dark\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.ConfirmBeforeScript || !got.ConfirmBeforePackage {
		t.Fatal("missing keys should keep default true")
	}
}

func TestValidateCatalogURL(t *testing.T) {
	if err := ValidateCatalogURL(""); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCatalogURL("https://a/b"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCatalogURL("ftp://x"); err == nil {
		t.Fatal("expected error")
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip(err)
	}
	got, err := ExpandPath("~")
	if err != nil || got != filepath.Clean(home) {
		t.Fatalf("ExpandPath(~) = %q, %v", got, err)
	}
}
