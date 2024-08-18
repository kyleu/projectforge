package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func list(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "List.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter, helper.ImpLayout)
	if m.HasSearches() {
		g.AddImport(helper.ImpComponentsEdit)
	}
	if m.HasTag("count") {
		g.AddImport(helper.ImpAppUtil)
	}
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	g.AddImport(m.Imports.Supporting("viewlist")...)
	g.AddBlocks(exportViewListClass(m, args.Models, g), exportViewListBody(m, args.Models))
	return g.Render(linebreak)
}

func exportViewListClass(m *model.Model, models model.Models, g *golang.Template) *golang.Block {
	ret := golang.NewBlock("List", "struct")
	ret.W("{%% code type List struct {")
	ret.W("  layout.Basic")
	ret.W("  Models %s.%s", m.Package, m.ProperPlural())
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		relModel := models.Get(rel.Table)
		relCols := rel.SrcColumns(m)
		relNames := strings.Join(relCols.ProperNames(), "")
		g.AddImport(helper.AppImport(relModel.PackageWithGroup("")))
		ret.W(commonLine, relModel.ProperPlural(), relNames, relModel.Package, relModel.ProperPlural())
	})
	ret.W("  Params filter.ParamSet")
	ret.W("  Paths []string")
	if m.HasSearches() {
		ret.W("  SearchQuery string")
	}
	if m.HasTag("count") {
		ret.W("  Count int")
	}
	ret.W("} %%}")
	return ret
}

func exportViewListBody(m *model.Model, models model.Models) *golang.Block {
	ret := golang.NewBlock("ListBody", "func")
	suffix := ""
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		if relModel := models.Get(rel.Table); relModel.CanTraverseRelation() {
			relCols := rel.SrcColumns(m)
			relNames := strings.Join(relCols.ProperNames(), "")
			suffix += fmt.Sprintf(", p.%sBy%s", relModel.ProperPlural(), relNames)
		}
	})

	ret.W("{%% func (p *List) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	const twoInd = "    "
	links := m.Links.WithTags(false, "list")
	pth := stringLead + m.Package + ".Route(p.Paths...) %%}"
	ln := fmt.Sprintf(`<a href="%s/_new"><button>{%%%%= components.SVGButton("plus", ps) %%%%} New</button>%s`, pth, helper.TextEndAnchor)

	if m.HasSearches() {
		ret.W(`    <div class="right">{%%%%= edit.SearchForm("", "q", "Search %s", p.SearchQuery, ps) %%%%}</div>`, m.TitlePlural())
	}

	ret.W(`    <div class="right mrs large-buttons">`)
	if len(links) > 0 {
		for _, link := range m.Links {
			icon := ""
			if link.Icon != "" {
				icon = fmt.Sprintf("{%%%%= components.SVGButton(%q, ps) %%%%}", link.Icon)
			}
			if link.Dangerous {
				msg := "      <a class=%q data-message=%q href=%q><button type=\"button\">%s %s</button></a>"
				ret.W(msg, "link-confirm", "Are you sure?", link.URL, icon, link.Title)
			} else {
				ret.W("      <a href=%q><button type=\"button\">%s %s</button></a>", link.URL, icon, link.Title)
			}
		}
	}
	msg := `<a href="{%%%%s %s.Route(p.Paths...) %%%%}/_random"><button>{%%%%= components.SVGButton("gift", ps) %%%%} Random</button></a>`
	ret.W(`      {%%%%- if len(p.Models) > 1 -%%%%}`+msg+`{%%%%- endif -%%%%}`, m.Package)
	ret.W(`      ` + ln)
	ret.W(`    </div>`)
	ret.W("    %s{%%%%= components.SVGIcon(`%s`, ps) %%%%} {%%%%s ps.Title %%%%}%s", helper.TextH3Start, m.Icon, helper.TextH3End)

	if m.HasTag("count") {
		ret.W("    {%%- if p.Count > 0 -%%}")
		ret.W("    <em>{%%s util.StringPlural(p.Count, \"items\") %%}</em>")
		ret.W("    {%%- endif -%%}")
	}
	if m.HasSearches() {
		ret.W("    <div class=\"clear\"></div>")
		ret.W("    {%%- if p.SearchQuery != \"\" -%%}")
		ret.W("    <hr />")
		ret.W("    <em>Search results for [{%%s p.SearchQuery %%}]</em> (<a href=\"?\">clear</a>)")
		ret.W(twoInd + helper.TextEndIfDash)
	}
	ret.W("    {%%- if len(p.Models) == 0 -%%}")
	ret.W("    <div class=\"mt\"><em>No %s available</em></div>", m.TitlePluralLower())
	ret.W("    {%%- else -%%}")
	ret.W("    <div class=\"mt\">")
	ret.W("      {%%= Table(p.Models" + suffix + ", p.Params, as, ps, p.Paths...) %%}")
	ret.W("    </div>")
	ret.W(twoInd + helper.TextEndIfDash)
	ret.W("  </div>")
	ret.W(helper.TextEndFunc)
	return ret
}
