package gohelper

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

//nolint:gocognit
func forMapCol(g *golang.File, ret *golang.Block, indent int, m model.StringProvider, args *metamodel.Args, col *model.Column) error {
	ind := util.StringRepeat("\t", indent)
	catchErr := func(s string) {
		ret.W(ind + "if " + s + " != nil {")
		ret.W(ind + "\treturn nil, nil, " + s)
		ret.W(ind + "}")
	}
	colMP := func(msg string) error {
		mp, err := col.ToGoMapParse()
		if err != nil {
			return err
		}
		ret.WF(msg, col.Proper(), mp)
		return nil
	}
	parseCall := "ret%s, err := m.Parse%s(k, true, true)"
	parseMsg := "ret.%s, err = m.Parse%s(k, true, true)"
	switch {
	case col.Type.Key() == types.KeyAny:
		ret.WF(ind+"ret.%s = m[%q]", col.Proper(), col.Camel())
	case col.Type.Key() == types.KeyReference:
		ret.WF(ind+"tmp%s, err := m.ParseString(%q, true, true)", col.Proper(), col.Camel())
		catchErr("err")
		ref, _, err := helper.LoadRef(col, args.Models, args.Events, args.ExtraTypes)
		if err != nil {
			return errors.Wrap(err, "invalid ref")
		}
		ret.WF(ind+"%sArg := %s{}", col.Camel(), ref.LastAddr(ref.Pkg.Last() != m.PackageName()))
		ret.WF(ind+"err = util.FromJSON([]byte(tmp%s), %sArg)", col.Proper(), col.Camel())
		catchErr("err")
		ret.WF(ind+"ret.%s = %sArg", col.Proper(), col.Camel())
	case col.Type.Key() == types.KeyEnum:
		e, err := model.AsEnumInstance(col.Type, args.Enums)
		if err != nil {
			return err
		}
		if err := colMP(ind + parseCall); err != nil {
			return err
		}
		catchErr("err")
		var enumRef string
		if e.Simple() {
			enumRef = fmt.Sprintf("%s(ret%s)", e.Proper(), col.Proper())
		} else {
			enumRef = fmt.Sprintf("All%s.Get(ret%s, nil)", e.ProperPlural(), col.Proper())
		}
		if e.PackageWithGroup("") == m.PackageWithGroup("") {
			ret.WF(ind+"ret.%s = %s", col.Proper(), enumRef)
		} else {
			ret.WF(ind+"ret.%s = %s.%s", col.Proper(), e.Package, enumRef)
		}
	case col.Type.Key() == types.KeyJSON:
		if err := colMP(ind + "ret.%s, err = m.Parse%s(k, true, true)"); err != nil {
			return err
		}
	case col.Type.Key() == types.KeyList:
		if e, _ := model.AsEnumInstance(col.Type.ListType(), args.Enums); e != nil {
			if err := colMP(ind + parseCall); err != nil {
				return err
			}
			catchErr("err")
			eRef := e.Proper()
			if e.PackageWithGroup("") != m.PackageWithGroup("") {
				eRef = e.Package + "." + eRef
			}
			ret.WF(ind+"ret.%s = %sParse(nil, ret%s...)", col.Proper(), eRef, col.Proper())
		} else {
			if err := colMP(ind + parseMsg); err != nil {
				return err
			}
			catchErr("err")
		}
	case col.Type.Key() == types.KeyNumeric:
		g.AddImport(helper.ImpAppNumeric)
		ret.WF(ind+"ret.%s, err = numeric.FromAny(v)", col.Proper())
	case col.Nullable || col.Type.Scalar():
		if err := colMP(ind + parseMsg); err != nil {
			return err
		}
	default:
		if err := colMP(ind + "ret%s, e := m.Parse%s(k, true, true)"); err != nil {
			return err
		}
		catchErr("e")
		ret.WF(ind+"if ret%s != nil {", col.Proper())
		ret.WF(ind+"\tret.%s = *ret%s", col.Proper(), col.Proper())
		ret.W(ind + "}")
	}
	return nil
}
