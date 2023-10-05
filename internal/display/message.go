package display

import tea "github.com/charmbracelet/bubbletea"

type StatusTextMsg struct {
	Text string
}

type BufferSavedMsg struct {
	FileName     string
	BytesWritten int64
	Quit         bool
}

func TeaMsgCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
