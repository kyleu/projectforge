package gomodel

import (
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func modelFromMap(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	ret.W("\tvar err error")
	ret.W("\tif setPK {")
	cols := m.Columns.WithoutTag("created").WithoutTag("updated").WithoutTag(model.RevisionType)
	forCols(ret, 2, cols.PKs()...)
	ret.W("\t\t// $PF_SECTION_START(pkchecks)$")
	ret.W("\t\t// $PF_SECTION_END(pkchecks)$")
	ret.W("\t}")
	forCols(ret, 1, cols.NonPKs()...)
	ret.W("\t// $PF_SECTION_START(extrachecks)$")
	ret.W("\t// $PF_SECTION_END(extrachecks)$")
	ret.W("\treturn ret, nil")
	ret.W("}")

	return ret
}

func forCols(ret *golang.Block, indent int, cols ...*model.Column) {
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
			ret.W(ind+"tmp%s, err := m.ParseMap(%q, true, true)", col.Proper(), col.Camel())
			catchErr("err")
			ret.W(ind+"err = util.CycleJSON(tmp%s, ret.%s)", col.Proper(), col.Proper())
			catchErr("err")
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
}
