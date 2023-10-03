package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

func handleAction(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	start, _ := m.eb.GetSelectionRange()
	key := msg.String()
	handled := true
	setMode := ModeNormal

	switch key {

	case "x", "y", "s":
		n, err := m.eb.CopySelection()
		if err != nil {
			panic(err)
		}

		if key == "x" || key == "s" {
			// Delete byte under cursor
			m.eb.PreviewChange(&core.Change{Position: start, Removed: int64(n), Data: []byte{}})
			if key == "x" {
				m.eb.CommitChange()
			} else {
				setMode = ModeInsert
			}
		}

	case "p", "P":
		// Paste clipboard
		removed := 0
		clipboard := m.eb.Clipboard

		// If in visual mode, delete selection first
		if m.mode == ModeVisual {
			n, err := m.eb.CopySelection()
			if err != nil {
				panic(err)
			}
			removed = n
		}

		if key == "p" && m.mode == ModeNormal {
			// Paste after cursor
			start++
		}

		m.eb.PreviewChange(&core.Change{Position: start, Removed: int64(removed), Data: clipboard})
		m.eb.CommitChange()
		start += int64(len(clipboard)) - 1

	default:
		handled = false
	}

	if handled {
		m.eb.Cursor = start
		m.eb.SelectionStart = start
		m.SetMode(setMode)
	}

	return m, nil
}
