package view

import (
	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func history(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "History.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpLayout, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.Package))
	g.AddBlocks(exportViewHistoryClass(m), exportViewHistoryBody(m), exportViewHistoryTable(m))
	return g.Render(addHeader)
}

func exportViewHistoryClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("History", "struct")
	ret.W("{%% code type History struct {")
	ret.W("  layout.Basic")
	ret.W("  Model *%s.%s", m.Package, m.Proper())
	ret.W("  History *%s.%sHistory", m.Package, m.Proper())
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
		ret.W("        <tr>")
		ret.W("          <th class=\"shrink\">%s %s</th>", m.Title(), pk.Title())
		ret.W("          <td><a href=\"{%%%%s p.Model.WebPath() %%%%}\">{%%%%s p.History.%s%s %%%%}</a></td>", m.Proper(), pk.Proper())
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
	ret.W("{%% func HistoryTable(model *history.History, histories history.HistoryHistories, params filter.ParamSet, as *app.State, ps *cutil.PageState) %%}")
	ret.W("  {%%- code prms := params.Get(\"history_history\", nil, ps.Logger) -%%}")
	ret.W("  <table class=\"mt\">")
	ret.W("    <thead>")
	ret.W("      <tr>")
	addHeader := func(key string, title string, help string) {
		msg := "        {%%%%= components.TableHeaderSimple(\"%s_history\", \"%s\", \"%s\", \"%s\", prms, ps.URI, ps) %%%%}"
		ret.W(msg, m.Name, key, title, help)
	}
	addHeader("id", "ID", model.TypeUUID.Help())
	for _, pk := range m.PKs() {
		addHeader(m.Package+"_"+pk.Name, m.Title()+" "+pk.Title(), pk.Type.Help())
	}
	addHeader("c", "Changes", "Object changes")
	addHeader("created", "Created", model.TypeTimestamp.Help())
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, h := range histories -%%}")
	ret.W("      <tr>")
	ret.W("        <td><a href=\"{%%s model.WebPath() %%}/history/{%%s h.ID.String() %%}\">{%%s h.ID.String() %%}</a></td>")
	ret.W("        <td><a href=\"{%%s model.WebPath() %%}\">{%%s h.HistoryID %%}</a></td>")
	ret.W("        <td>{%%= components.DisplayDiffs(h.Changes) %%}</td>")
	ret.W("        <td>{%%= components.DisplayTimestamp(&h.Created) %%}</td>")
	ret.W("      </tr>")
	ret.W("      {%%- endfor -%%}")
	ret.W("    </tbody>")
	ret.W("  </table>")
	ret.W("{%% endfunc %%}")
	return ret
}
