package util_test

import (
	"os"
	"testing"

	"github.com/hizkifw/gex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestSwap(t *testing.T) {
	// Create two temporary files
	tmp1, err := os.CreateTemp("", "gex-test")
	assert.NoError(t, err)
	tmp2, err := os.CreateTemp("", "gex-test")
	assert.NoError(t, err)

	defer func() {
		// Remove the temporary files
		os.Remove(tmp1.Name())
		os.Remove(tmp2.Name())
	}()

	// Write some data to the files
	_, err = tmp1.WriteString("foo")
	assert.NoError(t, err)
	_, err = tmp2.WriteString("bar")
	assert.NoError(t, err)

	// Close the files
	err = tmp1.Close()
	assert.NoError(t, err)
	err = tmp2.Close()
	assert.NoError(t, err)

	// Swap the files
	err = util.SwapFile(tmp1.Name(), tmp2.Name())
	assert.NoError(t, err)

	// Check the contents
	b, err := os.ReadFile(tmp1.Name())
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(b))

	b, err = os.ReadFile(tmp2.Name())
	assert.NoError(t, err)
	assert.Equal(t, "foo", string(b))
}
