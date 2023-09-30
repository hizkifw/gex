package display

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	fgPrimaryColor   = lipgloss.Color("#eeeeee")
	fgSecondaryColor = lipgloss.Color("#999999")
	bgSelectedColor  = lipgloss.Color("#1E3A8A")
	bgCursorColor    = lipgloss.Color("#1D4ED8")

	addrStyle = lipgloss.NewStyle().
			Foreground(fgSecondaryColor).
			Align(lipgloss.Right).
			PaddingRight(1)
)

func MakeStyle(primary, selected, cursor bool) lipgloss.Style {
	style := lipgloss.NewStyle()

	if primary {
		style = style.Foreground(fgPrimaryColor)
	} else {
		style = style.Foreground(fgSecondaryColor)
	}

	if cursor {
		style = style.Background(bgCursorColor)
	} else if selected {
		style = style.Background(bgSelectedColor)
	}

	return style
}
