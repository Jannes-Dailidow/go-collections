package collections

import (
	"iter"
	"slices"
)

type Slice[T any] []T

func (s Slice[T]) Native() []T {
	return []T(s)
}

func (s Slice[T]) Values() Iter[T] {
	return func(yield func(T) bool) {
		for _, t := range s {
			if !yield(t) {
				return
			}
		}
	}
}

func (s Slice[T]) All() Iter2[int, T] {
	return func(yield func(int, T) bool) {
		for i, t := range s {
			if !yield(i, t) {
				return
			}
		}
	}
}

func (s Slice[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](s.Values())
}

func (s Slice[T]) Seq2() iter.Seq2[int, T] {
	return iter.Seq2[int, T](s.All())
}

// Backward iterates from the last element to the first, pairing each with its index.
// Lazy, and does not copy.
func (s Slice[T]) Backward() Iter2[int, T] {
	return Iter2[int, T](slices.Backward(s))
}
