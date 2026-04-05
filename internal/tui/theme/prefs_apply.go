package theme

import (
	"github.com/JaimeJunr/Homestead/internal/infrastructure/preferences"
)

var textScaleLevel int

func TextScaleLevel() int {
	return textScaleLevel
}

// ListVerticalReserve is subtracted from terminal height when sizing list views.
func ListVerticalReserve() int {
	switch textScaleLevel {
	case 2:
		return 7
	case 1:
		return 6
	default:
		return 4
	}
}

// ApplyPreferences updates palette, contrast, and text scale from saved preferences.
func ApplyPreferences(p preferences.Preferences) {
	p.Normalize()
	textScaleLevel = 0
	switch p.TextScale {
	case preferences.TextScaleLarge:
		textScaleLevel = 1
	case preferences.TextScaleXLarge:
		textScaleLevel = 2
	}

	var pal palette
	if p.HighContrast {
		if p.Theme == preferences.ThemeLight {
			pal = highContrastLightPal
		} else {
			pal = highContrastDarkPal
		}
	} else if p.Theme == preferences.ThemeLight {
		pal = lightPal
	} else {
		pal = darkPal
	}
	applyPalette(pal)

	if textScaleLevel > 0 {
		Help = Help.PaddingTop(textScaleLevel / 2).PaddingBottom(1 + textScaleLevel/2)
		Title = Title.MarginBottom(1 + textScaleLevel/2)
		ConfirmBox = ConfirmBox.Width(min(72, 60+textScaleLevel*4))
	}
}
