package action

import (
	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func onSlam(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()
	srcFiles, diffs, err := diffs(pm, true)
	if err != nil {
		return ret.WithError(err)
	}

	for _, f := range diffs {
		switch f.Status {
		case diff.StatusIdentical, diff.StatusMissing, diff.StatusSkipped:
			// noop
		case diff.StatusDifferent, diff.StatusNew:
			src := srcFiles.Get(f.Path)
			tgtFS := pm.PSvc.GetFilesystem(pm.Prj)
			err := tgtFS.WriteFile(f.Path, []byte(src.Content), src.Mode, true)
			if err != nil {
				return ret.WithError(errors.Wrapf(err, "unable to write updated content to [%s]", f.Path))
			}
		default:
			return ret.WithError(errors.Errorf("unhandled diff status [%s]", f.Status))
		}
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: diffs, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
