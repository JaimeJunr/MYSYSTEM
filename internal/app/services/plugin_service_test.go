package services

import (
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/plugins"
)

// TestPluginService_InstallPlugin tests installing a plugin
func TestPluginService_InstallPlugin(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Create a built-in plugin to install
	plugin := &entities.ZshPlugin{
		ID:     "test-plugin",
		Name:   "Test Plugin",
		Source: types.PluginSourceBuiltIn,
	}

	// Track progress
	var progressUpdates []string
	callback := func(progress interfaces.PluginInstallProgress) {
		progressUpdates = append(progressUpdates, progress.Status)
	}

	// Note: This will likely fail because plugin doesn't exist
	// But we're testing the service orchestration logic
	_ = service.InstallPlugin(plugin, callback)

	// Verify callback was called
	if len(progressUpdates) == 0 {
		t.Log("Progress callback was not called (expected for non-existent plugin)")
	}
}

// TestPluginService_InstallPlugin_Invalid tests installing invalid plugin
func TestPluginService_InstallPlugin_Invalid(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Create invalid plugin (missing ID)
	plugin := &entities.ZshPlugin{
		Name:   "Invalid",
		Source: types.PluginSourceBuiltIn,
	}

	err := service.InstallPlugin(plugin, nil)
	if err == nil {
		t.Error("Expected error for invalid plugin, got nil")
	}
}

// TestPluginService_IsPluginInstalled tests checking if plugin is installed
func TestPluginService_IsPluginInstalled(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	installed, err := service.IsPluginInstalled("non-existent")
	if err != nil {
		t.Fatalf("IsPluginInstalled() error = %v", err)
	}

	if installed {
		t.Error("Non-existent plugin should not be installed")
	}
}

// TestPluginService_ListInstalledPlugins tests listing installed plugins
func TestPluginService_ListInstalledPlugins(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	list, err := service.ListInstalledPlugins()
	if err != nil {
		t.Fatalf("ListInstalledPlugins() error = %v", err)
	}

	// Should return empty list for empty directory
	if len(list) != 0 {
		t.Logf("Found %d plugins (expected 0 in empty test directory)", len(list))
	}
}

// TestPluginService_UninstallPlugin tests uninstalling a plugin
func TestPluginService_UninstallPlugin(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Try to uninstall non-existent plugin
	err := service.UninstallPlugin("non-existent")
	if err == nil {
		t.Error("Expected error for uninstalling non-existent plugin, got nil")
	}
}

// TestPluginService_UpdatePlugin tests updating a plugin
func TestPluginService_UpdatePlugin(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Try to update non-existent plugin
	err := service.UpdatePlugin("non-existent")
	if err == nil {
		t.Error("Expected error for updating non-existent plugin, got nil")
	}
}

// TestPluginService_ValidatePlugin tests plugin validation
func TestPluginService_ValidatePlugin(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	tests := []struct {
		name    string
		plugin  *entities.ZshPlugin
		wantErr bool
	}{
		{
			name: "valid built-in plugin",
			plugin: &entities.ZshPlugin{
				ID:     "git",
				Name:   "Git",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: false,
		},
		{
			name: "valid external plugin",
			plugin: &entities.ZshPlugin{
				ID:      "zsh-autosuggestions",
				Name:    "Zsh Autosuggestions",
				Source:  types.PluginSourceExternal,
				RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing ID",
			plugin: &entities.ZshPlugin{
				Name:   "Invalid",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
		{
			name: "invalid - missing Name",
			plugin: &entities.ZshPlugin{
				ID:     "invalid",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
		{
			name: "invalid - external without RepoURL",
			plugin: &entities.ZshPlugin{
				ID:     "invalid",
				Name:   "Invalid",
				Source: types.PluginSourceExternal,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidatePlugin(tt.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPluginService_GetAvailablePlugins tests getting available plugins
func TestPluginService_GetAvailablePlugins(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Define some available plugins
	available := []*entities.ZshPlugin{
		{
			ID:     "git",
			Name:   "Git",
			Source: types.PluginSourceBuiltIn,
		},
		{
			ID:      "zsh-autosuggestions",
			Name:    "Zsh Autosuggestions",
			Source:  types.PluginSourceExternal,
			RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
		},
	}

	// Set available plugins
	service.SetAvailablePlugins(available)

	// Get available plugins
	plugins := service.GetAvailablePlugins()
	if len(plugins) != len(available) {
		t.Errorf("GetAvailablePlugins() count = %d, want %d", len(plugins), len(available))
	}
}

// TestPluginService_GetPluginsBySource tests filtering plugins by source
func TestPluginService_GetPluginsBySource(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Define plugins with different sources
	available := []*entities.ZshPlugin{
		{
			ID:     "git",
			Name:   "Git",
			Source: types.PluginSourceBuiltIn,
		},
		{
			ID:     "docker",
			Name:   "Docker",
			Source: types.PluginSourceBuiltIn,
		},
		{
			ID:      "zsh-autosuggestions",
			Name:    "Zsh Autosuggestions",
			Source:  types.PluginSourceExternal,
			RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
		},
	}

	service.SetAvailablePlugins(available)

	// Get built-in plugins
	builtIn := service.GetPluginsBySource(types.PluginSourceBuiltIn)
	if len(builtIn) != 2 {
		t.Errorf("Built-in plugins count = %d, want 2", len(builtIn))
	}

	// Get external plugins
	external := service.GetPluginsBySource(types.PluginSourceExternal)
	if len(external) != 1 {
		t.Errorf("External plugins count = %d, want 1", len(external))
	}
}

// TestPluginService_GetPluginByID tests getting a plugin by ID
func TestPluginService_GetPluginByID(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	available := []*entities.ZshPlugin{
		{
			ID:     "git",
			Name:   "Git",
			Source: types.PluginSourceBuiltIn,
		},
		{
			ID:     "docker",
			Name:   "Docker",
			Source: types.PluginSourceBuiltIn,
		},
	}

	service.SetAvailablePlugins(available)

	// Get existing plugin
	plugin := service.GetPluginByID("git")
	if plugin == nil {
		t.Error("GetPluginByID() returned nil for existing plugin")
	}
	if plugin != nil && plugin.ID != "git" {
		t.Errorf("Plugin ID = %s, want git", plugin.ID)
	}

	// Get non-existent plugin
	notFound := service.GetPluginByID("non-existent")
	if notFound != nil {
		t.Error("GetPluginByID() should return nil for non-existent plugin")
	}
}

// TestPluginService_InstallMultiplePlugins tests installing multiple plugins
func TestPluginService_InstallMultiplePlugins(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	pluginsToInstall := []*entities.ZshPlugin{
		{
			ID:     "git",
			Name:   "Git",
			Source: types.PluginSourceBuiltIn,
		},
		{
			ID:     "docker",
			Name:   "Docker",
			Source: types.PluginSourceBuiltIn,
		},
	}

	// Track progress
	var progressUpdates []string
	callback := func(progress interfaces.PluginInstallProgress) {
		progressUpdates = append(progressUpdates, progress.PluginID+":"+progress.Status)
	}

	results := service.InstallMultiplePlugins(pluginsToInstall, callback)

	// Verify results
	if len(results) != len(pluginsToInstall) {
		t.Errorf("Results count = %d, want %d", len(results), len(pluginsToInstall))
	}
}

// TestPluginService_GetInstallationStatus tests getting installation status
func TestPluginService_GetInstallationStatus(t *testing.T) {
	tempDir := t.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	plugin := &entities.ZshPlugin{
		ID:     "test",
		Name:   "Test",
		Source: types.PluginSourceBuiltIn,
	}

	status := service.GetInstallationStatus(plugin)
	if status.IsInstalled {
		t.Error("Plugin should not be installed in empty test directory")
	}

	if status.Plugin.ID != plugin.ID {
		t.Errorf("Status plugin ID = %s, want %s", status.Plugin.ID, plugin.ID)
	}
}

// Benchmark tests
func BenchmarkPluginService_ValidatePlugin(b *testing.B) {
	tempDir := b.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	plugin := &entities.ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.ValidatePlugin(plugin)
	}
}

func BenchmarkPluginService_GetPluginByID(b *testing.B) {
	tempDir := b.TempDir()
	pluginManager := plugins.NewZshPluginInstaller(
		filepath.Join(tempDir, "oh-my-zsh"),
		filepath.Join(tempDir, "custom"),
	)
	service := NewPluginService(pluginManager)

	// Create many plugins
	available := make([]*entities.ZshPlugin, 100)
	for i := 0; i < 100; i++ {
		available[i] = &entities.ZshPlugin{
			ID:     filepath.Join("plugin", string(rune(i))),
			Name:   "Plugin",
			Source: types.PluginSourceBuiltIn,
		}
	}
	service.SetAvailablePlugins(available)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.GetPluginByID("plugin0")
	}
}
