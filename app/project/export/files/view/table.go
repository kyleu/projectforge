package view

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
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
	imps, err := helper.EnumImports(m.Columns.WithoutDisplays(util.KeyDetail).Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
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
	ret.W("  <table>")
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
