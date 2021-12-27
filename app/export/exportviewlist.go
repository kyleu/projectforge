package export

import (
	"github.com/kyleu/projectforge/app/file"
)

func exportViewList(m *Model, args *Args) *file.File {
	g := NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "List.html")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/views/layout")
	g.AddBlocks(exportViewListClass(m), exportViewListBody(m))
	return g.Render()
}

func exportViewListClass(m *Model) *Block {
	ret := NewBlock("List", "struct")
	ret.W("{%% code type List struct {")
	ret.W("  layout.Basic")
	ret.WF("  Models %s.%s", m.Package, m.properPlural())
	ret.W("} %%}")
	return ret
}

func exportViewListBody(m *Model) *Block {
	ret := NewBlock("ListBody", "func")
	ret.W("{%% func (p *List) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.properPlural() + "</h3>")
	ret.W("    {%%= components.JSON(p.Models) %%}")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
