package action

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func onAudit(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeAudit, pm.Prj, pm.Cfg, pm.Logger)

	var err error
	switch f := pm.Cfg.GetStringOpt("fix"); f {
	case "remove":
		err = auditRemove(ctx, pm.Cfg.GetStringOpt("file"), pm, ret)
	case "header":
		err = auditHeader(ctx, pm.Cfg.GetStringOpt("file"), pm, ret)
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
	fs, _ := pm.PSvc.GetFilesystem(pm.Prj)
	errs := project.Validate(pm.Prj, fs, pm.MSvc.Deps(), pm.MSvc.Dangerous())
	lo.ForEach(errs, func(err *project.ValidationError, _ int) {
		ret = ret.WithError(errors.Errorf("%s: %s", err.Code, err.Message))
	})

	return ret
}

func auditRun(pm *PrjAndMods, ret *Result) error {
	timer := util.TimerStart()
	tgt, err := pm.PSvc.GetFilesystem(pm.Prj)
	if err != nil {
		return err
	}
	generated, err := getGeneratedFiles(tgt, pm.Prj.Ignore, pm.Logger)
	if err != nil {
		return err
	}
	src, err := getModuleFiles(pm)
	if err != nil {
		return err
	}

	audits := lo.FilterMap(generated, func(g string, _ int) (*diff.Diff, bool) {
		if !lo.Contains(src, g) && !isBad(g) {
			return &diff.Diff{Path: g, Status: diff.StatusDifferent}, true
		}
		return nil, false
	})
	empty, err := getEmptyFolders(tgt, pm.Prj.Ignore, pm.Logger)
	if err != nil {
		return err
	}
	lo.ForEach(empty, func(e string, _ int) {
		audits = append(audits, &diff.Diff{Path: e, Status: diff.StatusMissing})
	})

	if pm.EArgs != nil {
		lo.ForEach(pm.EArgs.Models, func(m *model.Model, _ int) {
			names := m.Columns.Names()
			if uniq := lo.Uniq(names); len(uniq) != len(names) {
				ret.Errors = append(ret.Errors, fmt.Sprintf("model [%s] contains duplicate column names", m.Name))
			}
		})
		lo.ForEach(pm.EArgs.Enums, func(e *enum.Enum, _ int) {
			keys := e.Values.Keys()
			if uniq := lo.Uniq(keys); len(uniq) != len(keys) {
				ret.Errors = append(ret.Errors, fmt.Sprintf("enum [%s] contains duplicate keys", e.Name))
			}
		})
	}

	sch, err := schemaCheck(pm)
	if err != nil {
		return err
	}
	ret.Errors = append(ret.Errors, sch...)
	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: audits, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return nil
}

var badFiles = []string{"client.js.map", "client.css.map", "file/header.go"}

func isBad(s string) bool {
	for _, bad := range badFiles {
		if strings.HasSuffix(s, bad) {
			return true
		}
	}
	return false
}
