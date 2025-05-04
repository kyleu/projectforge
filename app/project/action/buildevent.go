package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

func onBuild(ctx context.Context, pm *PrjAndMods) *Result {
	phaseStr := pm.Cfg.GetStringOpt("phase")
	if phaseStr == "" {
		if x, _ := pm.Cfg.GetStringArray("cmds", true); len(x) == 1 {
			phaseStr = x[0]
		} else {
			phaseStr = "build"
		}
	}

	_, span, logger := telemetry.StartSpan(ctx, "build:"+phaseStr, pm.Logger)
	span.Attribute("project", pm.Prj.Key)
	defer span.Complete()
	pm.Logger = logger

	ret := newResult(TypeBuild, pm.Prj, pm.Cfg, logger)
	ret.AddLog("building project [%s] in [%s] with phase [%s]", pm.Prj.Key, pm.Prj.Path, phaseStr)
	phase := AllBuilds.Get(phaseStr)
	if phase == nil {
		return ret.WithError(errors.Errorf("invalid phase [%s]", phaseStr))
	}
	timer := util.TimerStart()
	ret = phase.Run(ctx, pm, ret)
	ret.Duration = timer.End()
	return ret
}

func onDeps(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
	if up := ret.Args.GetStringOpt("upgrade"); up != "" {
		o := ret.Args.GetStringOpt("old")
		n := ret.Args.GetStringOpt("new")
		err := build.OnDepsUpgrade(ctx, pm.Prj, up, o, n, pm.PSvc, pm.Logger)
		if err != nil {
			return ret.WithError(err)
		}
	}
	upd := ret.Args.GetStringOpt("upd")
	showAll := pm.Cfg.GetBoolOpt("showAll")
	deps, err := build.LoadDeps(ctx, pm.Prj.Key, pm.Prj.Path, upd != util.BoolFalse, pm.FS, showAll, pm.Logger)
	ret.Data = deps
	if err != nil {
		return ret.WithError(err)
	}
	return ret
}

func onImports(_ context.Context, pm *PrjAndMods, r *Result) *Result {
	fixStr := r.Args.GetStringOpt("fix")
	fix := fixStr == util.BoolTrue
	fileStr := r.Args.GetStringOpt("file")
	t := util.TimerStart()
	logs, diffs, err := build.Imports(pm.Prj.Package, fix, fileStr, pm.FS, pm.Logger)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"imports"}, Status: util.OK, Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onIgnored(_ context.Context, pm *PrjAndMods, r *Result) *Result {
	ign, err := build.Ignored(pm.Prj, pm.FS, pm.Logger)
	r.Data = ign
	if err != nil {
		return r.WithError(err)
	}
	res := &module.Result{Keys: []string{"ignored"}, Status: util.OK}
	lo.ForEach(ign, func(x string, _ int) {
		if x != "app/file/header.go" && x != "doc/faq.md" {
			res.Diffs = append(res.Diffs, &diff.Diff{Path: x, Status: diff.StatusDifferent})
		}
	})
	r.Modules = append(r.Modules, res)
	return r
}

func onPackages(_ context.Context, pm *PrjAndMods, r *Result) *Result {
	pkgs, err := build.Packages(pm.Prj, pm.FS, pm.Cfg.GetBoolOpt("all"), pm.Logger)
	r.Data = pkgs
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onCleanup(_ context.Context, pm *PrjAndMods, r *Result) *Result {
	t := util.TimerStart()
	logs, diffs, err := build.Cleanup(pm.FS, pm.Logger)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"cleanup"}, Status: util.OK, Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onSize(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	t := util.TimerStart()
	d, logs, err := build.Size(ctx, pm.FS, r.Args.GetStringOpt("exec"), pm.Logger)
	r.Data = d
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"size"}, Status: util.OK, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onClientInstall(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
	return simpleProc(ctx, "npm install", util.StringFilePath(pm.Prj.Path, "client"), ret, pm.Logger)
}

func onFull(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	t := util.TimerStart()
	logs, err := build.Full(ctx, pm.Prj, pm.Logger)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"fullbuild"}, Status: util.OK, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onDeployments(_ context.Context, pm *PrjAndMods, r *Result) *Result {
	deploys := lo.Reject(pm.Prj.Info.Deployments, func(x string, _ int) bool {
		return strings.HasPrefix(x, "ts:")
	})
	if pm.Prj.Info == nil || len(deploys) == 0 {
		return r
	}
	t := util.TimerStart()
	fix := pm.Cfg.GetBoolOpt("fix")
	logs, diffs, err := build.Deployments(pm.Prj.Version, pm.FS, fix, pm.Cfg.GetStringOpt("path"), deploys)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"deployments"}, Status: util.OK, Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onCoverage(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	ret, logs, err := runCoverage(ctx, pm.FS, r.Args.GetStringOpt("scope"), pm.Logger)
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	r.Data = ret
	return r
}
