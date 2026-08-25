package collections

import "maps"

// Len returns the number of elements. O(1).
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty reports whether the set has no elements. O(1).
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Has reports whether value is present. O(1), and the whole point of the type.
func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

// Add stores value and reports whether it was new.
//
// Like the builtin assignment it wraps, this panics on a nil set. Build one with
// make(Set[T]) or [CollectSet] before writing to it.
func (s Set[T]) Add(value T) bool {
	if _, ok := s[value]; ok {
		return false
	}
	s[value] = struct{}{}
	return true
}

// Remove drops value and reports whether it was present.
func (s Set[T]) Remove(value T) bool {
	if _, ok := s[value]; !ok {
		return false
	}
	delete(s, value)
	return true
}

// Clear removes every element, keeping the allocated set.
func (s Set[T]) Clear() {
	clear(s)
}

// Clone returns a copy.
func (s Set[T]) Clone() Set[T] {
	return Set[T](maps.Clone(s))
}

// NewOrderedSet returns an empty [OrderedSet] ready to Add to. The zero value
// &OrderedSet[T]{} works too, since Add allocates on first use; this is just the
// spelling that says so.
func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{index: make(map[T]int)}
}

// NewOrderedMap returns an empty [OrderedMap] ready to Put to.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{index: make(map[K]int)}
}
