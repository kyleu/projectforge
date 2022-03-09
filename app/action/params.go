package action

import (
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
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
