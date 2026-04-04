package catalog

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

type manifestJSON struct {
	SchemaVersion int           `json:"schema_version"`
	Packages      []packageJSON `json:"packages"`
}

type packageJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Category    string `json:"category"`
	DownloadURL string `json:"download_url"`
	InstallCmd  string `json:"install_cmd"`
	CheckCmd    string `json:"check_cmd"`
	Notes       string `json:"notes"`
	ProjectURL  string `json:"project_url"`
}

func normalizeCategory(s string) types.PackageCategory {
	c := types.PackageCategory(strings.ToLower(strings.TrimSpace(s)))
	if c.IsValid() {
		return c
	}
	return types.PackageCategoryOther
}

// ParseManifest decodes JSON and returns packages when schema_version is supported.
func ParseManifest(data []byte) ([]entities.Package, int, error) {
	var m manifestJSON
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, 0, err
	}
	if !IsSchemaSupported(m.SchemaVersion) {
		return nil, m.SchemaVersion, ErrUnsupportedSchema
	}
	var out []entities.Package
	for _, p := range m.Packages {
		pkg := entities.Package{
			ID:          strings.TrimSpace(p.ID),
			Name:        strings.TrimSpace(p.Name),
			Description: strings.TrimSpace(p.Description),
			Version:     strings.TrimSpace(p.Version),
			Category:    normalizeCategory(p.Category),
			DownloadURL: strings.TrimSpace(p.DownloadURL),
			InstallCmd:  strings.TrimSpace(p.InstallCmd),
			CheckCmd:    strings.TrimSpace(p.CheckCmd),
			Notes:       strings.TrimSpace(p.Notes),
			ProjectURL:  strings.TrimSpace(p.ProjectURL),
		}
		if err := pkg.Validate(); err != nil {
			continue
		}
		out = append(out, pkg)
	}
	return out, m.SchemaVersion, nil
}

// ReadAndParseCacheFile reads the cache path and parses the manifest.
func ReadAndParseCacheFile(path string) ([]entities.Package, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	pkgs, sv, err := ParseManifest(data)
	if err != nil {
		return nil, sv, err
	}
	return pkgs, sv, nil
}
