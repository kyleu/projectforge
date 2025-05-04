package build

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const gomod = "go.mod"

func OnDepsUpgrade(ctx context.Context, prj *project.Project, up string, o string, n string, pSvc *project.Service, logger util.Logger) error {
	var deps Dependencies
	if up == "all" {
		pfs, err := pSvc.GetFilesystem(prj)
		if err != nil {
			return err
		}
		curr, err := LoadDeps(ctx, prj.Key, prj.Path, true, pfs, false, logger)
		if err != nil {
			return err
		}
		deps = lo.Reject(curr, func(x *Dependency, _ int) bool {
			return x.Version == x.Available
		})
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

func upgradeDeps(ctx context.Context, prj *project.Project, deps Dependencies, pSvc *project.Service, logger util.Logger) error {
	logger.Infof("upgrading [%d] dependencies for [%s]", len(deps), prj.Key)
	fs, err := pSvc.GetFilesystem(prj)
	if err != nil {
		return err
	}
	err = bumpGoMod(ctx, prj, fs, deps, logger)
	if err != nil {
		return err
	}
	return nil
}

func bumpGoMod(ctx context.Context, prj *project.Project, fs filesystem.FileLoader, deps Dependencies, logger util.Logger) error {
	bts, err := fs.ReadFile(gomod)
	if err != nil {
		return errors.Wrap(err, "unable to read [go.mod]")
	}
	str := string(bts)
	lines := util.StringSplitLines(str)

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

	final := util.StringJoin(lines, util.StringDetectLinebreak(str))
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
