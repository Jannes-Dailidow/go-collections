package collections

import (
	"errors"
	"strconv"
	"testing"
)

func TestKeyByPairsEachElementWithItsKey(t *testing.T) {
	got := CollectMap(ints(3).KeyBy(strconv.Itoa))
	if len(got) != 3 || got["2"] != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestKeyByXReturnsTheCallbackError(t *testing.T) {
	_, err := CollectMapX(ints(4).KeyByX(func(v int) (string, error) {
		if v == 2 {
			return "", errBoom
		}
		return strconv.Itoa(v), nil
	}))
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestEnumerateNumbersFromZero(t *testing.T) {
	src := Slice[string]{"a", "b", "c"}
	var positions []int
	for n, s := range src.Values().Enumerate() {
		positions = append(positions, n)
		if s != string(rune('a'+n)) {
			t.Fatalf("position %d carries %q", n, s)
		}
	}
	assertInts(t, positions, 0, 1, 2)
}

// The counter has to live inside the returned closure.
func TestEnumerateIsReIterable(t *testing.T) {
	i := ints(3).Enumerate()
	for pass := range 2 {
		var positions []int
		for n := range i {
			positions = append(positions, n)
		}
		if len(positions) != 3 || positions[0] != 0 || positions[2] != 2 {
			t.Fatalf("pass %d: got %v", pass, positions)
		}
	}
}

func TestEnumerateNumbersWhatSurvivedAFilter(t *testing.T) {
	var positions []int
	for n := range ints(10).Filter(func(v int) bool { return v%3 == 0 }).Enumerate() {
		positions = append(positions, n)
	}
	assertInts(t, positions, 0, 1, 2, 3)
}

func TestEnumerateOnAFallibleStreamCarriesTheAbort(t *testing.T) {
	_, err := CollectMapX(abortAt(2).Enumerate())
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestZipPairsByPosition(t *testing.T) {
	letters := Slice[string]{"a", "b", "c"}
	got := CollectMap(ints(3).Zip(letters.Values()))
	if got[0] != "a" || got[2] != "c" {
		t.Fatalf("got %v", got)
	}
}

func TestZipStopsAtTheShorterSide(t *testing.T) {
	short := Slice[string]{"a", "b"}
	if n := len(CollectMap(ints(9).Zip(short.Values()))); n != 2 {
		t.Fatalf("a short right side must cut the zip, got %d pairs", n)
	}
	long := Slice[string]{"a", "b", "c", "d"}
	if n := len(CollectMap(ints(2).Zip(long.Values()))); n != 2 {
		t.Fatalf("a short left side must cut the zip, got %d pairs", n)
	}
}

// iter.Pull holds live state, so it has to be created inside the returned closure.
func TestZipIsReIterable(t *testing.T) {
	letters := Slice[string]{"a", "b", "c"}
	i := ints(3).Zip(letters.Values())
	for pass := range 2 {
		got := CollectMap(i)
		if len(got) != 3 || got[1] != "b" {
			t.Fatalf("pass %d: got %v", pass, got)
		}
	}
}

func TestZipIsLazy(t *testing.T) {
	var pulls int
	letters := Slice[string]{"a", "b", "c", "d", "e"}
	var n int
	for range counting(1000, &pulls).Zip(letters.Values()) {
		n++
		if n == 2 {
			break
		}
	}
	if pulls != 2 {
		t.Fatalf("want 2 pulls from the left side, got %d", pulls)
	}
}

func TestCollapseLeavesThePairLayer(t *testing.T) {
	got := pairs(3).Collapse(func(k, v int) int { return k + v }).Native()
	assertInts(t, got, 0, 11, 22)
}

func TestCollapseXReturnsTheCallbackError(t *testing.T) {
	_, err := pairs(4).CollapseX(func(k, _ int) (int, error) {
		if k == 2 {
			return 0, errBoom
		}
		return k, nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairsCarriesEntriesIntoTheSingleValueLayer(t *testing.T) {
	got := Pairs(pairs(3)).Slice()
	if len(got) != 3 || got[1].Key != 1 || got[1].Value != 10 {
		t.Fatalf("got %v", got)
	}
}

// The point of Pairs is reaching the single-value operations from a pair stream.
func TestPairsThenASingleValueOperation(t *testing.T) {
	got := Pairs(pairs(5)).
		Filter(func(e Entry[int, int]) bool { return e.Value >= 20 }).
		Map(func(e Entry[int, int]) int { return e.Key }).
		Native()
	assertInts(t, got, 2, 3, 4)
}

func TestPairsXCarriesTheAbort(t *testing.T) {
	if _, err := PairsX(abortPairsAt(2)).Native(); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}
