package theme

import (
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/JaimeJunr/Homestead/internal/domain/types"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func StripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

var (
	Title                 lipgloss.Style
	Help                  lipgloss.Style
	ConfirmBox            lipgloss.Style
	Yes                   lipgloss.Style
	No                    lipgloss.Style
	Selected              lipgloss.Style
	ScriptScreenOuter     lipgloss.Style
	ScriptScreenAccent    lipgloss.Style
	ScriptLogArea         lipgloss.Style
	ScriptScreenFooterBar lipgloss.Style
)

func InstallerBreadcrumb(segment string) string {
	return "📦 Instaladores > " + segment
}

func InstallerPackageSectionTitle(c types.PackageCategory) string {
	switch c {
	case types.PackageCategoryIDE:
		return "💻 IDEs e Editores"
	case types.PackageCategoryTool:
		return "🔧 Ferramentas de Desenvolvimento"
	case types.PackageCategoryApp:
		return "📱 Aplicações"
	case types.PackageCategoryZshCore:
		return "🐚 Componentes Core (Zsh)"
	case types.PackageCategoryTerminal:
		return "🖥️ Emuladores de Terminal"
	case types.PackageCategoryShell:
		return "🐚 Shells Alternativos"
	case types.PackageCategoryAI:
		return "🤖 Integração com IA"
	case types.PackageCategoryGames:
		return "🎮 Games"
	case types.PackageCategorySysAdmin:
		return "🛡️ Administração de sistemas"
	case types.PackageCategoryUtilities:
		return "🧰 Utilitários"
	default:
		return "📦"
	}
}
