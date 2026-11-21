package store

import (
	"slices"
	"strings"
	"testing"
)

func TestDerange(t *testing.T) {
	for n := 2; n <= 8; n++ {
		ids := make([]int64, n)
		for i := range ids {
			ids[i] = int64(i + 1)
		}
		for range 500 {
			out := derange(ids)
			if hasFixedPoint(ids, out) {
				t.Fatalf("n=%d: %v -> %v has someone drawing themselves", n, ids, out)
			}
			sorted := slices.Sorted(slices.Values(out))
			if !slices.Equal(sorted, ids) {
				t.Fatalf("n=%d: %v is not a permutation of %v", n, out, ids)
			}
		}
	}
}

func TestDerangeDoesNotModifyInput(t *testing.T) {
	ids := []int64{10, 20, 30}
	derange(ids)
	if !slices.Equal(ids, []int64{10, 20, 30}) {
		t.Fatalf("input changed to %v", ids)
	}
}

func TestNewPublicID(t *testing.T) {
	seen := map[string]bool{}
	for range 1000 {
		id := newPublicID()
		if len(id) != publicIDLen || strings.Trim(id, publicIDAlphabet) != "" {
			t.Fatalf("public ID %q is not %d alphanumeric characters", id, publicIDLen)
		}
		if seen[id] {
			t.Fatalf("duplicate public ID %q", id)
		}
		seen[id] = true
	}
}
