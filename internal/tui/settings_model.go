package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

type settingsMode int

const (
	settingsBrowse settingsMode = iota
	settingsEdit
)

const (
	settingsRowCatalog = iota
	settingsRowTheme
	settingsRowHighContrast
	settingsRowReduceMotion
	settingsRowTextScale
	settingsRowScriptRoot
	settingsRowDotfiles
	settingsRowConfirmScript
	settingsRowConfirmPackage
	settingsRowSave
	settingsRowCancel
	settingsRowCount
)

type SettingsModel struct {
	draft              preferences.Preferences
	prefsPath          string
	catalogEnvOverride bool

	cursor     int
	mode       settingsMode
	input      textinput.Model
	editTarget string

	statusLine string
	width      int
	height     int
}

func (m *SettingsModel) IsEditing() bool {
	return m.mode == settingsEdit
}

func NewSettingsModel(p preferences.Preferences, prefsPath string, catalogEnvOverride bool) SettingsModel {
	ti := textinput.New()
	ti.CharLimit = 2048
	ti.Width = 72

	return SettingsModel{
		draft:              p,
		prefsPath:          prefsPath,
		catalogEnvOverride: catalogEnvOverride,
		cursor:             0,
		mode:               settingsBrowse,
		input:              ti,
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.mode == settingsEdit {
			return m.updateEdit(msg)
		}
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.statusLine = ""
			return m, nil
		case "down", "j":
			if m.cursor < settingsRowCount-1 {
				m.cursor++
			}
			m.statusLine = ""
			return m, nil
		case "enter":
			return m.activateRow()
		case "esc", "q":
			return m, func() tea.Msg { return btmsg.SettingsCancelled{} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m SettingsModel) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = settingsBrowse
		m.statusLine = ""
		return m, nil
	case "enter":
		val := strings.TrimSpace(m.input.Value())
		switch m.editTarget {
		case "catalog_url":
			if err := preferences.ValidateCatalogURL(val); err != nil {
				m.statusLine = err.Error()
				return m, nil
			}
			m.draft.CatalogURL = val
		case "script_root":
			exp, err := preferences.ExpandPath(val)
			if err != nil {
				m.statusLine = err.Error()
				return m, nil
			}
			if err := preferences.ValidateScriptRoot(exp); err != nil {
				m.statusLine = err.Error()
				return m, nil
			}
			m.draft.ScriptRoot = val
		case "dotfiles_repo":
			m.draft.DotfilesRepo = val
		}
		m.mode = settingsBrowse
		m.statusLine = ""
		return m, nil
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m SettingsModel) activateRow() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case settingsRowCatalog:
		if m.catalogEnvOverride {
			m.statusLine = "HOMESTEAD_CATALOG_URL está definida; o ficheiro é ignorado até limpar a variável."
			return m, nil
		}
		m.mode = settingsEdit
		m.editTarget = "catalog_url"
		m.input.SetValue(m.draft.CatalogURL)
		m.input.Focus()
		return m, textinput.Blink
	case settingsRowTheme:
		if m.draft.Theme == preferences.ThemeLight {
			m.draft.Theme = preferences.ThemeDark
		} else {
			m.draft.Theme = preferences.ThemeLight
		}
		m.statusLine = ""
		return m, nil
	case settingsRowHighContrast:
		m.draft.HighContrast = !m.draft.HighContrast
		m.statusLine = ""
		return m, nil
	case settingsRowReduceMotion:
		m.draft.ReduceMotion = !m.draft.ReduceMotion
		m.statusLine = ""
		return m, nil
	case settingsRowTextScale:
		switch m.draft.TextScale {
		case preferences.TextScaleNormal:
			m.draft.TextScale = preferences.TextScaleLarge
		case preferences.TextScaleLarge:
			m.draft.TextScale = preferences.TextScaleXLarge
		default:
			m.draft.TextScale = preferences.TextScaleNormal
		}
		m.statusLine = ""
		return m, nil
	case settingsRowScriptRoot:
		m.mode = settingsEdit
		m.editTarget = "script_root"
		m.input.SetValue(m.draft.ScriptRoot)
		m.input.Focus()
		return m, textinput.Blink
	case settingsRowDotfiles:
		m.mode = settingsEdit
		m.editTarget = "dotfiles_repo"
		m.input.SetValue(m.draft.DotfilesRepo)
		m.input.Focus()
		return m, textinput.Blink
	case settingsRowConfirmScript:
		m.draft.ConfirmBeforeScript = !m.draft.ConfirmBeforeScript
		return m, nil
	case settingsRowConfirmPackage:
		m.draft.ConfirmBeforePackage = !m.draft.ConfirmBeforePackage
		return m, nil
	case settingsRowSave:
		return m.trySave()
	case settingsRowCancel:
		return m, func() tea.Msg { return btmsg.SettingsCancelled{} }
	}
	return m, nil
}

func (m SettingsModel) trySave() (tea.Model, tea.Cmd) {
	if err := preferences.ValidateCatalogURL(m.draft.CatalogURL); err != nil {
		m.statusLine = err.Error()
		return m, nil
	}
	expRoot, err := preferences.ExpandPath(m.draft.ScriptRoot)
	if err != nil {
		m.statusLine = err.Error()
		return m, nil
	}
	if err := preferences.ValidateScriptRoot(expRoot); err != nil {
		m.statusLine = err.Error()
		return m, nil
	}
	m.draft.Normalize()
	p := m.draft
	return m, func() tea.Msg { return btmsg.SettingsSaved{Prefs: p} }
}

func (m SettingsModel) View() string {
	var b strings.Builder
	b.WriteString(theme.Title.Render("⚙️  Configurações") + "\n\n")
	b.WriteString(theme.Help.Render("Fluxo: preferências YAML → arranque do TUI e deste ecrã.") + "\n")
	if strings.TrimSpace(m.prefsPath) != "" {
		b.WriteString(theme.Help.Render("Ficheiro: "+m.prefsPath) + "\n")
	}
	b.WriteString(theme.Help.Render("Catálogo de instaladores: URL JSON; vazio = manifesto padrão do Homestead. HOMESTEAD_CATALOG_URL sobrepõe o valor gravado.") + "\n")
	b.WriteString(theme.Help.Render("Tema: paleta clara ou escura (Lipgloss) para todo o TUI.") + "\n")
	b.WriteString(theme.Help.Render("Acessibilidade: contraste alto, menos animação (spinner/progress), escala de texto (espaçamento e listas).") + "\n\n")
	if m.catalogEnvOverride {
		b.WriteString(theme.Help.Render("Agora: HOMESTEAD_CATALOG_URL está definida — a linha «URL» abaixo não é usada até limpar a variável.") + "\n\n")
	}

	row := func(i int, label, value string) {
		prefix := "  "
		if m.cursor == i && m.mode == settingsBrowse {
			prefix = theme.Selected.Render(" ▸ ")
		} else if m.cursor == i {
			prefix = theme.Selected.Render(" │ ")
		}
		line := prefix + label
		if value != "" {
			line += ": " + value
		}
		b.WriteString(line + "\n")
	}

	urlNote := m.draft.CatalogURL
	if strings.TrimSpace(urlNote) == "" {
		urlNote = "(padrão do repositório)"
	}
	row(settingsRowCatalog, "URL do catálogo (instaladores)", urlNote)

	themeLabel := "escuro"
	if m.draft.Theme == preferences.ThemeLight {
		themeLabel = "claro"
	}
	row(settingsRowTheme, "Tema (claro / escuro)", themeLabel)
	row(settingsRowHighContrast, "Contraste alto", boolLabel(m.draft.HighContrast))
	row(settingsRowReduceMotion, "Reduzir animações", boolLabel(m.draft.ReduceMotion))

	ts := "normal"
	switch m.draft.TextScale {
	case preferences.TextScaleLarge:
		ts = "grande"
	case preferences.TextScaleXLarge:
		ts = "muito grande"
	}
	row(settingsRowTextScale, "Tamanho do texto (espaçamento)", ts)

	sr := m.draft.ScriptRoot
	if strings.TrimSpace(sr) == "" {
		sr = "(directório de trabalho actual)"
	}
	row(settingsRowScriptRoot, "Raiz dos scripts", sr)

	df := m.draft.DotfilesRepo
	if strings.TrimSpace(df) == "" {
		df = preferences.DefaultDotfilesRepo() + " (padrão)"
	}
	row(settingsRowDotfiles, "Repo dotfiles", df)

	row(settingsRowConfirmScript, "Confirmar antes de executar script", boolLabel(m.draft.ConfirmBeforeScript))
	row(settingsRowConfirmPackage, "Confirmar antes de instalar pacote", boolLabel(m.draft.ConfirmBeforePackage))

	saveLabel := "Gravar e aplicar"
	if strings.TrimSpace(m.prefsPath) == "" {
		saveLabel += " (sem caminho — não grava disco)"
	}
	row(settingsRowSave, saveLabel, "")
	row(settingsRowCancel, "Cancelar (Esc)", "")

	b.WriteString("\n")
	if m.mode == settingsEdit {
		b.WriteString(theme.ScriptScreenAccent.Render("Editar (Enter: OK • Esc: voltar)") + "\n")
		b.WriteString(m.input.View() + "\n")
	} else {
		b.WriteString(theme.Help.Render("↑/↓ • Enter: alterar • Esc: sair sem gravar • ?: ajuda") + "\n")
	}
	if strings.TrimSpace(m.statusLine) != "" {
		b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ErrFg())).Render(m.statusLine) + "\n")
	}
	return b.String()
}

func boolLabel(v bool) string {
	if v {
		return "Sim"
	}
	return "Não"
}
