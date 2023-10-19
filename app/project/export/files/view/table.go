package view

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func table(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	if m.Columns.HasFormat(model.FmtCountry.Key) || m.Columns.HasFormat(model.FmtSI.Key) {
		g.AddImport(helper.ImpAppUtil)
	}
	lo.ForEach(m.Columns, func(c *model.Column, _ int) {
		if c.Type.Key() == types.KeyEnum {
			e, _ := model.AsEnumInstance(c.Type, args.Enums)
			g.AddImport(helper.AppImport("app/" + e.PackageWithGroup("")))
		}
	})
	vtf, err := exportViewTableFunc(m, args.Models, args.Enums, g)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(vtf)
	return g.Render(addHeader, linebreak)
}

func exportViewTableFunc(m *model.Model, models model.Models, enums enum.Enums, g *golang.Template) (*golang.Block, error) {
	xCols := m.Columns.ForDisplay("summary")
	firstCols := xCols.WithTag("list-first")
	restCols := xCols.WithoutTags("list-first")
	summCols := append(slices.Clone(firstCols), restCols...)
	ret := golang.NewBlock("Table", "func")
	suffix := ""
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			relCols := rel.SrcColumns(m)
			relNames := strings.Join(relCols.ProperNames(), "")
			g.AddImport(helper.AppImport("app/" + relModel.PackageWithGroup("")))
			suffix += fmt.Sprintf(", %sBy%s %s.%s", relModel.CamelPlural(), relNames, relModel.Package, relModel.ProperPlural())
		}
	})
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
		title := util.StringToTitle(col.Name)
		if col.HasTag("no-title") {
			title = ""
		}
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %s, prms, ps.URI, ps)", m.Package, col.Name, title, h)
		ret.W("        {%%= " + call + " %%}")
	}
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, model := range models -%%}")
	ret.W("      <tr>")
	lo.ForEach(summCols, func(col *model.Column, _ int) {
		viewTableColumn(g, ret, models, m, true, col, "model.", "", 4, enums)
	})
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
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool,
	col *model.Column, modelKey string, prefix string, indent int, enums enum.Enums,
) {
	ind := util.StringRepeat("  ", indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		switch {
		case col.PK && link:
			ret.W(ind+"<td><a href=%q>%s</a></td>", m.LinkURL(modelKey, enums), col.ToGoViewString(modelKey, true, false, enums, "table"))
		case col.HasTag("link") && link:
			ret.W(ind+"<td><a href=%q>%s</a></td>", m.LinkURL(modelKey, enums), col.ToGoViewString(modelKey, true, false, enums, "table"))
		case col.HasTag("grouped"):
			u := fmt.Sprintf("/%s/%s/%s", m.Route(), col.TitleLower(), col.ToGoViewString(modelKey, false, true, enums, "simple"))
			ret.W(ind+"<td><a href=%q>%s</a></td>", u, col.ToGoViewString(modelKey, true, false, enums, "table"))
		case col.HasTag("title"):
			ret.W(ind + "<td><strong>" + col.ToGoViewString(modelKey, true, false, enums, "table") + "</strong></td>")
		default:
			ret.W(ind + "<td>" + col.ToGoViewString(modelKey, true, false, enums, "table") + "</td>")
		}
		return
	}

	if types.IsString(col.Type) {
		g.AddImport(helper.ImpURL)
	}

	toStrings := getTableColumnStrings(m, modelKey, rels, models, prefix)

	ret.W(ind + "<td class=\"nowrap\">")
	if col.PK && link {
		ret.W(ind + "  <a href=\"" + m.LinkURL(modelKey, enums) + "\">" + col.ToGoViewString(modelKey, true, false, enums, "table") + toStrings + "</a>")
	} else {
		ret.W(ind + "  " + col.ToGoViewString(modelKey, true, false, enums, "table") + toStrings)
	}
	const l = "<a title=%q href=\"{%%%%s %s %%%%}\">{%%%%= components.SVGRef(%q, 18, 18, \"\", ps) %%%%}</a>"
	const msgNotNull = "%s  " + l
	const msg = "%s  {%%%% if %s%s != nil %%%%}" + l + "{%%%% endif %%%%}"
	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
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
	})
	ret.W(ind + "</td>")
}

func getTableColumnStrings(m *model.Model, modelKey string, rels model.Relations, models model.Models, prefix string) string {
	var toStrings string
	lo.ForEach(rels, func(rel *model.Relation, idx int) {
		relModel := models.Get(rel.Table)
		if !relModel.CanTraverseRelation() {
			return
		}
		srcCol := m.Columns.Get(rel.Src[0])
		tgtCol := relModel.Columns.Get(rel.Tgt[0])
		k := relModel.CamelPlural()
		if prefix != "" {
			k = prefix + relModel.ProperPlural()
		}
		relTitles := relModel.Columns.WithTag("title")
		if len(relTitles) == 0 {
			relTitles = relModel.PKs()
		}
		if len(relTitles) == 1 && relTitles[0].Name == tgtCol.Name {
			return
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
	})
	return toStrings
}
