package collections

import (
	"errors"
	"testing"
)

func TestFindReturnsTheFirstMatch(t *testing.T) {
	got := ints(6).Find(func(v int) bool { return v > 2 })
	if got == nil || *got != 3 {
		t.Fatalf("got %v", got)
	}
}

func TestFindReturnsNilWhenNothingMatches(t *testing.T) {
	if got := ints(6).Find(func(int) bool { return false }); got != nil {
		t.Fatalf("want nil, got %v", *got)
	}
}

func TestFindShortCircuits(t *testing.T) {
	var pulls int
	counting(1000, &pulls).Find(func(v int) bool { return v == 2 })
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

// v0.1 returned &s[i], a pointer into the caller's slice. An iterator has no
// addressable source, so this pointer is to a copy and writing through it must not
// reach the source.
func TestFindReturnsAPointerToACopy(t *testing.T) {
	src := Slice[int]{1, 2, 3}
	got := src.Values().Find(func(v int) bool { return v == 2 })
	*got = 99
	assertInts(t, src.Native(), 1, 2, 3)
}

func TestFindLastReturnsTheLastMatch(t *testing.T) {
	got := ints(6).FindLast(func(v int) bool { return v%2 == 0 })
	if got == nil || *got != 4 {
		t.Fatalf("got %v", got)
	}
}

func TestFindLastDrainsTheStream(t *testing.T) {
	var pulls int
	counting(5, &pulls).FindLast(func(v int) bool { return v == 0 })
	if pulls != 5 {
		t.Fatalf("FindLast has to reach the end, got %d pulls", pulls)
	}
}

func TestFirstAndLast(t *testing.T) {
	first, last := ints(4).First(), ints(4).Last()
	if first == nil || *first != 0 {
		t.Fatalf("first: %v", first)
	}
	if last == nil || *last != 3 {
		t.Fatalf("last: %v", last)
	}
}

func TestFirstAndLastOnAnEmptyStreamAreNil(t *testing.T) {
	if ints(0).First() != nil || ints(0).Last() != nil {
		t.Fatal("want nil on both")
	}
}

func TestFirstShortCircuits(t *testing.T) {
	var pulls int
	counting(1000, &pulls).First()
	if pulls != 1 {
		t.Fatalf("want 1 pull, got %d", pulls)
	}
}

func TestFindXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(6).FindX(func(v int) (bool, error) {
		if v == 2 {
			return false, errBoom
		}
		return false, nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

// A match found before the abort is a real answer, so it comes back without error.
func TestFindBeforeAnAbortSeesNoError(t *testing.T) {
	got, err := abortAt(3).Find(func(v int) bool { return v == 1 })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || *got != 1 {
		t.Fatalf("got %v", got)
	}
}

// FindLast cannot stop early, so the same stream does surface the abort.
func TestFindLastReachesTheAbort(t *testing.T) {
	_, err := abortAt(3).FindLast(func(v int) bool { return v == 1 })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairFindReturnsAnEntry(t *testing.T) {
	got := pairs(5).Find(func(_, v int) bool { return v == 30 })
	if got == nil || got.Key != 3 || got.Value != 30 {
		t.Fatalf("got %v", got)
	}
}

func TestPairFirstAndLast(t *testing.T) {
	first := pairs(3).First()
	last := pairs(3).Last()
	if first == nil || first.Key != 0 || last == nil || last.Key != 2 {
		t.Fatalf("got %v and %v", first, last)
	}
}

func TestFindOnACollection(t *testing.T) {
	s := Slice[string]{"ada", "bob", "cyd"}
	got := s.Values().Find(func(v string) bool { return v > "b" })
	if got == nil || *got != "bob" {
		t.Fatalf("got %v", got)
	}
}
