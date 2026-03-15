package entities

import (
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
}

// Validate checks if the package is valid
func (p *Package) Validate() error {
	if p.ID == "" {
		return types.ErrInvalidInput
	}
	if p.Name == "" {
		return types.ErrInvalidInput
	}
	if p.DownloadURL == "" && p.InstallCmd == "" {
		return types.ErrInvalidInput
	}
	if !p.Category.IsValid() {
		return types.ErrInvalidInput
	}
	return nil
}

// IsIDE returns true if this package is an IDE
func (p *Package) IsIDE() bool {
	return p.Category == types.PackageCategoryIDE
}

// IsTool returns true if this package is a development tool
func (p *Package) IsTool() bool {
	return p.Category == types.PackageCategoryTool
}
