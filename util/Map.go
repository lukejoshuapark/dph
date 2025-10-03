package util

func Map[T any, U any](arr []T, f func(T) U) []U {
	r := make([]U, len(arr))

	for i, v := range arr {
		r[i] = f(v)
	}

	return r
}
