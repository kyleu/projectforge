package view

import (
	"fmt"

	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

func table(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", "v" + m.Package}, "Table.html")
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpCutil, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/" + m.Package))
	g.AddBlocks(exportViewTableFunc(m, args.Models))
	return g.Render(addHeader)
}

func exportViewTableFunc(m *model.Model, models model.Models) *golang.Block {
	summCols := m.Columns.ForDisplay("summary")
	ret := golang.NewBlock("Table", "func")
	ret.W("{%% func Table(models " + m.Package + "." + m.ProperPlural() + ", params filter.ParamSet, as *app.State, ps *cutil.PageState) %%}")
	ret.W("  {%%- code prms := params.Get(\"" + m.Package + "\", nil, ps.Logger) -%%}")
	ret.W("  <table class=\"mt\">")
	ret.W("    <thead>")
	ret.W("      <tr>")
	for _, col := range summCols {
		call := fmt.Sprintf("components.TableHeaderSimple(%q, %q, %q, %q, prms, ps.URI, ps)", m.Package, col.Name, util.StringToTitle(col.Name), col.Help())
		ret.W("        {%%= " + call + " %%}")
	}
	ret.W("      </tr>")
	ret.W("    </thead>")
	ret.W("    <tbody>")
	ret.W("      {%%- for _, model := range models -%%}")
	ret.W("      <tr>")
	for _, col := range summCols {
		viewTableColumn(ret, models, m, true, col, "model.", 4)
	}
	ret.W("      </tr>")
	ret.W("      {%%- endfor -%%}")
	ret.W("    </tbody>")
	ret.W("  </table>")
	ret.W("{%% endfunc %%}")
	return ret
}

func viewTableColumn(ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, prefix string, indent int) {
	ind := util.StringRepeat("  ", indent)
	if rels := m.RelationsFor(col); len(rels) == 0 {
		switch {
		case col.PK && link:
			ret.W(ind+"<td><a href=%q>%s</a></td>", m.LinkURL(prefix), col.ToGoViewString(prefix))
		case col.HasTag("grouped"):
			u := fmt.Sprintf("/%s/%s/%s", m.Route(), col.TitleLower(), col.ToGoViewString(prefix))
			ret.W(ind+"<td><a href=%q>%s</a></td>", u, col.ToGoViewString(prefix))
		default:
			ret.W(ind + "<td>" + col.ToGoViewString(prefix) + "</td>")
		}
		return
	}
	ret.W(ind + "<td>")
	if col.PK && link {
		ret.W(ind + "  <div class=\"icon\"><a href=\"" + m.LinkURL(prefix) + "\">" + col.ToGoViewString(prefix) + "</a></div>")
	} else {
		ret.W(ind + "  <div class=\"icon\">" + col.ToGoViewString(prefix) + "</div>")
	}
	const msg = "%s  <a title=%q href=\"{%%%%s %s %%%%}\">{%%%%= components.SVGRefIcon(%q, ps) %%%%}</a>"
	for _, rel := range m.Relations {
		if slices.Contains(rel.Src, col.Name) {
			relModel := models.Get(rel.Table)
			ret.W(msg, ind, relModel.Title(), rel.WebPath(m, relModel, prefix), relModel.Icon)
		}
	}
	ret.W(ind + "</td>")
}
