package action

import (
	"context"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func onPreview(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypePreview, pm.Cfg, pm.Logger)
	timer := util.TimerStart()
	_, dfs, err := diffs(ctx, pm)
	if err != nil {
		return ret.WithError(err)
	}

	mr := &module.Result{Keys: pm.Mods.Keys(), Status: "OK", Diffs: dfs, Duration: timer.End()}
	ret.Modules = append(ret.Modules, mr)
	return ret
}
