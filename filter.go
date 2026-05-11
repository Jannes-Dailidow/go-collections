package slicest

func Filter[T any, S ~[]T](s S, fn func(T) bool) S {
	return FilterI(s, func(_ int, t T) bool {
		return fn(t)
	})
}

func FilterI[T any, S ~[]T](s S, fn func(int, T) bool) S {
	result, _ := FilterXI(s, func(i int, t T) (bool, error) {
		return fn(i, t), nil
	})
	return result
}

func FilterX[T any, S ~[]T](s S, fn func(T) (bool, error)) (S, error) {
	return FilterXI(s, func(_ int, t T) (bool, error) {
		return fn(t)
	})
}

func FilterXI[T any, S ~[]T](s S, fn func(int, T) (bool, error)) (S, error) {
	var result S
	for i, t := range s {
		ok, err := fn(i, t)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, t)
		}
	}
	return result, nil
}
