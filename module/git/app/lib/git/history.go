package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) History(ctx context.Context, hist *HistoryResult, logger util.Logger) (*Result, error) {
	err := gitHistory(ctx, s.Path, hist, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve history")
	}
	return NewResult(s.Key, ok, util.ValueMap{"history": hist}), nil
}

const hFDelimit, hLDelimit = "»¦«", "»¦¦¦«"

var historyFormat = fmt.Sprintf("%%H%s%%an%s%%ae%s%%cd%s%%B%s", hFDelimit, hFDelimit, hFDelimit, hFDelimit, hLDelimit)

func gitHistory(ctx context.Context, path string, hist *HistoryResult, logger util.Logger) error {
	// if hist.Commit != "" {
	//	curr := &HistoryEntry{SHA: hist.Commit}
	//	curr.Files = HistoryFiles{
	//		{Status: "OK", File: "foo.txt"},
	//	}
	//	hist.Entries = append(hist.Entries, curr)
	//	return nil
	//}

	cmd := "log --pretty=format:" + historyFormat
	if hist.Since != nil {
		cmd += fmt.Sprintf(" --since %v", hist.Since)
	}
	if hist.Limit > 0 {
		cmd += fmt.Sprintf(" --max-count %d", hist.Limit)
	}
	lo.ForEach(hist.Authors, func(author string, _ int) {
		cmd += fmt.Sprintf(" --author %s", author)
	})
	if hist.Path != "" {
		cmd += fmt.Sprintf(" -- %s", path)
	}

	out, err := gitCmd(ctx, cmd, path, logger)
	if err != nil {
		if isNoRepo(err) {
			return nil
		}
		return err
	}
	// hist.Debug = out
	res, err := ParseResultsDelimited(out)
	if err != nil {
		return err
	}
	hist.Entries = res

	return nil
}

func ParseResultsDelimited(output string) (HistoryEntries, error) {
	var commits HistoryEntries

	lines := util.StringSplitAndTrim(output, "»¦¦¦«")
	for _, line := range lines {
		parts := strings.Split(line, "»¦«")
		if len(parts) != 5 {
			return nil, errors.Errorf("line [%s] only has [%d] parts", line, len(parts))
		}
		commits = append(commits, &HistoryEntry{
			SHA:         parts[0],
			AuthorName:  parts[1],
			AuthorEmail: parts[2],
			Message:     parts[4],
			Occurred:    parts[3],
		})
	}
	return commits, nil
}
