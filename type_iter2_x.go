package collections

import "iter"

type Iter2X[K, V any] func(yield func(K, V) bool) error

func (i Iter2X[K, V]) Iter2(err *error) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		e := i(yield)
		if e != nil && err != nil {
			*err = e
		}
	}
}

func (i Iter2X[K, V]) Seq2(err *error) iter.Seq2[K, V] {
	return iter.Seq2[K, V](i.Iter2(err))
}

// Keys yields the keys of the stream, carrying any abort with them.
func (i Iter2X[K, V]) Keys() IterX[K] {
	return func(yield func(K) bool) error {
		return i(func(k K, _ V) bool {
			return yield(k)
		})
	}
}

// Values yields the values of the stream, carrying any abort with them.
func (i Iter2X[K, V]) Values() IterX[V] {
	return func(yield func(V) bool) error {
		return i(func(_ K, v V) bool {
			return yield(v)
		})
	}
}
