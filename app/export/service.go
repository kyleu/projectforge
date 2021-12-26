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
	var ret file.Files

	for _, m := range args.Models {
		fs := file.Files{exportModelFile(m, args), exportServiceFile(m, args), exportControllerFile(m, args)}
		ret = append(ret, fs...)
	}

	return ret, nil
}

func (s *Service) Inject(args *Args, files file.Files) error {
	for _, f := range files {
		err := s.InjectFile(f, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) InjectFile(f *file.File, args *Args) error {
	if args == nil {
		return nil
	}
	var err error
	switch f.FullPath() {
	case "app/controller/routes.go":
		err = injectRoutes(f, args)
	case "app/controller/menu.go":
		err = injectMenu(f, args)
	}
	return err
}
