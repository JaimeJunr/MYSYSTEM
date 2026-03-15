package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/app/services"
	"github.com/JaimeJunr/Homestead/internal/domain/interfaces"
)

// KeySelectAll is the standard key for "marcar todos" in multi-select lists (DRY: reuse in other parts)
const KeySelectAll = "a"

// ZshWizardView represents each view/step of the wizard (core is installed separately)
type ZshWizardView int

const (
	ZshWizardViewPlugins ZshWizardView = iota
	ZshWizardViewTools
	ZshWizardViewProjectConfig
	ZshWizardViewReview
)

// WizardItem represents a selectable item in the wizard
type WizardItem struct {
	ID          string
	Name        string
	Description string
	Selected    bool
}

// ZshWizardModel is the Bubbletea model for the Zsh configuration wizard
type ZshWizardModel struct {
	wizardService *services.WizardService
	state         *services.WizardState
	currentView   ZshWizardView

	// Items for each view
	coreItems    []WizardItem
	pluginItems  []WizardItem
	toolItems    []WizardItem
	projectItems []WizardItem

	// Current cursor position in list
	cursor int

	// Dimensions
	width  int
	height int

	// Done flag (user wants to apply or cancel)
	done      bool
	cancelled bool
}

// Styles for the wizard
var (
	wizardTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)

	wizardStepStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginBottom(1)

	wizardItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	wizardSelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("10")).
				Bold(true)

	wizardCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true)

	wizardHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	wizardPreviewStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2)

	wizardProgressStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("63"))
)

// NewZshWizardModel creates a new Zsh wizard model (assumes core is already installed)
func NewZshWizardModel(wizardService *services.WizardService) ZshWizardModel {
	state := wizardService.CreateNewWizard()
	// Pre-fill core so generated .zshrc assumes zsh + oh-my-zsh + powerlevel10k
	state.Selections.CoreComponents = []string{"zsh", "oh-my-zsh", "powerlevel10k"}

	m := ZshWizardModel{
		wizardService: wizardService,
		state:         state,
		currentView:   ZshWizardViewPlugins,
		cursor:        0,
	}

	m.initItems()

	return m
}

// initItems initializes the items for all wizard views
func (m *ZshWizardModel) initItems() {
	// Core components
	m.coreItems = []WizardItem{
		{ID: "zsh", Name: "Zsh", Description: "Z Shell - shell poderoso e configurável"},
		{ID: "oh-my-zsh", Name: "Oh My Zsh", Description: "Framework para gerenciar configuração Zsh"},
		{ID: "powerlevel10k", Name: "Powerlevel10k", Description: "Tema Zsh rápido e customizável"},
	}

	// Plugins
	m.pluginItems = []WizardItem{
		{ID: "git", Name: "git", Description: "Plugin para Git (built-in)"},
		{ID: "docker", Name: "docker", Description: "Plugin para Docker (built-in)"},
		{ID: "rails", Name: "rails", Description: "Plugin para Ruby on Rails (built-in)"},
		{ID: "z", Name: "z", Description: "Navegação rápida de diretórios (built-in)"},
		{ID: "sudo", Name: "sudo", Description: "Adiciona sudo com double ESC (built-in)"},
		{ID: "zsh-autosuggestions", Name: "zsh-autosuggestions", Description: "Sugestões automáticas baseadas no histórico"},
		{ID: "zsh-syntax-highlighting", Name: "zsh-syntax-highlighting", Description: "Destaque de sintaxe para comandos"},
		{ID: "fzf-zsh", Name: "fzf", Description: "Integração fuzzy finder"},
		{ID: "you-should-use", Name: "you-should-use", Description: "Lembra aliases existentes"},
		{ID: "zsh-completions", Name: "zsh-completions", Description: "Completions adicionais"},
		{ID: "zsh-history-substring-search", Name: "history-substring-search", Description: "Busca no histórico por substring"},
		{ID: "fast-syntax-highlighting", Name: "fast-syntax-highlighting", Description: "Syntax highlighting mais rápido"},
		{ID: "auto-notify", Name: "auto-notify", Description: "Notificações para comandos longos"},
		{ID: "zsh-vi-mode", Name: "zsh-vi-mode", Description: "Melhor modo Vi para Zsh"},
		{ID: "zsh-autocomplete", Name: "zsh-autocomplete", Description: "Autocomplete em tempo real"},
	}

	// Tools
	m.toolItems = []WizardItem{
		{ID: "nvm", Name: "NVM", Description: "Node Version Manager"},
		{ID: "bun", Name: "Bun", Description: "Runtime JavaScript/TypeScript rápido"},
		{ID: "sdkman", Name: "SDKMAN!", Description: "Gerenciador de SDKs para JVM"},
		{ID: "pnpm", Name: "pnpm", Description: "Gerenciador de pacotes Node.js eficiente"},
		{ID: "deno", Name: "Deno", Description: "Runtime seguro para JavaScript e TypeScript"},
		{ID: "angular-cli", Name: "Angular CLI", Description: "CLI para Angular"},
		{ID: "openvpn3", Name: "OpenVPN 3", Description: "Cliente VPN moderno"},
		{ID: "homebrew", Name: "Homebrew", Description: "Gerenciador de pacotes para Linux"},
	}

	// Project configs
	m.projectItems = []WizardItem{
		{ID: "include-project", Name: "Incluir configs de projeto", Description: "Carrega ~/.zsh/projects/ ao iniciar"},
		{ID: "separate-aliases", Name: "Aliases separados", Description: "Mantém aliases em ~/.zsh/general/aliases.zsh"},
		{ID: "separate-functions", Name: "Funções separadas", Description: "Mantém funções em ~/.zsh/general/functions.zsh"},
	}
}

// Init initializes the model
func (m ZshWizardModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates state
func (m ZshWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

// handleKey handles keyboard input
func (m ZshWizardModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle keys that work in all views (including Review, where getCurrentItems() is nil)
	switch key {
	case "ctrl+c":
		m.cancelled = true
		m.done = true
		return m, tea.Quit

	case "esc":
		if m.currentView > ZshWizardViewPlugins {
			m.currentView--
			_ = m.wizardService.PreviousStep(m.state)
			m.cursor = 0
		} else {
			m.cancelled = true
			m.done = true
		}
		return m, nil

	case "n", "tab", "right":
		if m.currentView == ZshWizardViewReview {
			m.done = true
			return m, nil
		}
		if m.currentView < ZshWizardViewReview {
			m.currentView++
			_ = m.wizardService.NextStep(m.state)
			m.cursor = 0
		}
		return m, nil

	case "enter":
		if m.currentView == ZshWizardViewReview {
			m.done = true
			return m, nil
		}
	}
	// From here on we need a list (not on Review)
	items := m.getCurrentItems()
	if items == nil {
		return m, nil
	}

	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case "down", "j":
		if m.cursor < len(*items)-1 {
			m.cursor++
		}
		return m, nil

	case " ": // Space to toggle selection
		if m.cursor >= 0 && m.cursor < len(*items) {
			m.toggleItem(items, m.cursor)
		}
		return m, nil

	case KeySelectAll: // Marcar todos (padrão do sistema para listas de seleção)
		m.selectAllInCurrentView(items)
		return m, nil

	case "enter":
		if m.cursor >= 0 && m.cursor < len(*items) {
			m.toggleItem(items, m.cursor)
		}
		return m, nil
	}

	return m, nil
}

// toggleItem toggles the selected state of an item
func (m *ZshWizardModel) toggleItem(items *[]WizardItem, index int) {
	if index < 0 || index >= len(*items) {
		return
	}

	item := &(*items)[index]
	item.Selected = !item.Selected

	// Update wizard state
	switch m.currentView {
	case ZshWizardViewPlugins:
		if item.Selected {
			m.wizardService.AddPlugin(m.state, item.ID)
		} else {
			m.wizardService.RemovePlugin(m.state, item.ID)
		}
	case ZshWizardViewTools:
		if item.Selected {
			m.wizardService.AddTool(m.state, item.ID)
		} else {
			m.wizardService.RemoveTool(m.state, item.ID)
		}
	case ZshWizardViewProjectConfig:
		if item.ID == "include-project" {
			m.wizardService.SetIncludeProjectConfig(m.state, item.Selected)
		}
	}
}

// getCurrentItems returns the current view's items
func (m *ZshWizardModel) getCurrentItems() *[]WizardItem {
	switch m.currentView {
	case ZshWizardViewPlugins:
		return &m.pluginItems
	case ZshWizardViewTools:
		return &m.toolItems
	case ZshWizardViewProjectConfig:
		return &m.projectItems
	}
	return nil
}

// selectAllInCurrentView marks all items in the current list and syncs to wizard state.
// Padrão reutilizável: outras telas de seleção múltipla podem usar a mesma tecla "a".
func (m *ZshWizardModel) selectAllInCurrentView(items *[]WizardItem) {
	if items == nil {
		return
	}
	for i := range *items {
		item := &(*items)[i]
		if item.Selected {
			continue
		}
		item.Selected = true
		switch m.currentView {
		case ZshWizardViewPlugins:
			m.wizardService.AddPlugin(m.state, item.ID)
		case ZshWizardViewTools:
			m.wizardService.AddTool(m.state, item.ID)
		case ZshWizardViewProjectConfig:
			if item.ID == "include-project" {
				m.wizardService.SetIncludeProjectConfig(m.state, true)
			}
		}
	}
}

// View renders the wizard UI
func (m ZshWizardModel) View() string {
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	switch m.currentView {
	case ZshWizardViewPlugins:
		return m.renderSelectionView("Plugins Zsh", "Selecione os plugins que deseja instalar", m.pluginItems)
	case ZshWizardViewTools:
		return m.renderSelectionView("Ferramentas de Desenvolvimento", "Selecione as ferramentas de desenvolvimento", m.toolItems)
	case ZshWizardViewProjectConfig:
		return m.renderSelectionView("Configurações de Projeto", "Opções de organização de configuração", m.projectItems)
	case ZshWizardViewReview:
		return m.renderReviewView()
	}

	return ""
}

// renderSelectionView renders a multi-select list view
func (m ZshWizardModel) renderSelectionView(title, subtitle string, items []WizardItem) string {
	var builder strings.Builder

	// Progress indicator
	progress := m.GetProgress()
	totalSteps := m.wizardService.GetTotalSteps()
	stepNum := int(m.currentView) + 1
	progressLine := wizardProgressStyle.Render(
		fmt.Sprintf("Etapa %d/%d (%d%%)", stepNum, totalSteps, progress),
	)

	// Title
	titleLine := wizardTitleStyle.Render(fmt.Sprintf("🔧 %s", title))
	subtitleLine := wizardStepStyle.Render(subtitle)

	builder.WriteString(progressLine + "\n\n")
	builder.WriteString(titleLine + "\n")
	builder.WriteString(subtitleLine + "\n\n")

	// Items list
	for i, item := range items {
		cursor := "  "
		if m.cursor == i {
			cursor = wizardCursorStyle.Render("> ")
		}

		checkbox := "[ ]"
		if item.Selected {
			checkbox = wizardSelectedItemStyle.Render("[✓]")
		}

		line := fmt.Sprintf("%s%s %s", cursor, checkbox, item.Name)
		if item.Description != "" {
			line += wizardStepStyle.Render(" - " + item.Description)
		}

		if m.cursor == i {
			builder.WriteString(wizardCursorStyle.Render(line) + "\n")
		} else {
			builder.WriteString(line + "\n")
		}
	}

	// Help (a: marcar todos = padrão reutilizável para seleção múltipla)
	help := "\n" + wizardHelpStyle.Render(
		"↑/↓: navegar • espaço: selecionar • a: marcar todos • n/→: próximo • esc: voltar • ctrl+c: sair",
	)
	builder.WriteString(help)

	return builder.String()
}

// renderReviewView renders the review/confirmation view
func (m ZshWizardModel) renderReviewView() string {
	var builder strings.Builder

	title := wizardTitleStyle.Render("✅ Revisão e Confirmação")
	builder.WriteString(title + "\n\n")

	// Preview of selections
	preview := m.wizardService.GeneratePreview(m.state)
	previewBox := wizardPreviewStyle.Render(preview)
	builder.WriteString(previewBox + "\n\n")

	help := wizardHelpStyle.Render(
		"enter/n: confirmar e aplicar • esc: voltar • ctrl+c: cancelar",
	)
	builder.WriteString(help)

	return builder.String()
}

// GetProgress returns the current progress percentage
func (m ZshWizardModel) GetProgress() int {
	return m.wizardService.GetProgress(m.state)
}

// GetSelections returns the current ConfigSelections
func (m ZshWizardModel) GetSelections() interfaces.ConfigSelections {
	return m.state.Selections
}

// IsDone returns true if the wizard is done
func (m ZshWizardModel) IsDone() bool {
	return m.done
}

// IsCancelled returns true if the wizard was cancelled
func (m ZshWizardModel) IsCancelled() bool {
	return m.cancelled
}
