package catalog

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

func TestParseManifest_EmptyPackages(t *testing.T) {
	raw := []byte(`{"schema_version":1,"packages":[]}`)
	pkgs, sv, err := ParseManifest(raw)
	if err != nil {
		t.Fatal(err)
	}
	if sv != 1 {
		t.Fatalf("schema version = %d", sv)
	}
	if len(pkgs) != 0 {
		t.Fatalf("len = %d", len(pkgs))
	}
}

func TestParseManifest_InvalidJSON(t *testing.T) {
	_, _, err := ParseManifest([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseManifest_SchemaTooNew(t *testing.T) {
	raw := []byte(`{"schema_version":99,"packages":[]}`)
	_, sv, err := ParseManifest(raw)
	if !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("err = %v want ErrUnsupportedSchema", err)
	}
	if sv != 99 {
		t.Fatalf("sv = %d", sv)
	}
}

func TestParseManifest_SchemaZero(t *testing.T) {
	raw := []byte(`{"schema_version":0,"packages":[]}`)
	_, _, err := ParseManifest(raw)
	if !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("err = %v", err)
	}
}

func TestParseManifest_UnknownCategoryBecomesOther(t *testing.T) {
	raw := []byte(`{
		"schema_version": 1,
		"packages": [
			{
				"id": "x-tool",
				"name": "X",
				"description": "d",
				"version": "1",
				"category": "not_a_real_category",
				"install_cmd": "true",
				"check_cmd": "true"
			}
		]
	}`)
	pkgs, _, err := ParseManifest(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("len = %d", len(pkgs))
	}
	if pkgs[0].Category != types.PackageCategoryOther {
		t.Fatalf("category = %q", pkgs[0].Category)
	}
}

func TestParseManifest_SkipsInvalidPackage(t *testing.T) {
	raw := []byte(`{
		"schema_version": 1,
		"packages": [
			{"id": "", "name": "bad", "category": "ide", "install_cmd": "x", "check_cmd": "x"},
			{
				"id": "good",
				"name": "Good",
				"description": "d",
				"version": "1",
				"category": "ide",
				"install_cmd": "true",
				"check_cmd": "true"
			}
		]
	}`)
	pkgs, _, err := ParseManifest(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 || pkgs[0].ID != "good" {
		t.Fatalf("pkgs = %#v", pkgs)
	}
}

func TestReadAndParseCacheFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "installer-catalog.json")
	raw := []byte(`{"schema_version":1,"packages":[{"id":"a","name":"A","description":"d","version":"1","category":"tool","install_cmd":"true","check_cmd":"true"}]}`)
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	pkgs, sv, err := ReadAndParseCacheFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if sv != 1 || len(pkgs) != 1 || pkgs[0].ID != "a" {
		t.Fatalf("sv=%d pkgs=%#v", sv, pkgs)
	}
}
