package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/homesteadcli"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/catalog"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/config"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/executor"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/installer"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/repository"
	"github.com/JaimeJunr/Homestead/internal/tui"
)

// version is set by release builds (-ldflags "-X main.version=...").
var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			homesteadcli.PrintHelp(os.Stdout)
			return
		case "export-profile":
			os.Exit(homesteadcli.RunExportProfile(os.Args[2:], version, os.Stdout, os.Stderr))
		case "shell-init":
			os.Exit(homesteadcli.RunShellInit(os.Args[2:], os.Stderr))
		}
	}

	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(version)
		return
	}

	// Dependency Injection (Manual Wiring)

	prefsPath, err := preferences.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Homestead: preferências: %v\n", err)
		os.Exit(1)
	}
	prefs, err := preferences.Load(prefsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Homestead: carregar preferências: %v\n", err)
		os.Exit(1)
	}

	profilePath, err := profilestate.DefaultPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Homestead: perfil: %v\n", err)
		os.Exit(1)
	}
	profState, err := profilestate.Load(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Homestead: carregar perfil: %v\n", err)
		os.Exit(1)
	}
	profPtr := &profState

	catalogEnvSet := strings.TrimSpace(os.Getenv("HOMESTEAD_CATALOG_URL")) != ""
	catalogURL := catalog.EffectiveCatalogURL(prefs.CatalogURL)

	// Infrastructure layer - Scripts
	scriptRepo := repository.NewInMemoryScriptRepository()
	scriptExecutor := executor.NewBashExecutorWithRoot(prefs.ScriptRoot)

	// Infrastructure layer - Packages/Installers
	packageRepo := repository.NewInMemoryPackageRepository()
	packageInstaller := installer.NewDefaultPackageInstallerWithRoot(prefs.ScriptRoot)

	// Application layer
	scriptService := services.NewScriptService(scriptRepo, scriptExecutor)
	installerService := services.NewInstallerService(packageRepo, packageInstaller)
	if pkgs, _, err := catalog.ReadAndParseCacheFile(catalog.CacheFilePath()); err == nil {
		_ = installerService.MergePackages(pkgs)
	}

	configManager := config.NewFileConfigManager("")
	configService := services.NewConfigService(configManager)

	dotfilesPath := prefs.DotfilesRepo
	if strings.TrimSpace(dotfilesPath) == "" {
		dotfilesPath = preferences.DefaultDotfilesRepo()
	}
	repoService, err := services.NewRepoService(dotfilesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Homestead: repositório dotfiles: %v\n", err)
		os.Exit(1)
	}

	// Presentation layer
	model := tui.NewModel(scriptService, installerService, configService, repoService, catalogURL, prefs, prefsPath, catalogEnvSet, profPtr, profilePath)

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
