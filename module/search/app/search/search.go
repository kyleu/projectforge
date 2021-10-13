package search

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
)

type Provider func(context.Context, *app.State, *Params) (Results, error)

func Search(ctx context.Context, as *app.State, params *Params) (Results, []error) {
	var allProviders []Provider
	// $PF_SECTION_START(search_functions)$
	testFunc := func(ctx context.Context, as *app.State, p *Params) (Results, error) {
		return Results{{URL: "/search?q=test", Title: "Test Result", Icon: "star", Match: p.Q + "!!!"}}, nil
	}
	allProviders = append(allProviders, testFunc)
	// $PF_SECTION_END(search_functions)$

	if len(allProviders) == 0 {
		return nil, []error{errors.New("no search providers configured")}
	}

	ret := Results{}
	var errs []error
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(allProviders))
	params.Q = strings.TrimSpace(params.Q)

	for _, p := range allProviders {
		f := p
		go func() {
			res, err := f(ctx, as, params)
			mu.Lock()
			if err != nil {
				errs = append(errs, err)
			}
			ret = append(ret, res...)
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	ret.Sort()
	return ret, errs
}
