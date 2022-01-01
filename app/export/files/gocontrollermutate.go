package files

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func controllerCreateForm(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.PackageProper()+"CreateForm", "func")
	ret.W("func %sCreateForm(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.create.form\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret := &%s{}", m.ClassRef())
	ret.W("\t\tps.Title = \"Create [" + m.Proper() + "]\"")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %q, \"Create\")", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreate(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.PackageProper()+"Create", "func")
	ret.W("func %sCreate(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.create\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret, err := %sFromForm(rc, true)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\terr = as.Services.%s.Add(ps.Context, nil, ret)", m.PackageProper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to save newly-created %s\")", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] created\", ret.String())")
	ret.W("\t\treturn flashAndRedir(true, msg, ret.WebPath(), rc, ps)")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEditForm(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.PackageProper()+"EditForm", "func")
	ret.W("func %sEditForm(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.edit.form\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = \"Edit [\" + ret.String() + \"]\"")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Edit{Model: ret}, ps, %q, ret.String())", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEdit(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.PackageProper()+"Edit", "func")
	ret.W("func %sEdit(rc *fasthttp.RequestCtx) {", m.PackageProper())
	ret.W("\tact(\"%s.edit\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tfrm, err := %sFromForm(rc, false)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	for _, pk := range m.Columns.PKs() {
		ret.W("\t\tfrm.%s = ret.%s", pk.Proper(), pk.Proper())
	}
	ret.W("\t\terr = as.Services.%s.Update(ps.Context, nil, ret)", m.PackageProper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to update %s [%%%%s]\", frm.String())", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] updated\", ret.String())")
	ret.W("\t\treturn flashAndRedir(true, msg, ret.WebPath(), rc, ps)")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerModelFromForm(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func %sFromForm(rc *fasthttp.RequestCtx, setPK bool) (*%s, error) {", m.Package, m.ClassRef())
	ret.W("\tfrm, err := cutil.ParseForm(rc)")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to parse form\")")
	ret.W("\t}")
	ret.W("\treturn %s.FromMap(frm, setPK)", m.Package)
	ret.W("}")
	return ret
}
