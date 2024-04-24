package export

import (
	"fmt"
	"projectforge.dev/projectforge/app/project/export/csfiles"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Files(p *project.Project, addHeader bool, linebreak string) (_ file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	return files.All(p, addHeader, linebreak)
}

func (s *Service) FilesCSharp(p *project.Project, addHeader bool, linebreak string) (_ file.Files, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				e = err
			} else {
				e = errors.New(fmt.Sprint(rec))
			}
		}
	}()
	return csfiles.CSAll(p, addHeader)
}
