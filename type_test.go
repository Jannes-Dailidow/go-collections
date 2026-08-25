package collections

import (
	"errors"
	"testing"
)

func TestIter2XKeysAndValuesYieldTheirHalf(t *testing.T) {
	keys, err := abortPairsAt(3).Keys().Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("Keys must carry the abort, got %v", err)
	}
	if keys != nil {
		t.Fatalf("the partial result must be discarded, got %v", keys)
	}

	got, err := pairs(3).X().Values().Native()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertInts(t, got, 0, 10, 20)
}

func TestIter2XKeysStopsOnAConsumerBreak(t *testing.T) {
	var err error
	var n int
	for range pairs(100).X().Keys().Iter(&err) {
		n++
		break
	}
	if err != nil {
		t.Fatalf("a break must not set err, got %v", err)
	}
	assertInts(t, []int{n}, 1)
}

// OrderedMap stores its order as []Entry, so these pin the behaviour that rename
// could have broken.
func TestOrderedMapKeepsInsertionOrder(t *testing.T) {
	entries := Slice[Entry[string, int]]{
		{Key: "c", Value: 1},
		{Key: "a", Value: 2},
		{Key: "b", Value: 3},
	}
	m := CollectOrderedMap(func(yield func(string, int) bool) {
		for _, e := range entries {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	})

	var order []string
	for k := range m.All() {
		order = append(order, k)
	}
	if len(order) != 3 || order[0] != "c" || order[2] != "b" {
		t.Fatalf("want c, a, b -- got %v", order)
	}
}

func TestOrderedMapReinsertKeepsPositionAndUpdatesValue(t *testing.T) {
	m := CollectOrderedMap(func(yield func(string, int) bool) {
		for _, e := range []Entry[string, int]{
			{Key: "a", Value: 1}, {Key: "b", Value: 2}, {Key: "a", Value: 9},
		} {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	})

	var keys []string
	var values []int
	for k, v := range m.All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
		t.Fatalf("a must keep its original position, got %v", keys)
	}
	assertInts(t, values, 9, 2)
}

func TestOrderedMapIsNilSafe(t *testing.T) {
	var m *OrderedMap[string, int]
	for range m.All() {
		t.Fatal("a nil OrderedMap must yield nothing")
	}
	for range m.Keys() {
		t.Fatal("a nil OrderedMap must yield no keys")
	}
	for range m.Values() {
		t.Fatal("a nil OrderedMap must yield no values")
	}
	if m.Map() != nil || m.Native() != nil {
		t.Fatal("a nil OrderedMap must convert to nil")
	}
}

func TestOrderedSetIsNilSafe(t *testing.T) {
	var s *OrderedSet[int]
	for range s.Values() {
		t.Fatal("a nil OrderedSet must yield nothing")
	}
	for range s.All() {
		t.Fatal("a nil OrderedSet must yield nothing")
	}
	if s.Slice() != nil || s.Native() != nil || s.Set() != nil {
		t.Fatal("a nil OrderedSet must convert to nil")
	}
}
