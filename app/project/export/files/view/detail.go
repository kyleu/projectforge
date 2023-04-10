package view

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const commonLine = "  %sBy%s %s.%s"

func detail(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Detail.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	rrs := args.Models.ReverseRelations(m.Name)
	if len(rrs) > 0 {
		g.AddImport(helper.ImpFilter)
	}
	for _, rel := range rrs {
		rm := args.Models.Get(rel.Table)
		g.AddImport(helper.AppImport("app/" + rm.PackageWithGroup("")))
		if rm.PackageWithGroup("") != m.PackageWithGroup("") {
			g.AddImport(helper.AppImport("views/" + rm.PackageWithGroup("v")))
		}
	}
	if len(rrs) > 0 || m.IsRevision() || m.IsHistory() || args.Audit(m) {
		g.AddImport(helper.ImpFilter)
	}
	if args.Audit(m) {
		g.AddImport(helper.AppImport("app/lib/audit"))
		g.AddImport(helper.AppImport("views/vaudit"))
	}
	vdb, err := exportViewDetailBody(g, m, args.Audit(m), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(exportViewDetailClass(m, args.Models, args.Audit(m), g), vdb)
	return g.Render(addHeader)
}

func exportViewDetailClass(m *model.Model, models model.Models, audit bool, g *golang.Template) *golang.Block {
	ret := golang.NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	rrs := models.ReverseRelations(m.Name)
	if m.IsHistory() {
		ret.W("  Histories %s.Histories", m.Package)
	}
	if m.Columns.HasFormat(model.FmtCountry) {
		g.AddImport(helper.ImpAppUtil)
	}
	for _, rel := range m.Relations {
		relModel := models.Get(rel.Table)
		relCols := rel.SrcColumns(m)
		relNames := strings.Join(relCols.ProperNames(), "")
		g.AddImport(helper.AppImport("app/" + relModel.PackageWithGroup("")))
		ret.W(commonLine, relModel.Proper(), relNames, "*"+relModel.Package, relModel.Proper())
	}

	if len(rrs) > 0 || m.IsRevision() || m.IsHistory() || audit {
		ret.W("  Params filter.ParamSet")
	}
	for _, rel := range rrs {
		rm := models.Get(rel.Table)
		rCols := rel.TgtColumns(rm)
		ret.W(commonLine, "Rel"+rm.ProperPlural(), strings.Join(rCols.ProperNames(), ""), rm.Package, rm.ProperPlural())
	}
	if m.IsRevision() {
		ret.W("  %s %s.%s", m.HistoryColumn().ProperPlural(), m.Package, m.ProperPlural())
	}
	if audit {
		ret.W("  AuditRecords audit.Records")
	}
	ret.W("} %%}")
	return ret
}

func exportViewDetailBody(g *golang.Template, m *model.Model, audit bool, models model.Models, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("DetailBody", "func")
	ret.W("{%% func (p *Detail) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\">")
	ret.W("      <a href=\"#modal-%s\"><button type=\"button\">JSON</button></a>", m.Camel())
	ret.W("      <a href=\"{%%s p.Model.WebPath() %%}/edit\"><button>{%%= components.SVGRef(\"edit\", 15, 15, \"icon\", ps) %%}Edit</button></a>")
	ret.W("    </div>")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} {%%s p.Model.TitleString() %%}</h3>")
	ret.W("    <div><a href=\"/" + m.Route() + "\"><em>" + m.Title() + "</em></a></div>")
	ret.W("    <table class=\"mt\">")
	ret.W("      <tbody>")
	for _, col := range m.Columns {
		ret.W("        <tr>")
		h, err := col.Help(enums)
		if err != nil {
			return nil, err
		}
		ret.W(`          <th class="shrink" title="%s">%s</th>`, h, col.Title())
		viewDetailColumn(g, ret, models, m, false, col, "p.Model.", 5)
		ret.W("        </tr>")
	}
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	if m.IsRevision() {
		if err := exportViewDetailRevisions(ret, m, enums); err != nil {
			return nil, err
		}
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
	exportViewDetailReverseRelations(ret, m, models, g)
	if audit {
		ret.W("  {%%- if len(p.AuditRecords) > 0 -%%}")
		ret.W("  <div class=\"card\">")
		ret.W("    <h3>Audits</h3>")
		ret.W("    {%%= vaudit.RecordTable(p.AuditRecords, p.Params, as, ps) %%}")
		ret.W("  </div>")
		ret.W("  {%%- endif -%%}")
	}
	ret.W("  {%%%%= components.JSONModal(%q, \"%s JSON\", p.Model, 1) %%%%}", m.Camel(), m.Title())
	ret.W("{%% endfunc %%}")
	return ret, nil
}

func exportViewDetailReverseRelations(ret *golang.Block, m *model.Model, models model.Models, g *golang.Template) {
	rels := models.ReverseRelations(m.Name)
	if len(rels) == 0 {
		return
	}
	g.AddImport(helper.ImpAppUtil)
	ret.W("  <div class=\"card\">")
	ret.W("    <h3 class=\"mb\">Relations</h3>")
	ret.W("    <ul class=\"accordion\">")
	for _, rel := range rels {
		tgt := models.Get(rel.Table)
		tgtCols := rel.TgtColumns(tgt)
		tgtName := fmt.Sprintf("%sBy%s", tgt.ProperPlural(), strings.Join(tgtCols.ProperNames(), ""))
		ret.W("      <li>")
		ret.W("        <input id=\"accordion-%s\" type=\"checkbox\" hidden />", tgtName)
		ret.W("        <label for=\"accordion-%s\">", tgtName)
		ret.W("          {%%= components.ExpandCollapse(3, ps) %%}")
		ret.W("          {%%%%= components.SVGRefIcon(`%s`, ps) %%%%}", tgt.Icon)
		msg := "          {%%%%d len(p.Rel%s) %%%%} {%%%%s util.StringPluralMaybe(\"%s\", len(p.Rel%s)) %%%%} by [%s]"
		ret.W(msg, tgtName, tgt.Title(), tgtName, strings.Join(tgtCols.Names(), ", "))
		ret.W("        </label>")
		ret.W("        <div class=\"bd\">")
		ret.W("          {%%%%- if len(p.Rel%s) == 0 -%%%%}", tgtName)
		ret.W("          <em>no related %s</em>", tgt.TitlePlural())
		ret.W("          {%%- else -%%}")
		ret.W("          <div class=\"overflow clear\">")
		var addons string
		if m.CanTraverseRelation() {
			for range tgt.Relations {
				addons += ", nil"
			}
		}
		if m.PackageWithGroup("") == tgt.PackageWithGroup("") {
			ret.W("            {%%%%= Table(p.Rel%s%s, p.Params, as, ps) %%%%}", tgtName, addons)
		} else {
			ret.W("            {%%%%= v%s.Table(p.Rel%s%s, p.Params, as, ps) %%%%}", tgt.Package, tgtName, addons)
		}
		ret.W("          </div>")
		ret.W("          {%%- endif -%%}")
		ret.W("        </div>")
		ret.W("      </li>")
	}
	ret.W("    </ul>")
	ret.W("  </div>")
}

func viewDetailColumn(g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, indent int) {
	ind := util.StringRepeat("  ", indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		switch {
		case col.PK && link:
			ret.W(ind+"<td><a href=%q>%s</a></td>", m.LinkURL(modelKey), col.ToGoViewString(modelKey, true, false))
		case col.HasTag("grouped"):
			u := fmt.Sprintf("/%s/%s/%s", m.Route(), col.TitleLower(), col.ToGoViewString(modelKey, false, true))
			ret.W(ind+"<td><a href=%q>%s</a></td>", u, col.ToGoViewString(modelKey, true, false))
		case col.HasTag("title"):
			ret.W(ind + "<td><strong>" + col.ToGoViewString(modelKey, true, false) + "</strong></td>")
		default:
			ret.W(ind + "<td>" + col.ToGoViewString(modelKey, true, false) + "</td>")
		}
		return
	}

	var toStrings string
	for _, rel := range rels {
		relModel := models.Get(rel.Table)
		lCols := rel.SrcColumns(m)
		lNames := strings.Join(lCols.ProperNames(), "")

		msg := "{%%%% if p.%sBy%s != nil %%%%} ({%%%%s p.%sBy%s.TitleString() %%%%}){%%%% endif %%%%}"
		toStrings += fmt.Sprintf(msg, relModel.Proper(), lNames, relModel.Proper(), lNames)
	}

	ret.W(ind + "<td class=\"nowrap\">")
	if col.PK && link {
		ret.W(ind + "  <a href=\"" + m.LinkURL(modelKey) + "\">" + col.ToGoViewString(modelKey, true, false) + toStrings + "</a>")
	} else {
		ret.W(ind + "  " + col.ToGoViewString(modelKey, true, false) + toStrings)
	}
	const l = "<a title=%q href=\"{%%%%s %s %%%%}\">{%%%%= components.SVGRef(%q, 18, 18, \"\", ps) %%%%}</a>"
	const msgNotNull = "%s  " + l
	const msg = "%s  {%%%% if %s%s != nil %%%%}" + l + "{%%%% endif %%%%}"
	for _, rel := range rels {
		if slices.Contains(rel.Src, col.Name) {
			switch col.Type.Key() {
			case types.KeyBool, types.KeyInt, types.KeyFloat:
				g.AddImport(helper.ImpFmt)
			}
			relModel := models.Get(rel.Table)
			if col.Nullable {
				ret.W(msg, ind, modelKey, col.Proper(), relModel.Title(), rel.WebPath(m, relModel, modelKey), relModel.Icon)
			} else {
				ret.W(msgNotNull, ind, relModel.Title(), rel.WebPath(m, relModel, modelKey), relModel.Icon)
			}
		}
	}
	ret.W(ind + "</td>")
}
