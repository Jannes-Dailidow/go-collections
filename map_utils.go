package slicest

func ToMap[T any, K comparable, V any, S ~[]T](s S, fn func(T) (K, V)) map[K]V {
	result := make(map[K]V, len(s))
	for _, t := range s {
		k, v := fn(t)
		result[k] = v
	}
	return result
}

func FromMap[K comparable, V any, T any, M ~map[K]V](m M, fn func(K, V) T) []T {
	result := make([]T, len(m))
	var i int
	for k, v := range m {
		t := fn(k, v)
		result[i] = t
		i++
	}
	return result
}

func MapKeys[K comparable, V any, M ~map[K]V](m M) []K {
	result := make([]K, len(m))
	var i int
	for k := range m {
		result[i] = k
		i++
	}
	return result
}

func MapValues[K comparable, V any, M ~map[K]V](m M) []V {
	result := make([]V, len(m))
	var i int
	for _, v := range m {
		result[i] = v
		i++
	}
	return result
}
