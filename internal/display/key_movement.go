package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func handleCursorMovement(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "up", "k":
		// Move cursor up
		if m.eb.Cursor >= int64(m.ncols) {
			m.MoveCursor(-int64(m.ncols))
		}

	case "down", "j":
		// Move cursor down
		m.MoveCursor(int64(m.ncols))

	case "left", "h":
		// Move cursor left
		if m.eb.Cursor > 0 {
			m.MoveCursor(-1)
		}

	case "right", "l":
		// Move cursor right
		m.MoveCursor(1)

	case "0", "home":
		// Move cursor to start of line
		m.MoveCursor(-(m.eb.Cursor % int64(m.ncols)))

	case "$", "end":
		// Move cursor to end of line
		m.MoveCursor(int64(m.ncols) - (m.eb.Cursor % int64(m.ncols)) - 1)

	case "g":
		m.SetMode(ModeCommand)
		m.cmdText.SetValue("goto ")

	case "G":
		m.SetCursor(m.eb.Size() - 1)

	case "ctrl+d", "pgdown":
		m.MoveCursor(int64(m.ncols) * int64(m.nrows))

	case "ctrl+u", "pgup":
		m.MoveCursor(-int64(m.ncols) * int64(m.nrows))
	}

	return m, nil
}
