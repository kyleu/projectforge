package export

import (
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func exportControllerFile(m *Model, args *Args) *file.File {
	g := NewGoFile("controller", []string{"app", "controller"}, m.Package)
	g.AddImport(ImportTypeExternal, "github.com/pkg/errors")
	g.AddImport(ImportTypeExternal, "github.com/valyala/fasthttp")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/controller/cutil")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/app/util")
	g.AddImport(ImportTypeApp, "{{{ .Package }}}/views/v"+m.Package)
	g.AddBlocks(controllerList(m), controllerDetail(m), controllerModelFromPath(m))
	return g.Render()
}

func controllerList(m *Model) *Block {
	ret := NewBlock(m.camel()+"List", "func")
	ret.WF("func %sList(rc *fasthttp.RequestCtx) {", m.packageProper())
	ret.WF("\tact(\"%s.list\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.WF("\t\tps.Title = %q", util.StringToPlural(m.packageProper()))
	ret.W("\t\tparams := cutil.ParamSetFromRequest(rc)")
	ret.WF("\t\tret, err := as.Services.%s.List(ps.Context, nil, params.Get(%q, nil, ps.Logger))", m.packageProper(), m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")
	ret.WF("\t\treturn render(rc, as, &v%s.List{Models: ret}, ps, %q)", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDetail(m *Model) *Block {
	ret := NewBlock(m.camel()+"Detail", "func")
	ret.WF("func %sDetail(rc *fasthttp.RequestCtx) {", m.packageProper())
	ret.WF("\tact(\"%s.detail\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.WF("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.WF("\t\tps.Title = %q", util.StringToPlural(m.packageProper()))
	ret.W("\t\tps.Data = ret")
	ret.WF("\t\treturn render(rc, as, &v%s.Detail{Model: ret}, ps, %q)", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerModelFromPath(m *Model) *Block {
	ret := NewBlock(m.camel()+"FromPath", "func")
	ret.WF("func %sFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*codegen.User, error) {", m.Package)
	for _, col := range m.Columns.PKs() {
		argFor(col, ret)
	}
	ret.WF("\treturn as.Services.%s.Get(ps.Context, nil, idArg, nameArg)", m.packageProper())
	ret.W("}")

	return ret
}

func argFor(col *Column, b *Block) {
	switch col.Type.Keys[0] {
	case "string":
		b.WF("\t%sArg, err := rcRequiredString(rc, %q, false)", col.camel(), col.camel())
		b.W("\tif err != nil {")
		b.WF("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.camel())
		b.W("\t}")
	case "uuid":
		b.WF("\t%sArgStr, err := rcRequiredString(rc, %q, false)", col.camel(), col.camel())
		b.W("\tif err != nil {")
		b.WF("\t\treturn nil, errors.Wrap(err, \"must provide [%s] as an argument\")", col.camel())
		b.W("\t}")
		b.WF("\t%sArgP := util.UUIDFromString(%sArgStr)", col.camel(), col.camel())
		b.WF("\tif %sArgP == nil {", col.camel())
		b.WF("\t\treturn nil, errors.Errorf(\"argument [%s] (%%%%s) is not a valid UUID\", %sArgStr)", col.camel(), col.camel())
		b.W("\t}")
		b.WF("\t%sArg := *%sArgP", col.camel(), col.camel())
	default:
		b.WF("\tERROR: unhandled arg type [%s]", col.Type.String())
	}
}
