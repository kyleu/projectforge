package files

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func ViewEdit(m *model.Model, args *model.Args) *file.File {
	g := golang.NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "Edit.html")
	for _, imp := range importsForTypes("webedit", m.Columns.Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/layout")
	g.AddBlocks(exportViewEditClass(m), exportViewEditBody(m))
	return g.Render()
}

func exportViewEditClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Edit", "struct")
	ret.W("{%% code type Edit struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	ret.W("  IsNew bool")
	ret.W("} %%}")
	return ret
}

func exportViewEditBody(m *model.Model) *golang.Block {
	editURL := "/" + m.Package
	for _, pk := range m.Columns.PKs() {
		editURL += "/{%% " + pk.ToGoString("p.Model.") + " %%}"
	}

	delMsg := fmt.Sprintf("Are you sure you wish to delete %s [{%%%%s p.Model.String() %%%%}]?", m.Proper())

	ret := golang.NewBlock("EditBody", "func")
	ret.W("{%% func (p *Edit) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\"><a href=\"{%%s p.Model.WebPath() %%}/delete\" onclick=\"return confirm('" + delMsg + "')\"><button>Delete</button></a></div>")
	ret.W("    {%%- if p.IsNew -%%}")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} New " + m.Proper() + "</h3>")
	ret.W("    {%%- else -%%}")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} Edit " + m.Proper() + " [{%%s p.Model.String() %%}]</h3>")
	ret.W("    {%%- endif -%%}")
	ret.W("    <form action=\"\" method=\"post\">")
	ret.W("      <table>")
	ret.W("        <tbody>")
	for _, col := range m.Columns {
		call := col.ToGoEditString("p.Model.")
		if col.PK {
			ret.W("          {%% if p.IsNew %%}" + call + "{%% endif %%}")
		} else {
			ret.W("          " + call)
		}
	}
	ret.W("          <tr><td colspan=\"2\"><button type=\"submit\">Save Changes</button></td></tr>")
	ret.W("        </tbody>")
	ret.W("      </table>")
	ret.W("    </form>")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
