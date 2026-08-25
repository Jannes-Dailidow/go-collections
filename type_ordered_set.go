package collections

import (
	"iter"
	"slices"
)

type OrderedSet[T comparable] struct {
	order []T
	index map[T]int
}

func CollectOrderedSet[T comparable](i Iter[T]) *OrderedSet[T] {
	result, _ := CollectOrderedSetX(i.X())
	return result
}

func CollectOrderedSetX[T comparable](i IterX[T]) (*OrderedSet[T], error) {
	result := &OrderedSet[T]{index: make(map[T]int)}
	err := i(func(t T) bool {
		if _, ok := result.index[t]; ok {
			return true
		}
		result.index[t] = len(result.order)
		result.order = append(result.order, t)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OrderedSet[T]) Values() Iter[T] {
	return func(yield func(T) bool) {
		if s == nil {
			return
		}
		for _, t := range s.order {
			if !yield(t) {
				return
			}
		}
	}
}

func (s *OrderedSet[T]) All() Iter2[int, T] {
	return func(yield func(int, T) bool) {
		if s == nil {
			return
		}
		for i, t := range s.order {
			if !yield(i, t) {
				return
			}
		}
	}
}

func (s *OrderedSet[T]) Slice() Slice[T] {
	if s == nil {
		return nil
	}
	return Slice[T](slices.Clone(s.order))
}

func (s *OrderedSet[T]) Native() []T {
	return []T(s.Slice())
}

func (s *OrderedSet[T]) Set() Set[T] {
	if s == nil {
		return nil
	}
	result := make(Set[T], len(s.order))
	for _, t := range s.order {
		result[t] = struct{}{}
	}
	return result
}

func (s *OrderedSet[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](s.Values())
}

// Backward iterates from the last element to the first, pairing each with its index.
// Lazy, and safe on a nil receiver.
func (s *OrderedSet[T]) Backward() Iter2[int, T] {
	return func(yield func(int, T) bool) {
		if s == nil {
			return
		}
		for i := len(s.order) - 1; i >= 0; i-- {
			if !yield(i, s.order[i]) {
				return
			}
		}
	}
}
