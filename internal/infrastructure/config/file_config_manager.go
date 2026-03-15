package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/templates"
	"gopkg.in/yaml.v3"
)

// ZshrcTemplateData holds the data for rendering zshrc.tmpl
type ZshrcTemplateData struct {
	GeneratedAt          string
	HasOhMyZsh           bool
	HasPowerlevel10k     bool
	Plugins              []string
	HasNVM               bool
	HasBun               bool
	HasSDKMAN            bool
	HasPNPM              bool
	HasDeno              bool
	HasHomebrew          bool
	HasPyenv             bool
	HasCargo             bool
	IncludeProjectConfig bool
}

// AliasesTemplateData holds the data for rendering aliases.tmpl
type AliasesTemplateData struct {
	Name        string
	GeneratedAt string
	Aliases     map[string]string
}

// FunctionsTemplateData holds the data for rendering functions.tmpl
type FunctionsTemplateData struct {
	Name        string
	GeneratedAt string
	Functions   map[string]string
}

// FileConfigManager implements ConfigManager interface using file-based storage
type FileConfigManager struct {
	configDir string
	loader    *templates.TemplateLoader
}

// NewFileConfigManager creates a new file-based configuration manager
func NewFileConfigManager(configDir string) interfaces.ConfigManager {
	return &FileConfigManager{
		configDir: configDir,
	}
}

// NewFileConfigManagerWithTemplates creates a file-based configuration manager using a TemplateLoader
func NewFileConfigManagerWithTemplates(configDir string, loader *templates.TemplateLoader) interfaces.ConfigManager {
	return &FileConfigManager{
		configDir: configDir,
		loader:    loader,
	}
}

// SaveConfig saves a shell configuration to a YAML file
func (fcm *FileConfigManager) SaveConfig(config *entities.ShellConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(fcm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	configFile := filepath.Join(fcm.configDir, config.ID+".yaml")
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig loads a shell configuration from a YAML file
func (fcm *FileConfigManager) LoadConfig(configName string) (*entities.ShellConfig, error) {
	configFile := filepath.Join(fcm.configDir, configName+".yaml")

	// Read file
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config '%s' not found", configName)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal YAML
	var config entities.ShellConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// DeleteConfig removes a configuration file
func (fcm *FileConfigManager) DeleteConfig(configName string) error {
	configFile := filepath.Join(fcm.configDir, configName+".yaml")

	if err := os.Remove(configFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config '%s' not found", configName)
		}
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}

// ListConfigs returns all available configuration names
func (fcm *FileConfigManager) ListConfigs() ([]string, error) {
	// Read directory
	entries, err := os.ReadDir(fcm.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	// Collect .yaml files
	var configs []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") {
			// Remove .yaml extension to get config name
			configName := strings.TrimSuffix(name, ".yaml")
			configs = append(configs, configName)
		}
	}

	return configs, nil
}

// filterInstalledPlugins returns only plugin names that exist under ~/.oh-my-zsh (built-in or custom).
// Evita "[oh-my-zsh] plugin 'X' not found" ao escrever apenas plugins instalados.
func filterInstalledPlugins(homeDir string, plugins []string) []string {
	zsh := filepath.Join(homeDir, ".oh-my-zsh")
	builtIn := filepath.Join(zsh, "plugins")
	custom := filepath.Join(zsh, "custom", "plugins")
	var out []string
	for _, name := range plugins {
		builtInPath := filepath.Join(builtIn, name)
		customPath := filepath.Join(custom, name)
		if dirExists(builtInPath) || dirExists(customPath) {
			out = append(out, name)
		}
	}
	return out
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// GenerateZshrc generates .zshrc content based on selections
func (fcm *FileConfigManager) GenerateZshrc(selections interfaces.ConfigSelections) (string, error) {
	homeDir, _ := os.UserHomeDir()

	if fcm.loader != nil {
		plugins := filterInstalledPlugins(homeDir, selections.Plugins)
		data := ZshrcTemplateData{
			GeneratedAt:          time.Now().Format(time.RFC3339),
			HasOhMyZsh:           sliceContains(selections.CoreComponents, "oh-my-zsh"),
			HasPowerlevel10k:     sliceContains(selections.CoreComponents, "powerlevel10k"),
			Plugins:              plugins,
			HasNVM:               sliceContains(selections.Tools, "nvm"),
			HasBun:               sliceContains(selections.Tools, "bun"),
			HasSDKMAN:            sliceContains(selections.Tools, "sdkman"),
			HasPNPM:              sliceContains(selections.Tools, "pnpm"),
			HasDeno:              sliceContains(selections.Tools, "deno"),
			HasHomebrew:          sliceContains(selections.Tools, "homebrew"),
			HasPyenv:             sliceContains(selections.Tools, "pyenv"),
			HasCargo:             sliceContains(selections.Tools, "cargo"),
			IncludeProjectConfig: selections.IncludeProjectConfig,
		}
		return fcm.loader.RenderTemplate("zshrc.tmpl", data)
	}

	// Fallback: string builder (used when no template loader is configured)
	var builder strings.Builder

	// Powerlevel10k instant prompt DEVE ser o primeiro bloco (evita warning de console output)
	hasP10k := sliceContains(selections.CoreComponents, "powerlevel10k")
	if hasP10k {
		builder.WriteString("# Powerlevel10k instant prompt - must stay at top\n")
		builder.WriteString("typeset -g POWERLEVEL9K_INSTANT_PROMPT=quiet\n")
		builder.WriteString("if [[ -r \"${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh\" ]]; then\n")
		builder.WriteString("  source \"${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh\"\n")
		builder.WriteString("fi\n\n")
	}

	builder.WriteString("# Generated by Homestead\n")
	builder.WriteString(fmt.Sprintf("# Generated at: %s\n\n", time.Now().Format(time.RFC3339)))

	// Oh My Zsh configuration
	if sliceContains(selections.CoreComponents, "oh-my-zsh") {
		builder.WriteString("# Path to oh-my-zsh installation\n")
		builder.WriteString("export ZSH=\"$HOME/.oh-my-zsh\"\n\n")
	}

	// Theme (sem instant prompt aqui; já está no topo)
	if hasP10k {
		builder.WriteString("# Theme\n")
		builder.WriteString("ZSH_THEME=\"powerlevel10k/powerlevel10k\"\n\n")
	}

	// Plugins: só os que existem no disco (evita "plugin 'X' not found")
	installedPlugins := filterInstalledPlugins(homeDir, selections.Plugins)
	if len(installedPlugins) > 0 {
		builder.WriteString("# Plugins\n")
		builder.WriteString("plugins=(")
		for i, plugin := range installedPlugins {
			if i > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(plugin)
		}
		builder.WriteString(")\n\n")
	}

	// Source Oh My Zsh
	if sliceContains(selections.CoreComponents, "oh-my-zsh") {
		builder.WriteString("# Source Oh My Zsh\n")
		builder.WriteString("source $ZSH/oh-my-zsh.sh\n\n")
	}

	// Tools - NVM
	if sliceContains(selections.Tools, "nvm") {
		builder.WriteString("# NVM (Node Version Manager)\n")
		builder.WriteString("export NVM_DIR=\"$HOME/.nvm\"\n")
		builder.WriteString("[ -s \"$NVM_DIR/nvm.sh\" ] && \\. \"$NVM_DIR/nvm.sh\"\n")
		builder.WriteString("[ -s \"$NVM_DIR/bash_completion\" ] && \\. \"$NVM_DIR/bash_completion\"\n\n")
	}

	// Tools - Bun
	if sliceContains(selections.Tools, "bun") {
		builder.WriteString("# Bun\n")
		builder.WriteString("export BUN_INSTALL=\"$HOME/.bun\"\n")
		builder.WriteString("export PATH=\"$BUN_INSTALL/bin:$PATH\"\n\n")
	}

	// Project configs
	if selections.IncludeProjectConfig {
		builder.WriteString("# Project-specific configurations\n")
		builder.WriteString("if [[ -f ~/.zsh/projects/ivt.zsh ]]; then\n")
		builder.WriteString("  source ~/.zsh/projects/ivt.zsh\n")
		builder.WriteString("fi\n\n")
	}

	// Custom aliases
	if len(selections.CustomAliases) > 0 {
		builder.WriteString("# Custom Aliases\n")
		for alias, command := range selections.CustomAliases {
			builder.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, command))
		}
		builder.WriteString("\n")
	}

	// Custom environment variables
	if len(selections.CustomEnvVars) > 0 {
		builder.WriteString("# Custom Environment Variables\n")
		for key, value := range selections.CustomEnvVars {
			builder.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, value))
		}
		builder.WriteString("\n")
	}

	// General aliases and functions
	builder.WriteString("# Source general aliases and functions\n")
	builder.WriteString("if [[ -f ~/.zsh/general/aliases.zsh ]]; then\n")
	builder.WriteString("  source ~/.zsh/general/aliases.zsh\n")
	builder.WriteString("fi\n\n")
	builder.WriteString("if [[ -f ~/.zsh/general/functions.zsh ]]; then\n")
	builder.WriteString("  source ~/.zsh/general/functions.zsh\n")
	builder.WriteString("fi\n\n")

	// Powerlevel10k config
	if sliceContains(selections.CoreComponents, "powerlevel10k") {
		builder.WriteString("# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh\n")
		builder.WriteString("[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh\n")
	}

	return builder.String(), nil
}

// GenerateAliasesFile generates aliases.zsh content
func (fcm *FileConfigManager) GenerateAliasesFile(config *entities.ShellConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config cannot be nil")
	}

	if fcm.loader != nil {
		data := AliasesTemplateData{
			Name:        config.Name,
			GeneratedAt: time.Now().Format(time.RFC3339),
			Aliases:     config.Aliases,
		}
		return fcm.loader.RenderTemplate("aliases.tmpl", data)
	}

	var builder strings.Builder

	// Header
	builder.WriteString("# Aliases - Generated by Homestead\n")
	builder.WriteString(fmt.Sprintf("# Config: %s (%s)\n", config.Name, config.ID))
	builder.WriteString(fmt.Sprintf("# Generated at: %s\n\n", time.Now().Format(time.RFC3339)))

	// Write aliases
	if len(config.Aliases) > 0 {
		for alias, command := range config.Aliases {
			builder.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, command))
		}
	} else {
		builder.WriteString("# No aliases defined\n")
	}

	return builder.String(), nil
}

// GenerateFunctionsFile generates functions.zsh content
func (fcm *FileConfigManager) GenerateFunctionsFile(config *entities.ShellConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config cannot be nil")
	}

	if fcm.loader != nil {
		data := FunctionsTemplateData{
			Name:        config.Name,
			GeneratedAt: time.Now().Format(time.RFC3339),
			Functions:   config.Functions,
		}
		return fcm.loader.RenderTemplate("functions.tmpl", data)
	}

	var builder strings.Builder

	// Header
	builder.WriteString("# Functions - Generated by Homestead\n")
	builder.WriteString(fmt.Sprintf("# Config: %s (%s)\n", config.Name, config.ID))
	builder.WriteString(fmt.Sprintf("# Generated at: %s\n\n", time.Now().Format(time.RFC3339)))

	// Write functions
	if len(config.Functions) > 0 {
		for name, body := range config.Functions {
			builder.WriteString(fmt.Sprintf("%s() {\n", name))
			builder.WriteString(body)
			builder.WriteString("\n}\n\n")
		}
	} else {
		builder.WriteString("# No functions defined\n")
	}

	return builder.String(), nil
}

// BackupExistingConfig creates a backup of existing .zshrc and related files
func (fcm *FileConfigManager) BackupExistingConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(homeDir, ".zsh", "backups", timestamp)

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Files to backup
	filesToBackup := []string{
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".zsh", "general", "aliases.zsh"),
		filepath.Join(homeDir, ".zsh", "general", "functions.zsh"),
	}

	// Backup each file if it exists
	for _, srcFile := range filesToBackup {
		if _, err := os.Stat(srcFile); err == nil {
			// File exists, backup it
			fileName := filepath.Base(srcFile)
			destFile := filepath.Join(backupDir, fileName)

			data, err := os.ReadFile(srcFile)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", srcFile, err)
			}

			if err := os.WriteFile(destFile, data, 0644); err != nil {
				return fmt.Errorf("failed to write backup %s: %w", destFile, err)
			}
		}
	}

	return nil
}

// ApplyConfig writes the generated configs to the filesystem
func (fcm *FileConfigManager) ApplyConfig(selections interfaces.ConfigSelections) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Backup existing config (best-effort: do not block apply if backup fails)
	_ = fcm.BackupExistingConfig()

	// Generate .zshrc
	zshrcContent, err := fcm.GenerateZshrc(selections)
	if err != nil {
		return fmt.Errorf("failed to generate .zshrc: %w", err)
	}

	// Write .zshrc
	zshrcPath := filepath.Join(homeDir, ".zshrc")
	if err := os.WriteFile(zshrcPath, []byte(zshrcContent), 0644); err != nil {
		return fmt.Errorf("failed to write .zshrc: %w", err)
	}

	// Create .zsh/general directory
	generalDir := filepath.Join(homeDir, ".zsh", "general")
	if err := os.MkdirAll(generalDir, 0755); err != nil {
		return fmt.Errorf("failed to create .zsh/general directory: %w", err)
	}

	// Create placeholder files for aliases and functions
	aliasesPath := filepath.Join(generalDir, "aliases.zsh")
	if _, err := os.Stat(aliasesPath); os.IsNotExist(err) {
		placeholder := "# General aliases - Add your aliases here\n"
		if err := os.WriteFile(aliasesPath, []byte(placeholder), 0644); err != nil {
			return fmt.Errorf("failed to write aliases.zsh: %w", err)
		}
	}

	functionsPath := filepath.Join(generalDir, "functions.zsh")
	if _, err := os.Stat(functionsPath); os.IsNotExist(err) {
		placeholder := "# General functions - Add your functions here\n"
		if err := os.WriteFile(functionsPath, []byte(placeholder), 0644); err != nil {
			return fmt.Errorf("failed to write functions.zsh: %w", err)
		}
	}

	return nil
}

// Helper function to check if a slice contains a string
func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
