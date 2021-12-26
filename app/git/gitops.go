package git

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func gitStatus(path string) ([]string, error) {
	out, err := gitCmd("status --porcelain", path)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return nil, nil
		}
		return nil, err
	}

	lines := util.StringSplitAndTrim(out, "\n")

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

func gitBranch(path string) string {
	out, err := gitCmd("branch --show-current", path)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "norepo"
		}
		return "error: " + err.Error()
	}
	return strings.TrimSpace(out)
}

func gitFetch(path string, dryRun bool) (string, error) {
	cmd := "fetch"
	if dryRun {
		cmd += " --dry-run"
	}
	out, err := gitCmd(cmd, path)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}

func gitCommit(path string, message string) (string, error) {
	_, err := gitCmd("add .", path)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "", nil
		}
		return "", err
	}
	out, err := gitCmd(fmt.Sprintf(`commit -m %q`, message), path)
	if err != nil {
		return "", err
	}
	return out, nil
}
