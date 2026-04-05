package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/domain/entities"
	"github.com/JaimeJunr/Homestead/internal/tui/cmds"
	"github.com/JaimeJunr/Homestead/internal/tui/items"
	btmsg "github.com/JaimeJunr/Homestead/internal/tui/msg"
	"github.com/JaimeJunr/Homestead/internal/tui/sysurl"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Iniciando..."
	}

	var body string
	switch m.state {
	case ViewMainMenu:
		body = m.mainMenu.View() + "\n" + theme.Help.Render("?: ajuda contextual · q: sair")

	case ViewScriptList:
		helpLine := "\n↑/↓: navegar • /: buscar • f: favorito • esc: voltar (ou limpar busca) • enter: executar • d: simular (dry-run) • ?: ajuda • q: sair"
		help := theme.Help.Render(helpLine)
		var feedback string
		if m.err != nil {
			feedback = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ErrFg())).
				Render("\n⚠ "+m.err.Error()) + "\n"
		}
		body = m.scriptList.View() + feedback + help

	case ViewInstallerCategories:
		body = m.installerList.View() + theme.Help.Render("\n↑/↓: navegar • enter: abrir • esc: voltar • ?: ajuda • q: sair")

	case ViewPackageList:
		help := theme.Help.Render("\n↑/↓: navegar • /: buscar • enter: confirmação • o: abrir URL • c: copiar URL • esc: voltar (ou limpar busca) • ?: ajuda • q: sair")
		toast := ""
		if strings.TrimSpace(m.keyboardToast) != "" {
			toast = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n"
		}
		body = m.packageList.View() + toast + help

	case ViewConfirmation:
		body = m.renderConfirmation()

	case ViewInstalling:
		body = m.renderInstallProgress()

	case ViewZshWizard:
		if m.zshWizard != nil {
			body = m.zshWizard.View()
		} else {
			body = "Iniciando wizard..."
		}

	case ViewZshApplying:
		body = m.renderZshApplyFeedback()

	case ViewZshRepoWizard:
		if m.zshRepoWizard != nil {
			body = m.zshRepoWizard.View()
			if strings.TrimSpace(m.keyboardToast) != "" {
				body += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n"
			}
		} else {
			body = "Iniciando Configurar Zsh..."
		}

	case ViewScriptOutput:
		body = m.renderScriptOutput()

	case ViewNativeMonitor:
		body = m.renderNativeMonitorView()

	case ViewSettings:
		if m.settingsModel != nil {
			body = m.settingsModel.View()
		} else {
			body = ""
		}

	default:
		body = ""
	}
	return m.maybeWrapHelp(body)
}

func (m Model) renderConfirmation() string {
	var title, description string

	switch item := m.selectedItem.(type) {
	case entities.Script:
		if item.NativeMonitor != "" {
			title = "Abrir monitor?"
			description = fmt.Sprintf("Você deseja abrir:\n\n  %s\n  %s", item.Name, item.Description)
		} else {
			title = "Executar script?"
			description = fmt.Sprintf("Você deseja executar:\n\n  %s\n  %s", item.Name, item.Description)
			if item.RequiresSudo {
				description += "\n\n⚠️  Este script requer permissões de administrador (sudo)"
			}
			if item.SupportsDryRun {
				description += "\n\nNa lista de scripts, tecla d: simulação (dry-run) antes de executar."
			}
		}
	case entities.Package:
		title = "Instalar pacote?"
		description = fmt.Sprintf("Você deseja instalar:\n\n  %s\n  %s\n  Versão: %s",
			item.Name, item.Description, item.Version)
		if kb := sysurl.PackageKeyboardURL(item); kb != "" {
			description += "\n\n🔗 Verificação (sem mouse: tecla o abre no navegador, c copia a URL):\n  " + kb
		}
		if strings.TrimSpace(item.UtilityScriptPath) != "" {
			description += "\n\n⚠️  Será executado o script Homestead deste utilitário (mesmo fluxo que os outros instaladores)."
			if item.RequiresSudo {
				description += "\nPode ser pedida senha de administrador (sudo)."
			} else {
				description += "\nSem sudo no wrapper; o script pode alterar só o seu utilizador."
			}
		} else if item.DownloadURL != "" {
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
		yesButton = theme.Selected.Render(" Sim ")
		noButton = theme.No.Render(" Não ")
	} else {
		yesButton = theme.Yes.Render(" Sim ")
		noButton = theme.Selected.Render(" Não ")
	}

	helpConfirm := "←/→: escolher • enter: confirmar • esc: voltar • ?: ajuda"
	if p, ok := m.selectedItem.(entities.Package); ok && sysurl.PackageKeyboardURL(p) != "" {
		helpConfirm = "o: abrir URL • c: copiar URL • " + helpConfirm
	}
	if s, ok := m.selectedItem.(entities.Script); ok && s.SupportsDryRun && s.NativeMonitor == "" {
		helpConfirm = "d na lista: dry-run • " + helpConfirm
	}
	toastLine := ""
	if strings.TrimSpace(m.keyboardToast) != "" {
		toastLine = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render(m.keyboardToast) + "\n\n"
	}
	content := theme.Title.Render(title) + "\n\n" +
		description + "\n\n" +
		yesButton + "  " + noButton + "\n\n" +
		toastLine +
		theme.Help.Render(helpConfirm)

	boxW := m.width - 8
	if boxW < 52 {
		boxW = 52
	}
	if boxW > 88 {
		boxW = 88
	}
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor())).
		Padding(1, 2).
		Width(boxW)
	box := boxStyle.Render(content)

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

func (m Model) renderInstallProgress() string {
	var pkg entities.Package
	if p, ok := m.selectedItem.(entities.Package); ok {
		pkg = p
	}

	title := theme.Title.Render(fmt.Sprintf("Instalando: %s", pkg.Name))

	statusIcons := map[string]string{
		"preparing":   "⏳",
		"downloading": "⬇️ ",
		"installing":  "⚙️ ",
		"complete":    "✅",
		"failed":      "❌",
	}

	icon := statusIcons[m.installStatus]
	status := fmt.Sprintf("%s %s", m.installStatusGlyph(icon), m.installMessage)
	progressBar := m.progress.ViewAs(m.installPercent)

	content := title + "\n\n" +
		status + "\n\n" +
		progressBar + "\n\n"

	if m.canAbort && !m.aborted {
		content += theme.Help.Render("⚠️  Pressione Ctrl+C para abortar (não recomendado) • ?: ajuda")
	} else if m.installStatus == "complete" {
		content += theme.Help.Render("Instalação concluída! Retornando ao menu...")
	} else if m.installStatus == "failed" {
		content += theme.Help.Render("❌ Instalação falhou. Retornando ao menu...")
	} else {
		content += theme.Help.Render("Aguarde... não feche o programa • ?: ajuda")
	}

	box := theme.ConfirmBox.Render(content)

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

func scriptOutputViewportWH(termW, termH int) (w, h int) {
	boxW := scriptOutputCardWidth(termW)
	w = boxW - 8
	if w < 28 {
		w = 28
	}
	h = termH - 20 - 2*theme.TextScaleLevel()
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
	vp.Style = theme.ScriptLogArea
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
	return lipgloss.NewStyle().Foreground(lipgloss.Color(theme.BorderColor())).Render(strings.Repeat("─", n))
}

func (m Model) renderScriptOutput() string {
	boxW := scriptOutputCardWidth(m.width)

	if m.scriptOutputPhase == "running" {
		accent := "📜 Executando script"
		wait := "Capturando saída…"
		note := "A saída aparecerá no painel abaixo quando o script terminar."
		sudoNote := "Scripts com sudo usam o terminal completo para pedir senha."
		head := theme.Title.Render("Homestead") + "\n" +
			theme.Help.Render("Gerenciador de Sistema") + "\n" +
			scriptOutputDivider(boxW) + "\n" +
			theme.ScriptScreenAccent.Render(accent) + "\n" +
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Render(m.scriptOutputTitle)
		body := "\n\n" + fmt.Sprintf("%s %s", m.waitGlyph(), theme.Help.Render(wait))
		body += "\n\n" + theme.Help.Render(note)
		body += "\n" + theme.Help.Render(sudoNote)
		body += "\n" + theme.Help.Render("?: ajuda")
		content := head + body
		box := theme.ScriptScreenOuter.Width(boxW)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
	}

	doneAccent := "📜 Saída do script"
	head := theme.Title.Render("Homestead") + "\n" +
		theme.Help.Render("Gerenciador de Sistema") + "\n" +
		scriptOutputDivider(boxW) + "\n" +
		theme.ScriptScreenAccent.Render(doneAccent)
	nameLine := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Render(m.scriptOutputTitle)
	if m.scriptOutputErr != nil {
		nameLine += "  " + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Render("· falhou")
	}
	view := m.scriptOutputView.View()
	footerText := "↑/↓  PgUp/PgDn  rolar · Enter / Esc / q  voltar · ?: ajuda"
	if m.scriptOutputErr != nil {
		footerText = "Ver mensagem de erro no fim do texto · " + footerText
	}
	footer := theme.ScriptScreenFooterBar.Width(max(12, boxW-8)).Render(footerText)
	content := head + "\n" + nameLine + "\n\n" + view + "\n" + footer
	box := theme.ScriptScreenOuter.Width(boxW)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box.Render(content))
}

func (m Model) renderZshApplyFeedback() string {
	title := theme.Title.Render("Configuração Zsh")

	switch m.zshApplyPhase {
	case "applying":
		status := fmt.Sprintf("%s Aplicando configuração Zsh...", m.waitGlyph())
		content := title + "\n\n" + status + "\n\n" + theme.Help.Render("Aguarde... • ?: ajuda")
		box := theme.ConfirmBox.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	case "success":
		content := title + "\n\n" +
			"✅ Configuração aplicada com sucesso.\n\n" +
			"O arquivo ~/.zshrc foi atualizado com os plugins e ferramentas selecionados.\n" +
			"Criados/atualizados: ~/.zsh/general/aliases.zsh e functions.zsh.\n\n" +
			"Verifique: cat ~/.zshrc\n\n" +
			"Não instala plugins externos (ex.: zsh-autosuggestions); apenas escreve o .zshrc.\n" +
			"Use Instaladores para Zsh/Oh My Zsh se ainda não estiverem instalados.\n\n" +
			theme.Help.Render("Retornando ao menu em 2s (ou Enter/Esc para voltar) • ?: ajuda")
		box := theme.ConfirmBox.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	case "error":
		errMsg := ""
		if m.zshApplyError != nil {
			errMsg = m.zshApplyError.Error()
		}
		content := title + "\n\n" +
			"❌ Erro ao aplicar configuração:\n\n" + errMsg + "\n\n" +
			theme.Help.Render("Pressione Enter ou Esc para voltar ao menu • ?: ajuda")
		box := theme.ConfirmBox.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	default:
		content := title + "\n\n" + m.waitGlyph() + " Aguarde..."
		box := theme.ConfirmBox.Render(content)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
	}
}

func (m Model) urlForKeyboardOpen() string {
	switch m.state {
	case ViewConfirmation:
		if p, ok := m.selectedItem.(entities.Package); ok {
			return sysurl.PackageKeyboardURL(p)
		}
	case ViewPackageList:
		if sel := m.packageList.SelectedItem(); sel != nil {
			if it, ok := sel.(items.PackageItem); ok {
				return sysurl.PackageKeyboardURL(it.Pkg)
			}
		}
	}
	return ""
}

func (m Model) handleURLShortcut(wantCopy bool) (Model, tea.Cmd) {
	url := m.urlForKeyboardOpen()
	if url != "" {
		if wantCopy {
			return m, cmds.CopyURL(url)
		}
		return m, cmds.OpenURL(url)
	}
	if m.state == ViewPackageList || (m.state == ViewConfirmation && isSelectedPackage(m.selectedItem)) {
		m.keyboardToast = "Este pacote não tem URL de projeto nem download."
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return btmsg.ClearKeyboardToast{} })
	}
	return m, nil
}

func isSelectedPackage(item interface{}) bool {
	_, ok := item.(entities.Package)
	return ok
}
