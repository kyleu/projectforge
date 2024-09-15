package action

import (
	"context"
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

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
	lo.ForEach(empty, func(e string, _ int) {
		audits = append(audits, &diff.Diff{Path: e, Status: diff.StatusMissing})
	})

	if pm.Prj.ExportArgs != nil {
		lo.ForEach(pm.Prj.ExportArgs.Models, func(m *model.Model, _ int) {
			names := m.Columns.Names()
			if uniq := lo.Uniq(names); len(uniq) != len(names) {
				ret.Errors = append(ret.Errors, fmt.Sprintf("model [%s] contains duplicate column names", m.Name))
			}
		})
		lo.ForEach(pm.Prj.ExportArgs.Enums, func(e *enum.Enum, _ int) {
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
	mr := &module.Result{Keys: pm.Mods.Keys(), Status: util.OK, Diffs: audits, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return nil
}
