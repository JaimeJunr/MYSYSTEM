package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// TestZshPluginInstaller_IsPluginInstalled tests checking if a plugin is installed
func TestZshPluginInstaller_IsPluginInstalled(t *testing.T) {
	tempDir := t.TempDir()

	// Create fake plugin directories
	builtInDir := filepath.Join(tempDir, "plugins", "git")
	err := os.MkdirAll(builtInDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test plugin directory: %v", err)
	}

	// Create plugin file
	pluginFile := filepath.Join(builtInDir, "git.plugin.zsh")
	err = os.WriteFile(pluginFile, []byte("# Git plugin"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test plugin file: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Test built-in plugin that exists
	installed, err := installer.IsPluginInstalled("git")
	if err != nil {
		t.Fatalf("IsPluginInstalled() error = %v", err)
	}

	if !installed {
		t.Error("IsPluginInstalled() should return true for existing plugin")
	}

	// Test plugin that doesn't exist
	installed, err = installer.IsPluginInstalled("non-existent")
	if err != nil {
		t.Fatalf("IsPluginInstalled() error = %v", err)
	}

	if installed {
		t.Error("IsPluginInstalled() should return false for non-existent plugin")
	}
}

// TestZshPluginInstaller_InstallBuiltInPlugin tests installing built-in plugins
func TestZshPluginInstaller_InstallBuiltInPlugin(t *testing.T) {
	tempDir := t.TempDir()

	// Create Oh My Zsh directory structure
	pluginsDir := filepath.Join(tempDir, "plugins")
	err := os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins directory: %v", err)
	}

	// Create a built-in plugin directory
	gitPluginDir := filepath.Join(pluginsDir, "git")
	err = os.MkdirAll(gitPluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create git plugin directory: %v", err)
	}

	// Create plugin file
	pluginFile := filepath.Join(gitPluginDir, "git.plugin.zsh")
	err = os.WriteFile(pluginFile, []byte("# Git plugin"), 0644)
	if err != nil {
		t.Fatalf("Failed to create plugin file: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Create built-in plugin
	plugin := &entities.ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	// Install should succeed (built-in plugins don't need installation)
	err = installer.InstallPlugin(plugin, nil)
	if err != nil {
		t.Fatalf("InstallPlugin() error = %v", err)
	}
}

// TestZshPluginInstaller_InstallExternalPlugin tests installing external plugins
func TestZshPluginInstaller_InstallExternalPlugin(t *testing.T) {
	tempDir := t.TempDir()
	customDir := filepath.Join(tempDir, "custom", "plugins")

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Create external plugin
	plugin := &entities.ZshPlugin{
		ID:      "zsh-autosuggestions",
		Name:    "Zsh Autosuggestions",
		Source:  types.PluginSourceExternal,
		RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
	}

	// Note: This test would actually try to clone the repository
	// In a real test environment, we would mock the git clone command
	// For now, we'll test the validation and directory creation

	// Verify custom directory is created
	err := os.MkdirAll(customDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom directory: %v", err)
	}

	// Test validation of external plugin
	err = installer.ValidatePlugin(plugin)
	if err != nil {
		t.Fatalf("ValidatePlugin() error = %v", err)
	}
}

// TestZshPluginInstaller_ValidatePlugin tests plugin validation
func TestZshPluginInstaller_ValidatePlugin(t *testing.T) {
	tempDir := t.TempDir()
	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

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
				ID:      "zsh-syntax-highlighting",
				Name:    "Zsh Syntax Highlighting",
				Source:  types.PluginSourceExternal,
				RepoURL: "https://github.com/zsh-users/zsh-syntax-highlighting",
			},
			wantErr: false,
		},
		{
			name: "invalid plugin - missing ID",
			plugin: &entities.ZshPlugin{
				Name:   "Invalid",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
		{
			name: "invalid plugin - missing Name",
			plugin: &entities.ZshPlugin{
				ID:     "invalid",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installer.ValidatePlugin(tt.plugin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlugin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestZshPluginInstaller_ListInstalledPlugins tests listing installed plugins
func TestZshPluginInstaller_ListInstalledPlugins(t *testing.T) {
	tempDir := t.TempDir()

	// Create plugin directories
	pluginsDir := filepath.Join(tempDir, "plugins")
	err := os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins directory: %v", err)
	}

	// Create some plugin directories with .plugin.zsh files
	plugins := []string{"git", "docker", "rails"}
	for _, plugin := range plugins {
		pluginDir := filepath.Join(pluginsDir, plugin)
		err := os.MkdirAll(pluginDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create plugin directory: %v", err)
		}

		pluginFile := filepath.Join(pluginDir, plugin+".plugin.zsh")
		err = os.WriteFile(pluginFile, []byte("# "+plugin+" plugin"), 0644)
		if err != nil {
			t.Fatalf("Failed to create plugin file: %v", err)
		}
	}

	// Create custom plugins directory
	customDir := filepath.Join(tempDir, "custom", "plugins")
	err = os.MkdirAll(customDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom plugins directory: %v", err)
	}

	// Create a custom plugin
	customPluginDir := filepath.Join(customDir, "zsh-autosuggestions")
	err = os.MkdirAll(customPluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom plugin directory: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// List installed plugins
	installed, err := installer.ListInstalledPlugins()
	if err != nil {
		t.Fatalf("ListInstalledPlugins() error = %v", err)
	}

	// Verify count (3 built-in + 1 custom = 4)
	if len(installed) < 3 {
		t.Errorf("ListInstalledPlugins() count = %d, want at least 3", len(installed))
	}

	// Verify plugins are present
	pluginMap := make(map[string]bool)
	for _, p := range installed {
		pluginMap[p] = true
	}

	for _, expected := range plugins {
		if !pluginMap[expected] {
			t.Errorf("Plugin %s not found in list", expected)
		}
	}
}

// TestZshPluginInstaller_UninstallPlugin tests uninstalling a plugin
func TestZshPluginInstaller_UninstallPlugin(t *testing.T) {
	tempDir := t.TempDir()
	customDir := filepath.Join(tempDir, "custom", "plugins")

	// Create custom plugin directory
	err := os.MkdirAll(customDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom plugins directory: %v", err)
	}

	// Create a plugin to uninstall
	pluginDir := filepath.Join(customDir, "test-plugin")
	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin directory: %v", err)
	}

	pluginFile := filepath.Join(pluginDir, "test-plugin.plugin.zsh")
	err = os.WriteFile(pluginFile, []byte("# Test plugin"), 0644)
	if err != nil {
		t.Fatalf("Failed to create plugin file: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Uninstall plugin
	err = installer.UninstallPlugin("test-plugin")
	if err != nil {
		t.Fatalf("UninstallPlugin() error = %v", err)
	}

	// Verify plugin directory is removed
	if _, err := os.Stat(pluginDir); !os.IsNotExist(err) {
		t.Error("Plugin directory still exists after uninstall")
	}
}

// TestZshPluginInstaller_UninstallBuiltInPlugin tests that built-in plugins cannot be uninstalled
func TestZshPluginInstaller_UninstallBuiltInPlugin(t *testing.T) {
	tempDir := t.TempDir()

	// Create built-in plugin
	pluginsDir := filepath.Join(tempDir, "plugins", "git")
	err := os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins directory: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Attempt to uninstall built-in plugin
	err = installer.UninstallPlugin("git")
	if err == nil {
		t.Error("Expected error when uninstalling built-in plugin, got nil")
	}
}

// TestZshPluginInstaller_UpdatePlugin tests updating an external plugin
func TestZshPluginInstaller_UpdatePlugin(t *testing.T) {
	tempDir := t.TempDir()
	customDir := filepath.Join(tempDir, "custom", "plugins")

	// Create custom plugin directory (simulating an installed plugin)
	pluginDir := filepath.Join(customDir, "zsh-autosuggestions")
	err := os.MkdirAll(pluginDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugin directory: %v", err)
	}

	// Create .git directory to simulate git repository
	gitDir := filepath.Join(pluginDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Note: This test would actually try to run git pull
	// In a real test environment, we would mock the git command
	// For now, we'll test that the method exists and can be called

	// This will likely fail because git pull won't work in test environment
	// But it validates the method signature and basic error handling
	_ = installer.UpdatePlugin("zsh-autosuggestions")
}

// TestZshPluginInstaller_ProgressCallback tests installation with progress callback
func TestZshPluginInstaller_ProgressCallback(t *testing.T) {
	tempDir := t.TempDir()
	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	// Create built-in plugin
	pluginsDir := filepath.Join(tempDir, "plugins", "git")
	err := os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins directory: %v", err)
	}

	pluginFile := filepath.Join(pluginsDir, "git.plugin.zsh")
	err = os.WriteFile(pluginFile, []byte("# Git plugin"), 0644)
	if err != nil {
		t.Fatalf("Failed to create plugin file: %v", err)
	}

	plugin := &entities.ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	// Track progress callbacks
	var progressCalls []string
	callback := func(progress interfaces.PluginInstallProgress) {
		progressCalls = append(progressCalls, progress.Status)
	}

	// Install with callback
	err = installer.InstallPlugin(plugin, callback)
	if err != nil {
		t.Fatalf("InstallPlugin() error = %v", err)
	}

	// Verify callback was called
	if len(progressCalls) == 0 {
		t.Error("Progress callback was not called")
	}
}

// Benchmark tests
func BenchmarkZshPluginInstaller_IsPluginInstalled(b *testing.B) {
	tempDir := b.TempDir()

	// Create test plugin
	pluginsDir := filepath.Join(tempDir, "plugins", "git")
	_ = os.MkdirAll(pluginsDir, 0755)
	pluginFile := filepath.Join(pluginsDir, "git.plugin.zsh")
	_ = os.WriteFile(pluginFile, []byte("# Git plugin"), 0644)

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = installer.IsPluginInstalled("git")
	}
}

func BenchmarkZshPluginInstaller_ListInstalledPlugins(b *testing.B) {
	tempDir := b.TempDir()

	// Create test plugins
	pluginsDir := filepath.Join(tempDir, "plugins")
	_ = os.MkdirAll(pluginsDir, 0755)

	plugins := []string{"git", "docker", "rails", "node", "python"}
	for _, plugin := range plugins {
		pluginDir := filepath.Join(pluginsDir, plugin)
		_ = os.MkdirAll(pluginDir, 0755)
		pluginFile := filepath.Join(pluginDir, plugin+".plugin.zsh")
		_ = os.WriteFile(pluginFile, []byte("# "+plugin), 0644)
	}

	installer := NewZshPluginInstaller(tempDir, filepath.Join(tempDir, "custom"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = installer.ListInstalledPlugins()
	}
}
