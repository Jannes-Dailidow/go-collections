package slicest

func Find[T any, S ~[]T](s S, fn func(T) bool) *T {
	return FindI(s, func(_ int, t T) bool {
		return fn(t)
	})
}

func FindI[T any, S ~[]T](s S, fn func(int, T) bool) *T {
	result, _ := FindXI(s, func(i int, t T) (bool, error) {
		return fn(i, t), nil
	})
	return result
}

func FindX[T any, S ~[]T](s S, fn func(T) (bool, error)) (*T, error) {
	return FindXI(s, func(_ int, t T) (bool, error) {
		return fn(t)
	})
}

func FindXI[T any, S ~[]T](s S, fn func(int, T) (bool, error)) (*T, error) {
	for i, t := range s {
		if ok, err := fn(i, t); err != nil {
			return nil, err
		} else if ok {
			return &s[i], nil
		}
	}
	return nil, nil
}
