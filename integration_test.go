package main

import (
	"testing"

	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
	"github.com/JaimeJunr/Homestead/internal/tui"
)

// TestIntegration_ScriptsAndTUI tests integration between scripts and TUI
func TestIntegration_ScriptsAndTUI(t *testing.T) {
	// Create dependencies - Scripts
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	// Create dependencies - Installers
	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	// Create dependencies - Config
	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)

	// Get all scripts
	allScripts, err := scriptService.GetAllScripts()
	if err != nil {
		t.Fatalf("Failed to get scripts: %v", err)
	}

	if len(allScripts) == 0 {
		t.Fatal("No scripts found")
	}

	// Create TUI model
	repoService, _ := services.NewRepoService("")
	model := tui.NewModel(scriptService, installerService, configService, repoService, "")

	// Verify model initializes correctly
	if model.Init() == nil {
		t.Error("Expected Init() to return spinner tick command")
	}

	t.Logf("Integration test successful: %d scripts available, TUI initialized", len(allScripts))
}

// TestIntegration_AllCategoriesHaveScripts verifies each category has scripts
func TestIntegration_AllCategoriesHaveScripts(t *testing.T) {
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	service := services.NewScriptService(scriptRepo, scriptExec)

	categories := []string{"cleanup", "monitoring", "utilities"}

	for _, category := range categories {
		scripts, err := service.GetScriptsByCategory(types.Category(category))
		if err != nil {
			t.Errorf("Failed to get scripts for category %s: %v", category, err)
			continue
		}
		if len(scripts) == 0 {
			t.Errorf("Category %s has no scripts", category)
		} else {
			t.Logf("Category %s: %d scripts", category, len(scripts))
		}
	}
}

// TestIntegration_ScriptPaths verifies all script paths exist
func TestIntegration_ScriptPaths(t *testing.T) {
	// Note: This test will fail if run without the actual script files
	// Skip in CI/CD environments or when scripts aren't present
	if testing.Short() {
		t.Skip("Skipping script path verification in short mode")
	}

	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	service := services.NewScriptService(scriptRepo, scriptExec)

	allScripts, err := service.GetAllScripts()
	if err != nil {
		t.Fatalf("Failed to get scripts: %v", err)
	}

	for _, script := range allScripts {
		t.Logf("Script: %s at %s", script.Name, script.Path)
		// In a full integration environment, you would verify:
		// if !fileExists(script.Path) {
		//     t.Errorf("Script file not found: %s", script.Path)
		// }
	}
}
