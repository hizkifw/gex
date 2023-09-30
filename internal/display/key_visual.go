package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func HandleKeypressVisual(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits visual mode
	case "esc":
		m.SetMode(ModeNormal)
		m.eb.SelectionStart = m.eb.Cursor
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)

	return m, nil
}
