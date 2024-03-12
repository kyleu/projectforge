package notebook

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/exec"
	"{{{ .Package }}}/app/util"
)

var baseURL = fmt.Sprintf("http://localhost:%d/", util.AppPort+10)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Status() string {
	rsp, err := http.DefaultClient.Get(baseURL)
	if err != nil {
		return "not-started"
	}
	if rsp.StatusCode != 200 {
		return "invalid"
	}
	return "running"
}

func (s *Service) Start(execSvc *exec.Service) error {
	if s.Status() == "running" {
		return errors.Errorf("can't start notebook, something is already listening on port [%d]", util.AppPort+10)
	}
	_, err := util.StartProcess("bin/dev.sh", "./notebook", nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
