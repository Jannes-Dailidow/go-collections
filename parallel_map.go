package slicest

import (
	"context"
	"runtime"
	"sync"
)

func ParallelMap[T, U any, S ~[]T](s S, fn func(T) U, maxWorkers int) []U {
	return ParallelMapI(s, func(_ int, t T) U {
		return fn(t)
	}, maxWorkers)
}

func ParallelMapI[T, U any, S ~[]T](s S, fn func(int, T) U, maxWorkers int) []U {
	var wg sync.WaitGroup
	result := make([]U, len(s))

	if maxWorkers <= 0 {
		maxWorkers = max(1, runtime.NumCPU())
	}
	semaphore := make(chan struct{}, maxWorkers)

	for i, t := range s {
		semaphore <- struct{}{}
		wg.Go(func() {
			defer func() { <-semaphore }()
			result[i] = fn(i, t)
		})
	}

	wg.Wait()

	return result
}

func ParallelMapX[T, U any, S ~[]T](s S, fn func(T) (U, error), maxWorkers int) ([]U, error) {
	return ParallelMapXIC(context.Background(), s, func(_ context.Context, _ int, t T) (U, error) {
		return fn(t)
	}, maxWorkers)
}

func ParallelMapXC[T, U any, S ~[]T](ctx context.Context, s S, fn func(T) (U, error), maxWorkers int) ([]U, error) {
	return ParallelMapXIC(ctx, s, func(_ context.Context, _ int, t T) (U, error) {
		return fn(t)
	}, maxWorkers)
}

func ParallelMapXI[T, U any, S ~[]T](s S, fn func(int, T) (U, error), maxWorkers int) ([]U, error) {
	return ParallelMapXIC(context.Background(), s, func(_ context.Context, i int, t T) (U, error) {
		return fn(i, t)
	}, maxWorkers)
}

func ParallelMapXIC[T, U any, S ~[]T](ctx context.Context, s S, fn func(context.Context, int, T) (U, error), maxWorkers int) ([]U, error) {
	var wg sync.WaitGroup
	result := make([]U, len(s))

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
				u, err := fn(ctx, i, t)
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
