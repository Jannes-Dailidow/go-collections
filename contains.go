package slicest

func ContainsValue[T comparable, S ~[]T](s S, value T) bool {
	return Contains(s, func(t T) bool {
		return t == value
	})
}

func Contains[T any, S ~[]T](s S, fn func(T) bool) bool {
	return ContainsI(s, func(_ int, t T) bool {
		return fn(t)
	})
}

func ContainsI[T any, S ~[]T](s S, fn func(int, T) bool) bool {
	result, _ := ContainsXI(s, func(i int, t T) (bool, error) {
		return fn(i, t), nil
	})
	return result
}

func ContainsX[T any, S ~[]T](s S, fn func(T) (bool, error)) (bool, error) {
	return ContainsXI(s, func(_ int, t T) (bool, error) {
		return fn(t)
	})
}

func ContainsXI[T any, S ~[]T](s S, fn func(int, T) (bool, error)) (bool, error) {
	for i, t := range s {
		result, err := fn(i, t)
		if result || err != nil {
			return result, err
		}
	}
	return false, nil
}
