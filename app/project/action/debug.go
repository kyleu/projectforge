package action

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func onDebug(pm *PrjAndMods) *Result {
	ret := newResult(TypeDebug, pm.Prj, pm.Cfg, pm.Logger)

	ret.AddWarn("Project:")
	ret.Log(util.ToJSON(pm.Prj))

	ret.AddWarn("Modules:")
	lo.ForEach(pm.Mods.Keys(), func(m string, _ int) {
		ret.Log(" - " + m)
	})

	return ret
}
