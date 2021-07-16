package action

import (
	"fmt"

	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/svg"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func onSVG(prj *project.Project, cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)

	src := fmt.Sprintf("%s/client/src/svg", prj.Path)
	tgt := fmt.Sprintf("%s/app/util/svg.go", prj.Path)

	count, err := svg.Run(src, tgt)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	ret.AddLog("creating [%d] SVGs for project [%s]", count, prj.Key)

	return ret
}
