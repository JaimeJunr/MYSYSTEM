package entities

import (
	"strings"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// Package represents a software package that can be installed
type Package struct {
	ID          string
	Name        string
	Description string
	Version     string // "latest" or specific version
	Category    types.PackageCategory
	DownloadURL string
	InstallCmd  string // Command to install after download
	CheckCmd    string // Command to check if already installed
	Notes       string // optional; extra text on install confirmation
	ProjectURL        string // official site or source repo
	UtilityScriptPath string // repo-relative bash installer; utilities category only
	RequiresSudo      bool
}

// Validate checks if the package is valid
func (p *Package) Validate() error {
	if p.ID == "" {
		return types.ErrInvalidInput
	}
	if p.Name == "" {
		return types.ErrInvalidInput
	}
	ut := strings.TrimSpace(p.UtilityScriptPath)
	if p.DownloadURL == "" && p.InstallCmd == "" && ut == "" {
		return types.ErrInvalidInput
	}
	if !p.Category.IsValid() {
		return types.ErrInvalidInput
	}
	if ut != "" && p.Category != types.PackageCategoryUtilities {
		return types.ErrInvalidInput
	}
	if ut != "" && (p.DownloadURL != "" || p.InstallCmd != "") {
		return types.ErrInvalidInput
	}
	if p.Category == types.PackageCategoryUtilities {
		if ut == "" || strings.TrimSpace(p.ProjectURL) == "" {
			return types.ErrInvalidInput
		}
	}
	return nil
}

// ResolveInstallKind returns the installer strategy implied by package fields.
func (p *Package) ResolveInstallKind() types.PackageInstallKind {
	if strings.TrimSpace(p.UtilityScriptPath) != "" {
		return types.InstallKindUtilityScript
	}
	if p.DownloadURL != "" {
		return types.InstallKindShellWithDownload
	}
	return types.InstallKindShellLocal
}

// IsIDE returns true if this package is an IDE
func (p *Package) IsIDE() bool {
	return p.Category == types.PackageCategoryIDE
}

// IsTool returns true if this package is a development tool
func (p *Package) IsTool() bool {
	return p.Category == types.PackageCategoryTool
}
