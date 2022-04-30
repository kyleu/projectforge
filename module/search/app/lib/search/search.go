package search

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/search/result"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Provider func(context.Context, *app.State, *Params, *zap.SugaredLogger) (result.Results, error)

func Search(ctx context.Context, as *app.State, params *Params) (result.Results, []error) {
	ctx, span, logger := telemetry.StartSpan(ctx, "search", as.Logger)
	defer span.Complete()

	if params.Q == "" {
		return nil, nil
	}
	var allProviders []Provider
	// $PF_SECTION_START(search_functions)$
	testFunc := func(ctx context.Context, as *app.State, p *Params, logger *zap.SugaredLogger) (result.Results, error) {
		return result.Results{{URL: "/search?q=test", Title: "Test Result", Icon: "star", Matches: nil}}, nil
	}
	allProviders = append(allProviders, testFunc)
	// $PF_SECTION_END(search_functions)$
	// $PF_INJECT_START(codegen)$
	// $PF_INJECT_END(codegen)$
	if len(allProviders) == 0 {
		return nil, []error{errors.New("no search providers configured")}
	}

	params.Q = strings.TrimSpace(params.Q)

	results, errs := util.AsyncCollect(allProviders, func(item Provider) (result.Results, error) {
		return item(ctx, as, params, logger)
	})

	ret := make(result.Results, 0, len(results)*len(results))
	for _, x := range results {
		ret = append(ret, x...)
	}

	ret.Sort()
	return ret, errs
}
