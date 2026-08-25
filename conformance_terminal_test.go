package collections

import (
	"errors"
	"testing"
)

// Terminal operations cannot be registered as pass-throughs, so they get their own
// table and their own invariant: an abort on the very first pull has to surface from
// every one of them. Any terminal operation that drops the error each returns fails
// here, which is the mistake this table exists to catch.
//
// The source aborts before yielding anything, because a short-circuiting operation
// that stops before reaching an abort is right not to report it -- see
// TestTakeBeforeAnAbortSeesNoError.
var terminalOps = []struct {
	name string
	run  func(IterX[int]) error
}{
	{"Slice", func(i IterX[int]) error { _, err := i.Slice(); return err }},
	{"Native", func(i IterX[int]) error { _, err := i.Native(); return err }},
	{"Each", func(i IterX[int]) error { return i.Each(func(int) {}) }},
	{"EachX", func(i IterX[int]) error { return i.EachX(func(int) error { return nil }) }},
	{"Find", func(i IterX[int]) error {
		_, err := i.Find(func(int) bool { return true })
		return err
	}},
	{"FindLast", func(i IterX[int]) error {
		_, err := i.FindLast(func(int) bool { return true })
		return err
	}},
	{"First", func(i IterX[int]) error { _, err := i.First(); return err }},
	{"Last", func(i IterX[int]) error { _, err := i.Last(); return err }},
	{"Contains", func(i IterX[int]) error {
		_, err := i.Contains(func(int) bool { return true })
		return err
	}},
	{"Every", func(i IterX[int]) error {
		_, err := i.Every(func(int) bool { return true })
		return err
	}},
	{"Index", func(i IterX[int]) error {
		_, err := i.Index(func(int) bool { return true })
		return err
	}},
	{"Count", func(i IterX[int]) error { _, err := i.Count(); return err }},
	{"CountBy", func(i IterX[int]) error {
		_, err := i.CountBy(func(int) bool { return true })
		return err
	}},
	{"IsEmpty", func(i IterX[int]) error { _, err := i.IsEmpty(); return err }},
	{"Reduce", func(i IterX[int]) error {
		_, err := i.Reduce(func(acc, v int) int { return acc + v })
		return err
	}},
	{"Fold", func(i IterX[int]) error {
		_, err := i.Fold(0, func(acc, v int) int { return acc + v })
		return err
	}},
	{"MinBy", func(i IterX[int]) error {
		_, err := i.MinBy(func(v int) int { return v })
		return err
	}},
	{"MaxBy", func(i IterX[int]) error {
		_, err := i.MaxBy(func(v int) int { return v })
		return err
	}},
	{"SumBy", func(i IterX[int]) error {
		_, err := i.SumBy(func(v int) int { return v })
		return err
	}},
	{"EqualBy", func(i IterX[int]) error {
		_, err := i.EqualBy(ints(4), func(v int) int { return v })
		return err
	}},
}

var terminalPairOps = []struct {
	name string
	run  func(Iter2X[int, int]) error
}{
	{"Each", func(i Iter2X[int, int]) error { return i.Each(func(int, int) {}) }},
	{"EachX", func(i Iter2X[int, int]) error {
		return i.EachX(func(int, int) error { return nil })
	}},
	{"Find", func(i Iter2X[int, int]) error {
		_, err := i.Find(func(int, int) bool { return true })
		return err
	}},
	{"FindLast", func(i Iter2X[int, int]) error {
		_, err := i.FindLast(func(int, int) bool { return true })
		return err
	}},
	{"First", func(i Iter2X[int, int]) error { _, err := i.First(); return err }},
	{"Last", func(i Iter2X[int, int]) error { _, err := i.Last(); return err }},
	{"Contains", func(i Iter2X[int, int]) error {
		_, err := i.Contains(func(int, int) bool { return true })
		return err
	}},
	{"Every", func(i Iter2X[int, int]) error {
		_, err := i.Every(func(int, int) bool { return true })
		return err
	}},
	{"Count", func(i Iter2X[int, int]) error { _, err := i.Count(); return err }},
	{"IsEmpty", func(i Iter2X[int, int]) error { _, err := i.IsEmpty(); return err }},
	{"Fold", func(i Iter2X[int, int]) error {
		_, err := i.Fold(0, func(acc, k, v int) int { return acc + k })
		return err
	}},
	{"MinBy", func(i Iter2X[int, int]) error {
		_, err := i.MinBy(func(k, _ int) int { return k })
		return err
	}},
	{"SumBy", func(i Iter2X[int, int]) error {
		_, err := i.SumBy(func(k, _ int) int { return k })
		return err
	}},
	{"EqualBy", func(i Iter2X[int, int]) error {
		_, err := i.EqualBy(pairs(3), func(k, _ int) int { return k })
		return err
	}},
	{"CollectMapX", func(i Iter2X[int, int]) error { _, err := CollectMapX(i); return err }},
}

func TestTerminalOpsSurfaceAnImmediateAbort(t *testing.T) {
	for _, op := range terminalOps {
		t.Run(op.name, func(t *testing.T) {
			if err := op.run(abortAt(0)); !errors.Is(err, errBoom) {
				t.Fatalf("want errBoom, got %v", err)
			}
		})
	}
	for _, op := range terminalPairOps {
		t.Run("pair/"+op.name, func(t *testing.T) {
			if err := op.run(abortPairsAt(0)); !errors.Is(err, errBoom) {
				t.Fatalf("want errBoom, got %v", err)
			}
		})
	}
}

// The mirror case: a clean stream must not invent an error.
func TestTerminalOpsReportNoErrorOnACleanStream(t *testing.T) {
	for _, op := range terminalOps {
		t.Run(op.name, func(t *testing.T) {
			if err := op.run(ints(4).X()); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
	for _, op := range terminalPairOps {
		t.Run("pair/"+op.name, func(t *testing.T) {
			if err := op.run(pairs(3).X()); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
