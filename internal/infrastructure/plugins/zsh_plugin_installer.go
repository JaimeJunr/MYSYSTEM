package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ZshPluginInstaller implements PluginManager interface
type ZshPluginInstaller struct {
	zshDir    string // Path to Oh My Zsh installation (e.g., ~/.oh-my-zsh)
	customDir string // Path to custom plugins directory (e.g., ~/.oh-my-zsh/custom)
}

// NewZshPluginInstaller creates a new Zsh plugin installer
func NewZshPluginInstaller(zshDir, customDir string) interfaces.PluginManager {
	return &ZshPluginInstaller{
		zshDir:    zshDir,
		customDir: customDir,
	}
}

// InstallPlugin installs a Zsh plugin with progress reporting
func (zpi *ZshPluginInstaller) InstallPlugin(plugin *entities.ZshPlugin, progressCallback interfaces.PluginProgressCallback) error {
	// Validate plugin
	if err := zpi.ValidatePlugin(plugin); err != nil {
		return fmt.Errorf("plugin validation failed: %w", err)
	}

	// Report progress
	if progressCallback != nil {
		progressCallback(interfaces.PluginInstallProgress{
			PluginID:   plugin.ID,
			PluginName: plugin.Name,
			Status:     "starting",
			Message:    fmt.Sprintf("Installing %s...", plugin.Name),
		})
	}

	// Built-in plugins don't need installation
	if plugin.IsBuiltIn() {
		// Just verify it exists
		installed, err := zpi.IsPluginInstalled(plugin.ID)
		if err != nil {
			return err
		}

		if !installed {
			return fmt.Errorf("built-in plugin '%s' not found in Oh My Zsh installation", plugin.ID)
		}

		if progressCallback != nil {
			progressCallback(interfaces.PluginInstallProgress{
				PluginID:    plugin.ID,
				PluginName:  plugin.Name,
				Status:      "complete",
				Message:     fmt.Sprintf("%s is already available (built-in)", plugin.Name),
				IsCompleted: true,
			})
		}

		return nil
	}

	// Check if already installed
	installed, err := zpi.IsPluginInstalled(plugin.ID)
	if err != nil {
		return err
	}

	if installed {
		if progressCallback != nil {
			progressCallback(interfaces.PluginInstallProgress{
				PluginID:    plugin.ID,
				PluginName:  plugin.Name,
				Status:      "complete",
				Message:     fmt.Sprintf("%s is already installed", plugin.Name),
				IsCompleted: true,
			})
		}
		return nil
	}

	// Install external/custom plugin
	if plugin.IsExternal() || plugin.IsCustom() {
		// Ensure custom plugins directory exists
		customPluginsDir := filepath.Join(zpi.customDir, "plugins")
		if err := os.MkdirAll(customPluginsDir, 0755); err != nil {
			return fmt.Errorf("failed to create custom plugins directory: %w", err)
		}

		if progressCallback != nil {
			progressCallback(interfaces.PluginInstallProgress{
				PluginID:   plugin.ID,
				PluginName: plugin.Name,
				Status:     "installing",
				Progress:   50,
				Message:    fmt.Sprintf("Cloning %s...", plugin.Name),
			})
		}

		// Get install command
		installCmd := plugin.GetInstallCommand()
		if installCmd == "" {
			return fmt.Errorf("no install command available for plugin '%s'", plugin.ID)
		}

		// Execute install command
		cmd := exec.Command("sh", "-c", installCmd)
		cmd.Dir = customPluginsDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install plugin '%s': %w\nOutput: %s", plugin.ID, err, string(output))
		}

		if progressCallback != nil {
			progressCallback(interfaces.PluginInstallProgress{
				PluginID:    plugin.ID,
				PluginName:  plugin.Name,
				Status:      "complete",
				Progress:    100,
				Message:     fmt.Sprintf("%s installed successfully", plugin.Name),
				IsCompleted: true,
			})
		}

		return nil
	}

	return fmt.Errorf("unsupported plugin source: %s", plugin.Source)
}

// IsPluginInstalled checks if a plugin is already installed
func (zpi *ZshPluginInstaller) IsPluginInstalled(pluginID string) (bool, error) {
	// Check built-in plugins directory
	builtInPath := filepath.Join(zpi.zshDir, "plugins", pluginID, pluginID+".plugin.zsh")
	if _, err := os.Stat(builtInPath); err == nil {
		return true, nil
	}

	// Check custom plugins directory
	customPath := filepath.Join(zpi.customDir, "plugins", pluginID)
	if _, err := os.Stat(customPath); err == nil {
		return true, nil
	}

	return false, nil
}

// UninstallPlugin removes an installed plugin
func (zpi *ZshPluginInstaller) UninstallPlugin(pluginID string) error {
	// Check if it's a built-in plugin
	builtInPath := filepath.Join(zpi.zshDir, "plugins", pluginID)
	if _, err := os.Stat(builtInPath); err == nil {
		return fmt.Errorf("cannot uninstall built-in plugin '%s'", pluginID)
	}

	// Remove from custom plugins directory
	customPath := filepath.Join(zpi.customDir, "plugins", pluginID)
	if _, err := os.Stat(customPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin '%s' not found", pluginID)
		}
		return fmt.Errorf("failed to check plugin '%s': %w", pluginID, err)
	}

	// Remove directory
	if err := os.RemoveAll(customPath); err != nil {
		return fmt.Errorf("failed to uninstall plugin '%s': %w", pluginID, err)
	}

	return nil
}

// ListInstalledPlugins returns a list of currently installed plugin IDs
func (zpi *ZshPluginInstaller) ListInstalledPlugins() ([]string, error) {
	var plugins []string

	// List built-in plugins
	builtInDir := filepath.Join(zpi.zshDir, "plugins")
	if entries, err := os.ReadDir(builtInDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				// Verify it's a valid plugin (has .plugin.zsh file)
				pluginFile := filepath.Join(builtInDir, entry.Name(), entry.Name()+".plugin.zsh")
				if _, err := os.Stat(pluginFile); err == nil {
					plugins = append(plugins, entry.Name())
				}
			}
		}
	}

	// List custom plugins
	customPluginsDir := filepath.Join(zpi.customDir, "plugins")
	if entries, err := os.ReadDir(customPluginsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				plugins = append(plugins, entry.Name())
			}
		}
	}

	return plugins, nil
}

// UpdatePlugin updates an external plugin to the latest version
func (zpi *ZshPluginInstaller) UpdatePlugin(pluginID string) error {
	// Only custom/external plugins can be updated
	customPath := filepath.Join(zpi.customDir, "plugins", pluginID)

	// Check if plugin exists
	if _, err := os.Stat(customPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("plugin '%s' not found", pluginID)
		}
		return fmt.Errorf("failed to check plugin '%s': %w", pluginID, err)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(customPath, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("plugin '%s' is not a git repository, cannot update", pluginID)
	}

	// Run git pull
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = customPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update plugin '%s': %w\nOutput: %s", pluginID, err, string(output))
	}

	return nil
}

// ValidatePlugin checks if a plugin installation is valid and functional
func (zpi *ZshPluginInstaller) ValidatePlugin(plugin *entities.ZshPlugin) error {
	// Use entity's built-in validation
	if err := plugin.Validate(); err != nil {
		return err
	}

	// Additional validation for external plugins
	if plugin.IsExternal() {
		// Check if RepoURL is accessible (basic validation)
		if plugin.RepoURL == "" && plugin.InstallCmd == "" {
			return types.ErrInvalidInput
		}
	}

	return nil
}
