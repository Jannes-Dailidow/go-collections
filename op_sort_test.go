package collections

import (
	"errors"
	"strings"
	"testing"
)

func TestIterSortBy(t *testing.T) {
	src := Slice[int]{3, 1, 2}
	got := src.Values().SortBy(func(v int) int { return v }).Native()
	assertInts(t, got, 1, 2, 3)
}

// Stable, so equal keys keep the order they arrived in.
func TestSortByIsStable(t *testing.T) {
	type tagged struct {
		id  int
		key int
	}
	src := Slice[tagged]{{1, 9}, {2, 5}, {3, 9}, {4, 5}}
	got := src.Values().SortBy(func(x tagged) int { return x.key }).Slice()
	var ids []int
	for _, x := range got {
		ids = append(ids, x.id)
	}
	assertInts(t, ids, 2, 4, 1, 3)
}

// The iterator form buffers a copy, so the source is left alone.
func TestIterSortByDoesNotTouchTheSource(t *testing.T) {
	src := Slice[int]{3, 1, 2}
	src.Values().SortBy(func(v int) int { return v }).Native()
	assertInts(t, src.Native(), 3, 1, 2)
}

// The slice form is documented to sort in place.
func TestSliceSortBySortsInPlace(t *testing.T) {
	src := Slice[int]{3, 1, 2}
	got := src.SortBy(func(v int) int { return v })
	assertInts(t, src.Native(), 1, 2, 3)
	assertInts(t, got.Native(), 1, 2, 3)
}

func TestSortFuncForKeysThatAreNotOrdered(t *testing.T) {
	src := Slice[string]{"bb", "a", "ccc"}
	got := src.Values().SortFunc(func(a, b string) int {
		return len(a) - len(b)
	}).Slice()
	if got[0] != "a" || got[2] != "ccc" {
		t.Fatalf("got %v", got)
	}
}

func TestIterReverse(t *testing.T) {
	assertInts(t, ints(4).Reverse().Native(), 3, 2, 1, 0)
	if got := ints(0).Reverse().Native(); len(got) != 0 {
		t.Fatalf("want nothing, got %v", got)
	}
}

func TestSliceReverseIsInPlace(t *testing.T) {
	src := Slice[int]{1, 2, 3}
	src.Reverse()
	assertInts(t, src.Native(), 3, 2, 1)
}

func TestSortByIsReIterable(t *testing.T) {
	i := Slice[int]{3, 1, 2}.Values().SortBy(func(v int) int { return v })
	assertInts(t, i.Native(), 1, 2, 3)
	assertInts(t, i.Native(), 1, 2, 3)
}

// Sorting has to see everything first, so an abort anywhere surfaces and nothing is
// yielded.
func TestSortByOnAnAbortingStream(t *testing.T) {
	got, err := abortAt(3).SortBy(func(v int) int { return v }).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != nil {
		t.Fatalf("nothing should be yielded, got %v", got)
	}
}

func TestReverseOnAnAbortingStream(t *testing.T) {
	if _, err := abortAt(3).Reverse().Native(); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairSortBy(t *testing.T) {
	got := CollectOrderedMap(pairs(4).SortBy(func(k, _ int) int { return -k }))
	var keys []int
	for k := range got.All() {
		keys = append(keys, k)
	}
	assertInts(t, keys, 3, 2, 1, 0)
}

func TestPairReverse(t *testing.T) {
	got := CollectOrderedMap(pairs(3).Reverse())
	assertInts(t, got.Keys().Native(), 2, 1, 0)
}

// Sorting an ordered collection must repair its index, or every positional method
// afterwards is wrong.
func TestOrderedSetSortByRepairsTheIndex(t *testing.T) {
	s := orderedInts(5, 3, 9)
	s.SortBy(func(v int) int { return v })
	assertInts(t, s.Native(), 3, 5, 9)
	if s.IndexOf(3) != 0 || s.IndexOf(5) != 1 || s.IndexOf(9) != 2 {
		t.Fatalf("index not repaired: 3 at %d, 5 at %d, 9 at %d",
			s.IndexOf(3), s.IndexOf(5), s.IndexOf(9))
	}
	if got := s.At(0); got == nil || *got != 3 {
		t.Fatalf("At disagrees with the new order: %v", got)
	}
}

func TestOrderedSetReverseRepairsTheIndex(t *testing.T) {
	s := orderedInts(1, 2, 3)
	s.Reverse()
	assertInts(t, s.Native(), 3, 2, 1)
	if s.IndexOf(3) != 0 || s.IndexOf(1) != 2 {
		t.Fatal("index not repaired")
	}
}

func TestOrderedMapSortByRepairsTheIndex(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)
	m.SortBy(func(k string, _ int) string { return k })
	assertInts(t, m.Values().Native(), 1, 2, 3)
	if m.IndexOf("a") != 0 || m.IndexOf("c") != 2 {
		t.Fatalf("index not repaired: a at %d, c at %d", m.IndexOf("a"), m.IndexOf("c"))
	}
	if got := m.At(0); got == nil || got.Key != "a" {
		t.Fatalf("At disagrees with the new order: %v", got)
	}
}

func TestOrderedMapSortByValue(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 3)
	m.Put("b", 1)
	m.SortBy(func(_ string, v int) int { return v })
	if got := m.At(0); got == nil || got.Key != "b" {
		t.Fatalf("got %v", got)
	}
}

func TestOrderedMapReverseRepairsTheIndex(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Reverse()
	if m.IndexOf("b") != 0 || m.IndexOf("a") != 1 {
		t.Fatal("index not repaired")
	}
}

func TestSortOnOrderedCollectionsIsNilSafe(t *testing.T) {
	var s *OrderedSet[int]
	var m *OrderedMap[string, int]
	if s.SortBy(func(v int) int { return v }) != nil || s.Reverse() != nil {
		t.Fatal("want nil")
	}
	if m.SortBy(func(k string, _ int) string { return k }) != nil || m.Reverse() != nil {
		t.Fatal("want nil")
	}
}

// Sorting composes with the rest of the chain.
func TestSortInTheMiddleOfAChain(t *testing.T) {
	src := Slice[string]{"pear", "fig", "apple", "kiwi"}
	got := src.Values().
		Filter(func(s string) bool { return len(s) > 3 }).
		SortBy(func(s string) string { return s }).
		Map(strings.ToUpper).
		Slice()
	if len(got) != 3 || got[0] != "APPLE" || got[2] != "PEAR" {
		t.Fatalf("got %v", got)
	}
}
