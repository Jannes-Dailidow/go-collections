package collections

import "testing"

func TestSliceLenAndIsEmpty(t *testing.T) {
	if got := (Slice[int]{1, 2}).Len(); got != 2 {
		t.Fatalf("want 2, got %d", got)
	}
	if !(Slice[int]{}).IsEmpty() || (Slice[int]{1}).IsEmpty() {
		t.Fatal("IsEmpty is wrong")
	}
	var nilSlice Slice[int]
	if !nilSlice.IsEmpty() || nilSlice.Len() != 0 {
		t.Fatal("a nil slice is empty")
	}
}

func TestSliceAt(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	if got := s.At(1); got == nil || *got != 2 {
		t.Fatalf("got %v", got)
	}
	for _, n := range []int{-1, 3, 99} {
		if got := s.At(n); got != nil {
			t.Fatalf("At(%d) must be nil, got %v", n, *got)
		}
	}
}

// Unlike the iterator's Find, this points into the slice.
func TestSliceAtPointsIntoTheSlice(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	*s.At(1) = 99
	assertInts(t, s.Native(), 1, 99, 3)
}

func TestSliceAppend(t *testing.T) {
	assertInts(t, (Slice[int]{1}).Append(2, 3).Native(), 1, 2, 3)
	var empty Slice[int]
	assertInts(t, empty.Append(1).Native(), 1)
}

func TestSliceInsert(t *testing.T) {
	s := Slice[int]{1, 4}
	assertInts(t, s.Insert(1, 2, 3).Native(), 1, 2, 3, 4)
	assertInts(t, s.Insert(0, 0).Native(), 0, 1, 4)
	assertInts(t, (Slice[int]{1, 2}).Insert(2, 3).Native(), 1, 2, 3)
}

// slices.Insert panics out of range; this returns the slice untouched instead.
func TestSliceInsertOutOfRangeIsANoOp(t *testing.T) {
	s := Slice[int]{1, 2}
	for _, n := range []int{-1, 3} {
		assertInts(t, s.Insert(n, 9).Native(), 1, 2)
	}
}

func TestSliceDeleteAt(t *testing.T) {
	assertInts(t, (Slice[int]{1, 2, 3}).DeleteAt(1).Native(), 1, 3)
	s := Slice[int]{1, 2}
	for _, n := range []int{-1, 2} {
		assertInts(t, s.DeleteAt(n).Native(), 1, 2)
	}
}

func TestSliceCloneIsIndependent(t *testing.T) {
	s := Slice[int]{1, 2, 3}
	c := s.Clone()
	c[0] = 99
	assertInts(t, s.Native(), 1, 2, 3)
	assertInts(t, c.Native(), 99, 2, 3)
}

func TestSliceGrowAddsCapacity(t *testing.T) {
	s := Slice[int]{1, 2}
	grown := s.Grow(10)
	if cap(grown) < 12 {
		t.Fatalf("want capacity for 12, got %d", cap(grown))
	}
	assertInts(t, grown.Native(), 1, 2)
}

func TestSliceBackward(t *testing.T) {
	s := Slice[string]{"a", "b", "c"}
	var order []string
	var indexes []int
	for n, v := range s.Backward() {
		indexes = append(indexes, n)
		order = append(order, v)
	}
	assertInts(t, indexes, 2, 1, 0)
	if order[0] != "c" || order[2] != "a" {
		t.Fatalf("got %v", order)
	}
}

// Backward is lazy and does not copy, so a break stops it early.
func TestSliceBackwardIsLazy(t *testing.T) {
	s := Slice[int]{1, 2, 3, 4}
	var seen int
	for range s.Backward() {
		seen++
		break
	}
	if seen != 1 {
		t.Fatalf("want 1, got %d", seen)
	}
}
