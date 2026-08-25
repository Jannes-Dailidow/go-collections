package collections

// Entry is a key/value pair. It carries a pair stream into the single-value layer,
// and [OrderedMap] stores its order as a slice of these.
type Entry[K, V any] struct {
	Key   K
	Value V
}
