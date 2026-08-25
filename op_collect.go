package collections

// SetBy collects the stream into a [Set] keyed by fn. Terminal.
//
// This is the fluent form of [CollectSet], and the reason it exists: Slice[T] and
// Iter[T] declare T as any, so neither can grow a Set() of its own elements. Passing
// a key function moves the comparable constraint onto the key, where a method may
// introduce it.
//
//	s.Values().SetBy(func(t T) T { return t })   // the identity case
//	CollectSet(s.Values())                       // the same thing, as a function
func (i Iter[T]) SetBy[K comparable](fn func(T) K) Set[K] {
	result := make(Set[K])
	for t := range i {
		result[fn(t)] = struct{}{}
	}
	return result
}

// MapBy collects the stream into a [Map], with fn splitting each element into a key
// and a value. A repeated key keeps the last value. Terminal.
func (i Iter[T]) MapBy[K comparable, V any](fn func(T) (K, V)) Map[K, V] {
	result := make(Map[K, V])
	for t := range i {
		k, v := fn(t)
		result[k] = v
	}
	return result
}

// SetBy collects the stream into a [Set] keyed by fn, or returns any abort. Terminal.
func (i IterX[T]) SetBy[K comparable](fn func(T) K) (Set[K], error) {
	result := make(Set[K])
	err := i.each(func(t T) (bool, error) {
		result[fn(t)] = struct{}{}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MapBy collects the stream into a [Map], or returns any abort. Terminal.
func (i IterX[T]) MapBy[K comparable, V any](fn func(T) (K, V)) (Map[K, V], error) {
	result := make(Map[K, V])
	err := i.each(func(t T) (bool, error) {
		k, v := fn(t)
		result[k] = v
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SetBy collects the slice into a [Set] keyed by fn, without going through an
// iterator.
func (s Slice[T]) SetBy[K comparable](fn func(T) K) Set[K] {
	result := make(Set[K], len(s))
	for _, t := range s {
		result[fn(t)] = struct{}{}
	}
	return result
}

// MapBy collects the slice into a [Map], with fn splitting each element into a key
// and a value. A repeated key keeps the last value. It is what v0.1 called ToMap.
func (s Slice[T]) MapBy[K comparable, V any](fn func(T) (K, V)) Map[K, V] {
	result := make(Map[K, V], len(s))
	for _, t := range s {
		k, v := fn(t)
		result[k] = v
	}
	return result
}
