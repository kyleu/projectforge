package action

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Params struct {
	Span       trace.Span
	ProjectKey string
	T          Type
	Cfg        util.ValueMap
	MSvc       *module.Service
	PSvc       *project.Service
	Logger     *zap.SugaredLogger
}
