package collections

import "testing"

func TestIndexReturnsThePosition(t *testing.T) {
	if got := ints(6).Index(func(v int) bool { return v == 3 }); got != 3 {
		t.Fatalf("want 3, got %d", got)
	}
}

func TestIndexReturnsMinusOneWhenNothingMatches(t *testing.T) {
	if got := ints(6).Index(func(int) bool { return false }); got != -1 {
		t.Fatalf("want -1, got %d", got)
	}
}

func TestIndexShortCircuits(t *testing.T) {
	var pulls int
	counting(1000, &pulls).Index(func(v int) bool { return v == 2 })
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

// Position counts the stream, so a Filter ahead of it renumbers what survives.
func TestIndexCountsTheStreamNotTheSource(t *testing.T) {
	got := ints(10).
		Filter(func(v int) bool { return v%2 == 0 }).
		Index(func(v int) bool { return v == 6 })
	if got != 3 {
		t.Fatalf("6 is the fourth even number, so want 3, got %d", got)
	}
}

func TestIndexValue(t *testing.T) {
	src := Slice[string]{"ada", "bob", "cyd"}
	if got := IndexValue(src.Values(), "cyd"); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
	if got := IndexValue(src.Values(), "nope"); got != -1 {
		t.Fatalf("want -1, got %d", got)
	}
}

func TestIndexOnAnEmptyStream(t *testing.T) {
	if got := ints(0).Index(func(int) bool { return true }); got != -1 {
		t.Fatalf("want -1, got %d", got)
	}
}

func TestIndexAfterDropIsRelativeToWhatRemains(t *testing.T) {
	if got := ints(10).Drop(5).Index(func(v int) bool { return v == 7 }); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
}
