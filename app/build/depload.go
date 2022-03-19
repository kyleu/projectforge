package build

import (
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/util"
)

func LoadDeps(path string) (Dependencies, error) {
	_, out, err := util.RunProcessSimple("go list -m -u all", path)
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
