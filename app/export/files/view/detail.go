package view

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

const commonLine = "  %s %s.%s"

func detail(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "Detail.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.Package))
	rrs := args.Models.ReverseRelations(m.Name)
	if len(rrs) > 0 {
		g.AddImport(helper.ImpFilter)
	}
	for _, rel := range rrs {
		rm := args.Models.Get(rel.Table)
		g.AddImport(helper.AppImport("views/v"+rm.Package), helper.AppImport("app/"+rm.Package))
	}
	if m.IsRevision() || m.IsHistory() {
		g.AddImport(helper.ImpFilter)
	}
	g.AddBlocks(exportViewDetailClass(m, args.Models, g), exportViewDetailBody(m, args.Models))
	return g.Render(addHeader)
}

func exportViewDetailClass(m *model.Model, models model.Models, g *golang.Template) *golang.Block {
	ret := golang.NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	rrs := models.ReverseRelations(m.Name)
	if m.IsHistory() {
		ret.W("  Histories %s.%sHistories", m.Package, m.Proper())
	}
	for _, rel := range m.Relations {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			g.AddImport(helper.AppImport("app/" + relModel.Package))
			ret.W(commonLine, relModel.ProperPlural(), relModel.Package, relModel.ProperPlural())
		}
	}

	if len(rrs) > 0 || m.IsRevision() || m.IsHistory() {
		ret.W("  Params filter.ParamSet")
	}
	for _, rel := range rrs {
		rm := models.Get(rel.Table)
		rCols := rel.TgtColumns(rm)
		ret.W("  %sBy%s %s.%s", rm.ProperPlural(), strings.Join(rCols.ProperNames(), ""), rm.Package, rm.ProperPlural())
	}
	if m.IsRevision() {
		ret.W(commonLine, m.HistoryColumn().ProperPlural(), m.Package, m.ProperPlural())
	}
	ret.W("} %%}")
	return ret
}

func exportViewDetailBody(m *model.Model, models model.Models) *golang.Block {
	ret := golang.NewBlock("DetailBody", "func")
	ret.W("{%% func (p *Detail) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\">")
	ret.W("      <a href=\"#modal-%s\"><button type=\"button\">JSON</button></a>", m.Camel())
	ret.W("      <a href=\"{%%s p.Model.WebPath() %%}/edit\"><button>Edit</button></a>")
	ret.W("    </div>")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} {%%s p.Model.TitleString() %%}</h3>")
	ret.W("    <div><a href=\"/" + m.Route() + "\"><em>" + m.Title() + "</em></a></div>")
	ret.W("    <table class=\"mt\">")
	ret.W("      <tbody>")
	for _, col := range m.Columns {
		ret.W("        <tr>")
		ret.W(`          <th class="shrink" title="%s">%s</th>`, col.Help(), col.Title())
		viewTableColumn(ret, models, m, false, col, "p.Model.", "p.", 5)
		ret.W("        </tr>")
	}
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	if m.IsRevision() {
		exportViewDetailRevisions(ret, m)
	}
	if m.IsHistory() {
		ret.W("  {%%- if len(p.Histories) > 0 -%%}")
		ret.W("  <div class=\"card\">")
		ret.W("    <h3>Histories</h3>")
		ret.W("    {%%= HistoryTable(p.Model, p.Histories, p.Params, as, ps) %%}")
		ret.W("  </div>")
		ret.W("  {%%- endif -%%}")
	}
	ret.W("  {%%- comment %%}$PF_SECTION_START(extra)${%% endcomment -%%}")
	ret.W("  {%%- comment %%}$PF_SECTION_END(extra)${%% endcomment -%%}")
	exportViewDetailRelations(ret, m, models)
	ret.W("  {%%%%= components.JSONModal(%q, \"%s JSON\", p.Model, 1) %%%%}", m.Camel(), m.Title())
	ret.W("{%% endfunc %%}")
	return ret
}

func exportViewDetailRelations(ret *golang.Block, m *model.Model, models model.Models) {
	for _, rel := range models.ReverseRelations(m.Name) {
		tgt := models.Get(rel.Table)
		tgtCols := rel.TgtColumns(tgt)
		tgtName := fmt.Sprintf("%sBy%s", tgt.ProperPlural(), strings.Join(tgtCols.ProperNames(), ""))
		ret.W("  {%%%%- if len(p.%s) > 0 -%%%%}", tgtName)
		ret.W("  <div class=\"card\">")
		const msg = "    <h3>{%%%%= components.SVGRefIcon(`%s`, ps) %%%%} Related %s by [%s]</h3>"
		ret.W(msg, tgt.Icon, tgt.TitlePluralLower(), strings.Join(rel.TgtColumns(tgt).TitlesLower(), ", "))
		var addons string
		if m.CanTraverseRelation() {
			for range tgt.Relations {
				addons += ", nil"
			}
		}
		ret.W("    <div class=\"overflow clear\">")
		ret.W("      {%%%%= v%s.Table(p.%s%s, p.Params, as, ps) %%%%}", tgt.Package, tgtName, addons)
		ret.W("    </div>")
		ret.W("  </div>")
		ret.W("  {%%- endif -%%}")
	}
}
