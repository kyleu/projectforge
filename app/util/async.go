package util

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func AsyncCollect[T any, R any](items []T, f func(x T) (R, error)) ([]R, []error) {
	ret := make([]R, 0, len(items))
	var errs []error
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	lo.ForEach(items, func(x T, _ int) {
		i := x
		go func() {
			r, err := f(i)
			mu.Lock()
			if err == nil {
				ret = append(ret, r)
			} else {
				errs = append(errs, errors.Wrapf(err, "error running async function for item [%v]", i))
			}
			mu.Unlock()
			wg.Done()
		}()
	})
	wg.Wait()
	return ret, errs
}

func AsyncCollectMap[T any, K comparable, R any](items []T, k func(x T) K, f func(x T) (R, error)) (map[K]R, map[K]error) {
	ret := make(map[K]R, len(items))
	errs := map[K]error{}
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	lo.ForEach(items, func(x T, _ int) {
		i := x
		go func() {
			key := k(i)
			r, err := f(i)
			mu.Lock()
			if err == nil {
				ret[key] = r
			} else {
				errs[key] = errors.Wrapf(err, "error running async function for item [%v]", key)
			}
			mu.Unlock()
			wg.Done()
		}()
	})
	wg.Wait()
	return ret, errs
}

func AsyncRateLimit[T any, R any](items []T, f func(x T) (R, error), maxConcurrent int, timeout time.Duration) ([]R, []error) {
	ret := make([]R, 0, len(items))
	errs := make([]error, 0)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	limit := make(chan struct{}, maxConcurrent)
	defer close(limit)

	for idx, item := range items {
		select {
		case limit <- struct{}{}:
			wg.Add(1)
			go func(item T, idx int) {
				defer wg.Done()
				defer func() { <-limit }()

				r, err := f(item)
				mu.Lock()
				defer mu.Unlock()
				if err == nil {
					ret = append(ret, r)
				} else {
					errs = append(errs, errors.Wrapf(err, "error running async function for item [%v]", item))
				}
			}(item, idx)
		case <-time.After(timeout):
			errs = append(errs, errors.Errorf("job timed out after [%v]", timeout))
			return ret, errs
		}
	}

	wg.Wait()
	return ret, errs
}
