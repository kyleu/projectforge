package export

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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

func (s *Service) Files(ctx context.Context, args *model.Args, addHeader bool, logger *zap.SugaredLogger) (f file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	f, e = files.All(ctx, args, addHeader, logger)
	return
}

func (s *Service) Inject(ctx context.Context, args *model.Args, fs file.Files, logger *zap.SugaredLogger) (e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	e = inject.All(ctx, args, fs, logger)
	return
}
