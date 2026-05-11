package slicest

func Reduce[T any, S ~[]T, U any](s S, fn func(T, U) U) U {
	return ReduceI(s, func(_ int, t T, u U) U {
		return fn(t, u)
	})
}

func ReduceI[T any, S ~[]T, U any](s S, fn func(int, T, U) U) U {
	result, _ := ReduceXI(s, func(i int, t T, u U) (U, error) {
		return fn(i, t, u), nil
	})
	return result
}

func ReduceX[T any, S ~[]T, U any](s S, fn func(T, U) (U, error)) (U, error) {
	return ReduceXI(s, func(_ int, t T, u U) (U, error) {
		return fn(t, u)
	})
}

func ReduceD[T any, S ~[]T, U any](s S, init U, fn func(T, U) U) U {
	result, _ := ReduceXD(s, init, func(t T, u U) (U, error) {
		return fn(t, u), nil
	})
	return result
}

func ReduceXD[T any, S ~[]T, U any](s S, init U, fn func(T, U) (U, error)) (U, error) {
	return ReduceXDI(s, init, func(_ int, t T, u U) (U, error) {
		return fn(t, u)
	})
}

func ReduceXI[T any, S ~[]T, U any](s S, fn func(int, T, U) (U, error)) (U, error) {
	var zero U
	return ReduceXDI(s, zero, fn)
}

func ReduceDI[T any, S ~[]T, U any](s S, init U, fn func(int, T, U) U) U {
	result, _ := ReduceXDI(s, init, func(i int, t T, u U) (U, error) {
		return fn(i, t, u), nil
	})
	return result
}

func ReduceXDI[T any, S ~[]T, U any](s S, init U, fn func(int, T, U) (U, error)) (U, error) {
	var zero U
	for i, t := range s {
		var err error
		init, err = fn(i, t, init)
		if err != nil {
			return zero, err
		}
	}
	return init, nil
}
