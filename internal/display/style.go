package display

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hizkifw/gex/pkg/core"
)

var (
	fgPrimaryColor   = lipgloss.Color("#eeeeee")
	fgSecondaryColor = lipgloss.Color("#999999")
	fgDirtyColor     = lipgloss.Color("#ffff00")
	bgSelectedColor  = lipgloss.Color("#1E3A8A")
	bgCursorColor    = lipgloss.Color("#1D4ED8")

	addrStyle = lipgloss.NewStyle().
			Foreground(fgSecondaryColor).
			Align(lipgloss.Right).
			PaddingRight(1)
)

func MakeStyle(primary bool, activeRegions []core.Region) lipgloss.Style {
	style := lipgloss.NewStyle()

	if primary {
		style = style.Foreground(fgPrimaryColor)
	} else {
		style = style.Foreground(fgSecondaryColor)
	}

	for _, r := range activeRegions {
		switch r.Type {
		case core.RegionTypeSelection:
			style = style.Background(bgSelectedColor)
		case core.RegionTypeCursor:
			style = style.Background(bgCursorColor)
		case core.RegionTypeDirty:
			style = style.Foreground(fgDirtyColor)
		}
	}

	return style
}
