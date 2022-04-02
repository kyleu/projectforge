package action

import (
	"context"

	"projectforge.dev/projectforge/app/util"
)

func onDebug(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeDebug, pm.Cfg, pm.Logger)

	ret.AddWarn("Project:")
	ret.AddLog(util.ToJSON(pm.Prj))

	ret.AddWarn("Modules:")
	for _, m := range pm.Mods.Keys() {
		ret.AddLog(" - " + m)
	}

	return ret
}
