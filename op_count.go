package collections

// Count returns the number of elements. Terminal, and drains the stream; a
// collection's Len is the free version.
func (i Iter[T]) Count() int {
	result, _ := i.X().CountByX(func(T) (bool, error) {
		return true, nil
	})
	return result
}

// CountBy returns the number of elements fn accepts. Terminal, and drains the stream.
func (i Iter[T]) CountBy(fn func(T) bool) int {
	result, _ := i.X().CountByX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// CountByX returns the number of elements fn accepts, and aborts on the first error.
// Terminal.
func (i Iter[T]) CountByX(fn func(T) (bool, error)) (int, error) {
	return i.X().CountByX(fn)
}

// IsEmpty reports whether the stream yields nothing. Terminal, but pulls one element
// at most.
func (i Iter[T]) IsEmpty() bool {
	result, _ := i.X().IsEmpty()
	return result
}

// Count returns the number of pairs. Terminal, and drains the stream.
func (i Iter2[K, V]) Count() int {
	result, _ := i.X().CountByX(func(K, V) (bool, error) {
		return true, nil
	})
	return result
}

// CountBy returns the number of pairs fn accepts. Terminal, and drains the stream.
func (i Iter2[K, V]) CountBy(fn func(K, V) bool) int {
	result, _ := i.X().CountByX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
	return result
}

// CountByX returns the number of pairs fn accepts, and aborts on the first error.
// Terminal.
func (i Iter2[K, V]) CountByX(fn func(K, V) (bool, error)) (int, error) {
	return i.X().CountByX(fn)
}

// IsEmpty reports whether the stream yields nothing. Terminal, but pulls one pair at
// most.
func (i Iter2[K, V]) IsEmpty() bool {
	result, _ := i.X().IsEmpty()
	return result
}

// Count returns the number of elements, and any abort. Terminal, and drains the
// stream.
func (i IterX[T]) Count() (int, error) {
	return i.CountByX(func(T) (bool, error) {
		return true, nil
	})
}

// CountBy returns the number of elements fn accepts, and any abort. Terminal.
func (i IterX[T]) CountBy(fn func(T) bool) (int, error) {
	return i.CountByX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// CountByX returns the number of elements fn accepts. An error from fn aborts, as
// does an abort from upstream. Terminal.
func (i IterX[T]) CountByX(fn func(T) (bool, error)) (int, error) {
	var result int
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if ok {
			result++
		}
		return true, nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}

// IsEmpty reports whether the stream yields nothing. Terminal, but pulls one element
// at most.
func (i IterX[T]) IsEmpty() (bool, error) {
	result := true
	err := i.each(func(T) (bool, error) {
		result = false
		return false, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

// Count returns the number of pairs, and any abort. Terminal, and drains the stream.
func (i Iter2X[K, V]) Count() (int, error) {
	return i.CountByX(func(K, V) (bool, error) {
		return true, nil
	})
}

// CountBy returns the number of pairs fn accepts, and any abort. Terminal.
func (i Iter2X[K, V]) CountBy(fn func(K, V) bool) (int, error) {
	return i.CountByX(func(k K, v V) (bool, error) {
		return fn(k, v), nil
	})
}

// CountByX returns the number of pairs fn accepts. An error from fn aborts, as does
// an abort from upstream. Terminal.
func (i Iter2X[K, V]) CountByX(fn func(K, V) (bool, error)) (int, error) {
	var result int
	err := i.each(func(k K, v V) (bool, error) {
		ok, err := fn(k, v)
		if err != nil {
			return false, err
		}
		if ok {
			result++
		}
		return true, nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}

// IsEmpty reports whether the stream yields nothing. Terminal, but pulls one pair at
// most.
func (i Iter2X[K, V]) IsEmpty() (bool, error) {
	result := true
	err := i.each(func(K, V) (bool, error) {
		result = false
		return false, nil
	})
	if err != nil {
		return false, err
	}
	return result, nil
}
