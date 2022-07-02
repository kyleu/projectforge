package data

import (
	"fmt"
	"math"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/project/export/model"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func exportColumn(f *Field, ex []any) (*model.Column, error) {
	typ, tags, err := typeFor(f.Type, ex)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse type for field [%s]", f.Name)
	}

	key := util.StringToSnake(f.Name)
	if key == "type" {
		key = "typ"
	}
	var vals bool
	for _, e := range ex {
		if fmt.Sprint(e) != "" {
			vals = true
			break
		}
	}
	var isPK bool
	if key == "id" {
		if vals {
			isPK = true
		} else {
			key = "empty_id"
		}
	}
	ret := &model.Column{Name: key, Type: typ, PK: isPK, Nullable: false, Search: f.Unique, Tags: tags, HelpString: f.Description}
	return ret, nil
}

func typeFor(t string, ex []any) (*types.Wrapped, []string, error) {
	var tags []string
	isRef := strings.HasPrefix(t, "ref|")
	if isRef {
		tags = append(tags, "reference")
		t = strings.TrimPrefix(t, "ref|")
	}
	isList := strings.HasPrefix(t, "list|")
	if isList {
		tags = append(tags, "list")
		t = strings.TrimPrefix(t, "list|")
		n, nt, err := typeFor(t, ex)
		if err != nil {
			return nil, nil, err
		}
		return types.NewList(n), append(tags, nt...), nil
	}

	switch t {
	case "bool":
		return types.NewBool(), tags, nil
	case "int", "float":
		var isFloat bool
		for _, e := range ex {
			if fl, ok := e.(float64); ok {
				if fl != math.Round(fl) || fl > 1000000000000 {
					isFloat = true
					break
				}
			}
		}
		if isFloat {
			return types.NewFloat(64), tags, nil
		}
		return types.NewInt(64), tags, nil
	case "string":
		return types.NewString(), tags, nil
	case "ulong":
		return types.NewInt(64), tags, nil
	default:
		return types.NewString(), append(tags, "unhandled", t), nil
	}
}
