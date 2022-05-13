package controller

import (
	"fmt"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func controllerList(m *model.Model, grp *model.Column, models model.Models, g *golang.File) *golang.Block {
	ret := blockFor(m, grp, "list")
	meth := "List"
	grpArgs := ""
	if grp != nil {
		meth = "GetBy" + grp.Title()
		grpArgs = ", " + grp.Camel() + "Arg"
		controllerArgFor(grp, ret, "\"\"", 2)
	}

	ret.W("\t\tps.Title = %sDefaultTitle", m.Camel())
	ret.W("\t\tparams := cutil.ParamSetFromRequest(rc)")
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", " + incDel
	}
	ret.W("\t\tprms := params.Get(%q, nil, ps.Logger).Sanitize(%q)", m.Package, m.Package)
	ret.W("\t\tret, err := as.Services.%s.%s(ps.Context, nil%s, prms%s)", m.Proper(), meth, grpArgs, suffix)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")

	var toStrings string
	for _, rel := range m.Relations {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			srcCol := m.Columns.Get(rel.Src[0])

			//g.AddImport(helper.AppImport("app/" + relModel.Package))
			//ret.W("\t\tvar %s %s.%s", relModel.Plural(), relModel.Package, relModel.ProperPlural())

			for _, imp := range helper.ImportsForTypes("go", srcCol.Type) {
				g.AddImport(imp)
			}

			ret.W("\t\t%sIDs := make([]%s, 0, len(ret))", relModel.Camel(), model.ToGoType(srcCol.Type, srcCol.Nullable, m.Package))
			ret.W("\t\tfor _, x := range ret {")
			ret.W("\t\t\t%sIDs = append(%sIDs, x.%s)", relModel.Camel(), relModel.Camel(), srcCol.Proper())
			ret.W("\t\t}")
			suffix := ""
			if relModel.IsSoftDelete() {
				suffix = ", false"
			}
			ret.W("\t\t%s, err := as.Services.%s.GetMultiple(ps.Context, nil%s, %sIDs...)", relModel.Plural(), relModel.Proper(), suffix, relModel.Camel())
			ret.W("\t\tif err != nil {")
			ret.W("\t\t\treturn \"\", err")
			ret.W("\t\t}")

			toStrings += fmt.Sprintf(", %s: %s", relModel.ProperPlural(), relModel.Plural())
		}
	}
	ret.W("\t\treturn render(rc, as, &v%s.List{Models: ret%s, Params: params}, ps, %q%s)", m.Package, toStrings, m.Package, grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}
