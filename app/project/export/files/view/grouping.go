package view

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	golang2 "projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

func Grouping(m *model2.Model, grp *model2.Column, addHeader bool) (*file.File, error) {
	g := golang2.NewGoTemplate([]string{"views", "v" + m.Package}, fmt.Sprintf("%s.html", grp.ProperPlural()))
	g.AddImport(helper.ImpApp, helper.ImpAppUtil, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	g.AddBlocks(exportViewGroupedClass(grp), exportViewGroupedBody(m, grp))
	return g.Render(addHeader)
}

func exportViewGroupedClass(grp *model2.Column) *golang2.Block {
	ret := golang2.NewBlock(grp.ProperPlural(), "struct")
	ret.W("{%%%% code type %s struct {", grp.ProperPlural())
	ret.W("  layout.Basic")
	ret.W("  %s []*util.KeyValInt", grp.ProperPlural())
	ret.W("} %%}")
	return ret
}

func exportViewGroupedBody(m *model2.Model, grp *model2.Column) *golang2.Block {
	ret := golang2.NewBlock(fmt.Sprintf("%sBody", grp.ProperPlural()), "func")
	ret.W("{%%%% func (p *%s) Body(as *app.State, ps *cutil.PageState) %%%%}", grp.ProperPlural())
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>{%%%%= components.SVGRefIcon(`%s`, ps) %%%%} %s %s</h3>", m.Icon, m.Title(), grp.ProperPlural())
	ret.W("    <div class=\"mt\">")
	ret.W("      {%%- if len(p." + grp.ProperPlural() + ") == 0 -%%}")
	ret.W("      <em>No %s available</em>", grp.Plural())
	ret.W("      {%%- else -%%}")
	ret.W("      <table>")
	ret.W("        <thead>")
	ret.W("          <tr>")
	ret.W("            <th>%s</th>", grp.Proper())
	ret.W("            <th>%s Count</th>", m.Title())
	ret.W("          </tr>")
	ret.W("        </thead>")
	ret.W("        <tbody>")
	ret.W("          {%%- for _, x := range p." + grp.ProperPlural() + " -%%}")
	ret.W("          <tr>")
	ret.W("            <td><a href=\"/" + m.Camel() + "/" + grp.Camel() + "/{%%s x.Key %%}\">{%%s x.Key %%}</a></td>")
	ret.W("            <td>{%%d x.Count %%}</td>")
	ret.W("          </tr>")
	ret.W("          {%%- endfor -%%}")
	ret.W("        </tbody>")
	ret.W("      </table>")
	ret.W("      {%%- endif -%%}")
	ret.W("    </div>")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}
