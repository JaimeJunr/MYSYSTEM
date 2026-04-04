package catalog

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

//go:embed installer-catalog.json
var embeddedInstallerCatalog []byte

// EmbeddedCatalogBytes is the installer manifest baked into the binary at build time.
func EmbeddedCatalogBytes() []byte {
	return embeddedInstallerCatalog
}

// SupportedSchemaVersion is the only manifest schema this client applies.
const SupportedSchemaVersion = 1

// ErrUnsupportedSchema means schema_version is not SupportedSchemaVersion.
var ErrUnsupportedSchema = errors.New("installer catalog schema_version not supported")

// IsSchemaSupported reports whether v can be applied by this client.
func IsSchemaSupported(v int) bool {
	return v == SupportedSchemaVersion
}

// DefaultCatalogURL is the raw GitHub URL for the versioned manifest on main.
const DefaultCatalogURL = "https://raw.githubusercontent.com/JaimeJunr/Homestead/main/internal/infrastructure/catalog/installer-catalog.json"

// ResolveCatalogURL returns HOMESTEAD_CATALOG_URL when set, otherwise DefaultCatalogURL.
func ResolveCatalogURL() string {
	if u := strings.TrimSpace(os.Getenv("HOMESTEAD_CATALOG_URL")); u != "" {
		return u
	}
	return DefaultCatalogURL
}

// CacheFilePath returns the path for the on-disk catalog cache.
func CacheFilePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(".cache", "homestead", "installer-catalog.json")
	}
	return filepath.Join(dir, "homestead", "installer-catalog.json")
}

// WriteCache writes raw manifest bytes after a successful remote fetch.
func WriteCache(raw []byte) error {
	path := CacheFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}
