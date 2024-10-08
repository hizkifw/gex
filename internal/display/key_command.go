package display

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/util"
)

func handleCommand(m Model, command string, args []string) (Model, tea.Cmd) {
	// Execute the command
	switch command {
	case "q", "quit", "q!", "quit!":
		if m.eb.IsDirty() && !strings.HasSuffix(command, "!") {
			return m, TeaMsgCmd(StatusTextMsg{Text: "No write since last change (add ! to override)", Error: true})
		}
		return m, tea.Quit

	case "w", "write", "wq":
		// Save the buffer
		fileName := m.eb.Name
		overwrite := true
		if len(args) > 0 {
			fileName = args[0]
			overwrite = false
		}

		saveCmd := func() tea.Msg {
			var n int64
			var err error
			if overwrite {
				n, err = m.eb.Save(fileName)
			} else {
				n, err = m.eb.WriteToFile(fileName)
			}
			if err != nil {
				return StatusTextMsg{Text: "Error saving: " + err.Error(), Error: true}
			}

			return BufferSavedMsg{FileName: fileName, BytesWritten: n, Quit: command == "wq"}
		}

		return m, tea.Batch(TeaMsgCmd(StatusTextMsg{Text: "Saving " + fileName}), saveCmd)

	case "goto":
		// Go to a specific byte offset
		if len(args) == 0 {
			return m, TeaMsgCmd(StatusTextMsg{Text: "Usage: goto <offset(h)>"})
		}

		offsetHexStr := args[0]
		if len(offsetHexStr)%2 != 0 {
			offsetHexStr = "0" + offsetHexStr
		}
		offsetHex, err := hex.DecodeString(offsetHexStr)
		if err != nil {
			return m, TeaMsgCmd(StatusTextMsg{Text: "Invalid offset", Error: true})
		}
		if len(offsetHex) < 8 {
			offsetHex = append(make([]byte, 8-len(offsetHex)), offsetHex...)
		}

		offset := binary.BigEndian.Uint64(offsetHex)
		if offset > uint64(m.eb.Size()) {
			return m, TeaMsgCmd(StatusTextMsg{Text: fmt.Sprintf("Offset %xh out of range", offset), Error: true})
		}

		m.SetCursor(int64(offset))
		if m.prevMode != ModeVisual {
			m.eb.SelectionStart = m.eb.Cursor
		}

	case "set":
		// Set a option
		if len(args) < 2 {
			return m, TeaMsgCmd(StatusTextMsg{Text: "Usage: set <option> <value>"})
		}

		option := args[0]
		value := args[1]

		switch option {
		case "cols":
			// Set the number of columns
			cols, err := strconv.Atoi(value)
			if err != nil {
				return m, TeaMsgCmd(StatusTextMsg{Text: "Expected a number", Error: true})
			}
			m.ncols = cols
			m.ScrollToCursor()

		case "inspector.enabled":
			// Enable/disable the inspector
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				return m, TeaMsgCmd(StatusTextMsg{Text: "Expected either true or false", Error: true})
			}
			m.inspectorEnabled = enabled

		case "inspector.byteOrder":
			// Set the inspector byte order
			switch value {
			case "big", "be", "b":
				m.inspectorByteOrder = binary.BigEndian
			case "little", "le", "l":
				m.inspectorByteOrder = binary.LittleEndian
			default:
				return m, TeaMsgCmd(StatusTextMsg{Text: "Expected either b or l", Error: true})
			}

		default:
			return m, TeaMsgCmd(StatusTextMsg{Text: "Unknown option: " + option, Error: true})

		}

	default:
		return m, TeaMsgCmd(StatusTextMsg{Text: "Unknown command: " + command, Error: true})
	}

	return m, nil
}

func HandleKeypressCommand(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg.String() {

	// The "esc" key exits command mode
	case "esc":
		m.SetMode(m.prevMode)

	// The "up" and "down" keys cycle through the command history
	case "up", "down":
		if len(m.cmdHistory) == 0 {
			break
		}

		if msg.String() == "up" {
			m.cmdHistoryIndex--
		} else {
			m.cmdHistoryIndex++
		}
		m.cmdHistoryIndex = util.Clamp(m.cmdHistoryIndex, 0, len(m.cmdHistory))
		if m.cmdHistoryIndex == len(m.cmdHistory) {
			m.cmdText.SetValue("")
		} else {
			m.cmdText.SetValue(m.cmdHistory[m.cmdHistoryIndex])
			m.cmdText.SetCursor(len(m.cmdText.Value()))
		}

	// The "enter" key executes the command
	case "enter":
		m.cmdHistory = append(m.cmdHistory, m.cmdText.Value())
		split := strings.Split(m.cmdText.Value(), " ")
		command := split[0]
		args := []string{}
		if len(split) > 1 {
			args = split[1:]
		}
		m, cmd = handleCommand(m, command, args)
		m.SetMode(m.prevMode)

	// Pass the keypress to the command text input
	default:
		m.cmdText.Focus()
		m.cmdText, cmd = m.cmdText.Update(msg)

		// If the command text input is empty, exit command mode
		cmdVal := m.cmdText.Value()
		if cmdVal == "goto g" {
			m, cmd = handleCommand(m, "goto", []string{"00"})
			m.SetMode(m.prevMode)
		} else if cmdVal == "goto G" {
			m, cmd = handleCommand(m, "goto", []string{fmt.Sprintf("%x", m.eb.Size())})
			m.SetMode(m.prevMode)
		}
	}

	return m, cmd
}
