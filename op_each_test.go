package collections

import (
	"errors"
	"testing"
)

func TestEachVisitsEveryElement(t *testing.T) {
	var sum int
	ints(4).Each(func(v int) { sum += v })
	if sum != 6 {
		t.Fatalf("want 6, got %d", sum)
	}
}

func TestEachXStopsAtTheFirstError(t *testing.T) {
	var seen int
	err := ints(6).EachX(func(v int) error {
		seen++
		if v == 2 {
			return errBoom
		}
		return nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if seen != 3 {
		t.Fatalf("must stop at the failing element, saw %d", seen)
	}
}

func TestEachOnAFallibleStreamReturnsTheAbort(t *testing.T) {
	var seen int
	err := abortAt(2).Each(func(int) { seen++ })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if seen != 2 {
		t.Fatalf("want the 2 elements before the abort, saw %d", seen)
	}
}

func TestPairEachVisitsEveryPair(t *testing.T) {
	var keys, values int
	pairs(3).Each(func(k, v int) {
		keys += k
		values += v
	})
	if keys != 3 || values != 30 {
		t.Fatalf("got %d and %d", keys, values)
	}
}

func TestPairEachXStopsAtTheFirstError(t *testing.T) {
	err := pairs(6).EachX(func(k, _ int) error {
		if k == 2 {
			return errBoom
		}
		return nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestTapSeesEveryElementAndChangesNothing(t *testing.T) {
	var seen []int
	got := ints(4).Tap(func(v int) { seen = append(seen, v) }).Native()
	assertInts(t, got, 0, 1, 2, 3)
	assertInts(t, seen, 0, 1, 2, 3)
}

// Tap must stay lazy, so it only sees what the consumer actually pulls.
func TestTapOnlySeesWhatIsConsumed(t *testing.T) {
	var seen int
	ints(100).Tap(func(int) { seen++ }).Take(3).Native()
	if seen != 3 {
		t.Fatalf("want 3, got %d", seen)
	}
}

func TestTapInTheMiddleOfAChain(t *testing.T) {
	var kept int
	got := ints(6).
		Filter(func(v int) bool { return v%2 == 0 }).
		Tap(func(int) { kept++ }).
		Map(func(v int) int { return v * 10 }).
		Native()
	assertInts(t, got, 0, 20, 40)
	if kept != 3 {
		t.Fatalf("Tap must see only what survived the Filter, got %d", kept)
	}
}

func TestPairTapSeesEveryPair(t *testing.T) {
	var n int
	got := CollectMap(pairs(3).Tap(func(int, int) { n++ }))
	if len(got) != 3 || n != 3 {
		t.Fatalf("got %v and %d taps", got, n)
	}
}
