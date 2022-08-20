package view

import (
	"fmt"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func table(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	if m.Columns.HasFormat(model.FmtCountry) {
		g.AddImport(helper.ImpAppUtil)
	}
	g.AddBlocks(exportViewTableFunc(m, args.Models, g))
	return g.Render(addHeader)
}

func exportViewTableFunc(m *model.Model, models model.Models, g *golang.Template) *golang.Block {
	summCols := m.Columns.ForDisplay("summary")
	ret := golang.NewBlock("Table", "func")
	suffix := ""
	for _, rel := range m.Relations {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			g.AddImport(helper.AppImport("app/" + relModel.PackageWithGroup("")))
			msg := ", %s %s.%s"
			suffix += fmt.Sprintf(msg, relModel.Plural(), relModel.Package, relModel.ProperPlural())
		}
	}
	ret.W("{%% func Table(models " + m.Package + "." + m.ProperPlural() + suffix + ", params filter.ParamSet, as *app.State, ps *cutil.PageState) %%}")
	ret.W("  {%%- code prms := params.Get(\"" + m.Package + "\", nil, ps.Logger).Sanitize(\"" + m.Package + "\") -%%}")
	ret.W("  <table class=\"mt\">")
	ret.W("    <thead>")
	ret.W("      <tr>")
	for _, col := range summCols {
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name), col.Help())
		ret.W("        {%%= " + call + " %%}")
	}
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, model := range models -%%}")
	ret.W("      <tr>")
	for _, col := range summCols {
		viewTableColumn(ret, models, m, true, col, "model.", "", 4)
	}
	ret.W("      </tr>")
	ret.W("      {%%- endfor -%%}")
	ret.W("      {%%- if prms.HasNextPage(len(models) + prms.Offset) || prms.HasPreviousPage() -%%}")
	ret.W("      <tr>")
	ret.W("        <td colspan=\"%d\">{%%%%= components.Pagination(len(models) + prms.Offset, prms, ps.URI) %%%%}</td>", len(summCols))
	ret.W("      </tr>")
	ret.W("      {%%- endif -%%}")
	ret.W("    </tbody>")
	ret.W("  </table>")
	ret.W("{%% endfunc %%}")
	return ret
}

func viewTableColumn(ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, prefix string, indent int) {
	ind := util.StringRepeat("  ", indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		switch {
		case col.PK && link:
			ret.W(ind+"<td><a href=%q>%s</a></td>", m.LinkURL(modelKey), col.ToGoViewString(modelKey, true))
		case col.HasTag("grouped"):
			u := fmt.Sprintf("/%s/%s/%s", m.Route(), col.TitleLower(), col.ToGoViewString(modelKey, false))
			ret.W(ind+"<td><a href=%q>%s</a></td>", u, col.ToGoViewString(modelKey, true))
		case col.HasTag("title"):
			ret.W(ind + "<td><strong>" + col.ToGoViewString(modelKey, true) + "</strong></td>")
		default:
			ret.W(ind + "<td>" + col.ToGoViewString(modelKey, true) + "</td>")
		}
		return
	}

	var toStrings string
	for _, rel := range rels {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			k := relModel.Plural()
			if prefix != "" {
				k = prefix + relModel.ProperPlural()
			}
			get := fmt.Sprintf("%s.Get(%s%s)", k, modelKey, m.Columns.Get(rel.Src[0]).Proper())
			toStrings += "{%% if x := " + get + "; x != nil %%} ({%%s x.TitleString() %%}){%% endif %%}"
		}
	}

	ret.W(ind + "<td>")
	if col.PK && link {
		ret.W(ind + "  <div class=\"icon\"><a href=\"" + m.LinkURL(modelKey) + "\">" + col.ToGoViewString(modelKey, true) + toStrings + "</a></div>")
	} else {
		ret.W(ind + "  <div class=\"icon\">" + col.ToGoViewString(modelKey, true) + toStrings + "</div>")
	}
	const msg = "%s  <a title=%q href=\"{%%%%s %s %%%%}\">{%%%%= components.SVGRefIcon(%q, ps) %%%%}</a>"
	for _, rel := range rels {
		if slices.Contains(rel.Src, col.Name) {
			relModel := models.Get(rel.Table)
			ret.W(msg, ind, relModel.Title(), rel.WebPath(m, relModel, modelKey), relModel.Icon)
		}
	}
	ret.W(ind + "</td>")
}
