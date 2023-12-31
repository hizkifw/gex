package display

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/pkg/core"
	"github.com/hizkifw/gex/pkg/util"
)

func HandleKeypressInsert(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	key := msg.String()

	if key == "esc" {
		// Exit insert/replace mode
		m.eb.CommitChange()
		m.SetMode(ModeNormal)
		return m, nil
	}

	// Get the start position of the change
	start := m.eb.Cursor
	if m.eb.Preview != nil {
		start = m.eb.Preview.Position
	} else {
		m.eb.PreviewChange(&core.Change{Position: start, Removed: 0, Data: []byte{}})
	}

	// Pass inputs to the temporary text input
	m.tmpText, _ = m.tmpText.Update(msg)
	tmpInput := m.tmpText.Value()
	nBytes := int64(len(tmpInput))
	bufLen := m.eb.Size()
	removed := m.eb.Preview.Removed

	if m.activeColumn == ActiveColumnAscii {
		if m.mode == ModeReplace {
			removed = util.Min(nBytes, bufLen-start)
		}

		// Update the cursor position
		m.SetCursor(start + int64(m.tmpText.Position()))
		m.eb.SelectionStart = m.eb.Cursor

		// Set the preview change
		m.eb.PreviewChange(&core.Change{
			Position: start,
			Removed:  removed,
			Data:     []byte(tmpInput),
		})
	} else if m.activeColumn == ActiveColumnHex {
		// Count number of bytes in the temporary input
		if m.mode == ModeReplace {
			nBytes += nBytes % 2
			nBytes /= 2
			removed = util.Min(nBytes, bufLen-start)
		}

		// Get bytes from hex string
		b, _ := util.HexStringToBytes(tmpInput)

		// If moving left and right, update the text input again to move whole
		// byte instead of hex character.
		if key == "left" || key == "right" {
			m.tmpText, _ = m.tmpText.Update(msg)
		}

		// Update the cursor position
		m.SetCursor(start + int64(m.tmpText.Position()/2))
		m.eb.SelectionStart = m.eb.Cursor

		// Set the preview change
		m.eb.PreviewChange(&core.Change{
			Position: start,
			Removed:  removed,
			Data:     b,
		})
	}

	return m, nil
}
