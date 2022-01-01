package files

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func Controller(m *model.Model, args *model.Args) *file.File {
	g := golang.NewFile("controller", []string{"app", "controller"}, m.Package)
	for _, imp := range importsForTypes("parse", m.Columns.Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	g.AddImport(golang.ImportTypeInternal, "fmt")
	g.AddImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	g.AddImport(golang.ImportTypeExternal, "github.com/valyala/fasthttp")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/views/v"+m.Package)
	g.AddBlocks(
		controllerList(m), controllerDetail(m),
		controllerCreateForm(m), controllerCreate(m),
		controllerEditForm(m), controllerEdit(m),
		controllerModelFromPath(m), controllerModelFromForm(m),
	)
	return g.Render()
}

func controllerList(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Camel()+"List", "func")
	ret.W("func %sList(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.list\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tps.Title = %q", util.StringToPlural(m.PackageProper()))
	ret.W("\t\tparams := cutil.ParamSetFromRequest(rc)")
	ret.W("\t\tret, err := as.Services.%s.List(ps.Context, nil, params.Get(%q, nil, ps.Logger))", m.PackageProper(), m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.List{Models: ret, Params: params}, ps, %q)", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDetail(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Camel()+"Detail", "func")
	ret.W("func %sDetail(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.detail\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = ret.String()")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Detail{Model: ret}, ps, %q, ret.String())", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerModelFromPath(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.PackageProper()+"FromPath", "func")
	ret.W("func %sFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*%s, error) {", m.Package, m.ClassRef())
	pks := m.Columns.PKs()
	for _, col := range pks {
		controllerArgFor(col, ret)
	}
	args := make([]string, 0, len(pks))
	for _, x := range pks {
		args = append(args, x.Camel()+"Arg")
	}
	ret.W("\treturn as.Services.%s.Get(ps.Context, nil, %s)", m.PackageProper(), strings.Join(args, ", "))
	ret.W("}")

	return ret
}
