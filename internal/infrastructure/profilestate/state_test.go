package profilestate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordInstalledDedupes(t *testing.T) {
	var s State
	RecordInstalled(&s, "a")
	RecordInstalled(&s, "b")
	RecordInstalled(&s, "a")
	if len(s.InstalledPackageIDs) != 2 {
		t.Fatalf("got %v", s.InstalledPackageIDs)
	}
}

func TestToggleFavorite(t *testing.T) {
	var s State
	if !ToggleFavorite(&s, "x") {
		t.Fatal("expected favorite")
	}
	if !IsFavorite(&s, "x") {
		t.Fatal("expected IsFavorite")
	}
	if ToggleFavorite(&s, "x") {
		t.Fatal("expected unfavorite")
	}
	if IsFavorite(&s, "x") {
		t.Fatal("expected not favorite")
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.yaml")
	s := State{}
	RecordInstalled(&s, "pkg-a")
	ToggleFavorite(&s, "script-1")
	if err := Save(path, s); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.InstalledPackageIDs) != 1 || got.InstalledPackageIDs[0] != "pkg-a" {
		t.Fatalf("packages %+v", got.InstalledPackageIDs)
	}
	if len(got.FavoriteScriptIDs) != 1 || got.FavoriteScriptIDs[0] != "script-1" {
		t.Fatalf("favorites %+v", got.FavoriteScriptIDs)
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "profile.yaml")
	if err := Save(path, State{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
