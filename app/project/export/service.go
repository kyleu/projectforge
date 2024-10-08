package export

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Files(p *project.Project, linebreak string) (_ file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	return files.All(p, linebreak)
}
