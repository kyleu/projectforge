package checks

import (
	"strings"

	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func simpleOut(path string, cmd string, args []string, outCheck func(r *doctor.Result, out string) *doctor.Result) doctor.CheckFn {
	return func(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
		fullCmd := strings.Join(append([]string{cmd}, args...), " ")
		exitCode, out, err := util.RunProcessSimple(fullCmd, path)
		if err != nil {
			msg := "[%s] is not present on your computer"
			return r.WithError(doctor.NewError("missing", msg, cmd))
		}
		if exitCode != 0 {
			return r.WithError(doctor.NewError("exitcode", "[%s] returned [%d] as an exit code", fullCmd, exitCode))
		}
		return outCheck(r, out)
	}
}

func loadProject(r *doctor.Result, logger *zap.SugaredLogger) (*project.Project, *doctor.Result) {
	dir := "."
	fs := filesystem.NewFileSystem(dir, logger)
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

func checkMods(p *project.Project, r *doctor.Result) *doctor.Result {
	hasMod := func(key string) bool {
		for _, m := range p.Modules {
			if m == key {
				return true
			}
		}
		return false
	}
	if hasMod("desktop") && (!p.Build.Desktop) {
		r = r.WithError(doctor.NewError("config", "desktop module is enabled, but desktop build isn't set"))
	}
	if hasMod("ios") && (!p.Build.IOS) {
		r = r.WithError(doctor.NewError("config", "iOS module is enabled, but iOS build isn't set"))
	}
	if hasMod("android") && (!p.Build.Android) {
		r = r.WithError(doctor.NewError("config", "Android module is enabled, but Android build isn't set"))
	}
	return r
}
