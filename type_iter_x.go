package collections

import "iter"

type IterX[T any] func(yield func(T) bool) error

func (i IterX[T]) Iter(err *error) Iter[T] {
	return func(yield func(T) bool) {
		e := i(yield)
		if e != nil && err != nil {
			*err = e
		}
	}
}

func (i IterX[T]) Seq(err *error) iter.Seq[T] {
	return iter.Seq[T](i.Iter(err))
}

func (i IterX[T]) Slice() (Slice[T], error) {
	var result Slice[T]
	err := i(func(t T) bool {
		result = append(result, t)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i IterX[T]) Native() ([]T, error) {
	result, err := i.Slice()
	return []T(result), err
}
