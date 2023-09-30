package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

func HandleKeypressNormal(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	switch key {

	case "i", "a":
		// Enter insert mode

		if key == "a" {
			// Enter insert mode after cursor
			m.eb.Cursor++
		}

		m.SetMode(ModeInsert)
		m.eb.PreviewChange(&core.Change{Position: m.eb.Cursor, Removed: 0, Data: []byte{}})

	case "v":
		// Enter visual mode
		m.SetMode(ModeVisual)

	case ":":
		// Enter command mode
		m.SetMode(ModeCommand)

	case "ctrl+c", "q":
		// Exit program
		return m, tea.Quit

	case "u":
		// Undo last change
		m.eb.Undo()

	case "ctrl+r":
		// Redo last change
		m.eb.Redo()

	case "tab":
		// Toggle active column
		if m.activeColumn == ActiveColumnHex {
			m.activeColumn = ActiveColumnAscii
		} else {
			m.activeColumn = ActiveColumnHex
		}
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)
	m.eb.SelectionStart = m.eb.Cursor

	return m, nil
}
