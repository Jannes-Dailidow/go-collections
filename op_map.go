package collections

// Map yields fn applied to each element. Lazy.
func (i Iter[T]) Map[U any](fn func(T) U) Iter[U] {
	return func(yield func(U) bool) {
		for t := range i {
			if !yield(fn(t)) {
				return
			}
		}
	}
}

// MapX yields fn applied to each element, and aborts on the first error. Lazy.
func (i Iter[T]) MapX[U any](fn func(T) (U, error)) IterX[U] {
	return i.X().MapX(fn)
}

// Map yields fn applied to each pair, replacing both halves. Lazy.
func (i Iter2[K, V]) Map[K2, V2 any](fn func(K, V) (K2, V2)) Iter2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range i {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

// MapKeys yields each pair with its key replaced by fn. Lazy.
func (i Iter2[K, V]) MapKeys[K2 any](fn func(K, V) K2) Iter2[K2, V] {
	return func(yield func(K2, V) bool) {
		for k, v := range i {
			if !yield(fn(k, v), v) {
				return
			}
		}
	}
}

// MapValues yields each pair with its value replaced by fn. Lazy.
func (i Iter2[K, V]) MapValues[V2 any](fn func(K, V) V2) Iter2[K, V2] {
	return func(yield func(K, V2) bool) {
		for k, v := range i {
			if !yield(k, fn(k, v)) {
				return
			}
		}
	}
}

// MapX yields fn applied to each pair, and aborts on the first error. Lazy.
func (i Iter2[K, V]) MapX[K2, V2 any](fn func(K, V) (K2, V2, error)) Iter2X[K2, V2] {
	return i.X().MapX(fn)
}

// MapKeysX replaces each key with fn, and aborts on the first error. Lazy.
func (i Iter2[K, V]) MapKeysX[K2 any](fn func(K, V) (K2, error)) Iter2X[K2, V] {
	return i.X().MapKeysX(fn)
}

// MapValuesX replaces each value with fn, and aborts on the first error. Lazy.
func (i Iter2[K, V]) MapValuesX[V2 any](fn func(K, V) (V2, error)) Iter2X[K, V2] {
	return i.X().MapValuesX(fn)
}

// Map yields fn applied to each element, preserving any abort. Lazy.
func (i IterX[T]) Map[U any](fn func(T) U) IterX[U] {
	return i.MapX(func(t T) (U, error) {
		return fn(t), nil
	})
}

// MapX yields fn applied to each element. An error from fn aborts the stream, as
// does an abort from upstream. Lazy.
func (i IterX[T]) MapX[U any](fn func(T) (U, error)) IterX[U] {
	return func(yield func(U) bool) error {
		return i.each(func(t T) (bool, error) {
			u, err := fn(t)
			if err != nil {
				return false, err
			}
			return yield(u), nil
		})
	}
}

// Map yields fn applied to each pair, preserving any abort. Lazy.
func (i Iter2X[K, V]) Map[K2, V2 any](fn func(K, V) (K2, V2)) Iter2X[K2, V2] {
	return i.MapX(func(k K, v V) (K2, V2, error) {
		k2, v2 := fn(k, v)
		return k2, v2, nil
	})
}

// MapKeys yields each pair with its key replaced by fn, preserving any abort. Lazy.
func (i Iter2X[K, V]) MapKeys[K2 any](fn func(K, V) K2) Iter2X[K2, V] {
	return i.MapKeysX(func(k K, v V) (K2, error) {
		return fn(k, v), nil
	})
}

// MapValues yields each pair with its value replaced by fn, preserving any abort.
// Lazy.
func (i Iter2X[K, V]) MapValues[V2 any](fn func(K, V) V2) Iter2X[K, V2] {
	return i.MapValuesX(func(k K, v V) (V2, error) {
		return fn(k, v), nil
	})
}

// MapX yields fn applied to each pair. An error from fn aborts the stream, as does
// an abort from upstream. Lazy.
func (i Iter2X[K, V]) MapX[K2, V2 any](fn func(K, V) (K2, V2, error)) Iter2X[K2, V2] {
	return func(yield func(K2, V2) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			k2, v2, err := fn(k, v)
			if err != nil {
				return false, err
			}
			return yield(k2, v2), nil
		})
	}
}

// MapKeysX replaces each key with fn. An error from fn aborts the stream, as does an
// abort from upstream. Lazy.
func (i Iter2X[K, V]) MapKeysX[K2 any](fn func(K, V) (K2, error)) Iter2X[K2, V] {
	return func(yield func(K2, V) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			k2, err := fn(k, v)
			if err != nil {
				return false, err
			}
			return yield(k2, v), nil
		})
	}
}

// MapValuesX replaces each value with fn. An error from fn aborts the stream, as
// does an abort from upstream. Lazy.
func (i Iter2X[K, V]) MapValuesX[V2 any](fn func(K, V) (V2, error)) Iter2X[K, V2] {
	return func(yield func(K, V2) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			v2, err := fn(k, v)
			if err != nil {
				return false, err
			}
			return yield(k, v2), nil
		})
	}
}
