package controller

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func controllerCreateForm(m *model.Model, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "create", "form")
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
	}
	var decls []string
	if grp != nil {
		decls = append(decls, fmt.Sprintf("%s: %sArg", grp.Proper(), grp.Camel()))
	}
	ret.W("\t\tret := &%s{%s}", m.ClassRef(), strings.Join(decls, ", "))
	if grp == nil {
		ret.W("\t\tps.Title = \"Create [" + m.Proper() + "]\"")
	} else {
		ret.W("\t\tps.Title = fmt.Sprintf(\"Create ["+m.Proper()+"] for %s [%%%%s]\", %sArg)", grp.TitleLower(), grp.Camel())
	}
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %q%s, \"Create\")", m.Package, m.Package, grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreateFormRandom(m *model.Model) *golang.Block {
	ret := blockFor(m, nil, "create", "form", "random")
	ret.W("\t\tret := %s.Random()", m.Package)
	ret.W("\t\tps.Title = \"Create Random [" + m.Proper() + "]\"")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %q, \"Create\")", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreate(m *model.Model, g *golang.File, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "create")
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
	}
	ret.W("\t\tret, err := %sFromForm(rc, true)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp)
	ret.W("\t\terr = as.Services.%s.Create(ps.Context, nil, ret)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to save newly-created %s\")", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] created\", ret.String())")
	ret.W("\t\treturn flashAndRedir(true, msg, ret.WebPath(), rc, ps)")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEditForm(m *model.Model, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "edit", "form")
	if m.IsRevision() {
		ret.W("\t\trc.SetUserValue(\"includeDeleted\", true)")
	}
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	checkGrp(ret, grp)
	ret.W("\t\tps.Title = \"Edit [\" + ret.String() + \"]\"")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Edit{Model: ret}, ps, %q%s, ret.String())", m.Package, m.Package, grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEdit(m *model.Model, g *golang.File, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "edit")
	if m.IsRevision() {
		ret.W("\t\trc.SetUserValue(\"includeDeleted\", true)")
	}
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	checkGrp(ret, grp)
	ret.W("\t\tfrm, err := %sFromForm(rc, false)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp, "frm")
	for _, pk := range m.PKs() {
		ret.W("\t\tfrm.%s = ret.%s", pk.Proper(), pk.Proper())
	}
	ret.W("\t\terr = as.Services.%s.Update(ps.Context, nil, frm)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to update %s [%%%%s]\", frm.String())", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] updated\", frm.String())")
	ret.W("\t\treturn flashAndRedir(true, msg, frm.WebPath(), rc, ps)")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDelete(m *model.Model, g *golang.File, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "delete")
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	checkGrp(ret, grp)
	pkCamels := make([]string, 0, len(m.PKs()))
	for _, pk := range m.PKs() {
		pkCamels = append(pkCamels, "ret."+pk.Proper())
	}
	ret.W("\t\terr = as.Services.%s.Delete(ps.Context, nil, %s)", m.Proper(), strings.Join(pkCamels, ", "))
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to delete %s [%%%%s]\", ret.String())", m.TitleLower())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] deleted\", ret.String())")
	ret.W("\t\treturn flashAndRedir(true, msg, \"/%s\", rc, ps)", m.Camel())
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
