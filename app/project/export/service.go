package export

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files"
	"projectforge.dev/projectforge/app/project/export/model"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Files(p *project.Project, args *model.Args, addHeader bool, linebreak string) (_ file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	return files.All(p, args, addHeader, linebreak)
}
