package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHexView renders the hex dump.
func (m Model) RenderHexView() (string, error) {
	// Calculate the selection start and end positions
	selectionStart, selectionEnd := m.eb.GetSelectionRange()

	r := m.eb.ReadSeeker()
	offset := int64(m.viewRow * m.ncols)

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}

	// Read view from the underlying buffer
	buf := make([]byte, m.nrows*m.ncols)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	var sbAddr strings.Builder
	var sbHex strings.Builder
	var sbAscii strings.Builder
	for row := 0; row < m.nrows; row++ {
		// Address column
		sbAddr.WriteString(addrStyle.Render(fmt.Sprintf("%08x", m.ncols*(row+m.viewRow))))

		for col := 0; col < m.ncols; col++ {
			i := row*m.ncols + col
			pos := int64(i) + offset
			if col%8 == 0 {
				sbHex.WriteString(" ")
			}

			// Highlight selection
			selected := pos >= selectionStart && pos <= selectionEnd
			cursor := pos == m.eb.Cursor
			styleHex := MakeStyle(m.activeColumn == ActiveColumnHex, selected, cursor)
			styleAscii := MakeStyle(m.activeColumn == ActiveColumnAscii, selected, cursor)

			// Hex column
			if i >= n {
				sbHex.WriteString(styleHex.Render("  "))
			} else {
				sbHex.WriteString(styleHex.Render(fmt.Sprintf("%02x ", buf[i])))
			}

			// ASCII column
			if i >= n {
				sbAscii.WriteString(" ")
			} else if buf[i] >= 32 && buf[i] <= 126 {
				sbAscii.WriteString(styleAscii.Render(string(buf[i])))
			} else {
				sbAscii.WriteString(styleAscii.Render("."))
			}
		}

		if row < m.nrows-1 {
			sbAddr.WriteString("\n")
			sbHex.WriteString(" \n")
			sbAscii.WriteString("\n")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, sbAddr.String(), sbHex.String(), sbAscii.String()), nil
}

// CalculateViewSize calculates the number of rows and columns that can fit in
// the given width and height.
func CalculateViewSize(width, height int) (ncols, nrows int) {
	// 8 chars for the address + 2 padding
	// 3 chars for each hex value + 1 padding every 8 chars
	// 1 char for each ASCII value
	ncols = (width - 8 - 2) / (3 + 1)

	// Round down to nearest multiple of 8
	ncols = ncols - ncols%8

	// Allocate 2 rows for bottom status bar
	nrows = height - 2
	return
}
