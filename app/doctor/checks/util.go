package checks

import (
	"context"
	"strings"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func simpleOut(path string, cmd string, args []string, outCheck func(ctx context.Context, r *doctor.Result, out string) *doctor.Result) doctor.CheckFn {
	return func(ctx context.Context, r *doctor.Result, logger util.Logger) *doctor.Result {
		fullCmd := strings.Join(append([]string{cmd}, args...), " ")
		exitCode, out, err := telemetry.RunProcessSimple(ctx, fullCmd, path, logger)
		if err != nil {
			return r.WithError(doctor.NewError("missing", "[%s] is not present on your computer", cmd))
		}
		if exitCode != 0 {
			return r.WithError(doctor.NewError("exitcode", "[%s] returned [%d] as an exit code", fullCmd, exitCode))
		}
		return outCheck(ctx, r, out)
	}
}

func loadRootProject(r *doctor.Result) (*project.Project, *doctor.Result) {
	const dir = "."
	fs := filesystem.NewFileSystem(dir)
	if !fs.Exists(project.ConfigFilename) {
		return nil, r.WithError(doctor.NewError("missing", "no project found in [%s]", dir))
	}
	b, err := fs.ReadFile(project.ConfigFilename)
	if err != nil {
		return nil, r.WithError(doctor.NewError("missing", "unable to read project from [%s]", dir))
	}
	p := &project.Project{}
	err = util.FromJSON(b, p)
	if err != nil {
		return nil, r.WithError(doctor.NewError("invalid", "unable to parse project JSON from [%s]", dir))
	}
	return p, r
}
