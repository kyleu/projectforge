package git

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Magic(ctx context.Context, prj *project.Project, message string, logger util.Logger) (*Result, error) {
	var logs []string
	add := func(msg string, args ...any) {
		logs = append(logs, fmt.Sprintf(msg, args...))
	}

	statRet, err := s.Status(ctx, prj, logger)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get status for project [%s]", prj.Key)
	}
	add("OK")

	data := util.ValueMap{"branch": statRet.DataString("branch"), "magic": true, "commitMessage": message, "logs": logs}

	if d := statRet.DataStringArray("dirty"); len(d) > 0 {
		data["dirtyCount"] = len(d)
	}

	return NewResult(prj, "OK", data), nil
}
