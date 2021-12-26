package cmd

import (
	"context"
	"strings"

	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/export"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
)

func runToCompletion(ctx context.Context, projectKey string, t action.Type, cfg util.ValueMap) *action.Result {
	fs := filesystem.NewFileSystem(_flags.ConfigDir, _logger)
	mSvc := module.NewService(fs, _logger)
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
