package util

func Filter[T any](arr []T, f func(T) bool) []T {
	var result []T

	for _, item := range arr {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}
