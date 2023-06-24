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

func LoadDepsEasyMode(key string, fs filesystem.FileLoader) (Dependencies, error) {
	bytes, err := fs.ReadFile(gomod)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read [go.mod] for project [%s]", key)
	}
	lines := strings.Split(string(bytes), "\n")
	ret := make(Dependencies, 0, len(lines))
	for _, line := range lines {
		if !(strings.HasPrefix(line, "\t") && strings.Contains(line, " v")) {
			continue
		}
		start := strings.Index(line, " v")
		if start == -1 {
			return nil, errors.Errorf("project [%s] does not contain a version in [%s]", key, line)
		}
		dep := &Dependency{Key: strings.TrimSpace(line[:start])}
		start++
		offset := strings.Index(line[start:], " ")
		if offset == -1 {
			offset = len(line) - start
		}
		dep.Version = line[start : start+offset]
		ret = append(ret, dep)
	}
	return ret, nil
}

func LoadDeps(
	ctx context.Context, key string, path string, includeUpdates bool, fs filesystem.FileLoader, showAll bool, logger util.Logger,
) (Dependencies, error) {
	actual, err := LoadDepsEasyMode(key, fs)
	if err != nil {
		return nil, err
	}

	cmd := "go list -m all"
	if includeUpdates {
		cmd = "go list -m -u all"
	}
	_, out, err := telemetry.RunProcessSimple(ctx, cmd, path, logger)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	ret := make(Dependencies, 0, len(lines))
	for _, l := range lines {
		if strings.HasPrefix(l, "go list -m:") {
			e := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(l, "; to add it:"), "go list -m:"))
			return nil, errors.New("need [go mod tidy]:" + e)
		}
		d := &Dependency{}
		parts := util.StringSplitAndTrim(l, " ")
		if len(parts) == 0 {
			continue
		}
		if len(parts) > 0 {
			d.Key = parts[0]
		}
		if len(parts) > 1 {
			d.Version = parts[1]
		}
		if len(parts) > 2 {
			d.Available = strings.TrimSuffix(strings.TrimPrefix(parts[2], "["), "]")
		}
		if actual.Get(d.Key) != nil || showAll {
			ret = append(ret, d)
		}
	}
	err = loadReferences(ctx, path, ret, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load reference graph")
	}
	return ret, nil
}

func loadReferences(ctx context.Context, path string, deps Dependencies, logger util.Logger) error {
	_, out, err := telemetry.RunProcessSimple(ctx, "go mod graph", path, logger)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		src, dest := util.StringSplit(line, ' ', true)
		if src == "" || dest == "" {
			continue
		}
		src, _ = util.StringSplit(src, '@', true)
		dest, _ = util.StringSplit(dest, '@', true)
		curr := deps.Get(dest)
		if curr != nil {
			curr.AddRef(src)
			// } else {
			//	return errors.Errorf("missing dependency entry for [%s] in path [%s]", dest, path)
		}
	}
	return nil
}

func LoadDepsMap(projects project.Projects, minVersions int, pSvc *project.Service) (map[string]map[string][]string, error) {
	ret := map[string]map[string][]string{}
	for _, prj := range projects {
		deps, err := LoadDepsEasyMode(prj.Key, pSvc.GetFilesystem(prj))
		if err != nil {
			return nil, err
		}
		lo.ForEach(deps, func(dep *Dependency, _ int) {
			curr := lo.ValueOr(ret, dep.Key, map[string][]string{})
			vrs := curr[dep.Version]
			if !lo.Contains(vrs, prj.Key) {
				vrs = append(vrs, prj.Key)
				curr[dep.Version] = vrs
			}
			ret[dep.Key] = curr
		})
	}
	return lo.OmitBy(ret, func(k string, v map[string][]string) bool {
		return len(v) < minVersions
	}), nil
}
