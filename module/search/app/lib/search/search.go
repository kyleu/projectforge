package search

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/search/result"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Provider func(context.Context, *Params, *app.State, *cutil.PageState, util.Logger) (result.Results, error)

func Search(ctx context.Context, params *Params, as *app.State, page *cutil.PageState) (result.Results, []error) {
	ctx, span, logger := telemetry.StartSpan(ctx, "search", page.Logger)
	defer span.Complete()

	if params.Q == "" {
		return nil, nil
	}
	var allProviders []Provider
	// $PF_SECTION_START(search_functions)$
	// add your custom search functions here
	testFunc := func(ctx context.Context, p *Params, as *app.State, page *cutil.PageState, logger util.Logger) (result.Results, error) {
		if p.Q == "test" {
			return result.Results{{URL: "/about", Title: "Test Result", Icon: "star", Matches: nil}}, nil
		}
		return nil, nil
	}
	allProviders = append(allProviders, testFunc)
	// $PF_SECTION_END(search_functions)${{{ if .HasModule "export" }}}

	allProviders = append(allProviders, generatedSearch()...){{{ end }}}
	if len(allProviders) == 0 {
		return nil, []error{errors.New("no search providers configured")}
	}

	params.Q = strings.TrimSpace(params.Q)
	results, errs := util.AsyncCollect(allProviders, func(p Provider) ([]*result.Result, error) {
		return p(ctx, params, as, page, logger)
	})
	var ret result.Results = lo.Flatten(results)
	return ret.Sort(), errs
}
