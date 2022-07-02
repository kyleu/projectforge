package view

import (
	"fmt"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func exportViewDetailRevisions(ret *golang.Block, m *model.Model) {
	hc := m.HistoryColumns(false)
	ret.W("  {%%%%- if len(p.%s) > 1 -%%%%}", hc.Col.ProperPlural())
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>%s</h3>", hc.Col.ProperPlural())
	ret.W("    {%%%%- code prms := p.Params.Get(%q, nil, ps.Logger).Sanitize(%q) -%%%%}", m.Package, m.Package)
	ret.W("    <table class=\"mt\">")
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
