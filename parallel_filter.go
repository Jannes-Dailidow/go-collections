package slicest

import (
	"context"
	"runtime"
	"sync"
)

func ParallelFilter[T any, S ~[]T](s S, fn func(T) bool, maxWorkers int) S {
	return ParallelFilterI(s, func(_ int, t T) bool {
		return fn(t)
	}, maxWorkers)
}

func ParallelFilterI[T any, S ~[]T](s S, fn func(int, T) bool, maxWorkers int) S {
	var wg sync.WaitGroup
	var result S
	keep := make([]bool, len(s))

	if maxWorkers <= 0 {
		maxWorkers = max(1, runtime.NumCPU())
	}
	semaphore := make(chan struct{}, maxWorkers)

	for i, t := range s {
		semaphore <- struct{}{}
		wg.Go(func() {
			defer func() { <-semaphore }()
			keep[i] = fn(i, t)
		})
	}

	wg.Wait()

	for i, b := range keep {
		if b {
			result = append(result, s[i])
		}
	}

	return result
}

func ParallelFilterX[T any, S ~[]T](s S, fn func(T) (bool, error), maxWorkers int) (S, error) {
	return ParallelFilterXI(s, func(_ int, t T) (bool, error) {
		return fn(t)
	}, maxWorkers)
}

func ParallelFilterXC[T any, S ~[]T](ctx context.Context, s S, fn func(context.Context, T) (bool, error), maxWorkers int) (S, error) {
	return ParallelFilterXIC(ctx, s, func(ctx context.Context, _ int, t T) (bool, error) {
		return fn(ctx, t)
	}, maxWorkers)
}

func ParallelFilterXI[T any, S ~[]T](s S, fn func(int, T) (bool, error), maxWorkers int) (S, error) {
	return ParallelFilterXIC(context.Background(), s, func(_ context.Context, i int, t T) (bool, error) {
		return fn(i, t)
	}, maxWorkers)
}

func ParallelFilterXIC[T any, S ~[]T](ctx context.Context, s S, fn func(context.Context, int, T) (bool, error), maxWorkers int) (S, error) {
	var wg sync.WaitGroup
	var result S
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
				b, err := fn(ctx, i, t)
				if err != nil {
					cancel(err)
					return
				}

				keep[i] = b
			})
		}
	}

	wg.Wait()

	if ctx.Err() != nil {
		return nil, context.Cause(ctx)
	}

	for i, b := range keep {
		if b {
			result = append(result, s[i])
		}
	}
	return result, nil
}
