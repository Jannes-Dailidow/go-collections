package collections

// Sum adds the elements together. Zero on an empty stream. Terminal.
//
// It is a function rather than a method because a method cannot tighten T to
// [Number]. For the fluent form use [Iter.SumBy] with an identity key.
func Sum[T Number](i Iter[T]) T {
	return i.SumBy(func(t T) T {
		return t
	})
}

// SumX is the fallible counterpart of [Sum].
func SumX[T Number](i IterX[T]) (T, error) {
	return i.SumBy(func(t T) T {
		return t
	})
}

// SumBy adds up fn applied to each element. Zero on an empty stream. Terminal.
//
// Where the value can fail to compute, MapX into the fallible family first and Sum
// that: i.MapX(parse) then SumX.
func (i Iter[T]) SumBy[N Number](fn func(T) N) N {
	result, _ := i.X().SumBy(fn)
	return result
}

// SumBy adds up fn applied to each pair. Zero on an empty stream. Terminal.
func (i Iter2[K, V]) SumBy[N Number](fn func(K, V) N) N {
	result, _ := i.X().SumBy(fn)
	return result
}

// SumBy adds up fn applied to each element, and returns any abort. Terminal.
func (i IterX[T]) SumBy[N Number](fn func(T) N) (N, error) {
	var result N
	err := i.each(func(t T) (bool, error) {
		result += fn(t)
		return true, nil
	})
	if err != nil {
		var zero N
		return zero, err
	}
	return result, nil
}

// SumBy adds up fn applied to each pair, and returns any abort. Terminal.
func (i Iter2X[K, V]) SumBy[N Number](fn func(K, V) N) (N, error) {
	var result N
	err := i.each(func(k K, v V) (bool, error) {
		result += fn(k, v)
		return true, nil
	})
	if err != nil {
		var zero N
		return zero, err
	}
	return result, nil
}
