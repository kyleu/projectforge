package action

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

const projectKey = "project"

func onMerge(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeMerge, pm.Prj, pm.Cfg, pm.Logger)

	to := pm.Cfg.GetStringOpt("to")
	if to == "" {
		to = projectKey
	}
	file := pm.Cfg.GetStringOpt("file")

	timer := util.TimerStart()
	_, dfs, err := diffs(ctx, pm)
	if err != nil {
		return ret.WithError(err)
	}

	if file != "" {
		var matched diff.Diffs
		for _, d := range dfs {
			if d.Path == file {
				matched = append(matched, d)
			}
		}
		if len(matched) == 0 {
			return ret.WithError(errors.Errorf("no file [%s] with a difference to merge", file))
		}
		dfs = matched
	}

	for _, f := range dfs {
		switch to {
		case "module":
			switch f.Status {
			case diff.StatusIdentical, diff.StatusSkipped, diff.StatusNew:
				// noop
			case diff.StatusMissing, diff.StatusDifferent:
				ret = mergeLeft(pm, f, ret)
			default:
				return ret.WithError(errors.Errorf("unhandled diff status [%s]", f.Status))
			}
		case projectKey:
			switch f.Status {
			case diff.StatusIdentical, diff.StatusSkipped, diff.StatusMissing:
				// noop
			case diff.StatusNew, diff.StatusDifferent:
				ret = mergeRight(pm, f, ret)
			default:
				return ret.WithError(errors.Errorf("unhandled diff status [%s]", f.Status))
			}
		default:
			return ret.WithError(errors.Errorf("invalid \"to\" destination [%s]", to))
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func mergeLeft(pm *PrjAndMods, d *diff.Diff, ret *Result) *Result {
	logs, err := pm.MSvc.UpdateFile(pm.Mods, d, pm.Logger)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("merged [%s] from project to module: %s", d.Path, strings.Join(logs, ", "))
	return ret
}

func mergeRight(pm *PrjAndMods, d *diff.Diff, ret *Result) *Result {
	loader := pm.PSvc.GetFilesystem(pm.Prj)
	if !loader.Exists(d.Path) {
		return ret.WithError(errors.Errorf("no file found at path [%s] for project [%s]", d.Path, pm.Prj.Key))
	}
	b, err := loader.ReadFile(d.Path)
	if err != nil {
		return ret.WithError(err)
	}
	stat, err := loader.Stat(d.Path)
	if err != nil {
		return ret.WithError(err)
	}

	newContent, err := diff.Apply(b, d)
	if err != nil {
		return ret.WithError(err)
	}
	msg := ""
	if len(newContent) > 0 {
		if bytes.Equal(b, newContent) {
			msg = fmt.Sprintf("no changes required to [%s] for project [%s]", d.Path, pm.Prj.Key)
		} else {
			msg = fmt.Sprintf("wrote [%d] bytes to [%s] for project [%s]", len(newContent), d.Path, pm.Prj.Key)
			err = loader.WriteFile(d.Path, newContent, stat.Mode(), true)
			if err != nil {
				return ret.WithError(err)
			}
		}
	}
	ret.AddLog("merged [%s] from module to project: %s", d.Path, msg)
	return ret
}
