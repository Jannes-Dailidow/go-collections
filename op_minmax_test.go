package collections

import (
	"errors"
	"testing"
)

func TestMinAndMax(t *testing.T) {
	src := Slice[int]{3, 1, 4, 1, 5}
	min, max := Min(src.Values()), Max(src.Values())
	if min == nil || *min != 1 {
		t.Fatalf("min: %v", min)
	}
	if max == nil || *max != 5 {
		t.Fatalf("max: %v", max)
	}
}

func TestMinAndMaxOnAnEmptyStreamAreNil(t *testing.T) {
	if Min(ints(0)) != nil || Max(ints(0)) != nil {
		t.Fatal("want nil on both")
	}
}

func TestMinAndMaxOnStrings(t *testing.T) {
	src := Slice[string]{"cyd", "ada", "bob"}
	if got := Min(src.Values()); got == nil || *got != "ada" {
		t.Fatalf("got %v", got)
	}
}

func TestMinByAndMaxByUseTheKey(t *testing.T) {
	type user struct {
		name string
		age  int
	}
	src := Slice[user]{{"ada", 36}, {"bob", 24}, {"cyd", 41}}
	youngest := src.Values().MinBy(func(u user) int { return u.age })
	oldest := src.Values().MaxBy(func(u user) int { return u.age })
	if youngest == nil || youngest.name != "bob" {
		t.Fatalf("youngest: %v", youngest)
	}
	if oldest == nil || oldest.name != "cyd" {
		t.Fatalf("oldest: %v", oldest)
	}
}

// A tie goes to the first element, so the result is predictable.
func TestMinByTieGoesToTheFirst(t *testing.T) {
	type tagged struct {
		id  int
		key int
	}
	src := Slice[tagged]{{1, 5}, {2, 5}, {3, 9}}
	got := src.Values().MinBy(func(x tagged) int { return x.key })
	if got == nil || got.id != 1 {
		t.Fatalf("want the first of the tie, got %v", got)
	}
	max := src.Values().MaxBy(func(x tagged) int { return -x.key })
	if max == nil || max.id != 1 {
		t.Fatalf("MaxBy must tie the same way, got %v", max)
	}
}

func TestMinByOnAnAbortingStream(t *testing.T) {
	_, err := abortAt(2).MinBy(func(v int) int { return v })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairMinByAndMaxBy(t *testing.T) {
	smallest := pairs(5).MinBy(func(_, v int) int { return v })
	largest := pairs(5).MaxBy(func(_, v int) int { return v })
	if smallest == nil || smallest.Key != 0 {
		t.Fatalf("smallest: %v", smallest)
	}
	if largest == nil || largest.Key != 4 || largest.Value != 40 {
		t.Fatalf("largest: %v", largest)
	}
}

// MinBy returns a copy, like Find.
func TestMinByReturnsAPointerToACopy(t *testing.T) {
	src := Slice[int]{3, 1, 2}
	got := src.Values().MinBy(func(v int) int { return v })
	*got = 99
	assertInts(t, src.Native(), 3, 1, 2)
}

// The documented route for a key that can fail.
func TestMinByOverAFallibleKeyViaKeyByX(t *testing.T) {
	src := Slice[string]{"3", "1", "nope"}
	_, err := src.Values().KeyByX(func(s string) (int, error) {
		if s == "nope" {
			return 0, errBoom
		}
		return len(s), nil
	}).MinBy(func(k int, _ string) int { return k })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}
