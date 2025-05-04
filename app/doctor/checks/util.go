package checks

import (
	"context"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func simpleOut(path string, cmd string, args []string, outChecks ...func(ctx context.Context, r *doctor.Result, out string) *doctor.Result) doctor.CheckFn {
	return func(ctx context.Context, r *doctor.Result, logger util.Logger) *doctor.Result {
		fullCmd := util.StringJoin(append([]string{cmd}, args...), " ")
		exitCode, out, err := telemetry.RunProcessSimple(ctx, fullCmd, path, logger)
		if err != nil {
			return r.WithError(doctor.NewError("missing", "[%s] is not present on your computer", cmd))
		}
		if exitCode != 0 {
			return r.WithError(doctor.NewError("exitcode", "[%s] returned [%d] as an exit code", fullCmd, exitCode))
		}
		lo.ForEach(outChecks, func(outCheck func(ctx context.Context, r *doctor.Result, out string) *doctor.Result, _ int) {
			r = outCheck(ctx, r, out)
		})
		return r
	}
}

func loadRootProject(r *doctor.Result) (*project.Project, filesystem.FileLoader, *doctor.Result) {
	const dir = "."
	fs, _ := filesystem.NewFileSystem(dir, false, "")
	if !fs.Exists(project.ConfigDir) {
		return nil, nil, r.WithError(doctor.NewError("missing", "no project found in [%s]", dir))
	}
	b, err := fs.ReadFile(project.ConfigDir + "/project.json")
	if err != nil {
		return nil, nil, r.WithError(doctor.NewError("missing", "unable to read project from [%s]", dir))
	}
	p, err := util.FromJSONObj[*project.Project](b)
	if err != nil {
		return nil, nil, r.WithError(doctor.NewError("invalid", "unable to parse project JSON from [%s]", dir))
	}
	return p, fs, r
}
