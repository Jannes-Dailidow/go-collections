package collections

// IndexValue returns the position of the first element equal to value, or -1.
// Terminal, and short-circuits.
//
// It is a function rather than a method because a method cannot tighten T to
// comparable.
func IndexValue[T comparable](i Iter[T], value T) int {
	return i.Index(func(t T) bool {
		return t == value
	})
}

// IndexValueX is the fallible counterpart of [IndexValue].
func IndexValueX[T comparable](i IterX[T], value T) (int, error) {
	return i.Index(func(t T) bool {
		return t == value
	})
}

// Index returns the position of the first element fn accepts, or -1. Terminal, and
// short-circuits.
//
// Position is counted in the stream, not in whatever produced it, so a Filter ahead
// of an Index renumbers what survives. There is no pair form: a Map and a Set have
// no order for a position to mean anything in, and the ordered types reach this
// through Values.
func (i Iter[T]) Index(fn func(T) bool) int {
	result, _ := i.X().IndexX(func(t T) (bool, error) {
		return fn(t), nil
	})
	return result
}

// IndexX returns the position of the first element fn accepts, or -1, and aborts on
// the first error. Terminal, and short-circuits.
func (i Iter[T]) IndexX(fn func(T) (bool, error)) (int, error) {
	return i.X().IndexX(fn)
}

// Index returns the position of the first element fn accepts, or -1, and any abort
// reached before it. Terminal, and short-circuits.
func (i IterX[T]) Index(fn func(T) bool) (int, error) {
	return i.IndexX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// IndexX returns the position of the first element fn accepts, or -1. An error from
// fn aborts, as does an abort from upstream reached before a match. Terminal, and
// short-circuits.
func (i IterX[T]) IndexX(fn func(T) (bool, error)) (int, error) {
	result := -1
	var n int
	err := i.each(func(t T) (bool, error) {
		ok, err := fn(t)
		if err != nil {
			return false, err
		}
		if ok {
			result = n
			return false, nil
		}
		n++
		return true, nil
	})
	if err != nil {
		return -1, err
	}
	return result, nil
}
