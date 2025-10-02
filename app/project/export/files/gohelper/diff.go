package gohelper

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func DiffBlock(g *golang.File, cols model.Columns, m metamodel.StringProvider, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("Diff"+m.Proper(), "func")

	ret.WF("func (%s *%s) Diff(%sx *%s) util.Diffs {", m.FirstLetter(), m.Proper(), m.FirstLetter(), m.Proper())
	ret.W("\tvar diffs util.Diffs")
	for _, col := range cols.NotDerived() {
		err := columnDiff(g, m, col, enums, ret)
		if err != nil {
			return nil, err
		}
	}
	ret.W("\treturn diffs")
	ret.W("}")
	return ret, nil
}

func columnDiff(g *golang.File, m metamodel.StringProvider, col *model.Column, enums enum.Enums, ret *golang.Block) error {
	key := col.CamelNoReplace()
	if col.HasTag("updated") {
		return nil
	}
	l := fmt.Sprintf("%s.%s", m.FirstLetter(), col.Proper())
	r := fmt.Sprintf("%sx.%s", m.FirstLetter(), col.Proper())
	switch col.Type.Key() {
	case types.KeyAny, types.KeyJSON, types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference, types.KeyNumeric:
		ret.WF("\tdiffs = append(diffs, util.DiffObjects(%s, %s, %q)...)", l, r, col.Camel())
	case types.KeyOrderedMap:
		ret.WF("\tdiffs = append(diffs, util.DiffObjects(%s.Map, %s.Map, %q)...)", l, r, col.Camel())
	case types.KeyBool, types.KeyInt, types.KeyFloat:
		g.AddImport(helper.ImpFmt)
		ret.WF("\tif %s != %s {", l, r)
		ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s)))", key, l, r)
		ret.W("\t}")
	case types.KeyEnum:
		e, err := model.AsEnumInstance(col.Type, enums)
		if err != nil {
			return err
		}
		ret.WF("\tif %s != %s {", l, r)
		if e.Simple() {
			ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, string(%s), string(%s)))", key, l, r)
		} else {
			ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.Key, %s.Key))", key, l, r)
		}
		ret.W("\t}")
	case types.KeyString:
		ret.WF("\tif %s != %s {", l, r)
		ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, %s, %s))", key, l, r)
		ret.W("\t}")
	case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned, types.KeyUUID:
		if col.Nullable {
			msg := "\tif (%s == nil && %s != nil) || (%s != nil && %s == nil) || (%s != nil && %s != nil && *%s != *%s) {"
			line := fmt.Sprintf(msg, l, r, l, r, l, r, l, r)
			ret.W(line)
			g.AddImport(helper.ImpFmt)
			ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s))) //nolint:gocritic // it's nullable", key, l, r)
		} else {
			ret.WF("\tif %s != %s {", l, r)
			ret.WF("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.String(), %s.String()))", key, l, r)
		}
		ret.W("\t}")
	default:
		return errors.Errorf("unhandled diff type [%s]", col.Type.Key())
	}
	return nil
}
