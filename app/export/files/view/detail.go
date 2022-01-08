package view

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func detail(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "Detail.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.Package))
	if m.IsRevision() {
		g.AddImport(helper.ImpFilter)
	}
	g.AddBlocks(exportViewDetailClass(m), exportViewDetailBody(m))
	return g.Render()
}

func exportViewDetailClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	if m.IsRevision() {
		ret.W("  %s %s.%s", m.HistoryColumn().ProperPlural(), m.Package, m.ProperPlural())
		ret.W("  Params filter.ParamSet")
	}
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
		ret.W(`          <th class="shrink" title="%s">%s</th>`, col.Help(), col.Title())
		ret.W("          <td>" + col.ToGoViewString("p.Model.") + "</td>")
		ret.W("        </tr>")
	}
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	if m.IsRevision() {
		exportViewDetailRevisions(ret, m)
	}
	ret.W("{%% endfunc %%}")
	return ret
}

func exportViewDetailRevisions(ret *golang.Block, m *model.Model) {
	hc := m.HistoryColumns(false)
	ret.W("  {%%%%- if len(p.%s) > 1 -%%%%}", hc.Col.ProperPlural())
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>%s</h3>", hc.Col.ProperPlural())
	ret.W("    {%%%%- code prms := p.Params.Get(%q, nil, ps.Logger) -%%%%}", m.Package)
	ret.W("    <table>")
	ret.W("      <thead>")
	ret.W("        <tr>")
	addHeader := func(col *model.Column) {
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name), col.Help())
		ret.W("          {%%= " + call + " %%}")
	}
	for _, pk := range m.PKs() {
		addHeader(pk)
	}
	addHeader(hc.Col)
	for _, c := range m.Columns.WithTag("created") {
		addHeader(c)
	}
	ret.W("        </tr>")
	ret.W("      </thead>")
	ret.W("      <tbody>")
	ret.W("        {%%- for _, model := range p." + hc.Col.ProperPlural() + " -%%}")
	ret.W("        <tr>")
	linkURL := m.LinkURL("model.") + "/" + hc.Col.Camel() + "/" + hc.Col.ToGoViewString("model.")
	addView := func(col *model.Column) {
		if col.PK || col.HasTag(model.RevisionType) {
			ret.W("          <td><a href=\"" + linkURL + "\">" + col.ToGoViewString("model.") + "</a></td>")
		} else {
			ret.W("          <td>" + col.ToGoViewString("model.") + "</td>")
		}
	}
	for _, pk := range m.PKs() {
		addView(pk)
	}
	addView(hc.Col)
	for _, c := range m.Columns.WithTag("created") {
		addView(c)
	}
	ret.W("        </tr>")
	ret.W("        {%%- endfor -%%}")
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	ret.W("  {%%- endif -%%}")
}
