package view

import (
	"fmt"
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func table(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	if lo.ContainsBy(m.Columns, func(x *model.Column) bool {
		return !x.Type.Scalar() || x.Type.Key() == types.KeyString
	}) {
		g.AddImport(helper.ImpComponentsView)
	}
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	if m.Columns.HasFormat(model.FmtCountry.Key) || m.Columns.HasFormat(model.FmtSI.Key) {
		g.AddImport(helper.ImpAppUtil)
	}
	imps, err := helper.EnumImports(m.Columns.WithoutDisplays(util.KeyDetail).Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("viewtable")...)
	vtf, err := exportViewTableFunc(m, args.Acronyms, args.Models, args.Enums, g)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(vtf)
	return g.Render(linebreak)
}

func exportViewTableFunc(m *model.Model, acronyms []string, models model.Models, enums enum.Enums, g *golang.Template) (*golang.Block, error) {
	xCols := m.Columns.ForDisplay("summary")
	firstCols := xCols.WithTag("list-first")
	restCols := xCols.WithoutTags("list-first")
	summCols := append(slices.Clone(firstCols), restCols...)
	ret := golang.NewBlock("Table", "func")
	var suffix string
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			relCols := rel.SrcColumns(m)
			relNames := util.StringJoin(relCols.ProperNames(), "")
			g.AddImport(helper.AppImport(relModel.PackageWithGroup("")))
			suffix += fmt.Sprintf(", %sBy%s %s.%s", relModel.CamelPlural(), relNames, relModel.Package, relModel.ProperPlural())
		}
	})
	mt := m.Package + "." + m.ProperPlural() + suffix
	ret.W("{%% func Table(models " + mt + ", params filter.ParamSet, as *app.State, ps *cutil.PageState, paths ...string) %%}")
	ret.W("  {%%- code prms := params.Sanitized(\"" + m.Package + "\", ps.Logger) -%%}")
	ret.W(`  <div class="overflow clear">`)
	ret.W("    <table>")
	ret.W("      <thead>")
	ret.W("        <tr>")
	for _, col := range summCols {
		h, err := col.HelpString(enums)
		if err != nil {
			return nil, err
		}
		title := util.StringToTitle(col.Name, acronyms...)
		if col.HasTag("no-title") {
			title = ""
		}
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %s, prms, ps.URI, ps)", m.Package, col.Name, title, h)
		ret.W("          {%%= " + call + helper.TextTmplEnd)
	}
	ret.W("        </tr>")
	ret.W("      </thead>")
	ret.W("      <tbody>")
	ret.W("        {%%- for _, model := range models -%%}")
	ret.W("        <tr>")
	lo.ForEach(summCols, func(col *model.Column, _ int) {
		viewTableColumn(ret, models, m, true, col, helper.TextModelPrefix, "", 5, enums)
	})
	ret.W("        </tr>")
	ret.W("        {%%- endfor -%%}")
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	ret.W("  {%%- if prms.HasNextPage(len(models) + prms.Offset) || prms.HasPreviousPage() -%%}")
	ret.W("  <hr />")
	ret.W("  {%%= components.Pagination(len(models) + prms.Offset, prms, ps.URI) %%}")
	ret.W(`  <div class="clear"></div>`)
	ret.W(ind1 + helper.TextEndIfDash)
	ret.W(helper.TextEndFunc)
	return ret, nil
}
