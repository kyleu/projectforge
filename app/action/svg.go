package action

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/svg"
)

func onSVG(ctx context.Context, pm *PrjAndMods) *Result {
	ret := newResult(TypeSVG, pm.Prj, pm.Cfg, pm.Logger)

	src := fmt.Sprintf("%s/client/src/svg", pm.Prj.Path)
	tgt := fmt.Sprintf("%s/app/util/svg.go", pm.Prj.Path)

	fs := pm.PSvc.GetFilesystem(pm.Prj)
	count, err := svg.Run(fs, src, tgt, pm.Logger)
	if err != nil {
		return errorResult(err, TypeSVG, pm.Cfg, pm.Logger)
	}

	ret.AddLog("creating [%d] SVGs for project [%s]", count, pm.Prj.Key)

	return ret
}
