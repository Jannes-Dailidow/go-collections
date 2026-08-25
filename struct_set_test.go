package collections

import "testing"

func TestSetAddReportsNewness(t *testing.T) {
	s := make(Set[int])
	if !s.Add(1) {
		t.Fatal("first Add must report true")
	}
	if s.Add(1) {
		t.Fatal("a repeat Add must report false")
	}
	if s.Len() != 1 {
		t.Fatalf("want 1, got %d", s.Len())
	}
}

func TestSetRemoveReportsPresence(t *testing.T) {
	s := CollectSet(ints(3))
	if !s.Remove(1) {
		t.Fatal("Remove of a present value must report true")
	}
	if s.Remove(1) {
		t.Fatal("Remove of an absent value must report false")
	}
	if s.Len() != 2 {
		t.Fatalf("want 2, got %d", s.Len())
	}
}

func TestSetHasLenAndIsEmpty(t *testing.T) {
	s := CollectSet(ints(3))
	if !s.Has(2) || s.Has(9) {
		t.Fatal("Has is wrong")
	}
	if s.IsEmpty() || s.Len() != 3 {
		t.Fatal("len or IsEmpty is wrong")
	}
	var nilSet Set[int]
	if !nilSet.IsEmpty() || nilSet.Has(1) {
		t.Fatal("a nil set is empty and holds nothing")
	}
}

func TestSetClearAndClone(t *testing.T) {
	s := CollectSet(ints(3))
	c := s.Clone()
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("want empty after Clear")
	}
	if c.Len() != 3 {
		t.Fatal("the clone must be untouched")
	}
	s.Add(9)
	if s.Len() != 1 {
		t.Fatal("must still be usable after Clear")
	}
}

func TestSetToOrderedSetHoldsEverything(t *testing.T) {
	s := CollectSet(ints(3))
	got := s.OrderedSet()
	if got.Len() != 3 || !got.Has(2) {
		t.Fatalf("got %v", got.Native())
	}
}

func TestNewOrderedSetAndZeroValueBothWork(t *testing.T) {
	made := NewOrderedSet[int]()
	made.Add(1)
	if made.Len() != 1 {
		t.Fatal("NewOrderedSet is not usable")
	}

	zero := &OrderedSet[int]{}
	zero.Add(1)
	zero.Add(1)
	if zero.Len() != 1 {
		t.Fatal("the zero value must allocate on first Add")
	}
}

func TestNewOrderedMapAndZeroValueBothWork(t *testing.T) {
	made := NewOrderedMap[string, int]()
	made.Put("a", 1)
	if made.Len() != 1 {
		t.Fatal("NewOrderedMap is not usable")
	}

	zero := &OrderedMap[string, int]{}
	zero.Put("a", 1)
	zero.Put("a", 2)
	if zero.Len() != 1 {
		t.Fatal("the zero value must allocate on first Put")
	}
	if v, _ := zero.Get("a"); v != 2 {
		t.Fatalf("want 2, got %d", v)
	}
}
