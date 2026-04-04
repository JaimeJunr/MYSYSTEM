package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
	"github.com/JaimeJunr/Homestead/internal/tui"
)

func main() {
	// Dependency Injection (Manual Wiring)

	// Infrastructure layer - Scripts
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExecutor := executor.NewBashExecutor()

	// Infrastructure layer - Packages/Installers
	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstaller()

	// Application layer
	scriptService := services.NewScriptService(scriptRepo, scriptExecutor)
	installerService := services.NewInstallerService(packageRepo, packageInstaller)
	if pkgs, _, err := catalog.ReadAndParseCacheFile(catalog.CacheFilePath()); err == nil {
		_ = installerService.MergePackages(pkgs)
	}

	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)
	repoService, _ := services.NewRepoService("~/.config/" + services.DefaultRepoDirName)

	catalogURL := catalog.ResolveCatalogURL()

	// Presentation layer
	model := tui.NewModel(scriptService, installerService, configService, repoService, catalogURL)

	// Create the TUI program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar Homestead: %v\n", err)
		os.Exit(1)
	}
}
