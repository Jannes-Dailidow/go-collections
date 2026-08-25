package collections

import (
	"errors"
	"testing"
)

func nested() Iter[Iter[int]] {
	return func(yield func(Iter[int]) bool) {
		for _, n := range []int{2, 0, 3} {
			if !yield(ints(n)) {
				return
			}
		}
	}
}

func TestFlattenConcatenates(t *testing.T) {
	assertInts(t, Flatten(nested()).Native(), 0, 1, 0, 1, 2)
}

func TestFlattenSkipsEmptyInnerStreams(t *testing.T) {
	empty := Iter[Iter[int]](func(yield func(Iter[int]) bool) {
		for range 3 {
			if !yield(ints(0)) {
				return
			}
		}
	})
	if got := Flatten(empty).Native(); len(got) != 0 {
		t.Fatalf("want nothing, got %v", got)
	}
}

func TestFlattenSlices(t *testing.T) {
	src := Slice[Slice[int]]{{1, 2}, {}, {3}}
	assertInts(t, FlattenSlices(src.Values()).Native(), 1, 2, 3)
}

func TestFlattenIsLazy(t *testing.T) {
	var built int
	lazy := Iter[Iter[int]](func(yield func(Iter[int]) bool) {
		for range 1000 {
			built++
			if !yield(ints(2)) {
				return
			}
		}
	})
	assertInts(t, Flatten(lazy).Take(3).Native(), 0, 1, 0)
	if built != 2 {
		t.Fatalf("want 2 inner streams built, got %d", built)
	}
}

func TestFlatMapExpandsEachElement(t *testing.T) {
	got := ints(3).FlatMap(func(v int) Iter[int] {
		return Slice[int]{v, v * 10}.Values()
	}).Native()
	assertInts(t, got, 0, 0, 1, 10, 2, 20)
}

// A break inside an inner stream has to stop the outer stream too.
func TestFlatMapStopsTheOuterStreamOnBreak(t *testing.T) {
	var outer int
	var n int
	for range ints(100).Tap(func(int) { outer++ }).FlatMap(func(v int) Iter[int] {
		return Slice[int]{v, v}.Values()
	}) {
		n++
		if n == 3 {
			break
		}
	}
	if outer != 2 {
		t.Fatalf("want the outer stream stopped after 2 elements, got %d", outer)
	}
}

func TestFlatMapXPropagatesAnInnerAbort(t *testing.T) {
	_, err := ints(3).FlatMapX(func(v int) (IterX[int], error) {
		if v == 1 {
			return abortAt(1), nil
		}
		return ints(1).X(), nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("an inner abort must surface, got %v", err)
	}
}

func TestFlatMapXPropagatesTheCallbackError(t *testing.T) {
	_, err := ints(3).FlatMapX(func(v int) (IterX[int], error) {
		if v == 1 {
			return nil, errBoom
		}
		return ints(1).X(), nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestFlatMapOnAFallibleStream(t *testing.T) {
	got, err := ints(2).X().FlatMap(func(v int) Iter[int] {
		return Slice[int]{v, v}.Values()
	}).Native()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertInts(t, got, 0, 0, 1, 1)
}
