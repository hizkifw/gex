package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderView renders a hex dump and ASCII view of the given ReadSeeker.
func RenderView(r io.ReadSeeker, ncols, nrows, startRow int, selectionStart, selectionEnd int64) (string, error) {
	offset := int64(startRow * ncols)
	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}

	// Read view from the underlying buffer
	buf := make([]byte, ncols*nrows)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	var sbAddr strings.Builder
	var sbHex strings.Builder
	var sbAscii strings.Builder
	for row := 0; row < nrows; row++ {
		// Address column
		sbAddr.WriteString(addrStyle.Render(fmt.Sprintf("%08x", ncols*(row+startRow))))

		for col := 0; col < ncols; col++ {
			i := row*ncols + col
			pos := int64(i) + offset
			if col%8 == 0 {
				sbHex.WriteString(" ")
			}

			// Highlight selection
			var styleHex *lipgloss.Style
			var styleAscii *lipgloss.Style
			if pos >= selectionStart && pos <= selectionEnd {
				styleHex = &hexSelectedStyle
				styleAscii = &asciiSelectedStyle
			} else {
				styleHex = &hexNormalStyle
				styleAscii = &asciiNormalStyle
			}

			// Hex column
			if i >= n {
				sbHex.WriteString(styleHex.Render("  "))
			} else {
				sbHex.WriteString(styleHex.Render(fmt.Sprintf("%02x", buf[i])))
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

		if row < nrows-1 {
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
