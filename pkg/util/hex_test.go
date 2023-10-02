package util_test

import (
	"testing"

	"github.com/hizkifw/gex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestHexStringToBytes(t *testing.T) {
	assert := assert.New(t)

	var matrix = []struct {
		inp       string
		expBytes  []byte
		expParsed []bool
	}{
		{
			inp:       "0123456789abcdef",
			expBytes:  []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			expParsed: []bool{true, true, true, true, true, true, true, true},
		},
		{
			inp:       "12a",
			expBytes:  []byte{0x12, 0xa0},
			expParsed: []bool{true, false},
		},
		{
			inp:       "1za0",
			expBytes:  []byte{0x0, 0xa0},
			expParsed: []bool{false, true},
		},
	}

	for _, test := range matrix {
		b, p := util.HexStringToBytes(test.inp)
		assert.Equal(test.expBytes, b)
		assert.Equal(test.expParsed, p)
	}
}
