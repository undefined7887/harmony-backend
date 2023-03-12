package util

import (
	"sort"

	"golang.org/x/exp/constraints"

	"github.com/samber/lo"
)

func Map[T, R any](collection []T, cb func(item T) R) []R {
	return lo.Map(collection, func(item T, index int) R {
		return cb(item)
	})
}

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
