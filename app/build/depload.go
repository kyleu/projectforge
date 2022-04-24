package build

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/project"

	"projectforge.dev/projectforge/app/util"
)

func LoadDeps(path string, includeUpdates bool) (Dependencies, error) {
	cmd := "go list -m all"
	if includeUpdates {
		cmd = "go list -m -u all"
	}
	_, out, err := util.RunProcessSimple(cmd, path)
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
		ret = append(ret, d)
	}
	err = loadReferences(path, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load reference graph")
	}

	return ret, nil
}

func loadReferences(path string, deps Dependencies) error {
	_, out, err := util.RunProcessSimple("go mod graph", path)
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
		if curr == nil {
			return errors.Errorf("missing dependency entry for [%s]", dest)
		}
		curr.AddRef(src)
	}
	return nil
}

func LoadDepsMap(projects project.Projects, minVersions int) (map[string]map[string][]string, error) {
	ret := map[string]map[string][]string{}
	for _, prj := range projects {
		deps, err := LoadDeps(prj.Path, false)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		for _, dep := range deps {
			curr, ok := ret[dep.Key]
			if !ok {
				curr = map[string][]string{}
			}
			vrs := curr[dep.Version]
			if !slices.Contains(vrs, prj.Key) {
				vrs = append(vrs, prj.Key)
				curr[dep.Version] = vrs
			}
			ret[dep.Key] = curr
		}
	}
	for k, v := range ret {
		if len(v) < minVersions {
			delete(ret, k)
		}
	}
	return ret, nil
}
