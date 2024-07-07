package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) History(ctx context.Context, args *HistoryArgs, logger util.Logger) (*Result, error) {
	hist, err := gitHistory(ctx, s.Path, args, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve history")
	}
	return NewResult(s.Key, ok, util.ValueMap{"history": hist}), nil
}

const hFDelimit, hLDelimit = "»¦«", "»¦¦¦«"

var historyFormat = fmt.Sprintf("%%H%s%%an%s%%ae%s%%cd%s%%B%s", hFDelimit, hFDelimit, hFDelimit, hFDelimit, hLDelimit)

func gitHistory(ctx context.Context, path string, args *HistoryArgs, logger util.Logger) (*HistoryResult, error) {
	cmd := "log --pretty=format:" + historyFormat
	if args.Since != nil {
		cmd += fmt.Sprintf(" --since %v", args.Since)
	}
	if args.Limit > 0 {
		cmd += fmt.Sprintf(" --max-count %d", args.Limit)
	}
	lo.ForEach(args.Authors, func(author string, _ int) {
		cmd += fmt.Sprintf(" --author %s", author)
	})
	if args.Path != "" {
		cmd += fmt.Sprintf(" -- %s", path)
	}

	out, err := gitCmd(ctx, cmd, path, logger)
	if err != nil {
		if isNoRepo(err) {
			return nil, nil
		}
		return nil, err
	}
	res, err := ParseResultsDelimited(out)
	if err != nil {
		return nil, err
	}
	var dbg any
	if args.Debug {
		dbg = out
	}
	return &HistoryResult{Args: args, Entries: res, Debug: dbg}, nil
}

func ParseResultsDelimited(output string) (HistoryEntries, error) {
	var commits HistoryEntries

	lines := util.StringSplitAndTrim(output, "»¦¦¦«")
	for _, line := range lines {
		parts := strings.Split(line, "»¦«")
		if len(parts) != 5 {
			return nil, errors.Errorf("line [%s] only has [%d] parts", line, len(parts))
		}
		occ, err := util.TimeFromVerbose(parts[3])
		if err != nil {
			return nil, err
		}
		commits = append(commits, &HistoryEntry{SHA: parts[0], AuthorName: parts[1], AuthorEmail: parts[2], Message: parts[4], Occurred: *occ})
	}
	return commits, nil
}
