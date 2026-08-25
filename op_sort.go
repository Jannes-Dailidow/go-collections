package collections

import (
	"cmp"
	"slices"
)

// SortBy yields the elements ordered by fn. Stable, so equal keys keep their
// original order.
//
// This buffers the whole stream on the first pull, which is as lazy as sorting can
// be: nothing can be yielded until everything has been seen. It is the only family
// here that is not incremental, and it will not terminate on an endless stream.
func (i Iter[T]) SortBy[K cmp.Ordered](fn func(T) K) Iter[T] {
	return i.SortFunc(func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})
}

// SortFunc yields the elements ordered by compare, for keys that are not
// cmp.Ordered. Stable. Buffers the whole stream.
func (i Iter[T]) SortFunc(compare func(a, b T) int) Iter[T] {
	return func(yield func(T) bool) {
		result := i.Slice()
		slices.SortStableFunc(result, compare)
		for _, t := range result {
			if !yield(t) {
				return
			}
		}
	}
}

// Reverse yields the elements back to front. Buffers the whole stream; on a [Slice]
// prefer [Slice.Backward], which does not.
func (i Iter[T]) Reverse() Iter[T] {
	return func(yield func(T) bool) {
		result := i.Slice()
		for n := len(result) - 1; n >= 0; n-- {
			if !yield(result[n]) {
				return
			}
		}
	}
}

// SortBy yields the pairs ordered by fn. Stable. Buffers the whole stream.
func (i Iter2[K, V]) SortBy[K2 cmp.Ordered](fn func(K, V) K2) Iter2[K, V] {
	return i.SortFunc(func(a, b Entry[K, V]) int {
		return cmp.Compare(fn(a.Key, a.Value), fn(b.Key, b.Value))
	})
}

// SortFunc yields the pairs ordered by compare, which receives them as [Entry]
// values. Stable. Buffers the whole stream.
func (i Iter2[K, V]) SortFunc(compare func(a, b Entry[K, V]) int) Iter2[K, V] {
	return func(yield func(K, V) bool) {
		var result []Entry[K, V]
		for k, v := range i {
			result = append(result, Entry[K, V]{Key: k, Value: v})
		}
		slices.SortStableFunc(result, compare)
		for _, e := range result {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
}

// Reverse yields the pairs back to front. Buffers the whole stream.
func (i Iter2[K, V]) Reverse() Iter2[K, V] {
	return func(yield func(K, V) bool) {
		var result []Entry[K, V]
		for k, v := range i {
			result = append(result, Entry[K, V]{Key: k, Value: v})
		}
		for n := len(result) - 1; n >= 0; n-- {
			if !yield(result[n].Key, result[n].Value) {
				return
			}
		}
	}
}

// SortBy yields the elements ordered by fn, preserving any abort. Stable. Buffers the
// whole stream, so an abort anywhere in it surfaces before anything is yielded.
func (i IterX[T]) SortBy[K cmp.Ordered](fn func(T) K) IterX[T] {
	return i.SortFunc(func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})
}

// SortFunc yields the elements ordered by compare, preserving any abort. Stable.
func (i IterX[T]) SortFunc(compare func(a, b T) int) IterX[T] {
	return func(yield func(T) bool) error {
		result, err := i.Slice()
		if err != nil {
			return err
		}
		slices.SortStableFunc(result, compare)
		for _, t := range result {
			if !yield(t) {
				return nil
			}
		}
		return nil
	}
}

// Reverse yields the elements back to front, preserving any abort. Buffers the whole
// stream.
func (i IterX[T]) Reverse() IterX[T] {
	return func(yield func(T) bool) error {
		result, err := i.Slice()
		if err != nil {
			return err
		}
		for n := len(result) - 1; n >= 0; n-- {
			if !yield(result[n]) {
				return nil
			}
		}
		return nil
	}
}

// SortBy yields the pairs ordered by fn, preserving any abort. Stable.
func (i Iter2X[K, V]) SortBy[K2 cmp.Ordered](fn func(K, V) K2) Iter2X[K, V] {
	return i.SortFunc(func(a, b Entry[K, V]) int {
		return cmp.Compare(fn(a.Key, a.Value), fn(b.Key, b.Value))
	})
}

// SortFunc yields the pairs ordered by compare, preserving any abort. Stable.
func (i Iter2X[K, V]) SortFunc(compare func(a, b Entry[K, V]) int) Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		var result []Entry[K, V]
		err := i.each(func(k K, v V) (bool, error) {
			result = append(result, Entry[K, V]{Key: k, Value: v})
			return true, nil
		})
		if err != nil {
			return err
		}
		slices.SortStableFunc(result, compare)
		for _, e := range result {
			if !yield(e.Key, e.Value) {
				return nil
			}
		}
		return nil
	}
}

// Reverse yields the pairs back to front, preserving any abort. Buffers the whole
// stream.
func (i Iter2X[K, V]) Reverse() Iter2X[K, V] {
	return func(yield func(K, V) bool) error {
		var result []Entry[K, V]
		err := i.each(func(k K, v V) (bool, error) {
			result = append(result, Entry[K, V]{Key: k, Value: v})
			return true, nil
		})
		if err != nil {
			return err
		}
		for n := len(result) - 1; n >= 0; n-- {
			if !yield(result[n].Key, result[n].Value) {
				return nil
			}
		}
		return nil
	}
}

// SortBy sorts the slice in place and returns it, so it can be chained. Stable.
//
// In place means the caller's backing array, matching slices.Sort rather than the
// iterator versions above, which buffer a copy. Clone first to keep the original.
func (s Slice[T]) SortBy[K cmp.Ordered](fn func(T) K) Slice[T] {
	slices.SortStableFunc(s, func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})
	return s
}

// SortFunc sorts the slice in place by compare and returns it. Stable.
func (s Slice[T]) SortFunc(compare func(a, b T) int) Slice[T] {
	slices.SortStableFunc(s, compare)
	return s
}

// Reverse reverses the slice in place and returns it. For a view that does not
// mutate, use [Slice.Backward].
func (s Slice[T]) Reverse() Slice[T] {
	slices.Reverse(s)
	return s
}

// SortBy rewrites the insertion order by fn and returns the map, so it can be
// chained. Stable, and safe on a nil receiver.
//
// After this the order is no longer insertion order. Everything positional --
// [OrderedMap.At], [OrderedMap.IndexOf], [OrderedMap.Backward] -- follows the new one.
func (m *OrderedMap[K, V]) SortBy[K2 cmp.Ordered](fn func(K, V) K2) *OrderedMap[K, V] {
	return m.SortFunc(func(a, b Entry[K, V]) int {
		return cmp.Compare(fn(a.Key, a.Value), fn(b.Key, b.Value))
	})
}

// SortFunc rewrites the insertion order by compare and returns the map. Stable, and
// safe on a nil receiver.
func (m *OrderedMap[K, V]) SortFunc(compare func(a, b Entry[K, V]) int) *OrderedMap[K, V] {
	if m == nil {
		return nil
	}
	slices.SortStableFunc(m.order, compare)
	for n, e := range m.order {
		m.index[e.Key] = n
	}
	return m
}

// Reverse reverses the order and returns the map. Safe on a nil receiver.
func (m *OrderedMap[K, V]) Reverse() *OrderedMap[K, V] {
	if m == nil {
		return nil
	}
	slices.Reverse(m.order)
	for n, e := range m.order {
		m.index[e.Key] = n
	}
	return m
}

// SortBy rewrites the order by fn and returns the set, so it can be chained. Stable,
// and safe on a nil receiver.
//
// After this the order is no longer first-seen. [OrderedSet.At] and
// [OrderedSet.IndexOf] follow the new one.
func (s *OrderedSet[T]) SortBy[K cmp.Ordered](fn func(T) K) *OrderedSet[T] {
	return s.SortFunc(func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})
}

// SortFunc rewrites the order by compare and returns the set. Stable, and safe on a
// nil receiver.
func (s *OrderedSet[T]) SortFunc(compare func(a, b T) int) *OrderedSet[T] {
	if s == nil {
		return nil
	}
	slices.SortStableFunc(s.order, compare)
	for n, value := range s.order {
		s.index[value] = n
	}
	return s
}

// Reverse reverses the order and returns the set. Safe on a nil receiver.
func (s *OrderedSet[T]) Reverse() *OrderedSet[T] {
	if s == nil {
		return nil
	}
	slices.Reverse(s.order)
	for n, value := range s.order {
		s.index[value] = n
	}
	return s
}
