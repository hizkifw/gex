package display

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
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
	eb            *core.EditorBuffer
	width, height int
	nrows, ncols  int
	viewRow       int
	mode          EditingMode
	activeColumn  ActiveColumn

	ResponsiveCols bool

	// Command mode text field
	cmdText textinput.Model
	// Temporary buffer for inputs
	tmpText textinput.Model
}

func NewModel() Model {
	m := Model{
		eb:           core.NewEditorBuffer("", EmptyReadSeeker),
		width:        0,
		height:       0,
		nrows:        0,
		ncols:        16,
		viewRow:      0,
		mode:         ModeNormal,
		activeColumn: ActiveColumnHex,

		ResponsiveCols: false,

		cmdText: textinput.New(),
		tmpText: textinput.New(),
	}
	m.SetMode(ModeNormal)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Get the current terminal size on resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.cmdText.Width = msg.Width

		cols, rows := CalculateViewSize(m.width, m.height)
		m.nrows = rows
		if m.ResponsiveCols {
			m.ncols = cols
		}
		m.ScrollToCursor()

	// Handle keypresses
	case tea.KeyMsg:

		switch m.mode {
		case ModeNormal:
			return HandleKeypressNormal(m, msg)

		case ModeInsert:
			return HandleKeypressInsert(m, msg)

		case ModeVisual:
			return HandleKeypressVisual(m, msg)

		case ModeReplace:
			return HandleKeypressInsert(m, msg)

		case ModeCommand:
			return HandleKeypressCommand(m, msg)
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

	return fmt.Sprintf("%s\n%s\n%s", v, m.mode, m.cmdText.View())
}

func (m *Model) LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", name, err)
	}
	m.eb = core.NewEditorBuffer(name, f)
	return nil
}

// ScrollToCursor scrolls the view so that the cursor is visible.
func (m *Model) ScrollToCursor() {
	viewStart := int64(m.viewRow) * int64(m.ncols)
	viewEnd := viewStart + int64(m.ncols)*int64(m.nrows)
	if m.eb.Cursor < viewStart {
		m.viewRow = int(m.eb.Cursor / int64(m.ncols))
	} else if m.eb.Cursor >= viewEnd {
		m.viewRow = int(m.eb.Cursor/int64(m.ncols)) - m.nrows + 1
	}
}

// SetMode sets the editing mode.
func (m *Model) SetMode(mode EditingMode) {
	m.mode = mode

	m.cmdText.SetValue("")
	m.tmpText.SetValue("")
	if mode == ModeCommand {
		m.cmdText.Prompt = ":"
		m.cmdText.Focus()
	} else {
		m.cmdText.Prompt = ""
		m.cmdText.Blur()
	}

	if mode == ModeInsert || mode == ModeReplace {
		m.tmpText.Focus()
	} else {
		m.tmpText.Blur()
	}
}

// StatusMessage sets the status message.
func (m *Model) StatusMessage(msg string) {
	m.cmdText.SetValue(msg)
}
