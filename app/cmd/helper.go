package cmd

import (
	"context"
	"strings"

	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/log"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
)

func runToCompletion(ctx context.Context, projectKey string, t action.Type, cfg util.ValueMap) error {
	logger, err := log.InitLogging(true)
	if err != nil {
		return err
	}
	mSvc := module.NewService(logger)
	pSvc := project.NewService(logger)
	result := action.Apply(ctx, nil, projectKey, t, cfg, mSvc, pSvc, logger.With("service", "runner"))
	if len(result.Errors) > 0 {
		return result.AsError()
	}
	return nil
}

func extractConfig(args []string) ([]string, util.ValueMap) {
	var retArgs []string
	retMap := util.ValueMap{}
	for _, arg := range args {
		l, r := util.SplitString(arg, '=', true)
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
