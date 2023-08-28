package build

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

//nolint:gocognit
func SetDepsMap(ctx context.Context, projects project.Projects, dep *Dependency, pSvc *project.Service, logger util.Logger) (string, error) {
	logger.Infof("upgrading dependency [%s] to [%s]", dep.Key, dep.Version)
	var affected int

	_, errs := util.AsyncCollect(projects, func(p *project.Project) (any, error) {
		t := util.TimerStart()
		pfs, err := pSvc.GetFilesystem(p)
		if err != nil {
			return nil, err
		}
		var matched bool
		bytes, err := pfs.ReadFile(gomod)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read [go.mod] for project [%s]", p.Key)
		}
		str := string(bytes)
		lines := util.StringSplitLines(str)
		ret := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.Contains(line, dep.Key+" ") {
				start := strings.Index(line, " v")
				if start == -1 {
					return nil, errors.Errorf("project [%s] does not contain a version in [%s]", p.Key, line)
				}
				start++
				offset := strings.Index(line[start:], " ")
				if offset == -1 {
					offset = len(line) - start
				}
				curr := line[start : start+offset]
				if curr == dep.Version {
					ret = append(ret, line)
				} else {
					matched = true
					ret = append(ret, strings.Replace(line, curr, dep.Version, 1))
				}
			} else {
				ret = append(ret, line)
			}
		}
		if matched {
			affected++
			content := strings.Join(ret, util.StringDetectLinebreak(str))
			err = pfs.WriteFile(gomod, []byte(content), filesystem.DefaultMode, true)
			if err != nil {
				return nil, errors.Wrap(err, "unable to write [go.mod]")
			}
			_, _, err = telemetry.RunProcessSimple(ctx, "go mod tidy", p.Path, logger)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to run [go mod tidy] in path [%s]", p.Path)
			}
			logger.Infof("completed upgrade of [%s] in [%s]", p.Key, util.MicrosToMillis(t.End()))
		}
		return nil, nil
	})
	if len(errs) > 0 {
		return "", util.ErrorMerge(errs...)
	}
	return fmt.Sprintf("upgraded [%s] to [%s] in [%d] projects", dep.Key, dep.Version, affected), nil
}

func SetDepsProject(ctx context.Context, prjs project.Projects, key string, pSvc *project.Service, logger util.Logger) (string, error) {
	t := util.TimerStart()

	tgt := prjs.Get(key)
	if tgt == nil {
		return "", errors.Errorf("no project found with key [%s]", key)
	}
	logger.Infof("upgrading all dependencies for project [%s]", key)
	var affected int

	curr, err := LoadDepsMap(prjs, 2, pSvc)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	pfs, err := pSvc.GetFilesystem(tgt)
	if err != nil {
		return "", err
	}
	bytes, err := pfs.ReadFile(gomod)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read [go.mod] for project [%s]", key)
	}
	str := string(bytes)
	lines := util.StringSplitLines(str)
	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		hit, pline, errChild := setDepProcessLine(line, curr, key)
		if errChild != nil {
			return "", errChild
		}
		ret = append(ret, pline)
		if hit {
			affected++
		}
	}
	if affected > 0 {
		content := strings.Join(ret, util.StringDetectLinebreak(str))
		err = pfs.WriteFile(gomod, []byte(content), filesystem.DefaultMode, true)
		if err != nil {
			return "", errors.Wrap(err, "unable to write [go.mod]")
		}
		_, _, err = telemetry.RunProcessSimple(ctx, "go mod tidy", tgt.Path, logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to run [go mod tidy] in path [%s]", tgt.Path)
		}
		logger.Infof("completed upgrade of project [%s] in [%s]", key, util.MicrosToMillis(t.End()))
	}
	return fmt.Sprintf("upgraded [%d] dependencies in project [%s]", affected, key), nil
}

func setDepProcessLine(line string, curr map[string]map[string][]string, key string) (bool, string, error) {
	if dep := ParseDependency(line); dep != nil {
		start := strings.Index(line, " v")
		if start == -1 {
			return false, "", errors.Errorf("project [%s] does not contain a version in [%s]", key, line)
		}
		start++
		offset := strings.Index(line[start:], " ")
		if offset == -1 {
			offset = len(line) - start
		}

		if existing, ok := curr[dep.Key]; ok {
			v := line[start : start+offset]
			newVer := ""
			newCount := 0
			for kx, vx := range existing {
				if len(vx) > newCount {
					newCount = len(vx)
					newVer = kx
				}
			}
			if v != dep.Version {
				return true, strings.Replace(line, v, newVer, 1), nil
			}
		}
	}
	return false, line, nil
}
