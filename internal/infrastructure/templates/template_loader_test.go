package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// TestTemplateLoader_LoadTemplate tests loading a template from filesystem
func TestTemplateLoader_LoadTemplate(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create a test template
	templateContent := `# Test Template
Hello, {{.Name}}!`

	templatePath := filepath.Join(tempDir, "test.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Load template
	tmpl, err := loader.LoadTemplate("test.tmpl")
	if err != nil {
		t.Fatalf("LoadTemplate() error = %v", err)
	}

	if tmpl == nil {
		t.Error("LoadTemplate() returned nil template")
	}
}

// TestTemplateLoader_LoadTemplate_NotFound tests loading non-existent template
func TestTemplateLoader_LoadTemplate_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	_, err := loader.LoadTemplate("non-existent.tmpl")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

// TestTemplateLoader_RenderTemplate tests rendering a template with data
func TestTemplateLoader_RenderTemplate(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create a test template
	templateContent := `# Test Template
Name: {{.Name}}
Value: {{.Value}}`

	templatePath := filepath.Join(tempDir, "render.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Render template
	data := map[string]interface{}{
		"Name":  "TestConfig",
		"Value": 42,
	}

	result, err := loader.RenderTemplate("render.tmpl", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	if !strings.Contains(result, "Name: TestConfig") {
		t.Error("Rendered template missing expected Name value")
	}

	if !strings.Contains(result, "Value: 42") {
		t.Error("Rendered template missing expected Value")
	}
}

// TestTemplateLoader_RenderZshrc tests rendering .zshrc template
func TestTemplateLoader_RenderZshrc(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create .zshrc template
	zshrcTemplate := `# Generated .zshrc
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME="{{.Theme}}"

plugins=({{range $i, $p := .Plugins}}{{if $i}} {{end}}{{$p}}{{end}})

source $ZSH/oh-my-zsh.sh`

	templatePath := filepath.Join(tempDir, "zshrc.tmpl")
	err := os.WriteFile(templatePath, []byte(zshrcTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create zshrc template: %v", err)
	}

	// Render with data
	data := map[string]interface{}{
		"Theme":   "powerlevel10k/powerlevel10k",
		"Plugins": []string{"git", "docker", "rails"},
	}

	result, err := loader.RenderTemplate("zshrc.tmpl", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Verify content
	if !strings.Contains(result, "ZSH_THEME=\"powerlevel10k/powerlevel10k\"") {
		t.Error("Rendered .zshrc missing theme")
	}

	if !strings.Contains(result, "plugins=(git docker rails)") {
		t.Error("Rendered .zshrc missing or incorrect plugins")
	}
}

// TestTemplateLoader_RenderAliases tests rendering aliases template
func TestTemplateLoader_RenderAliases(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create aliases template
	aliasesTemplate := `# Aliases
{{range $name, $cmd := .Aliases}}alias {{$name}}='{{$cmd}}'
{{end}}`

	templatePath := filepath.Join(tempDir, "aliases.tmpl")
	err := os.WriteFile(templatePath, []byte(aliasesTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create aliases template: %v", err)
	}

	// Render with data
	data := map[string]interface{}{
		"Aliases": map[string]string{
			"ll":   "ls -la",
			"la":   "ls -A",
			"grep": "grep --color=auto",
		},
	}

	result, err := loader.RenderTemplate("aliases.tmpl", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Verify aliases are present
	if !strings.Contains(result, "alias ll='ls -la'") {
		t.Error("Rendered aliases missing 'll' alias")
	}

	if !strings.Contains(result, "alias la='ls -A'") {
		t.Error("Rendered aliases missing 'la' alias")
	}
}

// TestTemplateLoader_RenderFunctions tests rendering functions template
func TestTemplateLoader_RenderFunctions(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create functions template
	functionsTemplate := `# Functions
{{range $name, $body := .Functions}}{{$name}}() {
{{$body}}
}

{{end}}`

	templatePath := filepath.Join(tempDir, "functions.tmpl")
	err := os.WriteFile(templatePath, []byte(functionsTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create functions template: %v", err)
	}

	// Render with data
	data := map[string]interface{}{
		"Functions": map[string]string{
			"db-connect": `    local database="${1:-default}"
    mysql -u root -p "$database"`,
			"greet": `    echo "Hello, $1!"`,
		},
	}

	result, err := loader.RenderTemplate("functions.tmpl", data)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Verify functions are present
	if !strings.Contains(result, "db-connect()") {
		t.Error("Rendered functions missing 'db-connect' function")
	}

	if !strings.Contains(result, "greet()") {
		t.Error("Rendered functions missing 'greet' function")
	}
}

// TestTemplateLoader_RenderWithShellConfig tests rendering with ShellConfig entity
func TestTemplateLoader_RenderWithShellConfig(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create a template that uses ShellConfig
	configTemplate := `# Configuration: {{.Name}}
# Scope: {{.Scope}}

# Plugins
plugins=({{range $i, $p := .Plugins}}{{if $i}} {{end}}{{$p}}{{end}})

# Aliases
{{range $name, $cmd := .Aliases}}alias {{$name}}='{{$cmd}}'
{{end}}`

	templatePath := filepath.Join(tempDir, "config.tmpl")
	err := os.WriteFile(templatePath, []byte(configTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create config template: %v", err)
	}

	// Create ShellConfig
	config := &entities.ShellConfig{
		ID:      "test-config",
		Name:    "Test Configuration",
		Scope:   types.ConfigScopeGeneral,
		Plugins: []string{"git", "docker"},
		Aliases: map[string]string{
			"ll": "ls -la",
			"la": "ls -A",
		},
	}

	// Render with ShellConfig
	result, err := loader.RenderTemplate("config.tmpl", config)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Verify content
	if !strings.Contains(result, "Configuration: Test Configuration") {
		t.Error("Rendered template missing config name")
	}

	if !strings.Contains(result, "Scope: general") {
		t.Error("Rendered template missing config scope")
	}

	if !strings.Contains(result, "plugins=(git docker)") {
		t.Error("Rendered template missing or incorrect plugins")
	}
}

// TestTemplateLoader_RenderWithConfigSelections tests rendering with ConfigSelections
func TestTemplateLoader_RenderWithConfigSelections(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create a template that uses ConfigSelections
	selectionsTemplate := `# Core Components
{{range .CoreComponents}}- {{.}}
{{end}}
# Plugins
{{range .Plugins}}- {{.}}
{{end}}
# Tools
{{range .Tools}}- {{.}}
{{end}}`

	templatePath := filepath.Join(tempDir, "selections.tmpl")
	err := os.WriteFile(templatePath, []byte(selectionsTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create selections template: %v", err)
	}

	// Create ConfigSelections
	selections := interfaces.ConfigSelections{
		CoreComponents: []string{"zsh", "oh-my-zsh", "powerlevel10k"},
		Plugins:        []string{"git", "docker"},
		Tools:          []string{"nvm", "bun"},
	}

	// Render with ConfigSelections
	result, err := loader.RenderTemplate("selections.tmpl", selections)
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	// Verify content
	if !strings.Contains(result, "- zsh") {
		t.Error("Rendered template missing zsh component")
	}

	if !strings.Contains(result, "- git") {
		t.Error("Rendered template missing git plugin")
	}

	if !strings.Contains(result, "- nvm") {
		t.Error("Rendered template missing nvm tool")
	}
}

// TestTemplateLoader_InvalidTemplate tests rendering with invalid template syntax
func TestTemplateLoader_InvalidTemplate(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create invalid template
	invalidTemplate := `# Invalid Template
{{.Name`

	templatePath := filepath.Join(tempDir, "invalid.tmpl")
	err := os.WriteFile(templatePath, []byte(invalidTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid template: %v", err)
	}

	// Attempt to load and render
	_, err = loader.RenderTemplate("invalid.tmpl", map[string]interface{}{"Name": "Test"})
	if err == nil {
		t.Error("Expected error for invalid template syntax, got nil")
	}
}

// TestTemplateLoader_MissingData tests rendering with missing template data
func TestTemplateLoader_MissingData(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create template requiring data
	templateContent := `# Template
Name: {{.Name}}
Value: {{.RequiredField}}`

	templatePath := filepath.Join(tempDir, "required.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Render with incomplete data
	data := map[string]interface{}{
		"Name": "Test",
		// RequiredField is missing
	}

	result, err := loader.RenderTemplate("required.tmpl", data)
	// Template should still render, but with empty value
	if err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	if !strings.Contains(result, "Name: Test") {
		t.Error("Rendered template missing Name")
	}

	// RequiredField should render as empty (Go template default behavior)
	if !strings.Contains(result, "Value: <no value>") && !strings.Contains(result, "Value: ") {
		t.Logf("Template output: %s", result)
	}
}

// TestTemplateLoader_ListTemplates tests listing all available templates
func TestTemplateLoader_ListTemplates(t *testing.T) {
	tempDir := t.TempDir()
	loader := NewTemplateLoader(tempDir)

	// Create multiple templates
	templates := []string{"zshrc.tmpl", "aliases.tmpl", "functions.tmpl"}
	for _, name := range templates {
		templatePath := filepath.Join(tempDir, name)
		err := os.WriteFile(templatePath, []byte("# "+name), 0644)
		if err != nil {
			t.Fatalf("Failed to create template %s: %v", name, err)
		}
	}

	// List templates
	list, err := loader.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(list) != len(templates) {
		t.Errorf("ListTemplates() count = %d, want %d", len(list), len(templates))
	}

	// Verify all templates are listed
	templateMap := make(map[string]bool)
	for _, name := range list {
		templateMap[name] = true
	}

	for _, expected := range templates {
		if !templateMap[expected] {
			t.Errorf("Template %s not found in list", expected)
		}
	}
}

// Benchmark tests
func BenchmarkTemplateLoader_LoadTemplate(b *testing.B) {
	tempDir := b.TempDir()
	loader := NewTemplateLoader(tempDir)

	templateContent := `# Benchmark Template
Hello, {{.Name}}!`

	templatePath := filepath.Join(tempDir, "bench.tmpl")
	_ = os.WriteFile(templatePath, []byte(templateContent), 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.LoadTemplate("bench.tmpl")
	}
}

func BenchmarkTemplateLoader_RenderTemplate(b *testing.B) {
	tempDir := b.TempDir()
	loader := NewTemplateLoader(tempDir)

	templateContent := `# Benchmark Template
Name: {{.Name}}
plugins=({{range $i, $p := .Plugins}}{{if $i}} {{end}}{{$p}}{{end}})`

	templatePath := filepath.Join(tempDir, "bench.tmpl")
	_ = os.WriteFile(templatePath, []byte(templateContent), 0644)

	data := map[string]interface{}{
		"Name":    "BenchConfig",
		"Plugins": []string{"git", "docker", "rails", "nvm"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.RenderTemplate("bench.tmpl", data)
	}
}
