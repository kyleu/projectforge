package export

import (
	"github.com/kyleu/projectforge/app/export/files"
	"github.com/kyleu/projectforge/app/export/inject"
	"github.com/kyleu/projectforge/app/export/model"
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

func (s *Service) Files(args *model.Args, addHeader bool) (file.Files, error) {
	return files.All(args, addHeader)
}

func (s *Service) Inject(args *model.Args, fs file.Files) error {
	return inject.All(args, fs)
}
