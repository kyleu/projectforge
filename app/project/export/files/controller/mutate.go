package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controllerCreateForm(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "create", "form")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
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
	ret.W("\t\treturn %sRender(rc, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %s%s, \"Create\")", prefix, m.Package, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreateFormRandom(m *model.Model, prefix string) *golang.Block {
	ret := blockFor(m, prefix, nil, "create", "form", "random")
	ret.W("\t\tret := %s.Random()", m.Package)
	ret.W("\t\tps.Title = \"Create Random %s\"", m.Proper())
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn %sRender(rc, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %s, \"Create\")", prefix, m.Package, m.Breadcrumbs())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreate(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "create")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.W("\t\tret, err := %sFromForm(rc, true)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp)
	ret.W("\t\terr = as.Services.%s.Create(ps.Context, nil, ps.Logger, ret)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to save newly-created %s\")", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] created\", ret.String())")
	ret.W("\t\treturn %sFlashAndRedir(true, msg, ret.WebPath(), rc, ps)", prefix)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEditForm(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "edit", "form")
	if m.IsRevision() {
		ret.W("\t\trc.SetUserValue(\"includeDeleted\", true)")
	}
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	ret.W("\t\tps.Title = \"Edit \" + ret.String()")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn %sRender(rc, as, &v%s.Edit{Model: ret}, ps, %s%s, ret.String())", prefix, m.Package, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEdit(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "edit")
	if m.IsRevision() {
		ret.W("\t\trc.SetUserValue(\"includeDeleted\", true)")
	}
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	ret.W("\t\tfrm, err := %sFromForm(rc, false)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp, "frm")
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		ret.W("\t\tfrm.%s = ret.%s", pk.Proper(), pk.Proper())
	})
	ret.W("\t\terr = as.Services.%s.Update(ps.Context, nil, frm, ps.Logger)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to update %s [%%%%s]\", frm.String())", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] updated\", frm.String())")
	ret.W("\t\treturn %sFlashAndRedir(true, msg, frm.WebPath(), rc, ps)", prefix)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDelete(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "delete")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	pkCamels := lo.Map(m.PKs(), func(pk *model.Column, _ int) string {
		return "ret." + pk.Proper()
	})
	ret.W("\t\terr = as.Services.%s.Delete(ps.Context, nil, %s, ps.Logger)", m.Proper(), strings.Join(pkCamels, ", "))
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to delete %s [%%%%s]\", ret.String())", m.TitleLower())
	ret.W("\t\t}")
	ret.W("\t\tmsg := fmt.Sprintf(\"" + m.Proper() + " [%%s] deleted\", ret.String())")
	ret.W("\t\treturn %sFlashAndRedir(true, msg, \"/%s\", rc, ps)", prefix, m.Route())
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
