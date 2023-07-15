package action

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/build"
)

func onCreate(ctx context.Context, params *Params) *Result {
	path := params.Cfg.GetStringOpt("path")
	if path == "" {
		path = "."
	}

	prj, _ := params.PSvc.Get(params.ProjectKey)
	if prj == nil {
		prj = project.NewProject(params.ProjectKey, path)
	}
	prj = projectFromCfg(prj, params.Cfg)
	ret := newResult(TypeCreate, prj, params.Cfg, params.Logger)
	if params.CLI {
		err := cliProject(ctx, prj, params.MSvc.Keys(), params.Logger)
		if err != nil {
			return ret.WithError(err)
		}
	}

	params.ProjectKey = prj.Key

	params.Logger.Info("Saving project...")
	err := params.PSvc.Save(prj, params.Logger)
	if err != nil {
		return ret.WithError(err)
	}

	params.Logger.Info("Generating project...")
	_, err = params.PSvc.Refresh(params.Logger)
	if err != nil {
		msg := fmt.Sprintf("unable to load newly created project from path [%s]", path)
		return errorResult(errors.Wrap(err, msg), TypeCreate, params.Cfg, params.Logger)
	}

	ctx, pm, err := getPrjAndMods(ctx, params)
	if err != nil {
		return errorResult(err, TypeCreate, params.Cfg, params.Logger)
	}
	retS := onGenerate(pm)
	ret = ret.Merge(retS)
	if ret.HasErrors() {
		return ret
	}

	params.Logger.Info("Building project...")
	ret.Modules = nil
	fullBuild(ctx, prj, ret, params.Logger)
	params.Logger.Info("Your project is ready! Run [bin/dev." + build.ScriptExtension + "] to start your application")
	return ret
}
