package action

import (
	"context"
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type AuditResult struct {
	Stats *CodeStats `json:"stats"`
}

func onAudit(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeAudit, pm.Prj, pm.Cfg, pm.Logger)

	var err error
	switch f := pm.Cfg.GetStringOpt("fix"); f {
	case "remove":
		err = auditRemove(ctx, pm.Cfg.GetStringOpt("file"), pm, ret)
	case "":
		// noop, run normal audit
	default:
		return errorResult(errors.Errorf("invalid fix type [%s]", f), TypeAudit, pm.Cfg, pm.Logger)
	}
	if err != nil {
		return errorResult(err, TypeAudit, pm.Cfg, pm.Logger)
	}
	err = auditRun(pm, ret)
	if err != nil {
		return errorResult(err, TypeAudit, pm.Cfg, pm.Logger)
	}
	errs := project.Validate(pm.Prj, pm.FS, pm.MSvc.Deps(), pm.MSvc.Dangerous())
	lo.ForEach(errs, func(err *project.ValidationError, _ int) {
		ret = ret.WithError(errors.Errorf("%s: %s", err.Code, err.Message))
	})

	return ret
}

func auditRun(pm *PrjAndMods, ret *Result) error {
	timer := util.TimerStart()

	var audits diff.Diffs

	ign := slices.Clone(pm.Prj.Ignore)
	if pm.Prj.HasModule("notebook") {
		ign = append(ign, "notebook")
	}
	empty, err := getEmptyFolders(pm.FS, ign, pm.Logger)
	if err != nil {
		return err
	}
	audits = append(audits, lo.Map(empty, func(e string, _ int) *diff.Diff {
		return &diff.Diff{Path: e, Status: diff.StatusMissing}
	})...)

	if pm.Prj.ExportArgs != nil {
		exportCheck(pm, ret)
	}

	sch, err := schemaCheck(pm)
	if err != nil {
		return err
	}
	ret.Errors = append(ret.Errors, sch...)

	st, err := runCodeStats(pm.Prj.Path, false)
	if err != nil {
		return err
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: util.OK, Diffs: audits, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)

	ret.Data = &AuditResult{Stats: st}
	return nil
}

func exportCheck(pm *PrjAndMods, ret *Result) {
	lo.ForEach(pm.Prj.ExportArgs.Enums, func(e *enum.Enum, _ int) {
		keys := e.Values.Keys()
		if uniq := lo.Uniq(keys); len(uniq) != len(keys) {
			ret.Errors = append(ret.Errors, fmt.Sprintf("enum [%s] contains duplicate keys", e.Name))
		}
	})
	lo.ForEach(pm.Prj.ExportArgs.Events, func(e *model.Event, _ int) {
		names := e.Columns.Names()
		if uniq := lo.Uniq(names); len(uniq) != len(names) {
			ret.Errors = append(ret.Errors, fmt.Sprintf("event [%s] contains duplicate column names", e.Name))
		}
	})
	lo.ForEach(pm.Prj.ExportArgs.Models, func(m *model.Model, _ int) {
		names := m.Columns.Names()
		if uniq := lo.Uniq(names); len(uniq) != len(names) {
			ret.Errors = append(ret.Errors, fmt.Sprintf("model [%s] contains duplicate column names", m.Name))
		}
	})
}

func getEmptyFolders(tgt filesystem.FileLoader, ignore []string, logger util.Logger, pth ...string) ([]string, error) {
	ret := &util.StringSlice{}
	pStr := util.StringPath(pth...)
	fc := len(tgt.ListFiles(pStr, nil, logger))
	ds := tgt.ListDirectories(pStr, ignore, logger)
	if fc == 0 && len(ds) == 0 {
		ret.Push(pStr)
	}
	for _, d := range ds {
		p := append(slices.Clone(pth), d)
		childRes, err := getEmptyFolders(tgt, ignore, logger, p...)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get empty folders for [%s/%s]", util.StringPath(p...), d)
		}
		ret.Push(childRes...)
	}
	return ret.Slice, nil
}

func auditRemove(ctx context.Context, fn string, pm *PrjAndMods, ret *Result) error {
	if fn == "" {
		paths := ret.Modules.Paths(false)
		if len(paths) == 0 {
			x := &Result{}
			err := auditRun(pm, x)
			if err != nil {
				return err
			}
			paths = x.Modules.Paths(false)
		}
		for _, pth := range paths {
			err := auditRemove(ctx, pth, pm, ret)
			if err != nil {
				return errors.Wrapf(err, "can't fix audit of [%s]", pth)
			}
		}
		return nil
	}
	ret.AddLog("removed [%s]", fn)
	return pm.FS.Remove(fn, pm.Logger)
}
