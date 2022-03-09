package controller

import (
	"fmt"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

func Grouping(m *model.Model, args *model.Args, grp *model.Column, addHeader bool) (*file.File, error) {
	name := m.Package + "by" + grp.Name
	g := golang.NewFile("controller", []string{"app", "controller"}, name)
	g.AddImport(helper.ImpFmt, helper.ImpErrors, helper.ImpFastHTTP, helper.ImpApp, helper.ImpCutil)
	g.AddImport(helper.AppImport("app/" + m.Package))
	g.AddImport(helper.AppImport("views/v" + m.Package))
	g.AddBlocks(
		controllerGrouped(m, grp), controllerList(m, grp), controllerDetail(args.Models, m, grp),
		controllerCreateForm(m, grp), controllerCreate(m, g, grp),
		controllerEditForm(m, grp), controllerEdit(m, g, grp), controllerDelete(m, g, grp),
	)
	return g.Render(addHeader)
}

func controllerGrouped(m *model.Model, grp *model.Column) *golang.Block {
	name := fmt.Sprintf("%s%sList", m.Proper(), grp.Proper())
	ret := golang.NewBlock(name, "func")
	ret.W("func %s(rc *fasthttp.RequestCtx) {", name)
	ret.W("\tact(\"%s.%s.list\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package, grp.Camel())
	ret.W("\t\tps.Title = \"[%s] by %s\"", m.ProperPlural(), grp.TitleLower())
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", " + incDel
	}
	ret.W("\t\tret, err := as.Services.%s.Get%s(ps.Context, nil%s)", m.Proper(), grp.ProperPlural(), suffix)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.%s{%s: ret}, ps, %q, %q)", m.Package, grp.ProperPlural(), grp.ProperPlural(), m.Package, grp.Camel())
	ret.W("\t})")
	ret.W("}")
	return ret
}
