package display

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
	"github.com/hizkifw/gex/pkg/util"
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
	prevMode      EditingMode
	mode          EditingMode
	activeColumn  ActiveColumn

	ResponsiveCols bool

	// Inspector
	inspectorEnabled   bool
	inspectorByteOrder binary.ByteOrder

	// Status bar error state
	statusError bool
	// Command mode text field
	cmdText textinput.Model
	// Temporary buffer for inputs
	tmpText textinput.Model
	// Command history
	cmdHistory      []string
	cmdHistoryIndex int

	// Statistics
	framesRendered  uint64
	renderTimeAvgNS float64
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

		inspectorEnabled:   true,
		inspectorByteOrder: binary.LittleEndian,

		statusError:     false,
		cmdText:         textinput.New(),
		tmpText:         textinput.New(),
		cmdHistory:      []string{},
		cmdHistoryIndex: 0,
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

	case StatusTextMsg:
		m.StatusMessage(msg.Text, msg.Error)

	case BufferSavedMsg:
		if msg.Quit {
			return m, tea.Quit
		}

		if err := m.eb.Reload(); err != nil {
			m.StatusMessage(fmt.Sprintf("Error reloading buffer: %s", err), true)
		} else {
			m.StatusMessage(fmt.Sprintf("Saved %d bytes to %s", msg.BytesWritten, msg.FileName), false)
		}
	}

	return m, nil
}

func (m Model) View() string {
	tStart := time.Now()

	// Hex view
	hexView, err := m.RenderHexView()
	if err != nil {
		hexView = err.Error()
	}

	// Status bar
	statusBar := m.RenderStatus()

	tEnd := time.Now()
	renderTime := float64(tEnd.Sub(tStart).Nanoseconds())
	m.renderTimeAvgNS = (m.renderTimeAvgNS*float64(m.framesRendered) + renderTime) / float64(m.framesRendered+1)
	m.framesRendered++

	return fmt.Sprintf("%s\n%s", hexView, statusBar)
}

func (m *Model) LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", name, err)
	}
	m.eb = core.NewEditorBuffer(name, f)
	return nil
}

// SetCursor sets the cursor position.
func (m *Model) SetCursor(pos int64) {
	m.eb.Cursor = util.Clamp(pos, 0, m.eb.Size()-1)
	m.ScrollToCursor()
}

// MoveCursor moves the cursor by the given amount.
func (m *Model) MoveCursor(amount int64) {
	m.SetCursor(m.eb.Cursor + amount)
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
	m.prevMode = m.mode
	m.mode = mode

	m.statusError = false
	m.cmdText.SetValue("")
	m.tmpText.SetValue("")
	if mode == ModeCommand {
		m.cmdText.Prompt = ":"
		m.cmdText.Focus()
		m.cmdHistoryIndex = len(m.cmdHistory)
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
func (m *Model) StatusMessage(msg string, isError bool) {
	m.cmdText.SetValue(msg)
	m.statusError = isError
}
