package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file/diff"
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
	err = auditRun(ctx, pm, ret)
	if err != nil {
		return errorResult(err, TypeAudit, pm.Cfg, pm.Logger)
	}

	fs := pm.PSvc.GetFilesystem(pm.Prj)
	errs := project.Validate(pm.Prj, pm.MSvc.Deps(), fs)
	for _, err := range errs {
		ret = ret.WithError(errors.Errorf("%s: %s", err.Code, err.Message))
	}

	return ret
}

func auditRun(ctx context.Context, pm *PrjAndMods, ret *Result) error {
	timer := util.TimerStart()
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	generated, err := getGeneratedFiles(tgt, pm.Prj.Ignore, pm.Logger)
	if err != nil {
		return err
	}

	src, err := getModuleFiles(ctx, pm)
	if err != nil {
		return err
	}

	audits := make(diff.Diffs, 0, len(generated))
	for _, g := range generated {
		if !slices.Contains(src, g) {
			if (!strings.HasSuffix(g, "client.js.map")) && (!strings.HasSuffix(g, "file/header.go")) {
				audits = append(audits, &diff.Diff{Path: g, Status: diff.StatusDifferent})
			}
		}
	}

	empty, err := getEmptyFolders(tgt, pm.Prj.Ignore, pm.Logger)
	if err != nil {
		return err
	}
	for _, e := range empty {
		audits = append(audits, &diff.Diff{Path: e, Status: diff.StatusMissing})
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
