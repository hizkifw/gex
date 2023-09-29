package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/internal/display"
	"github.com/hizkifw/gex/pkg/core"
)

type model struct {
	eb            *core.EditorBuffer
	width, height int
	mode          editingMode
}

type editingMode int

const (
	modeNormal editingMode = iota
	modeInsert
	modeVisual
	modeCommand
)

func initialModel(fname string) model {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	return model{
		eb:     core.NewEditorBuffer(fname, f),
		width:  80,
		height: 24,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func handleCursorMovement(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "up" and "k" keys move the cursor up
	case "up", "k":
		if m.eb.Cursor >= int64(16) {
			m.eb.Cursor -= int64(16)
		}

	// The "down" and "j" keys move the cursor down
	case "down", "j":
		m.eb.Cursor += int64(16)

	// The "left" and "h" keys move the cursor left
	case "left", "h":
		if m.eb.Cursor > 0 {
			m.eb.Cursor--
		}

	// The "right" and "l" keys move the cursor right
	case "right", "l":
		m.eb.Cursor++
	}
	return m, nil
}

func handleAction(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "x" key deletes the character under the cursor
	case "x":
		m.eb.PreviewChange(&core.Change{Position: m.eb.SelectionStart, Removed: 1 + m.eb.Cursor - m.eb.SelectionStart, Data: []byte{}})
		m.eb.CommitChange()
		m.eb.Cursor = m.eb.SelectionStart
		m.mode = modeNormal

	}

	return m, nil
}

func handleKeypressNormal(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "i" key enters insert mode
	case "i":
		m.mode = modeInsert
		m.eb.PreviewChange(&core.Change{Position: m.eb.Cursor, Removed: 0, Data: []byte{}})

	// The "v" key enters visual mode
	case "v":
		m.mode = modeVisual

	// The ":" key enters command mode
	case ":":
		m.mode = modeCommand

	// These keys should exit the program.
	case "ctrl+c", "q":
		return m, tea.Quit

	// The "u" key undoes the last change
	case "u":
		m.eb.Undo()

	// The "ctrl+r" key redoes the last change
	case "ctrl+r":
		m.eb.Redo()

	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)
	m.eb.SelectionStart = m.eb.Cursor

	return m, nil
}

func handleKeypressInsert(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits insert mode
	case "esc":
		m.mode = modeNormal
		m.eb.CommitChange()
	}

	return m, nil
}

func handleKeypressVisual(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits visual mode
	case "esc":
		m.mode = modeNormal
		m.eb.SelectionStart = m.eb.Cursor
	}

	m, _ = handleAction(m, msg)
	m, _ = handleCursorMovement(m, msg)

	return m, nil
}

func handleKeypressCommand(m model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {

	// The "esc" key exits command mode
	case "esc":
		m.mode = modeNormal
	}
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Get the current terminal size on resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Is it a key press?
	case tea.KeyMsg:

		switch m.mode {
		case modeNormal:
			return handleKeypressNormal(m, msg)

		case modeInsert:
			return handleKeypressInsert(m, msg)

		case modeVisual:
			return handleKeypressVisual(m, msg)

		case modeCommand:
			return handleKeypressCommand(m, msg)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := m.eb.Name + "\n\n"

	// The body
	nrows := m.height - 7
	if nrows < 0 {
		nrows = 0
	}
	v, err := display.RenderView(m.eb.ReadSeeker(), 16, nrows, 0, m.eb.SelectionStart, m.eb.Cursor)
	if err != nil {
		v = err.Error()
	}
	s += v

	// The footer
	s += fmt.Sprintf("\nCursor: %d, Changes: %d\n", m.eb.Cursor, len(m.eb.UndoStack))
	switch m.mode {
	case modeNormal:
		s += "NORMAL"
	case modeInsert:
		s += "INSERT"
	case modeVisual:
		s += "VISUAL"
	case modeCommand:
		s += "COMMAND"
	}

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel(os.Args[1]))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
