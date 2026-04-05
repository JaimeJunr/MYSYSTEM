package profilestate

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
)

func TestWriteExportJSON(t *testing.T) {
	var buf bytes.Buffer
	s := State{}
	RecordInstalled(&s, "alpha")
	ToggleFavorite(&s, "cleanup-1")
	p := preferences.DefaultPreferences()
	p.CatalogURL = "https://example.com/c.json"
	if err := WriteExport(&buf, "json", s, p, "testver"); err != nil {
		t.Fatal(err)
	}
	var doc ExportDoc
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("json: %v\n%s", err, buf.String())
	}
	if doc.Version != 1 {
		t.Fatalf("version %d", doc.Version)
	}
	if doc.HomesteadVersion != "testver" {
		t.Fatalf("ver %q", doc.HomesteadVersion)
	}
	if len(doc.InstalledPackages) != 1 || doc.InstalledPackages[0] != "alpha" {
		t.Fatalf("packages %+v", doc.InstalledPackages)
	}
	if len(doc.FavoriteScripts) != 1 || doc.FavoriteScripts[0] != "cleanup-1" {
		t.Fatalf("favorites %+v", doc.FavoriteScripts)
	}
	if doc.Preferences.CatalogURL != p.CatalogURL {
		t.Fatalf("prefs %+v", doc.Preferences)
	}
}

func TestWriteExportText(t *testing.T) {
	var buf bytes.Buffer
	s := State{}
	RecordInstalled(&s, "x")
	p := preferences.DefaultPreferences()
	if err := WriteExport(&buf, "text", s, p, "v"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "x") || !strings.Contains(out, "Homestead") {
		t.Fatalf("%q", out)
	}
}

func TestWriteExportBadFormat(t *testing.T) {
	err := WriteExport(&bytes.Buffer{}, "xml", State{}, preferences.DefaultPreferences(), "")
	if err == nil {
		t.Fatal("expected error")
	}
}
