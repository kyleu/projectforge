package action

import (
	"context"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

type Build struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitzero"`
	Description string `json:"description,omitzero"`
	Expensive   bool   `json:"expensive,omitzero"`

	Run func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result `json:"-"`
}

func simpleProc(ctx context.Context, cmd string, path string, ret *Result, logger util.Logger) *Result {
	exitCode, out, err := telemetry.RunProcessSimple(ctx, cmd, path, logger)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("build output for [%s]:\n%s", cmd, out)
	if exitCode != 0 {
		ret.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}
	return ret
}

func simpleBuild(key string, title string, cmd string, expensive bool) *Build {
	return &Build{Key: key, Title: title, Description: "Runs [" + cmd + "]", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, cmd, pm.Prj.Path, ret, pm.Logger)
	}, Expensive: expensive}
}

type Builds []*Build

func (b Builds) ForAllProjects() Builds {
	return lo.Filter(b, func(x *Build, _ int) bool {
		return x.Key != "lint"
	})
}

func (b Builds) Get(key string) *Build {
	return lo.FindOrElse(b, nil, func(x *Build) bool {
		return x.Key == key
	})
}

func fullBuild(ctx context.Context, prj *project.Project, r *Result, logger util.Logger) *Result {
	logs, err := build.Full(ctx, prj, logger)
	lo.ForEach(logs, func(l string, _ int) {
		r.Log(l)
	})
	if err != nil {
		return r.WithError(err)
	}
	return r
}
