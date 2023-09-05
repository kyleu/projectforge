package git

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type magicArgs struct {
	Ctx     context.Context //nolint:containedctx
	Prj     *project.Project
	DryRun  bool
	Dirty   int
	Ahead   int
	Behind  int
	Branch  string
	Message string
	Result  *Result
	Logger  util.Logger
}

func (s *Service) magicCommit(a *magicArgs, add func(string, ...any)) error {
	add("committing [%s] with message [%s]", util.StringPlural(a.Dirty, "file"), a.Message)
	if !a.DryRun {
		x, err := s.Commit(a.Ctx, a.Prj, a.Message, a.Logger)
		if err != nil {
			return errors.Wrap(err, "unable to commit")
		}
		a.Result.Data = x.Data.Merge(a.Result.Data)
	}
	return nil
}

func (s *Service) magicPull(a *magicArgs, add func(string, ...any)) error {
	add("pulling [%s] from [%s]", util.StringPlural(a.Behind, "commit"), a.Branch)
	if !a.DryRun {
		x, err := s.Pull(a.Ctx, a.Prj, a.Logger)
		if err != nil {
			return errors.Wrap(err, "unable to pull")
		}
		a.Result.Data = x.Data.Merge(a.Result.Data)
	}
	return nil
}

func (s *Service) magicPush(a *magicArgs, count int, add func(string, ...any)) error {
	add("pushing [%s] to [%s]", util.StringPlural(count, "commit"), a.Branch)
	if !a.DryRun {
		x, err := s.Push(a.Ctx, a.Prj, a.Logger)
		if err != nil {
			return errors.Wrap(err, "unable to push")
		}
		a.Result.Data = x.Data.Merge(a.Result.Data)
	}
	return nil
}

func (s *Service) magicStash(a *magicArgs, add func(string, ...any)) error {
	add("stashing [%s]", util.StringPlural(a.Dirty, "file"))
	if !a.DryRun {
		_, err := s.gitStash(a.Ctx, a.Prj, a.Logger)
		if err != nil {
			return errors.Wrap(err, "unable to apply stash")
		}
		a.Result.Data["stash"] = "applied"
	}
	return nil
}

func (s *Service) magicStashPop(a *magicArgs, add func(string, ...any)) error {
	add("restoring [%s] from stash", util.StringPlural(a.Dirty, "file"))
	if !a.DryRun {
		_, err := s.gitStashPop(a.Ctx, a.Prj, a.Logger)
		if err != nil {
			return errors.Wrap(err, "unable to pop stash")
		}
		a.Result.Data["stash"] = "popped"
	}
	return nil
}
