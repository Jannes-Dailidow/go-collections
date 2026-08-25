package collections

import "testing"

func newTestMap() Map[string, int] {
	return Map[string, int]{"a": 1, "b": 2}
}

func TestMapLenHasAndGet(t *testing.T) {
	m := newTestMap()
	if m.Len() != 2 || m.IsEmpty() {
		t.Fatal("len or IsEmpty is wrong")
	}
	if !m.Has("a") || m.Has("z") {
		t.Fatal("Has is wrong")
	}
	if v, ok := m.Get("b"); !ok || v != 2 {
		t.Fatalf("got %v %v", v, ok)
	}
	if _, ok := m.Get("z"); ok {
		t.Fatal("absent key reported present")
	}
}

func TestMapGetOr(t *testing.T) {
	m := newTestMap()
	if got := m.GetOr("a", 99); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}
	if got := m.GetOr("z", 99); got != 99 {
		t.Fatalf("want the fallback, got %d", got)
	}
}

func TestMapPutAndDelete(t *testing.T) {
	m := make(Map[string, int])
	m.Put("a", 1)
	m.Put("a", 2)
	if v, _ := m.Get("a"); v != 2 {
		t.Fatalf("Put must replace, got %d", v)
	}
	m.Delete("a")
	m.Delete("absent")
	if !m.IsEmpty() {
		t.Fatalf("want empty, got %v", m)
	}
}

func TestMapClearKeepsTheMapUsable(t *testing.T) {
	m := newTestMap()
	m.Clear()
	if !m.IsEmpty() {
		t.Fatal("want empty")
	}
	m.Put("c", 3)
	if m.Len() != 1 {
		t.Fatal("must still be usable after Clear")
	}
}

func TestMapCloneIsIndependent(t *testing.T) {
	m := newTestMap()
	c := m.Clone()
	c.Put("a", 99)
	if v, _ := m.Get("a"); v != 1 {
		t.Fatalf("the original changed: %d", v)
	}
}

func TestMapMergeLastWriteWins(t *testing.T) {
	m := newTestMap()
	m.Merge(Map[string, int]{"b": 20, "c": 3})
	if v, _ := m.Get("b"); v != 20 {
		t.Fatalf("want 20, got %d", v)
	}
	if v, _ := m.Get("c"); v != 3 {
		t.Fatalf("want 3, got %d", v)
	}
}

func TestMapMergeFuncResolvesCollisions(t *testing.T) {
	m := newTestMap()
	m.MergeFunc(Map[string, int]{"b": 20, "c": 3}, func(existing, incoming int) int {
		return existing + incoming
	})
	if v, _ := m.Get("b"); v != 22 {
		t.Fatalf("want 22, got %d", v)
	}
	if v, _ := m.Get("c"); v != 3 {
		t.Fatalf("a new key must not go through resolve, got %d", v)
	}
}

func TestMapInvertBy(t *testing.T) {
	m := newTestMap()
	got := m.InvertBy(func(_ string, v int) int { return v })
	if len(got) != 2 || got[1] != "a" || got[2] != "b" {
		t.Fatalf("got %v", got)
	}
}

func TestMapInvertByCanKeyOnTheKeyToo(t *testing.T) {
	m := newTestMap()
	got := m.InvertBy(func(k string, v int) string { return k + "!" })
	if got["a!"] != "a" {
		t.Fatalf("got %v", got)
	}
}

func TestMapToOrderedMapHoldsEverything(t *testing.T) {
	m := newTestMap()
	got := m.OrderedMap()
	if got.Len() != 2 || !got.Has("a") || !got.Has("b") {
		t.Fatalf("got %v", got.Native())
	}
}
