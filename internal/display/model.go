package display

import (
	"bytes"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

type EditingMode int

const (
	ModeNormal EditingMode = iota
	ModeInsert
	ModeVisual
	ModeCommand
)

var (
	EmptyReadSeeker = io.ReadSeeker(bytes.NewReader([]byte{}))
)

type Model struct {
	Eb            *core.EditorBuffer
	Width, Height int
	Nrows, Ncols  int
	ViewRow       int
	Mode          EditingMode

	ResponsiveCols bool
}

func NewModel() Model {
	return Model{
		Eb:             core.NewEditorBuffer("", EmptyReadSeeker),
		Width:          0,
		Height:         0,
		Nrows:          0,
		Ncols:          16,
		ViewRow:        0,
		ResponsiveCols: true,
		Mode:           ModeNormal,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Get the current terminal size on resize
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		cols, rows := CalculateViewSize(m.Width, m.Height)
		m.Nrows = rows
		if m.ResponsiveCols {
			m.Ncols = cols
		}

	// Handle keypresses
	case tea.KeyMsg:

		switch m.Mode {
		case ModeNormal:
			return handleKeypressNormal(m, msg)

		case ModeInsert:
			return handleKeypressInsert(m, msg)

		case ModeVisual:
			return handleKeypressVisual(m, msg)

		case ModeCommand:
			return handleKeypressCommand(m, msg)
		}
	}

	return m, nil
}

func (m Model) View() string {
	// View
	ncols, nrows := CalculateViewSize(m.Width, m.Height)
	sel1, sel2 := m.Eb.GetSelectionRange()
	v, err := RenderView(m.Eb.ReadSeeker(), ncols, nrows, m.ViewRow, sel1, sel2)
	if err != nil {
		v = err.Error()
	}

	// Status
	var status string
	switch m.Mode {
	case ModeNormal:
		status = "NORMAL"
	case ModeInsert:
		status = "INSERT"
	case ModeVisual:
		status = "VISUAL"
	case ModeCommand:
		status = "COMMAND"
	}

	return fmt.Sprintf("%s\n%s\n", v, status)
}

func (m *Model) LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", name, err)
	}
	m.Eb = core.NewEditorBuffer(name, f)
	return nil
}
