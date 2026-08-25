package collections

// Scan yields the running accumulator, one value per element. The initial value is
// not yielded, so the output is the same length as the input. Lazy.
func (i Iter[T]) Scan[U any](init U, fn func(U, T) U) Iter[U] {
	return func(yield func(U) bool) {
		acc := init
		for t := range i {
			acc = fn(acc, t)
			if !yield(acc) {
				return
			}
		}
	}
}

// ScanX yields the running accumulator and aborts on the first error. Lazy.
func (i Iter[T]) ScanX[U any](init U, fn func(U, T) (U, error)) IterX[U] {
	return i.X().ScanX(init, fn)
}

// Scan yields the running accumulator, preserving any abort. Lazy.
func (i IterX[T]) Scan[U any](init U, fn func(U, T) U) IterX[U] {
	return i.ScanX(init, func(acc U, t T) (U, error) {
		return fn(acc, t), nil
	})
}

// ScanX yields the running accumulator. An error from fn aborts the stream, as does
// an abort from upstream. Lazy.
func (i IterX[T]) ScanX[U any](init U, fn func(U, T) (U, error)) IterX[U] {
	return func(yield func(U) bool) error {
		acc := init
		return i.each(func(t T) (bool, error) {
			next, err := fn(acc, t)
			if err != nil {
				return false, err
			}
			acc = next
			return yield(acc), nil
		})
	}
}
