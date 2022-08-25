package data

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Derive(name string, pkg string, content string) (*model.Model, error) {
	println(util.ToJSON(types.NewList(types.NewString())))
	var js any
	err := util.FromJSON([]byte(content), &js)
	if err == nil {
		return deriveJSON(name, pkg, js)
	}
	return nil, errors.New("unhandled content")
}

func deriveJSON(name string, pkg string, js any) (*model.Model, error) {
	switch t := js.(type) {
	case map[string]any:
		return deriveJSONObject(name, pkg, t)
	default:
		return nil, errors.Errorf("unhandled type [%T]: %v", js, t)
	}
}

func deriveJSONObject(name string, pkg string, m map[string]any) (*model.Model, error) {
	if pkg == "" {
		pkg = name
	}
	cols, seedData, err := deriveExtractColumns("", m)
	if err != nil {
		return nil, err
	}
	ret := &model.Model{Name: name, Package: pkg, Icon: "star", Columns: cols, SeedData: seedData}
	return ret, nil
}

func deriveExtractColumns(prefix string, childMap util.ValueMap) (model.Columns, [][]any, error) {
	var ret model.Columns
	var sd [][]any
	for k, v := range childMap {
		col := &model.Column{Name: k, PK: k == "id"}
		switch t := v.(type) {
		case float64:
			if t == float64(int64(t)) {
				col.Type = types.NewInt(64)
			} else {
				col.Type = types.NewFloat(64)
			}
		case string:
			col.Type = types.NewString()
		case map[string]any:
			newCols, newSeed, err := deriveExtractColumns(k, t)
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, newCols...)
			sd = append(sd, newSeed...)
		default:
			return nil, nil, errors.Errorf("unhandled type [%T]: %v", v, t)
		}
		if col.Type != nil {
			ret = append(ret, col)
		}
	}
	return ret, sd, nil
}
