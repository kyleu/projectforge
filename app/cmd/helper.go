package cmd

import (
	"context"
	"strings"

	"projectforge.dev/app/export"
	"projectforge.dev/app/lib/filesystem"
	"projectforge.dev/app/module"
	"projectforge.dev/app/project"
	"projectforge.dev/app/util"
	"projectforge.dev/app/action"
)

func runToCompletion(ctx context.Context, projectKey string, t action.Type, cfg util.ValueMap) *action.Result {
	fs := filesystem.NewFileSystem(_flags.ConfigDir, _logger)
	mSvc := module.NewService(ctx, fs, _logger)
	pSvc := project.NewService(_logger)
	eSvc := export.NewService(_logger)
	logger := _logger.With("service", "runner")
	p := &action.Params{Span: nil, ProjectKey: projectKey, T: t, Cfg: cfg, MSvc: mSvc, PSvc: pSvc, ESvc: eSvc, CLI: true, Logger: logger}
	return action.Apply(ctx, p)
}

func extractConfig(args []string) ([]string, util.ValueMap) {
	var retArgs []string
	retMap := util.ValueMap{}
	for _, arg := range args {
		l, r := util.StringSplit(arg, '=', true)
		l = strings.TrimSpace(l)
		r = strings.TrimSpace(r)
		if r == "" {
			retArgs = append(retArgs, l)
		} else {
			retMap[l] = r
		}
	}
	return retArgs, retMap
}
