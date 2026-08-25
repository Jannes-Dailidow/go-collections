package collections

import (
	"slices"
	"testing"
)

func setOf(values ...int) Set[int] {
	result := make(Set[int])
	for _, v := range values {
		result.Add(v)
	}
	return result
}

// sorted so a Set's arbitrary iteration order does not make the assertion flaky
func sortedOf(s Set[int]) []int {
	result := s.Slice().Native()
	slices.Sort(result)
	return result
}

func TestSetUnion(t *testing.T) {
	assertInts(t, sortedOf(setOf(1, 2).Union(setOf(2, 3))), 1, 2, 3)
	assertInts(t, sortedOf(setOf().Union(setOf(1))), 1)
}

func TestSetIntersect(t *testing.T) {
	assertInts(t, sortedOf(setOf(1, 2, 3).Intersect(setOf(2, 3, 4))), 2, 3)
	if !setOf(1).Intersect(setOf(2)).IsEmpty() {
		t.Fatal("disjoint sets intersect to nothing")
	}
}

func TestSetDiff(t *testing.T) {
	assertInts(t, sortedOf(setOf(1, 2, 3).Diff(setOf(2))), 1, 3)
	if !setOf(1).Diff(setOf(1)).IsEmpty() {
		t.Fatal("want empty")
	}
}

// SymDiff is what v0.1 called Unique.
func TestSetSymDiff(t *testing.T) {
	assertInts(t, sortedOf(setOf(1, 2, 3).SymDiff(setOf(3, 4))), 1, 2, 4)
	if !setOf(1, 2).SymDiff(setOf(1, 2)).IsEmpty() {
		t.Fatal("identical sets have no symmetric difference")
	}
}

func TestSetSubsetAndSuperset(t *testing.T) {
	if !setOf(1, 2).IsSubset(setOf(1, 2, 3)) {
		t.Fatal("want subset")
	}
	if setOf(1, 9).IsSubset(setOf(1, 2, 3)) {
		t.Fatal("not a subset")
	}
	if !setOf().IsSubset(setOf(1)) {
		t.Fatal("the empty set is a subset of everything")
	}
	if !setOf(1, 2, 3).IsSuperset(setOf(1, 2)) {
		t.Fatal("want superset")
	}
	if !setOf(1, 2).IsSubset(setOf(1, 2)) {
		t.Fatal("a set is a subset of itself")
	}
}

func TestSetIsDisjoint(t *testing.T) {
	if !setOf(1, 2).IsDisjoint(setOf(3, 4)) {
		t.Fatal("want disjoint")
	}
	if setOf(1, 2).IsDisjoint(setOf(2, 3)) {
		t.Fatal("they share 2")
	}
	if !setOf().IsDisjoint(setOf(1)) {
		t.Fatal("the empty set is disjoint from everything")
	}
}

// The ordered versions have to produce a predictable order, which is the only reason
// to reach for them over Set.
func TestOrderedSetUnionKeepsLeftOrderFirst(t *testing.T) {
	got := orderedInts(3, 1).Union(orderedInts(1, 9))
	assertInts(t, got.Native(), 3, 1, 9)
}

func TestOrderedSetIntersectKeepsLeftOrder(t *testing.T) {
	got := orderedInts(5, 3, 9).Intersect(orderedInts(9, 3))
	assertInts(t, got.Native(), 3, 9)
}

func TestOrderedSetDiffKeepsLeftOrder(t *testing.T) {
	got := orderedInts(5, 3, 9).Diff(orderedInts(3))
	assertInts(t, got.Native(), 5, 9)
}

func TestOrderedSetSymDiffPutsLeftContributionFirst(t *testing.T) {
	got := orderedInts(5, 3).SymDiff(orderedInts(3, 9))
	assertInts(t, got.Native(), 5, 9)
}

func TestOrderedSetSubsetSupersetAndDisjoint(t *testing.T) {
	if !orderedInts(1, 2).IsSubset(orderedInts(3, 2, 1)) {
		t.Fatal("order must not affect subset")
	}
	if !orderedInts(3, 2, 1).IsSuperset(orderedInts(1)) {
		t.Fatal("want superset")
	}
	if !orderedInts(1).IsDisjoint(orderedInts(2)) {
		t.Fatal("want disjoint")
	}
	if orderedInts(1).IsDisjoint(orderedInts(1)) {
		t.Fatal("not disjoint")
	}
}

func TestOrderedSetOpsAreNilSafe(t *testing.T) {
	var nilSet *OrderedSet[int]
	if got := nilSet.Union(orderedInts(1)); got.Len() != 1 {
		t.Fatal("union with nil on the left")
	}
	if got := orderedInts(1).Union(nilSet); got.Len() != 1 {
		t.Fatal("union with nil on the right")
	}
	if !nilSet.Intersect(orderedInts(1)).IsEmpty() {
		t.Fatal("nil intersects to nothing")
	}
	if !nilSet.IsSubset(orderedInts(1)) {
		t.Fatal("nil is a subset")
	}
	if !nilSet.IsDisjoint(orderedInts(1)) {
		t.Fatal("nil is disjoint")
	}
}

// The two families must agree on content, differing only in order.
func TestSetAndOrderedSetAgreeOnContent(t *testing.T) {
	plain := setOf(5, 3, 9).Diff(setOf(3))
	ordered := orderedInts(5, 3, 9).Diff(orderedInts(3))
	assertInts(t, sortedOf(plain), 5, 9)
	assertInts(t, sortedOf(ordered.Set()), 5, 9)
}
