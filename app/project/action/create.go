package action

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

func onCreate(ctx context.Context, params *Params) *Result {
	path := util.OrDefault(params.Cfg.GetStringOpt("path"), ".")
	if params.CLI {
		if pth, err := promptString("Directory", "Choose the directory your project will live in", path, true); err != nil {
			return errorResult(err, params.T, params.Cfg, params.Logger)
		} else {
			path = pth
		}
	}

	_, _ = params.PSvc.Refresh(params.Logger)
	var prj *project.Project
	if params.ProjectKey == "" {
		prj = params.PSvc.ByPath(path)
		if prj != nil {
			prj.Name = strings.TrimSuffix(prj.Name, " (missing)")
		}
	} else {
		prj, _ = params.PSvc.Get(params.ProjectKey)
	}
	if prj == nil {
		prj = project.NewProject(util.Choose(params.ProjectKey == "", filepath.Base(path), params.ProjectKey), path)
	}
	if len(params.Cfg) > 0 {
		err := ProjectFromMap(prj, params.Cfg, true)
		if err != nil {
			return errorResult(err, params.T, params.Cfg, params.Logger)
		}
	}

	if params.CLI {
		err := cliProject(ctx, prj, params.MSvc.Modules(), params.Logger)
		if err != nil {
			return errorResult(err, params.T, params.Cfg, params.Logger)
		}
	}
	ret := newResult(TypeCreate, prj, params.Cfg, params.Logger)

	params.ProjectKey = prj.Key

	if errs := project.Validate(prj, nil, nil, nil); len(errs) > 0 {
		return ret.WithError(errors.New(errs.Error()))
	}

	params.Logger.Info("Saving project...")
	if err := params.PSvc.Save(prj, params.Logger); err != nil {
		return ret.WithError(err)
	}

	params.Logger.Info("Generating project...")
	if _, err := params.PSvc.Refresh(params.Logger); err != nil {
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
