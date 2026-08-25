package collections

// Number is the set of types Sum and SumBy accumulate over. Complex types are left
// out on purpose: they are not cmp.Ordered, so they cannot be used with MinBy and
// MaxBy, and admitting them here alone would be inconsistent.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
