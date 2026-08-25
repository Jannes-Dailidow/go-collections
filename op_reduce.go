package collections

// Reduce folds the stream into a single value, seeded with U's zero value. Terminal.
//
// The callback takes the accumulator first, matching [Iter.Scan] and the stdlib
// rather than the Reduce of v0.1, which took the element first.
func (i Iter[T]) Reduce[U any](fn func(U, T) U) U {
	var zero U
	result, _ := i.X().FoldX(zero, func(acc U, t T) (U, error) {
		return fn(acc, t), nil
	})
	return result
}

// Fold folds the stream into a single value, starting from init. It is what the root
// package calls ReduceD. Terminal.
func (i Iter[T]) Fold[U any](init U, fn func(U, T) U) U {
	result, _ := i.X().FoldX(init, func(acc U, t T) (U, error) {
		return fn(acc, t), nil
	})
	return result
}

// ReduceX folds the stream into a single value, seeded with U's zero value, and
// aborts on the first error. Terminal.
func (i Iter[T]) ReduceX[U any](fn func(U, T) (U, error)) (U, error) {
	var zero U
	return i.X().FoldX(zero, fn)
}

// FoldX folds the stream into a single value, starting from init, and aborts on the
// first error. Terminal.
func (i Iter[T]) FoldX[U any](init U, fn func(U, T) (U, error)) (U, error) {
	return i.X().FoldX(init, fn)
}

// Reduce folds the stream into a single value, seeded with U's zero value. Terminal.
func (i Iter2[K, V]) Reduce[U any](fn func(U, K, V) U) U {
	var zero U
	result, _ := i.X().FoldX(zero, func(acc U, k K, v V) (U, error) {
		return fn(acc, k, v), nil
	})
	return result
}

// Fold folds the stream into a single value, starting from init. Terminal.
func (i Iter2[K, V]) Fold[U any](init U, fn func(U, K, V) U) U {
	result, _ := i.X().FoldX(init, func(acc U, k K, v V) (U, error) {
		return fn(acc, k, v), nil
	})
	return result
}

// ReduceX folds the stream into a single value, seeded with U's zero value, and
// aborts on the first error. Terminal.
func (i Iter2[K, V]) ReduceX[U any](fn func(U, K, V) (U, error)) (U, error) {
	var zero U
	return i.X().FoldX(zero, fn)
}

// FoldX folds the stream into a single value, starting from init, and aborts on the
// first error. Terminal.
func (i Iter2[K, V]) FoldX[U any](init U, fn func(U, K, V) (U, error)) (U, error) {
	return i.X().FoldX(init, fn)
}

// Reduce folds the stream into a single value, seeded with U's zero value, and
// returns any abort. Terminal.
func (i IterX[T]) Reduce[U any](fn func(U, T) U) (U, error) {
	var zero U
	return i.FoldX(zero, func(acc U, t T) (U, error) {
		return fn(acc, t), nil
	})
}

// Fold folds the stream into a single value, starting from init, and returns any
// abort. Terminal.
func (i IterX[T]) Fold[U any](init U, fn func(U, T) U) (U, error) {
	return i.FoldX(init, func(acc U, t T) (U, error) {
		return fn(acc, t), nil
	})
}

// ReduceX folds the stream into a single value, seeded with U's zero value, and
// aborts on the first error. Terminal.
func (i IterX[T]) ReduceX[U any](fn func(U, T) (U, error)) (U, error) {
	var zero U
	return i.FoldX(zero, fn)
}

// FoldX folds the stream into a single value, starting from init. An error from fn
// aborts, as does an abort from upstream, and either discards the accumulator.
// Terminal.
func (i IterX[T]) FoldX[U any](init U, fn func(U, T) (U, error)) (U, error) {
	result := init
	err := i.each(func(t T) (bool, error) {
		next, e := fn(result, t)
		if e != nil {
			return false, e
		}
		result = next
		return true, nil
	})
	if err != nil {
		var zero U
		return zero, err
	}
	return result, nil
}

// Reduce folds the stream into a single value, seeded with U's zero value, and
// returns any abort. Terminal.
func (i Iter2X[K, V]) Reduce[U any](fn func(U, K, V) U) (U, error) {
	var zero U
	return i.FoldX(zero, func(acc U, k K, v V) (U, error) {
		return fn(acc, k, v), nil
	})
}

// Fold folds the stream into a single value, starting from init, and returns any
// abort. Terminal.
func (i Iter2X[K, V]) Fold[U any](init U, fn func(U, K, V) U) (U, error) {
	return i.FoldX(init, func(acc U, k K, v V) (U, error) {
		return fn(acc, k, v), nil
	})
}

// ReduceX folds the stream into a single value, seeded with U's zero value, and
// aborts on the first error. Terminal.
func (i Iter2X[K, V]) ReduceX[U any](fn func(U, K, V) (U, error)) (U, error) {
	var zero U
	return i.FoldX(zero, fn)
}

// FoldX folds the stream into a single value, starting from init. An error from fn
// aborts, as does an abort from upstream, and either discards the accumulator.
// Terminal.
func (i Iter2X[K, V]) FoldX[U any](init U, fn func(U, K, V) (U, error)) (U, error) {
	result := init
	err := i.each(func(k K, v V) (bool, error) {
		next, e := fn(result, k, v)
		if e != nil {
			return false, e
		}
		result = next
		return true, nil
	})
	if err != nil {
		var zero U
		return zero, err
	}
	return result, nil
}
