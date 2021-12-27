package export

import (
	"github.com/kyleu/projectforge/app/file"
)

func exportViewDetail(m *Model, args *Args) *file.File {
	g := NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "Detail.html")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/views/layout")
	g.AddBlocks(exportViewDetailClass(m), exportViewDetailBody(m))
	return g.Render()
}

func exportViewDetailClass(m *Model) *Block {
	ret := NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.WF("  Model *%s.%s", m.Package, m.proper())
	ret.W("} %%}")
	return ret
}

func exportViewDetailBody(m *Model) *Block {
	ret := NewBlock("DetailBody", "func")
	ret.W("{%% func (p *Detail) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.proper() + " [{%%s p.Model.String() %%}]</h3>")
	ret.W("    {%%= components.JSON(p.Model) %%}")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
