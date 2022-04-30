package build

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const gomod = "go.mod"

func OnDepsUpgrade(ctx context.Context, prj *project.Project, up string, o string, n string, pSvc *project.Service, logger *zap.SugaredLogger) error {
	var deps Dependencies
	if up == "all" {
		curr, err := LoadDeps(ctx, prj.Key, prj.Path, true, pSvc.GetFilesystem(prj), false, logger)
		if err != nil {
			return err
		}
		deps = make(Dependencies, 0, len(curr))
		for _, x := range curr {
			if x.Version != x.Available {
				deps = append(deps, x)
			}
		}
	} else {
		if o == "" {
			return errors.New("must provide [old] argument when upgrading")
		}
		if n == "" {
			return errors.New("must provide [new] argument when upgrading")
		}
		deps = Dependencies{{Key: up, Version: o, Available: n}}
	}
	err := upgradeDeps(ctx, prj, deps, pSvc, logger)
	if err != nil {
		return err
	}
	return nil
}

func upgradeDeps(ctx context.Context, prj *project.Project, deps Dependencies, pSvc *project.Service, logger *zap.SugaredLogger) error {
	logger.Infof("upgrading [%d] dependencies for [%s]", len(deps), prj.Key)
	fs := pSvc.GetFilesystem(prj)
	err := bumpGoMod(ctx, prj, fs, deps, logger)
	if err != nil {
		return err
	}
	return nil
}

func bumpGoMod(ctx context.Context, prj *project.Project, fs filesystem.FileLoader, deps Dependencies, logger *zap.SugaredLogger) error {
	bts, err := fs.ReadFile(gomod)
	if err != nil {
		return errors.Wrap(err, "unable to read [go.mod]")
	}
	lines := strings.Split(string(bts), "\n")

	for _, dep := range deps {
		logger.Infof("upgrading [%s] dependency [%s] from [%s] to [%s]", prj.Key, dep.Key, dep.Version, dep.Available)
		hit := false
		for idx, l := range lines {
			if !strings.HasPrefix(l, "\t"+dep.Key) {
				continue
			}
			hit = true
			vIdx := strings.Index(l, dep.Version)
			if vIdx == -1 {
				return errors.Errorf("dependency [%s] does not match [%s], the original version", dep.Key, dep.Version)
			}
			newLine := l[:vIdx] + dep.Available + l[vIdx+len(dep.Version):]
			lines[idx] = newLine
		}
		if !hit {
			return errors.Errorf("can't find [%s] in [go.mod]", dep.Key)
		}
	}

	final := strings.Join(lines, "\n")
	err = fs.WriteFile(gomod, []byte(final), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrap(err, "unable to write [go.mod]")
	}

	_, _, err = telemetry.RunProcessSimple(ctx, "go mod tidy", prj.Path, logger)
	if err != nil {
		return errors.Wrap(err, "unable to run [go mod tidy]")
	}
	return nil
}

func SetDepsMap(ctx context.Context, projects project.Projects, dep *Dependency, pSvc *project.Service, logger *zap.SugaredLogger) (string, error) {
	logger.Infof("upgrading dependency [%s] to [%s]", dep.Key, dep.Version)
	var affected int

	_, errs := util.AsyncCollect(projects, func(item *project.Project) (any, error) {
		t := util.TimerStart()
		fs := pSvc.GetFilesystem(item)
		var matched bool
		bytes, err := fs.ReadFile(gomod)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read [go.mod] for project [%s]", item.Key)
		}
		lines := strings.Split(string(bytes), "\n")
		ret := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.Contains(line, dep.Key+" ") {
				start := strings.Index(line, " v")
				if start == -1 {
					return nil, errors.Errorf("project [%s] does not contain a version in [%s]", item.Key, line)
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
			content := strings.Join(ret, "\n")
			err = fs.WriteFile(gomod, []byte(content), filesystem.DefaultMode, true)
			if err != nil {
				return nil, errors.Wrap(err, "unable to write [go.mod]")
			}
			_, _, err = telemetry.RunProcessSimple(ctx, "go mod tidy", item.Path, logger)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to run [go mod tidy] in path [%s]", item.Path)
			}
			logger.Infof("completed upgrade of [%s] in [%s]", item.Key, util.MicrosToMillis(t.End()))
		}
		return nil, nil
	})
	if len(errs) > 0 {
		return "", errs[0]
	}
	return fmt.Sprintf("upgraded [%s] to [%s] in [%d] projects", dep.Key, dep.Version, affected), nil
}
