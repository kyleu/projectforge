package files

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func ViewList(m *model.Model, args *model.Args) *file.File {
	g := golang.NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "List.html")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/lib/filter")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/layout")
	g.AddBlocks(exportViewListClass(m), exportViewListBody(m))
	return g.Render()
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
	ret.W("    <div class=\"right\"><a href=\"/" + m.Package + "/new\"><button>New</button></a></div>")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.ProperPlural() + "</h3>")
	ret.W("    <div class=\"mt\">")
	ret.W("      {%%= Table(p.Models, p.Params, as, ps) %%}")
	ret.W("    </div>")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
