package view

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func exportViewDetailRelations(g *golang.Template, m *model.Model, rels model.Relations, models model.Models) (*golang.Block, error) {
	ret := golang.NewBlock("DetailRelations", "func")
	ret.W("{%% func DetailRelations(as *app.State, p *Detail, ps *cutil.PageState) %%}")
	exportViewDetailReverseRelations(ret, m, rels, models, g)
	ret.W(helper.TextEndFunc)
	return ret, nil
}

func exportViewDetailReverseRelations(ret *golang.Block, m *model.Model, rels model.Relations, models model.Models, g *golang.Template) {
	g.AddImport(helper.ImpAppUtil)
	ret.WF("  {%%%%- code relationHelper := %s.%s{p.Model} -%%%%}", m.Package, m.ProperPlural())
	ret.W("  <div class=\"card\">")
	ret.W("    <h3 class=\"mb\">Relations</h3>")
	ret.W("    <ul class=\"accordion\">")
	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		exportViewDetailReverseRelation(ret, m, rel, models, g)
	})
	ret.W("    </ul>")
	ret.W("  </div>")
}

func exportViewDetailReverseRelation(ret *golang.Block, m *model.Model, rel *model.Relation, models model.Models, g *golang.Template) {
	g.AddImport(helper.ImpAppUtil)
	tgt := models.Get(rel.Table)
	tgtCols := rel.TgtColumns(tgt)
	tgtName := fmt.Sprintf("%sBy%s", tgt.ProperPlural(), util.StringJoin(tgtCols.ProperNames(), ""))
	ret.W("      <li>")
	extra := fmt.Sprintf("{%%%% if p.Params.Specifies(`%s`) %%%%} checked=\"checked\""+helper.TextEndIfExtra, tgt.Package)
	ret.WF("        <input id=\"accordion-%s\" type=\"checkbox\" hidden=\"hidden\"%s />", tgtName, extra)
	ret.WF("        <label for=\"accordion-%s\">", tgtName)
	ret.W("          {%%= components.ExpandCollapse(3, ps) %%}")
	ret.WF("          {%%%%= components.SVGInline(`%s`, 16, ps) %%%%}", tgt.Icon)
	msg := "          {%%%%s util.StringPlural(len(p.Rel%s), \"%s\") %%%%} by [%s]"
	ret.WF(msg, tgtName, tgt.Title(), util.StringJoin(tgtCols.Titles(), ", "))
	ret.W("        </label>")
	ret.W("        <div class=\"bd\"><div><div>")
	ret.WF("          {%%%%- if len(p.Rel%s) == 0 -%%%%}", tgtName)
	ret.WF("          <em>no related %s</em>", tgt.TitlePlural())
	ret.W("          {%%- else -%%}")
	ret.W("          <div class=\"overflow clear\">")
	var addons string
	lo.ForEach(tgt.Relations, func(r *model.Relation, _ int) {
		if len(r.Tgt) == 1 {
			if r.Table == m.Name {
				addons += ", relationHelper"
			} else {
				addons += ", nil"
			}
		}
	})
	if m.PackageWithGroup("") == tgt.PackageWithGroup("") {
		ret.WF("            {%%%%= Table(p.Rel%s%s, p.Params, as, ps) %%%%}", tgtName, addons)
	} else {
		ret.WF("            {%%%%= v%s.Table(p.Rel%s%s, p.Params, as, ps) %%%%}", tgt.Package, tgtName, addons)
	}
	ret.W("          </div>")
	ret.W(ind5 + helper.TextEndIfDash)
	ret.W("        </div></div></div>")
	ret.W("      </li>")
}
