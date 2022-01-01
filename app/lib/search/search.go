package search

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/lib/search/result"
)

type Provider func(context.Context, *app.State, *Params) (result.Results, error)

func Search(ctx context.Context, as *app.State, params *Params) (result.Results, []error) {
	if params.Q == "" {
		return nil, nil
	}
	var allProviders []Provider
	// $PF_SECTION_START(search_functions)$
	projectFunc := func(ctx context.Context, as *app.State, p *Params) (result.Results, error) {
		return as.Services.Projects.Search(p.Q)
	}
	moduleFunc := func(ctx context.Context, as *app.State, p *Params) (result.Results, error) {
		return as.Services.Modules.Search(p.Q)
	}
	allProviders = append(allProviders, projectFunc, moduleFunc)
	// $PF_SECTION_END(search_functions)$
	// $PF_INJECT_START(codegen)$
	// $PF_INJECT_END(codegen)$

	if len(allProviders) == 0 {
		return nil, []error{errors.New("no search providers configured")}
	}

	ret := result.Results{}
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
