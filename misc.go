package slicest

func Intersection[T comparable, S ~[]T](s1, s2 S) []T {
	m2 := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		m2[v] = struct{}{}
	}

	var result []T
	for _, t1 := range s1 {
		if _, ok := m2[t1]; ok {
			result = append(result, t1)
			delete(m2, t1)
		}
	}
	return result
}

func Unique[T comparable, S ~[]T](s1, s2 S) []T {
	m1 := make(map[T]struct{}, len(s1))
	m2 := make(map[T]struct{}, len(s2))

	for _, v := range s1 {
		m1[v] = struct{}{}
	}
	for _, v := range s2 {
		m2[v] = struct{}{}
	}

	var result []T
	for t1 := range m1 {
		if _, ok := m2[t1]; !ok {
			result = append(result, t1)
			m2[t1] = struct{}{}
		}
	}
	for t2 := range m2 {
		if _, ok := m1[t2]; !ok {
			result = append(result, t2)
		}
	}
	return result
}

func Deduplicate[T comparable, S ~[]T](s S) []T {
	return DeduplicateFunc(s, func(t T) T { return t })
}

func DeduplicateFunc[T any, K comparable, S ~[]T](s S, fn func(t T) K) []T {
	seen := make(map[K]struct{}, len(s))

	var result S
	for _, t := range s {
		key := fn(t)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, t)
		}
	}
	return result
}
