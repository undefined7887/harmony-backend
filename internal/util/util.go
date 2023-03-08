package util

import (
	"golang.org/x/exp/constraints"
	"sort"
)

func Sort[T constraints.Ordered](values []T) {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
}

func SortSequence[T constraints.Ordered](values ...T) {
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
}
