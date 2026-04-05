package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/profilestate"
	"github.com/JaimeJunr/Homestead/internal/tui/items"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

func (m *Model) loadScripts(category types.Category) {
	m.loadScriptsWithParent(category, ViewMainMenu)
}

func (m *Model) loadScriptsWithParent(category types.Category, parent ViewState) {
	m.scriptListParent = parent
	m.scriptListCategory = category
	scripts, err := m.scriptService.GetScriptsByCategory(category)
	if err != nil {
		scripts = []entities.Script{}
	}

	rowItems := make([]list.Item, len(scripts))
	for i, script := range scripts {
		fav := m.profile != nil && profilestate.IsFavorite(m.profile, script.ID)
		rowItems[i] = items.ScriptItem{Script: script, Favorite: fav}
	}

	delegate := list.NewDefaultDelegate()
	m.scriptList = list.New(rowItems, delegate, m.width, m.height-theme.ListVerticalReserve())

	categoryNames := map[types.Category]string{
		types.CategoryCleanup:    "🧹 Limpeza do Sistema",
		types.CategoryMonitoring: "📊 Monitoramento (Go · ~3s)",
		types.CategoryCheckup:    "🩺 Manutenção / Check-up (leitura)",
		types.CategoryInstall:    "📦 Instaladores",
		types.CategoryUtilities:  "🧰 Utilitários",
	}

	m.scriptList.Title = categoryNames[category]
	m.scriptList.FilterInput.Prompt = "Filtrar: "
	m.scriptList.SetShowHelp(false)
	m.scriptList.SetShowStatusBar(true)
	m.scriptList.SetStatusBarItemName("script", "scripts")
	m.scriptList.SetFilteringEnabled(true)
}

func (m *Model) loadPackages(category types.PackageCategory) {
	packages, err := m.installerService.GetPackagesByCategory(category)
	if err != nil {
		packages = []entities.Package{}
	}

	t := theme.InstallerBreadcrumb(theme.InstallerPackageSectionTitle(category))
	m.setPackageList(packages, category, &t)
}

func (m *Model) loadPackagesFromCategories(categories []types.PackageCategory) {
	packages, err := m.installerService.GetPackagesByCategories(categories)
	if err != nil {
		packages = []entities.Package{}
	}

	seg := theme.InstallerPackageSectionTitle(categories[0])
	if len(categories) != 1 {
		seg = "Múltiplas categorias"
	}
	title := theme.InstallerBreadcrumb(seg)
	m.setPackageList(packages, categories[0], &title)
}

func (m *Model) setPackageList(packages []entities.Package, category types.PackageCategory, titleOverride *string) {
	rowItems := make([]list.Item, len(packages))
	for i, pkg := range packages {
		rowItems[i] = items.PackageItem{Pkg: pkg}
	}

	delegate := list.NewDefaultDelegate()
	m.packageList = list.New(rowItems, delegate, m.width, m.height-theme.ListVerticalReserve())

	if titleOverride != nil {
		m.packageList.Title = *titleOverride
	} else {
		m.packageList.Title = theme.InstallerBreadcrumb(theme.InstallerPackageSectionTitle(category))
	}
	m.packageList.FilterInput.Prompt = "Filtrar: "
	m.packageList.SetShowHelp(false)
	m.packageList.SetShowStatusBar(true)
	m.packageList.SetStatusBarItemName("pacote", "pacotes")
	m.packageList.SetFilteringEnabled(true)
}

func (m *Model) loadInstallerCategories() {
	rowItems := []list.Item{
		items.InstallerCategoryItem{
			Heading: "💻 IDEs & Dev (CLI)",
			Desc:    "VS Code, Cursor, Claude Code, Antigravity e afins",
			Categories: []types.PackageCategory{
				types.PackageCategoryIDE,
			},
		},
		items.InstallerCategoryItem{
			Heading: "📱 Aplicações",
			Desc:    "Google Chrome, Insomnia e outras aplicações",
			Categories: []types.PackageCategory{
				types.PackageCategoryApp,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🧰 Utilitários",
			Desc:    "VPN, Flatpak, periféricos e pacotes nativos (mesmo fluxo que os outros instaladores)",
			Categories: []types.PackageCategory{
				types.PackageCategoryUtilities,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🔧 Ferramentas de desenvolvimento",
			Desc:    "GitHub CLI (gh), NVM, Bun, pnpm, Deno e afins",
			Categories: []types.PackageCategory{
				types.PackageCategoryTool,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🐚 Shells alternativos",
			Desc:    "Fish Shell e outros",
			Categories: []types.PackageCategory{
				types.PackageCategoryShell,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🖥️ Emuladores de Terminal",
			Desc:    "WezTerm, Kitty, Alacritty, Zash Terminal, Warp, Wave e outros",
			Categories: []types.PackageCategory{
				types.PackageCategoryTerminal,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🐚 Componentes Core (Zsh)",
			Desc:    "Zsh, Oh My Zsh, Powerlevel10k",
			Categories: []types.PackageCategory{
				types.PackageCategoryZshCore,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🎮 Games",
			Desc:    "Prism Launcher, Lutris",
			Categories: []types.PackageCategory{
				types.PackageCategoryGames,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🤖 Integração com IA",
			Desc:    "ShellGPT, Fish-AI e assistentes por shell",
			Categories: []types.PackageCategory{
				types.PackageCategoryAI,
			},
		},
		items.InstallerCategoryItem{
			Heading: "🛡️ Administração de sistemas",
			Desc:    "Cockpit, Webmin, Topgrade, integração AD, clientes SSH e diagnóstico",
			Categories: []types.PackageCategory{
				types.PackageCategorySysAdmin,
			},
		},
	}

	delegate := list.NewDefaultDelegate()
	m.installerList = list.New(rowItems, delegate, m.width, m.height-theme.ListVerticalReserve())
	m.installerList.Title = "📦 Instaladores"
	m.installerList.SetShowStatusBar(false)
}

func (m *Model) reloadScriptList() {
	idx := m.scriptList.Index()
	m.loadScriptsWithParent(m.scriptListCategory, m.scriptListParent)
	items := m.scriptList.Items()
	if len(items) == 0 {
		return
	}
	if idx < 0 {
		idx = 0
	}
	if idx >= len(items) {
		idx = len(items) - 1
	}
	m.scriptList.Select(idx)
}
