package gomodel

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func modelFromMap(g *golang.File, m *model.Model, enums enum.Enums, database string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	cols := m.Columns.WithoutTags("created", "updated", model.RevisionType)
	needsErr := lo.ContainsBy(cols, func(c *model.Column) bool {
		return c.NeedsErr(m.Name, database)
	})
	if needsErr {
		ret.W("\tvar err error")
	}
	ret.W("\tif setPK {")
	err := forCols(g, ret, 2, m, enums, cols.PKs()...)
	if err != nil {
		return nil, err
	}
	ret.W("\t\t// $PF_SECTION_START(pkchecks)$")
	ret.W("\t\t// $PF_SECTION_END(pkchecks)$")
	ret.W("\t}")
	err = forCols(g, ret, 1, m, enums, cols.NonPKs()...)
	if err != nil {
		return nil, err
	}
	ret.W("\t// $PF_SECTION_START(extrachecks)$")
	ret.W("\t// $PF_SECTION_END(extrachecks)$")
	ret.W("\treturn ret, nil")
	ret.W("}")

	return ret, nil
}

func forCols(g *golang.File, ret *golang.Block, indent int, m *model.Model, enums enum.Enums, cols ...*model.Column) error {
	ind := util.StringRepeat("\t", indent)
	catchErr := func(s string) {
		ret.W(ind + "if " + s + " != nil {")
		ret.W(ind + "\treturn nil, " + s)
		ret.W(ind + "}")
	}
	for _, col := range cols {
		switch {
		case col.Type.Key() == types.KeyAny:
			ret.W(ind+"ret.%s = m[%q]", col.Proper(), col.Camel())
		case col.Type.Key() == types.KeyReference:
			ret.W(ind+"tmp%s, err := m.ParseString(%q, true, true)", col.Proper(), col.Camel())
			catchErr("err")

			ref, err := model.AsRef(col.Type)
			if err != nil {
				return errors.Wrap(err, "invalid ref")
			}
			if ref.Pkg.Last() == m.Package {
				ret.W("\t%sArg := &%s{}", col.Camel(), ref.K)
			} else {
				ret.W("\t%sArg := &%s.%s{}", col.Camel(), ref.Pkg.Last(), ref.K)
				g.AddImport(golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
			ret.W(ind+"err = util.FromJSON([]byte(tmp%s), %sArg)", col.Proper(), col.Camel())
			catchErr("err")
			ret.W(ind+"ret.%s = %sArg", col.Proper(), col.Camel())
		case col.Type.Key() == types.KeyEnum:
			e, err := model.AsEnumInstance(col.Type, enums)
			if err != nil {
				return err
			}
			ret.W(ind+"ret%s, err := m.Parse%s(%q, true, true)", col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("err")
			enumRef := ""
			if e.Simple() {
				enumRef = fmt.Sprintf("%s(ret%s)", e.Proper(), col.Proper())
			} else {
				enumRef = fmt.Sprintf("All%s.Get(ret%s, nil)", e.ProperPlural(), col.Proper())
			}
			if e.PackageWithGroup("") == m.PackageWithGroup("") {
				ret.W(ind+"ret.%s = %s", col.Proper(), enumRef)
			} else {
				ret.W(ind+"ret.%s = %s.%s", col.Proper(), e.Package, enumRef)
			}
		case col.Nullable || col.Type.Scalar():
			ret.W(ind+"ret.%s, err = m.Parse%s(%q, true, true)", col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("err")
		default:
			ret.W(ind+"ret%s, e := m.Parse%s(%q, true, true)", col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("e")
			ret.W(ind+"if ret%s != nil {", col.Proper())
			ret.W(ind+"\tret.%s = *ret%s", col.Proper(), col.Proper())
			ret.W(ind + "}")
		}
	}
	return nil
}
