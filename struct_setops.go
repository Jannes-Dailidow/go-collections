package collections

// Set algebra, on the two types whose element parameter is already comparable. The
// [OrderedSet] versions keep the receiver's order first and append whatever the other
// side contributes, so the result is as predictable as the inputs.

// Union returns the elements in either set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := make(Set[T], len(s)+len(other))
	for value := range s {
		result[value] = struct{}{}
	}
	for value := range other {
		result[value] = struct{}{}
	}
	return result
}

// Intersect returns the elements in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	small, large := s, other
	if len(large) < len(small) {
		small, large = large, small
	}
	result := make(Set[T])
	for value := range small {
		if _, ok := large[value]; ok {
			result[value] = struct{}{}
		}
	}
	return result
}

// Diff returns the elements in s that are not in other.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	result := make(Set[T])
	for value := range s {
		if _, ok := other[value]; !ok {
			result[value] = struct{}{}
		}
	}
	return result
}

// SymDiff returns the elements in exactly one of the two sets. It is what the root
// package calls Unique.
func (s Set[T]) SymDiff(other Set[T]) Set[T] {
	result := make(Set[T])
	for value := range s {
		if _, ok := other[value]; !ok {
			result[value] = struct{}{}
		}
	}
	for value := range other {
		if _, ok := s[value]; !ok {
			result[value] = struct{}{}
		}
	}
	return result
}

// IsSubset reports whether every element of s is in other. True for an empty s.
func (s Set[T]) IsSubset(other Set[T]) bool {
	if len(s) > len(other) {
		return false
	}
	for value := range s {
		if _, ok := other[value]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset reports whether every element of other is in s.
func (s Set[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint reports whether the two sets share no element.
func (s Set[T]) IsDisjoint(other Set[T]) bool {
	small, large := s, other
	if len(large) < len(small) {
		small, large = large, small
	}
	for value := range small {
		if _, ok := large[value]; ok {
			return false
		}
	}
	return true
}

// Union returns the elements in either set, s's order first.
func (s *OrderedSet[T]) Union(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for value := range s.Values() {
		result.Add(value)
	}
	for value := range other.Values() {
		result.Add(value)
	}
	return result
}

// Intersect returns the elements in both sets, in s's order.
func (s *OrderedSet[T]) Intersect(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for value := range s.Values() {
		if other.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// Diff returns the elements in s that are not in other, in s's order.
func (s *OrderedSet[T]) Diff(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for value := range s.Values() {
		if !other.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// SymDiff returns the elements in exactly one of the two sets, s's contribution
// first.
func (s *OrderedSet[T]) SymDiff(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for value := range s.Values() {
		if !other.Has(value) {
			result.Add(value)
		}
	}
	for value := range other.Values() {
		if !s.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// IsSubset reports whether every element of s is in other. True for an empty s.
func (s *OrderedSet[T]) IsSubset(other *OrderedSet[T]) bool {
	if s.Len() > other.Len() {
		return false
	}
	for value := range s.Values() {
		if !other.Has(value) {
			return false
		}
	}
	return true
}

// IsSuperset reports whether every element of other is in s.
func (s *OrderedSet[T]) IsSuperset(other *OrderedSet[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint reports whether the two sets share no element.
func (s *OrderedSet[T]) IsDisjoint(other *OrderedSet[T]) bool {
	for value := range s.Values() {
		if other.Has(value) {
			return false
		}
	}
	return true
}
