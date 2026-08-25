package collections

// Len returns the number of pairs. O(1), and safe on a nil receiver.
func (m *OrderedMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.order)
}

// IsEmpty reports whether the map has no pairs. O(1), and safe on a nil receiver.
func (m *OrderedMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Has reports whether the key is present. O(1), and safe on a nil receiver.
func (m *OrderedMap[K, V]) Has(key K) bool {
	if m == nil {
		return false
	}
	_, ok := m.index[key]
	return ok
}

// Get returns the value for key and whether it was present. O(1), and safe on a nil
// receiver.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	pos, ok := m.index[key]
	if !ok {
		var zero V
		return zero, false
	}
	return m.order[pos].Value, true
}

// GetOr returns the value for key, or fallback when the key is absent.
func (m *OrderedMap[K, V]) GetOr(key K, fallback V) V {
	if value, ok := m.Get(key); ok {
		return value
	}
	return fallback
}

// IndexOf returns the position of key in the insertion order, or -1. O(1) through the
// internal index, which is why this is here and not on an iterator.
func (m *OrderedMap[K, V]) IndexOf(key K) int {
	if m == nil {
		return -1
	}
	if pos, ok := m.index[key]; ok {
		return pos
	}
	return -1
}

// At returns the pair at position n, or nil when n is out of range. O(1), and safe on
// a nil receiver.
//
// The Entry is a copy. Unlike [Slice.At], this cannot hand out a pointer into the
// order: writing a new Key through it would leave the index describing a key that is
// no longer there.
func (m *OrderedMap[K, V]) At(n int) *Entry[K, V] {
	if m == nil || n < 0 || n >= len(m.order) {
		return nil
	}
	result := m.order[n]
	return &result
}

// Put stores value under key. An existing key keeps its position and takes the new
// value, matching [CollectOrderedMap].
//
// This allocates on first use, so &OrderedMap[K, V]{} is ready to Put to. A nil
// pointer is not, and panics.
func (m *OrderedMap[K, V]) Put(key K, value V) {
	if m.index == nil {
		m.index = make(map[K]int)
	}
	if pos, ok := m.index[key]; ok {
		m.order[pos].Value = value
		return
	}
	m.index[key] = len(m.order)
	m.order = append(m.order, Entry[K, V]{Key: key, Value: value})
}

// Delete removes key and reports whether it was present.
//
// O(n): every later key moves down one position, and its index entry has to move with
// it. Deleting in a loop over a large map is quadratic; rebuild with a Filter and
// [CollectOrderedMap] instead.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	if m == nil {
		return false
	}
	pos, ok := m.index[key]
	if !ok {
		return false
	}
	m.order = append(m.order[:pos], m.order[pos+1:]...)
	delete(m.index, key)
	for n := pos; n < len(m.order); n++ {
		m.index[m.order[n].Key] = n
	}
	return true
}

// Clear removes every pair.
func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		return
	}
	m.order = nil
	clear(m.index)
}

// Clone returns a copy with the same order. Shallow: the values themselves are not
// copied.
func (m *OrderedMap[K, V]) Clone() *OrderedMap[K, V] {
	if m == nil {
		return nil
	}
	result := &OrderedMap[K, V]{
		order: make([]Entry[K, V], len(m.order)),
		index: make(map[K]int, len(m.index)),
	}
	copy(result.order, m.order)
	for key, pos := range m.index {
		result.index[key] = pos
	}
	return result
}
