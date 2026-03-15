package entities

import (
	"testing"
	"time"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// TestShellConfig_Validate tests the validation logic for ShellConfig
func TestShellConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ShellConfig
		wantErr bool
	}{
		{
			name: "valid general config",
			config: ShellConfig{
				ID:    "general-config",
				Name:  "General Configuration",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: false,
		},
		{
			name: "valid project config with all fields",
			config: ShellConfig{
				ID:     "ivt-config",
				Name:   "IVT Project Configuration",
				Scope:  types.ConfigScopeProject,
				Plugins: []string{"git", "docker"},
				Aliases: map[string]string{
					"performit": "cd $PERFORMIT_DIR",
				},
				Functions: map[string]string{
					"db-connect": "mysql -u root -p",
				},
				EnvVars: map[string]string{
					"IVT_DIR": "$HOME/ivt",
				},
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			config: ShellConfig{
				Name:  "Test Config",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
		{
			name: "missing Name",
			config: ShellConfig{
				ID:    "test",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
		{
			name: "invalid Scope",
			config: ShellConfig{
				ID:    "test",
				Name:  "Test",
				Scope: types.ConfigScope("invalid"),
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			config: ShellConfig{
				ID:    "",
				Name:  "Test",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
		{
			name: "empty Name",
			config: ShellConfig{
				ID:    "test",
				Name:  "",
				Scope: types.ConfigScopeGeneral,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ShellConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestShellConfig_AddPlugin tests adding plugins to config
func TestShellConfig_AddPlugin(t *testing.T) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{},
	}

	// Add first plugin
	config.AddPlugin("git")
	if len(config.Plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(config.Plugins))
	}
	if config.Plugins[0] != "git" {
		t.Errorf("Expected plugin 'git', got '%s'", config.Plugins[0])
	}

	// Add second plugin
	config.AddPlugin("docker")
	if len(config.Plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(config.Plugins))
	}

	// Try to add duplicate - should not add
	config.AddPlugin("git")
	if len(config.Plugins) != 2 {
		t.Errorf("Expected 2 plugins (no duplicates), got %d", len(config.Plugins))
	}
}

// TestShellConfig_RemovePlugin tests removing plugins from config
func TestShellConfig_RemovePlugin(t *testing.T) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker", "rails"},
	}

	// Remove existing plugin
	config.RemovePlugin("docker")
	if len(config.Plugins) != 2 {
		t.Errorf("Expected 2 plugins after removal, got %d", len(config.Plugins))
	}

	// Verify docker is gone
	for _, p := range config.Plugins {
		if p == "docker" {
			t.Error("Plugin 'docker' should have been removed")
		}
	}

	// Remove non-existent plugin - should be no-op
	config.RemovePlugin("non-existent")
	if len(config.Plugins) != 2 {
		t.Errorf("Expected 2 plugins (no change), got %d", len(config.Plugins))
	}
}

// TestShellConfig_HasPlugin tests checking if plugin exists
func TestShellConfig_HasPlugin(t *testing.T) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
	}

	if !config.HasPlugin("git") {
		t.Error("Expected HasPlugin('git') to be true")
	}

	if !config.HasPlugin("docker") {
		t.Error("Expected HasPlugin('docker') to be true")
	}

	if config.HasPlugin("rails") {
		t.Error("Expected HasPlugin('rails') to be false")
	}
}

// TestShellConfig_AddAlias tests adding aliases
func TestShellConfig_AddAlias(t *testing.T) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		Aliases: make(map[string]string),
	}

	// Add new alias
	config.AddAlias("ll", "ls -la")
	if config.Aliases["ll"] != "ls -la" {
		t.Errorf("Expected alias 'll'='ls -la', got '%s'", config.Aliases["ll"])
	}

	// Update existing alias
	config.AddAlias("ll", "ls -lah")
	if config.Aliases["ll"] != "ls -lah" {
		t.Errorf("Expected updated alias 'll'='ls -lah', got '%s'", config.Aliases["ll"])
	}
}

// TestShellConfig_AddFunction tests adding functions
func TestShellConfig_AddFunction(t *testing.T) {
	config := ShellConfig{
		ID:        "test",
		Name:      "Test",
		Scope:     types.ConfigScopeGeneral,
		Functions: make(map[string]string),
	}

	funcBody := `echo "Hello from function"`

	// Add new function
	config.AddFunction("greet", funcBody)
	if config.Functions["greet"] != funcBody {
		t.Errorf("Expected function body, got '%s'", config.Functions["greet"])
	}
}

// TestShellConfig_AddEnvVar tests adding environment variables
func TestShellConfig_AddEnvVar(t *testing.T) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		EnvVars: make(map[string]string),
	}

	// Add new env var
	config.AddEnvVar("PATH", "/usr/local/bin:$PATH")
	if config.EnvVars["PATH"] != "/usr/local/bin:$PATH" {
		t.Errorf("Expected env var PATH, got '%s'", config.EnvVars["PATH"])
	}
}

// TestShellConfig_IsGeneral tests scope checking
func TestShellConfig_IsGeneral(t *testing.T) {
	general := ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeGeneral,
	}

	if !general.IsGeneral() {
		t.Error("Expected IsGeneral() to be true")
	}

	project := ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeProject,
	}

	if project.IsGeneral() {
		t.Error("Expected IsGeneral() to be false for project scope")
	}
}

// TestShellConfig_IsProject tests project scope checking
func TestShellConfig_IsProject(t *testing.T) {
	project := ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeProject,
	}

	if !project.IsProject() {
		t.Error("Expected IsProject() to be true")
	}

	general := ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeGeneral,
	}

	if general.IsProject() {
		t.Error("Expected IsProject() to be false for general scope")
	}
}

// TestShellConfig_Touch tests updating modified time
func TestShellConfig_Touch(t *testing.T) {
	config := ShellConfig{
		ID:         "test",
		Name:       "Test",
		Scope:      types.ConfigScopeGeneral,
		ModifiedAt: time.Now().Add(-1 * time.Hour),
	}

	oldTime := config.ModifiedAt
	time.Sleep(10 * time.Millisecond)

	config.Touch()

	if !config.ModifiedAt.After(oldTime) {
		t.Error("Expected ModifiedAt to be updated")
	}
}

// Benchmark tests
func BenchmarkShellConfig_Validate(b *testing.B) {
	config := ShellConfig{
		ID:    "test",
		Name:  "Test",
		Scope: types.ConfigScopeGeneral,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkShellConfig_AddPlugin(b *testing.B) {
	config := ShellConfig{
		ID:      "test",
		Name:    "Test",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.AddPlugin("git")
		config.Plugins = []string{} // Reset
	}
}
