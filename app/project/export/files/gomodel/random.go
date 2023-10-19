package gomodel

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const nilKey = "nil"

func modelRandom(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"Random", "struct")
	ret.W("func Random() *%s {", m.Proper())
	ret.W("\treturn &%s{", m.Proper())
	maxColLength := m.Columns.MaxCamelLength() + 1
	for _, col := range m.Columns {
		rnd, err := randFor(col, m.PackageWithGroup(""), enums)
		if err != nil {
			return nil, err
		}
		ret.W("\t\t%s %s,", util.StringPad(col.Proper()+":", maxColLength), rnd)
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func randFor(col *model.Column, pkg string, enums enum.Enums) (string, error) {
	switch col.Type.Key() {
	case types.KeyAny:
		return types.KeyNil, nil
	case types.KeyBool:
		return "util.RandomBool()", nil
	case types.KeyEnum:
		et, err := model.AsEnumInstance(col.Type, enums)
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
		e, _ := model.AsEnumInstance(col.Type.ListType(), enums)
		if e != nil {
			return fmt.Sprintf("%s{All%s.Random(), All%s.Random()}", e.ProperPlural(), e.ProperPlural(), e.ProperPlural()), nil
		}
		return nilKey, nil
	case types.KeyMap, types.KeyValueMap:
		return "util.RandomValueMap(4)", nil
	case types.KeyReference:
		return nilKey, nil
	case types.KeyString:
		switch col.Format {
		case model.FmtHTML.Key:
			return "\"<h3>\" + util.RandomString(6) + \"</h3>\"", nil
		case model.FmtIcon.Key:
			return "util.RandomIcon()", nil
		case model.FmtURL.Key:
			return "\"https://\" + util.RandomString(6) + \".com/\" + util.RandomString(6)", nil
		default:
			return "util.RandomString(12)", nil
		}
	case types.KeyDate, types.KeyTimestamp:
		if col.HasTag("deleted") {
			return types.KeyNil, nil
		}
		if col.Nullable {
			return "util.TimeCurrentP()", nil
		}
		return "util.TimeCurrent()", nil
	case types.KeyUUID:
		if col.Nullable {
			return "util.UUIDP()", nil
		}
		return "util.UUID()", nil
	default:
		return "", errors.Errorf("unhandled x type [%s]", col.Type.String())
	}
}
