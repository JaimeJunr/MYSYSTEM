package entities

import (
	"fmt"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ZshPlugin represents a Zsh/Oh My Zsh plugin
type ZshPlugin struct {
	ID          string
	Name        string
	Description string
	Source      types.PluginSource // BuiltIn, External, Custom

	// Installation
	RepoURL    string // For external plugins (GitHub URL)
	InstallCmd string // Custom install command (optional)
	CheckCmd   string // Command to verify if installed (optional)

	// Configuration
	LoadOrder  int    // Order of loading (lower = earlier)
	ConfigFile string // Path to plugin-specific config file
}

// Validate checks if the plugin configuration is valid
func (zp *ZshPlugin) Validate() error {
	if zp.ID == "" {
		return types.ErrInvalidInput
	}
	if zp.Name == "" {
		return types.ErrInvalidInput
	}
	if !zp.Source.IsValid() {
		return types.ErrInvalidInput
	}

	// External plugins must have either RepoURL or InstallCmd
	if zp.Source == types.PluginSourceExternal {
		if zp.RepoURL == "" && zp.InstallCmd == "" {
			return types.ErrInvalidInput
		}
	}

	return nil
}

// IsBuiltIn returns true if this is a built-in Oh My Zsh plugin
func (zp *ZshPlugin) IsBuiltIn() bool {
	return zp.Source == types.PluginSourceBuiltIn
}

// IsExternal returns true if this is an external plugin
func (zp *ZshPlugin) IsExternal() bool {
	return zp.Source == types.PluginSourceExternal
}

// IsCustom returns true if this is a custom user plugin
func (zp *ZshPlugin) IsCustom() bool {
	return zp.Source == types.PluginSourceCustom
}

// NeedsInstallation returns true if this plugin requires installation
func (zp *ZshPlugin) NeedsInstallation() bool {
	// Built-in plugins come with Oh My Zsh, no installation needed
	if zp.IsBuiltIn() {
		return false
	}

	// External plugins always need installation
	if zp.IsExternal() {
		return true
	}

	// Custom plugins need installation if they have an InstallCmd
	if zp.IsCustom() && zp.InstallCmd != "" {
		return true
	}

	return false
}

// GetInstallCommand returns the command to install this plugin
func (zp *ZshPlugin) GetInstallCommand() string {
	// If explicit InstallCmd is provided, use it
	if zp.InstallCmd != "" {
		return zp.InstallCmd
	}

	// Built-in plugins don't need installation
	if zp.IsBuiltIn() {
		return ""
	}

	// For external plugins with RepoURL, generate git clone command
	if zp.IsExternal() && zp.RepoURL != "" {
		return fmt.Sprintf("git clone %s ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/%s", zp.RepoURL, zp.ID)
	}

	return ""
}

// GetCheckCommand returns the command to check if plugin is installed
func (zp *ZshPlugin) GetCheckCommand() string {
	// If explicit CheckCmd is provided, use it
	if zp.CheckCmd != "" {
		return zp.CheckCmd
	}

	// Generate default check command based on source
	if zp.IsBuiltIn() {
		return fmt.Sprintf("test -f $ZSH/plugins/%s/%s.plugin.zsh", zp.ID, zp.ID)
	}

	if zp.IsExternal() {
		return fmt.Sprintf("test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/%s", zp.ID)
	}

	return ""
}
