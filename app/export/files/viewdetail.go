package files

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func ViewDetail(m *model.Model, args *model.Args) *file.File {
	g := golang.NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "Detail.html")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/layout")
	g.AddBlocks(exportViewDetailClass(m), exportViewDetailBody(m))
	return g.Render()
}

func exportViewDetailClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	ret.W("} %%}")
	return ret
}

func exportViewDetailBody(m *model.Model) *golang.Block {
	ret := golang.NewBlock("DetailBody", "func")
	ret.W("{%% func (p *Detail) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\"><a href=\"{%%s p.Model.WebPath() %%}/edit\"><button>Edit</button></a></div>")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.Proper() + " [{%%s p.Model.String() %%}]</h3>")
	ret.W("    <table>")
	ret.W("      <tbody>")
	for _, col := range m.Columns {
		ret.W("        <tr>")
		ret.W("          <th class=\"shrink\">" + col.Proper() + "</th>")
		ret.W("          <td>" + col.ToGoViewString("p.Model.") + "</td>")
		ret.W("        </tr>")
	}
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
