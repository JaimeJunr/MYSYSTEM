package entities

import (
	"time"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ShellConfig represents a shell configuration (zshrc, bashrc, etc)
type ShellConfig struct {
	ID   string
	Name string
	Scope types.ConfigScope // General, Project, Tool

	// Content
	Plugins  []string          // Lista de plugins habilitados
	Aliases  map[string]string // name -> command
	Functions map[string]string // name -> body
	EnvVars  map[string]string // VAR -> value

	// Files
	SourcedFiles []string // Arquivos para source

	// Metadata
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Validate checks if the configuration is valid
func (sc *ShellConfig) Validate() error {
	if sc.ID == "" {
		return types.ErrInvalidInput
	}
	if sc.Name == "" {
		return types.ErrInvalidInput
	}
	if !sc.Scope.IsValid() {
		return types.ErrInvalidInput
	}
	return nil
}

// AddPlugin adds a plugin to the configuration (no duplicates)
func (sc *ShellConfig) AddPlugin(plugin string) {
	if sc.Plugins == nil {
		sc.Plugins = []string{}
	}

	// Check if already exists
	for _, p := range sc.Plugins {
		if p == plugin {
			return
		}
	}

	sc.Plugins = append(sc.Plugins, plugin)
	sc.Touch()
}

// RemovePlugin removes a plugin from the configuration
func (sc *ShellConfig) RemovePlugin(plugin string) {
	if sc.Plugins == nil {
		return
	}

	newPlugins := []string{}
	for _, p := range sc.Plugins {
		if p != plugin {
			newPlugins = append(newPlugins, p)
		}
	}

	sc.Plugins = newPlugins
	sc.Touch()
}

// HasPlugin checks if a plugin is enabled
func (sc *ShellConfig) HasPlugin(plugin string) bool {
	for _, p := range sc.Plugins {
		if p == plugin {
			return true
		}
	}
	return false
}

// AddAlias adds or updates an alias
func (sc *ShellConfig) AddAlias(name, command string) {
	if sc.Aliases == nil {
		sc.Aliases = make(map[string]string)
	}
	sc.Aliases[name] = command
	sc.Touch()
}

// AddFunction adds or updates a function
func (sc *ShellConfig) AddFunction(name, body string) {
	if sc.Functions == nil {
		sc.Functions = make(map[string]string)
	}
	sc.Functions[name] = body
	sc.Touch()
}

// AddEnvVar adds or updates an environment variable
func (sc *ShellConfig) AddEnvVar(name, value string) {
	if sc.EnvVars == nil {
		sc.EnvVars = make(map[string]string)
	}
	sc.EnvVars[name] = value
	sc.Touch()
}

// IsGeneral returns true if this is a general configuration
func (sc *ShellConfig) IsGeneral() bool {
	return sc.Scope == types.ConfigScopeGeneral
}

// IsProject returns true if this is a project-specific configuration
func (sc *ShellConfig) IsProject() bool {
	return sc.Scope == types.ConfigScopeProject
}

// IsTool returns true if this is a tool-specific configuration
func (sc *ShellConfig) IsTool() bool {
	return sc.Scope == types.ConfigScopeTool
}

// Touch updates the ModifiedAt timestamp
func (sc *ShellConfig) Touch() {
	sc.ModifiedAt = time.Now()
}
