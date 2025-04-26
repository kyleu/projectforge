package util

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func AsyncCollect[T any, R any](items []T, f func(x T) (R, error), loggers ...Logger) ([]R, []error) {
	ret := make([]R, 0, len(items))
	var errs []error
	mu := sync.Mutex{}
	size := len(items)
	wg := sync.WaitGroup{}
	wg.Add(size)
	var processed int

	lo.ForEach(items, func(x T, _ int) {
		i := x
		go func() {
			defer wg.Done()
			r, err := f(i)
			mu.Lock()
			defer mu.Unlock()
			processed++
			if err == nil {
				ret = append(ret, r)
			} else {
				errs = append(errs, errors.Wrapf(err, "error running async function for item [%v]", i))
			}
			for _, logger := range loggers {
				logger.Debugf("processed [%d/%d] items", processed, size)
			}
		}()
	})
	wg.Wait()
	return ret, errs
}

func AsyncCollectMap[T any, K comparable, R any](items []T, k func(x T) K, f func(x T) (R, error), loggers ...Logger) (map[K]R, map[K]error) {
	ret := make(map[K]R, len(items))
	errs := map[K]error{}
	mu := sync.Mutex{}
	size := len(items)
	wg := sync.WaitGroup{}
	wg.Add(size)
	var processed int

	lo.ForEach(items, func(x T, _ int) {
		i := x
		go func() {
			defer wg.Done()
			key := k(i)
			r, err := f(i)
			mu.Lock()
			defer mu.Unlock()
			processed++
			if err == nil {
				ret[key] = r
			} else {
				errs[key] = errors.Wrapf(err, "error running async function for item [%v]", key)
			}
			for _, logger := range loggers {
				logger.Debugf("processed [%d/%d] items", processed, size)
			}
		}()
	})
	wg.Wait()
	return ret, errs
}

func AsyncRateLimit[T any, R any](key string, items []T, f func(x T) (R, error), maxConcurrent int, timeout time.Duration, loggers ...Logger) ([]R, []error) {
	ret := make([]R, 0, len(items))
	errs := make([]error, 0)
	mu := sync.Mutex{}
	size := len(items)
	wg := sync.WaitGroup{}
	wg.Add(size)
	var processed int
	var started int
	prefix := fmt.Sprintf("[%s] ", key)
	log := func(msg string, args ...any) {
		for _, logger := range loggers {
			logger.Debugf(prefix+msg, args...)
		}
	}

	limit := make(chan struct{}, maxConcurrent)
	defer close(limit)
	log("starting to process [%d] items, [%d] at once)", size, maxConcurrent)

	for idx, item := range items {
		select {
		case limit <- struct{}{}:
			go func(item T, idx int) {
				defer wg.Done()
				defer func() { <-limit }()
				str := fmt.Sprint(item)
				if str == "" {
					str = "item"
				}
				started++
				log("[%d/%d] starting to process [%s]...", started, size, str)
				r, err := f(item)
				mu.Lock()
				defer mu.Unlock()
				processed++
				if err == nil {
					ret = append(ret, r)
				} else {
					errs = append(errs, errors.Wrapf(err, "error running async function for item [%v]", item))
				}
				log("[%d/%d] processing [%s] complete", processed, size, str)
			}(item, idx)
		case <-time.After(timeout):
			errs = append(errs, errors.Errorf("job timed out after [%v]", timeout))
			return ret, errs
		}
	}

	wg.Wait()
	return ret, errs
}
