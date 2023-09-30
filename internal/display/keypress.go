package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

func handleCursorMovement(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	case "up", "k":
		if m.Eb.Cursor >= int64(m.Ncols) {
			m.Eb.Cursor -= int64(m.Ncols)
		}

	case "down", "j":
		m.Eb.Cursor += int64(m.Ncols)

	case "left", "h":
		if m.Eb.Cursor > 0 {
			m.Eb.Cursor--
		}

	case "right", "l":
		m.Eb.Cursor++

	case "0":
		m.Eb.Cursor -= m.Eb.Cursor % int64(m.Ncols)

	case "$":
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

	case "ctrl+d":
		m.Eb.Cursor += int64(m.Ncols) * int64(m.Nrows)

	case "ctrl+u":
		m.Eb.Cursor -= int64(m.Ncols) * int64(m.Nrows)
	}

	// Cursor bounds check
	if m.Eb.Cursor < 0 {
		m.Eb.Cursor = 0
	} else if m.Eb.Cursor >= m.Eb.Size() {
		m.Eb.Cursor = m.Eb.Size() - 1
	}

	// Scroll cursor into view
	viewStart := int64(m.ViewRow) * int64(m.Ncols)
	viewEnd := viewStart + int64(m.Ncols)*int64(m.Nrows)
	if m.Eb.Cursor < viewStart {
		m.ViewRow = int(m.Eb.Cursor / int64(m.Ncols))
	} else if m.Eb.Cursor >= viewEnd {
		m.ViewRow = int(m.Eb.Cursor/int64(m.Ncols)) - m.Nrows + 1
	}

	return m, nil
}

func handleAction(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "x" key deletes the character under the cursor
	case "x":
		m.Eb.PreviewChange(&core.Change{Position: m.Eb.SelectionStart, Removed: 1 + m.Eb.Cursor - m.Eb.SelectionStart, Data: []byte{}})
		m.Eb.CommitChange()
		m.Eb.Cursor = m.Eb.SelectionStart
		m.Mode = ModeNormal

	}

	return m, nil
}

func handleKeypressNormal(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "i" key enters insert mode
	case "i":
		m.Mode = ModeInsert
		m.Eb.PreviewChange(&core.Change{Position: m.Eb.Cursor, Removed: 0, Data: []byte{}})

	// The "v" key enters visual mode
	case "v":
		m.Mode = ModeVisual

	// The ":" key enters command mode
	case ":":
		m.Mode = ModeCommand

	// These keys should exit the program.
	case "ctrl+c", "q":
		return m, tea.Quit

	// The "u" key undoes the last change
	case "u":
		m.Eb.Undo()

	// The "ctrl+r" key redoes the last change
	case "ctrl+r":
		m.Eb.Redo()

	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)
	m.Eb.SelectionStart = m.Eb.Cursor

	return m, nil
}

func handleKeypressInsert(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits insert mode
	case "esc":
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
