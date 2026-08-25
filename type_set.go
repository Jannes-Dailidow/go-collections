package collections

import "iter"

type Set[T comparable] map[T]struct{}

func CollectSet[T comparable](i Iter[T]) Set[T] {
	result, _ := CollectSetX(i.X())
	return result
}

func CollectSetX[T comparable](i IterX[T]) (Set[T], error) {
	result := make(Set[T])
	err := i(func(t T) bool {
		result[t] = struct{}{}
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Set[T]) Native() map[T]struct{} {
	return map[T]struct{}(s)
}

func (s Set[T]) Values() Iter[T] {
	return func(yield func(T) bool) {
		for t := range s {
			if !yield(t) {
				return
			}
		}
	}
}

func (s Set[T]) Slice() Slice[T] {
	result := make(Slice[T], 0, len(s))
	for t := range s {
		result = append(result, t)
	}
	return result
}

func (s Set[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](s.Values())
}

// OrderedSet copies the elements into an [OrderedSet]. The order is whatever the set
// iterated in, which is arbitrary; to get a meaningful one, sort afterwards or build
// with [CollectOrderedSet] from a source that already has an order.
func (s Set[T]) OrderedSet() *OrderedSet[T] {
	return CollectOrderedSet(s.Values())
}
