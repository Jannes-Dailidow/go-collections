package collections

import (
	"errors"
	"testing"
)

func TestFilterKeepsOnlyMatching(t *testing.T) {
	got := ints(6).Filter(func(v int) bool { return v%2 == 0 }).Native()
	assertInts(t, got, 0, 2, 4)
}

func TestFilterXKeepsOnlyMatching(t *testing.T) {
	got, err := ints(6).FilterX(func(v int) (bool, error) {
		return v%2 == 0, nil
	}).Native()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertInts(t, got, 0, 2, 4)
}

func TestFilterXReturnsTheCallbackError(t *testing.T) {
	got, err := ints(6).FilterX(func(v int) (bool, error) {
		if v == 3 {
			return false, errBoom
		}
		return true, nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != nil {
		t.Fatalf("the partial result must be discarded, got %v", got)
	}
}

func TestFilterOnACollection(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4}
	got := s.Values().Filter(func(v int) bool { return v > 2 }).Slice()
	assertInts(t, got.Native(), 3, 4)
}

func TestPairFilterKeepsOnlyMatching(t *testing.T) {
	m := Map[string, int]{"a": 1, "b": 2, "c": 3}
	got := CollectMap(m.All().Filter(func(_ string, v int) bool { return v > 1 }))
	if len(got) != 2 {
		t.Fatalf("want 2 pairs, got %v", got)
	}
	if _, ok := got["a"]; ok {
		t.Fatalf("a should have been filtered out: %v", got)
	}
}

// Filtering is what makes an OrderedSet's first-seen order worth preserving, so
// check the two compose.
func TestFilterPreservesOrderThroughAnOrderedSet(t *testing.T) {
	src := Slice[int]{5, 3, 5, 9, 3, 1}
	got := CollectOrderedSet(src.Values().Filter(func(v int) bool { return v > 2 }))
	assertInts(t, got.Native(), 5, 3, 9)
}
