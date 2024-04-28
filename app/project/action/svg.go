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
	if pm.Prj.IsCSharp() {
		tgt = fmt.Sprintf("%s/Util/Icons.cs", pm.Prj.Path)
	}

	if pm.Cfg.GetStringOpt("mode") == refreshMode {
		ret.AddLog("refreshing app SVG for project [%s]", pm.Prj.Key)
		err := svg.RefreshAppIcon(ctx, pm.Prj, pm.FS, pm.Logger)
		if err != nil {
			return errorResult(err, TypeSVG, pm.Cfg, pm.Logger)
		}
		return ret
	}

	count, err := svg.Run(pm.FS, tgt, pm.Logger, pm.Prj)
	if err != nil {
		return errorResult(err, TypeSVG, pm.Cfg, pm.Logger)
	}

	ret.AddLog("creating [%d] SVGs for project [%s]", count, pm.Prj.Key)

	return ret
}
