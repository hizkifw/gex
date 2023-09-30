package core

import (
	"io"

	"github.com/hizkifw/gex/pkg/util"
)

// Change represents a change that has been made to a buffer. Changes are
// described by a position, a length, and the data that was inserted at that
// position. Length bytes starting at the Position were removed, and the Data
// was inserted at the Position.
type Change struct {
	// The position at which the change was made.
	Position int64

	// The number of bytes that were removed at the position.
	Removed int64

	// The data inserted at the position. If longer than Length, the extra bytes
	// will replace the bytes at Position + Length.
	Data []byte
}

// ReadSeeker returns a ReadSeeker with the Change applied to the given
// ReadSeeker.
func (c *Change) ReadSeeker(r io.ReadSeeker) io.ReadSeeker {
	p, _ := r.Seek(0, io.SeekCurrent)
	return &changeReadSeeker{r, c, p}
}

// changeReadSeeker is a ReadSeeker that applies a Change to the underlying
// ReadSeeker.
type changeReadSeeker struct {
	r io.ReadSeeker
	c *Change
	p int64
}

// Read implements io.ReadSeeker
var _ io.ReadSeeker = &changeReadSeeker{}

func (r *changeReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.p = offset
	case io.SeekCurrent:
		r.p += offset
	case io.SeekEnd:
		bufEnd, err := r.r.Seek(0, io.SeekEnd)
		if err != nil {
			return r.p, err
		}
		viewEnd := bufEnd - r.c.Removed + int64(len(r.c.Data))
		r.p = viewEnd + offset
	}

	if r.p < 0 {
		return r.p, io.EOF
	}

	return r.p, nil
}

func (r *changeReadSeeker) Read(out []byte) (int, error) {
	dataLength := int64(len(r.c.Data))
	startPos := r.p

	if r.p < r.c.Position {
		// Read data before the change from the underlying buffer
		if _, err := r.r.Seek(r.p, io.SeekStart); err != nil {
			return int(r.p - startPos), err
		}
		maxRead := util.Clamp(int(r.c.Position-r.p), 0, len(out))
		n, err := r.r.Read(out[:maxRead])
		r.p += int64(n)
		if err != nil {
			return int(r.p - startPos), err
		}
	}
	if r.p >= r.c.Position && r.p < r.c.Position+dataLength {
		// Read within the new data
		n := copy(out[r.p-startPos:], r.c.Data[r.p-r.c.Position:])
		r.p += int64(n)
	}
	if r.p >= r.c.Position+dataLength {
		// Read underlying buffer after the new data, skipping the removed length
		if _, err := r.r.Seek(r.p+r.c.Removed-dataLength, io.SeekStart); err != nil {
			return int(r.p - startPos), err
		}
		n, err := r.r.Read(out[r.p-startPos:])
		r.p += int64(n)
		if err != nil {
			return int(r.p - startPos), err
		}
	}

	return int(r.p - startPos), nil
}
