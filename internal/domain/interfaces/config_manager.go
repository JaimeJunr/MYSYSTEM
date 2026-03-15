package interfaces

import "github.com/JaimeJunr/Homestead/internal/domain/entities"

// ConfigSelections represents the user's selections during wizard
type ConfigSelections struct {
	CoreComponents       []string          // zsh, oh-my-zsh, powerlevel10k
	Plugins              []string          // git, docker, rails, etc
	Tools                []string          // nvm, bun, homebrew, etc
	IncludeProjectConfig bool              // Include IVT/Performit configs
	CustomAliases        map[string]string // Additional custom aliases
	CustomFunctions      map[string]string // Additional custom functions
	CustomEnvVars        map[string]string // Additional environment variables
}

// ConfigManager defines the interface for managing shell configurations
type ConfigManager interface {
	// SaveConfig saves a shell configuration to persistent storage
	SaveConfig(config *entities.ShellConfig) error

	// LoadConfig loads a shell configuration by name
	LoadConfig(configName string) (*entities.ShellConfig, error)

	// DeleteConfig removes a configuration
	DeleteConfig(configName string) error

	// ListConfigs returns all available configuration names
	ListConfigs() ([]string, error)

	// GenerateZshrc generates a .zshrc file content based on selections
	GenerateZshrc(selections ConfigSelections) (string, error)

	// GenerateAliasesFile generates an aliases.zsh file content
	GenerateAliasesFile(config *entities.ShellConfig) (string, error)

	// GenerateFunctionsFile generates a functions.zsh file content
	GenerateFunctionsFile(config *entities.ShellConfig) (string, error)

	// BackupExistingConfig creates a backup of existing .zshrc and related files
	BackupExistingConfig() error

	// ApplyConfig writes the generated configs to the filesystem
	ApplyConfig(selections ConfigSelections) error
}
