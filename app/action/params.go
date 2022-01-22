package action

import (
	"github.com/kyleu/projectforge/app/export"
	"github.com/kyleu/projectforge/app/lib/filesystem"
	"github.com/kyleu/projectforge/app/lib/telemetry"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

type Params struct {
	Span       *telemetry.Span
	ProjectKey string
	T          Type
	Cfg        util.ValueMap
	RootFiles  filesystem.FileLoader
	MSvc       *module.Service
	PSvc       *project.Service
	ESvc       *export.Service
	CLI        bool
	Logger     *zap.SugaredLogger
}
