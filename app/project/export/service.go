package export

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files"
	"projectforge.dev/projectforge/app/project/export/inject"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Files(ctx context.Context, p *project.Project, args *model.Args, addHeader bool, logger util.Logger) (f file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	f, e = files.All(ctx, p, args, addHeader, logger)
	return
}

func (s *Service) Inject(ctx context.Context, args *model.Args, fs file.Files, logger util.Logger) (e error) {
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
