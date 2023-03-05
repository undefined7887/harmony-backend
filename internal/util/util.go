package util

import (
	"golang.org/x/exp/constraints"
	"sort"
)

func SortSequence[T constraints.Ordered](values ...T) []T {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	return values
}
