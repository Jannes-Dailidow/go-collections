package collections

import "iter"

type Map[K comparable, V any] map[K]V

func CollectMap[K comparable, V any](i Iter2[K, V]) Map[K, V] {
	result, _ := CollectMapX(i.X())
	return result
}

func CollectMapX[K comparable, V any](i Iter2X[K, V]) (Map[K, V], error) {
	result := make(Map[K, V])
	err := i(func(k K, v V) bool {
		result[k] = v
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m Map[K, V]) Native() map[K]V {
	return map[K]V(m)
}

func (m Map[K, V]) Keys() Iter[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func (m Map[K, V]) Values() Iter[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				return
			}
		}
	}
}

func (m Map[K, V]) All() Iter2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m Map[K, V]) Seq2() iter.Seq2[K, V] {
	return iter.Seq2[K, V](m.All())
}

// OrderedMap copies the pairs into an [OrderedMap]. The order is whatever the map
// iterated in, which is arbitrary; to get a meaningful one, sort afterwards or build
// with [CollectOrderedMap] from a source that already has an order.
func (m Map[K, V]) OrderedMap() *OrderedMap[K, V] {
	return CollectOrderedMap(m.All())
}
