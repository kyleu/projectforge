package view

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func history(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "History.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	if m.Columns.HasFormat(model.FmtCountry) {
		g.AddImport(helper.ImpAppUtil)
	}
	g.AddBlocks(exportViewHistoryClass(m), exportViewHistoryBody(m), exportViewHistoryTable(m))
	return g.Render(addHeader)
}

func exportViewHistoryClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("History", "struct")
	ret.W("{%% code type History struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	ret.W("  History *%s.History", m.Package)
	ret.W("} %%}")
	return ret
}

func exportViewHistoryBody(m *model.Model) *golang.Block {
	ret := golang.NewBlock("HistoryBody", "func")
	ret.W("{%% func (p *History) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    <h3>{%%= components.SVGRefIcon(`" + m.Icon + "`, ps) %%} " + m.Title() + " History [{%%s p.History.ID.String() %%}]</h3>")
	ret.W("    <table class=\"mt\">")
	ret.W("      <tbody>")

	ret.W("        <tr>")
	ret.W("          <th class=\"shrink\">ID</th>")
	ret.W("          <td>{%%s p.History.ID.String() %%}</td>")
	ret.W("        </tr>")
	for _, pk := range m.PKs() {
		x := pk.Clone()
		x.Name = m.Proper() + x.Proper()
		ret.W("        <tr>")
		ret.W("          <th class=\"shrink\">%s</th>", x.Title())
		ret.W("          <td><a href=\"{%%%%s p.Model.WebPath() %%%%}\">%s</a></td>", x.ToGoViewString("p.History.", true, false))
		ret.W("        </tr>")
	}
	ret.W("        <tr>")
	ret.W("          <th class=\"shrink\">Old</th>")
	ret.W("          <td>{%%= components.JSON(p.History.Old) %%}</td>")
	ret.W("        </tr>")
	ret.W("        <tr>")
	ret.W("          <th class=\"shrink\">New</th>")
	ret.W("          <td>{%%= components.JSON(p.History.New) %%}</td>")
	ret.W("        </tr>")
	ret.W("        <tr>")
	ret.W("          <th class=\"shrink\">Changes</th>")
	ret.W("          <td>{%%= components.DisplayDiffs(p.History.Changes) %%}</td>")
	ret.W("        </tr>")
	ret.W("        <tr>")
	ret.W("          <th class=\"shrink\">Created</th>")
	ret.W("          <td>{%%= components.DisplayTimestamp(&p.History.Created) %%}</td>")
	ret.W("        </tr>")
	ret.W("      </tbody>")
	ret.W("    </table>")
	ret.W("  </div>")
	ret.W("{%% endfunc %%}")
	return ret
}

func exportViewHistoryTable(m *model.Model) *golang.Block {
	ret := golang.NewBlock("HistoryTable", "struct")
	const decl = "{%%%% func HistoryTable(model *%s.%s, histories %s.Histories, params filter.ParamSet, as *app.State, ps *cutil.PageState) %%%%}"
	ret.W(decl, m.Package, m.Proper(), m.Package)
	ret.W("  {%%- code prms := params.Get(\"history_history\", nil, ps.Logger).Sanitize(\"history_history\") -%%}")
	ret.W("  <table class=\"mt\">")
	ret.W("    <thead>")
	ret.W("      <tr>")
	addHeader := func(key string, title string, help string) {
		const msg = "        {%%%%= components.TableHeaderSimple(\"%s_history\", \"%s\", \"%s\", \"%s\", prms, ps.URI, ps) %%%%}"
		ret.W(msg, m.Name, key, title, help)
	}
	addHeader("id", "ID", "System-generated history UUID identifier")
	for _, pk := range m.PKs() {
		addHeader(m.Package+"_"+pk.Name, m.Title()+" "+pk.Title(), model.Help(pk.Type, pk.Format))
	}
	addHeader("c", "Changes", "Object changes")
	addHeader("created", "Created", "Time when history was created")
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, h := range histories -%%}")
	ret.W("      <tr>")
	ret.W("        <td><a href=\"{%%s model.WebPath() %%}/history/{%%s h.ID.String() %%}\">{%%s h.ID.String() %%}</a></td>")
	for _, pk := range m.PKs() {
		x := pk.Clone()
		x.Name = m.Proper() + x.Proper()
		ret.W("        <td><a href=\"{%%s model.WebPath() %%}\">" + x.ToGoViewString("h.", true, false) + "</a></td>")
	}
	ret.W("        <td>{%%= components.DisplayDiffs(h.Changes) %%}</td>")
	ret.W("        <td>{%%= components.DisplayTimestamp(&h.Created) %%}</td>")
	ret.W("      </tr>")
	ret.W("      {%%- endfor -%%}")
	ret.W("    </tbody>")
	ret.W("  </table>")
	ret.W("{%% endfunc %%}")
	return ret
}
