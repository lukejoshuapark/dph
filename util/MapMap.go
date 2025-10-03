package util

func MapMap[T comparable, U any, V any](m map[T]U, f func(T, U) V) []V {
	r := make([]V, 0, len(m))

	for k, v := range m {
		r = append(r, f(k, v))
	}

	return r
}
