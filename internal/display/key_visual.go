package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func HandleKeypressVisual(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "esc":
		// Exit visual mode
		m.SetMode(ModeNormal)
		m.eb.SelectionStart = m.eb.Cursor

	case ":":
		// Enter command mode
		m.SetMode(ModeCommand)
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)

	return m, nil
}
