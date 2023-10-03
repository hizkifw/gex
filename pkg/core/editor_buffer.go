package core

import (
	"io"

	"golang.org/x/exp/slices"
)

// EditorBuffer represents a file or buffer that is open in the editor.
type EditorBuffer struct {
	// The name of the buffer. If the buffer is a file, this is the full path to
	// the file. If the buffer is not a file, this can be any string.
	Name string

	// The current cursor position.
	Cursor int64

	// The current selection start position. If there is no selection, this is
	// equal to Cursor.
	SelectionStart int64

	// The underlying buffer containing the actual data.
	Buffer io.ReadSeeker

	// Clipboard holds the current clipboard contents.
	Clipboard []byte

	// The undo stack. When changes are made to the buffer, they are pushed
	// here. This stack serves as the source of truth for the buffer's contents.
	UndoStack []Change

	// The redo stack. When undoing changes, they are popped from the undo stack
	// and pushed here.
	RedoStack []Change

	// Preview is a Change that is currently being edited. Once the user commits
	// the change, it will be pushed onto the UndoStack.
	Preview *Change

	// Regions is a list of user-defined regions in the buffer. This does not
	// include the selection and other internal regions.
	Regions []Region
}

// NewEditorBuffer creates a new EditorBuffer with the given name and buffer.
func NewEditorBuffer(name string, buffer io.ReadSeeker) *EditorBuffer {
	return &EditorBuffer{
		Name:      name,
		Buffer:    buffer,
		Clipboard: make([]byte, 0),
		UndoStack: make([]Change, 0),
		RedoStack: make([]Change, 0),
		Preview:   nil,
	}
}

// ReadSeeker returns a ReadSeeker with all changes applied to the underlying
// buffer.
func (b *EditorBuffer) ReadSeeker() io.ReadSeeker {
	r := io.ReadSeeker(b.Buffer)

	// Apply all changes in the undo stack
	for i := 0; i < len(b.UndoStack); i++ {
		r = b.UndoStack[i].ReadSeeker(r)
	}

	// Add the preview change
	if b.Preview != nil {
		r = b.Preview.ReadSeeker(r)
	}

	return r
}

// Undo undoes the last change.
func (b *EditorBuffer) Undo() bool {
	if len(b.UndoStack) == 0 {
		return false
	}

	// Move the last change from the undo stack to the redo stack
	chg := b.UndoStack[len(b.UndoStack)-1]
	b.UndoStack = b.UndoStack[:len(b.UndoStack)-1]
	b.RedoStack = append(b.RedoStack, chg)

	return true
}

// Redo redoes the last change.
func (b *EditorBuffer) Redo() bool {
	if len(b.RedoStack) == 0 {
		return false
	}

	// Move the last change from the redo stack to the undo stack
	chg := b.RedoStack[len(b.RedoStack)-1]
	b.RedoStack = b.RedoStack[:len(b.RedoStack)-1]
	b.UndoStack = append(b.UndoStack, chg)

	return true
}

// PreviewChange applies the given change to the preview buffer.
func (b *EditorBuffer) PreviewChange(chg *Change) {
	b.Preview = chg
}

// CommitChange commits the preview change to the buffer.
func (b *EditorBuffer) CommitChange() {
	if b.Preview == nil {
		return
	}

	b.UndoStack = append(b.UndoStack, *b.Preview)
	b.Preview = nil

	// Clear the redo stack
	b.RedoStack = make([]Change, 0)
}

// GetSelectionRange returns the start and end of the current selection.
func (b *EditorBuffer) GetSelectionRange() (int64, int64) {
	if b.Cursor < b.SelectionStart {
		return b.Cursor, b.SelectionStart
	}
	return b.SelectionStart, b.Cursor
}

// Size returns the size of the buffer.
func (b *EditorBuffer) Size() int64 {
	rs := b.ReadSeeker()
	size, _ := rs.Seek(0, io.SeekEnd)
	return size
}

// CopySelection copies the current selection to the clipboard.
func (b *EditorBuffer) CopySelection() (int, error) {
	start, end := b.GetSelectionRange()
	// Add 1 to the end because the range is inclusive
	b.Clipboard = make([]byte, end-start+1)
	rs := b.ReadSeeker()
	rs.Seek(start, io.SeekStart)
	return rs.Read(b.Clipboard)
}

// GetRegions returns a combined list of user-defined regions and internal
// regions.
func (b *EditorBuffer) GetRegions() []Region {
	// Get the list of dirty bytes
	chgs := b.UndoStack
	if b.Preview != nil {
		chgs = append(chgs, *b.Preview)
	}

	// Sort the changes by position
	slices.SortFunc(chgs, func(i, j Change) int {
		return int(i.Position - j.Position)
	})

	// List of consolidated dirty region ranges
	dirty := make([]Range, 0)
	for n, chg := range chgs {
		chStart := chg.Position
		chEnd := chg.Position + int64(len(chg.Data)) - 1

		if n == 0 {
			// First change
			dirty = append(dirty, Range{Start: chStart, End: chEnd})
			continue
		}

		// Check if the change overlaps with the previous change
		prev := dirty[len(dirty)-1]
		if chStart <= prev.End {
			// The change overlaps with the previous change, so merge them
			prev.End = chEnd
			dirty[len(dirty)-1] = prev
			continue
		}

		// The change does not overlap with the previous change, so add it to the
		// list of dirty regions
		dirty = append(dirty, Range{Start: chStart, End: chEnd})
	}

	// Allocate the list of regions, with enough capacity for the dirty regions,
	// user-defined regions, and selection and cursor regions.
	regions := make([]Region, 0, len(dirty)+len(b.Regions)+2)

	// Add the dirty regions
	for i := range dirty {
		regions = append(regions, Region{
			Type:  RegionTypeDirty,
			Range: dirty[i],
		})
	}

	// Add user-defined regions
	regions = append(regions, b.Regions...)

	// Add the selection region
	start, end := b.GetSelectionRange()
	regions = append(regions,
		Region{
			Type: RegionTypeSelection,
			Range: Range{
				Start: start,
				End:   end,
			},
		},
		Region{
			Type: RegionTypeCursor,
			Range: Range{
				Start: b.Cursor,
				End:   b.Cursor,
			},
		},
	)

	// Sort the regions by position
	slices.SortFunc(regions, func(i, j Region) int {
		return int(i.Start - j.Start)
	})

	return regions
}
