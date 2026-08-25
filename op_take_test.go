package collections

import (
	"errors"
	"testing"
)

// counting wraps a source so a test can assert how far it was pulled.
func counting(n int, pulls *int) Iter[int] {
	return func(yield func(int) bool) {
		for v := range n {
			*pulls++
			if !yield(v) {
				return
			}
		}
	}
}

func TestTakeStopsAtN(t *testing.T) {
	assertInts(t, ints(6).Take(2).Native(), 0, 1)
}

func TestTakeOfMoreThanAvailableYieldsEverything(t *testing.T) {
	assertInts(t, ints(3).Take(99).Native(), 0, 1, 2)
}

func TestTakeOfNoneYieldsNothing(t *testing.T) {
	for _, n := range []int{0, -1} {
		if got := ints(6).Take(n).Native(); len(got) != 0 {
			t.Fatalf("Take(%d) must yield nothing, got %v", n, got)
		}
	}
}

// Take has to stop pulling, not just stop yielding, or it is useless on an
// expensive source.
func TestTakeStopsPullingAtN(t *testing.T) {
	var pulls int
	counting(1000, &pulls).Take(3).Native()
	if pulls != 3 {
		t.Fatalf("want 3 pulls, got %d", pulls)
	}
}

// Stopping early means an abort further along the stream is never reached, so it
// must not surface.
func TestTakeBeforeAnAbortSeesNoError(t *testing.T) {
	got, err := abortAt(3).Take(2).Native()
	if err != nil {
		t.Fatalf("the abort was never reached, so err must be nil, got %v", err)
	}
	assertInts(t, got, 0, 1)
}

func TestTakeReachingAnAbortSurfacesIt(t *testing.T) {
	if _, err := abortAt(3).Take(99).Native(); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestDropSkipsTheFirstN(t *testing.T) {
	assertInts(t, ints(5).Drop(2).Native(), 2, 3, 4)
}

func TestDropOfMoreThanAvailableYieldsNothing(t *testing.T) {
	if got := ints(3).Drop(99).Native(); len(got) != 0 {
		t.Fatalf("want nothing, got %v", got)
	}
}

func TestTakeWhileStopsAtTheFirstReject(t *testing.T) {
	assertInts(t, ints(6).TakeWhile(func(v int) bool { return v < 3 }).Native(), 0, 1, 2)
}

// A later element that would pass must not come back once the run has ended.
func TestTakeWhileDoesNotResume(t *testing.T) {
	src := Slice[int]{1, 2, 9, 1, 2}
	got := src.Values().TakeWhile(func(v int) bool { return v < 5 }).Slice()
	assertInts(t, got.Native(), 1, 2)
}

func TestDropWhileYieldsFromTheFirstReject(t *testing.T) {
	assertInts(t, ints(6).DropWhile(func(v int) bool { return v < 3 }).Native(), 3, 4, 5)
}

// Once the prefix is over the predicate is done, even if a later element would have
// matched it.
func TestDropWhileOnlyExaminesThePrefix(t *testing.T) {
	src := Slice[int]{1, 2, 9, 1, 2}
	var calls int
	got := src.Values().DropWhile(func(v int) bool {
		calls++
		return v < 5
	}).Slice()
	assertInts(t, got.Native(), 9, 1, 2)
	if calls != 3 {
		t.Fatalf("the predicate must stop being called after the prefix, got %d calls", calls)
	}
}

// The conformance suite registers these as pass-throughs, which cannot expose a
// leaked counter: a Take that does not truncate never reaches its limit, and a
// DropWhile that drops nothing never flips its flag. Only a truncating
// configuration does, so these are tested here rather than there.
func TestTruncatingOpsAreReIterable(t *testing.T) {
	cases := []struct {
		name string
		op   func(Iter[int]) Iter[int]
		want []int
	}{
		{"Take", func(i Iter[int]) Iter[int] { return i.Take(2) }, []int{0, 1}},
		{"Drop", func(i Iter[int]) Iter[int] { return i.Drop(2) }, []int{2, 3}},
		{"TakeWhile", func(i Iter[int]) Iter[int] {
			return i.TakeWhile(func(v int) bool { return v < 2 })
		}, []int{0, 1}},
		{"DropWhile", func(i Iter[int]) Iter[int] {
			return i.DropWhile(func(v int) bool { return v < 2 })
		}, []int{2, 3}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := c.op(ints(4))
			assertInts(t, i.Native(), c.want...)
			assertInts(t, i.Native(), c.want...)
		})
	}
	for _, c := range cases {
		t.Run("X/"+c.name, func(t *testing.T) {
			// the same operation reached through the fallible family
			i := c.op(ints(4).X().Iter(nil))
			assertInts(t, i.Native(), c.want...)
			assertInts(t, i.Native(), c.want...)
		})
	}
}

func TestFallibleTruncatingOpsAreReIterable(t *testing.T) {
	cases := []struct {
		name string
		op   func(IterX[int]) IterX[int]
		want []int
	}{
		{"Take", func(i IterX[int]) IterX[int] { return i.Take(2) }, []int{0, 1}},
		{"Drop", func(i IterX[int]) IterX[int] { return i.Drop(2) }, []int{2, 3}},
		{"TakeWhile", func(i IterX[int]) IterX[int] {
			return i.TakeWhile(func(v int) bool { return v < 2 })
		}, []int{0, 1}},
		{"DropWhile", func(i IterX[int]) IterX[int] {
			return i.DropWhile(func(v int) bool { return v < 2 })
		}, []int{2, 3}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := c.op(ints(4).X())
			for pass := range 2 {
				got, err := i.Native()
				if err != nil {
					t.Fatalf("pass %d: unexpected error: %v", pass, err)
				}
				assertInts(t, got, c.want...)
			}
		})
	}
}

func TestTakeWhileXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(6).TakeWhileX(func(v int) (bool, error) {
		if v == 2 {
			return false, errBoom
		}
		return true, nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestDropWhileXReturnsTheCallbackError(t *testing.T) {
	_, err := ints(6).DropWhileX(func(v int) (bool, error) {
		if v == 2 {
			return false, errBoom
		}
		return true, nil
	}).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestPairTakeAndDrop(t *testing.T) {
	got := CollectOrderedMap(pairs(5).Drop(1).Take(2))
	var keys []int
	for k := range got.All() {
		keys = append(keys, k)
	}
	assertInts(t, keys, 1, 2)
}

func TestPairTakeWhileAndDropWhile(t *testing.T) {
	kept := CollectMap(pairs(5).TakeWhile(func(k, _ int) bool { return k < 2 }))
	if len(kept) != 2 {
		t.Fatalf("want 2 pairs, got %v", kept)
	}
	rest := CollectMap(pairs(5).DropWhile(func(k, _ int) bool { return k < 3 }))
	if len(rest) != 2 || rest[3] != 30 {
		t.Fatalf("want pairs 3 and 4, got %v", rest)
	}
}

// Take and Drop compose into pagination, which is the reason they exist.
func TestTakeAndDropPaginate(t *testing.T) {
	page := func(n int) []int {
		return ints(10).Drop(n * 3).Take(3).Native()
	}
	assertInts(t, page(0), 0, 1, 2)
	assertInts(t, page(2), 6, 7, 8)
	assertInts(t, page(3), 9)
}
