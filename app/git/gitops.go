package git

import (
	"sort"
	"strings"

	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

const (
	gitUnknown = "?? "
	gitMerge   = "M "
	gitDeleted = "D "
	gitMoved   = "AM "
)

func gitBranch(path string) string {
	out, err := gitCmd("branch --show-current", path)
	if err != nil {
		if errors.Is(err, noRepo) {
			return "norepo"
		}
		return "error: " + err.Error()
	}
	return strings.TrimSpace(out)
}

func gitStatus(path string) ([]string, error) {
	out, err := gitCmd("status --porcelain", path)
	if err != nil {
		if errors.Is(err, noRepo) {
			return nil, nil
		}
		return nil, err
	}

	lines := util.SplitAndTrim(out, "\n")

	dirty := make([]string, 0, len(lines))
	for _, line := range lines {
		if i := strings.Index(line, " "); i > -1 {
			line = line[i+1:]
		}
		dirty = append(dirty, strings.TrimSpace(line))
	}
	sort.Strings(dirty)

	return dirty, nil
}
