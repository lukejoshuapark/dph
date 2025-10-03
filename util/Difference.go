package util

func Difference[T comparable](a, b []T) []T {
	m := make(map[T]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}

	diff := make([]T, 0, len(a))
	seen := make(map[T]struct{}, len(a))
	for _, v := range a {
		if _, inB := m[v]; !inB {
			if _, already := seen[v]; !already {
				diff = append(diff, v)
				seen[v] = struct{}{}
			}
		}
	}

	return diff
}
