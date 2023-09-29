package core

import "golang.org/x/exp/constraints"

func clamp[T constraints.Integer](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
