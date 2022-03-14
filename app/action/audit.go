package action

import (
	"strings"

	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func onAudit(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	timer := util.TimerStart()

	tgt := pm.PSvc.GetFilesystem(pm.Prj)
	filenames, err := tgt.ListFilesRecursive("", pm.Prj.Ignore)
	if err != nil {
		return errorResult(err, pm.Cfg, pm.Logger)
	}

	var generated []string
	for _, fn := range filenames {
		b, e := tgt.PeekFile(fn, 1024)
		if e != nil {
			return errorResult(e, pm.Cfg, pm.Logger)
		}
		if file.ContainsHeader(string(b)) {
			generated = append(generated, fn)
		}
	}

	src, err := pm.MSvc.GetFilenames(pm.Mods)
	if err != nil {
		return errorResult(err, pm.Cfg, pm.Logger)
	}

	if pm.Mods.Get("export") != nil {
		args := &model.Args{}
		if argsX := pm.Prj.Info.ModuleArg("export"); argsX != nil {
			e := util.CycleJSON(argsX, &args)
			if e != nil {
				return errorResult(err, pm.Cfg, pm.Logger)
			}
		}
		args.Modules = pm.Mods.Keys()
		files, e := pm.ESvc.Files(args, true)
		if e != nil {
			return errorResult(e, pm.Cfg, pm.Logger)
		}
		for _, f := range files {
			src = append(src, f.FullPath())
		}
	}

	var audits []*diff.Diff
	for _, g := range generated {
		if !util.StringArrayContains(src, g) {
			if (!strings.HasSuffix(g, "client.js.map")) && (!strings.HasSuffix(g, "file/header.go")) {
				audits = append(audits, &diff.Diff{Path: g, Status: diff.StatusMissing})
			}
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: audits, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
