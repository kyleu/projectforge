package database

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const explainPrefix = "{{{ .ExplainPrefix }}}"

func (s *Service) Explain(ctx context.Context, q string, values []any, logger util.Logger) ([]string, error) {
	q = strings.TrimSpace(q)
	if !strings.HasPrefix(q, explainPrefix) {
		q = explainPrefix + q
	}
	res, err := s.db.QueryxContext(ctx, q, values...)
	if err != nil {
		return nil, errors.Wrap(err, "invalid explain result")
	}
	defer func() { _ = res.Close() }()
	var ret []string
	for res.Next() {
		x := map[string]any{}
		err = res.MapScan(x)
		if err != nil {
			return nil, errors.Wrap(err, "can't read results")
		}
		for _, v := range x {
			if s, ok := v.(string); ok {
				ret = append(ret, s)
			}
		}
	}

	return ret, nil
}
