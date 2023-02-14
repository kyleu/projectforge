package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

const projectKey = "project"

func onGenerate(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeGenerate, pm.Prj, pm.Cfg, pm.Logger)
	timer := util.TimerStart()

	to := pm.Cfg.GetStringOpt("to")
	if to == "" {
		to = projectKey
	}
	fl := pm.Cfg.GetStringOpt("file")

	srcFiles, dfs, err := diffs(ctx, pm)
	if err != nil {
		return ret.WithError(err)
	}
	dfs, err = limitDiffs(dfs, fl)
	if err != nil {
		return ret.WithError(err)
	}

	tgtFS := pm.PSvc.GetFilesystem(pm.Prj)
	for _, f := range dfs {
		switch f.Status {
		case diff.StatusIdentical, diff.StatusMissing, diff.StatusSkipped:
			// noop
		case diff.StatusDifferent, diff.StatusNew:
			switch to {
			case "module":
				ret = mergeToModule(pm, f, ret)
			case projectKey:
				ret = gen(pm, srcFiles, f, ret, tgtFS)
			}
		default:
			return ret.WithError(errors.Errorf("unhandled diff status [%s]", f.Status))
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: dfs, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func limitDiffs(dfs diff.Diffs, fl string) (diff.Diffs, error) {
	if fl == "" {
		return dfs, nil
	}
	var matched diff.Diffs
	for _, d := range dfs {
		if d.Path == fl {
			matched = append(matched, d)
		}
	}
	if len(matched) == 0 {
		return nil, errors.Errorf("no file [%s] with a difference to merge", fl)
	}
	return matched, nil
}

func gen(pm *PrjAndMods, srcFiles file.Files, f *diff.Diff, ret *Result, tgtFS filesystem.FileLoader) *Result {
	src := srcFiles.Get(f.Path)
	if src == nil {
		return ret.WithError(errors.Errorf("unable to read file from [%s]", f.Path))
	}
	if idx := strings.Index(src.Content, file.GenerateOncePattern); idx > -1 {
		src.Content = src.Content[strings.Index(src.Content[idx:], "\n")+idx+1:]
	}
	err := tgtFS.WriteFile(f.Path, []byte(src.Content), src.Mode, true)
	if err != nil {
		return ret.WithError(errors.Wrapf(err, "unable to write updated content to [%s]", f.Path))
	}
	return ret
}

func mergeToModule(pm *PrjAndMods, d *diff.Diff, ret *Result) *Result {
	logs, err := pm.MSvc.UpdateFile(pm.Mods, d, pm.Logger)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("merged [%s] from project to module: %s", d.Path, strings.Join(logs, ", "))
	return ret
}
