package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

type Coverage struct {
	Packages util.ValueMap `json:"packages"`
	SVG      string        `json:"svg,omitzero"`
}

func runCoverage(ctx context.Context, fs filesystem.FileLoader, scope string, logger util.Logger) (*Coverage, []string, error) {
	ex := &build.ExecHelper{}
	var logs []string

	if scope == "" {
		scope = "./app/..."
	}

	testCmd := "go test -race -coverprofile ./tmp/coverage.out " + scope
	testOut, err := ex.Cmd(ctx, "coverage-test", testCmd, fs.Root(), false, logger)
	if err != nil {
		return nil, ex.Logs, errors.Wrapf(err, "unable to run [%s] for test coverage", testCmd)
	}
	testLines := util.StringSplitLines(testOut)
	testMap := make(util.ValueMap, len(testLines))
	for _, x := range testLines {
		p, c := util.StringCut(x, ':', true)
		p = strings.TrimPrefix(p, "ok")
		p = strings.TrimSuffix(p, "coverage")
		p = strings.TrimSpace(p)
		c = strings.TrimSuffix(c, " of statements")
		if p != "" {
			testMap[p] = c
		}
	}

	ret := &Coverage{Packages: testMap}

	treemapCmd := "go-cover-treemap -padding 0 -coverprofile ./tmp/coverage.out > ./tmp/coverage.svg"
	if _, err := ex.Cmd(ctx, "coverage-run", treemapCmd, fs.Root(), false, logger); err != nil {
		return nil, ex.Logs, errors.Wrapf(err, "unable to run [%s] for treemap generation", testCmd)
	}

	b, err := fs.ReadFile("./tmp/coverage.svg")
	if err != nil {
		return nil, ex.Logs, errors.Wrap(err, "unable to read SVG from treemap generation")
	}
	ret.SVG = string(b)

	return ret, logs, nil
}
