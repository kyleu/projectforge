package gohelper

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

func BlockRandom(cols model.Columns, str StringProvider, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(str.Proper()+"Random", "struct")
	ret.WF("func Random%s() *%s {", str.Proper(), str.Proper())
	ret.WF("\treturn &%s{", str.Proper())
	nd := cols.NotDerived()
	maxColLength := nd.MaxCamelLength() + 1
	for _, col := range nd {
		rnd, err := randForCol(col, str.PackageWithGroup(""), enums)
		if err != nil {
			return nil, err
		}
		ret.WF("\t\t%s %s,", util.StringPad(col.Proper()+":", maxColLength), rnd)
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func BlockArrayRandom(str StringProvider) *golang.Block {
	ret := golang.NewBlock(str.Proper()+"ArrayRandom", "func")
	ret.WF("func (%s %s) Random() *%s {", str.FirstLetter(), str.ProperPlural(), str.Proper())
	ret.WF("\treturn util.RandomElement(%s)", str.FirstLetter())
	ret.W("}")
	return ret
}

func randForCol(col *model.Column, pkg string, enums enum.Enums) (string, error) {
	return randForType(col.Type, col.Format, col.Nullable, col.Tags, pkg, enums)
}

//nolint:gocognit, gocyclo, cyclop
func randForType(t *types.Wrapped, format string, nullable bool, tags []string, pkg string, enums enum.Enums) (string, error) {
	switch t.Key() {
	case types.KeyAny, types.KeyJSON:
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
		if types.Bits(t) == 64 {
			return "util.RandomInt64(10000)", nil
		}
		if types.Bits(t) == 32 {
			return "util.RandomInt32(10000)", nil
		}
		if types.Bits(t) == 16 {
			return "util.RandomInt16(10000)", nil
		}
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
	case types.KeyOrderedMap:
		return "nil", nil
	case types.KeyNumeric:
		return "numeric.Random(4)", nil
	case types.KeyReference:
		return nilKey, nil
	case types.KeyString:
		switch format {
		case model.FmtHTML.Key:
			return fmt.Sprintf("%q + util.RandomString(6) + %q", helper.TextH3Start, helper.TextH3End), nil
		case model.FmtIcon.Key:
			return "util.RandomIcon()", nil
		case model.FmtImage.Key:
			return "\"https://via.placeholder.com/320x180\"", nil
		case model.FmtURL.Key:
			return "util.RandomURL().String()", nil
		default:
			return "util.RandomString(12)", nil
		}
	case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned:
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
