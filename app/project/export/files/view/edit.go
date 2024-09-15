package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func edit(m *model.Model, p *project.Project, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"views", m.PackageWithGroup("v")}, "Edit.html")
	lo.ForEach(helper.ImportsForTypes("webedit", "", m.Columns.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpApp, helper.ImpComponents, helper.ImpComponentsEdit, helper.ImpCutil, helper.ImpLayout, helper.ImpAppUtil)
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))

	imps, err := helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("viewedit")...)
	veb, err := exportViewEditBody(m, p, args)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(exportViewEditClass(m), veb)
	return g.Render(linebreak)
}

func exportViewEditClass(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Edit", "struct")
	ret.W("{%% code type Edit struct {")
	ret.W("  layout.Basic")
	ret.W("  Model %s", m.Pointer())
	ret.W("  Paths []string")
	ret.W("  IsNew bool")
	ret.W("} %%}")
	return ret
}

func exportViewEditBody(m *model.Model, p *project.Project, args *model.Args) (*golang.Block, error) {
	delMsg := fmt.Sprintf("Are you sure you wish to delete %s [{%%%%s p.Model.String() %%%%}]?", m.TitleLower())
	ret := golang.NewBlock("EditBody", "func")
	ret.W("{%% func (p *Edit) Body(as *app.State, ps *cutil.PageState) %%}")
	ret.W("  <div class=\"card\">")
	ret.W("    {%%- if p.IsNew -%%}")
	ret.W("    <div class=\"right\"><a href=\"?prototype=random\"><button>Random</button></a></div>")
	ret.W("    %s%s New %s%s", helper.TextH3Start, iconRef(m), m.Title(), helper.TextH3End)
	ret.W("    {%%- else -%%}")
	delPrefix := "    <div class=\"right\"><a class=\"link-confirm\" href=\"{%%s p.Model.WebPath(p.Paths...) %%}/delete\" data-message=\""
	ret.W(delPrefix + delMsg + `"><button>{%%= components.SVGButton("times", ps) %%} Delete</button></a></div>`)
	ret.W("    %s%s Edit %s [{%%%%s p.Model.String() %%%%}]%s", helper.TextH3Start, iconRef(m), m.Title(), helper.TextH3End)
	ret.W("    " + helper.TextEndIfDash)
	rt := fmt.Sprintf("{%%%%s util.Choose(p.IsNew, %s.Route(p.Paths...) + `/_new`, p.Model.WebPath(p.Paths...) + `/edit`) %%%%}", m.Package)
	ret.W("    <form action=%q class=\"mt\" method=\"post\">", rt)
	ret.W("      <table class=\"mt expanded\">")
	ret.W("        <tbody>")
	editCols := m.Columns.NotDerived().WithoutTags("created", "updated")
	for _, col := range editCols {
		id := ""
		if len(m.RelationsFor(col)) > 0 {
			id = "input-" + col.Camel()
		}
		call, err := col.ToGoEditString("p.Model.", col.Format, id, args.Enums)
		if err != nil {
			return nil, err
		}
		if col.PK {
			ret.W("          {%% if p.IsNew %%}" + call + helper.TextEndIf)
		} else {
			ret.W(ind5 + call)
		}
	}
	ret.W("          <tr><td colspan=\"2\"><button type=\"submit\">Save Changes</button></td></tr>")
	ret.W("        </tbody>")
	ret.W("      </table>")
	ret.W("    </form>")
	ret.W("  </div>")

	var scripts []string
	for _, rel := range m.Relations {
		if len(rel.Src) != 1 {
			continue
		}
		relScript, err := exportViewEditRelation(m, rel, p, args)
		if err != nil {
			return nil, err
		}
		if relScript != "" {
			scripts = append(scripts, relScript)
		}
	}
	if len(scripts) > 0 {
		ret.W("  <script>")
		ret.W("    document.addEventListener(\"DOMContentLoaded\", function() {")
		for _, x := range scripts {
			ret.W(x)
		}
		ret.W("    });")
		ret.W("  </script>")
	}

	ret.W(helper.TextEndFunc)
	return ret, nil
}

func exportViewEditRelation(m *model.Model, rel *model.Relation, p *project.Project, args *model.Args) (string, error) {
	relModel := args.Models.Get(rel.Table)
	if !relModel.HasTag("search") {
		return "", nil
	}
	if !relModel.CanTraverseRelation() {
		return "", nil
	}
	src := m.Columns.Get(rel.Src[0])
	tgt := relModel.Columns.Get(rel.Tgt[0])

	var title string
	if titles := relModel.Columns.Searches(); len(titles) > 0 {
		var refs []string
		lo.ForEach(titles, func(title *model.Column, _ int) {
			if !title.PK {
				refs = append(refs, fmt.Sprintf("(o[%q] || %q)", title.Camel(), "[no "+title.TitleLower()+"]"))
			}
		})
		title = fmt.Sprintf(`(o) => %s + " (" + o[%q] + ")"`, strings.Join(refs, " + \" / \" + "), tgt.Camel())
	} else {
		title = fmt.Sprintf("(o) => o[%q]", tgt.Camel())
	}

	val := fmt.Sprintf("(o) => o[%q]", tgt.Camel())

	msg := `      %s.autocomplete(document.getElementById("input-%s"), "/%s?%s.l=10", "q", %s, %s);`
	return fmt.Sprintf(msg, p.CleanKey(), src.Camel(), relModel.Route(), relModel.Camel(), title, val), nil
}
