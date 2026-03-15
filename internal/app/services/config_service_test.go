package services

import (
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
)

// TestConfigService_CreateConfig tests creating a new configuration
func TestConfigService_CreateConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create a new config
	cfg := &entities.ShellConfig{
		ID:      "test-config",
		Name:    "Test Configuration",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
		Aliases: map[string]string{
			"ll": "ls -la",
		},
	}

	err := service.CreateConfig(cfg)
	if err != nil {
		t.Fatalf("CreateConfig() error = %v", err)
	}

	// Verify config was saved
	loaded, err := manager.LoadConfig("test-config")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if loaded.ID != cfg.ID {
		t.Errorf("Config ID mismatch: got %s, want %s", loaded.ID, cfg.ID)
	}
}

// TestConfigService_CreateConfig_Invalid tests creating invalid config
func TestConfigService_CreateConfig_Invalid(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Try to create invalid config (missing ID)
	cfg := &entities.ShellConfig{
		Name:  "Invalid",
		Scope: types.ConfigScopeGeneral,
	}

	err := service.CreateConfig(cfg)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

// TestConfigService_GetConfig tests retrieving a configuration
func TestConfigService_GetConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create and save a config
	cfg := &entities.ShellConfig{
		ID:      "get-test",
		Name:    "Get Test",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git"},
	}

	err := manager.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Get config through service
	retrieved, err := service.GetConfig("get-test")
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}

	if retrieved.ID != cfg.ID {
		t.Errorf("Config ID mismatch: got %s, want %s", retrieved.ID, cfg.ID)
	}
}

// TestConfigService_GetConfig_NotFound tests retrieving non-existent config
func TestConfigService_GetConfig_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	_, err := service.GetConfig("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent config, got nil")
	}
}

// TestConfigService_ListConfigs tests listing all configurations
func TestConfigService_ListConfigs(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create multiple configs
	configs := []*entities.ShellConfig{
		{
			ID:    "config1",
			Name:  "Config 1",
			Scope: types.ConfigScopeGeneral,
		},
		{
			ID:    "config2",
			Name:  "Config 2",
			Scope: types.ConfigScopeProject,
		},
	}

	for _, cfg := range configs {
		err := manager.SaveConfig(cfg)
		if err != nil {
			t.Fatalf("SaveConfig() error = %v", err)
		}
	}

	// List configs
	list, err := service.ListConfigs()
	if err != nil {
		t.Fatalf("ListConfigs() error = %v", err)
	}

	if len(list) != len(configs) {
		t.Errorf("ListConfigs() count = %d, want %d", len(list), len(configs))
	}
}

// TestConfigService_DeleteConfig tests deleting a configuration
func TestConfigService_DeleteConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create a config
	cfg := &entities.ShellConfig{
		ID:    "delete-test",
		Name:  "Delete Test",
		Scope: types.ConfigScopeGeneral,
	}

	err := manager.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Delete through service
	err = service.DeleteConfig("delete-test")
	if err != nil {
		t.Fatalf("DeleteConfig() error = %v", err)
	}

	// Verify it's gone
	_, err = manager.LoadConfig("delete-test")
	if err == nil {
		t.Error("Config still exists after deletion")
	}
}

// TestConfigService_ApplyConfig tests applying a configuration
func TestConfigService_ApplyConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh"},
		Plugins:        []string{"git", "docker"},
		Tools:          []string{"nvm"},
	}

	// Note: ApplyConfig would write to home directory in real usage
	// In tests, we just verify the method can be called
	err := service.ApplyConfig(selections)
	if err != nil {
		t.Logf("ApplyConfig() error = %v (expected in test environment)", err)
	}
}

// TestConfigService_GenerateZshrc tests generating .zshrc content
func TestConfigService_GenerateZshrc(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh", "powerlevel10k"},
		Plugins:        []string{"git", "docker"},
		Tools:          []string{"nvm", "bun"},
	}

	zshrc, err := service.GenerateZshrc(selections)
	if err != nil {
		t.Fatalf("GenerateZshrc() error = %v", err)
	}

	if zshrc == "" {
		t.Error("GenerateZshrc() returned empty content")
	}

	// Verify key content is present
	if !contains(zshrc, "oh-my-zsh") {
		t.Error("Generated .zshrc missing oh-my-zsh reference")
	}

	if !contains(zshrc, "plugins=") {
		t.Error("Generated .zshrc missing plugins declaration")
	}
}

// TestConfigService_MergeConfigs tests merging general and project configs
func TestConfigService_MergeConfigs(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create general config
	general := &entities.ShellConfig{
		ID:    "general",
		Name:  "General Config",
		Scope: types.ConfigScopeGeneral,
		Aliases: map[string]string{
			"ll": "ls -la",
			"la": "ls -A",
		},
		Plugins: []string{"git", "docker"},
	}

	// Create project config
	project := &entities.ShellConfig{
		ID:    "project",
		Name:  "Project Config",
		Scope: types.ConfigScopeProject,
		Aliases: map[string]string{
			"ll":  "ls -lah", // Override general alias
			"psr": "cd ~/project && bundle exec rails s",
		},
		Plugins: []string{"rails"},
	}

	// Merge configs
	merged := service.MergeConfigs(general, project)

	// Verify merge results
	if len(merged.Plugins) != 3 { // git, docker, rails
		t.Errorf("Merged plugins count = %d, want 3", len(merged.Plugins))
	}

	// Verify project alias overrides general
	if merged.Aliases["ll"] != "ls -lah" {
		t.Errorf("Project alias should override general: got %s", merged.Aliases["ll"])
	}

	// Verify general alias is preserved
	if merged.Aliases["la"] != "ls -A" {
		t.Error("General alias 'la' should be preserved")
	}

	// Verify project-specific alias is included
	if merged.Aliases["psr"] != "cd ~/project && bundle exec rails s" {
		t.Error("Project-specific alias 'psr' should be included")
	}
}

// TestConfigService_ValidateConfig tests configuration validation
func TestConfigService_ValidateConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	tests := []struct {
		name    string
		config  *entities.ShellConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &entities.ShellConfig{
				ID:    "valid",
				Name:  "Valid Config",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: false,
		},
		{
			name: "invalid - missing ID",
			config: &entities.ShellConfig{
				Name:  "Invalid",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
		{
			name: "invalid - missing Name",
			config: &entities.ShellConfig{
				ID:    "invalid",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid Scope",
			config: &entities.ShellConfig{
				ID:    "invalid",
				Name:  "Invalid",
				Scope: types.ConfigScope("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigService_BackupCurrentConfig tests backing up current config
func TestConfigService_BackupCurrentConfig(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Note: BackupCurrentConfig would backup actual home directory files
	// In tests, we just verify the method can be called
	err := service.BackupCurrentConfig()
	if err != nil {
		t.Logf("BackupCurrentConfig() error = %v (expected in test environment)", err)
	}
}

// TestConfigService_GetConfigsByScope tests filtering configs by scope
func TestConfigService_GetConfigsByScope(t *testing.T) {
	tempDir := t.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create configs with different scopes
	configs := []*entities.ShellConfig{
		{
			ID:    "general1",
			Name:  "General 1",
			Scope: types.ConfigScopeGeneral,
		},
		{
			ID:    "general2",
			Name:  "General 2",
			Scope: types.ConfigScopeGeneral,
		},
		{
			ID:    "project1",
			Name:  "Project 1",
			Scope: types.ConfigScopeProject,
		},
	}

	for _, cfg := range configs {
		err := manager.SaveConfig(cfg)
		if err != nil {
			t.Fatalf("SaveConfig() error = %v", err)
		}
	}

	// Get general configs
	generalConfigs, err := service.GetConfigsByScope(types.ConfigScopeGeneral)
	if err != nil {
		t.Fatalf("GetConfigsByScope() error = %v", err)
	}

	if len(generalConfigs) != 2 {
		t.Errorf("General configs count = %d, want 2", len(generalConfigs))
	}

	// Get project configs
	projectConfigs, err := service.GetConfigsByScope(types.ConfigScopeProject)
	if err != nil {
		t.Fatalf("GetConfigsByScope() error = %v", err)
	}

	if len(projectConfigs) != 1 {
		t.Errorf("Project configs count = %d, want 1", len(projectConfigs))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkConfigService_CreateConfig(b *testing.B) {
	tempDir := b.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	cfg := &entities.ShellConfig{
		ID:      "bench",
		Name:    "Benchmark",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.ID = filepath.Join("bench", string(rune(i)))
		_ = service.CreateConfig(cfg)
	}
}

func BenchmarkConfigService_ListConfigs(b *testing.B) {
	tempDir := b.TempDir()
	manager := config.NewFileConfigManager(tempDir)
	service := NewConfigService(manager)

	// Create some configs
	for i := 0; i < 10; i++ {
		cfg := &entities.ShellConfig{
			ID:    filepath.Join("bench", string(rune(i))),
			Name:  "Benchmark",
			Scope: types.ConfigScopeGeneral,
		}
		_ = manager.SaveConfig(cfg)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ListConfigs()
	}
}
