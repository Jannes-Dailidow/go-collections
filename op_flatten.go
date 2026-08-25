package collections

// Flatten concatenates a stream of streams. Lazy.
//
// It is a function rather than a method because a method receiver cannot destructure
// its own T into Iter[U].
func Flatten[T any](i Iter[Iter[T]]) Iter[T] {
	return func(yield func(T) bool) {
		for inner := range i {
			for t := range inner {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlattenSlices concatenates a stream of slices, which is the shape that turns up
// when a Map produced one slice per element. Lazy.
func FlattenSlices[T any](i Iter[Slice[T]]) Iter[T] {
	return func(yield func(T) bool) {
		for inner := range i {
			for _, t := range inner {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// FlatMap yields every element of every stream fn returns. Lazy.
func (i Iter[T]) FlatMap[U any](fn func(T) Iter[U]) Iter[U] {
	return func(yield func(U) bool) {
		for t := range i {
			for u := range fn(t) {
				if !yield(u) {
					return
				}
			}
		}
	}
}

// FlatMapX yields every element of every stream fn returns, and aborts on the first
// error, whether fn's or an inner stream's. Lazy.
func (i Iter[T]) FlatMapX[U any](fn func(T) (IterX[U], error)) IterX[U] {
	return i.X().FlatMapX(fn)
}

// FlatMap yields every element of every stream fn returns, preserving any abort.
// Lazy.
func (i IterX[T]) FlatMap[U any](fn func(T) Iter[U]) IterX[U] {
	return i.FlatMapX(func(t T) (IterX[U], error) {
		return fn(t).X(), nil
	})
}

// FlatMapX yields every element of every stream fn returns. An error from fn aborts
// the stream, as does an abort from an inner stream or from upstream. Lazy.
func (i IterX[T]) FlatMapX[U any](fn func(T) (IterX[U], error)) IterX[U] {
	return func(yield func(U) bool) error {
		return i.each(func(t T) (bool, error) {
			inner, err := fn(t)
			if err != nil {
				return false, err
			}
			stopped := false
			if err := inner(func(u U) bool {
				if !yield(u) {
					stopped = true
					return false
				}
				return true
			}); err != nil {
				return false, err
			}
			return !stopped, nil
		})
	}
}
