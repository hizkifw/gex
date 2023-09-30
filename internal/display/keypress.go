package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

func handleCursorMovement(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "up", "k":
		// Move cursor up
		if m.Eb.Cursor >= int64(m.Ncols) {
			m.Eb.Cursor -= int64(m.Ncols)
		}

	case "down", "j":
		// Move cursor down
		m.Eb.Cursor += int64(m.Ncols)

	case "left", "h":
		// Move cursor left
		if m.Eb.Cursor > 0 {
			m.Eb.Cursor--
		}

	case "right", "l":
		// Move cursor right
		m.Eb.Cursor++

	case "0", "home":
		// Move cursor to start of line
		m.Eb.Cursor -= m.Eb.Cursor % int64(m.Ncols)

	case "$", "end":
		// Move cursor to end of line
		m.Eb.Cursor += int64(m.Ncols) - (m.Eb.Cursor % int64(m.Ncols)) - 1

	case "w":
		m.Eb.Cursor += int64(2)

	case "W":
		m.Eb.Cursor += int64(4)

	case "b":
		m.Eb.Cursor -= int64(2)

	case "B":
		m.Eb.Cursor -= int64(4)

	case "g":
		m.Eb.Cursor = 0

	case "G":
		m.Eb.Cursor = m.Eb.Size() - 1

	case "ctrl+d", "pgdown":
		m.Eb.Cursor += int64(m.Ncols) * int64(m.Nrows)

	case "ctrl+u", "pgup":
		m.Eb.Cursor -= int64(m.Ncols) * int64(m.Nrows)
	}

	// Cursor bounds check
	if m.Eb.Cursor < 0 {
		m.Eb.Cursor = 0
	} else if m.Eb.Cursor >= m.Eb.Size() {
		m.Eb.Cursor = m.Eb.Size() - 1
	}

	// Scroll cursor into view
	m.ScrollToCursor()

	return m, nil
}

func handleAction(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	start, _ := m.Eb.GetSelectionRange()
	key := msg.String()
	handled := true

	switch key {

	case "x", "y":
		n, err := m.Eb.CopySelection()
		if err != nil {
			panic(err)
		}

		if key == "x" {
			// Delete byte under cursor
			m.Eb.PreviewChange(&core.Change{Position: start, Removed: int64(n), Data: []byte{}})
			m.Eb.CommitChange()
		}

	case "p", "P":
		// Paste clipboard
		removed := 0
		clipboard := m.Eb.Clipboard

		// If in visual mode, delete selection first
		if m.Mode == ModeVisual {
			n, err := m.Eb.CopySelection()
			if err != nil {
				panic(err)
			}
			removed = n
		}

		if key == "p" && m.Mode == ModeNormal {
			// Paste after cursor
			start++
		}

		m.Eb.PreviewChange(&core.Change{Position: start, Removed: int64(removed), Data: clipboard})
		m.Eb.CommitChange()
		start += int64(len(clipboard)) - 1

	default:
		handled = false
	}

	if handled {
		m.Eb.Cursor = start
		m.Eb.SelectionStart = start
		m.Mode = ModeNormal
	}

	return m, nil
}

func handleKeypressNormal(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "i":
		// Enter insert mode
		m.Mode = ModeInsert
		m.Eb.PreviewChange(&core.Change{Position: m.Eb.Cursor, Removed: 0, Data: []byte{}})

	case "a":
		// Enter insert mode after cursor
		m.Eb.Cursor++
		m.Mode = ModeInsert
		m.Eb.PreviewChange(&core.Change{Position: m.Eb.Cursor, Removed: 0, Data: []byte{}})

	case "v":
		// Enter visual mode
		m.Mode = ModeVisual

	case ":":
		// Enter command mode
		m.Mode = ModeCommand

	case "ctrl+c", "q":
		// Exit program
		return m, tea.Quit

	case "u":
		// Undo last change
		m.Eb.Undo()

	case "ctrl+r":
		// Redo last change
		m.Eb.Redo()

	case "tab":
		// Toggle active column
		if m.ActiveColumn == ActiveColumnHex {
			m.ActiveColumn = ActiveColumnAscii
		} else {
			m.ActiveColumn = ActiveColumnHex
		}
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)
	m.Eb.SelectionStart = m.Eb.Cursor

	return m, nil
}

func handleKeypressInsert(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "esc":
		// Exit insert mode
		m.Mode = ModeNormal
		m.Eb.CommitChange()
	}

	return m, nil
}

func handleKeypressVisual(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits visual mode
	case "esc":
		m.Mode = ModeNormal
		m.Eb.SelectionStart = m.Eb.Cursor
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)

	return m, nil
}

func handleKeypressCommand(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits command mode
	case "esc":
		m.Mode = ModeNormal
	}
	return m, nil
}
