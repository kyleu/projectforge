package action

import (
	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func onMerge(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)

	to, _ := pm.Cfg.GetString("to", true)
	if to == "" {
		to = "project"
	}
	file, _ := pm.Cfg.GetString("file", true)

	start := util.TimerStart()
	_, dfs, err := diffs(pm, true)
	if err != nil {
		return ret.WithError(err)
	}

	if file != "" {
		var matched []*diff.Diff
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
		case "project":
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

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}

func mergeLeft(pm *PrjAndMods, f *diff.Diff, ret *Result) *Result {
	if err := pm.MSvc.UpdateFile(pm.Mods, f); err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("merged [%s] from project to module", f.Path)
	return ret
}

func mergeRight(pm *PrjAndMods, f *diff.Diff, ret *Result) *Result {
	ret.AddLog("merged [%s] from module to project", f.Path)
	return ret
}
