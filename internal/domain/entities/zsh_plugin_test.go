package entities

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// TestZshPlugin_Validate tests the validation logic for ZshPlugin
func TestZshPlugin_Validate(t *testing.T) {
	tests := []struct {
		name    string
		plugin  ZshPlugin
		wantErr bool
	}{
		{
			name: "valid built-in plugin",
			plugin: ZshPlugin{
				ID:          "git",
				Name:        "Git Plugin",
				Description: "Git aliases and functions",
				Source:      types.PluginSourceBuiltIn,
				CheckCmd:    "test -f $ZSH/plugins/git/git.plugin.zsh",
			},
			wantErr: false,
		},
		{
			name: "valid external plugin with repo",
			plugin: ZshPlugin{
				ID:          "zsh-autosuggestions",
				Name:        "Zsh Autosuggestions",
				Description: "Fish-like autosuggestions",
				Source:      types.PluginSourceExternal,
				RepoURL:     "https://github.com/zsh-users/zsh-autosuggestions",
				InstallCmd:  "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
				CheckCmd:    "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions",
			},
			wantErr: false,
		},
		{
			name: "valid custom plugin",
			plugin: ZshPlugin{
				ID:          "my-plugin",
				Name:        "My Plugin",
				Description: "Custom plugin",
				Source:      types.PluginSourceCustom,
				ConfigFile:  "~/.zsh/plugins/my-plugin.zsh",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			plugin: ZshPlugin{
				Name:   "Test",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
		{
			name: "missing Name",
			plugin: ZshPlugin{
				ID:     "test",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
		{
			name: "invalid Source",
			plugin: ZshPlugin{
				ID:     "test",
				Name:   "Test",
				Source: types.PluginSource("invalid"),
			},
			wantErr: true,
		},
		{
			name: "external plugin without RepoURL or InstallCmd",
			plugin: ZshPlugin{
				ID:     "test",
				Name:   "Test",
				Source: types.PluginSourceExternal,
			},
			wantErr: true,
		},
		{
			name: "empty ID",
			plugin: ZshPlugin{
				ID:     "",
				Name:   "Test",
				Source: types.PluginSourceBuiltIn,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.plugin.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZshPlugin.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestZshPlugin_IsBuiltIn tests checking if plugin is built-in
func TestZshPlugin_IsBuiltIn(t *testing.T) {
	builtIn := ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	if !builtIn.IsBuiltIn() {
		t.Error("Expected IsBuiltIn() to be true")
	}

	external := ZshPlugin{
		ID:     "ext",
		Name:   "External",
		Source: types.PluginSourceExternal,
	}

	if external.IsBuiltIn() {
		t.Error("Expected IsBuiltIn() to be false for external plugin")
	}
}

// TestZshPlugin_IsExternal tests checking if plugin is external
func TestZshPlugin_IsExternal(t *testing.T) {
	external := ZshPlugin{
		ID:     "ext",
		Name:   "External",
		Source: types.PluginSourceExternal,
	}

	if !external.IsExternal() {
		t.Error("Expected IsExternal() to be true")
	}

	builtIn := ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	if builtIn.IsExternal() {
		t.Error("Expected IsExternal() to be false for built-in plugin")
	}
}

// TestZshPlugin_IsCustom tests checking if plugin is custom
func TestZshPlugin_IsCustom(t *testing.T) {
	custom := ZshPlugin{
		ID:     "custom",
		Name:   "Custom",
		Source: types.PluginSourceCustom,
	}

	if !custom.IsCustom() {
		t.Error("Expected IsCustom() to be true")
	}

	builtIn := ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	if builtIn.IsCustom() {
		t.Error("Expected IsCustom() to be false for built-in plugin")
	}
}

// TestZshPlugin_NeedsInstallation tests checking if plugin needs installation
func TestZshPlugin_NeedsInstallation(t *testing.T) {
	tests := []struct {
		name   string
		plugin ZshPlugin
		want   bool
	}{
		{
			name: "built-in plugin does not need installation",
			plugin: ZshPlugin{
				ID:     "git",
				Source: types.PluginSourceBuiltIn,
			},
			want: false,
		},
		{
			name: "external plugin needs installation",
			plugin: ZshPlugin{
				ID:      "autosuggestions",
				Source:  types.PluginSourceExternal,
				RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
			},
			want: true,
		},
		{
			name: "custom plugin needs installation if has InstallCmd",
			plugin: ZshPlugin{
				ID:         "custom",
				Source:     types.PluginSourceCustom,
				InstallCmd: "cp my-plugin.zsh ~/.zsh/plugins/",
			},
			want: true,
		},
		{
			name: "custom plugin without InstallCmd doesn't need installation",
			plugin: ZshPlugin{
				ID:     "custom",
				Source: types.PluginSourceCustom,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plugin.NeedsInstallation(); got != tt.want {
				t.Errorf("ZshPlugin.NeedsInstallation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestZshPlugin_GetInstallCommand tests getting install command
func TestZshPlugin_GetInstallCommand(t *testing.T) {
	tests := []struct {
		name   string
		plugin ZshPlugin
		want   string
	}{
		{
			name: "plugin with explicit InstallCmd",
			plugin: ZshPlugin{
				ID:         "test",
				InstallCmd: "custom install command",
			},
			want: "custom install command",
		},
		{
			name: "external plugin without InstallCmd generates git clone",
			plugin: ZshPlugin{
				ID:      "autosuggestions",
				Source:  types.PluginSourceExternal,
				RepoURL: "https://github.com/zsh-users/zsh-autosuggestions",
			},
			want: "git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/autosuggestions",
		},
		{
			name: "built-in plugin returns empty",
			plugin: ZshPlugin{
				ID:     "git",
				Source: types.PluginSourceBuiltIn,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plugin.GetInstallCommand(); got != tt.want {
				t.Errorf("ZshPlugin.GetInstallCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestZshPlugin_GetCheckCommand tests getting check command
func TestZshPlugin_GetCheckCommand(t *testing.T) {
	tests := []struct {
		name   string
		plugin ZshPlugin
		want   string
	}{
		{
			name: "plugin with explicit CheckCmd",
			plugin: ZshPlugin{
				ID:       "test",
				CheckCmd: "test -f /path/to/plugin",
			},
			want: "test -f /path/to/plugin",
		},
		{
			name: "built-in plugin without CheckCmd generates default",
			plugin: ZshPlugin{
				ID:     "git",
				Source: types.PluginSourceBuiltIn,
			},
			want: "test -f $ZSH/plugins/git/git.plugin.zsh",
		},
		{
			name: "external plugin without CheckCmd generates default",
			plugin: ZshPlugin{
				ID:     "autosuggestions",
				Source: types.PluginSourceExternal,
			},
			want: "test -d ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/autosuggestions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plugin.GetCheckCommand(); got != tt.want {
				t.Errorf("ZshPlugin.GetCheckCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkZshPlugin_Validate(b *testing.B) {
	plugin := ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = plugin.Validate()
	}
}

func BenchmarkZshPlugin_IsBuiltIn(b *testing.B) {
	plugin := ZshPlugin{
		ID:     "git",
		Name:   "Git",
		Source: types.PluginSourceBuiltIn,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = plugin.IsBuiltIn()
	}
}
