package view

import (
	"fmt"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func exportViewDetailRevisions(ret *golang.Block, m *model.Model, enums enum.Enums) error {
	hc := m.HistoryColumns(false)
	ret.W("  {%%%%- if len(p.%s) > 1 -%%%%}", hc.Col.ProperPlural())
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>%s</h3>", hc.Col.ProperPlural())
	ret.W("    {%%%%- code prms := p.Params.Get(%q, nil, ps.Logger).Sanitize(%q) -%%%%}", m.Package, m.Package)
	ret.W("    <table class=\"mt\">")
	ret.W("      <thead>")
	ret.W("        <tr>")
	addHeader := func(col *model.Column) error {
		h, err := col.Help(enums)
		if err != nil {
			return err
		}
		msg := "components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)"
		call := fmt.Sprintf(msg, m.Package, col.Name, util.StringToTitle(col.Name), h)
		ret.W("          {%%= " + call + " %%}")
		return nil
	}
	for _, pk := range m.PKs() {
		if err := addHeader(pk); err != nil {
			return err
		}
	}
	if err := addHeader(hc.Col); err != nil {
		return err
	}
	for _, c := range m.Columns.WithTag("created") {
		if err := addHeader(c); err != nil {
			return err
		}
	}
	ret.W("        </tr>")
	ret.W("      </thead>")
	ret.W("      <tbody>")
	ret.W("        {%%- for _, model := range p." + hc.Col.ProperPlural() + " -%%}")
	ret.W("        <tr>")
	linkURL := m.LinkURL("model.") + "/" + hc.Col.Camel() + "/" + hc.Col.ToGoViewString("model.", false, true)
	addView := func(col *model.Column) {
		if col.PK || col.HasTag(model.RevisionType) {
			ret.W("          <td><a href=\"" + linkURL + "\">" + col.ToGoViewString("model.", true, false) + "</a></td>")
		} else {
			ret.W("          <td>" + col.ToGoViewString("model.", true, false) + "</td>")
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
	return nil
}
