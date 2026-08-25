package collections

// Take yields at most the first n elements and then stops pulling. Lazy.
func (i Iter[T]) Take(n int) Iter[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		var taken int
		for t := range i {
			if !yield(t) {
				return
			}
			taken++
			if taken >= n {
				return
			}
		}
	}
}

// Drop discards the first n elements and yields the rest. Lazy.
func (i Iter[T]) Drop(n int) Iter[T] {
	return func(yield func(T) bool) {
		var dropped int
		for t := range i {
			if dropped < n {
				dropped++
				continue
			}
			if !yield(t) {
				return
			}
		}
	}
}

// TakeWhile yields elements until fn rejects one, then stops pulling. Lazy.
func (i Iter[T]) TakeWhile(fn func(T) bool) Iter[T] {
	return func(yield func(T) bool) {
		for t := range i {
			if !fn(t) || !yield(t) {
				return
			}
		}
	}
}

// DropWhile discards elements until fn rejects one, and yields that one and every
// element after it. Lazy.
func (i Iter[T]) DropWhile(fn func(T) bool) Iter[T] {
	return func(yield func(T) bool) {
		dropping := true
		for t := range i {
			if dropping {
				if fn(t) {
					continue
				}
				dropping = false
			}
			if !yield(t) {
				return
			}
		}
	}
}

// TakeWhileX yields elements until fn rejects one, and aborts on the first error.
// Lazy.
func (i Iter[T]) TakeWhileX(fn func(T) (bool, error)) IterX[T] {
	return i.X().TakeWhileX(fn)
}

// DropWhileX discards elements until fn rejects one, and aborts on the first error.
// Lazy.
func (i Iter[T]) DropWhileX(fn func(T) (bool, error)) IterX[T] {
	return i.X().DropWhileX(fn)
}

// Take yields at most the first n pairs and then stops pulling. Lazy.
func (i Iter2[K, V]) Take(n int) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		var taken int
		for k, v := range i {
			if !yield(k, v) {
				return
			}
			taken++
			if taken >= n {
				return
			}
		}
	}
}

// Drop discards the first n pairs and yields the rest. Lazy.
func (i Iter2[K, V]) Drop(n int) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		var dropped int
		for k, v := range i {
			if dropped < n {
				dropped++
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// TakeWhile yields pairs until fn rejects one, then stops pulling. Lazy.
func (i Iter2[K, V]) TakeWhile(fn func(K, V) bool) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range i {
			if !fn(k, v) || !yield(k, v) {
				return
			}
		}
	}
}

// DropWhile discards pairs until fn rejects one, and yields that one and every pair
// after it. Lazy.
func (i Iter2[K, V]) DropWhile(fn func(K, V) bool) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		dropping := true
		for k, v := range i {
			if dropping {
				if fn(k, v) {
					continue
				}
				dropping = false
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

// TakeWhileX yields pairs until fn rejects one, and aborts on the first error. Lazy.
func (i Iter2[K, V]) TakeWhileX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return i.X().TakeWhileX(fn)
}

// DropWhileX discards pairs until fn rejects one, and aborts on the first error.
// Lazy.
func (i Iter2[K, V]) DropWhileX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return i.X().DropWhileX(fn)
}

// Take yields at most the first n elements and then stops pulling, preserving any
// abort reached before that. Lazy.
func (i IterX[T]) Take(n int) IterX[T] {
	return func(yield func(T) bool) error {
		if n <= 0 {
			return nil
		}
		var taken int
		return i.each(func(t T) (bool, error) {
			if !yield(t) {
				return false, nil
			}
			taken++
			return taken < n, nil
		})
	}
}

// Drop discards the first n elements and yields the rest, preserving any abort. Lazy.
func (i IterX[T]) Drop(n int) IterX[T] {
	return func(yield func(T) bool) error {
		var dropped int
		return i.each(func(t T) (bool, error) {
			if dropped < n {
				dropped++
				return true, nil
			}
			return yield(t), nil
		})
	}
}

// TakeWhile yields elements until fn rejects one, preserving any abort. Lazy.
func (i IterX[T]) TakeWhile(fn func(T) bool) IterX[T] {
	return i.TakeWhileX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// DropWhile discards elements until fn rejects one, preserving any abort. Lazy.
func (i IterX[T]) DropWhile(fn func(T) bool) IterX[T] {
	return i.DropWhileX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// TakeWhileX yields elements until fn rejects one. An error from fn aborts the
// stream, as does an abort from upstream reached before the rejection. Lazy.
func (i IterX[T]) TakeWhileX(fn func(T) (bool, error)) IterX[T] {
	return func(yield func(T) bool) error {
		return i.each(func(t T) (bool, error) {
			ok, err := fn(t)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			return yield(t), nil
		})
	}
}

// DropWhileX discards elements until fn rejects one. An error from fn aborts the
// stream, as does an abort from upstream. Lazy.
func (i IterX[T]) DropWhileX(fn func(T) (bool, error)) IterX[T] {
	return func(yield func(T) bool) error {
		dropping := true
		return i.each(func(t T) (bool, error) {
			if dropping {
				ok, err := fn(t)
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
				dropping = false
			}
			return yield(t), nil
		})
	}
}

// Take yields at most the first n pairs and then stops pulling, preserving any abort
// reached before that. Lazy.
func (i Iter2X[K, V]) Take(n int) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		if n <= 0 {
			return nil
		}
		var taken int
		return i.each(func(k K, v V) (bool, error) {
			if !yield(k, v) {
				return false, nil
			}
			taken++
			return taken < n, nil
		})
	}
}

// Drop discards the first n pairs and yields the rest, preserving any abort. Lazy.
func (i Iter2X[K, V]) Drop(n int) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		var dropped int
		return i.each(func(k K, v V) (bool, error) {
			if dropped < n {
				dropped++
				return true, nil
			}
			return yield(k, v), nil
		})
	}
}

// TakeWhile yields pairs until fn rejects one, preserving any abort. Lazy.
func (i Iter2X[K, V]) TakeWhile(fn func(K, V) bool) Iter2X[K, V] {
	return i.TakeWhileX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// DropWhile discards pairs until fn rejects one, preserving any abort. Lazy.
func (i Iter2X[K, V]) DropWhile(fn func(K, V) bool) Iter2X[K, V] {
	return i.DropWhileX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// TakeWhileX yields pairs until fn rejects one. An error from fn aborts the stream,
// as does an abort from upstream reached before the rejection. Lazy.
func (i Iter2X[K, V]) TakeWhileX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		return i.each(func(k K, v V) (bool, error) {
			ok, err := fn(k, v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			return yield(k, v), nil
		})
	}
}

// DropWhileX discards pairs until fn rejects one. An error from fn aborts the
// stream, as does an abort from upstream. Lazy.
func (i Iter2X[K, V]) DropWhileX(fn func(K, V) (bool, error)) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		dropping := true
		return i.each(func(k K, v V) (bool, error) {
			if dropping {
				ok, err := fn(k, v)
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
				dropping = false
			}
			return yield(k, v), nil
		})
	}
}
