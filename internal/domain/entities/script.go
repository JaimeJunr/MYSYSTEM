package entities

import (
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// Script represents a system maintenance script
type Script struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Category     types.Category
	RequiresSudo bool
}

// Validate checks if the script entity is valid
func (s *Script) Validate() error {
	if s.ID == "" {
		return types.ErrInvalidInput
	}
	if s.Name == "" {
		return types.ErrInvalidInput
	}
	if s.Path == "" {
		return types.ErrInvalidInput
	}
	if !s.Category.IsValid() {
		return types.ErrInvalidInput
	}
	return nil
}

// IsCleanup returns true if the script is a cleanup script
func (s *Script) IsCleanup() bool {
	return s.Category == types.CategoryCleanup
}

// IsMonitoring returns true if the script is a monitoring script
func (s *Script) IsMonitoring() bool {
	return s.Category == types.CategoryMonitoring
}

// IsInstall returns true if the script is an install script
func (s *Script) IsInstall() bool {
	return s.Category == types.CategoryInstall
}
