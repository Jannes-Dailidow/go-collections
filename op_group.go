package collections

// GroupBy collects the stream into a [Map] of slices, keyed by fn. Terminal.
//
// It is a function rather than a method, and not by choice: a method on Iter[T]
// returning Map[K, Slice[T]] is an instantiation cycle, because Map's own Values
// method hands back Iter[Slice[T]], whose GroupBy hands back Map[K, Slice[Slice[T]]],
// and so on. See the Technical decisions section of the README.
func GroupBy[T any, K comparable](i Iter[T], fn func(T) K) Map[K, Slice[T]] {
	result := make(Map[K, Slice[T]])
	for t := range i {
		key := fn(t)
		result[key] = append(result[key], t)
	}
	return result
}

// GroupByX is the fallible counterpart of [GroupBy].
func GroupByX[T any, K comparable](i IterX[T], fn func(T) K) (Map[K, Slice[T]], error) {
	result := make(Map[K, Slice[T]])
	err := i.each(func(t T) (bool, error) {
		key := fn(t)
		result[key] = append(result[key], t)
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GroupByOrdered collects the stream into an [OrderedMap] of slices, keyed by fn, with
// the groups in the order their keys were first seen. Terminal.
func GroupByOrdered[T any, K comparable](i Iter[T], fn func(T) K) *OrderedMap[K, Slice[T]] {
	result := NewOrderedMap[K, Slice[T]]()
	for t := range i {
		key := fn(t)
		group, _ := result.Get(key)
		result.Put(key, append(group, t))
	}
	return result
}

// GroupByOrderedX is the fallible counterpart of [GroupByOrdered].
func GroupByOrderedX[T any, K comparable](i IterX[T], fn func(T) K) (*OrderedMap[K, Slice[T]], error) {
	result := NewOrderedMap[K, Slice[T]]()
	err := i.each(func(t T) (bool, error) {
		key := fn(t)
		group, _ := result.Get(key)
		result.Put(key, append(group, t))
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Partition splits the stream into the elements fn accepts and the rest. It is the
// two-bucket case of [GroupBy], without needing a key. Terminal.
func (i Iter[T]) Partition(fn func(T) bool) (accepted, rejected Slice[T]) {
	for t := range i {
		if fn(t) {
			accepted = append(accepted, t)
			continue
		}
		rejected = append(rejected, t)
	}
	return accepted, rejected
}

// PartitionX splits the stream into the elements fn accepts and the rest, and aborts
// on the first error. Terminal.
func (i Iter[T]) PartitionX(fn func(T) (bool, error)) (accepted, rejected Slice[T], err error) {
	return i.X().PartitionX(fn)
}

// Partition splits the stream into the elements fn accepts and the rest, plus any
// abort. Terminal.
func (i IterX[T]) Partition(fn func(T) bool) (accepted, rejected Slice[T], err error) {
	return i.PartitionX(func(t T) (bool, error) {
		return fn(t), nil
	})
}

// PartitionX splits the stream into the elements fn accepts and the rest. An error
// from fn aborts, as does an abort from upstream, and either discards both halves.
// Terminal.
func (i IterX[T]) PartitionX(fn func(T) (bool, error)) (accepted, rejected Slice[T], err error) {
	err = i.each(func(t T) (bool, error) {
		ok, e := fn(t)
		if e != nil {
			return false, e
		}
		if ok {
			accepted = append(accepted, t)
			return true, nil
		}
		rejected = append(rejected, t)
		return true, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return accepted, rejected, nil
}
