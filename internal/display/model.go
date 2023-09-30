package display

import (
	"bytes"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
)

type EditingMode string
type ActiveColumn int

const (
	ModeNormal  EditingMode = "NORMAL"
	ModeInsert  EditingMode = "INSERT"
	ModeVisual  EditingMode = "VISUAL"
	ModeReplace EditingMode = "REPLACE"
	ModeCommand EditingMode = "COMMAND"

	ActiveColumnHex ActiveColumn = iota
	ActiveColumnAscii
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
	ActiveColumn  ActiveColumn

	ResponsiveCols bool
}

func NewModel() Model {
	return Model{
		Eb:           core.NewEditorBuffer("", EmptyReadSeeker),
		Width:        0,
		Height:       0,
		Nrows:        0,
		Ncols:        16,
		ViewRow:      0,
		Mode:         ModeNormal,
		ActiveColumn: ActiveColumnHex,

		ResponsiveCols: false,
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
		m.ScrollToCursor()

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
	v, err := m.RenderHexView()
	if err != nil {
		v = err.Error()
	}

	return fmt.Sprintf("%s\n%s\n", v, m.Mode)
}

func (m *Model) LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", name, err)
	}
	m.Eb = core.NewEditorBuffer(name, f)
	return nil
}

// ScrollToCursor scrolls the view so that the cursor is visible.
func (m *Model) ScrollToCursor() {
	viewStart := int64(m.ViewRow) * int64(m.Ncols)
	viewEnd := viewStart + int64(m.Ncols)*int64(m.Nrows)
	if m.Eb.Cursor < viewStart {
		m.ViewRow = int(m.Eb.Cursor / int64(m.Ncols))
	} else if m.Eb.Cursor >= viewEnd {
		m.ViewRow = int(m.Eb.Cursor/int64(m.Ncols)) - m.Nrows + 1
	}
}
