package git

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Magic(ctx context.Context, prj *project.Project, message string, dryRun bool, logger util.Logger) (*Result, error) {
	args, err := s.magicArgsFor(ctx, prj, message, dryRun, logger)
	if err != nil {
		return nil, err
	}

	var logs []string
	add := func(msg string, args ...any) {
		logs = append(logs, fmt.Sprintf(msg, args...))
	}

	switch {
	case args.Ahead == 0 && args.Behind == 0:
		if args.Dirty > 0 {
			args.Result.Status = "commit"
			if err = s.magicCommit(args, add); err != nil {
				return nil, err
			}
			if err = s.magicPush(args, 1, add); err != nil {
				return nil, err
			}
		}
	case args.Ahead == 0 && args.Behind > 0:
		if args.Dirty > 0 {
			if err = s.magicStash(args, add); err != nil {
				return nil, err
			}
		}
		args.Result.Status = "pull"
		if err = s.magicPull(args, add); err != nil {
			return nil, err
		}
		if args.Dirty > 0 {
			args.Result.Status += ", commit"
			if err = s.magicStashPop(args, add); err != nil {
				return nil, err
			}
			if err = s.magicCommit(args, add); err != nil {
				return nil, err
			}
			if err = s.magicPush(args, 1, add); err != nil {
				return nil, err
			}
		}
	case args.Ahead > 0 && args.Behind == 0:
		if args.Dirty == 0 {
			args.Result.Status = "push"
		} else {
			args.Result.Status = "commit, push"
			if err = s.magicCommit(args, add); err != nil {
				return nil, err
			}
		}
		if err = s.magicPush(args, args.Ahead+1, add); err != nil {
			return nil, err
		}
	case args.Ahead > 0 && args.Behind > 0:
		args.Result.Status = "conflicting commits"
		return nil, errors.New("TODO: handle conflicting commits")
	default:
		return nil, errors.New("invalid git state")
	}

	args.Result.Data["logs"] = logs
	return args.Result, nil
}
