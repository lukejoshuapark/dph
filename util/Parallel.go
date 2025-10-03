package util

import (
	"errors"
	"sync"
)

func Parallel[T any](inputs []T, f func(T) error) error {
	errc := make(chan error, len(inputs))
	wg := &sync.WaitGroup{}
	wg.Add(len(inputs))

	go func() {
		wg.Wait()
		close(errc)
	}()

	for _, input := range inputs {
		go func(input T) {
			defer wg.Done()
			errc <- f(input)
		}(input)
	}

	errs := make([]error, 0, len(inputs))
	for err := range errc {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
