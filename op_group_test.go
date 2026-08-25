package collections

import (
	"errors"
	"testing"
)

func TestGroupBy(t *testing.T) {
	got := GroupBy(ints(7), func(v int) int { return v % 3 })
	if got.Len() != 3 {
		t.Fatalf("want 3 groups, got %v", got)
	}
	assertInts(t, got[0].Native(), 0, 3, 6)
	assertInts(t, got[1].Native(), 1, 4)
}

func TestGroupByKeepsWithinGroupOrder(t *testing.T) {
	src := Slice[int]{5, 1, 3, 2}
	got := GroupBy(src.Values(), func(v int) int { return v % 2 })
	assertInts(t, got[1].Native(), 5, 1, 3)
}

func TestGroupByOnAnEmptyStream(t *testing.T) {
	if got := GroupBy(ints(0), func(v int) int { return v }); !got.IsEmpty() {
		t.Fatalf("want empty, got %v", got)
	}
}

func TestGroupByX(t *testing.T) {
	if _, err := GroupByX(abortAt(3), func(v int) int { return v }); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	got, err := GroupByX(ints(4).X(), func(v int) int { return v % 2 })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Len() != 2 {
		t.Fatalf("got %v", got)
	}
}

// The ordered form exists for exactly this: the groups come out in the order their
// keys were first seen, which GroupBy cannot promise.
func TestGroupByOrderedKeepsFirstSeenKeyOrder(t *testing.T) {
	src := Slice[string]{"pear", "fig", "plum", "kiwi", "fig"}
	got := GroupByOrdered(src.Values(), func(s string) byte { return s[0] })
	var keys []byte
	for k := range got.All() {
		keys = append(keys, k)
	}
	if len(keys) != 3 || keys[0] != 'p' || keys[1] != 'f' || keys[2] != 'k' {
		t.Fatalf("got %q", keys)
	}
	group, _ := got.Get('f')
	if len(group) != 2 {
		t.Fatalf("want 2 figs, got %v", group)
	}
}

func TestGroupByOrderedX(t *testing.T) {
	if _, err := GroupByOrderedX(abortAt(3), func(v int) int { return v }); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPartitionSplitsInTwo(t *testing.T) {
	accepted, rejected := ints(6).Partition(func(v int) bool { return v%2 == 0 })
	assertInts(t, accepted.Native(), 0, 2, 4)
	assertInts(t, rejected.Native(), 1, 3, 5)
}

func TestPartitionWhenEverythingGoesOneWay(t *testing.T) {
	accepted, rejected := ints(3).Partition(func(int) bool { return true })
	assertInts(t, accepted.Native(), 0, 1, 2)
	if len(rejected) != 0 {
		t.Fatalf("want nothing rejected, got %v", rejected)
	}
}

func TestPartitionXDiscardsBothHalvesOnError(t *testing.T) {
	accepted, rejected, err := ints(6).PartitionX(func(v int) (bool, error) {
		if v == 3 {
			return false, errBoom
		}
		return true, nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if accepted != nil || rejected != nil {
		t.Fatalf("both halves must be discarded, got %v and %v", accepted, rejected)
	}
}

func TestPartitionOnAnAbortingStream(t *testing.T) {
	_, _, err := abortAt(2).Partition(func(int) bool { return true })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

// Partition and GroupBy answer the same question two ways, so they must agree.
func TestPartitionAgreesWithGroupBy(t *testing.T) {
	accepted, rejected := ints(6).Partition(func(v int) bool { return v%2 == 0 })
	grouped := GroupBy(ints(6), func(v int) bool { return v%2 == 0 })
	assertInts(t, accepted.Native(), grouped[true].Native()...)
	assertInts(t, rejected.Native(), grouped[false].Native()...)
}
