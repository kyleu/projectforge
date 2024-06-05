package gomodel

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func ModelDiff(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"diff")
	g.AddImport(helper.ImpAppUtil)

	mdiff, err := modelDiffBlock(g, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mdiff)

	return g.Render(addHeader, linebreak)
}

func modelDiffBlock(g *golang.File, m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("Diff"+m.Proper(), "func")

	ret.W("func (%s *%s) Diff(%sx *%s) util.Diffs {", m.FirstLetter(), m.Proper(), m.FirstLetter(), m.Proper())
	ret.W("\tvar diffs util.Diffs")
	for _, col := range m.Columns.NotDerived() {
		if col.HasTag("updated") {
			continue
		}
		l := fmt.Sprintf("%s.%s", m.FirstLetter(), col.Proper())
		r := fmt.Sprintf("%sx.%s", m.FirstLetter(), col.Proper())
		switch col.Type.Key() {
		case types.KeyAny, types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
			ret.W("\tdiffs = append(diffs, util.DiffObjects(%s, %s, %q)...)", l, r, col.Camel())
		case types.KeyBool, types.KeyInt, types.KeyFloat:
			g.AddImport(helper.ImpFmt)
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s)))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyEnum:
			e, err := model.AsEnumInstance(col.Type, enums)
			if err != nil {
				return nil, err
			}
			ret.W("\tif %s != %s {", l, r)
			if e.Simple() {
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, string(%s), string(%s)))", col.Camel(), l, r)
			} else {
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.Key, %s.Key))", col.Camel(), l, r)
			}
			ret.W("\t}")
		case types.KeyString:
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s, %s))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyDate, types.KeyTimestamp, types.KeyUUID:
			if col.Nullable {
				msg := "\tif (%s == nil && %s != nil) || (%s != nil && %s == nil) || (%s != nil && %s != nil && *%s != *%s) {"
				line := fmt.Sprintf(msg, l, r, l, r, l, r, l, r)
				ret.W(line)
				g.AddImport(helper.ImpFmt)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s))) //nolint:gocritic // it's nullable", col.Camel(), l, r)
			} else {
				ret.W("\tif %s != %s {", l, r)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.String(), %s.String()))", col.Camel(), l, r)
			}
			ret.W("\t}")
		default:
			return nil, errors.Errorf("unhandled diff type [%s]", col.Type.Key())
		}
	}
	ret.W("\treturn diffs")
	ret.W("}")
	return ret, nil
}
