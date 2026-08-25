package collections

// Dedup yields the first occurrence of each element and drops the repeats. Lazy,
// and holds every distinct element seen so far.
//
// It is a function rather than a method because a method cannot tighten T to
// comparable. For the fluent form use [Iter.DedupBy] with an identity key.
func Dedup[T comparable](i Iter[T]) Iter[T] {
	return i.DedupBy(func(t T) T {
		return t
	})
}

// DedupX is the fallible counterpart of [Dedup].
func DedupX[T comparable](i IterX[T]) IterX[T] {
	return i.DedupBy(func(t T) T {
		return t
	})
}

// Compact drops the zero-valued elements. Lazy.
func Compact[T comparable](i Iter[T]) Iter[T] {
	return i.CompactBy(func(t T) T {
		return t
	})
}

// CompactX is the fallible counterpart of [Compact].
func CompactX[T comparable](i IterX[T]) IterX[T] {
	return i.CompactBy(func(t T) T {
		return t
	})
}

// DedupBy yields the first element for each distinct key and drops the rest. Lazy,
// and holds every distinct key seen so far.
func (i Iter[T]) DedupBy[K comparable](fn func(T) K) Iter[T] {
	return func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for t := range i {
			k := fn(t)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			if !yield(t) {
				return
			}
		}
	}
}

// CompactBy drops the elements whose key is the zero value. Lazy.
func (i Iter[T]) CompactBy[K comparable](fn func(T) K) Iter[T] {
	return func(yield func(T) bool) {
		var zero K
		for t := range i {
			if fn(t) == zero {
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

// DedupBy yields the first pair for each distinct key and drops the rest. Lazy, and
// holds every distinct key seen so far.
func (i Iter2[K, V]) DedupBy[K2 comparable](fn func(K, V) K2) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		seen := make(map[K2]struct{})
		for k, v := range i {
			k2 := fn(k, v)
			if _, ok := seen[k2]; ok {
				continue
			}
			seen[k2] = struct{}{}
			if !yield(k, v) {
				return
			}
		}
	}
}

// CompactBy drops the pairs whose key is the zero value. Lazy.
func (i Iter2[K, V]) CompactBy[K2 comparable](fn func(K, V) K2) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		var zero K2
		for k, v := range i {
			if fn(k, v) == zero {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// DedupBy yields the first element for each distinct key, preserving any abort.
// Lazy, and holds every distinct key seen so far.
func (i IterX[T]) DedupBy[K comparable](fn func(T) K) IterX[T] {
	return func(yield func(T) bool) error {
		seen := make(map[K]struct{})
		return i.each(func(t T) (bool, error) {
			k := fn(t)
			if _, ok := seen[k]; ok {
				return true, nil
			}
			seen[k] = struct{}{}
			return yield(t), nil
		})
	}
}

// CompactBy drops the elements whose key is the zero value, preserving any abort.
// Lazy.
func (i IterX[T]) CompactBy[K comparable](fn func(T) K) IterX[T] {
	return func(yield func(T) bool) error {
		var zero K
		return i.each(func(t T) (bool, error) {
			if fn(t) == zero {
				return true, nil
			}
			return yield(t), nil
		})
	}
}

// DedupBy yields the first pair for each distinct key, preserving any abort. Lazy,
// and holds every distinct key seen so far.
func (i Iter2X[K, V]) DedupBy[K2 comparable](fn func(K, V) K2) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		seen := make(map[K2]struct{})
		return i.each(func(k K, v V) (bool, error) {
			k2 := fn(k, v)
			if _, ok := seen[k2]; ok {
				return true, nil
			}
			seen[k2] = struct{}{}
			return yield(k, v), nil
		})
	}
}

// CompactBy drops the pairs whose key is the zero value, preserving any abort. Lazy.
func (i Iter2X[K, V]) CompactBy[K2 comparable](fn func(K, V) K2) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		var zero K2
		return i.each(func(k K, v V) (bool, error) {
			if fn(k, v) == zero {
				return true, nil
			}
			return yield(k, v), nil
		})
	}
}
