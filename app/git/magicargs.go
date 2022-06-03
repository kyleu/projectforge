package git

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) magicArgsFor(ctx context.Context, prj *project.Project, message string, dryRun bool, logger util.Logger) (*magicArgs, error) {
	statRet, err := s.Status(ctx, prj, logger)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get status for project [%s]", prj.Key)
	}
	branch := statRet.DataString("branch")

	data := util.ValueMap{"branch": branch, "magic": true, "commitMessage": message}

	dirty := statRet.DataStringArray("dirty")
	dirtyCount := len(dirty)
	if len(dirty) > 0 {
		data["dirtyCount"] = dirtyCount
	}

	ahead := statRet.DataInt("commitsAhead")
	behind := statRet.DataInt("commitsBehind")

	ret := NewResult(prj, "no changes needed", data)

	return &magicArgs{
		Ctx: ctx, Prj: prj, DryRun: dryRun, Dirty: dirtyCount, Ahead: ahead, Behind: behind,
		Branch: branch, Message: message, Result: ret, Logger: logger,
	}, nil
}
