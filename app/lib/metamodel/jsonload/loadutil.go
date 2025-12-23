package jsonload

import (
	"encoding/json/jsontext"

	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/util"
)

func fromMD[T any](x map[string]jsontext.Value, k string) (T, error) {
	t, ok := x[k]
	if !ok {
		return *new(T), nil
	}
	ret, err := util.FromJSONObj[T](t)
	if err != nil {
		return *new(T), err
	}
	return ret, nil
}

func extractPath(sch *jsonschema.Schema) (string, string, []string, error) {
	id := util.Str(sch.ID())
	var prefix util.Str
	if idx := id.Index("://"); idx > -1 {
		prefix = id.Substring(0, idx+3)
		id = id.Substring(idx+3, id.Length())
	}
	if id.HasSuffix(".json") {
		id = id.TrimSuffix(".json")
	}
	if id.HasSuffix(".schema") {
		id = id.TrimSuffix(".schema")
	}

	parts := id.Split("/")
	partsLen := len(parts)
	n := parts[partsLen-1]
	var p util.Str
	var g util.Strings
	if partsLen > 1 {
		p = parts[partsLen-2]
		g = parts[:partsLen-2]
		if len(g) > 0 {
			g[0] = prefix + g[0]
		}
	}
	return n.String(), p.String(), g.Strings(), nil
}
