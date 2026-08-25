package collections

import (
	"context"
	"runtime"
	"sync"
)

// ParallelMap applies fn to every element concurrently and returns the results in
// their original order. A maxWorkers of zero or less means runtime.NumCPU.
//
// The parallel family lives on [Slice] because it has to know how many elements there
// are before it can start: the results go into a preallocated slice so that order
// survives. The [Iter] forms materialize first and then hand over to these.
func (s Slice[T]) ParallelMap[U any](fn func(T) U, maxWorkers int) Slice[U] {
	result, _ := s.ParallelMapXC(context.Background(),
		func(_ context.Context, t T) (U, error) {
			return fn(t), nil
		}, maxWorkers)
	return result
}

// ParallelMapX applies fn to every element concurrently. The first error cancels the
// work not yet started and is returned; the partial result is discarded.
func (s Slice[T]) ParallelMapX[U any](fn func(T) (U, error), maxWorkers int) (Slice[U], error) {
	return s.ParallelMapXC(context.Background(),
		func(_ context.Context, t T) (U, error) {
			return fn(t)
		}, maxWorkers)
}

// ParallelMapXC applies fn to every element concurrently, passing each call a context
// that is cancelled by the first error or by the caller's own cancellation.
func (s Slice[T]) ParallelMapXC[U any](ctx context.Context, fn func(context.Context, T) (U, error), maxWorkers int) (Slice[U], error) {
	var wg sync.WaitGroup
	result := make(Slice[U], len(s))

	if maxWorkers <= 0 {
		maxWorkers = max(1, runtime.NumCPU())
	}
	semaphore := make(chan struct{}, maxWorkers)

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

sLoop:
	for i, t := range s {
		select {
		case <-ctx.Done():
			break sLoop
		case semaphore <- struct{}{}:
			wg.Go(func() {
				defer func() { <-semaphore }()
				u, err := fn(ctx, t)
				if err != nil {
					cancel(err)
					return
				}

				result[i] = u
			})
		}
	}

	wg.Wait()

	if ctx.Err() != nil {
		return nil, context.Cause(ctx)
	}

	return result, nil
}

// ParallelMap drains the stream and applies fn to every element concurrently,
// returning the results in stream order.
func (i Iter[T]) ParallelMap[U any](fn func(T) U, maxWorkers int) Slice[U] {
	return i.Slice().ParallelMap(fn, maxWorkers)
}

// ParallelMapX drains the stream and applies fn concurrently, stopping at the first
// error.
func (i Iter[T]) ParallelMapX[U any](fn func(T) (U, error), maxWorkers int) (Slice[U], error) {
	return i.Slice().ParallelMapX(fn, maxWorkers)
}

// ParallelMapXC drains the stream and applies fn concurrently under ctx.
func (i Iter[T]) ParallelMapXC[U any](ctx context.Context, fn func(context.Context, T) (U, error), maxWorkers int) (Slice[U], error) {
	return i.Slice().ParallelMapXC(ctx, fn, maxWorkers)
}

// ParallelMap drains the stream and applies fn to every element concurrently. An
// abort while draining is returned before any work starts.
func (i IterX[T]) ParallelMap[U any](fn func(T) U, maxWorkers int) (Slice[U], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelMap(fn, maxWorkers), nil
}

// ParallelMapX drains the stream and applies fn concurrently, stopping at the first
// error from either.
func (i IterX[T]) ParallelMapX[U any](fn func(T) (U, error), maxWorkers int) (Slice[U], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelMapX(fn, maxWorkers)
}

// ParallelMapXC drains the stream and applies fn concurrently under ctx.
func (i IterX[T]) ParallelMapXC[U any](ctx context.Context, fn func(context.Context, T) (U, error), maxWorkers int) (Slice[U], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelMapXC(ctx, fn, maxWorkers)
}
