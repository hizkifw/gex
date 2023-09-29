package display

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Style for the hex dump address
	addrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#999999")).
			Width(9).
			Align(lipgloss.Right).
			PaddingRight(1)

	// Style for the hex dump hex values
	hexNormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eeeeee")).
			Width(3).
			PaddingRight(1)

	// Style for the hex dump hex values that are part of the selection
	hexSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#222222")).
				Background(lipgloss.Color("#eeeeee")).
				Width(3).
				PaddingRight(1)

	// Style for the hex dump ASCII values
	asciiNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#999999")).
				Width(1)

	// Style for the hex dump ASCII values that are part of the selection
	asciiSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#222222")).
				Background(lipgloss.Color("#eeeeee")).
				Width(1)
)
