package profilestate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
)

// ExportPreferences mirrors migration-relevant prefs (no confirm flags).
type ExportPreferences struct {
	CatalogURL   string `json:"catalog_url,omitempty"`
	ScriptRoot   string `json:"script_root,omitempty"`
	DotfilesRepo string `json:"dotfiles_repo,omitempty"`
	Theme        string `json:"theme,omitempty"`
}

// ExportDoc is the JSON export shape for cloning setup on another machine.
type ExportDoc struct {
	Version           int               `json:"version"`
	ExportedAt        string            `json:"exported_at"`
	HomesteadVersion  string            `json:"homestead_version"`
	Preferences       ExportPreferences `json:"preferences"`
	InstalledPackages []string          `json:"installed_packages"`
	FavoriteScripts   []string          `json:"favorite_scripts"`
}

func prefsToExport(p preferences.Preferences) ExportPreferences {
	return ExportPreferences{
		CatalogURL:   strings.TrimSpace(p.CatalogURL),
		ScriptRoot:   strings.TrimSpace(p.ScriptRoot),
		DotfilesRepo: strings.TrimSpace(p.DotfilesRepo),
		Theme:        strings.TrimSpace(p.Theme),
	}
}

// WriteExport writes JSON or plain-text profile export to w.
func WriteExport(w io.Writer, format string, s State, p preferences.Preferences, appVersion string) error {
	s.Normalize()
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "json", "":
		doc := ExportDoc{
			Version:           1,
			ExportedAt:        time.Now().UTC().Format(time.RFC3339),
			HomesteadVersion:  strings.TrimSpace(appVersion),
			Preferences:       prefsToExport(p),
			InstalledPackages: append([]string(nil), s.InstalledPackageIDs...),
			FavoriteScripts:   append([]string(nil), s.FavoriteScriptIDs...),
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(doc); err != nil {
			return fmt.Errorf("json encode: %w", err)
		}
		return nil
	case "text", "txt":
		return writeExportText(w, s, p, appVersion)
	default:
		return fmt.Errorf("formato desconhecido: %q (use json ou text)", format)
	}
}

func writeExportText(w io.Writer, s State, p preferences.Preferences, appVersion string) error {
	ep := prefsToExport(p)
	_, err := fmt.Fprintf(w, "Homestead — export de perfil\nExportado (UTC): %s\nVersão do app: %s\n\n",
		time.Now().UTC().Format(time.RFC3339), strings.TrimSpace(appVersion))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "[Preferências]\ncatalog_url: %s\nscript_root: %s\ndotfiles_repo: %s\ntheme: %s\n\n",
		ep.CatalogURL, ep.ScriptRoot, ep.DotfilesRepo, ep.Theme)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "[Pacotes instalados via Homestead (registo no TUI)]")
	if err != nil {
		return err
	}
	if len(s.InstalledPackageIDs) == 0 {
		_, err = fmt.Fprintln(w, "(nenhum)")
	} else {
		for _, id := range s.InstalledPackageIDs {
			if _, err = fmt.Fprintf(w, "- %s\n", id); err != nil {
				return err
			}
		}
	}
	_, err = fmt.Fprintln(w, "\n[Scripts favoritos]")
	if err != nil {
		return err
	}
	if len(s.FavoriteScriptIDs) == 0 {
		_, err = fmt.Fprintln(w, "(nenhum)")
	} else {
		for _, id := range s.FavoriteScriptIDs {
			if _, err = fmt.Fprintf(w, "- %s\n", id); err != nil {
				return err
			}
		}
	}
	return nil
}
