package collections

import (
	"context"
	"runtime"
	"sync"
)

// ParallelFilter runs fn on every element concurrently and returns the accepted ones
// in their original order. A maxWorkers of zero or less means runtime.NumCPU.
func (s Slice[T]) ParallelFilter(fn func(T) bool, maxWorkers int) Slice[T] {
	result, _ := s.ParallelFilterXC(context.Background(),
		func(_ context.Context, t T) (bool, error) {
			return fn(t), nil
		}, maxWorkers)
	return result
}

// ParallelFilterX runs fn on every element concurrently. The first error cancels the
// work not yet started and is returned.
func (s Slice[T]) ParallelFilterX(fn func(T) (bool, error), maxWorkers int) (Slice[T], error) {
	return s.ParallelFilterXC(context.Background(),
		func(_ context.Context, t T) (bool, error) {
			return fn(t)
		}, maxWorkers)
}

// ParallelFilterXC runs fn on every element concurrently, passing each call a context
// that is cancelled by the first error or by the caller's own cancellation.
func (s Slice[T]) ParallelFilterXC(ctx context.Context, fn func(context.Context, T) (bool, error), maxWorkers int) (Slice[T], error) {
	var wg sync.WaitGroup
	keep := make([]bool, len(s))

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
				ok, err := fn(ctx, t)
				if err != nil {
					cancel(err)
					return
				}

				keep[i] = ok
			})
		}
	}

	wg.Wait()

	if ctx.Err() != nil {
		return nil, context.Cause(ctx)
	}

	var result Slice[T]
	for i, ok := range keep {
		if ok {
			result = append(result, s[i])
		}
	}
	return result, nil
}

// ParallelFilter drains the stream and runs fn on every element concurrently.
func (i Iter[T]) ParallelFilter(fn func(T) bool, maxWorkers int) Slice[T] {
	return i.Slice().ParallelFilter(fn, maxWorkers)
}

// ParallelFilterX drains the stream and runs fn concurrently, stopping at the first
// error.
func (i Iter[T]) ParallelFilterX(fn func(T) (bool, error), maxWorkers int) (Slice[T], error) {
	return i.Slice().ParallelFilterX(fn, maxWorkers)
}

// ParallelFilterXC drains the stream and runs fn concurrently under ctx.
func (i Iter[T]) ParallelFilterXC(ctx context.Context, fn func(context.Context, T) (bool, error), maxWorkers int) (Slice[T], error) {
	return i.Slice().ParallelFilterXC(ctx, fn, maxWorkers)
}

// ParallelFilter drains the stream and runs fn on every element concurrently. An
// abort while draining is returned before any work starts.
func (i IterX[T]) ParallelFilter(fn func(T) bool, maxWorkers int) (Slice[T], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelFilter(fn, maxWorkers), nil
}

// ParallelFilterX drains the stream and runs fn concurrently, stopping at the first
// error from either.
func (i IterX[T]) ParallelFilterX(fn func(T) (bool, error), maxWorkers int) (Slice[T], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelFilterX(fn, maxWorkers)
}

// ParallelFilterXC drains the stream and runs fn concurrently under ctx.
func (i IterX[T]) ParallelFilterXC(ctx context.Context, fn func(context.Context, T) (bool, error), maxWorkers int) (Slice[T], error) {
	result, err := i.Slice()
	if err != nil {
		return nil, err
	}
	return result.ParallelFilterXC(ctx, fn, maxWorkers)
}
