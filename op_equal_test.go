package collections

import (
	"errors"
	"strings"
	"testing"
)

func TestEqualByMatchesIdenticalStreams(t *testing.T) {
	id := func(v int) int { return v }
	if !ints(4).EqualBy(ints(4), id) {
		t.Fatal("want true")
	}
}

func TestEqualByRejectsDifferentContent(t *testing.T) {
	id := func(v int) int { return v }
	if ints(4).EqualBy(ints(4).Map(func(v int) int { return v + 1 }), id) {
		t.Fatal("want false")
	}
}

// Both directions of a length mismatch, since a naive implementation only catches one.
func TestEqualByRejectsBothLengthMismatches(t *testing.T) {
	id := func(v int) int { return v }
	if ints(3).EqualBy(ints(5), id) {
		t.Fatal("a longer right side must not compare equal")
	}
	if ints(5).EqualBy(ints(3), id) {
		t.Fatal("a longer left side must not compare equal")
	}
}

func TestEqualByOnTwoEmptyStreamsIsTrue(t *testing.T) {
	if !ints(0).EqualBy(ints(0), func(v int) int { return v }) {
		t.Fatal("want true")
	}
}

// The key is what makes this useful: equality on a projection rather than the whole
// element.
func TestEqualByComparesOnlyTheKey(t *testing.T) {
	left := Slice[string]{"Ada", "Bob"}
	right := Slice[string]{"ada", "BOB"}
	if !left.Values().EqualBy(right.Values(), strings.ToLower) {
		t.Fatal("want true when compared case-insensitively")
	}
	if left.Values().EqualBy(right.Values(), func(s string) string { return s }) {
		t.Fatal("want false when compared exactly")
	}
}

func TestEqualByIsOrderSensitive(t *testing.T) {
	left := Slice[int]{1, 2, 3}
	right := Slice[int]{3, 2, 1}
	if left.Values().EqualBy(right.Values(), func(v int) int { return v }) {
		t.Fatal("want false: same elements, different order")
	}
}

func TestEqualByShortCircuits(t *testing.T) {
	var pulls int
	counting(1000, &pulls).EqualBy(ints(1000).Map(func(v int) int { return v + 1 }),
		func(v int) int { return v })
	if pulls != 1 {
		t.Fatalf("the first pair already differs, so want 1 pull, got %d", pulls)
	}
}

func TestEqualByOnAFallibleStream(t *testing.T) {
	got, err := ints(4).X().EqualBy(ints(4), func(v int) int { return v })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got {
		t.Fatal("want true")
	}

	if _, err := abortAt(2).EqualBy(ints(4), func(v int) int { return v }); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairEqualBy(t *testing.T) {
	key := func(k, v int) int { return k*100 + v }
	if !pairs(3).EqualBy(pairs(3), key) {
		t.Fatal("want true")
	}
	if pairs(3).EqualBy(pairs(4), key) {
		t.Fatal("want false")
	}
}

// An OrderedSet has an order, so this is meaningful on it; a Set does not.
func TestEqualByOnOrderedSets(t *testing.T) {
	left := CollectOrderedSet(Slice[int]{3, 1, 3, 2}.Values())
	right := CollectOrderedSet(Slice[int]{3, 1, 2}.Values())
	if !left.Values().EqualBy(right.Values(), func(v int) int { return v }) {
		t.Fatal("want true: dedup leaves the same first-seen order")
	}
}
