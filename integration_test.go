package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
	"github.com/JaimeJunr/Homestead/internal/tui"
)

func skipIfShortIntegration(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration tests are skipped with go test -short")
	}
}

// TestIntegration_ScriptsAndTUI tests integration between scripts and TUI
func TestIntegration_ScriptsAndTUI(t *testing.T) {
	skipIfShortIntegration(t)

	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	scriptService := services.NewScriptService(scriptRepo, scriptExec)

	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()
	installerService := services.NewInstallerService(packageRepo, packageInstaller)

	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)

	allScripts, err := scriptService.GetAllScripts()
	if err != nil {
		t.Fatalf("Failed to get scripts: %v", err)
	}

	if len(allScripts) == 0 {
		t.Fatal("No scripts found")
	}

	repoService, _ := services.NewRepoService("")
	prefs := preferences.DefaultPreferences()
	prof := &profilestate.State{}
	model := tui.NewModel(scriptService, installerService, configService, repoService, "", prefs, "", false, prof, "")

	if model.Init() == nil {
		t.Error("Expected Init() to return a non-nil command batch")
	}

	t.Logf("Integration test successful: %d scripts available, TUI initialized", len(allScripts))
}

// TestIntegration_AllCategoriesHaveScripts verifies each category has scripts
func TestIntegration_AllCategoriesHaveScripts(t *testing.T) {
	skipIfShortIntegration(t)

	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	service := services.NewScriptService(scriptRepo, scriptExec)

	categories := []string{"cleanup", "monitoring", "checkup", "utilities"}

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

// TestIntegration_ScriptPaths verifies bash script paths exist on disk (repo root = test cwd).
func TestIntegration_ScriptPaths(t *testing.T) {
	skipIfShortIntegration(t)

	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExec := executor.NewBashExecutor()
	service := services.NewScriptService(scriptRepo, scriptExec)

	allScripts, err := service.GetAllScripts()
	if err != nil {
		t.Fatalf("Failed to get scripts: %v", err)
	}

	for _, script := range allScripts {
		if script.Path == "" {
			continue
		}
		full := filepath.Join(root, filepath.FromSlash(script.Path))
		if _, err := os.Stat(full); err != nil {
			t.Errorf("script %s (%s): %v", script.ID, full, err)
		}
	}
}
