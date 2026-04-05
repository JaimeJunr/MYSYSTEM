package preferences

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ThemeDark  = "dark"
	ThemeLight = "light"
)

const (
	TextScaleNormal  = "normal"
	TextScaleLarge   = "large"
	TextScaleXLarge  = "xlarge"
)

type Preferences struct {
	CatalogURL           string `yaml:"catalog_url,omitempty"`
	ScriptRoot           string `yaml:"script_root,omitempty"`
	DotfilesRepo         string `yaml:"dotfiles_repo,omitempty"`
	Theme                string `yaml:"theme,omitempty"`
	TextScale            string `yaml:"text_scale,omitempty"`
	HighContrast         bool   `yaml:"high_contrast"`
	ReduceMotion         bool   `yaml:"reduce_motion"`
	ConfirmBeforeScript  bool   `yaml:"confirm_before_script"`
	ConfirmBeforePackage bool   `yaml:"confirm_before_package"`
}

type rawPreferences struct {
	CatalogURL           string `yaml:"catalog_url"`
	ScriptRoot           string `yaml:"script_root"`
	DotfilesRepo         string `yaml:"dotfiles_repo"`
	Theme                string `yaml:"theme"`
	TextScale            string `yaml:"text_scale"`
	HighContrast         *bool  `yaml:"high_contrast"`
	ReduceMotion         *bool  `yaml:"reduce_motion"`
	ConfirmBeforeScript  *bool  `yaml:"confirm_before_script"`
	ConfirmBeforePackage *bool  `yaml:"confirm_before_package"`
}

func DefaultPreferences() Preferences {
	return Preferences{
		Theme:                ThemeDark,
		TextScale:            TextScaleNormal,
		ConfirmBeforeScript:  true,
		ConfirmBeforePackage: true,
	}
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(dir, "homestead", "preferences.yaml"), nil
}

func Load(path string) (Preferences, error) {
	p := DefaultPreferences()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return p, nil
		}
		return p, fmt.Errorf("read preferences: %w", err)
	}
	var raw rawPreferences
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return DefaultPreferences(), fmt.Errorf("parse preferences: %w", err)
	}
	p = fromRaw(raw)
	p.Normalize()
	return p, nil
}

func fromRaw(r rawPreferences) Preferences {
	p := DefaultPreferences()
	p.CatalogURL = strings.TrimSpace(r.CatalogURL)
	p.ScriptRoot = strings.TrimSpace(r.ScriptRoot)
	p.DotfilesRepo = strings.TrimSpace(r.DotfilesRepo)
	if r.Theme != "" {
		p.Theme = strings.TrimSpace(r.Theme)
	}
	if strings.TrimSpace(r.TextScale) != "" {
		p.TextScale = strings.TrimSpace(r.TextScale)
	}
	if r.HighContrast != nil {
		p.HighContrast = *r.HighContrast
	}
	if r.ReduceMotion != nil {
		p.ReduceMotion = *r.ReduceMotion
	}
	if r.ConfirmBeforeScript != nil {
		p.ConfirmBeforeScript = *r.ConfirmBeforeScript
	}
	if r.ConfirmBeforePackage != nil {
		p.ConfirmBeforePackage = *r.ConfirmBeforePackage
	}
	return p
}

func (p *Preferences) Normalize() {
	if p.Theme != ThemeLight && p.Theme != ThemeDark {
		p.Theme = ThemeDark
	}
	switch p.TextScale {
	case TextScaleNormal, TextScaleLarge, TextScaleXLarge:
	default:
		p.TextScale = TextScaleNormal
	}
}

func Save(path string, p Preferences) error {
	p.Normalize()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("preferences dir: %w", err)
	}
	data, err := yaml.Marshal(&p)
	if err != nil {
		return fmt.Errorf("marshal preferences: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write preferences: %w", err)
	}
	return nil
}

func DefaultDotfilesRepo() string {
	return "~/.config/homestead-dotfiles"
}

func ExpandPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", nil
	}
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home dir: %w", err)
		}
		if path == "~" {
			return filepath.Clean(home), nil
		}
		return filepath.Clean(filepath.Join(home, path[2:])), nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func ValidateScriptRoot(expandedRoot string) error {
	if expandedRoot == "" {
		return nil
	}
	st, err := os.Stat(expandedRoot)
	if err != nil {
		return fmt.Errorf("raiz dos scripts: %w", err)
	}
	if !st.IsDir() {
		return fmt.Errorf("raiz dos scripts: não é um diretório")
	}
	scriptsDir := filepath.Join(expandedRoot, "scripts")
	st2, err := os.Stat(scriptsDir)
	if err != nil || !st2.IsDir() {
		return fmt.Errorf("raiz dos scripts: precisa conter um diretório scripts/")
	}
	return nil
}

func ValidateCatalogURL(u string) error {
	u = strings.TrimSpace(u)
	if u == "" {
		return nil
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("URL do catálogo inválida")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL do catálogo: use http ou https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL do catálogo: host em falta")
	}
	return nil
}
