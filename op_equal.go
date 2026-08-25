package collections

import "iter"

// EqualBy reports whether both streams yield the same keys in the same order, and
// end together. Terminal, and short-circuits on the first difference.
//
// This is positional, so it only means something when both sides have a defined
// order. Reached from a [Map] or a [Set] through Values it compares two arbitrary
// orderings and will usually report false.
func (i Iter[T]) EqualBy[K comparable](other Iter[T], fn func(T) K) bool {
	next, stop := iter.Pull(other.Seq())
	defer stop()
	for t := range i {
		u, ok := next()
		if !ok || fn(t) != fn(u) {
			return false
		}
	}
	_, ok := next()
	return !ok
}

// EqualBy reports whether both streams yield the same pair keys in the same order,
// and end together. Terminal, and short-circuits on the first difference.
func (i Iter2[K, V]) EqualBy[K2 comparable](other Iter2[K, V], fn func(K, V) K2) bool {
	next, stop := iter.Pull2(other.Seq2())
	defer stop()
	for k, v := range i {
		k2, v2, ok := next()
		if !ok || fn(k, v) != fn(k2, v2) {
			return false
		}
	}
	_, _, ok := next()
	return !ok
}

// EqualBy reports whether both streams yield the same keys in the same order, and
// end together, plus any abort. The other side is infallible, for the same reason
// [IterX.Zip] takes an infallible side. Terminal, and short-circuits.
func (i IterX[T]) EqualBy[K comparable](other Iter[T], fn func(T) K) (bool, error) {
	next, stop := iter.Pull(other.Seq())
	defer stop()
	result := true
	err := i.each(func(t T) (bool, error) {
		u, ok := next()
		if !ok || fn(t) != fn(u) {
			result = false
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	if !result {
		return false, nil
	}
	_, ok := next()
	return !ok, nil
}

// EqualBy reports whether both streams yield the same pair keys in the same order,
// and end together, plus any abort. Terminal, and short-circuits.
func (i Iter2X[K, V]) EqualBy[K2 comparable](other Iter2[K, V], fn func(K, V) K2) (bool, error) {
	next, stop := iter.Pull2(other.Seq2())
	defer stop()
	result := true
	err := i.each(func(k K, v V) (bool, error) {
		k2, v2, ok := next()
		if !ok || fn(k, v) != fn(k2, v2) {
			result = false
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	if !result {
		return false, nil
	}
	_, _, ok := next()
	return !ok, nil
}
