package util

import (
	"crypto/rand"
	"fmt"
	"os"
)

// SwapFile swaps the two files. This is done by renaming a to a temporary name,
// renaming b to a, and renaming the temporary name to b.
//
// In the future this could be made atomic on Linux using the renameat2 syscall.
func SwapFile(a, b string) error {
	rnd := make([]byte, 4)
	if _, err := rand.Read(rnd); err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s~%x", a, rnd)

	// Rename a to a temporary name
	if err := os.Rename(a, tmp); err != nil {
		return err
	}

	// Rename b to a
	if err := os.Rename(b, a); err != nil {
		return err
	}

	// Rename the temporary name to b
	if err := os.Rename(tmp, b); err != nil {
		return err
	}

	return nil
}
