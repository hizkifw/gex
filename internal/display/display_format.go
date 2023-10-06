package display

import (
	"encoding/binary"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hizkifw/gex/pkg/core"
	"github.com/hizkifw/gex/pkg/util"
)

// RenderHexView renders the hex dump.
func (m Model) RenderHexView() (string, error) {
	// Get the list of regions
	regions := m.eb.GetRegions()

	r := m.eb.ReadSeeker()
	offset := int64(m.viewRow * m.ncols)

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}

	// Read view from the underlying buffer, plus some extra bytes to make sure
	// the inspector can read ahead
	buf := make([]byte, (m.nrows*m.ncols)+8)
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

			// Check for any active regions at this position
			activeRegions := core.GetActiveRegions(regions, pos)

			// Highlight selection
			isEditing := m.mode == ModeInsert || m.mode == ModeReplace
			styleHex := MakeStyle(m.activeColumn == ActiveColumnHex, isEditing, activeRegions)
			styleAscii := MakeStyle(m.activeColumn == ActiveColumnAscii, isEditing, activeRegions)

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

	// Inspector
	var inspTable string
	if m.inspectorEnabled && m.eb.Cursor >= offset && m.eb.Cursor < offset+int64(m.nrows*m.ncols) {
		var sbInsK strings.Builder
		var sbInsV strings.Builder
		inspOffset := int64(m.eb.Cursor) - offset
		insp := util.Inspect(buf[inspOffset:], m.inspectorByteOrder)
		for i, r := range insp {
			sbInsK.WriteString(r.Key)
			sbInsV.WriteString(r.Val)
			if i < len(insp)-1 {
				sbInsK.WriteString("  \n")
				sbInsV.WriteString("\n")
			}
		}
		byteOrderDisp := " LE "
		if m.inspectorByteOrder == binary.BigEndian {
			byteOrderDisp = " BE "
		}
		inspTable =
			padLeftStyle.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.JoinHorizontal(lipgloss.Top,
						windowTitleStyle.Render("Inspector"),
						statusBarStyle.Render(byteOrderDisp),
					),
					windowStyle.Render(
						lipgloss.JoinHorizontal(lipgloss.Top,
							sbInsK.String(),
							sbInsV.String(),
						),
					),
				),
			)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		sbAddr.String(),
		sbHex.String(),
		sbAscii.String(),
		inspTable,
	), nil
}

// RenderStatus renders the status bar.
func (m Model) RenderStatus() string {
	var sb strings.Builder
	sb.WriteString(statusStyle[m.mode].Render(string(m.mode)))

	// Dirty indicator
	sb.WriteString(statusBarStyle.Render(" "))
	if m.eb.IsDirty() {
		sb.WriteString(statusBarStyle.Render("*"))
	}

	// File name
	fname := path.Base(m.eb.Name)
	if fname == "." {
		fname = "[No Name]"
	}
	sb.WriteString(statusBarStyle.Render(fname))

	sb.WriteString("\n")
	if m.statusError {
		sb.WriteString(textErrorStyle.Render(m.cmdText.View()))
	} else {
		sb.WriteString(m.cmdText.View())
	}

	return sb.String()
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
