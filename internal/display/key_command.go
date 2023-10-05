package display

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleCommand(m Model, command string, args []string) (Model, tea.Cmd) {
	// Execute the command
	switch command {
	case "q", "quit", "q!":
		if m.eb.IsDirty() && command != "q!" {
			m.StatusMessage("No write since last change (add ! to override)")
			return m, nil
		}
		return m, tea.Quit

	case "w", "write", "wq":
		// Save the buffer
		m.StatusMessage("Saving...")
		fileName := m.eb.Name
		if len(args) > 0 {
			fileName = args[0]
		}

		return m, func() tea.Msg {
			n, err := m.eb.Save(fileName)
			if err != nil {
				return StatusTextMsg{Text: "Error saving: " + err.Error()}
			}

			return BufferSavedMsg{FileName: fileName, BytesWritten: n, Quit: command == "wq"}
		}

	case "goto":
		// Go to a specific byte offset
		if len(args) == 0 {
			m.StatusMessage("Usage: goto <offset(h)>")
			return m, nil
		}

		offsetHexStr := args[0]
		if len(offsetHexStr)%2 != 0 {
			offsetHexStr = "0" + offsetHexStr
		}
		offsetHex, err := hex.DecodeString(offsetHexStr)
		if err != nil {
			m.StatusMessage("Invalid offset")
			return m, nil
		}
		if len(offsetHex) < 8 {
			offsetHex = append(make([]byte, 8-len(offsetHex)), offsetHex...)
		}

		offset := binary.BigEndian.Uint64(offsetHex)
		if offset > uint64(m.eb.Size()) {
			m.StatusMessage(fmt.Sprintf("Offset %xh out of range", offset))
			return m, nil
		}

		m.SetCursor(int64(offset))
		if m.prevMode != ModeVisual {
			m.eb.SelectionStart = m.eb.Cursor
		}

	case "set":
		// Set a option
		if len(args) < 2 {
			m.StatusMessage("Usage: set <option> <value>")
			return m, nil
		}

		option := args[0]
		value := args[1]

		switch option {
		case "cols":
			// Set the number of columns
			cols, err := strconv.Atoi(value)
			if err != nil {
				m.StatusMessage("Invalid value")
				return m, nil
			}
			m.ncols = cols

		case "inspector.enabled":
			// Enable/disable the inspector
			enabled, err := strconv.ParseBool(value)
			if err != nil {
				m.StatusMessage("Invalid value")
				return m, nil
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
				m.StatusMessage("Expected either b or l")
			}

		default:
			m.StatusMessage("Unknown option: " + option)

		}

	default:
		m.StatusMessage("Unknown command: " + command)
	}

	return m, nil
}

func HandleKeypressCommand(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	switch msg.String() {

	// The "esc" key exits command mode
	case "esc":
		m.SetMode(m.prevMode)

	// The "enter" key executes the command
	case "enter":
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
		if cmdVal == "" {
			m.SetMode(m.prevMode)
		} else if cmdVal == "goto g" {
			m, cmd = handleCommand(m, "goto", []string{"00"})
			m.SetMode(m.prevMode)
		} else if cmdVal == "goto G" {
			m, cmd = handleCommand(m, "goto", []string{fmt.Sprintf("%x", m.eb.Size())})
			m.SetMode(m.prevMode)
		}
	}

	return m, cmd
}
