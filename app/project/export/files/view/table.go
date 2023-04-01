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

func table(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	if m.Columns.HasFormat(model.FmtCountry) {
		g.AddImport(helper.ImpAppUtil)
	}
	vtf, err := exportViewTableFunc(m, args.Models, args.Enums, g)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(vtf)
	return g.Render(addHeader)
}

func exportViewTableFunc(m *model.Model, models model.Models, enums enum.Enums, g *golang.Template) (*golang.Block, error) {
	summCols := m.Columns.ForDisplay("summary")
	ret := golang.NewBlock("Table", "func")
	suffix := ""
	for _, rel := range m.Relations {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			relCols := rel.SrcColumns(m)
			relNames := strings.Join(relCols.ProperNames(), "")
			g.AddImport(helper.AppImport("app/" + relModel.PackageWithGroup("")))
			msg := ", %sBy%s %s.%s"
			suffix += fmt.Sprintf(msg, relModel.CamelPlural(), relNames, relModel.Package, relModel.ProperPlural())
		}
	}
	ret.W("{%% func Table(models " + m.Package + "." + m.ProperPlural() + suffix + ", params filter.ParamSet, as *app.State, ps *cutil.PageState) %%}")
	ret.W("  {%%- code prms := params.Get(\"" + m.Package + "\", nil, ps.Logger).Sanitize(\"" + m.Package + "\") -%%}")
	ret.W("  <table class=\"mt\">")
	ret.W("    <thead>")
	ret.W("      <tr>")
	for _, col := range summCols {
		h, err := col.Help(enums)
		if err != nil {
			return nil, err
		}
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name), h)
		ret.W("        {%%= " + call + " %%}")
	}
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, model := range models -%%}")
	ret.W("      <tr>")
	for _, col := range summCols {
		viewTableColumn(g, ret, models, m, true, col, "model.", "", 4)
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
	return ret, nil
}

func viewTableColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, prefix string, indent int,
) {
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
		if !relModel.CanTraverseRelation() {
			continue
		}
		srcCol := m.Columns.Get(rel.Src[0])
		tgtCol := m.Columns.Get(rel.Tgt[0])
		k := relModel.CamelPlural()
		if prefix != "" {
			k = prefix + relModel.ProperPlural()
		}
		if srcCol.Nullable && !tgtCol.Nullable {
			get := fmt.Sprintf("%sBy%s.Get(*%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
			toStrings += "{%% if " + modelKey + srcCol.Proper() + " != nil %%}"
			toStrings += "{%% if x := " + get + "; x != nil %%} ({%%s x.TitleString() %%}){%% endif %%}"
			toStrings += "{%% endif %%}"
		} else {
			get := fmt.Sprintf("%sBy%s.Get(%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
			toStrings += "{%% if x := " + get + "; x != nil %%} ({%%s x.TitleString() %%}){%% endif %%}"
		}
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
