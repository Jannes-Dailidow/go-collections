package collections

// each runs fn over the stream and merges a callback error with a stream abort, so
// that every fallible operation is spared the out-of-band capture. Returning a
// non-nil error stops iteration; the accompanying bool is ignored.
//
// Every operation on [IterX] is built on this.
func (i IterX[T]) each(fn func(T) (bool, error)) error {
	var ferr error
	err := i(func(t T) bool {
		ok, e := fn(t)
		if e != nil {
			ferr = e
			return false
		}
		return ok
	})
	if ferr != nil {
		return ferr
	}
	return err
}

// each is the [Iter2X] counterpart of [IterX.each].
func (i Iter2X[K, V]) each(fn func(K, V) (bool, error)) error {
	var ferr error
	err := i(func(k K, v V) bool {
		ok, e := fn(k, v)
		if e != nil {
			ferr = e
			return false
		}
		return ok
	})
	if ferr != nil {
		return ferr
	}
	return err
}
