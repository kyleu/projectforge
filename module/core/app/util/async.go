package util

import (
	"sync"

	"github.com/pkg/errors"
)

func AsyncCollect[T any, R any](items []T, f func(item T) (R, error)) ([]R, []error) {
	ret := make([]R, 0, len(items))
	var errs []error
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	for _, item := range items {
		item := item
		go func() {
			r, err := f(item)
			mu.Lock()
			if err == nil {
				ret = append(ret, r)
			} else {
				errs = append(errs, errors.Wrapf(err, "error running async function for item [%v]", item))
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return ret, errs
}
