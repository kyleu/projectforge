package codegen

import (
	"github.com/kyleu/projectforge/app/project"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Test(prj *project.Project) (interface{}, error) {
	return "Codegen!", nil
}
