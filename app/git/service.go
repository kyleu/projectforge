package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const ok = "OK"

type Service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "build")
	return &Service{logger: logger}
}

func (s Service) Status(ctx context.Context, prj *project.Project, logger *zap.SugaredLogger) (*Result, error) {
	ctx, span, logger := telemetry.StartSpan(ctx, "git.status:"+prj.Key, logger)
	defer span.Complete()

	dirty, err := gitStatus(prj.Path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git status")
	}
	branch := gitBranch(prj.Path)
	data := util.ValueMap{"branch": branch}
	if len(dirty) > 0 {
		data["dirty"] = dirty
	}
	status := ok
	if len(dirty) > 0 {
		status = fmt.Sprintf("[%d] changes", len(dirty))
	}
	return NewResult(prj, status, data), nil
}

func (s Service) CreateRepo(prj *project.Project) (*Result, error) {
	return NewResult(prj, "TODO", util.ValueMap{"TODO": "Create Repo"}), nil
}

func (s Service) Magic(prj *project.Project) (*Result, error) {
	return NewResult(prj, "TODO", util.ValueMap{"TODO": "Magic!"}), nil
}

func (s Service) Fetch(prj *project.Project) (*Result, error) {
	x, err := gitFetch(prj.Path, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch")
	}
	count := 0
	lines := strings.Split(x, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "   ") {
			count++
		}
	}
	status := ok
	fetched := "no updates"
	if count > 0 {
		status = fmt.Sprintf("[%d] %s fetched", count, util.StringPluralMaybe("update", count))
		fetched = status
	}

	return NewResult(prj, status, util.ValueMap{"updates": fetched}), nil
}

func (s Service) Commit(prj *project.Project, msg string) (*Result, error) {
	result, err := gitCommit(prj.Path, msg)
	if err != nil {
		return nil, err
	}

	return NewResult(prj, ok, util.ValueMap{"commit": result}), nil
}
