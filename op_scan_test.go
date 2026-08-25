package collections

import (
	"errors"
	"testing"
)

func TestScanYieldsTheRunningAccumulator(t *testing.T) {
	got := ints(5).Scan(0, func(acc, v int) int { return acc + v }).Native()
	assertInts(t, got, 0, 1, 3, 6, 10)
}

// One value out per value in, and the seed is not one of them.
func TestScanDoesNotYieldTheInitialValue(t *testing.T) {
	got := ints(3).Scan(100, func(acc, v int) int { return acc + v }).Native()
	assertInts(t, got, 100, 101, 103)
	if len(got) != 3 {
		t.Fatalf("want one value per element, got %d", len(got))
	}
}

func TestScanOnAnEmptyStreamYieldsNothing(t *testing.T) {
	if got := ints(0).Scan(7, func(acc, v int) int { return acc + v }).Native(); len(got) != 0 {
		t.Fatalf("want nothing, got %v", got)
	}
}

// The accumulator has to live inside the returned closure.
func TestScanIsReIterable(t *testing.T) {
	i := ints(4).Scan(0, func(acc, v int) int { return acc + v })
	assertInts(t, i.Native(), 0, 1, 3, 6)
	assertInts(t, i.Native(), 0, 1, 3, 6)
}

func TestScanIsLazy(t *testing.T) {
	var pulls int
	got := counting(1000, &pulls).Scan(0, func(acc, v int) int { return acc + v }).Take(3).Native()
	assertInts(t, got, 0, 1, 3)
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

func TestScanXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(5).ScanX(0, func(acc, v int) (int, error) {
		if v == 3 {
			return 0, errBoom
		}
		return acc + v, nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestScanCanChangeType(t *testing.T) {
	got := ints(3).Scan(Slice[int]{}, func(acc Slice[int], v int) Slice[int] {
		return append(acc, v*10)
	}).Slice()
	if len(got) != 3 || len(got[2]) != 3 || got[2][2] != 20 {
		t.Fatalf("got %v", got)
	}
}
