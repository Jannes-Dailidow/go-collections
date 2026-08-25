package collections

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func seq(n int) Slice[int] {
	result := make(Slice[int], n)
	for i := range result {
		result[i] = i
	}
	return result
}

// Order is the whole reason the results go into a preallocated slice.
func TestParallelMapPreservesOrder(t *testing.T) {
	got := seq(200).ParallelMap(func(v int) int { return v * 2 }, 8)
	if got.Len() != 200 {
		t.Fatalf("want 200, got %d", got.Len())
	}
	for i, v := range got {
		if v != i*2 {
			t.Fatalf("position %d holds %d", i, v)
		}
	}
}

func TestParallelMapOnAnEmptySlice(t *testing.T) {
	if got := (Slice[int]{}).ParallelMap(func(v int) int { return v }, 4); got.Len() != 0 {
		t.Fatalf("want empty, got %v", got)
	}
}

// A worker limit of one has to serialize, which the semaphore guarantees.
func TestParallelMapRespectsAWorkerLimitOfOne(t *testing.T) {
	var inFlight, peak atomic.Int64
	seq(50).ParallelMap(func(v int) int {
		n := inFlight.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		inFlight.Add(-1)
		return v
	}, 1)
	if peak.Load() != 1 {
		t.Fatalf("one worker means one at a time, saw %d", peak.Load())
	}
}

func TestParallelMapNeverExceedsTheWorkerLimit(t *testing.T) {
	var inFlight, peak atomic.Int64
	seq(200).ParallelMap(func(v int) int {
		n := inFlight.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		inFlight.Add(-1)
		return v
	}, 4)
	if peak.Load() > 4 {
		t.Fatalf("want at most 4 in flight, saw %d", peak.Load())
	}
}

// Zero or less means NumCPU, so it still has to produce the right answer.
func TestParallelMapWithNoWorkerLimit(t *testing.T) {
	for _, workers := range []int{0, -1} {
		got := seq(100).ParallelMap(func(v int) int { return v + 1 }, workers)
		if got.Len() != 100 || got[99] != 100 {
			t.Fatalf("workers=%d gave %v", workers, got[:5])
		}
	}
}

func TestParallelMapXReturnsTheFirstErrorAndDiscardsTheResult(t *testing.T) {
	got, err := seq(100).ParallelMapX(func(v int) (int, error) {
		if v == 50 {
			return 0, errBoom
		}
		return v, nil
	}, 4)
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != nil {
		t.Fatal("the partial result must be discarded")
	}
}

// The error cancels work not yet started, so not every element is visited.
func TestParallelMapXStopsStartingNewWork(t *testing.T) {
	var started atomic.Int64
	_, err := seq(1000).ParallelMapX(func(v int) (int, error) {
		started.Add(1)
		if v == 0 {
			return 0, errBoom
		}
		return v, nil
	}, 2)
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if started.Load() == 1000 {
		t.Fatal("the cancellation did not stop anything")
	}
}

func TestParallelMapXCHonoursACancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := seq(10).ParallelMapXC(ctx, func(context.Context, int) (int, error) {
		return 0, nil
	}, 2)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
}

// The callback receives the cancellable context, so it can stop its own work.
func TestParallelMapXCPassesTheContextToTheCallback(t *testing.T) {
	var sawLive atomic.Bool
	_, err := seq(4).ParallelMapXC(context.Background(),
		func(ctx context.Context, v int) (int, error) {
			if ctx.Err() == nil {
				sawLive.Store(true)
			}
			return v, nil
		}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sawLive.Load() {
		t.Fatal("the callback never saw a live context")
	}
}

func TestParallelFilterPreservesOrder(t *testing.T) {
	got := seq(100).ParallelFilter(func(v int) bool { return v%10 == 0 }, 8)
	assertInts(t, got.Native(), 0, 10, 20, 30, 40, 50, 60, 70, 80, 90)
}

func TestParallelFilterXReturnsTheFirstError(t *testing.T) {
	got, err := seq(100).ParallelFilterX(func(v int) (bool, error) {
		if v == 50 {
			return false, errBoom
		}
		return true, nil
	}, 4)
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if got != nil {
		t.Fatal("the partial result must be discarded")
	}
}

func TestParallelEachVisitsEveryElement(t *testing.T) {
	var sum atomic.Int64
	seq(100).ParallelEach(func(v int) { sum.Add(int64(v)) }, 8)
	if sum.Load() != 4950 {
		t.Fatalf("want 4950, got %d", sum.Load())
	}
}

func TestParallelEachXReturnsTheFirstError(t *testing.T) {
	err := seq(100).ParallelEachX(func(v int) error {
		if v == 10 {
			return errBoom
		}
		return nil
	}, 4)
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

// The iterator forms materialize first and then hand over.
func TestParallelOnAnIterator(t *testing.T) {
	got := ints(50).ParallelMap(func(v int) int { return v * 2 }, 4)
	if got.Len() != 50 || got[49] != 98 {
		t.Fatalf("got %v", got.Len())
	}
	kept := ints(50).ParallelFilter(func(v int) bool { return v < 3 }, 4)
	assertInts(t, kept.Native(), 0, 1, 2)
}

// A stream that aborts while draining never reaches the worker pool.
func TestParallelOnAnAbortingStreamFailsBeforeAnyWork(t *testing.T) {
	var ran atomic.Int64
	_, err := abortAt(3).ParallelMap(func(v int) int {
		ran.Add(1)
		return v
	}, 4)
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if ran.Load() != 0 {
		t.Fatalf("no work should have started, %d calls ran", ran.Load())
	}
}

func TestParallelEachOnAFallibleStream(t *testing.T) {
	if err := abortAt(2).ParallelEach(func(int) {}, 4); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	var sum atomic.Int64
	if err := ints(10).X().ParallelEach(func(v int) { sum.Add(int64(v)) }, 4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.Load() != 45 {
		t.Fatalf("want 45, got %d", sum.Load())
	}
}

// Parallel and sequential must agree on the answer.
func TestParallelMapAgreesWithMap(t *testing.T) {
	double := func(v int) int { return v * 2 }
	parallel := seq(100).ParallelMap(double, 8)
	sequential := seq(100).Values().Map(double).Slice()
	assertInts(t, parallel.Native(), sequential.Native()...)
}
