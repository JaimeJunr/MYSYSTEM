package theme

import (
	"github.com/charmbracelet/lipgloss"
)

func init() {
	applyPalette(darkPal)
}

const (
	VariantDark  = "dark"
	VariantLight = "light"
)

type palette struct {
	titleFg, helpFg, borderFg, yesFg, noFg, selectedBg, selectedFg string
	scriptLogBg, scriptLogFg, footerFg, footerBg, errFg             string
}

var darkPal = palette{
	titleFg: "205", helpFg: "241", borderFg: "63", yesFg: "10", noFg: "9",
	selectedBg: "63", selectedFg: "230",
	scriptLogBg: "236", scriptLogFg: "252", footerFg: "241", footerBg: "235",
	errFg: "9",
}

var lightPal = palette{
	titleFg: "125", helpFg: "240", borderFg: "62", yesFg: "22", noFg: "160",
	selectedBg: "62", selectedFg: "230",
	scriptLogBg: "252", scriptLogFg: "235", footerFg: "240", footerBg: "254",
	errFg: "160",
}

// Palettes for high contrast (WCAG-oriented terminal colours).
var highContrastDarkPal = palette{
	titleFg: "15", helpFg: "15", borderFg: "15", yesFg: "10", noFg: "9",
	selectedBg: "15", selectedFg: "0",
	scriptLogBg: "0", scriptLogFg: "15", footerFg: "15", footerBg: "0",
	errFg: "9",
}

var highContrastLightPal = palette{
	titleFg: "0", helpFg: "0", borderFg: "0", yesFg: "22", noFg: "160",
	selectedBg: "0", selectedFg: "231",
	scriptLogBg: "231", scriptLogFg: "0", footerFg: "0", footerBg: "231",
	errFg: "160",
}

var currentBorderFg, currentErrFg string

func BorderColor() string {
	return currentBorderFg
}

func ErrFg() string {
	return currentErrFg
}

func applyPalette(p palette) {
	currentBorderFg = p.borderFg
	currentErrFg = p.errFg

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(p.titleFg)).
		MarginBottom(1)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.helpFg))

	ConfirmBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.borderFg)).
		Padding(1, 2).
		Width(60)

	Yes = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.yesFg)).
		Bold(true)

	No = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.noFg)).
		Bold(true)

	Selected = lipgloss.NewStyle().
		Background(lipgloss.Color(p.selectedBg)).
		Foreground(lipgloss.Color(p.selectedFg)).
		Bold(true).
		Padding(0, 1)

	ScriptScreenOuter = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.borderFg)).
		Padding(1, 2)

	ScriptScreenAccent = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.borderFg)).
		Bold(true)

	ScriptLogArea = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(p.borderFg)).
		Padding(0, 1).
		Background(lipgloss.Color(p.scriptLogBg)).
		Foreground(lipgloss.Color(p.scriptLogFg))

	ScriptScreenFooterBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.footerFg)).
		Background(lipgloss.Color(p.footerBg)).
		Padding(0, 1)
}

func ErrColor(variant string) string {
	if currentErrFg != "" {
		return currentErrFg
	}
	if variant == VariantLight {
		return lightPal.errFg
	}
	return darkPal.errFg
}
