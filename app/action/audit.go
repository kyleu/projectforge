package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func onAudit(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeAudit, pm.Prj, pm.Cfg, pm.Logger)

	var err error
	switch f := pm.Cfg.GetStringOpt("fix"); f {
	case "remove":
		err = auditRemove(ctx, pm, ret)
	case "header":
		err = auditHeader(ctx, pm, ret)
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
	return ret
}

func auditRun(ctx context.Context, pm *PrjAndMods, ret *Result) error {
	timer := util.TimerStart()
	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	filenames, err := tgt.ListFilesRecursive("", pm.Prj.Ignore)
	if err != nil {
		return err
	}

	var generated []string
	for _, fn := range filenames {
		b, e := tgt.PeekFile(fn, 1024)
		if e != nil {
			return e
		}
		if file.ContainsHeader(string(b)) {
			generated = append(generated, fn)
		}
	}

	src, err := pm.MSvc.GetFilenames(pm.Mods)
	if err != nil {
		return err
	}

	if pm.Mods.Get("export") != nil {
		args := &model.Args{}
		if argsX := pm.Prj.Info.ModuleArg("export"); argsX != nil {
			e := util.CycleJSON(argsX, &args)
			if e != nil {
				return err
			}
		}
		args.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(ctx, args, true, pm.Logger)
		if e != nil {
			return err
		}
		for _, f := range files {
			src = append(src, f.FullPath())
		}
	}

	var audits []*diff.Diff
	for _, g := range generated {
		if !slices.Contains(src, g) {
			if (!strings.HasSuffix(g, "client.js.map")) && (!strings.HasSuffix(g, "file/header.go")) {
				audits = append(audits, &diff.Diff{Path: g, Status: diff.StatusMissing})
			}
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: audits, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return nil
}

func auditRemove(ctx context.Context, pm *PrjAndMods, ret *Result) error {
	path := pm.Cfg.GetStringOpt("file")
	if path == "" {
		return errors.New("must specify [file] as argument")
	}
	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	ret.AddLog("removed file [%s]", path)
	return tgt.Remove(path)
}

func auditHeader(ctx context.Context, pm *PrjAndMods, ret *Result) error {
	path := pm.Cfg.GetStringOpt("file")
	if path == "" {
		return errors.New("must specify [file] as argument")
	}
	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	content, err := tgt.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	var hit bool
	fixed := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, file.HeaderContent) {
			hit = true
		} else {
			fixed = append(fixed, line)
		}
	}
	if hit {
		final := strings.Join(fixed, "\n")
		err = tgt.WriteFile(path, []byte(final), filesystem.DefaultMode, true)
		if err != nil {
			return errors.Wrapf(err, "can't overwrite file at path [%s]", path)
		}
		ret.AddLog("removed header from file [%s]", path)
	} else {
		ret.AddLog("no header to remove from file [%s]", path)
	}
	return nil
}
