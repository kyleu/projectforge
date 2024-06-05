package gomodel

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const nilKey = "nil"

func modelRandom(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"Random", "struct")
	ret.W("func Random() *%s {", m.Proper())
	ret.W("\treturn &%s{", m.Proper())
	cols := m.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength() + 1
	for _, col := range cols {
		rnd, err := randForCol(col, m.PackageWithGroup(""), enums)
		if err != nil {
			return nil, err
		}
		ret.W("\t\t%s %s,", util.StringPad(col.Proper()+":", maxColLength), rnd)
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func randForCol(col *model.Column, pkg string, enums enum.Enums) (string, error) {
	return randForType(col.Type, col.Format, col.Nullable, col.Tags, pkg, enums)
}

func randForType(t *types.Wrapped, format string, nullable bool, tags []string, pkg string, enums enum.Enums) (string, error) {
	switch t.Key() {
	case types.KeyAny:
		return types.KeyNil, nil
	case types.KeyBool:
		return "util.RandomBool()", nil
	case types.KeyEnum:
		et, err := model.AsEnumInstance(t, enums)
		if err != nil {
			return "", err
		}
		if pkg == et.PackageWithGroup("") {
			if et.Simple() {
				return fmt.Sprintf("%s(util.RandomString(12))", et.Proper()), nil
			}
			return fmt.Sprintf("All%s.Random()", et.ProperPlural()), nil
		}
		if et.Simple() {
			return fmt.Sprintf("%s(util.RandomString(12))", et.Proper()), nil
		}
		return fmt.Sprintf("%s.All%s.Random()", et.Package, et.ProperPlural()), nil
	case types.KeyInt:
		return "util.RandomInt(10000)", nil
	case types.KeyFloat:
		return "util.RandomFloat(1000)", nil
	case types.KeyList:
		lt := t.ListType()
		switch lt.Key() {
		case types.KeyString, types.KeyInt, types.KeyFloat:
			x, _ := randForType(lt, "", false, nil, pkg, enums)
			return fmt.Sprintf("[]%s{%s, %s}", lt.Key(), x, x), nil
		case types.KeyEnum:
			e, _ := model.AsEnumInstance(lt, enums)
			if e != nil {
				eRef := e.ProperPlural()
				xRef := "All" + e.ProperPlural()
				if e.PackageWithGroup("") != pkg {
					eRef = e.Package + "." + eRef
					xRef = e.Package + "." + xRef
				}

				return fmt.Sprintf("%s{%s.Random(), %s.Random()}", eRef, xRef, xRef), nil
			}
			return nilKey, nil
		}
		return nilKey, nil
	case types.KeyMap, types.KeyValueMap:
		return "util.RandomValueMap(4)", nil
	case types.KeyReference:
		return nilKey, nil
	case types.KeyString:
		switch format {
		case model.FmtHTML.Key:
			return fmt.Sprintf("%q + util.RandomString(6) + %q", helper.TextH3Start, helper.TextH3End), nil
		case model.FmtIcon.Key:
			return "util.RandomIcon()", nil
		case model.FmtImage.Key:
			return "\"http://via.placeholder.com/320x180\"", nil
		case model.FmtURL.Key:
			return "util.RandomURL().String()", nil
		default:
			return "util.RandomString(12)", nil
		}
	case types.KeyDate, types.KeyTimestamp:
		if slices.Contains(tags, "deleted") {
			return types.KeyNil, nil
		}
		if nullable {
			return "util.TimeCurrentP()", nil
		}
		return "util.TimeCurrent()", nil
	case types.KeyUUID:
		if nullable {
			return "util.UUIDP()", nil
		}
		return "util.UUID()", nil
	default:
		return "", errors.Errorf("unhandled random type [%s]", t.String())
	}
}
