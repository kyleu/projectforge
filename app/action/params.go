package action

import (
	"projectforge.dev/app/lib/filesystem"
	"projectforge.dev/app/lib/telemetry"
	"projectforge.dev/app/module"
	"projectforge.dev/app/project"
	"projectforge.dev/app/util"
	"go.uber.org/zap"
	"projectforge.dev/app/export"
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
