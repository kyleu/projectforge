package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func list(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "List.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter, helper.ImpLayout)
	if m.HasSearches() {
		g.AddImport(helper.ImpComponentsEdit)
	}
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	g.AddBlocks(exportViewListClass(m, args.Models, g), exportViewListBody(m, args.Models))
	return g.Render(addHeader, linebreak)
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
	if m.HasSearches() {
		ret.W("  SearchQuery string")
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
	if !m.HasSearches() {
		ret.W("    <div class=\"right\"><a href=\"/%s/_new\"><button>New</button></a></div>", m.Route())
		ret.W("    %s%s{%%%%s ps.Title %%%%}%s", helper.TextH3Start, svgRef(m.Icon), helper.TextH3End)
	} else {
		ret.W(`    <div class="right">{%%%%= edit.SearchForm("", "q", "Search %s", p.SearchQuery, ps) %%%%}</div>`, m.TitlePlural())
		ret.W(`    <div class="right mrs large-buttons">`)
		ret.W(`      {%%%%- if len(p.Models) > 0 -%%%%}<a href="/%s/_random"><button>Random</button></a>{%%%%- endif -%%%%}`, m.Route())
		ret.W(`      <a href="/%s/_new"><button>New</button>`+helper.TextEndAnchor, m.Route())
		ret.W(`    </div>`)
		ret.W("    %s{%%%%= components.SVGRefIcon(`%s`, ps) %%%%}{%%%%s ps.Title %%%%}%s", helper.TextH3Start, m.Icon, helper.TextH3End)
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
	ret.W("      {%%= Table(p.Models" + suffix + ", p.Params, as, ps) %%}")
	ret.W("    </div>")
	ret.W(twoInd + helper.TextEndIfDash)
	ret.W("  </div>")
	ret.W(helper.TextEndFunc)
	return ret
}
