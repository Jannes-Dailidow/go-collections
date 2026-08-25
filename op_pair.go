package collections

import "iter"

// KeyBy pairs every element with a key derived from it, entering the pair layer. K
// is unconstrained; tighten it when collecting, with [CollectMap] or
// [CollectOrderedMap]. Lazy.
func (i Iter[T]) KeyBy[K any](fn func(T) K) Iter2[K, T] {
	return func(yield func(K, T) bool) {
		for t := range i {
			if !yield(fn(t), t) {
				return
			}
		}
	}
}

// KeyByX pairs every element with a key derived from it, and aborts on the first
// error. Lazy.
func (i Iter[T]) KeyByX[K any](fn func(T) (K, error)) Iter2X[K, T] {
	return i.X().KeyByX(fn)
}

// Enumerate pairs every element with its position in the stream. It is what the root
// package's I suffix became. Lazy.
func (i Iter[T]) Enumerate() Iter2[int, T] {
	return func(yield func(int, T) bool) {
		var n int
		for t := range i {
			if !yield(n, t) {
				return
			}
			n++
		}
	}
}

// Zip pairs each element with the element at the same position in other, and stops
// as soon as either side runs out. Lazy.
func (i Iter[T]) Zip[U any](other Iter[U]) Iter2[T, U] {
	return func(yield func(T, U) bool) {
		next, stop := iter.Pull(other.Seq())
		defer stop()
		for t := range i {
			u, ok := next()
			if !ok || !yield(t, u) {
				return
			}
		}
	}
}

// Collapse turns every pair into a single value, leaving the pair layer. It is the
// FromMap of v0.1. Lazy.
func (i Iter2[K, V]) Collapse[T any](fn func(K, V) T) Iter[T] {
	return func(yield func(T) bool) {
		for k, v := range i {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

// CollapseX turns every pair into a single value, and aborts on the first error.
// Lazy.
func (i Iter2[K, V]) CollapseX[T any](fn func(K, V) (T, error)) IterX[T] {
	return i.X().CollapseX(fn)
}

// KeyBy pairs every element with a key derived from it, preserving any abort. Lazy.
func (i IterX[T]) KeyBy[K any](fn func(T) K) Iter2X[K, T] {
	return i.KeyByX(func(t T) (K, error) {
		return fn(t), nil
	})
}

// KeyByX pairs every element with a key derived from it. An error from fn aborts the
// stream, as does an abort from upstream. Lazy.
func (i IterX[T]) KeyByX[K any](fn func(T) (K, error)) Iter2X[K, T] {
	return func(yield func(K, T) bool) error {
		return i.each(func(t T) (bool, error) {
			k, err := fn(t)
			if err != nil {
				return false, err
			}
			return yield(k, t), nil
		})
	}
}

// Enumerate pairs every element with its position in the stream, preserving any
// abort. Lazy.
func (i IterX[T]) Enumerate() Iter2X[int, T] {
	return func(yield func(int, T) bool) error {
		var n int
		return i.each(func(t T) (bool, error) {
			if !yield(n, t) {
				return false, nil
			}
			n++
			return true, nil
		})
	}
}

// Zip pairs each element with the element at the same position in other, stopping as
// soon as either side runs out and preserving any abort. The other side is
// infallible: zipping two fallible streams would have two aborts to reconcile, and
// there is no single answer to which one wins. Lazy.
func (i IterX[T]) Zip[U any](other Iter[U]) Iter2X[T, U] {
	return func(yield func(T, U) bool) error {
		next, stop := iter.Pull(other.Seq())
		defer stop()
		return i.each(func(t T) (bool, error) {
			u, ok := next()
			if !ok {
				return false, nil
			}
			return yield(t, u), nil
		})
	}
}

// Collapse turns every pair into a single value, preserving any abort. Lazy.
func (i Iter2X[K, V]) Collapse[T any](fn func(K, V) T) IterX[T] {
	return i.CollapseX(func(k K, v V) (T, error) {
		return fn(k, v), nil
	})
}

// CollapseX turns every pair into a single value. An error from fn aborts the
// stream, as does an abort from upstream. Lazy.
func (i Iter2X[K, V]) CollapseX[T any](fn func(K, V) (T, error)) IterX[T] {
	return func(yield func(T) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			t, err := fn(k, v)
			if err != nil {
				return false, err
			}
			return yield(t), nil
		})
	}
}

// Pairs carries a pair stream into the single-value layer as [Entry] values, so the
// single-value operations apply to it. Lazy.
//
// It is a function rather than a method to break an instantiation cycle: Iter2[K, V]
// having a method that returns Iter[Entry[K, V]], while Iter[T] has KeyBy returning
// Iter2[K, T], lets the compiler chase Iter -> Iter2 -> Iter forever. A free
// function is only instantiated where it is called, so the cycle never forms. The
// alternative was making KeyBy a function, and KeyBy earns the method more.
func Pairs[K, V any](i Iter2[K, V]) Iter[Entry[K, V]] {
	return func(yield func(Entry[K, V]) bool) {
		for k, v := range i {
			if !yield(Entry[K, V]{Key: k, Value: v}) {
				return
			}
		}
	}
}

// PairsX is the fallible counterpart of [Pairs].
func PairsX[K, V any](i Iter2X[K, V]) IterX[Entry[K, V]] {
	return func(yield func(Entry[K, V]) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			return yield(Entry[K, V]{Key: k, Value: v}), nil
		})
	}
}
