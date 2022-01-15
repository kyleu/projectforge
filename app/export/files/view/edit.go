package view

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func edit(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "Edit.html")
	for _, imp := range helper.ImportsForTypes("webedit", m.Columns.Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.Package))
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
	for _, pk := range m.PKs() {
		editURL += "/{%% " + pk.ToGoString("p.Model.") + " %%}"
	}

	delMsg := fmt.Sprintf("Are you sure you wish to delete %s [{%%%%s p.Model.String() %%%%}]?", m.Proper())

	ret := golang.NewBlock("EditBody", "func")
	ret.W("{%% func (p *Edit) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    {%%- if p.IsNew -%%}")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} New " + m.Proper() + "</h3>")
	ret.W("    {%%- else -%%}")
	ret.W("    <div class=\"right\"><a href=\"{%%s p.Model.WebPath() %%}/delete\" onclick=\"return confirm('" + delMsg + "')\"><button>Delete</button></a></div>")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} Edit " + m.Proper() + " [{%%s p.Model.String() %%}]</h3>")
	ret.W("    {%%- endif -%%}")
	ret.W("    <form action=\"\" class=\"mt\" method=\"post\">")
	ret.W("      <table>")
	ret.W("        <tbody>")
	editCols := m.Columns.WithoutTag("created").WithoutTag("updated")
	for _, col := range editCols {
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
