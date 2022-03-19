package action

import (
	"projectforge.dev/projectforge/app/build"
)

func onDeps(pm *PrjAndMods, ret *Result) *Result {
	if up, _ := ret.Args.GetString("upgrade", true); up != "" {
		o, _ := ret.Args.GetString("old", true)
		n, _ := ret.Args.GetString("new", true)
		err := build.OnDepsUpgrade(pm.Prj, up, o, n, pm.PSvc, pm.Logger)
		if err != nil {
			return ret.WithError(err)
		}
	}
	deps, err := build.LoadDeps(pm.Prj.Path)
	ret.Data = deps
	if err != nil {
		return ret.WithError(err)
	}
	return ret
}
