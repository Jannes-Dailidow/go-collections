package collections

import "iter"

type OrderedMap[K comparable, V any] struct {
	order []Entry[K, V]
	index map[K]int
}

func CollectOrderedMap[K comparable, V any](i Iter2[K, V]) *OrderedMap[K, V] {
	result, _ := CollectOrderedMapX(i.X())
	return result
}

func CollectOrderedMapX[K comparable, V any](i Iter2X[K, V]) (*OrderedMap[K, V], error) {
	result := &OrderedMap[K, V]{index: make(map[K]int)}
	err := i(func(k K, v V) bool {
		if pos, ok := result.index[k]; ok {
			result.order[pos].Value = v
			return true
		}
		result.index[k] = len(result.order)
		result.order = append(result.order, Entry[K, V]{Key: k, Value: v})
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *OrderedMap[K, V]) All() Iter2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for _, e := range m.order {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Keys() Iter[K] {
	return func(yield func(K) bool) {
		if m == nil {
			return
		}
		for _, e := range m.order {
			if !yield(e.Key) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Values() Iter[V] {
	return func(yield func(V) bool) {
		if m == nil {
			return
		}
		for _, e := range m.order {
			if !yield(e.Value) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) Map() Map[K, V] {
	if m == nil {
		return nil
	}
	result := make(Map[K, V], len(m.order))
	for _, e := range m.order {
		result[e.Key] = e.Value
	}
	return result
}

func (m *OrderedMap[K, V]) Native() map[K]V {
	return map[K]V(m.Map())
}

func (m *OrderedMap[K, V]) Seq2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](m.All())
}

// Backward iterates the pairs from the last insertion to the first. Lazy, and safe
// on a nil receiver.
func (m *OrderedMap[K, V]) Backward() Iter2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for i := len(m.order) - 1; i >= 0; i-- {
			if !yield(m.order[i].Key, m.order[i].Value) {
				return
			}
		}
	}
}
