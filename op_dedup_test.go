package collections

import (
	"errors"
	"strings"
	"testing"
)

func TestDedupKeepsTheFirstOccurrence(t *testing.T) {
	src := Slice[int]{3, 1, 3, 2, 1, 3}
	assertInts(t, Dedup(src.Values()).Native(), 3, 1, 2)
}

func TestDedupByUsesTheKey(t *testing.T) {
	src := Slice[string]{"Ada", "ada", "Bob", "BOB", "cyd"}
	got := src.Values().DedupBy(strings.ToLower).Slice()
	if len(got) != 3 || got[0] != "Ada" || got[1] != "Bob" {
		t.Fatalf("the first spelling of each key must win, got %v", got)
	}
}

func TestDedupIsLazy(t *testing.T) {
	var pulls int
	got := Dedup(counting(1000, &pulls)).Take(3).Native()
	assertInts(t, got, 0, 1, 2)
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

// The seen set has to live inside the returned closure, or a second pass finds
// everything already seen and yields nothing.
func TestDedupIsReIterable(t *testing.T) {
	src := Slice[int]{1, 1, 2}
	i := Dedup(src.Values())
	assertInts(t, i.Native(), 1, 2)
	assertInts(t, i.Native(), 1, 2)
}

func TestDedupXPreservesTheAbort(t *testing.T) {
	if _, err := DedupX(abortAt(3)).Native(); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestCompactDropsZeroValues(t *testing.T) {
	src := Slice[int]{0, 1, 0, 2}
	assertInts(t, Compact(src.Values()).Native(), 1, 2)
}

func TestCompactOnStrings(t *testing.T) {
	src := Slice[string]{"a", "", "b"}
	got := Compact(src.Values()).Slice()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("got %v", got)
	}
}

// CompactBy tests the key, not the element, which is what makes it useful on
// structs.
func TestCompactByTestsTheKey(t *testing.T) {
	type user struct {
		name string
		id   int
	}
	src := Slice[user]{{"ada", 1}, {"nobody", 0}, {"bob", 2}}
	got := src.Values().CompactBy(func(u user) int { return u.id }).Slice()
	if len(got) != 2 || got[0].name != "ada" || got[1].name != "bob" {
		t.Fatalf("got %v", got)
	}
}

func TestPairDedupByKeepsTheFirstPairPerKey(t *testing.T) {
	src := func(yield func(string, int) bool) {
		for _, e := range []Entry[string, int]{
			{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "a", Value: 9},
		} {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
	got := CollectOrderedMap(Iter2[string, int](src).DedupBy(func(k string, _ int) string {
		return k
	}))
	// unlike CollectOrderedMap on its own, which would update a to 9
	if got.Native()["a"] != 1 {
		t.Fatalf("the first pair for a key must win, got %v", got.Native())
	}
}

func TestPairCompactByDropsZeroKeys(t *testing.T) {
	got := CollectMap(pairs(4).CompactBy(func(_, v int) int { return v }))
	if len(got) != 3 {
		t.Fatalf("the pair whose value is 0 must be dropped, got %v", got)
	}
}

// Dedup and CollectOrderedSet do the same job, so the choice between them is lazy
// versus eager, not one of behaviour.
func TestDedupAndCollectOrderedSetAgree(t *testing.T) {
	src := Slice[int]{5, 3, 5, 9, 3}
	lazy := Dedup(src.Values()).Native()
	eager := CollectOrderedSet(src.Values()).Native()
	assertInts(t, lazy, eager...)
}
