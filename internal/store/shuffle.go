package store

import (
	"math/rand/v2"
	"slices"
)

// derange returns a random reordering of ids where no element stays at its own
// index, so out[i] is who ids[i] buys a gift for and nobody draws themselves.
// ids must contain at least two distinct elements.
//
// It reshuffles until it hits a valid ordering. Roughly 1 in e random
// permutations qualifies, so this takes about three tries on average, and
// every valid assignment is equally likely.
func derange(ids []int64) []int64 {
	out := slices.Clone(ids)
	for {
		rand.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
		if !hasFixedPoint(ids, out) {
			return out
		}
	}
}

func hasFixedPoint(a, b []int64) bool {
	for i := range a {
		if a[i] == b[i] {
			return true
		}
	}
	return false
}
