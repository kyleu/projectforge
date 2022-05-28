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

	ahead := statRet.DataInt("commitsAhead")
	behind := statRet.DataInt("commitsBehind")

	ret := NewResult(prj, "no changes needed", data)

	if ahead > 0 {
		if behind > 0 {
			ret.Status = "conflicting commits"
		} else {
			add("pushing [%d] commits", ahead)
			x, err := s.Push(ctx, prj, logger)
			if err != nil {
				return nil, errors.Wrap(err, "unable to push")
			}
			ret.Data = x.Data.Merge(ret.Data)
		}
	} else if behind > 0 {
		add("pulling [%d] commits", behind)
		x, err := s.Pull(ctx, prj, logger)
		if err != nil {
			return nil, errors.Wrap(err, "unable to pull")
		}
		ret.Data = x.Data.Merge(ret.Data)
	}

	return ret, nil
}
