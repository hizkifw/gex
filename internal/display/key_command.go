package display

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleKeypressCommand(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	split := strings.Split(m.cmdText.Value(), " ")
	command := split[0]
	args := []string{}
	if len(split) > 1 {
		args = split[1:]
	}

	switch msg.String() {

	// The "esc" key exits command mode
	case "esc":
		m.SetMode(ModeNormal)
		return m, nil

	// The "enter" key executes the command
	case "enter":
		m.SetMode(ModeNormal)

	// Pass the keypress to the command text input
	default:
		m.cmdText.Focus()
		m.cmdText, cmd = m.cmdText.Update(msg)
		return m, cmd
	}

	// Execute the command
	switch command {
	case "q", "quit", "q!":
		return m, tea.Quit

	case "w", "write":
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

			return BufferSavedMsg{FileName: fileName, BytesWritten: n}
		}

	default:
		m.StatusMessage("Unknown command: " + command)
	}

	return m, cmd
}
