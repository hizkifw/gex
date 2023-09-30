package display

import (
	tea "github.com/charmbracelet/bubbletea"
)

func HandleKeypressCommand(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	command := m.cmdText.Value()

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
	case "q":
		return m, tea.Quit
	default:
		m.StatusMessage("Unknown command: " + command)
	}

	return m, cmd
}
