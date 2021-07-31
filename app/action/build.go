package action

import (
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func onBuild(prj *project.Project, cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	ret.AddLog("building project [%s] in [%s]", prj.Key, prj.Path)

	exitCode, out, err := util.RunProcessSimple("make build", prj.Path)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		ret.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}

	return ret
}
