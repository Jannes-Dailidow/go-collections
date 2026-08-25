package collections

import (
	"context"
	"runtime"
	"sync"
)

// ParallelEach runs fn on every element concurrently, for side effects only. A
// maxWorkers of zero or less means runtime.NumCPU.
//
// fn is called from many goroutines at once, so whatever it touches has to be safe
// for that. Nothing here synchronizes it for you.
func (s Slice[T]) ParallelEach(fn func(T), maxWorkers int) {
	s.ParallelEachXC(context.Background(),
		func(_ context.Context, t T) error {
			fn(t)
			return nil
		}, maxWorkers)
}

// ParallelEachX runs fn on every element concurrently. The first error cancels the
// work not yet started and is returned.
func (s Slice[T]) ParallelEachX(fn func(T) error, maxWorkers int) error {
	return s.ParallelEachXC(context.Background(),
		func(_ context.Context, t T) error {
			return fn(t)
		}, maxWorkers)
}

// ParallelEachXC runs fn on every element concurrently, passing each call a context
// that is cancelled by the first error or by the caller's own cancellation.
func (s Slice[T]) ParallelEachXC(ctx context.Context, fn func(context.Context, T) error, maxWorkers int) error {
	var wg sync.WaitGroup

	if maxWorkers <= 0 {
		maxWorkers = max(1, runtime.NumCPU())
	}
	semaphore := make(chan struct{}, maxWorkers)

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

sLoop:
	for _, t := range s {
		select {
		case <-ctx.Done():
			break sLoop
		case semaphore <- struct{}{}:
			wg.Go(func() {
				defer func() { <-semaphore }()
				if err := fn(ctx, t); err != nil {
					cancel(err)
				}
			})
		}
	}

	wg.Wait()

	if ctx.Err() != nil {
		return context.Cause(ctx)
	}

	return nil
}

// ParallelEach drains the stream and runs fn on every element concurrently.
func (i Iter[T]) ParallelEach(fn func(T), maxWorkers int) {
	i.Slice().ParallelEach(fn, maxWorkers)
}

// ParallelEachX drains the stream and runs fn concurrently, stopping at the first
// error.
func (i Iter[T]) ParallelEachX(fn func(T) error, maxWorkers int) error {
	return i.Slice().ParallelEachX(fn, maxWorkers)
}

// ParallelEachXC drains the stream and runs fn concurrently under ctx.
func (i Iter[T]) ParallelEachXC(ctx context.Context, fn func(context.Context, T) error, maxWorkers int) error {
	return i.Slice().ParallelEachXC(ctx, fn, maxWorkers)
}

// ParallelEach drains the stream and runs fn on every element concurrently. An abort
// while draining is returned before any work starts.
func (i IterX[T]) ParallelEach(fn func(T), maxWorkers int) error {
	result, err := i.Slice()
	if err != nil {
		return err
	}
	result.ParallelEach(fn, maxWorkers)
	return nil
}

// ParallelEachX drains the stream and runs fn concurrently, stopping at the first
// error from either.
func (i IterX[T]) ParallelEachX(fn func(T) error, maxWorkers int) error {
	result, err := i.Slice()
	if err != nil {
		return err
	}
	return result.ParallelEachX(fn, maxWorkers)
}

// ParallelEachXC drains the stream and runs fn concurrently under ctx.
func (i IterX[T]) ParallelEachXC(ctx context.Context, fn func(context.Context, T) error, maxWorkers int) error {
	result, err := i.Slice()
	if err != nil {
		return err
	}
	return result.ParallelEachXC(ctx, fn, maxWorkers)
}
