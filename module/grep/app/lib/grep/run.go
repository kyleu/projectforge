package grep

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func Run(ctx context.Context, req *Request, logger util.Logger) (*Response, error) {
	cmd := req.ToCommand()

	_, out, err := util.RunProcessSimple(cmd, ".")
	if err != nil {
		return nil, err
	}

	lines := util.StringSplitLines(out)
	ret := &Response{Request: req}
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, e := util.FromJSONMap([]byte(line))
		if e != nil {
			return nil, e
		}
		data := m.GetMapOpt("data")
		switch m.GetStringOpt("type") {
		case "begin", "end":
			// noop
		case "match":
			subMatches, _ := data.GetMapArray("submatches", true)
			for _, x := range subMatches {
				ret.Matches = append(ret.Matches, &Match{
					File:    strings.TrimPrefix(strings.TrimPrefix(data.GetMapOpt("path").GetStringOpt("text"), req.Path), "/"),
					Offset:  data.GetIntOpt("absolute_offset"),
					LineNum: data.GetIntOpt("line_number"),
					Text:    data.GetMapOpt("lines").GetStringOpt("text"),
					Match:   x.GetMapOpt("match").GetStringOpt("text"),
					Start:   x.GetIntOpt("start"),
					End:     x.GetIntOpt("end"),
				})
			}
		case "summary":
			stats := data.GetMapOpt("stats")
			ret.BytesSearched = stats.GetIntOpt("bytes_searched")
			ret.ElapsedNanos = stats.GetMapOpt("elapsed").GetIntOpt("nanos")
		default:
			return nil, errors.Errorf("unhandled type [%s]", m.GetStringOpt("type"))
		}
	}

	return ret, nil
}
