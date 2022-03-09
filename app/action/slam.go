package action

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func onSlam(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()
	srcFiles, dfs, err := diffs(pm)
	if err != nil {
		return ret.WithError(err)
	}

	tgtFS := pm.PSvc.GetFilesystem(pm.Prj)
	for _, f := range dfs {
		switch f.Status {
		case diff.StatusIdentical, diff.StatusMissing, diff.StatusSkipped:
			// noop
		case diff.StatusDifferent, diff.StatusNew:
			src := srcFiles.Get(f.Path)
			err := tgtFS.WriteFile(f.Path, []byte(src.Content), src.Mode, true)
			if err != nil {
				return ret.WithError(errors.Wrapf(err, "unable to write updated content to [%s]", f.Path))
			}
		default:
			return ret.WithError(errors.Errorf("unhandled diff status [%s]", f.Status))
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: dfs, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
