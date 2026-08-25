package collections

// Filter yields the elements fn accepts. Lazy.
func (i Iter[T]) Filter(fn func(T) bool) Iter[T] {
	return func(yield func(T) bool) {
		for t := range i {
			if fn(t) && !yield(t) {
				return
			}
		}
	}
}

// FilterX yields the elements fn accepts, and aborts on the first error. Lazy.
func (i Iter[T]) FilterX(fn func(T) (bool, error)) IterX[T] {
	return i.X().FilterX(fn)
}

// Filter yields the pairs fn accepts. Lazy.
func (i Iter2[K, V]) Filter(fn func(K, V) bool) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range i {
			if fn(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// FilterX yields the pairs fn accepts, and aborts on the first error. Lazy.
func (i Iter2[K, V]) FilterX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return i.X().FilterX(fn)
}

// Filter yields the elements fn accepts, preserving any abort. Lazy.
func (i IterX[T]) Filter(fn func(T) bool) IterX[T] {
	return i.FilterX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// FilterX yields the elements fn accepts. An error from fn aborts the stream, as
// does an abort from upstream. Lazy.
func (i IterX[T]) FilterX(fn func(T) (bool, error)) IterX[T] {
	return func(yield func(T) bool) error {
		return i.each(func(t T) (bool, error) {
			ok, err := fn(t)
			if err != nil {
				return false, err
			}
			if !ok {
				return true, nil
			}
			return yield(t), nil
		})
	}
}

// Filter yields the pairs fn accepts, preserving any abort. Lazy.
func (i Iter2X[K, V]) Filter(fn func(K, V) bool) Iter2X[K, V] {
	return i.FilterX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// FilterX yields the pairs fn accepts. An error from fn aborts the stream, as does
// an abort from upstream. Lazy.
func (i Iter2X[K, V]) FilterX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			ok, err := fn(k, v)
			if err != nil {
				return false, err
			}
			if !ok {
				return true, nil
			}
			return yield(k, v), nil
		})
	}
}
