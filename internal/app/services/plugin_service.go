package services

import (
	"fmt"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// PluginService provides high-level plugin management operations
type PluginService struct {
	pluginManager    interfaces.PluginManager
	availablePlugins []*entities.ZshPlugin
}

// NewPluginService creates a new plugin service
func NewPluginService(pluginManager interfaces.PluginManager) *PluginService {
	return &PluginService{
		pluginManager:    pluginManager,
		availablePlugins: make([]*entities.ZshPlugin, 0),
	}
}

// InstallPlugin installs a Zsh plugin with progress reporting
func (ps *PluginService) InstallPlugin(plugin *entities.ZshPlugin, progressCallback interfaces.PluginProgressCallback) error {
	// Validate plugin
	if err := ps.ValidatePlugin(plugin); err != nil {
		return fmt.Errorf("invalid plugin: %w", err)
	}

	// Install plugin
	if err := ps.pluginManager.InstallPlugin(plugin, progressCallback); err != nil {
		return fmt.Errorf("failed to install plugin '%s': %w", plugin.ID, err)
	}

	return nil
}

// IsPluginInstalled checks if a plugin is installed
func (ps *PluginService) IsPluginInstalled(pluginID string) (bool, error) {
	installed, err := ps.pluginManager.IsPluginInstalled(pluginID)
	if err != nil {
		return false, fmt.Errorf("failed to check if plugin '%s' is installed: %w", pluginID, err)
	}

	return installed, nil
}

// ListInstalledPlugins returns list of installed plugins
func (ps *PluginService) ListInstalledPlugins() ([]string, error) {
	plugins, err := ps.pluginManager.ListInstalledPlugins()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed plugins: %w", err)
	}

	return plugins, nil
}

// UninstallPlugin removes an installed plugin
func (ps *PluginService) UninstallPlugin(pluginID string) error {
	if err := ps.pluginManager.UninstallPlugin(pluginID); err != nil {
		return fmt.Errorf("failed to uninstall plugin '%s': %w", pluginID, err)
	}

	return nil
}

// UpdatePlugin updates a plugin to the latest version
func (ps *PluginService) UpdatePlugin(pluginID string) error {
	if err := ps.pluginManager.UpdatePlugin(pluginID); err != nil {
		return fmt.Errorf("failed to update plugin '%s': %w", pluginID, err)
	}

	return nil
}

// ValidatePlugin validates a plugin configuration
func (ps *PluginService) ValidatePlugin(plugin *entities.ZshPlugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}

	// Use plugin manager's validation
	if err := ps.pluginManager.ValidatePlugin(plugin); err != nil {
		return err
	}

	return nil
}

// SetAvailablePlugins sets the list of available plugins
func (ps *PluginService) SetAvailablePlugins(plugins []*entities.ZshPlugin) {
	ps.availablePlugins = plugins
}

// GetAvailablePlugins returns all available plugins
func (ps *PluginService) GetAvailablePlugins() []*entities.ZshPlugin {
	return ps.availablePlugins
}

// GetPluginsBySource returns plugins filtered by source
func (ps *PluginService) GetPluginsBySource(source types.PluginSource) []*entities.ZshPlugin {
	var filtered []*entities.ZshPlugin
	for _, plugin := range ps.availablePlugins {
		if plugin.Source == source {
			filtered = append(filtered, plugin)
		}
	}
	return filtered
}

// GetPluginByID returns a plugin by ID
func (ps *PluginService) GetPluginByID(id string) *entities.ZshPlugin {
	for _, plugin := range ps.availablePlugins {
		if plugin.ID == id {
			return plugin
		}
	}
	return nil
}

// InstallMultiplePlugins installs multiple plugins with progress reporting
func (ps *PluginService) InstallMultiplePlugins(plugins []*entities.ZshPlugin, progressCallback interfaces.PluginProgressCallback) []PluginInstallResult {
	results := make([]PluginInstallResult, len(plugins))

	for i, plugin := range plugins {
		err := ps.InstallPlugin(plugin, progressCallback)
		results[i] = PluginInstallResult{
			Plugin:  plugin,
			Success: err == nil,
			Error:   err,
		}
	}

	return results
}

// GetInstallationStatus returns the installation status of a plugin
func (ps *PluginService) GetInstallationStatus(plugin *entities.ZshPlugin) PluginInstallationStatus {
	installed, err := ps.IsPluginInstalled(plugin.ID)

	return PluginInstallationStatus{
		Plugin:      plugin,
		IsInstalled: installed,
		Error:       err,
	}
}

// PluginInstallResult represents the result of a plugin installation
type PluginInstallResult struct {
	Plugin  *entities.ZshPlugin
	Success bool
	Error   error
}

// PluginInstallationStatus represents the installation status of a plugin
type PluginInstallationStatus struct {
	Plugin      *entities.ZshPlugin
	IsInstalled bool
	Error       error
}
