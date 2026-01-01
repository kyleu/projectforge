package action

import (
	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/util"
)

type Params struct {
	ProjectKey string
	T          Type
	Cfg        util.ValueMap
	MSvc       *module.Service
	PSvc       *project.Service
	XSvc       *exec.Service
	SSvc       *websocket.Service
	ESvc       *export.Service
	CLI        bool
	Logger     util.Logger
}
