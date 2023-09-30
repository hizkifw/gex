package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func HandleKeypressInsert(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	if key == "esc" {
		// Exit insert mode
		m.SetMode(ModeNormal)
		m.eb.CommitChange()
		return m, nil
	}

	if m.activeColumn == ActiveColumnAscii {
		// Handle inserting characters.

	} else if m.activeColumn == ActiveColumnHex {
		// Handle inserting hex characters.

	}

	return m, nil
}
