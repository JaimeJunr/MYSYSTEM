package entities

import (
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

const (
	NativeMonitorNone         = ""
	NativeMonitorBattery      = "battery"
	NativeMonitorMemory       = "memory"
	NativeMonitorDisk         = "disk"
	NativeMonitorLoad         = "load"
	NativeMonitorNetwork      = "network"
	NativeMonitorThermal      = "thermal"
	NativeMonitorSystemdUser  = "systemd-user"
)

func ValidNativeMonitor(k string) bool {
	switch k {
	case NativeMonitorBattery, NativeMonitorMemory, NativeMonitorDisk,
		NativeMonitorLoad, NativeMonitorNetwork, NativeMonitorThermal,
		NativeMonitorSystemdUser:
		return true
	default:
		return false
	}
}

// Script represents a system maintenance script
type Script struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Category     types.Category
	RequiresSudo bool
	SupportsDryRun bool
	// If set, TUI shows a native panel and Path is ignored.
	NativeMonitor string
}

// Validate checks if the script entity is valid
func (s *Script) Validate() error {
	if s.ID == "" {
		return types.ErrInvalidInput
	}
	if s.Name == "" {
		return types.ErrInvalidInput
	}
	if !s.Category.IsValid() {
		return types.ErrInvalidInput
	}
	if s.NativeMonitor != "" {
		if !ValidNativeMonitor(s.NativeMonitor) {
			return types.ErrInvalidInput
		}
		return nil
	}
	if s.Path == "" {
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

// IsUtilities returns true if the script is a desktop utility installer (Flatpak, pacotes, etc.)
func (s *Script) IsUtilities() bool {
	return s.Category == types.CategoryUtilities
}
