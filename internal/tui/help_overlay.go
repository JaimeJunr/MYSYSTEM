package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/tui/theme"
)

func (m Model) contextualHelpLines() []string {
	title := "Atalhos deste ecrã"
	var lines []string
	switch m.state {
	case ViewMainMenu:
		lines = []string{
			"↑ / ↓ ou j / k — mover",
			"Enter — abrir",
			"q — sair (menu principal)",
			"Ctrl+C — sair",
			"? — esta ajuda",
		}
	case ViewScriptList:
		lines = []string{
			"↑ / ↓ — mover",
			"/ — filtrar (outra vez Esc limpa)",
			"f — favorito",
			"d — próxima execução em simulação (dry-run)",
			"Enter — executar ou abrir monitor",
			"Esc — voltar",
			"q — sair (no menu principal)",
			"? — esta ajuda",
		}
	case ViewInstallerCategories:
		lines = []string{
			"↑ / ↓ — mover",
			"Enter — abrir categoria",
			"Esc — voltar ao menu",
			"q — sair",
			"? — esta ajuda",
		}
	case ViewPackageList:
		lines = []string{
			"↑ / ↓ — mover",
			"/ — filtrar",
			"Enter — instalar (com confirmação se activo)",
			"o — abrir URL no navegador",
			"c — copiar URL",
			"Esc — voltar",
			"q — sair",
			"? — esta ajuda",
		}
	case ViewConfirmation:
		lines = []string{
			"← / → ou h / l — Sim / Não",
			"Enter — confirmar escolha",
			"Esc — cancelar",
			"o / c — abrir ou copiar URL (se o pacote tiver)",
			"? — esta ajuda",
		}
	case ViewScriptOutput:
		if m.scriptOutputPhase == "running" {
			lines = []string{
				"Aguarde o fim do script…",
				"(sudo pode usar o terminal completo)",
				"? — esta ajuda",
			}
		} else {
			lines = []string{
				"↑ / ↓ / PgUp / PgDn — rolar saída",
				"Enter / Esc / q — voltar",
				"? — esta ajuda",
			}
		}
	case ViewNativeMonitor:
		lines = []string{
			"r — atualizar já",
			"Enter / Esc / q — voltar",
			"? — esta ajuda",
		}
	case ViewInstalling:
		lines = []string{
			"Ctrl+C — abortar (quando disponível)",
			"Aguarde o fim da instalação",
			"? — esta ajuda",
		}
	case ViewZshWizard:
		if m.zshWizard != nil {
			switch m.zshWizard.currentView {
			case ZshWizardViewReview:
				lines = []string{
					"Enter ou n — confirmar e aplicar",
					"Esc — voltar",
					"Ctrl+C — cancelar",
					"? — esta ajuda",
				}
			default:
				lines = []string{
					"↑ / ↓ — mover",
					"Espaço — alternar selecção",
					"a — marcar todos",
					"n / Tab / → — próximo passo",
					"Esc — voltar",
					"Ctrl+C — sair",
					"? — esta ajuda",
				}
			}
		} else {
			lines = []string{"? — esta ajuda"}
		}
	case ViewZshApplying:
		lines = []string{
			"Enter / Esc — voltar (quando aparecer)",
			"? — esta ajuda",
		}
	case ViewZshRepoWizard:
		lines = []string{
			"Esc — voltar um passo ou sair",
			"y / s, n — sim / não (repo existente)",
			"g / u — GitHub automático / colar URL",
			"b / r — backup / restaurar (painel configurado)",
			"Enter — enviar URL ou nome (campos de texto)",
			"Ctrl+C — sair",
			"? — esta ajuda",
		}
	case ViewSettings:
		lines = []string{
			"↑ / ↓ — mover",
			"Enter — editar / alternar / gravar",
			"Esc — cancelar edição ou sair sem gravar",
			"? — ajuda (na edição de URL/caminho, ? é texto normal)",
		}
	default:
		lines = []string{"? — fechar ajuda"}
	}

	out := make([]string, 0, len(lines)+4)
	out = append(out, title+":", "")
	out = append(out, lines...)
	out = append(out, "", "Esc ou ? — fechar")
	return out
}

func (m Model) maybeWrapHelp(body string) string {
	if !m.helpOpen {
		return body
	}
	return renderHelpOverlay(m.width, m.height, body, m.contextualHelpLines())
}

func renderHelpOverlay(termW, termH int, background string, lines []string) string {
	panelW := termW - 6
	if panelW < 44 {
		panelW = 44
	}
	if panelW > 78 {
		panelW = 78
	}

	content := theme.Title.Render("❓ Ajuda") + "\n\n" +
		theme.Help.Render(strings.Join(lines, "\n"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderColor())).
		Padding(1, 2).
		Width(panelW).
		Render(content)

	bgLines := strings.Split(background, "\n")
	maxBg := termH - lipgloss.Height(box) - 2
	if maxBg < 3 {
		maxBg = 3
	}
	if len(bgLines) > maxBg {
		bgLines = bgLines[:maxBg]
	}
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	prefix := ""
	if termH > 12 && len(bgLines) > 0 {
		for i := range bgLines {
			bgLines[i] = dim.Render(theme.StripANSI(bgLines[i]))
		}
		prefix = strings.Join(bgLines, "\n") + "\n\n"
	}
	stack := prefix + lipgloss.Place(termW, lipgloss.Height(box)+2, lipgloss.Center, lipgloss.Top, box)
	return lipgloss.Place(termW, termH, lipgloss.Left, lipgloss.Top, stack)
}

func (m Model) waitGlyph() string {
	if m.prefs.ReduceMotion {
		return "…"
	}
	return m.spinner.View()
}

func (m Model) installStatusGlyph(icon string) string {
	if m.prefs.ReduceMotion {
		if icon != "" {
			return icon
		}
		return "…"
	}
	if icon != "" {
		return icon
	}
	return m.spinner.View()
}

func (m Model) suppressHelpHotkey() bool {
	if m.state == ViewSettings && m.settingsModel != nil && m.settingsModel.IsEditing() {
		return true
	}
	if m.state == ViewZshRepoWizard && m.zshRepoWizard != nil && m.zshRepoWizard.isTextInputView() {
		return true
	}
	return false
}
