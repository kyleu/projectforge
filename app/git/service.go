package git

import (
	"github.com/kyleu/projectforge/app/project"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "build")
	return &Service{logger: logger}
}

func (s Service) GetStatus(prj *project.Project) (*Status, error) {
	x, err := gitStatus(prj.Path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git status")
	}
	branch := gitBranch(prj.Path)
	return &Status{Key: prj.Key, Project: prj, Branch: branch, Dirty: x}, nil
}

func (s Service) GetStatusAll(prjs project.Projects) (Statuses, error) {
	ret := make(Statuses, 0, len(prjs))
	for _, prj := range prjs {
		s, err := s.GetStatus(prj)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get status for project [%s]", prj.Key)
		}
		ret = append(ret, s)
	}
	return ret, nil
}
