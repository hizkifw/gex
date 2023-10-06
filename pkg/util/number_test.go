package util_test

import (
	"testing"

	"github.com/hizkifw/gex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	t.Run("Test Min with integers", func(t *testing.T) {
		result := util.Min(5, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("Test Min with negative integers", func(t *testing.T) {
		result := util.Min(-10, -5)
		assert.Equal(t, -10, result)
	})
}

func TestMax(t *testing.T) {
	t.Run("Test Max with integers", func(t *testing.T) {
		result := util.Max(5, 10)
		assert.Equal(t, 10, result)
	})

	t.Run("Test Max with negative integers", func(t *testing.T) {
		result := util.Max(-10, -5)
		assert.Equal(t, -5, result)
	})
}

func TestClamp(t *testing.T) {
	t.Run("Test Clamp with integers within range", func(t *testing.T) {
		result := util.Clamp(5, 0, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("Test Clamp with integers below range", func(t *testing.T) {
		result := util.Clamp(-5, 0, 10)
		assert.Equal(t, 0, result)
	})

	t.Run("Test Clamp with integers above range", func(t *testing.T) {
		result := util.Clamp(15, 0, 10)
		assert.Equal(t, 10, result)
	})
}
