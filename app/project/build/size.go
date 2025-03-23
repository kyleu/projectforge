package build

import (
	"cmp"
	"context"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func Size(ctx context.Context, fs filesystem.FileLoader, path string, logger util.Logger) (any, []string, error) {
	ex := &ExecHelper{}

	if path == "" {
		cmd := "go build -o ./tmp/size_test"
		_, err := ex.Cmd(ctx, "project build", cmd, fs.Root(), logger)
		if err != nil {
			return nil, ex.Logs, errors.Wrapf(err, "unable to run [%s]", cmd)
		}
		path = "./tmp/size_test"
	}
	cmd := "go tool nm -size " + path
	out, err := ex.Cmd(ctx, "project size", cmd, fs.Root(), logger)
	if err != nil {
		return nil, ex.Logs, errors.Wrapf(err, "unable to run [%s]", cmd)
	}

	lines := util.StringSplitLines(out)

	ret := SizeResultMap{}
	other := &SizeResult{Name: "Other", Type: "X"}
	runtime := &SizeResult{Name: "Runtime", Type: "X"}
	ret.Add("runtime", runtime)
	ret.Add("other", other)
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := util.StringSplitAndTrim(line, " ")
		if len(split) == 3 && split[1] == "U" {
			continue
		}
		if len(split) < 4 {
			return nil, []string{line}, errors.Errorf("found [%d] parts in line [%s]", len(split), line)
		}
		szStr, typ, n := split[1], split[2], strings.Join(split[3:], " ")
		sz, err := util.ParseInt(szStr, "", true)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to process size [%s] line [%s]", szStr, line)
		}
		sr := &SizeResult{Name: n, Type: typ, Size: sz}
		insert := true
		if strings.HasPrefix(n, "$") {
			other.Size += sz
			insert = false
		}
		if strings.HasPrefix(n, "_") {
			runtime.Size += sz
			insert = false
		}
		if strings.HasPrefix(n, "go:") {
			runtime.Size += sz
			insert = false
		}
		if idx := strings.LastIndex(n, ".."); idx > -1 {
			sr.Type, sr.Name = n[:idx], n[idx+2:]
		} else if idx := strings.LastIndex(n, "."); idx > -1 {
			sr.Type, sr.Name = n[:idx], n[idx+1:]
		}
		if insert {
			ret.Add(sr.Type, sr)
		}
	}

	ret = ret.Flatten()
	retSizes := make([]*util.KeyVal[int], 0, len(ret))
	for k, v := range ret.TotalSizes() {
		retSizes = append(retSizes, &util.KeyVal[int]{Key: k, Val: v})
	}
	slices.SortFunc(retSizes, func(l *util.KeyVal[int], r *util.KeyVal[int]) int {
		return cmp.Compare(l.Val, r.Val)
	})
	retSizeStrings := lo.Map(retSizes, func(x *util.KeyVal[int], _ int) string {
		return x.String()
	})
	return retSizeStrings, retSizeStrings, nil
}
