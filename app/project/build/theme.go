package build

import (
	"context"

	"github.com/muesli/gamut"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ThemeRebuild(ctx context.Context, prj *project.Project, pSvc *project.Service, logger util.Logger) (*theme.Theme, error) {
	color := prj.Theme.Base
	if color == "" {
		color = prj.Theme.Light.NavBackground
	}
	thm := theme.ColorTheme(color, gamut.Hex(color))
	prj.Theme = thm
	err := pSvc.Save(prj, logger)
	if err != nil {
		return nil, err
	}
	return thm, nil
}
