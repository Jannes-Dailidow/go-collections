package collections

import "iter"

type Iter[T any] func(yield func(T) bool)

func FromSeq[T any](seq iter.Seq[T]) Iter[T] {
	return Iter[T](seq)
}

func (i Iter[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](i)
}

func (i Iter[T]) X() IterX[T] {
	return func(yield func(T) bool) error {
		i(yield)
		return nil
	}
}

func (i Iter[T]) Slice() Slice[T] {
	result, _ := i.X().Slice()
	return result
}

func (i Iter[T]) Native() []T {
	return []T(i.Slice())
}
