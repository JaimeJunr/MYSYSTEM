package services

import (
	"fmt"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ConfigService provides high-level configuration management operations
type ConfigService struct {
	configManager interfaces.ConfigManager
}

// NewConfigService creates a new configuration service
func NewConfigService(configManager interfaces.ConfigManager) *ConfigService {
	return &ConfigService{
		configManager: configManager,
	}
}

// CreateConfig creates and saves a new configuration
func (cs *ConfigService) CreateConfig(config *entities.ShellConfig) error {
	// Validate config
	if err := cs.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save config
	if err := cs.configManager.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

// GetConfig retrieves a configuration by ID
func (cs *ConfigService) GetConfig(configID string) (*entities.ShellConfig, error) {
	config, err := cs.configManager.LoadConfig(configID)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration '%s': %w", configID, err)
	}

	return config, nil
}

// ListConfigs returns all available configurations
func (cs *ConfigService) ListConfigs() ([]string, error) {
	configs, err := cs.configManager.ListConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}

	return configs, nil
}

// DeleteConfig deletes a configuration by ID
func (cs *ConfigService) DeleteConfig(configID string) error {
	if err := cs.configManager.DeleteConfig(configID); err != nil {
		return fmt.Errorf("failed to delete configuration '%s': %w", configID, err)
	}

	return nil
}

// ApplyConfig applies a configuration to the system
func (cs *ConfigService) ApplyConfig(selections interfaces.ConfigSelections) error {
	// Validate selections
	if len(selections.CoreComponents) == 0 {
		return fmt.Errorf("at least one core component must be selected")
	}

	// Apply configuration
	if err := cs.configManager.ApplyConfig(selections); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	return nil
}

// GenerateZshrc generates .zshrc content from selections
func (cs *ConfigService) GenerateZshrc(selections interfaces.ConfigSelections) (string, error) {
	zshrc, err := cs.configManager.GenerateZshrc(selections)
	if err != nil {
		return "", fmt.Errorf("failed to generate .zshrc: %w", err)
	}

	return zshrc, nil
}

// MergeConfigs merges general and project-specific configurations
// Project config takes precedence over general config
func (cs *ConfigService) MergeConfigs(general, project *entities.ShellConfig) *entities.ShellConfig {
	merged := &entities.ShellConfig{
		ID:    "merged",
		Name:  "Merged Configuration",
		Scope: types.ConfigScopeGeneral,
	}

	// Merge plugins (combine both, deduplicate)
	pluginMap := make(map[string]bool)
	if general != nil {
		for _, plugin := range general.Plugins {
			pluginMap[plugin] = true
		}
	}
	if project != nil {
		for _, plugin := range project.Plugins {
			pluginMap[plugin] = true
		}
	}
	for plugin := range pluginMap {
		merged.Plugins = append(merged.Plugins, plugin)
	}

	// Merge aliases (project overrides general)
	merged.Aliases = make(map[string]string)
	if general != nil {
		for key, value := range general.Aliases {
			merged.Aliases[key] = value
		}
	}
	if project != nil {
		for key, value := range project.Aliases {
			merged.Aliases[key] = value // Override
		}
	}

	// Merge functions (project overrides general)
	merged.Functions = make(map[string]string)
	if general != nil {
		for key, value := range general.Functions {
			merged.Functions[key] = value
		}
	}
	if project != nil {
		for key, value := range project.Functions {
			merged.Functions[key] = value // Override
		}
	}

	// Merge environment variables (project overrides general)
	merged.EnvVars = make(map[string]string)
	if general != nil {
		for key, value := range general.EnvVars {
			merged.EnvVars[key] = value
		}
	}
	if project != nil {
		for key, value := range project.EnvVars {
			merged.EnvVars[key] = value // Override
		}
	}

	// Merge sourced files (combine both)
	if general != nil {
		merged.SourcedFiles = append(merged.SourcedFiles, general.SourcedFiles...)
	}
	if project != nil {
		merged.SourcedFiles = append(merged.SourcedFiles, project.SourcedFiles...)
	}

	return merged
}

// ValidateConfig validates a shell configuration
func (cs *ConfigService) ValidateConfig(config *entities.ShellConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Use entity's built-in validation
	if err := config.Validate(); err != nil {
		return err
	}

	return nil
}

// BackupCurrentConfig backs up the current shell configuration
func (cs *ConfigService) BackupCurrentConfig() error {
	if err := cs.configManager.BackupExistingConfig(); err != nil {
		return fmt.Errorf("failed to backup current configuration: %w", err)
	}

	return nil
}

// GetConfigsByScope returns configurations filtered by scope
func (cs *ConfigService) GetConfigsByScope(scope types.ConfigScope) ([]*entities.ShellConfig, error) {
	// Get all config names
	configNames, err := cs.configManager.ListConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}

	// Load and filter by scope
	var configs []*entities.ShellConfig
	for _, name := range configNames {
		config, err := cs.configManager.LoadConfig(name)
		if err != nil {
			// Skip configs that fail to load
			continue
		}

		if config.Scope == scope {
			configs = append(configs, config)
		}
	}

	return configs, nil
}

// UpdateConfig updates an existing configuration
func (cs *ConfigService) UpdateConfig(config *entities.ShellConfig) error {
	// Validate config
	if err := cs.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Update timestamps
	config.Touch()

	// Save config
	if err := cs.configManager.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}

	return nil
}

// ExportConfig exports a configuration to a specific location
func (cs *ConfigService) ExportConfig(configID, exportPath string) error {
	// Load config
	config, err := cs.configManager.LoadConfig(configID)
	if err != nil {
		return fmt.Errorf("failed to load configuration '%s': %w", configID, err)
	}

	// For now, just validate that we can load it
	// In a full implementation, we would write to exportPath
	if config == nil {
		return fmt.Errorf("loaded configuration is nil")
	}

	return nil
}
