package collections

import (
	"errors"
	"testing"
)

func TestCount(t *testing.T) {
	if got := ints(7).Count(); got != 7 {
		t.Fatalf("want 7, got %d", got)
	}
	if got := ints(0).Count(); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}
}

func TestCountBy(t *testing.T) {
	if got := ints(10).CountBy(func(v int) bool { return v%3 == 0 }); got != 4 {
		t.Fatalf("want 4, got %d", got)
	}
}

func TestCountByXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(6).CountByX(func(v int) (bool, error) {
		if v == 3 {
			return false, errBoom
		}
		return true, nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestIsEmpty(t *testing.T) {
	if !ints(0).IsEmpty() {
		t.Fatal("an empty stream is empty")
	}
	if ints(1).IsEmpty() {
		t.Fatal("a stream of one is not empty")
	}
}

// IsEmpty must not drain what it is asked about.
func TestIsEmptyPullsOneElementAtMost(t *testing.T) {
	var pulls int
	counting(1000, &pulls).IsEmpty()
	if pulls != 1 {
		t.Fatalf("want 1 pull, got %d", pulls)
	}
}

func TestIsEmptyAfterAFilterThatRemovesEverything(t *testing.T) {
	if !ints(10).Filter(func(int) bool { return false }).IsEmpty() {
		t.Fatal("want empty")
	}
}

func TestPairCountAndIsEmpty(t *testing.T) {
	if got := pairs(4).Count(); got != 4 {
		t.Fatalf("want 4, got %d", got)
	}
	if got := pairs(4).CountBy(func(k, _ int) bool { return k > 1 }); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
	if pairs(1).IsEmpty() {
		t.Fatal("not empty")
	}
}

// Count drains; a collection's own length does not. The two must agree.
func TestCountAgreesWithTheSourceLength(t *testing.T) {
	src := Slice[int]{4, 5, 6}
	if got := src.Values().Count(); got != len(src) {
		t.Fatalf("want %d, got %d", len(src), got)
	}
}
