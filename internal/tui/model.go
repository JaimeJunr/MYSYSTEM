package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

// ViewState represents different views in the TUI
type ViewState int

const (
	ViewMainMenu ViewState = iota
	ViewScriptList
	ViewPackageList
	ViewConfirmation
	ViewExecuting
	ViewInstalling
	ViewZshWizard
	ViewZshApplying
)

// Model is the main TUI model
type Model struct {
	scriptService    *services.ScriptService
	installerService *services.InstallerService
	configService    *services.ConfigService
	state            ViewState
	mainMenu         list.Model
	scriptList       list.Model
	packageList      list.Model
	selectedMenu     int
	selectedItem     interface{} // Can be Script or Package
	confirmYes       bool        // true = yes selected, false = no selected
	width            int
	height           int
	err              error

	// Installation progress
	progress       progress.Model
	spinner        spinner.Model
	installStatus  string
	installMessage string
	installPercent float64
	canAbort       bool
	aborted        bool

	// Zsh wizard
	zshWizard *ZshWizardModel

	// Zsh core: when true, "Configurar Zsh" is shown in menu (oh-my-zsh installed)
	zshCoreInstalled bool
	zshCoreChecked   bool

	// Zsh apply feedback: phase "applying" | "success" | "error"
	zshApplyPhase string
	zshApplyError error
}

// menuAction identifies the main menu action (for dynamic menu with optional "Configurar Zsh")
const (
	menuActionCleanup     = "cleanup"
	menuActionMonitoring  = "monitoring"
	menuActionInstallers  = "installers"
	menuActionZshConfig   = "zsh_config"
	menuActionMigration  = "migration"
	menuActionSettings    = "settings"
	menuActionQuit       = "quit"
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
)

// getMainMenuItems returns menu items; "Configurar Zsh" only when zsh core is installed
func getMainMenuItems(zshCoreInstalled bool) []list.Item {
	items := []list.Item{
		menuItem{title: "🧹 Limpeza do Sistema", desc: "Scripts de limpeza e manutenção", action: menuActionCleanup},
		menuItem{title: "📊 Monitoramento", desc: "Informações do sistema", action: menuActionMonitoring},
		menuItem{title: "📦 Instaladores", desc: "Instalar ferramentas e aplicações (IDEs, Zsh, Oh My Zsh, etc.)", action: menuActionInstallers},
	}
	if zshCoreInstalled {
		items = append(items, menuItem{title: "⚙️  Configurar Zsh", desc: "Plugins, temas e .zshrc", action: menuActionZshConfig})
	}
	items = append(items,
		menuItem{title: "🔄 Migração", desc: "Exportar/Importar configurações (em breve)", action: menuActionMigration},
		menuItem{title: "⚙️  Configurações", desc: "Configurar a ferramenta (em breve)", action: menuActionSettings},
		menuItem{title: "❌ Sair", desc: "Fechar Homestead", action: menuActionQuit},
	)
	return items
}

// NewModel creates the TUI model with dependencies injected
func NewModel(scriptService *services.ScriptService, installerService *services.InstallerService, configService *services.ConfigService) Model {
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height-4)
		if m.scriptList.Items() != nil {
			m.scriptList.SetSize(msg.Width, msg.Height-4)
		}
		if m.packageList.Items() != nil {
			m.packageList.SetSize(msg.Width, msg.Height-4)
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
		// In Zsh apply result screen, Enter/Esc return to menu
		if m.state == ViewZshApplying && (m.zshApplyPhase == "success" || m.zshApplyPhase == "error") {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.state = ViewMainMenu
				m.zshApplyPhase = ""
				m.zshApplyError = nil
				return m, nil
			}
		}
		// When in Zsh wizard, let the wizard handle all keys (don't consume enter/esc here)
		if m.state != ViewZshWizard {
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
				if m.state == ViewScriptList || m.state == ViewPackageList || m.state == ViewConfirmation {
					m.state = ViewMainMenu
					m.confirmYes = false
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

	// Update the appropriate list based on state
	var cmd tea.Cmd
	switch m.state {
	case ViewMainMenu:
		m.mainMenu, cmd = m.mainMenu.Update(msg)
	case ViewScriptList:
		m.scriptList, cmd = m.scriptList.Update(msg)
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
			m.state = ViewPackageList
			m.selectedMenu = 2
			m.loadPackagesFromCategories([]types.PackageCategory{types.PackageCategoryIDE, types.PackageCategoryZshCore})
		case menuActionZshConfig:
			m.state = ViewZshWizard
			wizardService := services.NewWizardService()
			wizard := NewZshWizardModel(wizardService)
			wizard.width = m.width
			wizard.height = m.height
			m.zshWizard = &wizard
		case menuActionQuit:
			return m, tea.Quit
		case menuActionMigration, menuActionSettings:
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
		}

	case ViewPackageList:
		// Show confirmation for package installation
		selected := m.packageList.SelectedItem()
		if pkgItem, ok := selected.(packageItem); ok {
			m.selectedItem = pkgItem.pkg
			m.state = ViewConfirmation
			m.confirmYes = false
		}

	case ViewConfirmation:
		if m.confirmYes {
			// User confirmed - execute action
			switch item := m.selectedItem.(type) {
			case entities.Script:
				// Execute script
				m.state = ViewExecuting
				return m, tea.Sequence(
					tea.ExitAltScreen,
					executeScript(m.scriptService, item.ID),
				)
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
			// User cancelled - go back
			m.state = ViewMainMenu
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
		types.PackageCategoryIDE:     "💻 IDEs e Editores",
		types.PackageCategoryTool:    "🔧 Ferramentas de Desenvolvimento",
		types.PackageCategoryApp:     "📱 Aplicações",
		types.PackageCategoryZshCore: "🐚 Componentes Core (Zsh)",
	}
	title := "📦 Instaladores"
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
		types.PackageCategoryIDE:     "💻 IDEs e Editores",
		types.PackageCategoryTool:    "🔧 Ferramentas de Desenvolvimento",
		types.PackageCategoryApp:     "📱 Aplicações",
		types.PackageCategoryZshCore: "🐚 Componentes Core (Zsh)",
	}
	if titleOverride != nil {
		m.packageList.Title = *titleOverride
	} else {
		m.packageList.Title = categoryNames[category]
	}
	m.packageList.SetShowStatusBar(false)
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
		return m.scriptList.View() + help

	case ViewPackageList:
		help := helpStyle.Render("\n↑/↓: navegar • enter: instalar • esc: voltar • q: sair")
		return m.packageList.View() + help

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

	default:
		return "Executando..."
	}
}

// renderConfirmation renders the confirmation dialog
func (m Model) renderConfirmation() string {
	var title, description string

	switch item := m.selectedItem.(type) {
	case entities.Script:
		title = "Executar Script?"
		description = fmt.Sprintf("Você deseja executar:\n\n  %s\n  %s", item.Name, item.Description)
		if item.RequiresSudo {
			description += "\n\n⚠️  Este script requer permissões de administrador (sudo)"
		}
	case entities.Package:
		title = "Instalar Pacote?"
		description = fmt.Sprintf("Você deseja instalar:\n\n  %s\n  %s\n  Versão: %s",
			item.Name, item.Description, item.Version)
		description += "\n\n⚠️  O download será iniciado e a instalação executada"
	default:
		return "Erro: tipo desconhecido"
	}

	yesButton := " Não "
	noButton := " Sim "

	if m.confirmYes {
		yesButton = selectedStyle.Render(" Sim ")
		noButton = noStyle.Render(" Não ")
	} else {
		yesButton = yesStyle.Render(" Sim ")
		noButton = selectedStyle.Render(" Não ")
	}

	content := titleStyle.Render(title) + "\n\n" +
		description + "\n\n" +
		yesButton + "  " + noButton + "\n\n" +
		helpStyle.Render("←/→: escolher • enter: confirmar • esc: cancelar")

	box := confirmBoxStyle.Render(content)

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

// executeScript creates a command that executes a script and quits
func executeScript(service *services.ScriptService, scriptID string) tea.Cmd {
	return func() tea.Msg {
		err := service.ExecuteScript(scriptID)
		if err != nil {
			return tea.Quit()
		}
		return tea.Quit()
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
