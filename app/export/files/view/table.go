package view

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func table(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.Package))
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
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name), col.Help())
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
