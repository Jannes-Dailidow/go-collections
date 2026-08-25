package collections

import "slices"

// Chunk groups the stream into batches of n, the last one short if the stream does
// not divide evenly. A non-positive n yields nothing. Lazy.
//
// Chunk and [Window] are functions rather than methods, and this one is not a choice
// between two designs the way [Pairs] was. A method on Iter[T] returning
// Iter[Slice[T]] is an instantiation cycle on its own: instantiating Iter[T] pulls in
// Iter[Slice[T]], which pulls in Iter[Slice[Slice[T]]], without end.
//
// Each batch is a fresh slice, so a consumer may keep it.
func Chunk[T any](i Iter[T], n int) Iter[Slice[T]] {
	return func(yield func(Slice[T]) bool) {
		if n <= 0 {
			return
		}
		var batch Slice[T]
		for t := range i {
			batch = append(batch, t)
			if len(batch) == n {
				if !yield(batch) {
					return
				}
				batch = nil
			}
		}
		if len(batch) > 0 {
			yield(batch)
		}
	}
}

// ChunkX is the fallible counterpart of [Chunk]. A batch cut short by an abort is
// discarded rather than yielded.
func ChunkX[T any](i IterX[T], n int) IterX[Slice[T]] {
	return func(yield func(Slice[T]) bool) error {
		if n <= 0 {
			return nil
		}
		var batch Slice[T]
		stopped := false
		err := i.each(func(t T) (bool, error) {
			batch = append(batch, t)
			if len(batch) < n {
				return true, nil
			}
			full := batch
			batch = nil
			if !yield(full) {
				stopped = true
				return false, nil
			}
			return true, nil
		})
		if err != nil {
			return err
		}
		if !stopped && len(batch) > 0 {
			yield(batch)
		}
		return nil
	}
}

// Window slides a window of n over the stream, yielding one slice per position.
// Nothing is yielded until n elements have arrived, so a stream shorter than n yields
// nothing. A non-positive n yields nothing. Lazy.
//
// Each window is a copy, so a consumer may keep it.
func Window[T any](i Iter[T], n int) Iter[Slice[T]] {
	return func(yield func(Slice[T]) bool) {
		if n <= 0 {
			return
		}
		buf := make(Slice[T], 0, n)
		for t := range i {
			if len(buf) == n {
				copy(buf, buf[1:])
				buf = buf[:n-1]
			}
			buf = append(buf, t)
			if len(buf) == n && !yield(slices.Clone(buf)) {
				return
			}
		}
	}
}

// WindowX is the fallible counterpart of [Window].
func WindowX[T any](i IterX[T], n int) IterX[Slice[T]] {
	return func(yield func(Slice[T]) bool) error {
		if n <= 0 {
			return nil
		}
		buf := make(Slice[T], 0, n)
		return i.each(func(t T) (bool, error) {
			if len(buf) == n {
				copy(buf, buf[1:])
				buf = buf[:n-1]
			}
			buf = append(buf, t)
			if len(buf) == n {
				return yield(slices.Clone(buf)), nil
			}
			return true, nil
		})
	}
}
