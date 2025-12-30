package controller

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func controllerList(g *golang.File, m *model.Model, grp *model.Column, models model.Models, enums enum.Enums, prefix string) (*golang.Block, error) {
	ret := blockFor(m, prefix, grp, "list")
	meth := "List"
	var grpArgs string
	if grp != nil {
		meth = helper.TextGetBy + grp.Title()
		grpArgs = helper.TextCommaSpace + grp.Camel() + helper.TextArg
		controllerArgFor(grp, ret, `""`, 2)
	}

	var suffix string
	if m.IsSoftDelete() {
		suffix = helper.TextCommaSpace + incDel
	}
	if m.HasSearches() {
		g.AddImport(helper.ImpStrings)
		ret.W("\t\tq := strings.TrimSpace(cutil.QueryStringString(ps.URI, \"q\"))")
	}
	ret.WF("\t\tprms := ps.Params.Sanitized(%q, ps.Logger)", m.Package)
	if grpArgs == "" && m.HasSearches() {
		ret.WF("\t\tvar ret %s.%s", m.Package, m.ProperPlural())
		ret.W("\t\tvar err error")
		ret.W("\t\tif q == \"\" {")
		ret.WF("\t\t\tret, err = as.Services.%s.%s(ps.Context, nil, prms%s, ps.Logger)", m.Proper(), meth, suffix)
		ret.WE(3, `""`)
		ret.W("\t\t} else {")
		ret.WF("\t\t\tret, err = as.Services.%s.Search(ps.Context, q, nil, prms%s, ps.Logger)", m.Proper(), suffix)
		ret.WE(3, `""`)
		ret.W("\t\t\tif len(ret) == 1 && !cutil.IsContentTypeJSON(cutil.GetContentType(r)) {")
		if len(m.Group) == 0 {
			ret.W("\t\t\t\treturn FlashAndRedir(true, \"single result found\", ret[0].WebPath(), ps)")
		} else {
			g.AddImport(helper.ImpAppController)
			ret.W("\t\t\t\treturn controller.FlashAndRedir(true, \"single result found\", ret[0].WebPath(), ps)")
		}
		ret.W("\t\t\t}")
		ret.W("\t\t}")
	} else {
		ret.WF("\t\tret, err := as.Services.%s.%s(ps.Context, nil%s, prms%s, ps.Logger)", m.Proper(), meth, grpArgs, suffix)
		ret.WE(2, `""`)
	}
	if m.HasTag("count") {
		ret.WF("\t\tcount, err := as.Services.%s.Count(ps.Context, nil%s, \"\"%s, ps.Logger)", m.Proper(), grpArgs, suffix)
		ret.WE(2, `""`)
	}
	ret.WF("\t\tps.SetTitleAndData(%q, ret)", m.TitlePlural())
	var toStrings string
	for _, rel := range m.Relations {
		relModel := models.Get(rel.Table)
		if !relModel.CanTraverseRelation() {
			continue
		}
		srcCol := m.Columns.Get(rel.Src[0])
		lo.ForEach(helper.ImportsForTypes("go", "", srcCol.Type), func(imp *model.Import, _ int) {
			g.AddImport(imp)
		})
		tgtCol := relModel.Columns.Get(rel.Tgt[0])
		lo.ForEach(helper.ImportsForTypes("go", "", srcCol.Type), func(imp *model.Import, _ int) {
			g.AddImport(imp)
		})
		gt, err := model.ToGoType(srcCol.Type, srcCol.Nullable, m.Package, enums)
		if err != nil {
			return nil, err
		}
		g.AddImport(helper.ImpLo)
		ret.WF("\t\t%sIDsBy%s := lo.Map(ret, func(x %s, _ int) %s {", relModel.Camel(), srcCol.Proper(), m.Pointer(), gt)
		ret.WF("\t\t\treturn x.%s", srcCol.Proper())
		ret.W("\t\t})")
		var sfx string
		if relModel.IsSoftDelete() {
			sfx = ", false"
		}
		c := fmt.Sprintf("%sIDsBy%s", relModel.Camel(), srcCol.Proper())
		if srcCol.Nullable && !srcCol.Type.Scalar() && !tgtCol.Nullable {
			g.AddImport(helper.ImpAppUtil)
			c = "util.ArrayDereference(" + c + ")"
		}
		call := "\t\t%sBy%s, err := as.Services.%s.GetMultiple(ps.Context, nil, nil%s, ps.Logger, %s...)"
		ret.WF(call, relModel.CamelPlural(), srcCol.Proper(), relModel.Proper(), sfx, c)
		ret.WE(2, `""`)

		toStrings += fmt.Sprintf(", %sBy%s: %sBy%s", relModel.ProperPlural(), srcCol.Proper(), relModel.CamelPlural(), srcCol.Proper())
	}
	var searchSuffix string
	if m.HasSearches() {
		searchSuffix += ", SearchQuery: q"
	}
	if m.HasTag("count") {
		searchSuffix += ", Count: count"
	}
	ret.WF("\t\tpage := &v%s.List{Models: ret%s, Params: ps.Params%s}", m.Package, toStrings, searchSuffix)
	render := "\t\treturn %sRender(r, as, page, ps, %s%s)"
	var bc string
	if grp != nil {
		bc = grp.BC()
	}
	ret.WF(render, prefix, m.Breadcrumbs(), bc)
	ret.W("\t})")
	ret.W("}")
	return ret, nil
}
