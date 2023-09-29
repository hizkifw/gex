package core

import (
	"io"
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

	// The undo stack. When changes are made to the buffer, they are pushed
	// here. This stack serves as the source of truth for the buffer's contents.
	UndoStack []Change

	// The redo stack. When undoing changes, they are popped from the undo stack
	// and pushed here.
	RedoStack []Change

	// Preview is a Change that is currently being edited. Once the user commits
	// the change, it will be pushed onto the UndoStack.
	Preview *Change
}

// NewEditorBuffer creates a new EditorBuffer with the given name and buffer.
func NewEditorBuffer(name string, buffer io.ReadSeeker) *EditorBuffer {
	return &EditorBuffer{
		Name:      name,
		Buffer:    buffer,
		UndoStack: make([]Change, 0),
		RedoStack: make([]Change, 0),
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
