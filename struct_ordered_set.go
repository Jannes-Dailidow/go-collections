package collections

// Len returns the number of elements. O(1), and safe on a nil receiver.
func (s *OrderedSet[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.order)
}

// IsEmpty reports whether the set has no elements. O(1), and safe on a nil receiver.
func (s *OrderedSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Has reports whether value is present. O(1), and safe on a nil receiver.
func (s *OrderedSet[T]) Has(value T) bool {
	if s == nil {
		return false
	}
	_, ok := s.index[value]
	return ok
}

// IndexOf returns the position of value in first-seen order, or -1. O(1) through the
// internal index.
func (s *OrderedSet[T]) IndexOf(value T) int {
	if s == nil {
		return -1
	}
	if pos, ok := s.index[value]; ok {
		return pos
	}
	return -1
}

// At returns the element at position n, or nil when n is out of range. O(1), and safe
// on a nil receiver.
//
// The element is a copy, for the same reason as [OrderedMap.At]: writing through a
// pointer into the order would leave the index describing a value that is no longer
// there.
func (s *OrderedSet[T]) At(n int) *T {
	if s == nil || n < 0 || n >= len(s.order) {
		return nil
	}
	result := s.order[n]
	return &result
}

// Add stores value and reports whether it was new. An existing value keeps its
// original position.
//
// This allocates on first use, so &OrderedSet[T]{} is ready to Add to. A nil pointer
// is not, and panics.
func (s *OrderedSet[T]) Add(value T) bool {
	if s.index == nil {
		s.index = make(map[T]int)
	}
	if _, ok := s.index[value]; ok {
		return false
	}
	s.index[value] = len(s.order)
	s.order = append(s.order, value)
	return true
}

// Remove drops value and reports whether it was present.
//
// O(n), for the same reason as [OrderedMap.Delete]: the positions after it all move.
func (s *OrderedSet[T]) Remove(value T) bool {
	if s == nil {
		return false
	}
	pos, ok := s.index[value]
	if !ok {
		return false
	}
	s.order = append(s.order[:pos], s.order[pos+1:]...)
	delete(s.index, value)
	for n := pos; n < len(s.order); n++ {
		s.index[s.order[n]] = n
	}
	return true
}

// Clear removes every element.
func (s *OrderedSet[T]) Clear() {
	if s == nil {
		return
	}
	s.order = nil
	clear(s.index)
}

// Clone returns a copy with the same order.
func (s *OrderedSet[T]) Clone() *OrderedSet[T] {
	if s == nil {
		return nil
	}
	result := &OrderedSet[T]{
		order: make([]T, len(s.order)),
		index: make(map[T]int, len(s.index)),
	}
	copy(result.order, s.order)
	for value, pos := range s.index {
		result.index[value] = pos
	}
	return result
}
