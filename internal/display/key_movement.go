package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func handleCursorMovement(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "up", "k":
		// Move cursor up
		if m.eb.Cursor >= int64(m.ncols) {
			m.eb.Cursor -= int64(m.ncols)
		}

	case "down", "j":
		// Move cursor down
		m.eb.Cursor += int64(m.ncols)

	case "left", "h":
		// Move cursor left
		if m.eb.Cursor > 0 {
			m.eb.Cursor--
		}

	case "right", "l":
		// Move cursor right
		m.eb.Cursor++

	case "0", "home":
		// Move cursor to start of line
		m.eb.Cursor -= m.eb.Cursor % int64(m.ncols)

	case "$", "end":
		// Move cursor to end of line
		m.eb.Cursor += int64(m.ncols) - (m.eb.Cursor % int64(m.ncols)) - 1

	case "g":
		m.eb.Cursor = 0

	case "G":
		m.eb.Cursor = m.eb.Size() - 1

	case "ctrl+d", "pgdown":
		m.eb.Cursor += int64(m.ncols) * int64(m.nrows)

	case "ctrl+u", "pgup":
		m.eb.Cursor -= int64(m.ncols) * int64(m.nrows)
	}

	// Cursor bounds check
	if m.eb.Cursor < 0 {
		m.eb.Cursor = 0
	} else if m.eb.Cursor >= m.eb.Size() {
		m.eb.Cursor = m.eb.Size() - 1
	}

	// Scroll cursor into view
	m.ScrollToCursor()

	return m, nil
}
