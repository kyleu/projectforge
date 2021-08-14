package action

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func onBuild(pm *PrjAndMods) *Result {
	ret := newResult(pm.Cfg, pm.Logger)
	ret.AddLog("building project [%s] in [%s]", pm.Prj.Key, pm.Prj.Path)

	exitCode, out, err := util.RunProcessSimple("make build", pm.Prj.Path)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("build output for [" + pm.Prj.Key + "]:\n" + out)
	if exitCode != 0 {
		ret.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}

	return ret
}
