package collections

import (
	"errors"
	"testing"
)

// Every lazy operation has to hold the same handful of invariants on every iterator
// type. Rather than assert them once per method, each operation is registered here
// configured as a pass-through -- Filter accepting everything, Map applying identity
// -- so the expected output is always the source. Adding an operation should mean
// adding a line to a table, not writing another six tests.
//
// The invariants:
//
//   - it passes every element through in order
//   - it stays lazy, so it terminates on an infinite source
//   - a consumer break leaves the error nil
//   - an error from the callback aborts the stream
//   - an abort from upstream propagates
//   - an empty source yields nothing
//   - it can be iterated more than once, so any state it carries lives inside the
//     returned closure rather than beside it

var errBoom = errors.New("boom")

// srcLen is the length of the standard source. Take is registered with exactly this
// n, not a larger one: a Take that leaks its counter across iterations only misbehaves
// once the counter has been reached, so a generous n would hide the bug.
const srcLen = 4

// pairLen is the same for the standard pair source.
const pairLen = 3

// --- sources ------------------------------------------------------------

func ints(n int) Iter[int] {
	return func(yield func(int) bool) {
		for v := range n {
			if !yield(v) {
				return
			}
		}
	}
}

func infinite() Iter[int] {
	return func(yield func(int) bool) {
		for v := 0; ; v++ {
			if !yield(v) {
				return
			}
		}
	}
}

// abortAt yields 0..at-1 and then aborts.
func abortAt(at int) IterX[int] {
	return func(yield func(int) bool) error {
		for v := range at {
			if !yield(v) {
				return nil
			}
		}
		return errBoom
	}
}

func pairs(n int) Iter2[int, int] {
	return func(yield func(int, int) bool) {
		for v := range n {
			if !yield(v, v*10) {
				return
			}
		}
	}
}

func infinitePairs() Iter2[int, int] {
	return func(yield func(int, int) bool) {
		for v := 0; ; v++ {
			if !yield(v, v*10) {
				return
			}
		}
	}
}

func abortPairsAt(at int) Iter2X[int, int] {
	return func(yield func(int, int) bool) error {
		for v := range at {
			if !yield(v, v*10) {
				return nil
			}
		}
		return errBoom
	}
}

// --- operation tables ---------------------------------------------------

// plainOps are pass-through operations that keep an infallible stream infallible.
var plainOps = []struct {
	name string
	op   func(Iter[int]) Iter[int]
}{
	{"Filter", func(i Iter[int]) Iter[int] { return i.Filter(func(int) bool { return true }) }},
	{"Map", func(i Iter[int]) Iter[int] { return i.Map(func(v int) int { return v }) }},
	{"Take", func(i Iter[int]) Iter[int] { return i.Take(srcLen) }},
	{"Drop", func(i Iter[int]) Iter[int] { return i.Drop(0) }},
	{"TakeWhile", func(i Iter[int]) Iter[int] { return i.TakeWhile(func(int) bool { return true }) }},
	{"DropWhile", func(i Iter[int]) Iter[int] { return i.DropWhile(func(int) bool { return false }) }},
	{"Tap", func(i Iter[int]) Iter[int] { return i.Tap(func(int) {}) }},
	{"DedupBy", func(i Iter[int]) Iter[int] { return i.DedupBy(func(v int) int { return v }) }},
	// v+1 is never the zero value, so nothing is compacted away
	{"CompactBy", func(i Iter[int]) Iter[int] { return i.CompactBy(func(v int) int { return v + 1 }) }},
}

// fallibleOps are pass-through operations on an already-fallible stream.
var fallibleOps = []struct {
	name string
	op   func(IterX[int]) IterX[int]
}{
	{"Filter", func(i IterX[int]) IterX[int] { return i.Filter(func(int) bool { return true }) }},
	{"Map", func(i IterX[int]) IterX[int] { return i.Map(func(v int) int { return v }) }},
	{"FilterX", func(i IterX[int]) IterX[int] {
		return i.FilterX(func(int) (bool, error) { return true, nil })
	}},
	{"MapX", func(i IterX[int]) IterX[int] {
		return i.MapX(func(v int) (int, error) { return v, nil })
	}},
	{"Take", func(i IterX[int]) IterX[int] { return i.Take(srcLen) }},
	{"Drop", func(i IterX[int]) IterX[int] { return i.Drop(0) }},
	{"TakeWhile", func(i IterX[int]) IterX[int] { return i.TakeWhile(func(int) bool { return true }) }},
	{"DropWhile", func(i IterX[int]) IterX[int] { return i.DropWhile(func(int) bool { return false }) }},
	{"TakeWhileX", func(i IterX[int]) IterX[int] {
		return i.TakeWhileX(func(int) (bool, error) { return true, nil })
	}},
	{"DropWhileX", func(i IterX[int]) IterX[int] {
		return i.DropWhileX(func(int) (bool, error) { return false, nil })
	}},
	{"Tap", func(i IterX[int]) IterX[int] { return i.Tap(func(int) {}) }},
	{"DedupBy", func(i IterX[int]) IterX[int] { return i.DedupBy(func(v int) int { return v }) }},
	{"CompactBy", func(i IterX[int]) IterX[int] {
		return i.CompactBy(func(v int) int { return v + 1 })
	}},
}

// liftingOps take an infallible stream into the fallible family. Each is registered
// twice: once with a callback that never fails, once with one that fails on 2.
var liftingOps = []struct {
	name string
	ok   func(Iter[int]) IterX[int]
	fail func(Iter[int]) IterX[int]
}{
	{
		"FilterX",
		func(i Iter[int]) IterX[int] { return i.FilterX(func(int) (bool, error) { return true, nil }) },
		func(i Iter[int]) IterX[int] {
			return i.FilterX(func(v int) (bool, error) {
				if v == 2 {
					return false, errBoom
				}
				return true, nil
			})
		},
	},
	{
		"MapX",
		func(i Iter[int]) IterX[int] { return i.MapX(func(v int) (int, error) { return v, nil }) },
		func(i Iter[int]) IterX[int] {
			return i.MapX(func(v int) (int, error) {
				if v == 2 {
					return 0, errBoom
				}
				return v, nil
			})
		},
	},
	{
		"TakeWhileX",
		func(i Iter[int]) IterX[int] {
			return i.TakeWhileX(func(int) (bool, error) { return true, nil })
		},
		func(i Iter[int]) IterX[int] {
			return i.TakeWhileX(func(v int) (bool, error) {
				if v == 2 {
					return false, errBoom
				}
				return true, nil
			})
		},
	},
	{
		"DropWhileX",
		func(i Iter[int]) IterX[int] {
			return i.DropWhileX(func(int) (bool, error) { return false, nil })
		},
		// keeps dropping until it fails, so the error arrives while still dropping
		func(i Iter[int]) IterX[int] {
			return i.DropWhileX(func(v int) (bool, error) {
				if v == 2 {
					return false, errBoom
				}
				return true, nil
			})
		},
	},
}

// pairOps are pass-through operations on an infallible pair stream.
var pairOps = []struct {
	name string
	op   func(Iter2[int, int]) Iter2[int, int]
}{
	{"Filter", func(i Iter2[int, int]) Iter2[int, int] {
		return i.Filter(func(int, int) bool { return true })
	}},
	{"Map", func(i Iter2[int, int]) Iter2[int, int] {
		return i.Map(func(k, v int) (int, int) { return k, v })
	}},
	{"MapKeys", func(i Iter2[int, int]) Iter2[int, int] {
		return i.MapKeys(func(k, _ int) int { return k })
	}},
	{"MapValues", func(i Iter2[int, int]) Iter2[int, int] {
		return i.MapValues(func(_, v int) int { return v })
	}},
	{"Take", func(i Iter2[int, int]) Iter2[int, int] { return i.Take(pairLen) }},
	{"Drop", func(i Iter2[int, int]) Iter2[int, int] { return i.Drop(0) }},
	{"TakeWhile", func(i Iter2[int, int]) Iter2[int, int] {
		return i.TakeWhile(func(int, int) bool { return true })
	}},
	{"DropWhile", func(i Iter2[int, int]) Iter2[int, int] {
		return i.DropWhile(func(int, int) bool { return false })
	}},
	{"Tap", func(i Iter2[int, int]) Iter2[int, int] { return i.Tap(func(int, int) {}) }},
	{"DedupBy", func(i Iter2[int, int]) Iter2[int, int] {
		return i.DedupBy(func(k, _ int) int { return k })
	}},
	{"CompactBy", func(i Iter2[int, int]) Iter2[int, int] {
		return i.CompactBy(func(k, _ int) int { return k + 1 })
	}},
}

// falliblePairOps are pass-through operations on an already-fallible pair stream.
var falliblePairOps = []struct {
	name string
	op   func(Iter2X[int, int]) Iter2X[int, int]
}{
	{"Filter", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.Filter(func(int, int) bool { return true })
	}},
	{"Map", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.Map(func(k, v int) (int, int) { return k, v })
	}},
	{"MapKeys", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.MapKeys(func(k, _ int) int { return k })
	}},
	{"MapValues", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.MapValues(func(_, v int) int { return v })
	}},
	{"FilterX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.FilterX(func(int, int) (bool, error) { return true, nil })
	}},
	{"MapX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.MapX(func(k, v int) (int, int, error) { return k, v, nil })
	}},
	{"MapKeysX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.MapKeysX(func(k, _ int) (int, error) { return k, nil })
	}},
	{"MapValuesX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.MapValuesX(func(_, v int) (int, error) { return v, nil })
	}},
	{"Take", func(i Iter2X[int, int]) Iter2X[int, int] { return i.Take(pairLen) }},
	{"Drop", func(i Iter2X[int, int]) Iter2X[int, int] { return i.Drop(0) }},
	{"TakeWhile", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.TakeWhile(func(int, int) bool { return true })
	}},
	{"DropWhile", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.DropWhile(func(int, int) bool { return false })
	}},
	{"TakeWhileX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.TakeWhileX(func(int, int) (bool, error) { return true, nil })
	}},
	{"DropWhileX", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.DropWhileX(func(int, int) (bool, error) { return false, nil })
	}},
	{"Tap", func(i Iter2X[int, int]) Iter2X[int, int] { return i.Tap(func(int, int) {}) }},
	{"DedupBy", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.DedupBy(func(k, _ int) int { return k })
	}},
	{"CompactBy", func(i Iter2X[int, int]) Iter2X[int, int] {
		return i.CompactBy(func(k, _ int) int { return k + 1 })
	}},
}

// --- single-value invariants --------------------------------------------

func TestPassesEveryElementThrough(t *testing.T) {
	for _, op := range plainOps {
		t.Run(op.name, func(t *testing.T) {
			assertInts(t, op.op(ints(srcLen)).Native(), 0, 1, 2, 3)
		})
	}
	for _, op := range fallibleOps {
		t.Run("X/"+op.name, func(t *testing.T) {
			got, err := op.op(ints(srcLen).X()).Native()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertInts(t, got, 0, 1, 2, 3)
		})
	}
	for _, op := range liftingOps {
		t.Run("lift/"+op.name, func(t *testing.T) {
			got, err := op.ok(ints(srcLen)).Native()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertInts(t, got, 0, 1, 2, 3)
		})
	}
}

// A non-lazy operation drains its source, so this hangs rather than failing. The
// test timeout is the backstop.
func TestStaysLazy(t *testing.T) {
	for _, op := range plainOps {
		t.Run(op.name, func(t *testing.T) {
			assertStopsAfterThree(t, countBreaking(op.op(infinite())))
		})
	}
	for _, op := range fallibleOps {
		t.Run("X/"+op.name, func(t *testing.T) {
			var err error
			assertStopsAfterThree(t, countBreaking(op.op(infinite().X()).Iter(&err)))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestConsumerBreakIsNotAnError(t *testing.T) {
	for _, op := range fallibleOps {
		t.Run(op.name, func(t *testing.T) {
			var err error
			for range op.op(ints(100).X()).Iter(&err) {
				break
			}
			if err != nil {
				t.Fatalf("a break must not set err, got %v", err)
			}
		})
	}
}

func TestUpstreamAbortPropagates(t *testing.T) {
	for _, op := range fallibleOps {
		t.Run(op.name, func(t *testing.T) {
			got, err := op.op(abortAt(2)).Native()
			if !errors.Is(err, errBoom) {
				t.Fatalf("want errBoom, got %v", err)
			}
			if got != nil {
				t.Fatalf("an abort must discard the partial result, got %v", got)
			}
		})
	}
}

func TestCallbackErrorAborts(t *testing.T) {
	for _, op := range liftingOps {
		t.Run(op.name, func(t *testing.T) {
			var seen int
			counted := Iter[int](func(yield func(int) bool) {
				for v := range 6 {
					seen++
					if !yield(v) {
						return
					}
				}
			})
			if _, e := op.fail(counted).Native(); !errors.Is(e, errBoom) {
				t.Fatalf("want errBoom, got %v", e)
			}
			if seen > 3 {
				t.Fatalf("iteration must stop at the failing element, saw %d", seen)
			}
		})
	}
}

func TestEmptySourceYieldsNothing(t *testing.T) {
	for _, op := range plainOps {
		t.Run(op.name, func(t *testing.T) {
			if got := op.op(ints(0)).Native(); len(got) != 0 {
				t.Fatalf("want empty, got %v", got)
			}
		})
	}
	for _, op := range fallibleOps {
		t.Run("X/"+op.name, func(t *testing.T) {
			got, err := op.op(ints(0).X()).Native()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("want empty, got %v", got)
			}
		})
	}
}

// An operation that keeps its counter outside the returned closure passes on the
// first iteration and yields nothing on the second. Take and Drop are the first
// operations in the package that can get this wrong.
func TestIsReIterable(t *testing.T) {
	for _, op := range plainOps {
		t.Run(op.name, func(t *testing.T) {
			i := op.op(ints(srcLen))
			assertInts(t, i.Native(), 0, 1, 2, 3)
			assertInts(t, i.Native(), 0, 1, 2, 3)
		})
	}
	for _, op := range fallibleOps {
		t.Run("X/"+op.name, func(t *testing.T) {
			i := op.op(ints(srcLen).X())
			for pass := range 2 {
				got, err := i.Native()
				if err != nil {
					t.Fatalf("pass %d: unexpected error: %v", pass, err)
				}
				assertInts(t, got, 0, 1, 2, 3)
			}
		})
	}
	for _, op := range pairOps {
		t.Run("pair/"+op.name, func(t *testing.T) {
			i := op.op(pairs(3))
			assertPairs(t, i)
			assertPairs(t, i)
		})
	}
}

// --- pair invariants ----------------------------------------------------

func TestPairsPassEveryPairThrough(t *testing.T) {
	for _, op := range pairOps {
		t.Run(op.name, func(t *testing.T) {
			assertPairs(t, op.op(pairs(3)))
		})
	}
	for _, op := range falliblePairOps {
		t.Run("X/"+op.name, func(t *testing.T) {
			var err error
			assertPairs(t, op.op(pairs(3).X()).Iter2(&err))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestPairsStayLazy(t *testing.T) {
	for _, op := range pairOps {
		t.Run(op.name, func(t *testing.T) {
			var n int
			for range op.op(infinitePairs()) {
				n++
				if n == 3 {
					break
				}
			}
			assertStopsAfterThree(t, n)
		})
	}
}

func TestPairsUpstreamAbortPropagates(t *testing.T) {
	for _, op := range falliblePairOps {
		t.Run(op.name, func(t *testing.T) {
			_, err := CollectMapX(op.op(abortPairsAt(2)))
			if !errors.Is(err, errBoom) {
				t.Fatalf("want errBoom, got %v", err)
			}
		})
	}
}

func TestPairsConsumerBreakIsNotAnError(t *testing.T) {
	for _, op := range falliblePairOps {
		t.Run(op.name, func(t *testing.T) {
			var err error
			for range op.op(pairs(100).X()).Iter2(&err) {
				break
			}
			if err != nil {
				t.Fatalf("a break must not set err, got %v", err)
			}
		})
	}
}

// --- assertions ---------------------------------------------------------

func assertInts(t *testing.T, got []int, want ...int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("want %v, got %v", want, got)
		}
	}
}

func assertPairs(t *testing.T, i Iter2[int, int]) {
	t.Helper()
	var n int
	for k, v := range i {
		if k != n || v != n*10 {
			t.Fatalf("pair %d: got (%d, %d)", n, k, v)
		}
		n++
	}
	if n != 3 {
		t.Fatalf("want 3 pairs, got %d", n)
	}
}

func countBreaking(i Iter[int]) int {
	var n int
	for range i {
		n++
		if n == 3 {
			break
		}
	}
	return n
}

func assertStopsAfterThree(t *testing.T, n int) {
	t.Helper()
	if n != 3 {
		t.Fatalf("want 3 elements before the break, got %d", n)
	}
}
