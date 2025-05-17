package data

import (
	"fmt"
	"math"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func exportColumn(f *Field, ex []any) (*model.Column, error) {
	typ, tags, err := TypeForString(f.Type, ex)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse type for field [%s]", f.Name)
	}

	key := util.StringToSnake(f.Name)
	if key == "type" {
		key = "typ"
	}
	vals := lo.ContainsBy(ex, func(e any) bool {
		return fmt.Sprint(e) != ""
	})
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

func TypeForString(t string, ex []any) (*types.Wrapped, []string, error) {
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
		n, nt, err := TypeForString(t, ex)
		if err != nil {
			return nil, nil, err
		}
		return types.NewList(n), append(tags, nt...), nil
	}

	switch t {
	case types.KeyBool:
		return types.NewBool(), tags, nil
	case types.KeyInt, types.KeyFloat:
		isFloat := lo.ContainsBy(ex, func(e any) bool {
			if fl, err := util.Cast[float64](e); err == nil {
				if fl != math.Round(fl) || fl > 1000000000000 {
					return true
				}
			}
			return false
		})
		if isFloat {
			return types.NewFloat(64), tags, nil
		}
		return types.NewInt(64), tags, nil
	case types.KeyString:
		return types.NewString(), tags, nil
	case "ulong":
		return types.NewInt(64), tags, nil
	default:
		return types.NewString(), append(tags, "unhandled", t), nil
	}
}
