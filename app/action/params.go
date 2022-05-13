package action

import (
	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Params struct {
	ProjectKey string
	T          Type
	Cfg        util.ValueMap
	RootFiles  filesystem.FileLoader
	MSvc       *module.Service
	PSvc       *project.Service
	ESvc       *export.Service
	CLI        bool
	Logger     util.Logger
}
