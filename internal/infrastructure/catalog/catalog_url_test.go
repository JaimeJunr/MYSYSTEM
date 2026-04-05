package catalog

import "testing"

func TestEffectiveCatalogURL(t *testing.T) {
	t.Setenv("HOMESTEAD_CATALOG_URL", "")
	got := EffectiveCatalogURL("")
	if got != DefaultCatalogURL {
		t.Fatalf("empty override: got %q", got)
	}
	got = EffectiveCatalogURL("https://example.com/m.json")
	if got != "https://example.com/m.json" {
		t.Fatalf("file override: got %q", got)
	}
	t.Setenv("HOMESTEAD_CATALOG_URL", "https://env.example/cat.json")
	got = EffectiveCatalogURL("https://file.example/ignored.json")
	if got != "https://env.example/cat.json" {
		t.Fatalf("env wins: got %q", got)
	}
}
