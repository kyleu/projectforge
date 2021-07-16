package project

import (
	"io/ioutil"
	"path/filepath"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const ConfigFilename = ".projectforge.json"

type Service struct {
	filesystems map[string]filesystem.FileLoader
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{filesystems: map[string]filesystem.FileLoader{}, logger: logger}
}

func (s *Service) GetFilesystem(prj *Project) filesystem.FileLoader {
	curr, ok := s.filesystems[prj.Key]
	if ok {
		return curr
	}

	fs := filesystem.NewFileSystem(prj.Path, s.logger)

	s.filesystems[prj.Key] = fs
	return fs
}

func (s *Service) Root() (*Project, error) {
	return s.Load(".")
}

func (s *Service) Load(path string) (*Project, error) {
	b, err := ioutil.ReadFile(filepath.Join(path, ConfigFilename))
	if err != nil {
		return nil, err
	}

	ret := &Project{}
	err = util.FromJSON(b, &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load project from [%s]", ConfigFilename)
	}
	ret.Path = path
	return ret, nil
}
