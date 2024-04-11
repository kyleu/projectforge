package database

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func QueryFieldDescs(fields util.FieldDescs, query string, argOffset int) (string, []any, error) {
	if !strings.Contains(query, ":") {
		return "", nil, errors.New("search query [%s] does not contain [:], and should not be passed to this method")
	}
	parts := util.StringSplitAndTrim(query, ",")
	if len(parts) == 0 {
		return "", nil, nil
	}

	var wcs []string
	var vals []any
	for _, part := range parts {
		col, q := util.StringSplit(part, ':', true)
		if q == "" {
			continue
		}
		desc := fields.Get(col)
		if desc == nil {
			return "", nil, errors.Errorf("invalid field [%s]", col)
		}

		var v any
		switch desc.Type {
		case "string":
			v = "%" + q + "%"
		default:
			var err error
			v, err = desc.Parse(q)
			if err != nil {
				return "", nil, errors.Wrapf(err, "invalid search term [%s] for field [%s] of type [%s]", q, desc.Key, desc.Type)
			}
		}
		if v == nil {
			return "", nil, errors.Errorf("unable to parse [%s] value for search field [%s]", desc.Type, desc.Key)
		}
		wc := fmt.Sprintf("%s like $%d", desc.Key, len(wcs)+1+argOffset)
		wcs = append(wcs, wc)
		vals = append(vals, v)
	}
	wc := "(" + strings.Join(wcs, " and ") + ")"
	return wc, vals, nil
}
