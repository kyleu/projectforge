package action

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
)

func onPreview(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	start := util.TimerStart()
	_, dfs, err := diffs(pm)
	if err != nil {
		return ret.WithError(err)
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: dfs, Duration: util.TimerEnd(start)}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
