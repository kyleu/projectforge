package action

import (
	"projectforge.dev/projectforge/app/util"
)

func onDebug(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)

	ret.AddWarn("Project:")
	ret.AddLog(util.ToJSON(pm.Prj))

	ret.AddWarn("Modules:")
	for _, m := range pm.Mods.Keys() {
		ret.AddLog(" - " + m)
	}

	return ret
}
