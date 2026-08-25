package collections

// Each runs fn over every element. Terminal.
func (i Iter[T]) Each(fn func(T)) {
	for t := range i {
		fn(t)
	}
}

// EachX runs fn over every element and stops at the first error. Terminal.
func (i Iter[T]) EachX(fn func(T) error) error {
	return i.X().EachX(fn)
}

// Tap runs fn on every element and passes it through unchanged, for logging and
// metrics in the middle of a chain. Lazy.
func (i Iter[T]) Tap(fn func(T)) Iter[T] {
	return func(yield func(T) bool) {
		for t := range i {
			fn(t)
			if !yield(t) {
				return
			}
		}
	}
}

// Each runs fn over every pair. Terminal.
func (i Iter2[K, V]) Each(fn func(K, V)) {
	for k, v := range i {
		fn(k, v)
	}
}

// EachX runs fn over every pair and stops at the first error. Terminal.
func (i Iter2[K, V]) EachX(fn func(K, V) error) error {
	return i.X().EachX(fn)
}

// Tap runs fn on every pair and passes it through unchanged. Lazy.
func (i Iter2[K, V]) Tap(fn func(K, V)) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range i {
			fn(k, v)
			if !yield(k, v) {
				return
			}
		}
	}
}

// Each runs fn over every element and returns any abort. Terminal.
func (i IterX[T]) Each(fn func(T)) error {
	return i.EachX(func(t T) error {
		fn(t)
		return nil
	})
}

// EachX runs fn over every element and stops at the first error, whether fn's or the
// stream's. Terminal.
func (i IterX[T]) EachX(fn func(T) error) error {
	return i.each(func(t T) (bool, error) {
		if err := fn(t); err != nil {
			return false, err
		}
		return true, nil
	})
}

// Tap runs fn on every element and passes it through unchanged, preserving any
// abort. Lazy.
func (i IterX[T]) Tap(fn func(T)) IterX[T] {
	return func(yield func(T) bool) error {
		return i.each(func(t T) (bool, error) {
			fn(t)
			return yield(t), nil
		})
	}
}

// Each runs fn over every pair and returns any abort. Terminal.
func (i Iter2X[K, V]) Each(fn func(K, V)) error {
	return i.EachX(func(k K, v V) error {
		fn(k, v)
		return nil
	})
}

// EachX runs fn over every pair and stops at the first error, whether fn's or the
// stream's. Terminal.
func (i Iter2X[K, V]) EachX(fn func(K, V) error) error {
	return i.each(func(k K, v V) (bool, error) {
		if err := fn(k, v); err != nil {
			return false, err
		}
		return true, nil
	})
}

// Tap runs fn on every pair and passes it through unchanged, preserving any abort.
// Lazy.
func (i Iter2X[K, V]) Tap(fn func(K, V)) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			fn(k, v)
			return yield(k, v), nil
		})
	}
}
