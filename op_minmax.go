package collections

import "cmp"

// Min returns the smallest element, or nil on an empty stream. Terminal, and drains
// the stream.
//
// It is a function rather than a method because a method cannot tighten T to
// cmp.Ordered. For the fluent form use [Iter.MinBy] with an identity key.
func Min[T cmp.Ordered](i Iter[T]) *T {
	return i.MinBy(func(t T) T {
		return t
	})
}

// MinX is the fallible counterpart of [Min].
func MinX[T cmp.Ordered](i IterX[T]) (*T, error) {
	return i.MinBy(func(t T) T {
		return t
	})
}

// Max returns the largest element, or nil on an empty stream. Terminal, and drains
// the stream.
func Max[T cmp.Ordered](i Iter[T]) *T {
	return i.MaxBy(func(t T) T {
		return t
	})
}

// MaxX is the fallible counterpart of [Max].
func MaxX[T cmp.Ordered](i IterX[T]) (*T, error) {
	return i.MaxBy(func(t T) T {
		return t
	})
}

// MinBy returns the element with the smallest key, or nil on an empty stream. The
// first element wins a tie. Terminal, and drains the stream.
//
// The key function is infallible. Where a key can fail, use [Iter.KeyByX] to carry
// the failure into a pair stream and take the minimum of that.
func (i Iter[T]) MinBy[K cmp.Ordered](fn func(T) K) *T {
	result, _ := i.X().MinBy(fn)
	return result
}

// MaxBy returns the element with the largest key, or nil on an empty stream. The
// first element wins a tie. Terminal, and drains the stream.
func (i Iter[T]) MaxBy[K cmp.Ordered](fn func(T) K) *T {
	result, _ := i.X().MaxBy(fn)
	return result
}

// MinBy returns the pair with the smallest key, or nil on an empty stream. Terminal,
// and drains the stream.
func (i Iter2[K, V]) MinBy[K2 cmp.Ordered](fn func(K, V) K2) *Entry[K, V] {
	result, _ := i.X().MinBy(fn)
	return result
}

// MaxBy returns the pair with the largest key, or nil on an empty stream. Terminal,
// and drains the stream.
func (i Iter2[K, V]) MaxBy[K2 cmp.Ordered](fn func(K, V) K2) *Entry[K, V] {
	result, _ := i.X().MaxBy(fn)
	return result
}

// MinBy returns the element with the smallest key, or nil on an empty stream, and any
// abort. Terminal, and drains the stream.
func (i IterX[T]) MinBy[K cmp.Ordered](fn func(T) K) (*T, error) {
	var result *T
	var best K
	err := i.each(func(t T) (bool, error) {
		k := fn(t)
		if result == nil || k < best {
			best = k
			result = &t
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MaxBy returns the element with the largest key, or nil on an empty stream, and any
// abort. Terminal, and drains the stream.
func (i IterX[T]) MaxBy[K cmp.Ordered](fn func(T) K) (*T, error) {
	var result *T
	var best K
	err := i.each(func(t T) (bool, error) {
		k := fn(t)
		if result == nil || k > best {
			best = k
			result = &t
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MinBy returns the pair with the smallest key, or nil on an empty stream, and any
// abort. Terminal, and drains the stream.
func (i Iter2X[K, V]) MinBy[K2 cmp.Ordered](fn func(K, V) K2) (*Entry[K, V], error) {
	var result *Entry[K, V]
	var best K2
	err := i.each(func(k K, v V) (bool, error) {
		k2 := fn(k, v)
		if result == nil || k2 < best {
			best = k2
			result = &Entry[K, V]{Key: k, Value: v}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MaxBy returns the pair with the largest key, or nil on an empty stream, and any
// abort. Terminal, and drains the stream.
func (i Iter2X[K, V]) MaxBy[K2 cmp.Ordered](fn func(K, V) K2) (*Entry[K, V], error) {
	var result *Entry[K, V]
	var best K2
	err := i.each(func(k K, v V) (bool, error) {
		k2 := fn(k, v)
		if result == nil || k2 > best {
			best = k2
			result = &Entry[K, V]{Key: k, Value: v}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
