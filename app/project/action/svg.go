package action

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/project/svg"
)

const refreshMode = "refresh"

func onSVG(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeSVG, pm.Prj, pm.Cfg, pm.Logger)

	tgt := fmt.Sprintf("%s/app/util/svg.go", pm.Prj.Path)

	fs, err := pm.PSvc.GetFilesystem(pm.Prj)
	if err != nil {
		return ret.WithError(err)
	}

	if pm.Cfg.GetStringOpt("mode") == refreshMode {
		ret.AddLog("refreshing app SVG for project [%s]", pm.Prj.Key)
		err = svg.RefreshAppIcon(ctx, pm.Prj, fs, pm.Logger)
		if err != nil {
			return errorResult(err, TypeSVG, pm.Cfg, pm.Logger)
		}
		return ret
	}

	count, err := svg.Run(fs, tgt, pm.Logger)
	if err != nil {
		return errorResult(err, TypeSVG, pm.Cfg, pm.Logger)
	}

	ret.AddLog("creating [%d] SVGs for project [%s]", count, pm.Prj.Key)

	return ret
}
