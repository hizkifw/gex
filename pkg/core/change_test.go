package core_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/hizkifw/gex/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestChange_ReadSeeker(t *testing.T) {
	assert := assert.New(t)

	// Generate large buffer
	var buf bytes.Buffer
	for i := 0; i < 1024; i++ {
		buf.WriteString("0123456789abcdef")
	}
	largeBuffer := buf.Bytes()

	var matrix = []struct {
		chg      *core.Change
		inp      []byte
		expected []byte
	}{
		{
			// Case 1 - Replace without changing length
			chg:      &core.Change{Position: 1, Removed: 5, Data: []byte("hello")},
			inp:      []byte("0123456789"),
			expected: []byte("0hello6789"),
		},
		{
			// Case 2 - Replace part of the buffer and extend with new data
			chg:      &core.Change{Position: 10, Removed: 2, Data: []byte("hello")},
			inp:      []byte("0123456789abcdefghij0123456789"),
			expected: []byte("0123456789hellocdefghij0123456789"),
		},
		{
			// Case 3 - Remove part of the buffer, shortening it
			chg:      &core.Change{Position: 3, Removed: 2, Data: []byte{}},
			inp:      []byte("asdfghjkl;"),
			expected: []byte("asdhjkl;"),
		},
		{
			// Case 4 - Insert at the end of the buffer
			chg:      &core.Change{Position: 5, Removed: 0, Data: []byte(", world")},
			inp:      []byte("hello"),
			expected: []byte("hello, world"),
		},
		{
			// Case 5 - Insert at the beginning of the buffer
			chg:      &core.Change{Position: 0, Removed: 0, Data: []byte("hello, ")},
			inp:      []byte("world"),
			expected: []byte("hello, world"),
		},
		{
			// Case 6 - Large change
			chg:      &core.Change{Position: 0, Removed: 10, Data: largeBuffer},
			inp:      []byte("0123456789"),
			expected: largeBuffer,
		},
	}

	for _, m := range matrix {
		inp := bytes.NewReader(m.inp)
		r := m.chg.ReadSeeker(inp)

		// Test reading from different seek positions
		for start := 0; start < len(m.expected); start++ {
			r.Seek(int64(start), io.SeekStart)
			actual, err := io.ReadAll(r)
			if err != io.EOF {
				assert.NoError(err)
			}
			assert.Equal(m.expected[start:], actual)
		}
	}
}

func TestChange_Stacked(t *testing.T) {
	assert := assert.New(t)

	var matrix = []struct {
		changes  []core.Change
		inp      []byte
		expected []byte
	}{
		{
			changes: []core.Change{
				{Position: 0, Removed: 1, Data: []byte("a")},
				{Position: 1, Removed: 1, Data: []byte("bc")},
				{Position: 0, Removed: 2, Data: []byte("ZY")},
			},
			inp:      []byte("0123456789"),
			expected: []byte("ZYc23456789"),
		},
	}

	for _, m := range matrix {
		t.Logf("Test case: %#v", m.changes)
		inp := bytes.NewReader(m.inp)
		r := io.ReadSeeker(inp)
		for i := 0; i < len(m.changes); i++ {
			r = m.changes[i].ReadSeeker(r)
		}

		actual, err := io.ReadAll(r)
		if err != io.EOF {
			assert.NoError(err)
		}
		assert.Equal(m.expected, actual)
	}
}
