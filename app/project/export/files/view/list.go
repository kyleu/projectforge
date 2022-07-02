package view

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	golang2 "projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

func list(m *model2.Model, args *model2.Args, addHeader bool) (*file.File, error) {
	g := golang2.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "List.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	g.AddBlocks(exportViewListClass(m, args.Models, g), exportViewListBody(m, args.Models))
	return g.Render(addHeader)
}

func exportViewListClass(m *model2.Model, models model2.Models, g *golang2.Template) *golang2.Block {
	ret := golang2.NewBlock("List", "struct")
	ret.W("{%% code type List struct {")
	ret.W("  layout.Basic")
	ret.W("  Models %s.%s", m.Package, m.ProperPlural())
	for _, rel := range m.Relations {
		relModel := models.Get(rel.Table)
		g.AddImport(helper.AppImport("app/" + relModel.PackageWithGroup("")))
		ret.W(commonLine, relModel.ProperPlural(), relModel.Package, relModel.ProperPlural())
	}
	ret.W("  Params filter.ParamSet")
	ret.W("} %%}")
	return ret
}

func exportViewListBody(m *model2.Model, models model2.Models) *golang2.Block {
	ret := golang2.NewBlock("ListBody", "func")

	suffix := ""
	for _, rel := range m.Relations {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			suffix += ", p." + relModel.ProperPlural()
		}
	}

	ret.W("{%% func (p *List) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\"><a href=\"/%s/new\"><button>New</button></a></div>", m.Route())
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.TitlePlural() + "</h3>")
	ret.W("    {%%- if len(p.Models) == 0 -%%}")
	ret.W("    <div class=\"mt\"><em>No %s available</em></div>", m.TitlePluralLower())
	ret.W("    {%%- else -%%}")
	ret.W("    <div class=\"overflow clear\">")
	ret.W("      {%%= Table(p.Models" + suffix + ", p.Params, as, ps) %%}")
	ret.W("    </div>")
	ret.W("    {%%- endif -%%}")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
