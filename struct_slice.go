package collections

import "slices"

// Len returns the number of elements. O(1), unlike [Iter.Count], which drains a
// stream to find out.
func (s Slice[T]) Len() int {
	return len(s)
}

// IsEmpty reports whether the slice has no elements. O(1).
func (s Slice[T]) IsEmpty() bool {
	return len(s) == 0
}

// At returns the element at n, or nil when n is out of range. O(1), and the reason
// this is on the slice rather than on an iterator.
//
// The pointer is into the slice, so writing through it changes the element.
func (s Slice[T]) At(n int) *T {
	if n < 0 || n >= len(s) {
		return nil
	}
	return &s[n]
}

// Append returns the slice with values added, so a chain does not have to break out
// to the builtin. Like append, it may or may not reuse the backing array.
func (s Slice[T]) Append(values ...T) Slice[T] {
	return append(s, values...)
}

// Insert returns the slice with values inserted at n. Out-of-range n returns the
// slice unchanged, rather than panicking the way slices.Insert does.
func (s Slice[T]) Insert(n int, values ...T) Slice[T] {
	if n < 0 || n > len(s) {
		return s
	}
	return slices.Insert(s, n, values...)
}

// DeleteAt returns the slice with the element at n removed. Out-of-range n returns
// the slice unchanged.
func (s Slice[T]) DeleteAt(n int) Slice[T] {
	if n < 0 || n >= len(s) {
		return s
	}
	return slices.Delete(s, n, n+1)
}

// Clone returns a copy with its own backing array. Shallow: the elements themselves
// are not copied.
func (s Slice[T]) Clone() Slice[T] {
	return slices.Clone(s)
}

// Grow returns a slice with room for n more elements without another allocation.
// There is no Grow on [Map] or [Set]: Go offers no way to add capacity to a map after
// it is made, so size those with make instead.
func (s Slice[T]) Grow(n int) Slice[T] {
	return slices.Grow(s, n)
}
