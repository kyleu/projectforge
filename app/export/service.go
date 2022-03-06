package export

import (
	"go.uber.org/zap"
	"projectforge.dev/app/export/files"
	"projectforge.dev/app/export/inject"
	"projectforge.dev/app/export/model"

	"projectforge.dev/app/file"
)

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "export")
	return &Service{logger: logger}
}

func (s *Service) Files(args *model.Args, addHeader bool) (file.Files, error) {
	return files.All(args, addHeader)
}

func (s *Service) Inject(args *model.Args, fs file.Files) error {
	return inject.All(args, fs)
}
