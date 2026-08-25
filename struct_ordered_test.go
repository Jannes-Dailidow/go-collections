package collections

import "testing"

func orderedInts(values ...int) *OrderedSet[int] {
	result := NewOrderedSet[int]()
	for _, v := range values {
		result.Add(v)
	}
	return result
}

func TestOrderedSetKeepsFirstSeenPosition(t *testing.T) {
	s := orderedInts(5, 3, 5, 9)
	assertInts(t, s.Native(), 5, 3, 9)
	if s.IndexOf(3) != 1 || s.IndexOf(9) != 2 {
		t.Fatalf("positions are wrong: %v", s.Native())
	}
	if s.IndexOf(99) != -1 {
		t.Fatal("absent value must be -1")
	}
}

func TestOrderedSetAt(t *testing.T) {
	s := orderedInts(5, 3, 9)
	if got := s.At(1); got == nil || *got != 3 {
		t.Fatalf("got %v", got)
	}
	for _, n := range []int{-1, 3} {
		if s.At(n) != nil {
			t.Fatalf("At(%d) must be nil", n)
		}
	}
}

// At hands back a copy, because writing a new value through a pointer into the order
// would leave the index describing something that is no longer there.
func TestOrderedSetAtReturnsACopy(t *testing.T) {
	s := orderedInts(5, 3, 9)
	*s.At(0) = 99
	assertInts(t, s.Native(), 5, 3, 9)
	if s.IndexOf(5) != 0 {
		t.Fatal("the index must be intact")
	}
}

// Removal has to repair every later position.
func TestOrderedSetRemoveRepairsTheIndex(t *testing.T) {
	s := orderedInts(1, 2, 3, 4)
	if !s.Remove(2) {
		t.Fatal("want true")
	}
	assertInts(t, s.Native(), 1, 3, 4)
	if s.IndexOf(1) != 0 || s.IndexOf(3) != 1 || s.IndexOf(4) != 2 {
		t.Fatalf("indexes not repaired: 3 is at %d, 4 is at %d", s.IndexOf(3), s.IndexOf(4))
	}
	if s.Remove(99) {
		t.Fatal("removing an absent value must report false")
	}
}

func TestOrderedSetClearAndClone(t *testing.T) {
	s := orderedInts(1, 2, 3)
	c := s.Clone()
	s.Clear()
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatal("want empty")
	}
	assertInts(t, c.Native(), 1, 2, 3)
	s.Add(9)
	if s.Len() != 1 || s.IndexOf(9) != 0 {
		t.Fatal("must be usable after Clear")
	}
}

func TestOrderedSetCloneIsIndependent(t *testing.T) {
	s := orderedInts(1, 2)
	c := s.Clone()
	c.Add(3)
	if s.Len() != 2 {
		t.Fatal("the original grew")
	}
	if c.Len() != 3 {
		t.Fatal("the clone did not")
	}
}

func TestOrderedSetStructuralMethodsAreNilSafe(t *testing.T) {
	var s *OrderedSet[int]
	if s.Len() != 0 || !s.IsEmpty() || s.Has(1) || s.IndexOf(1) != -1 || s.At(0) != nil {
		t.Fatal("a nil OrderedSet must answer emptily")
	}
	if s.Clone() != nil {
		t.Fatal("cloning nil gives nil")
	}
	s.Clear()
	if s.Remove(1) {
		t.Fatal("removing from nil reports false")
	}
}

func TestOrderedMapGetPutAndPosition(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("c", 1)
	m.Put("a", 2)
	m.Put("c", 9)
	if m.Len() != 2 {
		t.Fatalf("want 2, got %d", m.Len())
	}
	if v, ok := m.Get("c"); !ok || v != 9 {
		t.Fatalf("Put must update in place, got %v", v)
	}
	if m.IndexOf("c") != 0 || m.IndexOf("a") != 1 {
		t.Fatal("an updated key must keep its position")
	}
	if m.GetOr("z", 7) != 7 {
		t.Fatal("GetOr must fall back")
	}
}

func TestOrderedMapAt(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	got := m.At(1)
	if got == nil || got.Key != "b" || got.Value != 2 {
		t.Fatalf("got %v", got)
	}
	if m.At(2) != nil || m.At(-1) != nil {
		t.Fatal("out of range must be nil")
	}
}

func TestOrderedMapAtReturnsACopy(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 1)
	got := m.At(0)
	got.Key = "zzz"
	got.Value = 99
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatal("writing through At must not reach the map")
	}
	if m.IndexOf("a") != 0 {
		t.Fatal("the index must be intact")
	}
}

func TestOrderedMapDeleteRepairsTheIndex(t *testing.T) {
	m := NewOrderedMap[string, int]()
	for i, k := range []string{"a", "b", "c", "d"} {
		m.Put(k, i)
	}
	if !m.Delete("b") {
		t.Fatal("want true")
	}
	if m.Len() != 3 {
		t.Fatalf("want 3, got %d", m.Len())
	}
	if m.IndexOf("a") != 0 || m.IndexOf("c") != 1 || m.IndexOf("d") != 2 {
		t.Fatalf("indexes not repaired: c at %d, d at %d", m.IndexOf("c"), m.IndexOf("d"))
	}
	if m.Has("b") {
		t.Fatal("b is still there")
	}
	if m.Delete("absent") {
		t.Fatal("deleting an absent key reports false")
	}
}

// Delete then Put must not corrupt the index, which is the combination most likely
// to expose an off-by-one in the repair loop.
func TestOrderedMapDeleteThenPut(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Delete("a")
	m.Put("c", 3)
	var keys []string
	for k := range m.All() {
		keys = append(keys, k)
	}
	if len(keys) != 2 || keys[0] != "b" || keys[1] != "c" {
		t.Fatalf("got %v", keys)
	}
	if m.IndexOf("b") != 0 || m.IndexOf("c") != 1 {
		t.Fatalf("b at %d, c at %d", m.IndexOf("b"), m.IndexOf("c"))
	}
}

func TestOrderedMapBackward(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	var keys []string
	for k := range m.Backward() {
		keys = append(keys, k)
	}
	if len(keys) != 3 || keys[0] != "c" || keys[2] != "a" {
		t.Fatalf("got %v", keys)
	}
}

func TestOrderedMapStructuralMethodsAreNilSafe(t *testing.T) {
	var m *OrderedMap[string, int]
	if m.Len() != 0 || !m.IsEmpty() || m.Has("a") || m.IndexOf("a") != -1 || m.At(0) != nil {
		t.Fatal("a nil OrderedMap must answer emptily")
	}
	if _, ok := m.Get("a"); ok {
		t.Fatal("nil holds nothing")
	}
	if m.GetOr("a", 5) != 5 {
		t.Fatal("GetOr on nil must fall back")
	}
	if m.Clone() != nil {
		t.Fatal("cloning nil gives nil")
	}
	m.Clear()
	if m.Delete("a") {
		t.Fatal("deleting from nil reports false")
	}
	for range m.Backward() {
		t.Fatal("nil yields nothing backwards")
	}
}
