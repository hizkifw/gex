package core

import (
	"fmt"
	"io"
	"os"

	"github.com/hizkifw/gex/pkg/util"
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

// Reload reloads the buffer from the file that is backing it.
func (b *EditorBuffer) Reload() error {
	// Close the existing buffer if it is a file
	if f, ok := b.Buffer.(*os.File); ok {
		f.Close()
	}

	f, err := os.Open(b.Name)
	if err != nil {
		return err
	}

	b.Buffer = f
	b.UndoStack = make([]Change, 0)
	b.RedoStack = make([]Change, 0)
	b.Preview = nil

	return nil
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

	chg := *b.Preview
	b.Preview = nil
	if chg.Removed == 0 && len(chg.Data) == 0 {
		return
	}

	b.UndoStack = append(b.UndoStack, chg)
	b.RedoStack = make([]Change, 0)
}

// IsDirty returns true if the buffer contains unsaved changes.
func (b *EditorBuffer) IsDirty() bool {
	return len(b.UndoStack) > 0 || b.Preview != nil
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

	chgRanges := make([]Range, len(chgs))
	for i := range chgs {
		offset := int64(0)
		for j := i + 1; j < len(chgs); j++ {
			if chgs[j].Position < chgs[i].Position {
				offset += int64(len(chgs[j].Data)) - chgs[j].Removed
			}
		}

		chgRanges[i] = Range{
			Start: chgs[i].Position + offset,
			End:   chgs[i].Position + offset + int64(len(chgs[i].Data)) - 1,
		}
	}

	// Sort the change ranges by position
	slices.SortFunc(chgRanges, func(i, j Range) int {
		return int(i.Start - j.Start)
	})

	// List of consolidated dirty region ranges
	dirty := make([]Range, 0)
	for n, rng := range chgRanges {
		if n == 0 {
			// First change
			dirty = append(dirty, Range{Start: rng.Start, End: rng.End})
			continue
		}

		// Check if the change overlaps with the previous change
		prev := dirty[len(dirty)-1]
		if rng.Start <= prev.End {
			// The change overlaps with the previous change, so merge them
			prev.End = rng.End
			dirty[len(dirty)-1] = prev
			continue
		}

		// The change does not overlap with the previous change, so add it to the
		// list of dirty regions
		dirty = append(dirty, Range{Start: rng.Start, End: rng.End})
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

// WriteToFile writes the buffer contents to the given file. Do not call this
// with the same file that is backing the buffer. To safely save the buffer to
// the same file, use Save.
func (b *EditorBuffer) WriteToFile(filename string) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	rs := b.ReadSeeker()
	_, err = rs.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return io.Copy(f, rs)
}

// SaveInPlace will modify the edited file in-place. This will only work if the
// changes do not change the size of the file.
func (b *EditorBuffer) SaveInPlace() error {
	// Check if the changes modify the file size
	for _, chg := range b.UndoStack {
		if chg.Removed != int64(len(chg.Data)) {
			return fmt.Errorf("change at position %d modifies file size", chg.Position)
		}
	}

	// Open the file for writing
	f, err := os.OpenFile(b.Name, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer f.Close()

	// Apply the changes to the file
	for _, chg := range b.UndoStack {
		_, err = f.Seek(chg.Position, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek to position %d: %w", chg.Position, err)
		}

		_, err = f.Write(chg.Data)
		if err != nil {
			return fmt.Errorf("failed to write data: %w", err)
		}
	}

	return nil
}

// Save saves the buffer to the file that is backing it. It will also create a
// backup file. If fileName is empty, the buffer's name will be used.
func (b *EditorBuffer) Save(fileName string) (int64, error) {
	if fileName == "" {
		fileName = b.Name
	}

	// Create the backup file
	backupFilename := fileName + "~"

	// Write the buffer to the backup file
	n, err := b.WriteToFile(backupFilename)
	if err != nil {
		return n, err
	}

	// Swap the backup file with the original file
	err = util.SwapFile(fileName, backupFilename)

	return n, err
}
