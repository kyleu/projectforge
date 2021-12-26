package export

import (
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/file"
)

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "export")
	return &Service{logger: logger}
}

func (s *Service) Export(args *Args) (file.Files, error) {
	s.logger.Debugf("starting file export")
	var ret file.Files

	for _, m := range args.Models {
		modelFile := exportModelFile(m, args)
		svcFile := exportServiceFile(m, args)
		ret = append(ret, modelFile, svcFile)
	}

	return ret, nil
}
