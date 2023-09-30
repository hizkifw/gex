package util

import "golang.org/x/exp/constraints"

func Min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Clamp[T constraints.Integer](v, min, max T) T {
	return Max(Min(v, max), min)
}
