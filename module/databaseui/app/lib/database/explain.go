package database

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const explainPrefix = "explain "

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
		if len(x) != 1 {
			return nil, errors.New("return from explain contains more than one column")
		}
		for _, v := range x {
			s, ok := v.(string)
			if !ok {
				return nil, errors.Errorf("explain column is [%T], not [string]", v)
			}
			ret = append(ret, s)
		}
	}

	return ret, nil
}
