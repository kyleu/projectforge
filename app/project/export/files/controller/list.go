package controller

import (
	"fmt"

	"projectforge.dev/projectforge/app/project/export/files/helper"
	golang2 "projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

func controllerList(m *model2.Model, grp *model2.Column, models model2.Models, g *golang2.File, prefix string) *golang2.Block {
	ret := blockFor(m, prefix, grp, "list")
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
	ret.W("\t\tret, err := as.Services.%s.%s(ps.Context, nil%s, prms%s, ps.Logger)", m.Proper(), meth, grpArgs, suffix)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = %q", m.TitlePlural())
	ret.W("\t\tps.Data = ret")

	var toStrings string
	for _, rel := range m.Relations {
		relModel := models.Get(rel.Table)
		if !relModel.CanTraverseRelation() {
			continue
		}
		srcCol := m.Columns.Get(rel.Src[0])
		for _, imp := range helper.ImportsForTypes("go", srcCol.Type) {
			g.AddImport(imp)
		}

		ret.W("\t\t%sIDs := make([]%s, 0, len(ret))", relModel.Camel(), model2.ToGoType(srcCol.Type, srcCol.Nullable, m.Package))
		ret.W("\t\tfor _, x := range ret {")
		ret.W("\t\t\t%sIDs = append(%sIDs, x.%s)", relModel.Camel(), relModel.Camel(), srcCol.Proper())
		ret.W("\t\t}")
		suffix := ""
		if relModel.IsSoftDelete() {
			suffix = ", false"
		}
		call := "\t\t%s, err := as.Services.%s.GetMultiple(ps.Context, nil%s, ps.Logger, %sIDs...)"
		ret.W(call, relModel.Plural(), relModel.Proper(), suffix, relModel.Camel())
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", err")
		ret.W("\t\t}")

		toStrings += fmt.Sprintf(", %s: %s", relModel.ProperPlural(), relModel.Plural())
	}
	render := "\t\treturn %sRender(rc, as, &v%s.List{Models: ret%s, Params: params}, ps, %s%s)"
	ret.W(render, prefix, m.Package, toStrings, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}
