package collections

import "maps"

// Len returns the number of pairs. O(1).
func (m Map[K, V]) Len() int {
	return len(m)
}

// IsEmpty reports whether the map has no pairs. O(1).
func (m Map[K, V]) IsEmpty() bool {
	return len(m) == 0
}

// Has reports whether the key is present. O(1), and the reason this is on the map
// rather than on an iterator, where the same question costs a scan.
func (m Map[K, V]) Has(key K) bool {
	_, ok := m[key]
	return ok
}

// Get returns the value for key and whether it was present.
func (m Map[K, V]) Get(key K) (V, bool) {
	value, ok := m[key]
	return value, ok
}

// GetOr returns the value for key, or fallback when the key is absent.
func (m Map[K, V]) GetOr(key K, fallback V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return fallback
}

// Put stores value under key, replacing whatever was there.
//
// Like the builtin assignment it wraps, this panics on a nil map. Build one with
// make(Map[K, V]) or a Collect function before writing to it.
func (m Map[K, V]) Put(key K, value V) {
	m[key] = value
}

// Delete removes key. Absent keys are not an error.
func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

// Clear removes every pair, keeping the allocated map.
func (m Map[K, V]) Clear() {
	clear(m)
}

// Clone returns a copy. Shallow: the values themselves are not copied.
func (m Map[K, V]) Clone() Map[K, V] {
	return Map[K, V](maps.Clone(m))
}

// Merge copies every pair from other into m, last write winning. Use [Map.MergeFunc]
// to decide collisions instead.
func (m Map[K, V]) Merge(other Map[K, V]) {
	maps.Copy(m, other)
}

// MergeFunc copies every pair from other into m, calling resolve with the existing
// and incoming values for each key already present.
func (m Map[K, V]) MergeFunc(other Map[K, V], resolve func(existing, incoming V) V) {
	for key, incoming := range other {
		if existing, ok := m[key]; ok {
			m[key] = resolve(existing, incoming)
			continue
		}
		m[key] = incoming
	}
}

// InvertBy turns the map inside out, keying it by fn applied to each pair. Colliding
// new keys keep the last one, in whatever order the map iterated, so invert on a key
// you know to be unique.
//
// It takes a key function rather than inverting V directly because V is declared as
// any and a method cannot tighten it to comparable.
func (m Map[K, V]) InvertBy[V2 comparable](fn func(K, V) V2) Map[V2, K] {
	result := make(Map[V2, K], len(m))
	for key, value := range m {
		result[fn(key, value)] = key
	}
	return result
}
