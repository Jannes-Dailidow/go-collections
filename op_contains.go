package collections

// ContainsValue reports whether the stream holds value. Terminal, and
// short-circuits.
//
// It is a function rather than a method because a method cannot tighten T to
// comparable.
func ContainsValue[T comparable](i Iter[T], value T) bool {
	return i.Contains(func(t T) bool {
		return t == value
	})
}

// ContainsValueX is the fallible counterpart of [ContainsValue].
func ContainsValueX[T comparable](i IterX[T], value T) (bool, error) {
	return i.Contains(func(t T) bool {
		return t == value
	})
}

// Contains reports whether fn accepts any element. False on an empty stream.
// Terminal, and short-circuits.
//
// There is no Some: it would be a second name for this one operation.
func (i Iter[T]) Contains(fn func(T) bool) bool {
	result, _ := i.X().ContainsX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// ContainsX reports whether fn accepts any element, and aborts on the first error.
// Terminal, and short-circuits.
func (i Iter[T]) ContainsX(fn func(T) (bool, error)) (bool, error) {
	return i.X().ContainsX(fn)
}

// Every reports whether fn accepts every element. True on an empty stream, which is
// the vacuous answer rather than a special case. Terminal, and short-circuits on the
// first rejection.
func (i Iter[T]) Every(fn func(T) bool) bool {
	result, _ := i.X().EveryX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// EveryX reports whether fn accepts every element, and aborts on the first error.
// Terminal, and short-circuits on the first rejection.
func (i Iter[T]) EveryX(fn func(T) (bool, error)) (bool, error) {
	return i.X().EveryX(fn)
}

// Contains reports whether fn accepts any pair. False on an empty stream. Terminal,
// and short-circuits.
func (i Iter2[K, V]) Contains(fn func(K, V) bool) bool {
	result, _ := i.X().ContainsX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
	return result
}

// ContainsX reports whether fn accepts any pair, and aborts on the first error.
// Terminal, and short-circuits.
func (i Iter2[K, V]) ContainsX(fn func(K, V) (bool, error)) (bool, error) {
	return i.X().ContainsX(fn)
}

// Every reports whether fn accepts every pair. True on an empty stream. Terminal,
// and short-circuits on the first rejection.
func (i Iter2[K, V]) Every(fn func(K, V) bool) bool {
	result, _ := i.X().EveryX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
	return result
}

// EveryX reports whether fn accepts every pair, and aborts on the first error.
// Terminal, and short-circuits on the first rejection.
func (i Iter2[K, V]) EveryX(fn func(K, V) (bool, error)) (bool, error) {
	return i.X().EveryX(fn)
}

// Contains reports whether fn accepts any element, and any abort reached before the
// match. Terminal, and short-circuits.
func (i IterX[T]) Contains(fn func(T) bool) (bool, error) {
	return i.ContainsX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// ContainsX reports whether fn accepts any element. An error from fn aborts, as does
// an abort from upstream reached before a match. Terminal, and short-circuits.
func (i IterX[T]) ContainsX(fn func(T) (bool, error)) (bool, error) {
	var result bool
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if ok {
			result = true
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// Every reports whether fn accepts every element, and any abort reached before a
// rejection. Terminal, and short-circuits.
func (i IterX[T]) Every(fn func(T) bool) (bool, error) {
	return i.EveryX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// EveryX reports whether fn accepts every element. An error from fn aborts, as does
// an abort from upstream reached before a rejection. Terminal, and short-circuits.
func (i IterX[T]) EveryX(fn func(T) (bool, error)) (bool, error) {
	result := true
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if !ok {
			result = false
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// Contains reports whether fn accepts any pair, and any abort reached before the
// match. Terminal, and short-circuits.
func (i Iter2X[K, V]) Contains(fn func(K, V) bool) (bool, error) {
	return i.ContainsX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// ContainsX reports whether fn accepts any pair. An error from fn aborts, as does an
// abort from upstream reached before a match. Terminal, and short-circuits.
func (i Iter2X[K, V]) ContainsX(fn func(K, V) (bool, error)) (bool, error) {
	var result bool
	err := i.each(func(k K, v V) (bool, error) {
		ok, err := fn(k, v)
		if err != nil {
			return false, err
		}
		if ok {
			result = true
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// Every reports whether fn accepts every pair, and any abort reached before a
// rejection. Terminal, and short-circuits.
func (i Iter2X[K, V]) Every(fn func(K, V) bool) (bool, error) {
	return i.EveryX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// EveryX reports whether fn accepts every pair. An error from fn aborts, as does an
// abort from upstream reached before a rejection. Terminal, and short-circuits.
func (i Iter2X[K, V]) EveryX(fn func(K, V) (bool, error)) (bool, error) {
	result := true
	err := i.each(func(k K, v V) (bool, error) {
		ok, err := fn(k, v)
		if err != nil {
			return false, err
		}
		if !ok {
			result = false
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}
