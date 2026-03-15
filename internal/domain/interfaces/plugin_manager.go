package interfaces

import "github.com/JaimeJunr/Homestead/internal/domain/entities"

// PluginInstallProgress represents the progress of a plugin installation
type PluginInstallProgress struct {
	PluginID    string
	PluginName  string
	Status      string // "cloning", "installing", "complete", "failed"
	Progress    int    // 0-100
	Message     string
	Error       error
	IsCompleted bool
}

// PluginProgressCallback is called to report plugin installation progress
type PluginProgressCallback func(progress PluginInstallProgress)

// PluginManager defines the interface for managing Zsh plugins
type PluginManager interface {
	// InstallPlugin installs a Zsh plugin with progress reporting
	InstallPlugin(plugin *entities.ZshPlugin, progressCallback PluginProgressCallback) error

	// IsPluginInstalled checks if a plugin is already installed
	IsPluginInstalled(pluginID string) (bool, error)

	// UninstallPlugin removes an installed plugin
	UninstallPlugin(pluginID string) error

	// ListInstalledPlugins returns a list of currently installed plugin IDs
	ListInstalledPlugins() ([]string, error)

	// UpdatePlugin updates an external plugin to the latest version
	UpdatePlugin(pluginID string) error

	// ValidatePlugin checks if a plugin installation is valid and functional
	ValidatePlugin(plugin *entities.ZshPlugin) error
}
