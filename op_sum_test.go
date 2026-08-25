package collections

import (
	"errors"
	"strconv"
	"testing"
)

func TestSum(t *testing.T) {
	if got := Sum(ints(5)); got != 10 {
		t.Fatalf("want 10, got %d", got)
	}
}

func TestSumOnAnEmptyStreamIsZero(t *testing.T) {
	if got := Sum(ints(0)); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}
}

func TestSumOnFloats(t *testing.T) {
	src := Slice[float64]{1.5, 2.25}
	if got := Sum(src.Values()); got != 3.75 {
		t.Fatalf("want 3.75, got %v", got)
	}
}

// A named type with a numeric underlying type satisfies Number through its ~ terms.
func TestSumOnANamedNumericType(t *testing.T) {
	type cents int64
	src := Slice[cents]{100, 250}
	if got := Sum(src.Values()); got != 350 {
		t.Fatalf("want 350, got %v", got)
	}
}

func TestSumByUsesTheKey(t *testing.T) {
	type line struct {
		item string
		qty  int
	}
	src := Slice[line]{{"a", 2}, {"b", 3}}
	if got := src.Values().SumBy(func(l line) int { return l.qty }); got != 5 {
		t.Fatalf("want 5, got %d", got)
	}
}

func TestSumByOnAnAbortingStream(t *testing.T) {
	_, err := abortAt(2).SumBy(func(v int) int { return v })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairSumBy(t *testing.T) {
	if got := pairs(4).SumBy(func(_, v int) int { return v }); got != 60 {
		t.Fatalf("want 60, got %d", got)
	}
}

// The documented route when the value can fail to compute.
func TestSumOverFallibleValuesViaMapX(t *testing.T) {
	ok := Slice[string]{"1", "2", "3"}
	got, err := SumX(ok.Values().MapX(strconv.Atoi))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 6 {
		t.Fatalf("want 6, got %d", got)
	}

	bad := Slice[string]{"1", "nope"}
	if _, err := SumX(bad.Values().MapX(strconv.Atoi)); err == nil {
		t.Fatal("want a parse error")
	}
}
