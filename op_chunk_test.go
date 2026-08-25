package collections

import (
	"errors"
	"testing"
)

func TestChunkGroupsIntoBatches(t *testing.T) {
	got := Chunk(ints(6), 2).Slice()
	if len(got) != 3 {
		t.Fatalf("want 3 batches, got %v", got)
	}
	assertInts(t, got[0].Native(), 0, 1)
	assertInts(t, got[2].Native(), 4, 5)
}

func TestChunkYieldsAShortFinalBatch(t *testing.T) {
	got := Chunk(ints(5), 2).Slice()
	if len(got) != 3 {
		t.Fatalf("want 3 batches, got %v", got)
	}
	assertInts(t, got[2].Native(), 4)
}

func TestChunkOfNonPositiveNYieldsNothing(t *testing.T) {
	for _, n := range []int{0, -1} {
		if got := Chunk(ints(6), n).Slice(); len(got) != 0 {
			t.Fatalf("Chunk(%d) must yield nothing, got %v", n, got)
		}
	}
}

// Batches must not share a backing array, or a consumer that keeps one watches it
// change underneath.
func TestChunkBatchesAreIndependent(t *testing.T) {
	got := Chunk(ints(4), 2).Slice()
	got[0][0] = 99
	assertInts(t, got[1].Native(), 2, 3)
	if got[0][0] != 99 {
		t.Fatal("the first batch should have kept the write")
	}
}

func TestChunkIsLazy(t *testing.T) {
	var pulls int
	got := Chunk(counting(1000, &pulls), 2).Take(2).Slice()
	if len(got) != 2 {
		t.Fatalf("want 2 batches, got %v", got)
	}
	if pulls != 4 {
		t.Fatalf("want 4 pulls, got %d", pulls)
	}
}

func TestChunkIsReIterable(t *testing.T) {
	i := Chunk(ints(5), 2)
	for pass := range 2 {
		got := i.Slice()
		if len(got) != 3 || len(got[0]) != 2 || len(got[2]) != 1 {
			t.Fatalf("pass %d: got %v", pass, got)
		}
	}
}

func TestChunkXDiscardsTheBatchCutShortByAnAbort(t *testing.T) {
	_, err := ChunkX(abortAt(5), 2).Native()
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}

func TestChunkXYieldsWholeBatchesBeforeAnAbort(t *testing.T) {
	var batches int
	err := ChunkX(abortAt(5), 2).Each(func(Slice[int]) { batches++ })
	if !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
	if batches != 2 {
		t.Fatalf("want the 2 complete batches before the abort, got %d", batches)
	}
}

func TestWindowSlidesOverTheStream(t *testing.T) {
	got := Window(ints(4), 2).Slice()
	if len(got) != 3 {
		t.Fatalf("want 3 windows, got %v", got)
	}
	assertInts(t, got[0].Native(), 0, 1)
	assertInts(t, got[1].Native(), 1, 2)
	assertInts(t, got[2].Native(), 2, 3)
}

func TestWindowShorterThanNYieldsNothing(t *testing.T) {
	if got := Window(ints(2), 3).Slice(); len(got) != 0 {
		t.Fatalf("want nothing, got %v", got)
	}
}

func TestWindowOfExactlyNYieldsOne(t *testing.T) {
	got := Window(ints(3), 3).Slice()
	if len(got) != 1 {
		t.Fatalf("want 1 window, got %v", got)
	}
	assertInts(t, got[0].Native(), 0, 1, 2)
}

// Windows overlap, so sharing the buffer would be visible immediately.
func TestWindowsAreCopies(t *testing.T) {
	got := Window(ints(4), 2).Slice()
	got[0][0] = 99
	assertInts(t, got[1].Native(), 1, 2)
}

func TestWindowIsReIterable(t *testing.T) {
	i := Window(ints(4), 2)
	for pass := range 2 {
		if got := i.Slice(); len(got) != 3 {
			t.Fatalf("pass %d: want 3 windows, got %v", pass, got)
		}
	}
}

func TestWindowIsLazy(t *testing.T) {
	var pulls int
	got := Window(counting(1000, &pulls), 2).Take(2).Slice()
	if len(got) != 2 {
		t.Fatalf("want 2 windows, got %v", got)
	}
	if pulls != 3 {
		t.Fatalf("want 3 pulls for 2 windows of 2, got %d", pulls)
	}
}

func TestWindowXCarriesTheAbort(t *testing.T) {
	if _, err := WindowX(abortAt(4), 2).Native(); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}
