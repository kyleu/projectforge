package notebook

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

var (
	BaseURL       = fmt.Sprintf("http://localhost:%d/", util.AppPort+10)
	FavoritePages = util.NewOrderedMap[string](false, 10)
)

type Service struct {
	FS filesystem.FileLoader
}

func NewService() *Service {
	fs, _ := filesystem.NewFileSystem("notebook/docs", false, "")
	return &Service{FS: fs}
}

func (s *Service) Status() string {
	rsp, err := http.DefaultClient.Get(BaseURL)
	if err != nil {
		return "not-started"
	}
	if rsp.StatusCode != 200 {
		return "invalid"
	}
	return "running"
}

func (s *Service) Start() error {
	if s.Status() == "running" {
		return errors.Errorf("can't start notebook, something is already listening on port [%d]", util.AppPort+10)
	}
	_, err := util.StartProcess("bin/dev.sh", "./notebook", nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
