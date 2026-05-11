package slicest

func IndexValue[T comparable, S ~[]T](s S, value T) int {
	return Index(s, func(t T) bool {
		return t == value
	})
}

func Index[T any, S ~[]T](s S, fn func(T) bool) int {
	return IndexI(s, func(_ int, t T) bool {
		return fn(t)
	})
}

func IndexI[T any, S ~[]T](s S, fn func(int, T) bool) int {
	result, _ := IndexXI(s, func(i int, t T) (bool, error) {
		return fn(i, t), nil
	})
	return result
}

func IndexX[T any, S ~[]T](s S, fn func(T) (bool, error)) (int, error) {
	return IndexXI(s, func(_ int, t T) (bool, error) {
		return fn(t)
	})
}

func IndexXI[T any, S ~[]T](s S, fn func(int, T) (bool, error)) (int, error) {
	for i, t := range s {
		if ok, err := fn(i, t); err != nil {
			return -1, err
		} else if ok {
			return i, nil
		}
	}
	return -1, nil
}
