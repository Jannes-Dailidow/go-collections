package slicest

func Map[T, U any, S ~[]T](s S, fn func(T) U) []U {
	result, _ := MapXI(s, func(_ int, t T) (U, error) {
		return fn(t), nil
	})
	return result
}

func MapX[T, U any, S ~[]T](s S, fn func(T) (U, error)) ([]U, error) {
	return MapXI(s, func(_ int, t T) (U, error) {
		return fn(t)
	})
}

func MapI[T, U any, S ~[]T](s S, fn func(int, T) U) []U {
	result, _ := MapXI(s, func(i int, t T) (U, error) {
		return fn(i, t), nil
	})
	return result
}

func MapXI[T, U any, S ~[]T](s S, fn func(int, T) (U, error)) ([]U, error) {
	result := make([]U, len(s))
	for i, t := range s {
		out, err := fn(i, t)
		if err != nil {
			return nil, err
		}
		result[i] = out
	}
	return result, nil
}
