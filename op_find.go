package collections

// Find returns the first element fn accepts, or nil. Terminal, and short-circuits.
//
// The pointer is to a copy, not into the source. The v0.1 Find returned &s[i], a
// pointer into the slice it was given; an iterator has no addressable source to point
// into, so writing through this pointer changes nothing.
func (i Iter[T]) Find(fn func(T) bool) *T {
	result, _ := i.X().FindX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// FindX returns the first element fn accepts, or nil, and aborts on the first error.
// Terminal, and short-circuits.
func (i Iter[T]) FindX(fn func(T) (bool, error)) (*T, error) {
	return i.X().FindX(fn)
}

// FindLast returns the last element fn accepts, or nil. Terminal, and drains the
// stream: there is no way to know an element is the last one without reaching the
// end.
func (i Iter[T]) FindLast(fn func(T) bool) *T {
	result, _ := i.X().FindLastX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// FindLastX returns the last element fn accepts, or nil, and aborts on the first
// error. Terminal, and drains the stream.
func (i Iter[T]) FindLastX(fn func(T) (bool, error)) (*T, error) {
	return i.X().FindLastX(fn)
}

// First returns the first element, or nil on an empty stream. Terminal, and
// short-circuits.
func (i Iter[T]) First() *T {
	result, _ := i.X().FindX(func(T) (bool, error) {
		return true, nil
	})
	return result
}

// Last returns the last element, or nil on an empty stream. Terminal, and drains the
// stream.
func (i Iter[T]) Last() *T {
	result, _ := i.X().FindLastX(func(T) (bool, error) {
		return true, nil
	})
	return result
}

// Find returns the first pair fn accepts, or nil. Terminal, and short-circuits.
func (i Iter2[K, V]) Find(fn func(K, V) bool) *Entry[K, V] {
	result, _ := i.X().FindX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
	return result
}

// FindX returns the first pair fn accepts, or nil, and aborts on the first error.
// Terminal, and short-circuits.
func (i Iter2[K, V]) FindX(fn func(K, V) (bool, error)) (*Entry[K, V], error) {
	return i.X().FindX(fn)
}

// FindLast returns the last pair fn accepts, or nil. Terminal, and drains the stream.
func (i Iter2[K, V]) FindLast(fn func(K, V) bool) *Entry[K, V] {
	result, _ := i.X().FindLastX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
	return result
}

// FindLastX returns the last pair fn accepts, or nil, and aborts on the first error.
// Terminal, and drains the stream.
func (i Iter2[K, V]) FindLastX(fn func(K, V) (bool, error)) (*Entry[K, V], error) {
	return i.X().FindLastX(fn)
}

// First returns the first pair, or nil on an empty stream. Terminal, and
// short-circuits.
func (i Iter2[K, V]) First() *Entry[K, V] {
	result, _ := i.X().FindX(func(K, V) (bool, error) {
		return true, nil
	})
	return result
}

// Last returns the last pair, or nil on an empty stream. Terminal, and drains the
// stream.
func (i Iter2[K, V]) Last() *Entry[K, V] {
	result, _ := i.X().FindLastX(func(K, V) (bool, error) {
		return true, nil
	})
	return result
}

// Find returns the first element fn accepts, or nil, and any abort reached before it.
// Terminal, and short-circuits.
func (i IterX[T]) Find(fn func(T) bool) (*T, error) {
	return i.FindX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// FindX returns the first element fn accepts, or nil. An error from fn aborts, as
// does an abort from upstream reached before a match. Terminal, and short-circuits.
func (i IterX[T]) FindX(fn func(T) (bool, error)) (*T, error) {
	var result *T
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if ok {
			result = &t
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindLast returns the last element fn accepts, or nil. Terminal, and drains the
// stream, so an abort anywhere in it surfaces.
func (i IterX[T]) FindLast(fn func(T) bool) (*T, error) {
	return i.FindLastX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// FindLastX returns the last element fn accepts, or nil. Terminal, and drains the
// stream.
func (i IterX[T]) FindLastX(fn func(T) (bool, error)) (*T, error) {
	var result *T
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if ok {
			result = &t
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// First returns the first element, or nil on an empty stream. Terminal, and
// short-circuits.
func (i IterX[T]) First() (*T, error) {
	return i.FindX(func(T) (bool, error) {
		return true, nil
	})
}

// Last returns the last element, or nil on an empty stream. Terminal, and drains the
// stream.
func (i IterX[T]) Last() (*T, error) {
	return i.FindLastX(func(T) (bool, error) {
		return true, nil
	})
}

// Find returns the first pair fn accepts, or nil, and any abort reached before it.
// Terminal, and short-circuits.
func (i Iter2X[K, V]) Find(fn func(K, V) bool) (*Entry[K, V], error) {
	return i.FindX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// FindX returns the first pair fn accepts, or nil. An error from fn aborts, as does
// an abort from upstream reached before a match. Terminal, and short-circuits.
func (i Iter2X[K, V]) FindX(fn func(K, V) (bool, error)) (*Entry[K, V], error) {
	var result *Entry[K, V]
	err := i.each(func(k K, v V) (bool, error) {
		ok, err := fn(k, v)
		if err != nil {
			return false, err
		}
		if ok {
			result = &Entry[K, V]{Key: k, Value: v}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindLast returns the last pair fn accepts, or nil. Terminal, and drains the stream.
func (i Iter2X[K, V]) FindLast(fn func(K, V) bool) (*Entry[K, V], error) {
	return i.FindLastX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// FindLastX returns the last pair fn accepts, or nil. Terminal, and drains the
// stream.
func (i Iter2X[K, V]) FindLastX(fn func(K, V) (bool, error)) (*Entry[K, V], error) {
	var result *Entry[K, V]
	err := i.each(func(k K, v V) (bool, error) {
		ok, err := fn(k, v)
		if err != nil {
			return false, err
		}
		if ok {
			result = &Entry[K, V]{Key: k, Value: v}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// First returns the first pair, or nil on an empty stream. Terminal, and
// short-circuits.
func (i Iter2X[K, V]) First() (*Entry[K, V], error) {
	return i.FindX(func(K, V) (bool, error) {
		return true, nil
	})
}

// Last returns the last pair, or nil on an empty stream. Terminal, and drains the
// stream.
func (i Iter2X[K, V]) Last() (*Entry[K, V], error) {
	return i.FindLastX(func(K, V) (bool, error) {
		return true, nil
	})
}
