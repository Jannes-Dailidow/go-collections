package collections

import (
	"errors"
	"testing"
)

func TestContains(t *testing.T) {
	if !ints(4).Contains(func(v int) bool { return v == 2 }) {
		t.Fatal("want true")
	}
	if ints(4).Contains(func(v int) bool { return v == 9 }) {
		t.Fatal("want false")
	}
}

func TestContainsOnAnEmptyStreamIsFalse(t *testing.T) {
	if ints(0).Contains(func(int) bool { return true }) {
		t.Fatal("want false")
	}
}

func TestContainsShortCircuits(t *testing.T) {
	var pulls int
	counting(1000, &pulls).Contains(func(v int) bool { return v == 1 })
	if pulls != 2 {
		t.Fatalf("want 2 pulls, got %d", pulls)
	}
}

func TestContainsValue(t *testing.T) {
	src := Slice[string]{"ada", "bob"}
	if !ContainsValue(src.Values(), "bob") {
		t.Fatal("want true")
	}
	if ContainsValue(src.Values(), "cyd") {
		t.Fatal("want false")
	}
}

func TestEvery(t *testing.T) {
	if !ints(4).Every(func(v int) bool { return v < 4 }) {
		t.Fatal("want true")
	}
	if ints(4).Every(func(v int) bool { return v < 2 }) {
		t.Fatal("want false")
	}
}

// Vacuous truth, not a special case.
func TestEveryOnAnEmptyStreamIsTrue(t *testing.T) {
	if !ints(0).Every(func(int) bool { return false }) {
		t.Fatal("want true on an empty stream")
	}
}

func TestEveryShortCircuitsOnTheFirstRejection(t *testing.T) {
	var pulls int
	counting(1000, &pulls).Every(func(v int) bool { return v < 2 })
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

func TestEveryXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(6).EveryX(func(v int) (bool, error) {
		if v == 2 {
			return true, errBoom
		}
		return true, nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairContainsAndEvery(t *testing.T) {
	if !pairs(4).Contains(func(_, v int) bool { return v == 20 }) {
		t.Fatal("want true")
	}
	if !pairs(4).Every(func(k, _ int) bool { return k < 4 }) {
		t.Fatal("want true")
	}
	if pairs(4).Every(func(k, _ int) bool { return k < 2 }) {
		t.Fatal("want false")
	}
}

func TestContainsOnASetUsesTheIteratorPath(t *testing.T) {
	s := CollectSet(ints(4))
	if !s.Values().Contains(func(v int) bool { return v == 2 }) {
		t.Fatal("want true")
	}
}
