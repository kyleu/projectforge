package controller

import (
	"fmt"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	golang2 "projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

func Grouping(m *model2.Model, args *model2.Args, grp *model2.Column, addHeader bool) (*file.File, error) {
	name := m.Package + "by" + grp.Name
	g := golang2.NewFile("controller", []string{"app", "controller"}, name)
	g.AddImport(helper.ImpFmt, helper.ImpErrors, helper.ImpFastHTTP, helper.ImpApp, helper.ImpCutil)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	g.AddImport(helper.AppImport("views/" + m.PackageWithGroup("v")))
	var prefix string
	if len(m.Group) > 0 {
		prefix = defaultPrefix
	}
	g.AddBlocks(
		controllerGrouped(m, grp, prefix), controllerList(m, grp, args.Models, g, prefix), controllerDetail(args.Models, m, grp, prefix),
		controllerCreateForm(m, grp, prefix), controllerCreate(m, g, grp, prefix),
		controllerEditForm(m, grp, prefix), controllerEdit(m, g, grp, prefix), controllerDelete(m, g, grp, prefix),
	)
	return g.Render(addHeader)
}

func controllerGrouped(m *model2.Model, grp *model2.Column, prefix string) *golang2.Block {
	name := fmt.Sprintf("%s%sList", m.Proper(), grp.Proper())
	ret := golang2.NewBlock(name, "func")
	ret.W("func %s(rc *fasthttp.RequestCtx) {", name)
	ret.W("\tAct(\"%s.%s.list\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package, grp.Camel())
	ret.W("\t\tps.Title = \"[%s] by %s\"", m.ProperPlural(), grp.TitleLower())
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", " + incDel
	}
	ret.W("\t\tret, err := as.Services.%s.Get%s(ps.Context, nil%s, ps.Logger)", m.Proper(), grp.ProperPlural(), suffix)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn %sRender(rc, as, &v%s.%s{%s: ret}, ps, %s, %q)", prefix, m.Package, grp.ProperPlural(), grp.ProperPlural(), m.Breadcrumbs(), grp.Camel())
	ret.W("\t})")
	ret.W("}")
	return ret
}
