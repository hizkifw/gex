package display

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hizkifw/gex/pkg/core"
)

var (
	fgPrimaryColor    = lipgloss.Color("#eeeeee")
	fgSecondaryColor  = lipgloss.Color("#999999")
	fgDirtyColor      = lipgloss.Color("#4ade80")
	bgSelectedColor   = lipgloss.Color("#1e3a8a")
	bgCursorColor     = lipgloss.Color("#1d4ed8")
	bgEditingColor    = lipgloss.Color("#7e22ce")
	bgStatusModeColor = lipgloss.Color("#444444")
	bgStatusBarColor  = lipgloss.Color("#222222")
	bgErrorColor      = lipgloss.Color("#ff5555")

	addrStyle = lipgloss.NewStyle().
			Foreground(fgSecondaryColor).
			Align(lipgloss.Right).
			PaddingRight(1)

	textErrorStyle = lipgloss.NewStyle().
			Foreground(bgErrorColor).
			Bold(true)
	statusDefaultStyle = lipgloss.NewStyle().
				Foreground(fgPrimaryColor).
				Background(bgStatusModeColor).
				Padding(0, 1, 0, 1).
				Bold(true)
	statusEditingStyle = lipgloss.NewStyle().
				Foreground(fgPrimaryColor).
				Background(bgEditingColor).
				PaddingLeft(1).
				PaddingRight(1).
				Bold(true)
	statusBarStyle = lipgloss.NewStyle().
			Foreground(fgPrimaryColor).
			Background(bgStatusBarColor)
	windowStyle = lipgloss.NewStyle().
			Foreground(fgPrimaryColor).
			Background(bgStatusBarColor).
			BorderForeground(fgPrimaryColor).
			Padding(1, 2, 1, 2)
	windowTitleStyle = lipgloss.NewStyle().
				Foreground(fgPrimaryColor).
				Background(bgStatusModeColor).
				Padding(0, 1, 0, 1).
				Bold(true)

	padLeftStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	statusStyle = map[EditingMode]lipgloss.Style{
		ModeNormal:  statusDefaultStyle,
		ModeVisual:  statusEditingStyle,
		ModeInsert:  statusEditingStyle,
		ModeReplace: statusEditingStyle,
		ModeCommand: statusDefaultStyle,
	}
)

func MakeStyle(primary bool, isEditing bool, activeRegions []core.Region) lipgloss.Style {
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
			if isEditing {
				style = style.Background(bgEditingColor)
			} else {
				style = style.Background(bgCursorColor)
			}
		case core.RegionTypeDirty:
			style = style.Foreground(fgDirtyColor)
		}
	}

	return style
}
