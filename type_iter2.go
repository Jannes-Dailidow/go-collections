package collections

import "iter"

type Iter2[K, V any] func(yield func(K, V) bool)

func FromSeq2[K, V any](seq iter.Seq2[K, V]) Iter2[K, V] {
	return Iter2[K, V](seq)
}

func (i Iter2[K, V]) Seq2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](i)
}

func (i Iter2[K, V]) X() Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		i(yield)
		return nil
	}
}

func (i Iter2[K, V]) Keys() Iter[K] {
	return func(yield func(K) bool) {
		for k := range i {
			if !yield(k) {
				return
			}
		}
	}
}

func (i Iter2[K, V]) Values() Iter[V] {
	return func(yield func(V) bool) {
		for _, v := range i {
			if !yield(v) {
				return
			}
		}
	}
}
