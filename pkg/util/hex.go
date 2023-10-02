package util

import (
	"encoding/hex"
)

// HexStringToBytes converts a hex string to a byte slice. The second return
// value is a slice of booleans indicating whether each byte was successfully
// parsed. If a byte could not be parsed, it will be set to 0.
func HexStringToBytes(s string) ([]byte, []bool) {
	// Pad end of string with 0 if necessary
	incompleteEnd := false
	if len(s)%2 != 0 {
		s += "0"
		incompleteEnd = true
	}

	sz := len(s) / 2
	b := make([]byte, sz)
	parsed := make([]bool, sz)

	for i := 0; i < sz; i++ {
		parsed[i] = true
		n, err := hex.Decode(b[i:i+1], []byte(s[i*2:i*2+2]))
		if err != nil || n != 1 {
			parsed[i] = false
		}
	}

	if incompleteEnd {
		parsed[sz-1] = false
	}

	return b, parsed
}
