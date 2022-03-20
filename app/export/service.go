package export

import (
	"context"

	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/export/files"
	"projectforge.dev/projectforge/app/export/inject"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "export")
	return &Service{logger: logger}
}

func (s *Service) Files(ctx context.Context, args *model.Args, addHeader bool, logger *zap.SugaredLogger) (file.Files, error) {
	return files.All(ctx, args, addHeader, logger)
}

func (s *Service) Inject(ctx context.Context, args *model.Args, fs file.Files, logger *zap.SugaredLogger) error {
	return inject.All(ctx, args, fs, logger)
}
