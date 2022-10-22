package controller

import (
	"fmt"


	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controllerList(m *model.Model, grp *model.Column, models model.Models, enums enum.Enums, g *golang.File, prefix string) (*golang.Block, error) {
	ret := blockFor(m, prefix, grp, "list")
	meth := "List"
	grpArgs := ""
	if grp != nil {
		meth = "GetBy" + grp.Title()
		grpArgs = ", " + grp.Camel() + "Arg"
		controllerArgFor(grp, ret, "\"\"", 2)
	}

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
		for _, imp := range helper.ImportsForTypes("go", enums, srcCol.Type) {
			g.AddImport(imp)
		}
		tgtCol := relModel.Columns.Get(rel.Tgt[0])
		for _, imp := range helper.ImportsForTypes("go", enums, srcCol.Type) {
			g.AddImport(imp)
		}
		gt, err := model.ToGoType(srcCol.Type, srcCol.Nullable, m.Package, enums)
		if err != nil {
			return nil, err
		}

		ret.W("\t\t%sIDs := make([]%s, 0, len(ret))", relModel.Camel(), gt)
		ret.W("\t\tfor _, x := range ret {")
		ret.W("\t\t\t%sIDs = append(%sIDs, x.%s)", relModel.Camel(), relModel.Camel(), srcCol.Proper())
		ret.W("\t\t}")
		suffix := ""
		if relModel.IsSoftDelete() {
			suffix = ", false"
		}
		c := fmt.Sprintf("%sIDs", relModel.Camel())
		if srcCol.Nullable && !tgtCol.Nullable {
			c = "util.ArrayDefererence(" + c + ")"
		}
		call := "\t\t%s, err := as.Services.%s.GetMultiple(ps.Context, nil%s, ps.Logger, %s...)"
		ret.W(call, relModel.Plural(), relModel.Proper(), suffix, c)
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", err")
		ret.W("\t\t}")

		toStrings += fmt.Sprintf(", %s: %s", relModel.ProperPlural(), relModel.Plural())
	}
	render := "\t\treturn %sRender(rc, as, &v%s.List{Models: ret%s, Params: params}, ps, %s%s)"
	ret.W(render, prefix, m.Package, toStrings, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret, nil
}
