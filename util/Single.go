package util

import "errors"

var ErrNoElementFound = errors.New("no element found")
var ErrTooManyElements = errors.New("multiple elements found")

func Single[T any](arr []T, f func(T) bool) (T, error) {
	ind := -1

	for i, v := range arr {
		if f(v) {
			if ind >= 0 {
				var zero T
				return zero, ErrTooManyElements
			}

			ind = i
		}
	}

	if ind < 0 {
		var zero T
		return zero, ErrNoElementFound
	}

	return arr[ind], nil
}
