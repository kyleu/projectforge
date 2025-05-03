package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

const commonLine, ind5 = "  %sBy%s %s.%s", "          "

func iconRef(m *model.Model) string {
	if icons := m.Columns.WithFormat("icon"); len(icons) == 1 {
		return "{%%= components.SVGIcon(p.Model." + icons[0].IconDerived() + ", ps) %%}"
	}
	return "{%%= components.SVGIcon(`" + m.Icon + "`, ps) %%}"
}

func detail(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Detail.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout)
	if lo.ContainsBy(m.Columns, func(x *model.Column) bool {
		return !x.Type.Scalar() || x.Type.Key() == types.KeyString
	}) {
		g.AddImport(helper.ImpComponentsView)
	}
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	rrs := args.Models.WithController().ReverseRelations(m.Name)
	if len(rrs) > 0 {
		g.AddImport(helper.ImpFilter)
	}
	if m.Columns.HasFormat(model.FmtSI.Key) {
		g.AddImport(helper.ImpAppUtil)
	}
	lo.ForEach(rrs, func(rel *model.Relation, _ int) {
		rm := args.Models.Get(rel.Table)
		g.AddImport(helper.AppImport(rm.PackageWithGroup("")))
		if rm.PackageWithGroup("") != m.PackageWithGroup("") {
			g.AddImport(helper.ViewImport(rm.PackageWithGroup("v")))
		}
	})
	imps, err := helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("viewdetail")...)
	lo.ForEach(m.Columns, func(c *model.Column, _ int) {
		switch c.Type.Key() {
		case types.KeyEnum:
			e, _ := model.AsEnumInstance(c.Type, args.Enums)
			g.AddImport(helper.AppImport(e.PackageWithGroup("")))
		case types.KeyList:
			if t := c.Type.ListType(); t != nil && t.Key() == types.KeyEnum {
				e, _ := model.AsEnumInstance(t, args.Enums)
				g.AddImport(helper.AppImport(e.PackageWithGroup("")))
			}
		}
	})
	if len(rrs) > 0 || args.Audit(m) {
		g.AddImport(helper.ImpFilter)
	}
	if args.Audit(m) {
		g.AddImport(helper.AppImport("lib/audit"))
		g.AddImport(helper.ViewImport("vaudit"))
	}
	vdb, err := exportViewDetailBody(m, rrs, args.Audit(m))
	if err != nil {
		return nil, err
	}
	vdt, err := exportViewDetailTable(m, args.Models, args.Enums)
	if err != nil {
		return nil, err
	}

	g.AddBlocks(exportViewDetailClass(m, rrs, args.Models, args.Audit(m), g), vdb, vdt)
	if len(rrs) > 0 {
		vdr, err := exportViewDetailRelations(g, m, rrs, args.Models)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(vdr)
	}
	return g.Render(linebreak)
}

func exportViewDetailClass(m *model.Model, rrs model.Relations, models model.Models, audit bool, g *golang.Template) *golang.Block {
	ret := golang.NewBlock("Detail", "struct")
	ret.W("{%% code type Detail struct {")
	ret.W("  layout.Basic")
	ret.WF("  Model %s", m.Pointer())
	if m.Columns.HasFormat(model.FmtCountry.Key) {
		g.AddImport(helper.ImpAppUtil)
	}
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		relModel := models.Get(rel.Table)
		relCols := rel.SrcColumns(m)
		relNames := strings.Join(relCols.ProperNames(), "")
		g.AddImport(helper.AppImport(relModel.PackageWithGroup("")))
		ret.WF(commonLine, relModel.Proper(), relNames, "*"+relModel.Package, relModel.Proper())
	})

	if len(rrs) > 0 || audit {
		ret.W("  Params filter.ParamSet")
	}
	lo.ForEach(rrs, func(rel *model.Relation, _ int) {
		rm := models.Get(rel.Table)
		rCols := rel.TgtColumns(rm)
		ret.WF(commonLine, "Rel"+rm.ProperPlural(), strings.Join(rCols.ProperNames(), ""), rm.Package, rm.ProperPlural())
	})
	if audit {
		ret.W("  AuditRecords audit.Records")
	}
	ret.W("  Paths []string")
	ret.W("} %%}")
	return ret
}

func exportViewDetailBody(m *model.Model, rrs model.Relations, audit bool) (*golang.Block, error) {
	ret := golang.NewBlock("DetailBody", "func")
	ret.W("{%% func (p *Detail) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <div class=\"right\">")
	links := m.Links.WithTags(true, "detail")
	for _, link := range links {
		paths := lo.Map(m.PKs(), func(pk *model.Column, _ int) string {
			return "{%%s p.Model." + model.ToGoString(pk.Type, pk.Nullable, pk.Proper(), true) + helper.TextTmplEnd
		})
		u := link.URL
		u = strings.ReplaceAll(u, "{}", strings.Join(paths, "/"))
		u = strings.ReplaceAll(u, "[]", "{%%s p.Model.WebPath()"+helper.TextTmplEnd)
		var icon string
		if link.Icon != "" {
			icon = fmt.Sprintf("{%%%%= components.SVGButton(%q, ps) %%%%} ", link.Icon)
		}
		ret.WF("      <a href=%q><button type=\"button\">%s%s</button></a>", u, icon, link.Title)
	}
	ret.WF(`      <a href="#modal-%s"><button type="button" title="JSON">{%%%%= components.SVGButton("code", ps) %%%%}</button></a>`, m.Camel())
	ret.W(`      <a href="{%%s p.Model.WebPath(p.Paths...) %%}/edit" title="Edit"><button>{%%= components.SVGButton("edit", ps) %%}</button></a>`)
	ret.W("    </div>")
	ret.WF("    %s%s {%%%%s p.Model.TitleString() %%%%}%s", helper.TextH3Start, iconRef(m), helper.TextH3End)
	ret.WF("    <div><a href=\"{%%%%s %s.Route(p.Paths...) %%%%}\"><em>%s</em></a></div>", m.Package, m.Title())
	ret.W("    {%%= DetailTable(p, ps) %%}")
	ret.W("  </div>")
	ret.W("  {%%- comment %%}$PF_SECTION_START(extra)${%% endcomment -%%}")
	ret.W("  {%%- comment %%}$PF_SECTION_END(extra)${%% endcomment -%%}")
	if len(rrs) > 0 {
		ret.W("  {%%= DetailRelations(as, p, ps) %%}")
	}
	if audit {
		ret.W("  {%%- if len(p.AuditRecords) > 0 -%%}")
		ret.W("  <div class=\"card\">")
		ret.WF("    %sAudits%s", helper.TextH3Start, helper.TextH3End)
		ret.W("    {%%= vaudit.RecordTable(p.AuditRecords, p.Params, as, ps) %%}")
		ret.W("  </div>")
		ret.W(ind1 + helper.TextEndIfDash)
	}
	ret.WF("  {%%%%= components.JSONModal(%q, \"%s JSON\", p.Model, 1) %%%%}", m.Camel(), m.Title())
	ret.W(helper.TextEndFunc)
	return ret, nil
}

func exportViewDetailTable(m *model.Model, models model.Models, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("DetailTable", "func")
	ret.W("{%% func DetailTable(p *Detail, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"mt overflow full-width\">")
	ret.W("    <table>")
	ret.W("      <tbody>")
	for _, col := range m.Columns {
		if col.HasTag("debug-only") {
			ret.W(`        {%%- if as.Debug -%%}`)
		}
		ret.W("        <tr>")
		h, err := col.Help(enums)
		if err != nil {
			return nil, err
		}
		hlp := h
		if !strings.HasPrefix(hlp, "\"") {
			hlp = "\"{%%s " + hlp + " %%}\""
		}
		ret.WF(`          <th class="shrink" title=%s>%s</th>`, hlp, col.Title())
		viewDetailColumn(ret, models, m, false, col, "p.Model.", 5, enums)
		ret.W("        </tr>")
		if col.HasTag("debug-only") {
			ret.W(ind5 + helper.TextEndIfDash)
		}
	}
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	ret.W(helper.TextEndFunc)
	return ret, nil
}
