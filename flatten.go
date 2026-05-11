package slicest

func Flatten[T any, S ~[][]T](ss S) []T {
	var length int
	for _, s := range ss {
		length += len(s)
	}

	result := make([]T, 0, length)
	for _, s := range ss {
		result = append(result, s...)
	}
	return result
}
