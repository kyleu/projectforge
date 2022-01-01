package files

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func ViewTable(m *model.Model, args *model.Args) *file.File {
	g := golang.NewGoTemplate(m.Package, []string{"views", "v" + m.Package}, "Table.html")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/lib/filter")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/components")
	g.AddBlocks(exportViewTableFunc(m))
	return g.Render()
}

func exportViewTableFunc(m *model.Model) *golang.Block {
	linkURL := m.LinkURL("model.")
	ret := golang.NewBlock("Table", "func")
	ret.W("{%% func Table(models " + m.Package + "." + m.ProperPlural() + ", params filter.ParamSet, as *app.State, ps *cutil.PageState) %%}")
	ret.W("  {%%- code prms := params.Get(\"" + m.Package + "\", nil, ps.Logger) -%%}")
	ret.W("  <table>")
	ret.W("    <thead>")
	ret.W("      <tr>")
	for _, col := range m.Columns {
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name))
		ret.W("        {%%= " + call + " %%}")
	}
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, model := range models -%%}")
	ret.W("      <tr>")
	for _, col := range m.Columns {
		if col.PK {
			ret.W("        <td><a href=\"" + linkURL + "\">" + col.ToGoViewString("model.") + "</a></td>")
		} else {
			ret.W("        <td>" + col.ToGoViewString("model.") + "</td>")
		}
	}
	ret.W("      </tr>")
	ret.W("      {%%- endfor -%%}")
	ret.W("    </tbody>")
	ret.W("  </table>")
	ret.W("{%% endfunc %%}")
	return ret
}
