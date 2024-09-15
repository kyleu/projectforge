package action

import (
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func onPreview(pm *PrjAndMods) *Result {
	ret := newResult(TypePreview, pm.Prj, pm.Cfg, pm.Logger)
	timer := util.TimerStart()
	_, dfs, err := diffs(pm)
	if err != nil {
		return ret.WithError(err)
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: util.OK, Diffs: dfs, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
