package collections

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestReduceSeedsWithTheZeroValue(t *testing.T) {
	if got := ints(5).Reduce(func(acc, v int) int { return acc + v }); got != 10 {
		t.Fatalf("want 10, got %d", got)
	}
}

func TestFoldStartsFromInit(t *testing.T) {
	if got := ints(5).Fold(100, func(acc, v int) int { return acc + v }); got != 110 {
		t.Fatalf("want 110, got %d", got)
	}
}

func TestReduceCanChangeType(t *testing.T) {
	got := ints(3).Reduce(func(acc string, v int) string {
		return acc + strconv.Itoa(v)
	})
	if got != "012" {
		t.Fatalf("want 012, got %q", got)
	}
}

func TestFoldOnAnEmptyStreamReturnsInit(t *testing.T) {
	if got := ints(0).Fold(7, func(acc, v int) int { return acc + v }); got != 7 {
		t.Fatalf("want 7, got %d", got)
	}
}

// An abort discards the accumulator rather than returning a partial fold, matching
// the collectors.
func TestFoldXDiscardsTheAccumulatorOnError(t *testing.T) {
	got, err := ints(6).FoldX(0, func(acc, v int) (int, error) {
		if v == 3 {
			return 0, errBoom
		}
		return acc + v, nil
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != 0 {
		t.Fatalf("want the zero value, got %d", got)
	}
}

func TestReduceXOnAnAbortingStream(t *testing.T) {
	_, err := abortAt(2).Reduce(func(acc, v int) int { return acc + v })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairFoldSeesBothHalves(t *testing.T) {
	got := pairs(4).Fold(0, func(acc, k, v int) int { return acc + k + v })
	if got != 66 {
		t.Fatalf("want 66, got %d", got)
	}
}

func TestPairReduceBuildsAString(t *testing.T) {
	got := pairs(3).Reduce(func(acc string, k, v int) string {
		return acc + strconv.Itoa(k) + ":" + strconv.Itoa(v) + " "
	})
	if strings.TrimSpace(got) != "0:0 1:10 2:20" {
		t.Fatalf("got %q", got)
	}
}

// Fold is how the collections are built, so it must be able to build one.
func TestFoldCanBuildACollection(t *testing.T) {
	got := ints(4).Fold(Slice[int]{}, func(acc Slice[int], v int) Slice[int] {
		return append(acc, v*2)
	})
	assertInts(t, got.Native(), 0, 2, 4, 6)
}

// Reduce and Scan must agree: the last value Scan emits is what Reduce returns.
func TestScanEndsWhereReduceLands(t *testing.T) {
	add := func(acc, v int) int { return acc + v }
	scanned := ints(6).Scan(0, add).Last()
	reduced := ints(6).Reduce(add)
	if scanned == nil || *scanned != reduced {
		t.Fatalf("Scan ended at %v, Reduce returned %d", scanned, reduced)
	}
}
