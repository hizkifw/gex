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

	case "R":
		// Enter replace mode
		m.SetMode(ModeReplace)
		m.eb.PreviewChange(&core.Change{Position: m.eb.Cursor, Removed: 0, Data: []byte{}})

	case ":":
		// Enter command mode
		m.SetMode(ModeCommand)

	case "ctrl+c":
		// Tell user how to exit the program
		m.StatusMessage("Press :q! to quit without saving", false)

	case "u":
		// Undo last change
		if !m.eb.Undo() {
			m.StatusMessage("Nothing to undo", false)
		}

	case "ctrl+r":
		// Redo last change
		if !m.eb.Redo() {
			m.StatusMessage("Already at newest change", false)
		}

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
