package view

import (
	"projectforge.dev/app/export/files/helper"
	"projectforge.dev/app/export/golang"
	"projectforge.dev/app/export/model"
	"projectforge.dev/app/file"
)

func list(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "List.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.Package))
	g.AddBlocks(exportViewListClass(m), exportViewListBody(m))
	return g.Render(addHeader)
}

func exportViewListClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("List", "struct")
	ret.W("{%% code type List struct {")
	ret.W("  layout.Basic")
	ret.W("  Models %s.%s", m.Package, m.ProperPlural())
	ret.W("  Params filter.ParamSet")
	ret.W("} %%}")
	return ret
}

func exportViewListBody(m *model.Model) *golang.Block {
	ret := golang.NewBlock("ListBody", "func")
	ret.W("{%% func (p *List) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\"><a href=\"/%s/new\"><button>New</button></a></div>", m.Route())
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.TitlePlural() + "</h3>")
	ret.W("    {%%- if len(p.Models) == 0 -%%}")
	ret.W("    <div class=\"mt\"><em>No %s available</em></div>", m.TitlePluralLower())
	ret.W("    {%%- else -%%}")
	ret.W("    {%%= Table(p.Models, p.Params, as, ps) %%}")
	ret.W("    {%%- endif -%%}")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
