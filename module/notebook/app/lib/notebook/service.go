package notebook

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

var FavoritePages = util.NewOrderedMap[string](false, 10)

type Service struct {
	BaseURL string                `json:"baseURL"`
	FS      filesystem.FileLoader `json:"-"`
}

func NewService() *Service {
	baseURL := util.GetEnv("notebook_base_url", fmt.Sprintf("http://localhost:%d", util.AppPort+10))
	baseURL = strings.TrimSuffix(baseURL, "/")
	fs, _ := filesystem.NewFileSystem("notebook/docs", false, "")
	return &Service{BaseURL: baseURL, FS: fs}
}

func (s *Service) Status() string {
	rsp, err := http.DefaultClient.Get(s.BaseURL)
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
