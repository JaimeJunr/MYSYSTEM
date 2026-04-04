package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
	"github.com/JaimeJunr/Homestead/internal/monitoring"
)

// ViewState represents different views in the TUI
type ViewState int

const (
	ViewMainMenu ViewState = iota
	ViewScriptList
	ViewInstallerCategories
	ViewPackageList
	ViewConfirmation
	ViewScriptOutput
	ViewNativeMonitor
	ViewInstalling
	ViewZshWizard
	ViewZshApplying
	ViewZshRepoWizard
)

// Model is the main TUI model
type Model struct {
	scriptService    *services.ScriptService
	installerService *services.InstallerService
	configService    *services.ConfigService
	repoService      *services.RepoService
	state            ViewState
	mainMenu         list.Model
	scriptList       list.Model
	installerList    list.Model
	packageList      list.Model
	selectedMenu     int
	selectedItem     interface{} // Can be Script or Package
	confirmYes       bool        // true = yes selected, false = no selected
	confirmReturn    ViewState   // tela para voltar se cancelar a confirmação (lista de pacotes/scripts)
	confirmReturnOK  bool        // se false, cancelar volta ao menu principal
	width            int
	height           int
	err              error
	keyboardToast    string // feedback para o/c (abrir/copiar URL) sem mouse

	// Installation progress
	progress       progress.Model
	spinner        spinner.Model
	installStatus  string
	installMessage string
	installPercent float64
	canAbort       bool
	aborted        bool

	// Zsh plugins wizard (Plugins e temas Zsh)
	zshWizard *ZshWizardModel

	// Zsh repo wizard (Configurar Zsh - backup/migração via repositório)
	zshRepoWizard *ZshRepoModel

	// Zsh core: when true, "Plugins e temas Zsh" is shown in menu (oh-my-zsh installed)
	zshCoreInstalled bool
	zshCoreChecked   bool

	// Zsh apply feedback: phase "applying" | "success" | "error"
	zshApplyPhase string
	zshApplyError error

	// Script output (in-TUI); phase "running" | "done"
	scriptOutputView   viewport.Model
	scriptOutputPhase  string
	scriptOutputTitle  string
	scriptOutputErr    error

	// Monitores integrados (bateria / memória)
	nativeMonitorKind  string
	nativeBattery      *monitoring.BatterySnapshot
	nativeBatteryErr   error
	nativeMemory       *monitoring.MemorySnapshot
	nativeMemoryErr    error
}

// menuAction identifies the main menu action
const (
	menuActionCleanup      = "cleanup"
	menuActionMonitoring   = "monitoring"
	menuActionInstallers   = "installers"
	menuActionZshPlugins   = "zsh_plugins"   // Plugins e temas Zsh (wizard local)
	menuActionZshRepo   = "zsh_repo"   // Configurar Zsh (repo backup/migração)
	menuActionSettings  = "settings"
	menuActionQuit         = "quit"
)

// menuItem represents a menu option
type menuItem struct {
	title  string
	desc   string
	action string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

// scriptItem wraps a script for the list
type scriptItem struct {
	script entities.Script
}

func (i scriptItem) Title() string       { return i.script.Name }
func (i scriptItem) Description() string { return i.script.Description }
func (i scriptItem) FilterValue() string { return i.script.Name }

// packageItem wraps a package for the list
type packageItem struct {
	pkg entities.Package
}

func (i packageItem) Title() string       { return i.pkg.Name }
func (i packageItem) Description() string { return i.pkg.Description }
func (i packageItem) FilterValue() string { return i.pkg.Name }

// installerCategoryItem represents a logical group of packages (e.g. IDEs, Terminais)
type installerCategoryItem struct {
	title      string
	desc       string
	categories []types.PackageCategory
}

func (i installerCategoryItem) Title() string       { return i.title }
func (i installerCategoryItem) Description() string { return i.desc }
func (i installerCategoryItem) FilterValue() string { return i.title }

// progressMsg is sent when installation progress updates
type progressMsg interfaces.InstallProgress

// installCompleteMsg is sent when installation completes
type installCompleteMsg struct {
	err error
}

// zshCoreInstalledMsg is sent when the check for oh-my-zsh installation completes
type zshCoreInstalledMsg struct {
	installed bool
}

// zshApplyResultMsg is sent when ApplyConfig finishes
type zshApplyResultMsg struct {
	Err error
}

// zshApplyReturnToMenuMsg is sent after a delay to return to main menu
type zshApplyReturnToMenuMsg struct{}

// scriptCapturedMsg carries stdout/stderr after ExecuteScriptCapture
type scriptCapturedMsg struct {
	output string
	err    error
}

// scriptExecFinishedMsg is sent after tea.ExecProcess (sudo scripts)
type scriptExecFinishedMsg struct {
	err error
}

type urlActionDoneMsg struct {
	err  error
	verb string // "open" | "copy"
}

type clearKeyboardToastMsg struct{}

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

	yesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	noStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 1)

	// Script output screen — same accent (63) + title (205) as confirmações e listas
	scriptScreenOuterStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2)

	scriptScreenAccentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("63")).
				Bold(true)

	scriptLogAreaStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(0, 1).
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("252"))

	scriptScreenFooterBarStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Background(lipgloss.Color("235")).
					Padding(0, 1)
)

// getMainMenuItems returns menu items; "Plugins e temas Zsh" only when zsh core is installed; "Configurar Zsh" always
func getMainMenuItems(zshCoreInstalled bool) []list.Item {
	items := []list.Item{
		menuItem{title: "🧹 Limpeza do Sistema", desc: "Scripts de limpeza e manutenção", action: menuActionCleanup},
		menuItem{title: "📊 Monitoramento", desc: "Informações do sistema", action: menuActionMonitoring},
		menuItem{title: "📦 Instaladores", desc: "Instalar ferramentas e aplicações (IDEs, Zsh, Oh My Zsh, etc.)", action: menuActionInstallers},
	}
	if zshCoreInstalled {
		items = append(items, menuItem{title: "🔧 Plugins e temas Zsh", desc: "Plugins, temas e .zshrc local", action: menuActionZshPlugins})
	}
	items = append(items,
		menuItem{title: "⚙️  Configurar Zsh", desc: "Repositório de config: backup e migração entre máquinas", action: menuActionZshRepo},
		menuItem{title: "⚙️  Configurações", desc: "Configurar a ferramenta (em breve)", action: menuActionSettings},
		menuItem{title: "❌ Sair", desc: "Fechar Homestead", action: menuActionQuit},
	)
	return items
}

// NewModel creates the TUI model with dependencies injected
func NewModel(scriptService *services.ScriptService, installerService *services.InstallerService, configService *services.ConfigService, repoService *services.RepoService) Model {
	mainItems := getMainMenuItems(false) // will refresh when zsh core check completes
	mainList := list.New(mainItems, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "Homestead - Gerenciador de Sistema"
	mainList.SetShowStatusBar(false)
	mainList.SetFilteringEnabled(false)

	// Progress bar
	prog := progress.New(progress.WithDefaultGradient())

	// Spinner
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	return Model{
		scriptService:    scriptService,
		installerService: installerService,
		configService:    configService,
		repoService:      repoService,
		state:            ViewMainMenu,
		mainMenu:         mainList,
		progress:         prog,
		spinner:          spin,
		confirmYes:       false, // Default to "No"
	}
}

// checkZshCoreInstalled runs in a Cmd to detect if oh-my-zsh is installed (for menu)
func checkZshCoreInstalled(installerService *services.InstallerService) tea.Cmd {
	return func() tea.Msg {
		installed, _ := installerService.IsPackageInstalled("oh-my-zsh")
		return zshCoreInstalledMsg{installed: installed}
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, checkZshCoreInstalled(m.installerService))
}

// Update handles messages and updates state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if m.state == ViewScriptOutput && m.scriptOutputPhase == "done" {
			var vcmd tea.Cmd
			m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
			return m, vcmd
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height-4)
		if m.scriptList.Items() != nil {
			m.scriptList.SetSize(msg.Width, msg.Height-4)
		}
		if m.installerList.Items() != nil {
			m.installerList.SetSize(msg.Width, msg.Height-4)
		}
		if m.packageList.Items() != nil {
			m.packageList.SetSize(msg.Width, msg.Height-4)
		}
		if m.state == ViewScriptOutput {
			m.syncScriptOutputViewport()
		}
		return m, nil

	case progressMsg:
		m.installStatus = msg.Status
		m.installMessage = msg.Message
		m.installPercent = float64(msg.Progress) / 100.0
		m.canAbort = msg.CanAbort

		// If completed, return to main menu after a delay
		if msg.IsCompleted {
			return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return installCompleteMsg{err: msg.Error}
			})
		}
		return m, nil

	case installCompleteMsg:
		m.state = ViewMainMenu
		m.aborted = false
		// Re-check zsh core so "Configurar Zsh" appears if user just installed it
		return m, checkZshCoreInstalled(m.installerService)

	case zshCoreInstalledMsg:
		m.zshCoreChecked = true
		m.zshCoreInstalled = msg.installed
		m.mainMenu.SetItems(getMainMenuItems(m.zshCoreInstalled))
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.state == ViewScriptOutput {
			if m.scriptOutputPhase == "done" {
				switch msg.String() {
				case "enter", "esc", "q":
					m.state = m.confirmReturn
					m.scriptOutputPhase = ""
					m.scriptOutputTitle = ""
					m.scriptOutputErr = nil
					return m, nil
				}
				var vcmd tea.Cmd
				m.scriptOutputView, vcmd = m.scriptOutputView.Update(msg)
				return m, vcmd
			}
			// running: wait for async capture or ExecProcess
			return m, nil
		}
		if m.state == ViewNativeMonitor {
			switch msg.String() {
			case "enter", "esc", "q":
				m.state = m.confirmReturn
				m.nativeMonitorKind = ""
				m.nativeBattery, m.nativeMemory = nil, nil
				m.nativeBatteryErr, m.nativeMemoryErr = nil, nil
				return m, nil
			case "r":
				return m, m.nativeMonitorLoadCmd()
			}
			return m, nil
		}
		if m.state == ViewScriptList && m.err != nil {
			m.err = nil
		}
		// In Zsh apply result screen, Enter/Esc return to menu
		if m.state == ViewZshApplying && (m.zshApplyPhase == "success" || m.zshApplyPhase == "error") {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.state = ViewMainMenu
				m.zshApplyPhase = ""
				m.zshApplyError = nil
				return m, nil
			}
		}
		// When in Zsh wizard or Zsh repo wizard, let the wizard handle keys (don't consume enter/o/c here)
		if m.state != ViewZshWizard && m.state != ViewZshRepoWizard {
			switch msg.String() {
			case "ctrl+c", "q":
				if m.state == ViewMainMenu {
					return m, tea.Quit
				}
				if m.state == ViewInstalling && m.canAbort {
					// Allow abort during download
					m.aborted = true
					m.installMessage = "Instalação abortada pelo usuário"
					m.state = ViewMainMenu
					return m, nil
				}
			case "esc":
				switch m.state {
				case ViewScriptList:
					m.state = ViewMainMenu
					m.confirmYes = false
					return m, nil
				case ViewConfirmation:
					if m.confirmReturnOK {
						m.state = m.confirmReturn
					} else {
						m.state = ViewMainMenu
					}
					m.confirmReturnOK = false
					m.confirmYes = false
					return m, nil
				case ViewPackageList:
					// Volta para as categorias de instaladores
					m.state = ViewInstallerCategories
					return m, nil
				case ViewInstallerCategories:
					m.state = ViewMainMenu
					return m, nil
				}
			case "left", "h":
				if m.state == ViewConfirmation {
					m.confirmYes = false
					return m, nil
				}
			case "right", "l":
				if m.state == ViewConfirmation {
					m.confirmYes = true
					return m, nil
				}
			case "o", "O":
				return m.handleURLShortcut(false)
			case "c", "C":
				return m.handleURLShortcut(true)
			case "enter":
				return m.handleEnter()
			}
		}

	case zshApplyResultMsg:
		if m.state == ViewZshApplying {
			if msg.Err != nil {
				m.zshApplyPhase = "error"
				m.zshApplyError = msg.Err
			} else {
				m.zshApplyPhase = "success"
				m.zshApplyError = nil
			}
			return m, tea.Tick(time.Second*2, func(time.Time) tea.Msg {
				return zshApplyReturnToMenuMsg{}
			})
		}
		return m, nil

	case zshApplyReturnToMenuMsg:
		if m.state == ViewZshApplying {
			m.state = ViewMainMenu
			m.zshApplyPhase = ""
			m.zshApplyError = nil
		}
		return m, nil

	case scriptCapturedMsg:
		if m.state != ViewScriptOutput {
			return m, nil
		}
		m.scriptOutputPhase = "done"
		m.scriptOutputErr = msg.err
		text := stripANSI(msg.output)
		if strings.TrimSpace(text) == "" {
			text = "(sem saída no stdout/stderr)"
		}
		if msg.err != nil {
			text += "\n\n──\n" + msg.err.Error()
		}
		m.scriptOutputView.SetContent(text)
		m.scriptOutputView.GotoTop()
		return m, nil

	case scriptExecFinishedMsg:
		m.state = m.confirmReturn
		m.scriptOutputPhase = ""
		m.scriptOutputTitle = ""
		m.scriptOutputErr = nil
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
		}
		return m, nil

	case nativeMonitorReloadMsg:
		if m.state != ViewNativeMonitor || msg.kind != m.nativeMonitorKind {
			return m, nil
		}
		switch msg.kind {
		case entities.NativeMonitorBattery:
			m.nativeBattery = msg.battery
			m.nativeBatteryErr = msg.err
		case entities.NativeMonitorMemory:
			m.nativeMemory = msg.memory
			m.nativeMemoryErr = msg.err
		}
		return m, nativeMonitorScheduleTick()

	case nativeMonitorTickMsg:
		if m.state != ViewNativeMonitor {
			return m, nil
		}
		return m, m.nativeMonitorLoadCmd()

	case urlActionDoneMsg:
		if msg.err != nil {
			m.keyboardToast = fmt.Sprintf("⚠ %v", msg.err)
		} else if msg.verb == "copy" {
			m.keyboardToast = "URL copiada para a área de transferência."
		} else {
			m.keyboardToast = "URL aberta no navegador (app padrão)."
		}
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearKeyboardToastMsg{} })

	case clearKeyboardToastMsg:
		m.keyboardToast = ""
		return m, nil
	}

	// Delegate to ZshWizard when in wizard state
	if m.state == ViewZshWizard && m.zshWizard != nil {
		newWizard, cmd := m.zshWizard.Update(msg)
		wizard := newWizard.(ZshWizardModel)
		m.zshWizard = &wizard

		// Check if wizard is done
		if wizard.IsDone() || wizard.IsCancelled() {
			if wizard.IsCancelled() {
				m.state = ViewMainMenu
				m.zshWizard = nil
				return m, cmd
			}
			// Done and not cancelled: apply config and show feedback
			selections := wizard.GetSelections()
			m.zshWizard = nil
			m.state = ViewZshApplying
			m.zshApplyPhase = "applying"
			m.zshApplyError = nil
			return m, applyZshConfigCmd(m.configService, selections)
		}

		return m, cmd
	}

	// Delegate to ZshRepoWizard when in repo wizard state
	if m.state == ViewZshRepoWizard && m.zshRepoWizard != nil {
		newRepo, cmd := m.zshRepoWizard.Update(msg)
		repoWizard := newRepo.(ZshRepoModel)
		m.zshRepoWizard = &repoWizard

		if repoWizard.IsDone() || repoWizard.IsCancelled() {
			m.state = ViewMainMenu
			m.zshRepoWizard = nil
			return m, cmd
		}
		return m, cmd
	}

	// Update the appropriate list based on state
	var cmd tea.Cmd
	switch m.state {
	case ViewMainMenu:
		m.mainMenu, cmd = m.mainMenu.Update(msg)
	case ViewScriptList:
		m.scriptList, cmd = m.scriptList.Update(msg)
	case ViewInstallerCategories:
		m.installerList, cmd = m.installerList.Update(msg)
	case ViewPackageList:
		m.packageList, cmd = m.packageList.Update(msg)
	}

	return m, cmd
}

// handleEnter handles the enter key based on current state
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case ViewMainMenu:
		selected := m.mainMenu.SelectedItem()
		item, ok := selected.(menuItem)
		if !ok {
			return m, nil
		}
		switch item.action {
		case menuActionCleanup:
			m.state = ViewScriptList
			m.selectedMenu = 0
			m.loadScripts(types.CategoryCleanup)
		case menuActionMonitoring:
			m.state = ViewScriptList
			m.selectedMenu = 1
			m.loadScripts(types.CategoryMonitoring)
		case menuActionInstallers:
			m.state = ViewInstallerCategories
			m.selectedMenu = 2
			m.loadInstallerCategories()
		case menuActionZshPlugins:
			m.state = ViewZshWizard
			wizardService := services.NewWizardService()
			wizard := NewZshWizardModel(wizardService)
			wizard.width = m.width
			wizard.height = m.height
			m.zshWizard = &wizard
		case menuActionZshRepo:
			m.state = ViewZshRepoWizard
			repoWizard := NewZshRepoModel(m.repoService, m.configService)
			repoWizard.width = m.width
			repoWizard.height = m.height
			m.zshRepoWizard = &repoWizard
		case menuActionQuit:
			return m, tea.Quit
		case menuActionSettings:
			// Em breve
			return m, nil
		default:
			return m, nil
		}

	case ViewScriptList:
		// Show confirmation for script execution
		selected := m.scriptList.SelectedItem()
		if scriptItem, ok := selected.(scriptItem); ok {
			m.selectedItem = scriptItem.script
			m.state = ViewConfirmation
			m.confirmYes = false
			m.confirmReturn = ViewScriptList
			m.confirmReturnOK = true
		}

	case ViewPackageList:
		// Show confirmation for package installation
		selected := m.packageList.SelectedItem()
		if pkgItem, ok := selected.(packageItem); ok {
			m.selectedItem = pkgItem.pkg
			m.state = ViewConfirmation
			m.confirmYes = false
			m.confirmReturn = ViewPackageList
			m.confirmReturnOK = true
		}

	case ViewInstallerCategories:
		// Ao escolher uma categoria, carregamos a lista de pacotes
		selected := m.installerList.SelectedItem()
		if catItem, ok := selected.(installerCategoryItem); ok && len(catItem.categories) > 0 {
			m.state = ViewPackageList
			m.loadPackagesFromCategories(catItem.categories)
		}

	case ViewConfirmation:
		if m.confirmYes {
			// User confirmed - execute action
			switch item := m.selectedItem.(type) {
			case entities.Script:
				if item.NativeMonitor != "" {
					m.nativeMonitorKind = item.NativeMonitor
					m.nativeBattery, m.nativeMemory = nil, nil
					m.nativeBatteryErr, m.nativeMemoryErr = nil, nil
					m.state = ViewNativeMonitor
					return m, m.nativeMonitorLoadCmd()
				}
				m.scriptOutputTitle = item.Name
				m.scriptOutputPhase = "running"
				m.scriptOutputErr = nil
				m.scriptOutputView = newScriptOutputViewport(m.width, m.height)
				m.state = ViewScriptOutput
				if item.RequiresSudo {
					cmd, err := m.scriptService.ScriptInteractiveCommand(item.ID)
					if err != nil {
						m.state = m.confirmReturn
						m.scriptOutputPhase = ""
						m.scriptOutputTitle = ""
						m.err = err
						return m, nil
					}
					return m, tea.ExecProcess(cmd, func(execErr error) tea.Msg {
						return scriptExecFinishedMsg{err: execErr}
					})
				}
				return m, runScriptCaptureCmd(m.scriptService, item.ID)
			case entities.Package:
				// Install package
				m.state = ViewInstalling
				m.installStatus = "preparing"
				m.installMessage = "Preparando instalação..."
				m.installPercent = 0
				m.canAbort = false
				m.aborted = false
				return m, installPackage(m.installerService, item.ID)
			}
		} else {
			if m.confirmReturnOK {
				m.state = m.confirmReturn
			} else {
				m.state = ViewMainMenu
			}
			m.confirmReturnOK = false
		}
	}

	return m, nil
}

// loadScripts loads scripts for the selected category
func (m *Model) loadScripts(category types.Category) {
	scripts, err := m.scriptService.GetScriptsByCategory(category)
	if err != nil {
		scripts = []entities.Script{}
	}

	items := make([]list.Item, len(scripts))
	for i, script := range scripts {
		items[i] = scriptItem{script: script}
	}

	delegate := list.NewDefaultDelegate()
	m.scriptList = list.New(items, delegate, m.width, m.height-4)

	categoryNames := map[types.Category]string{
		types.CategoryCleanup:    "🧹 Limpeza do Sistema",
		types.CategoryMonitoring: "📊 Monitoramento",
		types.CategoryInstall:    "📦 Instaladores",
	}

	m.scriptList.Title = categoryNames[category]
	m.scriptList.SetShowStatusBar(false)
}

// loadPackages loads packages for the selected category
func (m *Model) loadPackages(category types.PackageCategory) {
	packages, err := m.installerService.GetPackagesByCategory(category)
	if err != nil {
		packages = []entities.Package{}
	}

	m.setPackageList(packages, category, nil)
}

// loadPackagesFromCategories loads packages from multiple categories (e.g. IDE + Zsh core for Instaladores)
func (m *Model) loadPackagesFromCategories(categories []types.PackageCategory) {
	packages, err := m.installerService.GetPackagesByCategories(categories)
	if err != nil {
		packages = []entities.Package{}
	}

	categoryNames := map[types.PackageCategory]string{
		types.PackageCategoryIDE:       "💻 IDEs e Editores",
		types.PackageCategoryTool:      "🔧 Ferramentas de Desenvolvimento",
		types.PackageCategoryApp:       "📱 Aplicações",
		types.PackageCategoryZshCore:   "🐚 Componentes Core (Zsh)",
		types.PackageCategoryTerminal:  "🖥️ Emuladores de Terminal",
		types.PackageCategoryShell:     "🐚 Shells Alternativos",
		types.PackageCategoryAI:        "🤖 Integração com IA",
		types.PackageCategoryGames:     "🎮 Games",
		types.PackageCategorySysAdmin:  "🛡️ Administração de sistemas",
	}
	title := "📦 Instaladores (IDEs, Shells, Terminais, IA)"
	if len(categories) == 1 {
		title = categoryNames[categories[0]]
	}
	m.setPackageList(packages, categories[0], &title)
}

func (m *Model) setPackageList(packages []entities.Package, category types.PackageCategory, titleOverride *string) {
	items := make([]list.Item, len(packages))
	for i, pkg := range packages {
		items[i] = packageItem{pkg: pkg}
	}

	delegate := list.NewDefaultDelegate()
	m.packageList = list.New(items, delegate, m.width, m.height-4)

	categoryNames := map[types.PackageCategory]string{
		types.PackageCategoryIDE:       "💻 IDEs e Editores",
		types.PackageCategoryTool:      "🔧 Ferramentas de Desenvolvimento",
		types.PackageCategoryApp:       "📱 Aplicações",
		types.PackageCategoryZshCore:   "🐚 Componentes Core (Zsh)",
		types.PackageCategoryTerminal:  "🖥️ Emuladores de Terminal",
		types.PackageCategoryShell:     "🐚 Shells Alternativos",
		types.PackageCategoryAI:        "🤖 Integração com IA",
		types.PackageCategoryGames:     "🎮 Games",
		types.PackageCategorySysAdmin:  "🛡️ Administração de sistemas",
	}
	if titleOverride != nil {
		m.packageList.Title = *titleOverride
	} else {
		m.packageList.Title = categoryNames[category]
	}
	m.packageList.SetShowStatusBar(false)
}

// loadInstallerCategories inicializa a lista de categorias dentro de "Instaladores"
func (m *Model) loadInstallerCategories() {
	items := []list.Item{
		installerCategoryItem{
			title: "💻 IDEs & Dev (CLI)",
			desc:  "VS Code, Cursor, Claude Code, Antigravity e afins",
			categories: []types.PackageCategory{
				types.PackageCategoryIDE,
			},
		},
		installerCategoryItem{
			title: "📱 Aplicações",
			desc:  "Google Chrome, Insomnia e outras aplicações",
			categories: []types.PackageCategory{
				types.PackageCategoryApp,
			},
		},
		installerCategoryItem{
			title: "🔧 Ferramentas de desenvolvimento",
			desc:  "GitHub CLI (gh), NVM, Bun, pnpm, Deno e afins",
			categories: []types.PackageCategory{
				types.PackageCategoryTool,
			},
		},
		installerCategoryItem{
			title: "🐚 Shells alternativos",
			desc:  "Fish Shell e outros",
			categories: []types.PackageCategory{
				types.PackageCategoryShell,
			},
		},
		installerCategoryItem{
			title: "🖥️ Emuladores de Terminal",
			desc:  "WezTerm, Kitty, Alacritty, Zash Terminal, Warp, Wave e outros",
			categories: []types.PackageCategory{
				types.PackageCategoryTerminal,
			},
		},
		installerCategoryItem{
			title: "🐚 Componentes Core (Zsh)",
			desc:  "Zsh, Oh My Zsh, Powerlevel10k",
			categories: []types.PackageCategory{
				types.PackageCategoryZshCore,
			},
		},
		installerCategoryItem{
			title: "🎮 Games",
			desc:  "Prism Launcher, Lutris",
			categories: []types.PackageCategory{
				types.PackageCategoryGames,
			},
		},
		installerCategoryItem{
			title: "🤖 Integração com IA",
			desc:  "ShellGPT, Fish-AI e assistentes por shell",
			categories: []types.PackageCategory{
				types.PackageCategoryAI,
			},
		},
		installerCategoryItem{
			title: "🛡️ Administração de sistemas",
			desc:  "Cockpit, Webmin, Topgrade, integração AD, clientes SSH e diagnóstico",
			categories: []types.PackageCategory{
				types.PackageCategorySysAdmin,
			},
		},
	}

	delegate := list.NewDefaultDelegate()
	m.installerList = list.New(items, delegate, m.width, m.height-4)
	m.installerList.Title = "📦 Instaladores"
	m.installerList.SetShowStatusBar(false)
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Iniciando..."
	}

	switch m.state {
	case ViewMainMenu:
		return m.mainMenu.View()

	case ViewScriptList:
		help := helpStyle.Render("\n↑/↓: navegar • enter: executar • esc: voltar • q: sair")
		var feedback string
		if m.err != nil {
			feedback = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Render("\n⚠ "+m.err.Error()) + "\n"
		}
		return m.scriptList.View() + feedback + help

	case ViewInstallerCategories:
		help := helpStyle.Render("\n↑/↓: navegar • enter: abrir categoria • esc: voltar • q: sair")
		return m.installerList.View() + help

	case ViewPackageList:
		help := helpStyle.Render("\n↑/↓: navegar • enter: confirmação • o: abrir URL • c: copiar URL • esc: voltar • q: sair")
		toast := ""
		if strings.TrimSpace(m.keyboardToast) != "" {
			toast = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n"
		}
		return m.packageList.View() + toast + help

	case ViewConfirmation:
		return m.renderConfirmation()

	case ViewInstalling:
		return m.renderInstallProgress()

	case ViewZshWizard:
		if m.zshWizard != nil {
			return m.zshWizard.View()
		}
		return "Iniciando wizard..."

	case ViewZshApplying:
		return m.renderZshApplyFeedback()

	case ViewZshRepoWizard:
		if m.zshRepoWizard != nil {
			body := m.zshRepoWizard.View()
			if strings.TrimSpace(m.keyboardToast) != "" {
				body += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n"
			}
			return body
		}
		return "Iniciando Configurar Zsh..."

	case ViewScriptOutput:
		return m.renderScriptOutput()

	case ViewNativeMonitor:
		return m.renderNativeMonitorView()

	default:
		return ""
	}
}

// renderConfirmation renders the confirmation dialog
func (m Model) renderConfirmation() string {
	var title, description string

	switch item := m.selectedItem.(type) {
	case entities.Script:
		if item.NativeMonitor != "" {
			title = "Abrir monitor?"
			description = fmt.Sprintf("Você deseja abrir:\n\n  %s\n  %s", item.Name, item.Description)
		} else {
			title = "Executar Script?"
			description = fmt.Sprintf("Você deseja executar:\n\n  %s\n  %s", item.Name, item.Description)
			if item.RequiresSudo {
				description += "\n\n⚠️  Este script requer permissões de administrador (sudo)"
			}
		}
	case entities.Package:
		title = "Instalar pacote?"
		description = fmt.Sprintf("Você deseja instalar:\n\n  %s\n  %s\n  Versão: %s",
			item.Name, item.Description, item.Version)
		if kb := PackageKeyboardURL(item); kb != "" {
			description += "\n\n🔗 Verificação (sem mouse: tecla o abre no navegador, c copia a URL):\n  " + kb
		}
		if item.DownloadURL != "" {
			description += "\n\n⚠️  Será feito download do arquivo e em seguida os comandos de instalação."
		} else {
			description += "\n\n⚠️  Comandos serão executados no terminal; pode ser pedida senha de administrador (sudo)."
		}
		if strings.TrimSpace(item.Notes) != "" {
			description += "\n\n── Informações e avisos ──\n" + strings.TrimSpace(item.Notes)
		}
	default:
		return "Erro: tipo desconhecido"
	}

	var yesButton, noButton string
	if m.confirmYes {
		yesButton = selectedStyle.Render(" Sim ")
		noButton = noStyle.Render(" Não ")
	} else {
		yesButton = yesStyle.Render(" Sim ")
		noButton = selectedStyle.Render(" Não ")
	}

	helpConfirm := "←/→: escolher • enter: confirmar • esc: voltar"
	if p, ok := m.selectedItem.(entities.Package); ok && PackageKeyboardURL(p) != "" {
		helpConfirm = "o: abrir URL • c: copiar URL • " + helpConfirm
	}
	toastLine := ""
	if strings.TrimSpace(m.keyboardToast) != "" {
		toastLine = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n\n"
	}
	content := titleStyle.Render(title) + "\n\n" +
		description + "\n\n" +
		yesButton + "  " + noButton + "\n\n" +
		toastLine +
		helpStyle.Render(helpConfirm)

	boxW := m.width - 8
	if boxW < 52 {
		boxW = 52
	}
	if boxW > 88 {
		boxW = 88
	}
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(boxW)
	box := boxStyle.Render(content)

	// Center the box
	verticalPadding := (m.height - lipgloss.Height(box)) / 2
	horizontalPadding := (m.width - lipgloss.Width(box)) / 2

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			PaddingTop(verticalPadding).
			PaddingLeft(horizontalPadding).
			Render(box),
	)
}

// renderInstallProgress renders the installation progress view
func (m Model) renderInstallProgress() string {
	var pkg entities.Package
	if p, ok := m.selectedItem.(entities.Package); ok {
		pkg = p
	}

	title := titleStyle.Render(fmt.Sprintf("Instalando: %s", pkg.Name))

	statusIcons := map[string]string{
		"preparing":   "⏳",
		"downloading": "⬇️ ",
		"installing":  "⚙️ ",
		"complete":    "✅",
		"failed":      "❌",
	}

	icon := statusIcons[m.installStatus]
	if icon == "" {
		icon = m.spinner.View()
	}

	status := fmt.Sprintf("%s %s", icon, m.installMessage)
	progressBar := m.progress.ViewAs(m.installPercent)

	content := title + "\n\n" +
		status + "\n\n" +
		progressBar + "\n\n"

	if m.canAbort && !m.aborted {
		content += helpStyle.Render("⚠️  Pressione Ctrl+C para abortar (não recomendado)")
	} else if m.installStatus == "complete" {
		content += helpStyle.Render("Instalação concluída! Retornando ao menu...")
	} else if m.installStatus == "failed" {
		content += helpStyle.Render("❌ Instalação falhou. Retornando ao menu...")
	} else {
		content += helpStyle.Render("Aguarde... não feche o programa")
	}

	box := confirmBoxStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func scriptOutputCardWidth(termW int) int {
	boxW := termW - 8
	if boxW < 52 {
		boxW = 52
	}
	if boxW > 88 {
		boxW = 88
	}
	return boxW
}

// Área útil do viewport (cabeçalho + rodapé do cartão consomem linhas)
func scriptOutputViewportWH(termW, termH int) (w, h int) {
	boxW := scriptOutputCardWidth(termW)
	w = boxW - 8
	if w < 28 {
		w = 28
	}
	h = termH - 20
	if h < 8 {
		h = 8
	}
	if termW < 20 || termH < 16 {
		w, h = 64, 12
	}
	return w, h
}

func newScriptOutputViewport(termW, termH int) viewport.Model {
	w, h := scriptOutputViewportWH(termW, termH)
	vp := viewport.New(w, h)
	vp.Style = scriptLogAreaStyle
	return vp
}

func (m *Model) syncScriptOutputViewport() {
	if m.width < 8 || m.height < 8 {
		return
	}
	w, h := scriptOutputViewportWH(m.width, m.height)
	m.scriptOutputView.Width = w
	m.scriptOutputView.Height = h
}

func scriptOutputDivider(width int) string {
	n := width - 4
	if n < 12 {
		n = 12
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(strings.Repeat("─", n))
}

func (m Model) renderScriptOutput() string {
	boxW := scriptOutputCardWidth(m.width)

	if m.scriptOutputPhase == "running" {
		head := titleStyle.Render("Homestead") + "\n" +
			helpStyle.Render("Gerenciador de Sistema") + "\n" +
			scriptOutputDivider(boxW) + "\n" +
			scriptScreenAccentStyle.Render("📜 Executando script") + "\n" +
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Render(m.scriptOutputTitle)
		body := "\n\n" + fmt.Sprintf("%s %s", m.spinner.View(), helpStyle.Render("Capturando saída…"))
		body += "\n\n" + helpStyle.Render("A saída aparecerá no painel abaixo quando o script terminar.")
		body += "\n" + helpStyle.Render("Scripts com sudo usam o terminal completo para pedir senha.")
		content := head + body
		box := scriptScreenOuterStyle.Width(boxW)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
	}

	head := titleStyle.Render("Homestead") + "\n" +
		helpStyle.Render("Gerenciador de Sistema") + "\n" +
		scriptOutputDivider(boxW) + "\n" +
		scriptScreenAccentStyle.Render("📜 Saída do script")
	nameLine := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Render(m.scriptOutputTitle)
	if m.scriptOutputErr != nil {
		nameLine += "  " + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render("· falhou")
	}
	view := m.scriptOutputView.View()
	footerText := "↑/↓  PgUp/PgDn  rolar · Enter / Esc / q  voltar"
	if m.scriptOutputErr != nil {
		footerText = "Ver mensagem de erro no fim do texto · " + footerText
	}
	footer := scriptScreenFooterBarStyle.Width(max(12, boxW-8)).Render(footerText)
	content := head + "\n" + nameLine + "\n\n" + view + "\n" + footer
	box := scriptScreenOuterStyle.Width(boxW)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

func runScriptCaptureCmd(service *services.ScriptService, scriptID string) tea.Cmd {
	return func() tea.Msg {
		out, err := service.ExecuteScriptCapture(scriptID)
		return scriptCapturedMsg{output: out, err: err}
	}
}

// installPackage creates a command that installs a package
func installPackage(service *services.InstallerService, packageID string) tea.Cmd {
	return func() tea.Msg {
		progressChan := make(chan interfaces.InstallProgress, 10)

		go func() {
			err := service.InstallPackage(packageID, func(progress interfaces.InstallProgress) {
				progressChan <- progress
			})
			if err != nil {
				progressChan <- interfaces.InstallProgress{
					Status:      "failed",
					Message:     err.Error(),
					IsCompleted: true,
					Error:       err,
				}
			}
			close(progressChan)
		}()

		// Return first progress update
		for progress := range progressChan {
			return progressMsg(progress)
		}

		return installCompleteMsg{err: nil}
	}
}

// applyZshConfigCmd runs ConfigService.ApplyConfig and sends zshApplyResultMsg
func applyZshConfigCmd(configService *services.ConfigService, selections interfaces.ConfigSelections) tea.Cmd {
	return func() tea.Msg {
		err := configService.ApplyConfig(selections)
		return zshApplyResultMsg{Err: err}
	}
}

// renderZshApplyFeedback renders the Zsh apply state (applying / success / error)
func (m Model) renderZshApplyFeedback() string {
	title := titleStyle.Render("Configuração Zsh")

	switch m.zshApplyPhase {
	case "applying":
		status := fmt.Sprintf("%s Aplicando configuração Zsh...", m.spinner.View())
		content := title + "\n\n" + status + "\n\n" + helpStyle.Render("Aguarde...")
		box := confirmBoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	case "success":
		content := title + "\n\n" +
			"✅ Configuração aplicada com sucesso.\n\n" +
			"O arquivo ~/.zshrc foi atualizado com os plugins e ferramentas selecionados.\n" +
			"Criados/atualizados: ~/.zsh/general/aliases.zsh e functions.zsh.\n\n" +
			"Verifique: cat ~/.zshrc\n\n" +
			"Não instala plugins externos (ex.: zsh-autosuggestions); apenas escreve o .zshrc.\n" +
			"Use Instaladores para Zsh/Oh My Zsh se ainda não estiverem instalados.\n\n" +
			helpStyle.Render("Retornando ao menu em 2s (ou Enter/Esc para voltar)")
		box := confirmBoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	case "error":
		errMsg := ""
		if m.zshApplyError != nil {
			errMsg = m.zshApplyError.Error()
		}
		content := title + "\n\n" +
			"❌ Erro ao aplicar configuração:\n\n" + errMsg + "\n\n" +
			helpStyle.Render("Pressione Enter ou Esc para voltar ao menu")
		box := confirmBoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	default:
		content := title + "\n\n" + m.spinner.View() + " Aguarde..."
		box := confirmBoxStyle.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
}

func (m Model) urlForKeyboardOpen() string {
	switch m.state {
	case ViewConfirmation:
		if p, ok := m.selectedItem.(entities.Package); ok {
			return PackageKeyboardURL(p)
		}
	case ViewPackageList:
		if sel := m.packageList.SelectedItem(); sel != nil {
			if it, ok := sel.(packageItem); ok {
				return PackageKeyboardURL(it.pkg)
			}
		}
	}
	return ""
}

func (m Model) handleURLShortcut(wantCopy bool) (Model, tea.Cmd) {
	url := m.urlForKeyboardOpen()
	if url != "" {
		if wantCopy {
			return m, copyURLTeaCmd(url)
		}
		return m, openURLTeaCmd(url)
	}
	if m.state == ViewPackageList || (m.state == ViewConfirmation && isSelectedPackage(m.selectedItem)) {
		m.keyboardToast = "Este pacote não tem URL de projeto nem download."
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearKeyboardToastMsg{} })
	}
	return m, nil
}

func isSelectedPackage(item interface{}) bool {
	_, ok := item.(entities.Package)
	return ok
}

func openURLTeaCmd(url string) tea.Cmd {
	return func() tea.Msg {
		err := OpenURL(url)
		return urlActionDoneMsg{verb: "open", err: err}
	}
}

func copyURLTeaCmd(url string) tea.Cmd {
	return func() tea.Msg {
		err := CopyURLToClipboard(url)
		return urlActionDoneMsg{verb: "copy", err: err}
	}
}
